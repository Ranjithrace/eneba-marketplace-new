# Database

PostgreSQL stores marketplace content, accounts, sessions, offers, orders and demonstration keys. Prices use integer EUR cents.

For a **new empty database**, apply the following files in order:

1. `001_schema.sql`: relational schema and `pg_trgm` search extension.
2. `002_seed.sql`: 20 products, 3 sellers, 60 offers, 3 homepage sections and 120 demonstration keys.
3. `003_marketplace_migration.sql`: session and checkout/idempotency changes.
4. `004_artwork_metadata.sql`: product cover-image URL metadata.

These are initial seed counts. Successful checkout consumes available keys, so current available stock will differ.

See the root README for PowerShell/Docker commands and [Render + Neon](../docs/RENDER_FREE_DEPLOYMENT.md) for hosted initialization. Stop on every SQL error. Never point test commands at the hosted database.

## Existing data

Do not repeat the schema/seed process on a database with accounts or orders. Back up data, inspect migration history and apply only the required forward migration. Preserve the Docker volume and do not run `docker compose down -v` unless deliberately deleting that database.

The root `apply-migration.ps1` and `apply-artwork-migration.ps1` scripts were developed for earlier Windows upgrades; they are not needed for a fresh clone. Review their target and backup behavior before using them on any existing installation.

Checkout uses transactions and row locking to protect available stock. Orders store per-user idempotency keys and request hashes. Allocated keys are associated with order items; order queries are scoped to the authenticated user.

Use only non-redeemable demonstration keys. Keep database connection strings and backups outside Git.
