package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/emilia-miki/music-browser/music_browser/environment"
	"github.com/emilia-miki/music-browser/music_browser/explorer"
	"github.com/emilia-miki/music-browser/music_browser/explorer/grpc_backend"
	"github.com/emilia-miki/music-browser/music_browser/explorer/local_backend"
	"github.com/emilia-miki/music-browser/music_browser/explorer/spotify_backend"
	graphql_schema "github.com/emilia-miki/music-browser/music_browser/graphql"
	"github.com/emilia-miki/music-browser/music_browser/logger"
	"github.com/graphql-go/handler"
	_ "github.com/lib/pq"
)

type GetRequestType = string

const (
	GetRequestSearchArtists  GetRequestType = "artists"
	GetRequestSearchAlbums                  = "albums"
	GetRequestGetArtist                     = "artist"
	GetRequestGetAlbum                      = "album"
	GetRequestGetTrack                      = "track"
	GetRequestGetTrackStream                = "track-stream"
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

func sendOkResponse(
	response http.ResponseWriter,
	bytes []byte,
	mime string,
) error {
	response.Header().Add("Content-Type", mime)
	_, err := response.Write(bytes)
	if err != nil {
		return err
	}

	if mime == "application/json" {
		log.Println("sending an OK response with the following content: " +
			string(bytes))
	} else {
		log.Printf("sending an OK response (%d bytes)\n", len(bytes))
	}
	return nil
}

func extractType(url string) GetRequestType {
	urlTrimmed, trimmed := strings.CutPrefix(url, "http://")
	if !trimmed {
		urlTrimmed, _ = strings.CutPrefix(url, "https://")
	}

	splits := strings.Split(urlTrimmed, "/")

	if len(splits) == 3 && splits[0] == "open.spotify.com" && (splits[1] == GetRequestGetArtist ||
		splits[1] == GetRequestGetAlbum ||
		splits[1] == GetRequestGetTrack) {
		return splits[1]
	}

	if len(splits) == 3 && splits[0] == "music.youtube.com" &&
		splits[1] == "channel" {
		return GetRequestGetArtist
	}

	if len(splits) == 3 && splits[0] == "music.youtube.com" &&
		splits[1] == "browse" {
		return GetRequestGetAlbum
	}

	if len(splits) == 2 && splits[0] == "music.youtube.com" &&
		strings.HasPrefix(splits[1], "watch?v=") {
		return GetRequestGetTrack
	}

	if len(splits) == 1 && strings.HasSuffix(splits[0], ".bandcamp.com") {
		return GetRequestGetArtist
	}

	if len(splits) == 3 && strings.HasSuffix(splits[0], ".bandcamp.com") {
		if splits[1] == GetRequestGetAlbum || splits[1] == GetRequestGetTrack {
			return splits[1]
		}
	}

	return ""
}

func requestHandler(response http.ResponseWriter, request *http.Request) {
	log.Printf("received a %s request with uri %s\n",
		request.Method, request.URL.String())

	params := request.URL.Query()
	if request.Method == "GET" {
		reqBackend := params.Get("backend")
		reqType := params.Get("type")
		reqQuery := params.Get("query")
		reqUrl := params.Get("url")

		if reqType != GetRequestSearchArtists &&
			reqType != GetRequestSearchAlbums &&
			reqType != GetRequestGetTrackStream {
			reqType = extractType(reqUrl)
		}

		var bytes []byte
		var mime string = "application/json"
		var err error
		if reqType == GetRequestSearchArtists {
			artists, err := explorerObj.SearchArtists(reqBackend, reqQuery)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}

			bytes, err = json.Marshal(artists)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}
		} else if reqType == GetRequestSearchAlbums {
			albums, err := explorerObj.SearchAlbums(reqBackend, reqQuery)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}

			bytes, err = json.Marshal(albums)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}
		} else if reqType == GetRequestGetArtist {
			artist, err := explorerObj.GetArtist(reqUrl)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}

			bytes, err = json.Marshal(artist)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}
		} else if reqType == GetRequestGetAlbum {
			album, err := explorerObj.GetAlbum(reqUrl)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}

			bytes, err = json.Marshal(album)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}
		} else if reqType == GetRequestGetTrack {
			track, err := explorerObj.GetTrack(reqUrl)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}

			bytes, err = json.Marshal(track)
			if err != nil {
				err = sendInternalServerErrorResponse(response, err)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				return
			}
		} else if reqType == GetRequestGetTrackStream {
			bytes, mime, err = explorerObj.GetTrackStream(reqUrl)
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

		err = sendOkResponse(response, bytes, mime)
		if err != nil {
			err = sendInternalServerErrorResponse(response, err)
			if err != nil {
				log.Println("ERROR: " + err.Error())
			}
		}
	} else if request.Method == "POST" {
		url := params.Get("url")
		err := explorerObj.DownloadTrack(url)
		if err != nil {
			err = sendInternalServerErrorResponse(response, err)
			if err != nil {
				log.Println("ERROR: " + err.Error())
			}
			return
		}

		err = sendOkResponse(response, []byte{'{', '}'}, "application/json")
		if err != nil {
			err = sendInternalServerErrorResponse(response, err)
			if err != nil {
				log.Println("ERROR: " + err.Error())
			}
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
	logger.Initialize()

	logger.Info.Println("Parsing environment variables")
	app := environment.GetApplication()
	secrets := environment.GetSecrets()
	backends := make(map[string]explorer.Backend)

	logger.Info.Println("Initializing Spotify backend")
	var err error
	backends[explorer.SpotifyBackendName] = spotify_backend.New(secrets.Spotify)
	if err != nil {
		log.Fatalln(err)
	}
	logger.Info.Println("Initializing Bandcamp backend")
	backends[explorer.BandcampBackendName], err = grpc_backend.New(
		explorer.BandcampBackendName,
		fmt.Sprintf("localhost:%d", app.Ports.BandcampAPI))
	if err != nil {
		log.Fatalln(err)
	}
	logger.Info.Println("Initializing Youtube Music backend")
	backends[explorer.YtMusicBackendName], err = grpc_backend.New(
		explorer.YtMusicBackendName,
		fmt.Sprintf("localhost:%d", app.Ports.YTMusicAPI))
	if err != nil {
		log.Fatalln(err)
	}
	logger.Info.Println("Initializing Local backend")
	backends[explorer.LocalBackendName], err = local_backend.New(
		app.ConnectionStrings.PostgreSQL)
	if err != nil {
		log.Fatalln(err)
	}

	logger.Info.Println("Initializing explorer")
	explorerObj, err = explorer.New(
		backends,
		app.ConnectionStrings.Redis,
		app.ConnectionStrings.PostgreSQL,
	)
	if err != nil {
		log.Fatalln(err)
	}

	logger.Info.Println("Initializing GraphQL")
	schema := graphql_schema.NewSchema(explorerObj)

	logger.Info.Println("Setting up request handlers")
	http.HandleFunc("/", requestHandler)
	http.Handle("/graphql", handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	}))

	logger.Info.Printf("Server listening on port %d\n", app.Ports.Server)
	fmt.Printf("Server listening on port %d\n", app.Ports.Server)
	err = http.ListenAndServe(fmt.Sprintf(":%d", app.Ports.Server), nil)
	if err != nil {
		log.Fatalf("Unable to start the server on port %d\n", app.Ports.Server)
	}

	logger.Inform()
}
