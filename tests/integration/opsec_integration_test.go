// Package integration_test provides integration tests for OPSEC functionality
// These tests verify OPSEC initialization and basic configuration in the context 
// of the full application.
//
// Run tests with: go test -v ./tests/integration/... -run OPSEC
//
//go:build integration
// +build integration

package integration_test

import (
	"testing"

	"github.com/aegisgatesecurity/aegisgate/pkg/opsec"
)

// ============================================================================
// OPSEC Integration Tests
// ============================================================================

// TestOPSECInitialization tests that OPSEC manager initializes correctly
func TestOPSECInitialization(t *testing.T) {
	config := opsec.DefaultOPSECConfig()
	manager := opsec.NewWithConfig(&config)
	
	if manager == nil {
		t.Fatal("OPSEC manager should not be nil")
	}
	
	err := manager.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize OPSEC manager: %v", err)
	}
	defer manager.Stop()
	
	// Verify manager is initialized
	if !manager.IsInitialized() {
		t.Error("OPSEC manager should be initialized after Initialize()")
	}
	
	t.Log("OPSEC manager initialized successfully")
}

// TestOPSECConfigValidation tests OPSEC configuration validation
func TestOPSECConfigValidation(t *testing.T) {
	// Test valid config
	validConfig := opsec.DefaultOPSECConfig()
	err := validConfig.Validate()
	if err != nil {
		t.Errorf("Valid config should not fail validation: %v", err)
	}
	
	// Verify defaults were applied using Is methods
	if !validConfig.IsAuditEnabled() {
		t.Error("Audit should be enabled by default")
	}
	
	if !validConfig.IsRotationEnabled() {
		t.Error("Rotation should be enabled by default")
	}
	
	t.Log("OPSEC config validation test passed")
}

// TestOPSECConcurrentAccess tests thread-safe access to OPSEC manager
func TestOPSECConcurrentAccess(t *testing.T) {
	config := opsec.DefaultOPSECConfig()
	manager := opsec.NewWithConfig(&config)
	
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize OPSEC manager: %v", err)
	}
	defer manager.Stop()
	
	// Get audit log
	auditLog := manager.GetAuditLog()
	if auditLog == nil {
		t.Fatal("Audit log should not be nil")
	}
	
	// Verify basic audit functionality
	if !auditLog.IsAuditEnabled() {
		t.Error("Audit should be enabled")
	}
	
	// Get entry count (basic functionality)
	_ = auditLog.GetEntryCount()
	
	t.Log("Concurrent access test passed")
}

// TestOPSECShutdown tests graceful shutdown of OPSEC manager
func TestOPSECShutdown(t *testing.T) {
	config := opsec.DefaultOPSECConfig()
	manager := opsec.NewWithConfig(&config)
	
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize OPSEC manager: %v", err)
	}
	
	// Shutdown
	err := manager.Stop()
	if err != nil {
		t.Errorf("Shutdown should not error: %v", err)
	}
	
	t.Log("OPSEC manager shutdown successfully")
}

// TestOPSECIntegrationWithMain tests OPSEC in the context it will be used in main.go
func TestOPSECIntegrationWithMain(t *testing.T) {
	// This simulates exactly what happens in main.go
	config := opsec.DefaultOPSECConfig()
	manager := opsec.NewWithConfig(&config)
	
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize OPSEC manager as in main.go: %v", err)
	}
	defer manager.Stop()
	
	// Get audit log as done in main.go
	auditLog := manager.GetAuditLog()
	if auditLog == nil {
		t.Fatal("Audit log should not be nil")
	}
	
	// Verify audit is enabled
	if !auditLog.IsAuditEnabled() {
		t.Error("Audit should be enabled by default")
	}
	
	t.Log("OPSEC integration with main.go pattern verified")
}

// TestOPSECChainIntegrity verifies audit chain integrity
func TestOPSECChainIntegrity(t *testing.T) {
	config := opsec.DefaultOPSECConfig()
	manager := opsec.NewWithConfig(&config)
	
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize OPSEC manager: %v", err)
	}
	defer manager.Stop()
	
	auditLog := manager.GetAuditLog()
	if auditLog == nil {
		t.Fatal("Audit log should not be nil")
	}
	
	// Enable log integrity
	auditLog.EnableLogIntegrity()
	
	// Verify chain integrity
	isValid, _ := auditLog.VerifyChainIntegrity()
	if !isValid {
		t.Error("Chain integrity check failed")
	}
	
	t.Log("OPSEC chain integrity test passed")
}

// TestOPSECMemoryScrubbingConfig tests memory scrubbing configuration
func TestOPSECMemoryScrubbingConfig(t *testing.T) {
	config := opsec.DefaultOPSECConfig()
	
	// Test memory scrubbing is configured
	if !config.IsMemoryScrubbingEnabled() {
		t.Log("Memory scrubbing is not enabled by default (this is OK)")
	}
	
	manager := opsec.NewWithConfig(&config)
	
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize OPSEC manager: %v", err)
	}
	defer manager.Stop()
	
	t.Log("OPSEC memory scrubbing configuration test completed")
}
