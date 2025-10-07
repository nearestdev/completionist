package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GATEOPENERZ/completionist-api/internal/config"
	"github.com/GATEOPENERZ/completionist-api/internal/database"
	"github.com/GATEOPENERZ/completionist-api/internal/handler"
	"github.com/GATEOPENERZ/completionist-api/internal/repository"
	"github.com/GATEOPENERZ/completionist-api/internal/routes"
	"github.com/GATEOPENERZ/completionist-api/internal/services/jikan"
	"github.com/GATEOPENERZ/completionist-api/internal/services/rawg"
	"github.com/GATEOPENERZ/completionist-api/internal/services/steam"
	"github.com/GATEOPENERZ/completionist-api/internal/services/tmdb"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db.DB, "migrations"); err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	listRepo := repository.NewListRepository(db)
	socialRepo := repository.NewSocialRepository(db)
	postsRepo := repository.NewPostsRepository(db)
	steamRepo := repository.NewSteamRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)

	jk := jikan.New()
	tm := tmdb.New(cfg.TMDBApiKey)
	st := steam.New(cfg.SteamWebAPIKey, cfg.PublicBaseURL)
	rg := rawg.New(cfg.RAWGApiKey)

	appHandler := handler.NewHandler(
		userRepo, mediaRepo, listRepo, socialRepo, postsRepo, steamRepo,
		jk, tm, st, rg, attachmentRepo, cfg,
	)

	router := routes.NewRouter(appHandler, cfg)

	fmt.Printf("Server starting on port %s\n", cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, router); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}