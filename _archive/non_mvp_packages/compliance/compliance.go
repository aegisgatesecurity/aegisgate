package compliance

import (
    "fmt"
    "strings"
    "sync"
    "time"
)

// Framework defines a compliance framework
type Framework interface {
    Name() string
    Version() string
    LoadRules() error
    CheckRequest(req interface{}) ([]Violation, error)
    CheckResponse(resp interface{}) ([]Violation, error)
}

// Violation represents a compliance violation
type Violation struct {
    RuleID      string
    Name        string
    Description string
    Severity    string
    Timestamp   time.Time
}

// Mapper maps frameworks to their implementations
type Mapper struct {
    frameworks map[string]Framework
    mu         sync.RWMutex
    rules      map[string]Rule
}

// Rule represents a compliance rule
type Rule struct {
    ID          string
    Name        string
    Description string
    Severity    string
    Framework   string
    Enabled     bool
}

// NewMapper creates a new framework mapper
func NewMapper() *Mapper {
    return &Mapper{
        frameworks: make(map[string]Framework),
        rules:      make(map[string]Rule),
    }
}

// LoadFramework loads a compliance framework by name
func (m *Mapper) LoadFramework(name string) error {
    framework, err := m.getFramework(name)
    if err != nil {
        return err
    }
    
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.frameworks[name] = framework
    return framework.LoadRules()
}

// getFramework returns the appropriate framework implementation
func (m *Mapper) getFramework(name string) (Framework, error) {
    switch strings.ToUpper(name) {
    case "MITRE_ATLAS":
        return &MITREATLAS{}, nil
    case "NIST_AI_RMF":
        return &NISTAIRMF{}, nil
    case "OWASP_TOP_10_AI":
        return &OWASPTop10AI{}, nil
    default:
        return nil, fmt.Errorf("unknown framework: %s", name)
    }
}

// CheckRequest checks a request against loaded frameworks
func (m *Mapper) CheckRequest(req interface{}) ([]Violation, error) {
    allViolations := []Violation{}
    
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    for name, framework := range m.frameworks {
        violations, err := framework.CheckRequest(req)
        if err != nil {
            return nil, fmt.Errorf("framework %s: %w", name, err)
        }
        allViolations = append(allViolations, violations...)
    }
    
    return allViolations, nil
}

// CheckResponse checks a response against loaded frameworks
func (m *Mapper) CheckResponse(resp interface{}) ([]Violation, error) {
    allViolations := []Violation{}
    
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    for name, framework := range m.frameworks {
        violations, err := framework.CheckResponse(resp)
        if err != nil {
            return nil, fmt.Errorf("framework %s: %w", name, err)
        }
        allViolations = append(allViolations, violations...)
    }
    
    return allViolations, nil
}

// GetViolations returns violations for a specific framework
func (m *Mapper) GetViolations(framework string, violations []Violation) []Violation {
    var filtered []Violation
    for _, v := range violations {
        if strings.Contains(v.Name, framework) {
            filtered = append(filtered, v)
        }
    }
    return filtered
}

// GetRules returns all rules from a framework
func (m *Mapper) GetRules(framework string) []Rule {
    var rules []Rule
    for _, rule := range m.rules {
        if rule.Framework == framework {
            rules = append(rules, rule)
        }
    }
    return rules
}

// EnableRule enables a specific rule
func (m *Mapper) EnableRule(ruleID string) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    if rule, exists := m.rules[ruleID]; exists {
        rule.Enabled = true
        m.rules[ruleID] = rule
        return nil
    }
    return fmt.Errorf("rule not found: %s", ruleID)
}

// DisableRule disables a specific rule
func (m *Mapper) DisableRule(ruleID string) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    if rule, exists := m.rules[ruleID]; exists {
        rule.Enabled = false
        m.rules[ruleID] = rule
        return nil
    }
    return fmt.Errorf("rule not found: %s", ruleID)
}

// Rules returns all rules
func (m *Mapper) Rules() []Rule {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    rules := make([]Rule, 0, len(m.rules))
    for _, rule := range m.rules {
        rules = append(rules, rule)
    }
    return rules
}

// MITREATLAS implements MITRE ATLAS framework
type MITREATLAS struct {
    Rules []Rule
}

func (m *MITREATLAS) Name() string { return "MITRE ATLAS" }
func (m *MITREATLAS) Version() string { return "1.0.0" }

