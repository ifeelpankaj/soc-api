-- name: CreateVisitor :one
INSERT INTO visitors (full_name, phone_number, email, photo_url, metadata)
VALUES ($1, $2, $3, $4, COALESCE(sqlc.narg('metadata'), '{}'::jsonb))
RETURNING *;

-- name: GetVisitor :one
SELECT *
FROM visitors
WHERE id = $1;

-- name: CreateVisitorInvite :one
INSERT INTO visitor_invites (society_id, flat_id, created_by, purpose, token_hash, expires_at, metadata)
VALUES ($1, $2, $3, $4, $5, $6, COALESCE(sqlc.narg('metadata'), '{}'::jsonb))
RETURNING *;

-- name: GetVisitorInviteByID :one
SELECT *
FROM visitor_invites
WHERE id = $1
  AND society_id = $2;

-- name: ListVisitorInviteHistory :many
SELECT
    vi.id,
    vi.society_id,
    vi.flat_id,
    vi.created_by,
    vi.purpose,
    vi.status,
    vi.expires_at,
    vi.used_at,
    vi.metadata,
    vi.created_at,
    vi.updated_at,
    sl.short_code,
    sl.expires_at AS short_link_expires_at,
    ve.id AS entry_id,
    ve.status AS entry_status,
    ve.checked_in_at,
    ve.checked_out_at,
    ve.qr_expires_at,
    ve.qr_used_at,
    v.full_name AS visitor_full_name,
    v.phone_number AS visitor_phone_number,
    v.email AS visitor_email,
    v.photo_url AS visitor_photo_url
FROM visitor_invites vi
LEFT JOIN short_links sl
    ON sl.resource_type = 'visitor_invite'
   AND sl.resource_id = vi.id
   AND sl.revoked_at IS NULL
LEFT JOIN visitor_entries ve
    ON ve.invite_id = vi.id
LEFT JOIN visitors v
    ON v.id = ve.visitor_id
WHERE vi.society_id = $1
  AND vi.flat_id = $2
  AND (sqlc.narg('status')::visitor_invite_status IS NULL OR vi.status = sqlc.narg('status')::visitor_invite_status)
  AND (
    sqlc.narg('search')::text IS NULL
    OR vi.purpose::text ILIKE '%' || sqlc.narg('search')::text || '%'
    OR vi.status::text ILIKE '%' || sqlc.narg('search')::text || '%'
    OR v.full_name ILIKE '%' || sqlc.narg('search')::text || '%'
    OR v.phone_number ILIKE '%' || sqlc.narg('search')::text || '%'
    OR v.email ILIKE '%' || sqlc.narg('search')::text || '%'
    OR sl.short_code ILIKE '%' || sqlc.narg('search')::text || '%'
  )
ORDER BY vi.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountVisitorInviteHistory :one
SELECT COUNT(*)
FROM visitor_invites
WHERE visitor_invites.society_id = $1
  AND visitor_invites.flat_id = $2
  AND (sqlc.narg('status')::visitor_invite_status IS NULL OR visitor_invites.status = sqlc.narg('status')::visitor_invite_status)
  AND (
    sqlc.narg('search')::text IS NULL
    OR visitor_invites.purpose::text ILIKE '%' || sqlc.narg('search')::text || '%'
    OR visitor_invites.status::text ILIKE '%' || sqlc.narg('search')::text || '%'
    OR EXISTS (
      SELECT 1
      FROM visitor_entries ve
      JOIN visitors v ON v.id = ve.visitor_id
      WHERE ve.invite_id = visitor_invites.id
        AND (
          v.full_name ILIKE '%' || sqlc.narg('search')::text || '%'
          OR v.phone_number ILIKE '%' || sqlc.narg('search')::text || '%'
          OR v.email ILIKE '%' || sqlc.narg('search')::text || '%'
        )
    )
    OR EXISTS (
      SELECT 1
      FROM short_links sl
      WHERE sl.resource_type = 'visitor_invite'
        AND sl.resource_id = visitor_invites.id
        AND sl.revoked_at IS NULL
        AND sl.short_code ILIKE '%' || sqlc.narg('search')::text || '%'
    )
  );

-- name: GetVisitorInviteHistoryByID :one
SELECT
    vi.id,
    vi.society_id,
    vi.flat_id,
    vi.created_by,
    vi.purpose,
    vi.status,
    vi.expires_at,
    vi.used_at,
    vi.metadata,
    vi.created_at,
    vi.updated_at,
    sl.short_code,
    sl.expires_at AS short_link_expires_at,
    ve.id AS entry_id,
    ve.status AS entry_status,
    ve.checked_in_at,
    ve.checked_out_at,
    ve.qr_expires_at,
    ve.qr_used_at,
    v.full_name AS visitor_full_name,
    v.phone_number AS visitor_phone_number,
    v.email AS visitor_email,
    v.photo_url AS visitor_photo_url
