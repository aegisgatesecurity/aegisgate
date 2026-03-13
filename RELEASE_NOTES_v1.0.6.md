# AegisGate v1.0.6 Release Notes

**Release Date**: March 13, 2026
**Version**: v1.0.6

---

## Overview

This release focuses on **security hardening**, **infrastructure improvements**, and **regulatory compliance enhancements**. Version 1.0.6 introduces new CI/CD security scanning capabilities, removes deprecated test files to reduce confusion, and expands the integration test suite for production readiness.

---

## Highlights

### CI/CD Security Enhancements
- **GitHub Actions CI Workflow** - Added comprehensive security scanning pipeline
  - **DAST (Dynamic Application Security Testing)** - Automated runtime vulnerability scanning
  - **Trivy Integration** - Container and dependency vulnerability scanning
  - **Gosec** - Go security checker for source code analysis
  - **CodeQL** - Semantic code analysis engine

### Infrastructure & Compliance
- **Trivy Configuration (.trivy.yaml)** - Added Trivy vulnerability scanner configuration
- **CODEOWNERS** - Implemented GitHub CODEOWNERS for automatic review assignments

### Removed Deprecated Files
- `pkg/compliance/atlas_clean.go.deleted` - Removed obsolete compliance file
- `pkg/proxy/proxy_benchmark_test.go` - Simplified benchmark tests
- `pkg/security/audit_test_incompatible.go.bak` - Removed incompatible test backup
- `pkg/security/csrf_test_incompatible.go.bak` - Removed incompatible test backup
- `pkg/security/panic_recovery_test_incompatible.go.bak` - Removed incompatible test backup

### New Integration Tests
- **ai_api_extended_test.go** - Extended AI API integration tests
- **atlas_patterns_test.go** - ATLAS compliance pattern tests
- **edge_cases_test.go** - Edge case handling tests
- **production_scenarios_test.go** - Production deployment scenarios

### New Rate Limiting Module
- **pkg/resilience/ratelimit/** - New dedicated rate limiting package
  - Token bucket algorithm implementation
  - Configurable rate limits per endpoint
  - Thread-safe concurrent request handling

---

## Dependencies Updated

### Go Module Updates (go.mod)
- Updated to latest compatible dependencies
- Improved security vulnerability coverage
- Enhanced module graph for better dependency resolution

---

## Breaking Changes

**None** - This is a backward-compatible release.

---

## Security Notes

1. All deprecated/incompatible test backup files have been removed to prevent confusion
2. CI/CD pipeline now includes automated security scanning before each release
3. Container images will be scanned for known CVEs using Trivy

---

## Migration Guide

### From v1.0.5
No migration steps required. Simply update to v1.0.6.

### From Earlier Versions
- Review [MIGRATION.md](MIGRATION.md) for v1.0.4+ migration instructions
- License keys from v1.0.3 and earlier require regeneration (see v1.0.4 migration)

---

## Testing

This release includes:
- **14 new integration tests** covering production scenarios
- **DAST automated scanning** in CI pipeline
- **Container vulnerability scanning** with Trivy
- **Source code security analysis** with Gosec and CodeQL

Run tests locally:
```bash
go test ./... -v
go test ./tests/integration/...
```

---

## Contributing

This release was made possible by contributions from the AegisGate community. Thank you!

---

## Links

- [Documentation](https://github.com/aegisgatesecurity/aegisgate/tree/main/docs)
- [Changelog](CHANGELOG.md)
- [Security Policy](SECURITY.md)
- [License](LICENSE)

---

## Contact

- **Security Issues**: security@aegisgate.io
- **General Support**: support@aegisgate.io
- **Sales**: sales@aegisgate.io

---

**Copyright © 2025-2026 AegisGate Security. All rights reserved.**