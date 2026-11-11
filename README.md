# Eneba-inspired internship marketplace (independent demo)

React + TypeScript / Go / PostgreSQL. A small buyer-facing marketplace with a purple/yellow interface inspired by Eneba. This is **not an official Eneba site**. Prices, sellers, checkout and digital keys are simulated; no real money or real redeemable keys.

## Features

- Backend-managed banner, best sellers and deals; each block loads independently, so one failed request does not hide the others.
- Catalog with search (including typos), platform/region/price filters, sorting and pagination.
- Product pages with multiple seller offers and live available stock.
- Registration, login, logout, HttpOnly session cookies and private order history.
- Persistent browser cart, quantity controls, simulated approved/declined payments, transactional key allocation and idempotent checkout.
- Title-specific cover art references stored in `products.cover_image_url` and returned by the API. These references **do not grant reuse rights** (see `docs/ARTWORK.md`). Neutral fallback if unavailable.

## Requirements

Windows PowerShell and Docker Desktop with Compose; network access on first dependency resolution. No standalone Go/Node installation is needed for Docker workflow.

**Development port settings (retain current `.env`):** frontend `5173`, API `8081`, PostgreSQL `55432`. Docker services use private internal ports independent of these values.

## Existing Windows installation — safe review-fix update

You already have a working installation under `D:\Projects\eneba-marketplace-new`. Do **not** extract over it, and do **not** rerun `001_schema.sql` or `002_seed.sql` on your existing database.

1. Extract the new ZIP to `D:\Projects\review-fix-preview`, leaving a top folder `eneba-marketplace-new`.
2. Run `powershell.exe -NoProfile -ExecutionPolicy Bypass -File 'D:\Projects\review-fix-preview\eneba-marketplace-new\apply-review-fixes.ps1'` (read-only preview).
3. Review differences/backup scope; then run same command with `-Apply`. The installer does **not** touch `.env`, database, orders, keys, backend, or dependency lockfiles.
4. Run `powershell.exe -NoProfile -ExecutionPolicy Bypass -File '.\apply-artwork-migration.ps1'` from your LIVE project for a **read-only preview**. Only after checking that output, rerun the same command with `-Apply`: this creates a PostgreSQL backup and fills **only empty product image URLs** from `database/004_artwork_metadata.sql`.
5. Run `docker compose up --build -d frontend` from the live folder (and inspect `docker compose ps`); frontend-only rebuild. Backend code does not need to change.
6. Check the homepage independently, store covers, and complete the acceptance scenarios in `docs/ACCEPTANCE.md`.

Your **existing complete** `frontend/package-lock.json` and `backend/go.sum` should be preserved. The previous ZIP contained incomplete *template* copies. This offline workspace could not recover your Windows-generated files. See `docs/LOCKFILES.md` to validate them, ensure they are committed to GitHub, and test a fresh clone. The ZIP's bundled placeholders are **not valid final lockfiles** and must not overwrite verified local files.

## Fresh GitHub clone — not yet verified

**Before pushing** the repository, confirm the real `package-lock.json` and `go.sum` are included (they were generated during the successful Windows Docker run). CI deliberately uses strict `npm ci` and `go mod verify`; it no longer silently repairs invalid files. See `docs/LOCKFILES.md`.

With a fresh clone containing complete locks, in PowerShell:

```powershell
Copy-Item .env.example .env
# Edit .env as needed, ensuring the ports are free before first startup.
docker compose config --quiet
docker compose up --build -d db
Get-Content '.\database\001_schema.sql' -Raw | docker compose exec -T db psql -X -v ON_ERROR_STOP=1 -U eneba -d eneba_marketplace
Get-Content '.\database\002_seed.sql' -Raw | docker compose exec -T db psql -X -v ON_ERROR_STOP=1 -U eneba -d eneba_marketplace
Get-Content '.\database\003_marketplace_migration.sql' -Raw | docker compose exec -T db psql -X -v ON_ERROR_STOP=1 -U eneba -d eneba_marketplace
Get-Content '.\database\004_artwork_metadata.sql' -Raw | docker compose exec -T db psql -X -v ON_ERROR_STOP=1 -U eneba -d eneba_marketplace
docker compose up --build -d
```

For schema/seed commands, check each exit code before continuing; the `-v ON_ERROR_STOP=1` option stops SQL on errors. The above is **only for a brand-new empty database**. Never use `docker compose down -v` on a project whose data you wish to retain.

Open `http://localhost:5173`, API health `http://localhost:8081/api/health`.

## Reproducible checks

```powershell
# Run inside frontend Docker image after correct lockfiles are present:
docker compose run --rm --no-deps frontend npm run test
# Docker builds also execute frontend tests/typecheck/Vite build and Go tests/build:
docker compose up --build -d
```

For Go/PostgreSQL integration tests, use a **separate disposable test database** and `TEST_DATABASE_URL`; never point tests at your live order/key data. The complete local acceptance checklist is in `docs/ACCEPTANCE.md`.

## Source layout

- `frontend/src/app/`: header, navigation and routing; `features/`: home, catalog, products, auth, cart, orders; `shared/`: API/types and artwork URL validation.
- `backend/cmd/api`: startup; `backend/internal/`: modular Go catalog, homepage, auth and order packages.
- `database/`: schema, seed and forward-only migrations (`003`, `004`).
- `.github/workflows/ci.yml`: strict automated checks requiring committed complete dependency locks.
- `compose.production.yaml`, `frontend/Dockerfile.production`, `deployment/nginx.conf`: hosting configuration (not yet deployed); see `docs/DEPLOYMENT.md`.

## Security and explicit scope limits

Passwords are bcrypt hashed; session tokens are stored as hashes and issued in HttpOnly cookies. Checkout validates current prices/stock and locks offer rows before assigning unique keys. Simulated approval and decline require no payment details. No seller dashboard, real payment gateway, email reset or production-grade anti-abuse system is included. Review public asset rights and apply appropriate security protections before deployment.

## Deliverables

The assignment requires a GitHub repository, public deployment URL and **full** AI chat transcript (both prompts and responses) or share link. Do not substitute a conversation summary. Record actual tests and unresolved limitations in the submission. **A successful local screenshot does not prove a fresh clone or public deployment works.**
