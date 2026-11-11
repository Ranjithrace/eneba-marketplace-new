# Database

- `001_schema.sql`: original Stage 2 relational schema and `pg_trgm` extension.
- `002_seed.sql`: original Stage 2 products (20), sellers (3), offers (60), homepage sections (3), demonstration keys (120).
- `003_marketplace_migration.sql`: non-destructive changes for Stage 2 to complete app: hashed session tokens (`auth_sessions`), unique checkout idempotency keys and request hashes on `orders`, schema migration marker.

**Verified existing Windows database:** preserve all records and the Docker volume. Do **not** rerun the initial schema/seed scripts as part of this upgrade. Run `apply-migration.ps1`, which creates a dated logical dump in `D:\Projects` and then applies migration `003` transactionally. The migration is safe to retry if the first attempt completed; it never deletes existing keys, users, offers, orders or products.

For a **fresh, separate database**, run 001, then 002, then 003 in that order. Do not use real redeemable game keys: seeded keys are `DEMO-ONLY-*`. Prices are integer EUR cents.

After the migration:

```powershell
'\dt' | docker compose exec -T db psql -X -U eneba -d eneba_marketplace
'SELECT COUNT(*) FROM products;' | docker compose exec -T db psql -X -U eneba -d eneba_marketplace
```

`orders.idempotency_key` has a unique `(user_id, idempotency_key)` index. Each `digital_keys.order_item_id` is unique to guarantee a delivered key belongs to at most one order item; checkout also acquires row locks to guard stock.
