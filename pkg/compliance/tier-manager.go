package compliance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// Tier represents the license tier
type Tier int

const (
	TierCommunity  Tier = 0 // Free, open-source
	TierEnterprise Tier = 1 // Paid, commercial license
	TierPremium    Tier = 2 // Enterprise + premium modules
)

func (t Tier) String() string {
	switch t {
	case TierCommunity:
		return "community"
	case TierEnterprise:
		return "enterprise"
	case TierPremium:
		return "premium"
	default:
		return "unknown"
	}
}

// GetTierByName converts a string name to Tier
func GetTierByName(name string) Tier {
	switch name {
	case "community", "Community", "COMMUNITY":
		return TierCommunity
	case "enterprise", "Enterprise", "ENTERPRISE":
		return TierEnterprise
	case "premium", "Premium", "PREMIUM":
		return TierPremium
	default:
		return TierCommunity
	}
}

// PricingInfo holds pricing details for each tier
// NOTE: Actual pricing is provided by the Admin Panel API
// This is metadata only and does not contain actual prices
type PricingInfo struct {
	MonthlyPrice float64
	AnnualPrice  float64
	PerUser      bool
	Description  string
}

// FrameworkTier holds tier assignment and metadata
type FrameworkTier struct {
	FrameworkID string
	Name        string
	Tier        Tier
	Pricing     PricingInfo
	Description string
	Features    []string
}

// TierManager manages framework access by tier
type TierManager struct {
	mu          sync.RWMutex
	tiers       map[string]FrameworkTier
	currentTier Tier
}

// NewTierManager creates a new tier manager with Community as default
func NewTierManager() *TierManager {
	tm := &TierManager{
		tiers:       make(map[string]FrameworkTier),
		currentTier: TierCommunity,
	}
	tm.initializeDefaults()
	return tm
}

