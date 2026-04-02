package auth_test

import (
	"context"
	"testing"
	"github.com/aegisgatesecurity/aegisgate/pkg/shared/auth"
)

func TestRBACManager_CheckPermission(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		resource string
		action   string
		expect   bool
	} {
		{
			name:     "Admin Access",
			userID:   "admin",
			resource: "dashboard",
			action:   "read",
			expect:   true,
		},
		{
			name:     "User Access Denied",
			userID:   "user",
			resource: "admin-panel",
			action:   "write",
			expect:   false,
		},
	}

	// Initialize RBAC Manager with test permissions
	rbac := auth.NewRBACManager()
	rbac.AddPermission("admin", "dashboard", "read")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rbac.CheckPermission(context.Background(), tt.userID, tt.resource, tt.action)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if got != tt.expect {
				t.Errorf("Expected: %v, got: %v", tt.expect, got)
			}
		})
	}
}

func TestRBACManager_AddPermission(t *testing.T) {
	rbac := auth.NewRBACManager()
	rbac.AddPermission("user", "profile", "read")

	got, err := rbac.CheckPermission(context.Background(), "user", "profile", "read")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !got {
		t.Errorf("Expected permission to be granted, but it was denied")
	}
}