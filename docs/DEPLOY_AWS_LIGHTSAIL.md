# Deploy Guide: AWS (Lightsail API + Amplify Web + RDS Postgres + S3)

This is the AWS deploy path for Completionist. The Go API runs in Docker on a single Amazon Lightsail instance behind Caddy for TLS. The Next.js frontend is hosted on AWS Amplify, built straight from this repo. PostgreSQL is managed by Amazon RDS, and media uploads go to Amazon S3.

The backend applies its SQL migrations at startup via golang-migrate (see [apps/api/internal/database/database.go](../apps/api/internal/database/database.go)), so there is no separate migration step.

> **Next.js version note.** Amplify Hosting documents managed Next.js SSR support through Next.js 15, and `apps/web` is on Next.js 16. Confirm Amplify's current Next.js support before the first deploy. If 16 is not yet supported, either pin the web app to Next.js 15 or, as a stopgap, run the Next server in a container on the Lightsail box behind the same Caddy. The API, RDS, and S3 parts of this guide are unaffected.

## 0. Architecture

```
        ┌────────────────────────────┐
        │  AWS Amplify (Next.js web)  │  CDN + TLS, built from GitHub
        └─────────────┬──────────────┘
                      │ HTTPS  (calls api.<domain>)
   users ── HTTPS ──► │
                      ▼
        ┌────────────────────────────┐
        │  Lightsail instance         │
        │  Caddy (TLS) -> Go API:8080 │  Docker, Dockerfile.api
        └──────┬───────────────┬──────┘
               │               │
               ▼               ▼
        ┌─────────────┐  ┌─────────────┐
        │ RDS Postgres │  │ S3 (media/) │
        └─────────────┘  └─────────────┘
```

## 1. Cost sketch

| Resource | Plan | Approx. monthly |
|---|---|---|
| Lightsail instance | 2 GB RAM / 2 vCPU / 60 GB SSD | $12 |
| Lightsail static IP | attached to a running instance | $0 |
| RDS Postgres | `db.t4g.micro`, single-AZ, 20 GB gp3 | $14 to 16 |
| Amplify Hosting | low traffic (build minutes + served GB) | $0 to 2 |
| S3 | media for a low-traffic project | $1 to 3 |
| Data transfer | Lightsail bundle includes 2 TB | $0 |
| **Total** | | **~$27 to 33** |

RDS is the main line item versus self-hosting Postgres on the box. You pay for it to get managed automated backups, point-in-time recovery, and snapshots, so there is no backup script to run and no cron to babysit.

## 2. S3 bucket and IAM (media)

1. AWS console, S3, **Create bucket**. Pick one region and use it for everything (Lightsail, RDS, S3).
2. Block all public access: **ON**. The API serves media through the app, so the bucket stays private.
3. Create one prefix: `media/`.

Create an IAM user with a scoped policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    { "Effect": "Allow", "Action": ["s3:PutObject", "s3:GetObject", "s3:DeleteObject"], "Resource": "arn:aws:s3:::<your-bucket>/*" },
    { "Effect": "Allow", "Action": ["s3:ListBucket"], "Resource": "arn:aws:s3:::<your-bucket>" }
  ]
}
```

Generate an access key. These map to `.env.production`:

- `STORAGE_DRIVER=s3`
- `AWS_REGION`, `AWS_BUCKET`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`

The Go API reads them in [apps/api/internal/config/config.go](../apps/api/internal/config/config.go).

## 3. RDS Postgres

1. AWS console, RDS, **Create database**, Standard create, **PostgreSQL 16**.
2. Template: Free tier or Dev/Test. Instance: **`db.t4g.micro`**. Storage: **20 GB gp3**, single-AZ to start.
3. Set the master username, a strong password, and an initial database name (e.g. `completionist`).
4. Connectivity: place it in the **default VPC** of your region, **Public access: No**.
5. Keep automated backups on (7 day retention is fine). Backups plus snapshots are what replace the old `pg_dump` cron.

Build the connection string for `.env.production`:

```
DATABASE_URL=postgres://<master_user>:<password>@<rds-endpoint>:5432/<db_name>?sslmode=require
```

`sslmode=require` matters: RDS encrypts connections in transit. The API runs its migrations against this database at startup. Networking is covered in the next step, since the Lightsail box has to reach this private endpoint.

## 4. Lightsail instance and VPC peering to RDS

1. Lightsail, **Create instance**, same region as RDS and S3. Linux/Unix, OS Only, **Ubuntu 24.04 LTS**.
2. Plan: **$12 / 2 GB / 2 vCPU**. Attach a **Static IP** under Networking.
3. IPv4 firewall: allow `22/tcp` (ideally your IP only), `80/tcp`, `443/tcp`.
4. **Enable Lightsail VPC peering**: Lightsail home, Account, Advanced, "Enable VPC peering" for the RDS region. This peers the Lightsail network with your region's default VPC, where RDS lives.
5. On the **RDS security group**, add an inbound rule: PostgreSQL (`5432`) from the Lightsail peered range (the region's Lightsail VPC CIDR, or the instance's private IP). The box can now reach the private RDS endpoint over TLS.

