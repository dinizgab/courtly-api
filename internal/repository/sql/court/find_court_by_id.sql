SELECT
    c.id as court_id,
    c.company_id,
    c.name,
    c.description,
    c.sport_type,
    c.hourly_price,
    c.is_active,
    c.capacity,
    cs.day_of_week,
    cs.opening_time,
    cs.closing_time
FROM
    courts c
JOIN court_schedules cs ON c.id = cs.court_id
WHERE
    c.id = $1
