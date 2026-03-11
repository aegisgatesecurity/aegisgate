package inspector

import (
	"fmt"
	"strings"
)

// Request represents an HTTP request
type Request struct {
	Method      string
	URL         string
	Headers     map[string]string
	Body        []byte
	Timestamp   int64
}

// Response represents an HTTP response
type Response struct {
	StatusCode  int
	Headers     map[string]string
	Body        []byte
	Timestamp   int64
}

// Policy defines an inspection policy
type Policy struct {
	ID          string
	Name        string
	Rules       []string
	Severity    string
	Enabled     bool
}

// Violation represents a policy violation
type Violation struct {
	RuleID      string
	PolicyID    string
	Severity    string
	Description string
	Timestamp   int64
}

// Inspector handles request/response inspection
type Inspector struct {
	policies        map[string]*Policy
	enabledPolicies map[string]bool
	rules           map[string]func(*Request, *Response) ([]Violation, error)
	mapper          InspectorMapper
}

// InspectorMapper interface for compliance mapping
type InspectorMapper interface {
	LoadFramework(fw string) error
	GetRules(fw string) []string
	MapViolation(violation Violation, framework string) bool
}

// NewInspector creates a new inspector
func NewInspector(mapper InspectorMapper) *Inspector {
	inspector := &Inspector{
		policies:        make(map[string]*Policy),
		enabledPolicies: make(map[string]bool),
		rules:           make(map[string]func(*Request, *Response) ([]Violation, error)),
		mapper:          mapper,
	}
	
	// Register example rules
	inspector.AddRule("sql_injection", inspector.sqlInjectionRule)
	inspector.AddRule("xss_detection", inspector.xssDetectionRule)
	
	return inspector
}

// AddPolicy adds a policy to the inspector
func (i *Inspector) AddPolicy(policy *Policy) {
	i.policies[policy.ID] = policy
}

// EnablePolicy enables a policy
func (i *Inspector) EnablePolicy(policyID string) error {
	if _, exists := i.policies[policyID]; !exists {
		return fmt.Errorf("policy not found: %s", policyID)
	}
	i.enabledPolicies[policyID] = true
	return nil
}

// DisablePolicy disables a policy
func (i *Inspector) DisablePolicy(policyID string) error {
	if _, exists := i.enabledPolicies[policyID]; !exists {
		return fmt.Errorf("policy not enabled: %s", policyID)
	}
	delete(i.enabledPolicies, policyID)
	return nil
}

// InspectRequest inspects an incoming request
func (i *Inspector) InspectRequest(req *Request) ([]Violation, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	
	var violations []Violation
	for policyID := range i.enabledPolicies {
		policy := i.policies[policyID]
		if policy != nil && policy.Enabled {
			for _, rule := range policy.Rules {
				if ruleFunc, exists := i.rules[rule]; exists {
					res, err := ruleFunc(req, nil)
					if err != nil {
						return nil, err
					}
					violations = append(violations, res...)
				}
			}
		}
	}
	return violations, nil
}

// InspectResponse inspects an outgoing response
func (i *Inspector) InspectResponse(resp *Response) ([]Violation, error) {
	if resp == nil {
		return nil, fmt.Errorf("response cannot be nil")
	}
	
	var violations []Violation
	for policyID := range i.enabledPolicies {
		policy := i.policies[policyID]
		if policy != nil && policy.Enabled {
			for _, rule := range policy.Rules {
				if ruleFunc, exists := i.rules[rule]; exists {
					res, err := ruleFunc(nil, resp)
					if err != nil {
						return nil, err
					}
					violations = append(violations, res...)
				}
			}
		}
	}
	return violations, nil
}

// AddRule adds an inspection rule
func (i *Inspector) AddRule(ruleID string, fn func(*Request, *Response) ([]Violation, error)) {
	i.rules[ruleID] = fn
}

// GetEnabledPolicies returns list of enabled policy IDs
func (i *Inspector) GetEnabledPolicies() []string {
	var policies []string
	for policyID := range i.enabledPolicies {
		policies = append(policies, policyID)
	}
	return policies
}

// GetPolicies returns all policies
func (i *Inspector) GetPolicies() map[string]*Policy {
	return i.policies
}

// sqlInjectionRule detects SQL injection attempts
func (i *Inspector) sqlInjectionRule(req *Request, resp *Response) ([]Violation, error) {
	if req == nil {
		return nil, nil
	}
	
	body := string(req.Body)
	suspiciousPatterns := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "UNION"}
	
	var violations []Violation
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(strings.ToUpper(body), pattern) {
			violations = append(violations, Violation{
				RuleID:      "sql_injection",
				Severity:    "high",
				Description: "Potential SQL injection detected",
				Timestamp:   req.Timestamp,
			})
			break
		}
	}
	return violations, nil
}

// xssDetectionRule detects XSS attempts
func (i *Inspector) xssDetectionRule(req *Request, resp *Response) ([]Violation, error) {
	if req == nil {
		return nil, nil
	}
	
	body := string(req.Body)
	suspiciousPatterns := []string{"<script>", "javascript:", "onerror=", "onload="}
	
	var violations []Violation
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(strings.ToLower(body), pattern) {
			violations = append(violations, Violation{
				RuleID:      "xss_detection",
				Severity:    "high",
				Description: "Potential XSS attack detected",
				Timestamp:   req.Timestamp,
			})
			break
		}
	}
	return violations, nil
}
