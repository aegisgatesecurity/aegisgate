Set-Location 'C:\Users\Administrator\Desktop\Testing\aegisgate'

# Remove temp files
Remove-Item -Force -ErrorAction SilentlyContinue fix_flaky.py
Remove-Item -Force -ErrorAction SilentlyContinue read_test.py
Remove-Item -Force -ErrorAction SilentlyContinue commit_msg.txt
Remove-Item -Force -ErrorAction SilentlyContinue PROJECT_MEMORY_ANCHOR.md

# Stage and commit
git add pkg/proxy/http2_test.go
git commit -m "Fix flaky TestStreamLimiterConcurrentAccess by adding proper synchronization

The test was failing intermittently (99/100 streams) due to a race
between wg.Wait() returning and internal map updates being visible.
Added:
- acquired channel to signal when each AllowStream completes
- Drain the channel to ensure all streams are fully registered
- Small sleep to allow mutex state to settle after goroutines exit"

# Push
git push

Write-Host "Done!"
