package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/emilia-miki/music-browser/music_browser/backend"
	bandcamp_backend "github.com/emilia-miki/music-browser/music_browser/backend/bandcamp"
	local_backend "github.com/emilia-miki/music-browser/music_browser/backend/local"
	"github.com/emilia-miki/music-browser/music_browser/backend/music_downloader"
	"github.com/emilia-miki/music-browser/music_browser/backend/music_explorer_cache"
	spotify_backend "github.com/emilia-miki/music-browser/music_browser/backend/spotify"
	yt_music_backend "github.com/emilia-miki/music-browser/music_browser/backend/youtube_music"
	"github.com/emilia-miki/music-browser/music_browser/environment"
	graphql_schema "github.com/emilia-miki/music-browser/music_browser/graphql"
	"github.com/emilia-miki/music-browser/music_browser/music_api"
	"github.com/graphql-go/handler"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var spotify backend.MusicExplorer
var ytMusic backend.MusicExplorer
var bandcamp backend.MusicExplorer
var local backend.MusicExplorer

var downloader backend.MusicDownloader

type Error struct {
	Message string
}

func serializeResponseBody(obj interface{}) []byte {
	json, err := json.Marshal(obj)
	if err != nil {
		log.Fatalf("Error serializing response body %s: %s\n", obj, err)
	}

	return json
}

func sendBadRequestErrorResponse(response http.ResponseWriter, message string) {
	errorResponse := Error{Message: message}
	json := serializeResponseBody(errorResponse)

	response.WriteHeader(http.StatusBadRequest)
	response.Write(json)
}

func sendOkJsonResponse(response http.ResponseWriter, json []byte) {
	response.WriteHeader(http.StatusOK)
	response.Write(json)
}

func getRequestHandler(response http.ResponseWriter, params url.Values) {
	var backend backend.MusicExplorer
	backendString := params.Get("backend")
	if backendString == "spotify" {
		backend = spotify
	} else if backendString == "bandcamp" {
		backend = bandcamp
	} else if backendString == "yt-music" {
		backend = ytMusic
	} else if backendString == "local" {
		backend = local
	} else {
		sendBadRequestErrorResponse(response, "Invalid backend")
	}

	query := params.Get("query")
	if query == "" {
		sendBadRequestErrorResponse(response, "You must provide a query string")
		return
	}

	var results interface{}
	searchType := params.Get("search-type")
	if searchType == "artist" {
		results = backend.SearchArtists(query)
	} else if searchType == "album" {
		results = backend.SearchAlbums(query)
	} else {
		sendBadRequestErrorResponse(response, "Invalid search type")
	}

	json := serializeResponseBody(results)
	sendOkJsonResponse(response, json)
}

func postRequestHandler(response http.ResponseWriter, params url.Values) {
	url := params.Get("url")
	if url == "" {
		sendBadRequestErrorResponse(response, "You must provide a song URL")
		return
	}

	downloader.DownloadTrack(url)
	response.WriteHeader(http.StatusOK)
}

func requestHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		getRequestHandler(response, request.URL.Query())
	} else if request.Method == "POST" {
		postRequestHandler(response, request.URL.Query())
	} else {
		sendBadRequestErrorResponse(response,
			"Only GET and POST methods are allowed on this endpoint")
	}
}

func main() {
	app := environment.GetApplication()
	secrets := environment.GetSecrets()

	redisOpt, err := redis.ParseURL(app.ConnectionStrings.Redis)
	if err != nil {
		log.Fatalf("Error parsing connection string for Redis: %s\n", err)
	}
	redisClient := redis.NewClient(redisOpt)
	defer redisClient.Close()
	cache := music_explorer_cache.MusicExplorerCache{RedisClient: redisClient}

	spotify = &spotify_backend.MusicExplorer{
		Cache:   cache,
		Secrets: secrets.Spotify,
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	ytMusicApiConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", app.Ports.YTMusicAPI), opts...)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %s", err)
	}
	defer ytMusicApiConn.Close()

	ytMusic = &yt_music_backend.MusicExplorer{
		Cache:  cache,
		Client: music_api.NewMusicApiClient(ytMusicApiConn),
	}

	bandcampApiConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", app.Ports.BandcampAPI), opts...)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %s", err)
	}
	defer bandcampApiConn.Close()

	bandcamp = &bandcamp_backend.MusicExplorer{
		Cache:  cache,
		Client: music_api.NewMusicApiClient(bandcampApiConn),
	}

	postgresDB, err := sql.Open("postgres", app.ConnectionStrings.PostgreSQL)
	defer postgresDB.Close()
	local = &local_backend.MusicExplorer{PostgresDB: postgresDB}

	downloader = &music_downloader.MusicDownloader{
		YtDlpConn:    ytMusicApiConn,
		PostgresConn: postgresDB,
	}

	http.HandleFunc("/", requestHandler)
	schema := graphql_schema.GetSchema(map[string]backend.MusicExplorer{
		"spotify":  spotify,
		"bandcamp": bandcamp,
		"yt-music": ytMusic,
		"local":    local,
	})
	http.Handle("/graphql", handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	}))

	log.Printf("Server listening on port %d\n", app.Ports.Server)
	err = http.ListenAndServe(fmt.Sprintf(":%d", app.Ports.Server), nil)
	if err != nil {
		log.Fatalf("Unable to start the server on port %d\n", app.Ports.Server)
	}
}
