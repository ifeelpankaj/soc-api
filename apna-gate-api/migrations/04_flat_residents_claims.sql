-- +migrate Up

CREATE TYPE flat_status AS ENUM (
    'vacant',
    'occupied',
    'blocked'
);

CREATE TYPE flat_resident_role AS ENUM (
    'owner',
    'tenant',
    'family'
);

CREATE TYPE flat_resident_status AS ENUM (
    'active',
    'inactive',
    'moved_out'
);

CREATE TYPE flat_claim_status AS ENUM (
    'pending',
    'approved',
    'rejected',
    'cancelled'
);

CREATE TABLE flats (
    id BIGSERIAL PRIMARY KEY,

    society_id BIGINT NOT NULL REFERENCES societies(id) ON DELETE CASCADE,

    block VARCHAR(50),
    floor VARCHAR(50),
    flat_number VARCHAR(50) NOT NULL,

    status flat_status NOT NULL DEFAULT 'vacant',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT flats_flat_number_check CHECK (flat_number <> ''),
    CONSTRAINT flats_block_check CHECK (block IS NULL OR block <> ''),
    CONSTRAINT flats_floor_check CHECK (floor IS NULL OR floor <> '')
);

CREATE UNIQUE INDEX uniq_flats_society_block_flat_number
ON flats (
    society_id,
    COALESCE(block, ''),
    flat_number
);

CREATE UNIQUE INDEX uniq_flats_id_society_id
ON flats (id, society_id);

CREATE INDEX idx_flats_society_id
ON flats (society_id);

CREATE INDEX idx_flats_status
ON flats (status);

CREATE INDEX idx_flats_active
ON flats (is_active);


CREATE TABLE flat_residents (
    id BIGSERIAL PRIMARY KEY,

    society_id BIGINT NOT NULL,
    flat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    role flat_resident_role NOT NULL DEFAULT 'family',
    status flat_resident_status NOT NULL DEFAULT 'active',

    is_primary BOOLEAN NOT NULL DEFAULT FALSE,

    moved_in_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    moved_out_at TIMESTAMPTZ,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_flat_residents_flat_society
        FOREIGN KEY (flat_id, society_id)
        REFERENCES flats(id, society_id)
        ON DELETE CASCADE,

    CONSTRAINT flat_residents_move_out_check CHECK (
        moved_out_at IS NULL OR moved_out_at >= moved_in_at
    )
);

CREATE UNIQUE INDEX uniq_active_flat_resident
ON flat_residents (flat_id, user_id)
WHERE status = 'active';

CREATE UNIQUE INDEX uniq_primary_resident_per_flat
ON flat_residents (flat_id)
WHERE is_primary = TRUE
AND status = 'active';

CREATE INDEX idx_flat_residents_society_id
ON flat_residents (society_id);

CREATE INDEX idx_flat_residents_flat_id
ON flat_residents (flat_id);

CREATE INDEX idx_flat_residents_user_id
ON flat_residents (user_id);

CREATE INDEX idx_flat_residents_status
ON flat_residents (status);


CREATE TABLE flat_claim_requests (
    id BIGSERIAL PRIMARY KEY,

    society_id BIGINT NOT NULL,
    flat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    requested_role flat_resident_role NOT NULL DEFAULT 'owner',
    requested_primary BOOLEAN NOT NULL DEFAULT TRUE,

    status flat_claim_status NOT NULL DEFAULT 'pending',

    note TEXT,
    rejection_reason TEXT,

    reviewed_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,

    cancelled_at TIMESTAMPTZ,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_flat_claims_flat_society
        FOREIGN KEY (flat_id, society_id)
        REFERENCES flats(id, society_id)
        ON DELETE CASCADE,

    CONSTRAINT flat_claim_note_check CHECK (
        note IS NULL OR length(note) <= 500
    ),

    CONSTRAINT flat_claim_rejection_reason_check CHECK (
        rejection_reason IS NULL OR length(rejection_reason) <= 500
    )
);

CREATE UNIQUE INDEX uniq_pending_flat_claim_per_user_flat
ON flat_claim_requests (flat_id, user_id)
WHERE status = 'pending';

CREATE INDEX idx_flat_claim_requests_society_id
ON flat_claim_requests (society_id);

CREATE INDEX idx_flat_claim_requests_flat_id
ON flat_claim_requests (flat_id);

CREATE INDEX idx_flat_claim_requests_user_id
ON flat_claim_requests (user_id);

CREATE INDEX idx_flat_claim_requests_status
ON flat_claim_requests (status);

CREATE INDEX idx_flat_claim_requests_pending_admin_view
ON flat_claim_requests (society_id, status, created_at)
WHERE status = 'pending';


CREATE TRIGGER flats_updated_at
BEFORE UPDATE ON flats
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER flat_residents_updated_at
BEFORE UPDATE ON flat_residents
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER flat_claim_requests_updated_at
BEFORE UPDATE ON flat_claim_requests
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


-- +migrate Down

DROP TRIGGER IF EXISTS flat_claim_requests_updated_at ON flat_claim_requests;
DROP TRIGGER IF EXISTS flat_residents_updated_at ON flat_residents;
DROP TRIGGER IF EXISTS flats_updated_at ON flats;

DROP TABLE IF EXISTS flat_claim_requests;
DROP TABLE IF EXISTS flat_residents;
DROP TABLE IF EXISTS flats;

DROP TYPE IF EXISTS flat_claim_status;
DROP TYPE IF EXISTS flat_resident_status;
DROP TYPE IF EXISTS flat_resident_role;
DROP TYPE IF EXISTS flat_status;