package spotify

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/emilia-miki/music-browser/backend/music_explorer_cache"
	"github.com/emilia-miki/music-browser/environment"
	"github.com/emilia-miki/music-browser/models"
)

type artist struct {
	Name         string `json:"name"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

func (a artist) mapToModel() models.Artist {
	return models.Artist{
		Name: a.Name,
		Url:  a.ExternalUrls.Spotify,
	}
}

func mapArtistsToModel(items []artist) []models.Artist {
	result := make([]models.Artist, len(items))
	for i, item := range items {
		result[i] = item.mapToModel()
	}

	return result
}

type album struct {
	Name         string `json:"name"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Images []struct {
		Url    string `json:"url"`
		Height uint32 `json:"height"`
		Width  uint32 `json:"width"`
	} `json:"images"`
	Tags        []string `json:"genres"`
	ReleaseDate string   `json:"release_date"`
	Artists     []artist `json:"artists"`
	NumTracks   uint8    `json:"total_tracks"`
}

func (a album) mapToModel() models.Album {
	artists := mapArtistsToModel(a.Artists)

	imageUrl := ""
	if len(a.Images) > 0 {
		// IMPROVEMENT: take the image with the best resolution instead of the first one
		imageUrl = a.Images[0].Url
	}

	return models.Album{
		Name:        a.Name,
		Url:         a.ExternalUrls.Spotify,
		ImageUrl:    imageUrl,
		Tags:        a.Tags,
		ReleaseDate: a.ReleaseDate,
		Artists:     artists,
		NumTracks:   a.NumTracks,
	}
}

func mapAlbumsToModel(items []album) []models.Album {
	result := make([]models.Album, len(items))
	for i, item := range items {
		result[i] = item.mapToModel()
	}

	return result
}

type track struct {
	Name         string `json:"name"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Album album `json:"album"`
}

func (t track) mapToModel() models.Track {
	return models.Track{
		Name:  t.Name,
		Url:   t.ExternalUrls.Spotify,
		Album: t.Album.mapToModel(),
	}
}

func mapTracksToModel(items []track) []models.Track {
	result := make([]models.Track, len(items))
	for i, item := range items {
		result[i] = item.mapToModel()
	}

	return result
}

type SearchResponse struct {
	Artists struct {
		Items []artist `json:"items"`
	} `json:"artists"`
	Albums struct {
		Items []album `json:"items"`
	} `json:"albums"`
	Tracks struct {
		Items []track `json:"items"`
	} `json:"tracks"`
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

	log.Println(req.URL)
	log.Println(req.Header)
	log.Println(formData.Encode())

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

	log.Println(resp.Status)
	log.Println(resp.Header)
	log.Println(string(respBody))

	me.accessToken = &accessToken{}
	json.Unmarshal(respBody, me.accessToken)
}

func (me *MusicExplorer) refreshAccessToken() {
	me.accessToken = nil
	me.getAccessToken()
}

func (me *MusicExplorer) search(searchType string, query string) SearchResponse {
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

	var jsonResponse SearchResponse
	err = json.Unmarshal(respBody, &jsonResponse)
	if err != nil {
		log.Fatalf("Error parsing tracks from spotify's response JSON: %s\n", err)
	}

	return jsonResponse
}

func (me *MusicExplorer) SearchArtists(query string) []models.Artist {
	jsonResponse := me.search("artist", query)
	return mapArtistsToModel(jsonResponse.Artists.Items)
}

func (me *MusicExplorer) SearchAlbums(query string) []models.Album {
	jsonResponse := me.search("album", query)
	return mapAlbumsToModel(jsonResponse.Albums.Items)
}

func (me *MusicExplorer) SearchTracks(query string) []models.Track {
	jsonResponse := me.search("track", query)
	return mapTracksToModel(jsonResponse.Tracks.Items)
}
