CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TYPE item_type AS ENUM ('manga', 'game', 'book', 'movie', 'series', 'music');
CREATE TYPE item_status AS ENUM ('planning', 'current', 'completed', 'paused', 'dropped');
CREATE TYPE priority_level AS ENUM ('low', 'medium', 'high');
CREATE TYPE post_type AS ENUM ('review', 'discussion', 'general');
CREATE TYPE challenge_type AS ENUM ('global', 'personal');
CREATE TYPE challenge_frequency AS ENUM ('daily', 'weekly', 'monthly', 'seasonal', 'infinite');
CREATE TYPE criteria_type AS ENUM ('count_items', 'specific_item', 'genre_count');
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(320) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  xp BIGINT NOT NULL DEFAULT 0,
  level INT NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT uq_users_username UNIQUE (username),
  CONSTRAINT uq_users_email UNIQUE (email)
);
CREATE TABLE media_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  item_type item_type NOT NULL,
  source VARCHAR(50),
  external_id VARCHAR(255),
  title VARCHAR(255) NOT NULL,
  description TEXT,
  cover_image_url VARCHAR(2048),
  release_date DATE,
  genres TEXT[],
  critic_rating_value REAL,
  critic_rating_count INT,
  user_rating_external REAL,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(source, external_id)
);
CREATE TABLE user_list_items (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  media_item_id UUID NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
  status item_status NOT NULL,
  progress_data JSONB DEFAULT '{}',
  rating INT CHECK (rating >= 1 AND rating <= 5),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, media_item_id)
);
CREATE TABLE wishlist_items (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  media_item_id UUID NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
  priority priority_level NOT NULL,
  price_cents INT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, media_item_id)
);
CREATE INDEX idx_user_list_items_user_id ON user_list_items(user_id);
CREATE INDEX idx_user_list_items_progress_data ON user_list_items USING GIN (progress_data);
CREATE INDEX idx_wishlist_items_user_id ON wishlist_items(user_id);
CREATE INDEX idx_media_items_source_external_id ON media_items(source, external_id);
CREATE TABLE game_achievements (
  id BIGSERIAL PRIMARY KEY,
  media_item_id UUID NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  icon_url VARCHAR(2048),
  external_api_name VARCHAR(255),
  UNIQUE(media_item_id, name)
);
CREATE TABLE user_unlocked_achievements (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  achievement_id BIGINT NOT NULL REFERENCES game_achievements(id) ON DELETE CASCADE,
  unlocked_at TIMESTAMPTZ NOT NULL,
  UNIQUE(user_id, achievement_id)
);
CREATE TABLE user_follows (
  id BIGSERIAL PRIMARY KEY,
  follower_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  followed_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(follower_id, followed_id),
  CHECK (follower_id != followed_id)
);
CREATE INDEX idx_user_follows_follower_id ON user_follows(follower_id);
CREATE INDEX idx_user_follows_followed_id ON user_follows(followed_id);
CREATE TABLE posts (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  media_item_id UUID REFERENCES media_items(id) ON DELETE CASCADE,
  title VARCHAR(255),
  content TEXT NOT NULL,
  post_type post_type NOT NULL DEFAULT 'general',
  rating INT CHECK (rating >= 1 AND rating <= 5),
  is_spoiler BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_media_item_id ON posts(media_item_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX idx_posts_type ON posts(post_type);
CREATE TABLE post_likes (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, post_id)
);
CREATE INDEX idx_post_likes_user_id ON post_likes(user_id);
CREATE INDEX idx_post_likes_post_id ON post_likes(post_id);
CREATE TABLE comments (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  parent_comment_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_parent_comment_id ON comments(parent_comment_id);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);
CREATE TABLE comment_likes (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  comment_id BIGINT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, comment_id)
);
CREATE INDEX idx_comment_likes_user_id ON comment_likes(user_id);
CREATE INDEX idx_comment_likes_comment_id ON comment_likes(comment_id);
CREATE TABLE steam_accounts (
  user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  steam_id VARCHAR(32) UNIQUE NOT NULL,
  persona VARCHAR(255) NOT NULL,
  avatar VARCHAR(2048) NOT NULL
);
CREATE TABLE IF NOT EXISTS lastfm_accounts (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(255) NOT NULL,
    session_key VARCHAR(255) NOT NULL,
    subscriber INT NOT NULL
);
CREATE TABLE attachments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  kind VARCHAR(50) NOT NULL,
  value TEXT NOT NULL,
  storage_provider VARCHAR(50),
  content_type VARCHAR(255),
  size_bytes BIGINT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE entity_attachments (
  id BIGSERIAL PRIMARY KEY,
  entity_table VARCHAR(100) NOT NULL,
  entity_pk TEXT NOT NULL,
  attachment_id UUID NOT NULL REFERENCES attachments(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_entity_attachments_entity ON entity_attachments(entity_table, entity_pk);
CREATE INDEX idx_entity_attachments_attachment ON entity_attachments(attachment_id);
CREATE TABLE xp_actions (
    action_key VARCHAR(50) PRIMARY KEY,
    xp_amount INT NOT NULL
);
CREATE TABLE level_definitions (
    level INT PRIMARY KEY,
    xp_required BIGINT NOT NULL
);
CREATE TABLE user_xp_history (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  amount INT NOT NULL,
  source_type VARCHAR(50) NOT NULL,
  source_id VARCHAR(255),
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_user_xp_history_user_id ON user_xp_history(user_id);
CREATE INDEX idx_user_xp_history_source ON user_xp_history(source_type, source_id);
CREATE TABLE seasons (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE challenges (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type challenge_type NOT NULL,
    frequency challenge_frequency NOT NULL,
    xp_reward INT NOT NULL,
    criteria_type criteria_type NOT NULL,
    criteria_metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE user_challenges (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    challenge_id BIGINT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    season_id BIGINT REFERENCES seasons(id) ON DELETE SET NULL,
    current_progress INT NOT NULL DEFAULT 0,
    target_progress INT NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, challenge_id, season_id)
);
CREATE INDEX idx_user_challenges_user ON user_challenges(user_id);
CREATE INDEX idx_user_challenges_completed ON user_challenges(is_completed);
CREATE INDEX idx_challenges_frequency ON challenges(frequency);

CREATE TABLE user_ranks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    season_id BIGINT REFERENCES seasons(id) ON DELETE CASCADE,
    current_rank INT NOT NULL DEFAULT 1,
    current_elo INT NOT NULL DEFAULT 0,
    peak_rank INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, season_id)
);
CREATE INDEX idx_user_ranks_user ON user_ranks(user_id);
CREATE INDEX idx_user_ranks_season ON user_ranks(season_id);

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);

CREATE TABLE direct_messages (
  id BIGSERIAL PRIMARY KEY,
  sender_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  receiver_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_dm_conversation ON direct_messages(sender_id, receiver_id);
CREATE INDEX idx_dm_receiver ON direct_messages(receiver_id, is_read);
CREATE INDEX idx_dm_created_at ON direct_messages(created_at DESC);

CREATE TYPE room_type AS ENUM ('chat', 'music', 'watch');

CREATE TABLE rooms (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  room_type room_type NOT NULL,
  creator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  is_public BOOLEAN NOT NULL DEFAULT TRUE,
  max_members INT DEFAULT 50,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_rooms_type ON rooms(room_type);
CREATE INDEX idx_rooms_public ON rooms(is_public);
CREATE INDEX idx_rooms_created_at ON rooms(created_at DESC);

CREATE TABLE room_members (
  id BIGSERIAL PRIMARY KEY,
  room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  is_moderator BOOLEAN NOT NULL DEFAULT FALSE,
  joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(room_id, user_id)
);
CREATE INDEX idx_room_members_room ON room_members(room_id);
CREATE INDEX idx_room_members_user ON room_members(user_id);

CREATE TABLE room_messages (
  id BIGSERIAL PRIMARY KEY,
  room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_room_messages_room ON room_messages(room_id, created_at DESC);
CREATE INDEX idx_room_messages_user ON room_messages(user_id);

CREATE TABLE room_state (
  room_id BIGINT PRIMARY KEY REFERENCES rooms(id) ON DELETE CASCADE,
  current_media_url TEXT,
  current_media_title VARCHAR(500),
  current_position_ms BIGINT DEFAULT 0,
  is_playing BOOLEAN DEFAULT FALSE,
  updated_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);