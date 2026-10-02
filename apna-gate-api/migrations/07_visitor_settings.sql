-- +migrate Up

CREATE TYPE visitor_approval_mode AS ENUM (
    'mandatory', --mandatory = every public/guard visitor needs resident approval

    'optional', --optional  = auto approve by default

    'hybrid' --hybrid    = check flat_visitor_settings by purpose
);

CREATE TYPE visitor_purpose AS ENUM (
    'guest', --family, friend, personal visitor
    'delivery', -- food, parcel, courier
    'cab', -- pickup/drop
    'service', -- electrician, plumber, internet, appliance repair
    'maintenance', -- society maintenance staff visit
    'staff' ,-- maid, cook, driver, nurse, regular helper
    'other' -- any other purpose not covered above
);

CREATE TABLE society_visitor_settings (
    id BIGSERIAL PRIMARY KEY,

    society_id BIGINT NOT NULL
        REFERENCES societies(id) ON DELETE CASCADE,

    approval_mode visitor_approval_mode NOT NULL DEFAULT 'hybrid',

    default_visit_duration_minutes INTEGER NOT NULL DEFAULT 360,
    grace_period_minutes INTEGER NOT NULL DEFAULT 60,
    qr_expiry_minutes INTEGER NOT NULL DEFAULT 1440,

    allow_resident_pre_approval BOOLEAN NOT NULL DEFAULT TRUE,
    allow_public_qr_entry BOOLEAN NOT NULL DEFAULT TRUE,
    allow_guard_entry BOOLEAN NOT NULL DEFAULT TRUE,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    updated_by BIGINT
        REFERENCES users(id) ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_society_visitor_settings_society
        UNIQUE (society_id),

    CONSTRAINT chk_society_default_visit_duration
        CHECK (default_visit_duration_minutes > 0),

    CONSTRAINT chk_society_grace_period
        CHECK (grace_period_minutes >= 0),

    CONSTRAINT chk_society_qr_expiry
        CHECK (qr_expiry_minutes > 0)
);

CREATE TABLE flat_visitor_settings (
    id BIGSERIAL PRIMARY KEY,

    society_id BIGINT NOT NULL
        REFERENCES societies(id) ON DELETE CASCADE,

    flat_id BIGINT NOT NULL
        REFERENCES flats(id) ON DELETE CASCADE,

    purpose visitor_purpose NOT NULL,

    approval_required BOOLEAN NOT NULL DEFAULT TRUE,

    default_visit_duration_minutes INTEGER,

    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    updated_by BIGINT
        REFERENCES users(id) ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_flat_visitor_settings_flat_purpose
        UNIQUE (flat_id, purpose),

    CONSTRAINT chk_flat_visit_duration
        CHECK (
            default_visit_duration_minutes IS NULL
            OR default_visit_duration_minutes > 0
        )
);

CREATE INDEX idx_flat_visitor_settings_society_id
    ON flat_visitor_settings (society_id);

CREATE INDEX idx_flat_visitor_settings_flat_id
    ON flat_visitor_settings (flat_id);

-- +migrate Down

DROP INDEX IF EXISTS idx_flat_visitor_settings_flat_id;
DROP INDEX IF EXISTS idx_flat_visitor_settings_society_id;

DROP TABLE IF EXISTS flat_visitor_settings;
DROP TABLE IF EXISTS society_visitor_settings;

DROP TYPE IF EXISTS visitor_purpose;
DROP TYPE IF EXISTS visitor_approval_mode;