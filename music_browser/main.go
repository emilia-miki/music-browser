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

type GetRequestType = string

const (
	GetRequestSearchArtists GetRequestType = "artists"
	GetRequestSearchAlbums                 = "albums"
	GetRequestGetArtist                    = "artist"
	GetRequestGetAlbum                     = "album"
	GetRequestGetTrack                     = "track"
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

func sendInternalServerErrorResponse(
	response http.ResponseWriter, err error,
) error {
	errorResponse := Error{Message: err.Error()}
	jsonBytes, err := json.Marshal(errorResponse)
	if err != nil {
		return err
	}

	response.WriteHeader(http.StatusInternalServerError)
	_, err = response.Write(jsonBytes)
	if err != nil {
		return err
	}

	log.Println(
		"sending an InternalServerError response with the following content: " +
			string(jsonBytes))
	return nil
}

func sendOkJsonResponse(response http.ResponseWriter, jsonBytes []byte) error {
	_, err := response.Write(jsonBytes)
	if err != nil {
		return err
	}

	log.Println("sending an OK response with the following content: " +
		string(jsonBytes))
	return nil
}

func requestHandler(response http.ResponseWriter, request *http.Request) {
	log.Printf("received a %s request with uri %s\n",
		request.Method, request.URL.String())

	params := request.URL.Query()
	if request.Method == "GET" {
		backendName := params.Get("backend")
		getRequestType := params.Get("type")
		query := params.Get("query")
		url := params.Get("url")

		var jsonBytes []byte
		if getRequestType == GetRequestGetArtist {
			artist, err := explorerObj.GetArtist(url)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}

			jsonBytes, err = json.Marshal(artist)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}
		} else if getRequestType == GetRequestGetAlbum {
			album, err := explorerObj.GetAlbum(url)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}

			jsonBytes, err = json.Marshal(album)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}
		} else if getRequestType == GetRequestGetTrack {
			track, err := explorerObj.GetTrack(url)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}

			jsonBytes, err = json.Marshal(track)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}
		} else if getRequestType == GetRequestSearchArtists {
			artists, err := explorerObj.SearchArtists(backendName, query)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}

			jsonBytes, err = json.Marshal(artists)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}
		} else if getRequestType == GetRequestSearchAlbums {
			albums, err := explorerObj.SearchAlbums(backendName, query)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}

			jsonBytes, err = json.Marshal(albums)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}
		} else {
			err := sendBadRequestErrorResponse(response, "Invalid type")
			if err != nil {
				log.Println("ERROR: " + err.Error())
			}
		}

		err := sendOkJsonResponse(response, jsonBytes)
		if err != nil {
			log.Println("ERROR: " + err.Error())
		}
	} else if request.Method == "POST" {
		url := params.Get("url")
		track, err := explorerObj.GetTrack(url)
		if err != nil {
			log.Println("ERROR: " + err.Error())
		}

		err = explorerObj.DownloadTrack(track)
		if err != nil {
			log.Println("ERROR: " + err.Error())
		}

		err = sendOkJsonResponse(response, []byte{})
		if err != nil {
			log.Println("ERROR: " + err.Error())
		}
	} else {
		err := sendBadRequestErrorResponse(response,
			"Only GET and POST methods are allowed on this endpoint")
		if err != nil {
			log.Println("ERROR: " + err.Error())
		}
	}
}

func main() {
	app := environment.GetApplication()
	secrets := environment.GetSecrets()
	backends := make(map[string]explorer.Backend)

	var err error
	backends[explorer.SpotifyBackendName] = spotify_backend.New(secrets.Spotify)
	if err != nil {
		log.Fatalln(err)
	}
	backends[explorer.BandcampBackendName], err = grpc_backend.New(
		explorer.BandcampBackendName,
		fmt.Sprintf("localhost:%d", app.Ports.BandcampAPI))
	if err != nil {
		log.Fatalln(err)
	}
	backends[explorer.YtMusicBackendName], err = grpc_backend.New(
		explorer.YtMusicBackendName,
		fmt.Sprintf("localhost:%d", app.Ports.YTMusicAPI))
	if err != nil {
		log.Fatalln(err)
	}
	backends[explorer.LocalBackendName], err = local_backend.New(
		app.ConnectionStrings.PostgreSQL)
	if err != nil {
		log.Fatalln(err)
	}

	explorerObj, err = explorer.New(
		backends,
		app.ConnectionStrings.Redis,
		app.ConnectionStrings.PostgreSQL,
	)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/", requestHandler)
	schema := graphql_schema.NewSchema(explorerObj)
	http.Handle("/graphql", handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	}))

	fmt.Printf("Server listening on port %d\n", app.Ports.Server)
	err = http.ListenAndServe(fmt.Sprintf(":%d", app.Ports.Server), nil)
	if err != nil {
		log.Fatalf("Unable to start the server on port %d\n", app.Ports.Server)
	}
}
