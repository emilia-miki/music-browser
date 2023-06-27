package test_case

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/emilia-miki/music-browser/integration_test/logger"
	"github.com/emilia-miki/music-browser/integration_test/music_api"
)

type BackendName = string

const (
	SpotifyBackendName  BackendName = "spotify"
	YtMusicBackendName              = "yt-music"
	BandcampBackendName             = "bandcamp"
	LocalBackendName                = "local"
)

type ReqType = string

const (
	ReqArtists     ReqType = "artists"
	ReqAlbums              = "albums"
	ReqTrackStream         = "track-stream"
)

type GetType = int

const (
	GetTypeArtist GetType = iota
	GetTypeAlbum
	GetTypeTrack
)

type TestCase struct {
	requestMethod string

	backendName BackendName
	reqType     ReqType
	query       string

	stringifier func(b []byte) string
	unmarshaler func(b []byte) (interface{}, error)
	validator   func(interface{}) ([]string, error)

	url string

	expected []byte

	success   bool
	isChecked bool
}

func NewSearch(
	backendName BackendName,
	searchType ReqType,
	query string,
) TestCase {
	var tc TestCase
	tc.requestMethod = "GET"
	tc.reqType = searchType
	tc.backendName = backendName
	tc.query = query
	tc.stringifier = jsonStringify
	if searchType == ReqArtists {
		tc.unmarshaler = unmarshalArtists
		tc.validator = validateArtists
	} else {
		tc.unmarshaler = unmarshalAlbums
		tc.validator = validateArtists
	}
	return tc
}

func NewGet(getType GetType, url string) TestCase {
	var tc TestCase
	tc.requestMethod = "GET"
	tc.url = url
	tc.stringifier = jsonStringify
	if getType == GetTypeArtist {
		tc.unmarshaler = unmarshalArtistWithAlbums
		tc.validator = validateArtistWithAlbums
	} else if getType == GetTypeAlbum {
		tc.unmarshaler = unmarshalAlbumWithTracks
		tc.validator = validateAlbumWithTracks
	} else {
		tc.unmarshaler = unmarshalTrackWithAlbumAndArtist
		tc.validator = validateTrackWithAlbumAndArtist
	}
	return tc
}

func NewGetStream(url string, expected []byte) TestCase {
	var tc TestCase
	tc.requestMethod = "GET"
	tc.url = url
	tc.reqType = ReqTrackStream
	tc.expected = expected
	tc.stringifier = binaryStringify
	tc.unmarshaler = func(b []byte) (interface{}, error) { return b, nil }
	tc.validator = getTrackStreamValidator(expected)
	return tc
}

func NewDownload(
	url string,
) TestCase {
	var tc TestCase
	tc.requestMethod = "POST"
	tc.url = url
	tc.stringifier = jsonStringify
	tc.unmarshaler = func([]byte) (interface{}, error) { return nil, nil }
	tc.validator = func(interface{}) ([]string, error) { return nil, nil }
	return tc
}

func (tc TestCase) BuildRequest(baseUrl string) (*http.Request, error) {
	params := tc.getParams()
	mappedKVPairs := make([]string, len(params))

	i := 0
	for k, v := range params {
		escaped := url.QueryEscape(v)
		mappedKVPairs[i] = fmt.Sprintf("%s=%s", k, escaped)
		i += 1
	}
	qs := strings.Join(mappedKVPairs, "&")

	url := fmt.Sprintf("%s/?%s", baseUrl, qs)
	return http.NewRequest(tc.requestMethod, url, nil)
}

func (tc TestCase) Check(response []byte) (TestCase, error) {
	tc.isChecked = true

	respString := tc.stringifier(response)
	logger.Info.Printf("Backend response:\n%s", respString)

	obj, err := tc.unmarshaler(response)
	if err != nil {
		return tc, fmt.Errorf(
			"Unable to unmarshal response: %s. Response body:\n%s",
			err, string(response))
	}

	wrns, err := tc.validator(obj)

	if len(wrns) > 0 {
		logger.Warning.Println(
			"Potential problems detected in object:\n" + respString)
	}
	for _, w := range wrns {
		logger.Warning.Println(w)
	}

	if err != nil {
		logger.Error.Println("Errors detected in object:\n" + respString)
		return tc, err
	}

	tc.success = true
	return tc, nil
}

func (tc TestCase) DidPass() bool {
	return tc.success
}

