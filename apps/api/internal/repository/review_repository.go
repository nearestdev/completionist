package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/nearestdev/completionist/internal/models"
)

type ReviewRepository struct {
	DB *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) *ReviewRepository {
	return &ReviewRepository{DB: db}
}

func (r *ReviewRepository) Create(review *models.CompletionReview) error {
	return r.DB.QueryRowx(`
		INSERT INTO completion_reviews (user_id, media_item_id, rating, tags, favorite_quote, review_text, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, media_item_id, rating, tags, favorite_quote, review_text, completed_at, created_at, updated_at
	`, review.UserID, review.MediaItemID, review.Rating, pq.Array(review.Tags), review.FavoriteQuote, review.ReviewText, review.CompletedAt).StructScan(review)
}

func (r *ReviewRepository) GetByUserAndMedia(userID int64, mediaItemID uuid.UUID) (*models.CompletionReview, error) {
	var review models.CompletionReview
	err := r.DB.QueryRowx(`
		SELECT id, user_id, media_item_id, rating, tags, favorite_quote, review_text, completed_at, created_at, updated_at
		FROM completion_reviews WHERE user_id = $1 AND media_item_id = $2
	`, userID, mediaItemID).StructScan(&review)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) GetByUserID(userID int64, limit, offset int) ([]models.ReviewWithUser, error) {
	var reviews []models.ReviewWithUser
	err := r.DB.Select(&reviews, `
		SELECT cr.id, cr.user_id, cr.media_item_id, cr.rating, cr.tags, cr.favorite_quote, cr.review_text, cr.completed_at, cr.created_at, cr.updated_at, u.username
		FROM completion_reviews cr
		JOIN users u ON u.id = cr.user_id
		WHERE cr.user_id = $1
		ORDER BY cr.created_at DESC LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	return reviews, err
}

func (r *ReviewRepository) GetByMediaID(mediaItemID uuid.UUID, limit, offset int) ([]models.ReviewWithUser, error) {
	var reviews []models.ReviewWithUser
	err := r.DB.Select(&reviews, `
		SELECT cr.id, cr.user_id, cr.media_item_id, cr.rating, cr.tags, cr.favorite_quote, cr.review_text, cr.completed_at, cr.created_at, cr.updated_at, u.username
		FROM completion_reviews cr
		JOIN users u ON u.id = cr.user_id
		WHERE cr.media_item_id = $1
		ORDER BY cr.created_at DESC LIMIT $2 OFFSET $3
	`, mediaItemID, limit, offset)
	return reviews, err
}

func (r *ReviewRepository) Update(id int64, rating *int, tags []string, favoriteQuote, reviewText *string) (*models.CompletionReview, error) {
	var review models.CompletionReview
	err := r.DB.QueryRowx(`
		UPDATE completion_reviews SET
			rating = COALESCE($2, rating),
			tags = COALESCE($3, tags),
			favorite_quote = COALESCE($4, favorite_quote),
			review_text = COALESCE($5, review_text),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, media_item_id, rating, tags, favorite_quote, review_text, completed_at, created_at, updated_at
	`, id, rating, pq.Array(tags), favoriteQuote, reviewText).StructScan(&review)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM completion_reviews WHERE id = $1`, id)
	return err
}
