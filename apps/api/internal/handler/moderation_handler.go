package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nearestdev/completionist/internal/httpx"
	"github.com/nearestdev/completionist/internal/middleware"
	"github.com/nearestdev/completionist/internal/models"
)

func (h *Handler) ReportContent(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	var in models.NewReport
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	item := &models.ModerationItem{
		EntityType: in.EntityType,
		EntityID:   in.EntityID,
		ReportedBy: &userID,
		Reason:     &in.Reason,
	}
	if err := h.ModerationRepo.CreateItem(item); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to submit report")
		return
	}
	httpx.JSON(w, http.StatusCreated, item)
}

func (h *Handler) AdminGetModerationQueue(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limit, offset := parsePagination(r)
	items, err := h.ModerationRepo.GetQueue(status, limit, offset)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get moderation queue")
		return
	}
	httpx.JSON(w, http.StatusOK, items)
}

func (h *Handler) AdminReviewModerationItem(w http.ResponseWriter, r *http.Request) {
	reviewerID, _ := middleware.UserIDFromContext(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid moderation item ID")
		return
	}
	var in struct {
		Status     models.ModerationStatus `json:"status"`
		ReviewNote *string                 `json:"reviewNote,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	if err := h.ModerationRepo.UpdateStatus(id, in.Status, reviewerID, in.ReviewNote); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to update moderation item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AdminBanUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var in struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	if err := h.ModerationRepo.BanUser(userID, in.Reason); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to ban user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AdminUnbanUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	if err := h.ModerationRepo.UnbanUser(userID); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to unban user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AdminGetBanAppeals(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limit, offset := parsePagination(r)
	appeals, err := h.ModerationRepo.GetBanAppeals(status, limit, offset)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to get ban appeals")
		return
	}
	httpx.JSON(w, http.StatusOK, appeals)
}

func (h *Handler) AdminReviewBanAppeal(w http.ResponseWriter, r *http.Request) {
	reviewerID, _ := middleware.UserIDFromContext(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid appeal ID")
		return
	}
	var in struct {
		Status     models.BanAppealStatus `json:"status"`
		ReviewNote *string                `json:"reviewNote,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	if err := h.ModerationRepo.UpdateBanAppeal(id, in.Status, reviewerID, in.ReviewNote); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to update appeal")
		return
	}
	if in.Status == models.BanAppealAccepted {
		var appeal models.BanAppeal
		_ = h.ModerationRepo.DB.QueryRowx(`SELECT user_id FROM ban_appeals WHERE id = $1`, id).StructScan(&appeal)
		_ = h.ModerationRepo.UnbanUser(appeal.UserID)
	}
	w.WriteHeader(http.StatusNoContent)
}
