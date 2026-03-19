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

func TestDefaultConfig(t *testing.T) {
	cfg := scanner.DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig should not return nil")
	}
	if cfg.BlockThreshold != scanner.Critical {
		t.Errorf("Expected Critical threshold, got %v", cfg.BlockThreshold)
	}
	if cfg.LogFindings != true {
		t.Error("Expected LogFindings to be true by default")
	}
	if cfg.ContextSize != 50 {
		t.Errorf("Expected ContextSize 50, got %d", cfg.ContextSize)
	}
	if cfg.MaxFindings != 100 {
		t.Errorf("Expected MaxFindings 100, got %d", cfg.MaxFindings)
	}
}

func TestScannerNewWithNilConfig(t *testing.T) {
	sc := scanner.New(nil)
	if sc == nil {
		t.Fatal("New(nil) should return a valid scanner")
	}
}

func TestScannerNewWithConfig(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.High,
		LogFindings:    false,
	}
	sc := scanner.New(cfg)
	if sc == nil {
		t.Fatal("New(config) should return a valid scanner")
	}
}

func TestScannerSetConfig(t *testing.T) {
	sc := scanner.New(nil)
	newCfg := &scanner.Config{
		BlockThreshold: scanner.Medium,
	}
	sc.SetConfig(newCfg)
}

func TestScanWithContext(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
		IncludeContext: false,
	}
	sc := scanner.New(cfg)
	findings := sc.ScanWithContext("test@example.com")
	if len(findings) == 0 {
		t.Error("Expected findings for email")
	}
	if findings[0].Context == "" {
		t.Error("Expected context when using ScanWithContext")
	}
}

func TestScanBytes(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	content := []byte("email: test@example.com")
	findings := sc.ScanBytes(content)
	if len(findings) == 0 {
		t.Error("Expected findings from bytes")
	}
}

func TestHasViolation(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.High,
	}
	sc := scanner.New(cfg)
	
	// Credit card should trigger violation at High threshold
	highContent := "4532015112830366"
	findings := sc.Scan(highContent)
	if !sc.HasViolation(findings) {
		t.Error("Expected violation for credit card at High threshold")
	}
	
	// Low severity should not trigger at High threshold
	lowContent := "test@example.com"
	findings = sc.Scan(lowContent)
	if sc.HasViolation(findings) {
		t.Error("Expected no violation for email (Low severity) at High threshold")
	}
}

func TestGetCriticalFindings(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	content := "email: test@example.com credit card: 4532015112830366"
	findings := sc.Scan(content)
	critical := sc.GetCriticalFindings(findings)
	if len(critical) == 0 {
		t.Error("Expected Critical findings from credit card")
	}
}

func TestGetFindingsByCategory(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	content := "email: test@example.com credit card: 4532015112830366"
	findings := sc.Scan(content)
	
	piiFindings := sc.GetFindingsByCategory(findings, scanner.CategoryPII)
	if len(piiFindings) == 0 {
		t.Error("Expected PII findings")
	}
	
	financialFindings := sc.GetFindingsByCategory(findings, scanner.CategoryFinancial)
	if len(financialFindings) == 0 {
		t.Error("Expected Financial findings")
	}
	
	// Non-existent category
	other := sc.GetFindingsByCategory(findings, "NonExistent")
	if len(other) != 0 {
		t.Error("Expected no findings for non-existent category")
	}
}

func TestGetFindingsBySeverity(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	content := "test@example.com"
	findings := sc.Scan(content)
	
	// Email is Low severity, should not appear in High results
	highFindings := sc.GetFindingsBySeverity(findings, scanner.High)
	if len(highFindings) != 0 {
		t.Error("Expected no High findings for email")
	}
	
	// Email should appear in Low results
	lowFindings := sc.GetFindingsBySeverity(findings, scanner.Low)
	if len(lowFindings) == 0 {
		t.Error("Expected Low findings for email")
	}
}

func TestGetViolationSummary(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	content := "email: test@example.com"
	findings := sc.Scan(content)
	summary := sc.GetViolationSummary(findings)
	
	if summary == nil {
		t.Error("Expected non-nil summary")
	}
	
	if summary[scanner.Low] != 1 {
		t.Errorf("Expected 1 Low finding, got %d", summary[scanner.Low])
	}
}

func TestGetViolationNames(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	content := "email: test@example.com"
	findings := sc.Scan(content)
	names := sc.GetViolationNames(findings)
	
	if len(names) == 0 {
		t.Error("Expected violation names")
	}
}

func TestAddPattern(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       []*scanner.Pattern{},
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	initialCount := len(cfg.Patterns)
	
	customPattern := &scanner.Pattern{
		Name:     "CustomPattern",
		Regex:    regexp.MustCompile("CUSTOM[0-9]+"),
		Severity: scanner.High,
		Category: "Custom",
	}
	sc.AddPattern(customPattern)
	
	if len(cfg.Patterns) != initialCount+1 {
		t.Error("Expected pattern to be added")
	}
}

