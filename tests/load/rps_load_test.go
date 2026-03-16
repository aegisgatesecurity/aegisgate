// Copyright 2024 AegisGate
// Comprehensive RPS Load Testing Framework
// Supports: 10k, 25k, 50k RPS with detailed metrics
//
// Run with: go test -bench=. -benchmem ./tests/load/...

//go:build !integration
// +build !integration

package load

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ============================================================================
// RPS Load Test Result - String Method
// ============================================================================

// String returns a human-readable summary of the result
func (r *RPSResult) String() string {
	return fmt.Sprintf(` ========================================
RPS Load Test Results: %s RPS
========================================
Configuration:
  Target RPS:     %d
  Duration:       %v
  Concurrency:    %d
  Clients:        %d

Throughput:
  Actual RPS:     %.2f
  Achievement:     %.1f%%

Requests:
  Total:          %d
  Success:        %d (%.2f%%)
  Errors:         %d (%.2f%%)
  Timeouts:       %d
  Conn Errors:    %d

Latency:
  Min:            %v
  Avg:            %v
  P50:            %v
  P90:            %v
  P95:            %v
  P99:            %v
  P999:           %v
  Max:            %v

Latency Distribution:
  <1ms:           %d
  1-5ms:          %d
  5-10ms:         %d
  10-50ms:        %d
  50-100ms:       %d
  100-500ms:      %d
  >500ms:         %d

Resources:
  Avg Memory:     %.2f MB
  Peak Memory:    %.2f MB
  Avg CPU:        %.2f%%
  GCs:            %d

Scaling Analysis Ready: %v
=======================================`,
		r.Config.Level,
		r.Config.TargetRPS,
		r.Config.Duration,
		r.Config.Concurrency,
		r.Config.NumClients,
		r.ActualRPS,
		r.RPSAchievement,
		r.TotalRequests,
		r.SuccessCount,
		float64(r.SuccessCount)/float64(r.TotalRequests)*100,
		r.ErrorCount,
		float64(r.ErrorCount)/float64(r.TotalRequests)*100,
		r.TimeoutCount,
		r.ConnectionError,
		r.MinLatency,
		r.AvgLatency,
		r.P50Latency,
		r.P90Latency,
		r.P95Latency,
		r.P99Latency,
		r.P999Latency,
		r.MaxLatency,
		r.LatencyBuckets["<1ms"],
		r.LatencyBuckets["1-5ms"],
		r.LatencyBuckets["5-10ms"],
		r.LatencyBuckets["10-50ms"],
		r.LatencyBuckets["50-100ms"],
		r.LatencyBuckets["100-500ms"],
		r.LatencyBuckets[">500ms"],
		r.AvgMemoryMB,
		r.PeakMemoryMB,
		r.AvgCPUPercent,
		r.NumGC,
		len(r.Samples) > 0,
	)
}

// ============================================================================
// RPS Load Test Runner
// ============================================================================

// RPSLoadTester is the main load testing engine
type RPSLoadTester struct {
	config RPSConfig
	client *http.Client
	stats  *RPSStats
}

// RPSStats holds running statistics during the test
type RPSStats struct {
	// Counters
	totalRequests  int64
	successCount   int64
	errorCount     int64
	timeoutCount   int64
	connErrorCount int64

	// Latency tracking
	latencies   []time.Duration
	latenciesMu sync.Mutex

	// Error tracking
	errorsMu    sync.Mutex
	errorCounts map[string]int64

	// Time series
	samplesMu sync.Mutex
	samples   []RPSTimeSeriesSample

	// Active connections
	activeConns int32
}

// NewRPSLoadTester creates a new RPS load tester
func NewRPSLoadTester(config RPSConfig) *RPSLoadTester {
	// Create HTTP client with connection pooling
	transport := &http.Transport{
		MaxIdleConns:        config.NumClients,
		MaxIdleConnsPerHost: config.NumClients / 4,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  false,
		DisableKeepAlives:   !config.KeepAlive,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.RequestTimeout,
	}

	return &RPSLoadTester{
		config: config,
		client: client,
		stats: &RPSStats{
			latencies:   make([]time.Duration, 0, 1000000),
			errorCounts: make(map[string]int64),
		},
	}
}

