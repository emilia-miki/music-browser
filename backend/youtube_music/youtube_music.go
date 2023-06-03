package youtube_music

import (
	"github.com/emilia-miki/music-browser/backend/music_explorer_cache"
	"github.com/emilia-miki/music-browser/models"
	"google.golang.org/grpc"
)

type MusicExplorer struct {
	Cache          music_explorer_cache.MusicExplorerCache
	YtMusicApiConn *grpc.ClientConn
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
