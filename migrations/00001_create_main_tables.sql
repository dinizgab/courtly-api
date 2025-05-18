-- +goose Up
-- +goose StatementBegin
create type booking_status as enum (
    'pending',
    'confirmed',
    'cancelled'
);

CREATE TABLE IF NOT EXISTS companies (
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name    varchar(100) NOT NULL,
    address varchar(255) NOT NULL,
    phone   varchar(20) NOT NULL,
    email   varchar(100) NOT NULL,
    slug    varchar(50) NOT NULL unique
);

CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         varchar(100) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    company_id    UUID references companies(id)
);

CREATE TABLE IF NOT EXISTS courts (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id   UUID NOT NULL REFERENCES companies(id) on delete cascade,
    name         varchar(150) NOT NULL,
    sport_type   varchar(50) NOT NULL,
    hourly_price NUMERIC(10,2) NOT NULL CHECK (hourly_price >= 0),
    is_active    BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS bookings (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    court_id         UUID NOT NULL references courts(id) on delete cascade,
    start_time       TIMESTAMPTZ NOT NULL,
    end_time         TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status           booking_status NOT NULL DEFAULT 'pending',
    guest_name       varchar(100) NOT NULL,
    guest_phone      varchar(20) NOT NULL,
    guest_email      varchar(100) NOT NULL,
    verification_code varchar(6) NOT NULL,
    CONSTRAINT chk_time_range CHECK (end_time > start_time)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  IF EXISTS bookings;
DROP TABLE  IF EXISTS courts;
DROP TABLE  IF EXISTS users;
DROP TABLE  IF EXISTS companies;
-- +goose StatementEnd
