SELECT
    id,
    court_id,
    start_time,
    end_time,
    created_at,
    status,
    guest_name,
    guest_phone,
    guest_email,
    verification_code
FROM
    bookings
WHERE
    id = $1;
