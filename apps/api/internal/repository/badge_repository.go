package repository

import (
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type BadgeRepository struct {
	DB *sqlx.DB
}

func NewBadgeRepository(db *sqlx.DB) *BadgeRepository {
	return &BadgeRepository{DB: db}
}

func (r *BadgeRepository) GetAll() ([]models.Badge, error) {
	var badges []models.Badge
	err := r.DB.Select(&badges, `SELECT id, code, name, description, icon, criteria_type, criteria_json, created_at, updated_at FROM badges ORDER BY id`)
	return badges, err
}

func (r *BadgeRepository) GetUserBadges(userID int64) ([]models.BadgeWithEarned, error) {
	var badges []models.BadgeWithEarned
	err := r.DB.Select(&badges, `
		SELECT b.id, b.code, b.name, b.description, b.icon, b.criteria_type, b.criteria_json, ub.earned_at
		FROM badges b
		LEFT JOIN user_badges ub ON ub.badge_id = b.id AND ub.user_id = $1
		ORDER BY b.id
	`, userID)
	return badges, err
}

func (r *BadgeRepository) AwardBadge(userID, badgeID int64) error {
	_, err := r.DB.Exec(`
		INSERT INTO user_badges (user_id, badge_id) VALUES ($1, $2)
		ON CONFLICT (user_id, badge_id) DO NOTHING
	`, userID, badgeID)
	return err
}

func (r *BadgeRepository) HasBadge(userID, badgeID int64) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM user_badges WHERE user_id = $1 AND badge_id = $2)`, userID, badgeID).Scan(&exists)
	return exists, err
}

func (r *BadgeRepository) GetCompletionCountByType(userID int64, itemType *string) (int, error) {
	var count int
	if itemType != nil {
		err := r.DB.QueryRow(`
			SELECT COUNT(*) FROM user_list_items uli
			JOIN media_items mi ON mi.id = uli.media_item_id
			WHERE uli.user_id = $1 AND uli.status = 'completed' AND mi.item_type = $2
		`, userID, *itemType).Scan(&count)
		return count, err
	}
	err := r.DB.QueryRow(`SELECT COUNT(*) FROM user_list_items WHERE user_id = $1 AND status = 'completed'`, userID).Scan(&count)
	return count, err
}
