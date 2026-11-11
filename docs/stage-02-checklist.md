# Stage 2 verification checkpoint

**Not verified on the user's Windows machine until they provide actual output.**

1. Preview and back up Stage 1 before applying changes. Do not overwrite `.env`.
2. Run schema and seed SQL against the existing PostgreSQL volume, without dropping it.
3. Rebuild updated Go and React containers, leaving the original project on port 8080 untouched.
4. Verify the SQL count query expects 20 products, 3 sellers, 60 offers, and 120 fictional keys on the sample database.
5. Verify `http://localhost:8081/api/products` returns `count: 20` and a product array.
6. Verify `http://localhost:5173` shows a preview list with sample prices.
7. Check `docker compose ps`: all 3 Stage 2 services should be running; PostgreSQL should be healthy.

The schema is intentionally only a first iteration. Authentication, transactions, and key allocation are not implemented until later stages. The presence of tables/seed keys does NOT mean checkout works.
