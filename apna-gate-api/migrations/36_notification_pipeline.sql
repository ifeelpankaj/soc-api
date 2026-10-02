-- +migrate Up
ALTER TABLE notification_outbox ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now();
ALTER TABLE notifications ADD COLUMN domain TEXT NOT NULL DEFAULT 'system'
  CHECK (domain IN ('visitor','maintenance','hub','system'));
UPDATE notifications SET domain = CASE
  WHEN type LIKE 'visitor%' OR type LIKE 'member_invite%' THEN 'visitor'
  WHEN type LIKE 'maintenance%' THEN 'maintenance'
  WHEN type LIKE 'hub.%' THEN 'hub'
  ELSE 'system' END;

CREATE TABLE notification_preferences (
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  society_id BIGINT NOT NULL REFERENCES societies(id) ON DELETE CASCADE,
  visitor_updates_push BOOLEAN NOT NULL DEFAULT true,
  maintenance_push BOOLEAN NOT NULL DEFAULT true,
  announcements_push BOOLEAN NOT NULL DEFAULT true,
  hub_replies_push BOOLEAN NOT NULL DEFAULT true,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, society_id)
);

CREATE TABLE push_deliveries (
  id BIGSERIAL PRIMARY KEY,
  notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
  device_token_id BIGINT REFERENCES device_tokens(id) ON DELETE SET NULL,
  provider TEXT NOT NULL CHECK (provider IN ('fcm','apns')),
  audience TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','processing','retry','sent','dead')),
  attempt_count INT NOT NULL DEFAULT 0,
  lifetime_attempt_count INT NOT NULL DEFAULT 0,
  replay_count INT NOT NULL DEFAULT 0,
  next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  lease_until TIMESTAMPTZ,
  lease_token UUID,
  provider_message_id TEXT,
  last_error_code TEXT,
  last_error TEXT,
  sent_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(notification_id, device_token_id)
);
CREATE INDEX push_deliveries_ready ON push_deliveries(next_attempt_at,id)
  WHERE status IN ('pending','retry','processing');

CREATE TABLE hub_notification_fanout_jobs (
  id BIGSERIAL PRIMARY KEY,
  society_id BIGINT NOT NULL REFERENCES societies(id) ON DELETE CASCADE,
  post_id BIGINT NOT NULL,
  actor_user_id BIGINT NOT NULL REFERENCES users(id),
  important BOOLEAN NOT NULL DEFAULT false,
  last_user_id BIGINT NOT NULL DEFAULT 0,
  lease_until TIMESTAMPTZ,
  lease_token UUID,
  completed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(society_id, post_id)
);
CREATE INDEX hub_notification_fanout_ready ON hub_notification_fanout_jobs(id)
  WHERE completed_at IS NULL;

-- Maintenance's legacy-named table is a domain event staging table. Its
-- lifecycle ends when the shared inbox and device work are materialized.
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION complete_legacy_maintenance_delivery() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 IF NEW.inbox_completed_at IS NOT NULL AND OLD.inbox_completed_at IS NULL THEN
  UPDATE maintenance_notification_deliveries SET completed_at=NEW.inbox_completed_at,last_error=NULL
  WHERE user_id=NEW.user_id AND event_key=NEW.event_key AND completed_at IS NULL;
 END IF;
 RETURN NEW;
END; $$;
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION complete_legacy_maintenance_delivery() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 IF NEW.completed_at IS NOT NULL AND OLD.completed_at IS NULL THEN
  UPDATE maintenance_notification_deliveries SET completed_at=NEW.completed_at,last_error=NULL
  WHERE user_id=NEW.user_id AND event_key=NEW.event_key AND completed_at IS NULL;
 END IF;
 RETURN NEW;
END; $$;
-- +migrate StatementEnd
DROP TABLE hub_notification_fanout_jobs;
DROP TABLE push_deliveries;
DROP TABLE notification_preferences;
ALTER TABLE notifications DROP COLUMN domain;
ALTER TABLE notification_outbox DROP COLUMN created_at;
