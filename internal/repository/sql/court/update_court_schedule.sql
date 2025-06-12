UPDATE court_schedules AS cs SET
	is_open = v.is_open,
	opening_time = v.opening_time,
	closing_time = v.closing_time
FROM (VALUES %s) AS v(id, is_open, opening_time, closing_time)
WHERE cs.id = v.id
