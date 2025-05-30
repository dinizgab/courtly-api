-- +goose Up
-- +goose StatementBegin
create type payment_status as enum (
    'pending',
    'paid',
    'expired',
    'mismatch', -- value_received != value_expected
    'refunded'  -- full refund
);

create table if not exists openpix_subaccounts (
    id uuid primary key default gen_random_uuid(),
    company_id uuid not null references companies(id) on delete cascade,
    pix_key varchar(120) not null unique,
    created_at timestamptz not null default now()
);

create table if not exists payments (
    id uuid primary key default gen_random_uuid(),

    booking_id uuid not null references bookings(id) on delete cascade,
    company_id uuid not null references companies(id) on delete cascade,
    subaccount_id uuid not null references openpix_subaccounts(id),
    correlation_id varchar(64) not null unique,
    charge_id varchar(64) not null unique,
    brcode text not null,

    value_total bigint not null,
    value_commission bigint not null default 300,
    value_company bigint generated always as (value_total - value_commission) stored, -- value_total - value_commission

    status payment_status not null default 'pending',

    expires_at timestamptz,
    paid_at timestamptz,

    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

CREATE INDEX payments_booking_idx   ON payments (booking_id);
CREATE INDEX payments_company_idx   ON payments (company_id);
CREATE INDEX payments_status_idx    ON payments (status) where status = 'pending';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists payments;
drop type if exists payment_status;
-- +goose StatementEnd
