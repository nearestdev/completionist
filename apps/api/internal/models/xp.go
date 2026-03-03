package models

import "time"

type UserXPHistory struct {
	ID          int64     `db:"id" json:"id"`
	UserID      int64     `db:"user_id" json:"userId"`
	Amount      int       `db:"amount" json:"amount"`
	SourceType  string    `db:"source_type" json:"sourceType"`
	SourceID    *string   `db:"source_id" json:"sourceId,omitempty"`
	Description *string   `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

const (
	XPSourceListCompletion = "list_completion"
	XPSourceReview         = "review"
	XPSourceComment        = "comment"
)