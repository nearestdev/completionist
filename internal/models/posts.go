package models

import (
	"time"

	"github.com/google/uuid"
)

type PostType string

const (
	PostTypeReview     PostType = "review"
	PostTypeDiscussion PostType = "discussion"
	PostTypeGeneral    PostType = "general"
)

type Post struct {
	ID          int64      `db:"id" json:"id"`
	UserID      int64      `db:"user_id" json:"userId"`
	MediaItemID *uuid.UUID `db:"media_item_id" json:"mediaItemId,omitempty"`
	Title       *string    `db:"title" json:"title,omitempty"`
	Content     string     `db:"content" json:"content"`
	PostType    PostType   `db:"post_type" json:"postType"`
	Rating      *int       `db:"rating" json:"rating,omitempty"`
	IsSpoiler   bool       `db:"is_spoiler" json:"isSpoiler"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
}

type NewPost struct {
	MediaItemID *uuid.UUID `json:"mediaItemId,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Content     string     `json:"content"`
	PostType    PostType   `json:"postType"`
	Rating      *int       `json:"rating,omitempty"`
	IsSpoiler   *bool      `json:"isSpoiler,omitempty"`
}

type UpdatePost struct {
	Title     *string `json:"title,omitempty"`
	Content   *string `json:"content,omitempty"`
	Rating    *int    `json:"rating,omitempty"`
	IsSpoiler *bool   `json:"isSpoiler,omitempty"`
}

type Comment struct {
	ID              int64     `db:"id" json:"id"`
	UserID          int64     `db:"user_id" json:"userId"`
	PostID          int64     `db:"post_id" json:"postId"`
	ParentCommentID *int64    `db:"parent_comment_id" json:"parentCommentId,omitempty"`
	Content         string    `db:"content" json:"content"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"`
}

type NewComment struct {
	PostID          int64  `json:"postId"`
	ParentCommentID *int64 `json:"parentCommentId,omitempty"`
	Content         string `json:"content"`
}

type UpdateComment struct {
	Content string `json:"content"`
}

type PostLike struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"userId"`
	PostID    int64     `db:"post_id" json:"postId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type CommentLike struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"userId"`
	CommentID int64     `db:"comment_id" json:"commentId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type PostWithDetails struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"userId"`
	Username      string     `json:"username"`
	MediaItemID   *uuid.UUID `json:"mediaItemId,omitempty"`
	MediaTitle    *string    `json:"mediaTitle,omitempty"`
	MediaCoverImage *string  `json:"mediaCoverImage,omitempty"`
	Title         *string    `json:"title,omitempty"`
	Content       string     `json:"content"`
	PostType      PostType   `json:"postType"`
	Rating        *int       `json:"rating,omitempty"`
	IsSpoiler     bool       `json:"isSpoiler"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	LikesCount    int64      `json:"likesCount"`
	CommentsCount int64      `json:"commentsCount"`
	IsLikedByUser bool       `json:"isLikedByUser"`
}

type CommentWithDetails struct {
	ID              int64                `json:"id"`
	UserID          int64                `json:"userId"`
	Username        string               `json:"username"`
	PostID          int64                `json:"postId"`
	ParentCommentID *int64               `json:"parentCommentId,omitempty"`
	Content         string               `json:"content"`
	CreatedAt       time.Time            `json:"createdAt"`
	UpdatedAt       time.Time            `json:"updatedAt"`
	LikesCount      int64                `json:"likesCount"`
	IsLikedByUser   bool                 `json:"isLikedByUser"`
	Replies         []CommentWithDetails `json:"replies,omitempty"`
}
