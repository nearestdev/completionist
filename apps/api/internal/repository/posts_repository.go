package repository

import (
	"database/sql"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostsRepository struct {
	DB *sqlx.DB
}

func NewPostsRepository(db *sqlx.DB) *PostsRepository {
	return &PostsRepository{DB: db}
}

func (r *PostsRepository) CreatePost(userID int64, np models.NewPost) (models.Post, error) {
	var p models.Post
	err := r.DB.QueryRowx(`
		INSERT INTO posts (user_id, media_item_id, title, content, post_type, rating, is_spoiler)
		VALUES ($1, $2, $3, $4, $5, $6, COALESCE($7, false))
		RETURNING id, user_id, media_item_id, title, content, post_type, rating, is_spoiler, created_at, updated_at
	`, userID, np.MediaItemID, np.Title, np.Content, np.PostType, np.Rating, np.IsSpoiler).
		StructScan(&p)
	return p, err
}

func (r *PostsRepository) GetPostsWithDetails(viewingUserID *int64, limit, offset int64) ([]models.PostWithDetails, error) {
	rows, err := r.DB.Queryx(`
		SELECT
			p.id,
			p.user_id,
			u.username,
			p.media_item_id,
			m.title AS media_title,
			m.cover_image_url AS media_cover_image,
			p.title,
			p.content,
			p.post_type,
			p.rating,
			p.is_spoiler,
			p.created_at,
			p.updated_at,
			(SELECT COUNT(*)::bigint FROM post_likes pl WHERE pl.post_id = p.id) AS likes_count,
			(SELECT COUNT(*)::bigint FROM comments c WHERE c.post_id = p.id) AS comments_count,
			CASE
				WHEN $1::bigint IS NULL THEN false
				ELSE EXISTS (SELECT 1 FROM post_likes pl2 WHERE pl2.post_id = p.id AND pl2.user_id = $1)
			END AS is_liked_by_user
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN media_items m ON m.id = p.media_item_id
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`, viewingUserID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.PostWithDetails{}
	for rows.Next() {
		var item models.PostWithDetails
		if err := rows.StructScan(&item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *PostsRepository) GetPostByIDWithDetails(postID int64, viewingUserID *int64) (*models.PostWithDetails, error) {
	var item models.PostWithDetails
	err := r.DB.QueryRowx(`
		SELECT
			p.id,
			p.user_id,
			u.username,
			p.media_item_id,
			m.title AS media_title,
			m.cover_image_url AS media_cover_image,
			p.title,
			p.content,
			p.post_type,
			p.rating,
			p.is_spoiler,
			p.created_at,
			p.updated_at,
			(SELECT COUNT(*)::bigint FROM post_likes pl WHERE pl.post_id = p.id) AS likes_count,
			(SELECT COUNT(*)::bigint FROM comments c WHERE c.post_id = p.id) AS comments_count,
			CASE
				WHEN $2::bigint IS NULL THEN false
				ELSE EXISTS (SELECT 1 FROM post_likes pl2 WHERE pl2.post_id = p.id AND pl2.user_id = $2)
			END AS is_liked_by_user
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN media_items m ON m.id = p.media_item_id
		WHERE p.id = $1
	`, postID, viewingUserID).StructScan(&item)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PostsRepository) GetPostsByUserID(userID int64, viewingUserID *int64, limit, offset int64) ([]models.PostWithDetails, error) {
	rows, err := r.DB.Queryx(`
		SELECT
			p.id,
			p.user_id,
			u.username,
			p.media_item_id,
			m.title AS media_title,
			m.cover_image_url AS media_cover_image,
			p.title,
			p.content,
			p.post_type,
			p.rating,
			p.is_spoiler,
			p.created_at,
			p.updated_at,
			(SELECT COUNT(*)::bigint FROM post_likes pl WHERE pl.post_id = p.id) AS likes_count,
			(SELECT COUNT(*)::bigint FROM comments c WHERE c.post_id = p.id) AS comments_count,
			CASE
				WHEN $2::bigint IS NULL THEN false
				ELSE EXISTS (SELECT 1 FROM post_likes pl2 WHERE pl2.post_id = p.id AND pl2.user_id = $2)
			END AS is_liked_by_user
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN media_items m ON m.id = p.media_item_id
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC
		LIMIT $3 OFFSET $4
	`, userID, viewingUserID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.PostWithDetails{}
	for rows.Next() {
		var item models.PostWithDetails
		if err := rows.StructScan(&item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *PostsRepository) UpdatePost(postID, userID int64, up models.UpdatePost) (models.Post, error) {
	var p models.Post
	err := r.DB.QueryRowx(`
		UPDATE posts
		SET
			title = COALESCE($1, title),
			content = COALESCE($2, content),
			rating = COALESCE($3, rating),
			is_spoiler = COALESCE($4, is_spoiler),
			updated_at = NOW()
		WHERE id = $5 AND user_id = $6
		RETURNING id, user_id, media_item_id, title, content, post_type, rating, is_spoiler, created_at, updated_at
	`, up.Title, up.Content, up.Rating, up.IsSpoiler, postID, userID).StructScan(&p)
	return p, err
}

func (r *PostsRepository) DeletePost(postID, userID int64) (int64, error) {
	res, err := r.DB.Exec(`DELETE FROM posts WHERE id = $1 AND user_id = $2`, postID, userID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *PostsRepository) CreatePostLike(userID, postID int64) (models.PostLike, error) {
	var like models.PostLike
	err := r.DB.QueryRowx(`
		INSERT INTO post_likes (user_id, post_id)
		VALUES ($1, $2)
		RETURNING id, user_id, post_id, created_at
	`, userID, postID).StructScan(&like)
	return like, err
}

func (r *PostsRepository) DeletePostLike(userID, postID int64) (int64, error) {
	res, err := r.DB.Exec(`DELETE FROM post_likes WHERE user_id = $1 AND post_id = $2`, userID, postID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *PostsRepository) CreateComment(userID int64, nc models.NewComment) (models.Comment, error) {
	var c models.Comment
	err := r.DB.QueryRowx(`
		INSERT INTO comments (user_id, post_id, parent_comment_id, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, post_id, parent_comment_id, content, created_at, updated_at
	`, userID, nc.PostID, nc.ParentCommentID, nc.Content).StructScan(&c)
	return c, err
}

func (r *PostsRepository) GetCommentsWithDetails(postID int64, viewingUserID *int64) ([]models.CommentWithDetails, error) {
	rows, err := r.DB.Queryx(`
		SELECT
			c.id,
			c.user_id,
			u.username,
			c.post_id,
			c.parent_comment_id,
			c.content,
			c.created_at,
			c.updated_at,
			(SELECT COUNT(*)::bigint FROM comment_likes cl WHERE cl.comment_id = c.id) AS likes_count,
			CASE
				WHEN $2::bigint IS NULL THEN false
				ELSE EXISTS (SELECT 1 FROM comment_likes cl2 WHERE cl2.comment_id = c.id AND cl2.user_id = $2)
			END AS is_liked_by_user
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1 AND c.parent_comment_id IS NULL
		ORDER BY c.created_at ASC
	`, postID, viewingUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.CommentWithDetails
	for rows.Next() {
		var base models.CommentWithDetails
		if err := rows.StructScan(&base); err != nil {
			return nil, err
		}
		replies, err := r.getCommentReplies(base.ID, viewingUserID)
		if err != nil {
			return nil, err
		}
		base.Replies = replies
		out = append(out, base)
	}
	return out, rows.Err()
}

func (r *PostsRepository) getCommentReplies(parentID int64, viewingUserID *int64) ([]models.CommentWithDetails, error) {
	rows, err := r.DB.Queryx(`
		SELECT
			c.id,
			c.user_id,
			u.username,
			c.post_id,
			c.parent_comment_id,
			c.content,
			c.created_at,
			c.updated_at,
			(SELECT COUNT(*)::bigint FROM comment_likes cl WHERE cl.comment_id = c.id) AS likes_count,
			CASE
				WHEN $2::bigint IS NULL THEN false
				ELSE EXISTS (SELECT 1 FROM comment_likes cl2 WHERE cl2.comment_id = c.id AND cl2.user_id = $2)
			END AS is_liked_by_user
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.parent_comment_id = $1
		ORDER BY c.created_at ASC
	`, parentID, viewingUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var replies []models.CommentWithDetails
	for rows.Next() {
		var child models.CommentWithDetails
		if err := rows.StructScan(&child); err != nil {
			return nil, err
		}
		replies = append(replies, child)
	}
	return replies, rows.Err()
}

func (r *PostsRepository) UpdateComment(commentID, userID int64, uc models.UpdateComment) (models.Comment, error) {
	var c models.Comment
	err := r.DB.QueryRowx(`
		UPDATE comments
		SET content = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, post_id, parent_comment_id, content, created_at, updated_at
	`, uc.Content, commentID, userID).StructScan(&c)
	return c, err
}

func (r *PostsRepository) DeleteComment(commentID, userID int64) (int64, error) {
	res, err := r.DB.Exec(`DELETE FROM comments WHERE id = $1 AND user_id = $2`, commentID, userID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *PostsRepository) CreateCommentLike(userID, commentID int64) (models.CommentLike, error) {
	var like models.CommentLike
	err := r.DB.QueryRowx(`
		INSERT INTO comment_likes (user_id, comment_id)
		VALUES ($1, $2)
		RETURNING id, user_id, comment_id, created_at
	`, userID, commentID).StructScan(&like)
	return like, err
}

func (r *PostsRepository) DeleteCommentLike(userID, commentID int64) (int64, error) {
	res, err := r.DB.Exec(`DELETE FROM comment_likes WHERE user_id = $1 AND comment_id = $2`, userID, commentID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *PostsRepository) FindMediaTitleAndCover(id *uuid.UUID) (*string, *string, error) {
	if id == nil {
		return nil, nil, nil
	}
	var title sql.NullString
	var cover sql.NullString
	if err := r.DB.QueryRow(`
		SELECT title, cover_image_url FROM media_items WHERE id = $1
	`, id).Scan(&title, &cover); err != nil {
		return nil, nil, err
	}
	var t *string
	var c *string
	if title.Valid {
		t = &title.String
	}
	if cover.Valid {
		c = &cover.String
	}
	return t, c, nil
}
