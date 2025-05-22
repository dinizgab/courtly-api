-- +goose Up
-- +goose StatementBegin
alter table bookings add column total_price decimal(10, 2) not null default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table bookings drop column total_price;
-- +goose StatementEnd
