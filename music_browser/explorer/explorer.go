package explorer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/emilia-miki/music-browser/music_browser/music_api"
)

const TTL_SECONDS = 3600 * 24
const TTL = time.Duration(TTL_SECONDS) * time.Second

type Backend interface {
	GetArtist(url string) (*music_api.ArtistWithAlbums, error)
	GetAlbum(url string) (*music_api.AlbumWithTracks, error)
	GetTrack(url string) (*music_api.TrackWithAlbumAndArtist, error)

	SearchArtists(query string) (*music_api.Artists, error)
	SearchAlbums(query string) (*music_api.Albums, error)
}

type Explorer struct {
	backends map[string]Backend

	redisClient *redis.Client
	pgDB        *sql.DB

	selectTrackUrlStmt    *sql.Stmt
	insertImageStmt       *sql.Stmt
	insertArtistStmt      *sql.Stmt
	insertAlbumStmt       *sql.Stmt
	insertTrackStmt       *sql.Stmt
	insertArtistAlbumStmt *sql.Stmt
	insertArtistTrackStmt *sql.Stmt
	insertAlbumTrackStmt  *sql.Stmt
}

func New(
	backends map[string]Backend,
	redisConnectionString string,
	postgresConnectionString string,
) (*Explorer, error) {
	var e Explorer

	e.backends = backends

	opts, err := redis.ParseURL(redisConnectionString)
	if err != nil {
		return nil, err
	}
	e.redisClient = redis.NewClient(opts)

	e.pgDB, err = sql.Open("postgres", postgresConnectionString)
	if err != nil {
		return nil, err
	}

	e.selectTrackUrlStmt, err = e.pgDB.Prepare(
		"SELECT url FROM track WHERE url = $1")
	if err != nil {
		return nil, err
	}
	e.insertImageStmt, err = e.pgDB.Prepare("INSERT INTO image " +
		"(url, path)" +
		"VALUES ($1, $2)")
	if err != nil {
		return nil, err
	}
	e.insertArtistStmt, err = e.pgDB.Prepare("INSERT INTO artist " +
		"(url, image_url, name)" +
		"VALUES ($1, $2, $3)")
	if err != nil {
		return nil, err
	}
	e.insertAlbumStmt, err = e.pgDB.Prepare("INSERT INTO album " +
		"(url, image_url, name, year)" +
		"VALUES ($1, $2, $3, $4)")
	if err != nil {
		return nil, err
	}
	e.insertTrackStmt, err = e.pgDB.Prepare("INSERT INTO track " +
		"(url, image_url, album_url, path, name, duration_seconds)" +
		"VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return nil, err
	}
	e.insertArtistAlbumStmt, err = e.pgDB.Prepare("INSERT INTO artist_album " +
		"(artist_url, album_url)" +
		"VALUES ($1, $2)")
	if err != nil {
		return nil, err
	}
	e.insertArtistTrackStmt, err = e.pgDB.Prepare("INSERT INTO artist_track " +
		"(artist_url, track_url)" +
		"VALUES ($1, $2)")
	if err != nil {
		return nil, err
	}
	e.insertAlbumTrackStmt, err = e.pgDB.Prepare("INSERT INTO album_track " +
		"(album_url, track_url)" +
		"VALUES ($1, $2)")
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (e *Explorer) Close() error {
	var err error
	err = e.redisClient.Close()
	if err != nil {
		return err
	}
	err = e.pgDB.Close()
	if err != nil {
		return err
	}

	err = e.selectTrackUrlStmt.Close()
	if err != nil {
		return err
	}
	err = e.insertImageStmt.Close()
	if err != nil {
		return err
	}
	err = e.insertArtistStmt.Close()
	if err != nil {
		return err
	}
	err = e.insertAlbumStmt.Close()
	if err != nil {
		return err
	}
	err = e.insertTrackStmt.Close()
	if err != nil {
		return err
	}
	err = e.insertArtistAlbumStmt.Close()
	if err != nil {
		return err
	}
	err = e.insertArtistTrackStmt.Close()
	if err != nil {
		return err
	}
	err = e.insertAlbumTrackStmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func extractBackendNameFromUrl(url string) string {
	splits := strings.Split(url, "/")
	if len(splits) < 3 {
		return "local"
	}

	name := splits[2]
	if name == "open.spotify.com" {
		return "spotify"
	} else if name == "bandcamp.com" {
		return "bandcamp"
	} else if name == "music.youtube.com" {
		return "yt-music"
	} else {
		return "local"
	}
}

func (e *Explorer) GetArtist(
	url string,
) (*music_api.ArtistWithAlbums, error) {
	backendName := extractBackendNameFromUrl(url)
	backend := e.backends[backendName]

	str, err := e.redisClient.GetEx(
		context.Background(),
		url,
		TTL,
	).Result()

	if err == nil {
		var artist *music_api.ArtistWithAlbums
		err = proto.Unmarshal([]byte(str), artist)
		if err != nil {
			return nil, err
		}
		return artist, nil
	} else if !errors.Is(err, redis.Nil) {
		return nil, err
	}

	artist, err := backend.GetArtist(url)
	if err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(artist)
	if err != nil {
		return nil, err
	}

	err = e.redisClient.SetEx(
		context.Background(),
		url,
		string(bytes),
		TTL,
	).Err()
	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (e *Explorer) GetAlbum(
	url string,
) (*music_api.AlbumWithTracks, error) {
	backendName := extractBackendNameFromUrl(url)
	backend := e.backends[backendName]

	str, err := e.redisClient.GetEx(
		context.Background(),
		url,
		TTL,
	).Result()

	if err == nil {
		var album *music_api.AlbumWithTracks
		err = proto.Unmarshal([]byte(str), album)
		if err != nil {
			return nil, err
		}
		return album, nil
	} else if !errors.Is(err, redis.Nil) {
		return nil, err
	}

	album, err := backend.GetAlbum(url)
	if err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(album)
	if err != nil {
		return nil, err
	}

	err = e.redisClient.SetEx(
		context.Background(),
		url,
		string(bytes),
		TTL,
	).Err()
	if err != nil {
		return nil, err
	}

	return album, nil
}

func (e *Explorer) GetTrack(
	url string,
) (*music_api.TrackWithAlbumAndArtist, error) {
	backendName := extractBackendNameFromUrl(url)
	log.Println("backend " + backendName)
	backend := e.backends[backendName]

	str, err := e.redisClient.GetEx(
		context.Background(),
		url,
		TTL,
	).Result()

	if err == nil {
		var track *music_api.TrackWithAlbumAndArtist
		err = proto.Unmarshal([]byte(str), track)
		if err != nil {
			return nil, err
		}

		return track, nil
	}

	track, err := backend.GetTrack(url)
	if err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(track)
	if err != nil {
		return nil, err
	}

	err = e.redisClient.SetEx(
		context.Background(),
		url,
		string(bytes),
		TTL,
	).Err()
	if err != nil {
		return nil, err
	}

	return track, nil
}

func (e *Explorer) SearchArtists(
	backendName string,
	query string,
) (*music_api.Artists, error) {
	backend := e.backends[backendName]
	key := fmt.Sprintf("%s:%s", backendName, query)

	str, err := e.redisClient.GetEx(
		context.Background(),
		key,
		TTL,
	).Result()

	if err == nil {
		var artists *music_api.Artists
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
		TTL,
	).Err()
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (e *Explorer) SearchAlbums(
	backendName string,
	query string,
) (*music_api.Albums, error) {
	backend := e.backends[backendName]
	key := fmt.Sprintf("%s:%s", backendName, query)

	str, err := e.redisClient.GetEx(
		context.Background(),
		key,
		TTL,
	).Result()

	if err == nil {
		var albums *music_api.Albums
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
		TTL,
	).Err()
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (e *Explorer) DownloadTrack(
	track *music_api.TrackWithAlbumAndArtist,
) error {
	cmd := exec.Command("yt-dlp",
		*track.Track.Url, "-x", "--audio-format", "opus", "--write-thumbnail")
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	re, err := regexp.Compile(`^\[\w+\] Destination: (.+)$`)
	if err != nil {
		return err
	}

	lines := strings.Split(string(out), "\n")

	var thumbnail string
	var song string
	for _, line := range lines {
		submatches := re.FindStringSubmatch(line)
		fileName := submatches[1]
		if strings.HasSuffix(fileName, ".webp") {
			thumbnail = fileName
		} else if strings.HasSuffix(fileName, ".opus") {
			song = fileName
		}
	}
	if thumbnail == "" || song == "" {
		return errors.New("thumbnail and song should have been set")
	}

	tx, err := e.pgDB.Begin()
	if err != nil {
		return err
	}

	var url string
	err = tx.Stmt(e.selectTrackUrlStmt).QueryRow(track.Track.Url).Scan(&url)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(e.insertImageStmt).Exec(track.Track.ImageUrl, thumbnail)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(e.insertArtistStmt).Exec(
		track.Artist.Url, track.Artist.ImageUrl, track.Artist.Name)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(e.insertAlbumStmt).Exec(
		track.Album.Url, track.Album.ImageUrl,
		track.Album.Name, track.Album.Year)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(e.insertArtistAlbumStmt).Exec(
		track.Artist.Url, track.Album.Url)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(e.insertTrackStmt).Exec(
		track.Track.Url, track.Track.ImageUrl, track.Track.AlbumUrl,
		song, track.Track.Name, track.Track.DurationSeconds)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(e.insertArtistTrackStmt).Exec(
		track.Artist.Url, track.Track.Url)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(e.insertAlbumTrackStmt).Exec(
		track.Track.AlbumUrl, track.Track.Url)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
