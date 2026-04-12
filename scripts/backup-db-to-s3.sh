#!/usr/bin/env bash
set -euo pipefail

# Nightly Postgres backup for the self-hosted db service in docker-compose.aws.yml.
# Dumps the database, gzips it into ./backups (mounted into the db container),
# and ships the archive to s3://${AWS_BUCKET}/db-backups/.
#
# Wire via cron (on the Lightsail host):
#   0 4 * * * /opt/completionist/scripts/backup-db-to-s3.sh >> /var/log/completionist-backup.log 2>&1

PROJECT_DIR="${PROJECT_DIR:-/opt/completionist}"
ENV_FILE="${PROJECT_DIR}/.env.production"

# shellcheck disable=SC1090
set -a; source "${ENV_FILE}"; set +a

: "${AWS_BUCKET:?AWS_BUCKET must be set in ${ENV_FILE}}"
: "${POSTGRES_USER:=completionist}"
: "${POSTGRES_DB:=completionist}"

STAMP=$(date -u +%Y%m%dT%H%M%SZ)
DUMP_NAME="completionist-${STAMP}.sql.gz"
HOST_BACKUP_DIR="${PROJECT_DIR}/backups"
mkdir -p "${HOST_BACKUP_DIR}"

docker compose --env-file "${ENV_FILE}" \
  -f "${PROJECT_DIR}/docker-compose.prod.yml" \
  -f "${PROJECT_DIR}/docker-compose.aws.yml" \
  exec -T db pg_dump -U "${POSTGRES_USER}" "${POSTGRES_DB}" \
  | gzip > "${HOST_BACKUP_DIR}/${DUMP_NAME}"

aws s3 cp "${HOST_BACKUP_DIR}/${DUMP_NAME}" "s3://${AWS_BUCKET}/db-backups/${DUMP_NAME}"

# Keep the 14 most recent local archives.
ls -1t "${HOST_BACKUP_DIR}"/completionist-*.sql.gz 2>/dev/null \
  | tail -n +15 \
  | xargs -r rm --

echo "[$(date -u +%FT%TZ)] backup ok: ${DUMP_NAME}"
