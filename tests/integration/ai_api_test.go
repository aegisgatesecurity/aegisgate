package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aegisgatesecurity/aegisgate/tests/integration/aiapifixtures"
)

func TestOpenAIChatCompletionsIntegration(t *testing.T) {
	fixture := aiapifixtures.NewOpenAIChatCompletionsFixture()
	defer fixture.Close()

	t.Run("Successful Request", func(t *testing.T) {
		// Create request with auth header
		client := &http.Client{}
		reqBody := "{\"model\": \"gpt-4\", \"messages\": [{\"role\": \"user\", \"content\": \"Hello\"}]}"
		req, _ := http.NewRequest("POST", fixture.URL(), strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-key")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		var response aiapifixtures.OpenAIChatResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if response.Model != "gpt-4" {
			t.Errorf("Expected model gpt-4, got %s", response.Model)
		}
	})
}

func TestOpenAIModelsIntegration(t *testing.T) {
	fixture := aiapifixtures.NewOpenAIModelsFixture()
	defer fixture.Close()
	req, _ := http.NewRequest("GET", fixture.URL()+"/v1/models", nil)
	req.Header.Set("Authorization", "Bearer valid-key")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	var response aiapifixtures.OpenAIModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(response.Data) != 3 {
		t.Errorf("Expected 3 models, got %d", len(response.Data))
	}
}

func TestAnthropicMessagesIntegration(t *testing.T) {
	fixture := aiapifixtures.NewAnthropicMessagesFixture()
	defer fixture.Close()
	// Create request with auth header
	client := &http.Client{}
	reqBody := "{\"model\": \"claude-3-opus\", \"max_tokens\": 1024, \"messages\": [{\"role\": \"user\", \"content\": \"Hello\"}]}"
	req, _ := http.NewRequest("POST", fixture.URL(), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "valid-key")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	var response aiapifixtures.AnthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if response.Model != "claude-3-opus" {
		t.Errorf("Expected model claude-3-opus, got %s", response.Model)
	}
}

func TestCohereChatIntegration(t *testing.T) {
	fixture := aiapifixtures.NewCohereChatFixture()
	defer fixture.Close()
	// Create request with auth header
	client := &http.Client{}
	reqBody := "{\"model\": \"command-r\", \"message\": \"Hello\"}"
	req, _ := http.NewRequest("POST", fixture.URL(), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-key")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	var response aiapifixtures.CohereChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if response.Text == "" {
		t.Error("Expected non-empty response text")
	}
}

// TestOpenAIRateLimiting tests OpenAI rate limiting response handling
func TestOpenAIRateLimiting(t *testing.T) {
	// Create a server that returns 429 status
	rateLimitedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"message":     "Rate limit exceeded for 'tokens' limit. Please retry after 60 seconds.",
				"type":        "rate_limit_error",
				"param":       "n",
				"code":        "rate_limit_exceeded",
				"retry_after": 60,
			},
		})
	}))
	defer rateLimitedServer.Close()

	client := &http.Client{}
	reqBody := "{\"model\": \"gpt-4\", \"messages\": [{\"role\": \"user\", \"content\": \"Hello\"}]}"
	req, _ := http.NewRequest("POST", rateLimitedServer.URL, strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", resp.StatusCode)
	}

	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter != "60" {
		t.Errorf("Expected Retry-After header '60', got '%s'", retryAfter)
	}
}

// TestAnthropicAuthenticationFailure tests Anthropic authentication failure handling
func TestAnthropicAuthenticationFailure(t *testing.T) {
	authFailedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"type":      "authentication_error",
			"message":   "API key missing or invalid",
			"request_id": "test-request-123",
		})
	}))
	defer authFailedServer.Close()

	client := &http.Client{}
	reqBody := "{\"model\": \"claude-3-opus\", \"max_tokens\": 1024, \"messages\": [{\"role\": \"user\", \"content\": \"Hello\"}]}"
	req, _ := http.NewRequest("POST", authFailedServer.URL, strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "invalid-key")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

