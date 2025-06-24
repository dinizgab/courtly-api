select id, correlation_id, booking_id, paid_at
from payments
where booking_id = $1;
