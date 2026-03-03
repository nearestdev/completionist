# Port Configuration

This project supports flexible port and host configuration for both backend and frontend services through a single `.env` file in the root directory.

The repository is organized as a small monorepo:
- `apps/api/` contains the Go backend module
- `apps/web/` contains the Next.js frontend

## How It Works

The configuration is separated into **base URLs** and **ports** for maximum flexibility:

**Backend:**
- `BACKEND_BASE_URL` - Host URL for the backend (default: `http://localhost`)
- `BACKEND_PORT` - Port for the Go API server (default: `8080`)
- Combined in code: `http://localhost:8080`

**Frontend:**
- `FRONTEND_BASE_URL` - Host URL for the frontend (default: `http://localhost`)
- `PORT` - Port for the Next.js frontend (default: `3000`)
- Combined in code: `http://localhost:3000`

**CORS:**
- `ALLOWED_ORIGINS` - Allowed CORS origins (default: auto-generated from `FRONTEND_BASE_URL:PORT`)

This separation allows you to:
- Change ports without touching URLs
- Use custom domains in production
- Keep configuration clean with just **2 port variables** (`BACKEND_PORT` and `PORT`)

## Environment Variables

### Root `.env` File

Both apps read from the root `.env` file, so you only need one configuration file:

```env
DATABASE_URL=postgres://postgres:testPassword@localhost:5434/completionist_db?sslmode=disable

BACKEND_BASE_URL=http://localhost
BACKEND_PORT=5000

FRONTEND_BASE_URL=http://localhost
PORT=4000

JWT_SECRET=your-jwt-secret
...
```

## Usage Examples

### Default Configuration

Simply omit all the URL/port variables - they'll use defaults:

```env
DATABASE_URL=postgres://postgres:testPassword@localhost:5434/completionist_db?sslmode=disable
JWT_SECRET=VgYsBfFni8y/TBtcjeby8BFnFq7kOwLj5NseAq21qZA=
...
```

Auto-generated:
- Backend: `http://localhost:8080`
- Frontend: `http://localhost:3000`
- Frontend API URL: `http://localhost:8080/api`
- CORS: `http://localhost:3000`

### Custom Ports (same host)

```env
DATABASE_URL=postgres://postgres:testPassword@localhost:5434/completionist_db?sslmode=disable

BACKEND_PORT=5000
PORT=4000

JWT_SECRET=VgYsBfFni8y/TBtcjeby8BFnFq7kOwLj5NseAq21qZA=
...
```

Auto-generated:
- Backend: `http://localhost:5000`
- Frontend: `http://localhost:4000`
- Frontend API URL: `http://localhost:5000/api`
- CORS: `http://localhost:4000`

### Production with Custom Domains

```env
BACKEND_BASE_URL=https://api.mydomain.com
BACKEND_PORT=443

FRONTEND_BASE_URL=https://app.mydomain.com
FRONTEND_PORT=443

JWT_SECRET=production-secret
...
```

Auto-generated:
- Backend: `https://api.mydomain.com:443`
- Frontend: `https://app.mydomain.com:443`
- Frontend API URL: `https://api.mydomain.com:443/api`
- CORS: `https://app.mydomain.com:443`

### Manual Override (Advanced)

You can still manually override final URLs:

```env
BACKEND_BASE_URL=http://localhost
BACKEND_PORT=5000
PUBLIC_BASE_URL=http://custom-backend-url:5000

FRONTEND_BASE_URL=http://localhost
FRONTEND_PORT=4000
ALLOWED_ORIGINS=http://localhost:4000,https://production.com
```

Manual values take precedence over auto-generated ones.

## Running the Services

### Backend
```bash
go -C apps/api run ./cmd/api
```

Starts on `BACKEND_BASE_URL:BACKEND_PORT`.

### Frontend
```bash
cd apps/web
bun run dev
```

Next.js automatically:
1. Reads `BACKEND_BASE_URL` and `BACKEND_PORT` from root `.env`
2. Builds API URL as `${BACKEND_BASE_URL}:${BACKEND_PORT}/api`
3. Runs dev server on `PORT` from `.env`
4. Serves the frontend from `FRONTEND_BASE_URL:PORT`

## Clean .env Example

Minimal configuration:

```env
DATABASE_URL=postgres://postgres:testPassword@localhost:5434/completionist_db?sslmode=disable

BACKEND_PORT=5000
PORT=4000

JWT_SECRET=VgYsBfFni8y/TBtcjeby8BFnFq7kOwLj5NseAq21qZA=

TMDB_API_KEY=7719e3fb899cd7b9a4cbf48c13109887

STEAM_WEB_API_KEY=B462958B98EF7F646215F777C6F7A6FA
STEAM_CALLBACK_REDIRECT_PATH=/settings
RAWG_API_KEY=a64777d5e95948cdb9216835d6b26973

GOOGLE_BOOKS_API_KEY=AIzaSyCiWu1Ai3s2AL2FwVYMFr3j65mXgjS5qnA

LASTFM_API_KEY=e5117a89c39743365978b0b36174edc5
LASTFM_API_SECRET=1bfe24b316d924a154ad9958cd54da5d
```

**That's it!** No need for `SERVER_ADDR`, `PUBLIC_BASE_URL`, `ALLOWED_ORIGINS`, or `FRONTEND_BASE_URL` - they're all auto-generated!
