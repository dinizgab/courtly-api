SELECT
    c.id,
    c.name,
    c.description,
    c.sport_type,
    c.hourly_price,
    c.is_active,
    c.capacity,
    cs.opening_time,
    cs.closing_time,
    cp.id AS photo_id,
    cp.path
FROM
    courts c
LEFT JOIN
    court_photos cp ON c.id = cp.court_id
LEFT JOIN court_schedules cs ON cs.court_id = c.id
    AND cs.day_of_week = EXTRACT(DOW FROM CURRENT_DATE)
WHERE
    company_id = $1
