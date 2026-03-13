# Fix script for AegisGate v1.0.4
Set-Location 'C:\Users\Administrator\Desktop\Testing\AegisGate'

Write-Host "=== AegisGate v1.0.4 Fix Push ===" -ForegroundColor Cyan

# Check current state
Write-Host "`nCurrent VERSION file:" -ForegroundColor Yellow
Get-Content VERSION
Write-Host "`nMain.go version line:" -ForegroundColor Yellow
Select-String -Path "cmd\aegisgate\main.go" -Pattern "const version"

# Stage and commit
Write-Host "`nStaging files..." -ForegroundColor Green
git add VERSION
git add "cmd\aegisgate\main.go"
git add ".github\workflows\release.yml"

# Commit
Write-Host "Committing..." -ForegroundColor Green
git commit -m "fix: Version consistency and workflow fixes

- Update VERSION to 1.0.4
- Update main.go version to 1.0.4
- Fix anchore/sbom-action from v1 to v0.16.1"

# Show commits
Write-Host "`nGit log:" -ForegroundColor Yellow
git log --oneline -3

# Update existing tag
Write-Host "`nUpdating tag v1.0.4..." -ForegroundColor Green
git tag -d v1.0.4
git tag -a v1.0.4 -m "AegisGate v1.0.4 - Production Helm Charts Release"

# Push
Write-Host "`nPushing to remote..." -ForegroundColor Green
git push origin main --force
git push origin v1.0.4 --force

Write-Host "`n=== DONE ===" -ForegroundColor Cyan