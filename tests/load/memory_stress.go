package load

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"
	"testing"
	"time"
)

// MemoryStressTest detects memory leaks under sustained load
type MemoryStressTest struct {
	TargetHost   string
	TargetPort   int
	Duration     time.Duration
	ConnInterval time.Duration
}

func NewMemoryStressTest(host string, port int) *MemoryStressTest {
	return &MemoryStressTest{
		TargetHost:   host,
		TargetPort:   port,
		Duration:     24 * time.Hour,
		ConnInterval: 100 * time.Millisecond,
	}
}

type MemoryResult struct {
	StartHeap    uint64
	EndHeap      uint64
	PeakHeap     uint64
	HeapGrowth   int64
	Connections  int
	Duration     time.Duration
	LeakDetected bool
	Samples      []MemorySample
}

type MemorySample struct {
	Timestamp time.Time
	HeapAlloc uint64
	HeapSys   uint64
	NumGC     uint32
}

func (m *MemoryStressTest) Run(ctx context.Context) (*MemoryResult, error) {
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	result := &MemoryResult{
		StartHeap: m1.HeapAlloc,
		PeakHeap:  m1.HeapAlloc,
		Samples:   make([]MemorySample, 0),
	}

	ctx, cancel := context.WithTimeout(ctx, m.Duration)
	defer cancel()

	var wg sync.WaitGroup
	connCount := 0
	ticker := time.NewTicker(m.ConnInterval)
	defer ticker.Stop()

	sampleTicker := time.NewTicker(30 * time.Second)
	defer sampleTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			goto done

		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				conn, err := net.DialTimeout("tcp",
					fmt.Sprintf("%s:%d", m.TargetHost, m.TargetPort),
					5*time.Second)
				if err == nil {
					time.Sleep(100 * time.Millisecond)
					_ = conn.Close()
				}
			}()
			connCount++

		case <-sampleTicker.C:
			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)
			result.Samples = append(result.Samples, MemorySample{
				Timestamp: time.Now(),
				HeapAlloc: ms.HeapAlloc,
				HeapSys:   ms.HeapSys,
				NumGC:     ms.NumGC,
			})
			if ms.HeapAlloc > result.PeakHeap {
				result.PeakHeap = ms.HeapAlloc
			}
		}
	}

done:
	wg.Wait()

	runtime.GC()
	runtime.ReadMemStats(&m2)
	result.EndHeap = m2.HeapAlloc
	result.HeapGrowth = int64(m2.HeapAlloc) - int64(m1.HeapAlloc)
	result.Connections = connCount
	result.Duration = m.Duration
	result.LeakDetected = m.detectLeak(result.Samples)

	return result, nil
}

func (m *MemoryStressTest) detectLeak(samples []MemorySample) bool {
	if len(samples) < 10 {
		return false
	}

	firstHalf := samples[:len(samples)/2]
	secondHalf := samples[len(samples)/2:]

	var firstAvg, secondAvg uint64
	for _, s := range firstHalf {
		firstAvg += s.HeapAlloc
	}
	for _, s := range secondHalf {
		secondAvg += s.HeapAlloc
	}

	firstAvg /= uint64(len(firstHalf))
	secondAvg /= uint64(len(secondHalf))

	return secondAvg > firstAvg*2
}

func ValidateMemoryResult(r *MemoryResult) error {
	if r.LeakDetected {
		return fmt.Errorf("potential memory leak detected: heap grew by %d bytes", r.HeapGrowth)
	}
	if r.HeapGrowth > 100*1024*1024 {
		return fmt.Errorf("excessive heap growth: %d bytes", r.HeapGrowth)
	}
	return nil
}

func BenchmarkMemoryStress(b *testing.B) {
	test := NewMemoryStressTest("localhost", 8443)
	test.Duration = 1 * time.Minute
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := test.Run(ctx)
		if err != nil {
			b.Fatalf("memory stress test failed: %v", err)
		}
	}
}
