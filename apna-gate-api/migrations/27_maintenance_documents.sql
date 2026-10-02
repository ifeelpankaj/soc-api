-- +migrate Up
-- sql-migrate runs this migration in one transaction. No concurrent bill write
-- can observe the brief, migration-only exception to the immutable-row trigger.
LOCK TABLE maintenance_bills IN ACCESS EXCLUSIVE MODE;
ALTER TABLE maintenance_bills ADD COLUMN issuer_snapshot JSONB;
ALTER TABLE maintenance_bills DISABLE TRIGGER maintenance_bill_immutable;
UPDATE maintenance_bills b SET issuer_snapshot=jsonb_build_object(
 'name',s.name,'society_code',s.society_code,
 'address_line1',s.address_line1,'address_line2',s.address_line2,
 'landmark',s.landmark,'city',s.city,'state',s.state,'pincode',s.pincode,'country',s.country,
 'email',s.email,'phone_number',s.phone_number,
 'captured_at',CURRENT_TIMESTAMP,'capture_source','rollout')
FROM societies s WHERE s.id=b.society_id;
ALTER TABLE maintenance_bills ENABLE TRIGGER maintenance_bill_immutable;
ALTER TABLE maintenance_bills ALTER COLUMN issuer_snapshot SET NOT NULL;
ALTER TABLE maintenance_bills ADD CONSTRAINT maintenance_issuer_snapshot_valid CHECK (
 jsonb_typeof(issuer_snapshot)='object' AND
 issuer_snapshot ?& ARRAY['name','society_code','captured_at','capture_source'] AND
 COALESCE(length(issuer_snapshot->>'name'),0)>0 AND
 COALESCE(length(issuer_snapshot->>'society_code'),0)>0 AND
 COALESCE(length(issuer_snapshot->>'captured_at'),0)>0 AND
 COALESCE(issuer_snapshot->>'capture_source','') IN ('rollout','issuance')
);

-- Always derive the snapshot from the society, including inserts made by old
-- application instances. Callers cannot provide an alternative issuer.
-- +migrate StatementBegin
CREATE FUNCTION capture_maintenance_issuer() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 SELECT jsonb_build_object(
  'name',s.name,'society_code',s.society_code,
  'address_line1',s.address_line1,'address_line2',s.address_line2,
  'landmark',s.landmark,'city',s.city,'state',s.state,'pincode',s.pincode,'country',s.country,
  'email',s.email,'phone_number',s.phone_number,
  'captured_at',CURRENT_TIMESTAMP,'capture_source','issuance')
 INTO NEW.issuer_snapshot FROM societies s WHERE s.id=NEW.society_id;
 RETURN NEW;
END; $$;
-- +migrate StatementEnd
CREATE TRIGGER maintenance_bill_capture_issuer BEFORE INSERT ON maintenance_bills
FOR EACH ROW EXECUTE FUNCTION capture_maintenance_issuer();

-- +migrate Down
DROP TRIGGER maintenance_bill_capture_issuer ON maintenance_bills;
DROP FUNCTION capture_maintenance_issuer();
ALTER TABLE maintenance_bills DROP COLUMN issuer_snapshot;
