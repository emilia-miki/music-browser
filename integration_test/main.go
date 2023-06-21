package main

import (
	"bufio"
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

	music_browser "github.com/emilia-miki/music-browser/music_browser"
	"github.com/emilia-miki/music-browser/music_browser/explorer"
	"github.com/emilia-miki/music-browser/music_browser/music_api"
)

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
		"localhost:3333/?backend=%s&type=%s&url=%s&query=%s", b, t, u, q)
	resp, err := http.Get(requestUri)
	if err != nil {
		return nil, errors.New("Unable to send a GET request: " + err.Error())
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(
			"Unable to read the response body: " + err.Error())
	}

	return bytes, nil
}

func validateArtist(artist *music_api.Artist) error {
	if artist == nil {
		err := errors.New("artist is nil")
		log.Println("ERROR: " + err.Error())
		return err
	}

	if artist.Name == nil {
		err := errors.New("artist is nil")
		log.Println("ERROR: " + err.Error())
		return err
	}

	if artist.Url == nil {
		log.Println("WARNING: artist.Url is nil")
	}

	if artist.ImageUrl == nil {
		log.Println("WARNING: artist.ImageUrl is nil")
	}

	return nil
}

func validateArtists(artists *music_api.Artists) error {
	if artists == nil {
		err := errors.New("artists is nil")
		log.Println("ERROR: " + err.Error())
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
		err := errors.New("album is nil")
		log.Println("ERROR: " + err.Error())
		return err
	}

	if album.Name == nil {
		err := errors.New("album is nil")
		log.Println("ERROR: " + err.Error())
		return err
	}

	if album.Url == nil {
		log.Println("WARNING: album.Url is nil")
	}

	if album.ImageUrl == nil {
		log.Println("WARNING: album.ImageUrl is nil")
	}

	if album.ArtistUrl == nil {
		log.Println("WARNING: album.ArtistUrl is nil")
	}

	if album.Year == nil {
		log.Println("WARNING: album.Year is nil")
	}

	return nil
}

func validateAlbums(albums *music_api.Albums) error {
	if albums == nil {
		err := errors.New("albums is nil")
		log.Println("ERROR: " + err.Error())
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
		err := errors.New("track is nil")
		log.Println("ERROR: " + err.Error())
		return err
	}

	if track.Name == nil {
		err := errors.New("track.Name is nil")
		log.Println("ERROR: " + err.Error())
		return err
	}

	if track.ImageUrl == nil {
		log.Println("WARNING: track.ImageUrl is nil")
	}

	if track.ArtistUrl == nil {
		log.Println("WARNING: track.ArtistUrl is nil")
	}

	if track.AlbumUrl == nil {
		log.Println("WARNING: track.AlbumUrl is nil")
	}

	if track.DurationSeconds == nil {
		log.Println("WARNING: track.DurationSeconds is nil")
	}

	return nil
}

