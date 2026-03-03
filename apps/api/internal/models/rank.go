package models

import (
	"time"
)

type Rank int

const (
	RankScribe     Rank = 1
	RankChronicler Rank = 2
	RankCurator    Rank = 3
	RankPreserver  Rank = 4
	RankWarden     Rank = 5
	RankOracle     Rank = 6
)

func (r Rank) String() string {
	switch r {
	case RankScribe:
		return "Scribe"
	case RankChronicler:
		return "Chronicler"
	case RankCurator:
		return "Curator"
	case RankPreserver:
		return "Preserver"
	case RankWarden:
		return "Warden"
	case RankOracle:
		return "Oracle"
	default:
		return "Unknown"
	}
}

type UserRank struct {
	ID          int64     `db:"id" json:"id"`
	UserID      int64     `db:"user_id" json:"userId"`
	SeasonID    *int64    `db:"season_id" json:"seasonId,omitempty"` 
	CurrentRank Rank      `db:"current_rank" json:"currentRank"`
	CurrentElo  int       `db:"current_elo" json:"currentElo"`
	PeakRank    Rank      `db:"peak_rank" json:"peakRank"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}
