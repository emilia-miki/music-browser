package spotify

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/emilia-miki/music-browser/music_browser/backend/music_explorer_cache"
	"github.com/emilia-miki/music-browser/music_browser/environment"
	"github.com/emilia-miki/music-browser/music_browser/music_api"
)

type image struct {
	Url    string `json:"url"`
	Height uint32 `json:"height"`
	Width  uint32 `json:"width"`
}

type item struct {
	Name         *string `json:"name"`
	ExternalUrls *struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Images []image `json:"images"`
	Id     *string `json:"id"`
}

type artistLink struct {
	Name string
	Url  string
}

type track struct {
	Name         string `json:"name"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	DurationMilliSeconds uint32 `json:"duration_ms"`
	Artists              []struct {
		Name         string `json:"name"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
	} `json:"artists"`
}

type album struct {
	Name         string `json:"name"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Images      []image `json:"images"`
	ReleaseDate string  `json:"release_date"`
	Tracks      struct {
		Items []track `json:"items"`
	} `json:"tracks"`
	Artists []struct {
		Name         string `json:"name"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
	} `json:"artists"`
}

type artist struct {
	Name         string `json:"name"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Images []image `json:"images"`
}

func getBestImage(images []image) string {
	maxHeight := uint32(0)
	url := ""
	for _, image := range images {
		if image.Height > maxHeight {
			maxHeight = image.Height
			url = image.Url
		}
	}

	return url
}

func parseTracks(
	items []track,
	album string,
	albumUrl string,
	albumImageUrl string) []*music_api.Track {
	tracks := make([]*music_api.Track, len(items))
	for i, item := range items {
		artists := make([]*music_api.ArtistLink, len(item.Artists))
		for j, artist := range item.Artists {
			artists[j] = &music_api.ArtistLink{
				Name: artist.Name,
				Url:  artist.ExternalUrls.Spotify,
			}
		}
		tracks[i] = &music_api.Track{
			Name:            item.Name,
			Url:             item.ExternalUrls.Spotify,
			ImageUrl:        albumImageUrl,
			DurationSeconds: item.DurationMilliSeconds / 1000,
			Album:           album,
			AlbumUrl:        albumUrl,
			Artists:         artists,
		}
	}

	return tracks
}

func (me *MusicExplorer) mapItemsToAlbums(items []item) []*music_api.Album {
	albums := make([]*music_api.Album, len(items))
	for i, item := range items {
		req, _ := http.NewRequest("GET",
			"https://api.spotify.com/v1/albums/"+*item.Id, nil)
		req.Header.Set("Authorization", "Bearer "+me.accessToken.AccessToken)
		client := http.DefaultClient
		resp, _ := client.Do(req)
		jsonString, _ := ioutil.ReadAll(resp.Body)

		var album album
		json.Unmarshal(jsonString, &album)

		image := getBestImage(album.Images)

		year, _ := strconv.ParseUint(
			strings.Split(album.ReleaseDate, "-")[0], 10, 32)
		durationMilliSeconds := uint32(0)
		for _, track := range album.Tracks.Items {
			durationMilliSeconds += track.DurationMilliSeconds
		}
		tracks := parseTracks(
			album.Tracks.Items, album.Name, album.ExternalUrls.Spotify, image)
		artists := make([]*music_api.ArtistLink, len(album.Artists))
		for j, artist := range album.Artists {
			artists[j] = &music_api.ArtistLink{
				Name: artist.Name,
				Url:  artist.ExternalUrls.Spotify,
			}
		}
		albums[i] = &music_api.Album{
			Name:            album.Name,
			Url:             album.ExternalUrls.Spotify,
			ImageUrl:        image,
			Year:            uint32(year),
			DurationSeconds: durationMilliSeconds / 1000,
			Tracks:          tracks,
			Artists:         artists,
		}
	}

	return albums
}

func (me *MusicExplorer) mapItemsToArtists(items []item) []*music_api.Artist {
	artists := make([]*music_api.Artist, len(items))
	for i, artist := range items {
		req, _ := http.NewRequest("GET",
			"https://api.spotify.com/v1/artists/"+*artist.Id+"/albums", nil)
		req.Header.Set("Authorization", "Bearer "+me.accessToken.AccessToken)
		client := http.DefaultClient
		resp, _ := client.Do(req)
		jsonAlbums, _ := ioutil.ReadAll(resp.Body)

		var jsonResponse struct {
			Items []item `json:"items"`
		}
		json.Unmarshal(jsonAlbums, &jsonResponse)
		artists[i] = &music_api.Artist{
			Name:     *artist.Name,
			Url:      artist.ExternalUrls.Spotify,
			ImageUrl: getBestImage(artist.Images),
			Albums:   me.mapItemsToAlbums(jsonResponse.Items),
		}
	}

	return artists
}

type searchResponse struct {
	Artists struct {
		Items []item `json:"items"`
	} `json:"artists"`
	Albums struct {
		Items []item `json:"items"`
	} `json:"albums"`
}

type accessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint32 `json:"expires_in"`
}

type MusicExplorer struct {
	Cache       music_explorer_cache.MusicExplorerCache
	Secrets     environment.SpotifySecrets
	accessToken *accessToken
}

func (me *MusicExplorer) getAccessToken() {
	if me.accessToken != nil {
		return
	}

	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	body := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", body)
	if err != nil {
		log.Fatalf("Failed to create POST request for getting an access token: %s\n", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(
		[]byte(me.Secrets.ClientId+":"+me.Secrets.ClientSecret)))

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request %s: %s\n", req.URL, err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %s\n", err)
	}
	resp.Body.Close()

	me.accessToken = &accessToken{}
	json.Unmarshal(respBody, me.accessToken)
}

func (me *MusicExplorer) refreshAccessToken() {
	me.accessToken = nil
	me.getAccessToken()
}

func (me *MusicExplorer) search(searchType string, query string) searchResponse {
	// TODO: USE REDIS CACHE HERE

	me.getAccessToken()

	formData := url.Values{}
	formData.Set("type", searchType)
	formData.Set("q", query)
	params := formData.Encode()

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/search?"+params, nil)
	if err != nil {
		log.Fatalf("Failed to create POST request for getting an access token: %s\n", err)
	}

	req.Header.Set("Authorization", "Bearer "+me.accessToken.AccessToken)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request %s: %s\n", req.URL, err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %s\n", err)
	}
	resp.Body.Close()

	var jsonResponse searchResponse
	err = json.Unmarshal(respBody, &jsonResponse)
	if err != nil {
		log.Fatalf("Error parsing tracks from spotify's response JSON: %s\n", err)
	}

	return jsonResponse
}

func (me *MusicExplorer) SearchArtists(query string) []*music_api.Artist {
	jsonResponse := me.search("artist", query)
	return me.mapItemsToArtists(jsonResponse.Artists.Items)
}

func (me *MusicExplorer) SearchAlbums(query string) []*music_api.Album {
	jsonResponse := me.search("album", query)
	return me.mapItemsToAlbums(jsonResponse.Albums.Items)
}