func (tc TestCase) String() string {
	params := tc.getParams()
	mappedKVPairs := make([]string, len(params))

	i := 0
	for k, v := range params {
		mappedKVPairs[i] = fmt.Sprintf(`%s="%s"`, k, v)
		i += 1
	}

	sort.Strings(mappedKVPairs)

	if !tc.isChecked {
		return fmt.Sprintf("%s(%s)",
			tc.requestMethod, strings.Join(mappedKVPairs, ", "))
	}

	var state string
	if tc.success {
		state = "SUCCESS"
	} else {
		state = "FAILURE"
	}

	return fmt.Sprintf("%s: %s(%s)",
		state, tc.requestMethod, strings.Join(mappedKVPairs, ", "))
}

func (tc TestCase) getParams() map[string]string {
	params := make(map[string]string, 8)

	if tc.backendName != "" {
		params["backend"] = tc.backendName
	}

	if tc.url != "" {
		params["url"] = tc.url
	}

	if tc.reqType != "" {
		params["type"] = tc.reqType
	}

	if tc.query != "" {
		params["query"] = tc.query
	}

	return params
}

func unmarshalArtistWithAlbums(b []byte) (interface{}, error) {
	artist := new(music_api.ArtistWithAlbums)
	err := json.Unmarshal(b, artist)
	return artist, err
}

func unmarshalAlbumWithTracks(b []byte) (interface{}, error) {
	album := new(music_api.AlbumWithTracks)
	err := json.Unmarshal(b, album)
	return album, err
}

func unmarshalTrackWithAlbumAndArtist(b []byte) (interface{}, error) {
	track := new(music_api.TrackWithAlbumAndArtist)
	err := json.Unmarshal(b, track)
	return track, err
}

func unmarshalArtists(b []byte) (interface{}, error) {
	artists := new(music_api.Artists)
	err := json.Unmarshal(b, artists)
	return artists, err
}

func unmarshalAlbums(b []byte) (interface{}, error) {
	albums := new(music_api.Albums)
	err := json.Unmarshal(b, albums)
	return albums, err
}

func validateArtist(
	artist *music_api.Artist,
) (warnings []string, err error) {
	var errs []error
	defer func() { err = errors.Join(errs...) }()

	if artist == nil {
		errs = append(errs, errors.New("Artist is nil"))
		return
	}

	if artist.Name == nil {
		errs = append(errs, errors.New("Artist.Name is nil"))
	}

	if artist.Url == nil {
		warnings = append(warnings, "Artist.Url is nil")
	}

	if artist.ImageUrl == nil {
		warnings = append(warnings, "Artist.ImageUrl is nil")
	}

	return
}

func validateAlbum(
	album *music_api.Album,
) (warnings []string, err error) {
	var errs []error
	defer func() { err = errors.Join(errs...) }()

	if album == nil {
		errs = append(errs, errors.New("Album is nil"))
		return
	}

	if album.Name == nil {
		errs = append(errs, errors.New("Album.Name is nil"))
	}

	if album.Url == nil {
		warnings = append(warnings, "Album.Url is nil")
	}

	if album.ImageUrl == nil {
		warnings = append(warnings, "Album.ImageUrl is nil")
	}

	if album.ArtistUrl == nil {
		warnings = append(warnings, "Album.ArtistUrl is nil")
	}

	if album.Year == nil {
		warnings = append(warnings, "Album.Year is nil")
	}

	return
}

func validateTrack(
	track *music_api.Track,
) (warnings []string, err error) {
	var errs []error
	defer func() { err = errors.Join(errs...) }()

	if track == nil {
		errs = append(errs, errors.New("Track is nil"))
		return
	}

	if track.Name == nil {
		errs = append(errs, errors.New("Track.Name is nil"))
	}

	if track.ImageUrl == nil {
		warnings = append(warnings, "Track.ImageUrl is nil")
	}

	if track.ArtistUrl == nil {
		warnings = append(warnings, "Track.ArtistUrl is nil")
	}

	if track.AlbumUrl == nil {
		warnings = append(warnings, "Track.AlbumUrl is nil")
	}

	if track.DurationSeconds == nil {
		warnings = append(warnings, "Track.DurationSeconds is nil")
	}

	return
}

func validateArtists(
	artistsInterface interface{},
) (warnings []string, err error) {
	var errs []error
	defer func() { err = errors.Join(errs...) }()

	artists := artistsInterface.(*music_api.Artists)
	if artists == nil {
		errs = append(errs, errors.New("Artists is nil"))
		return
	}

	for _, artist := range artists.Artists {
		wrns, err := validateArtist(artist)
		warnings = append(warnings, wrns...)
		errs = append(errs, err)
	}

	return
}

