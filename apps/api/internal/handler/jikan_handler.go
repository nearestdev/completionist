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
	_ "github.com/nearestdev/completionist/internal/services/jikan"
)

// @Summary      Search Anime (Jikan)
// @Description  Search for anime using the Jikan API (MyAnimeList)
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        q query string true "Search query"
// @Param        page query int false "Page number"
// @Success      200  {object}  jikan.AnimeSearchResponse
// @Failure      400  {object}  map[string]string
// @Failure      502  {object}  map[string]string
// @Router       /search/anime [get]
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

// @Summary      Search Manga (Jikan)
// @Description  Search for manga using the Jikan API (MyAnimeList)
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        q query string true "Search query"
// @Param        page query int false "Page number"
// @Success      200  {object}  jikan.MangaSearchResponse
// @Failure      400  {object}  map[string]string
// @Failure      502  {object}  map[string]string
// @Router       /search/manga [get]
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
	Status       models.ItemStatus `json:"status"`
	ProgressData *json.RawMessage  `json:"progressData,omitempty" swaggertype:"string"`
	Rating       *int              `json:"rating,omitempty"`
}

// @Summary      Add Anime from Jikan
// @Description  Imports an anime from Jikan/MAL into the user's list
// @Tags         Lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        mal_id path int true "MyAnimeList ID"
// @Param        request body addFromJikanBody true "List details"
// @Success      201  {object}  models.UserListItem
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lists/jikan/anime/{mal_id} [post]
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

// @Summary      Add Manga from Jikan
// @Description  Imports a manga from Jikan/MAL into the user's list
// @Tags         Lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        mal_id path int true "MyAnimeList ID"
// @Param        request body addFromJikanBody true "List details"
// @Success      201  {object}  models.UserListItem
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lists/jikan/manga/{mal_id} [post]
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
