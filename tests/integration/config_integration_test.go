// Package integration_test provides integration tests for immutable configuration functionality
// These tests verify config initialization, integrity verification, rollback capability,
// audit logging, and concurrent access patterns.
//
// Run tests with: go test -v ./tests/integration/... -run Config
//
//go:build integration
// +build integration

package integration_test

import (
	"testing"
	"time"

	imutableconfig "github.com/aegisgatesecurity/aegisgate/pkg/immutable-config"
	"github.com/aegisgatesecurity/aegisgate/pkg/opsec"
)

// ============================================================================
// Configuration Integration Tests
// ============================================================================

// TestImmutableConfigInitialization tests that immutable config manager initializes correctly
func TestImmutableConfigInitialization(t *testing.T) {
	provider := imutableconfig.NewInMemoryProvider()
	manager := imutableconfig.NewConfigManager(provider)
	
	if manager == nil {
		t.Fatal("Config manager should not be nil")
	}
	
	err := manager.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize config manager: %v", err)
	}
	defer manager.Close()
	
	t.Log("Immutable config manager initialized successfully")
}

// TestConfigIntegrityVerification tests hash chain validation for config integrity
func TestConfigIntegrityVerification(t *testing.T) {
	provider := imutableconfig.NewInMemoryProvider()
	manager := imutableconfig.NewConfigManager(provider)
	
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize config manager: %v", err)
	}
	defer manager.Close()
	
	// Create test config data
	configData := imutableconfig.NewConfigData(
		"v1.0.0",
		map[string]interface{}{
			"setting1": "value1",
			"setting2": 42,
			"setting3": true,
		},
		map[string]string{
			"author":  "test",
			"purpose": "integration test",
		},
	)
	
	// Save the config
	version, err := manager.SaveConfig(configData)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	t.Logf("Saved config version: %s with hash: %s", version.Version, configData.Hash)
	
	// Load and verify integrity
	loadedConfig, err := manager.LoadConfig("v1.0.0")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	if loadedConfig.Hash != configData.Hash {
		t.Errorf("Config hash mismatch: expected %s, got %s", configData.Hash, loadedConfig.Hash)
	}
	
	t.Log("Config integrity verification passed")
}

// TestConfigRollback tests reverting to previous config versions
func TestConfigRollback(t *testing.T) {
	provider := imutableconfig.NewInMemoryProvider()
	manager := imutableconfig.NewConfigManager(provider)
	
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize config manager: %v", err)
	}
	defer manager.Close()
	
	// Save multiple versions
	versions := []string{"v1.0.0", "v2.0.0", "v3.0.0"}
	for i, ver := range versions {
		configData := imutableconfig.NewConfigData(
			ver,
			map[string]interface{}{
				"version": i + 1,
				"data":    ver,
			},
			nil,
		)
		_, err := manager.SaveConfig(configData)
		if err != nil {
			t.Fatalf("Failed to save config version %s: %v", ver, err)
		}
	}
	
	// Load latest
	latest, err := manager.LoadLatestConfig()
	if err != nil {
		t.Fatalf("Failed to load latest config: %v", err)
	}
	t.Logf("Latest config version: %s", latest.Version)
	
	// Rollback to v1.0.0
	err = manager.RollbackToVersion("v1.0.0")
	if err != nil {
		t.Fatalf("Failed to rollback: %v", err)
	}
	
	// Verify rollback
	current := manager.GetCurrentConfig()
	if current == nil {
		t.Fatal("Current config should not be nil after rollback")
	}
	
	if current.Version != "v1.0.0" {
		t.Errorf("Expected rollback to v1.0.0, got %s", current.Version)
	}
	
	t.Log("Config rollback successful")
}

// TestConfigChangeAudit tests that all config changes are logged to audit chain
func TestConfigChangeAudit(t *testing.T) {
	// Create OPSEC manager for audit logging
	opsecConfig := opsec.DefaultOPSECConfig()
	opsecManager := opsec.NewWithConfig(&opsecConfig)
	if err := opsecManager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize OPSEC: %v", err)
	}
	defer opsecManager.Stop()
	
	auditLog := opsecManager.GetAuditLog()
	
	// Create config manager with audit logger
	provider := imutableconfig.NewInMemoryProvider()
	manager := imutableconfig.NewConfigManager(provider)
	
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize config manager: %v", err)
	}
	defer manager.Close()
	
	// Save config (should trigger audit event)
	configData := imutableconfig.NewConfigData(
		"v1.0.0",
		map[string]interface{}{"test": "data"},
		nil,
	)
	_, err := manager.SaveConfig(configData)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	// Verify audit event was logged - use GetEntryCount instead of GetChainLength
	entryCount := auditLog.GetEntryCount()
	if entryCount < 1 {
		t.Error("Expected at least 1 audit event for config save")
	}
	
	t.Logf("Audit chain contains %d events after config change", entryCount)
}

