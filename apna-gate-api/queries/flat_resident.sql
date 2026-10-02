-- name: AddFlatResident :one
INSERT INTO flat_residents (society_id, flat_id, user_id, role, status, is_primary, metadata, created_by)
VALUES ($1, $2, $3, $4, 'active', $5, COALESCE(sqlc.narg('metadata'), '{}'::jsonb), $6)
RETURNING *;

-- name: GetFlatResident :one
SELECT
    fr.*,
    u.full_name AS user_name,
    u.email AS user_email,
    u.phone_number AS user_phone,
    f.flat_number AS flat_number,
    f.block AS block,
    f.floor AS floor,
    s.name AS society_name,
    s.society_code AS society_code
FROM flat_residents fr
JOIN users u ON u.id = fr.user_id
JOIN flats f ON f.id = fr.flat_id
JOIN societies s ON s.id = fr.society_id
WHERE (sqlc.narg('id')::bigint IS NULL OR fr.id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR fr.society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('flat_id')::bigint IS NULL OR fr.flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('user_id')::bigint IS NULL OR fr.user_id = sqlc.narg('user_id')::bigint)
  AND (sqlc.narg('role')::flat_resident_role IS NULL OR fr.role = sqlc.narg('role')::flat_resident_role)
  AND (sqlc.narg('status')::flat_resident_status IS NULL OR fr.status = sqlc.narg('status')::flat_resident_status)
  AND (sqlc.narg('is_primary')::bool IS NULL OR fr.is_primary = sqlc.narg('is_primary')::bool)
ORDER BY fr.id DESC
LIMIT 1;

-- name: ListFlatResidents :many
SELECT
    fr.*,
    u.full_name AS user_name,
    u.email AS user_email,
    u.phone_number AS user_phone,
    f.flat_number AS flat_number,
    f.block AS block,
    f.floor AS floor,
    s.name AS society_name,
    s.society_code AS society_code
FROM flat_residents fr
JOIN users u ON u.id = fr.user_id
JOIN flats f ON f.id = fr.flat_id
JOIN societies s ON s.id = fr.society_id
WHERE (sqlc.narg('id')::bigint IS NULL OR fr.id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR fr.society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('flat_id')::bigint IS NULL OR fr.flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('user_id')::bigint IS NULL OR fr.user_id = sqlc.narg('user_id')::bigint)
  AND (sqlc.narg('role')::flat_resident_role IS NULL OR fr.role = sqlc.narg('role')::flat_resident_role)
  AND (sqlc.narg('status')::flat_resident_status IS NULL OR fr.status = sqlc.narg('status')::flat_resident_status)
  AND (sqlc.narg('is_primary')::bool IS NULL OR fr.is_primary = sqlc.narg('is_primary')::bool)
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
          sqlc.arg('search_mode')::text IN ('', 'all', 'society')
          AND (
              s.name ILIKE '%' || sqlc.arg('search')::text || '%'
              OR s.society_code ILIKE '%' || sqlc.arg('search')::text || '%'
          )
      )
      OR (
          sqlc.arg('search_mode')::text IN ('', 'all', 'flat')
          AND (
              f.flat_number ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(f.block, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(f.floor, '') ILIKE '%' || sqlc.arg('search')::text || '%'
          )
      )
      OR (
          sqlc.arg('search_mode')::text IN ('', 'all')
          AND (
              fr.role::text ILIKE '%' || sqlc.arg('search')::text || '%'
              OR fr.status::text ILIKE '%' || sqlc.arg('search')::text || '%'
          )
      )
  )
ORDER BY fr.moved_in_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: RemoveFlatResident :exec
UPDATE flat_residents
SET status = 'inactive', is_primary = FALSE, updated_at = NOW()
WHERE (sqlc.narg('id')::bigint IS NULL OR id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('flat_id')::bigint IS NULL OR flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('user_id')::bigint IS NULL OR user_id = sqlc.narg('user_id')::bigint)
  AND status = 'active';

-- name: MoveOutFlatResident :one
UPDATE flat_residents
SET status = 'moved_out', is_primary = FALSE, moved_out_at = NOW(), updated_at = NOW()
WHERE (sqlc.narg('id')::bigint IS NULL OR id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('flat_id')::bigint IS NULL OR flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('user_id')::bigint IS NULL OR user_id = sqlc.narg('user_id')::bigint)
  AND status = 'active'
RETURNING *;

-- name: ClearPrimaryFlatResident :exec
UPDATE flat_residents
SET is_primary = FALSE,
    role = 'family',
    updated_at = NOW()
WHERE society_id = $1
  AND flat_id = $2
  AND status = 'active'
  AND is_primary = TRUE;

-- name: SetPrimaryFlatResident :one
UPDATE flat_residents
SET is_primary = TRUE,
    role = 'owner',
    updated_at = NOW()
WHERE society_id = $1
  AND flat_id = $2
  AND id = $3
  AND status = 'active'
RETURNING *;

-- name: UpdateFlatResidentRole :one
UPDATE flat_residents
SET role = $1, updated_at = NOW()
WHERE (sqlc.narg('id')::bigint IS NULL OR id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('flat_id')::bigint IS NULL OR flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('user_id')::bigint IS NULL OR user_id = sqlc.narg('user_id')::bigint)
  AND status = 'active'
RETURNING *;

-- name: CountActiveFlatResidents :one
SELECT COUNT(*)
FROM flat_residents
WHERE society_id = $1 AND flat_id = $2 AND status = 'active';

-- name: CountActivePrimaryFlatResidents :one
SELECT COUNT(*)
FROM flat_residents
WHERE society_id = $1 AND flat_id = $2 AND status = 'active' AND is_primary = TRUE;
