select id, correlation_id, booking_id, paid_at, value_total
from payments
where booking_id = $1;
