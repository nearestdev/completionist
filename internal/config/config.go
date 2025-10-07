package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL             string
	JWTSecret               string
	ServerAddr              string
	TMDBApiKey              string
	SteamWebAPIKey          string
	PublicBaseURL           string
	RAWGApiKey              string
	AllowedOrigins          string
	FrontendBaseURL         string
	SteamCallbackRedirectPath string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	return &Config{
		DatabaseURL:             os.Getenv("DATABASE_URL"),
		JWTSecret:               os.Getenv("JWT_SECRET"),
		ServerAddr:              os.Getenv("SERVER_ADDR"),
		TMDBApiKey:              os.Getenv("TMDB_API_KEY"),
		SteamWebAPIKey:          os.Getenv("STEAM_WEB_API_KEY"),
		PublicBaseURL:           os.Getenv("PUBLIC_BASE_URL"),
		RAWGApiKey:              os.Getenv("RAWG_API_KEY"),
		AllowedOrigins:          os.Getenv("ALLOWED_ORIGINS"),
		FrontendBaseURL:         os.Getenv("FRONTEND_BASE_URL"),
		SteamCallbackRedirectPath: os.Getenv("STEAM_CALLBACK_REDIRECT_PATH"),
	}, nil
}