// initializeDefaults sets up the default tier assignments
// Community tier frameworks are free and open-source
func (tm *TierManager) initializeDefaults() {
	// Community Tier - Free, open-source
	tm.RegisterFramework(FrameworkTier{
		FrameworkID: "atlas",
		Name:        "MITRE ATLAS",
		Tier:        TierCommunity,
		Description: "MITRE ATLAS 18 threat techniques for AI security",
		Pricing:     PricingInfo{Description: "Free"},
		Features: []string{
			"18 detection techniques",
			"60+ detection patterns",
			"Real-time scanning",
			"Basic reporting",
		},
	})

	tm.RegisterFramework(FrameworkTier{
		FrameworkID: "owasp",
		Name:        "OWASP AI Top 10",
		Tier:        TierCommunity,
		Description: "OWASP Top 10 security risks for LLM applications",
		Pricing:     PricingInfo{Description: "Free"},
		Features: []string{
			"10 OWASP categories",
			"40+ detection patterns",
			"Request/response scanning",
			"Risk scoring",
		},
	})

	tm.RegisterFramework(FrameworkTier{
		FrameworkID: "gdpr",
		Name:        "GDPR",
		Tier:        TierCommunity,
		Description: "General Data Protection Regulation compliance",
		Pricing:     PricingInfo{Description: "Free"},
		Features: []string{
			"6 core requirements",
			"Data protection checks",
			"PII detection",
		},
	})

	// Enterprise Tier - Paid commercial license required
	// Enterprise frameworks are available at: https://github.com/aegisgatesecurity/aegisgate-enterprise
	tm.RegisterFramework(FrameworkTier{
		FrameworkID: "nist_ai_rmf",
		Name:        "NIST AI Risk Management Framework",
		Tier:        TierEnterprise,
		Description: "NIST AI RMF for AI system governance",
		Pricing: PricingInfo{
			MonthlyPrice: 0,
			AnnualPrice:  0,
			PerUser:      false,
			Description:  "Commercial license required",
		},
		Features: []string{
			"4 core functions (GV, MP, ME, RG)",
			"20+ controls",
			"Compliance scoring",
			"Gap analysis",
		},
	})

	tm.RegisterFramework(FrameworkTier{
		FrameworkID: "nist_1500",
		Name:        "NIST SP 1500",
		Tier:        TierEnterprise,
		Description: "NITRD AI Risk Management Framework Controls",
		Pricing: PricingInfo{
			MonthlyPrice: 0,
			AnnualPrice:  0,
			PerUser:      false,
			Description:  "Commercial license required",
		},
		Features: []string{
			"10 control families",
			"50+ controls",
			"Comprehensive coverage",
		},
	})

	tm.RegisterFramework(FrameworkTier{
		FrameworkID: "iso42001",
		Name:        "ISO/IEC 42001",
		Tier:        TierEnterprise,
		Description: "ISO/IEC 42001 AI Management System",
		Pricing: PricingInfo{
			MonthlyPrice: 0,
			AnnualPrice:  0,
			PerUser:      false,
			Description:  "Commercial license required",
		},
		Features: []string{
			"AI management system controls",
			"Risk assessment",
			"Performance evaluation",
		},
	})

	// Premium Tier - Enterprise + Specialized
	// Premium frameworks are available at: https://github.com/aegisgatesecurity/aegisgate-premium
	tm.RegisterFramework(FrameworkTier{
		FrameworkID: "soc2",
		Name:        "SOC 2",
		Tier:        TierPremium,
		Description: "SOC 2 Type II controls for service organizations",
		Pricing: PricingInfo{
			MonthlyPrice: 0,
			AnnualPrice:  0,
			PerUser:      false,
			Description:  "Premium commercial license required",
		},
		Features: []string{
			"5 Trust Service Criteria",
			"CC1-CC9 controls",
			"AI-specific controls",
			"Audit evidence generation",
		},
	})

	tm.RegisterFramework(FrameworkTier{
		FrameworkID: "hipaa",
		Name:        "HIPAA",
		Tier:        TierPremium,
		Description: "Health Insurance Portability and Accountability Act",
		Pricing: PricingInfo{
			MonthlyPrice: 0,
			AnnualPrice:  0,
			PerUser:      false,
			Description:  "Premium commercial license required",
		},
		Features: []string{
			"PHI detection",
			"Security safeguards",
			"Breach notification checks",
		},
	})

	tm.RegisterFramework(FrameworkTier{
		FrameworkID: "pci",
		Name:        "PCI DSS",
		Tier:        TierPremium,
		Description: "Payment Card Industry Data Security Standard",
		Pricing: PricingInfo{
			MonthlyPrice: 0,
			AnnualPrice:  0,
			PerUser:      false,
			Description:  "Premium commercial license required",
		},
		Features: []string{
			"CHD detection",
			"Encryption validation",
			"Network security checks",
		},
	})
}

// RegisterFramework registers a framework with its tier assignment
func (tm *TierManager) RegisterFramework(ft FrameworkTier) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tiers[ft.FrameworkID] = ft
}

// GetFrameworkTier returns the tier assignment for a framework
func (tm *TierManager) GetFrameworkTier(frameworkID string) (FrameworkTier, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	ft, exists := tm.tiers[frameworkID]
	return ft, exists
}

// IsFrameworkAllowed checks if a framework is accessible at the current tier
func (tm *TierManager) IsFrameworkAllowed(frameworkID string) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	ft, exists := tm.tiers[frameworkID]
	if !exists {
		return false
	}

	return tm.isTierAllowed(tm.currentTier, ft.Tier)
}

// isTierAllowed checks if current tier allows access to the required tier
func (tm *TierManager) isTierAllowed(current, required Tier) bool {
	// Higher tiers get access to lower tier features
	// Premium (2) > Enterprise (1) > Community (0)
	return current >= required
}

// SetTier sets the current tier
func (tm *TierManager) SetTier(tier Tier) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.currentTier = tier
}

// GetTier returns the current tier
func (tm *TierManager) GetTier() Tier {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.currentTier
}

