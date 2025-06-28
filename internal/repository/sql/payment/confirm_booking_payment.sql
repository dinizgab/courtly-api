with payment_confirmed as (
    update payments
    set status = 'paid',
        paid_at = $2,
        updated_at = now()
    where correlation_id = $1
        and status = 'pending'
    returning booking_id
)
update bookings b
set status = 'confirmed'
from payment_confirmed pc
where b.id = pc.booking_id
    and b.status = 'pending'

