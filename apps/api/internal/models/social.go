package models

import "time"

type UserFollow struct {
	ID         int64     `db:"id" json:"id"`
	FollowerID int64     `db:"follower_id" json:"followerId"`
	FollowedID int64     `db:"followed_id" json:"followedId"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}

type NewUserFollow struct {
	FollowedID int64 `json:"followedId"`
}

type UserFollowResponse struct {
	UserID     int64     `db:"user_id" json:"userId"`
	Username   string    `db:"username" json:"username"`
	FollowedAt time.Time `db:"followed_at" json:"followedAt"`
}

type UserStatsResponse struct {
	TotalCompleted int64 `json:"totalCompleted"`
	TotalPlanning  int64 `json:"totalPlanning"`
	TotalCurrent   int64 `json:"totalCurrent"`
	TotalPaused    int64 `json:"totalPaused"`
	TotalDropped   int64 `json:"totalDropped"`
	FollowersCount int64 `json:"followersCount"`
	FollowingCount int64 `json:"followingCount"`
}

type EnhancedUserProfile struct {
	ID             int64             `json:"id"`
	Username       string            `json:"username"`
	Bio            *string           `json:"bio,omitempty"`
	FavoriteGenres []string          `json:"favoriteGenres,omitempty"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
	Stats          UserStatsResponse `json:"stats"`
	IsFollowing    *bool             `json:"isFollowing,omitempty"`
	IsFollowedBy   *bool             `json:"isFollowedBy,omitempty"`
}
