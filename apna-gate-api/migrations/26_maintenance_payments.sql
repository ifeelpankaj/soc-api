-- +migrate Up
ALTER TABLE maintenance_bills ADD CONSTRAINT maintenance_bill_tenant UNIQUE(id,society_id);
ALTER TABLE maintenance_bills ADD CONSTRAINT maintenance_bill_amount UNIQUE(id,society_id,total_paise);

CREATE TABLE maintenance_payment_settings_versions (
 society_id BIGINT NOT NULL REFERENCES societies(id) ON DELETE RESTRICT,
 version BIGINT NOT NULL CHECK(version>0), enabled BOOLEAN NOT NULL,
 upi_id VARCHAR(255) NOT NULL, payee_name VARCHAR(120) NOT NULL,
 created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY(society_id,version)
);
CREATE TABLE maintenance_payment_settings (
 society_id BIGINT PRIMARY KEY REFERENCES societies(id) ON DELETE RESTRICT,
 version BIGINT NOT NULL,
 FOREIGN KEY(society_id,version) REFERENCES maintenance_payment_settings_versions(society_id,version)
);
CREATE TABLE maintenance_payment_requests (
 id UUID PRIMARY KEY, society_id BIGINT NOT NULL, bill_id BIGINT NOT NULL,
 settings_version BIGINT NOT NULL, amount_paise BIGINT NOT NULL CHECK(amount_paise>0),
 reference VARCHAR(40) NOT NULL UNIQUE,
 state TEXT NOT NULL DEFAULT 'active' CHECK(state IN ('active','superseded','closed')),
 created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 FOREIGN KEY(bill_id,society_id,amount_paise) REFERENCES maintenance_bills(id,society_id,total_paise),
 FOREIGN KEY(society_id,settings_version) REFERENCES maintenance_payment_settings_versions(society_id,version),
 UNIQUE(id,society_id,bill_id), UNIQUE(id,society_id)
);
CREATE UNIQUE INDEX maintenance_one_active_request ON maintenance_payment_requests(society_id,bill_id) WHERE state='active';
CREATE TABLE maintenance_payment_claims (
 id BIGSERIAL PRIMARY KEY, society_id BIGINT NOT NULL, bill_id BIGINT NOT NULL,
 request_id UUID NOT NULL, submitted_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 reference VARCHAR(64) NOT NULL, payment_date DATE NOT NULL,
 status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','verified','rejected','cancelled')),
 reviewed_by BIGINT REFERENCES users(id) ON DELETE RESTRICT, reviewed_at TIMESTAMPTZ, reason TEXT,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 FOREIGN KEY(request_id,society_id,bill_id) REFERENCES maintenance_payment_requests(id,society_id,bill_id),
 UNIQUE(id,society_id,bill_id)
);
CREATE UNIQUE INDEX maintenance_one_pending_claim ON maintenance_payment_claims(society_id,bill_id) WHERE status='pending';
CREATE INDEX maintenance_claims_reference ON maintenance_payment_claims(society_id,reference,id DESC);
CREATE TABLE maintenance_payments (
 id BIGSERIAL PRIMARY KEY, society_id BIGINT NOT NULL, bill_id BIGINT NOT NULL,
 claim_id BIGINT, payer_id BIGINT REFERENCES users(id) ON DELETE RESTRICT,
 settings_version BIGINT NOT NULL, amount_paise BIGINT NOT NULL CHECK(amount_paise>0),
 reference VARCHAR(64) NOT NULL, credit_date DATE NOT NULL,
 receipt_number UUID NOT NULL UNIQUE,
 verified_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 verified_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 status TEXT NOT NULL DEFAULT 'verified' CHECK(status IN ('verified','reversed')),
 reversed_by BIGINT REFERENCES users(id) ON DELETE RESTRICT, reversed_at TIMESTAMPTZ, reversal_reason TEXT,
 CHECK((status='verified' AND reversed_by IS NULL AND reversed_at IS NULL AND reversal_reason IS NULL) OR
       (status='reversed' AND reversed_by IS NOT NULL AND reversed_at IS NOT NULL AND reversal_reason IS NOT NULL AND length(trim(reversal_reason))>0)),
 FOREIGN KEY(bill_id,society_id,amount_paise) REFERENCES maintenance_bills(id,society_id,total_paise),
 FOREIGN KEY(claim_id,society_id,bill_id) REFERENCES maintenance_payment_claims(id,society_id,bill_id),
 FOREIGN KEY(society_id,settings_version) REFERENCES maintenance_payment_settings_versions(society_id,version),
 UNIQUE(id,society_id,bill_id)
);
CREATE UNIQUE INDEX maintenance_one_settlement ON maintenance_payments(society_id,bill_id) WHERE status='verified';
CREATE UNIQUE INDEX maintenance_verified_reference ON maintenance_payments(society_id,reference) WHERE status='verified';
CREATE UNIQUE INDEX maintenance_verified_claim ON maintenance_payments(claim_id) WHERE claim_id IS NOT NULL;
CREATE TABLE maintenance_payment_reference_reservations (
 society_id BIGINT NOT NULL, reference VARCHAR(64) NOT NULL, bill_id BIGINT NOT NULL,
 claim_id BIGINT, payment_id BIGINT,
 PRIMARY KEY(society_id,reference), CHECK(num_nonnulls(claim_id,payment_id)=1),
 FOREIGN KEY(bill_id,society_id) REFERENCES maintenance_bills(id,society_id),
 FOREIGN KEY(claim_id,society_id,bill_id) REFERENCES maintenance_payment_claims(id,society_id,bill_id),
 FOREIGN KEY(payment_id,society_id,bill_id) REFERENCES maintenance_payments(id,society_id,bill_id)
);
CREATE TABLE maintenance_payment_ledger (
 id BIGSERIAL PRIMARY KEY, society_id BIGINT NOT NULL, bill_id BIGINT NOT NULL, payment_id BIGINT NOT NULL,
 kind TEXT NOT NULL CHECK(kind IN ('collection','reversal')), amount_paise BIGINT NOT NULL,
 actor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 CHECK((kind='collection' AND amount_paise>0) OR (kind='reversal' AND amount_paise<0)),
 FOREIGN KEY(payment_id,society_id,bill_id) REFERENCES maintenance_payments(id,society_id,bill_id),
 UNIQUE(payment_id,kind)
);
CREATE TABLE maintenance_payment_audit_events (
 id BIGSERIAL PRIMARY KEY, society_id BIGINT NOT NULL REFERENCES societies(id) ON DELETE RESTRICT,
 bill_id BIGINT, actor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 action TEXT NOT NULL, entity_id TEXT NOT NULL, details JSONB NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 FOREIGN KEY(bill_id,society_id) REFERENCES maintenance_bills(id,society_id)
);
CREATE TABLE maintenance_payment_reports (
 id BIGSERIAL PRIMARY KEY, society_id BIGINT NOT NULL, bill_id BIGINT NOT NULL, request_id UUID,
 reported_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 reference VARCHAR(64) NOT NULL, amount_paise BIGINT NOT NULL CHECK(amount_paise>0),
 payment_date DATE NOT NULL, explanation TEXT NOT NULL, fingerprint TEXT NOT NULL,
 status TEXT NOT NULL DEFAULT 'open' CHECK(status IN ('open','investigating','resolved')),
 resolution_note TEXT, updated_by BIGINT REFERENCES users(id) ON DELETE RESTRICT,
 updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 CHECK(status<>'resolved' OR (resolution_note IS NOT NULL AND length(trim(resolution_note))>0 AND updated_by IS NOT NULL)),
 FOREIGN KEY(bill_id,society_id) REFERENCES maintenance_bills(id,society_id),
 FOREIGN KEY(request_id,society_id,bill_id) REFERENCES maintenance_payment_requests(id,society_id,bill_id),
 UNIQUE(society_id,reported_by,fingerprint)
);
CREATE TABLE maintenance_payment_idempotency (
 society_id BIGINT NOT NULL REFERENCES societies(id) ON DELETE RESTRICT,
 actor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
 operation TEXT NOT NULL, key VARCHAR(128) NOT NULL, request_hash TEXT NOT NULL,
 response JSONB NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 PRIMARY KEY(society_id,actor_id,operation,key)
);

