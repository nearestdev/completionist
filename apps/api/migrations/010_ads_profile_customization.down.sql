ALTER TABLE users
    DROP COLUMN IF EXISTS profile_theme,
    DROP COLUMN IF EXISTS profile_html;

DROP INDEX IF EXISTS idx_ad_impressions_campaign_type;
DROP INDEX IF EXISTS idx_ad_impressions_created_at;
DROP INDEX IF EXISTS idx_ad_impressions_user_id;
DROP INDEX IF EXISTS idx_ad_impressions_campaign_id;
DROP TABLE IF EXISTS ad_impressions;

DROP INDEX IF EXISTS idx_ad_campaigns_targeting;
DROP INDEX IF EXISTS idx_ad_campaigns_date_range;
DROP INDEX IF EXISTS idx_ad_campaigns_placement;
DROP INDEX IF EXISTS idx_ad_campaigns_status;
DROP TABLE IF EXISTS ad_campaigns;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ad_placement') THEN
        DROP TYPE ad_placement;
    END IF;
END
$$;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ad_campaign_status') THEN
        DROP TYPE ad_campaign_status;
    END IF;
END
$$;
