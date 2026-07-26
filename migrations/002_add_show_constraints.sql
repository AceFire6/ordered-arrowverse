-- Add unique constraints required by ON CONFLICT specs in the seed SQL.
-- The legacy 001 schema declared slug but did not enforce uniqueness, so
-- the harvested seed needs an explicit unique key to upsert show rows.

create unique index show_slug_unique on show (slug);

---- create above / drop below ----

drop index show_slug_unique;