FROM visitor_invites vi
LEFT JOIN short_links sl
    ON sl.resource_type = 'visitor_invite'
   AND sl.resource_id = vi.id
   AND sl.revoked_at IS NULL
LEFT JOIN visitor_entries ve
    ON ve.invite_id = vi.id
LEFT JOIN visitors v
    ON v.id = ve.visitor_id
WHERE vi.society_id = $1
  AND vi.flat_id = $2
  AND vi.id = $3;

-- name: GetVisitorInviteByTokenHash :one
SELECT *
FROM visitor_invites
WHERE token_hash = $1;

-- name: MarkVisitorInviteUsed :one
UPDATE visitor_invites
SET status = 'used',
    used_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND status = 'active'
RETURNING *;

-- name: GetVisitorInviteForUpdate :one
SELECT *
FROM visitor_invites
WHERE id = $1
FOR UPDATE;

-- name: CancelVisitorInvite :one
UPDATE visitor_invites
SET status = 'cancelled',
    updated_at = NOW()
WHERE id = $1
  AND society_id = $2
  AND status = 'active'
RETURNING *;

-- name: ExpireOldVisitorInvites :exec
UPDATE visitor_invites
SET status = 'expired',
    updated_at = NOW()
WHERE status = 'active'
  AND expires_at < NOW();

-- name: CreateVisitorEntry :one
INSERT INTO visitor_entries (
    society_id,
    flat_id,
    visitor_id,
    invite_id,
    source,
    purpose,
    status,
    vehicle_number,
    vehicle_type,
    companions_count,
    companion_details,
    expected_at,
    expected_checkout_at,
    approved_by,
    approved_at,
    delivery_partner,
    service_provider,
    handled_by_guard_id,
    created_by,
    qr_token_hash,
    qr_expires_at,
    notes,
    metadata
)
VALUES (
    sqlc.arg('society_id'),
    sqlc.narg('flat_id'),
    sqlc.arg('visitor_id'),
    sqlc.narg('invite_id'),
    sqlc.arg('source'),
    sqlc.arg('purpose'),
    sqlc.arg('status'),
    sqlc.narg('vehicle_number'),
    sqlc.narg('vehicle_type'),
    sqlc.arg('companions_count'),
    COALESCE(sqlc.narg('companion_details'), '[]'::jsonb),
    sqlc.narg('expected_at'),
    sqlc.narg('expected_checkout_at'),
    sqlc.narg('approved_by'),
    sqlc.narg('approved_at'),
    sqlc.narg('delivery_partner'),
    sqlc.narg('service_provider'),
    sqlc.narg('handled_by_guard_id'),
    sqlc.narg('created_by'),
    sqlc.narg('qr_token_hash'),
    sqlc.narg('qr_expires_at'),
    sqlc.narg('notes'),
    COALESCE(sqlc.narg('metadata'), '{}'::jsonb)
)
RETURNING *;

-- name: GetVisitorEntry :one
SELECT
    ve.*,

    v.full_name AS visitor_full_name,
    v.phone_number AS visitor_phone_number,
    v.email AS visitor_email,
    v.photo_url AS visitor_photo_url,

    f.flat_number,
    f.block,
    f.floor,

    approved_user.full_name AS approved_by_full_name,
    approved_user.phone_number AS approved_by_phone_number,

    guard_user.full_name AS handled_by_guard_full_name,
    guard_user.phone_number AS handled_by_guard_phone_number

FROM visitor_entries ve

JOIN visitors v
    ON v.id = ve.visitor_id

LEFT JOIN flats f
    ON f.id = ve.flat_id

LEFT JOIN users approved_user
    ON approved_user.id = ve.approved_by

LEFT JOIN users guard_user
    ON guard_user.id = ve.handled_by_guard_id

WHERE ve.id = $1
  AND ve.society_id = $2;


-- name: GetVisitorEntryForUpdate :one
SELECT ve.*
FROM visitor_entries ve
WHERE ve.id = $1
  AND ve.society_id = $2
FOR UPDATE;

-- name: GetVisitorEntryByQRHash :one
SELECT
    ve.*,
    v.full_name AS visitor_full_name,
    v.phone_number AS visitor_phone_number,
    v.email AS visitor_email,
    v.photo_url AS visitor_photo_url,
    f.flat_number,
    f.block,
    f.floor,
    approved_user.full_name AS approved_by_full_name,
    approved_user.phone_number AS approved_by_phone_number,
    guard_user.full_name AS handled_by_guard_full_name,
    guard_user.phone_number AS handled_by_guard_phone_number
