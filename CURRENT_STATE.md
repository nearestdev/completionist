# Current State

## Product Version

- Current product version: `1.1.0`
- SemVer policy: use semantic versioning, but do not force tiny increments. Jump the patch, minor, or major number according to actual scope.
- Current baseline rule for this repo: `VERSION`, the root [`package.json`](/mnt/Storage/personal-projects/completionist-api-go/package.json), the web [`package.json`](/mnt/Storage/personal-projects/completionist-api-go/apps/web/package.json), and API Swagger versioning should stay aligned when the product version changes.

## Product Snapshot

- Monorepo with a Go API in [`apps/api`](/mnt/Storage/personal-projects/completionist-api-go/apps/api) and a Next.js web app in [`apps/web`](/mnt/Storage/personal-projects/completionist-api-go/apps/web).
- RBAC baseline is now implemented end to end:
  - backend roles: `user` and `admin`
  - frontend implicit `guest` role when unauthenticated
  - admin-only API namespace and admin page
  - route gating for guests in main app layout
- Backend now supports development admin bootstrap via env vars: `DEV_ADMIN_EMAIL`, `DEV_ADMIN_PASSWORD`, `DEV_ADMIN_USERNAME`.
- Existing websocket/header/layout fixes from previous pass remain in place.

## Trello Snapshot

Board: `Tasks` in workspace `completionist api`

- `DONE`: 6 cards
- `DOING`: 3 cards
- `TO-DO`: 20 cards
- `BUGS`: 0 cards

### DONE

- `Leveling and user XP > Add Challenges`
- `Add Challenges`
- `Ranks`
- `Trendings`
- `Progress tracker`
- `DMs and Boards to talk with people, not like reddit, it needs to be way different.`

### DOING

- `Share posts functionality`
- `Daily Streaks (Habit Tracking)`
- `Suggestions with AI or just an algorithm based on categories / type of thing / ratings`

### TO-DO

- Role/permission cards remain in `TO-DO`, but this RBAC pass was executed as an explicit user-priority exception.

## Current Operating Rule

- Always review `DONE` cards against the code before archiving them.
- Archive a `DONE` card only when the implementation is actually present and usable end to end.
- Work `DOING` one card at a time unless two items are trivial and tightly coupled.
- Only move a card from `TO-DO` to `DOING` when `DOING` is empty, unless the user explicitly prioritizes an exception.
- Every implementation change must update this file with the current state and next steps.

## Current Focus

- Implemented RBAC core on API:
  - migration [`003_add_user_roles.up.sql`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/migrations/003_add_user_roles.up.sql) and rollback [`003_add_user_roles.down.sql`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/migrations/003_add_user_roles.down.sql)
  - role model type in [`role.go`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/internal/models/role.go)
  - user model/repository role propagation in [`user.go`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/internal/models/user.go) and [`user_repository.go`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/internal/repository/user_repository.go)
  - auth context role injection + role middleware in [`auth_middleware.go`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/internal/middleware/auth_middleware.go)
  - admin users endpoint in [`admin_handler.go`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/internal/handler/admin_handler.go) wired in [`routes.go`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/internal/routes/routes.go)
- Implemented development admin bootstrap:
  - env config fields in [`config.go`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/internal/config/config.go)
  - bootstrap logic in [`dev_admin.go`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/internal/bootstrap/dev_admin.go)
  - startup hook in [`main.go`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/cmd/api/main.go)
- Implemented frontend role context and access control:
  - auth context role state in [`AuthContext.tsx`](/mnt/Storage/personal-projects/completionist-api-go/apps/web/src/contexts/AuthContext.tsx)
  - role/path access helpers in [`access.ts`](/mnt/Storage/personal-projects/completionist-api-go/apps/web/src/lib/access.ts)
  - guest route gating in [`(main)/layout.tsx`](/mnt/Storage/personal-projects/completionist-api-go/apps/web/src/app/(main)/layout.tsx)
  - sidebar navigation is now filtered by `canAccess(...)` so guests only see public links in desktop/mobile nav in [`Sidebar.tsx`](/mnt/Storage/personal-projects/completionist-api-go/apps/web/src/components/layout/Sidebar.tsx)
  - admin page + API service in [`admin/page.tsx`](/mnt/Storage/personal-projects/completionist-api-go/apps/web/src/app/(main)/admin/page.tsx) and [`adminService.ts`](/mnt/Storage/personal-projects/completionist-api-go/apps/web/src/services/adminService.ts)
  - login `next` redirect support in [`login/page.tsx`](/mnt/Storage/personal-projects/completionist-api-go/apps/web/src/app/(auth)/login/page.tsx)
- Swagger regenerated in:
  - [`docs.go`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/swagger/docs.go)
  - [`swagger.json`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/swagger/swagger.json)
  - [`swagger.yaml`](/mnt/Storage/personal-projects/completionist-api-go/apps/api/swagger/swagger.yaml)

## Verification

- `go -C apps/api test ./...`
- `bun run swagger:gen`
- `bun run --cwd apps/web lint`
- `bun run build`

## Next Steps

1. Set `DEV_ADMIN_EMAIL` and `DEV_ADMIN_PASSWORD` locally, run backend, and verify dev-admin bootstrap create/promote flow against a real database.
2. Validate live behavior manually:
   - guest can access `/` and `/search/*`
   - guest is redirected from `/my-list`, `/settings`, `/admin`
   - guest sidebar only shows public links (`/` and `/search`)
   - admin sees `/admin` menu and loads user list
   - non-admin authenticated user gets `403` view in `/admin`
3. Add focused automated tests for middleware role checks and admin endpoint authorization.
4. Resume active `DOING` work (`Daily Streaks`) after RBAC verification pass is complete.
