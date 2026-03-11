#!/bin/bash
# AegisGate Production Build Script
# Version: 0.2.0
# Build script for production deployment

set -e

echo "========================================"
echo "  AegisGate Production Build Pipeline"
echo "========================================"

# Build the application
echo "[1/5] Building AegisGate..."
go build -o aegisgate.exe -ldflags="-s -w" ./cmd/aegisgate/

# Generate SBOM
echo "[2/5] Generating Software Bill of Materials..."
syft dir . -o cyclonedx-json > sbom.json

# Run tests
echo "[3/5] Running comprehensive test suite..."
go test ./pkg/... -v -race -coverprofile=coverage.out

# Build Docker image
echo "[4/5] Building Docker image..."
docker build -f Dockerfile.production -t aegisgate:latest .

# Tag for registry
echo "[5/5] Tagging image for deployment..."
docker tag aegisgate:latest aegisgate:production-v0.2.0

echo "========================================"
echo "  Build Pipeline Complete!"
echo "========================================"
echo "Output files:"
echo "  - aegisgate.exe"
echo "  - sbom.json"
echo "  - coverage.out"
echo "  - Docker image: aegisgate:latest"
echo "========================================"
