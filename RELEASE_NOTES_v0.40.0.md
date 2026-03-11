# Release Notes - v0.40.0

## Version 0.40.0 - Test Coverage & Quality Assurance Release

**Release Date:** March 2026

---

## Overview

Version 0.40.0 is a major release focusing on **test coverage expansion and quality assurance** for the AegisGate AI Chatbot Security Gateway. This release significantly improves test coverage across critical security packages, introduces comprehensive ML detector documentation, and enhances the overall code quality withgolangci-lint v2 integration.

---

## Highlights

### 🚀 New Features

#### 1. ML Detector Documentation (`docs/ML_DETECTORS.md`)
Comprehensive technical documentation covering all ML detection subsystems:

- **Anomaly Detector** - Statistical traffic analysis with Z-score and IQR methods
- **Prompt Injection Detector** - 15+ pattern categories:
  - Direct instruction overrides
  - Role manipulation attacks
  - Jailbreak attempts (DAN, developer mode)
  - Prompt extraction attempts
  - Hidden token injection ([INST], |endoftext|)
  - Base64/obfuscated content
  - Context switching attacks
- **Token Smuggling Detector** - 7 LLM-specific token formats:
  - Llama2 instruction tokens
  - ChatML tokens
  - OpenAI tokens
  - Vicuna chat tokens
  - Anthropic Claude tokens
  - XML tag injection
  - Base64 encoded instructions
- **Unicode Attack Detector** - 7 obfuscation techniques:
  - Homoglyph attacks
  - Zero-width characters
  - RTL override characters
  - Unicode escape sequences
  - Mixed script detection
  - Fullwidth character attacks
  - Excessive whitespace/invisible characters
- **Context Manipulation Detector** - 8 attack vectors:
  - Conversation reset attempts
  - Memory manipulation
  - Persona override attempts
  - System prompt extraction
  - Constraint breaking
  - Output formatting manipulation
  - Privilege escalation
  - Conflicting instructions
- **Content Analyzer** - 8 PII pattern types:
  - Social Security Numbers
  - Credit card numbers
  - Email addresses
  - Phone numbers
  - IP addresses
  - API keys
  - Passwords
  - Private keys
- **Behavioral Analyzer** - 3 anomaly categories:
  - Request frequency anomalies (>10 req/s)
  - Path diversity anomalies (scanning)
  - Data volume anomalies (exfiltration)
- **Combined Detector** - Unified facade with weighted scoring:
  - Prompt injection: 35%
  - Token smuggling: 25%
  - Unicode attacks: 20%
  - Context manipulation: 20%
- **Prometheus Metrics Integration** - Complete metrics coverage

#### 2. golangci-lint v2 Integration
- Updated to golangci-lint v2.x with modern configuration
- New `.golangci.yml` configuration file
- Optimized CI pipeline with parallel execution
- Improved linting rules:
  - `errcheck` - Error handling
  - `gosimple` - Simplification suggestions
  - `govet` - Vet analysis
  - `ineffassign` - Unused assignments
  - `staticcheck` - Static analysis
  - `unused` - Unused code detection

#### 3. Multi-Platform Docker Support
- New `Dockerfile.multiplatform` for cross-platform builds
- `build_multiplatform.sh` script for automated builds
- Support for:
  - linux/amd64
  - linux/arm64
  - linux/arm/v7
  - windows/amd64
  - darwin/amd64
  - darwin/arm64

---

## Test Coverage Improvements

### Package Coverage Results

#### pkg/security: 9.7% → 78.4% ✅
- Added `security_coverage_test.go` with comprehensive tests
- Test categories:
  - CSRF protection tests
  - Security headers tests
  - XSS protection tests
  - Panic recovery tests
  - Audit logging tests

#### pkg/compliance: 22.1% → 33.9% ✅
- Added `compliance_coverage_test.go` with framework tests
- Test categories:
  - HIPAA compliance tests
  - PCI-DSS compliance tests
  - SOC2 framework tests
  - OWASP compliance tests
  - ATLAS framework tests

#### pkg/proxy: 28.0% → 31.2% ✅
- Added `proxy_coverage_test.go` with integration tests
- Test categories:
  - HTTP/1.1 proxy tests
  - HTTP/2 proxy tests
  - HTTP/3 proxy tests (now passing)
  - ML middleware integration tests
  - Tenant isolation tests