ALTER TABLE maintenance_payment_claims ADD CHECK(reference=upper(btrim(reference)) AND reference<>'');
ALTER TABLE maintenance_payments ADD CHECK(reference=upper(btrim(reference)) AND reference<>'');
ALTER TABLE maintenance_payment_reference_reservations ADD CHECK(reference=upper(btrim(reference)) AND reference<>'');
ALTER TABLE maintenance_payment_reports ADD CHECK(reference=upper(btrim(reference)) AND reference<>'');

-- Immutable snapshots, audit and ledger. Payments allow only a one-way reversal.
CREATE TRIGGER maintenance_settings_version_immutable BEFORE UPDATE OR DELETE ON maintenance_payment_settings_versions FOR EACH ROW EXECUTE FUNCTION reject_maintenance_bill_mutation();
CREATE TRIGGER maintenance_ledger_immutable BEFORE UPDATE OR DELETE ON maintenance_payment_ledger FOR EACH ROW EXECUTE FUNCTION reject_maintenance_bill_mutation();
CREATE TRIGGER maintenance_audit_immutable BEFORE UPDATE OR DELETE ON maintenance_payment_audit_events FOR EACH ROW EXECUTE FUNCTION reject_maintenance_bill_mutation();
-- +migrate StatementBegin
CREATE FUNCTION guard_maintenance_payment() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 IF TG_OP='DELETE' THEN RAISE EXCEPTION 'Payments cannot be deleted' USING ERRCODE='23514'; END IF;
 IF OLD.status<>'verified' OR NEW.status<>'reversed' OR
 (to_jsonb(OLD)-ARRAY['status','reversed_by','reversed_at','reversal_reason']) IS DISTINCT FROM
 (to_jsonb(NEW)-ARRAY['status','reversed_by','reversed_at','reversal_reason']) THEN
 RAISE EXCEPTION 'Only a one-way payment reversal is permitted' USING ERRCODE='23514'; END IF;
 RETURN NEW;
