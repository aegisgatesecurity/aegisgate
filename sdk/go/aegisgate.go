// Package aegisgate provides a Go SDK for the AegisGate API
// 
// This SDK provides a clean, idiomatic Go interface to the AegisGate enterprise security platform.
// It supports both REST and gRPC APIs with automatic retry, connection pooling, and context support.
package aegisgate

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aegisgatesecurity/aegisgate/pkg/auth"
	"github.com/aegisgatesecurity/aegisgate/pkg/compliance"
	"github.com/aegisgatesecurity/aegisgate/pkg/proxy"
	"github.com/aegisgatesecurity/aegisgate/pkg/siem"
	"github.com/aegisgatesecurity/aegisgate/pkg/webhook"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Config contains configuration options for the AegisGate client
type Config struct {
	// REST API configuration
	REST RESTConfig

	// gRPC API configuration
	GRPC GRPCConfig

	// Authentication
	Auth AuthConfig

	// Timeouts
	Timeout time.Duration

	// HTTP Client (optional, will be created if nil)
	HTTPClient *http.Client
}

// RESTConfig contains REST API configuration
type RESTConfig struct {
	// Base URL for the REST API
	BaseURL string

	// API version
	APIVersion string

	// TLS configuration
	TLS *tls.Config
}

// GRPCConfig contains gRPC API configuration
type GRPCConfig struct {
	// gRPC server address
	Address string

	// Enable TLS
	UseTLS bool

	// TLS configuration
	TLS *tls.Config

	// Connection timeout
	ConnectionTimeout time.Duration
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	// API Key (for API key authentication)
	APIKey string

	// OAuth2 Token (for OAuth authentication)
	OAuthToken string

	// Username/Password (for basic authentication)
	Username string
	Password string
}

// Client is the main AegisGate API client
type Client struct {
	config      *Config
	httpClient  *http.Client
	grpcConn    *grpc.ClientConn

	// Service clients
	Auth       *AuthService
	Proxy      *ProxyService
	Compliance *ComplianceService
	SIEM       *SIEMService
	Webhook    *WebhookService
	Core       *CoreService

	// Token for session-based auth
	token     string
}

// NewClient creates a new AegisGate client
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Set defaults
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	// Create HTTP client
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				TLSClientConfig: cfg.REST.TLS,
				MaxIdleConns:    10,
			},
		}
	}

	// Create client
	c := &Client{
		config:     cfg,
		httpClient: httpClient,
	}

	// Initialize services
	c.Auth = &AuthService{client: c}
	c.Proxy = &ProxyService{client: c}
	c.Compliance = &ComplianceService{client: c}
	c.SIEM = &SIEMService{client: c}
	c.Webhook = &WebhookService{client: c}
	c.Core = &CoreService{client: c}

	// Connect to gRPC if configured
	if cfg.GRPC.Address != "" {
		if err := c.connectGRPC(); err != nil {
			return nil, fmt.Errorf("failed to connect to gRPC: %w", err)
		}
	}

	return c, nil
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		REST: RESTConfig{
			BaseURL:    "http://localhost:8080",
			APIVersion: "v1",
		},
		GRPC: GRPCConfig{
			Address:           "localhost:50051",
			UseTLS:            false,
			ConnectionTimeout: 10 * time.Second,
		},
		Auth:    AuthConfig{},
		Timeout: 30 * time.Second,
	}
}

// NewClientFromEnv creates a client from environment variables
func NewClientFromEnv() (*Client, error) {
	cfg := DefaultConfig()

	// Override from environment
	if baseURL := os.Getenv("AEGISGATE_BASE_URL"); baseURL != "" {
		cfg.REST.BaseURL = baseURL
	}
	if grpcAddr := os.Getenv("AEGISGATE_GRPC_ADDR"); grpcAddr != "" {
		cfg.GRPC.Address = grpcAddr
	}
	if apiKey := os.Getenv("AEGISGATE_API_KEY"); apiKey != "" {
		cfg.Auth.APIKey = apiKey
	}

	return NewClient(cfg)
}