FROM visitor_entries ve
JOIN visitors v ON v.id = ve.visitor_id
LEFT JOIN flats f ON f.id = ve.flat_id
LEFT JOIN users approved_user ON approved_user.id = ve.approved_by
LEFT JOIN users guard_user ON guard_user.id = ve.handled_by_guard_id
WHERE ve.qr_token_hash = $1;

-- name: GetVisitorEntryByInviteID :one
SELECT
    ve.*,
    v.full_name AS visitor_full_name,
    v.phone_number AS visitor_phone_number,
    v.email AS visitor_email,
    v.photo_url AS visitor_photo_url,
    f.flat_number,
    f.block,
    f.floor,
    approved_user.full_name AS approved_by_full_name,
    approved_user.phone_number AS approved_by_phone_number,
    guard_user.full_name AS handled_by_guard_full_name,
    guard_user.phone_number AS handled_by_guard_phone_number
FROM visitor_entries ve
JOIN visitors v ON v.id = ve.visitor_id
LEFT JOIN flats f ON f.id = ve.flat_id
LEFT JOIN users approved_user ON approved_user.id = ve.approved_by
LEFT JOIN users guard_user ON guard_user.id = ve.handled_by_guard_id
WHERE ve.invite_id = $1;

-- name: ListVisitorEntries :many
SELECT
    ve.*,
    v.full_name AS visitor_full_name,
    v.phone_number AS visitor_phone_number,
    v.email AS visitor_email,
    v.photo_url AS visitor_photo_url,
    f.flat_number,
    f.block,
    f.floor,
    approved_user.full_name AS approved_by_full_name,
    approved_user.phone_number AS approved_by_phone_number,
    guard_user.full_name AS handled_by_guard_full_name,
    guard_user.phone_number AS handled_by_guard_phone_number
