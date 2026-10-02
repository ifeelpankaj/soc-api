-- name: ListMaintenanceReminderCandidates :many
SELECT b.id, b.due_date, b.timezone
FROM maintenance_bills b
JOIN maintenance_bill_balances v ON v.bill_id=b.id
JOIN maintenance_settings ms ON ms.society_id=b.society_id AND ms.enabled
WHERE b.society_id=sqlc.arg('society_id') AND b.id>sqlc.arg('after_id')
AND v.outstanding_amount_paise>0
AND v.payment_claim_status IS DISTINCT FROM 'pending'
AND (sqlc.arg('as_of')::timestamptz AT TIME ZONE b.timezone)::date IN (b.due_date-3,b.due_date,b.due_date+7)
ORDER BY b.id LIMIT 500;

-- name: EnqueueMaintenanceReminder :execrows
INSERT INTO maintenance_notification_deliveries(bill_id,user_id,event_type,event_key,audience,event_data)
SELECT b.id,r.user_id,'maintenance_payment_reminder',sqlc.arg('event_key'),'resident',sqlc.arg('event_data')::jsonb
FROM maintenance_bills b
JOIN maintenance_bill_balances v ON v.bill_id=b.id AND v.outstanding_amount_paise>0 AND v.payment_claim_status IS DISTINCT FROM 'pending'
JOIN maintenance_settings ms ON ms.society_id=b.society_id AND ms.enabled
JOIN societies s ON s.id=b.society_id AND s.status='active' AND s.deleted_at IS NULL
JOIN flats f ON f.id=b.flat_id AND f.society_id=b.society_id AND f.is_active
JOIN flat_residents r ON r.flat_id=f.id AND r.society_id=f.society_id AND r.status='active'
JOIN society_members m ON m.society_id=r.society_id AND m.user_id=r.user_id AND m.status='active' AND m.role<>'staff'
WHERE b.id=sqlc.arg('bill_id') AND b.society_id=sqlc.arg('society_id')
ON CONFLICT(user_id,event_key) DO NOTHING;

-- name: GetMaintenanceReminderDelivery :one
SELECT b.bill_number, b.billing_month, b.due_date, v.outstanding_amount_paise
FROM maintenance_notification_deliveries d
JOIN maintenance_bills b ON b.id=d.bill_id
JOIN maintenance_bill_balances v ON v.bill_id=b.id AND v.outstanding_amount_paise>0 AND v.payment_claim_status IS DISTINCT FROM 'pending'
JOIN maintenance_settings ms ON ms.society_id=b.society_id AND ms.enabled
WHERE d.user_id=sqlc.arg('user_id') AND d.event_key=sqlc.arg('event_key')
AND d.event_type='maintenance_payment_reminder' AND b.society_id=sqlc.arg('society_id');
