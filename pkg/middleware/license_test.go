// SPDX-License-Identifier: MIT
// =========================================================================
// AegisGate License Middleware Tests
// =========================================================================

package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aegisgatesecurity/aegisgate/pkg/core"
)

// ============================================================================
// Helper Functions
// ============================================================================

func newTestValidator(licenseKey string, adminPanelURL string) *LicenseValidator {
	config := &LicenseConfig{
		LicenseKey:    licenseKey,
		AdminPanelURL: adminPanelURL,
		CacheDuration: 1 * time.Minute,
		FailOpen:      true,
	}
	return NewLicenseValidator(config)
}

func mockAdminPanelServer(tierName string, valid bool, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if !valid {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"valid":   false,
				"status":  "invalid",
				"message": "License key not found or expired",
			})
			return
		}

		// Map tier names to limits
		tierLimits := map[string]map[string]interface{}{
			"community":    {"max_servers": 1, "max_users": 3, "rate_limit_per_minute": 60},
			"developer":   {"max_servers": 5, "max_users": 10, "rate_limit_per_minute": 600},
			"professional": {"max_servers": 25, "max_users": 50, "rate_limit_per_minute": 3000},
			"enterprise":  {"max_servers": -1, "max_users": -1, "rate_limit_per_minute": -1},
		}

		limits, ok := tierLimits[tierName]
		if !ok {
			limits = tierLimits["community"]
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid": true,
			"status": "valid",
			"license": map[string]interface{}{
				"tier_name":             tierName,
				"max_servers":          limits["max_servers"],
				"max_users":            limits["max_users"],
				"rate_limit_per_minute": limits["rate_limit_per_minute"],
				"expires_at":           time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
			},
		})
	}))
}

// ============================================================================
// Test Cases
// ============================================================================

// Test: Community tier (no license key)
func TestLicenseValidator_CommunityTier(t *testing.T) {
	validator := newTestValidator("", "http://localhost:8443")
	
	result, err := validator.Validate(context.Background())
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !result.Valid {
		t.Error("expected valid (community tier with no key)")
	}
	
	if result.Tier != core.TierCommunity {
		t.Errorf("expected community tier, got %v", result.Tier)
	}
	
	if result.RateLimit != 60 {
		t.Errorf("expected rate limit 60, got %d", result.RateLimit)
	}
}

// Test: Valid Developer license
func TestLicenseValidator_DeveloperLicense(t *testing.T) {
	mockServer := mockAdminPanelServer("developer", true, http.StatusOK)
	defer mockServer.Close()

	validator := newTestValidator("DEV-XXXX-XXXX", mockServer.URL)
	validator.ClearCache()

	result, err := validator.Validate(context.Background())
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !result.Valid {
		t.Error("expected valid license")
	}
	
	if result.Tier != core.TierDeveloper {
		t.Errorf("expected developer tier, got %v", result.Tier)
	}
	
	if result.MaxServers != 5 {
		t.Errorf("expected max servers 5, got %d", result.MaxServers)
	}
	
	if result.MaxUsers != 10 {
		t.Errorf("expected max users 10, got %d", result.MaxUsers)
	}
	
	if result.RateLimit != 600 {
		t.Errorf("expected rate limit 600, got %d", result.RateLimit)
	}
}

// Test: Valid Professional license
func TestLicenseValidator_ProfessionalLicense(t *testing.T) {
	mockServer := mockAdminPanelServer("professional", true, http.StatusOK)
	defer mockServer.Close()

	validator := newTestValidator("PRO-XXXX-XXXX", mockServer.URL)
	validator.ClearCache()

	result, err := validator.Validate(context.Background())
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !result.Valid {
		t.Error("expected valid license")
	}
	
	if result.Tier != core.TierProfessional {
		t.Errorf("expected professional tier, got %v", result.Tier)
	}
	
	if result.MaxServers != 25 {
		t.Errorf("expected max servers 25, got %d", result.MaxServers)
	}
	
	if result.RateLimit != 3000 {
		t.Errorf("expected rate limit 3000, got %d", result.RateLimit)
	}
}

