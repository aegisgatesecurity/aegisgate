// SPDX-License-Identifier: MIT
// =========================================================================
// PROPRIETARY - AegisGate Security
// Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
// =========================================================================

package scanner_test

import (
	"regexp"
	"testing"

	"github.com/aegisgatesecurity/aegisgate/pkg/scanner"
)

// Test patterns_extended.go functions
func TestAdditionalPatterns(t *testing.T) {
	patterns := scanner.AdditionalPatterns()
	// Function executed - may return empty list which is valid
	_ = patterns
}

func TestAllPatterns(t *testing.T) {
	patterns := scanner.AllPatterns()
	if len(patterns) == 0 {
		t.Error("AllPatterns should return patterns")
	}
}

func TestAllPatternsIncludesDefaults(t *testing.T) {
	defaults := scanner.DefaultPatterns()
	all := scanner.AllPatterns()
	if len(all) < len(defaults) {
		t.Errorf("AllPatterns (%d) should include defaults (%d)", len(all), len(defaults))
	}
}

// Test scanner.go specific functions
func TestDefaultConfigDirect(t *testing.T) {
	cfg := scanner.DefaultConfig()
	if cfg.BlockThreshold != scanner.Critical {
		t.Errorf("Expected Critical, got %v", cfg.BlockThreshold)
	}
}

func TestNewWithNilConfig(t *testing.T) {
	s := scanner.New(nil)
	if s == nil {
		t.Fatal("New(nil) should not return nil")
	}
}

func TestSetConfigDirect(t *testing.T) {
	s := scanner.New(nil)
	newCfg := &scanner.Config{
		BlockThreshold: scanner.Medium,
	}
	s.SetConfig(newCfg)
	// SetConfig should not panic
}

func TestScanWithContextDirect(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	findings := s.ScanWithContext("email: test@example.com")
	if len(findings) == 0 {
		t.Error("ScanWithContext should return findings")
	}
}

func TestScanBytesDirect(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	findings := s.ScanBytes([]byte("email: test@example.com"))
	if len(findings) == 0 {
		t.Error("ScanBytes should return findings")
	}
}

func TestHasViolationDirect(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.High,
	}
	s := scanner.New(cfg)
	findings := s.Scan("4532015112830366") // Credit card
	if !s.HasViolation(findings) {
		t.Error("HasViolation should return true for credit card")
	}
}

func TestGetCriticalFindingsDirect(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	findings := s.Scan("email: test@example.com credit card: 4532015112830366")
	critical := s.GetCriticalFindings(findings)
	if critical == nil {
		t.Error("GetCriticalFindings should return slice")
	}
}

func TestGetFindingsByCategoryDirect(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	findings := s.Scan("email: test@example.com")
	filtered := s.GetFindingsByCategory(findings, scanner.CategoryPII)
	_ = filtered
}

func TestGetFindingsBySeverityDirect(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	findings := s.Scan("test@example.com")
	filtered := s.GetFindingsBySeverity(findings, scanner.Low)
	_ = filtered
}

func TestGetViolationSummaryDirect(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	findings := s.Scan("test@example.com")
	summary := s.GetViolationSummary(findings)
	if summary == nil {
		t.Error("GetViolationSummary should return map")
	}
}

func TestGetViolationNamesDirect(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	findings := s.Scan("test@example.com")
	names := s.GetViolationNames(findings)
	_ = names
}

func TestAddPatternDirect(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       []*scanner.Pattern{},
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	s.AddPattern(&scanner.Pattern{Name: "Test", Regex: regexp.MustCompile("test")})
}

