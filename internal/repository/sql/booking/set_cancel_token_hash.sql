update bookings
set cancel_token_hash = $1
where id = $2
