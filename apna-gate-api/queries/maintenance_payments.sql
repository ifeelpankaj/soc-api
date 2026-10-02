-- name: LockMaintenancePaymentBill :one
SELECT * FROM maintenance_bills WHERE society_id=$1 AND id=$2 FOR UPDATE;

-- name: MaintenanceResidentAccess :one
SELECT EXISTS(SELECT 1 FROM maintenance_bills b JOIN flats f ON f.id=b.flat_id AND f.society_id=b.society_id AND f.is_active
JOIN flat_residents r ON r.flat_id=f.id AND r.society_id=f.society_id AND r.status='active'
JOIN society_members m ON m.user_id=r.user_id AND m.society_id=r.society_id AND m.status='active'
JOIN societies s ON s.id=b.society_id AND s.status='active' AND s.deleted_at IS NULL
WHERE b.society_id=$1 AND b.id=$2 AND r.user_id=$3)::boolean;

-- name: GetUPISettings :one
SELECT v.* FROM maintenance_payment_settings s JOIN maintenance_payment_settings_versions v USING(society_id,version) WHERE s.society_id=$1;
-- name: GetUPISettingsVersion :one
SELECT * FROM maintenance_payment_settings_versions WHERE society_id=$1 AND version=$2;
-- name: ListUPISettingsVersions :many
SELECT * FROM maintenance_payment_settings_versions WHERE society_id=$1 AND (sqlc.arg('before_version')::bigint=0 OR version<sqlc.arg('before_version')) ORDER BY version DESC LIMIT sqlc.arg('limit');
-- name: InsertUPISettingsVersion :one
INSERT INTO maintenance_payment_settings_versions(society_id,version,enabled,upi_id,payee_name,created_by) VALUES($1,$2,$3,$4,$5,$6) RETURNING *;
-- name: PointUPISettings :exec
INSERT INTO maintenance_payment_settings(society_id,version) VALUES($1,$2) ON CONFLICT(society_id) DO UPDATE SET version=excluded.version;
-- name: SupersedeUPIRequests :exec
UPDATE maintenance_payment_requests SET state='superseded' WHERE society_id=$1 AND state='active';
-- name: CloseUPIRequests :exec
UPDATE maintenance_payment_requests SET state='closed' WHERE society_id=$1 AND bill_id=$2 AND state<>'closed';
-- name: GetActiveUPIRequest :one
SELECT * FROM maintenance_payment_requests WHERE society_id=$1 AND bill_id=$2 AND state='active';
-- name: GetUPIRequest :one
SELECT * FROM maintenance_payment_requests WHERE society_id=$1 AND id=$2;
-- name: InsertUPIRequest :one
INSERT INTO maintenance_payment_requests(id,society_id,bill_id,settings_version,amount_paise,reference,created_by) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING *;

-- name: GetPendingUPIClaim :one
SELECT * FROM maintenance_payment_claims WHERE society_id=$1 AND bill_id=$2 AND status='pending';
-- name: GetUPIClaim :one
SELECT * FROM maintenance_payment_claims WHERE society_id=$1 AND id=$2;
-- name: InsertUPIClaim :one
INSERT INTO maintenance_payment_claims(society_id,bill_id,request_id,submitted_by,reference,payment_date) VALUES($1,$2,$3,$4,$5,$6) RETURNING *;
-- name: ReviewUPIClaim :one
UPDATE maintenance_payment_claims SET status=$3,reviewed_by=$4,reviewed_at=now(),reason=$5 WHERE society_id=$1 AND id=$2 AND status='pending' RETURNING *;
-- name: ListUPIClaims :many
SELECT * FROM maintenance_payment_claims c WHERE c.society_id=$1
AND (sqlc.arg('user_id')::bigint=0 OR c.submitted_by=sqlc.arg('user_id'))
AND (sqlc.arg('status')::text='' OR c.status=sqlc.arg('status'))
AND (sqlc.arg('reference')::text='' OR c.reference=sqlc.arg('reference'))
AND (sqlc.arg('before_id')::bigint=0 OR c.id<sqlc.arg('before_id'))
AND (sqlc.arg('bill_id')::bigint=0 OR c.bill_id=sqlc.arg('bill_id'))
AND EXISTS(
  SELECT 1 FROM maintenance_bills fb
  JOIN flats bf ON bf.id=fb.flat_id AND bf.society_id=fb.society_id AND bf.is_active
  WHERE fb.id=c.bill_id AND fb.society_id=c.society_id
  AND (sqlc.arg('flat_id')::bigint=0 OR fb.flat_id=sqlc.arg('flat_id'))
  AND (sqlc.narg('billing_month')::date IS NULL OR fb.billing_month=sqlc.narg('billing_month'))
  AND (sqlc.narg('block')::text IS NULL OR COALESCE(bf.block, '') ILIKE '%' || sqlc.narg('block')::text || '%')
  AND (sqlc.narg('flat_number')::text IS NULL OR bf.flat_number ILIKE '%' || sqlc.narg('flat_number')::text || '%')
  AND (
    sqlc.arg('search')::text = ''
    OR bf.flat_number ILIKE '%' || sqlc.arg('search')::text || '%'
    OR COALESCE(bf.block, '') ILIKE '%' || sqlc.arg('search')::text || '%'
  )
)
AND (sqlc.arg('user_id')::bigint=0 OR EXISTS(SELECT 1 FROM maintenance_bills b JOIN flats f ON f.id=b.flat_id AND f.is_active JOIN flat_residents r ON r.flat_id=b.flat_id AND r.society_id=b.society_id JOIN society_members m ON m.society_id=r.society_id AND m.user_id=r.user_id WHERE b.id=c.bill_id AND r.user_id=sqlc.arg('user_id') AND r.status='active' AND m.status='active'))
ORDER BY c.id DESC LIMIT sqlc.arg('limit');

