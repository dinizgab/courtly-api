with upd_payments as (
    update payments set
        refund_requested_at = $2,
        refunded_at = $2,
        end_to_end_id = $3,
        status = 'refunded'
    where booking_id = $1
    returning booking_id
), update_bookings as (
    update bookings set
        status = 'cancelled'
    where id = (select booking_id from upd_payments)
        and status = 'confirmed'
    returning id
)
select 1
