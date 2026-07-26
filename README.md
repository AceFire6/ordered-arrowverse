# Arrowverse Series Ordering

A Go rewrite of the long-running [arrowverse.info](https://arrowverse.info)
episode-ordering tracker. Lists the canonical air-date ordering of the
Arrowverse shows so crossover episodes never spoil anyone again.

The legacy Python (Quart + BeautifulSoup + Redis) implementation has been
fully replaced. The Go service is the supported runtime going forward.

## Features

- Filtered home view (hide shows, air-date range, newest-first toggle) served
  via htmx partial swaps
- `/api` endpoint returning the filtered episode list as JSON
- `/recent_episodes.atom` Atom feed for RSS subscribers
- `/legal/privacy-policy` and `/legal/cookie-policy` static pages
- Legacy `/newest_first/`, `/hide/...` URLs 301-redirect to the canonical
  home view
- CSRF protection (Echo middleware + htmx meta tag header injection)
- Conditional UA vs GA4 analytics + AdSense, switched on the legacy-site
  banner detection
- Healthcheck endpoints (`/health/ping` pings Postgres + Redis)
- Reproducible Docker build (multi-stage, distroless runtime)

## How It Works

The database holds shows and episodes as static rows seeded at build time
via `cmd/harvest`. The service ships:

- Episodes and show metadata served from Postgres (`pgx/v5`)
- Templates rendered with [`a-h/templ`](https://templ.guide)
- Static assets (CSS, JS, favicon) embedded via `embed.FS` in production;
  served from disk in development
- Optional Redis client initialised at startup (reserved for future
  scraper caching; today it only contributes to the health check)
- Echo middleware stack (request ID, structured request log, secure
  headers, gzip, recover)

Data sources and show lists:

* [arrow.fandom.com](https://arrow.fandom.com)
* [en.wikipedia.org](https://en.wikipedia.org)

## Currently Supported Series

* Arrow
* Batwoman
* Black Lightning
* Constantine
* DC's Legends of Tomorrow
* The Flash
* Freedom Fighters: The Ray
* Stargirl
* Supergirl
* Superman & Lois
* Vixen

## Development Setup

Prerequisites (managed by [mise](https://mise.jdx.dev)):

- Go 1.25+
- `task` (Taskfile runner)
- Docker (for `compose.yaml`)
- `pre-commit` (optional for the hook chain)

First-time setup:

1. `cp .local.env.example .local.env` (and `.env.example .env` if needed
   for compose)
2. `mise install` to fetch the toolchain
3. `task compose:dev` to bring up Postgres + Redis
4. `task db:migrations:dev` to apply the schema
5. `go run ./cmd/harvest` to (re)generate the seed SQL into
   `migrations/002_seed_arrowverse.sql` against the live API
6. `go run ./cmd/ordered-arrowverse` to start the service at
   `http://localhost:8080`

(All of the above is automated via the `dev` target.)

### Tasks

| Task | What it does |
| ---- | ------------ |
| `task dev` | Run migrations, then start templ + air in watch mode |
| `task templ:generate` | Regenerate templ components |
| `task sqlc:generate:dev` | Regenerate sqlc queries against the dev DB |
| `task lint:go` | gofumpt + betteralign + golangci-lint + templ fmt |
| `task lint:db` | sqlfluff, squawk, sqlc compile/vet |
| `task db:migrations:dev` | Apply pending migrations |
| `task db:migrations:new -- NAME` | Author a new migration |

### Environment variables

See `.local.env.example` for the full set. The application reads them
directly (no 1Password dependency). Notable entries:

| Variable | Required | Purpose |
| -------- | -------- | ------- |
| `DATABASE_URL` | yes | Postgres connection string |
| `REDIS_URL` | no (defaults) | Redis connection string |
| `SERVICE_NAME` | yes | Service identifier in logs |
| `ENVIRONMENT` | yes | `development`, `production`, etc. |
| `SITE_TITLE` / `SITE_HEADING` | yes | Page <title> and navbar brand |
| `OLD_SITE_HOST` / `NEW_SITE_URL` | no | Legacy Heroku banner switch |
| `PORT` | no (defaults 8080) | HTTP listen port |
| `LOG_LEVEL` | no (defaults INFO) | zerolog threshold |
| `DATABASE_LOG_LEVEL` | no (defaults WARN) | pgx log threshold |
| `ADMIN_EMAIL` | yes | Outbound user-agent contact |
| `SHUTDOWN_TIMEOUT` | no (defaults 5s) | Graceful-shutdown grace period |
| `DEMO_MODE_HOST_PREFIX` | no (defaults `demo.`) | Reserved by future middleware |

## Production Build

`Dockerfile` is multi-stage and writes a stripped static binary onto
`gcr.io/distroless/static-debian12:nonroot`. Build + run locally:

```sh
docker compose -f compose.production.yaml up --build
```

The compose stack exposes the service on port 8080 with persistent
volumes for Postgres and Redis. The distroless image has no shell, so
container-level healthchecks are intentionally omitted (orchestrators
should HTTP-probe `/health/ping`).

## URLs

| Path | Purpose |
| ---- | ------- |
| `/` | Filtered home view |
| `/api` | JSON export of the filtered view |
| `/recent_episodes.atom` | Atom feed of the 15 most recent episodes |
| `/newest_first` | Legacy 301 → `/?newest_first=true` |
| `/hide/<list>` | Legacy 301 → canonical home with `hide_show` applied |
| `/legal/privacy-policy` | Privacy notice |
| `/legal/cookie-policy` | Cookie notice |
| `/ads.txt` | AdSense verification file |
| `/health/ping` | DB + Redis healthcheck JSON |

## CI / CD

- `.github/workflows/lint.yml` — `gofumpt`, `go vet`, `go build` on PRs
- `.github/workflows/test.yml` — Postgres + Redis services, tern
  migrations, `go test -race -count=1 ./...`
- `.github/workflows/build.yml` — Multi-arch Docker build on tag push
- `.github/workflows/lint_workflows.yml` — `actionlint` over the workflow
  set itself
- `.github/dependabot.yml` — gomod + github-actions ecosystems

Official actions are pinned to immutable commit SHAs (with tag-name
comments) per the project's supply-chain posture.

## License

MIT, like the upstream project.
