#!/bin/bash
# ============================================================================
# REMOVE PRIVATE KEYS FROM GIT HISTORY
# This script permanently removes .key files from git history using BFG
# ============================================================================

set -e

echo "=========================================="
echo "AegisGate - Remove Keys from Git History"
echo "=========================================="

# Check if repo exists
if [ ! -d ".git" ]; then
    echo "ERROR: Not a git repository. Run this from the AegisGate root directory."
    exit 1
fi

# Check for BFG or install it
if ! command -v bfg &> /dev/null; then
    echo "Installing BFG Repo-Cleaner..."
    
    # Try to download BFG
    BFG_JAR="/tmp/bfg-1.14.0.jar"
    
    if command -v java &> /dev/null; then
        if [ ! -f "$BFG_JAR" ]; then
            curl -L -o "$BFG_JAR" "https://repo1.maven.org/maven2/com/madgag/bfg/1.14.0/bfg-1.14.0.jar" 2>/dev/null || {
                echo "Failed to download BFG. Please install manually:"
                echo "  java -jar bfg-1.14.0.jar --delete-files '*.key'"
                exit 1
            }
        fi
        
        BFG_CMD="java -jar $BFG_JAR"
    else
        echo "ERROR: Java is required to run BFG Repo-Cleaner."
        echo "Please either:"
        echo "  1. Install Java: https://adoptium.net/"
        echo "  2. Or use GitHub's BFG UI to remove sensitive files"
        echo ""
        echo "Alternative: In GitHub web UI:"
        echo "  1. Go to repository Settings"
        echo "  2. Navigate to 'Repositories' section"  
        echo "  3. Use 'Rename branch' option to create new history"
        exit 1
    fi
else
    BFG_CMD="bfg"
fi

echo ""
echo "STEP 1: Backing up repository..."
git bundle create "/tmp/aegisgate-backup-$(date +%Y%m%d).bundle" --all 2>/dev/null || true
echo "Backup created at /tmp/aegisgate-backup-*.bundle"

echo ""
echo "STEP 2: Removing .key files from history..."
$BFG_CMD --delete-files '*.key' --no-blob-protection .

echo ""
echo "STEP 3: Cleaning up git reflog and garbage..."
git reflog expire --expire=now --all && git gc --prune=now --aggressive

echo ""
echo "STEP 4: Verifying removal..."
if git log --all --pretty=format:%H -- "*.key" 2>/dev/null | grep -q .; then
    echo "WARNING: Some .key references may still exist in history"
else
    echo "SUCCESS: No .key files found in git history"
fi

echo ""
echo "=========================================="
echo "IMPORTANT: Next Steps"
echo "=========================================="
echo "1. Force push to update remote:"
echo "   git push origin --force --all"
echo "   git push origin --force --tags"
echo ""
echo "2. Notify team members to re-clone the repository"
echo ""
echo "3. Rotate ALL certificates that were in these files"
echo ""
echo "4. Verify on GitHub:"
echo "   https://github.com/aegisgatesecurity/aegisgate/security"
echo "=========================================="
