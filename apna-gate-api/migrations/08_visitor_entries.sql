-- +migrate Up

CREATE TYPE visitor_source AS ENUM (
    'resident_link',
    'public_qr',
    'guard_entry',
    'quick_link'
);

CREATE TYPE visitor_status AS ENUM (
    'waiting_approval',
    'approved',
    'rejected',
    'checked_in',
    'checked_out',
    'cancelled',
    'expired',
    'auto_closed'
);

CREATE TYPE visitor_vehicle_type AS ENUM (
    'bike',
    'car',
    'auto',
    'cab',
    'truck',
    'other'
);

CREATE TYPE visitor_event_type AS ENUM (
    'created',
    'approved',
    'rejected',
    'checked_in',
    'checked_out',
    'cancelled',
    'expired',
    'auto_closed',
    'qr_generated',
    'qr_used'
);

CREATE TYPE visitor_invite_status AS ENUM (
    'active',
    'used',
    'expired',
    'cancelled'
);

CREATE TABLE visitors (
    id BIGSERIAL PRIMARY KEY,

    full_name VARCHAR(200) NOT NULL,
    phone_number VARCHAR(20),
    email VARCHAR(255),

    photo_url TEXT,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_visitors_contact_required
        CHECK (
            phone_number IS NOT NULL
            OR email IS NOT NULL
        )
);

CREATE INDEX idx_visitors_phone_number
    ON visitors (phone_number);

CREATE INDEX idx_visitors_email
    ON visitors (email);

CREATE TABLE visitor_invites (
    id BIGSERIAL PRIMARY KEY,

    society_id BIGINT NOT NULL
        REFERENCES societies(id) ON DELETE CASCADE,

    flat_id BIGINT NOT NULL
        REFERENCES flats(id) ON DELETE CASCADE,

    created_by BIGINT NOT NULL
        REFERENCES users(id) ON DELETE RESTRICT,

    purpose visitor_purpose NOT NULL,

    token_hash TEXT NOT NULL,

    status visitor_invite_status NOT NULL DEFAULT 'active',

    expires_at TIMESTAMPTZ NOT NULL,

    used_at TIMESTAMPTZ,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_visitor_invites_used_at
        CHECK (
            used_at IS NULL
            OR status = 'used'
        )
);

CREATE UNIQUE INDEX uq_visitor_invites_token_hash
    ON visitor_invites (token_hash);

CREATE INDEX idx_visitor_invites_society_flat_status
    ON visitor_invites (society_id, flat_id, status);

CREATE INDEX idx_visitor_invites_created_by
    ON visitor_invites (created_by);

CREATE INDEX idx_visitor_invites_expires_at
    ON visitor_invites (expires_at);

CREATE TABLE visitor_entries (
    id BIGSERIAL PRIMARY KEY,

    society_id BIGINT NOT NULL
        REFERENCES societies(id) ON DELETE CASCADE,

    flat_id BIGINT NOT NULL
        REFERENCES flats(id) ON DELETE CASCADE,

    visitor_id BIGINT NOT NULL
        REFERENCES visitors(id) ON DELETE RESTRICT,

    invite_id BIGINT
        REFERENCES visitor_invites(id) ON DELETE SET NULL,

    source visitor_source NOT NULL,
    purpose visitor_purpose NOT NULL,
    status visitor_status NOT NULL DEFAULT 'waiting_approval',

    vehicle_number VARCHAR(50),
    vehicle_type visitor_vehicle_type,

    companions_count INTEGER NOT NULL DEFAULT 0,
    companion_details JSONB NOT NULL DEFAULT '[]'::jsonb,

    expected_at TIMESTAMPTZ,
    expected_checkout_at TIMESTAMPTZ,

    checked_in_at TIMESTAMPTZ,
    checked_out_at TIMESTAMPTZ,
    auto_closed_at TIMESTAMPTZ,

    approved_by BIGINT
        REFERENCES users(id) ON DELETE SET NULL,

    rejected_by BIGINT
        REFERENCES users(id) ON DELETE SET NULL,

    handled_by_guard_id BIGINT
        REFERENCES users(id) ON DELETE SET NULL,

    created_by BIGINT
        REFERENCES users(id) ON DELETE SET NULL,

    qr_token_hash TEXT,
    qr_expires_at TIMESTAMPTZ,
    qr_used_at TIMESTAMPTZ,

    notes TEXT,
    rejection_reason TEXT,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_visitor_entries_companions_count
        CHECK (companions_count >= 0),

    CONSTRAINT chk_visitor_entries_companion_details_array
        CHECK (jsonb_typeof(companion_details) = 'array'),

    CONSTRAINT chk_visitor_entries_checkout_after_checkin
        CHECK (
            checked_out_at IS NULL
            OR checked_in_at IS NULL
            OR checked_out_at >= checked_in_at
        ),

    CONSTRAINT chk_visitor_entries_expected_checkout_after_expected
        CHECK (
            expected_checkout_at IS NULL
            OR expected_at IS NULL
            OR expected_checkout_at >= expected_at
        ),

    CONSTRAINT chk_visitor_entries_qr_expiry_requires_hash
        CHECK (
            qr_expires_at IS NULL
            OR qr_token_hash IS NOT NULL
        ),

    CONSTRAINT chk_visitor_entries_qr_used_requires_hash
        CHECK (
            qr_used_at IS NULL
            OR qr_token_hash IS NOT NULL
        )
);

