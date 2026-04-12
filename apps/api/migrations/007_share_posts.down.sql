DROP INDEX IF EXISTS idx_posts_shared_post_id;
ALTER TABLE posts DROP COLUMN IF EXISTS shared_post_id;

-- Note: PostgreSQL does not support removing values from an existing ENUM type.
-- The 'share' value in post_type will remain but will be unused after rollback.
