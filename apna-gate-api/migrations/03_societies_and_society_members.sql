-- +migrate Up

CREATE TYPE society_status AS ENUM (
    'pending',
    'active',
    'suspended',
    'rejected'
);

CREATE TYPE society_member_role AS ENUM (
    'owner',
    'admin',
    'staff',
    'resident'
);

CREATE TYPE society_member_status AS ENUM (
    'pending',
    'active',
    'suspended',
    'removed'
);

CREATE TABLE societies (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(200) NOT NULL,
    society_code VARCHAR(50) NOT NULL,

    email VARCHAR(255),
    phone_number VARCHAR(20),

    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    landmark VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    pincode VARCHAR(20),
    country VARCHAR(100) NOT NULL DEFAULT 'India',

    total_flats INT NOT NULL DEFAULT 0 CHECK (total_flats >= 0),
    total_blocks INT NOT NULL DEFAULT 0 CHECK (total_blocks >= 0),

    status society_status NOT NULL DEFAULT 'pending',

    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    approved_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMPTZ,

    rejected_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    rejected_at TIMESTAMPTZ,
    rejection_reason TEXT,

    suspended_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    suspended_at TIMESTAMPTZ,
    suspension_reason TEXT,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT chk_societies_approved_valid
        CHECK (
            status != 'active'
            OR (approved_by IS NOT NULL AND approved_at IS NOT NULL)
        ),

    CONSTRAINT chk_societies_rejected_valid
        CHECK (
            status != 'rejected'
            OR (rejected_by IS NOT NULL AND rejected_at IS NOT NULL)
        ),

    CONSTRAINT chk_societies_suspended_valid
        CHECK (
            status != 'suspended'
            OR (suspended_by IS NOT NULL AND suspended_at IS NOT NULL)
        )
);

CREATE TABLE society_members (
    id BIGSERIAL PRIMARY KEY,

    society_id BIGINT NOT NULL REFERENCES societies(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    role society_member_role NOT NULL,
    status society_member_status NOT NULL DEFAULT 'active',

    invited_by BIGINT REFERENCES users(id) ON DELETE SET NULL,

    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    removed_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    removed_at TIMESTAMPTZ,
    remove_reason TEXT,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_society_members_society_user
        UNIQUE (society_id, user_id),

    CONSTRAINT chk_society_members_removed_valid
        CHECK (
            status != 'removed'
            OR removed_at IS NOT NULL
        )
);

-- =====================================================
-- INDEXES: SOCIETIES
-- =====================================================

CREATE UNIQUE INDEX uq_societies_code_active
ON societies(society_code)
WHERE deleted_at IS NULL;

CREATE INDEX idx_societies_status
ON societies(status);

CREATE INDEX idx_societies_created_by
ON societies(created_by);

CREATE INDEX idx_societies_city_state
ON societies(city, state);

CREATE INDEX idx_societies_deleted_at
ON societies(deleted_at);

CREATE INDEX idx_societies_metadata
ON societies USING GIN(metadata);

-- Optional duplicate-name protection
CREATE UNIQUE INDEX uq_societies_name_city_state_active
ON societies(LOWER(name), LOWER(city), LOWER(state))
WHERE deleted_at IS NULL;

-- =====================================================
-- INDEXES: SOCIETY MEMBERS
-- =====================================================

CREATE INDEX idx_society_members_society_id
ON society_members(society_id);

CREATE INDEX idx_society_members_user_id
ON society_members(user_id);

CREATE INDEX idx_society_members_role
ON society_members(role);

CREATE INDEX idx_society_members_status
ON society_members(status);

CREATE INDEX idx_society_members_metadata
ON society_members USING GIN(metadata);

CREATE UNIQUE INDEX uq_society_members_one_active_owner
ON society_members(society_id)
WHERE role = 'owner' AND status = 'active';

CREATE INDEX idx_society_members_active_lookup
ON society_members(society_id, user_id, role)
WHERE status = 'active';

-- =====================================================
-- UPDATED_AT TRIGGERS
-- =====================================================

CREATE TRIGGER trg_societies_updated_at
BEFORE UPDATE ON societies
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trg_society_members_updated_at
BEFORE UPDATE ON society_members
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- +migrate Down

DROP TRIGGER IF EXISTS trg_society_members_updated_at ON society_members;
DROP TRIGGER IF EXISTS trg_societies_updated_at ON societies;

DROP INDEX IF EXISTS idx_society_members_active_lookup;
DROP INDEX IF EXISTS uq_society_members_one_active_owner;
DROP INDEX IF EXISTS idx_society_members_metadata;
DROP INDEX IF EXISTS idx_society_members_status;
DROP INDEX IF EXISTS idx_society_members_role;
DROP INDEX IF EXISTS idx_society_members_user_id;
DROP INDEX IF EXISTS idx_society_members_society_id;

DROP INDEX IF EXISTS uq_societies_name_city_state_active;
DROP INDEX IF EXISTS idx_societies_metadata;
DROP INDEX IF EXISTS idx_societies_deleted_at;
DROP INDEX IF EXISTS idx_societies_city_state;
DROP INDEX IF EXISTS idx_societies_created_by;
DROP INDEX IF EXISTS idx_societies_status;
DROP INDEX IF EXISTS uq_societies_code_active;

DROP TABLE IF EXISTS society_members;
DROP TABLE IF EXISTS societies;

DROP TYPE IF EXISTS society_member_status;
DROP TYPE IF EXISTS society_member_role;
DROP TYPE IF EXISTS society_status;