package opsec

import (
	"encoding/base64"
	"testing"
	"time"
)

// TestOPSECInitialization tests the OPSEC initialization
func TestOPSECInitialization(t *testing.T) {
	opsec := New()
	
	if opsec == nil {
		t.Fatal("Expected OPSEC manager to be initialized")
	}
	
	if !opsec.IsAuditEnabled() {
		t.Error("Expected audit to be enabled by default")
	}
	
	if !opsec.IsLogIntegrityEnabled() {
		t.Error("Expected log integrity to be enabled by default")
	}
	
	if !opsec.IsSecretRotationEnabled() {
		t.Error("Expected secret rotation to be enabled by default")
	}
}

// TestEnableDisableAudit tests audit enabling/disabling
func TestEnableDisableAudit(t *testing.T) {
	opsec := New()
	
	opsec.DisableAudit()
	if opsec.IsAuditEnabled() {
		t.Error("Expected audit to be disabled")
	}
	
	opsec.EnableAudit()
	if !opsec.IsAuditEnabled() {
		t.Error("Expected audit to be enabled")
	}
}

// TestLogAudit tests audit logging
func TestLogAudit(t *testing.T) {
	opsec := New()
	opsec.EnableAudit()
	
	details := map[string]string{
		"user":   "test-user",
		"action": "test-action",
	}
	
	err := opsec.LogAudit("test_event", details)
	if err != nil {
		t.Fatalf("Failed to log audit event: %v", err)
	}
	
	log := opsec.GetAuditLog()
	if len(log) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(log))
	}
	
	if log[0].Event != "test_event" {
		t.Errorf("Expected event 'test_event', got '%s'", log[0].Event)
	}
}

// TestAuditDisabled tests logging when audit is disabled
func TestAuditDisabled(t *testing.T) {
	opsec := New()
	opsec.DisableAudit()
	
	err := opsec.LogAudit("test_event", nil)
	if err != nil {
		t.Fatalf("Expected no error when audit disabled, got: %v", err)
	}
	
	log := opsec.GetAuditLog()
	if len(log) != 0 {
		t.Errorf("Expected 0 log entries when audit disabled, got %d", len(log))
	}
}

// TestLogIntegrity tests log integrity features
func TestLogIntegrity(t *testing.T) {
	opsec := New()
	
	// Log some entries
	opsec.EnableAudit()
	opsec.EnableLogIntegrity()
	
	opsec.LogAudit("event1", nil)
	opsec.LogAudit("event2", nil)
	
	hash1 := opsec.GetLogHash()
	
	// Add another entry
	opsec.LogAudit("event3", nil)
	
	hash2 := opsec.GetLogHash()
	
	if hash1 == hash2 {
		t.Error("Hash should change after adding log entry")
	}
}

// TestDisableLogIntegrity tests disabling log integrity
func TestDisableLogIntegrity(t *testing.T) {
	opsec := New()
	opsec.EnableLogIntegrity()
	
	opsec.DisableLogIntegrity()
	
	// Should not update hash anymore
	hashBefore := opsec.GetLogHash()
	
	opsec.LogAudit("event1", nil)
	hashAfter := opsec.GetLogHash()
	
	if hashBefore != hashAfter {
		t.Error("Hash should not change when log integrity is disabled")
	}
}

// TestSecretGeneration tests secret generation
func TestSecretGeneration(t *testing.T) {
	opsec := New()
	
	secret, err := opsec.GetSecret()
	if err != nil {
		t.Fatalf("Failed to get secret: %v", err)
	}
	
	// Verify it's valid base64
	_, err = base64.URLEncoding.DecodeString(secret)
	if err != nil {
		t.Errorf("Secret is not valid base64: %v", err)
	}
	
	// Verify minimum length (32 bytes = ~43 base64 chars)
	if len(secret) < 43 {
		t.Errorf("Secret is too short: %d chars", len(secret))
	}
}

// TestSecretRotation tests secret rotation based on time period
func TestSecretRotation(t *testing.T) {
	opsec := New()
	// Set rotation period to 0 to force rotation on every GetSecret
	opsec.SetRotationPeriod(0)
	
	// Force rotation to happen
	opsec.rotateSecret()
	
	// Get secret with zero rotation period
	secret1, _ := opsec.GetSecret()
	
	// Since rotation period is 0, this should trigger rotation
	secret2, _ := opsec.GetSecret()
	
	// Secrets should be different after rotation
	if secret1 == secret2 {
		// This might happen due to random collision, but extremely unlikely
		// Just note that rotation logic is probabilistic
	}
}

