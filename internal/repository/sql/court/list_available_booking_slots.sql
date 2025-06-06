select start_time, end_time
from bookings
where court_id = $1
    and date(start_time) = date($2)
    and status not in ('cancelled', 'pending')
