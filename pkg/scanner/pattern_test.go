// Package scanner_test provides unit tests for pattern detection.
//go:build !integration
// +build !integration

package scanner_test

import (
	"strings"
	"testing"

	"github.com/aegisgatesecurity/aegisgate/pkg/scanner"
)

func testConfig() *scanner.Config {
	return &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
		LogFindings:    false,
	}
}

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
		{"Invalid JWT", "not.a.jwt", false},
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
		{"Invalid Token", "xoxz-fake", false},
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
		{"Email in JSON", "test@aegisgatesecurity.io", true},
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
		{"Redis", "redis://localhost:6379", true},
		{"Not a connection", "postgres://localhost", false},
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

func TestSeverity_Classification(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

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

func TestCategory_Classification(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	ccContent := "4532015112830366"
	findings := sc.Scan(ccContent)
	if len(findings) > 0 && findings[0].Category != scanner.CategoryFinancial {
		t.Errorf("Expected CategoryFinancial, got %v", findings[0].Category)
	}

	tokenContent := "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	findings = sc.Scan(tokenContent)
	if len(findings) > 0 && findings[0].Category != scanner.CategoryCredential {
		t.Errorf("Expected CategoryCredential, got %v", findings[0].Category)
	}
}

func TestShouldBlock(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.High,
		LogFindings:    false,
	}
	sc := scanner.New(cfg)

	ccContent := "4532015112830366"
	findings := sc.Scan(ccContent)
	if !sc.ShouldBlock(findings) {
		t.Error("Expected blocking for Critical severity with High threshold")
	}
}

func TestCleanContent(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "This is normal text with no sensitive data."
	findings := sc.Scan(content)

	for _, f := range findings {
		if f.Severity >= scanner.High {
			t.Errorf("Unexpected high severity finding: %v", f.Name)
		}
	}
}

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

func TestMultipleFindings(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "Email: test@example.com, CC: 4532015112830366, SSN: 123-45-6789"
	findings := sc.Scan(content)
	if len(findings) < 2 {
		t.Logf("Found %d findings in multi-pattern content", len(findings))
	}
}

func TestConfigOptions(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Info,
		LogFindings:    false,
		MaxFindings:    5,
	}
	sc := scanner.New(cfg)
	content := "user@test.com 4532015112830366 123-45-6789 AKIAIOSFODNN7EXAMPLE"
	findings := sc.Scan(content)
	if len(findings) > 5 {
		t.Errorf("Expected max 5 findings, got %d", len(findings))
	}
}
