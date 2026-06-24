package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nearestdev/completionist/internal/httpx"
	"github.com/nearestdev/completionist/internal/middleware"
	"github.com/nearestdev/completionist/internal/models"
)

func (h *Handler) ListFranchises(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	limit, offset := parsePagination(r)

	if category != "" && h.FranchiseDiscovery != nil {
		franchises, err := h.FranchiseDiscovery.DiscoverByCategory(category)
		if err == nil && len(franchises) > 0 {
			httpx.JSON(w, http.StatusOK, franchises)
			return
		}
		dbFranchises, err := h.FranchiseRepo.GetByCategory(category, limit, offset)
		if err != nil {
			httpx.JSONError(w, http.StatusInternalServerError, "Failed to list franchises")
			return
		}
		httpx.JSON(w, http.StatusOK, dbFranchises)
		return
	}

	franchises, err := h.FranchiseRepo.GetAll(limit, offset)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to list franchises")
		return
	}
	if len(franchises) == 0 && h.FranchiseDiscovery != nil {
		discovered, err := h.FranchiseDiscovery.DiscoverMovieFranchises()
		if err == nil && len(discovered) > 0 {
			httpx.JSON(w, http.StatusOK, discovered)
			return
		}
	}
	httpx.JSON(w, http.StatusOK, franchises)
}

func (h *Handler) GetFranchise(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid franchise ID")
		return
	}

	franchise, err := h.FranchiseRepo.GetByID(id)
	if err != nil || franchise == nil {
		httpx.JSONError(w, http.StatusNotFound, "Franchise not found")
		return
	}

	var viewingUserID *int64
	if uid, ok := middleware.UserIDFromContext(r.Context()); ok {
		viewingUserID = &uid
	}

	items, err := h.FranchiseRepo.GetFranchiseItems(id, viewingUserID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get franchise items")
		return
	}

	completed := 0
	for _, item := range items {
		if item.IsCompleted {
			completed++
		}
	}

	result := models.FranchiseWithItems{
		Franchise:  *franchise,
		Items:      items,
		TotalItems: len(items),
		Completed:  completed,
	}
	httpx.JSON(w, http.StatusOK, result)
}

func (h *Handler) AdminCreateFranchise(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	var in models.NewFranchise
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	if in.Name == "" {
		httpx.JSONError(w, http.StatusBadRequest, "Name is required")
		return
	}
	f := &models.Franchise{Name: in.Name, Description: in.Description, CoverImageURL: in.CoverImageURL, CreatedBy: &userID}
	if err := h.FranchiseRepo.Create(f); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create franchise")
		return
	}
	httpx.JSON(w, http.StatusCreated, f)
}

func (h *Handler) AdminUpdateFranchise(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid franchise ID")
		return
	}
	var in models.NewFranchise
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	updated, err := h.FranchiseRepo.Update(id, &in.Name, in.Description, in.CoverImageURL)
	if err != nil || updated == nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to update franchise")
		return
	}
	httpx.JSON(w, http.StatusOK, updated)
}

func (h *Handler) AdminAddFranchiseItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	franchiseID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid franchise ID")
		return
	}
	var in models.NewFranchiseItem
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	item := &models.FranchiseItem{
		FranchiseID:  franchiseID,
		MediaItemID:  in.MediaItemID,
		Position:     in.Position,
		Relationship: in.Relationship,
	}
	if err := h.FranchiseRepo.AddItem(franchiseID, item); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to add franchise item")
		return
	}
	httpx.JSON(w, http.StatusCreated, item)
}

func (h *Handler) AdminRemoveFranchiseItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	franchiseID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid franchise ID")
		return
	}
	itemIDStr := chi.URLParam(r, "itemId")
	mediaItemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid media item ID")
		return
	}
	if err := h.FranchiseRepo.RemoveItem(franchiseID, mediaItemID); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to remove franchise item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
