CREATE TYPE payment_status as enum ('PENDING','AUTHORIZED','PARTIALLY_CAPTURED','CAPTURED','PARTIALLY_REFUNDED','REFUNDED','VOIDED');
CREATE TYPE payment_type as enum ('AUTHORIZATION','CAPTURE','REFUND','VOID');

CREATE EXTENSION "uuid-ossp";

CREATE TABLE IF NOT EXISTS payment_action
(
    id            UUID UNIQUE DEFAULT uuid_generate_v4(),
    amount        int          NOT NULL,
    payment_type  payment_type NOT NULL,
    response_code VARCHAR(4),
    created_at    timestamptz default now(),
    processed_at  timestamptz
);

CREATE TABLE IF NOT EXISTS payment
(
    id         UUID UNIQUE DEFAULT uuid_generate_v4(),
    amount     int            NOT NULL,
    currency   VARCHAR(3)     NOT NULL,
    status     payment_status NOT NULL,
    action_id  UUID references payment_action (id),
    created_at timestamptz default now(),
    updated_at timestamptz
);