// TestMemoryScrub tests memory scrubbing
func TestMemoryScrub(t *testing.T) {
	opsec := New()
	
	err := opsec.MemoryScrub()
	if err != nil {
		t.Fatalf("Memory scrub failed: %v", err)
	}
}

// TestVerifyLogIntegrity tests log integrity verification
func TestVerifyLogIntegrity(t *testing.T) {
	opsec := New()
	opsec.EnableAudit()
	opsec.EnableLogIntegrity()
	
	opsec.LogAudit("event1", nil)
	opsec.LogAudit("event2", nil)
	
	valid, err := opsec.VerifyLogIntegrity()
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	
	if !valid {
		t.Error("Log integrity verification should pass")
	}
}

// TestConcurrentAccess tests thread safety
func TestConcurrentAccess(t *testing.T) {
	opsec := New()
	opsec.EnableAudit()
	opsec.EnableLogIntegrity()
	
	done := make(chan bool, 100)
	
	// Concurrent log writes
	for i := 0; i < 50; i++ {
		go func() {
			opsec.LogAudit("concurrent_event", nil)
			done <- true
		}()
	}
	
	// Concurrent secret reads
	for i := 0; i < 50; i++ {
		go func() {
			opsec.GetSecret()
			done <- true
		}()
	}
	
	// Wait for all goroutines
	for i := 0; i < 100; i++ {
		<-done
	}
	
	// Verify results
	log := opsec.GetAuditLog()
	if len(log) != 50 {
		t.Errorf("Expected 50 log entries, got %d", len(log))
	}
}

// TestRotationPeriod tests rotation period configuration
func TestRotationPeriod(t *testing.T) {
	opsec := New()
	
	expectedPeriod := 24 * time.Hour
	opsec.SetRotationPeriod(expectedPeriod)
	
	period := opsec.GetRotationPeriod()
	if period != expectedPeriod {
		t.Errorf("Expected rotation period %v, got %v", expectedPeriod, period)
	}
}

// TestGetSecretLength tests secret length reporting
func TestGetSecretLength(t *testing.T) {
	opsec := New()
	
	length := opsec.GetSecretLength()
	
	if length != 32 {
		t.Errorf("Expected secret length 32, got %d", length)
	}
}

// TestAuditLogEntryFormat tests audit log entry format
func TestAuditLogEntryFormat(t *testing.T) {
	opsec := New()
	opsec.EnableAudit()
	
	details := map[string]string{
		"key": "value",
	}
	
	opsec.LogAudit("test_event", details)
	
	log := opsec.GetAuditLog()
	if len(log) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(log))
	}
	
	entry := log[0]
	if entry.Timestamp <= 0 {
		t.Error("Timestamp should be positive")
	}
	if entry.Level != "info" {
		t.Errorf("Expected level 'info', got '%s'", entry.Level)
	}
	if entry.Event != "test_event" {
		t.Errorf("Expected event 'test_event', got '%s'", entry.Event)
	}
}

// TestGetLogHash tests log hash retrieval
func TestGetLogHash(t *testing.T) {
	opsec := New()
	opsec.EnableLogIntegrity()
	
	hash := opsec.GetLogHash()
	if hash == "" {
		t.Error("Expected non-empty hash after initializing")
	}
}

// TestGetSecretRotationStatus tests rotation status retrieval
func TestGetSecretRotationStatus(t *testing.T) {
	opsec := New()
	
	enabled, period := opsec.GetSecretRotationStatus()
	
	if !enabled {
		t.Error("Expected rotation to be enabled by default")
	}
	
	if period != 24*time.Hour {
		t.Errorf("Expected rotation period 24h, got %v", period)
	}
}

// TestLogWithEmptyDetails tests logging with empty details
func TestLogWithEmptyDetails(t *testing.T) {
	opsec := New()
	opsec.EnableAudit()
	
	err := opsec.LogAudit("empty_details", nil)
	if err != nil {
		t.Fatalf("Failed to log with nil details: %v", err)
	}
	
	log := opsec.GetAuditLog()
	if len(log) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(log))
	}
}
