# Release v1.0.15 - Security Hardening Release

**Release Date:** March 19, 2026  
**Type:** Security-Focused Patch Release  
**Status:** ✅ Stable

---

## Release Highlights

This release is focused on addressing security vulnerabilities and hardening the AegisGate platform for production deployment. All users are encouraged to upgrade.

### Key Changes

| Category | Change | Impact |
|----------|--------|--------|
| 🛡️ **Security** | Upgraded Go to 1.25.8 | Fixed 22 CVEs in standard library |
| 🛡️ **Security** | Production Docker hardening | Distroless base, non-root, read-only FS |
| 🛡️ **Security** | SBOM published | Enterprise vulnerability scanning support |
| 📚 **Docs** | SECURITY_HARDENING.md | Sales enablement & compliance mapping |
| 🔧 **Build** | Dockerfile.production | New secure container build |

---

## Security Fixes

### CVE Resolution (22 Vulnerabilities Fixed)

| CVE ID | Go Version | Issue | Severity |
|--------|------------|-------|----------|
| GO-2026-4603 | 1.25.8 | XSS in html/template | Critical |
| GO-2026-4602 | 1.25.8 | FileInfo escape in os | Medium |
| GO-2026-4601 | 1.25.8 | IPv6 parsing bypass | Medium |
| GO-2026-4341 | 1.24.12 | Memory exhaustion in net/url | High |
| GO-2026-4340 | 1.24.12 | TLS handshake MITM | High |
| GO-2026-4337 | 1.24.13 | TLS session resumption | High |
| GO-2025-4012 | 1.24.8 | Cookie memory exhaustion | High |
| GO-2025-4013 | 1.24.8 | DSA certificate panic | Medium |
| GO-2025-4175 | 1.24.11 | Wildcard cert validation | Medium |
| GO-2025-4155 | 1.24.11 | x509 cert DoS | Medium |
| GO-2025-4008 | 1.24.8 | TLS ALPN info leak | Low |
| GO-2025-4009 | 1.24.8 | PEM quadratic complexity | Medium |
| GO-2025-4010 | 1.24.8 | IPv6 URL parsing bypass | Medium |
| GO-2025-4011 | 1.24.8 | ASN.1 memory exhaustion | Medium |
| GO-2025-4007 | 1.24.9 | x509 complexity DoS | Medium |
| GO-2025-3849 | 1.23.12 | database/sql bug | Low |
| GO-2025-3751 | 1.23.10 | Sensitive headers redirect | Low |
| GO-2025-3563 | 1.23.8 | Request smuggling | Low |
| GO-2025-3447 | 1.23.6 | ECDSA timing (ppc64le) | Low |
| GO-2025-3420 | 1.23.5 | Sensitive headers redirect | Low |
| GO-2025-3373 | 1.23.5 | IPv6 zone ID bypass | Low |
| GO-2025-3750 | 1.23.10 | Windows file race | Low |

---

## New Features

### 1. Production-Hardened Docker

**File:** `Dockerfile.production`

Security improvements:
- ✅ Distroless base image (zero CVE by design)
- ✅ Non-root user (UID 65532)
- ✅ Read-only filesystem with tmpfs
- ✅ Dropped ALL Linux capabilities
- ✅ Network isolation
- ✅ Resource limits

### 2. Secure Docker Compose

**File:** `docker-compose.production.yml`

Production-ready configuration with:
- seccomp profile
- Kernel hardening (sysctls)
- Resource limits per tier
- Read-only root filesystem
- No new privileges flag

### 3. Software Bill of Materials (SBOM)

**File:** `sbom.json`

- CycloneDX format
- All dependencies with checksums
- CPE identifiers for CVE scanning
- License compliance data

### 4. Security Hardening Documentation

**File:** `SECURITY_HARDENING.md`

Sales enablement document including:
- CVE details and risk assessment
- Docker security features explained
- Compliance mappings (CIS, PCI, SOC2, GDPR, NIST)
- Kubernetes PSS integration
- Pentesting recommendations

---

## Version Updates

| Component | Previous | Current |
|-----------|----------|---------|
| Go | 1.23.0 | 1.25.8 |
| Binary Version | v1.0.14 | v1.0.15 |
| README Badge | v1.0.14 | v1.0.15 |
| Docker Base | golang:1.23-alpine | golang:1.25.8-alpine |

---

## Testing Performed

### ✅ All Tests Pass

| Test Suite | Status |
|------------|--------|
| Unit Tests | ✅ 53 packages pass |
| Build | ✅ Binary compiles |
| Docker Build | ✅ 27.2MB image builds |
| Container Runtime | ✅ Starts successfully |
| Health Endpoint | ✅ /health responds |
| Version Endpoint | ✅ Returns v1.0.15 |

### Performance Metrics (Go 1.23 → 1.25)

| Metric | Notes |
|--------|-------|
| Binary Size | ~12MB (uncompressed) |
| Docker Image | 27.2MB |
| Build Time | ~45 seconds |
| Test Suite | 53 packages, all pass |

---

## Breaking Changes

**None.** This release is fully backward compatible.

The Go version upgrade (1.23 → 1.25) is a drop-in replacement - no code changes required.

---

## Upgrade Instructions

### Docker Users

```bash
# Pull latest
docker pull ghcr.io/aegisgatesecurity/aegisgate:v1.0.15

# Or build from source
docker build -f Dockerfile.production -t aegisgate:v1.0.15 .
```

### Binary Users

```bash
# Rebuild
go build -o aegisgate ./cmd/aegisgate

# Verify version
./aegisgate -version
# Output: AegisGate v1.0.15 (commit: xxx, date: 2026-03-19)
```

---

## Known Issues

None reported.

---

## Next Steps

- [ ] External penetration testing (recommended before enterprise sales)
- [ ] Security audit report
- [ ] Performance benchmarking under load

---

## Support

- **Documentation:** https://aegisgate.io/docs
- **Issues:** https://github.com/aegisgatesecurity/aegisgate/issues
- **Security:** security@aegisgatesecurity.io (48-hour response)

---

## Contributors

This release was prepared by the AegisGate Security Team.

---

## Changelog

- **v1.0.15** (2026-03-19): Security hardening release - Go 1.25.8, Docker hardening, SBOM
- **v1.0.14** (2026-03-18): Version sync, workflow fixes
- **v1.0.13** (2026-03-15): Feature additions

---

**License:** MIT  
**Copyright:** 2025-2026 AegisGate Security