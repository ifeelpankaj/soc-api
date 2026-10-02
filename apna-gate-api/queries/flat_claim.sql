-- name: SubmitFlatClaim :one
INSERT INTO flat_claim_requests (
    society_id, flat_id, user_id, requested_role, requested_primary, status, note, metadata
) VALUES (
    $1, $2, $3, $4, $5, 'pending', $6, COALESCE(sqlc.narg('metadata'), '{}'::jsonb)
)
RETURNING *;

-- name: GetFlatClaim :one
SELECT
    fc.*,
    u.full_name AS user_name,
    u.email AS user_email,
    u.phone_number AS user_phone,
    f.flat_number AS flat_number,
    f.block AS block,
    f.floor AS floor,
    f.status AS flat_status,
    s.name AS society_name,
    s.society_code AS society_code,
    reviewer.full_name AS reviewer_name,
    reviewer.email AS reviewer_email,
    reviewer.phone_number AS reviewer_phone
FROM flat_claim_requests fc
JOIN users u ON u.id = fc.user_id
LEFT JOIN users reviewer ON reviewer.id = fc.reviewed_by
JOIN flats f ON f.id = fc.flat_id
JOIN societies s ON s.id = fc.society_id
WHERE (sqlc.narg('id')::bigint IS NULL OR fc.id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR fc.society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('flat_id')::bigint IS NULL OR fc.flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('user_id')::bigint IS NULL OR fc.user_id = sqlc.narg('user_id')::bigint)
  AND (sqlc.narg('status')::flat_claim_status IS NULL OR fc.status = sqlc.narg('status')::flat_claim_status)
ORDER BY fc.id DESC
LIMIT 1;

-- name: ListFlatClaims :many
SELECT
    fc.*,
    u.full_name AS user_name,
    u.email AS user_email,
    u.phone_number AS user_phone,
    f.flat_number AS flat_number,
    f.block AS block,
    f.floor AS floor,
    f.status AS flat_status,
    s.name AS society_name,
    s.society_code AS society_code,
    reviewer.full_name AS reviewer_name,
    reviewer.email AS reviewer_email,
    reviewer.phone_number AS reviewer_phone
FROM flat_claim_requests fc
JOIN users u ON u.id = fc.user_id
LEFT JOIN users reviewer ON reviewer.id = fc.reviewed_by
JOIN flats f ON f.id = fc.flat_id
JOIN societies s ON s.id = fc.society_id
WHERE (sqlc.narg('id')::bigint IS NULL OR fc.id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR fc.society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('flat_id')::bigint IS NULL OR fc.flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('user_id')::bigint IS NULL OR fc.user_id = sqlc.narg('user_id')::bigint)
  AND (sqlc.narg('status')::flat_claim_status IS NULL OR fc.status = sqlc.narg('status')::flat_claim_status)
  AND (
      sqlc.arg('search')::text = ''
      OR (
          sqlc.arg('search_mode')::text IN ('', 'all', 'claimant')
          AND (
              u.full_name ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(u.email, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(u.phone_number, '') ILIKE '%' || sqlc.arg('search')::text || '%'
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
          sqlc.arg('search_mode')::text IN ('', 'all', 'reviewer')
          AND (
              COALESCE(reviewer.full_name, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(reviewer.email, '') ILIKE '%' || sqlc.arg('search')::text || '%'
              OR COALESCE(reviewer.phone_number, '') ILIKE '%' || sqlc.arg('search')::text || '%'
          )
      )
      OR (
          sqlc.arg('search_mode')::text IN ('', 'all')
          AND fc.status::text ILIKE '%' || sqlc.arg('search')::text || '%'
      )
  )
ORDER BY fc.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetFlatClaimStats :one
SELECT
    COUNT(*) AS total_claims,
    COUNT(*) FILTER (WHERE status = 'pending') AS pending_claims,
    COUNT(*) FILTER (WHERE status = 'approved') AS approved_claims,
    COUNT(*) FILTER (WHERE status = 'rejected') AS rejected_claims,
    COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled_claims
FROM flat_claim_requests fc
WHERE fc.society_id = sqlc.arg('society_id')::bigint;

-- name: ApproveFlatClaim :one
UPDATE flat_claim_requests
SET status = 'approved', reviewed_by = $3, reviewed_at = NOW(), updated_at = NOW()
WHERE id = $1 AND society_id = $2 AND status = 'pending'
RETURNING *;

-- name: RejectFlatClaim :one
UPDATE flat_claim_requests
SET status = 'rejected', reviewed_by = $3, reviewed_at = NOW(), rejection_reason = $4, updated_at = NOW()
WHERE id = $1 AND society_id = $2 AND status = 'pending'
RETURNING *;

-- name: CancelMyFlatClaim :one
UPDATE flat_claim_requests
SET status = 'cancelled', cancelled_at = NOW(), updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND status = 'pending'
RETURNING *;
