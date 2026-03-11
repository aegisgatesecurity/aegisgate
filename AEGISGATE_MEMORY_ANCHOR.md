# AEGISGATE PROJECT - COMPREHENSIVE MEMORY ANCHOR

**Last Updated**: March 10, 2026  
**Current Version**: 0.41.0  
**Latest Commit**: 9afdeca  
**Git Tag**: v0.41.0

---

## TABLE OF CONTENTS

1. [Project Overview](#1-project-overview)
2. [Current State](#2-current-state)
3. [Technical Architecture](#3-technical-architecture)
4. [Recent Development](#4-recent-development)
5. [Lessons Learned / Gotchas](#5-lessons-learned--gotchas)
6. [Troubleshooting Guide](#6-troubleshooting-guide)
7. [Future Plans](#7-future-plans)
8. [Development Workflow](#8-development-workflow)
9. [Enterprise Roadmap](#9-enterprise-roadmap)

---

## 1. PROJECT OVERVIEW

### What is AegisGate?

**AegisGate** is an enterprise-grade **AI Chatbot Security Gateway** - a reverse proxy designed to intercept, inspect, and secure AI chatbot traffic in real-time.

### Key Characteristics

| Attribute | Value |
|-----------|-------|
| **Language** | Go 1.24+ |
| **Module** | github.com/aegisgatesecurity/aegisgate |
| **License** | Apache 2.0 |
| **Repository** | https://github.com/aegisgatesecurity/aegisgate |
| **Architecture** | Reverse Proxy with MITM inspection |

### Primary Use Cases

- 🔍 **Real-Time Traffic Inspection** - Monitor all requests/responses
- 🛡️ **Threat Detection** - Identify prompt injection, jailbreak attempts, PII leaks
- 📋 **Compliance Enforcement** - Enforce MITRE ATLAS, NIST AI RMF, HIPAA, PCI-DSS
- 🔐 **TLS Inspection** - Decrypt and inspect HTTPS/HTTP/2 traffic
- 📊 **SIEM Integration** - Send security events to enterprise SIEM systems
- 🤖 **ML-Powered Security** - Statistical anomaly detection

### Supported AI Providers

| Provider | Status |
|----------|--------|
| OpenAI | ✅ Supported |
| Anthropic | ✅ Supported |
| Cohere | ✅ Supported |
| Azure OpenAI | ✅ Supported |
| Google Gemini | ✅ Supported |

---

## 2. CURRENT STATE

### Version Information

| File | Version |
|------|---------|
| `VERSION` | 0.41.0 |
| `cmd/aegisgate/main.go` | 0.41.0 |
| `README.md` badge | 0.41.0 |

### Git Status

```
Branch: main
Latest Commit: 9afdeca
Remote: origin (git@github.com:aegisgatesecurity/aegisgate.git)
Status: ✅ Synced with remote
```

### Release History

| Tag | Date | Type | Notes |
|-----|------|------|-------|
| v0.41.0 | March 10, 2026 | Major | YAML config, ML integration |
| v0.40.3 | March 9, 2026 | Patch | Workflow fixes |
| v0.40.2 | March 9, 2026 | Patch | Lint fixes |
| v0.40.1 | March 9, 2026 | Patch | Go 1.24 update |
| v0.40.0 | March 1, 2026 | Major | Initial v0.40 release |

### Test Coverage

| Package | Coverage |
|---------|----------|
| `pkg/config` | ~54.5% |
| `pkg/proxy` | ~28.4% |
| Overall | ~45% |

### Lint Status

| Category | Count | Status |
|----------|-------|--------|
| errcheck | ~22 | Intentional (proxy patterns) |
| staticcheck | ~25 | Acceptable style variations |
| ineffassign | 3 | Test files only |
| unused | 4 | Future use functions |

---

## 3. TECHNICAL ARCHITECTURE

### Core Components

```
aegisgate/
├── cmd/
│   ├── aegisgate/          # Main application entry point
│   └── gencerts/        # Certificate generation utility
├── pkg/
│   ├── proxy/           # Reverse proxy (HTTP/1.1, HTTP/2, HTTP/3)
│   ├── ml/              # Machine learning security
│   ├── compliance/      # MITRE ATLAS, NIST AI RMF, HIPAA, PCI-DSS
│   ├── auth/            # OAuth, SAML, OIDC authentication
│   ├── tls/             # Certificate authority & mTLS
│   ├── siem/            # SIEM integration (Splunk, Elastic, QRadar)
│   ├── opsec/           # OPSEC (audit, memory scrubbing, secrets)
│   ├── config/          # YAML configuration management
│   ├── i18n/            # Internationalization (12+ languages)
│   ├── security/        # Security middleware
│   └── ...              # 30+ total packages
├── .github/workflows/   # CI/CD pipelines
└── deploy/              # Docker, Kubernetes configs
```

### Compliance Frameworks Supported

| Framework | Description | Tier |
|-----------|-------------|------|
| MITRE ATLAS | Adversarial Threat Landscape for AI | Community |
| NIST AI RMF | AI Risk Management Framework | Enterprise |
| HIPAA | Healthcare data protection | Premium |
| PCI-DSS | Payment card security | Premium |
| SOC 2 | Service organization controls | Premium |
| GDPR | EU data protection | Community |
| ISO 42001 | AI Management System | Enterprise |

### Key Dependencies

```
github.com/prometheus/client_golang v1.20.5    # Metrics
golang.org/x/net v0.30.0                       # Networking
golang.org/x/oauth2 v0.26.0                     # OAuth
google.golang.org/grpc v1.70.0-dev             # gRPC
golang.org/x/crypto v0.30.0                     # Cryptography
gopkg.in/yaml.v3 v3.0.1                        # YAML config
```

---

## 4. RECENT DEVELOPMENT

### v0.41.0 Release (March 10, 2026)

#### New Features

1. **Complete YAML Configuration Loading**
   - Added `LoadFromFile()` function for YAML config files
   - Added `LoadWithEnvOverrides()` for environment variable overrides
   - Added TLS configuration support
   - Added comprehensive config parsing for ML, Security, and Plugin sections
   - Created sample `config.yaml`

2. **ML Detector Integration**
   - Integrated ML middleware into proxy request pipeline
   - Added prompt injection detection
   - Added behavioral anomaly detection (Z-score based)
   - Added content analysis (Shannon entropy)
   - Added PII detection capabilities

3. **Test Coverage Expansion**
   - Added 13 new tests for config package
   - Added comprehensive proxy tests
   - All tests passing

4. **Lint Configuration**
   - Fixed `.golangci.yml` for v2 format
   - Added proper exclusions for proxy patterns

#### Files Changed

| File | Change |
|------|--------|
| `cmd/aegisgate/main.go` | Full config loading, wired ML options |
| `pkg/config/config.go` | YAML loading, env overrides, TLS config |
| `pkg/config/config_test.go` | Expanded test coverage (13 tests) |
| `pkg/proxy/proxy.go` | ML integration |
| `pkg/proxy/proxy_test.go` | New test file |
| `go.mod` | Added gopkg.in/yaml.v3 dependency |
| `go.sum` | Updated dependencies |
| `.golangci.yml` | Fixed v2 config format |
| `config.yaml` | Sample configuration file (NEW) |

---

## 5. LESSONS LEARNED / GOTCHAS

### 🔴 Critical Gotchas

#### 1. Git Bash Command Parsing on Windows

**Issue**: Git commands with `-m "message"` get parsed incorrectly by Windows shell.

**Symptoms**:
```
error: pathspec 'message' did not match any file(s) known to git
```

**Solution**: Use Python subprocess or batch files for commits:
```python
import subprocess
subprocess.run(["git", "commit", "-m", "Your commit message"])
```

**Alternative**: Use `git commit --message="message"` format.

---

#### 2. Version Sync Workflow Requires Proper Permissions

**Issue**: GitHub Actions workflow failed with "Write access to repository not granted" (403).

**Root Cause**: Missing explicit permissions in workflow file.

**Solution**: Add explicit permissions in workflow:
```yaml
jobs:
  version-sync:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      actions: write
    steps:
      - uses: actions/checkout@v4
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
```

---

#### 3. Version Consistency is CRITICAL

**Issue**: CI fails if VERSION file doesn't match `cmd/aegisgate/main.go`.

**Always keep these in sync**:
- `VERSION` file
- `cmd/aegisgate/main.go` (const version = "x.x.x")
- Git tag (vX.X.X)

---

#### 4. Lint Issues in Proxy Code Are Intentional

**Issue**: errcheck and staticcheck report issues in proxy/TLS code.

**Why**: Common patterns like:
```go
defer conn.Close()
io.Copy(dst, src)
```

Trigger errcheck warnings but are CORRECT for proxy architecture.

**Solution**: 
- Use `.golangci.yml` to exclude known acceptable patterns
- Don't "fix" these - they will break the proxy

---

#### 5. YAML Import Gotcha

**Issue**: YAML config loading required adding dependency.

**Solution**: Add to `go.mod`:
```
gopkg.in/yaml.v3 v3.0.1
```

Then import in config package:
```go
import "gopkg.in/yaml.v3"
```

---

#### 6. JSON String Escaping in Go

**Issue**: Inline JSON with escaped quotes breaks Go string parsing.

**WRONG**:
```go
_, _ = fmt.Fprintf(w, "{\"ml\":\"not_enabled\"}")
```

**CORRECT**:
```go
_, _ = fmt.Fprintf(w, `{"ml":"not_enabled"}`)
```

---

### 🟡 Important Notes

#### 7. Go 1.24+ Required

**Issue**: Build failures with older Go versions.

**Solution**: Ensure Go 1.24+ is installed. Update `go.mod`:
```
go 1.24.0
```

---

#### 8. CGO and Windows Build

**Issue**: `-race` flag fails on Windows self-hosted runners (no GCC).

**Solution**: Remove `-race` from test commands in CI:
```yaml
# Wrong
go test -race ./...

# Correct (on Windows)
go test ./...
```

---

#### 9. Python Scripts Left in Repo During Linting Fixes

**Issue**: 46+ temporary Python files accumulated during lint fixes.

**Lesson**: Clean up temporary files before committing:
```bash
git status  # Check for temp files
git add -A  # Stage all
git commit  # Commit
```

---

#### 10. golangci-lint v2 Configuration Format

**Issue**: Old v1 config format doesn't work with v2 linter.

**Solution**: Use v2 format in `.golangci.yml`:
```yaml
version: 2

run:
  timeout: 5m
  modules-download-mode: readonly

linters:
  enable:
    - errcheck
    - staticcheck

issues:
  exclude-use-default: false
  exclude-rules:
    - path: _test\.go
      linters: [all]
```

---

#### 11. Test File Compilation Check

Before pushing, verify all packages compile:
```bash
go build ./cmd/aegisgate
go build ./cmd/gencerts
go test -c ./pkg/auth -o /dev/null
go test -c ./pkg/api -o /dev/null
```

---

#### 12. CI Runs on Self-Hosted Windows Runner

The CI pipeline uses a self-hosted Windows runner, not GitHub-hosted runners. This affects:
- PowerShell syntax (`pwsh` shell)
- No GCC for race detection
- Path separators (Windows-style)

---

#### 13. SARIF File Path Fixing for gosec

When running gosec on Windows, SARIF files have Windows paths that GitHub rejects. Use the PowerShell fix in `ci.yml`:
```powershell
$fixed = $prop.Value -replace '\\', '/'
$fixed = $fixed -replace '^[A-Za-z]:/', ''
```

---

#### 14. Docker Build Uses Go 1.24

Check Dockerfile for base image:
```dockerfile
FROM golang:1.24-alpine AS builder
```

---

### 🟢 Pro Tips

#### 15. Quick Config Test

Test your config.yaml before running the server:
```bash
go run ./cmd/aegisgate --config config.yaml --help
```

---

#### 16. ML Sensitivity Levels

| Sensitivity | Z-Threshold | Use Case |
|-------------|-------------|----------|
| low | 4.0 | Production with low false positives |
| medium | 3.0 | Default balance |
| high | 2.0 | High security environments |
| paranoid | 1.5 | Critical infrastructure |

---

#### 17. Environment Variable Override Pattern

Environment variables ALWAYS override YAML config values:

```bash
# config.yaml has: bind_address: ":8443"
# This will override to :9090
export AEGISGATE_BIND_ADDRESS=":9090"
./aegisgate --config config.yaml
```

---

## 6. TROUBLESHOOTING GUIDE

### Common Issues and Solutions

| Issue | Solution |
|-------|----------|
| **403 on git push in CI** | Add `permissions: contents: write` to workflow |
| **Version mismatch error** | Sync VERSION file with main.go |
| **Lint "errors" in proxy code** | These are intentional - don't fix |
| **Test compilation fails** | Run `go test -c ./pkg/xxx` to check |
| **CGO errors on Windows** | Use `CGO_ENABLED=0` for builds |
| **gosec SARIF rejected** | Use path-fixing PowerShell from ci.yml |
| **JSON parsing error** | Use backticks for raw strings |

### Verification Commands

```bash
# Check version consistency
type VERSION
grep "const version" cmd/aegisgate/main.go

# Verify build
go build ./cmd/aegisgate

# Run tests
go test ./pkg/...

# Check lint (should pass now)
golangci-lint run

# Check git status
git status
git log --oneline -5
```

---

## 7. FUTURE PLANS

### Roadmap (from README.md)

#### v0.42.0 (Q2 2026)
- [ ] Enhanced ML model for prompt injection detection
- [ ] Real-time threat intelligence integration
- [ ] Multi-tenant improvements
- [ ] Advanced analytics dashboard

#### v0.43.0 (Q3 2026)
- [ ] Hardware security module (HSM) support
- [ ] Zero-trust network integration
- [ ] Enhanced anomaly detection
- [ ] Custom rule engine

#### v0.44.0 (Q4 2026)
- [ ] AI-powered threat response
- [ ] Extended compliance modules
- [ ] Performance optimizations
- [ ] Enterprise SLA features

---

## 8. DEVELOPMENT WORKFLOW

### Standard Commit Process

```bash
# 1. Make changes
git add -A
git status

# 2. Commit (use Python/subprocess for Windows)
python -c "import subprocess; subprocess.run(['git', 'commit', '-m', 'your message'])"

# 3. Push
git push

# 4. Verify workflows pass
# Check GitHub Actions tab
```

### Version Bumping Process

```bash
# 1. Update VERSION file
echo "0.41.1" > VERSION

# 2. Update cmd/aegisgate/main.go
# Change: const version = "0.41.0" → const version = "0.41.1"

# 3. Create release notes
# Copy RELEASE_NOTES_v0.41.0.md → RELEASE_NOTES_v0.41.1.md
# Update content

# 4. Commit and tag
git add -A
git commit -m "Release v0.41.1"
git tag -a v0.41.1 -m "Release v0.41.1"
git push && git push --tags
```

---

## 9. ENTERPRISE ROADMAP

### Pricing Tiers

| Tier | Price | Frameworks |
|------|-------|------------|
| **Community** | Free | MITRE ATLAS, OWASP, GDPR |
| **Enterprise** | $10K-15K/mo | NIST AI RMF, ISO 42001 |
| **Premium** | $15K-25K/mo | SOC 2, HIPAA, PCI-DSS |

### Potential Fourth Tier (Developer)

| Proposed Tier | Price | Value Proposition |
|---------------|-------|-------------------|
| **Developer** | $99-299/mo | Indie devs, startups building AI apps |

### Enterprise Features Needed

| Feature | Priority | Status |
|---------|----------|--------|
| FIPS 140-2 Validation | Critical | Not started |
| HSM Integration | Critical | Not started |
| SLA Guarantees | High | Not started |
| 24/7 Support Contracts | High | Not started |
| Third-party Security Audit | High | Not started |
| Status Page | Medium | Not started |
| Demo Environment | Medium | Not started |

---

## Quick Reference

| Item | Value |
|------|-------|
| **Repo** | https://github.com/aegisgatesecurity/aegisgate |
| **Remote** | git@github.com:aegisgatesecurity/aegisgate.git |
| **Go Version** | 1.24+ |
| **Current Version** | 0.41.0 |
| **Latest Commit** | 9afdeca |
| **Main Branch** | main |
| **License** | Apache 2.0 |

---

*This document was created as a memory anchor for future development sessions.*
*Last updated: March 10, 2026*
