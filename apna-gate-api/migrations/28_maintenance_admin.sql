-- +migrate Up
ALTER TABLE maintenance_payments ADD COLUMN evidence_reference text NOT NULL DEFAULT ''
 CHECK (char_length(evidence_reference) <= 500 AND evidence_reference = btrim(evidence_reference));
-- The existing payment guard compares all immutable columns, including evidence_reference.

-- +migrate Down
ALTER TABLE maintenance_payments DROP COLUMN evidence_reference;
