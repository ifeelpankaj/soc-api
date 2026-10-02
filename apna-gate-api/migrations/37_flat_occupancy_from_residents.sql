-- +migrate Up
-- Occupancy follows active flat-resident records. A blocked flat retains its
-- separate operational state until an admin unblocks it.
-- +migrate StatementBegin
CREATE FUNCTION reconcile_flat_occupancy(p_society_id BIGINT, p_flat_id BIGINT)
RETURNS void LANGUAGE plpgsql AS $$
BEGIN
  -- Serialize changes to one flat before counting residents. The count runs in
  -- a fresh statement snapshot after any prior resident change has committed.
  PERFORM 1 FROM flats WHERE id = p_flat_id AND society_id = p_society_id FOR UPDATE;
  UPDATE flats f
  SET status = CASE WHEN EXISTS (
    SELECT 1 FROM flat_residents r
    WHERE r.society_id = p_society_id AND r.flat_id = p_flat_id AND r.status = 'active'
  ) THEN 'occupied'::flat_status ELSE 'vacant'::flat_status END,
  updated_at = now()
  WHERE f.id = p_flat_id AND f.society_id = p_society_id AND f.status <> 'blocked'
    AND f.status IS DISTINCT FROM CASE WHEN EXISTS (
      SELECT 1 FROM flat_residents r
      WHERE r.society_id = p_society_id AND r.flat_id = p_flat_id AND r.status = 'active'
    ) THEN 'occupied'::flat_status ELSE 'vacant'::flat_status END;
END;
$$;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE FUNCTION sync_flat_occupancy_from_residents()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  IF TG_OP = 'DELETE' THEN
    PERFORM reconcile_flat_occupancy(OLD.society_id, OLD.flat_id);
    RETURN OLD;
  END IF;
  IF TG_OP = 'UPDATE' AND (OLD.society_id, OLD.flat_id) IS DISTINCT FROM (NEW.society_id, NEW.flat_id) THEN
    PERFORM reconcile_flat_occupancy(OLD.society_id, OLD.flat_id);
  END IF;
  PERFORM reconcile_flat_occupancy(NEW.society_id, NEW.flat_id);
  RETURN NEW;
END;
$$;
-- +migrate StatementEnd

CREATE TRIGGER flat_residents_occupancy_sync
AFTER INSERT OR DELETE OR UPDATE OF status, flat_id, society_id ON flat_residents
FOR EACH ROW EXECUTE FUNCTION sync_flat_occupancy_from_residents();

UPDATE flats f SET status = CASE WHEN EXISTS (
  SELECT 1 FROM flat_residents r
  WHERE r.society_id = f.society_id AND r.flat_id = f.id AND r.status = 'active'
) THEN 'occupied'::flat_status ELSE 'vacant'::flat_status END
WHERE f.status <> 'blocked';

-- +migrate Down
DROP TRIGGER flat_residents_occupancy_sync ON flat_residents;
DROP FUNCTION sync_flat_occupancy_from_residents();
DROP FUNCTION reconcile_flat_occupancy(BIGINT, BIGINT);
