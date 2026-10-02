-- +migrate Up

ALTER TABLE plans
ADD COLUMN IF NOT EXISTS max_residents INT NOT NULL DEFAULT 100;

ALTER TABLE society_subscriptions
ADD COLUMN IF NOT EXISTS max_residents INT NOT NULL DEFAULT 100;

ALTER TABLE plans
DROP CONSTRAINT IF EXISTS plans_max_residents_check;

ALTER TABLE plans
ADD CONSTRAINT plans_max_residents_check CHECK (max_residents > 0);

ALTER TABLE society_subscriptions
DROP CONSTRAINT IF EXISTS society_subscriptions_max_residents_check;

ALTER TABLE society_subscriptions
ADD CONSTRAINT society_subscriptions_max_residents_check CHECK (max_residents > 0);

-- +migrate Down

ALTER TABLE society_subscriptions
DROP CONSTRAINT IF EXISTS society_subscriptions_max_residents_check;

ALTER TABLE plans
DROP CONSTRAINT IF EXISTS plans_max_residents_check;

ALTER TABLE society_subscriptions
DROP COLUMN IF EXISTS max_residents;

ALTER TABLE plans
DROP COLUMN IF EXISTS max_residents;
