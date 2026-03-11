// Phase 4 Build and SBOM Generation Script
// This script performs build validation and SBOM generation for AegisGate

# Check Go installation
$goVersion = go version
Write-Host "Go Version: $goVersion"

# Navigate to project directory
Set-Location C:\Users\Administrator\Desktop\Testing\aegisgate

# Update Go modules
Write-Host "Updating Go modules..."
go mod tidy

# Build the application
Write-Host "Building AegisGate application..."
go build -o aegisgate.exe ./src/cmd/aegisgate/

# Check if build succeeded
if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful: aegisgate.exe created"
} else {
    Write-Host "Build failed with exit code: $LASTEXITCODE"
    exit $LASTEXITCODE
}

# Run unit tests
Write-Host "Running unit tests..."
go test ./tests/unit/... -v

# Generate SBOM
Write-Host "Generating SBOM..."
syft dir . -o cyclonedx-json > sbom.json

# Check if SBOM was generated
if (Test-Path sbom.json) {
    Write-Host "SBOM generated successfully: sbom.json"
    $sbomSize = (Get-Item sbom.json).Length
    Write-Host "SBOM file size: $sbomSize bytes"
} else {
    Write-Host "SBOM generation failed"
    exit 1
}

Write-Host "Phase 4 Build and SBOM Generation completed successfully"
exit 0
