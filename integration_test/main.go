package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/emilia-miki/music-browser/integration_test/music_api"
)

var warningLogger *log.Logger
var errorLogger *log.Logger

type BackendName = string

const (
	SpotifyBackendName  BackendName = "spotify"
	YtMusicBackendName              = "yt-music"
	BandcampBackendName             = "bandcamp"
)

type GetRequestType = string

const (
	GetRequestSearchArtists GetRequestType = "artists"
	GetRequestSearchAlbums                 = "albums"
	GetRequestGetArtist                    = "artist"
	GetRequestGetAlbum                     = "album"
	GetRequestGetTrack                     = "track"
)

type Result struct {
	BackendName    BackendName
	GetRequestType GetRequestType
	Url            string
	Query          string
	Success        bool
}

func (r Result) String() string {
	var success string
	if r.Success {
		success = "SUCCESS "
	} else {
		success = "FAILURE "
	}

	return success + describeTest(
		r.BackendName, r.GetRequestType, r.Url, r.Query)
}

var Results []Result

func sendGetRequest(
	requestBackend string,
	requestType string,
	requestUrl string,
	requestQuery string,
) ([]byte, error) {
	b := url.QueryEscape(requestBackend)
	t := url.QueryEscape(requestType)
	u := url.QueryEscape(requestUrl)
	q := url.QueryEscape(requestQuery)

	requestUri := fmt.Sprintf(
		"http://localhost:3333/?backend=%s&type=%s&url=%s&query=%s", b, t, u, q)
	resp, err := http.Get(requestUri)
	if err != nil {
		return nil, errors.New("Unable to send a GET request: " + err.Error())
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(
			"Unable to read the response body: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Message string `json:"message"`
		}
		err = json.Unmarshal(bytes, &errorResponse)
		if err == nil {
			return nil, errors.New(fmt.Sprintf(
				"Error %d: %s", resp.StatusCode, errorResponse.Message))
		} else {
			return nil, errors.New(fmt.Sprintf(
				"Error %d: Unable to read error message from the response: %s",
				resp.StatusCode, err.Error(),
			))
		}
	}

	return bytes, nil
}

func describeTest(
	backendName string,
	requestType GetRequestType,
	url string,
	query string,
) string {
	var testDescription string

	if backendName != "" {
		testDescription += "backend=" + backendName + ", "
	}

	testDescription += "type=" + requestType + ", "

	if url != "" {
		testDescription += "url=" + url + ", "
	}

	if query != "" {
		testDescription += "query=" + query + ", "
	}

	testDescription, _ = strings.CutSuffix(testDescription, ", ")

	return testDescription
}

func getAndValidate(
	backendName string,
	requestType GetRequestType,
	url string,
	query string,
	unmarshaler func([]byte) (interface{}, error),
	validator func(interface{}) error,
) {
	testDescription := describeTest(backendName, requestType, url, query)
	var success bool
	defer func() {
		Results = append(Results, Result{
			BackendName:    backendName,
			GetRequestType: requestType,
			Url:            url,
			Query:          query,
			Success:        success,
		})
	}()

	log.Println("Testing GET request with " + testDescription)

	resp, err := sendGetRequest(backendName, requestType, url, query)
	if err != nil {
		errorLogger.Println(err)
		return
	}

	var buf bytes.Buffer
	err = json.Indent(&buf, resp, "", "\t")
	if err != nil {
		panic(err)
	}
	log.Println("Response:\n" + buf.String())

	obj, err := unmarshaler(resp)
	if err != nil {
		errorLogger.Println("Unable to unmarshal response: " + err.Error() +
			". The response was: " + string(resp))
		return
	}

	err = validator(obj)
	if err != nil {
		errorLogger.Println(err)
		return
	}

	success = true
}

