-- +goose Up
CREATE TABLE users
(
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    login         VARCHAR(256) NOT NULL UNIQUE,
    password_hash VARCHAR(256) NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE orders
(
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    number      TEXT        NOT NULL UNIQUE,
    user_id     BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status      TEXT        NOT NULL DEFAULT 'NEW' CHECK ( status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED') ),
    accrual     BIGINT CHECK (accrual IS NULL OR accrual >= 0),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()

);

CREATE TABLE withdrawals
(
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id      BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    order_number TEXT        NOT NULL,
    amount       BIGINT      NOT NULL CHECK ( amount > 0 ), -- cents --
    processed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_user_uploaded_at
    ON orders (user_id, uploaded_at DESC);

CREATE INDEX idx_withdrawals_user_processed_at
    ON withdrawals (user_id, processed_at DESC);

-- +goose Down
DROP TABLE withdrawals;
DROP TABLE orders;
DROP TABLE users;
