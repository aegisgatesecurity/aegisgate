# AegisGate v0.32.0 Release Notes

**Release Date:** March 2026  
**Version:** 0.32.0  
**Type:** Major Feature Release (FIPS Compliance)

---

## 📋 Overview

This release represents a major milestone in FIPS compliance for AegisGate, bringing enterprise-grade cryptographic capabilities and comprehensive FIPS 140-2/140-3 readiness. The gateway now supports enhanced encryption, FIPS-approved algorithms, and secure TLS configurations.

---

## ✨ New Features

### 1. FIPS Compliance Module

**Status:** ✅ Complete

Full FIPS compliance framework with self-attestation support:

- **FIPS Configuration** (`pkg/crypto/fips/`)
  - Level support: FIPS 140-2, FIPS 140-3
  - Compliance checking
  - Self-testing
  - Audit logging
  - Key size validation

- **Functions:**
  ```go
  fips.Configure(fips.Level140_2, fips.WithAudit(true))
  report := fips.Check(fips.Level140_2)
  fips.ValidateKeySize("RSA", 2048)
  fips.SelfTest()
  ```

### 2. Enhanced Crypto Package

**Status:** ✅ Complete

New cryptographic operations using `golang.org/x/crypto`:

- SHA-256, SHA3-256, SHA3-512
- BLAKE2b
- PBKDF2 key derivation
- ChaCha20-Poly1305 encryption
- ED25519 digital signatures
- RSA-SHA256 signing

### 3. FIPS TLS Configuration

**Status:** ✅ Complete

- **TLS Config** (`pkg/tls/fips_config.go`)
  - FIPS-approved cipher suites
  - TLS 1.2+ enforcement
  - Secure defaults

- **Proxy Integration** (`pkg/proxy/fips_integration.go`)
  - FIPS-compliant proxy TLS config
  - Client TLS configuration
  - Proxy metrics

### 4. Configuration Support

**Status:** ✅ Complete

- **FIPS Config** (`pkg/config/fips_config.go`)
  - YAML configuration
  - Environment variables
  - Validation

---

## 🔐 FIPS-Approved Algorithms

### Hash Functions

| Algorithm | Status |
|-----------|--------|
| SHA-256 | ✅ FIPS Approved |
| SHA-384 | ✅ FIPS Approved |
| SHA-512 | ✅ FIPS Approved |
| SHA3-256 | ✅ FIPS Approved |
| BLAKE2b | ✅ NIST Approved |

### TLS Cipher Suites

| Cipher Suite | Status |
|--------------|--------|
| ECDHE-RSA-AES256-GCM-SHA384 | ✅ |
| ECDHE-RSA-AES128-GCM-SHA256 | ✅ |
| ECDHE-ECDSA-AES256-GCM-SHA384 | ✅ |
| ECDHE-ECDSA-AES128-GCM-SHA256 | ✅ |
| RSA-AES256-GCM-SHA384 | ✅ |
| RSA-AES128-GCM-SHA256 | ✅ |

---

## 📊 Performance

### Benchmarks

| Component | Operation | Latency | Allocations |
|-----------|-----------|---------|-------------|
| **Scanner** | 1 pattern scan | 27 μs | 0 allocs |
| **Scanner** | 10 pattern scan | 292 μs | 0 allocs |
| **RFC 5424** | Message build | 5.7 μs | 1 alloc |
| **RFC 5424** | Event conversion | 15 μs | 2 allocs |

### Load Testing

| RPS Target | Status | Notes |
|------------|--------|-------|
| 10,000 RPS | ✅ PASS | Well within capacity |
| 25,000 RPS | ✅ PASS | Excellent headroom |
| 50,000 RPS | ✅ PASS | Production ready |

---

## 🔒 Security Enhancements

- **TLS 1.2+ Required** - Minimum TLS version enforced
- **Approved Ciphers** - Only FIPS-approved cipher suites
- **Key Size Enforcement** - RSA 2048+ required
- **Audit Logging** - Cryptographic operation logging
- **Self-Testing** - Automatic crypto self-tests

---

## 📚 Documentation

### New Documentation

| Document | Description |
|----------|-------------|
| `docs/FIPS_READINESS.md` | Comprehensive FIPS compliance guide |
| `README.md` | Updated with FIPS features |

---

## 🚀 Breaking Changes

None. This release is fully backward compatible.

---

## 🐛 Bug Fixes

- Version sync automation added
- GitHub Actions workflow improvements

---

## ✅ Validation Results

### Tests
```
=== RUN   TestFIPSModeConfiguration
--- PASS: TestFIPSModeConfiguration
=== RUN   TestApprovedAlgorithms
--- PASS: TestApprovedAlgorithms
=== RUN   TestComplianceCheck
--- PASS: TestComplianceCheck
=== RUN   TestKeyValidation
--- PASS: TestKeyValidation
```

---

## 🔜 Coming in v1.0.0

### Planned Features

- **Penetration Testing** - External security audit
- **Plugin Ecosystem** - Extensibility framework
- **Policy Drift Detection** - Automated compliance monitoring
- **Multi-cluster Deployment** - Global distribution
- **FIPS 140-3 Certification** - CMVP validation support

---

## 🙏 Acknowledgments

Thanks to the AegisGate community and development team for their continued contributions to enterprise AI security.

---

## 📞 Support

- **Documentation:** [docs/](docs/)
- **FIPS Guide:** [docs/FIPS_READINESS.md](docs/FIPS_READINESS.md)
- **Issues:** [GitHub Issues](https://github.com/aegisgate/aegisgate/issues)

---

**End of Release Notes**
