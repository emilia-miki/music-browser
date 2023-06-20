package spotify_backend

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/emilia-miki/music-browser/music_browser/environment"
	"github.com/emilia-miki/music-browser/music_browser/music_api"
)

// The main public type
type SpotifyBackend struct {
	secrets     environment.SpotifySecrets
	accessToken *accessToken
}

// The constructor
func New(
	secrets environment.SpotifySecrets,
) *SpotifyBackend {
	return &SpotifyBackend{
		secrets:     secrets,
		accessToken: nil,
	}
}

// Spotify API Access Token
type accessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint32 `json:"expires_in"`
}

// Public methods of SpotifyBackend

func (sb *SpotifyBackend) GetArtist(
	url string,
) (*music_api.ArtistWithAlbums, error) {
	var result music_api.ArtistWithAlbums

	id := extractIdFromUrl(url)
	uri := "https://api.spotify.com/v1/artists/" + id

	jsonData := sb.getFromApi(uri)

	var artistData struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Images []image `json:"images"`
		Name   string  `json:"name"`
	}
	json.Unmarshal(jsonData, &artistData)

	result.Artist = &music_api.Artist{
		Url:      artistData.ExternalUrls.Spotify,
		ImageUrl: getBestImageUrl(artistData.Images),
		Name:     artistData.Name,
	}

	uri = fmt.Sprintf("https://api.spotify.com/v1/artists/%s/albums", id)
	jsonData = sb.getFromApi(uri)

	var resp struct {
		Items []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Images  []image `json:"images"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
			}
			Name        string `json:"name"`
			ReleaseDate string `json:"release_date"`
		} `json:"items"`
	}
	json.Unmarshal(jsonData, &resp)

	result.Albums = &music_api.Albums{
		Albums: make([]*music_api.Album, len(resp.Items)),
	}
	for i, album := range resp.Items {
		artistUrls := make([]string, len(album.Artists))
		for i, artistLink := range album.Artists {
			artistUrls[i] = artistLink.ExternalUrls.Spotify
		}

		result.Albums.Albums[i] = &music_api.Album{
			Url:        album.ExternalUrls.Spotify,
			ImageUrl:   getBestImageUrl(album.Images),
			ArtistUrls: artistUrls,
			Name:       album.Name,
			Year:       getYearFromDate(album.ReleaseDate),
		}
	}

	return &result, nil
}

func (sb *SpotifyBackend) GetAlbum(
	url string,
) (*music_api.AlbumWithTracks, error) {
	var result music_api.AlbumWithTracks

	id := extractIdFromUrl(url)
	jsonData := sb.getFromApi("https://api.spotify.com/v1/albums/" + id)

	var albumData struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Images  []image `json:"images"`
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"artists"`
		Tracks struct {
			Items []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Images  []image `json:"images"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
				}
				Name                 string `json:"name"`
				DurationMilliseconds uint32 `json:"duration_ms"`
			} `json:"items"`
		} `json:"tracks"`
		Name        string `json:"name"`
		ReleaseDate string `json:"release_date"`
	}
	json.Unmarshal(jsonData, &albumData)

	artistUrls := make([]string, len(albumData.Artists))
	for i, artistLink := range albumData.Artists {
		artistUrls[i] = artistLink.ExternalUrls.Spotify
	}

	albumImageUrl := getBestImageUrl(albumData.Images)
	result.Album = &music_api.Album{
		Url:        albumData.ExternalUrls.Spotify,
		ImageUrl:   albumImageUrl,
		ArtistUrls: artistUrls,
		Name:       albumData.Name,
		Year:       getYearFromDate(albumData.ReleaseDate),
	}

	result.Tracks = &music_api.Tracks{
		Tracks: make([]*music_api.Track, len(albumData.Tracks.Items)),
	}
	for i, trackData := range albumData.Tracks.Items {
		artistUrls := make([]string, len(trackData.Artists))
		for i, artistLink := range trackData.Artists {
			artistUrls[i] = artistLink.ExternalUrls.Spotify
		}

		result.Tracks.Tracks[i] = &music_api.Track{
			Url:             trackData.ExternalUrls.Spotify,
			ImageUrl:        albumImageUrl,
			ArtistUrls:      artistUrls,
			AlbumUrl:        albumData.ExternalUrls.Spotify,
			Name:            trackData.Name,
			DurationSeconds: trackData.DurationMilliseconds / 1000,
		}
	}

	return &result, nil
}

func (sb *SpotifyBackend) GetTrack(
	url string,
) (*music_api.TrackWithAlbumAndArtists, error) {
	id := extractIdFromUrl(url)
	jsonResponse := sb.getFromApi("https://api.spotify.com/v1/tracks/" + id)

	var resp struct {
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"spotify"`
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"artists"`
		Album struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Images      []image `json:"images"`
			Name        string  `json:"name"`
			ReleaseDate string  `json:"release_date"`
		} `json:"album"`
		Name                 string `json:"name"`
		DurationMilliseconds uint32 `json:"duration_ms"`
	}
	json.Unmarshal(jsonResponse, &resp)

	imageUrl := getBestImageUrl(resp.Album.Images)

	artistList := make([]*music_api.Artist, len(resp.Artists))
	artistUrls := make([]string, len(resp.Artists))
	for i, artist := range resp.Artists {
		artistUrls[i] = artist.ExternalUrls.Spotify
		jsonResponse := sb.getFromApi("https://api.spotify.com/v1/artists/" + artist.Id)
		var respWithImages struct {
			Images []image `json:"images"`
		}
		json.Unmarshal(jsonResponse, &respWithImages)
		artistList[i] = &music_api.Artist{
			Url:      artist.ExternalUrls.Spotify,
			ImageUrl: getBestImageUrl(respWithImages.Images),
			Name:     artist.Name,
		}
	}
	artists := &music_api.Artists{
		Artists: artistList,
	}

	album := &music_api.Album{
		Url:        resp.Album.ExternalUrls.Spotify,
		ImageUrl:   imageUrl,
		ArtistUrls: artistUrls,
		Name:       resp.Album.Name,
		Year:       getYearFromDate(resp.Album.ReleaseDate),
	}

	track := &music_api.Track{
		Url:             url,
		ImageUrl:        imageUrl,
		ArtistUrls:      artistUrls,
		AlbumUrl:        resp.Album.ExternalUrls.Spotify,
		Name:            resp.Name,
		DurationSeconds: resp.DurationMilliseconds / 1000,
	}

	return &music_api.TrackWithAlbumAndArtists{
		Artists: artists,
		Album:   album,
		Track:   track,
	}, nil
}

func (sb *SpotifyBackend) SearchArtists(
	query string,
) (*music_api.Artists, error) {
	jsonResponse := sb.search("artist", query)
	var searchData struct {
		Artists struct {
			Items []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Images []image `json:"images"`
				Name   string  `json:"name"`
			} `json:"items"`
		} `json:"artists"`
	}
	json.Unmarshal(jsonResponse, &searchData)

	artists := make([]*music_api.Artist, len(searchData.Artists.Items))
	for i, artist := range searchData.Artists.Items {
		artists[i] = &music_api.Artist{
			Url:      artist.ExternalUrls.Spotify,
			ImageUrl: getBestImageUrl(artist.Images),
			Name:     artist.Name,
		}
	}

	return &music_api.Artists{
		Artists: artists,
	}, nil
}

func (sb *SpotifyBackend) SearchAlbums(
	query string,
) (*music_api.Albums, error) {
	jsonResponse := sb.search("album", query)
	var searchData struct {
		Albums struct {
			Items []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Images  []image `json:"images"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
				} `json:"artists"`
				Name        string `json:"name"`
				ReleaseDate string `json:"release_date"`
			} `json:"items"`
		} `json:"albums"`
	}
	json.Unmarshal(jsonResponse, &searchData)

	albums := make([]*music_api.Album, len(searchData.Albums.Items))
	for i, album := range searchData.Albums.Items {
		artistUrls := make([]string, len(album.Artists))
		for i, artistLink := range album.Artists {
			artistUrls[i] = artistLink.ExternalUrls.Spotify
		}

		albums[i] = &music_api.Album{
			Url:        album.ExternalUrls.Spotify,
			ImageUrl:   getBestImageUrl(album.Images),
			ArtistUrls: artistUrls,
			Name:       album.Name,
			Year:       getYearFromDate(album.ReleaseDate),
		}
	}

	return &music_api.Albums{
		Albums: albums,
	}, nil
}

