package models

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	ID              uuid.UUID `db:"id" json:"id"`
	Kind            string    `db:"kind" json:"kind"`
	Value           string    `db:"value" json:"value"`
	StorageProvider *string   `db:"storage_provider" json:"storageProvider,omitempty"`
	ContentType     *string   `db:"content_type" json:"contentType,omitempty"`
	SizeBytes       *int64    `db:"size_bytes" json:"sizeBytes,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"`
}

type EntityAttachment struct {
	ID            int64     `db:"id" json:"id"`
	EntityTable   string    `db:"entity_table" json:"entityTable"`
	EntityPK      string    `db:"entity_pk" json:"entityPk"`
	AttachmentID  uuid.UUID `db:"attachment_id" json:"attachmentId"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
}
