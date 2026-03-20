# ============================================================================
# REMOVE PRIVATE KEYS FROM GIT HISTORY
# This script permanently removes .key files from git history using BFG
# PowerShell version
# ============================================================================

param(
    [switch]$Force
)

$ErrorActionPreference = "Stop"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "AegisGate - Remove Keys from Git History"
Write-Host "==========================================" -ForegroundColor Cyan

# Check if repo exists
if (-not (Test-Path ".git")) {
    Write-Host "ERROR: Not a git repository. Run this from the AegisGate root directory." -ForegroundColor Red
    exit 1
}

# Step 1: Backup
Write-Host ""
Write-Host "STEP 1: Backing up repository..." -ForegroundColor Yellow
$backupName = "aegisgate-backup-$(Get-Date -Format 'yyyyMMdd-HHmmss').bundle"
$backupPath = Join-Path $env:TEMP $backupName
git bundle create $backupPath --all 2>$null
Write-Host "Backup created at: $backupPath" -ForegroundColor Green

# Step 2: Check for existing .key files
Write-Host ""
Write-Host "STEP 2: Finding .key files in git history..." -ForegroundColor Yellow
$keyFiles = git log --all --pretty=format:%H -- "*.key" 2>$null
if ($keyFiles) {
    Write-Host "Found the following commits with .key files:" -ForegroundColor Red
    $keyFiles | ForEach-Object { Write-Host "  $_" }
    Write-Host ""
} else {
    Write-Host "No .key files found in git history." -ForegroundColor Green
    Write-Host "No action needed." -ForegroundColor Green
    exit 0
}

# Step 3: Download and run BFG
Write-Host ""
Write-Host "STEP 3: Setting up BFG Repo-Cleaner..." -ForegroundColor Yellow

$BFGJar = Join-Path $env:TEMP "bfg-1.14.0.jar"

# Check if Java is available
$javaVersion = java -version 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Java is required to run BFG Repo-Cleaner." -ForegroundColor Red
    Write-Host ""
    Write-Host "Please either:" -ForegroundColor Yellow
    Write-Host "  1. Install Java: https://adoptium.net/" -ForegroundColor White
    Write-Host "  2. Use GitHub's BFG UI to remove sensitive files" -ForegroundColor White
    Write-Host ""
    Write-Host "Alternative: In GitHub web UI:" -ForegroundColor Yellow
    Write-Host "  1. Go to repository Settings" -ForegroundColor White
    Write-Host "  2. Navigate to 'Repositories' section" -ForegroundColor White
    Write-Host "  3. Use 'Rename branch' option to create new history" -ForegroundColor White
    exit 1
}

# Download BFG if not present
if (-not (Test-Path $BFGJar)) {
    Write-Host "Downloading BFG Repo-Cleaner..." -ForegroundColor Yellow
    try {
        Invoke-WebRequest -Uri "https://repo1.maven.org/maven2/com/madgag/bfg/1.14.0/bfg-1.14.0.jar" -OutFile $BFGJar
    } catch {
        Write-Host "Failed to download BFG. Please install manually." -ForegroundColor Red
        exit 1
    }
}

# Step 4: Run BFG to remove .key files
Write-Host ""
Write-Host "STEP 4: Removing .key files from history..." -ForegroundColor Yellow
Write-Host "This may take a few minutes..." -ForegroundColor Gray

try {
    java -jar $BFGJar --delete-files '*.key' --no-blob-protection .
    if ($LASTEXITCODE -ne 0) {
        throw "BFG exited with error code $LASTEXITCODE"
    }
} catch {
    Write-Host "BFG failed: $_" -ForegroundColor Red
    exit 1
}

# Step 5: Clean up git reflog and garbage
Write-Host ""
Write-Host "STEP 5: Cleaning up git reflog and garbage..." -ForegroundColor Yellow
git reflog expire --expire=now --all
git gc --prune=now --aggressive

# Step 6: Verify removal
Write-Host ""
Write-Host "STEP 6: Verifying removal..." -ForegroundColor Yellow
$remainingKeys = git log --all --pretty=format:%H -- "*.key" 2>$null
if ($remainingKeys) {
    Write-Host "WARNING: Some .key references may still exist in history" -ForegroundColor Red
} else {
    Write-Host "SUCCESS: No .key files found in git history" -ForegroundColor Green
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "IMPORTANT: Next Steps" -ForegroundColor Yellow
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "1. Force push to update remote:" -ForegroundColor Yellow
Write-Host "   git push origin --force --all" -ForegroundColor White
Write-Host "   git push origin --force --tags" -ForegroundColor White
Write-Host ""
Write-Host "2. Notify team members to re-clone the repository" -ForegroundColor Yellow
Write-Host ""
Write-Host "3. Rotate ALL certificates that were in these files" -ForegroundColor Yellow
Write-Host ""
Write-Host "4. Verify on GitHub:" -ForegroundColor Yellow
Write-Host "   https://github.com/aegisgatesecurity/aegisgate/security" -ForegroundColor White
Write-Host "==========================================" -ForegroundColor Cyan