-- name: ReserveUPIReference :exec
INSERT INTO maintenance_payment_reference_reservations(society_id,reference,bill_id,claim_id,payment_id) VALUES($1,$2,$3,$4,$5);
-- name: ReleaseUPIClaimReference :exec
DELETE FROM maintenance_payment_reference_reservations WHERE society_id=$1 AND claim_id=$2;
-- name: ReleaseUPIPaymentReference :exec
DELETE FROM maintenance_payment_reference_reservations WHERE society_id=$1 AND payment_id=$2;

-- name: GetActiveUPIPayment :one
SELECT * FROM maintenance_payments WHERE society_id=$1 AND bill_id=$2 AND status='verified';
-- name: GetUPIPayment :one
SELECT * FROM maintenance_payments WHERE society_id=$1 AND id=$2;
-- name: HasReversedUPIReference :one
SELECT EXISTS(SELECT 1 FROM maintenance_payments WHERE society_id=$1 AND reference=$2 AND status='reversed')::boolean;
-- name: InsertUPIPayment :one
INSERT INTO maintenance_payments(society_id,bill_id,claim_id,payer_id,settings_version,amount_paise,reference,credit_date,receipt_number,verified_by,evidence_reference)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING *;
-- name: ReverseUPIPayment :one
UPDATE maintenance_payments SET status='reversed',reversed_by=$3,reversed_at=now(),reversal_reason=$4 WHERE society_id=$1 AND id=$2 AND status='verified' RETURNING *;
-- name: ListUPIPayments :many
SELECT p.* FROM maintenance_payments p WHERE p.society_id=$1
AND (sqlc.arg('status')::text='' OR p.status=sqlc.arg('status'))
AND (sqlc.arg('reference')::text='' OR p.reference=sqlc.arg('reference'))
AND (sqlc.arg('before_id')::bigint=0 OR p.id<sqlc.arg('before_id'))
AND (sqlc.arg('bill_id')::bigint=0 OR p.bill_id=sqlc.arg('bill_id'))
AND EXISTS(
  SELECT 1 FROM maintenance_bills fb
  JOIN flats bf ON bf.id=fb.flat_id AND bf.society_id=fb.society_id AND bf.is_active
  WHERE fb.id=p.bill_id AND fb.society_id=p.society_id
  AND (sqlc.arg('flat_id')::bigint=0 OR fb.flat_id=sqlc.arg('flat_id'))
  AND (sqlc.narg('billing_month')::date IS NULL OR fb.billing_month=sqlc.narg('billing_month'))
  AND (sqlc.narg('block')::text IS NULL OR COALESCE(bf.block, '') ILIKE '%' || sqlc.narg('block')::text || '%')
  AND (sqlc.narg('flat_number')::text IS NULL OR bf.flat_number ILIKE '%' || sqlc.narg('flat_number')::text || '%')
  AND (
    sqlc.arg('search')::text = ''
    OR bf.flat_number ILIKE '%' || sqlc.arg('search')::text || '%'
    OR COALESCE(bf.block, '') ILIKE '%' || sqlc.arg('search')::text || '%'
  )
)
AND (sqlc.arg('user_id')::bigint=0 OR EXISTS(SELECT 1 FROM maintenance_bills b JOIN flats f ON f.id=b.flat_id AND f.is_active JOIN flat_residents r ON r.flat_id=b.flat_id AND r.society_id=b.society_id JOIN society_members m ON m.society_id=r.society_id AND m.user_id=r.user_id WHERE b.id=p.bill_id AND r.user_id=sqlc.arg('user_id') AND r.status='active' AND m.status='active'))
ORDER BY p.id DESC LIMIT sqlc.arg('limit');
-- name: InsertUPILedger :exec
INSERT INTO maintenance_payment_ledger(society_id,bill_id,payment_id,kind,amount_paise,actor_id) VALUES($1,$2,$3,$4,$5,$6);
-- name: InsertUPIAudit :exec
INSERT INTO maintenance_payment_audit_events(society_id,bill_id,actor_id,action,entity_id,details) VALUES($1,$2,$3,$4,$5,$6);
-- name: ListUPIAudit :many
SELECT * FROM maintenance_payment_audit_events WHERE society_id=$1 AND (sqlc.arg('bill_id')::bigint=0 OR bill_id=sqlc.arg('bill_id')) AND (sqlc.arg('before_id')::bigint=0 OR id<sqlc.arg('before_id')) ORDER BY id DESC LIMIT sqlc.arg('limit');

