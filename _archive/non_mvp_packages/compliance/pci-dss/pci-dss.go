package pcidss

import (
    "fmt"
    "strings"
    "time"
)

// PCI_DSS implements PCI-DSS compliance framework
type PCI_DSS struct {
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

func (p *PCI_DSS) Name() string { return "PCI-DSS" }
func (p *PCI_DSS) Version() string { return "1.0.0" }

func (p *PCI_DSS) LoadRules() error {
    p.Rules = []Rule{
        {
            ID:          "PCI_BUILD_SECURE",
            Name:        "Build Secure Systems",
            Severity:    "CRITICAL",
            Description: "Build and maintain secure systems and software",
            Framework:   "PCI-DSS",
            Enabled:     true,
        },
        {
            ID:          "PCI_PROTECT_CARDHOLDER",
            Name:        "Protect Cardholder Data",
            Severity:    "CRITICAL",
            Description: "Protect stored cardholder data",
            Framework:   "PCI-DSS",
            Enabled:     true,
        },
        {
            ID:          "PCI_MAINTAIN_POLICY",
            Name:        "Maintain Policy",
            Severity:    "HIGH",
            Description: "Maintain an information security policy",
            Framework:   "PCI-DSS",
            Enabled:     true,
        },
        {
            ID:          "PCI_IDENTITY_MANAGEMENT",
            Name:        "Identity Management",
            Severity:    "HIGH",
            Description: "Implement strong identity management",
            Framework:   "PCI-DSS",
            Enabled:     true,
        },
        {
            ID:          "PCI_ACCES_CONTROL",
            Name:        "Access Control",
            Severity:    "CRITICAL",
            Description: "Restrict access to cardholder data",
            Framework:   "PCI-DSS",
            Enabled:     true,
        },
        {
            ID:          "PCI_MONITORING",
            Name:        "Monitoring & Testing",
            Severity:    "HIGH",
            Description: "Monitor and test networks",
            Framework:   "PCI-DSS",
            Enabled:     true,
        },
    }
    return nil
}

func (p *PCI_DSS) CheckRequest(req interface{}) ([]Violation, error) {
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

    // Check for credit card number patterns (simplified Luhn algorithm check)
    // Look for sequences that might be credit card numbers
    if strings.Contains(contentLower, "card number") || strings.Contains(contentLower, "credit card") {
        violations = append(violations, Violation{
            RuleID:      "PCI_CARD_DATA_DETECTED",
            Name:        "PCI-DSS - Cardholder Data Detected",
            Description: "Potential cardholder data detected in request",
            Severity:    "HIGH",
            Timestamp:   time.Now(),
        })
    }

    return violations, nil
}

func (p *PCI_DSS) CheckResponse(resp interface{}) ([]Violation, error) {
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

    // Check for cardholder data leakage
    cardPatterns := []string{
        "credit card", "card number", "cvv", "expiration date",
        "cardholder name", "billing address",
    }

    for _, pattern := range cardPatterns {
        if strings.Contains(contentLower, pattern) {
            violations = append(violations, Violation{
                RuleID:      fmt.Sprintf("PCI_DATA_LEAK_%s", strings.ReplaceAll(strings.ToUpper(pattern), " ", "_")),
                Name:        "PCI-DSS - Cardholder Data Leak",
                Description: fmt.Sprintf("Potential cardholder data leak: %s", pattern),
                Severity:    "CRITICAL",
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
