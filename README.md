# Completionist

A social completion-tracking platform for games, movies, anime, books, and music. Track what you are playing, watching, and reading, mark your progress, and share it with other people who care about finishing things.

## What it is

Completionist is a full-stack app for people who keep backlogs. You add titles from real media catalogs (TMDB, Jikan, RAWG, Google Books, Last.fm, Steam), set their status and progress, and the platform builds your profile, ranks, streaks, and stats around them. On top of personal tracking it adds a social layer: posts, profiles, follows, realtime rooms and direct messages, seasonal challenges, and a moderated public space. It is aimed at hobbyists who want one place for their whole media backlog, and it doubles as a portfolio project demonstrating a layered Go API and a Next.js frontend.

## Key features

**Tracking and progress.** Add items across games, movies, TV, anime, manga, books, and music. Per-item status and progress with a validator that keeps progress values consistent. Separate wishlist (priority and optional target price) and a "next up" queue you can reorder.

**Social and posts.** User profiles, follows, and a post feed with reactions. Public, guest-readable pages for search and franchises.

**Realtime messaging.** Direct messages and shared rooms over a WebSocket hub. Rooms come in three types: chat, music, and watch, with synced playback state for music and watch rooms. Typing indicators, reactions, and pinned messages are delivered as realtime events.

**Challenges and seasons.** Time-boxed challenges and a season system that resets competitive standings on a schedule.

**Badges, streaks, ranks, and XP.** Earn badges, keep daily streaks alive, gain XP and levels, and climb an ELO-based rank ladder with six tiers (Scribe, Chronicler, Curator, Preserver, Warden, Oracle).

**Collections and reviews.** Group list items into named collections, and write reviews on the things you finish.

**Franchises.** Browse a franchise as a visual tree of its entries. Franchise pages are public.

**Subscriptions via Stripe.** Paid membership through Stripe Checkout and the Stripe billing portal. Members unlock subscription-gated features; non-members see ad slots.

**Moderation and admin.** Role-based admin panel for managing users, a moderation queue, bans, ad campaigns, and franchise entries. Optional LLM-backed content moderation, with a SQL-based path when no provider is configured. Audit logging records sensitive actions.

**External media integrations.** TMDB (movies and TV), Jikan (anime and manga), RAWG (games), Google Books (books), Last.fm (music scrobbles), and Steam (library import via OpenID).

## Architecture at a glance

Monorepo with two apps:

- **apps/api** is a Go 1.25 HTTP API on go-chi. Requests flow through a layered path: handler to service to repository, with raw SQL (sqlx) over PostgreSQL. Auth is JWT Bearer tokens with RBAC roles enforced in middleware. A gorilla/websocket hub handles realtime messaging, file uploads go to S3 (or local disk), and sensitive actions are written to an audit log. The REST surface is documented with swaggo annotations and served as Swagger UI.
- **apps/web** is a Next.js 16 App Router frontend in TypeScript and Tailwind. Every page is a client component; all calls go through one shared axios instance that attaches the JWT, and realtime updates arrive over a native browser WebSocket.

Schema migrations run automatically at API startup via golang-migrate. For a deeper write-up see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Tech stack

| Layer | Choice |
| --- | --- |
| API language | Go 1.25 |
| HTTP router | go-chi/chi v5 |
| DB access | sqlx over lib/pq (no ORM) |
| Database | PostgreSQL 16 |
| Migrations | golang-migrate |
| Auth | golang-jwt v5 (HS256), RBAC |
| Realtime | gorilla/websocket |
| File storage | AWS SDK v2 (S3) or local disk |
| Payments | stripe-go v81 |
| API docs | swaggo/swag + http-swagger |
| Frontend | Next.js 16, React 19, TypeScript, Tailwind v4 |
| HTTP client (web) | axios |
| Package manager | Bun |

## Monorepo layout

```
completionist-api-go/
├── apps/
│   ├── api/                  # Go backend
│   │   ├── cmd/api/main.go   # composition root
│   │   ├── internal/
│   │   │   ├── auth/         # JWT issue/validate
│   │   │   ├── bootstrap/    # dev admin seeding
│   │   │   ├── config/       # env + .env loading
│   │   │   ├── database/     # connect + run migrations
│   │   │   ├── handler/      # chi handlers
│   │   │   ├── httpx/        # JSON response helpers
│   │   │   ├── middleware/   # auth + RBAC
│   │   │   ├── models/       # domain + payload structs
│   │   │   ├── repository/   # raw SQL per domain
│   │   │   ├── routes/       # route wiring
│   │   │   ├── services/     # cross-domain + external API clients
│   │   │   └── websocket/    # realtime hub
│   │   ├── migrations/       # golang-migrate SQL pairs
│   │   └── swagger/          # generated OpenAPI docs
│   └── web/                  # Next.js frontend
│       └── src/{app,components,contexts,hooks,lib,services,types}
├── scripts/                  # dev.sh, backup-db-to-s3.sh
├── docs/                     # architecture, API, deploy guides
├── docker-compose.yml        # local Postgres
├── docker-compose.prod.yml   # production images
├── docker-compose.aws.yml    # self-hosted Postgres overlay
├── Dockerfile.api
├── Dockerfile.web
└── package.json              # root scripts (Bun)
```

## Quickstart

Prerequisites:

- Go 1.25+
- Bun 1.3.9+
- Docker with the Compose plugin

Install dependencies and the toolchain (installs JS deps, starts local Postgres, installs the golang-migrate and swag CLIs):

```bash
bun run setup
```

Run both apps in development (regenerates Swagger docs, then starts the API and the frontend together):

```bash
bun run dev
```

Or use the launcher script, which waits for Postgres to be ready before starting:

```bash
scripts/dev.sh
```

Ports:

- API: `8080`
- Web: `4000`
- Local Postgres: `5435` (host) mapped from `5432`

## Available scripts

Run from the repo root with `bun run <script>`.

| Script | What it does |
| --- | --- |
| `setup` | Install JS deps, start local Postgres, install migrate + swag CLIs |
| `dev` | Run API and frontend together (regenerates Swagger first) |
| `dev:backend` | Generate Swagger, then run the Go API in development |
| `dev:frontend` | Run Next.js on port 4000 |
| `build` | Build the Go binary and the Next.js app |
| `build:backend` | Compile the Go binary to `apps/api/bin/server` |
| `build:frontend` | Build the Next.js app |
| `start` | Run the compiled API and frontend together |
| `db:up` / `db:down` | Start / stop local Postgres |
| `db:reset` | Wipe the Postgres volume and start fresh |
| `db:logs` | Follow Postgres logs |
| `migrate:create` | Scaffold a new up/down SQL migration pair |
| `swagger:gen` | Regenerate Swagger docs from the `// @` annotations |

## API docs

The REST API is documented inline with swaggo annotations. Regenerate the OpenAPI spec with:

```bash
bun run swagger:gen
```

This is run automatically before each `bun run dev`. The generated Swagger UI is served by the API at `/swagger/index.html` on port 8080. For an endpoint-level reference see [docs/API.md](docs/API.md).

## Status

Active development. The core tracking, social, realtime, and subscription features are working end to end; expect ongoing changes as domains are refined. Deployment guides for Supabase + VPS and AWS Lightsail live in [docs/](docs/).
