package repository

import (
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type SteamRepository struct {
	DB *sqlx.DB
}

func NewSteamRepository(db *sqlx.DB) *SteamRepository {
	return &SteamRepository{DB: db}
}
func (r *SteamRepository) UpsertSteamAccount(a *models.SteamAccount) error {
	return r.DB.QueryRowx(`
		INSERT INTO steam_accounts (user_id, steam_id, persona, avatar)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET steam_id = EXCLUDED.steam_id, persona = EXCLUDED.persona, avatar = EXCLUDED.avatar
		RETURNING user_id, steam_id, persona, avatar
	`, a.UserID, a.SteamID, a.Persona, a.Avatar).StructScan(a)
}
func (r *SteamRepository) FindByUserID(userID int64) (*models.SteamAccount, error) {
	var a models.SteamAccount
	if err := r.DB.Get(&a, `SELECT user_id, steam_id, persona, avatar FROM steam_accounts WHERE user_id = $1`, userID); err != nil {
		return nil, err
	}
	return &a, nil
}
