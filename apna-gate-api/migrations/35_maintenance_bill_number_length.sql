-- +migrate Up
ALTER TABLE maintenance_bills ALTER COLUMN bill_number TYPE VARCHAR(180);

-- +migrate Down
ALTER TABLE maintenance_bills ALTER COLUMN bill_number TYPE VARCHAR(100);
