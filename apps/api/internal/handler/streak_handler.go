package handler

import (
	"net/http"
	"strconv"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetUserStreak(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	streak, err := h.StreakRepo.GetByUserID(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get streak")
		return
	}
	if streak == nil {
		httpx.JSON(w, http.StatusOK, map[string]interface{}{
			"userId": userID, "currentStreak": 0, "longestStreak": 0, "tier": "none",
		})
		return
	}
	httpx.JSON(w, http.StatusOK, streak)
}

func (h *Handler) GetStreakLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	streaks, err := h.StreakRepo.GetLeaderboard(limit)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get leaderboard")
		return
	}
	httpx.JSON(w, http.StatusOK, streaks)
}
