-- +goose Up
-- +goose StatementBegin
ALTER TABLE bookings ADD COLUMN company_id uuid,
ADD CONSTRAINT fk_bookings_company
FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE bookings
DROP CONSTRAINT fk_bookings_company,
DROP COLUMN company_id;
-- +goose StatementEnd