FROM visitor_entries ve
JOIN visitors v ON v.id = ve.visitor_id
LEFT JOIN flats f ON f.id = ve.flat_id
LEFT JOIN users approved_user ON approved_user.id = ve.approved_by
LEFT JOIN users guard_user ON guard_user.id = ve.handled_by_guard_id
WHERE ve.society_id = sqlc.arg('society_id')
  AND (sqlc.narg('flat_id')::bigint IS NULL OR ve.flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('status')::visitor_status IS NULL OR ve.status = sqlc.narg('status')::visitor_status)
  AND (sqlc.narg('source')::visitor_source IS NULL OR ve.source = sqlc.narg('source')::visitor_source)
  AND (sqlc.narg('purpose')::visitor_purpose IS NULL OR ve.purpose = sqlc.narg('purpose')::visitor_purpose)
  AND (sqlc.narg('block')::text IS NULL OR f.block = sqlc.narg('block')::text)
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR ve.created_at >= sqlc.narg('created_from')::timestamptz)
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR ve.created_at < sqlc.narg('created_to')::timestamptz)
  AND (
      sqlc.narg('event')::text IS NULL
      OR sqlc.narg('event_from')::timestamptz IS NULL
      OR sqlc.narg('event_to')::timestamptz IS NULL
      OR (
          sqlc.narg('event')::text = 'created'
          AND ve.created_at >= sqlc.narg('event_from')::timestamptz
          AND ve.created_at < sqlc.narg('event_to')::timestamptz
      )
      OR (
          sqlc.narg('event')::text = 'checked_in'
          AND ve.checked_in_at IS NOT NULL
          AND ve.checked_in_at >= sqlc.narg('event_from')::timestamptz
          AND ve.checked_in_at < sqlc.narg('event_to')::timestamptz
      )
      OR (
          sqlc.narg('event')::text = 'checked_out'
          AND ve.checked_out_at IS NOT NULL
          AND ve.checked_out_at >= sqlc.narg('event_from')::timestamptz
          AND ve.checked_out_at < sqlc.narg('event_to')::timestamptz
      )
      OR (
          sqlc.narg('event')::text = 'expected'
          AND ve.expected_at IS NOT NULL
          AND ve.expected_at >= sqlc.narg('event_from')::timestamptz
          AND ve.expected_at < sqlc.narg('event_to')::timestamptz
      )
      OR (
          sqlc.narg('event')::text = 'activity'
          AND (
              (ve.created_at >= sqlc.narg('event_from')::timestamptz AND ve.created_at < sqlc.narg('event_to')::timestamptz)
              OR (
                  ve.checked_in_at IS NOT NULL
                  AND ve.checked_in_at >= sqlc.narg('event_from')::timestamptz
                  AND ve.checked_in_at < sqlc.narg('event_to')::timestamptz
              )
              OR (
                  ve.checked_out_at IS NOT NULL
                  AND ve.checked_out_at >= sqlc.narg('event_from')::timestamptz
                  AND ve.checked_out_at < sqlc.narg('event_to')::timestamptz
              )
          )
      )
  )
  AND (
      COALESCE(sqlc.narg('search')::text, '') = ''
      OR v.full_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(v.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(f.flat_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(f.block, '') ILIKE '%' || sqlc.narg('search')::text || '%'
  )
ORDER BY ve.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountVisitorEntries :one
SELECT COUNT(*)
FROM visitor_entries ve
LEFT JOIN flats f ON f.id = ve.flat_id
WHERE ve.society_id = sqlc.arg('society_id')
  AND (sqlc.narg('flat_id')::bigint IS NULL OR ve.flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('status')::visitor_status IS NULL OR ve.status = sqlc.narg('status')::visitor_status)
  AND (sqlc.narg('source')::visitor_source IS NULL OR ve.source = sqlc.narg('source')::visitor_source)
  AND (sqlc.narg('purpose')::visitor_purpose IS NULL OR ve.purpose = sqlc.narg('purpose')::visitor_purpose)
  AND (sqlc.narg('block')::text IS NULL OR f.block = sqlc.narg('block')::text)
  AND (sqlc.narg('created_from')::timestamptz IS NULL OR ve.created_at >= sqlc.narg('created_from')::timestamptz)
  AND (sqlc.narg('created_to')::timestamptz IS NULL OR ve.created_at < sqlc.narg('created_to')::timestamptz)
  AND (
      sqlc.narg('event')::text IS NULL
      OR sqlc.narg('event_from')::timestamptz IS NULL
      OR sqlc.narg('event_to')::timestamptz IS NULL
      OR (
          sqlc.narg('event')::text = 'created'
          AND ve.created_at >= sqlc.narg('event_from')::timestamptz
          AND ve.created_at < sqlc.narg('event_to')::timestamptz
      )
      OR (
          sqlc.narg('event')::text = 'checked_in'
          AND ve.checked_in_at IS NOT NULL
          AND ve.checked_in_at >= sqlc.narg('event_from')::timestamptz
          AND ve.checked_in_at < sqlc.narg('event_to')::timestamptz
      )
      OR (
          sqlc.narg('event')::text = 'checked_out'
          AND ve.checked_out_at IS NOT NULL
          AND ve.checked_out_at >= sqlc.narg('event_from')::timestamptz
          AND ve.checked_out_at < sqlc.narg('event_to')::timestamptz
      )
      OR (
          sqlc.narg('event')::text = 'expected'
          AND ve.expected_at IS NOT NULL
          AND ve.expected_at >= sqlc.narg('event_from')::timestamptz
          AND ve.expected_at < sqlc.narg('event_to')::timestamptz
      )
      OR (
          sqlc.narg('event')::text = 'activity'
          AND (
              (ve.created_at >= sqlc.narg('event_from')::timestamptz AND ve.created_at < sqlc.narg('event_to')::timestamptz)
              OR (
                  ve.checked_in_at IS NOT NULL
                  AND ve.checked_in_at >= sqlc.narg('event_from')::timestamptz
                  AND ve.checked_in_at < sqlc.narg('event_to')::timestamptz
              )
              OR (
                  ve.checked_out_at IS NOT NULL
                  AND ve.checked_out_at >= sqlc.narg('event_from')::timestamptz
                  AND ve.checked_out_at < sqlc.narg('event_to')::timestamptz
              )
          )
      )
  )
  AND (
      COALESCE(sqlc.narg('search')::text, '') = ''
      OR EXISTS (
          SELECT 1
          FROM visitors v
          WHERE v.id = ve.visitor_id
            AND (
                v.full_name ILIKE '%' || sqlc.narg('search')::text || '%'
                OR COALESCE(v.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
            )
      )
      OR EXISTS (
          SELECT 1
          FROM flats f
          WHERE f.id = ve.flat_id
            AND (
                COALESCE(f.flat_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
                OR COALESCE(f.block, '') ILIKE '%' || sqlc.narg('search')::text || '%'
            )
      )
  );

-- name: GetVisitorEntryStats :one
SELECT
    COUNT(*) FILTER (
        WHERE ve.created_at >= (CURRENT_DATE AT TIME ZONE 'Asia/Kolkata') AT TIME ZONE 'Asia/Kolkata'
          AND ve.created_at < ((CURRENT_DATE + INTERVAL '1 day') AT TIME ZONE 'Asia/Kolkata') AT TIME ZONE 'Asia/Kolkata'
    )::bigint AS today_visitors,
    COUNT(*) FILTER (WHERE ve.status = 'checked_in')::bigint AS visitors_inside,
    COUNT(*) FILTER (WHERE ve.status = 'waiting_approval')::bigint AS pending_approvals,
    COUNT(*) FILTER (
        WHERE ve.checked_out_at IS NOT NULL
          AND ve.status = 'checked_out'
          AND ve.checked_out_at >= (CURRENT_DATE AT TIME ZONE 'Asia/Kolkata') AT TIME ZONE 'Asia/Kolkata'
          AND ve.checked_out_at < ((CURRENT_DATE + INTERVAL '1 day') AT TIME ZONE 'Asia/Kolkata') AT TIME ZONE 'Asia/Kolkata'
    )::bigint AS checked_out_today,
    COUNT(*) FILTER (
        WHERE ve.status = 'rejected'
          AND ve.updated_at >= (CURRENT_DATE AT TIME ZONE 'Asia/Kolkata') AT TIME ZONE 'Asia/Kolkata'
          AND ve.updated_at < ((CURRENT_DATE + INTERVAL '1 day') AT TIME ZONE 'Asia/Kolkata') AT TIME ZONE 'Asia/Kolkata'
    )::bigint AS rejected_today,
    COUNT(*) FILTER (
        WHERE ve.auto_closed_at IS NOT NULL
          AND ve.auto_closed_at >= (CURRENT_DATE AT TIME ZONE 'Asia/Kolkata') AT TIME ZONE 'Asia/Kolkata'
          AND ve.auto_closed_at < ((CURRENT_DATE + INTERVAL '1 day') AT TIME ZONE 'Asia/Kolkata') AT TIME ZONE 'Asia/Kolkata'
    )::bigint AS auto_closed_today
FROM visitor_entries ve
WHERE ve.society_id = $1;

-- name: GetVisitorEntryStatsInRange :one
SELECT
    COUNT(*) FILTER (WHERE ve.status = 'checked_in')::bigint AS visitors_inside,
    COUNT(*) FILTER (WHERE ve.status = 'waiting_approval')::bigint AS pending_approvals,
    COUNT(*) FILTER (
        WHERE ve.checked_out_at IS NOT NULL
          AND ve.status = 'checked_out'
          AND ve.checked_out_at >= sqlc.arg('event_from')::timestamptz
          AND ve.checked_out_at < sqlc.arg('event_to')::timestamptz
    )::bigint AS checked_out_in_range
FROM visitor_entries ve
WHERE ve.society_id = sqlc.arg('society_id');

-- name: GetVisitorEntryDailyStatsCreated :many
SELECT
    to_char((ve.created_at AT TIME ZONE 'Asia/Kolkata')::date, 'YYYY-MM-DD') AS stat_date,
    COUNT(*)::bigint AS count
FROM visitor_entries ve
WHERE ve.society_id = sqlc.arg('society_id')
  AND ve.created_at >= ((CURRENT_DATE AT TIME ZONE 'Asia/Kolkata')::date - (sqlc.arg('days')::int - 1)) AT TIME ZONE 'Asia/Kolkata'
  AND ve.created_at < ((CURRENT_DATE + INTERVAL '1 day') AT TIME ZONE 'Asia/Kolkata') AT TIME ZONE 'Asia/Kolkata'
GROUP BY (ve.created_at AT TIME ZONE 'Asia/Kolkata')::date
ORDER BY stat_date ASC;

-- name: CountWaitingAtGateVisitorEntries :one
SELECT COUNT(*)::bigint
FROM visitor_entries ve
WHERE ve.society_id = sqlc.arg('society_id')
  AND ve.status = 'approved'
  AND ve.checked_in_at IS NULL
  AND ve.source <> 'resident_link';

-- name: ListWaitingAtGateVisitorEntries :many
SELECT
    ve.*,
    v.full_name AS visitor_full_name,
    v.phone_number AS visitor_phone_number,
    v.email AS visitor_email,
    v.photo_url AS visitor_photo_url,
    f.flat_number,
    f.block,
    f.floor,
    approved_user.full_name AS approved_by_full_name,
    approved_user.phone_number AS approved_by_phone_number,
    guard_user.full_name AS handled_by_guard_full_name,
    guard_user.phone_number AS handled_by_guard_phone_number
FROM visitor_entries ve
JOIN visitors v ON v.id = ve.visitor_id
LEFT JOIN flats f ON f.id = ve.flat_id
LEFT JOIN users approved_user ON approved_user.id = ve.approved_by
LEFT JOIN users guard_user ON guard_user.id = ve.handled_by_guard_id
WHERE ve.society_id = sqlc.arg('society_id')
  AND ve.status = 'approved'
  AND ve.checked_in_at IS NULL
  AND ve.source <> 'resident_link'
  AND (
      COALESCE(sqlc.narg('search')::text, '') = ''
      OR v.full_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(v.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(f.flat_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(f.block, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(ve.vehicle_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(ve.delivery_partner, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR ve.purpose::text ILIKE '%' || sqlc.narg('search')::text || '%'
  )
ORDER BY ve.approved_at ASC NULLS LAST, ve.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountWaitingAtGateVisitorEntriesFiltered :one
SELECT COUNT(*)::bigint
FROM visitor_entries ve
JOIN visitors v ON v.id = ve.visitor_id
LEFT JOIN flats f ON f.id = ve.flat_id
WHERE ve.society_id = sqlc.arg('society_id')
  AND ve.status = 'approved'
  AND ve.checked_in_at IS NULL
  AND ve.source <> 'resident_link'
  AND (
      COALESCE(sqlc.narg('search')::text, '') = ''
      OR v.full_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(v.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(f.flat_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(f.block, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(ve.vehicle_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(ve.delivery_partner, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR ve.purpose::text ILIKE '%' || sqlc.narg('search')::text || '%'
  );

-- name: CountExpectedGuestEntries :one
SELECT COUNT(*)::bigint
FROM visitor_entries ve
JOIN visitors v ON v.id = ve.visitor_id
WHERE ve.society_id = sqlc.arg('society_id')
  AND ve.source = 'resident_link'
  AND ve.status = 'approved'
  AND ve.checked_in_at IS NULL
  AND COALESCE(ve.expected_at, ve.approved_at, ve.created_at) >= sqlc.arg('from_at')::timestamptz
  AND COALESCE(ve.expected_at, ve.approved_at, ve.created_at) < sqlc.arg('to_at')::timestamptz;

-- name: ListExpectedGuestEntries :many
SELECT
    ve.*,
    v.full_name AS visitor_full_name,
    v.phone_number AS visitor_phone_number,
    v.email AS visitor_email,
    v.photo_url AS visitor_photo_url,
    f.flat_number,
    f.block,
    f.floor,
    approved_user.full_name AS approved_by_full_name,
    approved_user.phone_number AS approved_by_phone_number,
    guard_user.full_name AS handled_by_guard_full_name,
    guard_user.phone_number AS handled_by_guard_phone_number
FROM visitor_entries ve
JOIN visitors v ON v.id = ve.visitor_id
LEFT JOIN flats f ON f.id = ve.flat_id
LEFT JOIN users approved_user ON approved_user.id = ve.approved_by
LEFT JOIN users guard_user ON guard_user.id = ve.handled_by_guard_id
WHERE ve.society_id = sqlc.arg('society_id')
  AND ve.source = 'resident_link'
  AND ve.status = 'approved'
  AND ve.checked_in_at IS NULL
  AND COALESCE(ve.expected_at, ve.approved_at, ve.created_at) >= sqlc.arg('from_at')::timestamptz
  AND COALESCE(ve.expected_at, ve.approved_at, ve.created_at) < sqlc.arg('to_at')::timestamptz
  AND (
      COALESCE(sqlc.narg('search')::text, '') = ''
      OR v.full_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(v.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(f.flat_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(f.block, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(ve.vehicle_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR ve.purpose::text ILIKE '%' || sqlc.narg('search')::text || '%'
  )
ORDER BY COALESCE(ve.expected_at, ve.approved_at, ve.created_at) ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountExpectedGuestEntriesFiltered :one
SELECT COUNT(*)::bigint
FROM visitor_entries ve
JOIN visitors v ON v.id = ve.visitor_id
LEFT JOIN flats f ON f.id = ve.flat_id
WHERE ve.society_id = sqlc.arg('society_id')
  AND ve.source = 'resident_link'
  AND ve.status = 'approved'
  AND ve.checked_in_at IS NULL
  AND COALESCE(ve.expected_at, ve.approved_at, ve.created_at) >= sqlc.arg('from_at')::timestamptz
  AND COALESCE(ve.expected_at, ve.approved_at, ve.created_at) < sqlc.arg('to_at')::timestamptz
  AND (
      COALESCE(sqlc.narg('search')::text, '') = ''
      OR v.full_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(v.phone_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(f.flat_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(f.block, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR COALESCE(ve.vehicle_number, '') ILIKE '%' || sqlc.narg('search')::text || '%'
      OR ve.purpose::text ILIKE '%' || sqlc.narg('search')::text || '%'
  );

-- name: ListSocietyPendingVisitorApprovals :many
SELECT
    ve.*,
    v.full_name AS visitor_full_name,
    v.phone_number AS visitor_phone_number,
    v.email AS visitor_email,
    v.photo_url AS visitor_photo_url,
    f.flat_number,
    f.block,
    f.floor,
    approved_user.full_name AS approved_by_full_name,
    approved_user.phone_number AS approved_by_phone_number,
    guard_user.full_name AS handled_by_guard_full_name,
    guard_user.phone_number AS handled_by_guard_phone_number,
    u.full_name AS primary_resident_name,
    fr.id AS primary_resident_id
FROM visitor_entries ve
JOIN visitors v ON v.id = ve.visitor_id
LEFT JOIN flats f ON f.id = ve.flat_id
LEFT JOIN users approved_user ON approved_user.id = ve.approved_by
LEFT JOIN users guard_user ON guard_user.id = ve.handled_by_guard_id
LEFT JOIN flat_residents fr
    ON fr.flat_id = ve.flat_id
   AND fr.society_id = ve.society_id
   AND fr.status = 'active'
   AND fr.is_primary = TRUE
LEFT JOIN users u ON u.id = fr.user_id
WHERE ve.society_id = sqlc.arg('society_id')
  AND ve.status = 'waiting_approval'
  AND (sqlc.narg('flat_id')::bigint IS NULL OR ve.flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('block')::text IS NULL OR f.block = sqlc.narg('block')::text)
ORDER BY ve.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSocietyPendingVisitorApprovals :one
SELECT COUNT(*)
FROM visitor_entries ve
LEFT JOIN flats f ON f.id = ve.flat_id
WHERE ve.society_id = sqlc.arg('society_id')
  AND ve.status = 'waiting_approval'
  AND (sqlc.narg('flat_id')::bigint IS NULL OR ve.flat_id = sqlc.narg('flat_id')::bigint)
  AND (sqlc.narg('block')::text IS NULL OR f.block = sqlc.narg('block')::text);

-- name: CountMemberVisitorApprovals :one
SELECT
    COUNT(*) FILTER (WHERE ve.approved_by = $2)::bigint AS approved_count,
    COUNT(*) FILTER (WHERE ve.rejected_by = $2)::bigint AS rejected_count
FROM visitor_entries ve
WHERE ve.society_id = $1
  AND (ve.approved_by = $2 OR ve.rejected_by = $2);

-- name: ListRecentVisitorEntriesByFlat :many
SELECT
    ve.*,
    v.full_name AS visitor_full_name,
    v.phone_number AS visitor_phone_number,
    v.email AS visitor_email,
    v.photo_url AS visitor_photo_url,
    f.flat_number,
    f.block,
    f.floor,
    approved_user.full_name AS approved_by_full_name,
    approved_user.phone_number AS approved_by_phone_number,
    guard_user.full_name AS handled_by_guard_full_name,
    guard_user.phone_number AS handled_by_guard_phone_number
FROM visitor_entries ve
JOIN visitors v ON v.id = ve.visitor_id
LEFT JOIN flats f ON f.id = ve.flat_id
LEFT JOIN users approved_user ON approved_user.id = ve.approved_by
LEFT JOIN users guard_user ON guard_user.id = ve.handled_by_guard_id
WHERE ve.society_id = $1
  AND ve.flat_id = $2
ORDER BY ve.created_at DESC
LIMIT $3;

-- name: ListPendingVisitorApprovals :many
SELECT
    ve.*,
    v.full_name AS visitor_full_name,
    v.phone_number AS visitor_phone_number,
    v.email AS visitor_email,
    v.photo_url AS visitor_photo_url,
    f.flat_number,
    f.block,
    f.floor,
    approved_user.full_name AS approved_by_full_name,
    approved_user.phone_number AS approved_by_phone_number,
    guard_user.full_name AS handled_by_guard_full_name,
    guard_user.phone_number AS handled_by_guard_phone_number
FROM visitor_entries ve
JOIN visitors v ON v.id = ve.visitor_id
LEFT JOIN flats f ON f.id = ve.flat_id
LEFT JOIN users approved_user ON approved_user.id = ve.approved_by
LEFT JOIN users guard_user ON guard_user.id = ve.handled_by_guard_id
WHERE ve.society_id = $1
  AND ve.flat_id = $2
  AND ve.status = 'waiting_approval'
ORDER BY ve.created_at DESC
LIMIT $3;

-- name: ApproveVisitorEntry :one
UPDATE visitor_entries
SET status = 'approved',
    approved_by = $3,
    approved_at = NOW(),
    qr_token_hash = $4,
    qr_expires_at = $5,
    updated_at = NOW()
WHERE id = $1
  AND society_id = $2
  AND status = 'waiting_approval'
RETURNING *;

-- name: MergeVisitorEntryMetadata :one
UPDATE visitor_entries
SET metadata = COALESCE(metadata, '{}'::jsonb) || sqlc.arg('metadata_patch')::jsonb,
    updated_at = NOW()
WHERE id = sqlc.arg('id')
  AND society_id = sqlc.arg('society_id')
RETURNING *;

-- name: RejectVisitorEntry :one
UPDATE visitor_entries
SET status = 'rejected',
    rejected_by = $3,
    rejection_reason = $4,
    updated_at = NOW()
WHERE id = $1
  AND society_id = $2
  AND status = 'waiting_approval'
RETURNING *;

-- name: GenerateVisitorEntryQR :one
UPDATE visitor_entries
SET qr_token_hash = $3,
    qr_expires_at = $4,
    updated_at = NOW()
WHERE id = $1
  AND society_id = $2
  AND status = 'approved'
RETURNING *;

-- name: CheckInVisitorEntry :one
UPDATE visitor_entries
SET status = 'checked_in',
    checked_in_at = NOW(),
    handled_by_guard_id = $3,
    qr_used_at = COALESCE(qr_used_at, NOW()),
    updated_at = NOW()
WHERE id = $1
  AND society_id = $2
  AND status = 'approved'
RETURNING *;

-- name: CheckOutVisitorEntry :one
UPDATE visitor_entries
SET status = 'checked_out',
    checked_out_at = NOW(),
    handled_by_guard_id = $3,
    updated_at = NOW()
WHERE id = $1
  AND society_id = $2
  AND status = 'checked_in'
RETURNING *;

-- name: AutoCloseExpiredVisitorEntries :exec
-- Checked-in visitors remain inside until explicit guard checkout.
UPDATE visitor_entries
SET updated_at = updated_at
WHERE false;

-- name: ExpireStaleWaitingApprovalEntries :exec
UPDATE visitor_entries
SET status = 'expired',
    updated_at = NOW()
WHERE status = 'waiting_approval'
  AND created_at < NOW() - INTERVAL '48 hours';

-- name: ExpireStaleApprovedEntries :exec
UPDATE visitor_entries
SET status = 'expired',
    updated_at = NOW()
WHERE status = 'approved'
  AND checked_in_at IS NULL
  AND (
    (expected_checkout_at IS NOT NULL AND expected_checkout_at < NOW())
    OR (expected_checkout_at IS NULL AND qr_expires_at IS NOT NULL AND qr_expires_at < NOW())
  );

-- name: UpdateVisitorProfile :one
UPDATE visitors
SET
    full_name = COALESCE(sqlc.narg('full_name'), full_name),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    email = COALESCE(sqlc.narg('email'), email),
    photo_url = COALESCE(sqlc.narg('photo_url'), photo_url),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateVisitorEntryDetails :one
UPDATE visitor_entries
SET
    flat_id = COALESCE(sqlc.narg('flat_id'), flat_id),
    vehicle_number = COALESCE(sqlc.narg('vehicle_number'), vehicle_number),
    vehicle_type = COALESCE(sqlc.narg('vehicle_type'), vehicle_type),
    companions_count = COALESCE(sqlc.narg('companions_count'), companions_count),
    companion_details = COALESCE(sqlc.narg('companion_details'), companion_details),
    notes = COALESCE(sqlc.narg('notes'), notes),
    updated_at = NOW()
WHERE id = $1
  AND society_id = $2
  AND status IN ('waiting_approval', 'approved')
RETURNING *;

-- name: CreateVisitorEntryEvent :one
INSERT INTO visitor_entry_events (visitor_entry_id, society_id, actor_user_id, event_type, message, metadata)
VALUES ($1, $2, sqlc.narg('actor_user_id'), $3, sqlc.narg('message'), COALESCE(sqlc.narg('metadata'), '{}'::jsonb))
RETURNING *;

-- name: ListVisitorEntryEvents :many
SELECT *
FROM visitor_entry_events
WHERE visitor_entry_id = $1
  AND society_id = $2
ORDER BY created_at DESC;
