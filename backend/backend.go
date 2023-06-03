package backend

import "github.com/emilia-miki/music-browser/models"

type MusicExplorer interface {
	SearchTracks(query string) []models.Track
	SearchAlbums(query string) []models.Album
	SearchArtists(query string) []models.Artist
}

type MusicDownloader interface {
	DownloadTrack(url string)
}
