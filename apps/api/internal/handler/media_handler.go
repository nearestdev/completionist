package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nearestdev/completionist/internal/httpx"
	"github.com/nearestdev/completionist/internal/middleware"
	"github.com/nearestdev/completionist/internal/models"
	"github.com/nearestdev/completionist/internal/services/googlebooks"
	"github.com/nearestdev/completionist/internal/services/jikan"
	"github.com/nearestdev/completionist/internal/services/rawg"
	"github.com/nearestdev/completionist/internal/services/tmdb"
)

// @Summary      Get Media Details
// @Description  Fetches media details (resolving from external if needed) and logs the view
// @Tags         Media
// @Accept       json
// @Produce      json
// @Param        source query string true "Source (JIKAN, TMDB)"
// @Param        external_id query string true "External ID"
// @Param        item_type query string true "Item Type (series, movie, manga, etc)"
// @Success      200  {object}  models.MediaItem
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /media/details [get]
func (h *Handler) GetMediaDetails(w http.ResponseWriter, r *http.Request) {
	source := r.URL.Query().Get("source")
	externalID := r.URL.Query().Get("external_id")
	itemTypeStr := r.URL.Query().Get("item_type")

	if source == "" || externalID == "" || itemTypeStr == "" {
		httpx.JSONError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	itemType := models.ItemType(itemTypeStr)

	var userID *int64
	if uid, ok := middleware.UserIDFromContext(r.Context()); ok {
		userID = &uid
	}

	mediaID, err := h.resolveMediaItem(r.Context(), source, externalID, itemType)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to resolve media: "+err.Error())
		return
	}

	metadata := map[string]interface{}{
		"source":      source,
		"external_id": externalID,
	}
	if err := h.AuditRepo.Log(userID, "media", mediaID.String(), "view", metadata); err != nil {
	}

	item, err := h.MediaRepo.GetByID(mediaID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to fetch media details: "+err.Error())
		return
	}

	httpx.JSON(w, http.StatusOK, item)
}

// @Summary      Get Trending Media
// @Description  Returns a list of trending media items based on list adds and recent views
// @Tags         Media
// @Accept       json
// @Produce      json
// @Param        limit query int false "Limit number of items (default 10)"
// @Success      200  {array}   models.MediaItem
// @Failure      500  {object}  map[string]string
// @Router       /media/trending [get]
func (h *Handler) GetTrendingMedia(w http.ResponseWriter, r *http.Request) {
	limit := 20

	trends, err := h.MediaRepo.GetTrending(limit)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to fetch trends")
		return
	}

	httpx.JSON(w, http.StatusOK, trends)
}

// @Summary      Get Media By ID
// @Description  Fetches a media item by its UUID
// @Tags         Media
// @Accept       json
// @Produce      json
// @Param        id path string true "Media Item ID (UUID)"
// @Success      200  {object}  models.MediaItem
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /media/{id} [get]
func (h *Handler) GetMediaByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		httpx.JSONError(w, http.StatusBadRequest, "Missing media ID")
		return
	}

	mediaID, err := uuid.Parse(idStr)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid media ID format")
		return
	}

	item, err := h.MediaRepo.GetByID(mediaID)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "Media not found")
		return
	}

	httpx.JSON(w, http.StatusOK, item)
}

func (h *Handler) resolveMediaItem(ctx context.Context, source, externalID string, itemType models.ItemType) (uuid.UUID, error) {
	var id uuid.UUID
	err := h.MediaRepo.DB.GetContext(ctx, &id, "SELECT id FROM media_items WHERE source = $1 AND external_id = $2", source, externalID)
	if err == nil {
		return id, nil
	}

	var newItem models.NewMediaItem
	switch source {
	case "JIKAN":
		malID, err := strconv.Atoi(externalID)
		if err != nil {
			return uuid.Nil, errors.New("invalid external_id")
		}

		if itemType == models.ItemTypeSeries || itemType == models.ItemTypeMovie {
			anime, err := h.Jikan.GetAnime(ctx, malID)
			if err != nil {
				return uuid.Nil, err
			}
			newItem = jikanAnimeToMediaItem(anime, source, externalID)
		} else {
			manga, err := h.Jikan.GetManga(ctx, malID)
			if err != nil {
				return uuid.Nil, err
			}
			newItem = jikanMangaToMediaItem(manga, source, externalID)
		}
	case "TMDB":
		tmdbID, err := strconv.Atoi(externalID)
		if err != nil {
			return uuid.Nil, errors.New("invalid external_id")
		}

		if itemType == models.ItemTypeMovie {
			movie, poster, err := h.TMDB.GetMovie(ctx, tmdbID)
			if err != nil {
				return uuid.Nil, err
			}
			newItem = tmdbMovieToMediaItem(movie, poster, source, externalID)
		} else {
			show, poster, err := h.TMDB.GetTV(ctx, tmdbID)
			if err != nil {
				return uuid.Nil, err
			}
			newItem = tmdbTVToMediaItem(show, poster, source, externalID)
		}
	case "RAWG":
		rawgID, err := strconv.Atoi(externalID)
		if err != nil {
			return uuid.Nil, errors.New("invalid external_id")
		}

		game, err := h.RAWG.GetGame(ctx, rawgID)
		if err != nil {
			return uuid.Nil, err
		}
		newItem = rawgGameToMediaItem(game, source, externalID)
	case "GOOGLE_BOOKS":
		book, err := h.GoogleBooks.GetVolume(ctx, externalID)
		if err != nil {
			return uuid.Nil, err
		}
		newItem = googleBookToMediaItem(book, source, externalID)
	default:
		return uuid.Nil, errors.New("unsupported source")
	}

	media, err := h.MediaRepo.FindOrCreate(newItem)
	if err != nil {
		return uuid.Nil, err
	}
	return media.ID, nil
}

