package grpc_backend

import (
	"context"
	"log"

	"github.com/emilia-miki/music-browser/music_browser/music_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcBackend struct {
	conn   *grpc.ClientConn
	client music_api.MusicApiClient
}

func New(uri string) (*GrpcBackend, error) {
	log.Println(uri)
	conn, err := grpc.Dial(
		uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	log.Println(conn)

	client := music_api.NewMusicApiClient(conn)
	log.Println(client)
	gb := &GrpcBackend{
		conn:   conn,
		client: client,
	}
	return gb, nil
}

func (gb *GrpcBackend) Close() error {
	return gb.conn.Close()
}

func (gb *GrpcBackend) GetArtist(
	url string,
) (*music_api.ArtistWithAlbums, error) {
	result, err := gb.client.GetArtist(
		context.Background(),
		&music_api.Url{
			Url: url,
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (gb *GrpcBackend) GetAlbum(
	url string,
) (*music_api.AlbumWithTracks, error) {
	result, err := gb.client.GetAlbum(
		context.Background(),
		&music_api.Url{
			Url: url,
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (gb *GrpcBackend) GetTrack(
	url string,
) (*music_api.TrackWithAlbumAndArtist, error) {
	log.Println("Getting track " + url + "!!!!!11111")
	result, err := gb.client.GetTrack(
		context.Background(),
		&music_api.Url{
			Url: url,
		},
	)

	if err != nil {
		log.Println("memes and memes and memes....")
		return nil, err
	}

	return result, nil
}

func (gb *GrpcBackend) SearchArtists(
	query string,
) (*music_api.Artists, error) {
	result, err := gb.client.SearchArtists(
		context.Background(),
		&music_api.Query{
			Query: query,
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (gb *GrpcBackend) SearchAlbums(
	query string,
) (*music_api.Albums, error) {
	result, err := gb.client.SearchAlbums(
		context.Background(),
		&music_api.Query{
			Query: query,
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
