-- +goose Up
-- +goose StatementBegin
ALTER TABLE courts ADD COLUMN description TEXT;
ALTER TABLE courts ADD COLUMN capacity INT;
ALTER TABLE courts ADD COLUMN opening_hour TIME;
ALTER TABLE courts ADD COLUMN closing_hour TIME;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE courts DROP COLUMN description;
ALTER TABLE courts DROP COLUMN capacity;
ALTER TABLE courts DROP COLUMN opening_hour;
ALTER TABLE courts DROP COLUMN closing_hour;
-- +goose StatementEnd
