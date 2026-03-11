package soc2

import (
    "fmt"
    "strings"
    "time"
)

// SOC2 implements SOC2 compliance framework
type SOC2 struct {
    Rules []Rule
}

type Rule struct {
    ID          string
    Name        string
    Description string
    Severity    string
    Framework   string
    Enabled     bool
}

func (s *SOC2) Name() string { return "SOC2" }
func (s *SOC2) Version() string { return "1.0.0" }

func (s *SOC2) LoadRules() error {
    s.Rules = []Rule{
        {
            ID:          "SOC2_SECURITY",
            Name:        "Security",
            Severity:    "CRITICAL",
            Description: "System is protected against unauthorized access",
            Framework:   "SOC2",
            Enabled:     true,
        },
        {
            ID:          "SOC2_AVAILABILITY",
            Name:        "Availability",
            Severity:    "HIGH",
            Description: "System is available for operation and use",
            Framework:   "SOC2",
            Enabled:     true,
        },
        {
            ID:          "SOC2_PROCESS_INTEGRITY",
            Name:        "Processing Integrity",
            Severity:    "HIGH",
            Description: "System processes are complete, valid, authorized",
            Framework:   "SOC2",
            Enabled:     true,
        },
        {
            ID:          "SOC2_CONFIDENTIALITY",
            Name:        "Confidentiality",
            Severity:    "HIGH",
            Description: "Confidential data is protected",
            Framework:   "SOC2",
            Enabled:     true,
        },
        {
            ID:          "SOC2_PRIVACY",
            Name:        "Privacy",
            Severity:    "HIGH",
            Description: "Personal information is protected and processed",
            Framework:   "SOC2",
            Enabled:     true,
        },
    }
    return nil
}

func (s *SOC2) CheckRequest(req interface{}) ([]Violation, error) {
    if req == nil {
        return []Violation{}, nil
    }

    violations := []Violation{}

    var content string
    switch v := req.(type) {
    case string:
        content = v
    case []byte:
        content = string(v)
    default:
        content = fmt.Sprintf("%v", v)
    }

    contentLower := strings.ToLower(content)

    // Check for security control bypass attempts
    securityIndicators := []string{
        "bypass", "override", "ignore", "disregard",
        "admin", "root", "superuser", "elevate",
    }

    for _, indicator := range securityIndicators {
        if strings.Contains(contentLower, indicator) {
            violations = append(violations, Violation{
                RuleID:      fmt.Sprintf("SOC2_SECURITY_%s", strings.ReplaceAll(strings.ToUpper(indicator), " ", "_")),
                Name:        "SOC2 - Security Control Indicator",
                Description: fmt.Sprintf("Security control indicator detected: %s", indicator),
                Severity:    "MEDIUM",
                Timestamp:   time.Now(),
            })
        }
    }

    return violations, nil
}

func (s *SOC2) CheckResponse(resp interface{}) ([]Violation, error) {
    if resp == nil {
        return []Violation{}, nil
    }

    violations := []Violation{}

    var content string
    switch v := resp.(type) {
    case string:
        content = v
    case []byte:
        content = string(v)
    default:
        content = fmt.Sprintf("%v", v)
    }

    contentLower := strings.ToLower(content)

    // Check for information disclosure
    disclosurePatterns := []string{
        "internal", "confidential", " proprietary",
        "trade secret", "restricted", "classified",
    }

    for _, pattern := range disclosurePatterns {
        if strings.Contains(contentLower, pattern) {
            violations = append(violations, Violation{
                RuleID:      fmt.Sprintf("SOC2_CONFIDENTIALITY_%s", strings.ReplaceAll(strings.ToUpper(pattern), " ", "_")),
                Name:        "SOC2 - Confidentiality Violation",
                Description: fmt.Sprintf("Potential confidentiality violation: %s", pattern),
                Severity:    "HIGH",
                Timestamp:   time.Now(),
            })
        }
    }

    return violations, nil
}

// Violation represents a compliance violation
type Violation struct {
    RuleID      string
    Name        string
    Description string
    Severity    string
    Timestamp   time.Time
}
