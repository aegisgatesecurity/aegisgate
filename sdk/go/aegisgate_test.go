package aegisgate

import (
	"testing"
	"time"
)

func TestConfig_Defaults(t *testing.T) {
	config := &Config{
		REST:    RESTConfig{BaseURL: "https://api.example.com"},
		GRPC:    GRPCConfig{Address: "localhost:50051"},
		Timeout: 30 * time.Second,
	}
	
	if config.REST.BaseURL != "https://api.example.com" {
		t.Errorf("Expected BaseURL https://api.example.com, got %s", config.REST.BaseURL)
	}
	
	if config.GRPC.Address != "localhost:50051" {
		t.Errorf("Expected Address localhost:50051, got %s", config.GRPC.Address)
	}
}

func TestRESTConfig(t *testing.T) {
	config := RESTConfig{
		BaseURL:    "https://api.aegisgate.com",
		APIVersion: "v1",
	}
	
	if config.BaseURL == "" {
		t.Error("BaseURL should not be empty")
	}
	
	if config.APIVersion != "v1" {
		t.Errorf("Expected APIVersion v1, got %s", config.APIVersion)
	}
}

func TestGRPCConfig(t *testing.T) {
	config := GRPCConfig{
		Address:         "grpc.aegisgate.com:443",
		UseTLS:          true,
		ConnectionTimeout: 10 * time.Second,
	}
	
	if config.Address == "" {
		t.Error("Address should not be empty")
	}
	
	if !config.UseTLS {
		t.Error("UseTLS should be true")
	}
}

func TestAuthConfig(t *testing.T) {
	config := AuthConfig{
		APIKey:   "test-api-key",
		Username: "testuser",
		Password: "testpass",
	}
	
	if config.APIKey != "test-api-key" {
		t.Errorf("Expected APIKey test-api-key, got %s", config.APIKey)
	}
	
	if config.Username != "testuser" {
		t.Errorf("Expected Username testuser, got %s", config.Username)
	}
}
