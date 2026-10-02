-- +migrate Up

CREATE TYPE billing_cycle AS ENUM (
    'monthly',
    'yearly'
);

CREATE TYPE subscription_status AS ENUM (
    'pending',
    'trial',
    'active',
    'expired',
    'cancelled'
);

CREATE TABLE plans (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,

    description TEXT,

    price_amount_paise BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'INR',
    billing_cycle billing_cycle NOT NULL,

    max_flats INT NOT NULL,
    max_admins INT NOT NULL DEFAULT 1,
    max_staff INT NOT NULL DEFAULT 2,
    max_residents INT NOT NULL DEFAULT 100,

    features JSONB NOT NULL DEFAULT '{}'::jsonb,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT plans_name_unique UNIQUE (name),
    CONSTRAINT plans_code_unique UNIQUE (code),

    CONSTRAINT plans_price_check CHECK (price_amount_paise >= 0),
    CONSTRAINT plans_max_flats_check CHECK (max_flats > 0),
    CONSTRAINT plans_max_admins_check CHECK (max_admins >= 0),
    CONSTRAINT plans_max_staff_check CHECK (max_staff >= 0),
    CONSTRAINT plans_max_residents_check CHECK (max_residents > 0),
    CONSTRAINT plans_currency_check CHECK (currency <> ''),
    CONSTRAINT plans_code_check CHECK (code <> ''),
    CONSTRAINT plans_name_check CHECK (name <> '')
);

CREATE TABLE society_subscriptions (
    id BIGSERIAL PRIMARY KEY,

    society_id BIGINT NOT NULL REFERENCES societies(id) ON DELETE CASCADE,
    plan_id BIGINT REFERENCES plans(id) ON DELETE SET NULL,

    status subscription_status NOT NULL DEFAULT 'pending',

    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,
    trial_ends_at TIMESTAMPTZ,

    -- Plan snapshot
    plan_name VARCHAR(100) NOT NULL,
    plan_code VARCHAR(50) NOT NULL,

    price_amount_paise BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'INR',
    billing_cycle billing_cycle NOT NULL,

    max_flats INT NOT NULL,
    max_admins INT NOT NULL,
    max_staff INT NOT NULL,
    max_residents INT NOT NULL,

    features JSONB NOT NULL DEFAULT '{}'::jsonb,

    -- Audit
    activated_at TIMESTAMPTZ,
    activated_by BIGINT REFERENCES users(id) ON DELETE SET NULL,

    expired_at TIMESTAMPTZ,

    cancelled_at TIMESTAMPTZ,
    cancelled_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    cancellation_reason TEXT,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT society_subscriptions_price_check CHECK (price_amount_paise >= 0),
    CONSTRAINT society_subscriptions_max_flats_check CHECK (max_flats > 0),
    CONSTRAINT society_subscriptions_max_admins_check CHECK (max_admins >= 0),
    CONSTRAINT society_subscriptions_max_staff_check CHECK (max_staff >= 0),
    CONSTRAINT society_subscriptions_max_residents_check CHECK (max_residents > 0),
    CONSTRAINT society_subscriptions_currency_check CHECK (currency <> ''),
    CONSTRAINT society_subscriptions_plan_name_check CHECK (plan_name <> ''),
    CONSTRAINT society_subscriptions_plan_code_check CHECK (plan_code <> ''),

    CONSTRAINT society_subscriptions_date_check CHECK (
        ends_at IS NULL
        OR starts_at IS NULL
        OR ends_at > starts_at
    ),

    CONSTRAINT society_subscriptions_trial_date_check CHECK (
        trial_ends_at IS NULL
        OR starts_at IS NULL
        OR trial_ends_at >= starts_at
    )
);

CREATE UNIQUE INDEX uniq_one_open_subscription_per_society
ON society_subscriptions (society_id)
WHERE status IN ('pending', 'trial', 'active');

CREATE INDEX idx_plans_active
ON plans (is_active);

CREATE INDEX idx_plans_code
ON plans (code);

CREATE INDEX idx_society_subscriptions_society_id
ON society_subscriptions (society_id);

CREATE INDEX idx_society_subscriptions_plan_id
ON society_subscriptions (plan_id);

CREATE INDEX idx_society_subscriptions_status
ON society_subscriptions (status);

CREATE INDEX idx_society_subscriptions_ends_at
ON society_subscriptions (ends_at);

CREATE INDEX idx_society_subscriptions_active_lookup
ON society_subscriptions (society_id, status, ends_at)
WHERE status IN ('trial', 'active');

CREATE TRIGGER plans_updated_at
BEFORE UPDATE ON plans
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER society_subscriptions_updated_at
BEFORE UPDATE ON society_subscriptions
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- +migrate Down

DROP TRIGGER IF EXISTS society_subscriptions_updated_at ON society_subscriptions;
DROP TRIGGER IF EXISTS plans_updated_at ON plans;

DROP INDEX IF EXISTS idx_society_subscriptions_active_lookup;
DROP INDEX IF EXISTS idx_society_subscriptions_ends_at;
DROP INDEX IF EXISTS idx_society_subscriptions_status;
DROP INDEX IF EXISTS idx_society_subscriptions_plan_id;
DROP INDEX IF EXISTS idx_society_subscriptions_society_id;
DROP INDEX IF EXISTS idx_plans_code;
DROP INDEX IF EXISTS idx_plans_active;
DROP INDEX IF EXISTS uniq_one_open_subscription_per_society;

DROP TABLE IF EXISTS society_subscriptions;
DROP TABLE IF EXISTS plans;

DROP TYPE IF EXISTS subscription_status;
DROP TYPE IF EXISTS billing_cycle;
