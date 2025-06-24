-- +goose Up
-- +goose StatementBegin
alter type payment_status add value 'refunding';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
