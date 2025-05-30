-- +goose Up
-- +goose StatementBegin
ALTER TABLE courts ADD COLUMN description TEXT NOT NULL DEFAULT '';
ALTER TABLE courts ADD COLUMN opening_time TIME NOT NULL DEFAULT '08:00:00';
ALTER TABLE courts ADD COLUMN closing_time TIME NOT NULL DEFAULT '22:00:00';
ALTER TABLE courts ADD COLUMN capacity INTEGER DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE courts DROP COLUMN description;
ALTER TABLE courts DROP COLUMN opening_time;
ALTER TABLE courts DROP COLUMN closing_time;
ALTER TABLE courts DROP COLUMN capacity;
-- +goose StatementEnd
