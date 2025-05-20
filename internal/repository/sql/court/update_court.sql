UPDATE courts SET
    name = $1,
    description = $2,
    sport_type = $3,
    hourly_price = $4,
    is_active = $5,
    opening_time = $6,
    closing_time = $7,
    capacity = $8
WHERE id = $9