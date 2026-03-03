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

// @Summary      Add item to user list
// @Description  Creates a new list item, potentially creating the media item if it doesn't exist
// @Tags         Lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateListItemPayload true "List item details"
// @Success      201  {object}  models.UserListItem
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lists [post]
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
		MediaItemID:  media.ID,
		Status:       payload.ListData.Status,
		ProgressData: payload.ListData.ProgressData,
		Rating:       payload.ListData.Rating,
	}

	item, err := h.ListRepo.CreateUserListItem(userID, newList)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create list item")
		return
	}

	if item.Status == models.StatusCompleted {
		desc := "Completed " + payload.MediaData.Title
		sid := strconv.FormatInt(item.ID, 10)
		_ = h.UserRepo.AddXP(userID, models.XPSourceListCompletion, &sid, &desc)

		meta := map[string]interface{}{
			"item_type":     string(media.ItemType),
			"genres":        media.Genres,
			"media_item_id": media.ID.String(),
		}
		_ = h.ChallengeService.NotifyAction(r.Context(), userID, "complete_item", meta)
	}

	httpx.JSON(w, http.StatusCreated, item)
}

// @Summary      Get user list items
// @Description  Retrieves all items in the logged-in user's list
// @Tags         Lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.UserListItem
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lists [get]
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

// @Summary      Update list item
// @Description  Updates status, progress, or rating of a list item
// @Tags         Lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        item_id path int true "List Item ID"
// @Param        request body models.UpdateUserListItem true "Update details"
// @Success      200  {object}  models.UserListItem
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lists/{item_id} [patch]
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

	existing, err := h.ListRepo.GetUserListItem(userID, itemID)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "Item not found")
		return
	}

	item, err := h.ListRepo.UpdateUserListItem(userID, itemID, upd)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to update list item")
		return
	}

	if existing.Status != models.StatusCompleted && item.Status == models.StatusCompleted {
		media, _ := h.MediaRepo.GetByID(item.MediaItemID)
		title := "item"
		if media != nil {
			title = media.Title
		}
		
		desc := "Completed " + title
		sid := strconv.FormatInt(item.ID, 10)
		_ = h.UserRepo.AddXP(userID, models.XPSourceListCompletion, &sid, &desc)

		if media != nil {
			meta := map[string]interface{}{
				"item_type":     string(media.ItemType),
				"genres":        media.Genres,
				"media_item_id": media.ID.String(),
			}
			_ = h.ChallengeService.NotifyAction(r.Context(), userID, "complete_item", meta)
		}
	}

	httpx.JSON(w, http.StatusOK, item)
}

// @Summary      Delete list item
// @Description  Removes an item from the user's list
// @Tags         Lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        item_id path int true "List Item ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lists/{item_id} [delete]
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

// @Summary      Add to wishlist
// @Description  Adds a media item to the user's wishlist
// @Tags         Wishlist
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.NewWishlistItem true "Wishlist item details"
// @Success      201  {object}  models.WishlistItem
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /wishlist [post]
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

// @Summary      Get wishlist
// @Description  Retrieves all items in the user's wishlist
// @Tags         Wishlist
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.WishlistItem
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /wishlist [get]
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

// @Summary      Remove from wishlist
// @Description  Removes an item from the user's wishlist
// @Tags         Wishlist
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        item_id path int true "Wishlist Item ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /wishlist/{item_id} [delete]
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