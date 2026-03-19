// SPDX-License-Identifier: MIT
// =========================================================================
// PROPRIETARY - AegisGate Security
// Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
// =========================================================================

package fuzzing

import (
	"context"
	"math/rand"
	"testing"
)

func TestInputFuzzerCreation(t *testing.T) {
	fuzzer := NewInputFuzzer("localhost", 8443)
	
	if fuzzer.TargetHost != "localhost" {
		t.Errorf("expected localhost, got %s", fuzzer.TargetHost)
	}
	
	if fuzzer.TargetPort != 8443 {
		t.Errorf("expected 8443, got %d", fuzzer.TargetPort)
	}
	
	if fuzzer.Iterations != 1000 {
		t.Errorf("expected 1000 iterations, got %d", fuzzer.Iterations)
	}
	
	if fuzzer.MaxInputSize != 4096 {
		t.Errorf("expected max input size 4096, got %d", fuzzer.MaxInputSize)
	}
}

func TestInputFuzzerCustomSettings(t *testing.T) {
	fuzzer := NewInputFuzzer("target.example.com", 9000)
	fuzzer.Iterations = 500
	fuzzer.Seed = 12345
	
	if fuzzer.TargetPort != 9000 {
		t.Errorf("expected port 9000, got %d", fuzzer.TargetPort)
	}
	
	if fuzzer.Iterations != 500 {
		t.Errorf("expected 500 iterations, got %d", fuzzer.Iterations)
	}
	
	if fuzzer.Seed != 12345 {
		t.Errorf("expected seed 12345, got %d", fuzzer.Seed)
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	// Test random byte generation
	rng := newRand(42) // Fixed seed for reproducibility
	
	for size := 1; size <= 100; size++ {
		result := generateRandomBytes(rng, size)
		
		if len(result) == 0 {
			t.Errorf("generated empty string for size %d", size)
		}
		
		if len(result) > size {
			t.Errorf("generated string longer than max: %d > %d", len(result), size)
		}
	}
}

func TestGenerateRandomBytesConsistency(t *testing.T) {
	// Same seed should produce same results
	rng1 := newRand(999)
	rng2 := newRand(999)
	
	result1 := generateRandomBytes(rng1, 100)
	result2 := generateRandomBytes(rng2, 100)
	
	if result1 != result2 {
		t.Error("same seed should produce identical random sequences")
	}
}

func TestFuzzerRunBasic(t *testing.T) {
	// Use non-routable IP to avoid actual network calls
	fuzzer := NewInputFuzzer("10.255.255.1", 1)
	fuzzer.Iterations = 10 // Small number for test
	fuzzer.MaxInputSize = 100
	
	ctx := context.Background()
	result, err := fuzzer.Run(ctx)
	
	// Should complete without panic
	if result == nil {
		t.Fatal("expected FuzzResult to be returned")
	}
	
	if result.Iterations != 10 {
		t.Errorf("expected 10 iterations, got %d", result.Iterations)
	}
	
	// Error is expected since host is unreachable
	t.Logf("Fuzzing completed: %d iterations, %d errors, status: %s", 
		result.Iterations, result.Errors, result.Status)
	_ = err // May or may not have error
}

func TestFuzzResultStructures(t *testing.T) {
	result := &FuzzResult{
		Iterations:    100,
		Crashes:       5,
		Errors:        10,
		Timeouts:      2,
		UniqueCrashes: []string{"crash1", "crash2"},
		Status:        "FAIL",
	}
	
	if result.Iterations != 100 {
		t.Errorf("expected 100 iterations, got %d", result.Iterations)
	}
	
	if result.Crashes != 5 {
		t.Errorf("expected 5 crashes, got %d", result.Crashes)
	}
	
	if result.Status != "FAIL" {
		t.Errorf("expected FAIL status, got %s", result.Status)
	}
}

func TestFuzzResultPassStatus(t *testing.T) {
	result := &FuzzResult{
		Iterations: 100,
		Crashes:   0,
		Status:    "PASS",
	}
	
	if result.Crashes > 0 {
		t.Error("expected no crashes for PASS status")
	}
	
	if result.Status != "PASS" {
		t.Errorf("expected PASS status, got %s", result.Status)
	}
}

// newRand creates a rand.Rand with the given seed
func newRand(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}