func validateTracks(tracks *music_api.Tracks) error {
	if tracks == nil {
		err := errors.New("tracks is nil")
		log.Println("ERROR: " + err.Error())
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
	artistWithAlbums *music_api.ArtistWithAlbums,
) error {
	if artistWithAlbums == nil {
		err := errors.New("artistWithAlbums is nil")
		log.Println("ERROR: " + err.Error())
		return err
	}

	if artistWithAlbums.Artist == nil {
		err := errors.New("artistWithAlbums.Artist is nil")
		log.Println("ERROR: " + err.Error())
		return err
	}

	if artistWithAlbums.Albums == nil {
		err := errors.New("artistWithAlbums.Albums is nil")
		log.Println("ERROR: " + err.Error())
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

func validateAlbumWithTracks(albumWithTracks *music_api.AlbumWithTracks) error {
	if albumWithTracks == nil {
		err := errors.New("albumWithTracks is nil")
		log.Println("ERROR: " + err.Error())
		return err
	}

	if albumWithTracks.Album == nil {
		err := errors.New("albumWithTracks.Album is nil")
		log.Println("ERROR: " + err.Error())
		return err
	}

	if albumWithTracks.Tracks == nil {
		err := errors.New("albumWithTracks.Tracks is nil")
		log.Println("ERROR: " + err.Error())
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

func startContainer(
	containerTool string,
	spotifyClientId string,
	spotifyClientSecret string,
) (*exec.Cmd, error) {
	cmd := exec.CommandContext(context.Background(),
		containerTool, "run", "-p", "3333:3333",
		"-e", fmt.Sprintf("SPOTIFY_CLIENT_ID=%s", spotifyClientId),
		"-e", fmt.Sprintf("SPOTIFY_CLIENT_SECRET=%s", spotifyClientSecret),
		"music-browser")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.New("unable to open pipe")
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
					"unexpected EOF encountered while capturing container output." +
						" Output so far: " + output)
			} else {
				return nil, errors.New(
					"unexpected error encounteres while capturing container " +
						"output. Output so far: " + output)
			}
		}

		line = scanner.Text()
		output += line
		if strings.Contains(line, "Server listening on port 3333") {
			go func() {
				for {
					res = scanner.Scan()
					if res == false {
						break
					}
				}
			}()
			break
		}
	}

	return cmd, nil
}

func main() {
	containerTool := ""

	_, err := exec.Command("which", "docker").Output()
	if err != nil {
		containerTool = "docker"
	}

	_, err = exec.Command("which", "podman").Output()
	if err == nil {
		containerTool = "podman"
	}

	if containerTool == "" {
		log.Fatalln("Neither podman nor docker found in PATH")
	}

	cmd := exec.Command(containerTool, "build", "-t", "music-browser", ".")
	cmd.Dir = "../"
	_, err = cmd.Output()
	if err != nil {
		log.Fatalln("Couldn't build the image: " + err.Error())
	}

	spotifyClientId := os.Getenv("SPOTIFY_CLIENT_ID")
	if spotifyClientId == "" {
		log.Fatalln("SPOTIFY_CLIENT_ID environment variable is not set")
	}

	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if spotifyClientSecret == "" {
		log.Fatalln("SPOTIFY_CLIENT_SECRET environment variable is not set")
	}

	cmd, err = startContainer(
		containerTool, spotifyClientId, spotifyClientSecret)
	if err != nil {
		log.Fatalln("Error starting container: " + err.Error())
	}

	resp, err := sendGetRequest(
		"",
		music_browser.GetRequestGetArtist,
		"https://open.spotify.com/browse/5kbidtcpyRRMdAQUnI1BG4",
		"")
	artist := new(music_api.Artist)
	err = json.Unmarshal(resp, artist)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateArtist(artist)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		"",
		music_browser.GetRequestGetArtist,
		"https://music.youtube.com/browse/UCevRJ36hDdpdG_4ujQfXfOA",
		"")
	artist = new(music_api.Artist)
	err = json.Unmarshal(resp, artist)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateArtist(artist)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		"",
		music_browser.GetRequestGetArtist,
		"https://macroblank.bandcamp.com",
		"")
	artist = new(music_api.Artist)
	err = json.Unmarshal(resp, artist)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateArtist(artist)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		"",
		music_browser.GetRequestGetAlbum,
		"https://open.spotify.com/browse/29p0sILcBTXJmhzqJPzcxB",
		"")
	album := new(music_api.Album)
	err = json.Unmarshal(resp, album)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateAlbum(album)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		"",
		music_browser.GetRequestGetAlbum,
		"https://music.youtube.com/browse/UCevRJ36hDdpdG_4ujQfXfOA",
		"")
	album = new(music_api.Album)
	err = json.Unmarshal(resp, album)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateAlbum(album)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		"",
		music_browser.GetRequestGetAlbum,
		"https://macroblank.bandcamp.com/album/--5",
		"")
	album = new(music_api.Album)
	err = json.Unmarshal(resp, album)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateAlbum(album)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		"",
		music_browser.GetRequestGetTrack,
		"https://open.spotify.com/browse/track/2R3enCnuZSLovqEu1LCl6t",
		"")
	track := new(music_api.Track)
	err = json.Unmarshal(resp, track)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateTrack(track)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		"",
		music_browser.GetRequestGetTrack,
		"https://music.youtube.com/watch?v=ziJUuC4y8s4",
		"")
	track = new(music_api.Track)
	err = json.Unmarshal(resp, track)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateTrack(track)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		"",
		music_browser.GetRequestGetTrack,
		"https://macroblank.bandcamp.com/track/--91",
		"")
	track = new(music_api.Track)
	err = json.Unmarshal(resp, track)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateTrack(track)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		explorer.SpotifyBackendName,
		music_browser.GetRequestSearchArtists,
		"",
		"Ne Obliviscaris")
	artists := new(music_api.Artists)
	err = json.Unmarshal(resp, artists)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateArtists(artists)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		explorer.YtMusicBackendName,
		music_browser.GetRequestSearchArtists,
		"",
		"Ne Obliviscaris")
	artists = new(music_api.Artists)
	err = json.Unmarshal(resp, artists)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateArtists(artists)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		explorer.BandcampBackendName,
		music_browser.GetRequestSearchArtists,
		"",
		"Macroblank")
	artists = new(music_api.Artists)
	err = json.Unmarshal(resp, artists)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateArtists(artists)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		explorer.SpotifyBackendName,
		music_browser.GetRequestSearchAlbums,
		"",
		"Exul")
	albums := new(music_api.Albums)
	err = json.Unmarshal(resp, albums)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateAlbums(albums)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		explorer.YtMusicBackendName,
		music_browser.GetRequestSearchAlbums,
		"",
		"Exul")
	albums = new(music_api.Albums)
	err = json.Unmarshal(resp, albums)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateAlbums(albums)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	resp, err = sendGetRequest(
		explorer.BandcampBackendName,
		music_browser.GetRequestSearchAlbums,
		"",
		"Death")
	albums = new(music_api.Albums)
	err = json.Unmarshal(resp, albums)
	if err != nil {
		log.Println("Error unmarshaling response: " + err.Error())
	}
	err = validateAlbums(albums)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		cmd.Wait()
		cmd.Cancel()
		cmd, err = startContainer(
			containerTool, spotifyClientId, spotifyClientSecret)
		if err != nil {
			log.Fatalln("Error restarting container: " + err.Error())
		}
	}

	cmd.Cancel()
	cmd.Wait()
	log.Println("Done!")
}
