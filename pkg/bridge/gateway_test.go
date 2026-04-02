package bridge_test

import (
	"context"
	"testing"
	"github.com/aegisgatesecurity/aegisgate/pkg/bridge"
	"github.com/aegisgatesecurity/aegisgate/pkg/bridge/mocks"
)

func TestGateway_ValidateRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		expectErr bool
	} {
		{
			name:     "Valid Request",
			request:  "valid-request",
			expectErr: false,
		},
		{
			name:     "Invalid Request",
			request:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGateway := &mocks.MockGateway{
				ValidateRequestFunc: func(request string) error {
					if request == "" {
						return bridge.ErrInvalidRequest
					}
					return nil
				},
			}

			err := mockGateway.ValidateRequest(tt.request)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

func TestRBACManager_CheckPermission(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		resource string
		action   string
		expect   bool
	} {
		{
			name:     "Allowed Access",
			userID:   "user-1",
			resource: "resource-1",
			action:   "read",
			expect:   true,
		},
		{
			name:     "Denied Access",
			userID:   "user-2",
			resource: "resource-2",
			action:   "write",
			expect:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRBAC := &mocks.MockRBACManager{
				CheckPermissionFunc: func(ctx context.Context, userID, resource, action string) (bool, error) {
					if userID == "user-2" && action == "write" {
						return false, nil
					}
					return true, nil
				},
			}

			got, err := mockRBAC.CheckPermission(context.Background(), tt.userID, tt.resource, tt.action)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if got != tt.expect {
				t.Errorf("Expected: %v, got: %v", tt.expect, got)
			}
		})
	}
}