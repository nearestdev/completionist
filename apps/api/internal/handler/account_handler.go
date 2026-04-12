package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
)

func (h *Handler) RequestAccountDeletion(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	var in struct {
		Reason *string `json:"reason,omitempty"`
	}
	_ = json.NewDecoder(r.Body).Decode(&in)

	req, err := h.ModerationRepo.CreateDeletionRequest(userID, in.Reason)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create deletion request")
		return
	}
	httpx.JSON(w, http.StatusCreated, req)
}

func (h *Handler) CancelAccountDeletion(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	if err := h.ModerationRepo.CancelDeletionRequest(userID); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to cancel deletion request")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ExportUserData(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	data, err := h.ModerationRepo.ExportUserData(userID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to export data")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=my-data.json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) SubmitBanAppeal(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	var in struct {
		AppealText string `json:"appealText"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.AppealText == "" {
		httpx.JSONError(w, http.StatusBadRequest, "Appeal text is required")
		return
	}
	appeal := &models.BanAppeal{UserID: userID, AppealText: in.AppealText}
	if err := h.ModerationRepo.CreateBanAppeal(appeal); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to submit appeal")
		return
	}
	httpx.JSON(w, http.StatusCreated, appeal)
}
