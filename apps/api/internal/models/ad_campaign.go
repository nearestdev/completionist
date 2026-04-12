package models

import (
	"encoding/json"
	"time"
)

type AdCampaignStatus string

const (
	AdStatusDraft     AdCampaignStatus = "draft"
	AdStatusActive    AdCampaignStatus = "active"
	AdStatusPaused    AdCampaignStatus = "paused"
	AdStatusCompleted AdCampaignStatus = "completed"
	AdStatusArchived  AdCampaignStatus = "archived"
)

type AdPlacement string

const (
	AdPlacementFeed    AdPlacement = "feed"
	AdPlacementSidebar AdPlacement = "sidebar"
	AdPlacementBanner  AdPlacement = "banner"
	AdPlacementProfile AdPlacement = "profile"
)

type AdCampaign struct {
	ID             int64            `db:"id" json:"id"`
	Name           string           `db:"name" json:"name"`
	AdvertiserName string           `db:"advertiser_name" json:"advertiserName"`
	ContentHTML    string           `db:"content_html" json:"contentHtml"`
	ImageURL       *string          `db:"image_url" json:"imageUrl,omitempty"`
	TargetURL      string           `db:"target_url" json:"targetUrl"`
	Placement      AdPlacement      `db:"placement" json:"placement"`
	Status         AdCampaignStatus `db:"status" json:"status"`
	BudgetCents    *int             `db:"budget_cents" json:"budgetCents,omitempty"`
	SpentCents     int              `db:"spent_cents" json:"spentCents"`
	StartAt        *time.Time       `db:"start_at" json:"startAt,omitempty"`
	EndAt          *time.Time       `db:"end_at" json:"endAt,omitempty"`
	TargetingRules json.RawMessage  `db:"targeting_rules" json:"targetingRules" swaggertype:"string"`
	CreatedAt      time.Time        `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time        `db:"updated_at" json:"updatedAt"`
}

type AdImpression struct {
	ID             int64     `db:"id" json:"id"`
	CampaignID     int64     `db:"campaign_id" json:"campaignId"`
	UserID         *int64    `db:"user_id" json:"userId,omitempty"`
	ImpressionType string    `db:"impression_type" json:"impressionType"`
	IPHash         *string   `db:"ip_hash" json:"ipHash,omitempty"`
	UserAgent      *string   `db:"user_agent" json:"userAgent,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
}

type NewAdCampaign struct {
	Name           string          `json:"name"`
	AdvertiserName string          `json:"advertiserName"`
	ContentHTML    string          `json:"contentHtml"`
	ImageURL       *string         `json:"imageUrl,omitempty"`
	TargetURL      string          `json:"targetUrl"`
	Placement      AdPlacement     `json:"placement"`
	BudgetCents    *int            `json:"budgetCents,omitempty"`
	StartAt        *time.Time      `json:"startAt,omitempty"`
	EndAt          *time.Time      `json:"endAt,omitempty"`
	TargetingRules json.RawMessage `json:"targetingRules,omitempty" swaggertype:"string"`
}

type AdCampaignStats struct {
	CampaignID  int64 `db:"campaign_id" json:"campaignId"`
	Views       int64 `db:"views" json:"views"`
	Clicks      int64 `db:"clicks" json:"clicks"`
	UniqueUsers int64 `db:"unique_users" json:"uniqueUsers"`
}
