-- +migrate Up

CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    society_id BIGINT,
    flat_id BIGINT,
    type VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    data JSONB NOT NULL DEFAULT '{}',
    event_key VARCHAR(255),
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT notifications_user_event_key_unique UNIQUE (user_id, event_key)
);

CREATE INDEX idx_notifications_user_created_at
ON notifications (user_id, created_at DESC, id DESC);

CREATE INDEX idx_notifications_user_unread
ON notifications (user_id, created_at DESC, id DESC)
WHERE read_at IS NULL;

-- +migrate Down

DROP TABLE IF EXISTS notifications;
