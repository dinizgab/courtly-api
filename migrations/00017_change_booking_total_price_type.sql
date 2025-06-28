-- +goose Up
-- +goose StatementBegin
alter table bookings
  alter column total_price type bigint using round(total_price * 100);

alter table courts
    alter column hourly_price type bigint using round(hourly_price * 100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table bookings
  alter column total_price type numeric(10, 2) using round(total_price / 100);

alter table courts
    alter column hourly_price type numeric(10, 2) using round(hourly_price / 100);
-- +goose StatementEnd
