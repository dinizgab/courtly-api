SELECT cancel_token_hash, cancel_token_expires_at
FROM bookings
WHERE id = $1;
