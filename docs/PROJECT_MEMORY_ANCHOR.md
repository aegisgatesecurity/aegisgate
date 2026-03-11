# AegisGate Project - Comprehensive Memory Anchor

**Last Updated:** 2026-02-26  
**Current Version:** 0.22.4  
**Repository:** C:\Users\Administrator\Desktop\Testing\aegisgate

---

## 📋 PROJECT OVERVIEW

**AegisGate** is an enterprise-grade AI security proxy platform written in Go. It provides comprehensive security scanning, threat intelligence, compliance management, and proxy services for AI/ML workloads.

### Core Purpose
- Secure AI/ML model inference and training traffic
- Provide compliance frameworks (NIST AI RMF, ISO/IEC 42001, OWASP AI Top 10, HIPAA, PCI-DSS)
- Threat intelligence integration (STIX 2.1, TAXII 2.1, MITRE ATLAS/ATT&CK)
- Enterprise SSO (SAML 2.0, OIDC)
- SIEM integration (Splunk, Elasticsearch, QRadar, etc.)

### Tech Stack
- **Language:** Go 1.24.0
- **Runtime:** Self-hosted Windows runner (no GCC, CGO disabled)
- **CI/CD:** GitHub Actions (version-sync.yml, release.yml)
- **Registry:** GHCR (GitHub Container Registry)
- **Platforms:** Linux (amd64, arm64, 386), Windows (amd64, arm64), macOS (amd64, arm64), FreeBSD (amd64)

---

## 🏗️ PROJECT STRUCTURE

```
aegisgate/
├── cmd/aegisgate/          # Main application entry point
├── cmd/debug/            # Debug utilities
├── pkg/auth/             # Authentication and authorization
├── pkg/certificate/      # Certificate management
├── pkg/compliance/       # Compliance frameworks (HIPAA, PCI, NIST, ISO)
├── pkg/config/           # Configuration management
├── pkg/core/             # Core module system, licensing, features
├── pkg/dashboard/        # Dashboard and security middleware
├── pkg/hash_chain/       # Hash chain validation
├── pkg/i18n/             # Internationalization (12 locales)
├── pkg/metrics/          # Performance metrics
├── pkg/ml/               # Machine learning components
├── pkg/pkiattest/        # PKI attestation (TPM 2.0)
├── pkg/proxy/            # HTTP/1.1 and HTTP/2 proxy with TLS MITM
├── pkg/reporting/        # Compliance reporting
├── pkg/sandbox/          # Feed-level sandboxing
├── pkg/scanner/          # Security scanning
├── pkg/security/         # CSRF, XSS, security headers
├── pkg/siem/             # SIEM integration (10+ platforms)
├── pkg/signature_verification/  # Digital signature verification
├── pkg/sso/              # SAML 2.0 and OIDC SSO
├── pkg/threatintel/      # STIX 2.1 / TAXII 2.1 threat intelligence
├── pkg/tls/              # TLS 1.3 support
├── pkg/trustdomain/      # Feed-specific trust domains
├── pkg/webhook/          # Webhook alerting system
├── pkg/websocket/        # WebSocket support
├── ui/frontend/          # Web UI (WCAG 2.1 AA compliant)
├── configs/              # Configuration templates
├── deploy/               # Deployment configurations
├── docs/                 # Documentation
└── tests/                # Test suites
```

---

## 📦 VERSION HISTORY (Recent)

| Version | Date | Key Changes |
|---------|------|-------------|
| **0.22.4** | 2026-02-26 | Fixed PowerShell exit code propagation (added exit 0) |
| **0.22.3** | 2026-02-26 | Fixed PowerShell OR operator (use call operator + LASTEXITCODE check) |
| **0.22.2** | 2026-02-26 | Fixed Syft download URL (amd64 vs x86_64) |
| **0.22.1** | 2026-02-26 | Removed -race flag (no GCC on Windows runner) |
| **0.22.0** | 2026-02-26 | Fixed PowerShell argument parsing (quoted -coverprofile) |
| **0.21.1** | 2026-02-25 | i18n enhancement (12 locales), test fixes |
| **0.21.0** | 2026-02-25 | Phase 2 complete (trust domains, sandbox, signatures) |
| **0.18.0** | 2026-02-23 | Threat intelligence, webhooks, enterprise SSO |
| **0.17.0** | 2026-02-23 | SIEM integration package |
| **0.16.0** | 2026-02-23 | Phase 1 production audit (CSRF, XSS, accessibility) |

