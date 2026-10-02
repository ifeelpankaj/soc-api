-- +migrate Up

CREATE TYPE device_platform AS ENUM (
    'ios',
    'android',
    'web'
);

CREATE TABLE device_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    platform device_platform NOT NULL,
    device_id TEXT,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT device_tokens_user_token_unique UNIQUE (user_id, token)
);

CREATE TRIGGER device_tokens_update_timestamp
BEFORE UPDATE ON device_tokens
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE INDEX idx_device_tokens_user_id ON device_tokens (user_id);

-- +migrate Down

DROP TABLE IF EXISTS device_tokens;
DROP TYPE IF EXISTS device_platform;
