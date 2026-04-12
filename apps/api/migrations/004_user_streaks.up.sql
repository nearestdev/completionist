DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'streak_tier') THEN
        CREATE TYPE streak_tier AS ENUM ('none', 'spark', 'fire', 'inferno');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS user_streaks (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_streak INT NOT NULL DEFAULT 0,
    longest_streak INT NOT NULL DEFAULT 0,
    last_activity_date DATE,
    tier streak_tier NOT NULL DEFAULT 'none',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
