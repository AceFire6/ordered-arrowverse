-- name: GetShowList :many
select slug,
       name,
       data_sources
from show
order by show_id;
