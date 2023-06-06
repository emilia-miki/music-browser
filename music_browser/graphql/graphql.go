package graphql

import (
	"github.com/emilia-miki/music-browser/music_browser/backend"
	"github.com/graphql-go/graphql"
)

var Backends map[string]backend.MusicExplorer

var backendEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "Backend",
	Values: graphql.EnumValueConfigMap{
		"SPOTIFY":  &graphql.EnumValueConfig{Value: "spotify"},
		"BANDCAMP": &graphql.EnumValueConfig{Value: "bandcamp"},
		"YT_MUSIC": &graphql.EnumValueConfig{Value: "yt-music"},
		"LOCAL":    &graphql.EnumValueConfig{Value: "local"},
	},
})

var artistLinkType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ArtistLink",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var trackType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Track",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"imageUrl": &graphql.Field{
			Type: graphql.String,
		},
		"durationSeconds": &graphql.Field{
			Type: graphql.Int,
		},
		"album": &graphql.Field{
			Type: graphql.String,
		},
		"albumUrl": &graphql.Field{
			Type: graphql.String,
		},
		"artists": &graphql.Field{
			Type: graphql.NewList(artistLinkType),
		},
	},
})

var albumType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Album",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"imageUrl": &graphql.Field{
			Type: graphql.String,
		},
		"year": &graphql.Field{
			Type: graphql.Int,
		},
		"durationSeconds": &graphql.Field{
			Type: graphql.Int,
		},
		"tracks": &graphql.Field{
			Type: graphql.NewList(trackType),
		},
		"artists": &graphql.Field{
			Type: graphql.NewList(artistLinkType),
		},
	},
})

var artistType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Artist",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"imageUrl": &graphql.Field{
			Type: graphql.String,
		},
		"albums": &graphql.Field{
			Type: graphql.NewList(albumType),
		},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"artists": &graphql.Field{
			Type: graphql.NewList(artistType),
			Args: graphql.FieldConfigArgument{
				"backend": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(backendEnum),
				},
				"query": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				backend := Backends[p.Args["backend"].(string)]
				query := p.Args["query"].(string)
				artists := backend.SearchArtists(query)
				return artists, nil
			},
		},
		"albums": &graphql.Field{
			Type: graphql.NewList(albumType),
			Args: graphql.FieldConfigArgument{
				"backend": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(backendEnum),
				},
				"query": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				backend := Backends[p.Args["backend"].(string)]
				query := p.Args["query"].(string)
				albums := backend.SearchAlbums(query)
				return albums, nil
			},
		},
	},
})

func GetSchema(backends map[string]backend.MusicExplorer) graphql.Schema {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	Backends = backends
	return schema
}
