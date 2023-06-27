package explorer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	// "github.com/emilia-miki/music-browser/music_browser/logger"
	"github.com/emilia-miki/music-browser/music_browser/music_api"
)

const TTL = 3600 * 24 * time.Second
const DefaultMusicDirectory = "music"

type BackendName = string

const (
	LocalBackendName    BackendName = "local"
	SpotifyBackendName              = "spotify"
	BandcampBackendName             = "bandcamp"
	YtMusicBackendName              = "yt-music"
)

type Backend interface {
	String() string

	GetArtist(url string) (*music_api.ArtistWithAlbums, error)
	GetAlbum(url string) (*music_api.AlbumWithTracks, error)
	GetTrack(url string) (*music_api.TrackWithAlbumAndArtist, error)

	SearchArtists(query string) (*music_api.Artists, error)
	SearchAlbums(query string) (*music_api.Albums, error)
}

type Explorer struct {
	backends map[BackendName]Backend

	redisClient *redis.Client
	pgDB        *sql.DB
	musicDir    *os.File

	selectLinkMap         *sql.Stmt
	insertLinkMap         *sql.Stmt
	selectImageUrlStmt    *sql.Stmt
	insertImageStmt       *sql.Stmt
	insertArtistStmt      *sql.Stmt
	insertAlbumStmt       *sql.Stmt
	selectTrackPathStmt   *sql.Stmt
	insertTrackStmt       *sql.Stmt
	insertArtistAlbumStmt *sql.Stmt
	insertArtistTrackStmt *sql.Stmt
	insertAlbumTrackStmt  *sql.Stmt
}

func New(
	backends map[BackendName]Backend,
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
	err = os.MkdirAll(DefaultMusicDirectory, os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, err
	}

	e.pgDB, err = sql.Open("postgres", postgresConnectionString)
	if err != nil {
		return nil, err
	}

	e.selectLinkMap, err = e.pgDB.Prepare(
		"SELECT yt_url FROM link_map WHERE sp_url = $1")
	if err != nil {
		return nil, err
	}

	e.insertLinkMap, err = e.pgDB.Prepare(
		"INSERT INTO link_map (sp_url, yt_url) VALUES ($1, $2)")
	if err != nil {
		return nil, err
	}

	e.selectImageUrlStmt, err = e.pgDB.Prepare(
		"SELECT url FROM image WHERE url = $1")
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

	e.selectTrackPathStmt, err = e.pgDB.Prepare("SELECT path FROM track " +
		"WHERE url = $1")
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
	var errs []error
	var err error

	err = e.redisClient.Close()
	errs = append(errs, err)

	err = e.pgDB.Close()
	errs = append(errs, err)

	err = e.selectLinkMap.Close()
	errs = append(errs, err)

	err = e.selectImageUrlStmt.Close()
	errs = append(errs, err)

	err = e.insertImageStmt.Close()
	errs = append(errs, err)

	err = e.insertArtistStmt.Close()
	errs = append(errs, err)

	err = e.insertAlbumStmt.Close()
	errs = append(errs, err)

	err = e.insertTrackStmt.Close()
	errs = append(errs, err)

	err = e.insertArtistAlbumStmt.Close()
	errs = append(errs, err)

	err = e.insertArtistTrackStmt.Close()
	errs = append(errs, err)

	err = e.insertAlbumTrackStmt.Close()
	errs = append(errs, err)

	return errors.Join(errs...)
}

func extractBackendName(url string) string {
	var trimmedUrl string
	var found bool
	trimmedUrl, found = strings.CutPrefix(url, "http://")
	if !found {
		trimmedUrl, found = strings.CutPrefix(url, "https://")
	}
	if !found {
		trimmedUrl = url
	}

	splits := strings.Split(trimmedUrl, "/")
	name := splits[0]
	if name == "open.spotify.com" {
		return SpotifyBackendName
	} else if strings.HasSuffix(name, "bandcamp.com") {
		return BandcampBackendName
	} else if name == "music.youtube.com" {
		return YtMusicBackendName
	} else {
		return LocalBackendName
	}
}

