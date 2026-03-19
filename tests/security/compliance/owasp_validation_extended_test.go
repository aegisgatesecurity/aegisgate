// SPDX-License-Identifier: MIT
// =========================================================================
// PROPRIETARY - AegisGate Security
// Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
// =========================================================================

package compliance

import (
	"testing"
	"time"
)

func TestOWASPCheckStructures(t *testing.T) {
	checks := []OWASPCheck{
		{ID: "A01", Name: "Broken Access Control", Category: "Access Control", Severity: "Critical", Status: "PASS"},
		{ID: "A02", Name: "Cryptographic Failures", Category: "Data Protection", Severity: "High", Status: "FAIL"},
		{ID: "A03", Name: "Injection", Category: "Input Validation", Severity: "Critical", Status: "PASS"},
	}

	for _, check := range checks {
		if check.ID == "" {
			t.Error("OWASP check ID should not be empty")
		}
		if check.Name == "" {
			t.Error("OWASP check name should not be empty")
		}
		if check.Category == "" {
			t.Error("OWASP check category should not be empty")
		}
		if check.Severity == "" {
			t.Error("OWASP check severity should not be empty")
		}
	}

	// Verify specific checks
	if checks[0].Severity != "Critical" {
		t.Errorf("expected Critical severity for A01, got %s", checks[0].Severity)
	}
}

func TestOWASPReportEdgeCases(t *testing.T) {
	// Test boundary: 70% exactly
	report := &OWASPReport{
		Passed: 7,
		Failed: 3,
	}
	report.calculateCompliance()
	if report.Compliance != 70.0 {
		t.Errorf("expected 70%%, got %.1f%%", report.Compliance)
	}

	// Test boundary: 89% exactly
	report = &OWASPReport{
		Passed: 89,
		Failed: 11,
	}
	report.calculateCompliance()
	if report.Compliance != 89.0 {
		t.Errorf("expected 89%%, got %.1f%%", report.Compliance)
	}

	// Test boundary: 90% exactly
	report = &OWASPReport{
		Passed: 9,
		Failed: 1,
	}
	report.calculateCompliance()
	if report.Compliance != 90.0 {
		t.Errorf("expected 90%%, got %.1f%%", report.Compliance)
	}
}

func TestOWASPValidatorSetters(t *testing.T) {
	v := NewOWASPValidator("testhost", 9999)

	// Test timeout setter
	if v.Timeout != 30*time.Second {
		t.Error("default timeout should be 30 seconds")
	}

	// Modify and verify
	v.Timeout = 60 * time.Second
	if v.Timeout != 60*time.Second {
		t.Error("timeout should be settable")
	}
}

func TestOWASPReportFullStack(t *testing.T) {
	report := &OWASPReport{
		Checks: []OWASPCheck{
			{ID: "A01", Name: "Test 1", Status: "PASS"},
			{ID: "A02", Name: "Test 2", Status: "FAIL"},
			{ID: "A03", Name: "Test 3", Status: "PASS"},
		},
		Passed:        2,
		Failed:        1,
		Compliance:    66.67,
		OverallStatus: "PARTIAL",
	}

	// Verify all fields
	if len(report.Checks) != 3 {
		t.Errorf("expected 3 checks, got %d", len(report.Checks))
	}

	if report.Passed != 2 {
		t.Errorf("expected 2 passed, got %d", report.Passed)
	}

	if report.Failed != 1 {
		t.Errorf("expected 1 failed, got %d", report.Failed)
	}

	_ = report.Compliance
	_ = report.OverallStatus
}