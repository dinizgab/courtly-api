INSERT INTO companies (name, address, phone, email, password_hash, cnpj, slug)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id
