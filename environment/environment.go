package environment

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Ports struct {
	Server      uint16
	YTMusicAPI  uint16
	BandcampAPI uint16
}

type ConnectionStrings struct {
	Redis      string
	PostgreSQL string
}

type Application struct {
	Ports             Ports
	ConnectionStrings ConnectionStrings
}

type SpotifySecrets struct {
	ClientId     string
	ClientSecret string
}

type Secrets struct {
	Spotify SpotifySecrets
}

func getPortFromEnvOrDefault(envVariable string, defaultValue uint16) uint16 {
	portString := os.Getenv(envVariable)
	if portString == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		log.Fatalf("Error parsing port %s: %s\n", portString, err)
	}

	return uint16(parsed)
}

func GetApplication() Application {
	var app Application

	app.Ports.Server = getPortFromEnvOrDefault("MUSIC_BROWSER_PORT", 3333)
	app.Ports.YTMusicAPI = getPortFromEnvOrDefault("YOUTUBE_MUSIC_API_PORT", 3334)
	app.Ports.BandcampAPI = getPortFromEnvOrDefault("MUSIC_BROWSER_PORT", 3335)
	app.ConnectionStrings.Redis = fmt.Sprintf("redis://localhost:%d/0", getPortFromEnvOrDefault("REDIS_PORT", 3336))
	app.ConnectionStrings.PostgreSQL = fmt.Sprintf("postgresql://postgres@localhost:%d/postgres?sslmode=disable", getPortFromEnvOrDefault("POSTGRES_PORT", 3337))

	return app
}

func GetSecrets() Secrets {
	var secrets Secrets

	secrets.Spotify.ClientId = os.Getenv("SPOTIFY_CLIENT_ID")
	secrets.Spotify.ClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")

	if secrets.Spotify.ClientId == "" || secrets.Spotify.ClientSecret == "" {
		log.Println("Spotify API secrets not provided")
	}

	return secrets
}
