-- name: CreatePendingSubscription :one
INSERT INTO society_subscriptions (
    society_id, plan_id, status, plan_name, plan_code, price_amount_paise, currency,
    billing_cycle, max_flats, max_admins, max_staff, max_residents, features, metadata, created_by
)
SELECT
    sqlc.arg('society_id')::bigint AS society_id,
    p.id AS plan_id,
    'pending'::subscription_status AS status,
    p.name AS plan_name,
    p.code AS plan_code,
    p.price_amount_paise AS price_amount_paise,
    p.currency AS currency,
    p.billing_cycle, p.max_flats, p.max_admins, p.max_staff, p.max_residents, p.features,
    COALESCE(sqlc.narg('metadata'), '{}'::jsonb) AS metadata,
    sqlc.arg('created_by')::bigint AS created_by
FROM plans p
WHERE p.id = sqlc.arg('plan_id')::bigint AND p.is_active = TRUE
RETURNING *;

-- name: CreateTrialSubscription :one
INSERT INTO society_subscriptions (
    society_id, plan_id, status, starts_at, ends_at, trial_ends_at,
    plan_name, plan_code, price_amount_paise, currency, billing_cycle,
    max_flats, max_admins, max_staff, max_residents, features, metadata, created_by
)
SELECT
    sqlc.arg('society_id')::bigint AS society_id,
    p.id AS plan_id,
    'trial'::subscription_status AS status,
    sqlc.arg('starts_at')::timestamptz AS starts_at,
    sqlc.narg('ends_at')::timestamptz AS ends_at,
    sqlc.arg('trial_ends_at')::timestamptz AS trial_ends_at,
    p.name, p.code, p.price_amount_paise, p.currency, p.billing_cycle,
    p.max_flats, p.max_admins, p.max_staff, p.max_residents, p.features,
    COALESCE(sqlc.narg('metadata'), '{}'::jsonb) AS metadata,
    sqlc.arg('created_by')::bigint AS created_by
FROM plans p
WHERE p.id = sqlc.arg('plan_id')::bigint AND p.is_active = TRUE
RETURNING *;

-- name: GetSubscription :one
SELECT
    ss.*,
    s.name AS society_name,
    s.society_code AS society_code,
    p.name AS current_plan_name,
    p.code AS current_plan_code
