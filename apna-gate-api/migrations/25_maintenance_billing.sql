-- +migrate Up
ALTER TABLE flats ADD COLUMN flat_type VARCHAR(50) CHECK (flat_type IS NULL OR length(trim(flat_type)) > 0);
ALTER TABLE flats ADD COLUMN area_sqft_hundredths BIGINT CHECK (area_sqft_hundredths > 0);

CREATE TABLE maintenance_settings (
    society_id BIGINT PRIMARY KEY REFERENCES societies(id) ON DELETE RESTRICT,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    config JSONB NOT NULL,
    updated_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE maintenance_type_rates (
    society_id BIGINT NOT NULL REFERENCES maintenance_settings(society_id) ON DELETE CASCADE,
    flat_type VARCHAR(50) NOT NULL,
    amount_paise BIGINT NOT NULL CHECK (amount_paise > 0),
    PRIMARY KEY (society_id, flat_type)
);
CREATE TABLE maintenance_billing_runs (
    id BIGSERIAL PRIMARY KEY,
    society_id BIGINT NOT NULL REFERENCES societies(id) ON DELETE RESTRICT,
    billing_month DATE NOT NULL CHECK (extract(day from billing_month) = 1),
    snapshot JSONB NOT NULL,
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (society_id, billing_month),
    UNIQUE (id, society_id, billing_month)
);
CREATE TABLE maintenance_bills (
    id BIGSERIAL PRIMARY KEY,
    run_id BIGINT NOT NULL,
    society_id BIGINT NOT NULL,
    flat_id BIGINT NOT NULL,
    billing_month DATE NOT NULL,
    bill_number VARCHAR(100) NOT NULL,
    due_date DATE NOT NULL,
    timezone VARCHAR(100) NOT NULL,
    total_paise BIGINT NOT NULL CHECK (total_paise > 0),
    snapshot JSONB NOT NULL,
    billed_party JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (run_id, society_id, billing_month) REFERENCES maintenance_billing_runs(id, society_id, billing_month) ON DELETE RESTRICT,
    FOREIGN KEY (flat_id, society_id) REFERENCES flats(id, society_id) ON DELETE RESTRICT,
    UNIQUE (society_id, flat_id, billing_month),
    UNIQUE (society_id, bill_number)
);
CREATE INDEX maintenance_bills_list ON maintenance_bills(society_id, id DESC);
CREATE INDEX maintenance_bills_flat ON maintenance_bills(society_id, flat_id, id DESC);
CREATE TABLE maintenance_bill_items (
    bill_id BIGINT NOT NULL REFERENCES maintenance_bills(id) ON DELETE RESTRICT,
    position INT NOT NULL,
    description TEXT NOT NULL,
    amount_paise BIGINT NOT NULL CHECK (amount_paise >= 0),
    PRIMARY KEY (bill_id, position)
);
CREATE TABLE maintenance_notification_deliveries (
    id BIGSERIAL PRIMARY KEY,
    bill_id BIGINT NOT NULL REFERENCES maintenance_bills(id) ON DELETE RESTRICT,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    attempts INT NOT NULL DEFAULT 0,
    available_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    lease_token UUID,
    completed_at TIMESTAMPTZ,
    last_error TEXT,
    UNIQUE (bill_id, user_id)
);
CREATE INDEX maintenance_delivery_pending ON maintenance_notification_deliveries(available_at) WHERE completed_at IS NULL;

-- Issued charges are append-only; future adjustments/payments use separate records.
-- +migrate StatementBegin
CREATE FUNCTION reject_maintenance_bill_mutation() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'Issued maintenance bills and items are immutable' USING ERRCODE = '23514';
END;
$$;
-- +migrate StatementEnd
CREATE TRIGGER maintenance_bill_immutable BEFORE UPDATE OR DELETE ON maintenance_bills
FOR EACH ROW EXECUTE FUNCTION reject_maintenance_bill_mutation();
CREATE TRIGGER maintenance_item_immutable BEFORE UPDATE OR DELETE ON maintenance_bill_items
FOR EACH ROW EXECUTE FUNCTION reject_maintenance_bill_mutation();

-- +migrate Down
DROP TABLE maintenance_notification_deliveries;
DROP TABLE maintenance_bill_items;
DROP TABLE maintenance_bills;
DROP FUNCTION reject_maintenance_bill_mutation();
DROP TABLE maintenance_billing_runs;
DROP TABLE maintenance_type_rates;
DROP TABLE maintenance_settings;
ALTER TABLE flats DROP COLUMN area_sqft_hundredths, DROP COLUMN flat_type;
