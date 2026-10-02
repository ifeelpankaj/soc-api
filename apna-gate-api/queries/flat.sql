-- name: CreateFlat :one
INSERT INTO flats (society_id, block, floor, flat_number, status, is_active, metadata, created_by, flat_type, area_sqft_hundredths)
VALUES ($1, $2, $3, $4, 'vacant', TRUE, COALESCE(sqlc.narg('metadata'), '{}'::jsonb), $5, sqlc.narg('flat_type'), sqlc.narg('area_sqft_hundredths'))
RETURNING *;

-- name: GetFlat :one
SELECT
    f.*,
    s.name AS society_name,
    s.society_code AS society_code
FROM flats f
JOIN societies s ON s.id = f.society_id
WHERE (sqlc.narg('id')::bigint IS NULL OR f.id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR f.society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('block')::text IS NULL OR f.block = sqlc.narg('block')::text)
  AND (sqlc.narg('floor')::text IS NULL OR f.floor = sqlc.narg('floor')::text)
  AND (sqlc.narg('flat_number')::text IS NULL OR f.flat_number = sqlc.narg('flat_number')::text)
  AND (sqlc.narg('status')::flat_status IS NULL OR f.status = sqlc.narg('status')::flat_status)
  AND (sqlc.narg('is_active')::bool IS NULL OR f.is_active = sqlc.narg('is_active')::bool)
ORDER BY f.id DESC
LIMIT 1;

-- name: ListFlats :many
SELECT
    f.*,
    s.name AS society_name,
    s.society_code AS society_code,
    COALESCE(primary_resident.full_name, '') AS primary_resident_name
FROM flats f
JOIN societies s ON s.id = f.society_id
LEFT JOIN LATERAL (
    SELECT u.full_name
    FROM flat_residents fr
    JOIN users u ON u.id = fr.user_id
    WHERE fr.society_id = f.society_id
      AND fr.flat_id = f.id
      AND fr.status = 'active'
      AND fr.is_primary = TRUE
    LIMIT 1
) primary_resident ON TRUE
WHERE (sqlc.narg('id')::bigint IS NULL OR f.id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR f.society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('block')::text IS NULL OR f.block = sqlc.narg('block')::text)
  AND (sqlc.narg('floor')::text IS NULL OR f.floor = sqlc.narg('floor')::text)
  AND (sqlc.narg('flat_number')::text IS NULL OR f.flat_number = sqlc.narg('flat_number')::text)
  AND (sqlc.narg('status')::flat_status IS NULL OR f.status = sqlc.narg('status')::flat_status)
  AND (sqlc.narg('is_active')::bool IS NULL OR f.is_active = sqlc.narg('is_active')::bool)
  AND (
      sqlc.arg('search')::text = ''
      OR f.flat_number ILIKE '%' || sqlc.arg('search')::text || '%'
      OR COALESCE(f.block, '') ILIKE '%' || sqlc.arg('search')::text || '%'
      OR COALESCE(f.floor, '') ILIKE '%' || sqlc.arg('search')::text || '%'
      OR f.status::text ILIKE '%' || sqlc.arg('search')::text || '%'
      OR s.name ILIKE '%' || sqlc.arg('search')::text || '%'
      OR s.society_code ILIKE '%' || sqlc.arg('search')::text || '%'
      OR COALESCE(primary_resident.full_name, '') ILIKE '%' || sqlc.arg('search')::text || '%'
  )
ORDER BY f.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: UpdateFlat :one
UPDATE flats
SET
    flat_type = COALESCE(sqlc.narg('flat_type'), flat_type),
    area_sqft_hundredths = COALESCE(sqlc.narg('area_sqft_hundredths'), area_sqft_hundredths),
    block = COALESCE(sqlc.narg('block'), block),
    floor = COALESCE(sqlc.narg('floor'), floor),
    flat_number = COALESCE(sqlc.narg('flat_number'), flat_number),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    metadata = COALESCE(sqlc.narg('metadata'), metadata),
    updated_at = NOW()
WHERE (sqlc.narg('id')::bigint IS NULL OR id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR society_id = sqlc.narg('society_id')::bigint)
RETURNING *;

-- name: DeactivateFlat :exec
UPDATE flats
SET is_active = FALSE, updated_at = NOW()
WHERE (sqlc.narg('id')::bigint IS NULL OR id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR society_id = sqlc.narg('society_id')::bigint);

