package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/nearestdev/completionist/internal/models"
)

type CollectionRepository struct {
	DB *sqlx.DB
}

func NewCollectionRepository(db *sqlx.DB) *CollectionRepository {
	return &CollectionRepository{DB: db}
}

func (r *CollectionRepository) Create(c *models.Collection) error {
	return r.DB.QueryRowx(`
		INSERT INTO collections (user_id, name, description, is_private)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, name, description, is_private, created_at, updated_at
	`, c.UserID, c.Name, c.Description, c.IsPrivate).StructScan(c)
}

func (r *CollectionRepository) GetByUserID(userID int64) ([]models.Collection, error) {
	var collections []models.Collection
	err := r.DB.Select(&collections, `
		SELECT id, user_id, name, description, is_private, created_at, updated_at
		FROM collections WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	return collections, err
}

func (r *CollectionRepository) GetByID(id int64) (*models.Collection, error) {
	var c models.Collection
	err := r.DB.QueryRowx(`
		SELECT id, user_id, name, description, is_private, created_at, updated_at
		FROM collections WHERE id = $1
	`, id).StructScan(&c)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CollectionRepository) Update(id int64, name *string, description *string, isPrivate *bool) (*models.Collection, error) {
	var c models.Collection
	err := r.DB.QueryRowx(`
		UPDATE collections SET
			name = COALESCE($2, name),
			description = COALESCE($3, description),
			is_private = COALESCE($4, is_private),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, name, description, is_private, created_at, updated_at
	`, id, name, description, isPrivate).StructScan(&c)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CollectionRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM collections WHERE id = $1`, id)
	return err
}

func (r *CollectionRepository) AssignItem(itemID, collectionID int64) error {
	_, err := r.DB.Exec(`UPDATE user_list_items SET collection_id = $2, updated_at = NOW() WHERE id = $1`, itemID, collectionID)
	return err
}

func (r *CollectionRepository) UnassignItem(itemID int64) error {
	_, err := r.DB.Exec(`UPDATE user_list_items SET collection_id = NULL, updated_at = NOW() WHERE id = $1`, itemID)
	return err
}

func (r *CollectionRepository) GetItemsByCollectionID(collectionID int64) ([]models.UserListItem, error) {
	var items []models.UserListItem
	err := r.DB.Select(&items, `
		SELECT id, user_id, media_item_id, status, progress_data, rating, item_type, collection_id, queue_position, created_at, updated_at
		FROM user_list_items WHERE collection_id = $1 ORDER BY created_at DESC
	`, collectionID)
	return items, err
}