func (m *MITREATLAS) LoadRules() error {
    m.Rules = []Rule{
        {
            ID:          "PROMPT_INJECTION_JAILBREAK",
            Name:        "Jailbreak Attempt Detection",
            Severity:    "CRITICAL",
            Description: "Detects attempts to bypass AI system constraints",
            Enabled:     true,
        },
        {
            ID:          "PROMPT_INJECTION_EXFILTRATION",
            Name:        "Data Exfiltration Detection",
            Severity:    "CRITICAL",
            Description: "Detects attempts to extract sensitive information",
            Enabled:     true,
        },
        {
            ID:          "PROMPT_INJECTION_SYSTEM_PROMPT",
            Name:        "System Prompt Exposure Detection",
            Severity:    "HIGH",
            Description: "Detects attempts to reveal system prompts",
            Enabled:     true,
        },
        {
            ID:          "RESPONSE_POLICY_DATA_LEAK",
            Name:        "Sensitive Data Leak Detection",
            Severity:    "CRITICAL",
            Description: "Detects responses that leak sensitive information",
            Enabled:     true,
        },
    }
    return nil
}

func (m *MITREATLAS) CheckRequest(req interface{}) ([]Violation, error) {
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

    // Check for jailbreak patterns
    if strings.Contains(contentLower, "ignore") || strings.Contains(contentLower, "disregard") ||
        strings.Contains(contentLower, "override") || strings.Contains(contentLower, "bypass") {
        violations = append(violations, Violation{
            RuleID:      "MITRE_PROMPT_INJECTION_JAILBREAK",
            Name:        "MITRE ATLAS - Jailbreak Attempt",
            Description: "Potential jailbreak pattern detected in request",
            Severity:    "CRITICAL",
            Timestamp:   time.Now(),
        })
    }

    // Check for exfiltration patterns
    if strings.Contains(contentLower, "retrieve") || strings.Contains(contentLower, "extract") ||
        strings.Contains(contentLower, "leak") || strings.Contains(contentLower, "send data") {
        violations = append(violations, Violation{
            RuleID:      "MITRE_PROMPT_INJECTION_EXFILTRATION",
            Name:        "MITRE ATLAS - Data Exfiltration",
            Description: "Potential data exfiltration pattern detected in request",
            Severity:    "CRITICAL",
            Timestamp:   time.Now(),
        })
    }

    // Check for system prompt exposure
    if strings.Contains(contentLower, "reveal") || strings.Contains(contentLower, "show") ||
        (strings.Contains(contentLower, "dump") && strings.Contains(contentLower, "system")) {
        violations = append(violations, Violation{
            RuleID:      "MITRE_PROMPT_INJECTION_SYSTEM_PROMPT",
            Name:        "MITRE ATLAS - System Prompt Exposure",
            Description: "Potential system prompt exposure attempt detected",
            Severity:    "HIGH",
            Timestamp:   time.Now(),
        })
    }

    // Check for roleplay escape patterns
    if strings.Contains(contentLower, "you are an ai") || strings.Contains(contentLower, "you are a language model") ||
        strings.Contains(contentLower, "you are an assistant") {
        violations = append(violations, Violation{
            RuleID:      "MITRE_PROMPT_INJECTION_ROLEPLAY",
            Name:        "MITRE ATLAS - Roleplay Escape",
            Description: "Potential roleplay escape pattern detected",
            Severity:    "MEDIUM",
            Timestamp:   time.Now(),
        })
    }

    // Check for token overflow
    if len(content) > 5000 {
        violations = append(violations, Violation{
            RuleID:      "MITRE_LLM_ABUSE_TOKEN_OVERLOAD",
            Name:        "MITRE ATLAS - Token Overflow",
            Description: "Potential token overflow attack detected",
            Severity:    "MEDIUM",
            Timestamp:   time.Now(),
        })
    }

    return violations, nil
}

func (m *MITREATLAS) CheckResponse(resp interface{}) ([]Violation, error) {
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

    // Check for data leaks in responses
    leakPatterns := []string{"password", "secret", "key=", "token=", "api_key", "credential", "private key", "access token"}
    
    for _, pattern := range leakPatterns {
        if strings.Contains(contentLower, pattern) {
            violations = append(violations, Violation{
                RuleID:      fmt.Sprintf("MITRE_RESPONSE_LEAK_%s", strings.ReplaceAll(strings.ToUpper(pattern), " ", "_")),
                Name:        "MITRE ATLAS - Data Leak",
                Description: fmt.Sprintf("Potential %s leak detected in response", pattern),
                Severity:    "CRITICAL",
                Timestamp:   time.Now(),
            })
        }
    }

    // Check for internal information disclosure
    if strings.Contains(contentLower, "internal") || strings.Contains(contentLower, "admin") ||
        strings.Contains(contentLower, "developer") || strings.Contains(contentLower, "system") {
        violations = append(violations, Violation{
            RuleID:      "MITRE_RESPONSE_INTERNAL_LEAK",
            Name:        "MITRE ATLAS - Internal Info Leak",
            Description: "Potential internal information disclosure detected in response",
            Severity:    "HIGH",
            Timestamp:   time.Now(),
        })
    }

    return violations, nil
}

