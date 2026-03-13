@echo off
cd "C:\Users\Administrator\Desktop\Testing\AegisGate"
echo Staging files...
git add VERSION
git add "cmd\aegisgate\main.go"
git add ".github\workflows\release.yml"
echo Committing...
git commit -m "fix: Update v1.0.4 - VERSION, main.go, workflow fixes"
echo Tagging...
git tag -d v1.0.4 2>nul
git tag -a v1.0.4 -m "v1.0.4"
echo Pushing...
git push origin main --force
git push origin v1.0.4 --force
echo Done!
pause