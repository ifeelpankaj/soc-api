-- +migrate Up
-- The legacy producer and shared outbox remain in the same transaction.
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION bridge_maintenance_outbox() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
 reminder boolean := NEW.event_type='maintenance_payment_reminder';
BEGIN
 INSERT INTO notification_outbox(user_id,society_id,flat_id,audience,event_key,payload)
 SELECT NEW.user_id,b.society_id,b.flat_id,NEW.audience,NEW.event_key,
 jsonb_build_object('Type',NEW.event_type,
 'Title',CASE WHEN reminder THEN 'Maintenance payment reminder' ELSE 'Maintenance update' END,
 'Body',CASE WHEN reminder THEN format('Your %s maintenance bill %s has INR %s outstanding. Due date: %s. Open the bill for details.',
 to_char(b.billing_month,'YYYY-MM'),b.bill_number,(v.outstanding_amount_paise::numeric/100)::numeric(22,2),to_char(b.due_date,'YYYY-MM-DD'))
 ELSE 'Open maintenance for bill '||b.bill_number||'.' END,
 'UserID',NEW.user_id,'SocietyID',b.society_id,'FlatID',b.flat_id,'EventKey',NEW.event_key,
 'Data',NEW.event_data||jsonb_build_object('type',NEW.event_type,'event_key',NEW.event_key,'bill_id',b.id::text,'society_id',b.society_id::text,'flat_id',b.flat_id::text)
 ||CASE WHEN reminder THEN jsonb_build_object('billing_month',to_char(b.billing_month,'YYYY-MM'),'due_date',to_char(b.due_date,'YYYY-MM-DD'),'outstanding_amount_paise',v.outstanding_amount_paise::text) ELSE '{}'::jsonb END)
 FROM maintenance_bills b JOIN maintenance_bill_balances v ON v.bill_id=b.id WHERE b.id=NEW.bill_id
 ON CONFLICT(user_id,event_key) DO NOTHING;
 RETURN NEW;
END; $$;
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION bridge_maintenance_outbox() RETURNS trigger LANGUAGE plpgsql AS $$
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