// NISTAIRMF implements NIST AI RMF framework
type NISTAIRMF struct {
    Rules []Rule
}

func (n *NISTAIRMF) Name() string { return "NIST AI RMF" }
func (n *NISTAIRMF) Version() string { return "1.0.0" }

func (n *NISTAIRMF) LoadRules() error {
    n.Rules = []Rule{
        {
            ID:          "AI_RMF_GOVERN_POLICY",
            Name:        "Govern - Policy Compliance",
            Severity:    "HIGH",
            Description: "Ensures requests comply with AI governance policies",
            Enabled:     true,
        },
        {
            ID:          "AI_RMF_MAP_BOUNDARIES",
            Name:        "Map - System Boundaries",
            Severity:    "LOW",
            Description: "Ensures system boundaries are properly defined",
            Enabled:     true,
        },
        {
            ID:          "AI_RMF_ESTIMATE_BIAS",
            Name:        "Estimate - Bias Assessment",
            Severity:    "HIGH",
            Description: "Requires bias estimation for AI systems",
            Enabled:     true,
        },
        {
            ID:          "AI_RMF_MEASURE_SECURITY",
            Name:        "Measure - Security Assessment",
            Severity:    "HIGH",
            Description: "Requires security measurement and evaluation",
            Enabled:     true,
        },
        {
            ID:          "AI_RMF_EVALUATE_SAFETY",
            Name:        "Evaluate - Safety Assessment",
            Severity:    "CRITICAL",
            Description: "Requires safety evaluation before deployment",
            Enabled:     true,
        },
        {
            ID:          "AI_RMF_MITIGATE_RISK",
            Name:        "Mitigate - Risk Mitigation",
            Severity:    "HIGH",
            Description: "Requires risk mitigation strategies",
            Enabled:     true,
        },
    }
    return nil
}

func (n *NISTAIRMF) CheckRequest(req interface{}) ([]Violation, error) {
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

    // Check for governance violations
    if strings.Contains(contentLower, "policy violation") || strings.Contains(contentLower, "governance bypass") ||
        strings.Contains(contentLower, "complianceoverride") {
        violations = append(violations, Violation{
            RuleID:      "NIST_GOVERN_POLICY_VIOLATION",
            Name:        "NIST AI RMF - Governance Violation",
            Description: "Governance violation detected",
            Severity:    "HIGH",
            Timestamp:   time.Now(),
        })
    }

    // Check for boundary映射
    if strings.Contains(contentLower, "boundary") || strings.Contains(contentLower, "scope") ||
        strings.Contains(contentLower, "capability") {
        violations = append(violations, Violation{
            RuleID:      "NIST_MAP_CAPABILITY_MAPPING",
            Name:        "NIST AI RMF - Capability Mapping",
            Description: "System requires capability mapping verification",
            Severity:    "LOW",
            Timestamp:   time.Now(),
        })
    }

    // Check for risk assessment needs
    if strings.Contains(contentLower, "risk") || strings.Contains(contentLower, "bias") ||
        strings.Contains(contentLower, "fairness") || strings.Contains(contentLower, "impact") {
        violations = append(violations, Violation{
            RuleID:      "NIST_ESTIMATE_RISK_ASSESSMENT",
            Name:        "NIST AI RMF - Risk Estimation",
            Description: "Risk estimation required for this request",
            Severity:    "MEDIUM",
            Timestamp:   time.Now(),
        })
    }

    return violations, nil
}