// Run executes the RPS load test
func (t *RPSLoadTester) Run(ctx context.Context) (*RPSResult, error) {
	result := &RPSResult{
		Config:         t.config,
		StartTime:      time.Now(),
		LatencyBuckets: make(map[string]int64),
		ErrorBreakdown: make(map[string]int64),
	}

	// Initialize latency buckets
	for _, b := range LatencyBuckets() {
		result.LatencyBuckets[b] = 0
	}

	// Start resource monitoring
	resourceChan := make(chan struct{})
	var resourcewg sync.WaitGroup
	resourcewg.Add(1)
	go t.monitorResources(ctx, resourceChan, &resourcewg)

	// Start the test
	fmt.Printf("Starting RPS load test: %s RPS\n", t.config.Level)
	fmt.Printf("Configuration: %d concurrent, %d clients, %v duration\n",
		t.config.Concurrency, t.config.NumClients, t.config.Duration)

	// Warmup phase
	fmt.Printf("Warmup phase: %v\n", t.config.WarmupDuration)
	time.Sleep(t.config.WarmupDuration)

	// Main test phase
	workerChan := make(chan struct{}, t.config.Concurrency)

	var wg sync.WaitGroup
	requestCounter := int64(0)
	targetRequests := int64(t.config.TargetRPS * int(t.config.Duration.Seconds()))

	// Create worker pool
	for i := 0; i < t.config.Concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case workerChan <- struct{}{}:
					// Process request
					atomic.AddInt64(&requestCounter, 1)
					t.makeRequest(ctx, result)
					<-workerChan

					// Check if we've sent enough requests
					if atomic.LoadInt64(&requestCounter) >= targetRequests {
						return
					}
				}
			}
		}(i)
	}

	// Wait for test duration or target requests
	testDuration := t.config.Duration
	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All workers completed")
	case <-time.After(testDuration):
		fmt.Println("Duration limit reached")
	}

	result.EndTime = time.Now()
	result.TotalDuration = result.EndTime.Sub(result.StartTime)

	// Stop resource monitoring
	close(resourceChan)
	resourcewg.Wait()

	// Calculate final metrics
	t.calculateMetrics(result)

	return result, nil
}

// makeRequest makes a single HTTP request and records metrics
func (t *RPSLoadTester) makeRequest(ctx context.Context, result *RPSResult) {
	atomic.AddInt64(&t.stats.totalRequests, 1)
	atomic.AddInt32(&t.stats.activeConns, 1)
	defer atomic.AddInt32(&t.stats.activeConns, -1)

	start := time.Now()

	// Create request
	url := fmt.Sprintf("http://%s:%d/health", t.config.TargetHost, t.config.TargetPort)
	if t.config.UseHTTPS {
		url = fmt.Sprintf("https://%s:%d/health", t.config.TargetHost, t.config.TargetPort)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		atomic.AddInt64(&t.stats.connErrorCount, 1)
		atomic.AddInt64(&t.stats.errorCount, 1)
		t.stats.errorsMu.Lock()
		t.stats.errorCounts["request_create_error"]++
		t.stats.errorsMu.Unlock()
		return
	}

	// Make request
	resp, err := t.client.Do(req)
	latency := time.Since(start)

	if err != nil {
		atomic.AddInt64(&t.stats.errorCount, 1)
		t.stats.errorsMu.Lock()

		if err == context.DeadlineExceeded {
			atomic.AddInt64(&t.stats.timeoutCount, 1)
			t.stats.errorCounts["timeout"]++
		} else {
			atomic.AddInt64(&t.stats.connErrorCount, 1)
			t.stats.errorCounts["connection_error"]++
		}
		t.stats.errorsMu.Unlock()
	} else {
		_ = resp.Body.Close()
		atomic.AddInt64(&t.stats.successCount, 1)
	}

	// Record latency
	t.stats.latenciesMu.Lock()
	t.stats.latencies = append(t.stats.latencies, latency)
	t.stats.latenciesMu.Unlock()

	// Update latency bucket
	t.updateLatencyBucket(latency, result)
}

