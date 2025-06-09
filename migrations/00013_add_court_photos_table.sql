-- +goose Up
-- +goose StatementBegin
create table court_photos (
    id uuid primary key default gen_random_uuid(),
    court_id uuid references courts(id) on delete cascade,
    path text not null,
    is_cover boolean default false,
    position int
);

create index on court_photos (court_id, position);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table court_photos;
-- +goose StatementEnd
