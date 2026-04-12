package repository

import (
	"database/sql"
	"time"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type StreakRepository struct {
	DB *sqlx.DB
}

func NewStreakRepository(db *sqlx.DB) *StreakRepository {
	return &StreakRepository{DB: db}
}

func (r *StreakRepository) GetByUserID(userID int64) (*models.UserStreak, error) {
	var s models.UserStreak
	err := r.DB.QueryRowx(`SELECT user_id, current_streak, longest_streak, last_activity_date, tier, updated_at FROM user_streaks WHERE user_id = $1`, userID).StructScan(&s)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StreakRepository) TouchStreak(userID int64, today time.Time) (*models.UserStreak, error) {
	date := today.Truncate(24 * time.Hour)
	var s models.UserStreak
	err := r.DB.QueryRowx(`
		INSERT INTO user_streaks (user_id, current_streak, longest_streak, last_activity_date, tier, updated_at)
		VALUES ($1, 1, 1, $2, 'spark', NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			current_streak = CASE
				WHEN user_streaks.last_activity_date = $2 THEN user_streaks.current_streak
				WHEN user_streaks.last_activity_date = $2::date - INTERVAL '1 day' THEN user_streaks.current_streak + 1
				ELSE 1
			END,
			longest_streak = GREATEST(
				user_streaks.longest_streak,
				CASE
					WHEN user_streaks.last_activity_date = $2 THEN user_streaks.current_streak
					WHEN user_streaks.last_activity_date = $2::date - INTERVAL '1 day' THEN user_streaks.current_streak + 1
					ELSE 1
				END
			),
			last_activity_date = $2,
			tier = CASE
				WHEN (CASE
					WHEN user_streaks.last_activity_date = $2 THEN user_streaks.current_streak
					WHEN user_streaks.last_activity_date = $2::date - INTERVAL '1 day' THEN user_streaks.current_streak + 1
					ELSE 1
				END) >= 100 THEN 'inferno'::streak_tier
				WHEN (CASE
					WHEN user_streaks.last_activity_date = $2 THEN user_streaks.current_streak
					WHEN user_streaks.last_activity_date = $2::date - INTERVAL '1 day' THEN user_streaks.current_streak + 1
					ELSE 1
				END) >= 10 THEN 'fire'::streak_tier
				WHEN (CASE
					WHEN user_streaks.last_activity_date = $2 THEN user_streaks.current_streak
					WHEN user_streaks.last_activity_date = $2::date - INTERVAL '1 day' THEN user_streaks.current_streak + 1
					ELSE 1
				END) >= 1 THEN 'spark'::streak_tier
				ELSE 'none'::streak_tier
			END,
			updated_at = NOW()
		RETURNING user_id, current_streak, longest_streak, last_activity_date, tier, updated_at
	`, userID, date).StructScan(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StreakRepository) GetLeaderboard(limit int) ([]models.UserStreak, error) {
	var streaks []models.UserStreak
	err := r.DB.Select(&streaks, `
		SELECT user_id, current_streak, longest_streak, last_activity_date, tier, updated_at
		FROM user_streaks
		ORDER BY current_streak DESC
		LIMIT $1
	`, limit)
	return streaks, err
}
