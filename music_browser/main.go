package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/emilia-miki/music-browser/music_browser/environment"
	"github.com/emilia-miki/music-browser/music_browser/explorer"
	"github.com/emilia-miki/music-browser/music_browser/explorer/grpc_backend"
	"github.com/emilia-miki/music-browser/music_browser/explorer/local_backend"
	"github.com/emilia-miki/music-browser/music_browser/explorer/spotify_backend"
	graphql_schema "github.com/emilia-miki/music-browser/music_browser/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/lib/pq"
)

type Error struct {
	Message string
}

var explorerObj *explorer.Explorer

func sendBadRequestErrorResponse(
	response http.ResponseWriter, message string,
) error {
	errorResponse := Error{Message: message}
	jsonBytes, err := json.Marshal(errorResponse)
	if err != nil {
		return err
	}

	response.WriteHeader(http.StatusBadRequest)
	response.Write(jsonBytes)

	return nil
}

func sendOkJsonResponse(response http.ResponseWriter, jsonBytes []byte) {
	response.WriteHeader(http.StatusOK)
	response.Write(jsonBytes)
}

func requestHandler(response http.ResponseWriter, request *http.Request) {
	params := request.URL.Query()
	if request.Method == "GET" {
		backendName := params.Get("backend")
		searchType := params.Get("type")
		query := params.Get("query")
		url := params.Get("url")

		var jsonBytes []byte
		if searchType == "artists" {
			artists, err := explorerObj.SearchArtists(backendName, query)
			if err != nil {
				log.Fatal(err)
			}

			jsonBytes, err = json.Marshal(artists)
			if err != nil {
				log.Fatal(err)
			}
		} else if searchType == "albums" {
			albums, err := explorerObj.SearchAlbums(backendName, query)
			if err != nil {
				log.Fatal(err)
			}

			jsonBytes, err = json.Marshal(albums)
			if err != nil {
				log.Fatal(err)
			}
		} else if searchType == "artist" {
			artist, err := explorerObj.GetArtist(url)
			if err != nil {
				log.Fatal(err)
			}

			jsonBytes, err = json.Marshal(artist)
			if err != nil {
				log.Fatal(err)
			}
		} else if searchType == "album" {
			album, err := explorerObj.GetAlbum(url)
			if err != nil {
				log.Fatal(err)
			}

			jsonBytes, err = json.Marshal(album)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := sendBadRequestErrorResponse(response, "Invalid type")
			if err != nil {
				log.Fatal(err)
			}
		}

		sendOkJsonResponse(response, jsonBytes)
	} else if request.Method == "POST" {
		url := params.Get("url")
		log.Println("the fuck???? " + url)
		track, err := explorerObj.GetTrack(url)
		if err != nil {
			log.Fatal(err)
		}

		err = explorerObj.DownloadTrack(track)
		if err != nil {
			log.Fatal(err)
		}

		response.WriteHeader(http.StatusOK)
	} else {
		err := sendBadRequestErrorResponse(response,
			"Only GET and POST methods are allowed on this endpoint")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	app := environment.GetApplication()
	secrets := environment.GetSecrets()
	backends := make(map[string]explorer.Backend)

	var err error
	backends["spotify"] = spotify_backend.New(secrets.Spotify)
	if err != nil {
		log.Fatal(err)
	}
	backends["bandcamp"], err = grpc_backend.New(
		fmt.Sprintf("localhost:%d", app.Ports.BandcampAPI))
	if err != nil {
		log.Fatal(err)
	}
	backends["yt-music"], err = grpc_backend.New(
		fmt.Sprintf("localhost:%d", app.Ports.YTMusicAPI))
	if err != nil {
		log.Fatal(err)
	}
	backends["local"], err = local_backend.New(
		app.ConnectionStrings.PostgreSQL)
	if err != nil {
		log.Fatal(err)
	}

	explorerObj, err = explorer.New(
		backends,
		app.ConnectionStrings.Redis,
		app.ConnectionStrings.PostgreSQL,
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", requestHandler)
	schema := graphql_schema.NewSchema(explorerObj)
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
