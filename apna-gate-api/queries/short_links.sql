-- name: CreateShortLink :one
INSERT INTO short_links (
    short_code,
    resource_type,
    resource_id,
    expires_at,
    created_by,
    metadata
)
VALUES (
    $1,
    $2,
    $3,
    sqlc.narg('expires_at'),
    sqlc.narg('created_by'),
    COALESCE(sqlc.narg('metadata'), '{}'::jsonb)
)
RETURNING *;

-- name: GetShortLinkByCode :one
SELECT *
FROM short_links
WHERE short_code = $1;

-- name: GetShortLinkByResource :one
SELECT *
FROM short_links
WHERE resource_type = $1
  AND resource_id = $2
  AND revoked_at IS NULL
ORDER BY created_at DESC
LIMIT 1;

-- name: RevokeShortLink :one
UPDATE short_links
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE short_code = $1
  AND revoked_at IS NULL
RETURNING *;