END; $$;
-- +migrate StatementEnd
CREATE TRIGGER maintenance_payment_guard BEFORE UPDATE OR DELETE ON maintenance_payments FOR EACH ROW EXECUTE FUNCTION guard_maintenance_payment();

-- +migrate StatementBegin
CREATE FUNCTION guard_maintenance_request() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 IF TG_OP='DELETE' THEN RAISE EXCEPTION 'Payment requests cannot be deleted' USING ERRCODE='23514'; END IF;
 IF (to_jsonb(OLD)-'state') IS DISTINCT FROM (to_jsonb(NEW)-'state') OR
 NOT ((OLD.state='active' AND NEW.state IN ('superseded','closed')) OR (OLD.state='superseded' AND NEW.state='closed')) THEN
 RAISE EXCEPTION 'Invalid payment request transition' USING ERRCODE='23514'; END IF;
 RETURN NEW;
END; $$;
-- +migrate StatementEnd
CREATE TRIGGER maintenance_request_guard BEFORE UPDATE OR DELETE ON maintenance_payment_requests FOR EACH ROW EXECUTE FUNCTION guard_maintenance_request();
-- +migrate StatementBegin
CREATE FUNCTION guard_maintenance_claim() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 IF TG_OP='DELETE' THEN RAISE EXCEPTION 'Payment claims cannot be deleted' USING ERRCODE='23514'; END IF;
 IF (to_jsonb(OLD)-ARRAY['status','reviewed_by','reviewed_at','reason']) IS DISTINCT FROM
 (to_jsonb(NEW)-ARRAY['status','reviewed_by','reviewed_at','reason']) OR OLD.status<>'pending' OR
 NEW.status NOT IN ('verified','rejected','cancelled') OR NEW.reviewed_by IS NULL OR NEW.reviewed_at IS NULL THEN
 RAISE EXCEPTION 'Invalid claim transition' USING ERRCODE='23514'; END IF;
 RETURN NEW;