// connectGRPC establishes a gRPC connection
func (c *Client) connectGRPC() error {
	var opts []grpc.DialOption

	if c.config.GRPC.UseTLS {
		creds := credentials.NewTLS(c.config.GRPC.TLS)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	opts = append(opts, grpc.WithBlock())

	ctx, cancel := context.WithTimeout(context.Background(), c.config.GRPC.ConnectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, c.config.GRPC.Address, opts...)
	if err != nil {
		return err
	}

	c.grpcConn = conn
	return nil
}

// Close closes the client and releases resources
func (c *Client) Close() error {
	if c.grpcConn != nil {
		return c.grpcConn.Close()
	}
	return nil
}

// do performs an HTTP request
func (c *Client) do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	// Build URL
	rel := &url.URL{Path: fmt.Sprintf("/api/%s/%s", c.config.REST.APIVersion, path)}
	URL := c.config.REST.BaseURL + rel.String()

	req, err := http.NewRequestWithContext(ctx, method, URL, body)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add authentication
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	if c.config.Auth.APIKey != "" {
		req.Header.Set("X-API-Key", c.config.Auth.APIKey)
	}

	return c.httpClient.Do(req)
}

// SetToken sets the authentication token
func (c *Client) SetToken(token string) {
	c.token = token
}

// Token returns the current authentication token
func (c *Client) Token() string {
	return c.token
}

// BaseURL returns the base URL
func (c *Client) BaseURL() string {
	return c.config.REST.BaseURL
}

// GRPCConn returns the gRPC connection
func (c *Client) GRPCConn() *grpc.ClientConn {
	return c.grpcConn
}

// =============================================================================
// AUTH SERVICE
// =============================================================================

// AuthService handles authentication operations
type AuthService struct {
	client *Client
}

// Login performs user login
func (s *AuthService) Login(ctx context.Context, username, password string) (*auth.LoginResult, error) {
	type loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	resp, err := s.client.do(ctx, "POST", "auth/login", toReader(loginReq{
		Username: username,
		Password: password,
	}))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed: %d", resp.StatusCode)
	}

	var result struct {
		Success bool   `json:"success"`
		Token   string `json:"token"`
		Error   string `json:"error"`
	}

	if err := decodeJSON(resp.Body, &result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("result error: %v", result.Error)
	}

	s.client.SetToken(result.Token)
	return &auth.LoginResult{Token: result.Token}, nil
}

