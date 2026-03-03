package models

import (
	"encoding/json"
	"time"
)

type ChallengeType string
type ChallengeFrequency string
type CriteriaType string

const (
	ChallengeTypeGlobal   ChallengeType = "global"
	ChallengeTypePersonal ChallengeType = "personal"

	FrequencyDaily    ChallengeFrequency = "daily"
	FrequencyWeekly   ChallengeFrequency = "weekly"
	FrequencyMonthly  ChallengeFrequency = "monthly"
	FrequencySeasonal ChallengeFrequency = "seasonal"
	FrequencyInfinite ChallengeFrequency = "infinite"

	CriteriaCountItems   CriteriaType = "count_items"
	CriteriaSpecificItem CriteriaType = "specific_item"
	CriteriaGenreCount   CriteriaType = "genre_count"
)

type Season struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	StartAt   time.Time `db:"start_at" json:"startAt"`
	EndAt     time.Time `db:"end_at" json:"endAt"`
	IsActive  bool      `db:"is_active" json:"isActive"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type Challenge struct {
	ID               int64              `db:"id" json:"id"`
	Title            string             `db:"title" json:"title"`
	Description      *string            `db:"description" json:"description,omitempty"`
	Type             ChallengeType      `db:"type" json:"type"`
	Frequency        ChallengeFrequency `db:"frequency" json:"frequency"`
	XPReward         int                `db:"xp_reward" json:"xpReward"`
	CriteriaType     CriteriaType       `db:"criteria_type" json:"criteriaType"`
	CriteriaMetadata json.RawMessage    `db:"criteria_metadata" json:"criteriaMetadata" swaggertype:"string"`
	CreatedAt        time.Time          `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time          `db:"updated_at" json:"updatedAt"`
}

type UserChallenge struct {
	ID              int64     `db:"id" json:"id"`
	UserID          int64     `db:"user_id" json:"userId"`
	ChallengeID     int64     `db:"challenge_id" json:"challengeId"`
	SeasonID        *int64    `db:"season_id" json:"seasonId,omitempty"`
	CurrentProgress int       `db:"current_progress" json:"currentProgress"`
	TargetProgress  int       `db:"target_progress" json:"targetProgress"`
	IsCompleted     bool      `db:"is_completed" json:"isCompleted"`
	CompletedAt     *time.Time `db:"completed_at" json:"completedAt,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"`

	ChallengeTitle       *string `db:"challenge_title" json:"challengeTitle,omitempty"`
	ChallengeDescription *string `db:"challenge_description" json:"challengeDescription,omitempty"`
	XPReward             *int    `db:"xp_reward" json:"xpReward,omitempty"`
}

type CriteriaCountMetadata struct {
	Target   int    `json:"target"`
	ItemType string `json:"itemType,omitempty"`
}

type CriteriaSpecificItemMetadata struct {
	MediaItemID string `json:"mediaItemId"`
	Action      string `json:"action"`
}

type CriteriaGenreCountMetadata struct {
	Genre  string `json:"genre"`
	Target int    `json:"target"`
}