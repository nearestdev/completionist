package handler

import (
	"net/http"
	"strconv"

	"github.com/GATEOPENERZ/completionist-api/internal/auth"
	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
)

func (h *Handler) LastFMAuth(w http.ResponseWriter, r *http.Request) {
	token, ok := middleware.AuthTokenFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	authURL := h.LastFM.GetAuthURL(token)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *Handler) LastFMCallback(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	authToken := r.URL.Query().Get("auth_token")
	if token == "" || authToken == "" {
		http.Redirect(w, r, h.Config.FrontendBaseURL+"/settings?error=lastfm_auth_failed", http.StatusFound)
		return
	}

	claims, err := auth.ValidateJWT(authToken)
	if err != nil {
		http.Redirect(w, r, h.Config.FrontendBaseURL+"/settings?error=invalid_auth_token", http.StatusFound)
		return
	}

	session, err := h.LastFM.GetSession(r.Context(), token)
	if err != nil {
		http.Redirect(w, r, h.Config.FrontendBaseURL+"/settings?error=lastfm_session_failed", http.StatusFound)
		return
	}

	account := &models.LastFMAccount{
		UserID:     claims.UserID,
		Username:   session.Session.Name,
		SessionKey: session.Session.Key,
		Subscriber: session.Session.Subscriber,
	}

	if err := h.LastFMRepo.Upsert(account); err != nil {
		http.Redirect(w, r, h.Config.FrontendBaseURL+"/settings?error=db_error", http.StatusFound)
		return
	}

	http.Redirect(w, r, h.Config.FrontendBaseURL+"/settings?success=lastfm_linked", http.StatusFound)
}

func (h *Handler) GetMyLastFMAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	account, err := h.LastFMRepo.FindByUserID(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "Last.fm account not found")
		return
	}

	httpx.JSON(w, http.StatusOK, account)
}

func (h *Handler) GetMyRecentTracks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	account, err := h.LastFMRepo.FindByUserID(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "Last.fm account not found")
		return
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 50 {
			limit = val
		}
	}

	tracks, err := h.LastFM.GetRecentTracks(r.Context(), account.Username, limit)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Failed to fetch recent tracks from Last.fm")
		return
	}

	httpx.JSON(w, http.StatusOK, tracks)
}