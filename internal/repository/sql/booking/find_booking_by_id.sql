SELECT
    b.id,
    b.court_id,
    b.start_time,
    b.end_time,
    b.created_at,
    b.status,
    b.guest_name,
    b.guest_phone,
    b.guest_email,
    b.verification_code,
    b.total_price,
    c.name AS name
FROM
    bookings b
JOIN courts c
    ON b.court_id = c.id
WHERE
    b.id = $1;
