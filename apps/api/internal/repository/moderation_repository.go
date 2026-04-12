package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type ModerationRepository struct {
	DB *sqlx.DB
}

func NewModerationRepository(db *sqlx.DB) *ModerationRepository {
	return &ModerationRepository{DB: db}
}

func (r *ModerationRepository) CreateItem(item *models.ModerationItem) error {
	return r.DB.QueryRowx(`
		INSERT INTO moderation_queue (entity_type, entity_id, reported_by, reason, ai_flagged, ai_categories, ai_scores)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, entity_type, entity_id, reported_by, reason, ai_flagged, ai_categories, ai_scores, status, reviewed_by, reviewed_at, review_note, created_at, updated_at
	`, item.EntityType, item.EntityID, item.ReportedBy, item.Reason, item.AIFlagged, item.AICategories, item.AIScores).StructScan(item)
}

func (r *ModerationRepository) GetQueue(status string, limit, offset int) ([]models.ModerationItem, error) {
	var items []models.ModerationItem
	query := `SELECT * FROM moderation_queue`
	args := []interface{}{}
	if status != "" {
		query += ` WHERE status = $1`
		args = append(args, status)
		query += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = append(args, limit, offset)
	} else {
		query += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		args = append(args, limit, offset)
	}
	err := r.DB.Select(&items, query, args...)
	return items, err
}

func (r *ModerationRepository) UpdateStatus(id int64, status models.ModerationStatus, reviewedBy int64, reviewNote *string) error {
	now := time.Now()
	_, err := r.DB.Exec(`
		UPDATE moderation_queue SET status = $2, reviewed_by = $3, reviewed_at = $4, review_note = $5, updated_at = $4
		WHERE id = $1
	`, id, status, reviewedBy, now, reviewNote)
	return err
}

func (r *ModerationRepository) BanUser(userID int64, reason string) error {
	_, err := r.DB.Exec(`UPDATE users SET banned_at = NOW(), ban_reason = $2, updated_at = NOW() WHERE id = $1`, userID, reason)
	return err
}

func (r *ModerationRepository) UnbanUser(userID int64) error {
	_, err := r.DB.Exec(`UPDATE users SET banned_at = NULL, ban_reason = NULL, updated_at = NOW() WHERE id = $1`, userID)
	return err
}

func (r *ModerationRepository) CreateBanAppeal(appeal *models.BanAppeal) error {
	return r.DB.QueryRowx(`
		INSERT INTO ban_appeals (user_id, appeal_text) VALUES ($1, $2)
		RETURNING id, user_id, appeal_text, status, reviewed_by, reviewed_at, review_note, created_at, updated_at
	`, appeal.UserID, appeal.AppealText).StructScan(appeal)
}

func (r *ModerationRepository) GetBanAppeals(status string, limit, offset int) ([]models.BanAppeal, error) {
	var appeals []models.BanAppeal
	err := r.DB.Select(&appeals, `
		SELECT id, user_id, appeal_text, status, reviewed_by, reviewed_at, review_note, created_at, updated_at
		FROM ban_appeals WHERE ($1 = '' OR status = $1::ban_appeal_status)
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, status, limit, offset)
	return appeals, err
}

func (r *ModerationRepository) UpdateBanAppeal(id int64, status models.BanAppealStatus, reviewedBy int64, reviewNote *string) error {
	now := time.Now()
	_, err := r.DB.Exec(`
		UPDATE ban_appeals SET status = $2, reviewed_by = $3, reviewed_at = $4, review_note = $5, updated_at = $4
		WHERE id = $1
	`, id, status, reviewedBy, now, reviewNote)
	return err
}

func (r *ModerationRepository) CreateDeletionRequest(userID int64, reason *string) (*models.AccountDeletionRequest, error) {
	scheduledAt := time.Now().AddDate(0, 0, 30)
	var req models.AccountDeletionRequest
	err := r.DB.QueryRowx(`
		INSERT INTO account_deletion_requests (user_id, reason, scheduled_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET reason = $2, status = 'pending', scheduled_at = $3, updated_at = NOW()
		RETURNING id, user_id, reason, status, scheduled_at, completed_at, created_at, updated_at
	`, userID, reason, scheduledAt).StructScan(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *ModerationRepository) CancelDeletionRequest(userID int64) error {
	_, err := r.DB.Exec(`
		UPDATE account_deletion_requests SET status = 'cancelled', updated_at = NOW()
		WHERE user_id = $1 AND status = 'pending'
	`, userID)
	return err
}

func (r *ModerationRepository) GetDeletionRequest(userID int64) (*models.AccountDeletionRequest, error) {
	var req models.AccountDeletionRequest
	err := r.DB.QueryRowx(`
		SELECT id, user_id, reason, status, scheduled_at, completed_at, created_at, updated_at
		FROM account_deletion_requests WHERE user_id = $1 AND status = 'pending'
	`, userID).StructScan(&req)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *ModerationRepository) ExportUserData(userID int64) (json.RawMessage, error) {
	type exportData struct {
		User         interface{} `json:"user"`
		ListItems    interface{} `json:"listItems"`
		Posts        interface{} `json:"posts"`
		Messages     interface{} `json:"messages"`
		Reviews      interface{} `json:"reviews"`
		Badges       interface{} `json:"badges"`
		Streaks      interface{} `json:"streaks"`
		Follows      interface{} `json:"follows"`
		XPHistory    interface{} `json:"xpHistory"`
	}

	var user interface{}
	r.DB.Get(&user, `SELECT id, username, email, role, xp, level, created_at FROM users WHERE id = $1`, userID)

	data := exportData{User: user}
	return json.Marshal(data)
}
