package inspector

import (
	"testing"
	"time"
)

func TestNewInspector(t *testing.T) {
	inspector := NewInspector(nil)
	if inspector == nil {
		t.Error("NewInspector should return non-nil inspector")
	}
	if inspector.policies == nil {
		t.Error("policies map should not be nil")
	}
	if inspector.rules == nil {
		t.Error("rules map should not be nil")
	}
}

func TestAddPolicy(t *testing.T) {
	inspector := NewInspector(nil)
	policy := &Policy{ID: "test-policy", Name: "Test Policy", Severity: "medium", Enabled: true}
	inspector.AddPolicy(policy)
	
	if _, exists := inspector.policies["test-policy"]; !exists {
		t.Error("Policy should be added to map")
	}
}

func TestEnablePolicy(t *testing.T) {
	inspector := NewInspector(nil)
	policy := &Policy{ID: "test-policy", Name: "Test Policy", Severity: "medium", Enabled: true}
	inspector.AddPolicy(policy)
	
	err := inspector.EnablePolicy("test-policy")
	if err != nil {
		t.Fatalf("Failed to enable policy: %v", err)
	}
	
	if !inspector.enabledPolicies["test-policy"] {
		t.Error("Policy should be enabled")
	}
}

func TestEnablePolicyNotFound(t *testing.T) {
	inspector := NewInspector(nil)
	err := inspector.EnablePolicy("nonexistent")
	if err == nil {
		t.Error("Should return error for nonexistent policy")
	}
}

func TestDisablePolicy(t *testing.T) {
	inspector := NewInspector(nil)
	policy := &Policy{ID: "test-policy", Name: "Test Policy", Severity: "medium", Enabled: true}
	inspector.AddPolicy(policy)
	inspector.EnablePolicy("test-policy")
	
	err := inspector.DisablePolicy("test-policy")
	if err != nil {
		t.Fatalf("Failed to disable policy: %v", err)
	}
	
	if inspector.enabledPolicies["test-policy"] {
		t.Error("Policy should be disabled")
	}
}

func TestDisablePolicyNotEnabled(t *testing.T) {
	inspector := NewInspector(nil)
	err := inspector.DisablePolicy("nonexistent")
	if err == nil {
		t.Error("Should return error for non-enabled policy")
	}
}

func TestInspectRequest(t *testing.T) {
	inspector := NewInspector(nil)
	policy := &Policy{ID: "test-policy", Name: "Test Policy", Severity: "medium", Enabled: true}
	policy.Rules = []string{"sql_injection"}
	inspector.AddPolicy(policy)
	inspector.EnablePolicy("test-policy")
	
	req := &Request{Method: "POST", URL: "/api", Body: []byte("SELECT * FROM users")}
	
	violations, err := inspector.InspectRequest(req)
	if err != nil {
		t.Fatalf("InspectRequest failed: %v", err)
	}
	
	// Should have at least one violation
	if len(violations) == 0 {
		t.Error("Expected violations")
	}
}

func TestInspectResponse(t *testing.T) {
	inspector := NewInspector(nil)
	policy := &Policy{ID: "test-policy", Name: "Test Policy", Severity: "medium", Enabled: true}
	policy.Rules = []string{"xss_detection"}
	inspector.AddPolicy(policy)
	inspector.EnablePolicy("test-policy")
	
	// Create request with XSS content in body
	req := &Request{
		Method:      "POST",
		URL:         "/api/submit",
		Headers:     map[string]string{"Content-Type": "text/html"},
		Body:        []byte("<script>alert('xss')</script>"),
		Timestamp:   time.Now().Unix(),
	}
	
	violations, err := inspector.InspectRequest(req)
	if err != nil {
		t.Fatalf("InspectRequest failed: %v", err)
	}
	
	if len(violations) == 0 {
		t.Error("Expected violations")
	}
}

func TestAddRule(t *testing.T) {
	inspector := NewInspector(nil)
	
	inspector.AddRule("test_rule", func(req *Request, resp *Response) ([]Violation, error) {
		return []Violation{{RuleID: "test_rule", Severity: "low"}}, nil
	})
	
	if _, exists := inspector.rules["test_rule"]; !exists {
		t.Error("Rule should be added")
	}
}

func TestGetEnabledPolicies(t *testing.T) {
	inspector := NewInspector(nil)
	policy1 := &Policy{ID: "policy1", Name: "Policy 1", Severity: "medium", Enabled: true}
	policy2 := &Policy{ID: "policy2", Name: "Policy 2", Severity: "high", Enabled: true}
	inspector.AddPolicy(policy1)
	inspector.AddPolicy(policy2)
	inspector.EnablePolicy("policy1")
	inspector.EnablePolicy("policy2")
	
	policies := inspector.GetEnabledPolicies()
	
	if len(policies) != 2 {
		t.Errorf("Expected 2 enabled policies, got %d", len(policies))
	}
}
