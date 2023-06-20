package local_backend

import (
	"database/sql"
	"github.com/emilia-miki/music-browser/music_browser/music_api"
	_ "github.com/lib/pq"
)

type LocalBackend struct {
	pgDB *sql.DB
}

func New(postgresConnectionString string) (*LocalBackend, error) {
	db, err := sql.Open("postgres", postgresConnectionString)
	if err != nil {
		return nil, err
	}

	return &LocalBackend{
		pgDB: db,
	}, nil
}

func (lb *LocalBackend) Close() {
	lb.pgDB.Close()
}

func (lb *LocalBackend) GetArtist(
	url string,
) (*music_api.ArtistWithAlbums, error) {
	return nil, nil
}

func (lb *LocalBackend) GetAlbum(
	url string,
) (*music_api.AlbumWithTracks, error) {
	return nil, nil
}

func (lb *LocalBackend) GetTrack(
	url string,
) (*music_api.TrackWithAlbumAndArtist, error) {
	return nil, nil
}

func (lb *LocalBackend) SearchArtists(
	query string,
) (*music_api.Artists, error) {
	return nil, nil
}

func (lb *LocalBackend) SearchAlbums(
	query string,
) (*music_api.Albums, error) {
	return nil, nil
}
