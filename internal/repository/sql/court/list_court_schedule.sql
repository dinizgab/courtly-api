SELECT
    id,
    is_open,
    day_of_week,
    opening_time,
    closing_time
FROM
    court_schedules
WHERE
    court_id = $1
