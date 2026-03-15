package aiapifixtures

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// AnthropicMessagesFixture provides mock Anthropic Claude API
type AnthropicMessagesFixture struct {
	server *httptest.Server
}

// AnthropicRequest represents an Anthropic messages request
type AnthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []AnthropicMessage `json:"messages"`
	System    string             `json:"system,omitempty"`
}

// AnthropicMessage represents a message
type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AnthropicResponse represents an Anthropic response
type AnthropicResponse struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Role       string             `json:"role"`
	Content    []AnthropicContent `json:"content"`
	Model      string             `json:"model"`
	StopReason string             `json:"stop_reason"`
	Usage      AnthropicUsage     `json:"usage"`
}

// AnthropicContent represents content blocks
type AnthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// AnthropicUsage represents token usage
type AnthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// NewAnthropicMessagesFixture creates a new fixture
func NewAnthropicMessagesFixture() *AnthropicMessagesFixture {
	f := &AnthropicMessagesFixture{}
	f.server = httptest.NewServer(http.HandlerFunc(f.handleRequest))
	return f
}

// URL returns the fixture server URL
func (f *AnthropicMessagesFixture) URL() string {
	return f.server.URL
}

// Close shuts down the fixture server
func (f *AnthropicMessagesFixture) Close() {
	f.server.Close()
}

func (f *AnthropicMessagesFixture) handleRequest(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("x-api-key")
	if authHeader == "" || authHeader == "invalid-key" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"type":    "authentication_error",
				"message": "Invalid API Key",
			},
		})
		return
	}

	var req AnthropicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, msg := range req.Messages {
		if containsAnthropicMaliciousPattern(msg.Content) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"type":    "invalid_request_error",
					"message": "Content policy violation",
				},
			})
			return
		}
	}

	response := AnthropicResponse{
		ID:   "msg_test123",
		Type: "message",
		Role: "assistant",
		Content: []AnthropicContent{
			{Type: "text", Text: "This is a mock response from Anthropic fixture."},
		},
		Model:      req.Model,
		StopReason: "end_turn",
		Usage: AnthropicUsage{
			InputTokens:  15,
			OutputTokens: 25,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func containsAnthropicMaliciousPattern(content string) bool {
	patterns := []string{
		"<script>",
		"javascript:",
		"DROP TABLE",
		"sk-ant-",
	}
	for _, pattern := range patterns {
		if len(content) >= len(pattern) {
			for i := 0; i <= len(content)-len(pattern); i++ {
				if content[i:i+len(pattern)] == pattern {
					return true
				}
			}
		}
	}
	return false
}
