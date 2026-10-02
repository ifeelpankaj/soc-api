-- name: CreatePlan :one
INSERT INTO plans (
    name, code, description, price_amount_paise, currency, billing_cycle,
    max_flats, max_admins, max_staff, max_residents, features, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, COALESCE(sqlc.narg('features'), '{}'::jsonb), TRUE
)
RETURNING *;

-- name: GetPlan :one
SELECT *
FROM plans
WHERE (sqlc.narg('id')::bigint IS NULL OR id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('code')::text IS NULL OR code = sqlc.narg('code')::text)
  AND (sqlc.narg('name')::text IS NULL OR name = sqlc.narg('name')::text)
  AND (sqlc.narg('billing_cycle')::billing_cycle IS NULL OR billing_cycle = sqlc.narg('billing_cycle')::billing_cycle)
  AND (sqlc.narg('is_active')::bool IS NULL OR is_active = sqlc.narg('is_active')::bool)
  AND (sqlc.narg('min_price_paise')::bigint IS NULL OR price_amount_paise >= sqlc.narg('min_price_paise')::bigint)
  AND (sqlc.narg('max_price_paise')::bigint IS NULL OR price_amount_paise <= sqlc.narg('max_price_paise')::bigint)
  AND (
      COALESCE(sqlc.narg('search')::text, '') = ''
      OR name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(description, '') ILIKE '%' || sqlc.narg('search')::text || '%'
  )
ORDER BY id DESC
LIMIT 1;

-- name: ListPlans :many
SELECT *
FROM plans
WHERE (sqlc.narg('id')::bigint IS NULL OR id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('code')::text IS NULL OR code = sqlc.narg('code')::text)
  AND (sqlc.narg('name')::text IS NULL OR name = sqlc.narg('name')::text)
  AND (sqlc.narg('billing_cycle')::billing_cycle IS NULL OR billing_cycle = sqlc.narg('billing_cycle')::billing_cycle)
  AND (sqlc.narg('is_active')::bool IS NULL OR is_active = sqlc.narg('is_active')::bool)
  AND (sqlc.narg('min_price_paise')::bigint IS NULL OR price_amount_paise >= sqlc.narg('min_price_paise')::bigint)
  AND (sqlc.narg('max_price_paise')::bigint IS NULL OR price_amount_paise <= sqlc.narg('max_price_paise')::bigint)
  AND (
      COALESCE(sqlc.narg('search')::text, '') = ''
      OR name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(description, '') ILIKE '%' || sqlc.narg('search')::text || '%'
  )
ORDER BY price_amount_paise ASC, id DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountPlans :one
SELECT COUNT(*)
FROM plans
WHERE (sqlc.narg('id')::bigint IS NULL OR id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('code')::text IS NULL OR code = sqlc.narg('code')::text)
  AND (sqlc.narg('name')::text IS NULL OR name = sqlc.narg('name')::text)
  AND (sqlc.narg('billing_cycle')::billing_cycle IS NULL OR billing_cycle = sqlc.narg('billing_cycle')::billing_cycle)
  AND (sqlc.narg('is_active')::bool IS NULL OR is_active = sqlc.narg('is_active')::bool)
  AND (sqlc.narg('min_price_paise')::bigint IS NULL OR price_amount_paise >= sqlc.narg('min_price_paise')::bigint)
  AND (sqlc.narg('max_price_paise')::bigint IS NULL OR price_amount_paise <= sqlc.narg('max_price_paise')::bigint)
  AND (
      COALESCE(sqlc.narg('search')::text, '') = ''
      OR name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(description, '') ILIKE '%' || sqlc.narg('search')::text || '%'
  );

-- name: UpdatePlan :one
UPDATE plans
SET
    name = COALESCE(sqlc.narg('name'), name),
    code = COALESCE(sqlc.narg('code'), code),
    description = COALESCE(sqlc.narg('description'), description),
    price_amount_paise = COALESCE(sqlc.narg('price_amount_paise'), price_amount_paise),
    currency = COALESCE(sqlc.narg('currency'), currency),
    billing_cycle = COALESCE(sqlc.narg('billing_cycle'), billing_cycle),
    max_flats = COALESCE(sqlc.narg('max_flats'), max_flats),
    max_admins = COALESCE(sqlc.narg('max_admins'), max_admins),
    max_staff = COALESCE(sqlc.narg('max_staff'), max_staff),
    max_residents = COALESCE(sqlc.narg('max_residents'), max_residents),
    features = COALESCE(sqlc.narg('features'), features),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: ActivatePlan :one
UPDATE plans
SET is_active = TRUE, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeactivatePlan :one
UPDATE plans
SET is_active = FALSE, updated_at = NOW()
WHERE id = $1
RETURNING *;
