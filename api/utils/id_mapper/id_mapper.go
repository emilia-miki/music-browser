package link_manager

import (
	"database/sql"
	"errors"

	"github.com/emilia-miki/music-browser/music_browser/dal/link_map_repository"
	"github.com/emilia-miki/music-browser/music_browser/explorer"
	"github.com/emilia-miki/music-browser/music_browser/music_api"
)

type LinkManager struct {
	spotify     explorer.Backend
	ytm         explorer.Backend
	linkMapRepo *link_map_repository.LinkMapRepository
}

func New(
	spotify explorer.Backend,
	ytm explorer.Backend,
	linkMapRepo *link_map_repository.LinkMapRepository,
) LinkManager {
	return LinkManager{
		spotify:     spotify,
		ytm:         ytm,
		linkMapRepo: linkMapRepo,
	}
}

func (lm *LinkManager) translateLink(id string) (string, error) {
	translated, err := lm.linkMapRepo.GetLinkMap(id)
	if err == nil {
		return translated, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	track, err := lm.spotify.GetTrack(id)
	if err != nil {
		return "", err
	}

	albums, err := lm.ytm.SearchAlbums(*track.Album.Name)
	if err != nil {
		return "", err
	}

	var album *music_api.AlbumWithTracks
	for _, a := range albums.Albums {
		if *a.Name == *track.Album.Name {
			album, err = lm.ytm.GetAlbum(*a.Id)
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
			translated = *t.Id
			err := lm.linkMapRepo.InsertLinkMap(id, translated)
			if err != nil {
				return "", err
			}

			return translated, nil
		}
	}

	return "", errors.New("Couldn't find this track on Youtube Music")
}
