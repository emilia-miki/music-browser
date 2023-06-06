package youtube_music

import (
	"context"

	"github.com/emilia-miki/music-browser/music_browser/backend/music_explorer_cache"
	"github.com/emilia-miki/music-browser/music_browser/music_api"
)

type MusicExplorer struct {
	Cache  music_explorer_cache.MusicExplorerCache
	Client music_api.MusicApiClient
}

func (me *MusicExplorer) SearchArtists(query string) []*music_api.Artist {
	// TODO: use redis cache

	searchResult, _ := me.Client.SearchArtists(context.Background(), &music_api.Query{
		Query: query,
	})

	return searchResult.Items
}

func (me *MusicExplorer) SearchAlbums(query string) []*music_api.Album {
	// TODO: use redis cache

	searchResult, _ := me.Client.SearchAlbums(context.Background(), &music_api.Query{
		Query: query,
	})

	return searchResult.Items
}
