package repository

import (
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AttachmentRepository struct {
	DB *sqlx.DB
}

func NewAttachmentRepository(db *sqlx.DB) *AttachmentRepository {
	return &AttachmentRepository{DB: db}
}

func (r *AttachmentRepository) Create(a *models.Attachment) error {
	return r.DB.QueryRowx(`
		INSERT INTO attachments (kind, value, storage_provider, content_type, size_bytes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, kind, value, storage_provider, content_type, size_bytes, created_at, updated_at
	`, a.Kind, a.Value, a.StorageProvider, a.ContentType, a.SizeBytes).StructScan(a)
}

func (r *AttachmentRepository) Link(entityTable, entityPk string, attachmentID uuid.UUID) (*models.EntityAttachment, error) {
	var ea models.EntityAttachment
	err := r.DB.QueryRowx(`
		INSERT INTO entity_attachments (entity_table, entity_pk, attachment_id)
		VALUES ($1, $2, $3)
		RETURNING id, entity_table, entity_pk, attachment_id, created_at
	`, entityTable, entityPk, attachmentID).StructScan(&ea)
	return &ea, err
}

func (r *AttachmentRepository) Unlink(entityTable, entityPk string, attachmentID uuid.UUID) (int64, error) {
	res, err := r.DB.Exec(`
		DELETE FROM entity_attachments
		WHERE entity_table = $1 AND entity_pk = $2 AND attachment_id = $3
	`, entityTable, entityPk, attachmentID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *AttachmentRepository) ListByEntity(entityTable, entityPk string) ([]models.Attachment, error) {
	out := []models.Attachment{}
	err := r.DB.Select(&out, `
		SELECT a.id, a.kind, a.value, a.storage_provider, a.content_type, a.size_bytes, a.created_at, a.updated_at
		FROM entity_attachments ea
		JOIN attachments a ON a.id = ea.attachment_id
		WHERE ea.entity_table = $1 AND ea.entity_pk = $2
		ORDER BY a.created_at DESC
	`, entityTable, entityPk)
	return out, err
}
