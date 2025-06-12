select
    id as photo_id,
    path
from court_photos
where court_id = $1
