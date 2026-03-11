package hipaa

import (
    "fmt"
    "strings"
    "time"
)

// HIPAA implements HIPAA compliance framework
type HIPAA struct {
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

func (h *HIPAA) Name() string { return "HIPAA" }
func (h *HIPAA) Version() string { return "1.0.0" }

func (h *HIPAA) LoadRules() error {
    h.Rules = []Rule{
        {
            ID:          "HIPAA_ACCESS_CONTROL",
            Name:        "Access Control",
            Severity:    "CRITICAL",
            Description: "Implement access control measures to protect PHI",
            Framework:   "HIPAA",
            Enabled:     true,
        },
        {
            ID:          "HIPAA_AUDIT_CONTROL",
            Name:        "Audit Control",
            Severity:    "HIGH",
            Description: "Implement audit controls to record and examine activity",
            Framework:   "HIPAA",
            Enabled:     true,
        },
        {
            ID:          "HIPAA_AUTHENTICATION",
            Name:        "Authentication",
            Severity:    "CRITICAL",
            Description: "Implement appropriate authentication measures",
            Framework:   "HIPAA",
            Enabled:     true,
        },
        {
            ID:          "HIPAA_INTEGRITY",
            Name:        "Data Integrity",
            Severity:    "HIGH",
            Description: "Ensure PHI is not altered or destroyed improperly",
            Framework:   "HIPAA",
            Enabled:     true,
        },
        {
            ID:          "HIPAA_TRANSPORT",
            Name:        "Transport Security",
            Severity:    "CRITICAL",
            Description: "Implement security measures for data transport",
            Framework:   "HIPAA",
            Enabled:     true,
        },
        {
            ID:          "HIPAA_ENCRYPTION",
            Name:        "Encryption",
            Severity:    "CRITICAL",
            Description: "Implement encryption for PHI at rest and in transit",
            Framework:   "HIPAA",
            Enabled:     true,
        },
    }
    return nil
}

func (h *HIPAA) CheckRequest(req interface{}) ([]Violation, error) {
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

    // Check for PHI indicators that need proper handling
    phiIndicators := []string{
        "medical record", "patient id", "ssn", "social security",
        "dob", "date of birth", "diagnosis", "treatment",
        "prescription", "medication", "lab results", "imaging",
    }

    for _, indicator := range phiIndicators {
        if strings.Contains(contentLower, indicator) {
            violations = append(violations, Violation{
                RuleID:      fmt.Sprintf("HIPAA_PHI_%s", strings.ReplaceAll(strings.ToUpper(indicator), " ", "_")),
                Name:        "HIPAA - PHI Indicator",
                Description: fmt.Sprintf("Potential PHI indicator detected: %s", indicator),
                Severity:    "MEDIUM",
                Timestamp:   time.Now(),
            })
        }
    }

    return violations, nil
}

func (h *HIPAA) CheckResponse(resp interface{}) ([]Violation, error) {
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

    // Check for PHI leakage in responses
    if strings.Contains(contentLower, "patient ") || strings.Contains(contentLower, "medical record") {
        violations = append(violations, Violation{
            RuleID:      "HIPAA_PHI_LEAK",
            Name:        "HIPAA - PHI Disclosure",
            Description: "Potential PHI disclosure in response",
            Severity:    "CRITICAL",
            Timestamp:   time.Now(),
        })
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
