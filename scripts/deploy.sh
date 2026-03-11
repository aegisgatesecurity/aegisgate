@echo off
REM AegisGate Deployment Script for Windows
REM Version: 0.1.0
REM Purpose: Set up development environment and build AegisGate

echo AegisGate v0.1.0 Deployment Script
echo ===================================

REM Check if Go is installed
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    exit /b 1
)

REM Get Go version
echo Checking Go version...
go version

REM Create required directories
echo Creating project directories...
if not exist "%~dp0aegisgate" mkdir "%~dp0aegisgate"
if not exist "%~dp0aegisgateconfig" mkdir "%~dp0aegisgateconfig"
if not exist "%~dp0aegisgatedocs" mkdir "%~dp0aegisgatedocs"
if not exist "%~dp0aegisgatescripts" mkdir "%~dp0aegisgatescripts"
if not exist "%~dp0aegisgatesrc" mkdir "%~dp0aegisgatesrc"
if not exist "%~dp0aegisgatesrccmd" mkdir "%~dp0aegisgatesrccmd"
if not exist "%~dp0aegisgatesrcpkg" mkdir "%~dp0aegisgatesrcpkg"
if not exist "%~dp0aegisgate	ests" mkdir "%~dp0aegisgate	ests"

REM Navigate to project directory
cd /d "%~dp0aegisgate"

REM Initialize git if not already done
if not exist ".git" (
    echo Initializing git repository...
    git init
) else (
    echo Git repository already initialized
)

REM Check if remote is configured
git remote -v | findstr "aegisgatesecurity" >nul 2>&1
if %errorlevel% neq 0 (
    echo Setting up GitHub remote...
    git remote add origin https://github.com/aegisgatesecurity/aegisgate.git
) else (
    echo GitHub remote already configured
)

REM Update go.mod if needed
if exist "go.mod" (
    echo Validating go.mod...
    type go.mod | findstr "aegisgatesecurity" >nul 2>&1
    if %errorlevel% neq 0 (
        echo Updating go.mod to use aegisgatesecurity/aegisgate...
        powershell -Command "(Get-Content go.mod) -replace 'github.com/block/goose/aegisgate', 'github.com/aegisgatesecurity/aegisgate' | Set-Content go.mod"
    ) else (
        echo go.mod already correctly configured
    )
) else (
    echo Initializing Go module...
    go mod init github.com/aegisgatesecurity/aegisgate
    echo go.mod created successfully
)

REM Run go mod tidy to resolve dependencies
echo Running go mod tidy...
go mod tidy

REM Build the project
echo Building AegisGate...
go build -o aegisgate.exe ./src/cmd/aegisgate

if exist "aegisgate.exe" (
    echo AegisGate built successfully!
    echo Running version check...
    .aegisgate.exe --version
    echo.
    echo Deployment completed successfully!
    echo ================================
    echo To run AegisGate, execute: .aegisgate.exe
) else (
    echo Build failed. Please check error messages above.
    exit /b 1
)

REM Run tests
echo.
echo Running tests...
go test ./... -v

echo.
echo ====================================
echo Deployment script completed!
echo Next steps:
echo 1. Configure AegisGate using config/aegisgate.yml
echo 2. Run 'make build' to rebuild
echo 3. Run 'make test' to verify
echo 4. Push changes: git add . && git commit -m "config: initial deployment setup" && git push -u origin main
