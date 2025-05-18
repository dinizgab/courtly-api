SELECT (id, company_id, name, sport_type, hourly_price, is_active)
FROM courts
WHERE id = $1;
