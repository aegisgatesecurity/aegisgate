// Package scanner_test provides unit tests for the pattern detection functionality.
//
//go:build !integration
// +build !integration

package scanner_test

import (
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

func TestPattern_GitHubToken(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	findings := sc.Scan(content)
	if len(findings) == 0 {
		t.Error("Expected to find GitHub token")
	}
}

func TestPattern_JWTToken(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"
	findings := sc.Scan(content)
	if len(findings) == 0 {
		t.Error("Expected to find JWT token")
	}
}

func TestPattern_SlackToken(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "xoxb-TESTONLY-000000000000-TESTINGEXAMPLE"
	findings := sc.Scan(content)
	if len(findings) == 0 {
		t.Error("Expected to find Slack token")
	}
}

func TestPattern_SSN(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "123-45-6789"
	findings := sc.Scan(content)
	if len(findings) == 0 {
		t.Error("Expected to find SSN")
	}
}

func TestPattern_EmailAddress(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "test@aegisgatesecurity.io"
	findings := sc.Scan(content)
	if len(findings) == 0 {
		t.Error("Expected to find email")
	}
}

func TestPattern_PrivateKeys(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "-----BEGIN RSA PRIVATE KEY-----"
	findings := sc.Scan(content)
	if len(findings) == 0 {
		t.Error("Expected to find private key")
	}
}

func TestPattern_PostgreSQL(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "postgresql://user:pass@localhost:5432/db"
	findings := sc.Scan(content)
	if len(findings) == 0 {
		t.Error("Expected to find PostgreSQL connection string")
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

func TestCleanContent(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	content := "This is normal text with no sensitive data"
	findings := sc.Scan(content)

	// Should not have any critical findings
	for _, f := range findings {
		if f.Severity >= scanner.High {
			t.Errorf("Unexpected high severity finding: %v", f.Name)
		}
	}
}

func TestJSONPayload(t *testing.T) {
	cfg := testConfig()
	sc := scanner.New(cfg)

	// Simple JSON without quotes inside
	content := "{\"action\": \"login\", \"userId\": 12345}"
	findings := sc.Scan(content)
	_ = findings
}
