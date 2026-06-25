# Development guide

This is the setup and day-to-day workflow guide for the Completionist monorepo. It covers what to install, how to get a local environment running, the environment variables, common database and docs tasks, and how to build for production.

The repo is a single bun workspace at the root with two apps:

- `apps/api` is the Go 1.25 HTTP API (chi router, sqlx over PostgreSQL, JWT auth, WebSocket realtime, S3 or local file storage).
- `apps/web` is the Next.js (App Router) + TypeScript + Tailwind frontend.

For deployment specifics see `docs/DEPLOY_AWS_LIGHTSAIL.md`. For a deeper explanation of how the host and port values fit together, see `docs/PORT_CONFIGURATION.md`.

## Prerequisites

| Tool | Version | Notes |
| --- | --- | --- |
| Go | 1.25+ | Builds and runs the API. Must be on `PATH`. |
| Bun | 1.3.9+ | Package manager and task runner for the whole repo. The root `package.json` pins `bun@1.3.9`. |
| Docker | recent, with the Compose plugin | Runs local Postgres. The dev script accepts either `docker compose` or the legacy `docker-compose`. |
| Node | 22+ | Not needed locally; Bun handles dev and `bun run build`. AWS Amplify builds the frontend on Node in the cloud. |

Two CLI tools are installed into your Go bin during setup, so you do not need to fetch them by hand:

- `migrate` (golang-migrate) for creating migration files. It is installed with the `postgres` build tag.
- `swag` (swaggo) for generating Swagger docs from the `// @` annotations in the handlers.

Make sure your Go bin directory (usually `~/go/bin`) is on your `PATH` so `migrate` resolves after setup.

## First-time setup

From the repo root:

```bash
bun run setup
```

That single command runs, in order:

1. `bun install` to install the JS workspace dependencies.
2. `bun run db:up` (`docker compose up -d`) to start the local Postgres container.
3. `bun run migrate:install-tool` to `go install` the golang-migrate CLI with the `postgres` tag.
4. `bun run swagger:install` to `go install` the swaggo `swag` CLI.

You do not run migrations as a separate step. The Go API applies all pending migrations from `apps/api/migrations/` automatically at startup (in both development and production), so the first `bun run dev` brings the schema up to date.

Before you start, create your env file. There is no tracked `.env`, only an example. Copy it and fill in the values you need:

```bash
cp .env.production.example .env
```

For local development most integration keys can stay blank. The only value worth setting early is `JWT_SECRET` (any non-empty string works locally). The database defaults already match the local Postgres container, so you do not need to set `DATABASE_URL` for local work unless you point at a different database.

## Running in development

The recommended launcher is the dev script, which owns the database lifecycle and waits for Postgres to be ready before starting the apps:

```bash
scripts/dev.sh
```

It starts the Postgres container (if it is not already up), polls `pg_isready` until the database accepts connections, prints the local URLs, then hands off to `bun run dev`. Ctrl+C stops the backend and frontend but leaves Postgres running so the next start is fast and your data survives. Use `scripts/dev.sh down` to stop the container (data preserved) and `scripts/dev.sh reset` to wipe the volume and start fresh.

If your database is already running, you can skip the script and run the apps directly:

```bash
bun run dev
```

`bun run dev` uses `concurrently` to run two named processes side by side:

- `BACKEND`: runs `bun run swagger:gen` first to regenerate the Swagger docs, then starts the Go API with `GO_ENV=development` via `go -C apps/api run ./cmd/api`. On startup it loads `.env`, applies migrations, and (in development) seeds a dev admin if the `DEV_ADMIN_*` variables are set.
- `FRONTEND`: runs the Next.js dev server with `PORT=4000`.

If you want to run just one side, the underlying scripts are `bun run dev:backend` and `bun run dev:frontend`.

Once both are up:

- Frontend: http://localhost:4000
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html

## Environment variables

Copy `.env.production.example` to `.env` for local work and to `.env.production` for a deployment. Never commit a filled-in env file and never paste real secret values into docs or chat.

The table below lists every key. "Required" means the app needs it to start or to serve its core function; everything marked optional gates a specific integration that stays dormant when left blank.

### Core

| Key | Required | Description |
| --- | --- | --- |
| `JWT_SECRET` | yes | HMAC secret used to sign and validate JWTs. Use a long random string. |
| `DATABASE_URL` | yes (prod) | Postgres connection string. For local dev the app falls back to the Docker Postgres defaults, so you usually leave it unset locally. In prod it points at the managed RDS Postgres endpoint (`sslmode=require`). |
| `GO_ENV` | no | `development` or `production`. Set automatically by the dev and start scripts; controls dev-only behavior such as admin seeding. |