CREATE INDEX idx_visitor_entries_society_status
    ON visitor_entries (society_id, status);

CREATE INDEX idx_visitor_entries_society_flat
    ON visitor_entries (society_id, flat_id);

CREATE INDEX idx_visitor_entries_visitor_id
    ON visitor_entries (visitor_id);

CREATE INDEX idx_visitor_entries_invite_id
    ON visitor_entries (invite_id);

CREATE INDEX idx_visitor_entries_created_at
    ON visitor_entries (created_at DESC);

CREATE INDEX idx_visitor_entries_expected_at
    ON visitor_entries (expected_at DESC);

CREATE INDEX idx_visitor_entries_checked_in_at
    ON visitor_entries (checked_in_at DESC);

CREATE INDEX idx_visitor_entries_active_inside
    ON visitor_entries (society_id, checked_in_at DESC)
    WHERE status = 'checked_in';

CREATE INDEX idx_visitor_entries_waiting_approval
    ON visitor_entries (society_id, flat_id, created_at DESC)
    WHERE status = 'waiting_approval';

CREATE UNIQUE INDEX uq_visitor_entries_qr_token_hash
    ON visitor_entries (qr_token_hash)
    WHERE qr_token_hash IS NOT NULL;

CREATE TABLE visitor_entry_events (
    id BIGSERIAL PRIMARY KEY,

    visitor_entry_id BIGINT NOT NULL
        REFERENCES visitor_entries(id) ON DELETE CASCADE,

    society_id BIGINT NOT NULL
        REFERENCES societies(id) ON DELETE CASCADE,

    actor_user_id BIGINT
        REFERENCES users(id) ON DELETE SET NULL,

    event_type visitor_event_type NOT NULL,

    message TEXT,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_visitor_entry_events_entry_id
    ON visitor_entry_events (visitor_entry_id, created_at DESC);

CREATE INDEX idx_visitor_entry_events_society_id
    ON visitor_entry_events (society_id, created_at DESC);

CREATE INDEX idx_visitor_entry_events_event_type
    ON visitor_entry_events (event_type);

-- +migrate Down

DROP INDEX IF EXISTS idx_visitor_entry_events_event_type;
DROP INDEX IF EXISTS idx_visitor_entry_events_society_id;
DROP INDEX IF EXISTS idx_visitor_entry_events_entry_id;

DROP TABLE IF EXISTS visitor_entry_events;

DROP INDEX IF EXISTS uq_visitor_entries_qr_token_hash;
DROP INDEX IF EXISTS idx_visitor_entries_waiting_approval;
DROP INDEX IF EXISTS idx_visitor_entries_active_inside;
DROP INDEX IF EXISTS idx_visitor_entries_checked_in_at;
DROP INDEX IF EXISTS idx_visitor_entries_expected_at;
DROP INDEX IF EXISTS idx_visitor_entries_created_at;
DROP INDEX IF EXISTS idx_visitor_entries_invite_id;
DROP INDEX IF EXISTS idx_visitor_entries_visitor_id;
DROP INDEX IF EXISTS idx_visitor_entries_society_flat;
DROP INDEX IF EXISTS idx_visitor_entries_society_status;

DROP TABLE IF EXISTS visitor_entries;

DROP INDEX IF EXISTS idx_visitor_invites_expires_at;
DROP INDEX IF EXISTS idx_visitor_invites_created_by;
DROP INDEX IF EXISTS idx_visitor_invites_society_flat_status;
DROP INDEX IF EXISTS uq_visitor_invites_token_hash;

DROP TABLE IF EXISTS visitor_invites;

DROP INDEX IF EXISTS idx_visitors_email;
DROP INDEX IF EXISTS idx_visitors_phone_number;

DROP TABLE IF EXISTS visitors;

DROP TYPE IF EXISTS visitor_invite_status;
DROP TYPE IF EXISTS visitor_event_type;
DROP TYPE IF EXISTS visitor_vehicle_type;
DROP TYPE IF EXISTS visitor_status;
DROP TYPE IF EXISTS visitor_source;