SELECT
    b.id,
    b.start_time,
    b.end_time,
    b.created_at,
    b.status,
    b.guest_name,
    b.guest_phone,
    b.guest_email,
    c.name,
    c.hourly_price
FROM
    bookings b
JOIN courts c
    ON c.id = b.court_id
    AND c.company_id = $1
WHERE
    b.company_id = $1
