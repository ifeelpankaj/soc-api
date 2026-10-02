-- +migrate Up

CREATE TYPE global_role AS ENUM (
    'user',
    'developer',
    'super_admin'
);

CREATE TYPE auth_provider AS ENUM (
    'email',
    'google',
    'apple',
    'phone'
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    -- Identity
    first_name VARCHAR(100),
    last_name VARCHAR(100),

    -- generated full display name
    full_name VARCHAR(200) NOT NULL,

    -- Authentication
    email VARCHAR(255),
    phone_number VARCHAR(20),
    password_hash TEXT,

    -- Authentication provider
    auth_provider auth_provider NOT NULL DEFAULT 'email',
    provider_id TEXT,

    -- Platform-level role only
    global_role global_role NOT NULL DEFAULT 'user',

    -- Verification
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    phone_verified BOOLEAN NOT NULL DEFAULT FALSE,

    -- Account state
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
    blocked_reason TEXT,

    -- Profile
    avatar_url TEXT,
    date_of_birth DATE,
    gender VARCHAR(20),

    -- Preferences
    timezone VARCHAR(100) NOT NULL DEFAULT 'Asia/Kolkata',
    language VARCHAR(20) NOT NULL DEFAULT 'en',

    -- Security
    last_login_at TIMESTAMPTZ,
    password_changed_at TIMESTAMPTZ,

    -- Soft delete
    deleted_at TIMESTAMPTZ,

    -- Metadata / extensibility
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    -- Audit timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints

    CONSTRAINT users_email_or_phone_required
    CHECK (
        email IS NOT NULL
        OR phone_number IS NOT NULL
    ),

    CONSTRAINT users_email_unique
    UNIQUE(email),

    CONSTRAINT users_phone_unique
    UNIQUE(phone_number)
);

-- =====================================================
-- update timestamp trigger
-- =====================================================

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER users_update_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- =====================================================
-- indexes
-- =====================================================

CREATE INDEX idx_users_email_lower
ON users (LOWER(email))
WHERE deleted_at IS NULL;

CREATE INDEX idx_users_phone
ON users (phone_number)
WHERE deleted_at IS NULL;

CREATE INDEX idx_users_global_role
ON users (global_role);

CREATE INDEX idx_users_active
ON users (is_active);

CREATE INDEX idx_users_blocked
ON users (is_blocked);

CREATE INDEX idx_users_created_at
ON users (created_at DESC);

CREATE INDEX idx_users_provider
ON users (auth_provider, provider_id);

CREATE INDEX idx_users_deleted_at
ON users (deleted_at);

CREATE INDEX idx_users_metadata_gin
ON users
USING GIN(metadata);

-- +migrate Down

DROP TRIGGER IF EXISTS users_update_timestamp ON users;

DROP INDEX IF EXISTS idx_users_metadata_gin;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_provider;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_blocked;
DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_users_global_role;
DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_users_email_lower;

DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS auth_provider;
DROP TYPE IF EXISTS global_role;

DROP FUNCTION IF EXISTS update_timestamp();