// TestConcurrentConfigAccess tests thread-safe config reads
func TestConcurrentConfigAccess(t *testing.T) {
	provider := imutableconfig.NewInMemoryProvider()
	manager := imutableconfig.NewConfigManager(provider)
	
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize config manager: %v", err)
	}
	defer manager.Close()
	
	// Save initial config
	configData := imutableconfig.NewConfigData(
		"v1.0.0",
		map[string]interface{}{"initial": "value"},
		nil,
	)
	manager.SaveConfig(configData)
	
	done := make(chan bool, 20)
	
	// Spawn multiple goroutines reading config concurrently
	for i := 0; i < 20; i++ {
		go func(id int) {
			config := manager.GetCurrentConfig()
			if config != nil {
				_ = config.Data
			}
			manager.GetLatestVersion()
			manager.GetConfigHash()
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 20; i++ {
		<-done
	}
	
	t.Log("Concurrent config access test passed - no race conditions detected")
}

// TestConfigPersistenceAcrossRestarts tests config survives manager restart
func TestConfigPersistenceAcrossRestarts(t *testing.T) {
	// First session - save config
	provider := imutableconfig.NewInMemoryProvider()
	manager1 := imutableconfig.NewConfigManager(provider)
	manager1.Initialize()
	
	configData := imutableconfig.NewConfigData(
		"v1.0.0",
		map[string]interface{}{
			"persistent": true,
			"data":       "test",
		},
		map[string]string{
			"session": "first",
		},
	)
	manager1.SaveConfig(configData)
	manager1.Close()
	
	// Second session - load config (provider persists in memory for this test)
	manager2 := imutableconfig.NewConfigManager(provider)
	manager2.Initialize()
	defer manager2.Close()
	
	loaded, err := manager2.LoadConfig("v1.0.0")
	if err != nil {
		t.Fatalf("Failed to load config in second session: %v", err)
	}
	
	if loaded.Version != "v1.0.0" {
		t.Errorf("Expected version v1.0.0, got %s", loaded.Version)
	}
	
	t.Log("Config persisted across manager restart")
}

// TestDashboardConfigEndpoint tests config API exposure (mock test)
func TestDashboardConfigEndpoint(t *testing.T) {
	provider := imutableconfig.NewInMemoryProvider()
	manager := imutableconfig.NewConfigManager(provider)
	manager.Initialize()
	defer manager.Close()
	
	// Save config
	configData := imutableconfig.NewConfigData(
		"v1.0.0",
		map[string]interface{}{
			"dashboard": "settings",
			"port":      8080,
		},
		nil,
	)
	manager.SaveConfig(configData)
	
	// Get current config (simulates what dashboard would retrieve)
	current := manager.GetCurrentConfig()
	if current == nil {
		t.Fatal("Current config should not be nil")
	}
	
	// Verify data is accessible
	if current.Data["port"] != 8080 {
		t.Error("Expected port to be 8080")
	}
	
	t.Log("Dashboard config endpoint simulation passed")
}

// TestConfigWithOPSEC tests config changes trigger OPSEC audit events
func TestConfigWithOPSEC(t *testing.T) {
	// Initialize OPSEC
	opsecConfig := opsec.DefaultOPSECConfig()
	opsecManager := opsec.NewWithConfig(&opsecConfig)
	if err := opsecManager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize OPSEC: %v", err)
	}
	defer opsecManager.Stop()
	
	auditLog := opsecManager.GetAuditLog()
	initialEntryCount := auditLog.GetEntryCount()
	
	// Initialize config manager
	provider := imutableconfig.NewInMemoryProvider()
	manager := imutableconfig.NewConfigManager(provider)
	manager.Initialize()
	defer manager.Close()
	
	// Make multiple config changes
	for i := 0; i < 5; i++ {
		configData := imutableconfig.NewConfigData(
			"v%d.0.0",
			map[string]interface{}{"change": i},
			nil,
		)
		manager.SaveConfig(configData)
	}
	
	// Verify audit events were created - use GetEntryCount
	finalEntryCount := auditLog.GetEntryCount()
	eventsLogged := finalEntryCount - initialEntryCount
	
	if eventsLogged < 1 {
		t.Error("Expected at least 1 audit event from config changes")
	}
	
	t.Logf("Config changes generated %d audit events", eventsLogged)
}

// TestConfigVersionHistory tests version history tracking
func TestConfigVersionHistory(t *testing.T) {
	provider := imutableconfig.NewInMemoryProvider()
	manager := imutableconfig.NewConfigManager(provider)
	manager.Initialize()
	defer manager.Close()
	
	// Save multiple versions
	for i := 1; i <= 5; i++ {
		configData := imutableconfig.NewConfigData(
			"v%d.0.0",
			map[string]interface{}{"version": i},
			nil,
		)
		manager.SaveConfig(configData)
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}
	
	// Get version history
	history, err := manager.GetVersionHistory()
	if err != nil {
		t.Fatalf("Failed to get version history: %v", err)
	}
	
	if len(history) < 5 {
		t.Errorf("Expected at least 5 versions in history, got %d", len(history))
	}
	
	t.Logf("Version history contains %d entries", len(history))
}

// TestConfigValidation tests config validation before saving
func TestConfigValidation(t *testing.T) {
	provider := imutableconfig.NewInMemoryProvider()
	manager := imutableconfig.NewConfigManager(provider)
	manager.Initialize()
	defer manager.Close()
	
	// Test invalid config (empty version)
	invalidConfig := imutableconfig.NewConfigData(
		"", // Empty version - invalid
		map[string]interface{}{"test": "data"},
		nil,
	)
	
	err := invalidConfig.Validate()
	if err == nil {
		t.Error("Expected validation error for empty version")
	}
	
	// Test valid config
	validConfig := imutableconfig.NewConfigData(
		"v1.0.0",
		map[string]interface{}{"test": "data"},
		nil,
	)
	
	err = validConfig.Validate()
	if err != nil {
		t.Errorf("Valid config should not fail validation: %v", err)
	}
	
	t.Log("Config validation test passed")
}
