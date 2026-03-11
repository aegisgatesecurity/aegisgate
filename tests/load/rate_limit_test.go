package load

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// RateLimitTest validates rate limiting under load
type RateLimitTest struct {
	TargetHost    string
	TargetPort    int
	BurstSize     int
	SustainedRate int
	Duration      time.Duration
}

func NewRateLimitTest(host string, port int) *RateLimitTest {
	return &RateLimitTest{
		TargetHost:    host,
		TargetPort:    port,
		BurstSize:     1000,
		SustainedRate: 100,
		Duration:      5 * time.Minute,
	}
}

type RateLimitResult struct {
	TotalRequests    int64
	AcceptedRequests int64
	RejectedRequests int64
	RateLimited      int64
	BurstPhase       PhaseResult
	SustainedPhase   PhaseResult
	Duration         time.Duration
}

type PhaseResult struct {
	Requests   int64
	Accepted   int64
	Rejected   int64
	Duration   time.Duration
	AvgLatency time.Duration
}

func (r *RateLimitTest) Run(ctx context.Context) (*RateLimitResult, error) {
	result := &RateLimitResult{}
	burstResult := r.runPhase(ctx, r.BurstSize, 0)
	result.BurstPhase = *burstResult
	result.TotalRequests += burstResult.Requests
	result.AcceptedRequests += burstResult.Accepted
	result.RejectedRequests += burstResult.Rejected

	time.Sleep(2 * time.Second)

	ctx, cancel := context.WithTimeout(ctx, r.Duration)
	defer cancel()

	sustainedResult := r.runSustainedPhase(ctx, r.SustainedRate)
	result.SustainedPhase = *sustainedResult
	result.TotalRequests += sustainedResult.Requests
	result.AcceptedRequests += sustainedResult.Accepted
	result.RejectedRequests += sustainedResult.Rejected
	result.RateLimited = sustainedResult.Rejected
	result.Duration = burstResult.Duration + sustainedResult.Duration

	return result, nil
}

func (r *RateLimitTest) runPhase(ctx context.Context, count int, delay time.Duration) *PhaseResult {
	var accepted, rejected int64
	var totalLatency int64
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reqStart := time.Now()
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", r.TargetHost, r.TargetPort), 5*time.Second)
			elapsed := time.Since(reqStart)
			atomic.AddInt64(&totalLatency, int64(elapsed))
			if err != nil {
				atomic.AddInt64(&rejected, 1)
				return
			}
			conn.Close()
			atomic.AddInt64(&accepted, 1)
		}()
		if delay > 0 {
			time.Sleep(delay)
		}
	}
	wg.Wait()
	duration := time.Since(start)
	total := accepted + rejected
	var avgLatency time.Duration
	if total > 0 {
		avgLatency = time.Duration(totalLatency / total)
	}
	return &PhaseResult{Requests: total, Accepted: accepted, Rejected: rejected, Duration: duration, AvgLatency: avgLatency}
}

func (r *RateLimitTest) runSustainedPhase(ctx context.Context, rate int) *PhaseResult {
	var accepted, rejected int64
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()
	start := time.Now()
	var wg sync.WaitGroup
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", r.TargetHost, r.TargetPort), 5*time.Second)
				if err != nil {
					atomic.AddInt64(&rejected, 1)
					return
				}
				conn.Close()
				atomic.AddInt64(&accepted, 1)
			}()
		}
	}
	wg.Wait()
	duration := time.Since(start)
	total := accepted + rejected
	return &PhaseResult{Requests: total, Accepted: accepted, Rejected: rejected, Duration: duration}
}

func ValidateRateLimit(r *RateLimitResult) error {
	if r.AcceptedRequests == 0 {
		return fmt.Errorf("no requests were accepted")
	}
	acceptRate := float64(r.AcceptedRequests) / float64(r.TotalRequests)
	if acceptRate < 0.5 {
		return fmt.Errorf("accept rate %.2f%% too low", acceptRate*100)
	}
	return nil
}

func BenchmarkRateLimit(b *testing.B) {
	test := NewRateLimitTest("localhost", 8443)
	test.Duration = 30 * time.Second
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := test.Run(ctx)
		if err != nil {
			b.Fatalf("rate limit test failed: %v", err)
		}
	}
}
