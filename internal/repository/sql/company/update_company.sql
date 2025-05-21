update companies
set
    name = $1,
    address = $2,
    phone = $3,
    email = $4,
    cnpj = $5,
    slug = $6
where id = $7
