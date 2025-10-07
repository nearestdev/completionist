package repository

import (
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type ListRepository struct {
	DB *sqlx.DB
}

func NewListRepository(db *sqlx.DB) *ListRepository {
	return &ListRepository{DB: db}
}

func (r *ListRepository) CreateUserListItem(userID int64, item models.NewUserListItem) (*models.UserListItem, error) {
	var out models.UserListItem
	err := r.DB.QueryRowx(`
		INSERT INTO user_list_items (user_id, media_item_id, status, progress, rating)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, media_item_id, status, progress, rating, created_at, updated_at
	`, userID, item.MediaItemID, item.Status, item.Progress, item.Rating).StructScan(&out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ListRepository) GetUserListItems(userID int64) ([]models.UserListItem, error) {
	items := []models.UserListItem{}
	err := r.DB.Select(&items, `
		SELECT id, user_id, media_item_id, status, progress, rating, created_at, updated_at
		FROM user_list_items
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ListRepository) UpdateUserListItem(userID, itemID int64, upd models.UpdateUserListItem) (*models.UserListItem, error) {
	var out models.UserListItem
	err := r.DB.QueryRowx(`
		UPDATE user_list_items
		SET status = $1, progress = $2, rating = $3, updated_at = NOW()
		WHERE id = $4 AND user_id = $5
		RETURNING id, user_id, media_item_id, status, progress, rating, created_at, updated_at
	`, upd.Status, upd.Progress, upd.Rating, itemID, userID).StructScan(&out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ListRepository) DeleteUserListItem(userID, itemID int64) (int64, error) {
	res, err := r.DB.Exec(`
		DELETE FROM user_list_items
		WHERE id = $1 AND user_id = $2
	`, itemID, userID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *ListRepository) CreateWishlistItem(userID int64, item models.NewWishlistItem) (*models.WishlistItem, error) {
	var out models.WishlistItem
	err := r.DB.QueryRowx(`
		INSERT INTO wishlist_items (user_id, media_item_id, priority, price_cents)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, media_item_id, priority, price_cents, created_at, updated_at
	`, userID, item.MediaItemID, item.Priority, item.PriceCents).StructScan(&out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ListRepository) GetWishlistItems(userID int64) ([]models.WishlistItem, error) {
	items := []models.WishlistItem{}
	err := r.DB.Select(&items, `
		SELECT id, user_id, media_item_id, priority, price_cents, created_at, updated_at
		FROM wishlist_items
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ListRepository) DeleteWishlistItem(userID, itemID int64) (int64, error) {
	res, err := r.DB.Exec(`
		DELETE FROM wishlist_items
		WHERE id = $1 AND user_id = $2
	`, itemID, userID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
