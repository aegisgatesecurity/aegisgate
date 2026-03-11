# Release Notes - v0.29.0

## Overview

**Release Date:** March 5, 2026
**Version:** v0.29.0
**Type:** Major Feature Release

---

## What's New in v0.29.0

This release introduces the **ATLAS Compliance Testing Framework** - a comprehensive test suite for validating MITRE ATT&CK for AI Systems compliance with 60+ pattern detection tests.

### 🚀 Major Features

#### 1. ATLAS Compliance Testing Framework

A complete testing infrastructure for validating AI/LLM security:

- **60+ Pattern Tests**: Comprehensive test coverage across 18 MITRE ATLAS techniques
- **Technique-Specific Tests**: Individual test functions for each ATLAS technique
- **Benign Request Validation**: Ensures legitimate requests are not blocked
- **Proxy Integration Testing**: End-to-end validation through the proxy

**Supported Techniques:**

| Technique ID | Name | Test Patterns |
|-------------|------|---------------|
| T1535 | Direct Prompt Injection | 5 patterns |
| T1484 | LLM Jailbreak | 5 patterns |
| T1632 | System Prompt Extraction | 5 patterns |
| T1589 | Training Data Exposure | 5 patterns |
| T1584 | Indirect Prompt Injection | 5 patterns |
| T1658 | Adversarial Examples | 5 patterns |
| T1648 | Serverless Compute Exploitation | 5 patterns |
| T1600 | Vector Database Poisoning | 3 patterns |
| T1613 | Content Injection | 3 patterns |
| T1563 | Plugin Exploitation | 3 patterns |
| T1622 | Defense Evasion | 3 patterns |
| T1606 | Forge Web Credentials | 2 patterns |
| T1621 | MFA Request Generation | 2 patterns |
| T1548 | Abuse Elevation Control | 2 patterns |
| T1490 | Inhibit System Recovery | 2 patterns |
| T1498 | Network DoS | 2 patterns |
| T1499 | Endpoint DoS | 2 patterns |
| T1602 | Config Repository Exfiltration | 2 patterns |

#### 2. Integration Test Infrastructure

New testing utilities and infrastructure:

- **Test Utilities** (`test_utils.go`):
  - `TestLogger` - Structured logging
  - `TestFixture` - Test data management
  - `MockLLMServer` - Configurable mock LLM server
  - `CoverageTracker` - Test coverage tracking
  - `BenchmarkResult` - Performance benchmarking
  - `GetAtlasChecker()` - Access to ATLAS compliance checker

- **Test Runner** (`test_runner.go`):
  - Configuration-based test execution
  - JSON report generation
  - Fixture-based test data loading
  - Test environment validation

- **Test Configuration** (`test_config.json`):
  - Configurable ATLAS settings
  - Server configuration
  - Test timeouts and retries

- **Test Fixtures** (`fixtures/atlas_fixtures.json`):
  - Structured test data for all ATLAS patterns
  - Expected results for validation

#### 3. Test Runner Scripts

Cross-platform test execution:

- **Windows** (`run_tests.bat`): Batch script for Windows environments
- **Linux/Mac** (`run_tests.sh`): Shell script for Unix-like systems
- **Makefile**: Make targets for common test operations

### 📝 Documentation

#### Comprehensive README Overhaul

The project README has been completely rewritten to include:

- **Architecture Diagrams**: Visual representation of system components
- **Quick Start Guides**: Step-by-step instructions for various deployment scenarios
- **Configuration Reference**: Complete configuration options with examples
- **Security Section**: TLS/mTLS, rate limiting, content scanning details
- **Compliance Section**: Framework support (MITRE ATLAS, OWASP, SOC 2, HIPAA, etc.)
- **ATLAS Testing Section**: Dedicated section for MITRE ATLAS testing
- **API Reference**: Complete endpoint documentation
- **Development Guide**: Building, testing, code quality
- **Deployment Guide**: Docker, Kubernetes, environment-specific configs

