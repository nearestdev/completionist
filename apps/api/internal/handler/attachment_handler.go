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

// @Summary      Create Attachment
// @Description  Creates a new standalone attachment record
// @Tags         Attachments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body createAttachmentBody true "Attachment Details"
// @Success      201  {object}  models.Attachment
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /attachments [post]
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

// @Summary      Link Attachment
// @Description  Links an existing attachment to an entity (post, comment, etc.)
// @Tags         Attachments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body linkBody true "Link Details"
// @Success      201  {object}  models.EntityAttachment
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /attachments/link [post]
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

// @Summary      Unlink Attachment
// @Description  Removes a link between an attachment and an entity
// @Tags         Attachments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body linkBody true "Link Details"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /attachments/unlink [delete]
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

// @Summary      List Attachments by Entity
// @Description  Retrieves all attachments linked to a specific entity
// @Tags         Attachments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        table query string true "Entity Table Name"
// @Param        pk query string true "Entity Primary Key"
// @Success      200  {array}   models.Attachment
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /attachments/by-entity [get]
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