// updateLatencyBucket categorizes latency into buckets
func (t *RPSLoadTester) updateLatencyBucket(latency time.Duration, result *RPSResult) {
	var bucket string
	switch {
	case latency < time.Millisecond:
		bucket = "<1ms"
	case latency < 5*time.Millisecond:
		bucket = "1-5ms"
	case latency < 10*time.Millisecond:
		bucket = "5-10ms"
	case latency < 50*time.Millisecond:
		bucket = "10-50ms"
	case latency < 100*time.Millisecond:
		bucket = "50-100ms"
	case latency < 500*time.Millisecond:
		bucket = "100-500ms"
	default:
		bucket = ">500ms"
	}
	result.LatencyBuckets[bucket]++
}

// monitorResources monitors CPU and memory usage during the test
func (t *RPSLoadTester) monitorResources(ctx context.Context, stopChan <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var totalMem float64
	var samples int
	var peakMem uint64

	for {
		select {
		case <-ctx.Done():
			return
		case <-stopChan:
			return
		case <-ticker.C:
			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)

			memMB := float64(ms.HeapAlloc) / 1024 / 1024
			if ms.HeapAlloc > peakMem {
				peakMem = ms.HeapAlloc
			}
			totalMem += memMB
			samples++

			sample := RPSTimeSeriesSample{
				Timestamp:   time.Now(),
				MemoryMB:    memMB,
				ActiveConns: int(atomic.LoadInt32(&t.stats.activeConns)),
			}

			t.stats.samplesMu.Lock()
			t.stats.samples = append(t.stats.samples, sample)
			t.stats.samplesMu.Unlock()
		}
	}
}

// calculateMetrics computes final metrics from collected data
func (t *RPSLoadTester) calculateMetrics(result *RPSResult) {
	// Basic counts
	result.TotalRequests = atomic.LoadInt64(&t.stats.totalRequests)
	result.SuccessCount = atomic.LoadInt64(&t.stats.successCount)
	result.ErrorCount = atomic.LoadInt64(&t.stats.errorCount)
	result.TimeoutCount = atomic.LoadInt64(&t.stats.timeoutCount)
	result.ConnectionError = atomic.LoadInt64(&t.stats.connErrorCount)

	// Throughput
	result.TotalDuration = result.Config.Duration
	result.ActualRPS = float64(result.TotalRequests) / result.TotalDuration.Seconds()
	result.TargetRPS = float64(result.Config.TargetRPS)
	result.RPSAchievement = (result.ActualRPS / result.TargetRPS) * 100

	// Latency calculations
	t.stats.latenciesMu.Lock()
	latencies := t.stats.latencies
	t.stats.latenciesMu.Unlock()

	if len(latencies) > 0 {
		// Sort for percentile calculations
		sorted := make([]time.Duration, len(latencies))
		copy(sorted, latencies)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

		result.MinLatency = sorted[0]
		result.MaxLatency = sorted[len(sorted)-1]

		var totalLatency int64
		for _, l := range sorted {
			totalLatency += int64(l)
		}
		result.AvgLatency = time.Duration(totalLatency / int64(len(sorted)))

		// Percentiles
		result.P50Latency = sorted[int(float64(len(sorted))*0.50)]
		result.P90Latency = sorted[int(float64(len(sorted))*0.90)]
		result.P95Latency = sorted[int(float64(len(sorted))*0.95)]
		result.P99Latency = sorted[int(float64(len(sorted))*0.99)]
		result.P999Latency = sorted[int(float64(len(sorted))*0.999)]
		if len(sorted) >= 10000 {
			result.P9999Latency = sorted[int(float64(len(sorted))*0.9999)]
		} else {
			result.P9999Latency = result.P999Latency
		}
	}

	// Error breakdown
	result.ErrorBreakdown = make(map[string]int64)
	t.stats.errorsMu.Lock()
	for k, v := range t.stats.errorCounts {
		result.ErrorBreakdown[k] = v
	}
	t.stats.errorsMu.Unlock()

	// Resource metrics
	t.stats.samplesMu.Lock()
	samples := t.stats.samples
	t.stats.samplesMu.Unlock()

	if len(samples) > 0 {
		var totalMem float64
		for _, s := range samples {
			totalMem += s.MemoryMB
		}
		result.AvgMemoryMB = totalMem / float64(len(samples))

		// Find peak memory
		for _, s := range samples {
			if s.MemoryMB > result.PeakMemoryMB {
				result.PeakMemoryMB = s.MemoryMB
			}
		}
		result.Samples = samples
	}

	// GC stats
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	result.NumGC = ms.NumGC
}

