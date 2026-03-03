package handler

import (
	"github.com/GATEOPENERZ/completionist-api/internal/config"
	"github.com/GATEOPENERZ/completionist-api/internal/repository"
	"github.com/GATEOPENERZ/completionist-api/internal/service/filestorage"
	"github.com/GATEOPENERZ/completionist-api/internal/services"
	"github.com/GATEOPENERZ/completionist-api/internal/services/challenges"
	"github.com/GATEOPENERZ/completionist-api/internal/services/googlebooks"
	"github.com/GATEOPENERZ/completionist-api/internal/services/jikan"
	"github.com/GATEOPENERZ/completionist-api/internal/services/lastfm"
	"github.com/GATEOPENERZ/completionist-api/internal/services/rawg"
	"github.com/GATEOPENERZ/completionist-api/internal/services/steam"
	"github.com/GATEOPENERZ/completionist-api/internal/services/tmdb"
	"github.com/GATEOPENERZ/completionist-api/internal/websocket"
)

type Handler struct {
	UserRepo         *repository.UserRepository
	MediaRepo        *repository.MediaRepository
	ListRepo         *repository.ListRepository
	SocialRepo       *repository.SocialRepository
	PostsRepo        *repository.PostsRepository
	SteamRepo        *repository.SteamRepository
	AttachmentRepo   *repository.AttachmentRepository
	LastFMRepo       *repository.LastFMRepository
	ChallengeRepo    *repository.ChallengeRepository
	MessagingRepo    *repository.MessagingRepository
	Jikan            *jikan.Client
	TMDB             *tmdb.Client
	Steam            *steam.Client
	RAWG             *rawg.Client
	GoogleBooks      *googlebooks.Client
	LastFM           *lastfm.Client
	ChallengeService *challenges.Service
	RankService      *services.RankService
	AuditRepo        *repository.AuditRepository
	WSHub            *websocket.Hub
	FileService      filestorage.Service
	Config           *config.Config
}

func NewHandler(
	userRepo *repository.UserRepository,
	mediaRepo *repository.MediaRepository,
	listRepo *repository.ListRepository,
	socialRepo *repository.SocialRepository,
	postsRepo *repository.PostsRepository,
	steamRepo *repository.SteamRepository,
	chalRepo *repository.ChallengeRepository,
	auditRepo *repository.AuditRepository,
	messagingRepo *repository.MessagingRepository,
	jk *jikan.Client,
	tm *tmdb.Client,
	st *steam.Client,
	rg *rawg.Client,
	attachmentRepo *repository.AttachmentRepository,
	gb *googlebooks.Client,
	lf *lastfm.Client,
	lastfmRepo *repository.LastFMRepository,
	chalService *challenges.Service,
	rankService *services.RankService,
	wsHub *websocket.Hub,
	fileService filestorage.Service,
	cfg *config.Config,
) *Handler {
	return &Handler{
		UserRepo:         userRepo,
		MediaRepo:        mediaRepo,
		ListRepo:         listRepo,
		SocialRepo:       socialRepo,
		PostsRepo:        postsRepo,
		SteamRepo:        steamRepo,
		ChallengeRepo:    chalRepo,
		AuditRepo:        auditRepo,
		MessagingRepo:    messagingRepo,
		Jikan:            jk,
		TMDB:             tm,
		Steam:            st,
		RAWG:             rg,
		AttachmentRepo:   attachmentRepo,
		GoogleBooks:      gb,
		LastFM:           lf,
		LastFMRepo:       lastfmRepo,
		ChallengeService: chalService,
		RankService:      rankService,
		WSHub:            wsHub,
		FileService:      fileService,
		Config:           cfg,
	}
}