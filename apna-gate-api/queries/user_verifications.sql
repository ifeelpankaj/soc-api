-- name: CreateUserVerification :one
INSERT INTO user_verifications (
    user_id,
    purpose,
    target,
    otp_hash,
    attempts,
    max_attempts,
    is_used,
    expires_at
)
VALUES (
    sqlc.arg('user_id'),
    sqlc.arg('purpose')::verification_purpose,
    sqlc.arg('target'),
    sqlc.arg('otp_hash'),
    sqlc.arg('attempts'),
    sqlc.arg('max_attempts'),
    sqlc.arg('is_used'),
    sqlc.arg('expires_at')
)
RETURNING
    id,
    user_id,
    purpose,
    target,
    otp_hash,
    attempts,
    max_attempts,
    is_used,
    expires_at,
    used_at,
    created_at;

-- name: GetActiveUserVerification :one
SELECT
    id,
    user_id,
    purpose,
    target,
    otp_hash,
    attempts,
    max_attempts,
    is_used,
    expires_at,
    used_at,
    created_at
FROM user_verifications
WHERE user_id = sqlc.arg('user_id')
  AND purpose = sqlc.arg('purpose')::verification_purpose
  AND target = sqlc.arg('target')
  AND is_used = false
  AND expires_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkUserVerificationUsed :exec
UPDATE user_verifications
SET is_used = true,
    used_at = NOW()
WHERE id = $1
  AND is_used = false;

-- name: IncrementUserVerificationAttempts :exec
UPDATE user_verifications
SET attempts = attempts + 1
WHERE id = $1;

-- name: DeleteActiveUserVerificationByPurpose :exec
UPDATE user_verifications
SET is_used = true,
    used_at = NOW()
WHERE user_id = sqlc.arg('user_id')
  AND purpose = sqlc.arg('purpose')::verification_purpose
  AND target = sqlc.arg('target')
  AND is_used = false;
-- name: MarkUserEmailVerified :exec
UPDATE users
SET
    email_verified = TRUE,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUsedOrExpiredUserVerifications :exec
DELETE FROM user_verifications
WHERE is_used = TRUE
   OR expires_at < NOW();
