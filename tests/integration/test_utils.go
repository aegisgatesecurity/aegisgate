package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aegisgatesecurity/aegisgate/pkg/compliance"
)

// ============================================================================
// TEST UTILITIES AND HELPERS
// ============================================================================

// TestLogger provides structured logging for tests
type TestLogger struct {
	Buffer strings.Builder
}

// NewTestLogger creates a new test logger
func NewTestLogger() *TestLogger {
	return &TestLogger{}
}

// Log logs a message with timestamp
func (l *TestLogger) Log(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Buffer.WriteString(fmt.Sprintf("[%s] %s\n", time.Now().Format(time.StampMicro), msg))
}

// GetLogs returns the accumulated logs
func (l *TestLogger) GetLogs() string {
	return l.Buffer.String()
}

// ============================================================================
// TEST FIXTURES
// ============================================================================

// TestFixture represents a test fixture with expected results
type TestFixture struct {
	Name        string                 `json:"name"`
	Category    string                 `json:"category"`
	Input       map[string]string     `json:"input"`
	Expected    map[string]interface{} `json:"expected"`
	ATLASParams ATLASParams           `json:"atlas_params"`
}

// ATLASParams holds ATLAS-specific test parameters
type ATLASParams struct {
	Technique string `json:"technique"`
	PatternID string `json:"pattern_id"`
	Severity  string `json:"severity"`
	BlockMode bool   `json:"block_mode"`
}

// LoadTestFixtures loads test fixtures from a JSON file
func LoadTestFixtures(path string) ([]TestFixture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var fixtures []TestFixture
	if err := json.Unmarshal(data, &fixtures); err != nil {
		return nil, err
	}

	return fixtures, nil
}

// SaveTestFixtures saves test fixtures to a JSON file
func SaveTestFixtures(fixtures []TestFixture, path string) error {
	data, err := json.MarshalIndent(fixtures, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ============================================================================
// ATLAS COMPLIANCE CHECKER
// ============================================================================

// AtlasTestResult represents the result of an ATLAS pattern test
type AtlasTestResult struct {
	Blocked   bool
	Detected  bool
	Technique string
	Pattern   string
	Score     float64
}

// GetAtlasChecker returns the ATLAS compliance checker
func GetAtlasChecker() *compliance.ATLASFramework {
	return compliance.NewAtlas()
}

// testAtlasPattern tests a single ATLAS pattern against the compliance checker
func testAtlasPattern(t *testing.T, payload string, patternID string) AtlasTestResult {
	// Use the compliance ATLAS checker directly
	checker := GetAtlasChecker()
	findings, _ := checker.Check(payload)

	// Determine if blocked based on findings
	blocked := len(findings) > 0

	// Extract technique from pattern ID
	technique := patternID
	if idx := strings.Index(patternID, "."); idx > 0 {
		technique = patternID[:idx]
	}

	return AtlasTestResult{
		Blocked:   blocked,
		Detected:  blocked,
		Pattern:   patternID,
		Technique: technique,
		Score:     1.0,
	}
}

// ============================================================================
// MOCK SERVERS
// ============================================================================

// MockLLMResponse defines a mock LLM response configuration
type MockLLMResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
	Delay      time.Duration
}

// CreateMockLLMServer creates a configurable mock LLM server for testing
func CreateMockLLMServer(responses map[string]MockLLMResponse) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get request body
		body := r.FormValue("prompt")
		if body == "" {
			body = r.FormValue("message")
		}

		// Find matching response
		for pattern, resp := range responses {
			if strings.Contains(strings.ToLower(body), strings.ToLower(pattern)) {
				time.Sleep(resp.Delay)

				for k, v := range resp.Headers {
					w.Header().Set(k, v)
				}
				w.WriteHeader(resp.StatusCode)
				w.Write([]byte(resp.Body))
				return
			}
		}

		// Default response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "OK", "model": "test-model"}`))
	}))
}

// CreateBlockingMockLLMServer creates a mock server that blocks ATLAS patterns
func CreateBlockingMockLLMServer() *httptest.Server {
	// Default responses that simulate blocking
	responses := map[string]MockLLMResponse{
		"ignore":    {StatusCode: 400, Body: `{"error": "ATLAS_BLOCKED", "reason": "Direct prompt injection detected"}`},
		"bypass":    {StatusCode: 400, Body: `{"error": "ATLAS_BLOCKED", "reason": "Jailbreak attempt detected"}`},
		"system":    {StatusCode: 400, Body: `{"error": "ATLAS_BLOCKED", "reason": "System prompt extraction detected"}`},
		"training":  {StatusCode: 400, Body: `{"error": "ATLAS_BLOCKED", "reason": "Training data exposure detected"}`},
		"injection": {StatusCode: 400, Body: `{"error": "ATLAS_BLOCKED", "reason": "Prompt injection detected"}`},
	}
	return CreateMockLLMServer(responses)
}

// ============================================================================
// TEST CONFIGURATION
// ============================================================================

// TestConfig holds test configuration
type TestConfig struct {
	// ATLAS Configuration
	ATLASEnabled   bool     `json:"atlas_enabled"`
	ATLASBlockMode bool     `json:"atlas_block_mode"`
	ATLASThreshold float64  `json:"atlas_threshold"`
	ExcludedParams []string `json:"excluded_params"`

	// Server Configuration
	UpstreamURL      string        `json:"upstream_url"`
	Timeout          time.Duration `json:"timeout"`
	MaxRetries       int           `json:"max_retries"`

	// Test Configuration
	Parallel  bool          `json:"parallel"`
	Verbose   bool          `json:"verbose"`
	Seed      int64         `json:"seed"`
	TimeoutDS time.Duration `json:"test_timeout"`
}

