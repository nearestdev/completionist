package handler

import (
	"net/http"
	"strconv"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
)

func (h *Handler) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	limit := 10
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 50 {
			limit = n
		}
	}

	items, err := h.SuggestionService.GetSuggestions(userID, limit)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get suggestions")
		return
	}
	httpx.JSON(w, http.StatusOK, items)
}
