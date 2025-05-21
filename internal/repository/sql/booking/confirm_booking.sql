UPDATE bookings SET
    status = 'confirmed'
WHERE
    id = $1 AND
    company_id = $2
