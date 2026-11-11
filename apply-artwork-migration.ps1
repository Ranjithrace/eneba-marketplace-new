param([switch]$Apply)
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
if ($root -ne 'D:\Projects\eneba-marketplace-new') { throw 'Run this script from the LIVE D:\Projects\eneba-marketplace-new folder.' }
Set-Location -LiteralPath $root
if (-not (Test-Path -LiteralPath '.env')) { throw '.env missing; refusing migration.' }
if (-not (Test-Path -LiteralPath 'database\004_artwork_metadata.sql')) { throw 'SQL migration missing.' }
$ports = Get-Content '.env' -Raw
foreach ($p in @('FRONTEND_PORT=5173','BACKEND_PORT=8081','POSTGRES_PORT=55432')) {
 if ($ports -notmatch ('(?m)^' + [regex]::Escape($p) + '\s*$')) { throw "Port setting does not match expected: $p" }
}
Write-Host 'Migration 004 will only set empty products.cover_image_url values for known demo titles.'
Write-Host 'No users, orders, keys, sellers or other product columns will be changed.'
if (-not $Apply) { Write-Host 'PREVIEW ONLY: no data changed. Re-run with -Apply after reviewing.'; exit 0 }
$cid = (docker compose ps -q db).Trim()
if ($LASTEXITCODE -ne 0 -or -not $cid) { throw 'PostgreSQL container is not running.' }
$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$backup = "D:\Projects\eneba-db-before-artwork-$stamp.dump"
Write-Host 'Creating database backup before any SQL changes...'
docker compose exec -T db pg_dump -U eneba -d eneba_marketplace -Fc -f /tmp/eneba-artwork-backup.dump
if ($LASTEXITCODE -ne 0) { throw 'pg_dump failed; database unchanged.' }
docker cp "${cid}:/tmp/eneba-artwork-backup.dump" "$backup"
if ($LASTEXITCODE -ne 0 -or -not (Test-Path -LiteralPath $backup)) { throw 'Could not transfer database backup; database unchanged.' }
if ((Get-Item -LiteralPath $backup).Length -lt 1024) { throw 'Backup too small; database unchanged.' }
docker compose exec -T db rm -f /tmp/eneba-artwork-backup.dump
if ($LASTEXITCODE -ne 0) { throw 'Could not clean temporary backup; database unchanged.' }
Write-Host "Verified database dump: $backup"
Get-Content -LiteralPath 'database\004_artwork_metadata.sql' -Raw | docker compose exec -T db psql -X -v ON_ERROR_STOP=1 -U eneba -d eneba_marketplace
if ($LASTEXITCODE -ne 0) { throw 'Migration failed; restore dump if required.' }
Write-Host 'MIGRATION 004 EXECUTED. Check the artwork counts before claiming it passed.'
