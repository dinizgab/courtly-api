SELECT
    id,
    company_id,
    name,
    description,
    sport_type,
    hourly_price,
    is_active,
    opening_time,
    closing_time,
    capacity
FROM
    courts
WHERE
    id = $1;
