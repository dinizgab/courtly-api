-- +goose Up
-- +goose StatementBegin
create extension if not exists btree_gist;

alter table bookings
add constraint no_overlapping_bookings
exclude using gist (
    court_id with =,
    tstzrange(start_time, end_time) with &&
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table bookings
drop constraint no_overlapping_bookings;
-- +goose StatementEnd
