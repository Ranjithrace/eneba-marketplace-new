# Architecture — complete marketplace implementation

The frontend is React/TypeScript built by Vite. `src/app/App.tsx` provides header/search and simple path routing. Individual features live in `src/features/home`, `catalog`, `product`, `auth`, `cart`, and `orders`. Shared DTO types, API calls, money formatting and navigation helpers are in `src/shared`. The client cart is stored in browser localStorage; session tokens never appear in JS (HttpOnly cookies).

`backend/cmd/api/main.go` configures the database and starts HTTP. `internal/api/router.go` attaches feature handlers, while `catalog`, `homepage`, `auth` and `orders` own their logic. `database` and `httpjson` provide shared infrastructure. PostgreSQL holds seeded content, accounts, sessions, offers, orders, order items and allocated demonstration keys.

HTTP routes are registered in `backend/internal/api/router.go`. Each homepage block has a separate Go endpoint; its content is returned to React. The product catalog uses a trigram-backed search plus filters and pagination. Checkout validates requests, locks offer rows in sorted order, allocates unused keys transactionally, and deduplicates identical idempotency keys per user. All order queries restrict by authenticated user ID.

`database/003_marketplace_migration.sql` only *adds* columns/indexes/session tables on top of Stage 2, preserving the verified data. No replacement of `.env` or PostgreSQL named volume occurs. Complete dependency lockfiles are committed; normal installation uses those files. See `docs/LOCKFILES.md`.

Testing: Go handler/validation unit tests, PostgreSQL journey integration tests (require an isolated `TEST_DATABASE_URL`), Node cart, money, navigation and idempotency tests, static formatting and build checks, and a GitHub Actions workflow. The author reported deployed browser checks passing. GitHub Actions checks passed on `deploy/render-free` (push run [#7](https://github.com/Ranjithrace/eneba-marketplace-new/actions/runs/43117103688) and PR run [#8](https://github.com/Ranjithrace/eneba-marketplace-new/actions/runs/43117119652)) and on `main` (run [#10](https://github.com/Ranjithrace/eneba-marketplace-new/actions/runs/43117879412)). The older run #4 remained stuck queued without jobs and is historical troubleshooting, not a failed test. See the root README for the recorded browser evidence.

## Public hosting

The live demo runs on Render with Neon PostgreSQL. `Dockerfile.render` compiles React and Go into one image. `backend/internal/web/site.go` serves the frontend alongside the API from one origin. Render terminates HTTPS; the application uses Secure cookies and Neon TLS. Local Docker Compose still runs three services. The Nginx production Compose configuration is an alternative, not the architecture of the submitted public URL.