---

## Technical Details

### Dependencies

No new external dependencies introduced in this release.

### Go Version

- **Minimum Required:** Go 1.24.0
- **Tested Against:** Go 1.24.0

### Build Improvements

- Test binary properly excluded from repository root
- Integration test coverage increased to 85%+
- All ATLAS patterns covered by automated tests

### Test Coverage

| Category | Coverage |
|----------|----------|
| ATLAS Techniques | 100% (18/18) |
| ATLAS Patterns | 100% (60+/60+) |
| Integration Tests | 85%+ |
| Unit Tests | 90%+ |

---

## Security

### ATLAS Compliance

This release adds comprehensive MITRE ATLAS compliance testing:

- **Detection Coverage**: 60+ attack patterns across 18 techniques
- **False Positive Prevention**: Benign request testing ensures legitimate traffic isn't blocked
- **Proxy Integration**: Tests validate end-to-end detection through the proxy
- **Continuous Validation**: Automated tests run with every build

### Vulnerability Detection

The ATLAS framework covers:

1. **Prompt Injection Attacks**
   - Direct prompt injection (T1535)
   - Indirect prompt injection (T1584)
   - Jailbreak attempts (T1484)

2. **Data Exfiltration**
   - System prompt extraction (T1632)
   - Training data exposure (T1589)
   - Config repository exfiltration (T1602)

3. **Adversarial Attacks**
   - Adversarial examples (T1658)
   - Defense evasion (T1622)

4. **Resource Attacks**
   - Serverless exploitation (T1648)
   - DoS attacks (T1498, T1499)

---

## Breaking Changes

**None** - This release is fully backward compatible with v0.28.2.

---

## Migration Guide

### Upgrading from v0.28.2

1. **Update your installation:**
```bash
git pull origin main
go build -o aegisgate ./cmd/aegisgate
```

2. **Run ATLAS tests to verify:**
```bash
go test -v -run Atlas ./tests/integration/...
```

3. **Update configuration (optional):**
```yaml
# Enable ATLAS compliance
atlas:
  enabled: true
  block_mode: true
  threshold: 0.75
```

---

## Project Statistics

| Metric | Value |
|--------|-------|
| Total Go Packages | 35+ |
| Go Files | 200+ |
| Test Files | 15+ |
| CI/CD Workflows | 8 |
| Compliance Frameworks | 14+ |
| ATLAS Techniques | 18 |
| ATLAS Patterns | 60+ |
| Test Coverage | 85%+ |

---

## New Files in This Release

```
tests/integration/
├── atlas_compliance_test.go    # Main ATLAS test suite (982 lines)
├── test_utils.go              # Test utilities (368 lines)
├── test_runner.go             # Test runner (253 lines)
├── test_config.json           # Test configuration
├── fixtures/
│   └── atlas_fixtures.json    # ATLAS test fixtures
├── Makefile                   # Test Makefile
├── run_tests.bat             # Windows test runner
├── run_tests.sh               # Linux/Mac test runner
└── README.md                  # Integration tests documentation
```

---

## Known Issues

None reported.

---

## Deprecations

No features deprecated in this release.

---

## Upcoming Features (Next Release)

- Enhanced ATLAS pattern detection with ML-based classification
- Real-time ATLAS pattern updates via threat intelligence
- Extended compliance report generation
- Performance benchmarking dashboard

---

## Support

- **Issues:** Report via [GitHub Issues](https://github.com/aegisgatesecurity/aegisgate/issues)
- **Discussions:** Use [GitHub Discussions](https://github.com/aegisgatesecurity/aegisgate/discussions)
- **Security:** See [SECURITY.md](SECURITY.md)

---

## Credits

Contributors to this release:
- AegisGate Development Team
- Security Researchers (ATLAS pattern validation)

---

*Thank you for using AegisGate!*

**Protect Your AI Applications with AegisGate** 🛡️
