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

// ConnectionFloodTest simulates 10K concurrent connections
type ConnectionFloodTest struct {
	TargetHost     string
	TargetPort     int
	MaxConnections int
	Duration       time.Duration
}

func NewConnectionFloodTest(host string, port int) *ConnectionFloodTest {
	return &ConnectionFloodTest{
		TargetHost:     host,
		TargetPort:     port,
		MaxConnections: 10000,
		Duration:       5 * time.Minute,
	}
}

func (c *ConnectionFloodTest) Run(ctx context.Context) (*FloodResult, error) {
	var activeConns int64
	var totalConns int64
	var failedConns int64
	var mu sync.Mutex
	errors := make([]error, 0)

	ctx, cancel := context.WithTimeout(ctx, c.Duration)
	defer cancel()

	var wg sync.WaitGroup
	sem := make(chan struct{}, c.MaxConnections)

	for i := 0; i < c.MaxConnections; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(id int) {
			defer wg.Done()
			defer func() { <-sem }()

			atomic.AddInt64(&activeConns, 1)
			atomic.AddInt64(&totalConns, 1)

			conn, err := net.DialTimeout("tcp",
				fmt.Sprintf("%s:%d", c.TargetHost, c.TargetPort),
				10*time.Second)
			if err != nil {
				atomic.AddInt64(&failedConns, 1)
				mu.Lock()
				errors = append(errors, fmt.Errorf("conn %d: %w", id, err))
				mu.Unlock()
				return
			}
			defer func() { _ = conn.Close() }()

			select {
			case <-ctx.Done():
			case <-time.After(c.Duration):
			}

			atomic.AddInt64(&activeConns, -1)
		}(i)
	}

	wg.Wait()

	return &FloodResult{
		TotalConnections:  totalConns,
		FailedConnections: failedConns,
		SuccessRate:       float64(totalConns-failedConns) / float64(totalConns) * 100,
		Errors:            errors,
		Duration:          c.Duration,
	}, nil
}

type FloodResult struct {
	TotalConnections  int64
	FailedConnections int64
	SuccessRate       float64
	Errors            []error
	Duration          time.Duration
}

func BenchmarkConnectionFlood(b *testing.B) {
	test := NewConnectionFloodTest("localhost", 8443)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := test.Run(ctx)
		if err != nil {
			b.Fatalf("flood test failed: %v", err)
		}
	}
}

func ValidateFloodResult(r *FloodResult) error {
	if r.SuccessRate < 99.0 {
		return fmt.Errorf("success rate %.2f%% below threshold 99%%", r.SuccessRate)
	}
	return nil
}
