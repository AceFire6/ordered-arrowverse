package db

import (
	"context"
	"time"
)

// Hand-maintained to mirror sqlc-generated output for the
// get_filtered_episode_list.sql query.
//
// sqlc is configured to emit a `_sqlc`-suffixed file per-query, but this
// filter query lives alongside Querier and shares GetEpisodesRow, so the
// generated file gets hand-edited and lives at filtered_episodes.go
// (no `_sqlc` suffix) to prevent `sqlc generate` from clobbering it.
//
// When the DB comes back online and `sqlc generate` is rerun, manually
// reconcile this file with the generated equivalent.
//
// See get_filtered_episode_list.sql for the source SQL.
const getEpisodesFiltered = `-- name: GetEpisodesFiltered :many
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
where (cardinality($1::text[]) = 0 or not (show.slug = any ($1)))
  and ($2::timestamptz is null or ep.air_date >= $2)
  and ($3::timestamptz is null or ep.air_date <= $3)
order by ordering
`

// GetEpisodesFiltered returns every episode optionally narrowed by show
// slugs to exclude and/or an inclusive air-date window. Pass an empty
// slice and nil times for an unfiltered query.
func (q *Queries) GetEpisodesFiltered(
	ctx context.Context,
	hideShowSlugs []string,
	fromDate *time.Time,
	toDate *time.Time,
) ([]GetEpisodesRow, error) {
	rows, err := q.db.Query(ctx, getEpisodesFiltered, hideShowSlugs, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []GetEpisodesRow
	for rows.Next() {
		var i GetEpisodesRow
		if err := rows.Scan(
			&i.ShowName,
			&i.ShowSlug,
			&i.EpisodeID,
			&i.Name,
			&i.AirDate,
			&i.Season,
			&i.Episode,
			&i.Description,
			&i.SourceLink,
			&i.Ordering,
			&i.ShowID,
			&i.CreatedAt,
			&i.ModifiedAt,
			&i.RuntimeMinutes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
