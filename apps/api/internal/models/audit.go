package models

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID         int64           `db:"id" json:"id"`
	UserID     *int64          `db:"user_id" json:"user_id,omitempty"`
	EntityType string          `db:"entity_type" json:"entity_type"`
	EntityID   string          `db:"entity_id" json:"entity_id"`
	Action     string          `db:"action" json:"action"`
	Metadata   json.RawMessage `db:"metadata" json:"metadata"`
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
}
