-- +goose Up
-- +goose StatementBegin
create table court_schedules (
    id uuid primary key default gen_random_uuid(),
    court_id uuid not null references courts(id) on delete cascade,
    day_of_week smallint not null,
    opening_time time not null,
    closing_time time not null,
    is_open boolean not null,
    UNIQUE (court_id, day_of_week)
);

alter table courts drop column opening_time;
alter table courts drop column closing_time;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table court_schedules;
alter table courts add column opening_time time not null default '08:00:00';
alter table courts add column closing_time time not null default '22:00:00';
-- +goose StatementEnd
