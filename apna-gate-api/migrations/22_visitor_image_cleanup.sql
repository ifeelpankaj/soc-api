-- +migrate Up
CREATE TABLE visitor_image_deletions (
    id BIGSERIAL PRIMARY KEY,
    file_id TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +migrate Down
DROP TABLE visitor_image_deletions;
