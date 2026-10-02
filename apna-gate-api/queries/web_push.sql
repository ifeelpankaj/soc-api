-- name: LockWebPushUser :one
SELECT u.id FROM users u
WHERE u.id = sqlc.arg(user_id) AND u.session_version = sqlc.arg(session_version)
AND u.is_active AND NOT u.is_blocked AND u.deleted_at IS NULL
AND (EXISTS (SELECT 1 FROM flat_residents r JOIN flats f ON f.id = r.flat_id JOIN societies s ON s.id = r.society_id WHERE r.user_id = u.id AND r.status = 'active' AND f.is_active AND s.status = 'active' AND s.deleted_at IS NULL)
 OR EXISTS (SELECT 1 FROM society_members m JOIN societies s ON s.id = m.society_id WHERE m.user_id = u.id AND m.role = 'resident' AND m.status = 'active' AND s.status = 'active' AND s.deleted_at IS NULL))
FOR SHARE OF u;

-- name: DeletePreviousWebInstallation :exec
DELETE FROM device_tokens WHERE platform = 'web'
AND (token = sqlc.arg(token) OR device_id = sqlc.arg(device_id));

-- name: InsertWebDeviceToken :one
INSERT INTO device_tokens(user_id, token, platform, device_id, session_version)
VALUES (sqlc.arg(user_id), sqlc.arg(token), 'web', sqlc.arg(device_id), sqlc.arg(session_version)) RETURNING *;

-- name: ListEligibleWebTokens :many
SELECT d.* FROM device_tokens d JOIN users u ON u.id = d.user_id
WHERE d.user_id = sqlc.arg(user_id) AND d.platform = 'web'
AND d.session_version = u.session_version
AND u.is_active AND NOT u.is_blocked AND u.deleted_at IS NULL
AND (EXISTS (SELECT 1 FROM flat_residents r JOIN flats f ON f.id = r.flat_id JOIN societies s ON s.id = r.society_id
  WHERE r.user_id = u.id AND r.status = 'active' AND f.is_active AND s.status = 'active' AND s.deleted_at IS NULL
  AND (sqlc.arg(society_id)::bigint = 0 OR r.society_id = sqlc.arg(society_id))
  AND (sqlc.arg(flat_id)::bigint = 0 OR r.flat_id = sqlc.arg(flat_id)))
 OR (sqlc.arg(flat_id)::bigint = 0 AND EXISTS (SELECT 1 FROM society_members m JOIN societies s ON s.id = m.society_id
  WHERE m.user_id = u.id AND m.role = 'resident' AND m.status = 'active' AND s.status = 'active' AND s.deleted_at IS NULL
  AND (sqlc.arg(society_id)::bigint = 0 OR m.society_id = sqlc.arg(society_id))))
) ORDER BY d.id;
