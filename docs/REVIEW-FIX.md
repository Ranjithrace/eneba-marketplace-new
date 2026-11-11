# Review findings and resolutions — November 11, 2026

1. **Incomplete dependency locks**: confirmed in downloaded visual-refresh snapshot (`package-lock.json` only root entry, `go.sum` empty). User's Windows build logs prove locally generated complete locks previously existed. Installer preserves those; CI now requires them strictly. See `docs/LOCKFILES.md`. **Not yet verified in current live folder or from a fresh GitHub clone.**
2. **README stale**: replaced with explicit existing-install and fresh-clone paths, fresh database migrations 001–004, current ports, testing scope and deployment limitations.
3. **Homepage sections coupled**: replaced `Promise.all` with three independent requests and individual errors/loading indicators. A failure in one leaves successful sections visible. Two unit tests cover error isolation and cleanup. Pending browser test using DevTools request blocking.
4. **Image mapping in frontend**: removed slug → Steam-ID table; product art now reads the backend-supplied `cover_image_url` only. Additive SQL migration fills blank URLs for known demo games and retains existing data. Two unit tests check URL validation. Artwork is still not rights-cleared for public deployment.

**Unchanged:** saturated-purple styling, component and route layout, backend business code, existing account/order/key records, `.env` and ports. Installer and SQL migration are separate, backed up and preview-first. There is no claim of direct access to the live Windows filesystem.
