# Deploy Guide: AWS Lightsail + ECR + S3 + Self-Hosted Postgres

This is the AWS-native deploy path. It runs the entire stack (`api`, `web`, and **a self-hosted Postgres**) on a single Amazon Lightsail instance, pulls images from Amazon ECR, and uses Amazon S3 for media uploads and nightly database backups. No managed database service is required.

For the alternative deploy that uses Supabase as the managed Postgres provider, see [DEPLOY_SUPABASE_VPS.md](DEPLOY_SUPABASE_VPS.md). The two guides share the same `Dockerfile.api` / `Dockerfile.web` / `docker-compose.prod.yml`; this guide layers an extra `docker-compose.aws.yml` override on top to add the `db` service.

The backend automatically runs SQL migrations at startup via golang-migrate (see [apps/api/internal/database/database.go](../apps/api/internal/database/database.go)), so there is no separate migration step.

## 0. Architecture

```
                        ┌──────────────────────────────────────────────┐
                        │             Lightsail instance               │
                        │       (Ubuntu 24.04, 2 GB / 2 vCPU)          │
                        │                                              │
   users ── HTTPS ──►   │   ┌────────┐   ┌────────┐   ┌────────────┐   │
                        │   │  web   │──►│  api   │──►│  db (pg16) │   │
                        │   │ :4000  │   │ :8080  │   │  :5432     │   │
                        │   └────────┘   └────────┘   └────────────┘   │
                        │                    │              │          │
                        └────────────────────┼──────────────┼──────────┘
                                             │              │
                                             ▼              ▼
                                     ┌──────────────┐  ┌──────────────┐
                                     │     S3       │  │     S3       │
                                     │ media/       │  │ db-backups/  │
                                     └──────────────┘  └──────────────┘

Image registry: ECR  →  pulled by Lightsail at deploy time.
```

## 1. Cost sketch

| Resource | Plan | Approx. monthly |
|---|---|---|
| Lightsail instance | 2 GB RAM / 2 vCPU / 60 GB SSD | $12 |
| Lightsail static IP | attached to a running instance | $0 |
| S3 (media + db backups) | low-traffic project | $1 to 3 |
| ECR | first 500 MB free, then $0.10/GB-month | $0 to 1 |
| Data transfer out | included in Lightsail bundle (2 TB) | $0 |
| **Total** | | **~$13 to 16** |

If the Lightsail box ever runs out of headroom, the upgrade path is to move Postgres to RDS (see the FAQ at the bottom).

## 2. Create the S3 bucket

1. AWS console → S3 → **Create bucket**.
2. Pick a region close to your Lightsail region (e.g. `us-east-1`). Use the same region for everything.
3. Block all public access: **ON**.
4. Bucket versioning: optional but recommended for the backups prefix.
5. Inside the bucket, create two prefixes (folders): `media/` and `db-backups/`.

Then create an IAM user with a scoped policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:PutObject", "s3:GetObject", "s3:DeleteObject"],
      "Resource": "arn:aws:s3:::<your-bucket>/*"
    },
    {
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": "arn:aws:s3:::<your-bucket>"
    }
  ]
}
```

Generate an access key for that user. You'll plug those values into `.env.production` as:

- `AWS_REGION`
- `AWS_BUCKET`
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `STORAGE_DRIVER=s3`

The Go API reads these in [apps/api/internal/config/config.go](../apps/api/internal/config/config.go).

## 3. Create the ECR repositories

Two repositories, one per image:

```bash
aws ecr create-repository --repository-name completionist-api --region <region>
aws ecr create-repository --repository-name completionist-web --region <region>
```

Optional but recommended: set a lifecycle policy to keep only the last 10 images per repo so storage doesn't drift upward.

> ECR is **optional**. If your local machine is slow or you don't want a registry, you can build directly on the Lightsail instance the way [DEPLOY_SUPABASE_VPS.md](DEPLOY_SUPABASE_VPS.md) does (`git pull && docker compose build`). ECR is the right call if you want immutable, SHA-tagged images and/or you build in CI.

## 4. Build and push images to ECR

From your dev machine, with Docker buildx and the AWS CLI logged in:

```bash
ACCOUNT_ID=<your-aws-account-id>
REGION=<your-region>
TAG=$(git rev-parse --short HEAD)
REGISTRY="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"

