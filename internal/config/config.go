package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)
type Config struct {
	Environment               string
	DatabaseURL               string
	JWTSecret                 string
	ServerAddr                string
	TMDBApiKey                string
	SteamWebAPIKey            string
	PublicBaseURL             string
	RAWGApiKey                string
	AllowedOrigins            string
	FrontendBaseURL           string
	SteamCallbackRedirectPath string
	GoogleBooksAPIKey         string
	LastFMAPIKey              string
	LastFMAPISecret           string
}
func Load() (*Config, error) {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}
	if err := godotenv.Load(); err != nil {
		if env == "development" {
			log.Println("Warning: .env file not found")
		}
	}
	return &Config{
		Environment:               env,
		DatabaseURL:               os.Getenv("DATABASE_URL"),
		JWTSecret:                 os.Getenv("JWT_SECRET"),
		ServerAddr:                os.Getenv("SERVER_ADDR"),
		TMDBApiKey:                os.Getenv("TMDB_API_KEY"),
		SteamWebAPIKey:            os.Getenv("STEAM_WEB_API_KEY"),
		PublicBaseURL:             os.Getenv("PUBLIC_BASE_URL"),
		RAWGApiKey:                os.Getenv("RAWG_API_KEY"),
		AllowedOrigins:            os.Getenv("ALLOWED_ORIGINS"),
		FrontendBaseURL:           os.Getenv("FRONTEND_BASE_URL"),
		SteamCallbackRedirectPath: os.Getenv("STEAM_CALLBACK_REDIRECT_PATH"),
		GoogleBooksAPIKey:         os.Getenv("GOOGLE_BOOKS_API_KEY"),
		LastFMAPIKey:              os.Getenv("LASTFM_API_KEY"),
		LastFMAPISecret:           os.Getenv("LASTFM_API_SECRET"),
	}, nil
}