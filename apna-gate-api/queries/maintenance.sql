-- name: MaintenanceLock :exec
SELECT pg_advisory_xact_lock(hashtextextended('maintenance/' || sqlc.arg('society_id')::bigint::text, 0));

-- name: GetMaintenanceSettings :one
SELECT config FROM maintenance_settings WHERE society_id=$1;

-- name: SaveMaintenanceSettings :exec
INSERT INTO maintenance_settings(society_id, enabled, config, updated_by) VALUES($1,$2,$3,$4)
ON CONFLICT(society_id) DO UPDATE SET enabled=excluded.enabled, config=excluded.config, updated_by=excluded.updated_by, updated_at=now();

-- name: DeleteMaintenanceTypeRates :exec
DELETE FROM maintenance_type_rates WHERE society_id=$1;

-- name: SaveMaintenanceTypeRate :exec
INSERT INTO maintenance_type_rates(society_id,flat_type,amount_paise) VALUES($1,$2,$3);

-- name: MaintenanceAdmin :one
SELECT EXISTS(SELECT 1 FROM society_members m JOIN societies s ON s.id=m.society_id WHERE m.society_id=$1 AND m.user_id=$2 AND m.status='active' AND m.role IN ('owner','admin') AND s.status='active' AND s.deleted_at IS NULL)::boolean;

-- name: ListMaintenanceFlats :many
SELECT f.id, f.flat_number, f.block, f.flat_type, f.area_sqft_hundredths,
 COALESCE((SELECT jsonb_build_object('user_id',u.id::text,'name',u.full_name) FROM flat_residents r JOIN users u ON u.id=r.user_id WHERE r.society_id=f.society_id AND r.flat_id=f.id AND r.status='active' AND r.is_primary ORDER BY r.id LIMIT 1), '{}'::jsonb)::jsonb AS billed_party
FROM flats f WHERE f.society_id=$1 AND f.is_active AND f.status::text=ANY(sqlc.arg('statuses')::text[]) ORDER BY f.id;

-- name: GetMaintenanceRun :one
SELECT r.id, r.snapshot, (SELECT count(*) FROM maintenance_bills b WHERE b.run_id=r.id)::bigint AS bill_count FROM maintenance_billing_runs r WHERE r.society_id=$1 AND r.billing_month=$2;

-- name: ListMaintenanceIssuedFlats :many
SELECT flat_id FROM maintenance_bills WHERE society_id=$1 AND billing_month=$2 ORDER BY flat_id;

-- name: AuditMaintenanceReconciliation :exec
INSERT INTO maintenance_billing_audit(society_id,actor_id,event_type,entity_id,details)
VALUES($1,$2,'billing_reconciled',$3,$4);

-- name: CreateMaintenanceRun :one
INSERT INTO maintenance_billing_runs(society_id,billing_month,snapshot,created_by) VALUES($1,$2,$3,$4) RETURNING id;

-- name: CreateMaintenanceBill :one
INSERT INTO maintenance_bills(run_id,society_id,flat_id,billing_month,bill_number,due_date,timezone,total_paise,snapshot,billed_party)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id;

-- name: CreateMaintenanceBillItem :exec
INSERT INTO maintenance_bill_items(bill_id,position,description,amount_paise) VALUES($1,$2,$3,$4);

-- name: EnqueueMaintenanceNotifications :exec
INSERT INTO maintenance_notification_deliveries(bill_id,user_id,event_key)
SELECT b.id,r.user_id,'maintenance_bill_generated:'||b.id FROM maintenance_bills b JOIN flat_residents r ON r.society_id=b.society_id AND r.flat_id=b.flat_id
JOIN society_members m ON m.society_id=r.society_id AND m.user_id=r.user_id AND m.status='active'
WHERE b.id=$1 AND r.status='active'
ON CONFLICT(user_id,event_key) DO NOTHING;

-- name: ListMaintenanceBills :many
SELECT sqlc.embed(b), v.paid_amount_paise,v.outstanding_amount_paise,v.payment_claim_status,v.payment_status, p.credit_date AS paid_on FROM maintenance_bills b
JOIN maintenance_bill_balances v ON v.bill_id=b.id
LEFT JOIN maintenance_payments p ON p.bill_id=b.id AND p.society_id=b.society_id AND p.status='verified'
JOIN flats bf ON bf.id=b.flat_id AND bf.society_id=b.society_id AND bf.is_active
WHERE b.society_id=sqlc.arg('society_id')
AND (sqlc.arg('bill_id')::bigint=0 OR b.id=sqlc.arg('bill_id'))
AND (sqlc.arg('flat_id')::bigint=0 OR b.flat_id=sqlc.arg('flat_id'))
AND (sqlc.arg('before_id')::bigint=0 OR b.id<sqlc.arg('before_id'))
AND (sqlc.narg('month')::date IS NULL OR b.billing_month=sqlc.narg('month'))
AND (sqlc.arg('status')::text='' OR v.payment_status=sqlc.arg('status'))
AND (sqlc.arg('display_status')::text=''
 OR (sqlc.arg('display_status')='outstanding' AND v.outstanding_amount_paise>0)
 OR sqlc.arg('display_status')=CASE WHEN v.payment_status='paid' THEN 'paid' WHEN v.payment_claim_status='pending' THEN 'pending_review' WHEN v.payment_claim_status='rejected' THEN 'rejected' ELSE v.payment_status END)
