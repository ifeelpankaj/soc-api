-- +migrate Up
CREATE TABLE maintenance_billing_audit (
 id BIGSERIAL PRIMARY KEY,
 society_id BIGINT NOT NULL REFERENCES societies(id),
 actor_id BIGINT REFERENCES users(id),
 event_type TEXT NOT NULL,
 entity_id BIGINT NOT NULL,
 details JSONB NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TRIGGER maintenance_billing_audit_immutable BEFORE UPDATE OR DELETE ON maintenance_billing_audit FOR EACH ROW EXECUTE FUNCTION reject_maintenance_bill_mutation();
CREATE TRIGGER maintenance_run_immutable BEFORE UPDATE OR DELETE ON maintenance_billing_runs FOR EACH ROW EXECUTE FUNCTION reject_maintenance_bill_mutation();
CREATE TRIGGER maintenance_idempotency_immutable BEFORE UPDATE OR DELETE ON maintenance_payment_idempotency FOR EACH ROW EXECUTE FUNCTION reject_maintenance_bill_mutation();
-- +migrate StatementBegin
CREATE FUNCTION audit_maintenance_billing() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 IF TG_TABLE_NAME='maintenance_billing_runs' THEN
  INSERT INTO maintenance_billing_audit(society_id,actor_id,event_type,entity_id,details) VALUES(NEW.society_id,NEW.created_by,'billing_issued',NEW.id,NEW.snapshot||jsonb_build_object('billing_month',NEW.billing_month,'outcome','issued'));
 ELSE
  IF TG_OP='UPDATE' AND OLD.first_enabled_month IS NOT NULL AND NEW.first_enabled_month IS DISTINCT FROM OLD.first_enabled_month THEN
   RAISE EXCEPTION 'First enabled month is immutable' USING ERRCODE='23514';
  END IF;
  INSERT INTO maintenance_billing_audit(society_id,actor_id,event_type,entity_id,details) VALUES(NEW.society_id,NEW.updated_by,'billing_settings_changed',NEW.society_id,NEW.config||jsonb_build_object('first_enabled_month',NEW.first_enabled_month,'outcome','saved'));
 END IF;
 RETURN NEW;
END; $$;
-- +migrate StatementEnd
CREATE TRIGGER maintenance_run_audit AFTER INSERT ON maintenance_billing_runs FOR EACH ROW EXECUTE FUNCTION audit_maintenance_billing();
CREATE TRIGGER maintenance_configuration_audit AFTER INSERT OR UPDATE ON maintenance_settings FOR EACH ROW EXECUTE FUNCTION audit_maintenance_billing();

-- +migrate StatementBegin
CREATE FUNCTION complete_legacy_maintenance_delivery() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 IF NEW.completed_at IS NOT NULL AND OLD.completed_at IS NULL THEN
  UPDATE maintenance_notification_deliveries SET completed_at=NEW.completed_at,last_error=NULL
  WHERE user_id=NEW.user_id AND event_key=NEW.event_key AND completed_at IS NULL;
 END IF;
 RETURN NEW;
END; $$;
-- +migrate StatementEnd
CREATE TRIGGER notification_outbox_legacy_completion AFTER UPDATE ON notification_outbox FOR EACH ROW EXECUTE FUNCTION complete_legacy_maintenance_delivery();

-- +migrate Down
DROP TRIGGER notification_outbox_legacy_completion ON notification_outbox;
DROP FUNCTION complete_legacy_maintenance_delivery();
DROP TRIGGER maintenance_configuration_audit ON maintenance_settings;
DROP TRIGGER maintenance_run_audit ON maintenance_billing_runs;
DROP FUNCTION audit_maintenance_billing();
DROP TRIGGER maintenance_idempotency_immutable ON maintenance_payment_idempotency;
DROP TRIGGER maintenance_run_immutable ON maintenance_billing_runs;
DROP TABLE maintenance_billing_audit;
