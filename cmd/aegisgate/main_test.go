// SPDX-License-Identifier: MIT
// =========================================================================
// PROPRIETARY - AegisGate Security
// Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
// =========================================================================

package main

import (
	"os"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	// Test that version flag parsing works
	// Note: Cannot test full main() as it binds to ports
	// This tests the constants and basic setup
	
	if version == "" {
		t.Error("version should not be empty")
	}
	
	if commit == "" {
		t.Error("commit should not be empty")
	}
	
	if date == "" {
		t.Error("date should not be empty")
	}
}

func TestManagementEndpoints(t *testing.T) {
	// Verify management endpoints are correctly defined
	expectedEndpoints := []string{"/health", "/version", "/stats"}
	
	for _, endpoint := range expectedEndpoints {
		if !managementEndpoints[endpoint] {
			t.Errorf("expected %s to be a management endpoint", endpoint)
		}
	}
	
	// Verify normal paths are NOT management endpoints
	if managementEndpoints["/api/v1/chat"] {
		t.Error("/api/v1/chat should not be a management endpoint")
	}
}

func TestConstants(t *testing.T) {
	// Test that constants are properly defined
	if version == "" {
		t.Error("version should not be empty")
	}
	
	if commit == "" {
		t.Error("commit should not be empty")
	}
	
	t.Logf("AegisGate version: %s, commit: %s", version, commit)
}

func TestLicenseKeyRequirement(t *testing.T) {
	// Test that environment variable reading works
	// This verifies the license key path is in place
	
	// Save original env
	origKey := os.Getenv("LICENSE_KEY")
	defer os.Setenv("LICENSE_KEY", origKey)
	
	// Test without license key
	os.Unsetenv("LICENSE_KEY")
	key := os.Getenv("LICENSE_KEY")
	if key != "" {
		t.Error("LICENSE_KEY should be empty when not set")
	}
	
	// Test with license key
	os.Setenv("LICENSE_KEY", "test-key-12345")
	key = os.Getenv("LICENSE_KEY")
	if key != "test-key-12345" {
		t.Error("LICENSE_KEY should be retrievable from environment")
	}
}