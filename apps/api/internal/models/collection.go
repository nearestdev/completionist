package models

import "time"

type Collection struct {
	ID          int64     `db:"id" json:"id"`
	UserID      int64     `db:"user_id" json:"userId"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description,omitempty"`
	IsPrivate   bool      `db:"is_private" json:"isPrivate"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

type NewCollection struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IsPrivate   *bool   `json:"isPrivate,omitempty"`
}

type UpdateCollection struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsPrivate   *bool   `json:"isPrivate,omitempty"`
}

type CollectionWithItems struct {
	Collection
	Items []UserListItem `json:"items"`
}
