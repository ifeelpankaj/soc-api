-- +migrate Up

CREATE TABLE short_links (
    id BIGSERIAL PRIMARY KEY,

    short_code VARCHAR(32) NOT NULL UNIQUE,

    resource_type VARCHAR(50) NOT NULL,
    resource_id BIGINT NOT NULL,

    expires_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,

    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_short_links_resource_type
        CHECK (resource_type IN ('visitor_invite', 'member_invite')),

    CONSTRAINT chk_short_links_code_not_empty
        CHECK (short_code <> '')
);

CREATE UNIQUE INDEX idx_short_links_code
    ON short_links (short_code);

CREATE INDEX idx_short_links_resource
    ON short_links (resource_type, resource_id);

CREATE INDEX idx_short_links_active
    ON short_links (short_code)
    WHERE revoked_at IS NULL;

CREATE TRIGGER short_links_updated_at
BEFORE UPDATE ON short_links
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- +migrate Down

DROP TRIGGER IF EXISTS short_links_updated_at ON short_links;

DROP TABLE IF EXISTS short_links;
