package music_explorer_cache

import (
	"github.com/emilia-miki/music-browser/models"
	"github.com/redis/go-redis/v9"
)

type MusicExplorerCache struct {
	RedisClient *redis.Client
}

func (*MusicExplorerCache) GetTracksFromCache(query string) (bool, []models.Track) {
	// TODO
	return false, []models.Track{}
}

func (*MusicExplorerCache) GetAlbumsFromCache(query string) (bool, []models.Album) {
	// TODO
	return false, []models.Album{}
}

func (*MusicExplorerCache) GetArtistsFromCache(query string) (bool, []models.Artist) {
	// TODO
	return false, []models.Artist{}
}
