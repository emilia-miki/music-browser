package grpc_music_explorer_with_cache

import (
	"context"

	"github.com/emilia-miki/music-browser/music_browser/music_api"
	"github.com/go-redis/cache/v9"
)

type GrpcMusicExplorerWithCache struct {
	Cache  *cache.Cache
	Client music_api.MusicApiClient
}

func (me *GrpcMusicExplorerWithCache) SearchArtists(query string) []*music_api.Artist {
	artists := new(*music_api.Artists)
	me.Cache.Once(&cache.Item{
		Key:   query,
		Value: artists,
		Do: func(item *cache.Item) (interface{}, error) {
			return me.Client.SearchArtists(
				context.Background(), &music_api.Query{Query: query})
		},
	})

	return (*artists).Items
}

func (me *GrpcMusicExplorerWithCache) SearchAlbums(query string) []*music_api.Album {
	albums := new(*music_api.Albums)
	me.Cache.Once(&cache.Item{
		Key:   query,
		Value: albums,
		Do: func(item *cache.Item) (interface{}, error) {
			return me.Client.SearchAlbums(
				context.Background(), &music_api.Query{Query: query})
		},
	})

	return (*albums).Items
}