func jikanAnimeToMediaItem(an *jikan.Anime, source, extID string) models.NewMediaItem {
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

	mdBytes, _ := json.Marshal(an)
	md := json.RawMessage(mdBytes)

	return models.NewMediaItem{
		ItemType:      models.ItemTypeSeries,
		Source:        &source,
		ExternalID:    &extID,
		Title:         an.Title,
		Description:   desc,
		CoverImageURL: cover,
		ReleaseDate:   release,
		Genres:        genres,
		Metadata:      &md,
	}
}

func jikanMangaToMediaItem(mg *jikan.Manga, source, extID string) models.NewMediaItem {
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

	mdBytes, _ := json.Marshal(mg)
	md := json.RawMessage(mdBytes)

	return models.NewMediaItem{
		ItemType:      models.ItemTypeManga,
		Source:        &source,
		ExternalID:    &extID,
		Title:         mg.Title,
		Description:   desc,
		CoverImageURL: cover,
		ReleaseDate:   release,
		Genres:        genres,
		Metadata:      &md,
	}
}

func tmdbMovieToMediaItem(mv *tmdb.Movie, poster *string, source, extID string) models.NewMediaItem {
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
	genres := make([]string, 0, len(mv.Genres))
	for _, g := range mv.Genres {
		genres = append(genres, g.Name)
	}

	return models.NewMediaItem{
		ItemType:      models.ItemTypeMovie,
		Source:        &source,
		ExternalID:    &extID,
		Title:         mv.Title,
		Description:   desc,
		CoverImageURL: poster,
		ReleaseDate:   release,
		Genres:        genres,
		Metadata:      nil,
	}
}

func tmdbTVToMediaItem(tv *tmdb.TV, poster *string, source, extID string) models.NewMediaItem {
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
	genres := make([]string, 0, len(tv.Genres))
	for _, g := range tv.Genres {
		genres = append(genres, g.Name)
	}

	return models.NewMediaItem{
		ItemType:      models.ItemTypeSeries,
		Source:        &source,
		ExternalID:    &extID,
		Title:         tv.Name,
		Description:   desc,
		CoverImageURL: poster,
		ReleaseDate:   release,
		Genres:        genres,
		Metadata:      nil,
	}
}

func rawgGameToMediaItem(game *rawg.Game, source, extID string) models.NewMediaItem {
	var release *time.Time
	if game.Released != "" {
		if t, err := time.Parse("2006-01-02", game.Released); err == nil {
			release = &t
		}
	}
	var cover *string
	if game.BackgroundImage != "" {
		c := game.BackgroundImage
		cover = &c
	}
	var desc *string
	if game.Description != "" {
		d := game.Description
		desc = &d
	}
	genres := make([]string, 0, len(game.Genres))
	for _, g := range game.Genres {
		genres = append(genres, g.Name)
	}

	return models.NewMediaItem{
		ItemType:      models.ItemTypeGame,
		Source:        &source,
		ExternalID:    &extID,
		Title:         game.Name,
		Description:   desc,
		CoverImageURL: cover,
		ReleaseDate:   release,
		Genres:        genres,
		Metadata:      nil,
	}
}

func googleBookToMediaItem(book *googlebooks.Book, source, extID string) models.NewMediaItem {
	var release *time.Time
	if book.VolumeInfo.PublishedDate != "" {
		layouts := []string{"2006-01-02", "2006-01", "2006"}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, book.VolumeInfo.PublishedDate); err == nil {
				release = &t
				break
			}
		}
	}
	var cover *string
	if book.VolumeInfo.ImageLinks.Thumbnail != "" {
		c := book.VolumeInfo.ImageLinks.Thumbnail
		cover = &c
	}
	var desc *string
	if book.VolumeInfo.Description != "" {
		d := book.VolumeInfo.Description
		desc = &d
	}

	return models.NewMediaItem{
		ItemType:      models.ItemTypeBook,
		Source:        &source,
		ExternalID:    &extID,
		Title:         book.VolumeInfo.Title,
		Description:   desc,
		CoverImageURL: cover,
		ReleaseDate:   release,
		Genres:        book.VolumeInfo.Categories,
		Metadata:      nil,
	}
}
