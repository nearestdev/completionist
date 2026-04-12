DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'moderation_status') THEN
        CREATE TYPE moderation_status AS ENUM ('pending', 'approved', 'rejected');
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'moderation_entity_type') THEN
        CREATE TYPE moderation_entity_type AS ENUM ('post', 'comment', 'room_message', 'direct_message');
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ban_appeal_status') THEN
        CREATE TYPE ban_appeal_status AS ENUM ('pending', 'accepted', 'denied');
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'deletion_request_status') THEN
        CREATE TYPE deletion_request_status AS ENUM ('pending', 'processing', 'completed', 'cancelled');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS moderation_queue (
    id BIGSERIAL PRIMARY KEY,
    entity_type moderation_entity_type NOT NULL,
    entity_id BIGINT NOT NULL,
    reported_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    reason TEXT,
    ai_flagged BOOLEAN NOT NULL DEFAULT FALSE,
    ai_categories JSONB DEFAULT '{}',
    ai_scores JSONB DEFAULT '{}',
    status moderation_status NOT NULL DEFAULT 'pending',
    reviewed_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    review_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_moderation_queue_status ON moderation_queue(status);
CREATE INDEX idx_moderation_queue_entity ON moderation_queue(entity_type, entity_id);
CREATE INDEX idx_moderation_queue_reported_by ON moderation_queue(reported_by);
CREATE INDEX idx_moderation_queue_created_at ON moderation_queue(created_at DESC);
CREATE INDEX idx_moderation_queue_ai_categories ON moderation_queue USING GIN (ai_categories);

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS banned_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS ban_reason TEXT;

CREATE INDEX IF NOT EXISTS idx_users_banned_at ON users(banned_at)
    WHERE banned_at IS NOT NULL;

CREATE TABLE IF NOT EXISTS ban_appeals (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    appeal_text TEXT NOT NULL,
    status ban_appeal_status NOT NULL DEFAULT 'pending',
    reviewed_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    review_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ban_appeals_user_id ON ban_appeals(user_id);
CREATE INDEX idx_ban_appeals_status ON ban_appeals(status);

CREATE TABLE IF NOT EXISTS account_deletion_requests (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT,
    status deletion_request_status NOT NULL DEFAULT 'pending',
    scheduled_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_account_deletion_requests_user UNIQUE (user_id)
);

CREATE INDEX idx_account_deletion_requests_status ON account_deletion_requests(status);
CREATE INDEX idx_account_deletion_requests_scheduled_at ON account_deletion_requests(scheduled_at)
    WHERE status = 'pending';
