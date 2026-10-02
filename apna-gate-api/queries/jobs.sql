-- name: ExpireWaitingVisitorEntriesForJob :execrows
UPDATE visitor_entries
SET status = 'expired', updated_at = NOW()
WHERE status = 'waiting_approval'
  AND created_at < NOW() - INTERVAL '48 hours';

-- name: ExpireApprovedVisitorEntriesForJob :execrows
UPDATE visitor_entries
SET status = 'expired', updated_at = NOW()
WHERE status = 'approved'
  AND checked_in_at IS NULL
  AND (
    (expected_checkout_at IS NOT NULL AND expected_checkout_at < NOW())
    OR (qr_expires_at IS NOT NULL AND qr_expires_at < NOW())
  );

-- name: ExpireVisitorInvitesForJob :execrows
UPDATE visitor_invites
SET status = 'expired', updated_at = NOW()
WHERE status = 'active'
  AND expires_at < NOW();

-- name: ExpireFlatMemberInvitesForJob :execrows
UPDATE flat_member_invites
SET status = 'expired', updated_at = NOW()
WHERE status = 'pending'
  AND expires_at < NOW();

-- name: ExpireSubscriptionsForJob :execrows
UPDATE society_subscriptions
SET status = 'expired', expired_at = NOW(), updated_at = NOW()
WHERE status IN ('trial', 'active')
  AND (
    (ends_at IS NOT NULL AND ends_at <= NOW())
    OR (status = 'trial' AND trial_ends_at IS NOT NULL AND trial_ends_at <= NOW())
  );

-- name: DeleteNotificationsBatch :execrows
WITH candidates AS (
    SELECT n.id
    FROM notifications n
    WHERE n.created_at < sqlc.arg('cutoff')
    ORDER BY n.created_at, n.id
    LIMIT sqlc.arg('batch_size')
    FOR UPDATE SKIP LOCKED
)
DELETE FROM notifications n
USING candidates c
WHERE n.id = c.id;

-- name: DeleteVerificationsBatch :execrows
WITH candidates AS (
    SELECT uv.id
    FROM user_verifications uv
    WHERE uv.is_used = TRUE OR uv.expires_at < NOW()
    ORDER BY uv.expires_at, uv.id
    LIMIT sqlc.arg('batch_size')
    FOR UPDATE SKIP LOCKED
)
DELETE FROM user_verifications uv
USING candidates c
WHERE uv.id = c.id;

-- name: ListVisitorEntryCleanupWindows :many
SELECT DISTINCT d.society_id, d.report_month
FROM monthly_visitor_report_deliveries d
JOIN visitor_entries ve
  ON ve.society_id = d.society_id
 AND ve.created_at >= (d.report_month::timestamp AT TIME ZONE sqlc.arg('job_timezone')::text)
 AND ve.created_at < ((d.report_month + INTERVAL '1 month')::timestamp AT TIME ZONE sqlc.arg('job_timezone')::text)
WHERE d.sent_at IS NOT NULL
  AND ve.status IN ('rejected', 'checked_out', 'cancelled', 'expired', 'auto_closed')
  AND ve.updated_at < sqlc.arg('cutoff')
ORDER BY d.report_month, d.society_id;

-- name: DeleteVisitorEntriesBatch :many
WITH candidates AS (
    SELECT ve.id
    FROM visitor_entries ve
    WHERE ve.society_id = sqlc.arg('society_id')
      AND ve.created_at >= sqlc.arg('period_start')
      AND ve.created_at < sqlc.arg('period_end')
      AND ve.updated_at < sqlc.arg('cutoff')
      AND ve.status IN ('rejected', 'checked_out', 'cancelled', 'expired', 'auto_closed')
      AND EXISTS (
          SELECT 1
          FROM monthly_visitor_report_deliveries d
          WHERE d.society_id = ve.society_id
            AND d.report_month = sqlc.arg('report_month')
            AND d.sent_at IS NOT NULL
      )
    ORDER BY ve.updated_at, ve.id
    LIMIT sqlc.arg('batch_size')
    FOR UPDATE SKIP LOCKED
)
DELETE FROM visitor_entries ve
USING candidates c
WHERE ve.id = c.id
RETURNING ve.visitor_id;

