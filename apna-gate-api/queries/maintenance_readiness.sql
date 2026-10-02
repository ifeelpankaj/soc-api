-- name: GetMaintenanceFirstMonth :one
SELECT first_enabled_month FROM maintenance_settings WHERE society_id=$1;

-- name: SetMaintenanceFirstMonth :exec
UPDATE maintenance_settings SET first_enabled_month=$2 WHERE society_id=$1 AND first_enabled_month IS NULL;

-- name: SaveMaintenanceReview :exec
INSERT INTO maintenance_preview_reviews(token,society_id,actor_id,billing_month,snapshot_hash,expires_at) VALUES($1,$2,$3,$4,$5,$6);

-- name: GetMaintenanceReview :one
SELECT snapshot_hash,expires_at FROM maintenance_preview_reviews WHERE token=$1 AND society_id=$2 AND actor_id=$3 AND billing_month=$4;

-- name: MaintenanceFlatAccess :one
SELECT EXISTS(SELECT 1 FROM flats f JOIN societies s ON s.id=f.society_id
WHERE f.id=$1 AND f.society_id=$2 AND s.status='active' AND s.deleted_at IS NULL AND (
 EXISTS(SELECT 1 FROM society_members m WHERE m.society_id=f.society_id AND m.user_id=$3 AND m.status='active' AND m.role IN ('owner','admin'))
 OR (f.is_active AND EXISTS(SELECT 1 FROM flat_residents r JOIN society_members m ON m.society_id=r.society_id AND m.user_id=r.user_id WHERE r.flat_id=f.id AND r.society_id=f.society_id AND r.user_id=$3 AND r.status='active' AND m.status='active' AND m.role='resident'))
))::boolean;

-- name: MaintenanceOutstandingBills :many
SELECT b.id,b.bill_number,b.billing_month,b.due_date,v.outstanding_amount_paise,v.payment_status
FROM maintenance_bills b JOIN maintenance_bill_balances v ON v.bill_id=b.id
WHERE b.society_id=$1 AND b.flat_id=$2 AND v.outstanding_amount_paise>0 ORDER BY b.billing_month,b.id;

-- name: MaintenancePendingClaimsCount :one
SELECT count(*)::bigint FROM maintenance_payment_claims WHERE status='pending';
