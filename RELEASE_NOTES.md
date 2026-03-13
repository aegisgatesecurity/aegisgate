# AegisGate v1.0.3 Release Notes

**Release Date**: March 13, 2026
**Version**: v1.0.3

---

## Security Critical - Immediate Action Required

This release addresses CRITICAL security vulnerabilities.

### Breaking Changes
Old license keys are NOT VALID. Generate new cryptographically signed licenses.

---

## Security Fixes

### Critical: DEV_MODE Bypass Removed
- Removed AEGISGATE_DEV_MODE environment variable bypass
- All licenses now require cryptographic validation

### Cryptographic License Validation
- HMAC-SHA256 for Developer and Professional tiers
- RSA-PKCS1v15 for Enterprise with hardware binding
- Offline/air-gapped validation support

### Unified 4-Tier Model
- Community (0): Free
- Developer (1): $29/mo  
- Professional (2): $99/mo
- Enterprise (3): Custom

---

## Features

### Added
- License generation tool (cmd/licensegen)
- Hardware ID binding for Enterprise
- Comprehensive documentation updates
- VERSION file for build tracking
- SECURITY.md and MIGRATION.md guides

### Removed (Security)
- docs/PRICING.md
- TODO.md  
- SESSION_MEMORY.md
- push.bat
- DEV_MODE bypass

---

## Migration

See [MIGRATION.md](MIGRATION.md) for detailed instructions.

**Community Users**: No action required
**Developer/Professional**: Generate new HMAC-signed key
**Enterprise**: Generate hardware-bound RSA key

---

## Contact

- Security: security@aegisgate.io
- Support: support@aegisgate.io
- Sales: sales@aegisgate.io
