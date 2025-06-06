select brcode, qr_code_image
from payments
where booking_id = $1
