package load

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"
	"testing"
	"time"
)

// LatencyBenchmark measures response times under load
type LatencyBenchmark struct {
	TargetHost      string
	TargetPort      int
	RequestCount    int
	Concurrency     int
	RequestInterval time.Duration
}

func NewLatencyBenchmark(host string, port int) *LatencyBenchmark {
	return &LatencyBenchmark{
		TargetHost:      host,
		TargetPort:      port,
		RequestCount:    10000,
		Concurrency:     100,
		RequestInterval: 1 * time.Millisecond,
	}
}

type LatencyResult struct {
	RequestsSent   int
	RequestsFailed int
	Latencies      []time.Duration
	MinLatency     time.Duration
	MaxLatency     time.Duration
	AvgLatency     time.Duration
	P50Latency     time.Duration
	P99Latency     time.Duration
	P999Latency    time.Duration
	TotalDuration  time.Duration
}

func (l *LatencyBenchmark) Run(ctx context.Context) (*LatencyResult, error) {
	latencies := make([]time.Duration, 0, l.RequestCount)
	var mu sync.Mutex
	var failed int
	var wg sync.WaitGroup

	requestChan := make(chan int, l.RequestCount)
	for i := 0; i < l.RequestCount; i++ {
		requestChan <- i
	}
	close(requestChan)

	start := time.Now()

	for i := 0; i < l.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requestChan {
				select {
				case <-ctx.Done():
					return
				default:
				}

				reqStart := time.Now()
				conn, err := net.DialTimeout("tcp",
					fmt.Sprintf("%s:%d", l.TargetHost, l.TargetPort),
					5*time.Second)
				if err != nil {
					mu.Lock()
					failed++
					mu.Unlock()
					continue
				}
				_ = conn.Close()
				elapsed := time.Since(reqStart)

				mu.Lock()
				latencies = append(latencies, elapsed)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	totalDuration := time.Since(start)

	return l.calculateResults(latencies, failed, totalDuration), nil
}

func (l *LatencyBenchmark) calculateResults(latencies []time.Duration, failed int, total time.Duration) *LatencyResult {
	if len(latencies) == 0 {
		return &LatencyResult{RequestsFailed: failed}
	}

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	var sum time.Duration
	min := latencies[0]
	max := latencies[0]
	for _, d := range latencies {
		sum += d
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}

	avg := sum / time.Duration(len(latencies))
	p50 := latencies[len(latencies)*50/100]
	p99 := latencies[len(latencies)*99/100]
	p999 := latencies[len(latencies)*999/1000]

	return &LatencyResult{
		RequestsSent:   len(latencies),
		RequestsFailed: failed,
		Latencies:      latencies,
		MinLatency:     min,
		MaxLatency:     max,
		AvgLatency:     avg,
		P50Latency:     p50,
		P99Latency:     p99,
		P999Latency:    p999,
		TotalDuration:  total,
	}
}

func ValidateLatency(r *LatencyResult) error {
	if r.P99Latency > 100*time.Millisecond {
		return fmt.Errorf("P99 latency %v exceeds 100ms threshold", r.P99Latency)
	}
	if r.AvgLatency > 50*time.Millisecond {
		return fmt.Errorf("average latency %v exceeds 50ms threshold", r.AvgLatency)
	}
	return nil
}

func BenchmarkLatency(b *testing.B) {
	bench := NewLatencyBenchmark("localhost", 8443)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bench.Run(ctx)
		if err != nil {
			b.Fatalf("latency benchmark failed: %v", err)
		}
	}
}
