package music_explorer_cache

import (
	"github.com/emilia-miki/music-browser/music_browser/music_api"
	"github.com/redis/go-redis/v9"
)

type MusicExplorerCache struct {
	RedisClient *redis.Client
}

func (*MusicExplorerCache) GetAlbumsFromCache(query string) (bool, []*music_api.Album) {
	// TODO
	return false, []*music_api.Album{}
}

func (*MusicExplorerCache) GetArtistsFromCache(query string) (bool, []*music_api.Artist) {
	// TODO
	return false, []*music_api.Artist{}
}
