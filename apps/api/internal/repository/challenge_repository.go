package repository

import (
	"time"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type ChallengeRepository struct {
	DB *sqlx.DB
}

func NewChallengeRepository(db *sqlx.DB) *ChallengeRepository {
	return &ChallengeRepository{DB: db}
}

func (r *ChallengeRepository) GetActiveSeason() (*models.Season, error) {
	var s models.Season
	err := r.DB.Get(&s, `
		SELECT id, name, start_at, end_at, is_active, created_at
		FROM seasons
		WHERE is_active = true
		  AND start_at <= NOW()
		  AND end_at >= NOW()
		LIMIT 1
	`)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ChallengeRepository) CreateSeason(s *models.Season) error {
	return r.DB.QueryRowx(`
		INSERT INTO seasons (name, start_at, end_at, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, start_at, end_at, is_active, created_at
	`, s.Name, s.StartAt, s.EndAt, s.IsActive).StructScan(s)
}

func (r *ChallengeRepository) CreateChallenge(c *models.Challenge) error {
	return r.DB.QueryRowx(`
		INSERT INTO challenges (title, description, type, frequency, xp_reward, criteria_type, criteria_metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, description, type, frequency, xp_reward, criteria_type, criteria_metadata, created_at, updated_at
	`, c.Title, c.Description, c.Type, c.Frequency, c.XPReward, c.CriteriaType, c.CriteriaMetadata).StructScan(c)
}

func (r *ChallengeRepository) GetChallengesByFrequency(freq models.ChallengeFrequency) ([]models.Challenge, error) {
	var challenges []models.Challenge
	err := r.DB.Select(&challenges, `
		SELECT id, title, description, type, frequency, xp_reward, criteria_type, criteria_metadata, created_at, updated_at
		FROM challenges
		WHERE frequency = $1
	`, freq)
	return challenges, err
}

func (r *ChallengeRepository) GetUserChallenges(userID int64, seasonID *int64) ([]models.UserChallenge, error) {
	var ucs []models.UserChallenge
	query := `
		SELECT 
			uc.id, uc.user_id, uc.challenge_id, uc.season_id, uc.current_progress, uc.target_progress, 
			uc.is_completed, uc.completed_at, uc.created_at, uc.updated_at,
			c.title AS challenge_title,
			c.description AS challenge_description,
			c.xp_reward
		FROM user_challenges uc
		JOIN challenges c ON uc.challenge_id = c.id
		WHERE uc.user_id = $1
	`
	args := []interface{}{userID}

	if seasonID != nil {
		query += " AND uc.season_id = $2"
		args = append(args, *seasonID)
	}

	err := r.DB.Select(&ucs, query, args...)
	return ucs, err
}

func (r *ChallengeRepository) FindUserChallenge(userID, challengeID int64, seasonID *int64) (*models.UserChallenge, error) {
	var uc models.UserChallenge
	query := `
		SELECT id, user_id, challenge_id, season_id, current_progress, target_progress, is_completed, completed_at, created_at, updated_at
		FROM user_challenges
		WHERE user_id = $1 AND challenge_id = $2
	`
	args := []interface{}{userID, challengeID}

	if seasonID != nil {
		query += " AND season_id = $3"
		args = append(args, *seasonID)
	} else {
		query += " AND season_id IS NULL"
	}

	err := r.DB.Get(&uc, query, args...)
	if err != nil {
		return nil, err
	}
	return &uc, nil
}

func (r *ChallengeRepository) UpsertUserChallenge(uc *models.UserChallenge) error {
	return r.DB.QueryRowx(`
		INSERT INTO user_challenges (user_id, challenge_id, season_id, current_progress, target_progress, is_completed, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, challenge_id, season_id) 
		DO UPDATE SET 
			current_progress = EXCLUDED.current_progress,
			target_progress = EXCLUDED.target_progress,
			is_completed = EXCLUDED.is_completed,
			completed_at = EXCLUDED.completed_at,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`, uc.UserID, uc.ChallengeID, uc.SeasonID, uc.CurrentProgress, uc.TargetProgress, uc.IsCompleted, uc.CompletedAt).
		Scan(&uc.ID, &uc.CreatedAt, &uc.UpdatedAt)
}

type LeaderboardEntry struct {
	UserID   int64  `db:"user_id" json:"userId"`
	Username string `db:"username" json:"username"`
	TotalXP  int64  `db:"total_xp" json:"totalXp"`
	Rank     int    `db:"rank" json:"rank"`
}

func (r *ChallengeRepository) GetSeasonLeaderboard(startAt, endAt time.Time, limit, offset int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	err := r.DB.Select(&entries, `
		SELECT 
			u.id AS user_id,
			u.username,
			COALESCE(SUM(h.amount), 0) AS total_xp,
			RANK() OVER (ORDER BY COALESCE(SUM(h.amount), 0) DESC) as rank
		FROM users u
		LEFT JOIN user_xp_history h ON h.user_id = u.id AND h.created_at BETWEEN $1 AND $2
		GROUP BY u.id
		ORDER BY total_xp DESC
		LIMIT $3 OFFSET $4
	`, startAt, endAt, limit, offset)
	return entries, err
}

func (r *ChallengeRepository) GetGlobalLeaderboard(limit, offset int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	err := r.DB.Select(&entries, `
		SELECT 
			id AS user_id,
			username,
			xp AS total_xp,
			RANK() OVER (ORDER BY xp DESC) as rank
		FROM users
		ORDER BY xp DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	return entries, err
}
