SELECT (id, name, address, phone, email, slug)
FROM companies
WHERE slug = $1
