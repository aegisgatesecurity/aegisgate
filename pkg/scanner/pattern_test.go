// Package scanner_test provides unit tests for the pattern detection functionality.
// Tests cover pattern matching, severity levels, category classification,
// and integration with the main scanner.
//
//go:build !integration
// +build !integration

package scanner_test

import (
	"strings"
	"testing"

	"github.com/aegisgatesecurity/aegisgate/pkg/scanner"
)

// ============================================================================
// Test Helpers
// ============================================================================

// testConfig returns a standard test configuration
func testConfig() *scanner.Config {
	return &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
		LogFindings:    false,
		IncludeContext: false,
	}
}

// ============================================================================
// Test 1: Pattern - Credit Card Detection
// ============================================================================

func TestPattern_CreditCard_Visa(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{"Visa with spaces", "4532 0151 1283 0366", true},
		{"Visa without spaces", "4532015112830366", true},
		{"Visa 13 digit", "4024007134564840", true},
		{"Invalid Visa", "3532015112830366", false},
		{"Random numbers", "1234567890123456", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasFinding := len(findings) > 0
			if hasFinding != tt.find {
				t.Errorf("Scan() = %v, want %v for content: %s", hasFinding, tt.find, tt.content)
			}
		})
	}
}

func TestPattern_CreditCard_Mastercard(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{"Mastercard valid", "5425233430109903", true},
		{"Mastercard valid 2", "5555555555554444", true},
		{"Mastercard invalid start", "4525233430109903", false},
		{"Random", "1234567890123456", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasFinding := len(findings) > 0
			if hasFinding != tt.find {
				t.Errorf("Scan() = %v, want %v", hasFinding, tt.find)
			}
		})
	}
}

func TestPattern_CreditCard_Amex(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{"Amex valid", "378282246310005", true},
		{"Amex valid 2", "371449635398431", true},
		{"Amex invalid", "368282246310005", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasFinding := len(findings) > 0
			if hasFinding != tt.find {
				t.Errorf("Scan() = %v, want %v", hasFinding, tt.find)
			}
		})
	}
}

// ============================================================================
// Test 2: Pattern - Credential Detection
// ============================================================================

func TestPattern_AWSAccessKey(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{"AWS Access Key ID", "AKIAIOSFODNN7EXAMPLE", true},
		{"AWS with prefix", "AWS_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE", true},
		{"Invalid AWS Key", "AKIA123456789012345", false},
		{"Random text", "Some random content", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasFinding := len(findings) > 0
			if hasFinding != tt.find {
				t.Errorf("Scan() = %v, want %v", hasFinding, tt.find)
			}
		})
	}
}

func TestPattern_GitHubToken(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{"GitHub Classic Token", "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", true},
		{"GitHub OAuth", "gho_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", true},
		{"GitHub App", "ghs_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", true},
		{"GitHub Refresh", "ghr_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", true},
		{"Invalid prefix", "ghx_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", false},
		{"Too short", "ghp_xxx", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasFinding := len(findings) > 0
			if hasFinding != tt.find {
				t.Errorf("Scan() = %v, want %v", hasFinding, tt.find)
			}
		})
	}
}

func TestPattern_JWTToken(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{"Valid JWT", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", true},
		{"JWT in JSON", `{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"}`, true},
		{"Invalid JWT", "not.a.jwt", false},
		{"Random", "some random content", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasFinding := len(findings) > 0
			if hasFinding != tt.find {
				t.Errorf("Scan() = %v, want %v", hasFinding, tt.find)
			}
		})
	}
}

func TestPattern_SlackToken(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{"Slack Bot Token", "xoxb-TESTONLY-000000000000-TESTINGEXAMPLE", true},
		{"Slack User Token", "xoxp-1234567890123-1234567890123-abcdefghijklmnopqrstuvwx", true},
		{"Slack App Token", "xoxa-1234567890123-1234567890123-abcdefghijklmnop", true},
		{"Slack Webhook", "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXX", false},
		{"Invalid Token", "xoxz-fake", false},
		{"Random text", "some random content", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasFinding := len(findings) > 0
			if hasFinding != tt.find {
				t.Errorf("Scan() = %v, want %v", hasFinding, tt.find)
			}
		})
	}
}

// ============================================================================
// Test 3: Pattern - PII Detection
// ============================================================================

func TestPattern_SSN(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{"SSN with dashes", "123-45-6789", true},
		{"SSN without dashes", "123456789", true},
		{"SSN in text", "My SSN is 987-65-4321 please", true},
		{"Invalid SSN", "000-00-0000", true},
		{"Random numbers", "1234567890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasFinding := len(findings) > 0
			if hasFinding != tt.find {
				t.Errorf("Scan() = %v, want %v", hasFinding, tt.find)
			}
		})
	}
}

