package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ItemType string

const (
	ItemTypeManga  ItemType = "manga"
	ItemTypeGame   ItemType = "game"
	ItemTypeBook   ItemType = "book"
	ItemTypeMovie  ItemType = "movie"
	ItemTypeSeries ItemType = "series"
	ItemTypeMusic  ItemType = "music"
)

type ItemStatus string

const (
	StatusPlanning  ItemStatus = "planning"
	StatusCurrent   ItemStatus = "current"
	StatusCompleted ItemStatus = "completed"
	StatusPaused    ItemStatus = "paused"
	StatusDropped   ItemStatus = "dropped"
)

type PriorityLevel string

const (
	PriorityLow    PriorityLevel = "low"
	PriorityMedium PriorityLevel = "medium"
	PriorityHigh   PriorityLevel = "high"
)

type MediaItem struct {
	ID                 uuid.UUID        `db:"id" json:"id"`
	ItemType           ItemType         `db:"item_type" json:"itemType"`
	Source             *string          `db:"source" json:"source,omitempty"`
	ExternalID         *string          `db:"external_id" json:"externalId,omitempty"`
	Title              string           `db:"title" json:"title"`
	Description        *string          `db:"description" json:"description,omitempty"`
	CoverImageURL      *string          `db:"cover_image_url" json:"coverImageUrl,omitempty"`
	ReleaseDate        *time.Time       `db:"release_date" json:"releaseDate,omitempty"`
	Genres             pq.StringArray   `db:"genres" json:"genres,omitempty" swaggertype:"array,string"`
	CriticRatingValue  *float32         `db:"critic_rating_value" json:"criticRatingValue,omitempty"`
	CriticRatingCount  *int             `db:"critic_rating_count" json:"criticRatingCount,omitempty"`
	UserRatingExternal *float32         `db:"user_rating_external" json:"userRatingExternal,omitempty"`
	Metadata           *json.RawMessage `db:"metadata" json:"metadata,omitempty" swaggertype:"string"`
	CreatedAt          time.Time        `db:"created_at" json:"createdAt"`
	UpdatedAt          time.Time        `db:"updated_at" json:"updatedAt"`
}

type ProgressData struct {
	Current              *int     `json:"current,omitempty"`
	Total                *int     `json:"total,omitempty"`
	Unit                 *string  `json:"unit,omitempty"`
	HoursPlayed          *float64 `json:"hoursPlayed,omitempty"`
	AchievementsUnlocked *int     `json:"achievementsUnlocked,omitempty"`
	AchievementsTotal    *int     `json:"achievementsTotal,omitempty"`
	Watched              *bool    `json:"watched,omitempty"`
	Season               *int     `json:"season,omitempty"`
	TracksListened       *int     `json:"tracksListened,omitempty"`
	TotalTracks          *int     `json:"totalTracks,omitempty"`
	AlbumsListened       *int     `json:"albumsListened,omitempty"`
	TotalAlbums          *int     `json:"totalAlbums,omitempty"`
}

type UserListItem struct {
	ID            int64            `db:"id" json:"id"`
	UserID        int64            `db:"user_id" json:"userId"`
	MediaItemID   uuid.UUID        `db:"media_item_id" json:"mediaItemId"`
	Status        ItemStatus       `db:"status" json:"status"`
	ProgressData  *json.RawMessage `db:"progress_data" json:"progressData,omitempty" swaggertype:"string"`
	Rating        *int             `db:"rating" json:"rating,omitempty"`
	ItemType      ItemType         `db:"item_type" json:"itemType"`
	CollectionID  *int64           `db:"collection_id" json:"collectionId,omitempty"`
	QueuePosition *int             `db:"queue_position" json:"queuePosition,omitempty"`
	CreatedAt     time.Time        `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time        `db:"updated_at" json:"updatedAt"`
	Media         *MediaItem       `db:"-" json:"media,omitempty"`
}

type WishlistItem struct {
	ID          int64         `db:"id" json:"id"`
	UserID      int64         `db:"user_id" json:"userId"`
	MediaItemID uuid.UUID     `db:"media_item_id" json:"mediaItemId"`
	Priority    PriorityLevel `db:"priority" json:"priority"`
	PriceCents  *int          `db:"price_cents" json:"priceCents,omitempty"`
	CreatedAt   time.Time     `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time     `db:"updated_at" json:"updatedAt"`
	Media       *MediaItem    `db:"-" json:"media,omitempty"`
}

type NewMediaItem struct {
	ItemType      ItemType         `json:"itemType"`
	Source        *string          `json:"source,omitempty"`
	ExternalID    *string          `json:"externalId,omitempty"`
	Title         string           `json:"title"`
	Description   *string          `json:"description,omitempty"`
	CoverImageURL *string          `json:"coverImageUrl,omitempty"`
	ReleaseDate   *time.Time       `json:"releaseDate,omitempty"`
	Genres        []string         `json:"genres,omitempty"`
	Metadata      *json.RawMessage `json:"metadata,omitempty" swaggertype:"string"`
}

type NewUserListItem struct {
	MediaItemID  uuid.UUID        `json:"mediaItemId"`
	Status       ItemStatus       `json:"status"`
	ProgressData *json.RawMessage `json:"progressData,omitempty" swaggertype:"string"`
	Rating       *int             `json:"rating,omitempty"`
}

type UpdateUserListItem struct {
	Status       ItemStatus       `json:"status"`
	ProgressData *json.RawMessage `json:"progressData,omitempty" swaggertype:"string"`
	Rating       *int             `json:"rating,omitempty"`
}

type NewWishlistItem struct {
	MediaItemID uuid.UUID     `json:"mediaItemId"`
	Priority    PriorityLevel `json:"priority"`
	PriceCents  *int          `json:"priceCents,omitempty"`
}

type CreateListItemPayload struct {
	MediaData NewMediaItem        `json:"mediaData"`
	ListData  NewUserListItemBody `json:"listData"`
}

type NewUserListItemBody struct {
	Status       ItemStatus       `json:"status"`
	ProgressData *json.RawMessage `json:"progressData,omitempty" swaggertype:"string"`
	Rating       *int             `json:"rating,omitempty"`
}
