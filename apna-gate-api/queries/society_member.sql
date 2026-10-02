-- name: AddSocietyMember :one
INSERT INTO society_members (society_id, user_id, role, status, invited_by, metadata)
VALUES ($1, $2, $3, 'active', $4, COALESCE(sqlc.narg('metadata'), '{}'::jsonb))
RETURNING *;

-- name: GetSocietyMember :one
SELECT
    sm.*,
    u.full_name AS user_full_name,
    u.email AS user_email,
    u.phone_number AS user_phone
FROM society_members sm
JOIN users u ON u.id = sm.user_id
WHERE (sqlc.narg('id')::bigint IS NULL OR sm.id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR sm.society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('user_id')::bigint IS NULL OR sm.user_id = sqlc.narg('user_id')::bigint)
  AND (sqlc.narg('role')::society_member_role IS NULL OR sm.role = sqlc.narg('role')::society_member_role)
  AND (sqlc.narg('status')::society_member_status IS NULL OR sm.status = sqlc.narg('status')::society_member_status)
  AND (sqlc.narg('email')::text IS NULL OR u.email = sqlc.narg('email')::text)
  AND (sqlc.narg('phone')::text IS NULL OR u.phone_number = sqlc.narg('phone')::text)
ORDER BY sm.id DESC
LIMIT 1;

-- name: ListSocietyMembers :many
SELECT
    sm.*,
    u.full_name AS user_full_name,
    u.email AS user_email,
    u.phone_number AS user_phone
