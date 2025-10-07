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
		RETURNING id, username, email, password_hash, created_at, updated_at`
	return r.DB.QueryRowx(query, user.Username, user.Email, user.PasswordHash).StructScan(user)
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = $1`
	if err := r.DB.Get(&user, query, username); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	if err := r.DB.Get(&user, query, email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
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