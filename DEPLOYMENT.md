# VPS Deployment

This repository includes a Docker Compose setup for running Paskihub Backend with PostgreSQL and Redis on one VPS.

## Files

- `Dockerfile`: multi-stage Go build for the backend binary.
- `docker-compose.yml`: app, PostgreSQL, Redis, networks, healthchecks, and persistent volumes.
- `.env.docker.example`: production-oriented environment template.

## First Deploy

1. Copy the template:

   ```bash
   cp .env.docker.example .env
   ```

2. Edit `.env` and replace every `change-me-*` value.

   Keep these internal Docker host values unless you use managed services:

   ```env
   DB_HOST=postgres
   DB_PORT=5432
   REDIS_HOST=redis
   REDIS_PORT=6379
   ```

3. Start the stack:

   ```bash
   docker compose up -d --build
   ```

4. Check logs:

   ```bash
   docker compose logs -f app
   ```

The app runs on `APP_PORT` from `.env`. By default, that is `3010`.

## Data Persistence

Docker named volumes persist database, Redis, uploads, and app logs:

- `paskihub_postgres_data`
- `paskihub_redis_data`
- `paskihub_uploads`
- `paskihub_logs`

Do not remove these volumes unless you intentionally want to delete data.

## Network Exposure

The backend port is published publicly:

```yaml
${APP_PORT:-3010}:${APP_PORT:-3010}
```

PostgreSQL and Redis are bound to localhost only:

```yaml
127.0.0.1:5432:5432
127.0.0.1:6379:6379
```

For production, place Nginx, Caddy, or another reverse proxy in front of the app and terminate HTTPS there.

## Useful Commands

```bash
docker compose ps
docker compose logs -f app
docker compose restart app
docker compose pull
docker compose up -d --build
docker compose down
```

Avoid `docker compose down -v` in production because it removes persistent volumes.
