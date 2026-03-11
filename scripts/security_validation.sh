#!/bin/bash
# AegisGate Security Validation Script
# Version: 0.2.0
# Purpose: Run security validation checks for Phase 5A

set -e

echo "AegisGate Security Validation Script"
echo "=================================="
echo "Date: $(date)"
echo ""

# Check Go installation
echo "[1/5] Checking Go installation..."
if command -v go &> /dev/null; then
    echo "✓ Go is installed"
    go version
else
    echo "✗ Go is not installed"
    exit 1
fi

# Run gosec static analysis
echo ""
echo "[2/5] Running gosec static analysis..."
gosec -quiet ./...

# Run govulncheck for vulnerability scanning
echo ""
echo "[3/5] Running govulncheck..."
govulncheck ./...

# Run security-focused unit tests
echo ""
echo "[4/5] Running OPSEC security tests..."
go test -v ./pkg/opsec/... -run "OPSEC"

# Run integration tests
echo ""
echo "[5/5] Running integration tests..."
go test -v ./tests/integration/... -run "TestHTTP|TestCore"

echo ""
echo "Security validation completed successfully!"
echo ""
echo "Next steps:"
echo "1. Review security findings report"
echo "2. Address any critical security issues"
echo "3. Generate security scan report"
echo "4. Document security assessment"
