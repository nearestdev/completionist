package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment               string
	DatabaseURL               string
	JWTSecret                 string
	BackendPort               string
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
	StorageDriver             string
	AWSRegion                 string
	AWSBucket                 string
	AWSAccessKeyID            string
	AWSSecretAccessKey        string
	DevAdminEmail             string
	DevAdminPassword          string
	DevAdminUsername          string
	StripeSecretKey           string
	StripeWebhookSecret       string
	StripePriceID             string
	ModerationEnabled         bool
	OpenAIModerationKey       string
	LLMProvider               string
	LLMAPIKey                 string
}

func Load() (*Config, error) {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	envFiles := []string{".env", "../.env", "../../.env"}
	envLoaded := false

	for _, envFile := range envFiles {
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Load(envFile); err != nil && env == "development" {
				log.Printf("Warning: could not load %s: %v", envFile, err)
			}
			envLoaded = true
			break
		}
	}

	if !envLoaded && env == "development" {
		log.Println("Warning: .env file not found")
	}

	backendBaseURL := os.Getenv("BACKEND_BASE_URL")
	if backendBaseURL == "" {
		backendBaseURL = "http://localhost"
	}

	frontendBaseURL := os.Getenv("FRONTEND_BASE_URL")
	if frontendBaseURL == "" {
		frontendBaseURL = "http://localhost"
	}

	backendPort := os.Getenv("BACKEND_PORT")
	if backendPort == "" {
		backendPort = "8080"
	}

	frontendPort := os.Getenv("PORT")
	if frontendPort == "" {
		frontendPort = "3000"
	}

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":" + backendPort
	}

	publicBaseURL := os.Getenv("PUBLIC_BASE_URL")
	if publicBaseURL == "" {
		publicBaseURL = backendBaseURL + ":" + backendPort
	}

	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins string
	defaultFrontendOrigin := frontendBaseURL + ":" + frontendPort
	secondaryFrontendOrigin := strings.Replace(defaultFrontendOrigin, "localhost", "127.0.0.1", 1)

	if allowedOriginsEnv == "" {
		allowedOrigins = defaultFrontendOrigin + "," + secondaryFrontendOrigin
	} else {
		allowedOrigins = defaultFrontendOrigin + "," + secondaryFrontendOrigin + "," + allowedOriginsEnv
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = frontendBaseURL + ":" + frontendPort
	}

	return &Config{
		Environment:               env,
		DatabaseURL:               os.Getenv("DATABASE_URL"),
		JWTSecret:                 os.Getenv("JWT_SECRET"),
		BackendPort:               backendPort,
		ServerAddr:                serverAddr,
		TMDBApiKey:                os.Getenv("TMDB_API_KEY"),
		SteamWebAPIKey:            os.Getenv("STEAM_WEB_API_KEY"),
		PublicBaseURL:             publicBaseURL,
		RAWGApiKey:                os.Getenv("RAWG_API_KEY"),
		AllowedOrigins:            allowedOrigins,
		FrontendBaseURL:           frontendURL,
		SteamCallbackRedirectPath: os.Getenv("STEAM_CALLBACK_REDIRECT_PATH"),
		GoogleBooksAPIKey:         os.Getenv("GOOGLE_BOOKS_API_KEY"),
		LastFMAPIKey:              os.Getenv("LASTFM_API_KEY"),
		LastFMAPISecret:           os.Getenv("LASTFM_API_SECRET"),
		StorageDriver:             os.Getenv("STORAGE_DRIVER"), // local or s3
		AWSRegion:                 os.Getenv("AWS_REGION"),
		AWSBucket:                 os.Getenv("AWS_BUCKET"),
		AWSAccessKeyID:            os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:        os.Getenv("AWS_SECRET_ACCESS_KEY"),
		DevAdminEmail:             os.Getenv("DEV_ADMIN_EMAIL"),
		DevAdminPassword:          os.Getenv("DEV_ADMIN_PASSWORD"),
		DevAdminUsername:          os.Getenv("DEV_ADMIN_USERNAME"),
		StripeSecretKey:           os.Getenv("STRIPE_SECRET_KEY"),
		StripeWebhookSecret:       os.Getenv("STRIPE_WEBHOOK_SECRET"),
		StripePriceID:             os.Getenv("STRIPE_PRICE_ID"),
		ModerationEnabled:         os.Getenv("MODERATION_ENABLED") == "true",
		OpenAIModerationKey:       os.Getenv("OPENAI_MODERATION_KEY"),
		LLMProvider:               os.Getenv("LLM_PROVIDER"),
		LLMAPIKey:                 os.Getenv("LLM_API_KEY"),
	}, nil
}
