@echo off

REM AegisGate Windows Deployment Script
REM This script sets up the development environment and builds AegisGate

echo ========================================
echo   AegisGate Deployment Script for Windows
echo ========================================
echo.

REM Check if Go is installed
echo [1/7] Checking Go installation...
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed!
    echo Please download and install Go from: https://golang.org/dl/
    pause
    exit /b 1
)
echo Go is installed.
go version

REM Verify the project directory structure
echo.
echo [2/7] Verifying project structure...
if not exist "go.mod" (
    echo ERROR: go.mod not found in current directory!
    echo Please run this script from the project root directory.
    pause
    exit /b 1
)

REM Directory structure for AegisGate
set PROJECT_ROOT=%cd%
set SRC_DIR=%PROJECT_ROOT%\src
set CONFIG_DIR=%PROJECT_ROOT%\config
set DOCS_DIR=%PROJECT_ROOT%\docs
set SCRIPTS_DIR=%PROJECT_ROOT%\scripts
set TESTS_DIR=%PROJECT_ROOT%\tests

echo Project root: %PROJECT_ROOT%
echo Source directory: %SRC_DIR%

REM Create directory structure if needed
echo.
echo [3/7] Creating directory structure...
if not exist "%SRC_DIR%\cmd\aegisgate" mkdir "%SRC_DIR%\cmd\aegisgate"
if not exist "%SRC_DIR%\pkg\certificate" mkdir "%SRC_DIR%\pkg\certificate"
if not exist "%SRC_DIR%\pkg\compliance" mkdir "%SRC_DIR%\pkg\compliance"
if not exist "%SRC_DIR%\pkg\config" mkdir "%SRC_DIR%\pkg\config"
if not exist "%SRC_DIR%\pkg\inspector" mkdir "%SRC_DIR%\pkg\inspector"
if not exist "%SRC_DIR%\pkg\metrics" mkdir "%SRC_DIR%\pkg\metrics"
if not exist "%SRC_DIR%\pkg\proxy" mkdir "%SRC_DIR%\pkg\proxy"
if not exist "%SRC_DIR%\pkg\scanner" mkdir "%SRC_DIR%\pkg\scanner"
if not exist "%SRC_DIR%\pkg\scanner\llm" mkdir "%SRC_DIR%\pkg\scanner\llm"
if not exist "%SRC_DIR%\pkg\scanner\matcher" mkdir "%SRC_DIR%\pkg\scanner\matcher"
if not exist "%SRC_DIR%\pkg\scanner\regex" mkdir "%SRC_DIR%\pkg\scanner\regex"
if not exist "%SRC_DIR%\pkg\tls" mkdir "%SRC_DIR%\pkg\tls"
if not exist "%CONFIG_DIR%" mkdir "%CONFIG_DIR%"
if not exist "%DOCS_DIR%" mkdir "%DOCS_DIR%"
if not exist "%SCRIPTS_DIR%" mkdir "%SCRIPTS_DIR%"
if not exist "%TESTS_DIR%\unit" mkdir "%TESTS_DIR%\unit"
if not exist "%TESTS_DIR%\integration" mkdir "%TESTS_DIR%\integration"

echo Directory structure created/verified.

REM Initialize Go modules if not already done
echo.
echo [4/7] Verifying Go modules...
if not exist "go.mod" (
    echo ERROR: go.mod file not found!
    echo Please ensure go.mod exists in project root.
    pause
    exit /b 1
)

REM Update go.mod if needed (replace placeholder)
echo Checking go.mod configuration...
type go.mod | findstr "module github.com/aegisgatesecurity/aegisgate" >nul
if %errorlevel% neq 0 (
    echo Updating go.mod to use correct module path...
    echo module github.com/aegisgatesecurity/aegisgate > go.mod
    echo. >> go.mod
    echo go 1.21 >> go.mod
    echo. >> go.mod
    echo require ( >> go.mod
    echo ) >> go.mod
)

echo Go modules verified.

REM Download dependencies (go mod tidy)
echo.
echo [5/7] Downloading dependencies...
go mod tidy
if %errorlevel% neq 0 (
    echo WARNING: Some dependencies failed to resolve
) else (
    echo Dependencies downloaded successfully.
)

REM Build the project
echo.
echo [6/7] Building AegisGate...
echo Building aegisgate.exe...
go build -o aegisgate.exe ./src/cmd/aegisgate/main.go
if %errorlevel% neq 0 (
    echo ERROR: Build failed!
    echo Please check the error messages above.
    pause
    exit /b 1
)

echo Build successful! aegisgate.exe created.

REM Generate SBOM (Security Bill of Materials)
echo.
echo [7/7] Generating SBOM (Security Bill of Materials)...
echo Note: SBOM generation requires syft tool (https://github.com/anchore/syft)
echo For now, generating a basic dependency listing...

echo ======================================== > sbom.txt
echo   AegisGate SBOM (Security Bill of Materials) >> sbom.txt
echo   Generated: %date% %time% >> sbom.txt
echo ======================================== >> sbom.txt
echo. >> sbom.txt
echo Dependencies: >> sbom.txt
go list -m all >> sbom.txt
echo. >> sbom.txt
echo Build artifacts: >> sbom.txt
echo - aegisgate.exe (main binary) >> sbom.txt
echo. >> sbom.txt
echo For full SBOM with syft: >> sbom.txt
echo syft . -o spdx-json ^> sbom.json >> sbom.txt

type sbom.txt

echo.
echo ========================================
echo   AegisGate Deployment Complete!
echo ========================================
echo.
echo Installation location: %PROJECT_ROOT%
echo Binary location: %PROJECT_ROOT%\aegisgate.exe
echo.
echo Next steps:
echo 1. Review the configuration in config\aegisgate.yml.example
echo 2. Copy config\aegisgate.yml.example to config\aegisgate.yml and customize
echo 3. Generate TLS certificates using the certificate management tools
echo 4. Run: .\aegisgate.exe --help for usage information
echo.
echo For testing: make test
echo For help: .\aegisgate.exe --help
echo.

pause
