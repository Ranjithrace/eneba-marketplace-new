# Dependency lockfiles — important submission gate

## Finding

The earlier ZIP includes a **partial** `frontend/package-lock.json` with only the root package, and an empty `backend/go.sum`. Those are not suitable for a fresh-clone reproducible build.

Your Windows log on November 11 showed `prepare-lockfiles.ps1` creating complete files and subsequently running `npm ci`, `npm run build`, Go tests and Docker builds. Those **actual locally generated files** are the correct ones to keep. This ZIP cannot fetch them from your Windows filesystem and does not claim it contains them.

The review-fix installer deliberately **never replaces** `frontend/package-lock.json` or `backend/go.sum`. To inspect the real local files:

```powershell
cd 'D:\Projects\eneba-marketplace-new'
$lock = Get-Content '.\frontend\package-lock.json' -Raw | ConvertFrom-Json
"NPM locked packages: $($lock.packages.PSObject.Properties.Count)"
"Go checksum lines: $((Get-Content '.\backend\go.sum' | Measure-Object -Line).Lines)"
```

Both should be nontrivial; 1 npm package / 0 Go checksum lines means the broken ZIP placeholders were used. Run the supplied `prepare-lockfiles.ps1` **only if genuinely missing**, on the local networked machine, then run the test/build commands. It formats frontend code, so review changes before committing.

**Before pushing to GitHub:** `git add frontend/package-lock.json backend/go.sum` and commit them. Confirm `git ls-files frontend/package-lock.json backend/go.sum` prints both. CI now uses `npm ci` and `go mod verify` and must **fail** if these inputs are missing or incomplete; it no longer regenerates the lockfiles on CI.

For a true fresh-clone test, clone into a different folder, copy `.env.example` to `.env`, run the new-database initialization in README, and build. This has **not** been executed in the authoring workspace or verified on the user's Windows machine.
