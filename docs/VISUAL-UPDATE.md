# Visual refresh — change log and local acceptance

## Source baseline

This ZIP contains a full project snapshot based on `eneba-marketplace-complete.zip` plus the **applied** `eneba-marketplace-search-fix.zip` source files. The user's latest confirmed Windows state has working catalog search with the SQL alias fix and passing frontend/Go builds. We did **not** remotely inspect local files or run a Windows deployment.

## Visual-only files

- `frontend/src/app/App.tsx`: functional two-row header, responsive search and navigation, independent Eneba-inspired demo identity.
- `frontend/src/shared/Icon.tsx`: self-contained outline icons (no extra package).
- `frontend/src/shared/coverArtwork.ts`: exact game-to-image-source mapping and explicit fallback.
- `frontend/src/features/catalog/{ProductCard,StorePage}.tsx`: portrait art, compact cards, condensed catalog heading.
- `frontend/src/features/home/HomePage.tsx`: anchor navigation into backend-driven collections; no fake achievements or ratings.
- `frontend/src/features/cart/CartPage.tsx`: correct product image when a known image is available.
- `frontend/src/features/product/ProductPage.tsx`: small copy clarification; API logic unchanged.
- `frontend/src/styles.css`: saturated purple design system and responsive layouts for catalog and all other pages.
- `frontend/tests/coverArtwork.test.mjs`: title-to-cover mapping and fallback regression checks.

Only the listed frontend files and documentation are copied by the upgrade script. The ZIP also includes unchanged backend/database references for a complete source snapshot; **do not drag-and-drop overwrite them onto the live installation**.

## Windows application procedure

1. Extract ZIP into `D:\Projects\eneba-visual-preview` **not** into the live project.
2. Run `powershell.exe -NoProfile -ExecutionPolicy Bypass -File "D:\Projects\eneba-visual-preview\eneba-marketplace-new\apply-visual-update.ps1"` to preview only.
3. Review the file list and any local-edit warnings. If expected, repeat with `-Apply`. The script backs up every overwritten file into a timestamped `D:\Projects\eneba-marketplace-visual-backup-*` directory and verifies copied hashes.
4. From `D:\Projects\eneba-marketplace-new`, run a one-off formatter using the existing frontend Docker image, bind-mounting **only** `frontend/src` so the formatted source is saved on Windows:

   ```powershell
   $src = (Resolve-Path '.\frontend\src').Path
   docker compose run --rm --no-deps -v "${src}:/app/src" frontend npm run format
   if ($LASTEXITCODE -ne 0) { throw 'Frontend formatting failed.' }
   docker compose up --build -d frontend
   docker compose exec -T frontend npm run format:check
   docker compose exec -T frontend npm run test
   docker compose exec -T frontend npm run build
   ```

   Do not proceed if an earlier command fails. This does not touch backend or DB containers, `.env`, or existing dependency lockfiles.

## Local acceptance (after applying patch)

- [ ] `docker compose up --build -d frontend` succeeds; backend and db unchanged.
- [ ] Desktop 1600+ px: saturated purple page, darker top navigation, 238px filter sidebar, **four** compact product columns in store.
- [ ] Tablet 900px: content fits, fewer columns, functional nav scrolling.
- [ ] Mobile 375/430px: no horizontal page overflow; filters and cards stack/flow sensibly.
- [ ] Game cover thumbnails load from Steam for matched titles. Minecraft or unavailable titles show a plain text fallback, never unrelated art.
- [ ] Search suggestions, filters, sorting, pagination, product offers, sign-in/out, cart quantities, simulated checkout and personal orders still function.
- [ ] Backend `8081`, frontend `5173`, PostgreSQL `55432` are preserved; `.env`, database, credentials and local dependency lockfiles unchanged.
- [ ] No fake review counts, discounts, stock badges or unsupported navigation controls.

External art is not licensed by this demo; see `docs/ARTWORK.md` before publishing.