func startContainer(
	containerTool string,
	spotifyClientId string,
	spotifyClientSecret string,
) (*exec.Cmd, error) {
	cmd := exec.CommandContext(
		context.Background(),
		containerTool, "run", "-p", "3333:3333",
		"-e", fmt.Sprintf("SPOTIFY_CLIENT_ID=%s", spotifyClientId),
		"-e", fmt.Sprintf("SPOTIFY_CLIENT_SECRET=%s", spotifyClientSecret),
		"music-browser")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.New("Unable to open stdout pipe")
	}

	cmd.Start()

	scanner := bufio.NewScanner(pipe)
	var output string
	var line string
	for {
		res := scanner.Scan()
		if res == false {
			err := scanner.Err()
			if err == nil {
				return nil, errors.New(
					"Unexpected EOF encountered while capturing container " +
						"output. Output so far: " + output)
			} else {
				return nil, errors.New(
					"Unexpected error encounteres while capturing container " +
						"output. Output so far: " + output)
			}
		}

		line = scanner.Text()
		output += line
		if strings.Contains(line, "Server listening on port 3333") {
			break
		}
	}

	go func() {
		for scanner.Scan() {
		}
		pipe.Close()
	}()

	return cmd, nil
}

func unmarshalArtistWithAlbums(b []byte) (interface{}, error) {
	artist := new(music_api.ArtistWithAlbums)
	err := json.Unmarshal(b, artist)
	return artist, err
}

func unmarshalAlbumWithTracks(b []byte) (interface{}, error) {
	album := new(music_api.AlbumWithTracks)
	err := json.Unmarshal(b, album)
	return album, err
}

func unmarshalTrackWithAlbumAndArtist(b []byte) (interface{}, error) {
	track := new(music_api.TrackWithAlbumAndArtist)
	err := json.Unmarshal(b, track)
	return track, err
}

func unmarshalArtists(b []byte) (interface{}, error) {
	artists := new(music_api.Artists)
	err := json.Unmarshal(b, artists)
	return artists, err
}

func unmarshalAlbums(b []byte) (interface{}, error) {
	albums := new(music_api.Albums)
	err := json.Unmarshal(b, albums)
	return albums, err
}

func validateArtist(artist *music_api.Artist) error {
	if artist == nil {
		err := errors.New("Artist is nil")
		errorLogger.Println(err)
		return err
	}

	if artist.Name == nil {
		err := errors.New("Artist is nil")
		errorLogger.Println(err)
		errorLogger.Writer().Write([]byte("IN OBJECT:\n"))
		b, err := json.MarshalIndent(artist, "", "\t")
		if err != nil {
			panic(err)
		}
		errorLogger.Writer().Write(b)
		errorLogger.Writer().Write([]byte{'\n'})
		return err
	}

	if artist.Url == nil {
		warningLogger.Println("Artist.Url is nil")
		warningLogger.Writer().Write([]byte("IN OBJECT:\n"))
		b, err := json.MarshalIndent(artist, "", "\t")
		if err != nil {
			panic(err)
		}
		warningLogger.Writer().Write(b)
		warningLogger.Writer().Write([]byte{'\n'})
	}

	if artist.ImageUrl == nil {
		warningLogger.Println("Artist.ImageUrl is nil")
		warningLogger.Writer().Write([]byte("IN OBJECT:\n"))
		b, err := json.MarshalIndent(artist, "", "\t")
		if err != nil {
			panic(err)
		}
		warningLogger.Writer().Write(b)
		warningLogger.Writer().Write([]byte{'\n'})
	}

	return nil
}

