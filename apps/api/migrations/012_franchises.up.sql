CREATE TABLE IF NOT EXISTS franchises (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    cover_image_url VARCHAR(2048),
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_franchises_name UNIQUE (name)
);

CREATE INDEX idx_franchises_name ON franchises(name);
CREATE INDEX idx_franchises_created_by ON franchises(created_by)
    WHERE created_by IS NOT NULL;

CREATE TABLE IF NOT EXISTS franchise_items (
    id BIGSERIAL PRIMARY KEY,
    franchise_id BIGINT NOT NULL REFERENCES franchises(id) ON DELETE CASCADE,
    media_item_id UUID NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
    position INT NOT NULL DEFAULT 0,
    relationship VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_franchise_items_franchise_media UNIQUE (franchise_id, media_item_id)
);

CREATE INDEX idx_franchise_items_franchise_id ON franchise_items(franchise_id);
CREATE INDEX idx_franchise_items_media_item_id ON franchise_items(media_item_id);
CREATE INDEX idx_franchise_items_position ON franchise_items(franchise_id, position);
