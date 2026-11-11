param([switch]$Apply)
$ErrorActionPreference = 'Stop'
$source = $PSScriptRoot
$target = 'D:\Projects\eneba-marketplace-new'
if (!(Test-Path -LiteralPath $target -PathType Container)) { throw "Working project not found: $target" }
if ($source.TrimEnd('\') -ieq $target.TrimEnd('\')) { throw 'Extract ZIP to a separate preview folder first.' }
$envFile = Join-Path $target '.env'
if (!(Test-Path -LiteralPath $envFile)) { throw 'Existing .env is missing. Refusing update.' }
$envText = Get-Content -LiteralPath $envFile -Raw
foreach ($required in @('FRONTEND_PORT=5173','BACKEND_PORT=8081','POSTGRES_PORT=55432')) {
    if ($envText -notmatch "(?m)^$required\s*$") { throw "Expected local setting not found: $required. No files changed." }
}
$allowed = Get-Content -LiteralPath (Join-Path $source 'upgrade-known-hashes.json') -Raw | ConvertFrom-Json
$files = @(Get-ChildItem -LiteralPath $source -Recurse -File -Force | Where-Object { $_.FullName -notmatch '[\\/]node_modules[\\/]|[\\/]dist[\\/]|[\\/].git[\\/]' } | Sort-Object FullName)
$changes = @()
foreach ($file in $files) {
    $rel = $file.FullName.Substring($source.Length + 1)
    if ($rel -eq 'upgrade-known-hashes.json' -or $rel -eq 'upgrade-delete-stale.json') { continue }
    if ($rel -eq '.env' -or $rel -match '(^|\\)postgres_data(\\|$)') { throw "Forbidden source file: $rel" }
    $dest = Join-Path $target $rel
    $newHash = (Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
    if (Test-Path -LiteralPath $dest -PathType Leaf) {
        $current = (Get-FileHash -LiteralPath $dest -Algorithm SHA256).Hash.ToLowerInvariant()
        if ($current -eq $newHash) { continue }
        $entry = $allowed.PSObject.Properties[$rel]
        if ($null -eq $entry -or @($entry.Value) -notcontains $current) { throw "Unexpected local changes in $rel. Refusing to overwrite anything. Merge/backup manually." }
    }
    $changes += [PSCustomObject]@{ Source = $file.FullName; Dest = $dest; Relative = $rel }
}
$stale = Get-Content -LiteralPath (Join-Path $source 'upgrade-delete-stale.json') -Raw | ConvertFrom-Json
$staleFiles = @()
foreach ($item in $stale.PSObject.Properties) {
    $path = Join-Path $target $item.Name
    if (Test-Path -LiteralPath $path -PathType Leaf) {
        $hash = (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash.ToLowerInvariant()
        if ($hash -notin @($item.Value)) { throw "Unexpected local changes in legacy file $($item.Name). Refusing to remove anything." }
        $staleFiles += [PSCustomObject]@{ Dest = $path; Relative = $item.Name }
    }
}
Write-Host "Preview source: $source"
Write-Host "Live project: $target"
Write-Host 'Existing .env and local ports will be preserved.'
Write-Host "Files to install or replace: $($changes.Count)"
$changes | ForEach-Object { Write-Host " + $($_.Relative)" }
$staleFiles | ForEach-Object { Write-Host " - obsolete unchanged Stage 2.5 file: $($_.Relative)" }
if (!$Apply) { Write-Host 'PREVIEW ONLY: no files, containers or database data changed.'; exit 0 }
$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$backup = "D:\Projects\eneba-marketplace-full-backup-$stamp"
New-Item -ItemType Directory -Path $backup -Force | Out-Null
foreach ($change in $changes) {
    if (Test-Path -LiteralPath $change.Dest -PathType Leaf) {
        $destination = Join-Path $backup $change.Relative
        New-Item -ItemType Directory -Path (Split-Path -Parent $destination) -Force | Out-Null
        Copy-Item -LiteralPath $change.Dest -Destination $destination -Force
    }
}
foreach ($old in $staleFiles) {
    $destination = Join-Path $backup $old.Relative
    New-Item -ItemType Directory -Path (Split-Path -Parent $destination) -Force | Out-Null
    Copy-Item -LiteralPath $old.Dest -Destination $destination -Force
}
Write-Host "Existing overwritten/removed files backed up to: $backup"
foreach ($change in $changes) {
    New-Item -ItemType Directory -Path (Split-Path -Parent $change.Dest) -Force | Out-Null
    Copy-Item -LiteralPath $change.Source -Destination $change.Dest -Force
}
foreach ($old in $staleFiles) { Remove-Item -LiteralPath $old.Dest -Force }
Write-Host "Updated $($changes.Count) files successfully. No .env, database rows, or Docker services changed."
Write-Host 'NEXT: run apply-migration.ps1, then prepare-lockfiles.ps1, then docker compose up --build -d.'
