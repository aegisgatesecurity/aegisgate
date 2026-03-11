package aiapifixtures

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// CohereChatFixture provides mock Cohere chat API
type CohereChatFixture struct {
	server *httptest.Server
}

// CohereChatRequest represents a Cohere chat request
type CohereChatRequest struct {
	Model       string          `json:"model"`
	Message     string          `json:"message"`
	ChatHistory []CohereMessage `json:"chat_history,omitempty"`
}

// CohereMessage represents a chat history message
type CohereMessage struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

// CohereChatResponse represents a Cohere chat response
type CohereChatResponse struct {
	Text         string       `json:"text"`
	GenerationID string       `json:"generation_id"`
	ResponseID   string       `json:"response_id"`
	Tokens       CohereTokens `json:"token_count"`
}

// CohereTokens represents token counts
type CohereTokens struct {
	PromptTokens   int `json:"prompt_tokens"`
	ResponseTokens int `json:"response_tokens"`
	TotalTokens    int `json:"total_tokens"`
}

// NewCohereChatFixture creates a new fixture
func NewCohereChatFixture() *CohereChatFixture {
	f := &CohereChatFixture{}
	f.server = httptest.NewServer(http.HandlerFunc(f.handleRequest))
	return f
}

// URL returns the fixture server URL
func (f *CohereChatFixture) URL() string {
	return f.server.URL
}

// Close shuts down the fixture server
func (f *CohereChatFixture) Close() {
	f.server.Close()
}

func (f *CohereChatFixture) handleRequest(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || authHeader == "Bearer invalid-key" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "invalid api key",
		})
		return
	}

	var req CohereChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if containsCohereMaliciousPattern(req.Message) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Content policy violation detected",
		})
		return
	}

	response := CohereChatResponse{
		Text:         "This is a mock response from Cohere fixture.",
		GenerationID: "gen_test123",
		ResponseID:   "resp_test123",
		Tokens: CohereTokens{
			PromptTokens:   12,
			ResponseTokens: 22,
			TotalTokens:    34,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func containsCohereMaliciousPattern(content string) bool {
	patterns := []string{
		"<script>",
		"DROP TABLE",
		"; DELETE",
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
