#!/bin/bash
# version-sync.sh - Ensures version consistency across all AegisGate files

set -e

VERSION_FILE="VERSION"
MAIN_GO="cmd/aegisgate/main.go"

echo "=== AegisGate Version Sync Script ==="
echo ""

# Read version from VERSION file
VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')

if [ -z "$VERSION" ]; then
    echo "ERROR: VERSION file is empty"
    exit 1
fi

echo "Version from VERSION file: $VERSION"

# Check version in main.go
MAIN_GO_VERSION=$(grep -o 'const version = "[^"]*"' "$MAIN_GO" | grep -o '[0-9.]*')

if [ "$VERSION" != "$MAIN_GO_VERSION" ]; then
    echo "ERROR: Version mismatch!"
    echo "  VERSION file: $VERSION"
    echo "  main.go:      $MAIN_GO_VERSION"
    exit 1
fi

echo "Version in main.go: $MAIN_GO_VERSION"
echo ""
echo "=== All version checks passed ==="
exit 0
