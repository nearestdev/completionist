DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'badge_criteria_type') THEN
        CREATE TYPE badge_criteria_type AS ENUM (
            'item_completion_count',
            'item_type_count',
            'streak_threshold',
            'franchise_completion'
        );
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS badges (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(255),
    criteria_type badge_criteria_type NOT NULL,
    criteria_json JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_badges_code UNIQUE (code)
);

CREATE TABLE IF NOT EXISTS user_badges (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_id BIGINT NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    earned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_user_badges_user_badge UNIQUE (user_id, badge_id)
);

CREATE INDEX idx_user_badges_user_id ON user_badges(user_id);
CREATE INDEX idx_user_badges_badge_id ON user_badges(badge_id);
CREATE INDEX idx_badges_criteria_type ON badges(criteria_type);
CREATE INDEX idx_badges_criteria_json ON badges USING GIN (criteria_json);

INSERT INTO badges (code, name, description, icon, criteria_type, criteria_json) VALUES
    ('first_steps', 'First Steps', 'Complete your first item', 'trophy', 'item_completion_count', '{"target": 1}'),
    ('bookworm', 'Bookworm', 'Read 10 books', 'book-open', 'item_type_count', '{"item_type": "book", "target": 10}'),
    ('binge_watcher', 'Binge Watcher', 'Watch 10 series', 'television', 'item_type_count', '{"item_type": "series", "target": 10}'),
    ('gamer', 'Gamer', 'Complete 10 games', 'game-controller', 'item_type_count', '{"item_type": "game", "target": 10}'),
    ('audiophile', 'Audiophile', 'Listen to 10 albums', 'music-notes', 'item_type_count', '{"item_type": "music", "target": 10}'),
    ('streak_spark', 'Streak Spark', 'Reach a 7-day streak', 'flame', 'streak_threshold', '{"target": 7}'),
    ('streak_fire', 'Streak Fire', 'Reach a 30-day streak', 'fire', 'streak_threshold', '{"target": 30}'),
    ('streak_inferno', 'Streak Inferno', 'Reach a 100-day streak', 'fire-simple', 'streak_threshold', '{"target": 100}'),
    ('completionist_bronze', 'Completionist Bronze', 'Complete 25 total items', 'medal', 'item_completion_count', '{"target": 25}'),
    ('completionist_silver', 'Completionist Silver', 'Complete 100 total items', 'medal', 'item_completion_count', '{"target": 100}'),
    ('completionist_gold', 'Completionist Gold', 'Complete 500 total items', 'medal', 'item_completion_count', '{"target": 500}'),
    ('franchise_master', 'Franchise Master', 'Complete all items in a franchise', 'star', 'franchise_completion', '{"target": 1}')
ON CONFLICT (code) DO NOTHING;
