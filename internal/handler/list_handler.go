package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) CreateListItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var payload models.CreateListItemPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	media, err := h.MediaRepo.FindOrCreate(payload.MediaData)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to upsert media item")
		return
	}
	newList := models.NewUserListItem{
		MediaItemID: media.ID,
		Status:      payload.ListData.Status,
		Progress:    payload.ListData.Progress,
		Rating:      payload.ListData.Rating,
	}
	item, err := h.ListRepo.CreateUserListItem(userID, newList)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create list item")
		return
	}
	httpx.JSON(w, http.StatusCreated, item)
}

func (h *Handler) GetMyListItems(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	items, err := h.ListRepo.GetUserListItems(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get list items")
		return
	}
	httpx.JSON(w, http.StatusOK, items)
}

func (h *Handler) UpdateMyListItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid item_id")
		return
	}
	var upd models.UpdateUserListItem
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	item, err := h.ListRepo.UpdateUserListItem(userID, itemID, upd)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to update list item")
		return
	}
	httpx.JSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteMyListItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid item_id")
		return
	}
	affected, err := h.ListRepo.DeleteUserListItem(userID, itemID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to delete list item")
		return
	}
	if affected == 0 {
		httpx.JSONError(w, http.StatusNotFound, "Not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AddToWishlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var payload models.NewWishlistItem
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	item, err := h.ListRepo.CreateWishlistItem(userID, payload)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create wishlist item")
		return
	}
	httpx.JSON(w, http.StatusCreated, item)
}

func (h *Handler) GetMyWishlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	items, err := h.ListRepo.GetWishlistItems(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get wishlist")
		return
	}
	httpx.JSON(w, http.StatusOK, items)
}

func (h *Handler) RemoveFromWishlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	idStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid item_id")
		return
	}
	affected, err := h.ListRepo.DeleteWishlistItem(userID, itemID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to delete wishlist item")
		return
	}
	if affected == 0 {
		httpx.JSONError(w, http.StatusNotFound, "Not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
