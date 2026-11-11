# Architecture and trade-offs

**Goal:** a demonstrable end-to-end buyer journey within a 7-day student assignment. UI is visually inspired by Eneba with a text-only Eneba-inspired demo wordmark, an independent demo label and explicit non-affiliation notice; public CDN artwork is referenced for titles, with the reuse/rights limitations documented in `docs/ARTWORK.md`. No official Eneba logo asset is copied.

- **Go modular monolith** instead of microservices: isolated packages for catalog, homepage, auth, orders; independently routed backend-managed homepage blocks; smaller operational footprint.
- **REST** instead of GraphQL: straightforward contracts for a limited set of pages.
- **SQL with the pgx driver**, not an ORM: SQL queries and atomic checkout transaction are inspectable and explainable. PostgreSQL's `pg_trgm` supports partial and approximate name matches; popularity falls back to seed ID when no purchases exist.
- **Database-backed HttpOnly cookies** instead of JWT in localStorage: revocable sessions, no client access to session token; server stores only token hashes. Passwords use bcrypt.
- **Cart persisted in browser localStorage**, authenticated order persisted on server. Browser price is display-only and never trusted for payment. A cart is not synced between browsers.
- **Checkout**: explicitly simulated approval/decline; offer locks serialize competing buyers, digital keys are allocated inside the same SQL transaction, an order is unique per user/idempotency key, and fingerprint stops reuse with a different payload. Quantity expands into separate order_items, one key each, matching the existing schema without destructive changes. Keys are fictional and cannot be redeemed.
- **Forward-only SQL migration** avoids deleting/reseeding Stage 2 rows. No destructive rollbacks or environment secret overwrites.
- **No seller administration**, real payments, password reset, email service, large catalogs or production-grade rate limiting: intentionally out of scope. Demo assets are CSS artwork to avoid relying on third-party game covers.
- **Production**: same-origin Nginx frontend/API and private PostgreSQL behind external HTTPS; it is a deployment recipe rather than a verified public service.
- **Lockfile limitation**: offline authoring environment cannot resolve transitive npm/go dependencies. `prepare-lockfiles.ps1` is a required one-time step on the networked Windows Docker environment before building; commit generated locks for reproducible installs.
