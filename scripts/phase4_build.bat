@echo off
echo Phase 4 Build and SBOM Generation Script
echo.

:: Check Go installation
go version
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Go installation not found
    exit /b 1
)

:: Navigate to project directory
cd /d C:\Users\Administrator\Desktop\Testing\aegisgate
echo Working directory: %CD%
echo.

:: Update Go modules
echo Updating Go modules...
go mod tidy
if %ERRORLEVEL% NEQ 0 (
    echo WARNING: go mod tidy completed with issues
)

:: Build the application
echo Building AegisGate application...
go build -o aegisgate.exe ./src/cmd/aegisgate/
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Build failed with exit code %ERRORLEVEL%
    exit /b 1
)

echo Build successful: aegisgate.exe created
echo.

:: Generate SBOM
echo Generating SBOM...
syft dir . -o cyclonedx-json > sbom.json
if %ERRORLEVEL% NEQ 0 (
    echo WARNING: SBOM generation completed with issues
)

:: Check if SBOM was generated
if exist sbom.json (
    for %%A in (sbom.json) do echo SBOM generated successfully: sbom.json (%%~zA bytes)
) else (
    echo WARNING: SBOM file not found after generation attempt
)

echo.
echo Phase 4 Build and SBOM Generation completed
exit /b 0