---

## ⚙️ WORKFLOWS

### 1. version-sync.yml
**Triggers:** Push/PR to main/master, paths: VERSION, CHANGELOG.md, go.mod

**Jobs:**
- `check-version-consistency`: Validates VERSION matches CHANGELOG, go.mod Go version
- `check-documentation`: Ensures README mentions current version

### 2. release.yml
**Triggers:** Push tag v*, or workflow_dispatch with version input

**Jobs:**
1. `version-check`: Extract and verify version consistency
2. `test`: Run Go tests (NO -race flag - no GCC)
3. `build-binaries`: Cross-compile for 8 platforms
4. `generate-sbom`: Create SBOM with Syft
5. `create-release`: GitHub release with binaries, checksums, SBOM

**Build Matrix:** linux/amd64, linux/arm64, linux/386, windows/amd64, windows/arm64, darwin/amd64, darwin/arm64, freebsd/amd64

---

## 🔧 POWERSHELL GOTCHAS (Critical Lessons)

### 1. Exit Code Propagation
**Problem:** GitHub Actions PowerShell runner propagates the last command's exit code as the script's exit code, even when no explicit exit statement.

**Solution:** Always end PowerShell scripts in GitHub Actions with explicit `exit 0` if the script should succeed:
```yaml
shell: pwsh
run: |
  # ... commands that might fail but are handled ...
  exit 0  # Ensure script succeeds
```

### 2. Call Operator for External Commands
**Problem:** `gh release view` returns exit code 1 when release doesn't exist, which propagates as script failure.

**Wrong:**
```powershell
$release = gh release view $version --json tagName 2>$null || "none"
```

**Correct:**
```powershell
$release = & gh release view $version --json tagName 2>$null
if ($LASTEXITCODE -eq 0 -and $release) {
  # Release exists
}
```

### 3. No `||` or `&&` Operators
PowerShell doesn't have bash-style `||` and `&&` operators. Use `;` and `if` statements instead.

### 4. Argument Splitting at `=`
**Problem:** PowerShell splits `-coverprofile=coverage.out` into `-Coverprofile` and `.out` (interpreted as package).

**Solution:** Quote arguments containing `=`:
```yaml
run: go test -v "-coverprofile=coverage.out" ./...
```

### 5. Syft Download URL
**Problem:** Syft uses `windows_amd64.zip`, not `windows_x86_64.zip` for asset naming.

**Correct URL Pattern:**
```powershell
$url = "https://github.com/anchore/syft/releases/download/$latest/syft_$versionNum" + "_windows_amd64.zip"
```

---

## 🚫 RUNNER LIMITATIONS

### No GCC (No CGO, No Race Detector)
The self-hosted Windows runner does NOT have GCC installed.

**Implications:**
- Cannot use `-race` flag in `go test`
- Cannot use CGO-enabled packages
- Must use `CGO_ENABLED=0` for cross-compilation

**Workaround for Race Detection:**
- Run race detector in a Linux container/builder
- Or install MinGW/MSYS2 on the Windows runner

---

## 🔄 RELEASE PROCESS

### Standard Release
1. Update `VERSION` file with new version number (NO v prefix: `0.22.4`)
2. Add changelog entry to `CHANGELOG.md`
3. Commit: `git commit -m "chore: release v0.X.Y"`
4. Tag: `git tag v0.X.Y` (WITH v prefix)
5. Push: `git push origin main --tags`
6. Release workflow triggers automatically

### Manual Release (workflow_dispatch)
1. Go to Actions → Release → Run workflow
2. Enter version (e.g., `v0.22.4`)
3. Click Run workflow