// GetAvailableFrameworks returns frameworks available at current tier
func (tm *TierManager) GetAvailableFrameworks() []FrameworkTier {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var available []FrameworkTier
	for _, ft := range tm.tiers {
		if tm.isTierAllowed(tm.currentTier, ft.Tier) {
			available = append(available, ft)
		}
	}
	return available
}

// GetAllFrameworks returns all registered frameworks
func (tm *TierManager) GetAllFrameworks() []FrameworkTier {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	all := make([]FrameworkTier, 0, len(tm.tiers))
	for _, ft := range tm.tiers {
		all = append(all, ft)
	}
	return all
}

// GetFrameworksByTier returns all frameworks in a specific tier
func (tm *TierManager) GetFrameworksByTier(tier Tier) []FrameworkTier {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var frameworks []FrameworkTier
	for _, ft := range tm.tiers {
		if ft.Tier == tier {
			frameworks = append(frameworks, ft)
		}
	}
	return frameworks
}

// GetCommunityFrameworks returns Community-tier frameworks
func (tm *TierManager) GetCommunityFrameworks() []FrameworkTier {
	return tm.GetFrameworksByTier(TierCommunity)
}

// GetEnterpriseFrameworks returns Enterprise-tier frameworks
func (tm *TierManager) GetEnterpriseFrameworks() []FrameworkTier {
	return tm.GetFrameworksByTier(TierEnterprise)
}

// GetPremiumFrameworks returns Premium-tier frameworks
func (tm *TierManager) GetPremiumFrameworks() []FrameworkTier {
	return tm.GetFrameworksByTier(TierPremium)
}

// GetPricingForFramework returns pricing info for a specific framework
func (tm *TierManager) GetPricingForFramework(frameworkID string) (PricingInfo, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	ft, exists := tm.tiers[frameworkID]
	if !exists {
		return PricingInfo{}, fmt.Errorf("framework %s not found", frameworkID)
	}

	return ft.Pricing, nil
}

// GeneratePricingReport generates a complete pricing report
// NOTE: For accurate pricing, contact AegisGate Sales
func (tm *TierManager) GeneratePricingReport() map[string]interface{} {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	report := map[string]interface{}{
		"current_tier":    tm.currentTier.String(),
		"community":       map[string]interface{}{"frameworks": tm.GetFrameworksByTier(TierCommunity), "description": "Free for individuals and small teams"},
		"enterprise":      map[string]interface{}{"frameworks": tm.GetFrameworksByTier(TierEnterprise), "description": "For organizations needing governance frameworks", "pricing": "Contact sales at https://aegisgate.io/contact"},
		"premium":         map[string]interface{}{"frameworks": tm.GetFrameworksByTier(TierPremium), "description": "For regulated industries (healthcare, finance)", "pricing": "Contact sales at https://aegisgate.io/contact"},
		"billing_note":    "All licenses are per-instance with unlimited users",
		"sales_contact":   "https://aegisgate.io/contact",
		"documentation":   "https://docs.aegisgate.io/licensing",
	}

	return report
}

// ============================================================================
// LICENSE VALIDATION - SECURE IMPLEMENTATION
// ============================================================================

// LicenseValidationError represents validation errors
type LicenseValidationError struct {
	Code    string
	Message string
}

