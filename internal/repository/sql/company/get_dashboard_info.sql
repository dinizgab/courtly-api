SELECT
    COALESCE(SUM(p.value_total), 0) AS total_earning,
    COALESCE(SUM(EXTRACT(EPOCH FROM (b.end_time - b.start_time)) / 3600.0), 0) AS total_booked_time,
    COALESCE(COUNT(b.id), 0) AS total_bookings,
    COALESCE(COUNT(b.guest_email), 0) AS total_guests
FROM
    bookings b
JOIN payments p on b.id = p.booking_id
WHERE
    b.company_id = $1
    AND b.start_time >= date_trunc('week', now())
    AND b.start_time < date_trunc('week', now() + INTERVAL '1 week')
    AND b.status = 'confirmed'
    AND p.status = 'paid'