---

## 📁 KEY FILES

| File | Purpose |
|------|---------|
| `VERSION` | Single source of truth for version |
| `CHANGELOG.md` | Version history (Keep a Changelog format) |
| `go.mod` | Go module definition (Go 1.24, golang.org/x/net) |
| `Makefile` | Build scripts with version references |
| `README.md` | Project documentation |
| `.github/workflows/version-sync.yml` | Version consistency checks |
| `.github/workflows/release.yml` | Release automation |

---

## 🏛️ ARCHITECTURE PHASES

### Phase 1 Complete (v0.16.0 - v0.17.0)
- Security middleware (CSRF, XSS, headers)
- WCAG 2.1 AA accessibility
- Performance benchmarks
- SIEM integration

### Phase 2 Complete (v0.19.0 - v0.21.0)
- Feed-specific trust domains
- Feed-level sandboxing
- PKI attestation (TPM 2.0)
- Digital signature verification
- Hash chain validation

### Current Development (v0.22.x)
- CI/CD stability improvements
- PowerShell compatibility fixes
- Self-hosted Windows runner adaptations

---

## 💡 PRO TIPS

1. **Always test PowerShell locally** if possible before committing workflow changes
2. **Use `2>$null`** to suppress stderr in PowerShell for expected failures
3. **Quote arguments with `=`** in PowerShell to prevent splitting
4. **Check self-hosted runner capabilities** before assuming features (e.g., GCC)
5. **Use `&` call operator** to capture output from commands that may fail
6. **Always end GitHub Actions PowerShell scripts with `exit 0`** if success is expected
7. **VERSION file should NOT have `v` prefix** (just `0.22.4`, not `v0.22.4`)
8. **Tag names SHOULD have `v` prefix** (`v0.22.4`, not `0.22.4`)

---

## 🐛 COMMON ISSUES & SOLUTIONS

| Issue | Cause | Solution |
|-------|-------|----------|
| Workflow exit code 1 | PowerShell exit code propagation | Add `exit 0` at end |
| `||` operator error | Invalid PowerShell syntax | Use `&` call + `$LASTEXITCODE` check |
| Cython/GCC errors | No C compiler on runner | Disable CGO, remove `-race` |
| `.out` package errors | PowerShell argument splitting | Quote arguments with `=` |
| Syft 404 error | Wrong architecture name | Use `amd64` not `x86_64` |

---

## 🔭 FUTURE CONSIDERATIONS

1. **Race Detection:** Consider adding a Linux runner with GCC for `-race` tests
2. **Docker Builds:** release.yml has `REGISTRY` and `IMAGE_NAME` but no Docker job
3. **SBOM Generation:** Currently generates CycloneDX and SPDX formats
4. **Multi-Architecture:** Already builds for 8+ platforms
5. **I18N:** 12 locale support (ar, de, en, es, fr, he, hi, ja, ko, pt, ru, zh)

---

## 📝 SESSION CONTEXT

### Recent Work (v0.22.0 - v0.22.4)
All 5 releases on 2026-02-26 were bug fixes for CI/CD issues on the self-hosted Windows runner:

1. **v0.22.0**: Fixed PowerShell argument parsing (`-coverprofile=coverage.out`)
2. **v0.22.1**: Removed `-race` flag (no GCC)
3. **v0.22.2**: Fixed Syft download URL
4. **v0.22.3**: Fixed PowerShell `||` operator issue
5. **v0.22.4**: Added explicit `exit 0` to prevent exit code propagation

### Key Commit Messages
- `fix: quote coverprofile argument to fix PowerShell parsing`
- `fix: remove race flag from release workflow - runner lacks GCC`
- `fix: correct Syft download URL for Windows`
- `fix: PowerShell || operator not valid, use call operator with $LASTEXITCODE`
- `fix: add explicit exit 0 to prevent PowerShell exit code propagation`

---

*This memory anchor documents the AegisGate project state as of 2026-02-26 for future reference and troubleshooting.*
