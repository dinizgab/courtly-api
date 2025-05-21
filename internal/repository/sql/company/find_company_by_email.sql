select email, password_hash
from companies
where email = $1