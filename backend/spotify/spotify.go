package spotify

import (
	"github.com/emilia-miki/music-browser/backend/music_explorer_cache"
	"github.com/emilia-miki/music-browser/environment"
	"github.com/emilia-miki/music-browser/models"
)

type accessToken struct {
	accessToken string
	tokenType   string
	expiresIn   uint32
}

type MusicExplorer struct {
	Cache       music_explorer_cache.MusicExplorerCache
	Secrets     environment.SpotifySecrets
	accessToken accessToken
}

func (*MusicExplorer) SearchTracks(query string) []models.Track {
	// TODO
	return []models.Track{}
}

func (*MusicExplorer) SearchAlbums(query string) []models.Album {
	// TODO
	return []models.Album{}
}

func (*MusicExplorer) SearchArtists(query string) []models.Artist {
	// TODO
	return []models.Artist{}
}
