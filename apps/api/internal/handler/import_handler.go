package handler

import (
	"net/http"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetConnections(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	accounts, err := h.ConnectedAcctRepo.GetByUserID(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get connections")
		return
	}
	httpx.JSON(w, http.StatusOK, accounts)
}

func (h *Handler) ConnectProvider(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	redirectURI := r.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = h.Config.FrontendBaseURL + "/settings/connections"
	}

	url, err := h.ImporterService.GetAuthURL(provider, redirectURI)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *Handler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	code := r.URL.Query().Get("code")
	if code == "" {
		httpx.JSONError(w, http.StatusBadRequest, "Missing authorization code")
		return
	}

	redirectURI := h.Config.PublicBaseURL + "/api/connections/" + provider + "/callback"

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	_, err := h.ImporterService.HandleCallback(userID, provider, code, redirectURI)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to connect account: "+err.Error())
		return
	}

	http.Redirect(w, r, h.Config.FrontendBaseURL+"/settings/connections?connected="+provider, http.StatusTemporaryRedirect)
}

func (h *Handler) DisconnectProvider(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	provider := chi.URLParam(r, "provider")

	if err := h.ConnectedAcctRepo.Delete(userID, provider); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to disconnect")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) TriggerSync(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	provider := chi.URLParam(r, "provider")

	job, err := h.ImporterService.TriggerSync(userID, provider)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to trigger sync: "+err.Error())
		return
	}
	httpx.JSON(w, http.StatusAccepted, job)
}

func (h *Handler) GetImportJobs(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	jobs, err := h.ConnectedAcctRepo.GetImportJobs(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get import jobs")
		return
	}
	httpx.JSON(w, http.StatusOK, jobs)
}
