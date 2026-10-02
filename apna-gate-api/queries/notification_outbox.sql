-- name: EnqueueNotificationOutbox :exec
INSERT INTO notification_outbox(user_id,society_id,flat_id,audience,event_key,payload) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT(user_id,event_key) DO NOTHING;

-- name: ClaimNotificationOutbox :one
UPDATE notification_outbox d SET attempts=d.attempts+1,available_at=now()+interval '2 minutes',lease_token=$1
WHERE d.id=(SELECT id FROM notification_outbox WHERE completed_at IS NULL AND available_at<=now() ORDER BY available_at,id FOR UPDATE SKIP LOCKED LIMIT 1) RETURNING d.*;

-- name: NotificationOutboxAccess :one
SELECT EXISTS(SELECT 1 FROM society_members m JOIN societies s ON s.id=m.society_id
WHERE m.society_id=sqlc.arg('society_id') AND m.user_id=sqlc.arg('user_id') AND m.status='active' AND s.status='active' AND s.deleted_at IS NULL
AND ((sqlc.arg('audience')::text='admin' AND m.role IN ('owner','admin'))
OR (sqlc.arg('audience')::text='hub_member' AND m.role IN ('owner','admin','resident') AND EXISTS(SELECT 1 FROM users hu WHERE hu.id=m.user_id AND hu.is_active AND NOT hu.is_blocked AND hu.deleted_at IS NULL))
OR (sqlc.arg('audience')::text='staff' AND m.role='staff')
OR (sqlc.arg('audience')::text IN ('resident','society_resident') AND m.role<>'staff' AND EXISTS(
 SELECT 1 FROM flat_residents r JOIN flats f ON f.id=r.flat_id AND f.society_id=r.society_id AND f.is_active
 WHERE r.society_id=m.society_id AND r.user_id=m.user_id AND r.status='active' AND (sqlc.arg('audience')::text='society_resident' OR r.flat_id=sqlc.narg('flat_id'))
))))::boolean;

-- name: MarkOutboxInbox :execrows
UPDATE notification_outbox SET inbox_completed_at=now() WHERE id=$1 AND lease_token=$2 AND available_at>now() AND completed_at IS NULL;

-- name: CompleteNotificationOutbox :exec
UPDATE notification_outbox SET completed_at=now(),push_completed_at=CASE WHEN sqlc.arg('delivered')::boolean THEN now() ELSE push_completed_at END,last_error=NULL
WHERE id=sqlc.arg('id') AND lease_token=sqlc.arg('lease_token') AND available_at>now();

-- name: RetryNotificationOutbox :exec
UPDATE notification_outbox SET available_at=now()+(LEAST(3600,30*power(2,LEAST(attempts,7)))::int*interval '1 second'),lease_token=NULL,last_error='delivery_failed'
WHERE id=$1 AND lease_token=$2;

-- name: NotificationBacklog :one
SELECT count(*)::bigint AS pending,COALESCE(extract(epoch FROM now()-min(available_at)),0)::double precision AS oldest_seconds FROM notification_outbox WHERE completed_at IS NULL;

-- name: OutboxInboxID :one
SELECT id FROM notifications WHERE user_id=$1 AND event_key=$2;
