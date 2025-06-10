SELECT
    c.id,
    c.name,
    c.description,
    c.sport_type,
    c.hourly_price,
    c.is_active,
    c.opening_time,
    c.closing_time,
    c.capacity,
    cp.id AS photo_id,
    cp.path
FROM
    courts c
JOIN
    court_photos cp ON c.id = cp.court_id
WHERE
    company_id = $1
