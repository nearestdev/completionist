DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ad_campaign_status') THEN
        CREATE TYPE ad_campaign_status AS ENUM ('draft', 'active', 'paused', 'completed', 'archived');
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ad_placement') THEN
        CREATE TYPE ad_placement AS ENUM ('feed', 'sidebar', 'banner', 'profile');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS ad_campaigns (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    advertiser_name VARCHAR(255) NOT NULL,
    content_html TEXT NOT NULL,
    image_url VARCHAR(2048),
    target_url VARCHAR(2048) NOT NULL,
    placement ad_placement NOT NULL DEFAULT 'feed',
    status ad_campaign_status NOT NULL DEFAULT 'draft',
    budget_cents INT,
    spent_cents INT NOT NULL DEFAULT 0,
    start_at TIMESTAMPTZ,
    end_at TIMESTAMPTZ,
    targeting_rules JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ad_campaigns_status ON ad_campaigns(status);
CREATE INDEX idx_ad_campaigns_placement ON ad_campaigns(placement);
CREATE INDEX idx_ad_campaigns_date_range ON ad_campaigns(start_at, end_at)
    WHERE status = 'active';
CREATE INDEX idx_ad_campaigns_targeting ON ad_campaigns USING GIN (targeting_rules);

CREATE TABLE IF NOT EXISTS ad_impressions (
    id BIGSERIAL PRIMARY KEY,
    campaign_id BIGINT NOT NULL REFERENCES ad_campaigns(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    impression_type VARCHAR(20) NOT NULL DEFAULT 'view',
    ip_hash VARCHAR(64),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ad_impressions_campaign_id ON ad_impressions(campaign_id);
CREATE INDEX idx_ad_impressions_user_id ON ad_impressions(user_id)
    WHERE user_id IS NOT NULL;
CREATE INDEX idx_ad_impressions_created_at ON ad_impressions(created_at DESC);
CREATE INDEX idx_ad_impressions_campaign_type ON ad_impressions(campaign_id, impression_type);

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS profile_html TEXT,
    ADD COLUMN IF NOT EXISTS profile_theme JSONB DEFAULT '{}';