aws ecr get-login-password --region "${REGION}" \
  | docker login --username AWS --password-stdin "${REGISTRY}"

# api
docker buildx build --platform linux/amd64 \
  -t "${REGISTRY}/completionist-api:${TAG}" \
  -t "${REGISTRY}/completionist-api:latest" \
  -f Dockerfile.api . --push

# web: build args bake the API URL at build time, so set them now
docker buildx build --platform linux/amd64 \
  --build-arg BACKEND_BASE_URL=https://api.<your-domain> \
  --build-arg BACKEND_PORT=443 \
  -t "${REGISTRY}/completionist-web:${TAG}" \
  -t "${REGISTRY}/completionist-web:latest" \
  -f Dockerfile.web . --push
```

`--platform linux/amd64` matters if you build on an Apple Silicon Mac. Lightsail instances are x86_64.

## 5. Create the Lightsail instance

1. AWS console → Lightsail → **Create instance**.
2. Region: same as your S3 bucket.
3. Platform: Linux/Unix → OS Only → **Ubuntu 24.04 LTS**.
4. Instance plan: **$12 / 2 GB / 2 vCPU / 60 GB SSD** minimum. The 1 GB plan is too small for Postgres + API + web + a Next.js production server side-by-side.
5. Name it (e.g. `completionist-prod`) and create.
6. Once running, attach a **Static IP** under Networking → Static IPs.
7. Open firewall ports under Networking → IPv4 Firewall: `22/tcp` (SSH, ideally restricted to your IP), `80/tcp`, `443/tcp`. You do **not** need to expose `4000` or `8080` once Caddy is in front; before that, temporarily allow `4000/tcp` and `8080/tcp` for smoke testing.

## 6. Install Docker on the instance

SSH in (use the Lightsail browser console or download the SSH key) and run the same install steps as the Supabase guide:

```bash
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin git awscli
sudo usermod -aG docker $USER
newgrp docker
```

Note that this also installs `awscli` (used by the backup script).

## 7. Configure the AWS CLI on the instance

Same IAM credentials you generated in step 2:

```bash
aws configure
# AWS Access Key ID: <key>
# AWS Secret Access Key: <secret>
# Default region name: <region>
# Default output format: json
```

This is what lets `scripts/backup-db-to-s3.sh` push dumps to S3.

## 8. Clone the repo and write `.env.production`

```bash
sudo mkdir -p /opt/completionist
sudo chown $USER:$USER /opt/completionist
git clone <YOUR_REPO_URL> /opt/completionist
cd /opt/completionist
cp .env.production.example .env.production
```

Edit `.env.production` and:

- **Comment out** the Supabase `DATABASE_URL` (Option A).
- **Uncomment** the Option B block and set `POSTGRES_PASSWORD` to a strong random value.
- Set `JWT_SECRET` to a strong random value.
- Set `BACKEND_BASE_URL`, `FRONTEND_BASE_URL`, `ALLOWED_ORIGINS`, `WEB_BACKEND_BASE_URL`, `WEB_BACKEND_PORT` to whatever your final domain (or static IP) will be.
- Set `STORAGE_DRIVER=s3` and fill in `AWS_REGION`, `AWS_BUCKET`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`.

The connection string `postgres://completionist:<POSTGRES_PASSWORD>@db:5432/completionist?sslmode=disable` works because `db` is the Compose service name and the connection never leaves the Docker network.

## 9. Pull from ECR and start the stack

```bash
ACCOUNT_ID=<your-aws-account-id>
REGION=<your-region>
REGISTRY="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"

aws ecr get-login-password --region "${REGION}" \
  | docker login --username AWS --password-stdin "${REGISTRY}"
```

