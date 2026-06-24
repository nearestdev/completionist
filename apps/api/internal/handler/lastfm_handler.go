package handler

import (
	"net/http"
	"strconv"

	"github.com/nearestdev/completionist/internal/auth"
	"github.com/nearestdev/completionist/internal/httpx"
	"github.com/nearestdev/completionist/internal/middleware"
	"github.com/nearestdev/completionist/internal/models"
)

// @Summary      Initiate Last.fm Auth
// @Description  Redirects the user to Last.fm to authorize the application
// @Tags         Last.fm
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      302
// @Failure      401  {object}  map[string]string
// @Router       /auth/lastfm [get]
func (h *Handler) LastFMAuth(w http.ResponseWriter, r *http.Request) {
	token, ok := middleware.AuthTokenFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	authURL := h.LastFM.GetAuthURL(token)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// @Summary      Last.fm Callback
// @Description  Handles the callback from Last.fm after authorization
// @Tags         Last.fm
// @Accept       json
// @Produce      json
// @Param        token query string true "Last.fm Token"
// @Param        auth_token query string true "App Auth Token"
// @Success      302
// @Failure      302  {string}  string "Redirects to settings with error/success"
// @Router       /auth/lastfm/callback [get]
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

// @Summary      Get My Last.fm Account
// @Description  Retrieves the linked Last.fm account details
// @Tags         Last.fm
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.LastFMAccount
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /me/lastfm [get]
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

// @Summary      Get Recent Tracks
// @Description  Retrieves recent tracks from the linked Last.fm account
// @Tags         Last.fm
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit (max 50)"
// @Success      200  {object}  lastfm.RecentTracksResponse
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      502  {object}  map[string]string
// @Router       /me/lastfm/recent [get]
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