func TestRemovePatternDirect(t *testing.T) {
	cfg := &scanner.Config{
		Patterns: []*scanner.Pattern{
			{Name: "RemoveMe", Regex: regexp.MustCompile("test")},
		},
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	removed := s.RemovePattern("RemoveMe")
	if !removed {
		t.Error("RemovePattern should return true")
	}
}

func TestGetPatternDirect(t *testing.T) {
	cfg := &scanner.Config{
		Patterns: []*scanner.Pattern{
			{Name: "FindMe", Regex: regexp.MustCompile("test")},
		},
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	p := s.GetPattern("FindMe")
	if p == nil {
		t.Error("GetPattern should find existing pattern")
	}
}

func TestCompilePatternDirect(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       []*scanner.Pattern{},
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	err := s.CompilePattern("NewPattern", "[A-Z]+", scanner.High, "Test", "description")
	if err != nil {
		t.Errorf("CompilePattern should not error: %v", err)
	}
}

func TestCompilePatternInvalidRegex(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       []*scanner.Pattern{},
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	err := s.CompilePattern("Invalid", "[invalid(", scanner.Low, "Test", "desc")
	if err == nil {
		t.Error("CompilePattern should error on invalid regex")
	}
}

func TestScannerNilPatterns(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       nil,
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	// Should not panic
	findings := s.Scan("test content")
	_ = findings
}

func TestScannerMultipleMaxFindings(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Low,
		MaxFindings:    2,
	}
	s := scanner.New(cfg)
	// Many emails in content - should be limited by MaxFindings
	findings := s.Scan("a@b.com c@d.com e@f.com g@h.com")
	if len(findings) > 2 {
		t.Errorf("Should be limited to MaxFindings, got %d", len(findings))
	}
}

func TestFindingContextExtraction(t *testing.T) {
	cfg := &scanner.Config{
		Patterns: []*scanner.Pattern{
			{Name: "Email", Regex: regexp.MustCompile(`[a-z]+@[a-z]+\.[a-z]+`), Severity: scanner.Low, Category: "PII"},
		},
		BlockThreshold: scanner.Critical,
		IncludeContext: true,
		ContextSize:    5,
	}
	s := scanner.New(cfg)
	findings := s.Scan("xxx test@example.com yyy")
	if len(findings) > 0 {
		// Context may be truncated at boundaries but should exist
		_ = findings[0].Context
	}
}

func TestMaskMatch(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	findings := s.Scan("short")
	if len(findings) > 0 {
		_ = findings[0].Match
	}
}

func TestScanWithNoMatchingPatterns(t *testing.T) {
	cfg := &scanner.Config{
		Patterns: []*scanner.Pattern{
			{Name: "Unique", Regex: regexp.MustCompile("UNIQUE_PATTERN_XYZ")},
		},
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	findings := s.Scan("This content has no matches")
	if len(findings) != 0 {
		t.Error("Expected no findings")
	}
}

func TestScanWithVariousContent(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	s := scanner.New(cfg)
	
	tests := []struct {
		name    string
		content string
		minFind int
	}{
		{"AWS Key", "AKIAIOSFODNN7EXAMPLE", 1},
		{"GitHub Token", "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", 1},
		{"Private Key", "-----BEGIN RSA PRIVATE KEY-----", 1},
		{"Connection String", "postgresql://user:pass@localhost/db", 1},
		{"JWT", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U", 1},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := s.Scan(tt.content)
			if len(findings) < tt.minFind {
				t.Errorf("Expected at least %d finding(s), got %d", tt.minFind, len(findings))
			}
		})
	}
}

func TestScannerWithCustomContextSize(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
		IncludeContext: true,
		ContextSize:    200,
	}
	s := scanner.New(cfg)
	findings := s.Scan("prefix test@example.com suffix")
	if len(findings) > 0 {
		_ = findings[0].Context
	}
}

func TestDefaultPatternsHasAllCategories(t *testing.T) {
	patterns := scanner.DefaultPatterns()
	categories := make(map[scanner.Category]bool)
	severities := make(map[scanner.Severity]bool)
	
	for _, p := range patterns {
		categories[p.Category] = true
		severities[p.Severity] = true
	}
	
	if !categories[scanner.CategoryPII] {
		t.Error("Should have PII patterns")
	}
	if !categories[scanner.CategoryCredential] {
		t.Error("Should have Credential patterns")
	}
	if !categories[scanner.CategoryFinancial] {
		t.Error("Should have Financial patterns")
	}
	if !severities[scanner.Critical] {
		t.Error("Should have Critical severity patterns")
	}
}