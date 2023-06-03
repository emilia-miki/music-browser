package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/emilia-miki/music-browser/backend"
	bandcamp_backend "github.com/emilia-miki/music-browser/backend/bandcamp"
	local_backend "github.com/emilia-miki/music-browser/backend/local"
	"github.com/emilia-miki/music-browser/backend/music_downloader"
	"github.com/emilia-miki/music-browser/backend/music_explorer_cache"
	spotify_backend "github.com/emilia-miki/music-browser/backend/spotify"
	yt_music_backend "github.com/emilia-miki/music-browser/backend/youtube_music"
	"github.com/emilia-miki/music-browser/environment"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
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
	switch params.Get("backend") {
	case "spotify":
		backend = spotify
		break
	case "yt-music":
		backend = ytMusic
		break
	case "bandcamp":
		backend = bandcamp
		break
	case "local":
		backend = local
		break
	default:
		sendBadRequestErrorResponse(response, "Invalid backend")
		return
	}

	query := params.Get("query")
	if query == "" {
		sendBadRequestErrorResponse(response, "You must provide a query string")
		return
	}

	var results interface{}
	switch params.Get("search-type") {
	case "track":
		results = backend.SearchTracks(query)
		break
	case "album":
		results = backend.SearchAlbums(query)
		break
	case "artist":
		results = backend.SearchArtists(query)
		break
	default:
		sendBadRequestErrorResponse(response, "Invalid search type")
		return
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
	log.Printf("Processing a %s request on %s\n", request.Method, request.RequestURI)

	switch request.Method {
	case "GET":
		getRequestHandler(response, request.URL.Query())
		break
	case "POST":
		postRequestHandler(response, request.URL.Query())
		break
	default:
		sendBadRequestErrorResponse(response, "Only GET and POST methods are allowed")
		break
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, "python3", "yt_music_api.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start Python script: %s", err)
	}

	ytMusicApiConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", app.Ports.YTMusicAPI), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %s", err)
	}
	defer ytMusicApiConn.Close()

	ytMusic = &yt_music_backend.MusicExplorer{
		Cache:          cache,
		YtMusicApiConn: ytMusicApiConn,
	}

	cmd = exec.CommandContext(ctx, "node", "bandcamp_api.js")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start Node.JS script: %s", err)
	}

	bandcampApiConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", app.Ports.BandcampAPI), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %s", err)
	}
	defer bandcampApiConn.Close()

	bandcamp = &bandcamp_backend.MusicExplorer{
		Cache:               cache,
		BandcampScraperConn: bandcampApiConn,
	}

	postgresDB, err := sql.Open("postgres", app.ConnectionStrings.PostgreSQL)
	defer postgresDB.Close()
	local = &local_backend.MusicExplorer{PostgresDB: postgresDB}

	downloader = &music_downloader.MusicDownloader{
		YtDlpConn:    ytMusicApiConn,
		PostgresConn: postgresDB,
	}

	http.HandleFunc("/", requestHandler)

	log.Printf("Server listening on port %d\n", app.Ports.Server)
	err = http.ListenAndServe(fmt.Sprintf(":%d", app.Ports.Server), nil)
	if err != nil {
		log.Fatalf("Unable to start the server on port %d\n", app.Ports.Server)
	}
}
