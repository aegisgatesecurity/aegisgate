// Package core provides license management for module activation.
package core

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// LicenseStatus represents the validity of a license.
type LicenseStatus int

const (
	LicenseStatusValid LicenseStatus = iota
	LicenseStatusExpired
	LicenseStatusInvalid
	LicenseStatusNotFound
	LicenseStatusTierMismatch
)

func (s LicenseStatus) String() string {
	switch s {
	case LicenseStatusValid:
		return "valid"
	case LicenseStatusExpired:
		return "expired"
	case LicenseStatusInvalid:
		return "invalid"
	case LicenseStatusNotFound:
		return "not_found"
	case LicenseStatusTierMismatch:
		return "tier_mismatch"
	default:
		return "unknown"
	}
}

// LicenseType represents the type of license.
type LicenseType string

const (
	LicenseTypeCommunity    LicenseType = "community"
	LicenseTypeProfessional LicenseType = "professional"
	LicenseTypeEnterprise   LicenseType = "enterprise"
	LicenseTypeEnterpriseAI LicenseType = "enterprise-ai"
	LicenseTypeCustom       LicenseType = "custom"
)

// License represents a parsed license.
type License struct {
	ID           string      `json:"id"`
	Type         LicenseType `json:"type"`
	Email        string      `json:"email"`
	Organization string      `json:"organization,omitempty"`
	Modules      []string    `json:"modules,omitempty"`
	Tiers        []Tier      `json:"tiers,omitempty"`
	IssuedAt     time.Time   `json:"issued_at"`
	ExpiresAt    time.Time   `json:"expires_at"`
	MaxServers   int         `json:"max_servers,omitempty"`
	Features     []string    `json:"features,omitempty"`
	Signature    string      `json:"signature,omitempty"`
}

// LicenseConfig contains configuration for the license manager.
type LicenseConfig struct {
	LicenseKey   string
	PublicKeyPEM string // For production: embedded public key
	GracePeriod  time.Duration
}

// LicenseManager handles license validation and module activation.
type LicenseManager struct {
	mu          sync.RWMutex
	config      LicenseConfig
	license     *License
	status      LicenseStatus
	validatedAt time.Time
	publicKey   *rsa.PublicKey
	overrides   map[string]bool // For development/testing
}

// NewLicenseManager creates a new license manager.
func NewLicenseManager(licenseKey string) *LicenseManager {
	lm := &LicenseManager{
		config: LicenseConfig{
			LicenseKey:  licenseKey,
			GracePeriod: 7 * 24 * time.Hour, // 7 days grace period
		},
		overrides: make(map[string]bool),
		status:    LicenseStatusNotFound,
	}

	// Validate license on creation
	if licenseKey != "" {
		_ = lm.validateLicense()
	}

	return lm
}

// SetLicenseKey sets or updates the license key.
func (lm *LicenseManager) SetLicenseKey(key string) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.config.LicenseKey = key
	lm.license = nil
	lm.status = LicenseStatusNotFound

	return lm.validateLicense()
}

// GetLicense returns the current license.
func (lm *LicenseManager) GetLicense() *License {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return lm.license
}

// GetStatus returns the current license status.
func (lm *LicenseManager) GetStatus() LicenseStatus {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return lm.status
}

// IsModuleLicensed checks if a specific module is licensed.
func (lm *LicenseManager) IsModuleLicensed(moduleID string, tier Tier) bool {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	// Development mode: check environment variable
	if os.Getenv("AEGISGATE_DEV_MODE") == "true" {
		return true
	}

	// Check overrides (for testing)
	if allowed, exists := lm.overrides[moduleID]; exists {
		return allowed
	}

	// Community is always free
	if tier <= TierCommunity {
		return true
	}

	// No license = only free tiers
	if lm.license == nil || lm.status != LicenseStatusValid {
		return false
	}

	// Check if license has explicit module permission
	for _, mod := range lm.license.Modules {
		if mod == moduleID || mod == "*" {
			return true
		}
	}

	// Check if license covers this tier
	for _, licensedTier := range lm.license.Tiers {
		if licensedTier >= tier {
			return true
		}
	}

	// Check license type tier permissions
	return lm.tierAllowedByLicenseType(tier)
}

// tierAllowedByLicenseType checks if the license type permits the tier.
func (lm *LicenseManager) tierAllowedByLicenseType(tier Tier) bool {
	if lm.license == nil {
		return false
	}

	switch lm.license.Type {
	case LicenseTypeCommunity:
		return tier <= TierCommunity

	case LicenseTypeProfessional:
		return tier <= TierProfessional

	case LicenseTypeEnterprise:
		return tier <= TierEnterprise

	case LicenseTypeEnterpriseAI:
		return tier <= TierEnterprise

	case LicenseTypeCustom:
		// Custom licenses must explicitly grant tiers
		return false

	default:
		return false
	}
}

// IsFeatureLicensed checks if a specific feature is licensed.
func (lm *LicenseManager) IsFeatureLicensed(feature string) bool {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if lm.license == nil {
		return false
	}

	for _, f := range lm.license.Features {
		if f == feature || f == "*" {
			return true
		}
	}

	return false
}