func validateArtists(artistsInterface interface{}) error {
	artists := artistsInterface.(*music_api.Artists)
	if artists == nil {
		err := errors.New("Artists is nil")
		errorLogger.Println(err)
		return err
	}

	for _, artist := range artists.Artists {
		err := validateArtist(artist)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateAlbum(album *music_api.Album) error {
	if album == nil {
		err := errors.New("Album is nil")
		errorLogger.Println(err)
		return err
	}

	if album.Name == nil {
		err := errors.New("Album is nil")
		errorLogger.Println(err)
		return err
	}

	if album.Url == nil {
		warningLogger.Println("Album.Url is nil")
		warningLogger.Writer().Write([]byte("IN OBJECT:\n"))
		b, err := json.MarshalIndent(album, "", "\t")
		if err != nil {
			panic(err)
		}
		warningLogger.Writer().Write(b)
		warningLogger.Writer().Write([]byte{'\n'})
	}

	if album.ImageUrl == nil {
		warningLogger.Println("Album.ImageUrl is nil")
		warningLogger.Writer().Write([]byte("IN OBJECT:\n"))
		b, err := json.MarshalIndent(album, "", "\t")
		if err != nil {
			panic(err)
		}
		warningLogger.Writer().Write(b)
		warningLogger.Writer().Write([]byte{'\n'})
	}

	if album.ArtistUrl == nil {
		warningLogger.Println("Album.ArtistUrl is nil")
		warningLogger.Writer().Write([]byte("IN OBJECT:\n"))
		b, err := json.MarshalIndent(album, "", "\t")
		if err != nil {
			panic(err)
		}
		warningLogger.Writer().Write(b)
		warningLogger.Writer().Write([]byte{'\n'})
	}

	if album.Year == nil {
		warningLogger.Println("Album.Year is nil")
		warningLogger.Writer().Write([]byte("IN OBJECT:\n"))
		b, err := json.MarshalIndent(album, "", "\t")
		if err != nil {
			panic(err)
		}
		warningLogger.Writer().Write(b)
		warningLogger.Writer().Write([]byte{'\n'})
	}

	return nil
}

func validateAlbums(albumsInterface interface{}) error {
	albums := albumsInterface.(*music_api.Albums)
	if albums == nil {
		err := errors.New("Albums is nil")
		errorLogger.Println(err)
		return err
	}

	for _, album := range albums.Albums {
		err := validateAlbum(album)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateTrack(track *music_api.Track) error {
	if track == nil {
		err := errors.New("Track is nil")
		errorLogger.Println(err)
		return err
	}

	if track.Name == nil {
		err := errors.New("Track.Name is nil")
		errorLogger.Println(err)
		errorLogger.Writer().Write([]byte("IN OBJECT:\n"))
		b, err := json.MarshalIndent(track, "", "\t")
		if err != nil {
			panic(err)
		}
		errorLogger.Writer().Write(b)
		errorLogger.Writer().Write([]byte{'\n'})
		return err
	}

	if track.ImageUrl == nil {
		warningLogger.Println("Track.ImageUrl is nil")
		warningLogger.Writer().Write([]byte("IN OBJECT:\n"))
		b, err := json.MarshalIndent(track, "", "\t")
		if err != nil {
			panic(err)
		}
		warningLogger.Writer().Write(b)
		warningLogger.Writer().Write([]byte{'\n'})
	}

	if track.ArtistUrl == nil {
		warningLogger.Println("Track.ArtistUrl is nil")
		warningLogger.Writer().Write([]byte("IN OBJECT:\n"))
		b, err := json.MarshalIndent(track, "", "\t")
		if err != nil {
			panic(err)
		}
		warningLogger.Writer().Write(b)
		warningLogger.Writer().Write([]byte{'\n'})
	}

	if track.AlbumUrl == nil {
		warningLogger.Println("Track.AlbumUrl is nil")
		warningLogger.Writer().Write([]byte("IN OBJECT:\n"))
		b, err := json.MarshalIndent(track, "", "\t")
		if err != nil {
			panic(err)
		}
		warningLogger.Writer().Write(b)
		warningLogger.Writer().Write([]byte{'\n'})
	}

	if track.DurationSeconds == nil {
		warningLogger.Println("Track.DurationSeconds is nil")
		warningLogger.Writer().Write([]byte("IN OBJECT:\n"))
		b, err := json.MarshalIndent(track, "", "\t")
		if err != nil {
			panic(err)
		}
		warningLogger.Writer().Write(b)
		warningLogger.Writer().Write([]byte{'\n'})
	}

	return nil
}

func validateTracks(tracks *music_api.Tracks) error {
	if tracks == nil {
		err := errors.New("Tracks is nil")
		errorLogger.Println(err)
		return err
	}

	for _, track := range tracks.Tracks {
		err := validateTrack(track)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateArtistWithAlbums(
	artistWithAlbumsInterface interface{},
) error {
	artistWithAlbums := artistWithAlbumsInterface.(*music_api.ArtistWithAlbums)
	if artistWithAlbums == nil {
		err := errors.New("ArtistWithAlbums is nil")
		errorLogger.Println(err)
		return err
	}

	if artistWithAlbums.Artist == nil {
		err := errors.New("ArtistWithAlbums.Artist is nil")
		errorLogger.Println(err)
		return err
	}

	if artistWithAlbums.Albums == nil {
		err := errors.New("ArtistWithAlbums.Albums is nil")
		errorLogger.Println(err)
		return err
	}

	err := validateArtist(artistWithAlbums.Artist)
	if err != nil {
		return err
	}

	err = validateAlbums(artistWithAlbums.Albums)
	if err != nil {
		return err
	}

	return nil
}

func validateAlbumWithTracks(
	albumWithTracksInterface interface{},
) error {
	albumWithTracks := albumWithTracksInterface.(*music_api.AlbumWithTracks)
	if albumWithTracks == nil {
		err := errors.New("AlbumWithTracks is nil")
		errorLogger.Println(err)
		return err
	}

	if albumWithTracks.Album == nil {
		err := errors.New("AlbumWithTracks.Album is nil")
		errorLogger.Println(err)
		return err
	}

	if albumWithTracks.Tracks == nil {
		err := errors.New("AlbumWithTracks.Tracks is nil")
		errorLogger.Println(err)
		return err
	}

	err := validateAlbum(albumWithTracks.Album)
	if err != nil {
		return err
	}

	err = validateTracks(albumWithTracks.Tracks)
	if err != nil {
		return err
	}

	return nil
}

func validateTrackWithAlbumAndArtist(
	trackWithAlbumAndArtistInterface interface{},
) error {
	trackWithAlbumAndArtist :=
		trackWithAlbumAndArtistInterface.(*music_api.TrackWithAlbumAndArtist)
	if trackWithAlbumAndArtist == nil {
		err := errors.New("TrackWithAlbumAndArtist is nil")
		errorLogger.Println(err)
		return err
	}

	if trackWithAlbumAndArtist.Track == nil {
		err := errors.New("TrackWithAlbumAndArtist.Track is nil")
		errorLogger.Println(err)
		return err
	}

	if trackWithAlbumAndArtist.Album == nil {
		err := errors.New("TrackWithAlbumAndArtist.Album is nil")
		errorLogger.Println(err)
		return err
	}

	if trackWithAlbumAndArtist.Artist == nil {
		err := errors.New("TrackWithAlbumAndArtist.Artist is nil")
		errorLogger.Println(err)
		return err
	}

	return nil
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetPrefix("\033[0m")
	warningLogger = log.New(os.Stdout, "\033[33mWARNING: ", log.Flags())
	errorLogger = log.New(os.Stderr, "\033[31mERROR: ", log.Flags())

	containerTool := ""

	_, err := exec.Command("which", "docker").Output()
	if err == nil {
		containerTool = "docker"
	}

	_, err = exec.Command("which", "podman").Output()
	if err == nil {
		containerTool = "podman"
	}

	if containerTool == "" {
		errorLogger.Fatalln("Neither podman nor docker found in PATH")
	}

	log.Println("Building the image with " + containerTool)
	cmd := exec.Command(containerTool, "build", "-t", "music-browser", ".")
	cmd.Dir = "../"

	b, err := cmd.Output()
	if err != nil {
		errorLogger.Fatalln("Couldn't build the image: " + err.Error())
		errorLogger.Fatalln(containerTool + " build output:")
		errorLogger.Fatalln(string(b))
	}

	spotifyClientId := os.Getenv("SPOTIFY_CLIENT_ID")
	if spotifyClientId == "" {
		errorLogger.Fatalln("SPOTIFY_CLIENT_ID environment variable is not set")
	}

	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if spotifyClientSecret == "" {
		errorLogger.Fatalln("SPOTIFY_CLIENT_SECRET environment variable is not set")
	}

	log.Println("Starting the container")
	cmd, err = startContainer(
		containerTool, spotifyClientId, spotifyClientSecret)
	if err != nil {
		errorLogger.Fatalln("Error starting container: " + err.Error())
	}

	log.Println("Running tests")

	getAndValidate(
		"",
		GetRequestGetArtist,
		"https://open.spotify.com/artist/5kbidtcpyRRMdAQUnI1BG4",
		"",
		unmarshalArtistWithAlbums,
		validateArtistWithAlbums,
	)

	getAndValidate(
		"",
		GetRequestGetArtist,
		"https://music.youtube.com/browse/UCevRJ36hDdpdG_4ujQfXfOA",
		"",
		unmarshalArtistWithAlbums,
		validateArtistWithAlbums,
	)

	getAndValidate(
		"",
		GetRequestGetArtist,
		"https://macroblank.bandcamp.com",
		"",
		unmarshalArtistWithAlbums,
		validateArtistWithAlbums,
	)

	getAndValidate(
		"",
		GetRequestGetAlbum,
		"https://open.spotify.com/album/29p0sILcBTXJmhzqJPzcxB",
		"",
		unmarshalAlbumWithTracks,
		validateAlbumWithTracks,
	)

	getAndValidate(
		"",
		GetRequestGetAlbum,
		"https://music.youtube.com/browse/MPREb_AOBS6uUAmis",
		"",
		unmarshalAlbumWithTracks,
		validateAlbumWithTracks,
	)

	getAndValidate(
		"",
		GetRequestGetAlbum,
		"https://macroblank.bandcamp.com/album/--5",
		"",
		unmarshalAlbumWithTracks,
		validateAlbumWithTracks,
	)

	getAndValidate(
		"",
		GetRequestGetTrack,
		"https://open.spotify.com/track/2R3enCnuZSLovqEu1LCl6t",
		"",
		unmarshalTrackWithAlbumAndArtist,
		validateTrackWithAlbumAndArtist,
	)

	getAndValidate(
		"",
		GetRequestGetTrack,
		"https://music.youtube.com/watch?v=ziJUuC4y8s4",
		"",
		unmarshalTrackWithAlbumAndArtist,
		validateTrackWithAlbumAndArtist,
	)

	getAndValidate(
		"",
		GetRequestGetTrack,
		"https://macroblank.bandcamp.com/track/--91",
		"",
		unmarshalTrackWithAlbumAndArtist,
		validateTrackWithAlbumAndArtist,
	)

	getAndValidate(
		SpotifyBackendName,
		GetRequestSearchArtists,
		"",
		"Ne Obliviscaris",
		unmarshalArtists,
		validateArtists,
	)

	getAndValidate(
		YtMusicBackendName,
		GetRequestSearchArtists,
		"",
		"Ne Obliviscaris",
		unmarshalArtists,
		validateArtists,
	)

	getAndValidate(
		BandcampBackendName,
		GetRequestSearchArtists,
		"",
		"Macroblank",
		unmarshalArtists,
		validateArtists,
	)

	getAndValidate(
		SpotifyBackendName,
		GetRequestSearchAlbums,
		"",
		"Exul",
		unmarshalAlbums,
		validateAlbums,
	)

	getAndValidate(
		YtMusicBackendName,
		GetRequestSearchAlbums,
		"",
		"Exul",
		unmarshalAlbums,
		validateAlbums,
	)

	getAndValidate(
		BandcampBackendName,
		GetRequestSearchAlbums,
		"",
		"Death",
		unmarshalAlbums,
		validateAlbums,
	)

	log.Print("Done! Shutting down the container\n\n")

	err = cmd.Process.Signal(os.Interrupt)
	if err != nil {
		errorLogger.Println(
			"Couldn't send SIGINT to the container process: " + err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if !ok {
			errorLogger.Println(
				"Couldn't shut down the container: " + err.Error())
		}
	}

	var passed int
	for _, result := range Results {
		if result.Success {
			passed += 1
			fmt.Printf("\033[32m%s\n", result)
		} else {
			fmt.Printf("\033[31m%s\n", result)
		}
	}
	fmt.Println("\033[0m")
	total := len(Results)
	failed := total - passed

	fmt.Println("Summary:")
	if failed == 0 {
		fmt.Print("\033[32m")
		fmt.Printf("%d/%d - ALL TESTS PASSED", passed, total)
	} else {
		fmt.Print("\033[31m")
		fmt.Printf("%d/%d - %d TESTS FAILED", passed, total, failed)
	}
	fmt.Println("\033[0m")
}