// ============================================================================
// Benchmark Functions
// ============================================================================

// BenchmarkRPS10K runs a 10k RPS load test
func BenchmarkRPS10K(b *testing.B) {
	config := DefaultRPSConfigs()[RPS10K]
	config.Duration = 10 * time.Second // Shortened for benchmark
	config.WarmupDuration = 2 * time.Second

	tester := NewRPSLoadTester(config)
	ctx := context.Background()

	b.ResetTimer()
	result, err := tester.Run(ctx)
	if err != nil {
		b.Fatalf("RPS load test failed: %v", err)
	}

	b.Logf("\n%s", result.String())
}

// BenchmarkRPS25K runs a 25k RPS load test
func BenchmarkRPS25K(b *testing.B) {
	config := DefaultRPSConfigs()[RPS25K]
	config.Duration = 10 * time.Second
	config.WarmupDuration = 2 * time.Second

	tester := NewRPSLoadTester(config)
	ctx := context.Background()

	b.ResetTimer()
	result, err := tester.Run(ctx)
	if err != nil {
		b.Fatalf("RPS load test failed: %v", err)
	}

	b.Logf("\n%s", result.String())
}

// BenchmarkRPS50K runs a 50k RPS load test
func BenchmarkRPS50K(b *testing.B) {
	config := DefaultRPSConfigs()[RPS50K]
	config.Duration = 10 * time.Second
	config.WarmupDuration = 2 * time.Second

	tester := NewRPSLoadTester(config)
	ctx := context.Background()

	b.ResetTimer()
	result, err := tester.Run(ctx)
	if err != nil {
		b.Fatalf("RPS load test failed: %v", err)
	}

	b.Logf("\n%s", result.String())
}

// BenchmarkRPSAll runs all RPS levels sequentially
func BenchmarkRPSAll(b *testing.B) {
	levels := []RPSLevel{RPS10K, RPS25K, RPS50K}

	for _, level := range levels {
		b.Run(fmt.Sprintf("RPS_%s", level), func(b *testing.B) {
			config := DefaultRPSConfigs()[level]
			config.Duration = 10 * time.Second
			config.WarmupDuration = 2 * time.Second

			tester := NewRPSLoadTester(config)
			ctx := context.Background()

			b.ResetTimer()
			result, err := tester.Run(ctx)
			if err != nil {
				b.Fatalf("RPS load test failed: %v", err)
			}

			b.Logf("\n%s", result.String())
		})
	}
}

// ============================================================================
// Stress Testing Functions
// ============================================================================

// StressTestRPS progressively increases RPS to find breaking point
func StressTestRPS(startRPS, maxRPS, stepRPS int, duration time.Duration) ([]*RPSResult, error) {
	results := []*RPSResult{}

	for rps := startRPS; rps <= maxRPS; rps += stepRPS {
		fmt.Printf("\n=== Testing %d RPS ===\n", rps)

		config := DefaultRPSConfigs()[RPS10K]
		config.TargetRPS = rps
		config.Concurrency = rps / 20 // Approximate concurrency for RPS
		if config.Concurrency < 10 {
			config.Concurrency = 10
		}
		config.Duration = duration
		config.WarmupDuration = 2 * time.Second

		tester := NewRPSLoadTester(config)
		ctx := context.Background()

		result, err := tester.Run(ctx)
		if err != nil {
			return results, err
		}

		results = append(results, result)

		// Check if error rate is too high
		errorRate := float64(result.ErrorCount) / float64(result.TotalRequests)
		if errorRate > 0.05 { // 5% error rate threshold
			fmt.Printf("Error rate exceeded threshold at %d RPS: %.2f%%\n", rps, errorRate*100)
			break
		}

		// Check if latency p99 is too high
		if result.P99Latency > 500*time.Millisecond {
			fmt.Printf("P99 latency exceeded threshold at %d RPS: %v\n", rps, result.P99Latency)
			break
		}
	}

	return results, nil
}

