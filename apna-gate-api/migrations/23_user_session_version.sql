-- +migrate Up
ALTER TABLE users ADD COLUMN session_version BIGINT NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE users DROP COLUMN session_version;
