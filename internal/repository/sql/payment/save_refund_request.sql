update payments set
    refund_requested_at = $2,
    refunded_at = $2,
    end_to_end_id = $3,
    status = 'refunded'
where booking_id = $1;
