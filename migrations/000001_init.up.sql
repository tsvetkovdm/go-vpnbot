CREATE SCHEMA vpnbot;
SET search_path TO vpnbot;

CREATE TABLE users (
    id          BIGSERIAL    PRIMARY KEY,
    version     BIGINT       NOT NULL DEFAULT 1,
    telegram_id BIGINT       NOT NULL UNIQUE,
    username    VARCHAR(32),
    first_name  VARCHAR(64)  NOT NULL,
    email       VARCHAR(255) NOT NULL UNIQUE,
    xray_uuid   UUID         NOT NULL UNIQUE,
    status      VARCHAR(20)  NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'blocked')),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE admins (
    id          BIGSERIAL   PRIMARY KEY,
    telegram_id BIGINT      NOT NULL UNIQUE,
    role        VARCHAR(20) NOT NULL DEFAULT 'admin' CHECK (role IN ('super', 'admin', 'support')),
    added_by    BIGINT      REFERENCES admins(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    removed_at  TIMESTAMPTZ 
);

CREATE TABLE plans (
    id         BIGSERIAL     PRIMARY KEY,
    version    BIGINT        NOT NULL DEFAULT 1,
    title      VARCHAR(100)  NOT NULL,
    is_active  BOOL          NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE TABLE plan_prices (
    id         BIGSERIAL     PRIMARY KEY,
    version    BIGINT        NOT NULL DEFAULT 1,
    plan_id    BIGINT        NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    duration   INTERVAL      NOT NULL CHECK (duration > interval '0'),
    price      NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    is_active  BOOL          NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uq_plan_prices_active_duration
    ON plan_prices (plan_id, duration) 
    WHERE is_active;

CREATE TABLE subscriptions(
    id                 BIGSERIAL    PRIMARY KEY,
    version            BIGINT       NOT NULL DEFAULT 1,
    user_id            BIGINT       NOT NULL UNIQUE REFERENCES users(id),
    expires_at         TIMESTAMPTZ  NOT NULL,
    anchor_day         SMALLINT     NOT NULL CHECK (anchor_day BETWEEN 1 AND 31),
    notified_3_day_at  TIMESTAMPTZ,
    notified_1_day_at  TIMESTAMPTZ,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE payments(
    id                   BIGSERIAL      PRIMARY KEY,
    user_id              BIGINT         NOT NULL REFERENCES users(id),
    plan_price_id        BIGINT         NOT NULL REFERENCES plan_prices(id),
    plan_title_snapshot  VARCHAR(100)   NOT NULL,
    price_snapshot       NUMERIC(10, 2) NOT NULL CHECK (price_snapshot >= 0),
    duration_snapshot    INTERVAL       NOT NULL CHECK (duration_snapshot > interval '0'),
    comment_code         SMALLINT       NOT NULL CHECK (comment_code BETWEEN 1000 AND 9999),
    status               VARCHAR(20)    NOT NULL DEFAULT 'created' CHECK (status IN ('created', 'pending', 'successful', 'rejected', 'expired', 'cancelled')),
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
                    'user_blocked', 'user_unblocked', 'user_deleted',
                    'payment_confirmed', 'payment_rejected',
                    'subscription_extended_manually',
                    'plan_created', 'plan_activated', 'plan_deactivated', 'plan_deleted',
                    'plan_price_added', 'plan_price_activated', 'plan_price_deactivated', 'plan_price_deleted',
                    'admin_added', 'admin_removed'
                )),
    target_type VARCHAR(20) NOT NULL CHECK (target_type IN ('user', 'payment', 'subscription', 'plan', 'plan_price', 'admin')),
    target_id   BIGINT      NOT NULL,
    comment     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE connections (
    id            BIGSERIAL    PRIMARY KEY,
    title         VARCHAR(64)  NOT NULL,
    inbound_tag   VARCHAR(64)  NOT NULL,
    link_template TEXT         NOT NULL CHECK (
                      link_template LIKE 'vless://{uuid}@%'
                      AND position('#' IN link_template) = 0
                  ),
    is_public     BOOL         NOT NULL DEFAULT false,
    is_active     BOOL         NOT NULL DEFAULT false,
    sort_order    INT          NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE access_grants (
    id                  BIGSERIAL   PRIMARY KEY,
    user_id             BIGINT      NOT NULL REFERENCES users(id),
    connection_id       BIGINT      NOT NULL REFERENCES connections(id),
    expires_at          TIMESTAMPTZ NOT NULL,
    granted_by_admin_id BIGINT      REFERENCES admins(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, connection_id)
);

CREATE TABLE feed_tokens (
    user_id    BIGINT      PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);