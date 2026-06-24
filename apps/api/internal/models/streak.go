package models

import "time"

type StreakTier string

const (
	StreakTierNone    StreakTier = "none"
	StreakTierSpark   StreakTier = "spark"
	StreakTierFire    StreakTier = "fire"
	StreakTierInferno StreakTier = "inferno"
)

func StreakTierFromCount(count int) StreakTier {
	switch {
	case count >= 100:
		return StreakTierInferno
	case count >= 10:
		return StreakTierFire
	case count >= 1:
		return StreakTierSpark
	default:
		return StreakTierNone
	}
}

type UserStreak struct {
	UserID           int64      `db:"user_id" json:"userId"`
	CurrentStreak    int        `db:"current_streak" json:"currentStreak"`
	LongestStreak    int        `db:"longest_streak" json:"longestStreak"`
	LastActivityDate *time.Time `db:"last_activity_date" json:"lastActivityDate,omitempty"`
	Tier             StreakTier `db:"tier" json:"tier"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updatedAt"`
}