// LicenseExpiringSoon checks if license expires within the given duration.
func (lm *LicenseManager) LicenseExpiringSoon(within time.Duration) bool {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if lm.license == nil {
		return false
	}

	return time.Until(lm.license.ExpiresAt) < within
}

// SetModuleOverride sets a development/testing override for a module.
func (lm *LicenseManager) SetModuleOverride(moduleID string, allowed bool) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	lm.overrides[moduleID] = allowed
}

// ClearOverrides removes all development overrides.
func (lm *LicenseManager) ClearOverrides() {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	lm.overrides = make(map[string]bool)
}

// validateLicense parses and validates the license key.
func (lm *LicenseManager) validateLicense() error {
	if lm.config.LicenseKey == "" {
		lm.status = LicenseStatusNotFound
		return fmt.Errorf("no license key configured")
	}

	// Decode license
	license, err := lm.parseLicense(lm.config.LicenseKey)
	if err != nil {
		lm.status = LicenseStatusInvalid
		return err
	}

	// Check expiration
	if time.Now().After(license.ExpiresAt) {
		lm.status = LicenseStatusExpired
		return fmt.Errorf("license expired at %s", license.ExpiresAt)
	}

	lm.license = license
	lm.status = LicenseStatusValid
	lm.validatedAt = time.Now()

	return nil
}

// parseLicense parses a license key into a License struct.
func (lm *LicenseManager) parseLicense(key string) (*License, error) {
	// License format: BASE64 encoded JSON
	// Production: Signed JWT or similar

	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		// Try raw URL encoding
		decoded, err = base64.RawURLEncoding.DecodeString(key)
		if err != nil {
			return nil, fmt.Errorf("invalid license encoding: %w", err)
		}
	}

	// Try to parse the entire decoded content as JSON first
	// (unsigned license). If that fails, check for signature format.
	var license License
	if err := json.Unmarshal(decoded, &license); err == nil {
		// Successfully parsed unsigned license
		return &license, nil
	}

	// Check for signed license format: payload.signature
	// The signature is appended after a period, but we need to find
	// where the JSON payload ends first
	decodedStr := string(decoded)

	// Find the last occurrence of "}" which marks the end of JSON
	lastBrace := strings.LastIndex(decodedStr, "}")
	if lastBrace == -1 {
		return nil, fmt.Errorf("invalid license format: no JSON object found")
	}

	jsonPayload := []byte(decodedStr[:lastBrace+1])
	signature := decodedStr[lastBrace+1:]

	// Remove leading period from signature if present
	signature = strings.TrimPrefix(signature, ".")

	// If there's a signature, verify it (for production)
	if signature != "" {
		// For production, implement proper signature verification
		// For now, we just parse the payload
		_ = signature
	}

	if err := json.Unmarshal(jsonPayload, &license); err != nil {
		return nil, fmt.Errorf("invalid license format: %w", err)
	}

	return &license, nil
}

// GenerateLicense generates a new license (for testing/admin purpose).
func GenerateLicense(licenseType LicenseType, email string, modules []string, tiers []Tier, validFor time.Duration) (string, error) {
	now := time.Now()
	license := License{
		ID:        generateLicenseID(),
		Type:      licenseType,
		Email:     email,
		Modules:   modules,
		Tiers:     tiers,
		IssuedAt:  now,
		ExpiresAt: now.Add(validFor),
	}

	data, err := json.Marshal(license)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// generateLicenseID creates a unique license ID.
func generateLicenseID() string {
	// Simple ID generation - use proper UUID in production
	return fmt.Sprintf("PAD-%d", time.Now().UnixNano())
}

// LicenseSummary provides a human-readable license summary.
type LicenseSummary struct {
	Type          string    `json:"type"`
	Email         string    `json:"email"`
	Organization  string    `json:"organization,omitempty"`
	Modules       []string  `json:"modules"`
	MaxTier       string    `json:"max_tier"`
	IssuedAt      time.Time `json:"issued_at"`
	ExpiresAt     time.Time `json:"expires_at"`
	DaysRemaining int       `json:"days_remaining"`
	Status        string    `json:"status"`
}

// Summary returns a human-readable license summary.
func (lm *LicenseManager) Summary() *LicenseSummary {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if lm.license == nil {
		return &LicenseSummary{
			Type:    string(LicenseTypeCommunity),
			MaxTier: TierCommunity.String(),
			Status:  LicenseStatusNotFound.String(),
		}
	}

	maxTier := TierCommunity
	for _, t := range lm.license.Tiers {
		if t > maxTier {
			maxTier = t
		}
	}

	return &LicenseSummary{
		Type:          string(lm.license.Type),
		Email:         lm.license.Email,
		Organization:  lm.license.Organization,
		Modules:       lm.license.Modules,
		MaxTier:       maxTier.String(),
		IssuedAt:      lm.license.IssuedAt,
		ExpiresAt:     lm.license.ExpiresAt,
		DaysRemaining: int(time.Until(lm.license.ExpiresAt).Hours() / 24),
		Status:        lm.status.String(),
	}
}

// SetPublicKey sets the public key for license signature verification.
func (lm *LicenseManager) SetPublicKey(pemData []byte) error {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return fmt.Errorf("failed to parse PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an RSA public key")
	}

	lm.publicKey = rsaPub
	return nil
}
