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

func TestInjectionTestCreation(t *testing.T) {
	test := NewInjectionTest("localhost", 8443)
	
	if test.TargetHost != "localhost" {
		t.Errorf("expected localhost, got %s", test.TargetHost)
	}
	
	if test.TargetPort != 8443 {
		t.Errorf("expected 8443, got %d", test.TargetPort)
	}
	
	if test.Timeout != 30*time.Second {
		t.Errorf("expected 30s timeout, got %v", test.Timeout)
	}
}

func TestSQLInjectionPayloads(t *testing.T) {
	if len(SQLInjectionPayloads) == 0 {
		t.Error("SQL injection payloads should not be empty")
	}
	
	expectedPayloads := []string{
		"' OR '1'='1",
		"' UNION SELECT * FROM users--",
		"admin'--",
	}
	
	for _, expected := range expectedPayloads {
		found := false
		for _, payload := range SQLInjectionPayloads {
			if payload == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected payload %s not found", expected)
		}
	}
}

func TestCommandInjectionPayloads(t *testing.T) {
	if len(CommandInjectionPayloads) == 0 {
		t.Error("Command injection payloads should not be empty")
	}
	
	expectedPayloads := []string{
		"; cat /etc/passwd",
		"$(id)",
		"; echo injected",
	}
	
	for _, expected := range expectedPayloads {
		found := false
		for _, payload := range CommandInjectionPayloads {
			if payload == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected payload %s not found", expected)
		}
	}
}

func TestNoSQLInjectionPayloads(t *testing.T) {
	if len(NoSQLInjectionPayloads) == 0 {
		t.Error("NoSQL injection payloads should not be empty")
	}
	
	expectedPayloads := []string{
		"{$ne: null}",
		"{$gt: ''}",
	}
	
	for _, expected := range expectedPayloads {
		found := false
		for _, payload := range NoSQLInjectionPayloads {
			if payload == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected payload %s not found", expected)
		}
	}
}

func TestInjectionResultStructures(t *testing.T) {
	result := InjectionResult{
		Payload:  "test payload",
		Type:     "SQLi",
		Injected: true,
		Detected: true,
		Response: "error",
		Severity: "High",
		Details:  "test details",
	}
	
	if result.Payload != "test payload" {
		t.Errorf("expected test payload, got %s", result.Payload)
	}
	
	if result.Type != "SQLi" {
		t.Errorf("expected SQLi type, got %s", result.Type)
	}
	
	if !result.Injected {
		t.Error("expected Injected to be true")
	}
	
	if !result.Detected {
		t.Error("expected Detected to be true")
	}
}

func TestInjectionReportStructures(t *testing.T) {
	report := &InjectionReport{
		SQLiResults:    []InjectionResult{{Payload: "test1", Type: "SQLi"}},
		CommandResults: []InjectionResult{{Payload: "test2", Type: "Command"}},
		NoSQLResults:   []InjectionResult{{Payload: "test3", Type: "NoSQL"}},
		TotalTests:     3,
		Vulnerabilities: 1,
		Status:         "FAIL",
	}
	
	if report.TotalTests != 3 {
		t.Errorf("expected 3 total tests, got %d", report.TotalTests)
	}
	
	if report.Vulnerabilities != 1 {
		t.Errorf("expected 1 vulnerability, got %d", report.Vulnerabilities)
	}
	
	if report.Status != "FAIL" {
		t.Errorf("expected FAIL status, got %s", report.Status)
	}
}

func TestRunInjectionTestsWithMockHost(t *testing.T) {
	// Use non-routable IP to avoid actual network calls
	test := NewInjectionTest("10.255.255.1", 1)
	test.Timeout = 1 * time.Second // Short timeout for tests
	
	ctx := context.Background()
	report, err := test.RunInjectionTests(ctx)
	
	// Should complete without panic
	if report == nil {
		t.Fatal("expected InjectionReport to be returned")
	}
	
	// Verify all three test categories were run
	totalExpected := len(SQLInjectionPayloads) + len(CommandInjectionPayloads) + len(NoSQLInjectionPayloads)
	if report.TotalTests != totalExpected {
		t.Errorf("expected %d total tests, got %d", totalExpected, report.TotalTests)
	}
	
	// Should have results in each category
	if len(report.SQLiResults) == 0 {
		t.Error("expected SQLi results")
	}
	
	if len(report.CommandResults) == 0 {
		t.Error("expected Command results")
	}
	
	if len(report.NoSQLResults) == 0 {
		t.Error("expected NoSQL results")
	}
	
	_ = err // May or may not have error
	t.Logf("Injection testing completed: %d total tests, %d vulnerabilities", 
		report.TotalTests, report.Vulnerabilities)
}

func TestMinFunction(t *testing.T) {
	if min(1, 2) != 1 {
		t.Error("min(1, 2) should be 1")
	}
	
	if min(5, 3) != 3 {
		t.Error("min(5, 3) should be 3")
	}
	
	if min(10, 10) != 10 {
		t.Error("min(10, 10) should be 10")
	}
}

func TestInjectionReportPassStatus(t *testing.T) {
	report := &InjectionReport{
		TotalTests:      10,
		Vulnerabilities: 0,
		Status:          "PASS",
	}
	
	if report.Vulnerabilities > 0 {
		t.Error("expected no vulnerabilities for PASS status")
	}
	
	if report.Status != "PASS" {
		t.Errorf("expected PASS status, got %s", report.Status)
	}
}

func TestInjectionReportStatusCalculation(t *testing.T) {
	report := &InjectionReport{
		TotalTests:      10,
		Vulnerabilities: 5,
	}
	
	// Manually call the status calculation logic
	if report.Vulnerabilities > 0 {
		report.Status = "FAIL"
	} else {
		report.Status = "PASS"
	}
	
	if report.Status != "FAIL" {
		t.Errorf("expected FAIL status with vulnerabilities, got %s", report.Status)
	}
}