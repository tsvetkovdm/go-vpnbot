CREATE SCHEMA vpnbot;
SET search_path TO vpnbot;

CREATE TABLE users (
    id          BIGSERIAL    PRIMARY KEY,
    version     BIGINT       NOT NULL DEFAULT 1,
    telegram_id BIGINT       NOT NULL UNIQUE,
    username    VARCHAR(32),
    email       VARCHAR(255) NOT NULL UNIQUE,
    xray_uuid   UUID         NOT NULL UNIQUE,
    status      VARCHAR(20)  NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'blocked')),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE admins (
    id          BIGSERIAL   PRIMARY KEY,
    telegram_id BIGINT      NOT NULL UNIQUE,
    is_super    BOOL        NOT NULL DEFAULT false,
    added_by    BIGINT      REFERENCES admins(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE plans (
    id         BIGSERIAL     PRIMARY KEY,
    version    BIGINT        NOT NULL DEFAULT 1,
    title      VARCHAR(100)  NOT NULL,
    price      NUMERIC(10,2) NOT NULL CHECK (price > 0),
    duration   INTERVAL      NOT NULL CHECK (duration > interval '0'),
    is_active  BOOL          NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE TABLE subscriptions(
    id                 BIGSERIAL    PRIMARY KEY,
    version            BIGINT       NOT NULL DEFAULT 1,
    user_id            BIGINT       NOT NULL UNIQUE REFERENCES users(id),
    plan_id            BIGINT       NOT NULL REFERENCES plans(id),
    status             VARCHAR(20)  NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'expired')),
    expires_at         TIMESTAMPTZ,
    billing_anchor_day SMALLINT CHECK (billing_anchor_day BETWEEN 1 AND 31),
    notified_3_day_at  TIMESTAMPTZ,
    notified_1_day_at  TIMESTAMPTZ,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE payments(
    id                   BIGSERIAL      PRIMARY KEY,
    user_id              BIGINT         NOT NULL REFERENCES users(id),
    plan_id              BIGINT                  REFERENCES plans(id) ON DELETE SET NULL,
    plan_title_snapshot  VARCHAR(100)   NOT NULL,
    price_snapshot       NUMERIC(10, 2) NOT NULL CHECK (price_snapshot > 0),
    duration_snapshot    INTERVAL       NOT NULL CHECK (duration_snapshot > interval '0'),
    status               VARCHAR(20)    NOT NULL DEFAULT 'created' CHECK (status IN ('created', 'pending', 'successful', 'rejected', 'expired', 'cancelled', 'failed')),
    resolved_by_admin_id BIGINT                  REFERENCES admins(id),
    resolved_at          TIMESTAMPTZ,
    created_at           TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uq_payments_one_unresolved_per_user
    ON payments (user_id)
    WHERE status IN ('created', 'pending');

CREATE TABLE admin_audit_log (
    id          BIGSERIAL   PRIMARY KEY,
    admin_id    BIGINT      NOT NULL REFERENCES admins(id),
    action      VARCHAR(50) NOT NULL CHECK (action IN (
                    'user_blocked', 'user_unblocked',
                    'payment_confirmed', 'payment_rejected',
                    'subscription_extended_manually',
                    'plan_created', 'plan_activated', 'plan_deactivated', 'plan_deleted',
                    'admin_added', 'admin_removed'
                )),
    target_type VARCHAR(20) NOT NULL CHECK (target_type IN ('user', 'payment', 'subscription', 'plan', 'admin')),
    target_id   BIGINT      NOT NULL,
    comment     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);