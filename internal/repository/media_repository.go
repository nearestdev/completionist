package repository

import (
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

	var md interface{}
	if newItem.Metadata != nil && len(*newItem.Metadata) > 0 {
		md = string(*newItem.Metadata)
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

func (r *MediaRepository) GetTrending(limit int) ([]models.MediaItem, error) {
	var items []models.MediaItem
	query := `
		WITH list_counts AS (
			SELECT media_item_id, COUNT(*) as count
			FROM user_list_items
			GROUP BY media_item_id
		),
		log_counts AS (
			SELECT entity_id::uuid as media_item_id, COUNT(*) as count
			FROM audit_logs
			WHERE entity_type = 'media' AND action = 'view' AND created_at > NOW() - INTERVAL '7 days'
			GROUP BY entity_id
		),
		combined_counts AS (
			SELECT COALESCE(l.media_item_id, g.media_item_id) as media_item_id, 
			       COALESCE(l.count, 0) + COALESCE(g.count, 0) as total_score
			FROM list_counts l
			FULL OUTER JOIN log_counts g ON l.media_item_id = g.media_item_id
		)
		SELECT 
			m.id, m.item_type, m.source, m.external_id, m.title, m.description, m.cover_image_url, 
			m.release_date, m.genres, m.critic_rating_value, m.critic_rating_count, 
			m.user_rating_external, m.metadata, m.created_at, m.updated_at
		FROM combined_counts c
		JOIN media_items m ON c.media_item_id = m.id
		ORDER BY c.total_score DESC
		LIMIT $1
	`
	err := r.DB.Select(&items, query, limit)
	if err != nil {
		return nil, err
	}
	return items, nil
}