-- name: LockCleanupVisitors :many
SELECT id FROM visitors
WHERE id = ANY(sqlc.arg('ids')::bigint[])
ORDER BY id FOR UPDATE;

-- name: LockOldOrphanVisitors :many
SELECT v.id FROM visitors v
WHERE v.updated_at < sqlc.arg('cutoff')
  AND NOT EXISTS (SELECT 1 FROM visitor_entries ve WHERE ve.visitor_id = v.id)
ORDER BY v.id LIMIT sqlc.arg('batch_size')
FOR UPDATE OF v SKIP LOCKED;

-- name: DeleteUnreferencedVisitors :many
DELETE FROM visitors v
WHERE v.id = ANY(sqlc.arg('ids')::bigint[])
  AND NOT EXISTS (SELECT 1 FROM visitor_entries ve WHERE ve.visitor_id = v.id)
RETURNING v.id, v.photo_url, v.photo_imagekit_file_id;

-- name: QueueVisitorImageDeletion :exec
INSERT INTO visitor_image_deletions (file_id) VALUES ($1)
ON CONFLICT (file_id) DO NOTHING;

-- name: ListVisitorImageDeletions :many
SELECT id, file_id FROM visitor_image_deletions
WHERE id > sqlc.arg('after_id')
ORDER BY id LIMIT sqlc.arg('batch_size');

-- name: CompleteVisitorImageDeletion :exec
DELETE FROM visitor_image_deletions WHERE id = $1;

-- name: ListVisitorInviteCleanupCandidateIDs :many
SELECT vi.id
FROM visitor_invites vi
WHERE vi.status IN ('used', 'expired', 'cancelled')
  AND vi.updated_at < sqlc.arg('cutoff')
  AND NOT EXISTS (
      SELECT 1 FROM visitor_entries ve WHERE ve.invite_id = vi.id
  )
ORDER BY vi.updated_at, vi.id
LIMIT sqlc.arg('batch_size')
FOR UPDATE SKIP LOCKED;

-- name: ListFlatMemberInviteCleanupCandidateIDs :many
SELECT fmi.id
FROM flat_member_invites fmi
WHERE fmi.status IN ('accepted', 'expired', 'cancelled')
  AND fmi.updated_at < sqlc.arg('cutoff')
ORDER BY fmi.updated_at, fmi.id
LIMIT sqlc.arg('batch_size')
FOR UPDATE SKIP LOCKED;

-- name: DeleteShortLinksForResources :execrows
DELETE FROM short_links
WHERE resource_type = sqlc.arg('resource_type')
  AND resource_id = ANY(sqlc.arg('resource_ids')::bigint[]);

-- name: DeleteVisitorInvitesByIDs :execrows
DELETE FROM visitor_invites vi
WHERE vi.id = ANY(sqlc.arg('ids')::bigint[])
  AND NOT EXISTS (
      SELECT 1 FROM visitor_entries ve WHERE ve.invite_id = vi.id
  );

-- name: DeleteFlatMemberInvitesByIDs :execrows
DELETE FROM flat_member_invites
WHERE id = ANY(sqlc.arg('ids')::bigint[]);

-- name: ListActiveSocietiesForVisitorReports :many
SELECT id, name, email
FROM societies
WHERE status = 'active'
  AND deleted_at IS NULL
ORDER BY id;

-- name: GetEarliestVisitorEntryCreatedAt :one
SELECT created_at
FROM visitor_entries
WHERE society_id = $1
ORDER BY created_at, id
LIMIT 1;

-- name: EnsureMonthlyVisitorReportDelivery :exec
INSERT INTO monthly_visitor_report_deliveries (society_id, report_month)
VALUES (sqlc.arg('society_id'), sqlc.arg('report_month'))
ON CONFLICT (society_id, report_month) DO NOTHING;

-- name: ListClaimableMonthlyVisitorReportDeliveries :many
SELECT d.*, s.name AS society_name, s.email AS society_email
FROM monthly_visitor_report_deliveries d
JOIN societies s ON s.id = d.society_id
WHERE (d.sent_at IS NULL
   OR (sqlc.arg('allow_resend')::boolean AND d.report_month = sqlc.arg('resend_month')))
  AND (d.processing_until IS NULL OR d.processing_until <= NOW())
  AND s.status = 'active'
  AND s.deleted_at IS NULL
