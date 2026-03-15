# FIPS 140-2/140-3 Readiness Report

**Document Version:** 1.0  
**AegisGate Version:** 0.31.0  
**Last Updated:** March 2026

---

## Executive Summary

This document describes the FIPS 140-2 and FIPS 140-3 compliance status of the AegisGate AI/LLM Security Gateway. AegisGate provides comprehensive FIPS readiness through its crypto module architecture.

---

## FIPS Compliance Status

### Current Status

| Requirement | Status | Notes |
|------------|--------|-------|
| FIPS 140-2 Readiness | ✅ Ready | Self-attestation supported |
| FIPS 140-3 Readiness | 🔄 In Progress | Framework in place |
| TLS 1.2+ | ✅ Compliant | Required minimum |
| TLS 1.3 | ✅ Compliant | Supported |
| Approved Cipher Suites | ✅ Compliant | All FIPS-approved |
| Key Size Requirements | ✅ Compliant | RSA 2048+ |

---

## Cryptographic Modules

### 1. Core FIPS Module (`pkg/crypto/fips/`)

```go
import "github.com/aegisgatesecurity/aegisgate/pkg/crypto/fips"

// Configure FIPS mode
fips.Configure(fips.Level140_2, fips.WithAudit(true))

// Run compliance check
report := fips.Check(fips.Level140_2)
fmt.Println(report.String())

// Validate key sizes
fips.ValidateKeySize("RSA", 2048)
```

### 2. Enhanced Crypto Module (`pkg/crypto/enhanced/`)

Provides additional cryptographic functions using `golang.org/x/crypto`:

- SHA-256, SHA3-256, SHA3-512
- BLAKE2b
- PBKDF2
- ChaCha20-Poly1305
- ED25519
- RSA-SHA256

### 3. TLS Configuration (`pkg/tls/fips_config.go`)

```go
import "github.com/aegisgatesecurity/aegisgate/pkg/tls"

// Get FIPS-compliant TLS config
config := tls.GetFIPSTLSConfig()
```

### 4. Proxy Integration (`pkg/proxy/fips_integration.go`)

```go
import "github.com/aegisgatesecurity/aegisgate/pkg/proxy"

// Configure proxy with FIPS
cfg := &proxy.ProxyTLSConfig{
    FIPSMode:      true,
    MinTLSVersion: tls.VersionTLS12,
}
```

---

## FIPS-Approved Algorithms

### Hash Functions

| Algorithm | Status | Use Case |
|-----------|--------|----------|
| SHA-256 | ✅ | Default hashing |
| SHA-384 | ✅ | Long messages |
| SHA-512 | ✅ | High security |
| SHA3-256 | ✅ | Alternative |
| BLAKE2b | ✅ | Fast hashing |

### Cipher Suites

| TLS 1.2 Cipher Suite | Status |
|---------------------|--------|
| ECDHE-RSA-AES256-GCM-SHA384 | ✅ |
| ECDHE-RSA-AES128-GCM-SHA256 | ✅ |
| ECDHE-ECDSA-AES256-GCM-SHA384 | ✅ |
| ECDHE-ECDSA-AES128-GCM-SHA256 | ✅ |
| RSA-AES256-GCM-SHA384 | ✅ |
| RSA-AES128-GCM-SHA256 | ✅ |

### Key Sizes

| Algorithm | Minimum Size | Status |
|-----------|--------------|--------|
| RSA | 2048 bits | ✅ |
| ECDSA | P-256 | ✅ |
| AES | 128 bits | ✅ |

---

## Configuration

### YAML Configuration

```yaml
fips:
  enabled: true
  level: "140-2"
  audit_logging: true
  approved_algorithms_only: true
  min_rsa_key_size: 2048
  min_tls_version: "1.2"
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| AEGISGATE_FIPS_ENABLED | Enable FIPS mode | false |
| AEGISGATE_FIPS_LEVEL | FIPS level (140-2/140-3) | 140-2 |
| AEGISGATE_FIPS_AUDIT_LOGGING | Enable audit logging | true |
| AEGISGATE_FIPS_MIN_RSA_KEY_SIZE | Minimum RSA key size | 2048 |

---

## Compliance Checks

Run compliance checks programmatically:

```go
// Full compliance report
report := fips.Check(fips.Level140_2)

// Individual checks
fips.ValidateKeySize("RSA", 2048)
fips.ValidateHashAlgorithm("SHA-256")

// Self-test
fips.SelfTest()
```

---

## Known Limitations

1. **Go Runtime**: Go's standard library crypto is not FIPS certified. For production FIPS environments, consider using a validated cryptographic module.

2. **Self-Attestation**: This release supports FIPS self-attestation. Official FIPS 140-3 certification requires NIST CMVP validation.

---

## Roadmap

### v0.32.0
- [ ] Enhanced key generation
- [ ] Hardware security module (HSM) integration

### v1.0.0
- [ ] FIPS 140-3 preparation
- [ ] CMVP validation support
- [ ] Third-party security audit

---

## References

- [FIPS 140-2](https://csrc.nist.gov/publications/detail/fips/140/2/final)
- [FIPS 140-3](https://csrc.nist.gov/publications/detail/fips/140/3/final)
- [SP 800-57](https://csrc.nist.gov/publications/detail/sp/800/57/part/1/rev-5/final)
- [NIST TLS Guidelines](https://csrc.nist.gov/publications/detail/sp/800/52/rev-2/final)

---

*This document is part of the AegisGate project. For support, please open an issue on GitHub.*