func validateAlbums(
	albumsInterface interface{},
) (warnings []string, err error) {
	var errs []error
	defer func() { err = errors.Join(errs...) }()

	albums := albumsInterface.(*music_api.Albums)
	if albums == nil {
		errs = append(errs, errors.New("Albums is nil"))
		return
	}

	for _, album := range albums.Albums {
		wrns, err := validateAlbum(album)
		warnings = append(warnings, wrns...)
		errs = append(errs, err)
	}

	return
}

func validateTracks(
	tracks *music_api.Tracks,
) (warnings []string, err error) {
	var errs []error
	defer func() { err = errors.Join(errs...) }()

	if tracks == nil {
		errs = append(errs, errors.New("Tracks is nil"))
		return
	}

	for _, track := range tracks.Tracks {
		wrns, err := validateTrack(track)
		warnings = append(warnings, wrns...)
		errs = append(errs, err)
	}

	return
}

func validateArtistWithAlbums(
	artistWithAlbumsInterface interface{},
) (warnings []string, err error) {
	var errs []error
	defer func() { err = errors.Join(errs...) }()

	artistWithAlbums := artistWithAlbumsInterface.(*music_api.ArtistWithAlbums)
	if artistWithAlbums == nil {
		errs = append(errs, errors.New("ArtistWithAlbums is nil"))
		return
	}

	var wrns []string

	wrns, err = validateArtist(artistWithAlbums.Artist)
	warnings = append(warnings, wrns...)
	errs = append(errs, err)

	wrns, err = validateAlbums(artistWithAlbums.Albums)
	warnings = append(warnings, wrns...)
	errs = append(errs, err)

	return
}

func validateAlbumWithTracks(
	albumWithTracksInterface interface{},
) (warnings []string, err error) {
	var errs []error
	defer func() { err = errors.Join(errs...) }()

	albumWithTracks := albumWithTracksInterface.(*music_api.AlbumWithTracks)
	if albumWithTracks == nil {
		errs = append(errs, errors.New("AlbumWithTracks is nil"))
		return
	}

	var wrns []string

	wrns, err = validateAlbum(albumWithTracks.Album)
	warnings = append(warnings, wrns...)
	errs = append(errs, err)

	wrns, err = validateTracks(albumWithTracks.Tracks)
	warnings = append(warnings, wrns...)
	errs = append(errs, err)

	return
}

func validateTrackWithAlbumAndArtist(
	trackWithAlbumAndArtistInterface interface{},
) (warnings []string, err error) {
	var errs []error
	defer func() { err = errors.Join(errs...) }()

	trackWithAlbumAndArtist :=
		trackWithAlbumAndArtistInterface.(*music_api.TrackWithAlbumAndArtist)
	if trackWithAlbumAndArtist == nil {
		errs = append(errs, errors.New("TrackWithAlbumAndArtist is nil"))
		return
	}

	var wrns []string

	wrns, err = validateArtist(trackWithAlbumAndArtist.Artist)
	warnings = append(warnings, wrns...)
	errs = append(errs, err)

	wrns, err = validateAlbum(trackWithAlbumAndArtist.Album)
	warnings = append(warnings, wrns...)
	errs = append(errs, err)

	wrns, err = validateTrack(trackWithAlbumAndArtist.Track)
	warnings = append(warnings, wrns...)
	errs = append(errs, err)

	return
}

func getTrackStreamValidator(
	expected []byte,
) func(interface{}) ([]string, error) {
	return func(bInterface interface{}) ([]string, error) {
		const threshold int = 95

		actual := bInterface.([]byte)

		if len(expected) != len(actual) {
			return nil, fmt.Errorf(
				"The byte stream lengths don't match:"+
					"expected length is %d b, but actual length is %d b",
				len(expected), len(actual),
			)
		}

		matchingBytesCount := 0
		for i := 0; i < len(expected); i++ {
			if expected[i] == actual[i] {
				matchingBytesCount += 1
			}
		}

		matchPercentage := 100 * matchingBytesCount / len(expected)
		if matchPercentage < threshold {
			return nil, fmt.Errorf(
				"The byte streams match %d, at least %d required%%",
				matchPercentage, threshold,
			)
		}

		return nil, nil
	}
}

func jsonStringify(b []byte) string {
	var buf bytes.Buffer
	err := json.Indent(&buf, b, "", "\t")
	if err != nil {
		logger.Error.Println("Unable to stringify JSON: " + err.Error())
		return "Unknown"
	}

	return buf.String()
}

func binaryStringify(b []byte) string {
	return fmt.Sprintf("byte array of length %d", len(b))
}
