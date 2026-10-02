-- name: CreateFlatMemberInvite :one
INSERT INTO flat_member_invites (
    society_id,
    flat_id,
    invited_by,
    role,
    phone,
    email,
    full_name,
    token_hash,
    expires_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetFlatMemberInviteByID :one
SELECT *
FROM flat_member_invites
WHERE id = $1
  AND society_id = $2;

-- name: GetFlatMemberInviteByIDAny :one
SELECT *
FROM flat_member_invites
WHERE id = $1;

-- name: GetFlatMemberInviteByTokenHash :one
SELECT *
FROM flat_member_invites
WHERE token_hash = $1;

-- name: ListPendingFlatMemberInvites :many
SELECT *
FROM flat_member_invites
WHERE society_id = $1
  AND flat_id = $2
  AND status = 'pending'
ORDER BY created_at DESC;

-- name: ListFlatMemberInviteHistory :many
SELECT
    fmi.id,
    fmi.society_id,
    fmi.flat_id,
    fmi.invited_by,
    fmi.role,
    fmi.phone,
    fmi.email,
    fmi.full_name,
    fmi.status,
    fmi.expires_at,
    fmi.created_at,
    fmi.updated_at,
    sl.short_code,
    sl.expires_at AS short_link_expires_at
FROM flat_member_invites fmi
LEFT JOIN short_links sl
    ON sl.resource_type = 'member_invite'
   AND sl.resource_id = fmi.id
   AND sl.revoked_at IS NULL
WHERE fmi.society_id = $1
  AND fmi.flat_id = $2
  AND (sqlc.narg('status')::flat_member_invite_status IS NULL OR fmi.status = sqlc.narg('status')::flat_member_invite_status)
  AND (
    sqlc.narg('search')::text IS NULL
    OR fmi.full_name ILIKE '%' || sqlc.narg('search')::text || '%'
    OR fmi.phone ILIKE '%' || sqlc.narg('search')::text || '%'
    OR fmi.email ILIKE '%' || sqlc.narg('search')::text || '%'
    OR fmi.role::text ILIKE '%' || sqlc.narg('search')::text || '%'
    OR fmi.status::text ILIKE '%' || sqlc.narg('search')::text || '%'
    OR sl.short_code ILIKE '%' || sqlc.narg('search')::text || '%'
  )
ORDER BY fmi.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountFlatMemberInviteHistory :one
SELECT COUNT(*)
FROM flat_member_invites
WHERE flat_member_invites.society_id = $1
  AND flat_member_invites.flat_id = $2
  AND (sqlc.narg('status')::flat_member_invite_status IS NULL OR flat_member_invites.status = sqlc.narg('status')::flat_member_invite_status)
  AND (
    sqlc.narg('search')::text IS NULL
    OR flat_member_invites.full_name ILIKE '%' || sqlc.narg('search')::text || '%'
    OR flat_member_invites.phone ILIKE '%' || sqlc.narg('search')::text || '%'
    OR flat_member_invites.email ILIKE '%' || sqlc.narg('search')::text || '%'
    OR flat_member_invites.role::text ILIKE '%' || sqlc.narg('search')::text || '%'
    OR flat_member_invites.status::text ILIKE '%' || sqlc.narg('search')::text || '%'
    OR EXISTS (
      SELECT 1
      FROM short_links sl
      WHERE sl.resource_type = 'member_invite'
        AND sl.resource_id = flat_member_invites.id
        AND sl.revoked_at IS NULL
        AND sl.short_code ILIKE '%' || sqlc.narg('search')::text || '%'
    )
  );

-- name: GetFlatMemberInviteHistoryByID :one
SELECT
    fmi.id,
    fmi.society_id,
    fmi.flat_id,
    fmi.invited_by,
    fmi.role,
    fmi.phone,
    fmi.email,
    fmi.full_name,
    fmi.status,
    fmi.expires_at,
    fmi.created_at,
    fmi.updated_at,
    sl.short_code,
    sl.expires_at AS short_link_expires_at
FROM flat_member_invites fmi
LEFT JOIN short_links sl
    ON sl.resource_type = 'member_invite'
   AND sl.resource_id = fmi.id
   AND sl.revoked_at IS NULL
WHERE fmi.society_id = $1
  AND fmi.flat_id = $2
  AND fmi.id = $3;

-- name: CancelFlatMemberInvite :one
UPDATE flat_member_invites
SET status = 'cancelled',
    updated_at = NOW()
WHERE id = $1
  AND society_id = $2
  AND flat_id = $3
  AND status = 'pending'
RETURNING *;

-- name: AcceptFlatMemberInvite :one
UPDATE flat_member_invites
SET status = 'accepted',
    updated_at = NOW()
WHERE id = $1
  AND status = 'pending'
  AND expires_at > NOW()
RETURNING *;

-- name: ExpireOldFlatMemberInvites :exec
UPDATE flat_member_invites
SET status = 'expired',
    updated_at = NOW()
WHERE status = 'pending'
  AND expires_at < NOW();
