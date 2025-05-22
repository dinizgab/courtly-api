SELECT
    b.start_time,
    b.end_time,
    c.name,
    c.hourly_price,
    co.address
FROM
    bookings b
JOIN courts c
    ON b.court_id = c.id
JOIN companies co
    ON c.company_id = co.id
WHERE
    b.id = $1
