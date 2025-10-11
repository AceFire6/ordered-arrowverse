-- name: GetEpisodes :many
select show.name as show_name,
       show.slug as show_slug,
       ep.episode_id,
       ep.name,
       ep.air_date,
       ep.season,
       ep.episode,
       ep.description,
       ep.source_link,
       ep.ordering,
       ep.show_id,
       ep.created_at,
       ep.modified_at,
       ep.runtime_minutes
from episode as ep
         inner join show on show.show_id = ep.show_id
order by ordering;
