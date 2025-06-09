SELECT
    id,
    court_id,
    start_time,
    end_time,
    created_at,
    status,
    guest_name,
    guest_email,
    guest_phone,
    verification_code
FROM
    bookings
WHERE
    court_id = $1
    and date(start_time) = current_date
