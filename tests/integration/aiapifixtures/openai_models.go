package aiapifixtures

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// OpenAIModelsFixture provides mock OpenAI models API
type OpenAIModelsFixture struct {
	server *httptest.Server
}

// OpenAIModelsResponse represents the models list response
type OpenAIModelsResponse struct {
	Data []OpenAIModel `json:"data"`
}

// OpenAIModel represents a model object
type OpenAIModel struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// NewOpenAIModelsFixture creates a new fixture
func NewOpenAIModelsFixture() *OpenAIModelsFixture {
	f := &OpenAIModelsFixture{}
	f.server = httptest.NewServer(http.HandlerFunc(f.handleRequest))
	return f
}

// URL returns the fixture server URL
func (f *OpenAIModelsFixture) URL() string {
	return f.server.URL
}

// Close shuts down the fixture server
func (f *OpenAIModelsFixture) Close() {
	f.server.Close()
}

func (f *OpenAIModelsFixture) handleRequest(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response := OpenAIModelsResponse{
		Data: []OpenAIModel{
			{ID: "gpt-4", Object: "model", Created: 1687882411, OwnedBy: "openai"},
			{ID: "gpt-4-turbo", Object: "model", Created: 1706037612, OwnedBy: "openai"},
			{ID: "gpt-3.5-turbo", Object: "model", Created: 1677610602, OwnedBy: "openai"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
