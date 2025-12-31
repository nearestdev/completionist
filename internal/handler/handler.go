package handler

import (
	"github.com/GATEOPENERZ/completionist-api/internal/config"
	"github.com/GATEOPENERZ/completionist-api/internal/repository"
	"github.com/GATEOPENERZ/completionist-api/internal/services/googlebooks"
	"github.com/GATEOPENERZ/completionist-api/internal/services/jikan"
	"github.com/GATEOPENERZ/completionist-api/internal/services/lastfm"
	"github.com/GATEOPENERZ/completionist-api/internal/services/rawg"
	"github.com/GATEOPENERZ/completionist-api/internal/services/steam"
	"github.com/GATEOPENERZ/completionist-api/internal/services/tmdb"
)

type Handler struct {
	UserRepo       *repository.UserRepository
	MediaRepo      *repository.MediaRepository
	ListRepo       *repository.ListRepository
	SocialRepo     *repository.SocialRepository
	PostsRepo      *repository.PostsRepository
	SteamRepo      *repository.SteamRepository
	AttachmentRepo *repository.AttachmentRepository
	LastFMRepo     *repository.LastFMRepository
	Jikan          *jikan.Client
	TMDB           *tmdb.Client
	Steam          *steam.Client
	RAWG           *rawg.Client
	GoogleBooks    *googlebooks.Client
	LastFM         *lastfm.Client
	Config         *config.Config
}

func NewHandler(
	userRepo *repository.UserRepository,
	mediaRepo *repository.MediaRepository,
	listRepo *repository.ListRepository,
	socialRepo *repository.SocialRepository,
	postsRepo *repository.PostsRepository,
	steamRepo *repository.SteamRepository,
	jk *jikan.Client,
	tm *tmdb.Client,
	st *steam.Client,
	rg *rawg.Client,
	attachmentRepo *repository.AttachmentRepository,
	gb *googlebooks.Client,
	lf *lastfm.Client,
	lastfmRepo *repository.LastFMRepository,
	cfg *config.Config,
) *Handler {
	return &Handler{
		UserRepo:       userRepo,
		MediaRepo:      mediaRepo,
		ListRepo:       listRepo,
		SocialRepo:     socialRepo,
		PostsRepo:      postsRepo,
		SteamRepo:      steamRepo,
		Jikan:          jk,
		TMDB:           tm,
		Steam:          st,
		RAWG:           rg,
		AttachmentRepo: attachmentRepo,
		GoogleBooks:    gb,
		LastFM:         lf,
		LastFMRepo:     lastfmRepo,
		Config:         cfg,
	}
}