END; $$;
-- +migrate StatementEnd
CREATE TRIGGER maintenance_claim_guard BEFORE UPDATE OR DELETE ON maintenance_payment_claims FOR EACH ROW EXECUTE FUNCTION guard_maintenance_claim();

CREATE VIEW maintenance_bill_balances AS
SELECT b.id AS bill_id, b.society_id,
 COALESCE(p.amount_paise,0)::bigint AS paid_amount_paise,
 (b.total_paise-COALESCE(p.amount_paise,0))::bigint AS outstanding_amount_paise,
 COALESCE((SELECT c.status FROM maintenance_payment_claims c WHERE c.bill_id=b.id ORDER BY c.id DESC LIMIT 1),'none')::text AS payment_claim_status,
 CASE WHEN p.id IS NOT NULL THEN 'paid' WHEN (now() AT TIME ZONE b.timezone)::date>b.due_date THEN 'overdue' ELSE 'unpaid' END::text AS payment_status
FROM maintenance_bills b LEFT JOIN maintenance_payments p ON p.bill_id=b.id AND p.society_id=b.society_id AND p.status='verified';

ALTER TABLE maintenance_notification_deliveries ADD COLUMN event_type TEXT NOT NULL DEFAULT 'maintenance_bill_generated';
ALTER TABLE maintenance_notification_deliveries ADD COLUMN event_key TEXT;
UPDATE maintenance_notification_deliveries SET event_key='maintenance_bill_generated:'||bill_id;
ALTER TABLE maintenance_notification_deliveries ALTER COLUMN event_key SET NOT NULL;
ALTER TABLE maintenance_notification_deliveries ALTER COLUMN event_key SET DEFAULT '';
ALTER TABLE maintenance_notification_deliveries ADD COLUMN audience TEXT NOT NULL DEFAULT 'resident' CHECK(audience IN ('resident','admin'));
ALTER TABLE maintenance_notification_deliveries ADD COLUMN event_data JSONB NOT NULL DEFAULT '{}';
ALTER TABLE maintenance_notification_deliveries DROP CONSTRAINT maintenance_notification_deliveries_bill_id_user_id_key;
ALTER TABLE maintenance_notification_deliveries ADD CONSTRAINT maintenance_delivery_event UNIQUE(user_id,event_key);

-- +migrate Down
DELETE FROM maintenance_notification_deliveries WHERE event_type<>'maintenance_bill_generated';
ALTER TABLE maintenance_notification_deliveries DROP CONSTRAINT maintenance_delivery_event;
ALTER TABLE maintenance_notification_deliveries ADD UNIQUE(bill_id,user_id);
ALTER TABLE maintenance_notification_deliveries DROP COLUMN event_type, DROP COLUMN event_key, DROP COLUMN audience, DROP COLUMN event_data;
DROP VIEW maintenance_bill_balances;
DROP TABLE maintenance_payment_idempotency,maintenance_payment_reports,maintenance_payment_audit_events,maintenance_payment_ledger,maintenance_payment_reference_reservations,maintenance_payments,maintenance_payment_claims,maintenance_payment_requests,maintenance_payment_settings,maintenance_payment_settings_versions;
DROP FUNCTION guard_maintenance_payment();
DROP FUNCTION guard_maintenance_request(),guard_maintenance_claim();
ALTER TABLE maintenance_bills DROP CONSTRAINT maintenance_bill_amount, DROP CONSTRAINT maintenance_bill_tenant;
