# AegisGate v1.0.3 Release Notes

**Release Date**: March 13, 2026
**Version**: v1.0.3
**Codename**: "Security Hardening"

---

## Security Critical - Immediate Action Required

This release addresses **CRITICAL security vulnerabilities** identified in the license validation system. 

### Breaking Changes

**OLD LICENSE KEYS WILL NOT WORK.** Users must generate new cryptographically signed licenses.

---

## Security Fixes

### Critical: DEV_MODE Bypass Removed
- **Issue**: Removed `AEGISGATE_DEV_MODE` environment variable bypass
- **Severity**: CRITICAL - Allowed unlimited feature access for Community users
- **Fixed**: All licenses now require cryptographic validation

### Cryptographic License Validation Implemented
- **HMAC-SHA256**: For Developer ($29/mo) and Professional ($99/mo) tiers
- **RSA-PKCS1v15**: For Enterprise tier with hardware binding
- **Hardware Fingerprinting**: Enterprise licenses bound to specific machine
- **Offline Validation**: Supports air-gapped deployments with local signature verification

### Unified 4-Tier License Model
Replaced inconsistent tier models with unified system:

| Tier | Value | Pricing | Target |
|------|-------|---------|--------|
| Community | 0 | Free | Evaluation, personal projects |
| Developer | 1 | $29/mo | Individual developers, MVPs |
| Professional | 2 | $99/mo | Production apps, small teams |
| Enterprise | 3 | Custom | Large organizations, regulated industries |

### Sensitive File Cleanup
Removed from repository to prevent information disclosure

---

## New Features

### License Generation Tool
Created `cmd/licensegen/main.go` for administrators

### Enhanced Documentation
- Comprehensive README.md with full feature documentation
- SECURITY.md documenting security considerations
- MIGRATION.md for license migration guide
- VERSION file for build tracking

### Enhanced .gitignore
Added security patterns for license secrets and production configs

---

## What's Changed

### Added
- Cryptographic license validation (HMAC + RSA)
- License generation tool (cmd/licensegen)
- Hardware ID binding for Enterprise
- VERSION file for build tracking
- Enhanced .gitignore with security patterns
- SECURITY.md documentation
- MIGRATION.md for upgrade guidance

### Changed
- **BREAKING**: License format changed to base64(JSON).base64(SIGNATURE)
- **BREAKING**: Old license keys no longer valid
- Unified tier model across all packages

### Removed
- **SECURITY**: DEV_MODE environment variable bypass
- **SECURITY**: docs/PRICING.md (information disclosure)
- **SECURITY**: TODO.md (contained vulnerability details)
- **SECURITY**: SESSION_MEMORY.md (debug info)
- **SECURITY**: push.bat (dangerous git workflow)

---

## Migration Guide

### For Community Users
No action required - Community tier remains free with no license key needed

### For Developer/Professional Users
Must generate new HMAC-signed license key

### For Enterprise Users  
Must generate hardware-bound RSA-signed license key

See MIGRATION.md for detailed instructions.

---

## Verification

```bash
# Check version
./aegisgate --version

# Expected output:
# AegisGate v1.0.3
```

---

## Contact

- Security Issues: security@aegisgate.io
- License Support: support@aegisgate.io
- Sales: sales@aegisgate.io

---

**IMPORTANT**: This is a security-critical release. Previous versions contain known vulnerabilities and should be upgraded immediately.
