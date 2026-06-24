package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/nearestdev/completionist/internal/models"
)

type ConnectedAccountRepository struct {
	DB *sqlx.DB
}

func NewConnectedAccountRepository(db *sqlx.DB) *ConnectedAccountRepository {
	return &ConnectedAccountRepository{DB: db}
}

func (r *ConnectedAccountRepository) Create(a *models.ConnectedAccount) error {
	return r.DB.QueryRowx(`
		INSERT INTO connected_accounts (user_id, provider, provider_user_id, provider_username, access_token, refresh_token, token_expires_at, scopes, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, provider, provider_user_id, provider_username, access_token, refresh_token, token_expires_at, scopes, last_synced_at, metadata, created_at, updated_at
	`, a.UserID, a.Provider, a.ProviderUserID, a.ProviderUsername, a.AccessToken, a.RefreshToken, a.TokenExpiresAt, a.Scopes, a.Metadata).StructScan(a)
}

func (r *ConnectedAccountRepository) GetByUserID(userID int64) ([]models.ConnectedAccount, error) {
	var accounts []models.ConnectedAccount
	err := r.DB.Select(&accounts, `
		SELECT id, user_id, provider, provider_user_id, provider_username, last_synced_at, metadata, created_at, updated_at
		FROM connected_accounts WHERE user_id = $1 ORDER BY provider
	`, userID)
	return accounts, err
}

func (r *ConnectedAccountRepository) GetByProvider(userID int64, provider string) (*models.ConnectedAccount, error) {
	var a models.ConnectedAccount
	err := r.DB.QueryRowx(`
		SELECT id, user_id, provider, provider_user_id, provider_username, access_token, refresh_token, token_expires_at, scopes, last_synced_at, metadata, created_at, updated_at
		FROM connected_accounts WHERE user_id = $1 AND provider = $2
	`, userID, provider).StructScan(&a)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ConnectedAccountRepository) Update(a *models.ConnectedAccount) error {
	_, err := r.DB.Exec(`
		UPDATE connected_accounts SET
			provider_user_id = $2, provider_username = $3, access_token = $4, refresh_token = $5,
			token_expires_at = $6, scopes = $7, last_synced_at = $8, metadata = $9, updated_at = NOW()
		WHERE id = $1
	`, a.ID, a.ProviderUserID, a.ProviderUsername, a.AccessToken, a.RefreshToken, a.TokenExpiresAt, a.Scopes, a.LastSyncedAt, a.Metadata)
	return err
}

func (r *ConnectedAccountRepository) Delete(userID int64, provider string) error {
	_, err := r.DB.Exec(`DELETE FROM connected_accounts WHERE user_id = $1 AND provider = $2`, userID, provider)
	return err
}

func (r *ConnectedAccountRepository) CreateImportJob(job *models.ImportJob) error {
	return r.DB.QueryRowx(`
		INSERT INTO import_jobs (user_id, connected_account_id, provider, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, connected_account_id, provider, status, total_items, imported_items, skipped_items, failed_items, error_log, started_at, completed_at, created_at, updated_at
	`, job.UserID, job.ConnectedAccountID, job.Provider, models.ImportJobPending).StructScan(job)
}

func (r *ConnectedAccountRepository) UpdateImportJob(job *models.ImportJob) error {
	_, err := r.DB.Exec(`
		UPDATE import_jobs SET
			status = $2, total_items = $3, imported_items = $4, skipped_items = $5, failed_items = $6,
			error_log = $7, started_at = $8, completed_at = $9, updated_at = NOW()
		WHERE id = $1
	`, job.ID, job.Status, job.TotalItems, job.ImportedItems, job.SkippedItems, job.FailedItems, job.ErrorLog, job.StartedAt, job.CompletedAt)
	return err
}

func (r *ConnectedAccountRepository) GetImportJobs(userID int64) ([]models.ImportJob, error) {
	var jobs []models.ImportJob
	err := r.DB.Select(&jobs, `
		SELECT id, user_id, connected_account_id, provider, status, total_items, imported_items, skipped_items, failed_items, error_log, started_at, completed_at, created_at, updated_at
		FROM import_jobs WHERE user_id = $1 ORDER BY created_at DESC LIMIT 20
	`, userID)
	return jobs, err
}