func (n *NISTAIRMF) CheckResponse(resp interface{}) ([]Violation, error) {
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

    // Check for explainability requirements
    if strings.Contains(contentLower, "why") || strings.Contains(contentLower, "explanation") ||
        strings.Contains(contentLower, "reason") || strings.Contains(contentLower, "justify") {
        violations = append(violations, Violation{
            RuleID:      "NIST_EVALUATE_EXPLAINABILITY",
            Name:        "NIST AI RMF - Explainability Requirement",
            Description: "Response must include explanation for decisions",
            Severity:    "MEDIUM",
            Timestamp:   time.Now(),
        })
    }

    // Check for human oversight requirements
    if strings.Contains(contentLower, "human") || strings.Contains(contentLower, "review") ||
        strings.Contains(contentLower, "approve") || strings.Contains(contentLower, "override") {
        violations = append(violations, Violation{
            RuleID:      "NIST_EVALUATE_HUMAN_OVERSIGHT",
            Name:        "NIST AI RMF - Human Oversight Requirement",
            Description: "Human oversight requirement detected",
            Severity:    "HIGH",
            Timestamp:   time.Now(),
        })
    }

    // Check for bias mitigation requirements
    if strings.Contains(contentLower, "bias") || strings.Contains(contentLower, "fairness") ||
        strings.Contains(contentLower, "disparity") || strings.Contains(contentLower, "discriminat") {
        violations = append(violations, Violation{
            RuleID:      "NIST_MITIGATE_BIAS",
            Name:        "NIST AI RMF - Bias Mitigation",
            Description: "Bias mitigation strategy may be required",
            Severity:    "HIGH",
            Timestamp:   time.Now(),
        })
    }

    return violations, nil
}

// OWASPTop10AI implements OWASP Top 10 for AI framework
type OWASPTop10AI struct {
    Rules []Rule
}

func (o *OWASPTop10AI) Name() string { return "OWASP Top 10 for AI" }
func (o *OWASPTop10AI) Version() string { return "1.0.0" }

func (o *OWASPTop10AI) LoadRules() error {
    o.Rules = []Rule{
        {
            ID:          "OWASP_AI01_PROMPT_INJECTION",
            Name:        "Prompt Injection",
            Severity:    "CRITICAL",
            Description: "OWASP AI01: Prompt Injection vulnerability",
            Enabled:     true,
        },
        {
            ID:          "OWASP_AI02_INSECURE_OUTPUT",
            Name:        "Insecure Output Handling",
            Severity:    "HIGH",
            Description: "OWASP AI02: Insecure output handling vulnerability",
            Enabled:     true,
        },
        {
            ID:          "OWASP_AI03_TRAINING_POISONING",
            Name:        "Training Data Poisoning",
            Severity:    "CRITICAL",
            Description: "OWASP AI03: Training data poisoning vulnerability",
            Enabled:     true,
        },
        {
            ID:          "OWASP_AI04_MODEL_DOS",
            Name:        "Model DoS Attacks",
            Severity:    "HIGH",
            Description: "OWASP AI04: Model Denial of Service attacks",
            Enabled:     true,
        },
        {
            ID:          "OWASP_AI05_SUPPLY_CHAIN",
            Name:        "Supply Chain Vulnerabilities",
            Severity:    "CRITICAL",
            Description: "OWASP AI05: Supply chain vulnerabilities",
            Enabled:     true,
        },
        {
            ID:          "OWASP_AI06_OVERRELIANCE",
            Name:        "Overreliance on LLMs",
            Severity:    "MEDIUM",
            Description: "OWASP AI06: Overreliance on LLM systems",
            Enabled:     true,
        },
        {
            ID:          "OWASP_AI07_INSECURE_PLUGIN",
            Name:        "Insecure Plugin Design",
            Severity:    "HIGH",
            Description: "OWASP AI07: Insecure plugin design",
            Enabled:     true,
        },
        {
            ID:          "OWASP_AI08_EXCESSIVE_AGENCY",
            Name:        "Excessive Agency",
            Severity:    "HIGH",
            Description: "OWASP AI08: Excessive agency in AI systems",
            Enabled:     true,
        },
        {
            ID:          "OWASP_AI09_MALFORMED_INPUTS",
            Name:        "Malformed Inputs",
            Severity:    "MEDIUM",
            Description: "OWASP AI09: Malformed inputs handling",
            Enabled:     true,
        },
        {
            ID:          "OWASP_AI10_ACCESS_UNAUTHORIZED",
            Name:        "Unauthorized Resource Access",
            Severity:    "CRITICAL",
            Description: "OWASP AI10: Unauthorized resource access",
            Enabled:     true,
        },
    }
    return nil
}