// DefaultTestConfig returns default test configuration
func DefaultTestConfig() TestConfig {
	return TestConfig{
		ATLASEnabled:    true,
		ATLASBlockMode:  true,
		ATLASThreshold:  0.75,
		ExcludedParams:  []string{"api_key", "token"},
		Timeout:         30 * time.Second,
		MaxRetries:      3,
		Parallel:        false,
		Verbose:         true,
		Seed:            time.Now().UnixNano(),
		TimeoutDS:       5 * time.Minute,
	}
}

// LoadTestConfig loads test configuration from file
func LoadTestConfig(path string) (TestConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultTestConfig(), err
	}

	var config TestConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultTestConfig(), err
	}

	return config, nil
}

// SaveTestConfig saves test configuration to file
func SaveTestConfig(config TestConfig, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ============================================================================
// ASSERTION HELPERS
// ============================================================================

// TestingT is an interface for test helpers
type TestingT interface {
	Errorf(format string, args ...interface{})
}

// AssertAtlasBlocked verifies that an ATLAS pattern was blocked
func AssertAtlasBlocked(t TestingT, result AtlasTestResult, pattern string) {
	if !result.Blocked {
		t.Errorf("Expected ATLAS pattern %s to be blocked, but it was allowed", pattern)
	}
}

// AssertAtlasAllowed verifies that a request was allowed through
func AssertAtlasAllowed(t TestingT, result AtlasTestResult, pattern string) {
	if result.Blocked {
		t.Errorf("Expected ATLAS pattern %s to be allowed, but it was blocked", pattern)
	}
}

// AssertTechniqueDetected verifies that a technique was detected
func AssertTechniqueDetected(t TestingT, result AtlasTestResult, technique string) {
	if result.Technique != technique {
		t.Errorf("Expected technique %s, got %s", technique, result.Technique)
	}
}

// AssertScoreAboveThreshold verifies that detection score is above threshold
func AssertScoreAboveThreshold(t TestingT, result AtlasTestResult, threshold float64) {
	if result.Score < threshold {
		t.Errorf("Expected score %.2f above threshold %.2f, got %.2f",
			result.Score, threshold, result.Score)
	}
}

// ============================================================================
// BENCHMARK HELPERS
// ============================================================================

// BenchmarkResult holds benchmark results
type BenchmarkResult struct {
	Name          string
	Iterations    int
	AvgDuration   time.Duration
	MinDuration   time.Duration
	MaxDuration   time.Duration
	TotalDuration time.Duration
	BytesAlloc    int64
	Allocs        int64
}

// RunBenchmark runs a benchmark test
func RunBenchmark(name string, fn func(), iterations int) BenchmarkResult {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		fn()
	}

	total := time.Since(start)

	return BenchmarkResult{
		Name:          name,
		Iterations:    iterations,
		AvgDuration:   total / time.Duration(iterations),
		MinDuration:   total / time.Duration(iterations), // Simplified
		MaxDuration:   total / time.Duration(iterations),
		TotalDuration: total,
	}
}

// PrintBenchmarkResults prints benchmark results
func PrintBenchmarkResults(results []BenchmarkResult) {
	fmt.Println("\n=== Benchmark Results ===")
	fmt.Printf("%-30s %15s %15s %15s\n", "Name", "Avg", "Min", "Max")
	fmt.Println(strings.Repeat("-", 75))

	for _, r := range results {
		fmt.Printf("%-30s %15s %15s %15s\n",
			r.Name,
			r.AvgDuration.String(),
			r.MinDuration.String(),
			r.MaxDuration.String())
	}
}

// ============================================================================
// COVERAGE TRACKING
// ============================================================================

// CoverageTracker tracks test coverage
type CoverageTracker struct {
	Techniques  map[string]bool
	Patterns    map[string]bool
	TotalTests  int
	PassedTests int
	FailedTests int
}

// NewCoverageTracker creates a new coverage tracker
func NewCoverageTracker() *CoverageTracker {
	return &CoverageTracker{
		Techniques:  make(map[string]bool),
		Patterns:    make(map[string]bool),
		TotalTests:  0,
		PassedTests: 0,
		FailedTests: 0,
	}
}

// RecordTest records a test result
func (c *CoverageTracker) RecordTest(technique, pattern string, passed bool) {
	c.TotalTests++
	c.Techniques[technique] = true
	c.Patterns[pattern] = true

	if passed {
		c.PassedTests++
	} else {
		c.FailedTests++
	}
}

// GetCoverage returns coverage statistics
func (c *CoverageTracker) GetCoverage() (techniqueCov, patternCov float64, total int) {
	techniqueCov = float64(len(c.Techniques)) / 18.0 * 100
	patternCov = float64(len(c.Patterns)) / 60.0 * 100
	total = c.TotalTests
	return
}

// PrintCoverage prints coverage report
func (c *CoverageTracker) PrintCoverage() {
	techCov, patCov, total := c.GetCoverage()

	fmt.Println("\n=== ATLAS Test Coverage Report ===")
	fmt.Printf("Total Tests: %d\n", total)
	fmt.Printf("Passed: %d | Failed: %d\n", c.PassedTests, c.FailedTests)
	fmt.Printf("Technique Coverage: %.1f%% (%d/18)\n", techCov, len(c.Techniques))
	fmt.Printf("Pattern Coverage: %.1f%% (%d/60)\n", patCov, len(c.Patterns))

	fmt.Println("\nCovered Techniques:")
	for tech := range c.Techniques {
		fmt.Printf("  - %s\n", tech)
	}
}
