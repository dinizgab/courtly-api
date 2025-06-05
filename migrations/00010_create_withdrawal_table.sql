-- +goose Up
-- +goose StatementBegin
create table if not exists withdrawals (
    id uuid primary key default gen_random_uuid(),
    company_id uuid references companies(id) on delete cascade,
    correlation_id uuid not null,
    value bigint not null,
    destination_alias varchar(255) not null,
    created_at timestamptz not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists withdrawals;
-- +goose StatementEnd
