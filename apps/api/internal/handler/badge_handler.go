package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nearestdev/completionist/internal/httpx"
)

func (h *Handler) ListBadges(w http.ResponseWriter, r *http.Request) {
	badges, err := h.BadgeRepo.GetAll()
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to list badges")
		return
	}
	httpx.JSON(w, http.StatusOK, badges)
}

func (h *Handler) GetUserBadges(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	badges, err := h.BadgeRepo.GetUserBadges(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get user badges")
		return
	}
	httpx.JSON(w, http.StatusOK, badges)
}
