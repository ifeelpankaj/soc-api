-- +migrate Up

DROP INDEX IF EXISTS idx_visitor_entries_waiting_at_gate;

CREATE INDEX IF NOT EXISTS idx_visitor_entries_waiting_at_gate
    ON visitor_entries (society_id, approved_at ASC)
    WHERE status = 'approved'
      AND checked_in_at IS NULL
      AND source <> 'resident_link';

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_visitors_full_name_trgm
    ON visitors USING gin (full_name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_visitors_phone_number_trgm
    ON visitors USING gin (phone_number gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_flats_flat_number_trgm
    ON flats USING gin (flat_number gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_societies_name_trgm
    ON societies USING gin (name gin_trgm_ops);

-- +migrate Down

DROP INDEX IF EXISTS idx_societies_name_trgm;
DROP INDEX IF EXISTS idx_flats_flat_number_trgm;
DROP INDEX IF EXISTS idx_visitors_phone_number_trgm;
DROP INDEX IF EXISTS idx_visitors_full_name_trgm;

CREATE INDEX IF NOT EXISTS idx_visitor_entries_waiting_at_gate
    ON visitor_entries (society_id, approved_at ASC)
    WHERE status = 'approved';
