# Requires Docker Desktop, runs in the EXISTING live project after the safe update.
# Resolves complete transitive dependency locks on the user's network and runs checks.
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$frontend = Join-Path $root 'frontend'
$backend = Join-Path $root 'backend'
if (!(Test-Path (Join-Path $root '.env'))) { throw 'Missing .env. The working local configuration must be preserved.' }
Write-Host 'Resolving and verifying frontend dependencies using Node 24...'
docker run --rm --mount "type=bind,source=$frontend,target=/app" --workdir /app node:24-alpine sh -c 'npm install --package-lock-only --no-audit --no-fund && npm ci --no-audit --no-fund && npm run format && npm run format:check && npm run test && npm run build'
if ($LASTEXITCODE -ne 0) { throw 'Frontend dependency or test/build verification failed.' }
Write-Host 'Resolving and verifying Go modules using Go 1.27...'
docker run --rm --mount "type=bind,source=$backend,target=/src" --workdir /src golang:1.27-alpine sh -c 'go mod tidy && go mod verify && go test ./... && go vet ./... && go build -o /tmp/eneba-api ./cmd/api'
if ($LASTEXITCODE -ne 0) { throw 'Backend dependency or test/build verification failed.' }
Write-Host 'Full package-lock.json and go.sum are now generated in the project folder.'
Write-Host 'Keep and commit these lockfiles. This script did not touch the database or .env.'
