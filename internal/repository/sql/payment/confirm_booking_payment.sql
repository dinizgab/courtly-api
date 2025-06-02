update payments
set status = 'paid',
    paid_at = $2,
    updated_at = now()
where correlation_id = $1
  and status = 'pending'