Install Docker on the instance (SSH in via the Lightsail console):

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
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin git
sudo usermod -aG docker $USER
newgrp docker
```

## 5. Deploy the API

```bash
sudo mkdir -p /opt/completionist
sudo chown $USER:$USER /opt/completionist
git clone <YOUR_REPO_URL> /opt/completionist
cd /opt/completionist
cp .env.production.example .env.production
```

Edit `.env.production`:

- `JWT_SECRET`: a strong random value.
- `DATABASE_URL`: the RDS string from step 3 (`sslmode=require`).
- `BACKEND_BASE_URL` and `PUBLIC_BASE_URL`: `https://api.<your-domain>` (finalized after Caddy, step 6).
- `FRONTEND_BASE_URL` and `ALLOWED_ORIGINS`: your Amplify URL (finalized after step 7).
- `STORAGE_DRIVER=s3` and the four `AWS_*` values from step 2.

Bring up the API (only `docker-compose.prod.yml`; there is no in-stack database):

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml up -d --build
```

Building on the box keeps it simple. If you prefer immutable images, push `Dockerfile.api` to ECR and point `api.image` at the ECR URI instead of building:

```bash
aws ecr create-repository --repository-name completionist-api --region <region>
docker buildx build --platform linux/amd64 -t <acct>.dkr.ecr.<region>.amazonaws.com/completionist-api:<tag> -f Dockerfile.api . --push
```

The API applies any pending migrations against RDS on startup.

## 6. HTTPS for the API with Caddy

The frontend gets TLS from Amplify. The API needs its own certificate. Put Caddy in front of it with an override `docker-compose.caddy.yml`:

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

volumes:
  caddy_data:
  caddy_config:
```

`Caddyfile`:

```
api.<your-domain> {
    reverse_proxy api:8080
}
```

Point `api.<your-domain>` at the Lightsail static IP (an A record in Route 53 or any DNS provider). Then:

```bash
docker compose --env-file .env.production \
  -f docker-compose.prod.yml \
  -f docker-compose.caddy.yml \
  up -d
```

Set `BACKEND_BASE_URL=https://api.<your-domain>` and `PUBLIC_BASE_URL=https://api.<your-domain>` in `.env.production`, then re-up. Once Caddy works, close `8080` in the Lightsail firewall; only `80` and `443` stay open.

## 7. Frontend on AWS Amplify

1. Amplify console, **Create new app**, **Host web app**, connect your GitHub repo and branch.
2. This is a monorepo: set the app root to **`apps/web`** (Amplify monorepo build settings, `appRoot: apps/web`). Amplify detects Next.js and runs the install plus `next build`.
3. Build-time environment variables (the Next build reads these via [apps/web/next.config.ts](../apps/web/next.config.ts) and bakes them into `NEXT_PUBLIC_*`):
   - `BACKEND_BASE_URL=https://api.<your-domain>`
   - `BACKEND_PORT=443`
4. Add a custom domain (e.g. `app.<your-domain>`). Amplify provisions the certificate and serves over its CDN.
5. Back in `.env.production` on the Lightsail box, set `FRONTEND_BASE_URL` and `ALLOWED_ORIGINS` to the Amplify domain, then re-up the API so CORS allows the frontend.

Re-read the Next.js version note at the top of this guide before the first Amplify build.

## 8. Verify

```bash
# API (from your machine, after Caddy + DNS):
curl -I https://api.<your-domain>/swagger/index.html

# Web: open https://app.<your-domain> in a browser.
# Confirm the network tab shows calls to https://api.<your-domain>/api/... ,
# and that login works (this exercises CORS end to end).
```

Confirm the schema landed in RDS (from the box, with psql pointed at DATABASE_URL):

```bash
psql "$DATABASE_URL" -c "\dt"
```

You should see ~25 tables (`users`, `posts`, `rooms`, `room_messages`, `direct_messages`, etc.).

## 9. Update and redeploy

- **API**: on the box, `git pull`, then `docker compose --env-file .env.production -f docker-compose.prod.yml -f docker-compose.caddy.yml up -d --build`. Migrations auto-apply. (Or push a new ECR tag and `pull` instead of `--build`.)
- **Web**: push to the connected branch. Amplify builds and deploys automatically, with atomic rollouts and one-click rollback.

---

## FAQ

### Why RDS instead of self-hosting Postgres on the box?

Managed automated backups, point-in-time recovery, and snapshots, with the database decoupled from the instance lifecycle. That removes the old `pg_dump`-to-S3 cron entirely. `db.t4g.micro` is the low-cost tier and is plenty for this workload; scale it up in place if you outgrow it.

### Why Amplify for the web instead of a container on Lightsail?

Managed builds from git, a global CDN, automatic TLS, and atomic deploys with rollback, none of which you operate by hand. The one caveat is the Next.js version note at the top; verify Amplify's Next 16 support before committing to it.

### Why is the API on Lightsail and not Amplify or App Runner?

The API is a long-lived Go process with a WebSocket hub, which wants a persistent server, not a frontend host or a request-scoped function. Lightsail is the cheapest predictable box for that. ECS/Fargate or App Runner are the next step if you outgrow a single instance.
