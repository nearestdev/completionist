CREATE TABLE IF NOT EXISTS collections (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_collections_user_id ON collections(user_id);
CREATE INDEX idx_collections_is_private ON collections(is_private);

ALTER TABLE user_list_items
    ADD COLUMN IF NOT EXISTS collection_id BIGINT REFERENCES collections(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_user_list_items_collection_id ON user_list_items(collection_id);

CREATE TABLE IF NOT EXISTS completion_reviews (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    media_item_id UUID NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 10),
    tags TEXT[],
    favorite_quote TEXT,
    review_text TEXT,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_completion_reviews_user_media UNIQUE (user_id, media_item_id)
);

CREATE INDEX idx_completion_reviews_user_id ON completion_reviews(user_id);
CREATE INDEX idx_completion_reviews_media_item_id ON completion_reviews(media_item_id);
CREATE INDEX idx_completion_reviews_rating ON completion_reviews(rating);
CREATE INDEX idx_completion_reviews_tags ON completion_reviews USING GIN (tags);

ALTER TABLE user_list_items
    ADD COLUMN IF NOT EXISTS queue_position INT;

CREATE INDEX IF NOT EXISTS idx_user_list_items_queue_position ON user_list_items(user_id, queue_position)
    WHERE queue_position IS NOT NULL;
