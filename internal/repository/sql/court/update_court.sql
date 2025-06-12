UPDATE courts SET
    name = $1,
    description = $2,
    sport_type = $3,
    hourly_price = $4,
    is_active = $5,
    capacity = $6
WHERE id = $7
