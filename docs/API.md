# Completionist REST API

The backend is a Go HTTP API (go-chi router, PostgreSQL via sqlx). This document is an orientation guide: it covers the base path, auth, response conventions, and the endpoint groups with a few representative routes each. It is not the full reference. For every route, request body, and response schema, use the generated Swagger UI (see [Full reference](#full-reference)).

## Base path and port

- The API listens on port `8080` (`SERVER_ADDR=:8080`).
- Every route is mounted under `/api`. So the login endpoint is `POST /api/login`, list items live at `/api/lists`, and so on.
- Locally that means `http://localhost:8080/api/...`. In production the host is fronted by Caddy and reached at `https://api.<domain>/api/...`.
- Uploaded files (when the local storage driver is active) are served from `/uploads/*`, outside the `/api` prefix.
- Swagger UI is served at `/swagger/index.html`, also outside `/api`.

Paths in the route tables below are written relative to `/api`. A row listing `POST /login` means `POST /api/login`.

## Authentication

Auth is a JWT Bearer token, signed HS256. There are no sessions or cookies on the API side; the token is the only credential.

### Obtaining a token

1. `POST /api/register` with `{ "username", "email", "password" }` creates the account and returns the user object.
2. `POST /api/login` with `{ "email", "password" }` returns `{ "token": "<jwt>" }`.

The token carries the user id and expires 24 hours after issue. There is no refresh endpoint; clients log in again when the token expires.

### Passing the token

Send it as a Bearer token on the `Authorization` header:

```
Authorization: Bearer <jwt>
```

The WebSocket endpoint cannot set custom headers from the browser, so `GET /api/ws` accepts the token as a query parameter instead: `/api/ws?auth_token=<jwt>`.

On every authenticated request the middleware validates the token, then does a live database lookup of the user to read the current role and ban state. Two consequences:

- A role change or ban takes effect on the next request, not on the next login.
- A banned user is blocked from all authenticated routes except `POST /api/me/appeal` and `GET /api/me/export`, which stay open so the user can appeal and export their data. A blocked request returns `403` with `{ "error": "Account suspended", "banReason": "<reason>" }`.

### Roles

Three roles exist: `user`, `member`, and `admin`. `member` is the paid subscription tier. Unauthenticated callers have no role and are treated as a guest with access to public routes only. Routes under `/api/admin` require the `admin` role and return `403` otherwise.

## Conventions

- Requests and responses are JSON. Send `Content-Type: application/json` on any request with a body.
- Success responses use the conventional status codes: `200` for reads and updates, `201` for creates, `204 No Content` for deletes.
- Errors use a single flat shape:

  ```json
  { "error": "Invalid request payload" }
  ```

  The HTTP status carries the category (`400` bad input, `401` missing or invalid token, `403` forbidden or banned, `404` not found, `500` server error). The ban response is the one case that adds a second field (`banReason`) alongside `error`.
- Pagination, where present, is `limit` and `offset` query parameters with server-side defaults (for example `/api/me/xp` defaults to `limit=20`, `offset=0`). Most collection endpoints return the full set for the current user and do not paginate.
- Path parameters are written in braces below, for example `{item_id}` or `{username}`.

## Endpoint groups

### Auth and current user

| Method | Path | Purpose |
| --- | --- | --- |
| POST | `/register` | Create a new account |
| POST | `/login` | Exchange email and password for a JWT |
| GET | `/me` | Current user's profile |
| PUT | `/me/profile` | Update the current user's profile |

Account lifecycle also lives here: `POST /me/delete` requests deletion, `DELETE /me/delete` cancels it, `GET /me/export` downloads the user's data, and `POST /me/appeal` submits a ban appeal.

### Users and profiles

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/users/search` | Search users by name |
| GET | `/users/{username}/profile` | Public profile for a user |
| GET | `/users/{id}/stats` | Aggregate completion stats for a user |
| GET | `/users/{id}/stats/heatmap` | Activity heatmap data for a user |

### Lists and tracking

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/lists` | All list items for the current user |
| POST | `/lists` | Add an item to the list (creates the media row if needed) |
| PATCH | `/lists/{item_id}` | Update status, progress, or rating |
| DELETE | `/lists/{item_id}` | Remove a list item |

Related sub-areas use the same pattern: the wishlist at `/wishlist` (items carry a priority and optional price) and the up-next queue at `/queue` (`PUT /queue/reorder` reorders it). There are also typed shortcuts that add an item straight from an external source, for example `POST /lists/tmdb/movie/{tmdb_id}` and `POST /lists/jikan/anime/{mal_id}`.

### Media and integrations

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/search/movies` | Search movies (TMDB) |
| GET | `/search/anime` | Search anime (Jikan) |
| GET | `/games/rawg/search` | Search games (RAWG) |
| GET | `/media/trending` | Trending media across sources |

Connected external accounts hang off this area too. Steam: `GET /auth/steam/login`, `POST /me/steam/attach`, `GET /me/steam/owned`, `GET /steam/achievements/{app_id}`. Last.fm: `GET /auth/lastfm`, `GET /me/lastfm/recent`. Generic OAuth connections are managed at `/connections` (`GET /connections`, `POST /connections/{provider}/connect`, `POST /connections/{provider}/sync`), and library imports are tracked at `GET /imports`.

### Posts and social

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/posts` | The post feed |
| POST | `/posts` | Create a post |
| POST | `/posts/{post_id}/like` | Like a post |
| POST | `/posts/{post_id}/comments` | Comment on a post |

Following lives alongside posts: `POST /users/{username}/follow`, `DELETE /users/{username}/unfollow`, `GET /users/{username}/followers`, and `GET /users/suggestions` for who to follow next. Comments have their own like and edit routes under `/comments/{comment_id}`.

### Messaging and rooms

| Method | Path | Purpose |
| --- | --- | --- |
| POST | `/messages/send` | Send a direct message |
| GET | `/messages/conversations` | List the user's conversations |
| POST | `/rooms` | Create a room (chat, music, or watch) |
| PATCH | `/rooms/{roomId}/state` | Update shared playback state for a room |

Realtime delivery for both DMs and rooms is the WebSocket at `GET /api/ws?auth_token=<jwt>`. The REST routes above persist and seed state; live events (new messages, joins, reactions, playback changes) arrive over the socket.

### Challenges, seasons, and leaderboards

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/challenges` | Active challenges |
| GET | `/me/challenges` | Current user's challenge progress |
| GET | `/seasons/current` | The active season |
| GET | `/seasons/leaderboard` | Season leaderboard |

### Badges, streaks, and ranks

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/badges` | All available badges |
| GET | `/users/{id}/badges` | Badges a user has earned |
| GET | `/streaks/leaderboard` | Streak leaderboard |
| GET | `/me/rank` | Current user's ELO rank for the active season |

XP history for the current user is at `GET /me/xp`.

### Collections

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/collections` | The user's collections |
| POST | `/collections` | Create a collection |
| PUT | `/lists/{item_id}/collection` | Assign a list item to a collection |
| DELETE | `/lists/{item_id}/collection` | Remove a list item from its collection |

### Reviews

| Method | Path | Purpose |
| --- | --- | --- |
| POST | `/reviews` | Write a review |
| PUT | `/reviews/{id}` | Edit a review |
| GET | `/reviews/media/{id}` | Reviews for a media item |
| GET | `/reviews/user/{id}` | Reviews written by a user |

### Franchises

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/franchises` | List franchises |
| GET | `/franchises/{id}` | A franchise and its entry tree |

Franchise listing and detail are public (reachable by guests). Creating and editing franchise entries is admin-only under `/admin/franchises`.

### Subscriptions

| Method | Path | Purpose |
| --- | --- | --- |
| POST | `/subscriptions/checkout` | Start a Stripe checkout session |
| GET | `/subscriptions/me` | The current user's subscription |
| POST | `/subscriptions/portal` | Open the Stripe billing portal |

Stripe calls back to `POST /api/webhooks/stripe` (public, verified by signature) to keep subscription state in sync.

### Moderation and admin

| Method | Path | Purpose |
| --- | --- | --- |
| POST | `/moderation/report` | Report a piece of content (any signed-in user) |
| GET | `/admin/moderation` | The moderation queue (admin) |
| POST | `/admin/users/{id}/ban` | Ban a user (admin) |
| GET | `/admin/ban-appeals` | Pending ban appeals (admin) |

Everything under `/admin` requires the `admin` role. It also covers user listing (`GET /admin/users`), ad campaigns (`/admin/ads`), and franchise management (`/admin/franchises`).

### Ads

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/ads` | Active ad for a placement (`?placement=`) |
| POST | `/ads/{id}/impression` | Record a view or click |

Both are public so ads can render for logged-out and non-member viewers.

## Full reference

The complete list of routes, with request bodies and response schemas, is the Swagger UI. The docs are generated from `// @` annotations in the handler source by swaggo, so they track the code.

- Open it at `http://localhost:8080/swagger/index.html` while the API is running.
- The OpenAPI title is "Completionist API" and the documents (`swagger.json`, `swagger.yaml`, `docs.go`) live in `apps/api/swagger/`.

Regenerate after changing any handler or its annotations:

```
bun run swagger:gen
```

That runs `swag init -g cmd/api/main.go` and rewrites the files in `apps/api/swagger/`. The dev launcher (`bun run dev`) regenerates them on startup, so a normal dev run already picks up annotation changes; run the command manually when you only need fresh docs without booting the server.
