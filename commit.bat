@echo off
cd /d C:\Users\Administrator\Desktop\Testing\aegisgate
git add -f pkg/proxy/http2_test.go
git commit -m "Fix flaky TestStreamLimiterConcurrentAccess by adding proper synchronization

The test was failing intermittently (99/100 streams) due to a race
between wg.Wait() returning and internal map updates being visible.
Added acquired channel to signal when each AllowStream completes,
drain the channel to ensure all streams are fully registered,
and small sleep to allow mutex state to settle after goroutines exit."
git push
echo Done!
