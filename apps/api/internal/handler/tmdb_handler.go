package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nearestdev/completionist/internal/httpx"
	"github.com/nearestdev/completionist/internal/middleware"
	"github.com/nearestdev/completionist/internal/models"
)

// @Summary      Search Movies (TMDb)
// @Description  Search for movies using The Movie Database API
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        q query string true "Search query"
// @Param        page query int false "Page number"
// @Success      200  {object}  tmdb.SearchResult
// @Failure      400  {object}  map[string]string
// @Failure      502  {object}  map[string]string
// @Router       /search/movies [get]
func (h *Handler) SearchMoviesTMDB(w http.ResponseWriter, r *http.Request) {
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
	res, err := h.TMDB.SearchMovies(r.Context(), q, page)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "TMDb error")
		return
	}
	httpx.JSON(w, http.StatusOK, res)
}

// @Summary      Search TV Shows (TMDb)
// @Description  Search for TV shows using The Movie Database API
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        q query string true "Search query"
// @Param        page query int false "Page number"
// @Success      200  {object}  tmdb.SearchResult
// @Failure      400  {object}  map[string]string
// @Failure      502  {object}  map[string]string
// @Router       /search/tv [get]
func (h *Handler) SearchTVTMDB(w http.ResponseWriter, r *http.Request) {
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
	res, err := h.TMDB.SearchTV(r.Context(), q, page)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "TMDb error")
		return
	}
	httpx.JSON(w, http.StatusOK, res)
}

type addFromTMDBBody struct {
	Status       models.ItemStatus `json:"status"`
	ProgressData *json.RawMessage  `json:"progressData,omitempty" swaggertype:"string"`
	Rating       *int              `json:"rating,omitempty"`
}

// @Summary      Add Movie from TMDb
// @Description  Imports a movie from TMDb into the user's list
// @Tags         Lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tmdb_id path int true "TMDb Movie ID"
// @Param        request body addFromTMDBBody true "List details"
// @Success      201  {object}  models.UserListItem
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lists/tmdb/movie/{tmdb_id} [post]
func (h *Handler) AddMovieFromTMDB(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "tmdb_id")
	tmdbID, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "invalid tmdb_id")
		return
	}
	var body addFromTMDBBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	mv, poster, err := h.TMDB.GetMovie(r.Context(), tmdbID)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Failed fetching from TMDb")
		return
	}
	var release *time.Time
	if mv.ReleaseDate != "" {
		if t, err := time.Parse("2006-01-02", mv.ReleaseDate); err == nil {
			release = &t
		}
	}
	var desc *string
	if mv.Overview != "" {
		d := mv.Overview
		desc = &d
	}
	var cover *string
	if poster != nil {
		cover = poster
	}
	genres := make([]string, 0, len(mv.Genres))
	for _, g := range mv.Genres {
		genres = append(genres, g.Name)
	}
	source := "TMDB"
	ext := strconv.Itoa(mv.ID)
	media, err := h.MediaRepo.FindOrCreate(models.NewMediaItem{
		ItemType:      models.ItemTypeMovie,
		Source:        &source,
		ExternalID:    &ext,
		Title:         mv.Title,
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
		MediaItemID:  media.ID,
		Status:       body.Status,
		ProgressData: body.ProgressData,
		Rating:       body.Rating,
	})
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "List create failed")
		return
	}
	httpx.JSON(w, http.StatusCreated, item)
}

// @Summary      Add TV Show from TMDb
// @Description  Imports a TV show from TMDb into the user's list
// @Tags         Lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tmdb_id path int true "TMDb TV ID"
// @Param        request body addFromTMDBBody true "List details"
// @Success      201  {object}  models.UserListItem
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lists/tmdb/tv/{tmdb_id} [post]
func (h *Handler) AddTVFromTMDB(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "tmdb_id")
	tmdbID, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "invalid tmdb_id")
		return
	}
	var body addFromTMDBBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	tv, poster, err := h.TMDB.GetTV(r.Context(), tmdbID)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Failed fetching from TMDb")
		return
	}
	var release *time.Time
	if tv.FirstAirDate != "" {
		if t, err := time.Parse("2006-01-02", tv.FirstAirDate); err == nil {
			release = &t
		}
	}
	var desc *string
	if tv.Overview != "" {
		d := tv.Overview
		desc = &d
	}
	var cover *string
	if poster != nil {
		cover = poster
	}
	genres := make([]string, 0, len(tv.Genres))
	for _, g := range tv.Genres {
		genres = append(genres, g.Name)
	}
	source := "TMDB"
	ext := strconv.Itoa(tv.ID)
	title := tv.Name
	media, err := h.MediaRepo.FindOrCreate(models.NewMediaItem{
		ItemType:      models.ItemTypeSeries,
		Source:        &source,
		ExternalID:    &ext,
		Title:         title,
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
		MediaItemID:  media.ID,
		Status:       body.Status,
		ProgressData: body.ProgressData,
		Rating:       body.Rating,
	})
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "List create failed")
		return
	}
	httpx.JSON(w, http.StatusCreated, item)
}
