package repository

import (
	"database/sql"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FranchiseRepository struct {
	DB *sqlx.DB
}

func NewFranchiseRepository(db *sqlx.DB) *FranchiseRepository {
	return &FranchiseRepository{DB: db}
}

func (r *FranchiseRepository) Create(f *models.Franchise) error {
	return r.DB.QueryRowx(`
		INSERT INTO franchises (name, description, cover_image_url, source, external_id, category, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, description, cover_image_url, source, external_id, category, created_by, created_at, updated_at
	`, f.Name, f.Description, f.CoverImageURL, f.Source, f.ExternalID, f.Category, f.CreatedBy).StructScan(f)
}

func (r *FranchiseRepository) GetByID(id int64) (*models.Franchise, error) {
	var f models.Franchise
	err := r.DB.QueryRowx(`
		SELECT id, name, description, cover_image_url, source, external_id, category, created_by, created_at, updated_at
		FROM franchises WHERE id = $1
	`, id).StructScan(&f)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FranchiseRepository) GetAll(limit, offset int) ([]models.Franchise, error) {
	var franchises []models.Franchise
	err := r.DB.Select(&franchises, `
		SELECT id, name, description, cover_image_url, source, external_id, category, created_by, created_at, updated_at
		FROM franchises ORDER BY name LIMIT $1 OFFSET $2
	`, limit, offset)
	return franchises, err
}

func (r *FranchiseRepository) GetByCategory(category string, limit, offset int) ([]models.Franchise, error) {
	var franchises []models.Franchise
	err := r.DB.Select(&franchises, `
		SELECT id, name, description, cover_image_url, source, external_id, category, created_by, created_at, updated_at
		FROM franchises WHERE category = $1 ORDER BY name LIMIT $2 OFFSET $3
	`, category, limit, offset)
	return franchises, err
}

func (r *FranchiseRepository) UpsertFromSource(f *models.Franchise) error {
	return r.DB.QueryRowx(`
		INSERT INTO franchises (name, description, cover_image_url, source, external_id, category)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (source, external_id) WHERE source IS NOT NULL AND external_id IS NOT NULL
		DO UPDATE SET name = $1, description = $2, cover_image_url = $3, updated_at = NOW()
		RETURNING id, name, description, cover_image_url, source, external_id, category, created_by, created_at, updated_at
	`, f.Name, f.Description, f.CoverImageURL, f.Source, f.ExternalID, f.Category).StructScan(f)
}

func (r *FranchiseRepository) Update(id int64, name *string, description *string, coverImageURL *string) (*models.Franchise, error) {
	var f models.Franchise
	err := r.DB.QueryRowx(`
		UPDATE franchises SET
			name = COALESCE($2, name),
			description = COALESCE($3, description),
			cover_image_url = COALESCE($4, cover_image_url),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, cover_image_url, created_by, created_at, updated_at
	`, id, name, description, coverImageURL).StructScan(&f)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FranchiseRepository) AddItem(franchiseID int64, item *models.FranchiseItem) error {
	return r.DB.QueryRowx(`
		INSERT INTO franchise_items (franchise_id, media_item_id, position, relationship)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (franchise_id, media_item_id) DO UPDATE SET position = $3, relationship = $4
		RETURNING id, franchise_id, media_item_id, position, relationship, created_at
	`, franchiseID, item.MediaItemID, item.Position, item.Relationship).StructScan(item)
}

func (r *FranchiseRepository) RemoveItem(franchiseID int64, mediaItemID uuid.UUID) error {
	_, err := r.DB.Exec(`DELETE FROM franchise_items WHERE franchise_id = $1 AND media_item_id = $2`, franchiseID, mediaItemID)
	return err
}

func (r *FranchiseRepository) GetFranchiseItems(franchiseID int64, viewingUserID *int64) ([]models.FranchiseItemWithMedia, error) {
	var items []models.FranchiseItemWithMedia
	err := r.DB.Select(&items, `
		SELECT
			fi.id, fi.franchise_id, fi.media_item_id, fi.position, fi.relationship, fi.created_at,
			mi.title, mi.cover_image_url, mi.item_type,
			CASE WHEN $2::bigint IS NOT NULL
				THEN EXISTS(SELECT 1 FROM user_list_items uli WHERE uli.user_id = $2 AND uli.media_item_id = fi.media_item_id AND uli.status = 'completed')
				ELSE false
			END AS is_completed
		FROM franchise_items fi
		JOIN media_items mi ON mi.id = fi.media_item_id
		WHERE fi.franchise_id = $1
		ORDER BY fi.position
	`, franchiseID, viewingUserID)
	return items, err
}

func (r *FranchiseRepository) GetByMediaItemID(mediaItemID uuid.UUID) ([]models.Franchise, error) {
	var franchises []models.Franchise
	err := r.DB.Select(&franchises, `
		SELECT f.id, f.name, f.description, f.cover_image_url, f.created_by, f.created_at, f.updated_at
		FROM franchises f
		JOIN franchise_items fi ON fi.franchise_id = f.id
		WHERE fi.media_item_id = $1
	`, mediaItemID)
	return franchises, err
}
