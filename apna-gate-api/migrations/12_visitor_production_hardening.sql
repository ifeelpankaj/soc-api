-- +migrate Up

CREATE UNIQUE INDEX IF NOT EXISTS uq_visitor_entries_invite_id
    ON visitor_entries (invite_id)
    WHERE invite_id IS NOT NULL;

ALTER TABLE society_visitor_settings
    ADD COLUMN IF NOT EXISTS allow_guard_on_behalf_approval BOOLEAN NOT NULL DEFAULT TRUE;

-- +migrate Down

ALTER TABLE society_visitor_settings
    DROP COLUMN IF EXISTS allow_guard_on_behalf_approval;

DROP INDEX IF EXISTS uq_visitor_entries_invite_id;
