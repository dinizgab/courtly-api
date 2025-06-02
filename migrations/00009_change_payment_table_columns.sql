-- +goose Up
-- +goose StatementBegin
alter table payments
drop column if exists subaccount_id,
drop column if exists charge_id,
add column if not exists payment_link_id uuid,
add column if not exists payment_link_url text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table payments
add column subaccount_id uuid not null references openpix_subaccounts(id),
add column charge_id varchar(64) not null unique,
drop column if exists payment_link_id,
drop column if exists payment_link_url,
-- +goose StatementEnd
