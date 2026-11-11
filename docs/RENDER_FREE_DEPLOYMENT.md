# Render Free + Neon Free Deployment

Live application: https://internship-game-marketplace.onrender.com

Source branch: `deploy/render-free`. Application deployment commit recorded during submission preparation: `075e202`.

The author reported successful public HTTPS functional and cold-visit checks. These are recorded results, not a new independent browser test performed for this documentation update. Deployment-branch CI run #4 remains unresolved as described in the root README.

## Architecture

- One Render Free Docker web service.
- React and TypeScript are compiled during the Docker build.
- Go serves the React application and all /api endpoints.
- Neon Free provides PostgreSQL.
- Render provides a public HTTPS .onrender.com address.
- The reviewer's laptop does not need Docker or the development computer.

Deployment files:
- Dockerfile.render
- render.yaml
- backend/internal/web/site.go

## Free hosting limitations

Render Free services sleep after inactivity and can take around
a minute to restart. The first reviewer visit may be slow.

Neon Free suspends idle database compute but retains database
storage subject to plan limits and service policies.

Free hosting does not guarantee uninterrupted availability.
Do not describe this setup as always-on.

## Deployment checklist

1. Document branding and image sources; artwork reuse permissions
   have not been verified. See docs/ARTWORK.md.
2. Check GitHub Actions for the exact deployment commit. Run #4 was queued with no jobs; do not label it passed. Document this limitation alongside completed local and public tests.
3. Ensure .env files, database dumps, passwords, and API keys
   are excluded from Git.
4. Do not upload local user accounts, orders, or session data.

## Neon setup

1. Create a Neon Free PostgreSQL project.
2. Prefer a European database region near Render Frankfurt.
3. Keep the connection string secret.
4. Use a connection string with TLS enabled
   (normally sslmode=require).
5. For a confirmed empty database, apply these scripts in order:

   database/001_schema.sql
   database/002_seed.sql
   database/003_marketplace_migration.sql
   database/004_artwork_metadata.sql

6. Check each script for SQL errors before continuing.
7. Verify the expected sample record counts:

   products: 20
   offers: 60
   digital_keys: 120
   homepage_sections: 3


Do not reapply initialization scripts to a database containing
real demonstration accounts or orders without reviewing them.

## Render setup

1. Connect the GitHub repository to Render.
2. Select the deploy/render-free branch and its render.yaml. The blueprint uses Dockerfile.render.
3. Confirm the service uses the Free plan.
4. Provide DATABASE_URL as a secret environment variable.
5. Keep COOKIE_SECURE=true.
6. Deploy and obtain the actual HTTPS service URL.
7. Configure PUBLIC_ORIGIN to that exact HTTPS origin.
8. Confirm /api/health reports a connected database.

The committed blueprint sets `autoDeployTrigger: checksPass`. A pending CI run may prevent automatic deployment of newer commits. The actual Render dashboard settings and deployed commit must be checked there; a Git push alone does not prove a new deployment.

Never commit DATABASE_URL or paste it into a public issue.

## Public acceptance checks

- Homepage and store render over HTTPS.
- Direct product URLs work after refreshing.
- Product catalogue, search, filters and seller offers work.
- Registration and login work using Secure cookies.
- Approved simulated checkout creates an order.
- Declined simulated checkout does not allocate a key.
- Different user accounts cannot access each other's orders.
- Anonymous requests to private endpoints are rejected.
- A reviewer can open the URL from a separate browser.
- Test a cold visit after the service has been idle.

The author reported that the main HTTPS account/checkout/privacy checks and a visit after inactivity passed. No exact cold-start duration is recorded in this document. See the root README for evidence and pending CI status.

## Submission

Provide the GitHub repository URL, working public HTTPS URL,
and the complete required AI conversation transcript or share link.

Provide the full genuine conversation history separately; this deployment guide is not a transcript.