### Test Execution
```bash
# Run all tests with coverage
go test -coverprofile=coverage.out ./...
go test -coverprofile=coverage.out ./pkg/security/...
go test -coverprofile=coverage.out ./pkg/compliance/...
go test -coverprofile=coverage.out ./pkg/proxy/...

# View coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Benchmark Tests
- Proxy benchmarks: `pkg/proxy/proxy_benchmark_test.go`
- Scanner benchmarks: `pkg/scanner/scanner_benchmark_test.go`
- Security middleware benchmarks: `pkg/security/middleware_bench_test.go`

Run benchmarks:
```bash
go test -bench=. -benchmem ./pkg/proxy/...
go test -bench=. -benchmem ./pkg/scanner/...
go test -bench=. -benchmem ./pkg/security/...
```

---

## Bug Fixes

### golangci-lint Path Exclusions
- **Issue**: Path exclusions in `.golangci.yml` not working with certain linters
- **Fix**: Used `--tests=false` flag in CI to exclude test files from linting
- **Result**: Clean lint runs with proper test file handling

### HTTP/3 Test Failures
- **Issue**: HTTP/3 integration tests were skipped or failing
- **Fix**: Updated test configuration and fixed compatibility issues
- **Result**: HTTP/3 tests now pass successfully

### Error Handling Improvements
- **Files**: `pkg/security/audit.go`, `pkg/tls/manager.go`, `cmd/aegisgate/main.go`
- **Fix**: Added proper error checking and handling
- **Result**: Improved reliability and error reporting

### Test Compilation Errors
- **Issue**: Multiple test files had compilation errors
- **Fixes**:
  - Fixed import paths
  - Corrected typos in test names
  - Aligned API calls with current implementations
- **Result**: All test files compile and run successfully

---

## Documentation Updates

### New Documentation
- **`docs/ML_DETECTORS.md`** - 397 lines of comprehensive ML detector documentation
  - Architecture overview
  - API reference for each detector
  - Configuration examples
  - Pattern definitions
  - Metrics documentation
  - Troubleshooting guide

### Updated Documentation
- **`README.md`** - Updated to v0.40.0 with:
  - Current version badge
  - Test coverage table
  - Latest features
  - NIST AI RMF compliance
  - MITRE ATLAS framework

---

## CI/CD Improvements

### GitHub Actions Workflow Updates

#### `.github/workflows/ci.yml`
- Updated to golangci-lint v2
- Parallel test execution
- Optimized build steps
- Better artifact handling

#### `.github/workflows/ml-pipeline.yml`
- Enhanced ML model testing
- Improved test coverage reporting
- Integration with codecov

### Build Pipeline Enhancements
- Multi-platform Docker builds
- SBOM generation
- Security scanning
- Version synchronization

---

## Dependencies

### No New External Dependencies
All functionality uses existing dependencies:
- `github.com/prometheus/client_golang` v1.20.5
- `golang.org/x/net` v0.30.0
- `golang.org/x/oauth2` v0.26.0
- `golang.org/x/crypto` v0.30.0
- `google.golang.org/grpc` v1.70.0-dev

### Go Version
- **Minimum**: Go 1.24.0
- **Recommended**: Go 1.24+

---

## Migration Guide

### Upgrading from v0.39.0

1. **Update Configuration**
   ```yaml
   ml:
     enabled: true
     sensitivity: 75
   ```

2. **Run Tests**
   ```bash
   go test -v ./...
   ```

3. **Check Linting**
   ```bash
   golangci-lint run ./...
   ```

4. **Update Docker Images**
   ```bash
   docker pull aegisgate/aegisgate:v0.40.0
   ```

---

## Performance

### ML Analysis Performance
- Analysis overhead: <10ms per request (typical)
- Memory usage: ~50MB baseline + ~10MB per 1000 active behavioral profiles
- Scalable horizontally with stateless detectors

### Test Execution
- Full test suite: ~2-5 minutes (depending on hardware)
- Coverage report generation: ~30 seconds

---

## Breaking Changes

### None

This release is fully backward compatible with v0.39.0.

---

## Security Enhancements

1. **Improved Error Handling**
   - Proper error propagation
   - Reduced silent failures
   - Better logging of security events

2. **CSRF Protection**
   - Enhanced token generation
   - Better validation logic

3. **Security Headers**
   - Comprehensive header configuration
   - CSP, HSTS, X-Frame-Options

---

## Compliance

### Supported Frameworks (v0.40.0)

| Framework | Controls | Evidence |
|-----------|----------|----------|
| **MITRE ATLAS** | 50+ | Automated |
| **NIST AI RMF** | Core Functions | Automated |
| **HIPAA** | 54 | Automated |
| **PCI-DSS** | 42 | Automated |
| **SOC2 Type II** | 33 | Automated |
| **OWASP Top 10** | 10 | Automated |

---

## Coming in Future Releases

- ML model training pipeline with custom models
- Real-time threat intelligence integration
- Enhanced anomaly detection with ensemble methods
- Custom pattern definitions via UI/API
- Advanced reporting dashboard
- Automated compliance reporting

---

## Contributors

- AegisGate Development Team

---

## Acknowledgments

Thank you to the open-source community for feedback, bug reports, and contributions.

---

## Links

- **Documentation:** https://aegisgatesecurity.io
- **GitHub Issues:** https://github.com/aegisgatesecurity/aegisgate/issues
- **Release Downloads:** https://github.com/aegisgatesecurity/aegisgate/releases
- **Helm Charts:** https://aegisgatesecurity.io

---

*For older release notes, see [RELEASE_NOTES_v0.39.0.md](RELEASE_NOTES_v0.39.0.md) and earlier.*