ORDER BY d.report_month, d.society_id
LIMIT sqlc.arg('batch_size');

-- name: GetMonthlyVisitorReportDeliveryForUpdate :one
SELECT *
FROM monthly_visitor_report_deliveries
WHERE society_id = sqlc.arg('society_id')
  AND report_month = sqlc.arg('report_month')
FOR UPDATE;

-- name: ClaimMonthlyVisitorReportDelivery :one
UPDATE monthly_visitor_report_deliveries
SET status = 'processing',
    attempt_count = attempt_count + 1,
    processing_until = sqlc.arg('processing_until'),
    sent_at = CASE WHEN sqlc.arg('allow_resend')::boolean THEN NULL ELSE sent_at END,
    last_error = NULL,
    updated_at = NOW()
WHERE society_id = sqlc.arg('society_id')
  AND report_month = sqlc.arg('report_month')
  AND (sent_at IS NULL OR sqlc.arg('allow_resend')::boolean)
  AND (processing_until IS NULL OR processing_until <= NOW())
RETURNING *;

-- name: CompleteMonthlyVisitorReportDelivery :execrows
UPDATE monthly_visitor_report_deliveries
SET status = 'sent',
    recipients = sqlc.arg('recipients'),
    provider_message_id = sqlc.arg('provider_message_id'),
    sent_at = NOW(),
    processing_until = NULL,
    last_error = NULL,
    updated_at = NOW()
WHERE society_id = sqlc.arg('society_id')
  AND report_month = sqlc.arg('report_month')
  AND (sent_at IS NULL OR sqlc.arg('allow_resend')::boolean);

-- name: FailMonthlyVisitorReportDelivery :execrows
UPDATE monthly_visitor_report_deliveries
SET status = 'failed',
    recipients = sqlc.arg('recipients'),
    processing_until = NULL,
    last_error = sqlc.arg('last_error'),
    updated_at = NOW()
WHERE society_id = sqlc.arg('society_id')
  AND report_month = sqlc.arg('report_month')
  AND (sent_at IS NULL OR sqlc.arg('allow_resend')::boolean);

-- name: ListMonthlyVisitorReportRecipients :many
SELECT DISTINCT LOWER(TRIM(u.email)) AS email
FROM society_members sm
JOIN users u ON u.id = sm.user_id
WHERE sm.society_id = $1
  AND sm.role = 'owner'
  AND sm.status = 'active'
  AND u.is_active = TRUE
  AND u.is_blocked = FALSE
  AND u.deleted_at IS NULL
  AND u.email IS NOT NULL
  AND TRIM(u.email) <> ''
ORDER BY email;

-- name: ListMonthlyVisitorReportRows :many
SELECT
    ve.id,
    v.full_name AS visitor_name,
    v.phone_number AS visitor_phone,
    v.email AS visitor_email,
    f.flat_number,
    f.block,
    f.floor,
    ve.source,
    ve.purpose,
    ve.status,
    ve.expected_at,
    ve.expected_checkout_at,
    ve.checked_in_at,
    ve.checked_out_at,
    ve.auto_closed_at,
    ve.vehicle_number,
    ve.vehicle_type,
    ve.companions_count,
    ve.companion_details,
    approver.full_name AS approver_name,
    guard_user.full_name AS guard_name,
    ve.rejection_reason,
    ve.notes,
    ve.created_at,
    ve.updated_at
FROM visitor_entries ve
JOIN visitors v ON v.id = ve.visitor_id
LEFT JOIN flats f
  ON f.id = ve.flat_id
 AND f.society_id = ve.society_id
LEFT JOIN users approver ON approver.id = ve.approved_by
LEFT JOIN users guard_user ON guard_user.id = ve.handled_by_guard_id
WHERE ve.society_id = sqlc.arg('society_id')
  AND ve.created_at >= sqlc.arg('period_start')
  AND ve.created_at < sqlc.arg('period_end')
ORDER BY ve.created_at, ve.id;
