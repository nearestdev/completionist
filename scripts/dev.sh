#!/usr/bin/env bash
set -euo pipefail

# Local dev launcher for the Completionist monorepo.
#
# Brings up the Postgres container, waits until it accepts connections, then
# runs the backend (Go) and frontend (Next.js) together via `bun run dev`.
# The Go backend loads .env and applies migrations on startup, so this script
# only owns the database lifecycle and the readiness gate before launching.
#
# Usage:
#   scripts/dev.sh         start db (if needed), then backend + frontend
#   scripts/dev.sh down    stop the Postgres container (volume/data preserved)
#   scripts/dev.sh reset   wipe the Postgres volume, then start everything fresh
#
# Ctrl+C stops the backend/frontend; the db container is left running so the
# next start is fast and data survives. Use `scripts/dev.sh down` to stop it.

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if [[ -t 1 ]]; then
  BOLD=$'\e[1m'; DIM=$'\e[2m'; RED=$'\e[31m'; GRN=$'\e[32m'; YLW=$'\e[33m'; CYN=$'\e[36m'; RST=$'\e[0m'
else
  BOLD=; DIM=; RED=; GRN=; YLW=; CYN=; RST=
fi
info() { printf '%s\n' "${CYN}${BOLD}==>${RST} ${BOLD}$1${RST}"; }
ok()   { printf '%s\n' "${GRN}ok${RST}  $1"; }
warn() { printf '%s\n' "${YLW}!${RST}   $1"; }
die()  { printf '%s\n' "${RED}error:${RST} $1" >&2; exit 1; }

case "${1:-}" in
  -h|--help|help)
    printf '%s\n' \
      "Usage: scripts/dev.sh [up|down|reset]" \
      "  up     (default) start Postgres, then backend + frontend" \
      "  down   stop the Postgres container (data preserved)" \
      "  reset  wipe the Postgres volume, then start everything fresh"
    exit 0
    ;;
esac

command -v docker >/dev/null 2>&1 || die "docker is not installed or not on PATH"
docker info >/dev/null 2>&1 || die "the Docker daemon is not running"

if docker compose version >/dev/null 2>&1; then
  DC=(docker compose)
elif command -v docker-compose >/dev/null 2>&1; then
  DC=(docker-compose)
else
  die "neither 'docker compose' nor 'docker-compose' is available"
fi

env_val() { grep -E "^$1=" .env 2>/dev/null | tail -n1 | cut -d= -f2- || true; }
FRONTEND_PORT="$(env_val PORT)";        FRONTEND_PORT="${FRONTEND_PORT:-4000}"
BACKEND_PORT="$(env_val BACKEND_PORT)"; BACKEND_PORT="${BACKEND_PORT:-5000}"

wait_for_db() {
  info "Waiting for Postgres to accept connections"
  local tries=0
  until "${DC[@]}" exec -T db pg_isready -U postgres -d completionist_db >/dev/null 2>&1; do
    tries=$((tries + 1))
    (( tries >= 60 )) && die "Postgres did not become ready within 60s"
    sleep 1
  done
  ok "Postgres is ready (localhost:5435)"
}

case "${1:-up}" in
  up) ;;
  reset)
    warn "Wiping the Postgres volume (all local data will be lost)"
    "${DC[@]}" down -v
    ;;
  down)
    "${DC[@]}" down
    ok "Postgres container stopped (data preserved)"
    exit 0
    ;;
  *)
    die "unknown command '$1' (use: up | down | reset)"
    ;;
esac

command -v go  >/dev/null 2>&1 || die "go is not installed or not on PATH"
command -v bun >/dev/null 2>&1 || die "bun is not installed or not on PATH"

if [[ ! -d node_modules ]]; then
  info "Installing JS dependencies (bun install)"
  bun install
fi

info "Starting Postgres (docker compose service: db)"
"${DC[@]}" up -d db
wait_for_db

printf '\n'
ok "Postgres   localhost:5435  ${DIM}(db: completionist_db)${RST}"
ok "Frontend   ${BOLD}http://localhost:${FRONTEND_PORT}${RST}"
ok "Backend    ${BOLD}http://localhost:${BACKEND_PORT}${RST}  ${DIM}(swagger: /swagger/index.html)${RST}"
printf '%s\n\n' "${DIM}DB keeps running after you quit; stop it with: scripts/dev.sh down${RST}"

info "Launching backend + frontend via 'bun run dev' (Ctrl+C to stop)"
exec bun run dev