// ============================================================================
// Scaling Analysis Functions
// ============================================================================

// AnalyzeScalingResults analyzes results from multiple RPS levels
func AnalyzeScalingResults(results []*RPSResult) string {
	if len(results) < 2 {
		return "Need at least 2 results to analyze scaling"
	}

	var analysis string
	analysis += "\n========================================\n"
	analysis += "Auto-Scaling Analysis\n"
	analysis += "========================================\n\n"

	// Calculate scaling efficiency
	for i := 1; i < len(results); i++ {
		prev := results[i-1]
		curr := results[i]

		rpsIncrease := float64(curr.Config.TargetRPS) / float64(prev.Config.TargetRPS)
		latencyIncrease := float64(curr.P99Latency) / float64(prev.P99Latency)
		memoryIncrease := curr.PeakMemoryMB / prev.PeakMemoryMB

		analysis += "Scale from " + prev.Config.Level.String() + " → " + curr.Config.Level.String() + " RPS:\n"
		analysis += fmt.Sprintf("  RPS increase:     %.1fx\n", rpsIncrease)
		analysis += fmt.Sprintf("  P99 latency:      %.2fx (%v → %v)\n",
			latencyIncrease, prev.P99Latency, curr.P99Latency)
		analysis += fmt.Sprintf("  Memory increase:  %.2fx (%.1fMB → %.1fMB)\n",
			memoryIncrease, prev.PeakMemoryMB, curr.PeakMemoryMB)

		// Scaling recommendation
		if latencyIncrease > 2.0 {
			analysis += "  ⚠️  Consider scaling: latency increased significantly\n"
		} else if latencyIncrease < 1.5 && memoryIncrease < 1.5 {
			analysis += "  ✅ Good scaling: linear performance\n"
		}

		analysis += "\n"
	}

	// Auto-scaling recommendations
	analysis += "Auto-Scaling Recommendations:\n"
	analysis += "-----------------------------\n"

	// Find the RPS level where latency starts degrading
	optimalRPS := results[0].Config.TargetRPS
	for i := 1; i < len(results); i++ {
		prevP99 := results[i-1].P99Latency
		currP99 := results[i].P99Latency

		if currP99 > prevP99*2 {
			optimalRPS = results[i-1].Config.TargetRPS
			break
		}
		optimalRPS = results[i].Config.TargetRPS
	}

	analysis += fmt.Sprintf("  Recommended trigger threshold: %d RPS\n", optimalRPS)
	analysis += "  Target P99 latency: <100ms\n"
	analysis += "  Scale-up: +25%% when RPS > threshold\n"
	analysis += "  Scale-down: -25%% when RPS < threshold * 0.7\n"

	return analysis
}

// ============================================================================
// Mock Server for Testing
// ============================================================================

// MockRPSTarget creates a simple HTTP server for testing
type MockRPSTarget struct {
	Server       *http.Server
	RequestCount int64
	LatencySum   time.Duration
	mu           sync.Mutex
}

func NewMockRPSTarget(port int) *MockRPSTarget {
	m := &MockRPSTarget{}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			elapsed := time.Since(start)
			atomic.AddInt64(&m.RequestCount, 1)
			m.mu.Lock()
			m.LatencySum += elapsed
			m.mu.Unlock()
		}()

		// Simulate minimal processing
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	m.Server = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return m
}

func (m *MockRPSTarget) Start() error {
	return m.Server.ListenAndServe()
}

func (m *MockRPSTarget) Stop() error {
	return m.Server.Close()
}

func (m *MockRPSTarget) Stats() (int64, time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.RequestCount, m.LatencySum
}
