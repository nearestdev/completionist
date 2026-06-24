package repository

import (
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type SocialRepository struct {
	DB *sqlx.DB
}

func NewSocialRepository(db *sqlx.DB) *SocialRepository {
	return &SocialRepository{DB: db}
}

func (r *SocialRepository) CreateUserFollow(followerID, followedID int64) (*models.UserFollow, error) {
	var follow models.UserFollow
	err := r.DB.QueryRowx(`
		INSERT INTO user_follows (follower_id, followed_id)
		VALUES ($1, $2)
		RETURNING id, follower_id, followed_id, created_at
	`, followerID, followedID).StructScan(&follow)
	if err != nil {
		return nil, err
	}
	return &follow, nil
}

func (r *SocialRepository) DeleteUserFollow(followerID, followedID int64) (int64, error) {
	res, err := r.DB.Exec(`
		DELETE FROM user_follows
		WHERE follower_id = $1 AND followed_id = $2
	`, followerID, followedID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *SocialRepository) GetUserFollowers(userID int64) ([]models.UserFollowResponse, error) {
	list := []models.UserFollowResponse{}
	err := r.DB.Select(&list, `
		SELECT
			u.id AS user_id,
			u.username,
			uf.created_at AS followed_at
		FROM user_follows uf
		JOIN users u ON uf.follower_id = u.id
		WHERE uf.followed_id = $1
		ORDER BY uf.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *SocialRepository) GetUserFollowing(userID int64) ([]models.UserFollowResponse, error) {
	list := []models.UserFollowResponse{}
	err := r.DB.Select(&list, `
		SELECT
			u.id AS user_id,
			u.username,
			uf.created_at AS followed_at
		FROM user_follows uf
		JOIN users u ON uf.followed_id = u.id
		WHERE uf.follower_id = $1
		ORDER BY uf.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *SocialRepository) CheckIsFollowing(followerID, followedID int64) (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM user_follows
			WHERE follower_id = $1 AND followed_id = $2
		)
	`, followerID, followedID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *SocialRepository) GetUserStats(userID int64) (models.UserStatsResponse, error) {
	var s models.UserStatsResponse
	type row struct {
		Completed int64 `db:"completed"`
		Planning  int64 `db:"planning"`
		Current   int64 `db:"current"`
		Paused    int64 `db:"paused"`
		Dropped   int64 `db:"dropped"`
	}
	var rcounts row
	if err := r.DB.Get(&rcounts, `
		SELECT
			COUNT(CASE WHEN status = 'completed' THEN 1 END) AS completed,
			COUNT(CASE WHEN status = 'planning' THEN 1 END)  AS planning,
			COUNT(CASE WHEN status = 'current' THEN 1 END)   AS current,
			COUNT(CASE WHEN status = 'paused' THEN 1 END)    AS paused,
			COUNT(CASE WHEN status = 'dropped' THEN 1 END)   AS dropped
		FROM user_list_items
		WHERE user_id = $1
	`, userID); err != nil {
		return s, err
	}
	var followers int64
	if err := r.DB.Get(&followers, `SELECT COUNT(*) FROM user_follows WHERE followed_id = $1`, userID); err != nil {
		return s, err
	}
	var following int64
	if err := r.DB.Get(&following, `SELECT COUNT(*) FROM user_follows WHERE follower_id = $1`, userID); err != nil {
		return s, err
	}
	s = models.UserStatsResponse{
		TotalCompleted: rcounts.Completed,
		TotalPlanning:  rcounts.Planning,
		TotalCurrent:   rcounts.Current,
		TotalPaused:    rcounts.Paused,
		TotalDropped:   rcounts.Dropped,
		FollowersCount: followers,
		FollowingCount: following,
	}
	return s, nil
}

func (r *SocialRepository) GetFollowSuggestions(userID, limit int64) ([]models.UserFollowResponse, error) {
	list := []models.UserFollowResponse{}
	err := r.DB.Select(&list, `
		SELECT
			u.id AS user_id,
			u.username,
			u.created_at AS followed_at
		FROM users u
		WHERE u.id <> $1
		  AND NOT EXISTS (
			  SELECT 1 FROM user_follows uf
			  WHERE uf.follower_id = $1 AND uf.followed_id = u.id
		  )
		ORDER BY u.created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	return list, nil
}
