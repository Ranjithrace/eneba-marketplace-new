# Production hosting preparation (NOT deployed or verified)

The production Compose file runs three services: private PostgreSQL 18, Go API, and Nginx serving compiled React and forwarding `/api` to Go. It binds Nginx to `127.0.0.1:${PUBLIC_PORT:-8088}` so a **separate trusted HTTPS reverse proxy** must terminate TLS and route to that local port. The server must have Docker Compose, a persistent volume and a domain with HTTPS; this configuration is **not** an automatically provisioned free-tier deployment.

1. Copy `.env.production.example` to `.env` **only on a separate production host**, replace the sample password and domain. **Never** overwrite your local `D:\Projects\eneba-marketplace-new\.env`.
2. Commit the actual full `frontend/package-lock.json` and `backend/go.sum` from the verified Windows installation. CI must use `npm ci` and `go mod verify` without re-generating them.
3. Run `docker compose --env-file .env -f compose.production.yaml up --build -d db`.
4. For a new database only, apply `database/001_schema.sql`, then `database/002_seed.sql`, then `database/003_marketplace_migration.sql` via `docker compose -f compose.production.yaml exec -T db psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"`. For an existing Stage 2 database, back it up and run `003_marketplace_migration.sql` if needed and `004_artwork_metadata.sql` after backup.
5. Start web + API using `docker compose --env-file .env -f compose.production.yaml up --build -d`. Configure a TLS reverse proxy for `PUBLIC_ORIGIN` to the local port 8088. Production uses HttpOnly `Secure` cookies, so HTTPS is required.
6. Verify real public URLs, privacy and checkout. Configure backups, monitoring and appropriate abuse protection before public launch. This is an educational demo, not a production payment platform.

Never expose port 5432 publicly. Keep `.env` and database dumps out of Git. Do not run `down -v` against persistent data.
