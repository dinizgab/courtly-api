SELECT
    c.id as court_id,
    c.company_id,
    c.name,
    c.description,
    c.sport_type,
    c.hourly_price,
    c.is_active,
    c.opening_time,
    c.closing_time,
    c.capacity,
    cp.id as photo_id,
    cp.path
FROM
    courts c
JOIN
    court_photos cp ON c.id = cp.court_id
WHERE
    c.id = $1;
