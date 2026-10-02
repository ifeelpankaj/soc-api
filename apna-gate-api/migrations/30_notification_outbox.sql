-- +migrate Up
CREATE TABLE notification_outbox (
 id BIGSERIAL PRIMARY KEY,
 user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 society_id BIGINT NOT NULL REFERENCES societies(id),
 flat_id BIGINT,
 audience TEXT NOT NULL CHECK (audience IN ('resident','admin','staff','society_resident')),
 event_key TEXT NOT NULL,
 payload JSONB NOT NULL,
 attempts INT NOT NULL DEFAULT 0,
 available_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 lease_token UUID,
 inbox_completed_at TIMESTAMPTZ,
 push_completed_at TIMESTAMPTZ,
 completed_at TIMESTAMPTZ,
 last_error TEXT,
 UNIQUE(user_id,event_key),
 FOREIGN KEY(flat_id,society_id) REFERENCES flats(id,society_id)
);
CREATE INDEX notification_outbox_ready ON notification_outbox(available_at,id) WHERE completed_at IS NULL;

-- Retain the legacy queue as a transactionally bridged producer for compatibility.
-- Existing event keys are preserved so inbox ON CONFLICT deduplication still applies.
-- +migrate StatementBegin
CREATE FUNCTION bridge_maintenance_outbox() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 INSERT INTO notification_outbox(user_id,society_id,flat_id,audience,event_key,payload)
 SELECT NEW.user_id,b.society_id,b.flat_id,NEW.audience,NEW.event_key,
 jsonb_build_object('Type',NEW.event_type,'Title','Maintenance update','Body','Open maintenance for bill '||b.bill_number||'.',
 'UserID',NEW.user_id,'SocietyID',b.society_id,'FlatID',b.flat_id,'EventKey',NEW.event_key,
 'Data',NEW.event_data||jsonb_build_object('type',NEW.event_type,'event_key',NEW.event_key,'bill_id',b.id::text,'society_id',b.society_id::text,'flat_id',b.flat_id::text))
 FROM maintenance_bills b WHERE b.id=NEW.bill_id
 ON CONFLICT(user_id,event_key) DO NOTHING;
 RETURN NEW;
END; $$;
-- +migrate StatementEnd
CREATE TRIGGER maintenance_outbox_bridge AFTER INSERT ON maintenance_notification_deliveries FOR EACH ROW EXECUTE FUNCTION bridge_maintenance_outbox();
INSERT INTO notification_outbox(user_id,society_id,flat_id,audience,event_key,payload,attempts,available_at,last_error)
SELECT d.user_id,b.society_id,b.flat_id,d.audience,d.event_key,
 jsonb_build_object('Type',d.event_type,'Title','Maintenance update','Body','Open maintenance for bill '||b.bill_number||'.',
 'UserID',d.user_id,'SocietyID',b.society_id,'FlatID',b.flat_id,'EventKey',d.event_key,
 'Data',d.event_data||jsonb_build_object('type',d.event_type,'event_key',d.event_key,'bill_id',b.id::text,'society_id',b.society_id::text,'flat_id',b.flat_id::text)),
 d.attempts,d.available_at,CASE WHEN d.last_error IS NOT NULL THEN 'delivery_failed' ELSE NULL END
FROM maintenance_notification_deliveries d JOIN maintenance_bills b ON b.id=d.bill_id WHERE d.completed_at IS NULL
ON CONFLICT(user_id,event_key) DO NOTHING;

-- +migrate Down
DROP TRIGGER maintenance_outbox_bridge ON maintenance_notification_deliveries;
DROP FUNCTION bridge_maintenance_outbox();
DROP TABLE notification_outbox;
