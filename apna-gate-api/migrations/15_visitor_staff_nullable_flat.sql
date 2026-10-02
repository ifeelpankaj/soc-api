-- +migrate Up
ALTER TABLE visitor_entries
    ALTER COLUMN flat_id DROP NOT NULL;

-- +migrate Down
-- Staff entries with NULL flat_id must be removed or assigned before reverting.
UPDATE visitor_entries SET flat_id = (
    SELECT f.id FROM flats f WHERE f.society_id = visitor_entries.society_id AND f.is_active = TRUE LIMIT 1
) WHERE flat_id IS NULL;

ALTER TABLE visitor_entries
    ALTER COLUMN flat_id SET NOT NULL;
