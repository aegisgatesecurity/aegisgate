package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aegisgatesecurity/aegisgate/pkg/scanner"
)

// Simple integration test using actual API
type TestServer struct {
	MockUpstream *httptest.Server
	Scanner      *scanner.Scanner
}

func setupTestServer(t *testing.T) *TestServer {
	ts := &TestServer{}

	// Setup mock upstream
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})
	ts.MockUpstream = httptest.NewServer(mux)

	// Setup scanner with default config
	ts.Scanner = scanner.New(scanner.DefaultConfig())

	return ts
}

func TestBasicIntegration(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.MockUpstream.Close()

	// Test scanner
	t.Run("Scanner Detection", func(t *testing.T) {
		testData := "Credit card: 4532015112830366"
		findings := ts.Scanner.Scan(testData)
		if len(findings) == 0 {
			t.Error("Expected credit card detection")
		}
	})

	// Test upstream
	t.Run("Upstream Health", func(t *testing.T) {
		resp, err := http.Get(ts.MockUpstream.URL + "/health")
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
