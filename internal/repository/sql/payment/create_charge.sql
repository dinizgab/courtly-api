insert into payments (
    company_id,
    booking_id,
    correlation_id,
    payment_link_id,
    payment_link_url,
    qr_code_image,
    brcode,
    value_total,
    expires_at
) values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)

