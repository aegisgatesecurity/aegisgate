$ErrorActionPreference = "Continue"
Set-Location "C:\Users\Administrator\Desktop\Testing\aegisgate"

Write-Host "=== Git Status ==="
git status --short

Write-Host "`n=== Git Remote ==="
git remote -v

Write-Host "`n=== Git Branch ==="
git branch

Write-Host "`n=== Files changed ==="
git diff --name-only HEAD 2>$null

Write-Host "`n=== Untracked files ==="
git ls-files --others --exclude-standard
