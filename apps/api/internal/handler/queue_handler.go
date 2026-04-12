package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetMyQueue(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	var items []struct {
		ID            int64  `db:"id" json:"id"`
		MediaItemID   string `db:"media_item_id" json:"mediaItemId"`
		QueuePosition *int   `db:"queue_position" json:"queuePosition"`
		Title         string `db:"title" json:"title"`
		CoverImageURL *string `db:"cover_image_url" json:"coverImageUrl,omitempty"`
		ItemType      string `db:"item_type" json:"itemType"`
	}
	err := h.ListRepo.DB.Select(&items, `
		SELECT uli.id, uli.media_item_id, uli.queue_position, mi.title, mi.cover_image_url, mi.item_type
		FROM user_list_items uli
		JOIN media_items mi ON mi.id = uli.media_item_id
		WHERE uli.user_id = $1 AND uli.queue_position IS NOT NULL
		ORDER BY uli.queue_position
	`, userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get queue")
		return
	}
	httpx.JSON(w, http.StatusOK, items)
}

func (h *Handler) AddToQueue(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	itemIDStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var maxPos *int
	_ = h.ListRepo.DB.QueryRow(`
		SELECT MAX(queue_position) FROM user_list_items WHERE user_id = $1 AND queue_position IS NOT NULL
	`, userID).Scan(&maxPos)

	nextPos := 1
	if maxPos != nil {
		nextPos = *maxPos + 1
	}

	_, err = h.ListRepo.DB.Exec(`UPDATE user_list_items SET queue_position = $2, updated_at = NOW() WHERE id = $1 AND user_id = $3`, itemID, nextPos, userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to add to queue")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RemoveFromQueue(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	itemIDStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	_, err = h.ListRepo.DB.Exec(`UPDATE user_list_items SET queue_position = NULL, updated_at = NOW() WHERE id = $1 AND user_id = $2`, itemID, userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to remove from queue")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ReorderQueue(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	var in struct {
		Items []struct {
			ID       int64 `json:"id"`
			Position int   `json:"position"`
		} `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	tx, err := h.ListRepo.DB.Beginx()
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback()

	for _, item := range in.Items {
		_, err := tx.Exec(`UPDATE user_list_items SET queue_position = $2, updated_at = NOW() WHERE id = $1 AND user_id = $3`, item.ID, item.Position, userID)
		if err != nil {
			httpx.JSONError(w, http.StatusInternalServerError, "Failed to reorder")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to commit reorder")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
