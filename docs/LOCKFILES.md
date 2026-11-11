# Dependency lockfiles

The submission branch includes `frontend/package-lock.json` and `backend/go.sum`. At the documentation review, the npm lockfile contained 77 package entries and the Go checksum file contained 28 nonempty lines. These counts confirm that the old empty/template files are no longer the committed versions; they are not a substitute for executing builds.

Use `npm ci` for frontend dependency installation and `go mod verify` for Go checksum validation. CI does not regenerate either file.

Earlier development ZIPs contained incomplete placeholders. That was a historical packaging issue: do not replace the current committed files with those copies.

For intentional dependency changes, regenerate the appropriate lockfile, review the diff and run the relevant checks before committing. Normal installation should use the committed files.

The author reported a successful fresh-clone Docker build on Windows. See the root README for reproducible setup and the distinction between local results and deployment-branch CI.
