with expired as (
    update payments
    set status = 'expired'
    where correlation_id = $1
    returning booking_id
)
update bookings
set status = 'cancelled'
where id in (select booking_id from expired)