DROP INDEX IF EXISTS idx_account_deletion_requests_scheduled_at;
DROP INDEX IF EXISTS idx_account_deletion_requests_status;
DROP TABLE IF EXISTS account_deletion_requests;

DROP INDEX IF EXISTS idx_ban_appeals_status;
DROP INDEX IF EXISTS idx_ban_appeals_user_id;
DROP TABLE IF EXISTS ban_appeals;

DROP INDEX IF EXISTS idx_users_banned_at;
ALTER TABLE users
    DROP COLUMN IF EXISTS ban_reason,
    DROP COLUMN IF EXISTS banned_at;

DROP INDEX IF EXISTS idx_moderation_queue_ai_categories;
DROP INDEX IF EXISTS idx_moderation_queue_created_at;
DROP INDEX IF EXISTS idx_moderation_queue_reported_by;
DROP INDEX IF EXISTS idx_moderation_queue_entity;
DROP INDEX IF EXISTS idx_moderation_queue_status;
DROP TABLE IF EXISTS moderation_queue;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'deletion_request_status') THEN
        DROP TYPE deletion_request_status;
    END IF;
END
$$;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ban_appeal_status') THEN
        DROP TYPE ban_appeal_status;
    END IF;
END
$$;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'moderation_entity_type') THEN
        DROP TYPE moderation_entity_type;
    END IF;
END
$$;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'moderation_status') THEN
        DROP TYPE moderation_status;
    END IF;
END
$$;
