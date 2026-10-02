-- +migrate Up
ALTER TABLE maintenance_settings ADD COLUMN first_enabled_month DATE CHECK (extract(day FROM first_enabled_month)=1);
UPDATE maintenance_settings s SET first_enabled_month=COALESCE(
 (SELECT min(b.billing_month) FROM maintenance_bills b WHERE b.society_id=s.society_id),
 date_trunc('month', now() AT TIME ZONE COALESCE(NULLIF(s.config->>'timezone',''),'Asia/Kolkata'))::date);

CREATE TABLE maintenance_preview_reviews (
 token UUID PRIMARY KEY,
 society_id BIGINT NOT NULL REFERENCES societies(id),
 actor_id BIGINT NOT NULL REFERENCES users(id),
 billing_month DATE NOT NULL CHECK (extract(day FROM billing_month)=1),
 snapshot_hash TEXT NOT NULL,
 expires_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX maintenance_preview_expiry ON maintenance_preview_reviews(expires_at);

-- Existing issued snapshots are untouched. Catch-up metadata lives in new snapshots.
-- +migrate Down
DROP TABLE maintenance_preview_reviews;
ALTER TABLE maintenance_settings DROP COLUMN first_enabled_month;
