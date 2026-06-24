package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/nearestdev/completionist/internal/models"
)

type AdRepository struct {
	DB *sqlx.DB
}

func NewAdRepository(db *sqlx.DB) *AdRepository {
	return &AdRepository{DB: db}
}

func (r *AdRepository) CreateCampaign(c *models.AdCampaign) error {
	return r.DB.QueryRowx(`
		INSERT INTO ad_campaigns (name, advertiser_name, content_html, image_url, target_url, placement, status, budget_cents, start_at, end_at, targeting_rules)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, name, advertiser_name, content_html, image_url, target_url, placement, status, budget_cents, spent_cents, start_at, end_at, targeting_rules, created_at, updated_at
	`, c.Name, c.AdvertiserName, c.ContentHTML, c.ImageURL, c.TargetURL, c.Placement, c.Status, c.BudgetCents, c.StartAt, c.EndAt, c.TargetingRules).StructScan(c)
}

func (r *AdRepository) GetActiveCampaigns(placement string) ([]models.AdCampaign, error) {
	var campaigns []models.AdCampaign
	err := r.DB.Select(&campaigns, `
		SELECT id, name, advertiser_name, content_html, image_url, target_url, placement, status, budget_cents, spent_cents, start_at, end_at, targeting_rules, created_at, updated_at
		FROM ad_campaigns
		WHERE status = 'active' AND ($1 = '' OR placement = $1::ad_placement)
		AND (start_at IS NULL OR start_at <= NOW())
		AND (end_at IS NULL OR end_at >= NOW())
		ORDER BY RANDOM() LIMIT 1
	`, placement)
	return campaigns, err
}

func (r *AdRepository) RecordImpression(imp *models.AdImpression) error {
	_, err := r.DB.Exec(`
		INSERT INTO ad_impressions (campaign_id, user_id, impression_type, ip_hash, user_agent)
		VALUES ($1, $2, $3, $4, $5)
	`, imp.CampaignID, imp.UserID, imp.ImpressionType, imp.IPHash, imp.UserAgent)
	return err
}

func (r *AdRepository) GetCampaignStats(campaignID int64) (*models.AdCampaignStats, error) {
	var stats models.AdCampaignStats
	err := r.DB.QueryRowx(`
		SELECT
			$1::bigint AS campaign_id,
			COUNT(*) FILTER (WHERE impression_type = 'view') AS views,
			COUNT(*) FILTER (WHERE impression_type = 'click') AS clicks,
			COUNT(DISTINCT user_id) AS unique_users
		FROM ad_impressions WHERE campaign_id = $1
	`, campaignID).StructScan(&stats)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *AdRepository) GetAllCampaigns(limit, offset int) ([]models.AdCampaign, error) {
	var campaigns []models.AdCampaign
	err := r.DB.Select(&campaigns, `
		SELECT id, name, advertiser_name, content_html, image_url, target_url, placement, status, budget_cents, spent_cents, start_at, end_at, targeting_rules, created_at, updated_at
		FROM ad_campaigns ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, limit, offset)
	return campaigns, err
}

func (r *AdRepository) UpdateCampaign(id int64, c *models.AdCampaign) error {
	_, err := r.DB.Exec(`
		UPDATE ad_campaigns SET
			name = $2, advertiser_name = $3, content_html = $4, image_url = $5, target_url = $6,
			placement = $7, status = $8, budget_cents = $9, start_at = $10, end_at = $11,
			targeting_rules = $12, updated_at = NOW()
		WHERE id = $1
	`, id, c.Name, c.AdvertiserName, c.ContentHTML, c.ImageURL, c.TargetURL, c.Placement, c.Status, c.BudgetCents, c.StartAt, c.EndAt, c.TargetingRules)
	return err
}
