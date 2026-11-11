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

## Existing installations and data safety

The original Windows installation was developed at
`D:\Projects\eneba-marketplace-new`.

Preserve its `.env`, PostgreSQL volume, accounts, orders and
digital keys. Do not rerun initialization scripts against an
existing populated database.

Never run `docker compose down -v` on databases whose
data must be retained.

Use a separate Git clone, Docker Compose project name and
host ports when reproducing or testing this application.

## Fresh GitHub clone — locally verified

A clean GitHub clone was built and tested on 2026-11-11.

Verified:
- All four database scripts applied to a new empty database.
- Initial counts: 20 products, 60 offers, 120 demo keys
  and 3 homepage sections.
- Go backend tests and production build.
- React frontend tests and production build.
- Backend health and PostgreSQL connectivity.
- Catalogue and typo-tolerant search.
- Homepage, store and product route responses.
- Production Docker Compose startup and shutdown.
- A separate single-container React + Go deployment build
  and local runtime check.

The clean-clone test used separate development ports:
5174 for React, 8082 for Go and 55433 for PostgreSQL.

The local production preview was tested on port 8088.
The Render-style single-container preview was tested
on port 8090.

These checks do not establish public deployment,
public HTTPS authentication or remote availability.

For a fresh installation, copy `.env.example` to `.env`,
select unused ports, and initialize a new empty database
using the four SQL scripts in `database/` in order.
Check each script's exit code before continuing.

Detailed free hosting instructions:
`docs/RENDER_FREE_DEPLOYMENT.md`.
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
