package handler

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetActiveAd(w http.ResponseWriter, r *http.Request) {
	placement := r.URL.Query().Get("placement")
	if placement == "" {
		placement = "feed"
	}
	campaigns, err := h.AdRepo.GetActiveCampaigns(placement)
	if err != nil || len(campaigns) == 0 {
		httpx.JSON(w, http.StatusOK, nil)
		return
	}
	httpx.JSON(w, http.StatusOK, campaigns[0])
}

func (h *Handler) RecordAdImpression(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	campaignID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}
	var in struct {
		Type string `json:"type"`
	}
	_ = json.NewDecoder(r.Body).Decode(&in)
	if in.Type == "" {
		in.Type = "view"
	}

	var userID *int64
	if uid, ok := middleware.UserIDFromContext(r.Context()); ok {
		userID = &uid
	}

	ipHash := fmt.Sprintf("%x", sha256.Sum256([]byte(r.RemoteAddr)))
	ua := r.UserAgent()

	imp := &models.AdImpression{
		CampaignID:     campaignID,
		UserID:         userID,
		ImpressionType: in.Type,
		IPHash:         &ipHash,
		UserAgent:      &ua,
	}
	_ = h.AdRepo.RecordImpression(imp)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AdminCreateAdCampaign(w http.ResponseWriter, r *http.Request) {
	var in models.NewAdCampaign
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	campaign := &models.AdCampaign{
		Name:           in.Name,
		AdvertiserName: in.AdvertiserName,
		ContentHTML:    in.ContentHTML,
		ImageURL:       in.ImageURL,
		TargetURL:      in.TargetURL,
		Placement:      in.Placement,
		Status:         models.AdStatusDraft,
		BudgetCents:    in.BudgetCents,
		StartAt:        in.StartAt,
		EndAt:          in.EndAt,
		TargetingRules: in.TargetingRules,
	}
	if err := h.AdRepo.CreateCampaign(campaign); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create campaign")
		return
	}
	httpx.JSON(w, http.StatusCreated, campaign)
}

func (h *Handler) AdminGetAdCampaigns(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	campaigns, err := h.AdRepo.GetAllCampaigns(limit, offset)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get campaigns")
		return
	}
	httpx.JSON(w, http.StatusOK, campaigns)
}

func (h *Handler) AdminUpdateAdCampaign(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}
	var campaign models.AdCampaign
	if err := json.NewDecoder(r.Body).Decode(&campaign); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	if err := h.AdRepo.UpdateCampaign(id, &campaign); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to update campaign")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AdminGetAdCampaignStats(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}
	stats, err := h.AdRepo.GetCampaignStats(id)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get campaign stats")
		return
	}
	httpx.JSON(w, http.StatusOK, stats)
}