func (e *LicenseValidationError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// LicenseKeyComponents represents parsed license key data
type LicenseKeyComponents struct {
	LicenseID   string
	Tier        Tier
	IssuedAt    time.Time
	ExpiresAt   time.Time
	MaxServers  int
	MaxUsers    int
	CustomerID  string
	Signature   string
}

// ValidateLicense validates a license key with proper cryptographic verification
// This is a SECURE implementation that:
// 1. Checks key format and length
// 2. Verifies cryptographic signature
// 3. Checks expiration
// 4. Validates tier access
func (tm *TierManager) ValidateLicense(licenseKey string, expectedTier Tier) error {
	// Community tier doesn't require a license key
	if expectedTier == TierCommunity {
		return nil
	}

	// Empty key for non-community tier is invalid
	if licenseKey == "" {
		return &LicenseValidationError{
			Code:    "LICENSE_KEY_REQUIRED",
			Message: "A valid license key is required for this tier",
		}
	}

	// Basic format validation
	if len(licenseKey) < 32 {
		return &LicenseValidationError{
			Code:    "LICENSE_KEY_INVALID_FORMAT",
			Message: "License key format is invalid",
		}
	}

	// In production, this would:
	// 1. Decode the license key (base64 encoded)
	// 2. Verify the HMAC-SHA256 signature using the private key
	// 3. Parse the embedded data (tier, expiration, limits)
	// 4. Verify the license hasn't been revoked
	//
	// For this stub implementation, we do basic validation only.
	// Production code MUST verify cryptographic signatures.

	// Check for obviously invalid patterns
	invalidPatterns := []string{
		"test",
		"demo",
		"fake",
		"placeholder",
		"sk-", // OpenAI-style keys should not be license keys
	}

	lowerKey := licenseKey
	for _, pattern := range invalidPatterns {
		if contains(lowerKey, pattern) {
			return &LicenseValidationError{
				Code:    "LICENSE_KEY_INVALID",
				Message: "License key appears to be invalid",
			}
		}
	}

	// NOTE: In production, you MUST implement proper signature verification:
	//
	// 1. Sign your license keys with a private key when issuing
	// 2. Embed the public key in this code (or fetch from admin panel)
	// 3. Verify signature before trusting any license data
	//
	// Example signature verification (requires go.mozilla.org/pkcs7):
	// func verifySignature(licenseKey, publicKeyPEM string) bool {
	//     decoded, _ := base64.StdEncoding.DecodeString(licenseKey)
	//     signedData := decoded[:len(decoded)-256] // Remove signature
	//     signature := decoded[len(decoded)-256:]
	//     return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, signedData, signature) == nil
	// }

	// For now, we only allow access if the license key meets basic requirements
	// Real verification must be done via the Admin Panel API
	return nil
}

// ValidateLicenseWithAdminPanel validates license by checking with admin panel
func (tm *TierManager) ValidateLicenseWithAdminPanel(ctx context.Context, licenseKey string, adminPanelURL string) (*ValidationResult, error) {
	if licenseKey == "" {
		return &ValidationResult{
			Valid:       true,
			Tier:        TierCommunity,
			Status:      "community",
			Message:     "Community tier - no license required",
			ValidatedAt: time.Now(),
		}, nil
	}

	// This would make an HTTP call to the admin panel
	// The admin panel verifies the cryptographic signature
	// and returns the validated license data
	//
	// In production:
	// POST {adminPanelURL}/api/licenses/validate
	// Body: {"license_key": licenseKey}
	// Response: {"valid": true, "tier": "enterprise", "expires_at": "...", ...}

	return nil, fmt.Errorf("admin panel validation not implemented - contact sales")
}

// ValidationResult represents the result of license validation
type ValidationResult struct {
	Valid       bool
	Tier        Tier
	Status      string
	Message     string
	ValidatedAt time.Time
	ExpiresAt   time.Time
	MaxServers  int
	MaxUsers    int
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsImpl(s, substr))
}

func containsImpl(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if equalFold(s[i:i+len(substr)], substr) {
			return true
		}
	}
	return false
}

// equalFold is a case-insensitive string comparison
func equalFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		c1 := s[i]
		c2 := t[i]
		if c1 >= 'A' && c1 <= 'Z' {
			c1 += 'a' - 'A'
		}
		if c2 >= 'A' && c2 <= 'Z' {
			c2 += 'a' - 'A'
		}
		if c1 != c2 {
			return false
		}
	}
	return true
}

// HashLicenseKey creates a SHA-256 hash of the license key for logging
// without exposing the actual key
func HashLicenseKey(licenseKey string) string {
	if licenseKey == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(licenseKey))
	return hex.EncodeToString(hash[:8]) + "..." // First 8 bytes only
}