AND (sqlc.arg('bill_number')::text='' OR b.bill_number=sqlc.arg('bill_number'))
AND (sqlc.narg('block')::text IS NULL OR COALESCE(bf.block, '') ILIKE '%' || sqlc.narg('block')::text || '%')
AND (sqlc.narg('flat_number')::text IS NULL OR bf.flat_number ILIKE '%' || sqlc.narg('flat_number')::text || '%')
AND (
    sqlc.arg('search')::text = ''
    OR bf.flat_number ILIKE '%' || sqlc.arg('search')::text || '%'
    OR COALESCE(bf.block, '') ILIKE '%' || sqlc.arg('search')::text || '%'
)
AND (NOT sqlc.arg('resident')::boolean OR EXISTS (
 SELECT 1 FROM flat_residents r JOIN society_members m ON m.society_id=r.society_id AND m.user_id=r.user_id
 JOIN flats f ON f.id=r.flat_id AND f.society_id=r.society_id AND f.is_active
 WHERE r.society_id=b.society_id AND r.flat_id=b.flat_id AND r.user_id=sqlc.arg('user_id') AND r.status='active' AND m.status='active'
))
ORDER BY b.id DESC LIMIT sqlc.arg('limit') OFFSET sqlc.arg('page_offset');

-- name: CountMaintenanceBills :one
SELECT count(*)::bigint FROM maintenance_bills b
JOIN maintenance_bill_balances v ON v.bill_id=b.id
JOIN flats bf ON bf.id=b.flat_id AND bf.society_id=b.society_id AND bf.is_active
WHERE b.society_id=sqlc.arg('society_id')
AND (sqlc.arg('bill_id')::bigint=0 OR b.id=sqlc.arg('bill_id'))
AND (sqlc.arg('flat_id')::bigint=0 OR b.flat_id=sqlc.arg('flat_id'))
AND (sqlc.arg('before_id')::bigint=0 OR b.id<sqlc.arg('before_id'))
AND (sqlc.narg('month')::date IS NULL OR b.billing_month=sqlc.narg('month'))
AND (sqlc.arg('status')::text='' OR v.payment_status=sqlc.arg('status'))
AND (sqlc.arg('display_status')::text=''
 OR (sqlc.arg('display_status')='outstanding' AND v.outstanding_amount_paise>0)
 OR sqlc.arg('display_status')=CASE WHEN v.payment_status='paid' THEN 'paid' WHEN v.payment_claim_status='pending' THEN 'pending_review' WHEN v.payment_claim_status='rejected' THEN 'rejected' ELSE v.payment_status END)
AND (sqlc.arg('bill_number')::text='' OR b.bill_number=sqlc.arg('bill_number'))
AND (sqlc.narg('block')::text IS NULL OR COALESCE(bf.block, '') ILIKE '%' || sqlc.narg('block')::text || '%')
AND (sqlc.narg('flat_number')::text IS NULL OR bf.flat_number ILIKE '%' || sqlc.narg('flat_number')::text || '%')
AND (
    sqlc.arg('search')::text = ''
    OR bf.flat_number ILIKE '%' || sqlc.arg('search')::text || '%'
    OR COALESCE(bf.block, '') ILIKE '%' || sqlc.arg('search')::text || '%'
)
AND (NOT sqlc.arg('resident')::boolean OR EXISTS (
 SELECT 1 FROM flat_residents r JOIN society_members m ON m.society_id=r.society_id AND m.user_id=r.user_id
 JOIN flats f ON f.id=r.flat_id AND f.society_id=r.society_id AND f.is_active
 WHERE r.society_id=b.society_id AND r.flat_id=b.flat_id AND r.user_id=sqlc.arg('user_id') AND r.status='active' AND m.status='active'
));

-- name: ListMaintenanceSocieties :many
SELECT m.society_id FROM maintenance_settings m JOIN societies s ON s.id=m.society_id WHERE m.enabled AND s.status='active' AND s.deleted_at IS NULL ORDER BY m.society_id;

-- name: ClaimMaintenanceDelivery :one
UPDATE maintenance_notification_deliveries d SET attempts=d.attempts+1, available_at=now()+interval '5 minutes', lease_token=sqlc.arg('lease_token')
WHERE d.id=(SELECT id FROM maintenance_notification_deliveries WHERE completed_at IS NULL AND available_at<=now() ORDER BY available_at,id FOR UPDATE SKIP LOCKED LIMIT 1)
RETURNING d.*;

-- name: GetMaintenanceDeliveryBill :one
SELECT b.* FROM maintenance_bills b JOIN societies s ON s.id=b.society_id
WHERE b.id=sqlc.arg('bill_id') AND s.status='active' AND s.deleted_at IS NULL AND (
(sqlc.arg('audience')::text='admin' AND EXISTS(SELECT 1 FROM society_members m WHERE m.society_id=b.society_id AND m.user_id=sqlc.arg('user_id') AND m.status='active' AND m.role IN ('owner','admin')))
OR (sqlc.arg('audience')::text='resident' AND EXISTS(
 SELECT 1 FROM flat_residents r JOIN society_members m ON m.society_id=r.society_id AND m.user_id=r.user_id
 JOIN flats f ON f.id=r.flat_id AND f.society_id=r.society_id AND f.is_active
 WHERE r.society_id=b.society_id AND r.flat_id=b.flat_id AND r.user_id=sqlc.arg('user_id') AND r.status='active' AND m.status='active'
)));

-- name: CompleteMaintenanceDelivery :exec
UPDATE maintenance_notification_deliveries SET completed_at=now(),last_error=NULL WHERE id=$1 AND lease_token=$2;

-- name: RetryMaintenanceDelivery :exec
UPDATE maintenance_notification_deliveries SET available_at=now()+(LEAST(3600,30*power(2,LEAST(attempts,7)))::int * interval '1 second'), last_error=sqlc.arg('last_error') WHERE id=sqlc.arg('id') AND lease_token=sqlc.arg('lease_token');
