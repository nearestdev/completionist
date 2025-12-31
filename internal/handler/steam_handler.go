package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GATEOPENERZ/completionist-api/internal/auth"
	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/go-chi/chi/v5"
)
func (h *Handler) SteamLogin(w http.ResponseWriter, r *http.Request) {
	tokenStr, ok := middleware.AuthTokenFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Authorization token missing from context")
		return
	}
	url := h.Steam.BuildOpenIDRedirect(tokenStr)
	http.Redirect(w, r, url, http.StatusFound)
}
func (h *Handler) SteamCallback(w http.ResponseWriter, r *http.Request) {
	redirectURL := h.Config.FrontendBaseURL + h.Config.SteamCallbackRedirectPath
	steamID, err := h.Steam.VerifyOpenID(r)
	if err != nil {
		http.Redirect(w, r, redirectURL+"?error=steam_verification_failed", http.StatusFound)
		return
	}
	authToken := r.URL.Query().Get("auth_token")
	if authToken == "" {
		http.Redirect(w, r, redirectURL+"?error=auth_failed", http.StatusFound)
		return
	}
	claims, err := auth.ValidateJWT(authToken)
	if err != nil {
		http.Redirect(w, r, redirectURL+"?error=invalid_auth_token", http.StatusFound)
		return
	}
	authUserID := claims.UserID
	summary, err := h.Steam.GetPlayerSummary(r.Context(), steamID)
	if err != nil {
		http.Redirect(w, r, redirectURL+"?error=steam_profile_fetch_failed", http.StatusFound)
		return
	}
	acc := models.SteamAccount{
		UserID:  authUserID,
		SteamID: steamID,
		Persona: summary.PersonaName,
		Avatar:  summary.AvatarFull,
	}
	if err := h.SteamRepo.UpsertSteamAccount(&acc); err != nil {
		http.Redirect(w, r, redirectURL+"?error=steam_link_failed", http.StatusFound)
		return
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
func (h *Handler) GetMySteamAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	acc, err := h.SteamRepo.FindByUserID(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "Steam not linked")
		return
	}
	httpx.JSON(w, http.StatusOK, acc)
}
func (h *Handler) GetUserSteamAccount(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.UserRepo.FindByUsername(username)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "User not found")
		return
	}
	acc, err := h.SteamRepo.FindByUserID(user.ID)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound, "Steam not linked for this user")
		return
	}
	httpx.JSON(w, http.StatusOK, acc)
}
func (h *Handler) GetMySteamOwnedGames(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	acc, err := h.SteamRepo.FindByUserID(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Steam not linked")
		return
	}
	g, err := h.Steam.GetOwnedGames(r.Context(), acc.SteamID, true, true)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Steam error")
		return
	}
	httpx.JSON(w, http.StatusOK, g)
}
func (h *Handler) GetMySteamAchievementsForApp(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	acc, err := h.SteamRepo.FindByUserID(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Steam not linked")
		return
	}
	appStr := chi.URLParam(r, "app_id")
	appID, err := strconv.Atoi(appStr)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid app_id")
		return
	}
	res, err := h.Steam.GetPlayerAchievements(r.Context(), acc.SteamID, appID)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Steam error")
		return
	}
	httpx.JSON(w, http.StatusOK, res)
}
func (h *Handler) GetSteamGameSchema(w http.ResponseWriter, r *http.Request) {
	appStr := chi.URLParam(r, "app_id")
	appID, err := strconv.Atoi(appStr)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid app_id")
		return
	}
	res, err := h.Steam.GetSchemaForGame(r.Context(), appID)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Steam error")
		return
	}
	httpx.JSON(w, http.StatusOK, res)
}
func (h *Handler) RAWGSearchGames(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}
	if q == "" {
		httpx.JSONError(w, http.StatusBadRequest, "q is required")
		return
	}
	res, err := h.RAWG.SearchGames(r.Context(), q, page)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "RAWG error")
		return
	}
	httpx.JSON(w, http.StatusOK, res)
}
func (h *Handler) RAWGGameAchievements(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "rawg_id")
	gameID, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid id")
		return
	}
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}
	res, err := h.RAWG.GameAchievements(r.Context(), gameID, page)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "RAWG error")
		return
	}
	httpx.JSON(w, http.StatusOK, res)
}
func (h *Handler) AttachSteamToUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var in models.AttachSteamRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	if in.SteamID == "" {
		httpx.JSONError(w, http.StatusBadRequest, "steamId is required")
		return
	}
	summary, err := h.Steam.GetPlayerSummary(r.Context(), in.SteamID)
	if err != nil {
		httpx.JSONError(w, http.StatusBadGateway, "Steam profile fetch failed")
		return
	}
	acc := models.SteamAccount{
		UserID:  userID,
		SteamID: in.SteamID,
		Persona: summary.PersonaName,
		Avatar:  summary.AvatarFull,
	}
	if err := h.SteamRepo.UpsertSteamAccount(&acc); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Linking Steam failed")
		return
	}
	httpx.JSON(w, http.StatusOK, acc)
}