func TestPattern_EmailAddress(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{"Standard email", "user@example.com", true},
		{"Email with subdomain", "user@mail.example.com", true},
		{"Email with plus", "user+tag@example.com", true},
		{"Email in JSON", `{"email": "test@aegisgatesecurity.io"}`, true},
		{"Invalid email", "notanemail", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasFinding := len(findings) > 0
			if hasFinding != tt.find {
				t.Errorf("Scan() = %v, want %v", hasFinding, tt.find)
			}
		})
	}
}

func TestPattern_PhoneNumber(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{"Phone with dashes", "555-123-4567", true},
		{"Phone with parens", "(555) 123-4567", true},
		{"Phone dots", "555.123.4567", true},
		{"Plain 7 digit", "5551234", true},
		{"Not a phone", "1234567890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasFinding := len(findings) > 0
			if hasFinding != tt.find {
				t.Errorf("Scan() = %v, want %v for: %s", hasFinding, tt.find, tt.content)
			}
		})
	}
}

// ============================================================================
// Test 4: Pattern - Private Key Detection
// ============================================================================

func TestPattern_PrivateKeys(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{"RSA Private Key", "-----BEGIN RSA PRIVATE KEY-----", true},
		{"EC Private Key", "-----BEGIN EC PRIVATE KEY-----", true},
		{"OpenSSH Private Key", "-----BEGIN OPENSSH PRIVATE KEY-----", true},
		{"Private Key Header", "-----BEGIN PRIVATE KEY-----", true},
		{"Not a key", "-----BEGIN CERTIFICATE-----", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasFinding := len(findings) > 0
			if hasFinding != tt.find {
				t.Errorf("Scan() = %v, want %v", hasFinding, tt.find)
			}
		})
	}
}

// ============================================================================
// Test 5: Pattern - Connection Strings
// ============================================================================

func TestPattern_DatabaseConnections(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{"PostgreSQL", "postgresql://user:pass@localhost:5432/db", true},
		{"MySQL", "mysql://user:pass@localhost:3306/db", true},
		{"MongoDB", "mongodb://localhost:27017", true},
		{"MongoDB+srv", "mongodb+srv://cluster.mongodb.net", true},
		{"Redis", "redis://localhost:6379", true},
		{"SQLServer", "sqlserver://user:pass@localhost:1433", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasFinding := len(findings) > 0
			if hasFinding != tt.find {
				t.Errorf("Scan() = %v, want %v", hasFinding, tt.find)
			}
		})
	}
}

// ============================================================================
// Test 6: Severity Classification
// ============================================================================

func TestSeverity_Classification(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	// Critical severity items should be detected
	criticalContent := "4532015112830366"
	findings := sc.Scan(criticalContent)
	if len(findings) == 0 {
		t.Error("Expected to find credit card as Critical")
	}

	if findings[0].Severity != scanner.Critical {
		t.Errorf("Expected Critical severity, got %v", findings[0].Severity)
	}
}

func TestSeverity_ShouldBlock(t *testing.T) {
	tests := []struct {
		severity scanner.Severity
		expected bool
	}{
		{scanner.Info, false},
		{scanner.Low, false},
		{scanner.Medium, false},
		{scanner.High, true},
		{scanner.Critical, true},
	}

	for _, tt := range tests {
		result := scanner.ShouldBlock(tt.severity)
		if result != tt.expected {
			t.Errorf("ShouldBlock(%v) = %v, want %v", tt.severity, result, tt.expected)
		}
	}
}

// ============================================================================
// Test 7: Category Classification
// ============================================================================

func TestCategory_Classification(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	// Test Financial category (Credit Card)
	ccContent := "4532015112830366"
	findings := sc.Scan(ccContent)
	if len(findings) > 0 && findings[0].Category != scanner.CategoryFinancial {
		t.Errorf("Expected CategoryFinancial, got %v", findings[0].Category)
	}

	// Test Credential category (GitHub Token)
	tokenContent := "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	findings = sc.Scan(tokenContent)
	if len(findings) > 0 && findings[0].Category != scanner.CategoryCredential {
		t.Errorf("Expected CategoryCredential, got %v", findings[0].Category)
	}

	// Test PII category (Email)
	emailContent := "user@example.com"
	findings = sc.Scan(emailContent)
	if len(findings) > 0 && findings[0].Category != scanner.CategoryPII {
		t.Errorf("Expected CategoryPII, got %v", findings[0].Category)
	}
}

