-- +goose Up
-- +goose StatementBegin
create type pix_key_type as enum ('cpf', 'cnpj', 'email', 'phone', 'random');

alter table openpix_subaccounts
    add column pix_key_type pix_key_type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table openpix_subaccounts
    drop column pix_key_type;
-- +goose StatementEnd
