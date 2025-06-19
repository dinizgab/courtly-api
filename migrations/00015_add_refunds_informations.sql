-- +goose Up
-- +goose StatementBegin
create extension if not exists pgcrypto;

alter table bookings
    add column cancel_token_hash text,
    add column cancel_token_expires_at timestamptz not null GENERATED ALWAYS AS ((start_time AT TIME ZONE 'UTC') - INTERVAL '3 hours') STORED;
alter table payments
    add column refund_requested_at timestamp,
    add column refunded_at timestamp,
    add column end_to_end_id text;

UPDATE bookings
SET cancel_token_hash = encode(gen_random_bytes(32), 'hex')
WHERE cancel_token_hash IS NULL;

alter table bookings
    alter column cancel_token_hash set not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table bookings
    drop column cancel_token_hash,
    drop column cancel_token_expires_at;
alter table payments
    drop column refund_requested_at,
    drop column refunded_at,
    drop column end_to_end_id;
-- +goose StatementEnd
