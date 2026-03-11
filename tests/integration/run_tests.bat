@echo off
REM ATLAS Compliance Test Runner for Windows

echo ========================================
echo ATLAS Compliance Test Suite
echo ========================================
echo.

REM Set working directory
cd /d "%~dp0.."

REM Check if go is available
where go >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ERROR: Go is not installed or not in PATH
    exit /b 1
)

REM Parse command line arguments
set TEST_TARGET=Atlas
set VERBOSE=
set COVERAGE=
set RACE=

:parse_args
if "%~1"=="" goto run_tests
if "%~1"=="-v" set VERBOSE=1 & shift & goto parse_args
if "%~1"=="-verbose" set VERBOSE=1 & shift & goto parse_args
if "%~1"=="-coverage" set COVERAGE=1 & shift & goto parse_args
if "%~1"=="-race" set RACE=1 & shift & goto parse_args
if "%~1"=="-" set TEST_TARGET=%~2 & shift & shift & goto parse_args
shift
goto parse_args

:run_tests
echo Running ATLAS compliance tests...
echo.

if defined VERBOSE (
    set GO_TEST_FLAGS=-v -count=1
) else (
    set GO_TEST_FLAGS=-v
)

if defined COVERAGE (
    set GO_TEST_FLAGS=%GO_TEST_FLAGS% -coverprofile=coverage.out -covermode=atomic
)

if defined RACE (
    set GO_TEST_FLAGS=%GO_TEST_FLAGS% -race
)

set GO_TEST_FLAGS=%GO_TEST_FLAGS% -run %TEST_TARGET% ./tests/integration/...

go test %GO_TEST_FLAGS%

set TEST_RESULT=%ERRORLEVEL%

echo.
if %TEST_RESULT% equ 0 (
    echo ========================================
    echo All tests passed!
    echo ========================================
) else (
    echo ========================================
    echo Tests failed with error code: %TEST_RESULT%
    echo ========================================
)

if defined COVERAGE (
    if exist coverage.out (
        echo.
        echo Coverage Report:
        go tool cover -func=coverage.out
    )
)

exit /b %TEST_RESULT%
