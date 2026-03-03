package repository

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type AuditRepository struct {
	DB *sqlx.DB
}

func NewAuditRepository(db *sqlx.DB) *AuditRepository {
	return &AuditRepository{DB: db}
}

func (r *AuditRepository) Log(userID *int64, entityType, entityID, action string, metadata map[string]interface{}) error {
	mdJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(`
		INSERT INTO audit_logs (user_id, entity_type, entity_id, action, metadata)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, entityType, entityID, action, mdJSON)
	return err
}
