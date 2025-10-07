package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) SearchAnimeJikan(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		httpx.JSONError(w, http.StatusBadRequest, "q is required")
		return
	}
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}
	res, err := h.Jikan.SearchAnime(r.Context(), q, page)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Jikan error")
		return
	}
	httpx.JSON(w, http.StatusOK, res)
}

func (h *Handler) SearchMangaJikan(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		httpx.JSONError(w, http.StatusBadRequest, "q is required")
		return
	}
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}
	res, err := h.Jikan.SearchManga(r.Context(), q, page)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Jikan error")
		return
	}
	httpx.JSON(w, http.StatusOK, res)
}

type addFromJikanBody struct {
	Status   models.ItemStatus `json:"status"`
	Progress *string           `json:"progress,omitempty"`
	Rating   *int              `json:"rating,omitempty"`
}

func (h *Handler) AddAnimeFromJikan(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "mal_id")
	malID, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "invalid mal_id")
		return
	}

	var body addFromJikanBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	an, err := h.Jikan.GetAnime(r.Context(), malID)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Failed fetching from Jikan")
		return
	}

	var release *time.Time
	if an.Aired.From != nil {
		release = an.Aired.From
	}
	var cover *string
	if an.Images.Webp.ImageURL != "" {
		c := an.Images.Webp.ImageURL
		cover = &c
	} else if an.Images.JPG.ImageURL != "" {
		c := an.Images.JPG.ImageURL
		cover = &c
	}
	var desc *string
	if an.Synopsis != "" {
		d := an.Synopsis
		desc = &d
	}
	genres := make([]string, 0, len(an.Genres))
	for _, g := range an.Genres {
		genres = append(genres, g.Name)
	}

	source := "JIKAN"
	ext := strconv.Itoa(an.MalID)

	media, err := h.MediaRepo.FindOrCreate(models.NewMediaItem{
		ItemType:      models.ItemTypeSeries,
		Source:        &source,
		ExternalID:    &ext,
		Title:         an.Title,
		Description:   desc,
		CoverImageURL: cover,
		ReleaseDate:   release,
		Genres:        genres,
		Metadata:      nil,
	})
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Media upsert failed")
		return
	}

	item, err := h.ListRepo.CreateUserListItem(userID, models.NewUserListItem{
		MediaItemID: media.ID,
		Status:      body.Status,
		Progress:    body.Progress,
		Rating:      body.Rating,
	})
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "List create failed")
		return
	}

	httpx.JSON(w, http.StatusCreated, item)
}

func (h *Handler) AddMangaFromJikan(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "mal_id")
	malID, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "invalid mal_id")
		return
	}

	var body addFromJikanBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	mg, err := h.Jikan.GetManga(r.Context(), malID)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Failed fetching from Jikan")
		return
	}

	var release *time.Time
	if mg.Published.From != nil {
		release = mg.Published.From
	}
	var cover *string
	if mg.Images.Webp.ImageURL != "" {
		c := mg.Images.Webp.ImageURL
		cover = &c
	} else if mg.Images.JPG.ImageURL != "" {
		c := mg.Images.JPG.ImageURL
		cover = &c
	}
	var desc *string
	if mg.Synopsis != "" {
		d := mg.Synopsis
		desc = &d
	}
	genres := make([]string, 0, len(mg.Genres))
	for _, g := range mg.Genres {
		genres = append(genres, g.Name)
	}

	source := "JIKAN"
	ext := strconv.Itoa(mg.MalID)

	media, err := h.MediaRepo.FindOrCreate(models.NewMediaItem{
		ItemType:      models.ItemTypeManga,
		Source:        &source,
		ExternalID:    &ext,
		Title:         mg.Title,
		Description:   desc,
		CoverImageURL: cover,
		ReleaseDate:   release,
		Genres:        genres,
		Metadata:      nil,
	})
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Media upsert failed")
		return
	}

	item, err := h.ListRepo.CreateUserListItem(userID, models.NewUserListItem{
		MediaItemID: media.ID,
		Status:      body.Status,
		Progress:    body.Progress,
		Rating:      body.Rating,
	})
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "List create failed")
		return
	}

	httpx.JSON(w, http.StatusCreated, item)
}
