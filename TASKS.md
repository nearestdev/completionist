# TASKS

> Working state of the project: open work, cleanup, and the ordered roadmap. Trello holds card status; this file mirrors it plus the work that never became a card.

## 1. Version & snapshot

- **Current product version:** `2.0.0`
- **Milestone reached:** `2.0.0`. Every feature has a backend (migrations, models, repos, services, handlers, routes) and a frontend (types, services, components, pages).
- **Climb strategy:** forward-only. Each release bumps patch or minor per §9. No resets, no regressions. Patch bumps for cleanup, tests, polish, or docs. Minor bumps for any card that ships a user-visible feature. Major bumps only for API breakage or the 2.0.0 cut.
- **SemVer policy:** semantic versioning with product-level judgment. The next number does not need to be sequential by one; `1.1.0` can jump to `1.3.0` or `1.5.0` if that better matches the scope of the release.
- **Files to keep in sync on version bumps:** [VERSION](VERSION), root [package.json](package.json), web [apps/web/package.json](apps/web/package.json), and the API Swagger version metadata under [apps/api/swagger/](apps/api/swagger/).

**Monorepo layout:**

- Go API: [apps/api](apps/api), `sqlx` + `lib/pq`, golang-migrate auto-runs on startup via [apps/api/internal/database/database.go](apps/api/internal/database/database.go)
- Next.js web: [apps/web](apps/web)
- Migrations: [apps/api/migrations/](apps/api/migrations/) (001 init, 002 seed XP, 003 roles)

**Deployment (single AWS-native path):**

- AWS: [docs/DEPLOY_AWS_LIGHTSAIL.md](docs/DEPLOY_AWS_LIGHTSAIL.md). Go API in Docker on a Lightsail instance behind Caddy ([docker-compose.prod.yml](docker-compose.prod.yml), api only), Next.js frontend on AWS Amplify, managed Postgres on RDS, media on S3. RDS handles backups (automated + PITR + snapshots), so there is no backup script.

**What exists end-to-end today:**

- RBAC: roles `user` / `admin` (backend), implicit `guest` (frontend). Admin namespace, admin page, route gating in [(main)/layout.tsx](apps/web/src/app/(main)/layout.tsx), sidebar filtering in [Sidebar.tsx](apps/web/src/components/layout/Sidebar.tsx), login `next` redirect.
- Dev admin bootstrap: `DEV_ADMIN_EMAIL`, `DEV_ADMIN_PASSWORD`, `DEV_ADMIN_USERNAME` wired in [config.go](apps/api/internal/config/config.go) + [dev_admin.go](apps/api/internal/bootstrap/dev_admin.go), invoked from [main.go](apps/api/cmd/api/main.go).
- Core social features per Trello DONE list, see §4.

---

## 2. Working rules (Trello)

1. Trello holds the work status.
2. Workspace `completionist api`, board `Tasks`.
3. Review `DONE` cards against the actual code before archiving them. An entry in `DONE` is a claim, not a guarantee.
4. Work `DOING` one card at a time. Do not pull from `TO-DO` while `DOING` still has active cards, unless I explicitly prioritize an exception.
5. Every meaningful change updates this file in the same task.
6. Version bumps touch all files listed in §1.

---

## 3. Trello live snapshot

Fetched live on 2026-04-11. Counts will drift as cards move.

| List  | Count | State |
|---    |---    |---    |
| DONE  | 6     | to verify and archive, see §4 |
| DOING | 3     | §5, must clear before pulling from TO-DO |
| TO-DO | 20    | §7, ordered by theme |
| BUGS  | 0     |  |

---

## 4. Bookkeeping cleanup (do first, it's cheap and unblocks prioritization)

These items are purely Trello hygiene. None of them is new engineering work. Doing them now stops future "what's already done?" confusion.

### 4.1 Verify DONE cards against code, then archive

Per rule 2.3, none of these should be archived until someone has spot-checked the implementation. The candidates:

- [#5 Leveling and user XP > Add Challenges](https://trello.com/c/BjdCjXOS)
- [#4 Add Challenges](https://trello.com/c/Gg1CR2H7)
- [#26 Ranks](https://trello.com/c/KfPD3eoC)
- [#8 Trendings](https://trello.com/c/NSwS3WH2)
- [#13 Progress tracker](https://trello.com/c/IvWch1b5)
- [#7 DMs and Boards to talk with people](https://trello.com/c/OpthOEoC): chat, rooms, DMs, reactions (known to exist, recent commits confirm)

Verification approach: for each card, grep the relevant model / handler / migration and confirm a user-visible entrypoint exists in the web app. Archive only what passes.

### 4.2 Stale TO-DO cards already implemented by the RBAC pass

The RBAC pass shipped role model, middleware, admin endpoint, and frontend gating. These TO-DO cards describe that exact work and were left behind:

- [#3 Add Roles to Users > Permission System](https://trello.com/c/HSWPOCV5): **archive**.
- [#1 Permission system](https://trello.com/c/V96CN5j1): **archive**.
- [#9 API Filtering](https://trello.com/c/Kuompad0): **archive** (role-gated routes exist in [routes.go](apps/api/internal/routes/routes.go)).
- [#2 Security](https://trello.com/c/AvaMbhAF): **rescope or archive**. The title is broad enough that there may be remaining security work worth a separate, narrower card (rate limiting, audit log review, secret rotation). Rewrite or split before archiving.

### 4.3 Duplicate cards

- [#14 Ad Campaigns](https://trello.com/c/m8pglZXz) and [#28 AD CAMPAIGNS](https://trello.com/c/rmfRFtze) are duplicates: **archive one** (keep the one with any description; both currently have empty descriptions, so keep the older `#14` and archive `#28`).

---

## 5. Active work (DOING list)

Per rule 2.4, these must be resolved before any new TO-DO card is pulled. Only one is actually spec'd; the other two should be scoped or moved back.

### 5.1 [#23 Daily Streaks (Habit Tracking)](https://trello.com/c/fJJNtADu): next active card

Well-specified and bounded. From the card:

- Streak counter for _active_ engagement (read 10 pages, watch 1 episode gives +1 streak).
- Visual: fire icon that grows with streak tier (Spark / Fire / Inferno).
- UI placement: top-right of header or user mini-profile.
- Hybrid plan: implement own logic first, leave a swap seam for an external habit-tracking API later.

Implementation shape (to be refined when picked up):

- **DB:** new migration `004_add_user_streaks.up.sql` + `.down.sql`, table `user_streaks (user_id PK, current_streak INT, longest_streak INT, last_activity_date DATE, tier TEXT)`. Drive tier from current_streak (Spark 1 to 9, Fire 10 to 99, Inferno 100+).
- **Trigger points:** any model update that counts as "progress", reuse existing XP trigger sites (completion / progress update / challenge progress). Single function `TouchStreak(userID, day)` that the existing flows call.
- **Handler:** `GET /api/users/:id/streak` (public, read-only) plus include streak in the user profile payload so the header can render without a second request.
- **Frontend:** streak pill component in the header/profile, three-tier icon. Read from the existing user context.
- **Tests:** unit test the streak transition matrix (0 to 1, +1 same day = no-op, +1 next day, skip one day = reset, grace window behavior if we decide on one).

### 5.2 [#12 Share posts functionality](https://trello.com/c/RFuAdAkT): needs scoping

Empty description. Before any code: decide what "share" means here.

- Share to an existing Post (repost / quote)?
- Share to a DM / room?
- Share via public link (OG tags, unfurl)?
- Share externally (Twitter / Mastodon intent URL)?

**Action:** either write a description on the card and keep it in DOING, or move back to TO-DO until scoped. Do not start coding until one of the options above is picked.

### 5.3 [#6 Suggestions with AI / algorithm](https://trello.com/c/a82StdGQ): needs scoping

Empty description. Also too broad for a single card.

**Action:** split into two tracks and move back to TO-DO:

- "Non-AI similarity suggestions": purely SQL / category / rating based, cheap, deterministic. Ship first.
- "LLM-based suggestions": pluggable provider, opt-in, rate-limited. Ship second.

---

## 6. Cross-cutting work not on Trello

These are real obligations but nobody has written a card for them. Left here so they don't fall off.

### 6.1 Tests for RBAC (from old CURRENT_STATE.md "Next Steps #4")

Missing coverage:

- Role middleware in [auth_middleware.go](apps/api/internal/middleware/auth_middleware.go): verify `user` cannot hit admin routes and gets 403.
- Admin endpoint in [admin_handler.go](apps/api/internal/handler/admin_handler.go): happy path + forbidden path.
- Frontend `canAccess` logic in [access.ts](apps/web/src/lib/access.ts): table-driven test of (role, path) to allowed.

### 6.2 Deploy validation (from old CURRENT_STATE.md "Next Steps #3")

Manual smoke test on the AWS stack (Lightsail API + Amplify web + RDS + S3):

- Guest can access `/` and `/search/*`.
- Guest is redirected from `/my-list`, `/settings`, `/admin`.
- Guest sidebar only shows public links.
- Admin sees `/admin` menu and loads user list.
- Non-admin authenticated user gets `403` view in `/admin`.

Not blocking, but should happen before the next release tag.

### 6.3 AWS deploy doc (just shipped, no card)

Already delivered: [docs/DEPLOY_AWS_LIGHTSAIL.md](docs/DEPLOY_AWS_LIGHTSAIL.md) (single AWS-native path: Lightsail API + Amplify web + RDS + S3) and the extended [.env.production.example](.env.production.example). Supabase was dropped as an option; the self-hosted Postgres overlay (`docker-compose.aws.yml`), the Supabase guide, the `Dockerfile.web` web image, and the `pg_dump` backup script were removed (RDS handles backups, Amplify hosts the frontend). No Trello card, it was an explicit priority exception.

### 6.4 Local dev launcher (just shipped, no card)

Already delivered: [scripts/dev.sh](scripts/dev.sh). One entrypoint for the whole local stack: brings up the Postgres container, waits for `pg_isready`, then runs backend + frontend via the existing `bun run dev`. Subcommands: `up` (default), `down`, `reset`. Backend still owns `.env` loading and migrations on startup; the script only owns the db lifecycle and the readiness gate.

---

## 7. Backlog (TO-DO list), grouped by theme

Still gated by rule 2.4: don't start any of these until DOING is empty and the §4 cleanup is done. Notes per card focus on dependencies and whether the codebase already has the pieces.

### 7.1 Gamification & retention

Thematically consistent with the existing XP / challenges / ranks work. Reuses the same user-activity plumbing. Best cluster to tackle after Daily Streaks.

- [#19 Badges & Achievements](https://trello.com/c/iDrzYcFl): needs a `badges` table + a rule engine that subscribes to the same hooks Streaks will use. Biggest retention win of the cluster.
- [#17 The "Life RPG" Stats Page](https://trello.com/c/KbjhTfOZ): mostly a read-side aggregate on existing tables (donut, heatmap, radar). No new writes. Can ship without touching the API surface much.
- [#18 "Sprint Mode" (Focus Timer)](https://trello.com/c/5nEJ2qDb): ties into progress updates. Small scope.
- [#16 The "Decision Paralysis" Randomizer](https://trello.com/c/e4ZZQunV): tiny, client-only. Good filler between bigger cards.

### 7.2 Data & workflow

Collection management and UX polish on the core list experience.

- [#15 Collections & Private Lists](https://trello.com/c/Mc5AxtWV): requires a new `collections` table + privacy flag on list items. Check if existing `user_list_items` can grow a `collection_id` nullable FK or if a join table is cleaner.
- [#22 "Next Up" Smart Queue](https://trello.com/c/Y0x342j6): ordering layer on the existing TODO list view. Thin.
- [#24 Rich Reviews & Tags](https://trello.com/c/lgba4h8v): `completion_reviews` table with rating, date, tags, favorite quote. Feeds into future stats page (#17).
- [#20 Franchise & Universe Mapping](https://trello.com/c/y37cR3MQ): content-graph work. Probably needs external data source; defer until `#21 Connector` exists.
- [#21 The "Connector" (Data Importing)](https://trello.com/c/7gqjDxiK): Steam / MAL / Goodreads / Trakt importers. Large. Each provider is its own mini-project. Highest impact on the new-user empty state, lowest payoff per day of work. Schedule deliberately.

### 7.3 Trust, safety, abuse

- [#10 Content upload filtering](https://trello.com/c/9OPpOqSM): AI-based moderation (no human mods). Needs a provider decision (OpenAI Moderation API, perspective API, local model). Cheap to ship behind a feature flag.
- [#11 Allow users to be banned / appeal / delete account](https://trello.com/c/gMloo6bk): ban state on user, soft-delete flow, data export for GDPR. Touches auth, admin page, and storage cleanup. Medium.

### 7.4 Monetization

- [#29 Plans / Memberships (remove ads, perks)](https://trello.com/c/DlPewjt1): prerequisite for any paid-tier work. Stripe or Lemon Squeezy. Blocks #27 and #14.
- [#27 Profile customization (HTML editor, paid assets)](https://trello.com/c/DfNiwJ9s): paywalled cosmetics. Depends on #29.
- [#14 Ad Campaigns](https://trello.com/c/m8pglZXz): house-ads infrastructure, impression / click tracking. Paired with #29 (members don't see ads). Dedupe first per §4.3.

### 7.5 Tech debt

- [#25 Review implementation and un-hardcode stuff](https://trello.com/c/kVkyMidd): audit sweep. Should be broken into concrete sub-cards when picked up; as a single card it is unbounded and will never finish cleanly.

---

## 8. Verification conventions

When finishing any of the above, the minimum check set is:

```bash
go -C apps/api test ./...
bun run swagger:gen
bun run --cwd apps/web lint
bun run build
docker build -f Dockerfile.api --target builder -t completionist-api-builder-test .
docker compose --env-file .env.production -f docker-compose.prod.yml config
```

The web app is built and verified by `bun run build` and by Amplify; there is no web Docker image. The production compose stack is api-only.

For UI-touching work, also manually exercise the feature in a browser against the dev server. Type-checks and tests verify code, not UX.

---

## 9. Phased release roadmap (1.1.0 to 2.0.0)

This is the plan for finishing every Trello card and cutting a feature-complete `2.0.0`. Releases are forward-only. Minor bumps ship user-visible features. Patch bumps ship stability/tests/docs only. Major (`2.0.0`) is cut exactly once, when everything below is green.

Each phase has a fixed **exit criteria** gate. A release does not ship, and this file is not updated with the new version, until its gate is met. Gates are strict: `bun run build`, `go -C apps/api test ./...`, the `docker compose -f docker-compose.prod.yml config` check, and a manual browser pass against the dev server.

### 9.0 Release table

| Version | Type  | Theme                            | Trello cards                                  |
|---------|-------|----------------------------------|-----------------------------------------------|
| 1.1.1   | patch | Cleanup + RBAC tests + deploy    | (cleanup only, no new cards)                  |
| 1.2.0   | minor | Daily Streaks                    | #23                                           |
| 1.3.0   | minor | Activity-driven gamification     | #19, #18                                      |
| 1.4.0   | minor | Stats page + randomizer          | #17, #16                                      |
| 1.5.0   | minor | Collections / reviews / queue    | #15, #24, #22                                 |
| 1.6.0   | minor | Sharing + v1 suggestions         | #12, #6 (non-AI half)                         |
| 1.7.0   | minor | Trust & safety                   | #10, #11                                      |
| 1.8.0   | minor | Monetization foundation          | #29                                           |
| 1.9.0   | minor | Ads + profile customization      | #14, #27                                      |
| 1.10.0  | minor | Importers batch 1                | #21 (Steam, MAL, AniList)                     |
| 1.11.0  | minor | Importers batch 2                | #21 (Letterboxd, Trakt, Goodreads)            |
| 1.12.0  | minor | Franchise graph                  | #20                                           |
| 1.13.0  | minor | AI Suggestions v2 (LLM)          | #6 (AI half)                                  |
| 1.14.0  | minor | Tech debt sweep                  | #25, #2 (rescoped)                            |
| **2.0.0** | **major** | **Feature-complete cut**  | (none)                                        |

Fifteen releases. Some are heavy (1.8.0 payments, 1.10.0/1.11.0 importers) and will likely take multiple sittings. Some are light (1.4.0, 1.5.0 subsets) and can land in a single session.

---

### 9.1 Phase 1.1.1: Cleanup + RBAC tests + deploy validation

**Bump:** patch. Zero new features; stability and hygiene only.

**Work items:**

1. Trello bookkeeping per §4:
   - Verify each DONE card against code with a spot-check grep. Archive the ones that pass. Move any failures back to TO-DO.
   - Archive stale RBAC cards #1, #3, #9. Rescope #2 Security to a narrower card (rate limiting / audit review / secret rotation) or archive if nothing concrete remains.
   - Archive ad-campaigns duplicate #28 (keep #14).
   - Update §3 of this file with the post-cleanup counts.
2. Scope the two unscoped DOING cards per §5.2, §5.3:
   - Write a description on #12 Share Posts committing to an interpretation (repost / DM share / public link / external intent), or move it back to TO-DO.
   - Split #6 AI Suggestions into "v1 (non-AI)" and "v2 (LLM)" and move both back to TO-DO. Only Daily Streaks stays in DOING.
3. RBAC tests per §6.1:
   - [auth_middleware.go](apps/api/internal/middleware/auth_middleware.go): table-driven test of (role, route) to allow/deny.
   - [admin_handler.go](apps/api/internal/handler/admin_handler.go): happy path + forbidden path.
   - [access.ts](apps/web/src/lib/access.ts): (role, path) to allowed table test.
4. Deploy validation per §6.2:
   - Stand up the AWS stack (Lightsail API + Amplify web + RDS + S3). Run the guest/admin/auth smoke checklist.
   - Document the result in §6.2.

**Exit gate:**

- DONE list contains only verified implementations. BUGS list still empty.
- Cleanup counts in §3 are accurate.
- `go -C apps/api test ./...` passes including new RBAC tests.
- One deploy path is smoke-tested and documented.
- DOING list has exactly one card: `#23 Daily Streaks`.

**Ship action:** bump all four version files to `1.1.1`, regenerate Swagger, commit, tag `v1.1.1`.

---

### 9.2 Phase 1.2.0: Daily Streaks

**Bump:** minor. First new user-visible feature of the climb.

**Trello cards:** [#23 Daily Streaks](https://trello.com/c/fJJNtADu).

**Work items:** see §5.1 for the full spec. Summary:

1. Migration `004_add_user_streaks.up.sql` / `.down.sql`, table `user_streaks (user_id PK, current_streak, longest_streak, last_activity_date, tier)`.
2. `TouchStreak(userID, day)` function called from existing activity hooks (completion, progress update, challenge progress).
3. Handler `GET /api/users/:id/streak`, and include the streak object in the user profile payload so the header renders without a second round-trip.
4. Frontend streak pill in the header; three-tier icon (Spark / Fire / Inferno).
5. Unit tests for the transition matrix: 0 to 1, same-day no-op, +1 next day, skip-day reset, grace-window behavior.

**Exit gate:**

- Migration runs cleanly up + down.
- Streak increments on a real activity write, end-to-end, in a browser.
- Transition-matrix unit tests pass.
- DOING list is empty.

**Ship action:** bump to `1.2.0`. Move #23 to DONE and archive after verification.

---

### 9.3 Phase 1.3.0: Activity-driven gamification

**Bump:** minor. Reuses the activity hook that Daily Streaks just installed.

**Trello cards:** [#19 Badges & Achievements](https://trello.com/c/iDrzYcFl), [#18 Sprint Mode (Focus Timer)](https://trello.com/c/5nEJ2qDb).

**Work items:**

1. **#19 Badges:**
   - Migration `005_add_badges.up.sql`: `badges` (definition table: code, name, description, icon, criteria_json) and `user_badges` (user_id, badge_id, earned_at).
   - Rule engine that subscribes to the same activity hook `TouchStreak` uses. Each badge has a predicate `(user, event) -> bool` driven by `criteria_json`.
   - Seed the initial badge set: Binger, Polyglot, Completionist, First Steps.
   - Badge shelf on the user profile page.
   - Notification (toast or in-app) when a badge is earned.
2. **#18 Sprint Mode:**
   - Client-only initially, no backend changes needed beyond an optional `last_sprint_at` column on the user.
   - 25-minute Pomodoro timer inside the completion-item card component.
   - On finish, prompt "Did you make progress?" and call the existing progress-update endpoint.

**Exit gate:**

- Starting a real activity earns at least one seeded badge end-to-end.
- Sprint timer runs and triggers the progress prompt.
- Unit tests for the badge rule engine pass.

**Ship action:** bump to `1.3.0`. Archive #19 and #18.

---

### 9.4 Phase 1.4.0: Stats page + randomizer

**Bump:** minor. Read-side; no new writes on the core data path.

**Trello cards:** [#17 Life RPG Stats Page](https://trello.com/c/KbjhTfOZ), [#16 Decision Paralysis Randomizer](https://trello.com/c/e4ZZQunV).

**Work items:**

1. **#17 Stats page:**
   - New route `/stats` on the web app.
   - Three aggregates on existing tables: time-donut (media_type to hours), heatmap calendar (daily completion counts), genre radar (category to count).
   - API handler returns all three in one payload.
   - Use Recharts or Chart.js on the frontend.
2. **#16 Randomizer:**
   - Floating button bottom-right of the main app layout.
   - On click, pull the user's TODO list, spin a fake reel, and land on one item.
   - Purely client-side.

**Exit gate:** stats page loads with real data; randomizer picks from the real TODO list.

**Ship action:** bump to `1.4.0`. Archive #17 and #16.

---

### 9.5 Phase 1.5.0: Collections / reviews / queue

**Bump:** minor. Largest data-model touch in the climb so far.

**Trello cards:** [#15 Collections & Private Lists](https://trello.com/c/Mc5AxtWV), [#24 Rich Reviews & Tags](https://trello.com/c/lgba4h8v), [#22 "Next Up" Smart Queue](https://trello.com/c/Y0x342j6).

**Work items:**

1. **#15 Collections:**
   - Migration `006_add_collections.up.sql`: `collections` table with owner, name, is_private. Add nullable `collection_id` FK to `user_list_items`.
   - CRUD handlers + frontend folder/tab UI under `/my-list`.
   - Private collections never appear on the social feed (enforce at the feed query layer).
2. **#24 Rich Reviews & Tags:**
   - Migration `007_add_completion_reviews.up.sql`: `completion_reviews` (user_id, item_id, rating, tags TEXT[], favorite_quote, completed_at).
   - Open a modal instead of a plain checkmark when clicking "Mark as Done".
   - Reviews feed into the #17 stats page from phase 1.4.0.
3. **#22 Next Up Smart Queue:**
   - Ordering layer: add `queue_position` on `user_list_items` or a separate `user_queue` table.
   - Drag-and-drop sidebar on the TODO list.
   - Auto-promotion: when an item moves to "Completed", the top of the queue becomes "In Progress".

**Exit gate:** collections create/read/update/delete works; leaving a review writes all fields; drag-reordering persists.

**Ship action:** bump to `1.5.0`. Archive #15, #24, #22.

---

### 9.6 Phase 1.6.0: Sharing + v1 suggestions

**Bump:** minor. First recommendation surface ships.

**Trello cards:** [#12 Share Posts](https://trello.com/c/RFuAdAkT), [#6 AI Suggestions v1](https://trello.com/c/a82StdGQ) (non-AI half, the v1 card per the 1.1.1 split).

**Work items:**

1. **#12 Share Posts:** implement whatever interpretation was committed to in phase 1.1.1 scoping. Most likely: share to DM/room (reuses existing chat infra) + shareable public link with OG tags.
2. **#6 AI Suggestions v1 (non-AI):**
   - SQL-driven similarity: "users who completed X also completed Y" + category match + rating overlap.
   - New endpoint `GET /api/suggestions/for-me`.
   - Home feed "Suggested for you" rail that calls it.
   - Deterministic, cheap, no external services.

**Exit gate:** a real user gets non-empty, non-random suggestions; share-to-DM round-trips through the chat system.

**Ship action:** bump to `1.6.0`. Archive #12; move the v2 (LLM) half of #6 forward, it ships in 1.13.0.

---

### 9.7 Phase 1.7.0: Trust & safety

**Bump:** minor. First safety-critical shipment; extra care.

**Trello cards:** [#10 Content upload filtering](https://trello.com/c/9OPpOqSM), [#11 Ban / appeal / delete](https://trello.com/c/gMloo6bk).

**Work items:**

1. **#10 Content filtering:**
   - Pick a provider (OpenAI Moderation API is free and good enough for v1).
   - Run every uploaded image and every post/comment body through the moderation endpoint before persisting.
   - Behind a feature flag so it can be toggled if the provider misbehaves.
   - Flagged content goes to an admin review queue (new table `moderation_queue`).
2. **#11 Ban / appeal / account deletion:**
   - Migration: `users.banned_at`, `users.banned_reason`, new table `ban_appeals` and `account_deletion_requests`.
   - Admin action to ban/unban from the existing admin page.
   - User-facing flow: see ban, submit appeal.
   - GDPR data-export endpoint (user downloads their data).
   - Soft-delete account flow with 30-day grace window.

**Exit gate:** uploading a deliberately-offensive test string gets flagged; admin can ban a user and that user is immediately blocked; data export produces a complete JSON.

**Ship action:** bump to `1.7.0`. Archive #10 and #11.

---

### 9.8 Phase 1.8.0: Monetization foundation

**Bump:** minor. Payments are risky; this is the most-likely place to stall.

**Trello cards:** [#29 Plans / Memberships](https://trello.com/c/DlPewjt1).

**Work items:**

1. Pick Stripe (most mature for our use case) or Lemon Squeezy (simpler but fewer features).
2. New `subscriptions` table, webhook endpoint, signature verification.
3. Add a `member` role between `user` and `admin`, fits the existing RBAC enum without breaking it.
4. Feature gating: one helper `isMember(user)` consumed wherever perks apply.
5. `/upgrade` page with plan cards and checkout redirect.
6. Handle the full webhook lifecycle: `checkout.completed`, `invoice.paid`, `customer.subscription.deleted`, `invoice.payment_failed`.
7. Test with Stripe test mode end-to-end.

**Exit gate:** a test user can upgrade, get `member` role, downgrade, and lose the role, all via real webhook events in test mode.

**Ship action:** bump to `1.8.0`. Archive #29. This release unblocks 1.9.0 and 1.11.0 perks.

---

### 9.9 Phase 1.9.0: Ads + profile customization

**Bump:** minor. Depends on 1.8.0.

**Trello cards:** [#14 Ad Campaigns](https://trello.com/c/m8pglZXz), [#27 Profile customization](https://trello.com/c/DfNiwJ9s).

**Work items:**

1. **#14 Ads:**
   - House-ads infrastructure: `ad_campaigns` and `ad_impressions` tables.
   - Admin panel to create/schedule campaigns.
   - Frontend: ad slot components that skip rendering for `member` role.
   - Click and impression tracking via a small beacon endpoint.
2. **#27 Profile customization:**
   - Sanitized HTML editor for profile bio (use DOMPurify server-side).
   - Asset library of paid backgrounds, badges, and fonts.
   - `isMember` check before applying custom CSS / HTML.
   - Free users get a static text bio.

**Exit gate:** a member account sees zero ads and can apply a customization; a free account sees ads and cannot.

**Ship action:** bump to `1.9.0`. Archive #14 and #27.

---

### 9.10 Phases 1.10.0 and 1.11.0: Data importers

**Bump:** minor x 2. This is the biggest card on the board and should be split.

**Trello card:** [#21 The "Connector"](https://trello.com/c/7gqjDxiK).

**1.10.0: Steam, MyAnimeList, AniList**

- OAuth flow for each provider (Steam uses OpenID, MAL + AniList use OAuth2).
- Per-provider importer with rate-limited, resumable pagination.
- Map provider entities to the existing `media_items` schema; create if missing, link to `user_list_items` on match.
- UI: "Connect your accounts" section in settings, per-provider connect button, sync status, last-synced-at.
- Background job queue for the import work (start simple: goroutine + channel, upgrade later if needed).

**1.11.0: Letterboxd, Trakt, Goodreads**

- Letterboxd has no public API; use the exportable CSV as the import path.
- Trakt has a clean OAuth + REST API.
- Goodreads OAuth was deprecated in 2020; use their CSV export or a scraping fallback. Decide at implementation time.
- Same architecture as 1.10.0, three more providers.

**Exit gate per release:** a real user account can connect the listed providers, kick off an import, and see library items populate without duplicates.

**Ship action:** bump to `1.10.0` after the first batch, `1.11.0` after the second. Only archive the single [#21](https://trello.com/c/7gqjDxiK) card after `1.11.0` ships.

---

### 9.11 Phase 1.12.0: Franchise graph

**Bump:** minor. Depends on 1.10.0/1.11.0 importers landing real external IDs.

**Trello card:** [#20 Franchise & Universe Mapping](https://trello.com/c/y37cR3MQ).

**Work items:**

1. New `franchises` and `franchise_items` tables. Each franchise is a directed graph (prequel, sequel, adaptation).
2. Seed from TMDB collections (movies), MAL relations (anime), and manual curation for the long tail.
3. Render a franchise tree on media detail pages.
4. Hook into the Completionist badge from #19: completing a full franchise unlocks a special badge.

**Exit gate:** viewing a media item that belongs to a franchise shows the full franchise tree and highlights completion status.

**Ship action:** bump to `1.12.0`. Archive #20.

---

### 9.12 Phase 1.13.0: AI Suggestions v2

**Bump:** minor. Opt-in and rate-limited from day one.

**Trello card:** [#6 AI Suggestions v2](https://trello.com/c/a82StdGQ) (LLM half from the 1.1.1 split).

**Work items:**

1. Provider abstraction: `SuggestionProvider` interface with an LLM implementation (Claude or OpenAI).
2. Prompt: pass the user's recent completions + ratings, get back 5 suggestions with rationale.
3. Rate limit per user: N calls/day, enforced in middleware.
4. Opt-in in user settings (default off for privacy and cost).
5. Fall back to the v1 non-AI suggestions when disabled, over-limit, or provider errors.
6. Log token usage for cost tracking.

**Exit gate:** an opted-in test user gets LLM suggestions; an opted-out user still gets v1 suggestions.

**Ship action:** bump to `1.13.0`. Archive the remaining half of #6.

---

### 9.13 Phase 1.14.0: Tech debt sweep

**Bump:** minor, broad but internal.

**Trello cards:** [#25 Review implementation, un-hardcode](https://trello.com/c/kVkyMidd), [#2 Security](https://trello.com/c/AvaMbhAF) (whatever remains after 1.1.1 rescoping).

**Work items:**

1. Audit sweep: every file with magic numbers, hardcoded URLs, hardcoded secrets, or inline config. Move to `config.go` / env vars / database seeds.
2. Re-read the whole codebase with a fresh eye for DRY opportunities, but only actually refactor where a second use already exists (no speculative abstractions).
3. Security pass: rate limiting on write endpoints, audit-log review, secret rotation runbook, dependency audit (`go mod tidy`, `bun audit`).
4. Final README refresh.

**Exit gate:** grep for hardcoded values returns only values that _should_ be hardcoded (constants, well-known magic numbers). Security checklist ticked off.

**Ship action:** bump to `1.14.0`. Archive #25 and whatever remains of #2.

---

### 9.14 Phase 2.0.0: Feature-complete cut

**Bump:** major. Cut exactly once, when everything above is green.

**Work items:**

1. Final smoke test on the AWS stack (Lightsail API + Amplify web + RDS + S3).
2. Verify TO-DO is empty, DOING is empty, BUGS is empty. If any card landed since the plan was written, decide: include in 2.0.0 or defer to 2.1.0.
3. Write `CHANGELOG.md` covering the entire 1.1.0 to 2.0.0 arc. One section per minor release, linking to the Trello card each one closed.
4. Regenerate Swagger and pin it as the `v2` API surface. From here on, breaking API changes require a `3.0.0`.
5. Tag `v2.0.0` and publish images to ECR with the `2.0.0` tag (not `latest`).
6. Update [README.md](README.md) with a "Project status: 2.0.0, feature-complete" banner.

**Exit gate (the only gate that matters):**

- Both deploys pass the smoke checklist in §6.2.
- `CHANGELOG.md` exists and covers the full arc.
- Trello board: TO-DO empty, DOING empty, BUGS empty. DONE either archived or held intentionally as history.
- Git tag `v2.0.0` is pushed.

**Ship action:** celebrate. Then open the `post-2.0` planning cycle as a separate effort, not as part of this plan.

---

## 10. How to use this file from here

- After every phase ships, update §1 (current version), §3 (Trello counts), and the relevant §9 section (mark phase done).
- Do not modify §9's ordering or scope without a paired note explaining why.
- If a card is added to Trello mid-climb, append it to the most thematically-appropriate phase and note the addition in that phase's **Work items** list.
- If a phase's scope balloons past a single release, split it (the importers are already split across 1.10 and 1.11 for this reason).
- Version bump files to touch every release: [VERSION](VERSION), [package.json](package.json), [apps/web/package.json](apps/web/package.json), Swagger metadata under [apps/api/swagger/](apps/api/swagger/). Regenerate Swagger with `bun run swagger:gen`.
