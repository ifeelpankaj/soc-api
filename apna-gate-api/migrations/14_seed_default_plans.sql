-- +migrate Up

INSERT INTO plans (
    name,
    code,
    description,
    price_amount_paise,
    currency,
    billing_cycle,
    max_flats,
    max_admins,
    max_staff,
    max_residents,
    features,
    is_active
)
VALUES
    (
        'Mobile Starter',
        'mobile_starter',
        'Starter plan for small societies with mobile-first operations.',
        0,
        'INR',
        'monthly',
        10,
        2,
        2,
        100,
        '{"visitors": true, "claims": true}'::jsonb,
        TRUE
    ),
    (
        'Mobile Growth',
        'mobile_growth',
        'Growth plan for larger societies with expanded limits.',
        0,
        'INR',
        'monthly',
        50,
        5,
        10,
        500,
        '{"visitors": true, "claims": true, "audit_logs": true}'::jsonb,
        TRUE
    )
ON CONFLICT (code) DO NOTHING;

-- +migrate Down

DELETE FROM plans WHERE code IN ('mobile_starter', 'mobile_growth');
