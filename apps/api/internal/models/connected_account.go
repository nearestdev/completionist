package models

import (
	"encoding/json"
	"time"
)

type ImportJobStatus string

const (
	ImportJobPending   ImportJobStatus = "pending"
	ImportJobRunning   ImportJobStatus = "running"
	ImportJobCompleted ImportJobStatus = "completed"
	ImportJobFailed    ImportJobStatus = "failed"
	ImportJobCancelled ImportJobStatus = "cancelled"
)

type ConnectedAccount struct {
	ID               int64            `db:"id" json:"id"`
	UserID           int64            `db:"user_id" json:"userId"`
	Provider         string           `db:"provider" json:"provider"`
	ProviderUserID   *string          `db:"provider_user_id" json:"providerUserId,omitempty"`
	ProviderUsername *string          `db:"provider_username" json:"providerUsername,omitempty"`
	AccessToken      *string          `db:"access_token" json:"-"`
	RefreshToken     *string          `db:"refresh_token" json:"-"`
	TokenExpiresAt   *time.Time       `db:"token_expires_at" json:"-"`
	Scopes           *string          `db:"scopes" json:"-"`
	LastSyncedAt     *time.Time       `db:"last_synced_at" json:"lastSyncedAt,omitempty"`
	Metadata         *json.RawMessage `db:"metadata" json:"metadata,omitempty" swaggertype:"string"`
	CreatedAt        time.Time        `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time        `db:"updated_at" json:"updatedAt"`
}

type ImportJob struct {
	ID                 int64           `db:"id" json:"id"`
	UserID             int64           `db:"user_id" json:"userId"`
	ConnectedAccountID *int64          `db:"connected_account_id" json:"connectedAccountId,omitempty"`
	Provider           string          `db:"provider" json:"provider"`
	Status             ImportJobStatus `db:"status" json:"status"`
	TotalItems         int             `db:"total_items" json:"totalItems"`
	ImportedItems      int             `db:"imported_items" json:"importedItems"`
	SkippedItems       int             `db:"skipped_items" json:"skippedItems"`
	FailedItems        int             `db:"failed_items" json:"failedItems"`
	ErrorLog           json.RawMessage `db:"error_log" json:"errorLog" swaggertype:"string"`
	StartedAt          *time.Time      `db:"started_at" json:"startedAt,omitempty"`
	CompletedAt        *time.Time      `db:"completed_at" json:"completedAt,omitempty"`
	CreatedAt          time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt          time.Time       `db:"updated_at" json:"updatedAt"`
}
