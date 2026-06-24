package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nearestdev/completionist/internal/httpx"
	"github.com/nearestdev/completionist/internal/middleware"
	"github.com/nearestdev/completionist/internal/models"
)

func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	var in models.NewCollection
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	if in.Name == "" {
		httpx.JSONError(w, http.StatusBadRequest, "Name is required")
		return
	}
	isPrivate := false
	if in.IsPrivate != nil {
		isPrivate = *in.IsPrivate
	}
	c := &models.Collection{UserID: userID, Name: in.Name, Description: in.Description, IsPrivate: isPrivate}
	if err := h.CollectionRepo.Create(c); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create collection")
		return
	}
	httpx.JSON(w, http.StatusCreated, c)
}

func (h *Handler) GetMyCollections(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	collections, err := h.CollectionRepo.GetByUserID(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get collections")
		return
	}
	httpx.JSON(w, http.StatusOK, collections)
}

func (h *Handler) GetCollection(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid collection ID")
		return
	}

	c, err := h.CollectionRepo.GetByID(id)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get collection")
		return
	}
	if c == nil {
		httpx.JSONError(w, http.StatusNotFound, "Collection not found")
		return
	}

	userID, _ := middleware.UserIDFromContext(r.Context())
	if c.IsPrivate && c.UserID != userID {
		httpx.JSONError(w, http.StatusNotFound, "Collection not found")
		return
	}

	items, err := h.CollectionRepo.GetItemsByCollectionID(id)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get collection items")
		return
	}
	httpx.JSON(w, http.StatusOK, models.CollectionWithItems{Collection: *c, Items: items})
}

func (h *Handler) UpdateCollection(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid collection ID")
		return
	}

	existing, err := h.CollectionRepo.GetByID(id)
	if err != nil || existing == nil || existing.UserID != userID {
		httpx.JSONError(w, http.StatusNotFound, "Collection not found")
		return
	}

	var in models.UpdateCollection
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	updated, err := h.CollectionRepo.Update(id, in.Name, in.Description, in.IsPrivate)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to update collection")
		return
	}
	httpx.JSON(w, http.StatusOK, updated)
}

func (h *Handler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid collection ID")
		return
	}

	existing, err := h.CollectionRepo.GetByID(id)
	if err != nil || existing == nil || existing.UserID != userID {
		httpx.JSONError(w, http.StatusNotFound, "Collection not found")
		return
	}

	if err := h.CollectionRepo.Delete(id); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to delete collection")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AssignItemToCollection(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var in struct {
		CollectionID int64 `json:"collectionId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if err := h.CollectionRepo.AssignItem(itemID, in.CollectionID); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to assign item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UnassignItemFromCollection(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	if err := h.CollectionRepo.UnassignItem(itemID); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to unassign item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
