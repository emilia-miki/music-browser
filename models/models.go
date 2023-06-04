package models

type Artist struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Album struct {
	Name        string   `json:"name"`
	Url         string   `json:"url"`
	ImageUrl    string   `json:"image_url"`
	Tags        []string `json:"tags"`
	ReleaseDate string   `json:"release_date"`
	Artists     []Artist `json:"artists"`
	NumTracks   uint8    `json:"num_tracks"`
}

type Track struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Album Album  `json:"album"`
}
