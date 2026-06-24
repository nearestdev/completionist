package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nearestdev/completionist/internal/models"
)

type RankService struct {
	db *sqlx.DB
}

func NewRankService(db *sqlx.DB) *RankService {
	return &RankService{db: db}
}

func (s *RankService) GetUserRank(userID int64, seasonID int64) (*models.UserRank, error) {
	var rank models.UserRank
	err := s.db.Get(&rank, "SELECT * FROM user_ranks WHERE user_id = $1 AND season_id = $2", userID, seasonID)
	if err == sql.ErrNoRows {
		return &models.UserRank{
			UserID:      userID,
			SeasonID:    &seasonID,
			CurrentRank: models.RankScribe,
			CurrentElo:  0,
			PeakRank:    models.RankScribe,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &rank, nil
}

func (s *RankService) UpdateElo(userID int64, seasonID int64, eloChange int) (*models.UserRank, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var rank models.UserRank
	err = tx.Get(&rank, "SELECT * FROM user_ranks WHERE user_id = $1 AND season_id = $2 FOR UPDATE", userID, seasonID)
	if err == sql.ErrNoRows {

		rank = models.UserRank{
			UserID:      userID,
			SeasonID:    &seasonID,
			CurrentRank: models.RankScribe,
			CurrentElo:  0,
			PeakRank:    models.RankScribe,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

	} else if err != nil {
		return nil, err
	}

	rank.CurrentElo += eloChange
	if rank.CurrentElo < 0 {
		rank.CurrentElo = 0
	}

	newRank := CalculateRank(rank.CurrentElo)
	rank.CurrentRank = newRank

	if newRank > rank.PeakRank {
		rank.PeakRank = newRank
	}

	rank.UpdatedAt = time.Now()

	query := `
		INSERT INTO user_ranks (user_id, season_id, current_rank, current_elo, peak_rank, created_at, updated_at)
		VALUES (:user_id, :season_id, :current_rank, :current_elo, :peak_rank, :created_at, :updated_at)
		ON CONFLICT (user_id, season_id) 
		DO UPDATE SET 
			current_rank = :current_rank,
			current_elo = :current_elo,
			peak_rank = :peak_rank,
			updated_at = :updated_at
	`
	_, err = tx.NamedExec(query, &rank)
	if err != nil {
		return nil, fmt.Errorf("failed to save user rank: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &rank, nil
}

func CalculateRank(elo int) models.Rank {
	switch {
	case elo >= 2500:
		return models.RankOracle
	case elo >= 1500:
		return models.RankWarden
	case elo >= 800:
		return models.RankPreserver
	case elo >= 400:
		return models.RankCurator
	case elo >= 150:
		return models.RankChronicler
	default:
		return models.RankScribe
	}
}
