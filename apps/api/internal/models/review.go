package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CompletionReview struct {
	ID            int64          `db:"id" json:"id"`
	UserID        int64          `db:"user_id" json:"userId"`
	MediaItemID   uuid.UUID      `db:"media_item_id" json:"mediaItemId"`
	Rating        int            `db:"rating" json:"rating"`
	Tags          pq.StringArray `db:"tags" json:"tags,omitempty" swaggertype:"array,string"`
	FavoriteQuote *string        `db:"favorite_quote" json:"favoriteQuote,omitempty"`
	ReviewText    *string        `db:"review_text" json:"reviewText,omitempty"`
	CompletedAt   *time.Time     `db:"completed_at" json:"completedAt,omitempty"`
	CreatedAt     time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time      `db:"updated_at" json:"updatedAt"`
}

type NewReview struct {
	MediaItemID   uuid.UUID  `json:"mediaItemId"`
	Rating        int        `json:"rating"`
	Tags          []string   `json:"tags,omitempty"`
	FavoriteQuote *string    `json:"favoriteQuote,omitempty"`
	ReviewText    *string    `json:"reviewText,omitempty"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
}

type UpdateReview struct {
	Rating        *int     `json:"rating,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	FavoriteQuote *string  `json:"favoriteQuote,omitempty"`
	ReviewText    *string  `json:"reviewText,omitempty"`
}

type ReviewWithUser struct {
	CompletionReview
	Username string `db:"username" json:"username"`
}
