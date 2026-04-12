package models

import (
	"encoding/json"
	"time"
)

type BadgeCriteriaType string

const (
	BadgeCriteriaItemCompletionCount BadgeCriteriaType = "item_completion_count"
	BadgeCriteriaItemTypeCount       BadgeCriteriaType = "item_type_count"
	BadgeCriteriaStreakThreshold     BadgeCriteriaType = "streak_threshold"
	BadgeCriteriaFranchiseCompletion BadgeCriteriaType = "franchise_completion"
)

type Badge struct {
	ID           int64             `db:"id" json:"id"`
	Code         string            `db:"code" json:"code"`
	Name         string            `db:"name" json:"name"`
	Description  *string           `db:"description" json:"description,omitempty"`
	Icon         *string           `db:"icon" json:"icon,omitempty"`
	CriteriaType BadgeCriteriaType `db:"criteria_type" json:"criteriaType"`
	CriteriaJSON json.RawMessage   `db:"criteria_json" json:"criteriaJson" swaggertype:"string"`
	CreatedAt    time.Time         `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time         `db:"updated_at" json:"updatedAt"`
}

type UserBadge struct {
	ID       int64     `db:"id" json:"id"`
	UserID   int64     `db:"user_id" json:"userId"`
	BadgeID  int64     `db:"badge_id" json:"badgeId"`
	EarnedAt time.Time `db:"earned_at" json:"earnedAt"`
}

type BadgeWithEarned struct {
	ID           int64             `db:"id" json:"id"`
	Code         string            `db:"code" json:"code"`
	Name         string            `db:"name" json:"name"`
	Description  *string           `db:"description" json:"description,omitempty"`
	Icon         *string           `db:"icon" json:"icon,omitempty"`
	CriteriaType BadgeCriteriaType `db:"criteria_type" json:"criteriaType"`
	CriteriaJSON json.RawMessage   `db:"criteria_json" json:"criteriaJson" swaggertype:"string"`
	EarnedAt     *time.Time        `db:"earned_at" json:"earnedAt,omitempty"`
}
