-- +migrate Up

WITH ranked_device_tokens AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY user_id, device_id
            ORDER BY last_seen_at DESC, updated_at DESC, id DESC
        ) AS row_number
    FROM device_tokens
    WHERE device_id IS NOT NULL
)
DELETE FROM device_tokens
WHERE id IN (
    SELECT id
    FROM ranked_device_tokens
    WHERE row_number > 1
);

CREATE UNIQUE INDEX idx_device_tokens_user_device_id_unique
ON device_tokens (user_id, device_id)
WHERE device_id IS NOT NULL;

-- +migrate Down

DROP INDEX IF EXISTS idx_device_tokens_user_device_id_unique;
