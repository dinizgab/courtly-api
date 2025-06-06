-- +goose Up
-- +goose StatementBegin
alter table payments add column qr_code_image text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table payments drop column qr_code_image;
-- +goose StatementEnd
