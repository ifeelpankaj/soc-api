-- name: UpsertDeviceTokenForDevice :one
WITH locked AS (
    SELECT pg_advisory_xact_lock((sqlc.arg('user_id') % 2147483647)::integer, hashtext(sqlc.arg('device_id')::text))
),
deleted AS (
    DELETE FROM device_tokens
    WHERE user_id = sqlc.arg('user_id')
      AND token = sqlc.arg('token')
      AND device_id IS DISTINCT FROM sqlc.arg('device_id')::text
    RETURNING 1
)
INSERT INTO device_tokens (
    user_id,
    token,
    platform,
    device_id,
    last_seen_at
)
SELECT
    sqlc.arg('user_id'),
    sqlc.arg('token'),
    sqlc.arg('platform'),
    sqlc.arg('device_id'),
    NOW()
FROM locked, (SELECT COUNT(*) FROM deleted) AS deleted_count
ON CONFLICT (user_id, device_id) WHERE device_id IS NOT NULL DO UPDATE SET
    token = EXCLUDED.token,
    platform = EXCLUDED.platform,
    last_seen_at = NOW(),
    updated_at = NOW()
RETURNING *;

-- name: UpsertDeviceToken :one
INSERT INTO device_tokens (
    user_id,
    token,
    platform,
    device_id,
    last_seen_at
)
VALUES ($1, $2, $3, $4, NOW())
ON CONFLICT (user_id, token) DO UPDATE SET
    platform = EXCLUDED.platform,
    device_id = EXCLUDED.device_id,
    last_seen_at = NOW(),
    updated_at = NOW()
RETURNING *;

-- name: DeleteDeviceToken :exec
DELETE FROM device_tokens
WHERE user_id = $1
  AND token = $2;

-- name: ListDeviceTokensByUserID :many
SELECT *
FROM device_tokens
WHERE user_id = $1
ORDER BY last_seen_at DESC;

-- name: DeleteDeviceTokenByValue :exec
DELETE FROM device_tokens
WHERE token = $1;