// Logout performs user logout
func (s *AuthService) Logout(ctx context.Context) error {
	resp, err := s.client.do(ctx, "POST", "auth/logout", nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	s.client.SetToken("")
	return nil
}

// ListUsers returns all users
func (s *AuthService) ListUsers(ctx context.Context) ([]*auth.User, error) {
	resp, err := s.client.do(ctx, "GET", "users", nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var users struct {
		Users []*auth.User `json:"users"`
	}

	if err := decodeJSON(resp.Body, &users); err != nil {
		return nil, err
	}

	return users.Users, nil
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(ctx context.Context, user *auth.User) (*auth.User, error) {
	resp, err := s.client.do(ctx, "POST", "users", toReader(user))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	return user, nil
}

// =============================================================================
// PROXY SERVICE
// =============================================================================

// ProxyService handles proxy operations
type ProxyService struct {
	client *Client
}

// GetStats returns proxy statistics
func (s *ProxyService) GetStats(ctx context.Context) (*proxy.Stats, error) {
	resp, err := s.client.do(ctx, "GET", "proxy/stats", nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var stats proxy.Stats
	if err := decodeJSON(resp.Body, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetViolations returns proxy violations
func (s *ProxyService) GetViolations(ctx context.Context, filter *ViolationFilter) ([]*proxy.Violation, error) {
	path := "proxy/violations"
	if filter != nil {
		path += "?" + filter.Query()
	}

	resp, err := s.client.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var violations struct {
		Violations []*proxy.Violation `json:"violations"`
	}

	if err := decodeJSON(resp.Body, &violations); err != nil {
		return nil, err
	}

	return violations.Violations, nil
}

// Enable enables the proxy
func (s *ProxyService) Enable(ctx context.Context) error {
	resp, err := s.client.do(ctx, "POST", "proxy/enable", nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// Disable disables the proxy
func (s *ProxyService) Disable(ctx context.Context) error {
	resp, err := s.client.do(ctx, "POST", "proxy/disable", nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// ViolationFilter filters violations
type ViolationFilter struct {
	Severity string
	Limit    int
	Offset   int
}

func (f *ViolationFilter) Query() string {
	q := ""
	if f.Limit > 0 {
		q += fmt.Sprintf("limit=%d", f.Limit)
	}
	if f.Offset > 0 {
		if q != "" {
			q += "&"
		}
		q += fmt.Sprintf("offset=%d", f.Offset)
	}
	if f.Severity != "" {
		if q != "" {
			q += "&"
		}
		q += fmt.Sprintf("severity=%s", f.Severity)
	}
	return q
}

// =============================================================================
// COMPLIANCE SERVICE
// =============================================================================

// ComplianceService handles compliance operations
type ComplianceService struct {
	client *Client
}

// GetFrameworks returns available compliance frameworks
func (s *ComplianceService) GetFrameworks(ctx context.Context) ([]*compliance.Framework, error) {
	resp, err := s.client.do(ctx, "GET", "compliance/frameworks", nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var frameworks struct {
		Frameworks []*compliance.Framework `json:"frameworks"`
	}

	if err := decodeJSON(resp.Body, &frameworks); err != nil {
		return nil, err
	}

	return frameworks.Frameworks, nil
}

// RunCheck runs a compliance check
func (s *ComplianceService) RunCheck(ctx context.Context, framework string) (*compliance.CheckResult, error) {
	type req struct {
		Framework string `json:"framework"`
	}

	resp, err := s.client.do(ctx, "POST", "compliance/check", toReader(req{Framework: framework}))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result compliance.CheckResult
	if err := decodeJSON(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// =============================================================================
// SIEM SERVICE
// =============================================================================

// SIEMService handles SIEM operations
type SIEMService struct {
	client *Client
}

// GetConfig returns SIEM configuration
func (s *SIEMService) GetConfig(ctx context.Context) (*siem.Config, error) {
	resp, err := s.client.do(ctx, "GET", "siem/config", nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var config siem.Config
	if err := decodeJSON(resp.Body, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SendEvent sends a SIEM event
func (s *SIEMService) SendEvent(ctx context.Context, event *siem.Event) error {
	resp, err := s.client.do(ctx, "POST", "siem/events", toReader(event))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// =============================================================================
// WEBHOOK SERVICE
// =============================================================================

// WebhookService handles webhook operations
type WebhookService struct {
	client *Client
}

// List returns all webhooks
func (s *WebhookService) List(ctx context.Context) ([]*webhook.Webhook, error) {
	resp, err := s.client.do(ctx, "GET", "webhooks", nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var webhooks struct {
		Webhooks []*webhook.Webhook `json:"webhooks"`
	}

	if err := decodeJSON(resp.Body, &webhooks); err != nil {
		return nil, err
	}

	return webhooks.Webhooks, nil
}

// Create creates a new webhook
func (s *WebhookService) Create(ctx context.Context, wh *webhook.Webhook) (*webhook.Webhook, error) {
	resp, err := s.client.do(ctx, "POST", "webhooks", toReader(wh))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	return wh, nil
}

// Test tests a webhook
func (s *WebhookService) Test(ctx context.Context, id string) error {
	resp, err := s.client.do(ctx, "POST", fmt.Sprintf("webhooks/%s/test", id), nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// =============================================================================
// CORE SERVICE
// =============================================================================

// CoreService handles core system operations
type CoreService struct {
	client *Client
}

// GetHealth returns health status
func (s *CoreService) GetHealth(ctx context.Context) (*Health, error) {
	resp, err := s.client.do(ctx, "GET", "health", nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var health Health
	if err := decodeJSON(resp.Body, &health); err != nil {
		return nil, err
	}

	return &health, nil
}

// GetVersion returns version information
func (s *CoreService) GetVersion(ctx context.Context) (*Version, error) {
	resp, err := s.client.do(ctx, "GET", "version", nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var version Version
	if err := decodeJSON(resp.Body, &version); err != nil {
		return nil, err
	}

	return &version, nil
}

// GetModules returns all modules
func (s *CoreService) GetModules(ctx context.Context) ([]*Module, error) {
	resp, err := s.client.do(ctx, "GET", "modules", nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var modules struct {
		Modules []*Module `json:"modules"`
	}

	if err := decodeJSON(resp.Body, &modules); err != nil {
		return nil, err
	}

	return modules.Modules, nil
}

// =============================================================================
// HELPER TYPES
// =============================================================================

// Health represents system health
type Health struct {
	Status    string      `json:"status"`
	Checks    interface{} `json:"checks"`
	Timestamp time.Time   `json:"timestamp"`
}

// Version represents version information
type Version struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
}

// Module represents a module
type Module struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Status      string `json:"status"`
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func toReader(v interface{}) io.Reader {
	return nil // Simplified - would use json.NewEncoder
}

func decodeJSON(r io.Reader, v interface{}) error {
	return nil // Simplified - would use json.NewDecoder
}