FROM society_members sm
JOIN users u ON u.id = sm.user_id
LEFT JOIN users invited_user ON invited_user.id = sm.invited_by
LEFT JOIN users removed_user ON removed_user.id = sm.removed_by
WHERE sm.society_id = sqlc.arg('society_id')
  AND (sqlc.narg('role')::society_member_role IS NULL OR sm.role = sqlc.narg('role')::society_member_role)
  AND (sqlc.narg('status')::society_member_status IS NULL OR sm.status = sqlc.narg('status')::society_member_status)
  AND (sqlc.narg('user_id')::bigint IS NULL OR sm.user_id = sqlc.narg('user_id')::bigint)
  AND (sqlc.narg('invited_by')::bigint IS NULL OR sm.invited_by = sqlc.narg('invited_by')::bigint)
  AND (sqlc.narg('removed_by')::bigint IS NULL OR sm.removed_by = sqlc.narg('removed_by')::bigint)
  AND (sqlc.narg('joined_from')::timestamptz IS NULL OR sm.joined_at >= sqlc.narg('joined_from')::timestamptz)
  AND (sqlc.narg('joined_to')::timestamptz IS NULL OR sm.joined_at <= sqlc.narg('joined_to')::timestamptz)
  AND (
      sqlc.arg('search')::text = ''
      OR (
          sqlc.arg('search_mode')::text IN ('', 'all', 'resident')
          AND (
              u.full_name ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(u.email, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(u.phone_number, '') ILIKE '%' || sqlc.arg('search')::text || '%'
          )
      )
      OR (
          sqlc.arg('search_mode')::text IN ('', 'all', 'invited_by')
          AND (
              COALESCE(invited_user.full_name, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(invited_user.email, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(invited_user.phone_number, '') ILIKE '%' || sqlc.arg('search')::text || '%'
          )
      )
      OR (
          sqlc.arg('search_mode')::text IN ('', 'all', 'removed_by')
          AND (
              COALESCE(removed_user.full_name, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(removed_user.email, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(removed_user.phone_number, '') ILIKE '%' || sqlc.arg('search')::text || '%'
          )
      )
      OR (
          sqlc.arg('search_mode')::text IN ('', 'all')
          AND (
              sm.role::text ILIKE '%' || sqlc.arg('search')::text || '%'
              OR sm.status::text ILIKE '%' || sqlc.arg('search')::text || '%'
          )
      )
  )
ORDER BY
    CASE WHEN sqlc.arg('sort_by') = 'role' AND sqlc.arg('sort_order') = 'asc' THEN sm.role END ASC,
    CASE WHEN sqlc.arg('sort_by') = 'role' AND sqlc.arg('sort_order') = 'desc' THEN sm.role END DESC,
    CASE WHEN sqlc.arg('sort_by') = 'status' AND sqlc.arg('sort_order') = 'asc' THEN sm.status END ASC,
    CASE WHEN sqlc.arg('sort_by') = 'status' AND sqlc.arg('sort_order') = 'desc' THEN sm.status END DESC,
    CASE WHEN sqlc.arg('sort_by') = 'joined_at' AND sqlc.arg('sort_order') = 'asc' THEN sm.joined_at END ASC,
    sm.joined_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSocietyMembers :one
SELECT COUNT(*)
FROM society_members sm
JOIN users u ON u.id = sm.user_id
LEFT JOIN users invited_user ON invited_user.id = sm.invited_by
LEFT JOIN users removed_user ON removed_user.id = sm.removed_by
WHERE sm.society_id = sqlc.arg('society_id')
  AND (sqlc.narg('role')::society_member_role IS NULL OR sm.role = sqlc.narg('role')::society_member_role)
  AND (sqlc.narg('status')::society_member_status IS NULL OR sm.status = sqlc.narg('status')::society_member_status)
  AND (sqlc.narg('user_id')::bigint IS NULL OR sm.user_id = sqlc.narg('user_id')::bigint)
  AND (sqlc.narg('invited_by')::bigint IS NULL OR sm.invited_by = sqlc.narg('invited_by')::bigint)
  AND (sqlc.narg('removed_by')::bigint IS NULL OR sm.removed_by = sqlc.narg('removed_by')::bigint)
  AND (sqlc.narg('joined_from')::timestamptz IS NULL OR sm.joined_at >= sqlc.narg('joined_from')::timestamptz)
  AND (sqlc.narg('joined_to')::timestamptz IS NULL OR sm.joined_at <= sqlc.narg('joined_to')::timestamptz)
  AND (
      sqlc.arg('search')::text = ''
      OR (
          sqlc.arg('search_mode')::text IN ('', 'all', 'resident')
          AND (
              u.full_name ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(u.email, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(u.phone_number, '') ILIKE '%' || sqlc.arg('search')::text || '%'
          )
      )
      OR (
          sqlc.arg('search_mode')::text IN ('', 'all', 'invited_by')
          AND (
              COALESCE(invited_user.full_name, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(invited_user.email, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(invited_user.phone_number, '') ILIKE '%' || sqlc.arg('search')::text || '%'
          )
      )
      OR (
          sqlc.arg('search_mode')::text IN ('', 'all', 'removed_by')
          AND (
              COALESCE(removed_user.full_name, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(removed_user.email, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(removed_user.phone_number, '') ILIKE '%' || sqlc.arg('search')::text || '%'
          )
      )
      OR (
          sqlc.arg('search_mode')::text IN ('', 'all')
          AND (
              sm.role::text ILIKE '%' || sqlc.arg('search')::text || '%'
              OR sm.status::text ILIKE '%' || sqlc.arg('search')::text || '%'
          )
      )
  );

-- name: ListSocietyMembershipsByUser :many
SELECT
    sm.id, sm.society_id, sm.user_id, sm.role, sm.status, sm.invited_by, sm.joined_at, sm.removed_by, sm.removed_at, sm.remove_reason, sm.metadata, sm.created_at, sm.updated_at,
    u.full_name AS user_full_name,
    u.email AS user_email,
    u.phone_number AS user_phone
FROM society_members sm
JOIN users u ON u.id = sm.user_id
WHERE sm.user_id = $1
  AND sm.status != 'removed'
ORDER BY sm.joined_at DESC;

-- name: ListMySocietiesByUser :many
SELECT
    sm.id, sm.society_id, sm.user_id, sm.role, sm.status, sm.invited_by, sm.joined_at, sm.removed_by, sm.removed_at, sm.remove_reason, sm.metadata, sm.created_at, sm.updated_at,
    u.full_name AS user_full_name,
    u.email AS user_email,
    u.phone_number AS user_phone,
    s.id AS society_row_id,
    s.name AS society_name,
    s.society_code AS society_code,
    s.email AS society_email,
    s.phone_number AS society_phone_number,
    s.address_line1 AS society_address_line1,
    s.address_line2 AS society_address_line2,
    s.landmark AS society_landmark,
    s.city AS society_city,
    s.state AS society_state,
    s.pincode AS society_pincode,
    s.country AS society_country,
    s.total_flats AS society_total_flats,
    s.total_blocks AS society_total_blocks,
    s.status AS society_status,
    s.created_by AS society_created_by,
    s.approved_by AS society_approved_by,
    s.approved_at AS society_approved_at,
    s.rejected_by AS society_rejected_by,
    s.rejected_at AS society_rejected_at,
    s.rejection_reason AS society_rejection_reason,
    s.suspended_by AS society_suspended_by,
    s.suspended_at AS society_suspended_at,
    s.suspension_reason AS society_suspension_reason,
    s.metadata AS society_metadata,
    s.created_at AS society_created_at,
    s.updated_at AS society_updated_at,
    s.deleted_at AS society_deleted_at
FROM society_members sm
JOIN users u ON u.id = sm.user_id
JOIN societies s ON s.id = sm.society_id
WHERE sm.user_id = $1
  AND sm.status != 'removed'
  AND s.deleted_at IS NULL
ORDER BY sm.joined_at DESC;

-- name: ChangeSocietyMemberRole :one
UPDATE society_members
SET role = $3, updated_at = NOW()
WHERE society_id = $1 AND user_id = $2
RETURNING *;

-- name: SuspendSocietyMember :one
UPDATE society_members
SET status = 'suspended', updated_at = NOW()
WHERE society_id = $1 AND user_id = $2 AND status = 'active'
RETURNING *;

-- name: ReactivateSocietyMember :one
UPDATE society_members
SET status = 'active', removed_by = NULL, removed_at = NULL, remove_reason = NULL, updated_at = NOW()
WHERE society_id = $1 AND user_id = $2 AND status IN ('suspended', 'removed')
RETURNING *;

-- name: RemoveSocietyMember :exec
UPDATE society_members
SET status = 'removed', removed_by = $3, removed_at = NOW(), remove_reason = $4, updated_at = NOW()
WHERE society_id = $1 AND user_id = $2 AND status != 'removed';

-- name: CountActiveOwners :one
SELECT COUNT(*)
FROM society_members
WHERE society_id = $1 AND role = 'owner' AND status = 'active';

-- name: DemoteActiveOwners :exec
UPDATE society_members
SET role = 'admin', updated_at = NOW()
WHERE society_id = $1 AND role = 'owner' AND status = 'active' AND user_id <> $2;

-- name: PromoteMemberToOwner :one
UPDATE society_members
SET role = 'owner', status = 'active', removed_by = NULL, removed_at = NULL, remove_reason = NULL, updated_at = NOW()
WHERE society_id = $1 AND user_id = $2
RETURNING *;

-- name: UpsertResidentSocietyMember :one
INSERT INTO society_members (society_id, user_id, role, status, invited_by, metadata)
VALUES ($1, $2, 'resident', 'active', $3, '{}'::jsonb)
ON CONFLICT (society_id, user_id) DO UPDATE
SET
    role = CASE
        WHEN society_members.role IN ('owner', 'admin', 'staff') THEN society_members.role
        ELSE 'resident'
    END,
    status = 'active',
    removed_by = NULL,
    removed_at = NULL,
    remove_reason = NULL,
    updated_at = NOW()
RETURNING *;
