-- +migrate Up

CREATE TYPE flat_member_invite_role AS ENUM (
    'family',
    'tenant'
);

CREATE TYPE flat_member_invite_status AS ENUM (
    'pending',
    'accepted',
    'expired',
    'cancelled'
);

CREATE TABLE flat_member_invites (
    id BIGSERIAL PRIMARY KEY,

    society_id BIGINT NOT NULL,
    flat_id BIGINT NOT NULL,
    invited_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    role flat_member_invite_role NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    full_name VARCHAR(200) NOT NULL,

    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    status flat_member_invite_status NOT NULL DEFAULT 'pending',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_flat_member_invites_flat_society
        FOREIGN KEY (flat_id, society_id)
        REFERENCES flats(id, society_id)
        ON DELETE CASCADE,

    CONSTRAINT flat_member_invites_full_name_check CHECK (full_name <> ''),

    CONSTRAINT flat_member_invites_contact_check CHECK (
        phone IS NOT NULL OR email IS NOT NULL
    )
);

CREATE UNIQUE INDEX uq_flat_member_invites_token_hash
    ON flat_member_invites (token_hash);

CREATE INDEX idx_flat_member_invites_flat_status
    ON flat_member_invites (flat_id, status);

CREATE INDEX idx_flat_member_invites_society_id
    ON flat_member_invites (society_id);

CREATE INDEX idx_flat_member_invites_expires_at
    ON flat_member_invites (expires_at)
    WHERE status = 'pending';

CREATE TRIGGER flat_member_invites_updated_at
BEFORE UPDATE ON flat_member_invites
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


-- +migrate Down

DROP TRIGGER IF EXISTS flat_member_invites_updated_at ON flat_member_invites;

DROP TABLE IF EXISTS flat_member_invites;

DROP TYPE IF EXISTS flat_member_invite_status;
DROP TYPE IF EXISTS flat_member_invite_role;
