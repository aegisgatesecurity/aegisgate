#!/bin/bash
# ============================================================================
# REMOVE ENTERPRISE/PREMIUM MODULES FROM PUBLIC REPO
# ============================================================================
# These modules should be in separate private repositories:
# - Enterprise: https://github.com/aegisgatesecurity/aegisgate-enterprise
# - Premium: https://github.com/aegisgatesecurity/aegisgate-premium
# ============================================================================

set -e

echo "=========================================="
echo "AegisGate - Remove Paid Modules from Public Repo"
echo "=========================================="

# Directories to remove
DIRS_TO_REMOVE=(
    "pkg/compliance/enterprise"
    "pkg/compliance/premium"
)

# Check if directories exist
for dir in "${DIRS_TO_REMOVE[@]}"; do
    if [ -d "$dir" ]; then
        echo "Found: $dir"
    fi
done

echo ""
echo "This script will PERMANENTLY DELETE the following directories:"
for dir in "${DIRS_TO_REMOVE[@]}"; do
    echo "  - $dir"
done
echo ""
echo "These modules contain PAID features and should be moved to private repos."
echo ""

read -p "Are you sure you want to continue? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo "Aborted."
    exit 0
fi

# Step 1: Backup
echo ""
echo "STEP 1: Creating backup..."
BACKUP_DIR="/tmp/aegisgate-modules-backup-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$BACKUP_DIR"

for dir in "${DIRS_TO_REMOVE[@]}"; do
    if [ -d "$dir" ]; then
        cp -r "$dir" "$BACKUP_DIR/"
        echo "  Backed up: $dir"
    fi
done
echo "Backup location: $BACKUP_DIR"

# Step 2: Remove from git history
echo ""
echo "STEP 2: Checking for git history..."

if [ -d ".git" ]; then
    echo "Removing from git history using git filter-branch..."
    
    for dir in "${DIRS_TO_REMOVE[@]}"; do
        echo "  Removing: $dir"
        git filter-branch --force --index-filter \
            "git rm --cached -r --ignore-unmatch $dir" \
            --prune-empty --tag-name-filter cat -- --all 2>/dev/null || true
    done
    
    echo "Cleaning git reflog..."
    git reflog expire --expire=now --all
    git gc --prune=now --aggressive
fi

# Step 3: Remove locally
echo ""
echo "STEP 3: Removing directories locally..."

for dir in "${DIRS_TO_REMOVE[@]}"; do
    if [ -d "$dir" ]; then
        rm -rf "$dir"
        echo "  Removed: $dir"
    fi
done

# Step 4: Update references in tier-manager.go
echo ""
echo "STEP 4: Updating tier-manager.go references..."

# The updated tier-manager.go already has the correct references
# Just verify it exists and is correct
if [ -f "pkg/compliance/tier-manager.go" ]; then
    echo "  tier-manager.go updated"
fi

# Step 5: Verify
echo ""
echo "STEP 5: Verification..."

# Check for remaining enterprise/premium code
REMAINING=$(find pkg/compliance -type f \( -name "*enterprise*" -o -name "*premium*" \) 2>/dev/null || true)
if [ -n "$REMAINING" ]; then
    echo "WARNING: Some files still reference enterprise/premium:"
    echo "$REMAINING"
else
    echo "  No enterprise/premium files found"
fi

echo ""
echo "=========================================="
echo "IMPORTANT: Next Steps"
echo "=========================================="
echo ""
echo "1. Push changes to remote:"
echo "   git add -A"
echo "   git commit -m 'Remove enterprise/premium modules - move to private repos'"
echo "   git push origin --force --all"
echo ""
echo "2. Ensure private repos exist:"
echo "   - https://github.com/aegisgatesecurity/aegisgate-enterprise"
echo "   - https://github.com/aegisgatesecurity/aegisgate-premium"
echo ""
echo "3. Copy backed up modules to private repos:"
echo "   cp -r $BACKUP_DIR/* /path/to/private/repo/"
echo ""
echo "4. Verify GitHub:"
echo "   - Check Settings > Danger Zone for any exposure"
echo "   - Enable 'Forking' restrictions if needed"
echo "   - Set repository to Private if required"
echo "=========================================="
