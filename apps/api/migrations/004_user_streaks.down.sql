DROP TABLE IF EXISTS user_streaks;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'streak_tier') THEN
        DROP TYPE streak_tier;
    END IF;
END
$$;
