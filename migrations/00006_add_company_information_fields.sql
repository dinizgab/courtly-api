-- +goose Up
-- +goose StatementBegin
ALTER TABLE companies
    ADD COLUMN password_hash varchar(255),
    ADD COLUMN cnpj varchar(14);

update companies set password_hash = '$2a$14$0RITU48n0y6u3ymb4WL3guPfOKaylRNcB7bEFpzSl0TR7VT35FeI2' where password_hash is null;
update companies set cnpj = '54234052000151' where id = '11111111-1111-1111-1111-111111111111'::uuid;
update companies set cnpj = '12345678000195' where id = '22222222-2222-2222-2222-222222222222'::uuid;

ALTER TABLE companies
    ALTER COLUMN cnpj SET NOT NULL,
    ALTER COLUMN password_hash SET NOT NULL;

ALTER TABLE companies
    ADD CONSTRAINT companies_email_unique UNIQUE (email),
    ADD CONSTRAINT companies_cnpj_unique UNIQUE (cnpj);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE companies
    DROP CONSTRAINT IF EXISTS companies_email_unique,
    DROP CONSTRAINT IF EXISTS companies_cnpj_unique,
    DROP COLUMN     IF EXISTS password_hash,
    DROP COLUMN     IF EXISTS cnpj;
-- +goose StatementEnd
