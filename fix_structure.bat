@echo off
cd /d "C:\Users\Administrator\Desktop\Testing\aegisgate"

echo Checking src\pkg directory...
if exist "src\pkg" (
    echo Moving files from src\pkg to pkg...
    move /Y "src\pkg\*" "pkg\"
    echo Cleaning up empty directories...
    rmdir /S /Q "src\pkg"
    rmdir /S /Q "src"
    echo Directory structure updated successfully.
) else (
    echo src\pkg directory does not exist. Structure may already be correct.
)

echo.
echo Current directory structure:
dir pkg
