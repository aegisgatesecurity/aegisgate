@echo off
REM AegisGate Quick Test Suite - Completes in under 5 minutes
setlocal enabledelayedexpansion

echo ========================================
echo AegisGate Quick Test Suite
echo ========================================
echo.

set PASSED=0
set FAILED=0

REM Test 1: Binary version check (10 seconds)
echo [1/5] Binary version check...
timeout /t 3 /nobreak >nul
cd /d C:\Users\Administrator\Desktop\Testing\AegisGate
aegisgate.exe --version >nul 2>&1
if %errorlevel% equ 0 (
    echo [1/5] PASS
    set /a PASSED+=1
) else (
    echo [1/5] FAIL
    set /a FAILED+=1
)

REM Test 2: Help command (5 seconds)
echo [2/5] Help command...
aegisgate.exe --help >nul 2>&1
if %errorlevel% equ 0 (
    echo [2/5] PASS
    set /a PASSED+=1
) else (
    echo [2/5] FAIL
    set /a FAILED+=1
)

REM Test 3: Quick unit tests (60 seconds max)
echo [3/5] Unit tests (pkg/core, pkg/config)...
cd /d C:\Users\Administrator\Desktop\Testing\AegisGate
go test ./pkg/core/... ./pkg/config/... -v -count=1 -timeout=30s >nul 2>&1
if %errorlevel% equ 0 (
    echo [3/5] PASS
    set /a PASSED+=1
) else (
    echo [3/5] FAIL
    set /a FAILED+=1
)

REM Test 4: Server starts in background (15 seconds - just verify it starts)
echo [4/5] Server startup (running in background)...
cd /d C:\Users\Administrator\Desktop\Testing\AegisGate
start /b aegisgate.exe -bind 0.0.0.0:9095 -tier community
timeout /t 5 /nobreak >nul
netstat -an | findstr "9095" >nul
if %errorlevel% equ 0 (
    echo [4/5] PASS
    set /a PASSED+=1
) else (
    echo [4/5] FAIL
    set /a FAILED+=1
)
taskkill /f /im aegisgate.exe >nul 2>&1

REM Test 5: Build verification (60 seconds)
echo [5/5] Build verification...
cd /d C:\Users\Administrator\Desktop\Testing\AegisGate
go build -o bin\aegisgate_test.exe .\cmd\aegisgate >nul 2>&1
if %errorlevel% equ 0 (
    echo [5/5] PASS
    set /a PASSED+=1
    del bin\aegisgate_test.exe >nul 2>&1
) else (
    echo [5/5] FAIL
    set /a FAILED+=1
)

echo.
echo ========================================
echo Results: %PASSED% passed, %FAILED% failed
echo ========================================

if %FAILED% gtr 0 (
    exit /b 1
) else (
    exit /b 0
)