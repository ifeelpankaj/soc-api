-- +migrate Up

ALTER TYPE visitor_event_type ADD VALUE IF NOT EXISTS 'guard_approved_on_behalf';

ALTER TABLE visitor_entries
    ADD COLUMN IF NOT EXISTS approved_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS delivery_partner VARCHAR(100),
    ADD COLUMN IF NOT EXISTS service_provider VARCHAR(100);

UPDATE visitor_entries
SET approved_at = updated_at
WHERE status IN ('approved', 'checked_in', 'checked_out')
  AND approved_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_visitor_entries_waiting_at_gate
    ON visitor_entries (society_id, approved_at ASC)
    WHERE status = 'approved';

CREATE INDEX IF NOT EXISTS idx_visitor_entries_delivery_partner
    ON visitor_entries (society_id, delivery_partner)
    WHERE delivery_partner IS NOT NULL;

-- +migrate Down

DROP INDEX IF EXISTS idx_visitor_entries_delivery_partner;
DROP INDEX IF EXISTS idx_visitor_entries_waiting_at_gate;

ALTER TABLE visitor_entries
    DROP COLUMN IF EXISTS service_provider,
    DROP COLUMN IF EXISTS delivery_partner,
    DROP COLUMN IF EXISTS approved_at;
