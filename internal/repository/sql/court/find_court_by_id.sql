SELECT
    c.id as court_id,
    c.company_id,
    c.name,
    c.description,
    c.sport_type,
    c.hourly_price,
    c.is_active,
    c.capacity
FROM
    courts c
WHERE
    c.id = $1