// Test: Valid Enterprise license
func TestLicenseValidator_EnterpriseLicense(t *testing.T) {
	mockServer := mockAdminPanelServer("enterprise", true, http.StatusOK)
	defer mockServer.Close()

	validator := newTestValidator("ENT-XXXX-XXXX", mockServer.URL)
	validator.ClearCache()

	result, err := validator.Validate(context.Background())
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !result.Valid {
		t.Error("expected valid license")
	}
	
	if result.Tier != core.TierEnterprise {
		t.Errorf("expected enterprise tier, got %v", result.Tier)
	}
	
	if result.MaxServers != -1 {
		t.Errorf("expected unlimited max servers (-1), got %d", result.MaxServers)
	}
	
	if result.MaxUsers != -1 {
		t.Errorf("expected unlimited max users (-1), got %d", result.MaxUsers)
	}
	
	if result.RateLimit != -1 {
		t.Errorf("expected unlimited rate limit (-1), got %d", result.RateLimit)
	}
}

// Test: License service unavailable (fail_open = true)
func TestLicenseValidator_ServiceUnavailable_FailOpen(t *testing.T) {
	validator := newTestValidator("TEST-KEY", "http://localhost:19999")
	validator.config.FailOpen = true

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := validator.Validate(ctx)
	
	t.Logf("Result: %+v, Error: %v", result, err)
}

// Test: Caching works correctly
func TestLicenseValidator_Caching(t *testing.T) {
	mockServer := mockAdminPanelServer("developer", true, http.StatusOK)
	defer mockServer.Close()

	validator := newTestValidator("DEV-KEY-12345", mockServer.URL)
	validator.ClearCache()

	// First validation
	result1, err := validator.Validate(context.Background())
	if err != nil {
		t.Fatalf("first validation failed: %v", err)
	}

	// Second validation should use cache
	result2, err := validator.Validate(context.Background())
	if err != nil {
		t.Fatalf("second validation failed: %v", err)
	}

	if result1.Tier != result2.Tier {
		t.Errorf("cached result differs: %v vs %v", result1.Tier, result2.Tier)
	}
}

// Test: ClearCache works
func TestLicenseValidator_ClearCache(t *testing.T) {
	mockServer := mockAdminPanelServer("developer", true, http.StatusOK)
	defer mockServer.Close()

	validator := newTestValidator("DEV-KEY", mockServer.URL)
	
	validator.Validate(context.Background())
	validator.ClearCache()
	
	result, err := validator.Validate(context.Background())
	if err != nil {
		t.Fatalf("validation after clear failed: %v", err)
	}
	
	if result == nil {
		t.Error("expected result after cache clear")
	}
}

// Test: GetTier method
func TestLicenseValidator_GetTier(t *testing.T) {
	mockServer := mockAdminPanelServer("professional", true, http.StatusOK)
	defer mockServer.Close()

	validator := newTestValidator("PRO-KEY", mockServer.URL)

	tier := validator.GetTier(context.Background())
	
	if tier != core.TierProfessional {
		t.Errorf("expected professional tier, got %v", tier)
	}
}

// Test: Context helpers
func TestContextHelpers(t *testing.T) {
	ctx := context.Background()
	
	// Test GetLicenseTierFromContext with no tier set
	tier := GetLicenseTierFromContext(ctx)
	if tier != core.TierCommunity {
		t.Errorf("expected default community tier, got %v", tier)
	}
	
	// Test GetLicenseRateLimitFromContext with no limit set
	rateLimit := GetLicenseRateLimitFromContext(ctx)
	if rateLimit != 60 {
		t.Errorf("expected default rate limit 60, got %d", rateLimit)
	}
	
	// Test with values set in context
	ctx = context.WithValue(ctx, "license_tier", core.TierDeveloper)
	ctx = context.WithValue(ctx, "license_rate_limit", 600)
	
	tier = GetLicenseTierFromContext(ctx)
	if tier != core.TierDeveloper {
		t.Errorf("expected developer tier from context, got %v", tier)
	}
	
	rateLimit = GetLicenseRateLimitFromContext(ctx)
	if rateLimit != 600 {
		t.Errorf("expected rate limit 600 from context, got %d", rateLimit)
	}
}

