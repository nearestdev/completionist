package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nearestdev/completionist/internal/httpx"
	"github.com/nearestdev/completionist/internal/middleware"
	"github.com/nearestdev/completionist/internal/models"
)

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	var in models.NewReview
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	if in.Rating < 1 || in.Rating > 10 {
		httpx.JSONError(w, http.StatusBadRequest, "Rating must be between 1 and 10")
		return
	}

	review := &models.CompletionReview{
		UserID:        userID,
		MediaItemID:   in.MediaItemID,
		Rating:        in.Rating,
		Tags:          pq.StringArray(in.Tags),
		FavoriteQuote: in.FavoriteQuote,
		ReviewText:    in.ReviewText,
		CompletedAt:   in.CompletedAt,
	}
	if err := h.ReviewRepo.Create(review); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create review")
		return
	}
	httpx.JSON(w, http.StatusCreated, review)
}

func (h *Handler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	limit, offset := parsePagination(r)
	reviews, err := h.ReviewRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get reviews")
		return
	}
	httpx.JSON(w, http.StatusOK, reviews)
}

func (h *Handler) GetMediaReviews(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	mediaItemID, err := uuid.Parse(idStr)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid media item ID")
		return
	}
	limit, offset := parsePagination(r)
	reviews, err := h.ReviewRepo.GetByMediaID(mediaItemID, limit, offset)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get reviews")
		return
	}
	httpx.JSON(w, http.StatusOK, reviews)
}

func (h *Handler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid review ID")
		return
	}
	_ = userID

	var in models.UpdateReview
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	updated, err := h.ReviewRepo.Update(id, in.Rating, in.Tags, in.FavoriteQuote, in.ReviewText)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to update review")
		return
	}
	if updated == nil {
		httpx.JSONError(w, http.StatusNotFound, "Review not found")
		return
	}
	httpx.JSON(w, http.StatusOK, updated)
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid review ID")
		return
	}
	if err := h.ReviewRepo.Delete(id); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to delete review")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parsePagination(r *http.Request) (int, int) {
	limit := 20
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return limit, offset
}
