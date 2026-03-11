# ATLAS Compliance Integration Tests

Comprehensive test suite for validating ATLAS (MITRE ATT&CK for AI Systems) compliance patterns.

## Overview

This test suite validates the AegisGate proxy's ability to detect and block 60+ ATLAS compliance patterns across 18 MITRE ATLAS techniques.

## Files

| File | Description |
|------|-------------|
| `atlas_compliance_test.go` | Main test file with 60+ ATLAS pattern tests |
| `test_utils.go` | Test utilities, fixtures, and helpers |
| `test_runner.go` | Test runner with configuration support |
| `test_config.json` | Default test configuration |
| `fixtures/atlas_fixtures.json` | Test data fixtures for ATLAS patterns |
| `Makefile` | Make targets for running tests |
| `run_tests.bat` | Windows test runner script |
| `run_tests.sh` | Linux/Mac test runner script |

## Quick Start

### Run All ATLAS Tests

```bash
# Using Go test directly
go test -v -run Atlas ./tests/integration/...

# Run the comprehensive test
go test -v -run TestAtlasComplianceComprehensive ./tests/integration/...

# Run specific technique tests
go test -v -run TestAtlasComplianceDirectPromptInjection ./tests/integration/...
```

## ATLAS Techniques Covered

| Technique | Name | Patterns |
|-----------|------|----------|
| T1535 | Direct Prompt Injection | 5 |
| T1484 | LLM Jailbreak | 5 |
| T1632 | System Prompt Extraction | 5 |
| T1589 | Training Data Exposure | 5 |
| T1584 | Indirect Prompt Injection | 5 |
| T1658 | Adversarial Examples | 5 |
| T1648 | Serverless Compute Exploitation | 5 |
| T1600 | Vector Database Poisoning | 3 |
| T1613 | Content Injection | 3 |
| T1563 | Plugin Exploitation | 3 |
| T1622 | Defense Evasion | 3 |
| T1606 | Forge Web Credentials | 2 |
| T1621 | MFA Request Generation | 2 |
| T1548 | Abuse Elevation Control | 2 |
| T1490 | Inhibit System Recovery | 2 |
| T1498 | Network DoS | 2 |
| T1499 | Endpoint DoS | 2 |
| T1602 | Config Repository Exfiltration | 2 |

**Total: 18 techniques, 60+ patterns**

## Test Results

```
=== RUN   TestAtlasComplianceComprehensive
    Testing 18 ATLAS techniques
=== RUN   TestAtlasComplianceComprehensive/T1535
    Technique T1535 coverage verified
...
--- PASS: TestAtlasComplianceComprehensive (0.00s)
```

### Test Categories

- **Technique Tests** - Tests for each MITRE ATLAS technique
- **Benign Request Tests** - Ensures legitimate requests are not blocked
- **Proxy Integration Tests** - Tests the proxy's ATLAS compliance integration

## Test Configuration

Edit `test_config.json` to customize test behavior:

```json
{
  "atlas_enabled": true,
  "atlas_block_mode": true,
  "atlas_threshold": 0.75,
  "excluded_params": ["api_key", "token", "secret"],
  "timeout": 30,
  "parallel": false
}
```

## Test Utilities

The `test_utils.go` provides:

- `TestLogger` - Structured logging
- `TestFixture` - Test data management
- `MockLLMServer` - Configurable mock LLM server
- `CoverageTracker` - Test coverage tracking
- `BenchmarkResult` - Performance benchmarking
- `GetAtlasChecker()` - Access to ATLAS compliance checker

## Running Specific Tests

```bash
# Test specific technique
go test -v -run TestAtlasComplianceDirectPromptInjection ./tests/integration/...

# Test specific pattern
go test -v -run T1535.001 ./tests/integration/...

# Run benchmarks
go test -bench=. ./tests/integration/...
```

## Test Development

### Adding New Tests

1. Add test cases to `atlas_compliance_test.go`
2. Add fixtures to `fixtures/atlas_fixtures.json`
3. Update the test count in coverage reports

### Adding New Techniques

1. Create new test function in `atlas_compliance_test.go`
2. Add technique to `TestAtlasComplianceComprehensive`
3. Update this README with new technique info

## Expected Test Behavior

- Tests use the `compliance.NewAtlas()` checker directly
- Each test case validates detection of specific ATLAS patterns
- Benign requests should NOT be blocked
- Malicious patterns SHOULD be detected

## Notes

- Tests use `httptest.Server` for mock upstream servers
- Tests utilize the existing `compliance.AtlasManager` for pattern detection
- The test suite validates both detection capabilities and false positive prevention
