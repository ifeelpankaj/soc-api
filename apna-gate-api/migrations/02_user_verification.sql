-- +migrate Up

CREATE TYPE verification_purpose AS ENUM (
    'email_verification',
    'phone_verification',
    'password_reset',
    'login_otp'
);

CREATE TABLE user_verifications (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    purpose verification_purpose NOT NULL,

    -- email or phone where OTP was sent
    target VARCHAR(255) NOT NULL,

    otp_hash TEXT NOT NULL,

    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,

    is_used BOOLEAN NOT NULL DEFAULT FALSE,

    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Only one active OTP per user + purpose + target
CREATE UNIQUE INDEX uq_active_user_verification
ON user_verifications (user_id, purpose, target)
WHERE is_used = FALSE;

CREATE INDEX idx_user_verifications_user_id
ON user_verifications(user_id);

CREATE INDEX idx_user_verifications_purpose
ON user_verifications(purpose);

CREATE INDEX idx_user_verifications_target
ON user_verifications(target);

CREATE INDEX idx_user_verifications_expires_at
ON user_verifications(expires_at);

-- +migrate Down

DROP INDEX IF EXISTS idx_user_verifications_expires_at;
DROP INDEX IF EXISTS idx_user_verifications_target;
DROP INDEX IF EXISTS idx_user_verifications_purpose;
DROP INDEX IF EXISTS idx_user_verifications_user_id;
DROP INDEX IF EXISTS uq_active_user_verification;

DROP TABLE IF EXISTS user_verifications;

DROP TYPE IF EXISTS verification_purpose;