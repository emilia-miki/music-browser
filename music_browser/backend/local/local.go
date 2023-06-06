package local

import (
	"database/sql"

	"github.com/emilia-miki/music-browser/music_browser/music_api"
)

type MusicExplorer struct {
	PostgresDB *sql.DB
}

func (*MusicExplorer) SearchAlbums(query string) []*music_api.Album {
	// TODO
	return []*music_api.Album{}
}

func (*MusicExplorer) SearchArtists(query string) []*music_api.Artist {
	// TODO
	return []*music_api.Artist{}
}
