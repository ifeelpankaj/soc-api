-- +migrate Up
ALTER TABLE notification_outbox DROP CONSTRAINT notification_outbox_audience_check;
ALTER TABLE notification_outbox ADD CONSTRAINT notification_outbox_audience_check
 CHECK (audience IN ('resident','admin','staff','society_resident','hub_member'));
ALTER TABLE notification_outbox ADD COLUMN push_enabled BOOLEAN NOT NULL DEFAULT true;

-- +migrate Down
DELETE FROM notification_outbox WHERE audience='hub_member';
ALTER TABLE notification_outbox DROP COLUMN push_enabled;
ALTER TABLE notification_outbox DROP CONSTRAINT notification_outbox_audience_check;
ALTER TABLE notification_outbox ADD CONSTRAINT notification_outbox_audience_check
 CHECK (audience IN ('resident','admin','staff','society_resident'));
