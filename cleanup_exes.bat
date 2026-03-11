@echo off
cd /d C:\Users\Administrator\Desktop\Testing\AegisGate
mkdir _archive 2>nul

echo Moving duplicate .exe files to archive...
move aegisgate_atlas_final.exe _archive\ 2>nul
move aegisgate_atlas_threshold_final.exe _archive\ 2>nul
move aegisgate_final.exe _archive\ 2>nul
move aegisgate_final_v2.exe _archive\ 2>nul
move aegisgate_fixed.exe _archive\ 2>nul
move aegisgate_original.exe _archive\ 2>nul
move aegisgate_production.exe _archive\ 2>nul
move aegisgate_secure.exe _archive\ 2>nul
move aegisgate_test.exe _archive\ 2>nul
move aegisgate_v0.2.exe _archive\ 2>nul
move aegisgate_v1.0.exe _archive\ 2>nul
move aegisgate_v1.1.exe _archive\ 2>nul
move aegisgate_v1.1_final.exe _archive\ 2>nul
move aegisgate_v1.2.exe _archive\ 2>nul
move aegisgate_v1.3.exe _archive\ 2>nul
move aegisgate_working.exe _archive\ 2>nul
move check.exe _archive\ 2>nul

echo Cleanup complete
echo Remaining .exe files:
dir aegisgate.exe

echo Archive contents:
dir _archive
pause