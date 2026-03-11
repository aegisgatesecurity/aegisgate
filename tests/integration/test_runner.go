package integration

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestRunnerConfig holds the test runner configuration
type TestRunnerConfig struct {
	ConfigPath      string
	FixturesPath    string
	OutputPath     string
	Parallel       bool
	Verbose        bool
	GenerateReport bool
	Filter         string
}

// TestReport represents the test execution report
type TestReport struct {
	Timestamp     time.Time       `json:"timestamp"`
	Duration      time.Duration   `json:"duration"`
	TotalTests    int             `json:"total_tests"`
	PassedTests   int             `json:"passed_tests"`
	FailedTests   int             `json:"failed_tests"`
	SkippedTests  int             `json:"skipped_tests"`
	Coverage      CoverageReport  `json:"coverage"`
	Results       []TestResult    `json:"results"`
}

// CoverageReport represents test coverage information
type CoverageReport struct {
	TechniquesCovered int      `json:"techniques_covered"`
	TechniquesTotal  int      `json:"techniques_total"`
	PatternsCovered int       `json:"patterns_covered"`
	PatternsTotal   int       `json:"patterns_total"`
	Techniques     []string  `json:"techniques"`
}

// TestResult represents an individual test result
type TestResult struct {
	Name       string        `json:"name"`
	Category   string        `json:"category"`
	Status     string        `json:"status"`
	Duration   time.Duration `json:"duration"`
	Techniques []string     `json:"techniques"`
	Error      string        `json:"error,omitempty"`
}

var (
	configPath      string
	fixturesPath    string
	outputPath      string
	parallel       bool
	verbose        bool
	generateReport bool
	filter         string
)

func init() {
	flag.StringVar(&configPath, "config", "test_config.json", "Path to test configuration")
	flag.StringVar(&fixturesPath, "fixtures", "fixtures/atlas_fixtures.json", "Path to test fixtures")
	flag.StringVar(&outputPath, "output", "test_report.json", "Path for test report output")
	flag.BoolVar(&parallel, "parallel", false, "Run tests in parallel")
	flag.BoolVar(&verbose, "verbose", false, "Verbose test output")
	flag.BoolVar(&generateReport, "report", true, "Generate test report")
	flag.StringVar(&filter, "filter", "", "Test filter pattern")
}

// TestMainWithConfig runs the test suite with configuration
func TestMainWithConfig(m *testing.M) {
	flag.Parse()

	// Set verbose flag
	if verbose {
		testing.Verbose()
	}

	// Run tests
	exitCode := m.Run()

	// Generate report if requested
	if generateReport && exitCode == 0 {
		generateTestReport()
	}

	os.Exit(exitCode)
}

// LoadConfig loads the test configuration
func LoadConfig(path string) (TestConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TestConfig{}, err
	}

	var config TestConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return TestConfig{}, err
	}

	return config, nil
}

// LoadFixtures loads the test fixtures
func LoadFixtures(path string) ([]TestFixture, error) {
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

// RunWithFixtures runs ATLAS tests using fixtures
func RunWithFixtures(t *testing.T, fixtures []TestFixture) {
	tracker := NewCoverageTracker()

	for _, fixture := range fixtures {
		name := fmt.Sprintf("%s/%s", fixture.Category, fixture.Name)

		if filter != "" && !strings.Contains(name, filter) {
			continue
		}

		t.Run(name, func(t *testing.T) {
			result := testAtlasPatternFromFixture(t, fixture)

			// Track coverage
			if fixture.ATLASParams.Technique != "" {
				expectedBlocked := false
				if eb, ok := fixture.Expected["blocked"].(bool); ok {
					expectedBlocked = eb
				}
				tracker.RecordTest(
					fixture.ATLASParams.Technique,
					fixture.ATLASParams.PatternID,
					result.Blocked == expectedBlocked,
				)
			}

			// Assert expected result
			expectedBlocked := fixture.Expected["blocked"].(bool)
			if result.Blocked != expectedBlocked {
				t.Errorf("Expected blocked=%v, got %v", expectedBlocked, result.Blocked)
			}
		})
	}

	// Print coverage at the end
	if verbose {
		tracker.PrintCoverage()
	}
}

// testAtlasPatternFromFixture tests a single ATLAS pattern from a fixture
func testAtlasPatternFromFixture(t *testing.T, fixture TestFixture) AtlasTestResult {
	// Extract prompt from fixture
	prompt := fixture.Input["prompt"]
	if prompt == "" {
		t.Fatal("Fixture missing 'prompt' in input")
	}

	patternID := fixture.ATLASParams.PatternID
	if patternID == "" {
		patternID = fixture.Category + ".001"
	}

	return testAtlasPattern(t, prompt, patternID)
}

// generateTestReport generates a JSON report of test results
func generateTestReport() {
	report := TestReport{
		Timestamp:     time.Now(),
		Duration:      0,
		TotalTests:    0,
		PassedTests:   0,
		FailedTests:   0,
		SkippedTests:  0,
		Coverage: CoverageReport{
			TechniquesTotal: 18,
			PatternsTotal:   60,
		},
		Results: []TestResult{},
	}

	// Save report
	data, _ := json.MarshalIndent(report, "", "  ")
	os.WriteFile(outputPath, data, 0644)

	fmt.Printf("Test report saved to: %s\n", outputPath)
}

// GetTestPaths returns the paths to test resources
func GetTestPaths() (string, string, string) {
	// Get the directory of the test file
	execDir, _ := os.Getwd()

	configPath = filepath.Join(execDir, "test_config.json")
	fixturesPath = filepath.Join(execDir, "fixtures", "atlas_fixtures.json")
	outputPath = filepath.Join(execDir, "test_report.json")

	return configPath, fixturesPath, outputPath
}

// ValidateTestEnvironment validates the test environment
func ValidateTestEnvironment() error {
	configPath, fixturesPath, _ := GetTestPaths()

	// Check config file
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", configPath)
	}

	// Check fixtures directory
	if _, err := os.Stat(fixturesPath); os.IsNotExist(err) {
		return fmt.Errorf("fixtures file not found: %s", fixturesPath)
	}

	return nil
}

// PrintTestSummary prints a summary of test results
func PrintTestSummary(tracker *CoverageTracker) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("ATLAS COMPLIANCE TEST SUMMARY")
	fmt.Println(strings.Repeat("=", 60))

	techCov, patCov, total := tracker.GetCoverage()

	fmt.Printf("Total Tests Run:    %d\n", total)
	fmt.Printf("Tests Passed:       %d\n", tracker.PassedTests)
	fmt.Printf("Tests Failed:       %d\n", tracker.FailedTests)
	fmt.Printf("\n")
	fmt.Printf("Technique Coverage: %.1f%% (%d/18)\n", techCov, len(tracker.Techniques))
	fmt.Printf("Pattern Coverage:   %.1f%% (%d/60)\n", patCov, len(tracker.Patterns))
	fmt.Println(strings.Repeat("=", 60))
}