func (e *Explorer) GetArtist(
	url string,
) (*music_api.ArtistWithAlbums, error) {
	if url == "" {
		return nil, errors.New("Empty url")
	}

	backendName := extractBackendName(url)
	backend, ok := e.backends[backendName]
	if !ok {
		return nil, errors.New("Invalid backend: " + backendName)
	}

	key := "artist:" + url
	str, err := e.redisClient.GetEx(
		context.Background(),
		key,
		TTL,
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
		key,
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
	if url == "" {
		return nil, errors.New("Empty url")
	}

	backendName := extractBackendName(url)
	backend, ok := e.backends[backendName]
	if !ok {
		return nil, errors.New("Invalid backend: " + backendName)
	}

	key := "album" + url
	str, err := e.redisClient.GetEx(
		context.Background(),
		key,
		TTL,
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
		key,
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
	if url == "" {
		return nil, errors.New("Empty url")
	}

	backendName := extractBackendName(url)
	backend, ok := e.backends[backendName]
	if !ok {
		return nil, errors.New("Invalid backend: " + backendName)
	}

	key := "track" + url

	str, err := e.redisClient.GetEx(
		context.Background(),
		url,
		TTL,
	).Result()

	if err == nil {
		track := new(music_api.TrackWithAlbumAndArtist)
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
		key,
		string(bytes),
		TTL,
	).Err()
	if err != nil {
		return nil, err
	}

	return track, nil
}

func (e *Explorer) SearchArtists(
	backendName BackendName,
	query string,
) (*music_api.Artists, error) {
	backend, ok := e.backends[backendName]
	if !ok {
		return nil, errors.New("Invalid backend: " + backendName)
	}

	if query == "" {
		return nil, errors.New("Empty query")
	}

	key := fmt.Sprintf("%s:artists:%s", backend, query)

	str, err := e.redisClient.GetEx(
		context.Background(),
		key,
		TTL,
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
		TTL,
	).Err()
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (e *Explorer) SearchAlbums(
	backendName BackendName,
	query string,
) (*music_api.Albums, error) {
	backend, ok := e.backends[backendName]
	if !ok {
		return nil, errors.New("Invalid backend: " + backendName)
	}

	if query == "" {
		return nil, errors.New("Empty query")
	}

	key := fmt.Sprintf("%s:albums:%s", backend, query)

	str, err := e.redisClient.GetEx(
		context.Background(),
		key,
		TTL,
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
		TTL,
	).Err()
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (e *Explorer) translateLink(url string) (string, error) {
	fmt.Printf("Translating sp_url %s to yt_url\n", url)

	var translatedUrl string
	err := e.selectLinkMap.QueryRow(url).Scan(&translatedUrl)
	fmt.Printf("Postgres returned err=%s, url=%s\n", err, translatedUrl)
	if err == nil {
		fmt.Printf("Link translation is cached: %s -> %s\n", url, translatedUrl)
		return translatedUrl, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	sp := e.backends[SpotifyBackendName]
	ytm := e.backends[YtMusicBackendName]

	track, err := sp.GetTrack(url)
	if err != nil {
		return "", err
	}

	albums, err := ytm.SearchAlbums(*track.Album.Name)
	if err != nil {
		return "", err
	}

	var album *music_api.AlbumWithTracks
	for _, a := range albums.Albums {
		if *a.Name == *track.Album.Name {
			album, err = ytm.GetAlbum(*a.Url)
			if err != nil {
				return "", err
			}
			break
		}
	}

	if album == nil {
		return "", errors.New("Couldn't find this track on Youtube Music")
	}

	for _, t := range album.Tracks.Tracks {
		if *track.Track.Name == *t.Name {
			fmt.Printf(
				"Translated sp_link %s to yt_link %s\n",
				url, *t.Url,
			)

			_, err := e.insertLinkMap.Exec(url, *t.Url)
			if err != nil {
				return "", err
			}

			return *t.Url, nil
		}
	}

	return "", errors.New("Couldn't find this track on Youtube Music")
}

func (e *Explorer) GetTrackStream(
	url string,
) (b []byte, mime string, err error) {
	if extractBackendName(url) == SpotifyBackendName {
		url, err = e.translateLink(url)
		if err != nil {
			return
		}
	}

	var path string
	row := e.selectTrackPathStmt.QueryRow(url)
	err = row.Scan(&path)
	if err != nil {
		return
	}

	b, err = os.ReadFile(fmt.Sprintf(
		"%s/%s", DefaultMusicDirectory, path))
	if err != nil {
		return
	}

	mime = "audio/ogg"
	return
}

func downloadImage(name string, url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("ERROR: " + err.Error())
	}

	ct := resp.Header.Get("Content-Type")
	ext, ok := strings.CutPrefix(ct, "image/")
	if !ok {
		return "", errors.New("Invalid Content-Type: " + ct)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	id := rand.Int() % 10000000000
	basePath := fmt.Sprintf("%s/%s [%d]", DefaultMusicDirectory, name, id)
	origPath := fmt.Sprintf("%s.%s", basePath, ext)
	newPath := fmt.Sprintf("%s.%s", basePath, "webp")

	err = ioutil.WriteFile(origPath, b, 0644)
	if err != nil {
		return "", err
	}

	if origPath != newPath {
		_, err := exec.Command("cwebp", origPath, "-o", newPath).Output()
		os.Remove(origPath)
		if err != nil {
			return "", err
		}
	}

	os.Remove(origPath)
	return newPath, nil
}

func (e *Explorer) ensureImageDownloaded(
	tx *sql.Tx,
	name string,
	imageUrl *string,
) (res sql.NullString, err error) {
	if imageUrl == nil {
		return
	}

	row := tx.Stmt(e.selectImageUrlStmt).QueryRow(*imageUrl)
	url := new(string)
	err = row.Scan(url)
	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		return
	}

	path, err := downloadImage(name, *imageUrl)
	if err != nil {
		return
	}

	_, err = tx.Stmt(e.insertImageStmt).Exec(*imageUrl, path)
	if err != nil {
		return
	}

	res = sql.NullString{
		String: *imageUrl,
		Valid:  true,
	}
	return
}

func (e *Explorer) DownloadTrack(url string) error {
	fmt.Println("downloading track " + url)
	backendName := extractBackendName(url)
	if backendName == SpotifyBackendName {
		var err error
		url, err = e.translateLink(url)
		if err != nil {
			return err
		}
		backendName = YtMusicBackendName
	}

	backend := e.backends[backendName]
	track, err := backend.GetTrack(url)
	if err != nil {
		return nil
	}

	if track.Track.Url == nil {
		track.Track.Url = &url
	}

	cmd := exec.Command("yt-dlp", *track.Track.Url,
		"-x", "--audio-format", "opus")
	cmd.Dir = DefaultMusicDirectory

	out, err := cmd.Output()
	if err != nil {
		return err
	}
	lines := string(out)

	if strings.Contains(lines, "has already been downloaded") {
		return nil
	}

	var trackPath string
	songRe, err := regexp.Compile(`\[download\] Destination: (.+)\.\w+`)
	if err != nil {
		return err
	}
	submatches := songRe.FindStringSubmatch(lines)
	if len(submatches) == 2 {
		trackPath = submatches[1] + ".opus"
	}

	if trackPath == "" {
		return errors.New("The song filename should have been set")
	}

	tx, err := e.pgDB.Begin()
	if err != nil {
		return err
	}

	artistExists := track.Artist.Url != nil
	if artistExists {
		url := *track.Artist.Url
		name := *track.Artist.Name

		imageUrl, err := e.ensureImageDownloaded(
			tx, *track.Artist.Name, track.Artist.ImageUrl)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(e.insertArtistStmt).Exec(url, imageUrl, name)
		if err != nil {
			return err
		}
	}

	albumExists := track.Album.Url != nil
	if albumExists {
		url := *track.Album.Url
		name := *track.Album.Name
		var year sql.NullInt32

		imageUrl, err := e.ensureImageDownloaded(
			tx, name, track.Album.ImageUrl)
		if err != nil {
			return err
		}

		if track.Album.Year != nil {
			year = sql.NullInt32{
				Int32: int32(*track.Album.Year),
				Valid: true,
			}
		}

		_, err = tx.Stmt(e.insertAlbumStmt).Exec(url, imageUrl, name, year)
		if err != nil {
			return err
		}
	}

	trackExists := track.Track.Url != nil && trackPath != ""
	if trackExists {
		url := *track.Track.Url
		var albumUrl sql.NullString
		path := trackPath
		name := *track.Track.Name
		var durationSeconds sql.NullInt32

		imageUrl, err := e.ensureImageDownloaded(
			tx, name, track.Track.ImageUrl)
		if err != nil {
			return nil
		}

		if albumExists {
			albumUrl = sql.NullString{
				String: *track.Album.Url,
				Valid:  true,
			}
		}

		if track.Track.DurationSeconds != nil {
			durationSeconds = sql.NullInt32{
				Int32: int32(*track.Track.DurationSeconds),
				Valid: true,
			}
		}

		_, err = tx.Stmt(e.insertTrackStmt).Exec(
			url, imageUrl, albumUrl, path, name, durationSeconds)
		if err != nil {
			return err
		}
	}

	if artistExists && albumExists {
		_, err = tx.Stmt(e.insertArtistAlbumStmt).Exec(
			*track.Artist.Url, *track.Album.Url)
		if err != nil {
			return err
		}
	}

	if artistExists && trackExists {
		_, err = tx.Stmt(e.insertArtistTrackStmt).Exec(
			*track.Artist.Url, *track.Track.Url)
		if err != nil {
			return err
		}
	}

	if albumExists && trackExists {
		_, err = tx.Stmt(e.insertAlbumTrackStmt).Exec(
			*track.Track.AlbumUrl, *track.Track.Url)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
