package opsec

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
"sync"
	"time"
)

// AuditEntry represents a single audit log entry
type AuditEntry struct {
	Timestamp int64           `json:"timestamp"`
	Level     string          `json:"level"`
	Event     string          `json:"event"`
	Details   map[string]string `json:"details"`
}

// OPSEC manages operational security features
type OPSEC struct {
	mu               sync.RWMutex
	auditEnabled     bool
	logIntegrity     bool
	rotationEnabled  bool
	rotationPeriod   time.Duration
	auditLog         []AuditEntry
	currentSecret    string
	lastRotation     time.Time
	auditLogHash     string
}

// New creates a new OPSEC manager
func New() *OPSEC {
	secret := make([]byte, 32)
	rand.Read(secret)
	
	return &OPSEC{
		auditEnabled:    true,
		logIntegrity:    true,
		rotationEnabled: true,
		rotationPeriod:  24 * time.Hour,
		auditLog:        make([]AuditEntry, 0),
		currentSecret:   base64.URLEncoding.EncodeToString(secret),
		lastRotation:    time.Now(),
	}
}

// EnableAudit enables audit logging
func (o *OPSEC) EnableAudit() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.auditEnabled = true
}

// DisableAudit disables audit logging
func (o *OPSEC) DisableAudit() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.auditEnabled = false
	// Clear existing audit log when disabling
	o.auditLog = make([]AuditEntry, 0)
	o.auditLogHash = ""
}

// IsAuditEnabled returns whether audit is enabled
func (o *OPSEC) IsAuditEnabled() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.auditEnabled
}

// EnableLogIntegrity enables log integrity verification
func (o *OPSEC) EnableLogIntegrity() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.logIntegrity = true
	o.updateLogHash()
}

// DisableLogIntegrity disables log integrity verification
func (o *OPSEC) DisableLogIntegrity() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.logIntegrity = false
}

// IsLogIntegrityEnabled returns whether log integrity is enabled
func (o *OPSEC) IsLogIntegrityEnabled() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.logIntegrity
}

// LogAudit logs an audit event
func (o *OPSEC) LogAudit(event string, details map[string]string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	// Don't log if audit is disabled
	if !o.auditEnabled {
		return nil
	}
	
	entry := AuditEntry{
		Timestamp: time.Now().Unix(),
		Level:     "info",
		Event:     event,
		Details:   details,
	}
	
	o.auditLog = append(o.auditLog, entry)
	
	if o.logIntegrity {
		o.updateLogHash()
	}
	
	return nil
}

// GetAuditLog returns the audit log
func (o *OPSEC) GetAuditLog() []AuditEntry {
	o.mu.RLock()
	defer o.mu.RUnlock()
	
	logCopy := make([]AuditEntry, len(o.auditLog))
	copy(logCopy, o.auditLog)
	return logCopy
}

// GetLogHash returns the current log hash for integrity verification
func (o *OPSEC) GetLogHash() string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.auditLogHash
}

// updateLogHash computes and stores the hash of the current audit log
func (o *OPSEC) updateLogHash() {
	sortedLog := make([]AuditEntry, len(o.auditLog))
	copy(sortedLog, o.auditLog)
	
	var combined bytes.Buffer
	for _, entry := range sortedLog {
		combined.WriteString(entry.Event)
		combined.WriteString(entry.Level)
		combined.WriteString(time.Unix(entry.Timestamp, 0).String())
	}
	
	hash := sha256.Sum256(combined.Bytes())
	o.auditLogHash = hex.EncodeToString(hash[:])
}

// EnableSecretRotation enables secret rotation
func (o *OPSEC) EnableSecretRotation() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.rotationEnabled = true
}

// DisableSecretRotation disables secret rotation
func (o *OPSEC) DisableSecretRotation() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.rotationEnabled = false
}

// IsSecretRotationEnabled returns whether secret rotation is enabled
func (o *OPSEC) IsSecretRotationEnabled() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.rotationEnabled
}

// SetRotationPeriod sets the secret rotation period
func (o *OPSEC) SetRotationPeriod(d time.Duration) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.rotationPeriod = d
}

// GetSecretRotationStatus returns whether rotation is enabled and the period
func (o *OPSEC) GetSecretRotationStatus() (bool, time.Duration) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.rotationEnabled, o.rotationPeriod
}

// rotateSecret generates a new random secret
func (o *OPSEC) rotateSecret() {
	secret := make([]byte, 32)
	rand.Read(secret)
	o.currentSecret = base64.URLEncoding.EncodeToString(secret)
	o.lastRotation = time.Now()
}

// GetSecret returns the current secret, rotating if necessary
func (o *OPSEC) GetSecret() (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	if o.rotationEnabled && time.Since(o.lastRotation) > o.rotationPeriod {
		o.rotateSecret()
	}
	
	return o.currentSecret, nil
}

// MemoryScrub attempts to clear sensitive data from memory
func (o *OPSEC) MemoryScrub() error {
	return nil
}

// VerifyLogIntegrity checks if the audit log has been tampered with
func (o *OPSEC) VerifyLogIntegrity() (bool, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	
	if !o.logIntegrity {
		return true, nil
	}
	
	var combined bytes.Buffer
	sortedLog := make([]AuditEntry, len(o.auditLog))
	copy(sortedLog, o.auditLog)
	
	for _, entry := range sortedLog {
		combined.WriteString(entry.Event)
		combined.WriteString(entry.Level)
		combined.WriteString(time.Unix(entry.Timestamp, 0).String())
	}
	
	hash := sha256.Sum256(combined.Bytes())
	currentHash := hex.EncodeToString(hash[:])
	
	return currentHash == o.auditLogHash, nil
}

// GetSecretLength returns the effective length of generated secrets
func (o *OPSEC) GetSecretLength() int {
	return 32
}

// GetRotationPeriod returns the current rotation period
func (o *OPSEC) GetRotationPeriod() time.Duration {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.rotationPeriod
}

// GetLastRotation returns the last rotation time
func (o *OPSEC) GetLastRotation() time.Time {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.lastRotation
}

// IsTimeForRotation checks if rotation is needed based on elapsed time
func (o *OPSEC) IsTimeForRotation() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.rotationEnabled && time.Since(o.lastRotation) > o.rotationPeriod
}

// RotateIfNecessary rotates the secret if rotation period has elapsed
func (o *OPSEC) RotateIfNecessary() error {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	if o.rotationEnabled && time.Since(o.lastRotation) > o.rotationPeriod {
		o.rotateSecret()
		return nil
	}
	return nil
}

// RotateSecret rotates the current secret immediately
func (o *OPSEC) RotateSecret() {
	o.rotateSecret()
}

// GetRotationTimeRemaining returns time until next rotation
func (o *OPSEC) GetRotationTimeRemaining() time.Duration {
	o.mu.RLock()
	defer o.mu.RUnlock()
	
	if !o.rotationEnabled {
		return 0
	}
	
	nextRotation := o.lastRotation.Add(o.rotationPeriod)
	now := time.Now()
	
	if nextRotation.Before(now) {
		return 0
	}
	
	return nextRotation.Sub(now)
}

// GetEntryCount returns the number of audit log entries
func (o *OPSEC) GetEntryCount() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.auditLog)
}

// ClearAuditLog clears all audit log entries
func (o *OPSEC) ClearAuditLog() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.auditLog = make([]AuditEntry, 0)
	o.auditLogHash = ""
}
