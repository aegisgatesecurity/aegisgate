// Package mcp - Mock implementations for testing
package mcp

import (
	"context"
)

// MockToolAuthorizer implements ToolAuthorizer interface for testing
type MockToolAuthorizer struct {
	AuthorizedTools   map[string]bool
	AuthorizationRate map[string]int
	MockDecision      *AuthorizationDecision
	MockError         error
	AuthorizeCalls    []AuthorizationCall
}

func NewMockToolAuthorizer() *MockToolAuthorizer {
	return &MockToolAuthorizer{
		AuthorizedTools:   make(map[string]bool),
		AuthorizationRate: make(map[string]int),
	}
}

func (m *MockToolAuthorizer) Authorize(ctx context.Context, call *AuthorizationCall) (*AuthorizationDecision, error) {
	m.AuthorizeCalls = append(m.AuthorizeCalls, *call)

	if m.MockError != nil {
		return nil, m.MockError
	}

	if m.MockDecision != nil {
		return m.MockDecision, nil
	}

	decision := &AuthorizationDecision{
		Allowed:     true,
		Reason:      "allowed by default",
		RiskScore:   0,
		MatchedRule: "default_policy",
	}

	if allowed, ok := m.AuthorizedTools[call.Name]; ok {
		decision.Allowed = allowed
		if !allowed {
			decision.Reason = "tool denied by policy"
		}
	}

	if rate, ok := m.AuthorizationRate[call.Name]; ok {
		decision.RiskScore = rate
	}

	return decision, nil
}

// MockPolicyEngine implements PolicyEngine interface for testing
type MockPolicyEngine struct {
	PolicyResults     map[string]*PolicyEvalResult
	MockResult        *PolicyEvalResult
	MockError         error
	EvaluateCalls     []PolicyEvalContext
}

func NewMockPolicyEngine() *MockPolicyEngine {
	return &MockPolicyEngine{
		PolicyResults: make(map[string]*PolicyEvalResult),
	}
}

func (m *MockPolicyEngine) Evaluate(ctx context.Context, eval *PolicyEvalContext) (*PolicyEvalResult, error) {
	m.EvaluateCalls = append(m.EvaluateCalls, *eval)

	if m.MockError != nil {
		return nil, m.MockError
	}

	if m.MockResult != nil {
		return m.MockResult, nil
	}

	result := &PolicyEvalResult{
		Allowed:      true,
		Reason:       "policy evaluation passed",
		MatchedRules: []string{},
		ModifiedRisk: eval.Parameters["risk_score"].(int),
	}

	key := eval.ToolName + ":" + eval.SessionID + ":" + eval.AgentID
	if res, ok := m.PolicyResults[key]; ok {
		result = res
	}

	return result, nil
}

// MockAuditLogger implements AuditLogger interface for testing
type MockAuditLogger struct {
	AuditEntries []*AuditEntry
	LoggedErrors []error
	LogError     error
}

func NewMockAuditLogger() *MockAuditLogger {
	return &MockAuditLogger{
		AuditEntries: make([]*AuditEntry, 0),
	}
}

func (m *MockAuditLogger) Log(ctx context.Context, entry *AuditEntry) error {
	if m.LogError != nil {
		return m.LogError
	}
	m.AuditEntries = append(m.AuditEntries, entry)
	return nil
}

// MockSessionManager implements SessionManager interface for testing
type MockSessionManager struct {
	Sessions       map[string]*Session
	LastError      error
	CreatedCalls   []struct{ AgentID string }
	DeletedCalls   []string
}

func NewMockSessionManager() *MockSessionManager {
	return &MockSessionManager{
		Sessions: make(map[string]*Session),
	}
}

func (m *MockSessionManager) CreateSession(ctx context.Context, agentID string) (*Session, error) {
	m.CreatedCalls = append(m.CreatedCalls, struct{ AgentID string }{AgentID: agentID})

	if m.LastError != nil {
		return nil, m.LastError
	}

	session := &Session{
		ID:      "session-" + agentID,
		AgentID: agentID,
	}
	m.Sessions[session.ID] = session
	return session, nil
}

func (m *MockSessionManager) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	if session, ok := m.Sessions[sessionID]; ok {
		return session, nil
	}
	return nil, m.LastError
}

func (m *MockSessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	m.DeletedCalls = append(m.DeletedCalls, sessionID)
	delete(m.Sessions, sessionID)
	return nil
}

// TestConnection is a test implementation of Connection
type TestConnection struct {
	ID      string
	AgentID string
}

func (c *TestConnection) GetAgentID() string {
	return c.AgentID
}

func (c *TestConnection) GetSessionID() string {
	return c.ID
}
