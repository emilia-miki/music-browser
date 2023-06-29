package explorer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/emilia-miki/music-browser/music_browser/music_api"
	"github.com/golang/protobuf/proto"
	"github.com/redis/go-redis/v9"
)

const ttl = 3600 * 24 * time.Second

type requestCode int

const (
	getArtistCode requestCode = iota
	getAlbumCode
	getTrackCode
	searchArtistsCode
	searchAlbumsCode
)

type Backend interface {
	GetArtist(id string) (*music_api.ArtistWithAlbums, error)
	GetAlbum(id string) (*music_api.AlbumWithTracks, error)
	GetTrack(id string) (*music_api.TrackWithAlbumAndArtist, error)

	SearchArtists(query string) (*music_api.Artists, error)
	SearchAlbums(query string) (*music_api.Albums, error)
}

type Explorer struct {
	redisClient  *redis.Client
	backends     map[string]Backend
	linkManager  *link_manager.LinkManager
	musicManager *music_manager.MusicManager
}

func New(
	redisClient *redis.Client,
	backends map[string]Backend,
	linkManager *link_manager.LinkManager,
	musicManager *music_manager.MusicManager,
) *Explorer {
	return &Explorer{
		redisClient:  redisClient,
		backends:     backends,
		linkManager:  linkManager,
		musicManager: musicManager,
	}
}

func (e *Explorer) GetArtist(
	backend Backend,
	id string,
) (*music_api.ArtistWithAlbums, error) {
	key := string(getArtistCode) + id
	str, err := e.redisClient.GetEx(
		context.Background(),
		key,
		ttl,
	).Result()

	if err == nil {
		artist := new(music_api.ArtistWithAlbums)
		err = proto.Unmarshal([]byte(str), artist)
		if err != nil {
			return nil, err
		}
		return artist, nil
	} else if !errors.Is(err, redis.Nil) {
		return nil, err
	}

	artist, err := backend.GetArtist(id)
	if err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(artist)
	if err != nil {
		return nil, err
	}

	err = e.redisClient.SetEx(
		context.Background(),
		key,
		string(bytes),
		ttl,
	).Err()
	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (e *Explorer) GetAlbum(
	backend Backend,
	id string,
) (*music_api.AlbumWithTracks, error) {
	key := string(getAlbumCode) + id
	str, err := e.redisClient.GetEx(
		context.Background(),
		key,
		ttl,
	).Result()

	if err == nil {
		album := new(music_api.AlbumWithTracks)
		err = proto.Unmarshal([]byte(str), album)
		if err != nil {
			return nil, err
		}
		return album, nil
	} else if !errors.Is(err, redis.Nil) {
		return nil, err
	}

	album, err := backend.GetAlbum(id)
	if err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(album)
	if err != nil {
		return nil, err
	}

	err = e.redisClient.SetEx(
		context.Background(),
		key,
		string(bytes),
		ttl,
	).Err()
	if err != nil {
		return nil, err
	}

	return album, nil
}

func (e *Explorer) GetTrack(
	backend Backend,
	id string,
) (*music_api.TrackWithAlbumAndArtist, error) {
	key := string(getTrackCode) + id
	str, err := e.redisClient.GetEx(
		context.Background(),
		key,
		ttl,
	).Result()

	if err == nil {
		track := new(music_api.TrackWithAlbumAndArtist)
		err = proto.Unmarshal([]byte(str), track)
		if err != nil {
			return nil, err
		}

		return track, nil
	}

	track, err := backend.GetTrack(id)
	if err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(track)
	if err != nil {
		return nil, err
	}

	err = e.redisClient.SetEx(
		context.Background(),
		key,
		string(bytes),
		ttl,
	).Err()
	if err != nil {
		return nil, err
	}

	return track, nil
}

func (e *Explorer) SearchArtists(
	backend Backend,
	query string,
) (*music_api.Artists, error) {
	key := string(searchArtistsCode) + query
	str, err := e.redisClient.GetEx(
		context.Background(),
		key,
		ttl,
	).Result()

	if err == nil {
		artists := new(music_api.Artists)
		err = proto.Unmarshal([]byte(str), artists)
		if err != nil {
			return nil, err
		}

		return artists, nil
	}

	artists, err := backend.SearchArtists(query)
	if err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(artists)
	if err != nil {
		return nil, err
	}

	err = e.redisClient.SetEx(
		context.Background(),
		key,
		string(bytes),
		ttl,
	).Err()
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (e *Explorer) SearchAlbums(
	backend Backend,
	query string,
) (*music_api.Albums, error) {
	key := string(searchAlbumsCode) + query
	str, err := e.redisClient.GetEx(
		context.Background(),
		key,
		ttl,
	).Result()

	if err == nil {
		albums := new(music_api.Albums)
		err = proto.Unmarshal([]byte(str), albums)
		if err != nil {
			return nil, err
		}

		return albums, nil
	}

	albums, err := backend.SearchAlbums(query)
	if err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(albums)
	if err != nil {
		return nil, err
	}

	err = e.redisClient.SetEx(
		context.Background(),
		key,
		string(bytes),
		ttl,
	).Err()
	if err != nil {
		return nil, err
	}

	return albums, nil
}
