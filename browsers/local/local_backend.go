package local_backend

import (
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

	"github.com/emilia-miki/music-browser/music_browser/dal/link_map_repository"
	"github.com/emilia-miki/music-browser/music_browser/dal/music_repository"
	"github.com/emilia-miki/music-browser/music_browser/music_api"
	_ "github.com/lib/pq"
)

const DefaultMusicDirectory = "music"

type LocalBackend struct {
	db          *sql.DB
	linkMapRepo *link_map_repository.LinkMapRepository
	musicRepo   *music_repository.MusicRepository
	mapLink     func(url string) (string, error)
}

func New(
	db *sql.DB,
	mapLink func(string) (string, error),
) (*LocalBackend, error) {
	var err error
	var lb LocalBackend

	lb.db = db
	lb.linkMapRepo = link_map_repository.New(lb.db)
	lb.musicRepo = music_repository.New(lb.db)
	lb.mapLink = mapLink

	err = os.MkdirAll(DefaultMusicDirectory, os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &lb, nil
}

func (lb *LocalBackend) String() string {
	return "local"
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

func (lb *LocalBackend) ensureImageDownloaded(
	tx *sql.Tx,
	name string,
	imageUrl *string,
) (res sql.NullString, err error) {
	if imageUrl == nil {
		return
	}

	exists, err := lb.musicRepo.ImageExists(*imageUrl)
	if err != nil || !exists {
		return
	}

	path, err := downloadImage(name, *imageUrl)
	if err != nil {
		return
	}

	err = lb.musicRepo.InsertImage(*imageUrl, path)
	if err != nil {
		return
	}

	res = sql.NullString{
		String: *imageUrl,
		Valid:  true,
	}
	return
}

func (lb *LocalBackend) DownloadTrack(
	track *music_api.TrackWithAlbumAndArtist,
) error {
	if track.Track.Url != nil {
		url, err := lb.mapLink(*track.Track.Url)
		if err != nil {
			return err
		}
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

	tx, err := lb.db.Begin()
	if err != nil {
		return err
	}

	artistExists := track.Artist.Url != nil
	if artistExists {
		url := *track.Artist.Url
		name := *track.Artist.Name

		imageUrl, err := lb.ensureImageDownloaded(
			tx, *track.Artist.Name, track.Artist.ImageUrl)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(lb.insertArtistStmt).Exec(url, imageUrl, name)
		if err != nil {
			return err
		}
	}

	albumExists := track.Album.Url != nil
	if albumExists {
		url := *track.Album.Url
		name := *track.Album.Name
		var year sql.NullInt32

		imageUrl, err := lb.ensureImageDownloaded(
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

		_, err = tx.Stmt(lb.insertAlbumStmt).Exec(url, imageUrl, name, year)
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

		imageUrl, err := lb.ensureImageDownloaded(
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

		_, err = tx.Stmt(lb.insertTrackStmt).Exec(
			url, imageUrl, albumUrl, path, name, durationSeconds)
		if err != nil {
			return err
		}
	}

	if artistExists && albumExists {
		_, err = tx.Stmt(lb.insertArtistAlbumStmt).Exec(
			*track.Artist.Url, *track.Album.Url)
		if err != nil {
			return err
		}
	}

	if artistExists && trackExists {
		_, err = tx.Stmt(lb.insertArtistTrackStmt).Exec(
			*track.Artist.Url, *track.Track.Url)
		if err != nil {
			return err
		}
	}

	if albumExists && trackExists {
		_, err = tx.Stmt(lb.insertAlbumTrackStmt).Exec(
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

func (lb *LocalBackend) GetTrackStream(
	url string,
) (b []byte, mime string, err error) {
	var path string
	row := lb.selectTrackPathStmt.QueryRow(url)
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
