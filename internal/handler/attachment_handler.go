package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GATEOPENERZ/completionist-api/internal/httpx"
	"github.com/GATEOPENERZ/completionist-api/internal/middleware"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/google/uuid"
)

type createAttachmentBody struct {
	Kind            string  `json:"kind"`
	Value           string  `json:"value"`
	StorageProvider *string `json:"storageProvider"`
	ContentType     *string `json:"contentType"`
	SizeBytes       *int64  `json:"sizeBytes"`
}

func (h *Handler) CreateAttachment(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.UserIDFromContext(r.Context()); !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var in createAttachmentBody
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	a := &models.Attachment{
		Kind:            in.Kind,
		Value:           in.Value,
		StorageProvider: in.StorageProvider,
		ContentType:     in.ContentType,
		SizeBytes:       in.SizeBytes,
	}
	if err := h.AttachmentRepo.Create(a); err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to create attachment")
		return
	}
	httpx.JSON(w, http.StatusCreated, a)
}

type linkBody struct {
	EntityTable  string    `json:"entityTable"`
	EntityPK     string    `json:"entityPk"`
	AttachmentID uuid.UUID `json:"attachmentId"`
}

func (h *Handler) LinkAttachment(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.UserIDFromContext(r.Context()); !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var in linkBody
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	ea, err := h.AttachmentRepo.Link(in.EntityTable, in.EntityPK, in.AttachmentID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to link attachment")
		return
	}
	httpx.JSON(w, http.StatusCreated, ea)
}

func (h *Handler) UnlinkAttachment(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.UserIDFromContext(r.Context()); !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var in linkBody
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	affected, err := h.AttachmentRepo.Unlink(in.EntityTable, in.EntityPK, in.AttachmentID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to unlink attachment")
		return
	}
	if affected == 0 {
		httpx.JSONError(w, http.StatusNotFound, "Link not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListAttachmentsByEntity(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.UserIDFromContext(r.Context()); !ok {
		httpx.JSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	table := r.URL.Query().Get("table")
	pk := r.URL.Query().Get("pk")
	if table == "" || pk == "" {
		httpx.JSONError(w, http.StatusBadRequest, "table and pk are required")
		return
	}
	items, err := h.AttachmentRepo.ListByEntity(table, pk)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, "Failed to list attachments")
		return
	}
	httpx.JSON(w, http.StatusOK, items)
}
