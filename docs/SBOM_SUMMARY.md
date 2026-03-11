# AegisGate MVP - SBOM Summary

## Software Bill of Materials

### MVP Dependencies (go.mod)

Module: github.com/aegisgatesecurity/aegisgate
Go Version: 1.23

### Local Packages (Replace Directives)

1. pkg/certificate - Certificate management
2. pkg/config - Environment-based configuration
3. pkg/proxy - Hardened reverse proxy
4. pkg/tls - TLS certificate manager

### External Dependencies

None - MVP uses only Go standard library:
- context
- crypto (rsa, tls, x509, rand)
- encoding (json, pem)
- fmt
- log/slog
- math/big
- net (http, url)
- os
- path/filepath
- sync (atomic, sync)
- syscall
- time

### Security Libraries (Stdlib)

- crypto/rand: Cryptographically secure random number generation
- crypto/rsa: RSA key generation and operations
- crypto/tls: TLS 1.2+ implementation
- crypto/x509: X.509 certificate parsing and validation

### Web Technologies (UI)

- HTML5
- CSS3 (Grid, Flexbox, Animations)
- Vanilla JavaScript (Fetch API)
- SVG Icons (embedded)

### No External Frameworks

MVP intentionally uses zero external dependencies:
- No frontend frameworks (React, Vue, Angular)
- No CSS frameworks (Bootstrap, Tailwind)
- No utility libraries (jQuery, Lodash)
- No external APIs or services

### Build Tools

- Go toolchain (go build)
- Standard Go testing (go test)
- No CI/CD dependencies for MVP

### SBOM Generation

Full SBOM available in sbom.json at project root.
Generated via: go mod vendor + analysis

### Vulnerability Scanning

Recommend scanning with:
- go mod tidy (dependency cleanup)
- go vet (static analysis)
- go test (runtime verification)

### License

All code: MIT License (see LICENSE file)
Standard library: BSD-style (Go project)

Generated: Phase 1 Week 13
Version: MVP v1.0.0

