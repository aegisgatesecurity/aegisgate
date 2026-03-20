# ============================================================================
# REMOVE ENTERPRISE/PREMIUM MODULES FROM PUBLIC REPO
# PowerShell version
# ============================================================================
# These modules should be in separate private repositories:
# - Enterprise: https://github.com/aegisgatesecurity/aegisgate-enterprise
# - Premium: https://github.com/aegisgatesecurity/aegisgate-premium
# ============================================================================

param(
    [switch]$Force,
    [switch]$SkipConfirmation
)

$ErrorActionPreference = "Stop"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "AegisGate - Remove Paid Modules from Public Repo"
Write-Host "==========================================" -ForegroundColor Cyan

# Directories to remove
$DirsToRemove = @(
    "pkg\compliance\enterprise",
    "pkg\compliance\premium"
)

# Check if directories exist
Write-Host ""
Write-Host "Checking for directories to remove..." -ForegroundColor Yellow
$foundDirs = @()
foreach ($dir in $DirsToRemove) {
    if (Test-Path $dir) {
        Write-Host "  Found: $dir" -ForegroundColor Green
        $foundDirs += $dir
    }
}

if ($foundDirs.Count -eq 0) {
    Write-Host "No enterprise/premium directories found." -ForegroundColor Green
    Write-Host "No action needed." -ForegroundColor Green
    exit 0
}

Write-Host ""
Write-Host "This script will PERMANENTLY DELETE the following directories:" -ForegroundColor Red
foreach ($dir in $foundDirs) {
    Write-Host "  - $dir" -ForegroundColor Red
}
Write-Host ""
Write-Host "These modules contain PAID features and should be moved to private repos." -ForegroundColor Yellow

# Confirmation prompt
if (-not $SkipConfirmation -and -not $Force) {
    Write-Host ""
    $confirm = Read-Host "Are you sure you want to continue? (yes/no)"
    if ($confirm -ne "yes") {
        Write-Host "Aborted." -ForegroundColor Yellow
        exit 0
    }
}

# Step 1: Backup
Write-Host ""
Write-Host "STEP 1: Creating backup..." -ForegroundColor Yellow
$backupDir = Join-Path $env:TEMP "aegisgate-modules-backup-$(Get-Date -Format 'yyyyMMdd-HHmmss')"
New-Item -ItemType Directory -Path $backupDir -Force | Out-Null

foreach ($dir in $foundDirs) {
    $backupPath = Join-Path $backupDir (Split-Path $dir -Leaf)
    Copy-Item -Path $dir -Destination $backupPath -Recurse -Force
    Write-Host "  Backed up: $dir -> $backupPath" -ForegroundColor Green
}
Write-Host "Backup location: $backupDir" -ForegroundColor Cyan

# Step 2: Remove from git history
Write-Host ""
Write-Host "STEP 2: Removing from git history..." -ForegroundColor Yellow

if (Test-Path ".git") {
    foreach ($dir in $foundDirs) {
        Write-Host "  Removing: $dir" -ForegroundColor Yellow
        git filter-branch --force --index-filter `
            "git rm --cached -r --ignore-unmatch $dir" `
            --prune-empty --tag-name-filter cat --all 2>$null
    }
    
    Write-Host "Cleaning git reflog..." -ForegroundColor Yellow
    git reflog expire --expire=now --all
    git gc --prune=now --aggressive
}

# Step 3: Remove locally
Write-Host ""
Write-Host "STEP 3: Removing directories locally..." -ForegroundColor Yellow

foreach ($dir in $foundDirs) {
    if (Test-Path $dir) {
        Remove-Item -Path $dir -Recurse -Force
        Write-Host "  Removed: $dir" -ForegroundColor Green
    }
}

# Step 4: Verify tier-manager.go references
Write-Host ""
Write-Host "STEP 4: Checking tier-manager.go references..." -ForegroundColor Yellow

if (Test-Path "pkg\compliance\tier-manager.go") {
    $tmContent = Get-Content "pkg\compliance\tier-manager.go" -Raw
    if ($tmContent -match 'aegisgate-enterprise|aegisgate-premium') {
        Write-Host "  tier-manager.go contains references to private repos (correct)" -ForegroundColor Green
    } else {
        Write-Host "  WARNING: tier-manager.go may need updating with private repo URLs" -ForegroundColor Yellow
    }
}

# Step 5: Final verification
Write-Host ""
Write-Host "STEP 5: Final verification..." -ForegroundColor Yellow

$remaining = Get-ChildItem -Path "pkg\compliance" -Recurse -File | Where-Object {
    $_.Name -match "enterprise|premium"
}

if ($remaining) {
    Write-Host "WARNING: Some files still reference enterprise/premium:" -ForegroundColor Red
    $remaining | ForEach-Object { Write-Host "  $_" }
} else {
    Write-Host "  No enterprise/premium files found" -ForegroundColor Green
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "IMPORTANT: Next Steps" -ForegroundColor Yellow
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "1. Push changes to remote:" -ForegroundColor Yellow
Write-Host "   git add -A" -ForegroundColor White
Write-Host "   git commit -m 'Remove enterprise/premium modules - move to private repos'" -ForegroundColor White
Write-Host "   git push origin --force --all" -ForegroundColor White
Write-Host ""
Write-Host "2. Ensure private repos exist:" -ForegroundColor Yellow
Write-Host "   - https://github.com/aegisgatesecurity/aegisgate-enterprise" -ForegroundColor White
Write-Host "   - https://github.com/aegisgatesecurity/aegisgate-premium" -ForegroundColor White
Write-Host ""
Write-Host "3. Copy backed up modules to private repos:" -ForegroundColor Yellow
Write-Host "   Copy-Item -Path '$backupDir\*' -Destination 'C:\path\to\private\repo\' -Recurse" -ForegroundColor White
Write-Host ""
Write-Host "4. Verify GitHub:" -ForegroundColor Yellow
Write-Host "   - Check Settings > Danger Zone for any exposure" -ForegroundColor White
Write-Host "   - Enable 'Forking' restrictions if needed" -ForegroundColor White
Write-Host "   - Set repository to Private if required" -ForegroundColor White
Write-Host "==========================================" -ForegroundColor Cyan
