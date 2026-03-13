$ErrorActionPreference = 'Stop'
Set-Location 'C:\Users\Administrator\Desktop\Testing\AegisGate'

# Initialize git if needed
if (-not (Test-Path '.git')) {
    Write-Host 'Initializing git repository...'
    git init
    git config user.email 'dev@aegisgate.com'
    git config user.name 'AegisGate Dev'
}

# Add files
Write-Host 'Adding files...'
git add README.md
git add RELEASE_NOTES_v1.0.4.md

# Commit
Write-Host 'Committing...'
git commit -m 'Release v1.0.4: Production-ready Helm charts with full validation'

# Show log
Write-Host 'Commit log:'
git log --oneline -3

# Create tag
Write-Host 'Creating tag v1.0.4...'
git tag -a v1.0.4 -m 'AegisGate v1.0.4 - Production Helm Charts Release'

# Push to remote
Write-Host 'Pushing to remote...'
git push origin master
git push origin v1.0.4

Write-Host 'Done!'