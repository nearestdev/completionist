ALTER TYPE post_type ADD VALUE IF NOT EXISTS 'share';

ALTER TABLE posts
    ADD COLUMN IF NOT EXISTS shared_post_id BIGINT REFERENCES posts(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_posts_shared_post_id ON posts(shared_post_id)
    WHERE shared_post_id IS NOT NULL;
