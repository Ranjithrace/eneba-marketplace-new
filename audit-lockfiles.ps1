# Verifies the structure of existing dependency lockfiles in Windows PowerShell 5.1+.
# Does not change files, dependencies, containers, or databases.
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$lockFile = Join-Path $root 'frontend\package-lock.json'
$sumFile = Join-Path $root 'backend\go.sum'
if (-not (Test-Path -LiteralPath $lockFile -PathType Leaf)) { throw 'Missing frontend/package-lock.json' }
if (-not (Test-Path -LiteralPath $sumFile -PathType Leaf)) { throw 'Missing backend/go.sum' }

$raw = Get-Content -LiteralPath $lockFile -Raw
# npm represents the root package as an empty property name inside packages.
# ConvertFrom-Json in Windows PowerShell 5.1 rejects empty property names.
# Rename that one property in memory only; do not edit package-lock.json.
$normalized = $raw -replace '("packages"\s*:\s*\{\s*)""\s*:', '${1}"__npm_root__":'
if ($normalized -eq $raw) { throw 'Could not locate npm root package entry in package-lock.json' }
$lock = $normalized | ConvertFrom-Json
if (-not $lock.packages -or -not $lock.lockfileVersion) { throw 'Invalid npm lockfile structure' }
$npmEntries = @($lock.packages.PSObject.Properties).Count

$goLines = @(Get-Content -LiteralPath $sumFile | Where-Object { $_.Trim() })
$sumLines = $goLines.Count
$invalidGoLines = @($goLines | Where-Object { $_ -notmatch '^\S+\s+\S+\s+h1:[A-Za-z0-9+/=]+$' })
Write-Host "npm lock entries: $npmEntries"
Write-Host "Go checksum lines: $sumLines"
if ($npmEntries -lt 10 -or $sumLines -lt 4 -or $invalidGoLines.Count -gt 0) {
    throw 'Dependency lockfiles appear incomplete or malformed. Refer to docs/LOCKFILES.md.'
}
Write-Host 'LOCKFILE STRUCTURE CHECK PASSED (installation/build verification is separate).'
