-- name: GetUserImage :one
SELECT avatar_url, avatar_imagekit_file_id, avatar_imagekit_file_path FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: LockUserImage :one
SELECT avatar_url, avatar_imagekit_file_id, avatar_imagekit_file_path FROM users WHERE id = $1 AND deleted_at IS NULL FOR UPDATE;

-- name: SetUserImage :execrows
UPDATE users SET avatar_url = sqlc.narg('url'), avatar_imagekit_file_id = sqlc.narg('file_id'),
 avatar_imagekit_file_path = sqlc.narg('file_path'), updated_at = NOW() WHERE id = sqlc.arg('id') AND deleted_at IS NULL;

-- name: GetVisitorImage :one
SELECT photo_url, photo_imagekit_file_id, photo_imagekit_file_path FROM visitors WHERE id = $1;

-- name: LockVisitorImage :one
SELECT photo_url, photo_imagekit_file_id, photo_imagekit_file_path FROM visitors WHERE id = $1 FOR UPDATE;

-- name: SetVisitorImage :execrows
UPDATE visitors SET photo_url = sqlc.narg('url'), photo_imagekit_file_id = sqlc.narg('file_id'),
 photo_imagekit_file_path = sqlc.narg('file_path'), updated_at = NOW() WHERE id = sqlc.arg('id');
