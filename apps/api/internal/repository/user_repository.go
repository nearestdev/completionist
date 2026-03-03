package repository

import (
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)
		RETURNING id, username, email, password_hash, xp, level, created_at, updated_at`
	return r.DB.QueryRowx(query, user.Username, user.Email, user.PasswordHash).StructScan(user)
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, xp, level, created_at, updated_at FROM users WHERE username = $1`
	if err := r.DB.Get(&user, query, username); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, xp, level, created_at, updated_at FROM users WHERE email = $1`
	if err := r.DB.Get(&user, query, email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, xp, level, created_at, updated_at FROM users WHERE id = $1`
	if err := r.DB.Get(&user, query, id); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SearchByUsername(query string, limit int) ([]models.UserFollowResponse, error) {
	var users []models.UserFollowResponse
	q := `
		SELECT id AS user_id, username, created_at AS followed_at 
		FROM users 
		WHERE username ILIKE $1 
		ORDER BY created_at DESC 
		LIMIT $2`
	if err := r.DB.Select(&users, q, "%"+query+"%", limit); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) AddXP(userID int64, sourceType string, sourceID *string, description *string) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var amount int
	err = tx.Get(&amount, "SELECT xp_amount FROM xp_actions WHERE action_key = $1", sourceType)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO user_xp_history (user_id, amount, source_type, source_id, description)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, amount, sourceType, sourceID, description)
	if err != nil {
		return err
	}

	var newXP int64
	err = tx.QueryRow(`
		UPDATE users 
		SET xp = xp + $1, updated_at = NOW()
		WHERE id = $2
		RETURNING xp
	`, amount, userID).Scan(&newXP)
	if err != nil {
		return err
	}

	var newLevel int
	err = tx.Get(&newLevel, `
		SELECT level 
		FROM level_definitions 
		WHERE xp_required <= $1 
		ORDER BY level DESC 
		LIMIT 1
	`, newXP)
	if err != nil {
		newLevel = 1
	}

	_, err = tx.Exec(`UPDATE users SET level = $1 WHERE id = $2`, newLevel, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) GetXPHistory(userID int64, limit, offset int) ([]models.UserXPHistory, error) {
	var history []models.UserXPHistory
	query := `
		SELECT id, user_id, amount, source_type, source_id, description, created_at 
		FROM user_xp_history 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`
	err := r.DB.Select(&history, query, userID, limit, offset)
	return history, err
}