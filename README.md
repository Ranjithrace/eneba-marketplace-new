# GameKey Demo — Eneba internship marketplace

An independent, Eneba-inspired digital game marketplace by **Ranjith Rajagopalan**, built with React, TypeScript, Go and PostgreSQL. Payments, sellers, prices and delivered keys are fictional. This is not an official Eneba website.

- **Live application:** https://internship-game-marketplace.onrender.com
- **Submission source:** [deploy/render-free](https://github.com/Ranjithrace/eneba-marketplace-new/tree/deploy/render-free)
- **Health endpoint:** https://internship-game-marketplace.onrender.com/api/health
- **Deployment guide:** [Render + Neon](docs/RENDER_FREE_DEPLOYMENT.md)

The public service uses free hosting. The first visit after inactivity may take time to start. Availability is subject to provider policies and limits.

## Try the application

Browse the catalogue, search for a game (including a partial title or a typo such as `minecrft`), filter results and compare seller offers. Register a demonstration account to add an offer to your cart and complete a simulated checkout. Select an approved or declined payment outcome; no real payment details are needed. Successful orders deliver non-redeemable demonstration keys and appear in the account's order history.

## Features

- Database-managed homepage banner, best sellers and deals, loaded independently.
- Catalogue search, suggestions, platform/region/price filters, sorting and pagination.
- Product pages with multiple seller offers and available stock.
- Registration, login, logout, HttpOnly session cookies and private orders.
- Browser-persisted cart with quantity controls.
- Simulated checkout with server-side price/stock validation, transactional key allocation and per-user idempotency.
- Product artwork URLs returned by the API, with fallback cards when artwork is unavailable.

## Architecture

React provides the interface. Go handles HTTP requests, authentication and marketplace logic. PostgreSQL stores products, offers, homepage content, accounts, sessions, orders and demo keys.

Locally, Docker Compose runs React/Vite, Go and PostgreSQL as three services. On Render, [Dockerfile.render](Dockerfile.render) builds React and Go into one image: Go serves both the compiled frontend and `/api` routes from the same HTTPS origin. Neon hosts PostgreSQL. See [architecture](docs/architecture.md).

## Run locally from a fresh clone

Requires Git and Docker Desktop with Docker Compose. Docker downloads the Go and Node build environments; standalone Go and Node installations are unnecessary.

These commands are for a **new folder and empty database**. Choose another folder/project name and unused ports if another copy is already running.

```powershell
git clone --branch deploy/render-free https://github.com/Ranjithrace/eneba-marketplace-new.git gamekey-demo
cd gamekey-demo
Copy-Item .env.example .env
```

The example uses frontend port `5173`, API port `8081` and database port `55432`. Adjust these in `.env` if occupied. The example credentials are for local development only.

Run the following block in PowerShell. It uses a separate Compose project name, waits for PostgreSQL and stops on command failure:

```powershell
docker compose -p gamekey-demo config --quiet
if ($LASTEXITCODE -ne 0) { throw 'Invalid Compose configuration.' }

docker compose -p gamekey-demo up -d --wait db
if ($LASTEXITCODE -ne 0) { throw 'Database startup failed.' }

$sqlFiles = @(
    'database/001_schema.sql',
    'database/002_seed.sql',
    'database/003_marketplace_migration.sql',
    'database/004_artwork_metadata.sql'
)
foreach ($file in $sqlFiles) {
    Get-Content -LiteralPath $file -Raw |
        docker compose -p gamekey-demo exec -T db sh -c 'psql -X -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"'
    if ($LASTEXITCODE -ne 0) { throw "Database initialization failed: $file" }
}

docker compose -p gamekey-demo up --build -d
if ($LASTEXITCODE -ne 0) { throw 'Application build/startup failed.' }

docker compose -p gamekey-demo ps
```

With the default example ports, open http://localhost:5173 and http://localhost:8081/api/health.

To stop this copy while retaining its database:

```powershell
docker compose -p gamekey-demo down
```

Do not rerun initialization scripts against a populated database or use `down -v` if you want to retain accounts, orders and keys. See [database notes](database/README.md).

## Tests and verification

Frontend checks after building the local image:

```powershell
docker compose -p gamekey-demo run --rm --no-deps frontend npm run format:check
docker compose -p gamekey-demo run --rm --no-deps frontend npm test
docker compose -p gamekey-demo run --rm --no-deps frontend npm run build
```

Check each command's exit code. Docker image builds also run frontend tests/build and Go tests/build. Database integration tests require an isolated `TEST_DATABASE_URL`; ordinary Docker builds do not establish that these integration tests ran. Never point integration tests at the public or working local database.

[GitHub Actions](.github/workflows/ci.yml) configures PostgreSQL, applies migrations 001–004, runs Go race tests/vet/build, and runs frontend formatting/tests/build. [Acceptance checklist](docs/ACCEPTANCE.md) describes additional browser scenarios.

### Recorded status — 12 November 2026

| Evidence | Status |
| --- | --- |
| Fresh clone, Docker builds and local production preview | Reported passed during the Windows development session |
| Public HTTPS registration/login, session persistence, approved/declined checkout, order persistence, stock updates, order privacy and idempotency | Reported passed by the author during deployed-site testing |
| Visit after inactivity | Reported passed by the author; no precise cold-start duration recorded here |
| Earlier `main` CI, commit `edb89f7` | [Passed](https://github.com/Ranjithrace/eneba-marketplace-new/actions/runs/43073825275) |
| Deployment commit `075e202` CI | [Run #4](https://github.com/Ranjithrace/eneba-marketplace-new/actions/runs/43096494005) remained queued with no jobs when checked; result unverified |

Run #4's UI cancellation failed, and authenticated force-cancellation returned HTTP 409 despite the displayed queued status. The exact cause has not been established. Earlier CI success does not prove that deployment-branch CI passed. Consult Actions for later runs; this table records the evidence available when these docs were updated.

The documentation update was based on repository inspection and the author's supplied test results; it was not a new execution of the application test suite.

Database-driven homepage sections; catalogue search with typo support, filters and sorting; multiple seller offers; registration/login; persistent cart; simulated approved/declined checkout; non-redeemable key delivery; and private order history.

| Path | Purpose |
| --- | --- |
| `frontend/src/app/` | Shared page shell and routing |
| `frontend/src/features/` | Home, catalogue, product, authentication, cart and order interfaces |
| `frontend/src/shared/` | API client, types and shared utilities |
| `backend/cmd/api/` | Application startup |
| `backend/internal/` | HTTP routing, feature logic, database access and static frontend serving |
| `database/` | Schema, seed and migrations 001–004 |
| `compose.yaml` | Local development services |
| `Dockerfile.render`, `render.yaml` | Render deployment configuration |
| `compose.production.yaml`, `deployment/` | Alternative Nginx-based Compose deployment |
| `docs/` | Architecture, deployment, tests and development history |

The stage checklists, visual-update/review-fix documents and `apply-*.ps1` helpers record earlier upgrade steps. They are not prerequisites for a fresh clone.

## Security, artwork and limits

Passwords use bcrypt; session tokens are hashed in storage and sent through HttpOnly cookies. The public deployment is configured for Secure cookies. Order access is restricted to the authenticated owner. Checkout validates server-side prices and stock and allocates keys within a transaction.

This educational demo has no real payment gateway, seller dashboard, email password reset or production-grade anti-abuse system. Demo key stock is finite and decreases with successful purchases.

Third-party cover images remain copyrighted by their owners; reuse permission has not been verified. Source references and fallback behavior are documented in [ARTWORK.md](docs/ARTWORK.md).

Secrets belong in local ignored `.env` files or hosting environment settings. Never commit credentials, session data or database dumps.

## Submission materials

Provide the GitHub submission branch, live URL and **full AI conversation history containing both prompts and responses**, as Markdown or accessible chat links. Include every relevant development conversation; a summary is not a replacement.

The full conversation history must be supplied separately with the submission unless added to this repository. This README does not claim that a transcript has been attached or that the submission has been sent.
