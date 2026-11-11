# Stage 1 — verification checklist

**Initial status:** Not run on the user's Windows computer. Do not mark checks as passed until the user supplies results.

- [ ] `D:\Projects\eneba-marketplace-new` did not exist before extraction; no existing files overwritten.
- [ ] Required host ports available (or configured alternate ports).
- [ ] Docker Compose configuration validates.
- [ ] Docker Compose services launch/build successfully.
- [ ] `docker compose ps` shows the expected services running, database healthy.
- [ ] `GET http://localhost:8080/api/health` returns HTTP 200 with `status: ok` and `database: connected`.
- [ ] `GET http://localhost:5173` returns HTTP 200 and the browser displays the Stage 1 page.
- [ ] `pg_isready` reports the PostgreSQL service accepts connections.

If ports were changed in `.env`, update the frontend/backend URLs accordingly.

## Stage completion boundary

Stage 1 intentionally does **not** implement authentication, catalog, search, checkout, real payment, or digital keys. There is no full application yet.
