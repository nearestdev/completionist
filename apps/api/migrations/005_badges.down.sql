DROP INDEX IF EXISTS idx_badges_criteria_json;
DROP INDEX IF EXISTS idx_badges_criteria_type;
DROP INDEX IF EXISTS idx_user_badges_badge_id;
DROP INDEX IF EXISTS idx_user_badges_user_id;
DROP TABLE IF EXISTS user_badges;
DROP TABLE IF EXISTS badges;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'badge_criteria_type') THEN
        DROP TYPE badge_criteria_type;
    END IF;
END
$$;
