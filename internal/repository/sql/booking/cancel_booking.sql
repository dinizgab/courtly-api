UPDATE bookings SET
    status = 'cancelled'
WHERE
    id = $1
