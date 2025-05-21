SELECT id, name, address, phone, email, cnpj, slug
FROM companies
WHERE id = $1
