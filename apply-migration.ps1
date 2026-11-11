# Must run from the LIVE project only, AFTER applying files and BEFORE rebuilding services.
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location -LiteralPath $root
if (!(Test-Path '.env')) { throw 'Existing .env not found; refusing to migrate.' }
if (!(Test-Path 'database/003_marketplace_migration.sql')) { throw 'Migration file missing.' }
$containerId = (docker compose ps -q db).Trim()
if ($LASTEXITCODE -ne 0 -or !$containerId) { throw 'The development database container must already be running.' }
$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$backup = "D:\Projects\eneba-db-before-full-update-$stamp.dump"
Write-Host 'Backing up the existing database before migration...'
docker compose exec -T db pg_dump -U eneba -d eneba_marketplace -Fc -f /tmp/eneba-upgrade-backup.dump
if ($LASTEXITCODE -ne 0) { throw 'Database backup command failed. No migration applied.' }
docker cp "${containerId}:/tmp/eneba-upgrade-backup.dump" "$backup"
if ($LASTEXITCODE -ne 0 -or !(Test-Path -LiteralPath $backup)) { throw 'Database backup transfer failed. No migration applied.' }
docker compose exec -T db rm -f /tmp/eneba-upgrade-backup.dump
if ($LASTEXITCODE -ne 0) { throw 'Could not clean temporary backup in container. No migration applied.' }
Write-Host "Backup file: $backup"
Write-Host 'Applying forward-only migration (transactional and retry-safe)...'
Get-Content 'database/003_marketplace_migration.sql' -Raw |
  docker compose exec -T db psql -X -v ON_ERROR_STOP=1 -U eneba -d eneba_marketplace
if ($LASTEXITCODE -ne 0) { throw 'Migration failed. Do not rebuild or re-run seed files.' }
Write-Host 'Migration applied. Existing products, offers, keys and orders were retained.'