### Host and port configuration

| Key | Required | Description |
| --- | --- | --- |
| `SERVER_ADDR` | no | Address the API listens on. Defaults to `:8080`. |
| `BACKEND_PORT` | no | API port advertised to other config and to the web build. Defaults to `8080`. |
| `BACKEND_BASE_URL` | no | Public base URL of the API (scheme and host, no port). Used to build absolute links, and set as an Amplify build env so the web app knows the API URL. |
| `PORT` | no | Frontend port. The scripts set this to `4000`. |
| `FRONTEND_BASE_URL` | no | Public base URL of the frontend (the Amplify domain in prod). |
| `PUBLIC_BASE_URL` | no | Public base URL used when building absolute media/upload links. |
| `ALLOWED_ORIGINS` | no | Comma-separated CORS allowlist for the API (the Amplify origin in prod). |

### File storage

| Key | Required | Description |
| --- | --- | --- |
| `STORAGE_DRIVER` | no | `local` (default) writes uploads to disk and serves them from `/uploads`. `s3` uses S3. |
| `AWS_REGION` | conditional | Required when `STORAGE_DRIVER=s3`. |
| `AWS_BUCKET` | conditional | S3 bucket for uploads. Required when `STORAGE_DRIVER=s3`. |
| `AWS_ACCESS_KEY_ID` | conditional | Required when `STORAGE_DRIVER=s3`. |
| `AWS_SECRET_ACCESS_KEY` | conditional | Required when `STORAGE_DRIVER=s3`. |

### Dev admin bootstrap (development only)

| Key | Required | Description |
| --- | --- | --- |
| `DEV_ADMIN_EMAIL` | no | If set in development, seeds an admin account with this email on startup. |
| `DEV_ADMIN_PASSWORD` | no | Password for the seeded dev admin. |
| `DEV_ADMIN_USERNAME` | no | Username for the seeded dev admin (example defaults to `admin`). |

### External media APIs (optional)

| Key | Required | Description |
| --- | --- | --- |
| `TMDB_API_KEY` | no | The Movie Database. Enables movie and TV search. |
| `STEAM_WEB_API_KEY` | no | Steam Web API. Enables Steam library import. |
| `STEAM_CALLBACK_REDIRECT_PATH` | no | Path the Steam OpenID flow returns to after login (example: `/settings`). |
| `RAWG_API_KEY` | no | RAWG. Enables game search and the user games section. |
| `GOOGLE_BOOKS_API_KEY` | no | Google Books. Enables book search. |
| `LASTFM_API_KEY` | no | Last.fm. Enables music search and scrobble history. |
| `LASTFM_API_SECRET` | no | Last.fm shared secret, paired with the API key. |

### Stripe (membership payments, optional)

| Key | Required | Description |
| --- | --- | --- |
| `STRIPE_SECRET_KEY` | no | Stripe secret key. Enables checkout and the billing portal. |
| `STRIPE_WEBHOOK_SECRET` | no | Verifies incoming Stripe webhook signatures. |
| `STRIPE_PRICE_ID` | no | The price the membership checkout charges. |

### Content moderation and LLM suggestions (optional)

| Key | Required | Description |
| --- | --- | --- |
| `MODERATION_ENABLED` | no | `true` turns on AI content moderation. Defaults to `false`. |
| `OPENAI_MODERATION_KEY` | no | Key for the OpenAI moderation endpoint, used when moderation is enabled. |
| `LLM_PROVIDER` | no | `openai` or `anthropic`. Empty disables LLM-powered suggestions; the SQL-based suggestion path still works. |
| `LLM_API_KEY` | no | API key for the chosen `LLM_PROVIDER`. |

## Database tasks

The local Postgres container is defined in `docker-compose.yml`:

- Host port `5435` (mapped from the container's `5432`)
- Database `completionist_db`, user `postgres`, password `testPassword`
- Data persists in the `postgres_data` Docker volume

Container lifecycle scripts (root `package.json`):

```bash
bun run db:up      # docker compose up -d
bun run db:down    # docker compose down (keeps the volume)
bun run db:logs    # follow Postgres logs
bun run db:reset   # docker compose down -v && up -d  (wipes all local data)
```

`bun run db:reset` destroys the volume, so the database comes back empty. On the next API start the migrations re-apply and you have a clean schema with no rows. Use it when you want a fresh database or when a half-applied migration left things in a bad state.

### Creating a migration

Migrations live in `apps/api/migrations/` and use sequential numbering with paired `up`/`down` SQL files (for example `012_add_widgets.up.sql` and `012_add_widgets.down.sql`). To scaffold a new pair:

```bash
bun run migrate:create -- add_widgets
```

The `--` passes the name through to the `migrate create` command, which assigns the next sequence number and creates both files. Write your forward SQL in the `.up.sql` file and the matching rollback in the `.down.sql` file.

You do not apply migrations manually. They run automatically when the API starts (the startup path tolerates "no change"), so creating the files and restarting the backend is enough to apply them locally.

## Regenerating Swagger docs

The API serves Swagger UI at `/swagger/index.html`, generated from the `// @` annotations in the handler files. The generated output lands in `apps/api/swagger/`.

```bash
bun run swagger:gen
```

`bun run dev` already runs this before starting the backend, so during normal development the docs stay current. Run it by hand after you change route annotations if you want to refresh the docs without restarting the full dev stack.

## Building for production

```bash
bun run build
```

This runs two steps in sequence:

- `build:backend`: compiles the Go API to `apps/api/bin/server`.
- `build:frontend`: runs the Next.js production build for `apps/web`.

To run the compiled artifacts locally:

```bash
bun run start
```

`start` runs the backend binary with `GO_ENV=production` and the frontend on `PORT=4000`, both via `concurrently`.

### Docker

`Dockerfile.api` is the only production image: a multi-stage Go build on `golang:1.25-alpine`, final image on `alpine:3.22` running as a non-root `app` user. It copies the compiled binary and the `migrations/` directory, exposes `8080`, and applies migrations at container startup like any other run of the binary. The frontend is built and served by AWS Amplify, so there is no production web image.

The Compose files:

- `docker-compose.yml`: local Postgres only, for development.
- `docker-compose.prod.yml`: the `api` service reading `.env.production`. It connects to managed RDS Postgres via `DATABASE_URL` and stores uploads in S3.

Typical production command (see the deploy doc for the full procedure):

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml up -d --build
```

## Ports

| Service | Port | Where |
| --- | --- | --- |
| API | 8080 | `SERVER_ADDR=:8080`. Local and inside the Docker network. Swagger UI at `/swagger/index.html`. |
| Frontend | 4000 | Next.js dev/start server, set by `PORT=4000` in the dev and start scripts. In production the frontend is hosted on AWS Amplify. |
| Local Postgres | 5435 | Host port mapped from the container's `5432` (`docker-compose.yml`). |

In production the Caddy overlay terminates TLS on `80`/`443` and reverse-proxies `api:8080`. The frontend is hosted on AWS Amplify with its own CDN and TLS. See the deploy doc for that setup.

## Running tests

The API has a Go test suite. Run it from the repo root:

```bash
go -C apps/api test ./...
```

The `-C apps/api` flag runs the command inside the API module so you do not need to change directories. The tests are standard-library only (the `testing` package plus `net/http/httptest`); they exercise pure logic and HTTP handlers in memory, so they need no database, Docker, or network access and run without any setup.

## Troubleshooting

**Port already in use.** If the API or frontend fails to bind, something else is holding `8080` or `4000`. Find and stop it (`lsof -i :8080` or `lsof -i :4000` on macOS/Linux), or change the port via `SERVER_ADDR`/`BACKEND_PORT` (API) or `PORT` (frontend) in your `.env`. If you change the API port, restart the backend; if you change ports for a Docker deploy, the web image must be rebuilt because the API URL is compile-time.

**Database is not up / connection refused.** The API cannot connect to Postgres on start. Check the container is running with `docker compose ps` (or `bun run db:logs` to tail it). If you started the apps with plain `bun run dev`, the database is not started for you; prefer `scripts/dev.sh`, which starts Postgres and waits for it to accept connections before launching. If Postgres is up but the API still cannot reach it, confirm your `DATABASE_URL` (or the default host port `5435`) matches the running container.

**Postgres did not become ready in time.** `scripts/dev.sh` polls for up to 60 seconds. A slow first start (image pull, volume init) can exceed that. Run `bun run db:logs` to watch the container come up, then rerun the script.

**`migrate` or `swag` not found.** These are installed into your Go bin by `bun run setup`. If the commands are not found, rerun `bun run migrate:install-tool` and `bun run swagger:install`, and make sure `~/go/bin` is on your `PATH`.

**A migration failed halfway.** The local fastest recovery is `bun run db:reset` to wipe the volume and let the migrations re-apply from scratch on the next API start. Only do this locally; it deletes all local data.

**Stale Swagger UI.** If `/swagger/index.html` does not reflect a route change, run `bun run swagger:gen` and reload. `bun run dev` regenerates on every backend start.
