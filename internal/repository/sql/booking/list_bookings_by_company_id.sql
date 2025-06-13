SELECT
    b.id,
    b.start_time,
    b.end_time,
    b.created_at,
    b.status,
    b.guest_name,
    b.guest_phone,
    b.guest_email,
    c.name
FROM
    bookings b
JOIN courts c
    ON c.id = b.court_id
WHERE
    b.company_id = $1
    and b.start_time >= coalesce($2, b.start_time)
    and b.end_time <= coalesce($3, b.end_time)
    and b.status = 'confirmed'
