package backend

import (
	"github.com/emilia-miki/music-browser/music_browser/music_api"
)

type MusicExplorer interface {
	SearchAlbums(query string) []*music_api.Album
	SearchArtists(query string) []*music_api.Artist
}

type MusicDownloader interface {
	DownloadTrack(url string)
}
