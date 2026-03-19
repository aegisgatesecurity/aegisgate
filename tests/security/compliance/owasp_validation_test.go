// SPDX-License-Identifier: MIT
// =========================================================================
// PROPRIETARY - AegisGate Security
// Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
// =========================================================================

package compliance

import (
	"context"
	"testing"
	"time"
)

func TestOWASPValidatorCreation(t *testing.T) {
	validator := NewOWASPValidator("localhost", 8443)
	
	if validator.TargetHost != "localhost" {
		t.Errorf("expected localhost, got %s", validator.TargetHost)
	}
	
	if validator.TargetPort != 8443 {
		t.Errorf("expected 8443, got %d", validator.TargetPort)
	}
	
	if validator.Timeout != 30*time.Second {
		t.Errorf("expected 30s timeout, got %v", validator.Timeout)
	}
}

func TestGetOWASPChecks(t *testing.T) {
	validator := NewOWASPValidator("localhost", 8443)
	checks := validator.GetOWASPChecks()
	
	if len(checks) != 10 {
		t.Errorf("expected 10 OWASP checks, got %d", len(checks))
	}
	
	// Verify critical categories are present
	criticalIDs := map[string]bool{
		"A01": false, // Broken Access Control
		"A03": false, // Injection
		"A07": false, // Authentication Failures
	}
	
	for _, check := range checks {
		if _, ok := criticalIDs[check.ID]; ok {
			criticalIDs[check.ID] = true
		}
	}
	
	for id, found := range criticalIDs {
		if !found {
			t.Errorf("missing critical check %s", id)
		}
	}
}

func TestOWASPReportCalculation(t *testing.T) {
	report := &OWASPReport{
		Passed: 9,
		Failed: 1,
	}
	
	report.calculateCompliance()
	
	if report.Compliance != 90.0 {
		t.Errorf("expected 90%% compliance, got %.1f%%", report.Compliance)
	}
	
	if report.OverallStatus != "COMPLIANT" {
		t.Errorf("expected COMPLIANT status (90%%+), got %s", report.OverallStatus)
	}
}

func TestOWASPReportPartial(t *testing.T) {
	report := &OWASPReport{
		Passed: 7,
		Failed: 3,
	}
	
	report.calculateCompliance()
	
	if report.Compliance != 70.0 {
		t.Errorf("expected 70%% compliance, got %.1f%%", report.Compliance)
	}
	
	if report.OverallStatus != "PARTIAL" {
		t.Errorf("expected PARTIAL status (70-89%%), got %s", report.OverallStatus)
	}
}

func TestOWASPReportNonCompliant(t *testing.T) {
	report := &OWASPReport{
		Passed: 2,
		Failed: 8,
	}
	
	report.calculateCompliance()
	
	if report.Compliance != 20.0 {
		t.Errorf("expected 20%% compliance, got %.1f%%", report.Compliance)
	}
	
	if report.OverallStatus != "NON-COMPLIANT" {
		t.Errorf("expected NON-COMPLIANT status, got %s", report.OverallStatus)
	}
}

func TestOWASPReportZero(t *testing.T) {
	report := &OWASPReport{}
	report.calculateCompliance()
	
	if report.Compliance != 0 {
		t.Errorf("expected 0%% compliance, got %.1f%%", report.Compliance)
	}
	
	if report.OverallStatus != "N/A" {
		t.Errorf("expected N/A status, got %s", report.OverallStatus)
	}
}

func TestOWASPReportAllPassed(t *testing.T) {
	report := &OWASPReport{
		Passed: 10,
		Failed: 0,
	}
	
	report.calculateCompliance()
	
	if report.Compliance != 100.0 {
		t.Errorf("expected 100%% compliance, got %.1f%%", report.Compliance)
	}
	
	if report.OverallStatus != "COMPLIANT" {
		t.Errorf("expected COMPLIANT status, got %s", report.OverallStatus)
	}
}

func TestRunValidationWithMockHost(t *testing.T) {
	// Use a non-routable IP to ensure connection fails (expected in test env)
	validator := NewOWASPValidator("10.255.255.1", 1)
	validator.Timeout = 1 * time.Second // Short timeout for tests
	
	ctx := context.Background()
	report, err := validator.RunValidation(ctx)
	
	// Connection should fail (host unreachable), which is expected
	// This tests the validation flow works correctly
	if err == nil {
		// If somehow connected, verify report structure
		if report == nil {
			t.Error("report should not be nil")
		}
	}
	
	// At minimum, verify we get a report back
	if report == nil {
		t.Fatal("expected report to be returned")
	}
	
	if len(report.Checks) != 10 {
		t.Errorf("expected 10 checks in report, got %d", len(report.Checks))
	}
	
	t.Logf("Validation completed: %.1f%% compliant (%d/%d passed)", 
		report.Compliance, report.Passed, report.Failed)
}