package local

import (
	"database/sql"

	"github.com/emilia-miki/music-browser/models"
)

type MusicExplorer struct {
	PostgresDB *sql.DB
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
