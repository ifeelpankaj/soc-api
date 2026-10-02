-- name: CreateNotification :one
INSERT INTO notifications (
    id,
    user_id,
    society_id,
    flat_id,
    type,
    title,
    body,
    data,
    event_key
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (user_id, event_key) DO NOTHING
RETURNING *;

-- name: ListNotificationsByUser :many
SELECT *
FROM notifications
WHERE user_id = $1
  AND (
    sqlc.narg('cursor_created_at')::timestamptz IS NULL
    OR created_at < sqlc.narg('cursor_created_at')::timestamptz
    OR (created_at = sqlc.narg('cursor_created_at')::timestamptz AND id < sqlc.narg('cursor_id')::uuid)
  )
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg('limit');

-- name: CountUnreadNotificationsByUser :one
SELECT COUNT(*)
FROM notifications
WHERE user_id = $1
  AND read_at IS NULL;

-- name: MarkNotificationRead :one
UPDATE notifications
SET read_at = COALESCE(read_at, NOW())
WHERE id = $1
  AND user_id = $2
RETURNING *;

-- name: MarkAllNotificationsRead :exec
UPDATE notifications
SET read_at = COALESCE(read_at, NOW())
WHERE user_id = $1
  AND read_at IS NULL;

-- name: GetOwnedNotification :one
SELECT * FROM notifications WHERE user_id = $1 AND id = $2;
