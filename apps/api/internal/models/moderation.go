package models

import (
	"encoding/json"
	"time"
)

type ModerationStatus string

const (
	ModerationStatusPending  ModerationStatus = "pending"
	ModerationStatusApproved ModerationStatus = "approved"
	ModerationStatusRejected ModerationStatus = "rejected"
)

type ModerationEntityType string

const (
	ModerationEntityPost          ModerationEntityType = "post"
	ModerationEntityComment       ModerationEntityType = "comment"
	ModerationEntityRoomMessage   ModerationEntityType = "room_message"
	ModerationEntityDirectMessage ModerationEntityType = "direct_message"
)

type ModerationItem struct {
	ID           int64                `db:"id" json:"id"`
	EntityType   ModerationEntityType `db:"entity_type" json:"entityType"`
	EntityID     int64                `db:"entity_id" json:"entityId"`
	ReportedBy   *int64               `db:"reported_by" json:"reportedBy,omitempty"`
	Reason       *string              `db:"reason" json:"reason,omitempty"`
	AIFlagged    bool                 `db:"ai_flagged" json:"aiFlagged"`
	AICategories json.RawMessage      `db:"ai_categories" json:"aiCategories" swaggertype:"string"`
	AIScores     json.RawMessage      `db:"ai_scores" json:"aiScores" swaggertype:"string"`
	Status       ModerationStatus     `db:"status" json:"status"`
	ReviewedBy   *int64               `db:"reviewed_by" json:"reviewedBy,omitempty"`
	ReviewedAt   *time.Time           `db:"reviewed_at" json:"reviewedAt,omitempty"`
	ReviewNote   *string              `db:"review_note" json:"reviewNote,omitempty"`
	CreatedAt    time.Time            `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time            `db:"updated_at" json:"updatedAt"`
}

type BanAppealStatus string

const (
	BanAppealPending  BanAppealStatus = "pending"
	BanAppealAccepted BanAppealStatus = "accepted"
	BanAppealDenied   BanAppealStatus = "denied"
)

type BanAppeal struct {
	ID         int64           `db:"id" json:"id"`
	UserID     int64           `db:"user_id" json:"userId"`
	AppealText string          `db:"appeal_text" json:"appealText"`
	Status     BanAppealStatus `db:"status" json:"status"`
	ReviewedBy *int64          `db:"reviewed_by" json:"reviewedBy,omitempty"`
	ReviewedAt *time.Time      `db:"reviewed_at" json:"reviewedAt,omitempty"`
	ReviewNote *string         `db:"review_note" json:"reviewNote,omitempty"`
	CreatedAt  time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time       `db:"updated_at" json:"updatedAt"`
}

type DeletionRequestStatus string

const (
	DeletionPending    DeletionRequestStatus = "pending"
	DeletionProcessing DeletionRequestStatus = "processing"
	DeletionCompleted  DeletionRequestStatus = "completed"
	DeletionCancelled  DeletionRequestStatus = "cancelled"
)

type AccountDeletionRequest struct {
	ID          int64                 `db:"id" json:"id"`
	UserID      int64                 `db:"user_id" json:"userId"`
	Reason      *string               `db:"reason" json:"reason,omitempty"`
	Status      DeletionRequestStatus `db:"status" json:"status"`
	ScheduledAt *time.Time            `db:"scheduled_at" json:"scheduledAt,omitempty"`
	CompletedAt *time.Time            `db:"completed_at" json:"completedAt,omitempty"`
	CreatedAt   time.Time             `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time             `db:"updated_at" json:"updatedAt"`
}

type NewReport struct {
	EntityType ModerationEntityType `json:"entityType"`
	EntityID   int64                `json:"entityId"`
	Reason     string               `json:"reason"`
}
