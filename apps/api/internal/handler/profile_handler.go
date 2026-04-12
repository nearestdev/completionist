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

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	role, _ := middleware.UserRoleFromContext(r.Context())
	if !models.UserRole(role).IsMemberOrAbove() {
		httpx.JSONError(w, http.StatusForbidden, "Membership required for profile customization")
		return
	}

	var in struct {
		ProfileHTML  *string          `json:"profileHtml,omitempty"`
		ProfileTheme *json.RawMessage `json:"profileTheme,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	_, err := h.UserRepo.DB.Exec(`
		UPDATE users SET profile_html = COALESCE($2, profile_html), profile_theme = COALESCE($3, profile_theme), updated_at = NOW()
		WHERE id = $1
	`, userID, in.ProfileHTML, in.ProfileTheme)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCustomProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var result struct {
		ProfileHTML  *string          `db:"profile_html" json:"profileHtml,omitempty"`
		ProfileTheme *json.RawMessage `db:"profile_theme" json:"profileTheme,omitempty"`
	}
	err = h.UserRepo.DB.QueryRowx(`SELECT profile_html, profile_theme FROM users WHERE id = $1`, userID).StructScan(&result)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "User not found")
		return
	}
	httpx.JSON(w, http.StatusOK, result)
}