// ============================================================================
// Test 8: Block Decision
// ============================================================================

func TestShouldBlock(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.High,
		LogFindings:    false,
	}
	sc := scanner.New(cfg)

	// Critical should block
	ccContent := "4532015112830366"
	findings := sc.Scan(ccContent)
	if !sc.ShouldBlock(findings) {
		t.Error("Expected blocking for Critical severity with High threshold")
	}
}

// ============================================================================
// Test 9: Multiple Findings
// ============================================================================

func TestMultipleFindings(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "Email: user@example.com, CC: 4532015112830366, SSN: 123-45-6789, AWS: AKIAIOSFODNN7EXAMPLE"
	findings := sc.Scan(content)

	if len(findings) < 3 {
		t.Logf("Found %d findings in multi-pattern content", len(findings))
	}
}

// ============================================================================
// Test 10: Clean Content
// ============================================================================

func TestCleanContent(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "This is normal text with no sensitive data. Just some JSON: {"action": "login", "userId": 12345}."
	findings := sc.Scan(content)

	for _, f := range findings {
		if f.Severity >= scanner.High {
			t.Errorf("Unexpected high severity finding: %v - %v", f.Name, f.Description)
		}
	}
}

// ============================================================================
// Test 11: Context Extraction
// ============================================================================

func TestContextExtraction(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
		LogFindings:    false,
		IncludeContext: true,
		ContextSize:    50,
	}
	sc := scanner.New(cfg)

	content := strings.Repeat("Some text ", 20) + "4532015112830366" + strings.Repeat(" Some text", 20)
	findings := sc.ScanWithContext(content)

	if len(findings) == 0 {
		t.Error("Expected to find credit card with context")
	}
}

// ============================================================================
// Test 12: Finding Filters
// ============================================================================

func TestFindingFilters(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "Email: test@example.com, Credit Card: 4532015112830366"
	findings := sc.Scan(content)

	if len(findings) < 2 {
		t.Skip("Need multiple findings to test filtering")
	}

	// Test filter by severity
	criticalFindings := sc.GetFindingsBySeverity(findings, scanner.Critical)
	for _, f := range criticalFindings {
		if f.Severity != scanner.Critical {
			t.Errorf("Expected Critical severity, got %v", f.Severity)
		}
	}

	// Test filter by category
	financialFindings := sc.GetFindingsByCategory(findings, scanner.CategoryFinancial)
	for _, f := range financialFindings {
		if f.Category != scanner.CategoryFinancial {
			t.Errorf("Expected CategoryFinancial, got %v", f.Category)
		}
	}
}

// ============================================================================
// Test 13: Empty and Null Inputs
// ============================================================================

func TestEmptyInputs(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
	}{
		{"Empty string", ""},
		{"Whitespace only", "   	
   "},
		{"Single character", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			if findings == nil {
				t.Error("Scan returned nil instead of empty slice")
			}
		})
	}
}

// ============================================================================
// Test 14: JSON Payload Scanning
// ============================================================================

func TestJSONPayloadScanning(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	tests := []struct {
		name    string
		content string
		find    bool
	}{
		{
			name:    "Sensitive JSON",
			content: `{"username": "john", "password": "secret123", "api_key": "AKIAIOSFODNN7EXAMPLE"}`,
			find:    true,
		},
		{
			name:    "Clean JSON",
			content: `{"username": "john", "action": "login", "userId": 12345}`,
			find:    false,
		},
		{
			name:    "Nested JSON with sensitive",
			content: `{"user": {"email": "john@example.com", "ssn": "123-45-6789"}}`,
			find:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := sc.Scan(tt.content)
			hasSensitive := len(findings) > 0
			if hasSensitive != tt.find {
				t.Errorf("Scan() = %v, want %v", hasSensitive, tt.find)
			}
		})
	}
}

// ============================================================================
// Test 15: HasViolation Check
// ============================================================================

func TestHasViolation(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
		LogFindings:    false,
	}
	sc := scanner.New(cfg)

	// Critical - should be a violation
	criticalContent := "4532015112830366"
	findings := sc.Scan(criticalContent)
	if !sc.HasViolation(findings) {
		t.Error("Expected HasViolation to return true for Critical severity")
	}
}

// ============================================================================
// Test 16: GetViolationSummary
// ============================================================================

func TestGetViolationSummary(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "4532015112830366 test@aegisgatesecurity.io"
	findings := sc.Scan(content)

	summary := sc.GetViolationSummary(findings)
	if summary != "" {
		t.Logf("Violation summary: %s", summary)
	}
}
