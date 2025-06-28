SELECT
    b.id,
    b.guest_name,
    b.guest_phone,
    b.guest_email,
    c.name,
    co.address,
    b.start_time,
    b.end_time,
    p.value_total,
    b.verification_code,
    b.cancel_token_hash
FROM
    bookings b
JOIN courts c
    ON b.court_id = c.id
JOIN companies co
    ON c.company_id = co.id
JOIN payments p
    ON b.id = p.booking_id
WHERE
    b.id = $1;
