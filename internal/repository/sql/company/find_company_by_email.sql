select id, email, password_hash
from companies
where email = $1
