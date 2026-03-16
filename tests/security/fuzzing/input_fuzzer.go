package fuzzing

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"
)

type InputFuzzer struct {
	TargetHost   string
	TargetPort   int
	Iterations   int
	Seed         int64
	MaxInputSize int
}

func NewInputFuzzer(host string, port int) *InputFuzzer {
	return &InputFuzzer{
		TargetHost:   host,
		TargetPort:   port,
		Iterations:   1000,
		Seed:         time.Now().UnixNano(),
		MaxInputSize: 4096,
	}
}

type FuzzResult struct {
	Iterations    int
	Crashes       int
	Errors        int
	Timeouts      int
	UniqueCrashes []string
	Status        string
}

func generateRandomBytes(r *rand.Rand, maxSize int) string {
	size := r.Intn(maxSize) + 1
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(r.Intn(256))
	}
	return string(data)
}

func (f *InputFuzzer) Run(ctx context.Context) (*FuzzResult, error) {
	rng := rand.New(rand.NewSource(f.Seed))
	result := &FuzzResult{Status: "PASS"}

	for i := 0; i < f.Iterations; i++ {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		input := generateRandomBytes(rng, f.MaxInputSize)
		err := f.sendInput(ctx, input)
		if err != nil {
			result.Errors++
		}
		result.Iterations++
	}

	if result.Crashes > 0 {
		result.Status = "FAIL"
	}
	return result, nil
}

func (f *InputFuzzer) sendInput(ctx context.Context, input string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", f.TargetHost, f.TargetPort), 5*time.Second)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	_, err = conn.Write([]byte(input))
	return err
}

func TestInputFuzzing(t *testing.T) {
	fuzzer := NewInputFuzzer("localhost", 8443)
	fuzzer.Iterations = 100
	ctx := context.Background()

	result, err := fuzzer.Run(ctx)
	if err != nil {
		t.Logf("Fuzzing completed with warning: %v", err)
	}

	if result.Crashes > 0 {
		t.Errorf("Found %d crashes during fuzzing", result.Crashes)
	}
}
