package repository

import (
	"encoding/json"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type MediaRepository struct {
	DB *sqlx.DB
}

func NewMediaRepository(db *sqlx.DB) *MediaRepository {
	return &MediaRepository{DB: db}
}

func (r *MediaRepository) FindOrCreate(newItem models.NewMediaItem) (*models.MediaItem, error) {
	if newItem.Source != nil && newItem.ExternalID != nil {
		var existing models.MediaItem
		err := r.DB.Get(&existing, `
			SELECT 
				id, item_type, source, external_id, title, description, cover_image_url, 
				release_date, genres, critic_rating_value, critic_rating_count, 
				user_rating_external, metadata, created_at, updated_at
			FROM media_items
			WHERE source = $1 AND external_id = $2
		`, newItem.Source, newItem.ExternalID)
		if err == nil {
			return &existing, nil
		}
	}

	var created models.MediaItem
	var genresArr pq.StringArray
	if len(newItem.Genres) > 0 {
		genresArr = pq.StringArray(newItem.Genres)
	}

	var md json.RawMessage
	if len(newItem.Metadata) > 0 {
		md = newItem.Metadata
	}

	err := r.DB.QueryRowx(`
		INSERT INTO media_items (
			item_type, source, external_id, title, description, cover_image_url, 
			release_date, genres, metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING 
			id, item_type, source, external_id, title, description, cover_image_url, 
			release_date, genres, critic_rating_value, critic_rating_count, 
			user_rating_external, metadata, created_at, updated_at
	`,
		newItem.ItemType,
		newItem.Source,
		newItem.ExternalID,
		newItem.Title,
		newItem.Description,
		newItem.CoverImageURL,
		newItem.ReleaseDate,
		genresArr,
		md,
	).StructScan(&created)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *MediaRepository) GetByID(id uuid.UUID) (*models.MediaItem, error) {
	var item models.MediaItem
	err := r.DB.Get(&item, `
		SELECT 
			id, item_type, source, external_id, title, description, cover_image_url, 
			release_date, genres, critic_rating_value, critic_rating_count, 
			user_rating_external, metadata, created_at, updated_at
		FROM media_items
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
