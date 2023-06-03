package models

type Track struct {
	Name        string
	Url         string
	ImageUrl    string
	Tags        []string
	ReleaseDate string
	Album       string
	Artist      string
}

type Album struct {
	Name        string
	Url         string
	ImageUrl    string
	Tags        []string
	ReleaseDate string
	Artist      string
	NumTracks   uint8
	NumMinutes  uint8
}

type Artist struct {
	Name     string
	Url      string
	ImageUrl string
	Tags     []string
	Genre    string
	Location string
}
