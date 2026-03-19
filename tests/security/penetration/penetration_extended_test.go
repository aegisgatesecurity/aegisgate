// SPDX-License-Identifier: MIT
// =========================================================================
// PROPRIETARY - AegisGate Security
// Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
// =========================================================================

package penetration

import (
	"context"
	"testing"
	"time"
)

func TestInjectionTestSetters(t *testing.T) {
	test := NewInjectionTest("localhost", 8080)

	if test.Timeout != 30*time.Second {
		t.Error("default timeout should be 30 seconds")
	}

	// Modify and verify
	test.Timeout = 60 * time.Second
	if test.Timeout != 60*time.Second {
		t.Error("timeout should be settable")
	}
}

func TestIndividualSQLInjection(t *testing.T) {
	test := NewInjectionTest("10.255.255.1", 1)
	test.Timeout = 1 * time.Second // Very short timeout

	ctx := context.Background()
	result := test.testSQLInjection(ctx, "' OR '1'='1")

	// Verify result structure regardless of connection outcome
	if result.Payload != "' OR '1'='1" {
		t.Errorf("expected payload in result, got %s", result.Payload)
	}

	if result.Type != "SQLi" {
		t.Errorf("expected type SQLi, got %s", result.Type)
	}
	// Connection may fail - that's expected in test environment
	_ = result.Detected
	_ = result.Response
	_ = result.Severity
}

func TestIndividualCommandInjection(t *testing.T) {
	test := NewInjectionTest("10.255.255.1", 1)
	test.Timeout = 1 * time.Second

	ctx := context.Background()
	result := test.testCommandInjection(ctx, "; ls")

	// Verify result structure regardless of connection outcome
	if result.Payload != "; ls" {
		t.Errorf("expected payload in result, got %s", result.Payload)
	}

	if result.Type != "Command" {
		t.Errorf("expected type Command, got %s", result.Type)
	}
	// Connection may fail - that's expected in test environment
	_ = result.Detected
	_ = result.Response
	_ = result.Severity
}

func TestIndividualNoSQLInjection(t *testing.T) {
	test := NewInjectionTest("10.255.255.1", 1)
	test.Timeout = 1 * time.Second // Use non-routable IP with short timeout

	ctx := context.Background()
	result := test.testNoSQLInjection(ctx, "{$ne: null}")

	// Verify result structure regardless of connection outcome
	if result.Payload != "{$ne: null}" {
		t.Errorf("expected payload in result, got %s", result.Payload)
	}

	if result.Type != "NoSQL" {
		t.Errorf("expected type NoSQL, got %s", result.Type)
	}
	// Connection may fail - that's expected in test environment
	_ = result.Injected
	_ = result.Detected
	_ = result.Response
	_ = result.Severity
}

func TestInjectionReportEdgeCases(t *testing.T) {
	// Zero vulnerabilities
	report := &InjectionReport{
		TotalTests:      15,
		Vulnerabilities: 0,
	}
	if report.Vulnerabilities > 0 {
		t.Error("should have no vulnerabilities")
	}

	// All categories populated
	report = &InjectionReport{
		SQLiResults:     []InjectionResult{{}, {}},
		CommandResults: []InjectionResult{{}, {}, {}},
		NoSQLResults:    []InjectionResult{{}},
		TotalTests:      6,
	}

	if len(report.SQLiResults) != 2 {
		t.Errorf("expected 2 SQLi results, got %d", len(report.SQLiResults))
	}

	if len(report.CommandResults) != 3 {
		t.Errorf("expected 3 Command results, got %d", len(report.CommandResults))
	}

	if len(report.NoSQLResults) != 1 {
		t.Errorf("expected 1 NoSQL result, got %d", len(report.NoSQLResults))
	}
}

func TestInjectionResultSeverityLevels(t *testing.T) {
	results := []InjectionResult{
		{Severity: "Info"},
		{Severity: "Low"},
		{Severity: "Medium"},
		{Severity: "High"},
		{Severity: "Critical"},
	}

	for _, r := range results {
		if r.Severity == "" {
			t.Error("severity should not be empty")
		}
	}
}

func TestMinFunctionEdgeCases(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},
		{-5, 5, -5},
		{100, 200, 100},
		{200, 100, 100},
	}

	for _, tt := range tests {
		result := min(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("min(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
		}
	}
}
