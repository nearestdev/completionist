# Deploy Guide: VPS + Supabase Postgres

This project can be deployed with:

- Supabase for PostgreSQL
- One VPS for the web and API containers
- Docker Compose for orchestration

The backend automatically runs SQL migrations at startup, so you do not need a separate migration step.

## 1. Create Supabase project

1. Create a new project in Supabase.
2. Open `Project Settings -> Database -> Connection string`.
3. Copy the **Session pooler** connection string (IPv4 friendly).
4. Ensure it includes `sslmode=require`.

Example format:

```text
postgres://postgres.<PROJECT_REF>:<PASSWORD>@aws-0-<REGION>.pooler.supabase.com:5432/postgres?sslmode=require
```

## 2. Prepare VPS

Use Ubuntu 24.04 (or similar). Then install Docker + Compose plugin:

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

## 3. Clone repo on VPS

```bash
git clone <YOUR_REPO_URL> completionist
cd completionist
```

## 4. Create production env file

```bash
cp .env.production.example .env.production
```

Edit `.env.production` and fill at least:

- `DATABASE_URL` (Supabase session pooler, ssl required)
- `JWT_SECRET`
- `BACKEND_BASE_URL`
- `FRONTEND_BASE_URL`
- `ALLOWED_ORIGINS`
- `WEB_BACKEND_BASE_URL`
- `WEB_BACKEND_PORT`

Notes:

- `WEB_BACKEND_*` is used at **web build time**.
- If backend URL changes later, rebuild/redeploy web.
- Keep `SERVER_ADDR=:8080` and `PORT=4000` unless you intentionally change ports.

## 5. Build and run

Important: pass `--env-file .env.production` so Compose can resolve build args from that file.

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml build
docker compose --env-file .env.production -f docker-compose.prod.yml up -d
```

Check status/logs:

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml ps
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f api
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f web
```

## 6. Verify deployment

From VPS:

```bash
curl -I http://localhost:4000
curl -I http://localhost:8080/swagger/index.html
```

From browser:

- Frontend: `http://<VPS_IP>:4000`
- API docs: `http://<VPS_IP>:8080/swagger/index.html`

## 7. Open firewall ports

Allow inbound:

- `4000/tcp` (frontend)
- `8080/tcp` (API)

If you add a reverse proxy for HTTPS later, you typically expose only `80/443`.

## 8. Update deployment

```bash
git pull
docker compose --env-file .env.production -f docker-compose.prod.yml build
docker compose --env-file .env.production -f docker-compose.prod.yml up -d
```

## 9. Optional: HTTPS + domain

Recommended next step is Caddy or Nginx in front of these two services:

- `app.yourdomain.com -> web:4000`
- `api.yourdomain.com -> api:8080`

Then update:

- `FRONTEND_BASE_URL=https://app.yourdomain.com`
- `BACKEND_BASE_URL=https://api.yourdomain.com`
- `ALLOWED_ORIGINS=https://app.yourdomain.com`
- `WEB_BACKEND_BASE_URL=https://api.yourdomain.com`

And redeploy with build + up.
