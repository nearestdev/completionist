package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		httpx.JSONError(w, http.StatusBadRequest, "q is required")
		return
	}

	res, err := h.GoogleBooks.Search(r.Context(), q)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Google Books API error")
		return
	}

	httpx.JSON(w, http.StatusOK, res)
}

type addFromGoogleBooksBody struct {
	Status   models.ItemStatus `json:"status"`
	Progress *string           `json:"progress,omitempty"`
	Rating   *int              `json:"rating,omitempty"`
}

func (h *Handler) AddBookFromGoogleBooks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	bookID := chi.URLParam(r, "book_id")
	if bookID == "" {
		httpx.JSONError(w, http.StatusBadRequest, "invalid book_id")
		return
	}

	var body addFromGoogleBooksBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	book, err := h.GoogleBooks.GetVolume(r.Context(), bookID)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Failed fetching from Google Books")
		return
	}

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

	source := "GOOGLE_BOOKS"
	media, err := h.MediaRepo.FindOrCreate(models.NewMediaItem{
		ItemType:      models.ItemTypeBook,
		Source:        &source,
		ExternalID:    &book.ID,
		Title:         book.VolumeInfo.Title,
		Description:   desc,
		CoverImageURL: cover,
		ReleaseDate:   release,
		Genres:        book.VolumeInfo.Categories,
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