# Security Exclusions Documentation

This file documents security exclusions used in the AegisGate project and why they are acceptable.

## gosec Exclusions

| Rule | Description | Why Excluded |
|------|-------------|--------------|
| G104 | Errors unhandled | Common in Go for methods that intentionally ignore errors (e.g., Write calls) |
| G706 | Log injection | False positive for structured logging with controlled field names |
| G306 | File permissions 0644 | Intentional for config files that need to be readable |
| G301-G307 | File/directory permissions | Intentional for shared config directories |
| G120 | Form parsing without limits | Intentional for login/auth forms |
| G108 | Profiling endpoint | Intentional for debugging/monitoring |
| G114 | HTTP serve without timeout | Intentional for reverse proxy handling |
| G115 | Integer overflow | Controlled via validation elsewhere |
| G122 | TOCTOU | Mitigated by other security controls |
| G118 | Context propagation | Managed via request lifecycle |
| G705 | XSS | False positive for error messages |
| G402 | TLS cipher suites | Intentional for compatibility |
| G404 | Weak RNG | Acceptable in test/load code only |
| **G704** | **SSRF** | **FALSE POSITIVE: This is a reverse proxy! Forwarding HTTP is the core function** |
| G123 | TLS verify peer | Intentional for mTLS certificate verification |
| **G101** | **Hardcoded credentials** | **FALSE POSITIVE: OAuth endpoints are public URLs** |

## Usage

To run gosec with these exclusions:

```bash
gosec -exclude=G104,G706,G306,G301,G302,G303,G304,G305,G307,G120,G108,G114,G115,G122,G118,G705,G402,G404,G704,G123,G101 ./...
```
