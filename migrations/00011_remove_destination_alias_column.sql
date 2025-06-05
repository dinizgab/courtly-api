-- +goose Up
-- +goose StatementBegin
alter table withdrawals
drop column if exists destination_alias;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table withdrawals
add column destination_alias varchar(255) not null;
-- +goose StatementEnd
