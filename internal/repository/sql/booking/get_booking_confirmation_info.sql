SELECT
    b.guest_name,
    b.guest_phone,
    b.guest_email,
    c.name,
    co.address,
    b.start_time,
    b.end_time,
    b.total_price,
    b.verification_code
FROM
    bookings b
JOIN courts c
    ON b.court_id = c.id
JOIN companies co
    ON c.company_id = co.id
WHERE
    b.id = $1;