If you're using ECR-built images, edit `docker-compose.prod.yml` (or add a second override) so `api.image` and `web.image` point at the ECR URIs you pushed in step 4 (e.g. `${REGISTRY}/completionist-api:${TAG}`). If you're building on the box instead, leave the file alone and add `--build` to the `up` command below.

Then bring it up using both compose files:

```bash
docker compose --env-file .env.production \
  -f docker-compose.prod.yml \
  -f docker-compose.aws.yml \
  pull

docker compose --env-file .env.production \
  -f docker-compose.prod.yml \
  -f docker-compose.aws.yml \
  up -d
```

The `db` service comes from [docker-compose.aws.yml](../docker-compose.aws.yml). Without that override file the stack falls back to expecting an external Postgres (the Supabase path).

golang-migrate runs at API startup automatically against the in-stack Postgres. The migrations create the `pgcrypto` extension before any `gen_random_uuid()` call, so vanilla `postgres:16-alpine` is sufficient; no custom image needed.

## 10. Verify

```bash
docker compose --env-file .env.production \
  -f docker-compose.prod.yml \
  -f docker-compose.aws.yml \
  ps

docker compose --env-file .env.production \
  -f docker-compose.prod.yml \
  -f docker-compose.aws.yml \
  logs -f api
```

The api logs should show migrations applied and the server listening on `:8080`. From the instance:

```bash
curl -I http://localhost:4000
curl -I http://localhost:8080/swagger/index.html
```

From your browser (using the static IP, before Caddy is in front):

- Frontend: `http://<STATIC_IP>:4000`
- API docs: `http://<STATIC_IP>:8080/swagger/index.html`

Confirm the schema landed by listing tables in the running db container:

```bash
docker compose --env-file .env.production \
  -f docker-compose.prod.yml \
  -f docker-compose.aws.yml \
  exec db psql -U completionist completionist -c "\dt"
```

You should see ~25 tables (`users`, `posts`, `rooms`, `room_messages`, `direct_messages`, etc.).

## 11. HTTPS with Caddy

Once the stack is up on `:4000` / `:8080`, put Caddy in front of it for automatic Let's Encrypt and clean URLs. Create `docker-compose.caddy.yml` (a third override):

```yaml
services:
  caddy:
    image: caddy:2-alpine
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - caddy_data:/data
      - caddy_config:/config
    depends_on:
      - api
      - web

volumes:
  caddy_data:
  caddy_config:
```

`Caddyfile`:

```
app.<your-domain> {
    reverse_proxy web:4000
}

api.<your-domain> {
    reverse_proxy api:8080
}
```

Then update `.env.production`:

- `BACKEND_BASE_URL=https://api.<your-domain>`
- `FRONTEND_BASE_URL=https://app.<your-domain>`
- `ALLOWED_ORIGINS=https://app.<your-domain>`
- `WEB_BACKEND_BASE_URL=https://api.<your-domain>`
- `WEB_BACKEND_PORT=443`

The web image bakes the API URL at build time (see the `BACKEND_BASE_URL` build-arg in [Dockerfile.web](../Dockerfile.web)), so after changing `WEB_BACKEND_BASE_URL` you must **rebuild and re-push the web image**, then `pull && up -d` on the box. Same gotcha as section 9 of [DEPLOY_SUPABASE_VPS.md](DEPLOY_SUPABASE_VPS.md).

Bring Caddy up:

```bash
docker compose --env-file .env.production \
  -f docker-compose.prod.yml \
  -f docker-compose.aws.yml \
  -f docker-compose.caddy.yml \
  up -d
```

Once verified working, close ports `4000` and `8080` in the Lightsail firewall. Only `80` and `443` need to be public.

## 12. Backups

Install the cron job:

```bash
sudo cp /opt/completionist/scripts/backup-db-to-s3.sh /opt/completionist/scripts/backup-db-to-s3.sh
sudo chmod +x /opt/completionist/scripts/backup-db-to-s3.sh
sudo crontab -e
# add:
0 4 * * * /opt/completionist/scripts/backup-db-to-s3.sh >> /var/log/completionist-backup.log 2>&1
```

