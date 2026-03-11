@echo off
REM AegisGate Production Build Script (Windows)
REM Version: 0.2.0

echo ========================================
echo   AegisGate Production Build Pipeline
echo ========================================

echo [1/5] Building AegisGate...
go build -o aegisgate.exe -ldflags="-s -w" ./cmd/aegisgate/

echo [2/5] Generating Software Bill of Materials...
syft dir . -o cyclonedx-json > sbom.json

echo [3/5] Running comprehensive test suite...
go test ./pkg/... -v -race -coverprofile=coverage.out

echo [4/5] Building Docker image...
docker build -f Dockerfile.production -t aegisgate:latest .

echo [5/5] Tagging image for deployment...
docker tag aegisgate:latest aegisgate:production-v0.2.0

echo ========================================
echo   Build Pipeline Complete!
echo ========================================
echo Output files:
echo   - aegisgate.exe
echo   - sbom.json  
echo   - coverage.out
echo   - Docker image: aegisgate:latest
echo ========================================
