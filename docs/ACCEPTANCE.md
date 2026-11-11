# Acceptance scenarios and recorded evidence

The root [README](../README.md#recorded-status--12-november-2026) records the author's reported local/public test results and successful GitHub Actions checks on both `main` and `deploy/render-free`. The unchecked items below are a reusable test checklist, **not a claim that every scenario has been executed**. In particular, stock-exhaustion, all responsive widths and request-failure injection need their own evidence if claimed.

GitHub CI evidence (12 November 2026): [main run #10](https://github.com/Ranjithrace/eneba-marketplace-new/actions/runs/43117879412) passed; [deployment push run #7](https://github.com/Ranjithrace/eneba-marketplace-new/actions/runs/43117103688) and [pull-request run #8](https://github.com/Ranjithrace/eneba-marketplace-new/actions/runs/43117119652) passed. [Run #4](https://github.com/Ranjithrace/eneba-marketplace-new/actions/runs/43096494005) was an earlier stuck queued run with no jobs; it is historical troubleshooting, not a failed test or the current CI result.

Upgrade-script preview/backup checks apply only to historical upgrade installations, not a fresh Git clone.

## Automated / infrastructural
- [ ] Update script preview reports no unexpected edits and writes nothing.
- [ ] Update creates a dated **file backup** outside the working folder; existing `.env` remains unchanged.
- [ ] Database dump created successfully; `003_marketplace_migration.sql` and `004_artwork_metadata.sql` report COMMIT without errors.
- [ ] Full locally generated `package-lock.json` and `go.sum` exist and `npm ci`, Go tests, Node tests/builds pass.
- [ ] `docker compose up --build -d` starts frontend, backend, and database on ports 5173, 8081, 55432; `docker compose ps` shows `db` healthy.
- [ ] `GET /api/health` reports `database: connected`; website returns HTTP 200.
- [ ] Seed data still contains 20 products and 60 offers; existing keys preserved; migration tables `auth_sessions` and columns `orders.idempotency_key` exist.

## Browser
- [ ] Homepage banner, best sellers and deals load from API; use DevTools to block one request and verify the other two still load; edit a homepage_sections heading in DB temporarily and verify the page changes (restore it afterwards).
- [ ] Header search works across the site with partial name `mine` and typo `minecrft`.
- [ ] Store: PC/Europe filter, two price bounds, sorting, pagination, empty results.
- [ ] Product page lists several sellers with distinct prices and live stock.
- [ ] Register a new account; sign out; sign back in; navigation updates; unauthenticated `/checkout` and `/orders` require login.
- [ ] Add multiple offers to cart, increase/decrease quantities; refresh browser and verify cart remains.
- [ ] Approve a test payment; record order number and delivered `DEMO-ONLY-...` key. Cart clears; order history shows it.
- [ ] Repeat the same checkout POST with identical idempotency key; expect the same order, **not a second allocation**.
- [ ] Decline another test checkout; order status `failed`, no key code delivered.
- [ ] Purchase remaining keys of an offer, verify sold-out state and that extra checkout is rejected without overselling.
- [ ] Sign in as a second user; verify first user's order URL returns 404, history does not show their keys.
- [ ] Inspect widths 360, 390, 768 and 1440px, keyboard navigation and no horizontal scrolling; navigation, search, product cards, seller offers, cart and order history are usable.

## Release review
- [ ] Real dependency locks committed; CI green; no test failures or secrets in repository.
- [ ] Trusted HTTPS domain, secure-cookie setting, strong secret and persistent DB backup strategy.
- [ ] Public URL independently tested end to end; share GitHub repo URL and **complete** AI transcript/chat link.
