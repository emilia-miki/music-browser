package graphql

import (
	"github.com/emilia-miki/music-browser/music_browser/explorer"
	"github.com/graphql-go/graphql"
)

var artistType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Artist",
	Fields: graphql.Fields{
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"image_url": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var artistsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Artists",
	Fields: graphql.Fields{
		"artists": &graphql.Field{
			Type: graphql.NewList(artistType),
		},
	},
})

var albumType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Album",
	Fields: graphql.Fields{
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"image_url": &graphql.Field{
			Type: graphql.String,
		},
		"artist_urls": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"year": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var albumsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Albums",
	Fields: graphql.Fields{
		"albums": &graphql.Field{
			Type: graphql.NewList(albumType),
		},
	},
})

var trackType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Track",
	Fields: graphql.Fields{
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"image_url": &graphql.Field{
			Type: graphql.String,
		},
		"artist_urls": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"album_url": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"duration_seconds": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var tracksType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Tracks",
	Fields: graphql.Fields{
		"tracks": &graphql.Field{
			Type: graphql.NewList(trackType),
		},
	},
})

var artistWithAlbumsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ArtistWithAlbums",
	Fields: graphql.Fields{
		"artist": &graphql.Field{
			Type: artistType,
		},
		"albums": &graphql.Field{
			Type: albumsType,
		},
	},
})

var albumWithTracksType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AlbumWithTracks",
	Fields: graphql.Fields{
		"album": &graphql.Field{
			Type: albumType,
		},
		"tracks": &graphql.Field{
			Type: tracksType,
		},
	},
})

func newQueryObject(explorer *explorer.Explorer) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"getArtist": &graphql.Field{
				Type: artistWithAlbumsType,
				Args: graphql.FieldConfigArgument{
					"url": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					url := p.Args["url"].(string)
					artist, err := explorer.GetArtist(url)
					if err != nil {
						return nil, err
					}

					return artist, nil
				},
			},
			"getAlbum": &graphql.Field{
				Type: albumWithTracksType,
				Args: graphql.FieldConfigArgument{
					"url": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					url := p.Args["url"].(string)
					album, err := explorer.GetAlbum(url)
					if err != nil {
						return nil, err
					}

					return album, nil
				},
			},
			"searchArtists": &graphql.Field{
				Type: graphql.NewList(artistType),
				Args: graphql.FieldConfigArgument{
					"backend": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"query": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					backendName := p.Args["backend"].(string)
					query := p.Args["query"].(string)
					artists, err := explorer.SearchArtists(backendName, query)
					if err != nil {
						return nil, err
					}

					return artists, nil
				},
			},
			"searchAlbums": &graphql.Field{
				Type: graphql.NewList(albumType),
				Args: graphql.FieldConfigArgument{
					"backend": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"query": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					backendName := p.Args["backend"].(string)
					query := p.Args["query"].(string)
					albums, err := explorer.SearchAlbums(backendName, query)
					if err != nil {
						return nil, err
					}

					return albums, nil
				},
			},
		},
	})
}

func NewSchema(explorer *explorer.Explorer) graphql.Schema {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: newQueryObject(explorer),
	})
	return schema
}
