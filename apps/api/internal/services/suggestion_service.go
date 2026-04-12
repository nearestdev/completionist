package services

import (
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type SuggestionProvider interface {
	GetSuggestions(userID int64, limit int) ([]models.MediaItem, error)
}

type SQLSuggestionProvider struct {
	db *sqlx.DB
}

func NewSQLSuggestionProvider(db *sqlx.DB) *SQLSuggestionProvider {
	return &SQLSuggestionProvider{db: db}
}

func (p *SQLSuggestionProvider) GetSuggestions(userID int64, limit int) ([]models.MediaItem, error) {
	var items []models.MediaItem
	err := p.db.Select(&items, `
		SELECT DISTINCT mi.*
		FROM media_items mi
		JOIN user_list_items other_uli ON other_uli.media_item_id = mi.id
		JOIN user_list_items my_uli ON my_uli.user_id = other_uli.user_id AND my_uli.user_id != $1
		WHERE other_uli.status = 'completed'
		AND my_uli.media_item_id IN (SELECT media_item_id FROM user_list_items WHERE user_id = $1 AND status = 'completed')
		AND mi.id NOT IN (SELECT media_item_id FROM user_list_items WHERE user_id = $1)
		ORDER BY mi.created_at DESC
		LIMIT $2
	`, userID, limit)
	return items, err
}

type SuggestionService struct {
	provider SuggestionProvider
}

func NewSuggestionService(provider SuggestionProvider) *SuggestionService {
	return &SuggestionService{provider: provider}
}

func (s *SuggestionService) GetSuggestions(userID int64, limit int) ([]models.MediaItem, error) {
	return s.provider.GetSuggestions(userID, limit)
}
