package models

import (
	"time"

	"github.com/google/uuid"
)

type Franchise struct {
	ID            int64     `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Description   *string   `db:"description" json:"description,omitempty"`
	CoverImageURL *string   `db:"cover_image_url" json:"coverImageUrl,omitempty"`
	Source        *string   `db:"source" json:"source,omitempty"`
	ExternalID    *string   `db:"external_id" json:"externalId,omitempty"`
	Category      *string   `db:"category" json:"category,omitempty"`
	CreatedBy     *int64    `db:"created_by" json:"createdBy,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
}

type FranchiseItem struct {
	ID           int64     `db:"id" json:"id"`
	FranchiseID  int64     `db:"franchise_id" json:"franchiseId"`
	MediaItemID  uuid.UUID `db:"media_item_id" json:"mediaItemId"`
	Position     int       `db:"position" json:"position"`
	Relationship *string   `db:"relationship" json:"relationship,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
}

type FranchiseWithItems struct {
	Franchise
	Items      []FranchiseItemWithMedia `json:"items"`
	TotalItems int                      `json:"totalItems"`
	Completed  int                      `json:"completed"`
}

type FranchiseItemWithMedia struct {
	FranchiseItem
	Title         string  `db:"title" json:"title"`
	CoverImageURL *string `db:"cover_image_url" json:"coverImageUrl,omitempty"`
	ItemType      string  `db:"item_type" json:"itemType"`
	IsCompleted   bool    `db:"is_completed" json:"isCompleted"`
}

type NewFranchise struct {
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	CoverImageURL *string `json:"coverImageUrl,omitempty"`
}

type NewFranchiseItem struct {
	MediaItemID  uuid.UUID `json:"mediaItemId"`
	Position     int       `json:"position"`
	Relationship *string   `json:"relationship,omitempty"`
}