-- name: GetUPIIdempotency :one
SELECT * FROM maintenance_payment_idempotency WHERE society_id=$1 AND actor_id=$2 AND operation=$3 AND key=$4;
-- name: InsertUPIIdempotency :exec
INSERT INTO maintenance_payment_idempotency(society_id,actor_id,operation,key,request_hash,response) VALUES($1,$2,$3,$4,$5,$6);

-- name: GetUPIReport :one
SELECT sqlc.embed(r), b.bill_number FROM maintenance_payment_reports r
JOIN maintenance_bills b ON b.id = r.bill_id AND b.society_id = r.society_id
WHERE r.society_id=$1 AND r.id=$2;
-- name: InsertUPIReport :one
INSERT INTO maintenance_payment_reports(society_id,bill_id,request_id,reported_by,reference,amount_paise,payment_date,explanation,fingerprint)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT(society_id,reported_by,fingerprint) DO NOTHING RETURNING *;
-- name: GetUPIReportByFingerprint :one
SELECT * FROM maintenance_payment_reports WHERE society_id=$1 AND reported_by=$2 AND fingerprint=$3;
-- name: UpdateUPIReport :one
UPDATE maintenance_payment_reports SET status=$3,resolution_note=$4,updated_by=$5,updated_at=now() WHERE society_id=$1 AND id=$2 RETURNING *;
-- name: ListUPIReports :many
SELECT sqlc.embed(r), b.bill_number FROM maintenance_payment_reports r
JOIN maintenance_bills b ON b.id = r.bill_id AND b.society_id = r.society_id
WHERE r.society_id=$1
AND (sqlc.arg('user_id')::bigint=0 OR r.reported_by=sqlc.arg('user_id'))
AND (sqlc.arg('status')::text='' OR r.status=sqlc.arg('status'))
AND (sqlc.arg('before_id')::bigint=0 OR r.id<sqlc.arg('before_id'))
AND (sqlc.arg('user_id')::bigint=0 OR EXISTS(SELECT 1 FROM maintenance_bills b JOIN flats f ON f.id=b.flat_id AND f.is_active JOIN flat_residents fr ON fr.flat_id=b.flat_id AND fr.society_id=b.society_id JOIN society_members m ON m.society_id=fr.society_id AND m.user_id=fr.user_id WHERE b.id=r.bill_id AND fr.user_id=sqlc.arg('user_id') AND fr.status='active' AND m.status='active'))
ORDER BY r.id DESC LIMIT sqlc.arg('limit');

-- name: UPICollectionSummary :one
SELECT count(*)::bigint AS billed_count, COALESCE(sum(b.total_paise),0)::bigint AS billed_paise,
COALESCE(sum(v.paid_amount_paise),0)::bigint AS collected_paise,
COALESCE(sum(v.outstanding_amount_paise),0)::bigint AS outstanding_paise,
COALESCE(sum(v.outstanding_amount_paise) FILTER(WHERE v.payment_status='overdue'),0)::bigint AS overdue_paise,
count(*) FILTER(WHERE v.payment_claim_status='pending')::bigint AS pending_claims
FROM maintenance_bills b JOIN maintenance_bill_balances v ON v.bill_id=b.id
WHERE b.society_id=$1 AND (sqlc.arg('flat_id')::bigint=0 OR b.flat_id=sqlc.arg('flat_id')) AND (sqlc.narg('month')::date IS NULL OR b.billing_month=sqlc.narg('month'));

-- name: EnqueueUPIEvent :exec
INSERT INTO maintenance_notification_deliveries(bill_id,user_id,event_type,event_key,audience,event_data)
SELECT sqlc.arg('bill_id'),m.user_id,sqlc.arg('event_type'),sqlc.arg('event_key'),sqlc.arg('audience'),sqlc.arg('event_data')
FROM society_members m WHERE m.society_id=sqlc.arg('society_id') AND m.status='active' AND
((sqlc.arg('audience')='admin' AND m.role IN ('owner','admin')) OR
(sqlc.arg('audience')='resident' AND EXISTS(SELECT 1 FROM flat_residents r JOIN maintenance_bills b ON b.flat_id=r.flat_id AND b.society_id=r.society_id WHERE b.id=sqlc.arg('bill_id') AND r.user_id=m.user_id AND r.status='active')))
ON CONFLICT(user_id,event_key) DO NOTHING;

-- name: GetPaymentBill :one
SELECT * FROM maintenance_bills WHERE society_id=$1 AND id=$2;
