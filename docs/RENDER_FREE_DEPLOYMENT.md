# Render Free + Neon Free Deployment

Status: Prepared and tested locally. Public deployment not yet verified.

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

## Before public deployment

1. Document branding and image sources; artwork reuse permissions
   have not been verified. See docs/ARTWORK.md.
2. Confirm GitHub Actions passes on the deployment branch.
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
2. Select the deploy/render-free branch and its render.yaml.
3. Confirm the service uses the Free plan.
4. Provide DATABASE_URL as a secret environment variable.
5. Keep COOKIE_SECURE=true.
6. Deploy and obtain the actual HTTPS service URL.
7. Configure PUBLIC_ORIGIN to that exact HTTPS origin.
8. Confirm /api/health reports a connected database.

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

Record actual results, known limitations and the final public URL.

## Submission

Provide the GitHub repository URL, working public HTTPS URL,
and the complete required AI conversation transcript or share link.

Do not claim public deployment or public acceptance tests passed
until they have actually been verified.
