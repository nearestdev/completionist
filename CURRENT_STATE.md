# Current State

## Product Version

- Current product version: `1.0.5`
- SemVer policy: use semantic versioning, but do not force tiny increments. Jump the patch, minor, or major number according to actual scope.
- Current baseline rule for this repo: `VERSION`, the root [`package.json`](/mnt/Storage/personal-projects/completionist-api-go/package.json), the web [`package.json`](/mnt/Storage/personal-projects/completionist-api-go/apps/web/package.json), and API Swagger versioning should stay aligned when the product version changes.

## Product Snapshot

- Monorepo with a Go API in [`apps/api`](/mnt/Storage/personal-projects/completionist-api-go/apps/api) and a Next.js web app in [`apps/web`](/mnt/Storage/personal-projects/completionist-api-go/apps/web).
- Backend already includes auth, users, lists/media, posts, messaging, websocket chat, attachments, challenges, rank/xp services, and external integrations for TMDB, Jikan, RAWG, Google Books, Steam, and Last.fm.
- Frontend already includes auth flows, search pages, profile pages, posts, messages, rooms, challenges, wishlists, my-list views, and trending UI sections.
- The worktree is currently dirty in `apps/web` and some package/lock files. Treat those edits as in-progress user work unless a future task explicitly says otherwise.

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

- The board still contains the backlog in `TO-DO`. Do not pull a new card into `DOING` until `DOING` is empty.

## Current Operating Rule

- Always review `DONE` cards against the code before archiving them.
- Archive a `DONE` card only when the implementation is actually present and usable end to end.
- Work `DOING` one card at a time unless two items are trivial and tightly coupled.
- Only move a card from `TO-DO` to `DOING` when `DOING` is empty.
- Every implementation change must update this file with the current state and next steps.

## Current Focus

- This pass established the workflow scaffolding, version baseline, Trello operating rules, and the global project-specific Codex skill.
- No Trello cards were moved or archived in this setup pass.
- The highest-priority active delivery item is `Daily Streaks (Habit Tracking)` because it is already in `DOING` and was touched most recently.

## Next Steps

1. Verify each `DONE` card against the current codebase and archive only the ones that are truly complete.
2. Finish `Daily Streaks (Habit Tracking)` before starting anything new from `TO-DO`.
3. Reassess `Share posts functionality` and the recommendations card after the streaks work is complete.
4. Once `DOING` is empty, pick a single high-leverage backlog item and move it into `DOING`.