FROM society_subscriptions ss
JOIN societies s ON s.id = ss.society_id
LEFT JOIN plans p ON p.id = ss.plan_id
LEFT JOIN users created_user ON created_user.id = ss.created_by
LEFT JOIN users activated_user ON activated_user.id = ss.activated_by
LEFT JOIN users cancelled_user ON cancelled_user.id = ss.cancelled_by
WHERE (sqlc.narg('id')::bigint IS NULL OR ss.id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR ss.society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('plan_id')::bigint IS NULL OR ss.plan_id = sqlc.narg('plan_id')::bigint)
  AND (sqlc.narg('status')::subscription_status IS NULL OR ss.status = sqlc.narg('status')::subscription_status)
  AND (sqlc.narg('plan_code')::text IS NULL OR ss.plan_code = sqlc.narg('plan_code')::text)
  AND (sqlc.narg('billing_cycle')::billing_cycle IS NULL OR ss.billing_cycle = sqlc.narg('billing_cycle')::billing_cycle)
  AND (
      sqlc.narg('is_active_only')::bool IS NULL
      OR sqlc.narg('is_active_only')::bool = FALSE
      OR (ss.status IN ('trial', 'active') AND (ss.ends_at IS NULL OR ss.ends_at > NOW()))
  )
  AND (sqlc.narg('starts_after')::timestamptz IS NULL OR ss.starts_at >= sqlc.narg('starts_after')::timestamptz)
  AND (sqlc.narg('starts_before')::timestamptz IS NULL OR ss.starts_at <= sqlc.narg('starts_before')::timestamptz)
  AND (sqlc.narg('ends_after')::timestamptz IS NULL OR ss.ends_at >= sqlc.narg('ends_after')::timestamptz)
  AND (sqlc.narg('ends_before')::timestamptz IS NULL OR ss.ends_at <= sqlc.narg('ends_before')::timestamptz)
  AND (sqlc.narg('expiring_before')::timestamptz IS NULL OR ss.ends_at <= sqlc.narg('expiring_before')::timestamptz)
  AND (
      sqlc.narg('expired_only')::bool IS NULL
      OR sqlc.narg('expired_only')::bool = FALSE
      OR (ss.ends_at IS NOT NULL AND ss.ends_at <= NOW())
  )
  AND (
      COALESCE(sqlc.narg('search')::text, '') = ''
      OR (
          COALESCE(sqlc.narg('search_mode')::text, 'all') IN ('', 'all', 'society')
          AND (
              s.name ILIKE '%' || sqlc.narg('search')::text || '%'
              OR s.society_code ILIKE '%' || sqlc.narg('search')::text || '%'
          )
      )
      OR (
          COALESCE(sqlc.narg('search_mode')::text, 'all') IN ('', 'all', 'plan')
          AND (
              ss.plan_name ILIKE '%' || sqlc.narg('search')::text || '%'
              OR ss.plan_code ILIKE '%' || sqlc.narg('search')::text || '%'
          )
      )
      OR (
          COALESCE(sqlc.narg('search_mode')::text, 'all') IN ('', 'all', 'action_user')
          AND (
              COALESCE(created_user.full_name, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(created_user.email, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(created_user.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(activated_user.full_name, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(activated_user.email, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(activated_user.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(cancelled_user.full_name, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(cancelled_user.email, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(cancelled_user.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
          )
      )
      OR (
          COALESCE(sqlc.narg('search_mode')::text, 'all') IN ('', 'all', 'resident_member')
          AND EXISTS (
              SELECT 1
              FROM society_members sm
              JOIN users member_user ON member_user.id = sm.user_id
              WHERE sm.society_id = ss.society_id
                AND sm.status != 'removed'
                AND (
                    member_user.full_name ILIKE '%' || sqlc.narg('search')::text || '%'
                    OR COALESCE(member_user.email, '') ILIKE '%' || sqlc.narg('search')::text || '%'
                    OR COALESCE(member_user.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
                )
          )
      )
      OR (
          COALESCE(sqlc.narg('search_mode')::text, 'all') IN ('', 'all')
          AND ss.status::text ILIKE '%' || sqlc.narg('search')::text || '%'
      )
  )
ORDER BY ss.id DESC
LIMIT 1;

-- name: ListSubscriptions :many
SELECT
    ss.*,
    s.name AS society_name,
    s.society_code AS society_code,
    p.name AS current_plan_name,
    p.code AS current_plan_code
FROM society_subscriptions ss
JOIN societies s ON s.id = ss.society_id
LEFT JOIN plans p ON p.id = ss.plan_id
LEFT JOIN users created_user ON created_user.id = ss.created_by
LEFT JOIN users activated_user ON activated_user.id = ss.activated_by
LEFT JOIN users cancelled_user ON cancelled_user.id = ss.cancelled_by
WHERE (sqlc.narg('id')::bigint IS NULL OR ss.id = sqlc.narg('id')::bigint)
  AND (sqlc.narg('society_id')::bigint IS NULL OR ss.society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('plan_id')::bigint IS NULL OR ss.plan_id = sqlc.narg('plan_id')::bigint)
  AND (sqlc.narg('status')::subscription_status IS NULL OR ss.status = sqlc.narg('status')::subscription_status)
  AND (sqlc.narg('plan_code')::text IS NULL OR ss.plan_code = sqlc.narg('plan_code')::text)
  AND (sqlc.narg('billing_cycle')::billing_cycle IS NULL OR ss.billing_cycle = sqlc.narg('billing_cycle')::billing_cycle)
  AND (
      sqlc.narg('is_active_only')::bool IS NULL
      OR sqlc.narg('is_active_only')::bool = FALSE
      OR (ss.status IN ('trial', 'active') AND (ss.ends_at IS NULL OR ss.ends_at > NOW()))
  )
  AND (sqlc.narg('starts_after')::timestamptz IS NULL OR ss.starts_at >= sqlc.narg('starts_after')::timestamptz)
  AND (sqlc.narg('starts_before')::timestamptz IS NULL OR ss.starts_at <= sqlc.narg('starts_before')::timestamptz)
  AND (sqlc.narg('ends_after')::timestamptz IS NULL OR ss.ends_at >= sqlc.narg('ends_after')::timestamptz)
  AND (sqlc.narg('ends_before')::timestamptz IS NULL OR ss.ends_at <= sqlc.narg('ends_before')::timestamptz)
  AND (sqlc.narg('expiring_before')::timestamptz IS NULL OR ss.ends_at <= sqlc.narg('expiring_before')::timestamptz)
  AND (
      sqlc.narg('expired_only')::bool IS NULL
      OR sqlc.narg('expired_only')::bool = FALSE
      OR (ss.ends_at IS NOT NULL AND ss.ends_at <= NOW())
  )
  AND (
      COALESCE(sqlc.narg('search')::text, '') = ''
      OR (
          COALESCE(sqlc.narg('search_mode')::text, 'all') IN ('', 'all', 'society')
          AND (
              s.name ILIKE '%' || sqlc.narg('search')::text || '%'
              OR s.society_code ILIKE '%' || sqlc.narg('search')::text || '%'
          )
      )
      OR (
          COALESCE(sqlc.narg('search_mode')::text, 'all') IN ('', 'all', 'plan')
          AND (
              ss.plan_name ILIKE '%' || sqlc.narg('search')::text || '%'
              OR ss.plan_code ILIKE '%' || sqlc.narg('search')::text || '%'
          )
      )
      OR (
          COALESCE(sqlc.narg('search_mode')::text, 'all') IN ('', 'all', 'action_user')
          AND (
              COALESCE(created_user.full_name, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(created_user.email, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(created_user.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(activated_user.full_name, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(activated_user.email, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(activated_user.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(cancelled_user.full_name, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(cancelled_user.email, '') ILIKE '%' || sqlc.narg('search')::text || '%'
              OR COALESCE(cancelled_user.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
          )
      )
      OR (
          COALESCE(sqlc.narg('search_mode')::text, 'all') IN ('', 'all', 'resident_member')
          AND EXISTS (
              SELECT 1
              FROM society_members sm
              JOIN users member_user ON member_user.id = sm.user_id
              WHERE sm.society_id = ss.society_id
                AND sm.status != 'removed'
                AND (
                    member_user.full_name ILIKE '%' || sqlc.narg('search')::text || '%'
                    OR COALESCE(member_user.email, '') ILIKE '%' || sqlc.narg('search')::text || '%'
                    OR COALESCE(member_user.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
                )
          )
      )
      OR (
          COALESCE(sqlc.narg('search_mode')::text, 'all') IN ('', 'all')
          AND ss.status::text ILIKE '%' || sqlc.narg('search')::text || '%'
      )
  )
ORDER BY ss.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ActivateSubscription :one
UPDATE society_subscriptions
SET status = 'active', starts_at = $2, ends_at = $3, activated_by = $4, activated_at = NOW(),
    cancelled_at = NULL, cancelled_by = NULL, cancellation_reason = NULL, expired_at = NULL,
    metadata = COALESCE(sqlc.narg('metadata'), metadata), updated_at = NOW()
WHERE id = $1 AND status IN ('pending', 'trial', 'expired')
RETURNING *;

-- name: RenewSubscription :one
UPDATE society_subscriptions
SET status = 'active', starts_at = $2, ends_at = $3, activated_by = $4, activated_at = NOW(),
    expired_at = NULL, metadata = COALESCE(sqlc.narg('metadata'), metadata), updated_at = NOW()
WHERE id = $1 AND status IN ('active', 'trial', 'expired')
RETURNING *;

-- name: CancelSubscription :one
UPDATE society_subscriptions
SET status = 'cancelled', cancelled_by = $2, cancelled_at = NOW(), cancellation_reason = $3,
    metadata = COALESCE(sqlc.narg('metadata'), metadata), updated_at = NOW()
WHERE id = $1 AND status != 'cancelled'
RETURNING *;

-- name: ExpireSubscription :one
UPDATE society_subscriptions
SET status = 'expired', expired_at = NOW(), updated_at = NOW()
WHERE id = $1 AND status IN ('trial', 'active')
RETURNING *;

-- name: ChangeSubscriptionPlan :one
UPDATE society_subscriptions ss
SET plan_id = p.id,
    plan_name = p.name,
    plan_code = p.code,
    price_amount_paise = p.price_amount_paise,
    currency = p.currency,
    billing_cycle = p.billing_cycle,
    max_flats = p.max_flats,
    max_admins = p.max_admins,
    max_staff = p.max_staff,
    max_residents = p.max_residents,
    features = p.features,
    updated_at = NOW()
FROM plans p
WHERE ss.id = $1 AND p.id = $2 AND p.is_active = TRUE
RETURNING ss.*;

-- name: GetSubscriptionStats :one
SELECT
    COUNT(*) AS total_subscriptions,
    COUNT(*) FILTER (WHERE status = 'pending') AS pending_subscriptions,
    COUNT(*) FILTER (WHERE status = 'trial') AS trial_subscriptions,
    COUNT(*) FILTER (WHERE status = 'active') AS active_subscriptions,
    COUNT(*) FILTER (WHERE status = 'expired') AS expired_subscriptions,
    COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled_subscriptions
FROM society_subscriptions ss
WHERE (sqlc.narg('society_id')::bigint IS NULL OR ss.society_id = sqlc.narg('society_id')::bigint)
  AND (sqlc.narg('plan_id')::bigint IS NULL OR ss.plan_id = sqlc.narg('plan_id')::bigint)
  AND (sqlc.narg('status')::subscription_status IS NULL OR ss.status = sqlc.narg('status')::subscription_status)
  AND (sqlc.narg('plan_code')::text IS NULL OR ss.plan_code = sqlc.narg('plan_code')::text)
  AND (sqlc.narg('billing_cycle')::billing_cycle IS NULL OR ss.billing_cycle = sqlc.narg('billing_cycle')::billing_cycle);

-- name: CountActiveFlatsForQuota :one
SELECT COUNT(*)
FROM flats
WHERE society_id = $1 AND is_active = TRUE;

-- name: CountActiveAdminsForQuota :one
SELECT COUNT(*)
FROM society_members
WHERE society_id = $1 AND status = 'active' AND role IN ('admin', 'owner');

-- name: CountActiveStaffForQuota :one
SELECT COUNT(*)
FROM society_members
WHERE society_id = $1 AND status = 'active' AND role = 'staff';

-- name: CountActiveResidentsForQuota :one
SELECT COUNT(*)
FROM society_members
WHERE society_id = $1 AND status = 'active' AND role = 'resident';

-- name: GetActiveSubscriptionForUpdate :one
SELECT *
FROM society_subscriptions
WHERE society_id = sqlc.arg('society_id')::bigint
  AND status IN ('trial', 'active')
  AND (ends_at IS NULL OR ends_at > NOW())
ORDER BY id DESC
LIMIT 1
FOR UPDATE;

-- name: ExpireDueSubscriptions :execrows
UPDATE society_subscriptions
SET status = 'expired', expired_at = NOW(), updated_at = NOW()
WHERE status IN ('trial', 'active')
  AND (
      (ends_at IS NOT NULL AND ends_at <= NOW())
      OR (status = 'trial' AND trial_ends_at IS NOT NULL AND trial_ends_at <= NOW())
  );
