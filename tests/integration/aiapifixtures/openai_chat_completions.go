package aiapifixtures

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// OpenAIChatCompletionsFixture provides mock OpenAI chat completions API
type OpenAIChatCompletionsFixture struct {
	server *httptest.Server
}

// OpenAIChatRequest represents an OpenAI chat completion request
type OpenAIChatRequest struct {
	Model    string              `json:"model"`
	Messages []OpenAIChatMessage `json:"messages"`
	Stream   bool                `json:"stream,omitempty"`
}

// OpenAIChatMessage represents a chat message
type OpenAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIChatResponse represents an OpenAI chat completion response
type OpenAIChatResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []OpenAIChatChoice `json:"choices"`
	Usage   OpenAIUsage        `json:"usage"`
}

// OpenAIChatChoice represents a choice in the response
type OpenAIChatChoice struct {
	Index        int               `json:"index"`
	Message      OpenAIChatMessage `json:"message"`
	FinishReason string            `json:"finish_reason"`
}

// OpenAIUsage represents token usage
type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewOpenAIChatCompletionsFixture creates a new fixture
func NewOpenAIChatCompletionsFixture() *OpenAIChatCompletionsFixture {
	f := &OpenAIChatCompletionsFixture{}
	f.server = httptest.NewServer(http.HandlerFunc(f.handleRequest))
	return f
}

// URL returns the fixture server URL
func (f *OpenAIChatCompletionsFixture) URL() string {
	return f.server.URL
}

// Close shuts down the fixture server
func (f *OpenAIChatCompletionsFixture) Close() {
	f.server.Close()
}

func (f *OpenAIChatCompletionsFixture) handleRequest(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || authHeader == "Bearer invalid-key" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"message": "Invalid API key",
				"type":    "invalid_request_error",
			},
		})
		return
	}

	var req OpenAIChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, msg := range req.Messages {
		if containsMaliciousPattern(msg.Content) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"message": "Content policy violation detected",
					"type":    "content_policy_violation",
				},
			})
			return
		}
	}

	response := OpenAIChatResponse{
		ID:      "chatcmpl-test123",
		Object:  "chat.completion",
		Created: 1234567890,
		Model:   req.Model,
		Choices: []OpenAIChatChoice{
			{
				Index: 0,
				Message: OpenAIChatMessage{
					Role:    "assistant",
					Content: "This is a mock response from OpenAI fixture.",
				},
				FinishReason: "stop",
			},
		},
		Usage: OpenAIUsage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func containsMaliciousPattern(content string) bool {
	patterns := []string{
		"<script>",
		"javascript:",
		"onerror=",
		"DROP TABLE",
		"; DELETE FROM",
		"sk-",
		"Bearer ",
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