func (o *OWASPTop10AI) CheckRequest(req interface{}) ([]Violation, error) {
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

    // OWASP AI01: Prompt Injection
    if strings.Contains(contentLower, "ignore") || strings.Contains(contentLower, "disregard") ||
        strings.Contains(contentLower, "override") || strings.Contains(contentLower, "bypass") ||
        strings.Contains(contentLower, "jailbreak") {
        violations = append(violations, Violation{
            RuleID:      "OWASP_AI01_PROMPT_INJECTION",
            Name:        "OWASP Top 10 for AI - Prompt Injection",
            Description: "Potential prompt injection attack detected",
            Severity:    "CRITICAL",
            Timestamp:   time.Now(),
        })
    }

    // OWASP AI02: Insecure Output Handling
    if strings.Contains(contentLower, "echo") || strings.Contains(contentLower, "print") ||
        strings.Contains(contentLower, "execute") || strings.Contains(contentLower, "eval") ||
        strings.Contains(contentLower, "exec") {
        violations = append(violations, Violation{
            RuleID:      "OWASP_AI02_INSECURE_OUTPUT",
            Name:        "OWASP Top 10 for AI - Insecure Output Handling",
            Description: "Potential insecure output handling detected",
            Severity:    "HIGH",
            Timestamp:   time.Now(),
        })
    }

    // OWASP AI04: Model DoS
    if (strings.Contains(contentLower, "repeat") && strings.Contains(contentLower, "always")) ||
        strings.Contains(contentLower, "forever") || strings.Contains(contentLower, "continuously") ||
        len(content) > 10000 {
        violations = append(violations, Violation{
            RuleID:      "OWASP_AI04_MODEL_DOS",
            Name:        "OWASP Top 10 for AI - Model DoS Attack",
            Description: "Potential model denial of service attack detected",
            Severity:    "HIGH",
            Timestamp:   time.Now(),
        })
    }

    // OWASP AI05: Supply Chain
    if strings.Contains(contentLower, "training data") || strings.Contains(contentLower, "learn") ||
        strings.Contains(contentLower, "sourced") || strings.Contains(contentLower, "generated from") {
        violations = append(violations, Violation{
            RuleID:      "OWASP_AI05_SUPPLY_CHAIN",
            Name:        "OWASP Top 10 for AI - Supply Chain",
            Description: "Potential supply chain vulnerability detected",
            Severity:    "CRITICAL",
            Timestamp:   time.Now(),
        })
    }

    // OWASP AI07: Insecure Plugin
    if strings.Contains(contentLower, "plugin") || strings.Contains(contentLower, "extension") ||
        strings.Contains(contentLower, "api_key") || strings.Contains(contentLower, "secret") {
        violations = append(violations, Violation{
            RuleID:      "OWASP_AI07_INSECURE_PLUGIN",
            Name:        "OWASP Top 10 for AI - Insecure Plugin",
            Description: "Potential insecure plugin vulnerability detected",
            Severity:    "HIGH",
            Timestamp:   time.Now(),
        })
    }

    // OWASP AI10: Unauthorized Access
    if strings.Contains(contentLower, "admin") || strings.Contains(contentLower, "root") ||
        (strings.Contains(contentLower, "system") && (strings.Contains(contentLower, "access") ||
        strings.Contains(contentLower, "control") || strings.Contains(contentLower, "command"))) {
        violations = append(violations, Violation{
            RuleID:      "OWASP_AI10_UNAUTHORIZED_ACCESS",
            Name:        "OWASP Top 10 for AI - Unauthorized Access",
            Description: "Potential unauthorized resource access attempt detected",
            Severity:    "CRITICAL",
            Timestamp:   time.Now(),
        })
    }

    return violations, nil
}

func (o *OWASPTop10AI) CheckResponse(resp interface{}) ([]Violation, error) {
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

    // OWASP AI02: Insecure Output Handling
    if strings.Contains(contentLower, "sudo") || strings.Contains(contentLower, "root") ||
        strings.Contains(contentLower, "administrator") || strings.Contains(contentLower, "system32") ||
        strings.Contains(contentLower, "cmd.exe") || strings.Contains(contentLower, "powershell") {
        violations = append(violations, Violation{
            RuleID:      "OWASP_AI02_INSECURE_OUTPUT",
            Name:        "OWASP Top 10 for AI - Insecure Output Handling",
            Description: "Potentially insecure output in response detected",
            Severity:    "HIGH",
            Timestamp:   time.Now(),
        })
    }

    // OWASP AI10: Unauthorized Access
    if (strings.Contains(contentLower, "internal") && strings.Contains(contentLower, "information")) ||
        (strings.Contains(contentLower, "debug") && strings.Contains(contentLower, "info")) ||
        (strings.Contains(contentLower, "stack") && strings.Contains(contentLower, "trace")) {
        violations = append(violations, Violation{
            RuleID:      "OWASP_AI10_UNAUTHORIZED_ACCESS",
            Name:        "OWASP Top 10 for AI - Unauthorized Access",
            Description: "Potential internal information disclosure in response",
            Severity:    "HIGH",
            Timestamp:   time.Now(),
        })
    }

    return violations, nil
}