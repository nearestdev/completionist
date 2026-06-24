package repository

import (
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type LastFMRepository struct {
	DB *sqlx.DB
}

func NewLastFMRepository(db *sqlx.DB) *LastFMRepository {
	return &LastFMRepository{DB: db}
}

func (r *LastFMRepository) Upsert(acc *models.LastFMAccount) error {
	query := `
        INSERT INTO lastfm_accounts (user_id, username, session_key, subscriber)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id) DO UPDATE SET
            username = EXCLUDED.username,
            session_key = EXCLUDED.session_key,
            subscriber = EXCLUDED.subscriber;
    `
	_, err := r.DB.Exec(query, acc.UserID, acc.Username, acc.SessionKey, acc.Subscriber)
	return err
}

func (r *LastFMRepository) FindByUserID(userID int64) (*models.LastFMAccount, error) {
	var acc models.LastFMAccount
	query := `SELECT user_id, username, session_key, subscriber FROM lastfm_accounts WHERE user_id = $1`
	err := r.DB.Get(&acc, query, userID)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}
