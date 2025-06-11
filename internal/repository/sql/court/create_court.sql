INSERT INTO courts(company_id,
                   name,
                   description,
                   sport_type,
                   hourly_price,
                   is_active,
                   capacity)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7)
RETURNING id;