// Test: HTTP Middleware routes correctly
func TestLicenseMiddleware_Routes(t *testing.T) {
	mockServer := mockAdminPanelServer("developer", true, http.StatusOK)
	defer mockServer.Close()
	
	validator := newTestValidator("DEV-KEY", mockServer.URL)
	originalGlobalValidator := globalValidator
	globalValidator = validator
	defer func() { globalValidator = originalGlobalValidator }()

	middleware := LicenseMiddleware()
	
	// Test that /health bypasses license check
	healthCalled := false
	healthHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		healthCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	
	handler := middleware(healthHandler)
	
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	
	if !healthCalled {
		t.Error("health endpoint should bypass middleware")
	}
}

// Test: Tier string conversion
func TestTierString(t *testing.T) {
	tests := [...]struct {
		tier    core.Tier
		want    string
	}{
		{core.TierCommunity, "Community"},
		{core.TierDeveloper, "Developer"},
		{core.TierProfessional, "Professional"},
		{core.TierEnterprise, "Enterprise"},
		{core.Tier(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.tier.String(); got != tt.want {
				t.Errorf("Tier.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: GetTierByName case insensitive
func TestGetTierByName(t *testing.T) {
	tests := [...]struct {
		name    string
		want    core.Tier
	}{
		{"community", core.TierCommunity},
		{"COMMUNITY", core.TierCommunity},
		{"Community", core.TierCommunity},
		{"developer", core.TierDeveloper},
		{"professional", core.TierProfessional},
		{"enterprise", core.TierEnterprise},
		{"unknown", core.TierCommunity},
		{"", core.TierCommunity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := core.GetTierByName(tt.name); got != tt.want {
				t.Errorf("GetTierByName(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

// Test: DefaultLicenseConfig
func TestDefaultLicenseConfig(t *testing.T) {
	config := DefaultLicenseConfig()
	
	if config.AdminPanelURL != "http://localhost:8443" {
		t.Errorf("expected admin panel URL, got %s", config.AdminPanelURL)
	}
	
	if config.CacheDuration != 5*time.Minute {
		t.Errorf("expected cache duration 5m, got %v", config.CacheDuration)
	}
	
	if !config.FailOpen {
		t.Error("expected fail-open by default")
	}
}

// Test: HTTP paths bypass license check
func TestLicenseMiddleware_HealthVersionStatsBypass(t *testing.T) {
	mockServer := mockAdminPanelServer("developer", true, http.StatusOK)
	defer mockServer.Close()

	validator := newTestValidator("DEV-KEY", mockServer.URL)
	originalGlobalValidator := globalValidator
	globalValidator = validator
	defer func() { globalValidator = originalGlobalValidator }()

	middleware := LicenseMiddleware()

	testPaths := []string{"/health", "/healthz", "/version", "/stats"}

	for _, path := range testPaths {
		called := false
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", path, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if !called {
			t.Errorf("path %s should bypass license check", path)
		}
	}
}

// Test: Tier marshaling
func TestTierJSON(t *testing.T) {
	// Test MarshalJSON
	data, err := json.Marshal(core.TierDeveloper)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	
	if !strings.Contains(string(data), "developer") {
		t.Errorf("expected developer in JSON, got %s", string(data))
	}

	// Test UnmarshalJSON
	var tier core.Tier
	err = json.Unmarshal([]byte("\"professional\""), &tier)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	
	if tier != core.TierProfessional {
		t.Errorf("expected professional tier, got %v", tier)
	}
}