// TestCohereTimeout tests Cohere API timeout scenarios
func TestCohereTimeout(t *testing.T) {
	timeoutServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response that times out
		time.Sleep(35 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer timeoutServer.Close()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	reqBody := "{\"model\": \"command-r\", \"message\": \"Hello\"}"
	req, _ := http.NewRequest("POST", timeoutServer.URL, strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key")

	_, err := client.Do(req)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

// TestOpenAIStreaming tests OpenAI streaming response handling
func TestOpenAIStreaming(t *testing.T) {
	streamingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Send SSE response
		fmt.Fprintf(w, "data: {\"id\":\"chatcmpl-123\",\"object\":\"chat.completion.chunk\",\"created\":1234567890,\"model\":\"gpt-4\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"Hello\"}}]}\n\n")
		fmt.Fprintf(w, "data: {\"id\":\"chatcmpl-123\",\"object\":\"chat.completion.chunk\",\"created\":1234567890,\"model\":\"gpt-4\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\" world\"},\"finish_reason\":\"stop\"}]}\n\n")
		fmt.Fprintf(w, "data: [DONE]\n\n")
		w.(http.Flusher).Flush()
	}))
	defer streamingServer.Close()

	client := &http.Client{}
	reqBody := "{\"model\": \"gpt-4\", \"messages\": [{\"role\": \"user\", \"content\": \"Hi\"}], \"stream\": true}"
	req, _ := http.NewRequest("POST", streamingServer.URL, strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/event-stream") {
		t.Errorf("Expected Content-Type to contain 'text/event-stream', got '%s'", contentType)
	}

	// Read the streaming response
	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	if n == 0 {
		t.Error("Expected streaming data, got empty response")
	}
}

// TestAnthropicStreaming tests Anthropic streaming response handling
func TestAnthropicStreaming(t *testing.T) {
	streamingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Send SSE response
		fmt.Fprintf(w, "data: {\"type\":\"content_block_delta\",\"delta\":{\"type\":\"text_delta\",\"text\":\"Hello\"}}\n\n")
		fmt.Fprintf(w, "data: {\"type\":\"content_block_delta\",\"delta\":{\"type\":\"text_delta\",\"text\":\" world\"}}\n\n")
		fmt.Fprintf(w, "data: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\"},\"usage\":{\"output_tokens\":10}}\n\n")
		fmt.Fprintf(w, "data: [DONE]\n\n")
		w.(http.Flusher).Flush()
	}))
	defer streamingServer.Close()

	client := &http.Client{}
	reqBody := "{\"model\": \"claude-3-opus\", \"max_tokens\": 1024, \"messages\": [{\"role\": \"user\", \"content\": \"Hi\"}], \"stream\": true}"
	req, _ := http.NewRequest("POST", streamingServer.URL, strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "test-key")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/event-stream") {
		t.Errorf("Expected Content-Type to contain 'text/event-stream', got '%s'", contentType)
	}

	// Read the streaming response
	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	if n == 0 {
		t.Error("Expected streaming data, got empty response")
	}
}

// TestAPIKeyValidation tests API key validation with invalid credentials
func TestAPIKeyValidation(t *testing.T) {
	fixture := aiapifixtures.NewOpenAIChatCompletionsFixture()
	defer fixture.Close()

	t.Run("Missing API Key", func(t *testing.T) {
		client := &http.Client{}
		reqBody := "{\"model\": \"gpt-4\", \"messages\": [{\"role\": \"user\", \"content\": \"Hello\"}]}"
		req, _ := http.NewRequest("POST", fixture.URL(), strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		// No Authorization header
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid API Key", func(t *testing.T) {
		client := &http.Client{}
		reqBody := "{\"model\": \"gpt-4\", \"messages\": [{\"role\": \"user\", \"content\": \"Hello\"}]}"
		req, _ := http.NewRequest("POST", fixture.URL(), strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid-key")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})
}
