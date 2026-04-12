DROP INDEX IF EXISTS idx_user_list_items_queue_position;
ALTER TABLE user_list_items DROP COLUMN IF EXISTS queue_position;

DROP INDEX IF EXISTS idx_completion_reviews_tags;
DROP INDEX IF EXISTS idx_completion_reviews_rating;
DROP INDEX IF EXISTS idx_completion_reviews_media_item_id;
DROP INDEX IF EXISTS idx_completion_reviews_user_id;
DROP TABLE IF EXISTS completion_reviews;

DROP INDEX IF EXISTS idx_user_list_items_collection_id;
ALTER TABLE user_list_items DROP COLUMN IF EXISTS collection_id;

DROP INDEX IF EXISTS idx_collections_is_private;
DROP INDEX IF EXISTS idx_collections_user_id;
DROP TABLE IF EXISTS collections;
