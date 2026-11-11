# GameKey Demo — Eneba internship marketplace

An independent Eneba-inspired marketplace by **Ranjith Rajagopalan**, built with React, TypeScript, Go and PostgreSQL.

## Review the submission

- **Live application:** https://internship-game-marketplace.onrender.com
- **Current submission source and full README:** [deploy/render-free](https://github.com/Ranjithrace/eneba-marketplace-new/tree/deploy/render-free)
- **Deployment guide:** [Render + Neon](https://github.com/Ranjithrace/eneba-marketplace-new/blob/deploy/render-free/docs/RENDER_FREE_DEPLOYMENT.md)

You are viewing `main`, which contains the earlier local-development implementation. The public deployment configuration and updated submission documentation are on **`deploy/render-free`**. Use that branch when reviewing or cloning the submission:

```powershell
git clone --branch deploy/render-free https://github.com/Ranjithrace/eneba-marketplace-new.git gamekey-demo
```

Follow that branch's README for Docker setup, database initialization, tests and architecture.

## What the demo includes

Database-driven homepage sections; catalogue search with typo support, filters and sorting; multiple seller offers; registration/login; persistent cart; simulated approved/declined checkout; non-redeemable key delivery; and private order history.

No real payments or redeemable keys are involved. This is not an official Eneba site. Third-party artwork reuse permissions remain unverified; see the submission branch's artwork documentation.

## Verification and hosting notes

The author reported successful local builds and public HTTPS account, checkout, privacy and cold-visit tests. Free hosting may delay the first visit after inactivity.

The earlier `main` CI run for `edb89f7` passed. At the documentation update on 12 November 2026, [deployment run #4](https://github.com/Ranjithrace/eneba-marketplace-new/actions/runs/43096494005) for `075e202` remained queued with no jobs; UI cancellation failed and force-cancellation returned HTTP 409. This is not a completed CI result. Check [Actions](https://github.com/Ranjithrace/eneba-marketplace-new/actions) for newer runs.

Full AI prompts/responses or accessible conversation links must accompany the assignment submission separately; this README is not that history.
