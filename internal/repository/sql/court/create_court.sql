INSERT INTO courts(
    company_id,
    name,
    sport_type,
    hourly_price,
    is_active
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);
