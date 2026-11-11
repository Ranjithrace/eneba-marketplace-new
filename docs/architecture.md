# Architecture — complete marketplace implementation

The frontend is React/TypeScript built by Vite. `src/app/App.tsx` provides header/search and simple path routing. Individual features live in `src/features/home`, `catalog`, `product`, `auth`, `cart`, and `orders`. Shared DTO types, API calls, money formatting and navigation helpers are in `src/shared`. The client cart is stored in browser localStorage; session tokens never appear in JS (HttpOnly cookies).

`backend/cmd/api/main.go` configures the database and starts HTTP. `internal/api/router.go` attaches feature handlers, while `catalog`, `homepage`, `auth` and `orders` own their logic. `database` and `httpjson` provide shared infrastructure. PostgreSQL holds seeded content, accounts, sessions, offers, orders, order items and allocated demonstration keys.

REST endpoints are described in the root README. Each homepage block has a separate Go endpoint; its content is returned to React. The product catalog uses a trigram-backed search plus filters and pagination. Checkout validates requests, locks offer rows in sorted order, allocates unused keys transactionally, and deduplicates identical idempotency keys per user. All order queries restrict by authenticated user ID.

`database/003_marketplace_migration.sql` only *adds* columns/indexes/session tables on top of Stage 2, preserving the verified data. No replacement of `.env` or PostgreSQL named volume occurs. Full lockfiles are generated on a networked Docker host using the included `prepare-lockfiles.ps1` script (offline authoring limitation).

Testing: Go handler/validation unit tests, PostgreSQL journey integration tests (require an isolated `TEST_DATABASE_URL`), Node cart, money, navigation and idempotency tests, static formatting and build checks, and a GitHub Actions workflow. Integration and browser verification remain required before public deployment.