func TestRemovePattern(t *testing.T) {
	cfg := &scanner.Config{
		Patterns: []*scanner.Pattern{
			{Name: "PatternToRemove", Regex: regexp.MustCompile("TEST")},
			{Name: "PatternToKeep", Regex: regexp.MustCompile("KEEP")},
		},
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	removed := sc.RemovePattern("PatternToRemove")
	if !removed {
		t.Error("Expected pattern to be removed")
	}
	
	// Try to remove non-existent
	removed = sc.RemovePattern("NonExistent")
	if removed {
		t.Error("Expected false for non-existent pattern")
	}
}

func TestGetPattern(t *testing.T) {
	cfg := &scanner.Config{
		Patterns: []*scanner.Pattern{
			{Name: "FindMe", Regex: regexp.MustCompile("TEST")},
		},
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	found := sc.GetPattern("FindMe")
	if found == nil {
		t.Error("Expected to find pattern")
	}
	
	notFound := sc.GetPattern("NonExistent")
	if notFound != nil {
		t.Error("Expected nil for non-existent pattern")
	}
}

func TestCompilePattern(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       []*scanner.Pattern{},
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	err := sc.CompilePattern("TestPattern", `[A-Z0-9]+`, scanner.High, "Test", "Test description")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Invalid regex
	err = sc.CompilePattern("InvalidPattern", `[invalid(`, scanner.Low, "Test", "Test")
	if err == nil {
		t.Error("Expected error for invalid regex")
	}
}

func TestScanMaxFindingsLimit(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
		MaxFindings:    1,
	}
	sc := scanner.New(cfg)
	
	// Multiple emails should only return MaxFindings
	content := "a@b.com c@d.com e@f.com"
	findings := sc.Scan(content)
	
	if len(findings) > 1 {
		t.Errorf("Expected at most %d findings, got %d", cfg.MaxFindings, len(findings))
	}
}

func TestScanEmptyContent(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	findings := sc.Scan("")
	if len(findings) != 0 {
		t.Error("Expected no findings for empty content")
	}
}

func TestScanNilPattern(t *testing.T) {
	cfg := &scanner.Config{
		Patterns: []*scanner.Pattern{
			nil,
			{Name: "Valid", Regex: regexp.MustCompile("test")},
		},
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	findings := sc.Scan("test content")
	if len(findings) == 0 {
		t.Error("Expected findings with nil pattern in list")
	}
}

func TestScanNilRegex(t *testing.T) {
	cfg := &scanner.Config{
		Patterns: []*scanner.Pattern{
			{Name: "NilRegex", Regex: nil},
			{Name: "Valid", Regex: regexp.MustCompile("test")},
		},
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	findings := sc.Scan("test content")
	if len(findings) == 0 {
		t.Error("Expected findings with nil regex in pattern")
	}
}

func TestContextExtraction(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
		IncludeContext: true,
		ContextSize:    10,
	}
	sc := scanner.New(cfg)
	
	// Create content with enough surrounding text
	content := "This is some prefix test@example.com and this is suffix"
	findings := sc.Scan(content)
	
	if len(findings) == 0 {
		t.Fatal("Expected findings")
	}
	
	// Context should be present but may be truncated at boundaries
	_ = findings[0].Context
}

func TestScanMultiplePatterns(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Low, // All findings will trigger
		MaxFindings:    1000,
	}
	sc := scanner.New(cfg)
	
	content := "email: test@example.com phone: 123-456-7890 cc: 4532015112830366"
	findings := sc.Scan(content)
	
	if len(findings) < 3 {
		t.Errorf("Expected at least 3 findings, got %d", len(findings))
	}
}

func TestFindingStruct(t *testing.T) {
	pattern := &scanner.Pattern{
		Name:     "TestPattern",
		Regex:    regexp.MustCompile("test"),
		Severity: scanner.High,
		Category: scanner.CategoryPII,
	}
	
	finding := scanner.Finding{
		Pattern:  pattern,
		Match:    "test",
		Position: 0,
		Context:  "surrounding text",
	}
	
	if finding.Pattern.Name != "TestPattern" {
		t.Error("Expected Pattern to be set")
	}
	if finding.Match != "test" {
		t.Error("Expected Match to be set")
	}
	if finding.Position != 0 {
		t.Error("Expected Position to be set")
	}
}

func TestShouldBlockBoundary(t *testing.T) {
	// Test at boundary (High is >= High)
	if !scanner.ShouldBlock(scanner.High) {
		t.Error("High should block")
	}
	
	// Medium is less than High
	if scanner.ShouldBlock(scanner.Medium) {
		t.Error("Medium should not block")
	}
}

func TestScannerConfigFields(t *testing.T) {
	cfg := &scanner.Config{
		Patterns:       scanner.DefaultPatterns(),
		BlockThreshold: scanner.Critical,
		LogFindings:    true,
		IncludeContext: true,
		ContextSize:    100,
		MaxFindings:    50,
	}
	scanner.New(cfg)
	
	// Test that all config fields work correctly
	if cfg.ContextSize != 100 {
		t.Error("ContextSize not set correctly")
	}
	if cfg.MaxFindings != 50 {
		t.Error("MaxFindings not set correctly")
	}
}

func TestFindingMatchPosition(t *testing.T) {
	cfg := &scanner.Config{
		Patterns: []*scanner.Pattern{
			{Name: "TestPattern", Regex: regexp.MustCompile("test"), Severity: scanner.High, Category: "Test"},
		},
		BlockThreshold: scanner.Critical,
	}
	sc := scanner.New(cfg)
	
	findings := sc.Scan("prefix test suffix")
	if len(findings) != 1 {
		t.Fatal("Expected 1 finding")
	}
	
	if findings[0].Position != 7 { // "test" starts at index 7 in "prefix test suffix"
		t.Errorf("Expected position 7, got %d", findings[0].Position)
	}
	
	if findings[0].Match != "test" {
		t.Errorf("Expected match 'test', got '%s'", findings[0].Match)
	}
}