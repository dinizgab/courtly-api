SELECT c.id, c.company_id, c.name, c.sport_type, c.hourly_price, c.is_active, count(b.id) as bookings_today
FROM courts c
LEFT JOIN bookings b ON b.court_id = c.id AND b.start_time = CURRENT_DATE
WHERE c.company_id = $1
GROUP BY c.id, c.company_id, c.name, c.sport_type, c.hourly_price, c.is_active