The script reads `.env.production`, runs `pg_dump` inside the `db` container, gzips the result into `/opt/completionist/backups/`, ships it to `s3://<AWS_BUCKET>/db-backups/`, and prunes anything older than the 14 most recent local archives. See [scripts/backup-db-to-s3.sh](../scripts/backup-db-to-s3.sh).

To test it manually:

```bash
sudo /opt/completionist/scripts/backup-db-to-s3.sh
aws s3 ls "s3://<AWS_BUCKET>/db-backups/"
```

To **restore** from a backup:

```bash
aws s3 cp s3://<AWS_BUCKET>/db-backups/completionist-<STAMP>.sql.gz .
gunzip -c completionist-<STAMP>.sql.gz \
  | docker compose --env-file .env.production \
      -f docker-compose.prod.yml \
      -f docker-compose.aws.yml \
      exec -T db psql -U completionist completionist
```

For a clean restore you'll usually want to drop and recreate the database first.

## 13. Update / redeploy

Each release:

1. Locally (or in CI): rebuild and push the images with a new SHA tag.
   ```bash
   TAG=$(git rev-parse --short HEAD)
   docker buildx build --platform linux/amd64 -t "${REGISTRY}/completionist-api:${TAG}" -f Dockerfile.api . --push
   docker buildx build --platform linux/amd64 -t "${REGISTRY}/completionist-web:${TAG}" -f Dockerfile.web . --push
   ```
2. SSH to the Lightsail instance, update the image tags in `docker-compose.prod.yml` (or whichever override holds them), and:
   ```bash
   docker compose --env-file .env.production \
     -f docker-compose.prod.yml \
     -f docker-compose.aws.yml \
     pull
   docker compose --env-file .env.production \
     -f docker-compose.prod.yml \
     -f docker-compose.aws.yml \
     up -d
   ```

The api applies any new migrations on startup. The `db` volume is preserved across updates.

---

## FAQ

### Why not Aurora or RDS?

Cost. Aurora is overkill for a single-product workload at this stage, and even RDS Postgres `db.t4g.micro` adds ~$13/mo for backups + multi-AZ failover you don't need yet. Self-hosted Postgres on the same Lightsail box costs $0 extra and has the same schema compatibility. The upgrade path is real and painless if you outgrow the box:

1. Spin up RDS Postgres in the same region.
2. Stop the `db` service: `docker compose ... stop db`.
3. Restore the latest S3 backup into RDS with `pg_restore`.
4. Point `DATABASE_URL` in `.env.production` at the RDS endpoint (use `sslmode=require`).
5. Drop `docker-compose.aws.yml` from the compose command (back to the Supabase-style path).
6. `docker compose --env-file .env.production -f docker-compose.prod.yml up -d`.

No Go code changes, no schema changes.

### Why not SQLite?

It was the original idea for this deploy and I considered it carefully. It's the wrong choice for this codebase:

- The schema uses **8 custom Postgres ENUM types**, **6+ JSONB columns with GIN indexes**, the **pgcrypto** extension, and `gen_random_uuid()`. None of these map cleanly to SQLite; every one would need a schema rewrite.
- The app has a **WebSocket chat system** (rooms, DMs, message reactions). SQLite is single-writer and would hit `SQLITE_BUSY` under any concurrent message load.
- A nightly `pg_dump` on a 2 GB Lightsail box already gives you the "DB-as-a-file" durability story you wanted, without sacrificing the features above.

Postgres-in-Docker hits the same cost target as SQLite ($0 extra) without any of the trade-offs.

### Why a separate `docker-compose.aws.yml` instead of editing `docker-compose.prod.yml`?

So the existing Supabase deploy path documented in [DEPLOY_SUPABASE_VPS.md](DEPLOY_SUPABASE_VPS.md) keeps working unchanged. Using `docker-compose.prod.yml` alone gives you the external-DB stack; adding `-f docker-compose.aws.yml` layers on the in-stack `db` service and the `api → db` healthcheck dependency. Both deploy targets stay first-class.
