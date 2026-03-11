@echo off
REM AegisGate Deployment Script for Windows
REM Version: 0.1.0
REM Purpose: Configure testing/dev environment for AegisGate

echo [+] Starting AegisGate Dev Environment Setup...
echo [i] This script will configure your Windows system for AegisGate development and testing.

REM Check if Go is installed
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [x] ERROR: Go is not installed or not in PATH
    echo [i] Please install Go from https://go.dev/dl/ and add it to your PATH
    exit /b 1
)

echo [✓] Go is installed
go version

REM Check if Git is installed
where git >nul 2>&1
if %errorlevel% neq 0 (
    echo [x] ERROR: Git is not installed or not in PATH
    echo [i] Please install Git from https://git-scm.com/download/win
    exit /b 1
)

echo [✓] Git is installed
git --version

REM Create project directories
echo [+] Creating project directories...
mkdir config >nul 2>&1
mkdir docs >nul 2>&1
mkdir scripts >nul 2>&1
mkdir tests >nul 2>&1
mkdir src >nul 2>&1
mkdir srccmd >nul 2>&1
mkdir srccmdaegisgate >nul 2>&1
mkdir srcpkg >nul 2>&1
mkdir srcpkgcertificate >nul 2>&1
mkdir srcpkgcompliance >nul 2>&1
mkdir srcpkgconfig >nul 2>&1
mkdir srcpkginspector >nul 2>&1
mkdir srcpkgmetrics >nul 2>&1
mkdir srcpkgproxy >nul 2>&1
mkdir srcpkgscanner >nul 2>&1
mkdir srcpkgscannerllm >nul 2>&1
mkdir srcpkgscannermatcher >nul 2>&1
mkdir srcpkgscanner
egex >nul 2>&1
mkdir srcpkg	ls >nul 2>&1

REM Initialize Go module if not already done
if not exist go.mod (
    echo [+] Initializing Go module...
    go mod init github.com/aegisgatesecurity/aegisgate
) else (
    echo [✓] Go module already initialized
)

REM Update go.mod if needed
if exist go.mod (
    findstr /c:"module github.com/aegisgatesecurity/aegisgate" go.mod > nul
    if %errorlevel% neq 0 (
        echo [i] Updating module path in go.mod...
        (echo module github.com/aegisgatesecurity/aegisgate) > go.mod
    ) else (
        echo [✓] Module path is correct
    )
)

REM Get dependencies
echo [+] Running go mod tidy...
go mod tidy

REM Create example config file
if not exist configaegisgate.yml (
    echo [+] Creating example config file...
    (
        echo # AegisGate Configuration
        echo server:
        echo   bind_address: "0.0.0.0:8443"
        echo   tls:
        echo     enabled: true
        echo     cert_file: "certs/server.crt"
        echo     key_file: "certs/server.key"
        echo     ca_file: "certs/ca.crt"
        echo     mode: "MITM"
        echo compliance:
        echo   enabled: true
        echo   frameworks:
        echo     - "MITRE_ATLAS"
        echo     - "NIST_AI_RMF"
        echo     - "OWASP_TOP_10_AI"
        echo logging:
        echo   level: "debug"
        echo   format: "json"
        echo   output: "stdout"
    ) > configaegisgate.yml.example

    echo [✓] Example config created at configaegisgate.yml.example
)

REM Build the project
echo [+] Building AegisGate...
go build -o aegisgate.exe ./src/cmd/aegisgate

if exist aegisgate.exe (
    echo [✓] Build successful!
    echo [i] Running version check...
    .aegisgate.exe --version
) else (
    echo [x] Build failed!
    exit /b 1
)

echo [+] Dev environment setup complete!
echo [i] Next steps:
echo    - Configure configaegisgate.yml with your settings
echo    - Generate TLS certificates (see docs/tls_setup.md)
echo    - Run tests: make test
echo    - Start dev server: .aegisgate.exe --config configaegisgate.yml