-- name: BlockFlat :one
UPDATE flats
SET status = 'blocked', updated_at = NOW()
WHERE (sqlc.narg('id')::bigint IS NULL OR id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR society_id = sqlc.narg('society_id')::bigint)
  AND is_active = TRUE
RETURNING *;

-- name: UnblockFlat :one
UPDATE flats
SET status = CASE WHEN EXISTS (
    SELECT 1 FROM flat_residents fr
    WHERE fr.society_id = flats.society_id AND fr.flat_id = flats.id AND fr.status = 'active'
) THEN 'occupied'::flat_status ELSE 'vacant'::flat_status END, updated_at = NOW()
WHERE (sqlc.narg('id')::bigint IS NULL OR id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR society_id = sqlc.narg('society_id')::bigint)
  AND status = 'blocked'
RETURNING *;

-- name: MarkFlatOccupied :one
UPDATE flats AS f
SET status = 'occupied', updated_at = NOW()
WHERE f.id = $1 AND f.society_id = $2 AND f.is_active = TRUE AND f.status != 'blocked'
  AND EXISTS (SELECT 1 FROM flat_residents fr WHERE fr.society_id = $2 AND fr.flat_id = $1 AND fr.status = 'active')
RETURNING *;

-- name: MarkFlatVacant :one
UPDATE flats AS f
SET status = 'vacant', updated_at = NOW()
WHERE f.id = $1 AND f.society_id = $2 AND f.status != 'blocked'
  AND NOT EXISTS (SELECT 1 FROM flat_residents fr WHERE fr.society_id = $2 AND fr.flat_id = $1 AND fr.status = 'active')
RETURNING *;

-- name: CountFlats :one
SELECT COUNT(*)
FROM flats f
JOIN societies s ON s.id = f.society_id
LEFT JOIN LATERAL (
    SELECT u.full_name
    FROM flat_residents fr
    JOIN users u ON u.id = fr.user_id
    WHERE fr.society_id = f.society_id
      AND fr.flat_id = f.id
      AND fr.status = 'active'
      AND fr.is_primary = TRUE
    LIMIT 1
) primary_resident ON TRUE
WHERE (sqlc.narg('id')::bigint IS NULL OR f.id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR f.society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('block')::text IS NULL OR f.block = sqlc.narg('block')::text)
  AND (sqlc.narg('floor')::text IS NULL OR f.floor = sqlc.narg('floor')::text)
  AND (sqlc.narg('flat_number')::text IS NULL OR f.flat_number = sqlc.narg('flat_number')::text)
  AND (sqlc.narg('status')::flat_status IS NULL OR f.status = sqlc.narg('status')::flat_status)
  AND (sqlc.narg('is_active')::bool IS NULL OR f.is_active = sqlc.narg('is_active')::bool)
  AND (
      sqlc.arg('search')::text = ''
      OR f.flat_number ILIKE '%' || sqlc.arg('search')::text || '%'
      OR COALESCE(f.block, '') ILIKE '%' || sqlc.arg('search')::text || '%'
      OR COALESCE(f.floor, '') ILIKE '%' || sqlc.arg('search')::text || '%'
      OR f.status::text ILIKE '%' || sqlc.arg('search')::text || '%'
      OR s.name ILIKE '%' || sqlc.arg('search')::text || '%'
      OR s.society_code ILIKE '%' || sqlc.arg('search')::text || '%'
      OR COALESCE(primary_resident.full_name, '') ILIKE '%' || sqlc.arg('search')::text || '%'
  );

-- name: GetFlatStats :one
SELECT
    sqlc.arg('society_id')::bigint AS society_id,
    COUNT(*) AS total_flats,
    COUNT(*) FILTER (WHERE status = 'vacant') AS vacant_flats,
    COUNT(*) FILTER (WHERE status = 'occupied') AS occupied_flats,
    COUNT(*) FILTER (WHERE status = 'blocked') AS blocked_flats,
    COUNT(*) FILTER (WHERE is_active = TRUE) AS active_flats,
    COUNT(*) FILTER (WHERE is_active = FALSE) AS inactive_flats
FROM flats
WHERE society_id = sqlc.arg('society_id')::bigint;
