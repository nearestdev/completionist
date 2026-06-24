package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/GATEOPENERZ/completionist-api/internal/bootstrap"
	"github.com/GATEOPENERZ/completionist-api/internal/config"
	"github.com/GATEOPENERZ/completionist-api/internal/database"
	"github.com/GATEOPENERZ/completionist-api/internal/handler"
	"github.com/GATEOPENERZ/completionist-api/internal/repository"
	"github.com/GATEOPENERZ/completionist-api/internal/routes"

	"github.com/GATEOPENERZ/completionist-api/internal/services"
	"github.com/GATEOPENERZ/completionist-api/internal/services/challenges"
	"github.com/GATEOPENERZ/completionist-api/internal/services/filestorage"
	"github.com/GATEOPENERZ/completionist-api/internal/services/googlebooks"
	"github.com/GATEOPENERZ/completionist-api/internal/services/importer"
	"github.com/GATEOPENERZ/completionist-api/internal/services/jikan"
	"github.com/GATEOPENERZ/completionist-api/internal/services/lastfm"
	"github.com/GATEOPENERZ/completionist-api/internal/services/rawg"
	"github.com/GATEOPENERZ/completionist-api/internal/services/steam"
	"github.com/GATEOPENERZ/completionist-api/internal/services/tmdb"
	"github.com/GATEOPENERZ/completionist-api/internal/websocket"
	_ "github.com/GATEOPENERZ/completionist-api/swagger"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// @title           Completionist API
// @version         2.0.0
// @description     This is the API server for the Completionist application.
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	migrationsPath := config.ResolveAppPath("migrations")
	if err := database.RunMigrations(db.DB, migrationsPath); err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	listRepo := repository.NewListRepository(db)
	socialRepo := repository.NewSocialRepository(db)
	postsRepo := repository.NewPostsRepository(db)
	steamRepo := repository.NewSteamRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)
	lastfmRepo := repository.NewLastFMRepository(db)
	challengeRepo := repository.NewChallengeRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	messagingRepo := repository.NewMessagingRepository(db)
	streakRepo := repository.NewStreakRepository(db)
	badgeRepo := repository.NewBadgeRepository(db)
	collectionRepo := repository.NewCollectionRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	moderationRepo := repository.NewModerationRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	adRepo := repository.NewAdRepository(db)
	connAcctRepo := repository.NewConnectedAccountRepository(db)
	franchiseRepo := repository.NewFranchiseRepository(db)
	rankService := services.NewRankService(db)

	if err := bootstrap.EnsureDevAdmin(cfg, userRepo); err != nil {
		log.Fatalf("could not ensure development admin user: %v", err)
	}

	wsHub := websocket.NewHub(messagingRepo)
	go wsHub.Run()

	var fileService filestorage.Service

	if cfg.StorageDriver == "s3" {
		awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
			awsConfig.WithRegion(cfg.AWSRegion),
			awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, "")),
		)
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}
		s3Client := s3.NewFromConfig(awsCfg)
		fileService = filestorage.NewS3Storage(s3Client, cfg.AWSBucket, cfg.AWSRegion)
		fmt.Printf("Initialized S3 Storage (Bucket: %s, Region: %s)\n", cfg.AWSBucket, cfg.AWSRegion)
	} else {
		uploadDir := config.ResolveAppPath("uploads")
		baseURL := fmt.Sprintf("http://localhost%s/uploads", cfg.ServerAddr)

		fileService, err = filestorage.NewLocalStorage(uploadDir, baseURL)

		if err != nil {
			log.Fatalf("could not initialize file storage: %v", err)
		}
		fmt.Printf("Initialized Local Storage (Dir: %s)\n", uploadDir)
	}

	challengeService := challenges.NewService(challengeRepo, userRepo, rankService)

	badgeService := services.NewBadgeService(badgeRepo)
	streakService := services.NewStreakService(streakRepo, badgeService)
	moderationService := services.NewModerationService(cfg.OpenAIModerationKey, cfg.ModerationEnabled)
	sqlSuggestionProvider := services.NewSQLSuggestionProvider(db)
	suggestionService := services.NewSuggestionService(sqlSuggestionProvider)
	stripeService := services.NewStripeService(cfg.StripeSecretKey, cfg.StripeWebhookSecret, cfg.StripePriceID, cfg.FrontendBaseURL, subscriptionRepo)
	importerService := importer.NewService(connAcctRepo, mediaRepo, listRepo)
	franchiseDiscovery := services.NewFranchiseDiscoveryService(franchiseRepo, cfg.TMDBApiKey)

	jk := jikan.New()
	tm := tmdb.New(cfg.TMDBApiKey)
	st := steam.New(cfg.SteamWebAPIKey, cfg.PublicBaseURL)
	rg := rawg.New(cfg.RAWGApiKey)
	gb := googlebooks.New(cfg.GoogleBooksAPIKey)
	lf := lastfm.New(cfg.LastFMAPIKey, cfg.LastFMAPISecret, cfg.PublicBaseURL)

	appHandler := handler.NewHandler(
		userRepo, mediaRepo, listRepo, socialRepo, postsRepo, steamRepo, challengeRepo, auditRepo, messagingRepo,
		streakRepo, badgeRepo, collectionRepo, reviewRepo, moderationRepo, subscriptionRepo, adRepo, connAcctRepo, franchiseRepo,
		jk, tm, st, rg, attachmentRepo, gb, lf, lastfmRepo, challengeService, rankService,
		streakService, badgeService, moderationService, suggestionService, stripeService, importerService, franchiseDiscovery,
		wsHub, fileService, cfg,
	)

	router := routes.NewRouter(appHandler, cfg)

	fmt.Printf("Server starting on port %s\n", cfg.ServerAddr)
	fmt.Printf("Swagger documentation available at http://localhost%s/swagger/index.html\n", cfg.ServerAddr)

	if err := http.ListenAndServe(cfg.ServerAddr, router); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