// Spotify API helper methods

func (sb *SpotifyBackend) search(searchType string, query string) []byte {
	sb.getAccessToken()

	formData := url.Values{}
	formData.Set("type", searchType)
	formData.Set("q", query)
	params := formData.Encode()

	return sb.getFromApi("https://api.spotify.com/v1/search?" + params)
}

func (sb *SpotifyBackend) getFromApi(url string) []byte {
	sb.getAccessToken()

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+sb.accessToken.AccessToken)
	client := http.DefaultClient
	resp, _ := client.Do(req)
	jsonString, _ := ioutil.ReadAll(resp.Body)
	return jsonString
}

func (sb *SpotifyBackend) refreshAccessToken() {
	sb.accessToken = nil
	sb.getAccessToken()
}

func (sb *SpotifyBackend) getAccessToken() {
	if sb.accessToken != nil {
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
		[]byte(sb.secrets.ClientId+":"+sb.secrets.ClientSecret)))

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

	sb.accessToken = new(accessToken)
	json.Unmarshal(respBody, sb.accessToken)
}

// Utility functions

func extractIdFromUrl(url string) string {
	splits := strings.Split(url, "/")
	return splits[len(splits)-1]
}

type image struct {
	Url    string `json:"url"`
	Height uint32 `json:"height"`
	Width  uint32 `json:"width"`
}

func getBestImageUrl(images []image) string {
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

func getYearFromDate(date string) uint32 {
	yearString, _, _ := strings.Cut(date, "-")
	year, _ := strconv.ParseUint(yearString, 10, 16)
	return uint32(year)
}
