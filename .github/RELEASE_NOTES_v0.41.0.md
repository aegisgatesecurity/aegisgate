# Release Notes - Padlock v0.41.0

## Release Date
March 10, 2026

## Overview
This release brings complete YAML configuration loading, full ML detector integration into the proxy pipeline, comprehensive test coverage expansion, and various bug fixes.

## What's New

### 1. Complete YAML Configuration Loading
- Added `LoadFromFile()` function to load configuration from YAML files
- Added `LoadWithEnvOverrides()` for environment variable overrides
- Added TLS configuration support
- Added comprehensive config parsing for ML, Security, and Plugin sections
- Sample `config.yaml` file included

### 2. ML Detector Integration
- Integrated ML middleware into the proxy request pipeline
- Added prompt injection detection
- Added behavioral anomaly detection (Z-score based)
- Added content analysis (Shannon entropy)
- Added PII detection capabilities

### 3. Test Coverage Expansion
- Added 13 new tests for config package
- Added comprehensive proxy tests
- All tests passing

### 4. Lint Configuration
- Fixed `.golangci.yml` for v2 format
- Added proper exclusions for proxy patterns

## Changes

### Modified Files
- `cmd/padlock/main.go` - Full config loading, wired ML options
- `pkg/config/config.go` - YAML loading, env overrides, TLS config
- `pkg/config/config_test.go` - Expanded test coverage
- `pkg/proxy/proxy.go` - ML integration
- `pkg/proxy/proxy_test.go` - New test file
- `go.mod` - Added gopkg.in/yaml.v3 dependency
- `go.sum` - Updated dependencies
- `.golangci.yml` - Fixed v2 config format

### New Files
- `config.yaml` - Sample configuration file
- `pkg/proxy/proxy_test.go` - Comprehensive proxy tests

## Upgrading

### From v0.40.x
No breaking changes. Simply update your binary and ensure you're using Go 1.24+.

### New Configuration
You can now use a YAML configuration file:
```bash
./padlock --config config.yaml
```

Or continue using environment variables (they override file values):
```bash
export PADLOCK_UPSTREAM=https://api.openai.com
./padlock --config config.yaml
```

## Known Issues
- None

## Contributors
- Thank you to all contributors!

## Next Release (v0.42.0 - Planned Q2 2026)
- Enhanced ML model for prompt injection detection
- Real-time threat intelligence integration
- Multi-tenant improvements
- Advanced analytics dashboard

## Links
- GitHub Release: https://github.com/jcolvin1056/padlock/releases/tag/v0.41.0
- Documentation: https://github.com/jcolvin1056/padlock#readme

---

**Full Changelog**: https://github.com/jcolvin1056/padlock/compare/v0.40.3...v0.41.0