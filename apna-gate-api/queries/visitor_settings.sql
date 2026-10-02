-- name: CreateDefaultSocietyVisitorSettings :exec
INSERT INTO society_visitor_settings (society_id, updated_by)
VALUES ($1, $2)
ON CONFLICT (society_id) DO NOTHING;

-- name: GetSocietyVisitorSettings :one
SELECT *
FROM society_visitor_settings
WHERE society_id = $1;

-- name: UpdateSocietyVisitorSettings :one
UPDATE society_visitor_settings
SET
    approval_mode = COALESCE(sqlc.narg('approval_mode')::visitor_approval_mode, approval_mode),
    default_visit_duration_minutes = COALESCE(sqlc.narg('default_visit_duration_minutes')::integer, default_visit_duration_minutes),
    grace_period_minutes = COALESCE(sqlc.narg('grace_period_minutes')::integer, grace_period_minutes),
    qr_expiry_minutes = COALESCE(sqlc.narg('qr_expiry_minutes')::integer, qr_expiry_minutes),
    allow_resident_pre_approval = COALESCE(sqlc.narg('allow_resident_pre_approval')::boolean, allow_resident_pre_approval),
    allow_public_qr_entry = COALESCE(sqlc.narg('allow_public_qr_entry')::boolean, allow_public_qr_entry),
    allow_guard_entry = COALESCE(sqlc.narg('allow_guard_entry')::boolean, allow_guard_entry),
    allow_guard_on_behalf_approval = COALESCE(sqlc.narg('allow_guard_on_behalf_approval')::boolean, allow_guard_on_behalf_approval),
    is_active = COALESCE(sqlc.narg('is_active')::boolean, is_active),
    updated_by = sqlc.arg('updated_by'),
    updated_at = NOW()
WHERE society_id = sqlc.arg('society_id')
RETURNING *;

-- name: CreateDefaultFlatVisitorSettings :exec
INSERT INTO flat_visitor_settings (
    society_id,
    flat_id,
    purpose,
    approval_required,
    updated_by
)
VALUES
    ($1, $2, 'guest', TRUE, $3),
    ($1, $2, 'delivery', FALSE, $3),
    ($1, $2, 'cab', FALSE, $3),
    ($1, $2, 'service', TRUE, $3),
    ($1, $2, 'maintenance', FALSE, $3),
    ($1, $2, 'staff', TRUE, $3),
    ($1, $2, 'other', TRUE, $3)
ON CONFLICT (flat_id, purpose) DO NOTHING;

-- name: ListFlatVisitorSettings :many
SELECT *
FROM flat_visitor_settings
WHERE society_id = $1
  AND flat_id = $2
ORDER BY CASE purpose
    WHEN 'guest' THEN 1
    WHEN 'delivery' THEN 2
    WHEN 'cab' THEN 3
    WHEN 'service' THEN 4
    WHEN 'maintenance' THEN 5
    WHEN 'staff' THEN 6
    ELSE 7
END;

-- name: GetFlatVisitorPurposeSetting :one
SELECT *
FROM flat_visitor_settings
WHERE society_id = $1
  AND flat_id = $2
  AND purpose = $3;

-- name: UpdateFlatVisitorPurposeSetting :one
UPDATE flat_visitor_settings
SET
    approval_required = COALESCE(sqlc.narg('approval_required')::boolean, approval_required),
    default_visit_duration_minutes = COALESCE(sqlc.narg('default_visit_duration_minutes')::integer, default_visit_duration_minutes),
    is_enabled = COALESCE(sqlc.narg('is_enabled')::boolean, is_enabled),
    updated_by = sqlc.arg('updated_by'),
    updated_at = NOW()
WHERE society_id = sqlc.arg('society_id')
  AND flat_id = sqlc.arg('flat_id')
  AND purpose = sqlc.arg('purpose')
RETURNING *;

-- name: DeleteFlatVisitorSettings :exec
DELETE FROM flat_visitor_settings
WHERE society_id = $1
  AND flat_id = $2;

-- name: ListSocietyFlatVisitorSettings :many
SELECT
    fvs.*,
    f.flat_number,
    f.block
FROM flat_visitor_settings fvs
JOIN flats f ON f.id = fvs.flat_id
WHERE fvs.society_id = sqlc.arg('society_id')
  AND (sqlc.narg('flat_id')::bigint IS NULL OR fvs.flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('block')::text IS NULL OR f.block = sqlc.narg('block')::text)
  AND (sqlc.narg('purpose')::visitor_purpose IS NULL OR fvs.purpose = sqlc.narg('purpose')::visitor_purpose)
ORDER BY f.block NULLS LAST, f.flat_number, CASE fvs.purpose
    WHEN 'guest' THEN 1
    WHEN 'delivery' THEN 2
    WHEN 'cab' THEN 3
    WHEN 'service' THEN 4
    WHEN 'maintenance' THEN 5
    WHEN 'staff' THEN 6
    ELSE 7
END
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSocietyFlatVisitorSettings :one
SELECT COUNT(*)
FROM flat_visitor_settings fvs
JOIN flats f ON f.id = fvs.flat_id
WHERE fvs.society_id = sqlc.arg('society_id')
  AND (sqlc.narg('flat_id')::bigint IS NULL OR fvs.flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('block')::text IS NULL OR f.block = sqlc.narg('block')::text)
  AND (sqlc.narg('purpose')::visitor_purpose IS NULL OR fvs.purpose = sqlc.narg('purpose')::visitor_purpose);
