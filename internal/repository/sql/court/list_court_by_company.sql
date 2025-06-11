SELECT
    c.id,
    c.company_id,
    c.name,
    c.sport_type,
    c.hourly_price,
    c.is_active,
    coalesce(b.bookings_today, 0) as bookings_today
FROM courts c
LEFT JOIN  (
    select court_id, count(*) as bookings_today
    from bookings
    where start_time::date = CURRENT_DATE
    group by court_id
) b on c.id = b.court_id
WHERE c.company_id = $1
