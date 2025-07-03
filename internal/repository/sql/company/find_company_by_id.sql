SELECT
    c.id,
    name,
    address,
    phone,
    email,
    cnpj,
    slug,
    pix_key,
    pix_key_type
FROM
    companies c
JOIN openpix_subaccounts os
    ON c.id = os.company_id
WHERE
    c.id = $1
