// Package bridge provides bridging interfaces for AegisGate and AegisGuard integration.
// This package defines the core interfaces for secure communication between
// the API gateway (AegisGate) and the agent security platform (AegisGuard).
package bridge

import "context"

// Gateway defines the interface for processing AI API requests and responses
// through the security bridge between AegisGate and AegisGuard.
type Gateway interface {
	// ValidateRequest validates an incoming request for security threats.
	ValidateRequest(request string) error

	// ProcessResponse processes and filters an outgoing response.
	ProcessResponse(response string) (string, error)
}

// RBACManager defines the interface for role-based access control
// integration with AegisGuard.
type RBACManager interface {
	// CheckPermission verifies if a user has permission to perform an action on a resource.
	CheckPermission(ctx context.Context, userID, resource, action string) (bool, error)
}

// DefaultGateway provides a basic implementation of the Gateway interface.
type DefaultGateway struct{}

// ValidateRequest implements Gateway.ValidateRequest.
func (g *DefaultGateway) ValidateRequest(request string) error {
	// Basic validation - no-op for stub
	return nil
}

// ProcessResponse implements Gateway.ProcessResponse.
func (g *DefaultGateway) ProcessResponse(response string) (string, error) {
	// Basic processing - no-op for stub
	return response, nil
}

// DefaultRBACManager provides a basic implementation of the RBACManager interface.
type DefaultRBACManager struct{}

// CheckPermission implements RBACManager.CheckPermission.
func (m *DefaultRBACManager) CheckPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	// Basic permission check - allow by default for stub
	return true, nil
}