# AegisGate 🔐

[![Go Version](https://img.shields.io/github/go-mod/go-version/aegisgatesecurity/aegisgate)](https://github.com/aegisgatesecurity/aegisgate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![License: Commercial](https://img.shields.io/badge/License-Commercial-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/aegisgatesecurity/aegisgate/build.yml)](https://github.com/aegisgatesecurity/aegisgate/actions)
[![Latest Release](https://img.shields.io/github/v/release/aegisgatesecurity/aegisgate)](https://github.com/aegisgatesecurity/aegisgate/releases/latest)

## Overview

AegisGate is a comprehensive enterprise-grade security platform written in Go, providing advanced TLS/SSL certificate management, PKI attestation, threat intelligence integration, and enterprise authentication capabilities. The project delivers a modular architecture designed for scalability, security, and extensibility.

## Version

**Current Release: v0.21.1** (i18n Enhancement Release)

## Features

### Core Security Features

- **TLS/SSL Certificate Management**
  - Automated certificate lifecycle management
  - Certificate issuance and renewal
  - Custom TLS implementations
  - Certificate validation and verification

- **PKI Attestation**
  - Public Key Infrastructure support
  - Certificate authority integration
  - Hardware security module (HSM) compatibility

- **Threat Intelligence Integration**
  - STIX 2.1 and TAXII 2.1 support
  - Real-time threat feeds
  - Indicator of Compromise (IOC) management
  - MITRE ATT&CK framework mapping

- **Enterprise Authentication**
  - SAML 2.0 Single Sign-On (SSO)
  - OpenID Connect (OIDC) provider
  - Multi-factor authentication support
  - Role-based access control (RBAC)

### Advanced Security Features

- **Digital Signature Verification**
  - RSA, ECDSA, and Ed25519 algorithms
  - Key management and rotation
  - Verification statistics and reporting

- **Hash Chain Validation**
  - Merkle tree integration
  - SHA256/SHA512 hash algorithms
  - Tamper detection mechanisms

- **Feed-level Security**
  - Feed-specific trust domains
  - Resource sandboxing with quotas
  - Policy-based governance

- **Webhook Alerting**
  - Real-time event notifications
  - Multiple authentication methods
  - Configurable retry with exponential backoff

### Internationalization (i18n)

Full internationalization support with 12 locales:

| Code | Language | Plural Rules |
|------|----------|--------------|
| ar | Arabic | Plural forms (6) |
| de | German | 1 form |
| en | English | 1 form |
| es | Spanish | 1 form |
| fr | French | 2 forms |
| he | Hebrew | 2 forms |
| hi | Hindi | 1 form |
| ja | Japanese | 1 form |
| ko | Korean | 1 form |
| pt | Portuguese | 2 forms |
| ru | Russian | 4 forms |
| zh | Chinese | 1 form |

## Architecture

### Package Structure

```
aegisgate/
├── cmd/                    # Application entry points
│   ├── aegisgate/           # Main application
│   └── debug/             # Debug utilities
├── pkg/                   # Core packages
│   ├── auth/             # Authentication & authorization
│   ├── certificate/      # TLS/SSL certificate handling
│   ├── compliance/       # Compliance management
│   ├── config/           # Configuration management
│   ├── core/             # Core functionality
│   ├── dashboard/        # Web dashboard
│   ├── hash_chain/       # Merkle tree & hash validation
│   ├── i18n/             # Internationalization
│   ├── immutable-config/ # Immutable configuration
│   ├── metrics/          # Metrics collection
│   ├── ml/               # Machine learning
│   ├── pkiattest/        # PKI attestation
│   ├── proxy/            # Proxy functionality
│   ├── reporting/        # Reporting engine
│   ├── sandbox/          # Feed-level sandboxing
│   ├── scanner/          # Security scanning
│   ├── security/         # Security utilities
│   ├── siem/             # SIEM integration
│   ├── signature_verification/ # Digital signatures
│   ├── sso/              # SSO (SAML/OIDC)
│   ├── threatintel/      # Threat intelligence
│   ├── tls/              # TLS implementations
│   ├── trustdomain/      # Trust domains
│   ├── webhook/          # Webhook alerting
│   └── websocket/        # WebSocket support
├── ui/                    # User interface
├── docs/                  # Documentation
├── configs/               # Configuration files
├── deploy/                # Deployment scripts
└── tests/                 # Test utilities
```

### Technology Stack

- **Language**: Go 1.24+
- **Web Framework**: Custom HTTP handlers with WebSocket support
- **Database**: Extensible storage backends
- **Security**: Go's standard library + golang.org/x/net

## Installation

### Prerequisites

- Go 1.24 or higher
- Git

### Build from Source

```bash
# Clone the repository
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# Build the binary
go build -o aegisgate ./cmd/aegisgate

# Run the application
./aegisgate
```

### Using Makefile

```bash
make build    # Build the binary
make test     # Run tests
make lint     # Run linter
make clean    # Clean artifacts
make all      # Build, test, and lint
```

### Docker

```bash
# Build Docker image
docker build -t aegisgate:latest .

# Run container
docker run -d -p 8080:8080 aegisgate:latest
```

## Configuration

AegisGate uses a flexible configuration system with support for:

- Environment variables
- Configuration files (YAML, JSON)
- Command-line flags

### Environment Variables

```bash
AEGISGATE_PORT=8080
AEGISGATE_LOG_LEVEL=info
AEGISGATE_CONFIG_PATH=/etc/aegisgate/config.yaml
```

### Configuration File Example

```yaml
server:
  port: 8080
  host: "0.0.0.0"

security:
  tls:
    enabled: true
    cert: "/path/to/cert.pem"
    key: "/path/to/key.pem"

auth:
  sso:
    enabled: true
    provider: "saml"

i18n:
  default_locale: "en"
  supported_locales:
    - "en"
    - "es"
    - "fr"
    - "de"
    - "zh"
    - "ja"
    - "ko"
    - "pt"
    - "ru"
    - "ar"
    - "he"
    - "hi"

threat_intel:
  enabled: true
  stix_version: "2.1"
  taxii_server: "https://taxii.example.com"
```

## Usage

### Basic Usage

```bash
# Start the server
./aegisgate serve

# Run with custom config
./aegisgate serve --config /path/to/config.yaml

# Check version
./aegisgate --version

# Show help
./aegisgate --help
```

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/health | GET | Health check |
| /api/v1/certificates | GET | List certificates |
| /api/v1/certificates | POST | Create certificate |
| /api/v1/auth/login | POST | Authentication |
| /api/v1/sso/saml | POST | SAML SSO |
| /api/v1/sso/oidc | POST | OIDC SSO |
| /api/v1/threatintel | GET | Threat intelligence |
| /api/v1/webhooks | GET/POST | Webhook management |
| /api/v1/metrics | GET | System metrics |

## Development

### Running Tests

```bash
go test ./...
go test -v ./pkg/...
```

### Code Quality

```bash
go fmt ./...
go vet ./...
golangci-lint run
```

### Adding New Features

1. Fork the repository
2. Create a feature branch
3. Implement changes with tests
4. Submit a pull request

## Internationalization

### Adding a New Locale

1. Create locale file: `pkg/i18n/locales/<locale>.json`
2. Define plural rules in `pkg/i18n/plural.go`
3. Register locale in `pkg/i18n/locales.go`
4. Update supported locales in configuration

### Supported Locales

Currently supported locales include: Arabic (ar), Chinese (zh), English (en), French (fr), German (de), Hebrew (he), Hindi (hi), Japanese (ja), Korean (ko), Portuguese (pt), Russian (ru), and Spanish (es).

## Security

### Reporting Security Issues

If you discover a security vulnerability, please report it via GitHub Security Advisories or contact the maintainers directly.

### Security Features

- TLS/SSL certificate pinning
- HMAC signature verification
- OAuth2 support for API security
- Session management with Redis backend option

## Documentation

- [API Documentation](docs/api.md)
- [Configuration Guide](docs/configuration.md)
- [Security Guide](docs/security.md)
- [Deployment Guide](docs/deployment.md)
- [Changelog](CHANGELOG.md)

## Roadmap

### Upcoming Features

- v0.22.0: Integration Testing and Performance Validation
- Enhanced SIEM integration
- Advanced machine learning for threat detection
- Expanded cloud provider support

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting pull requests.

### Development Setup

```bash
# Install dependencies
go mod download

# Run development server
go run ./cmd/aegisgate
```

## License

Dual licensed under:

- **MIT License** - See [LICENSE](LICENSE) for details
- **Commercial License** - Contact sales for enterprise licensing

## Support

- **GitHub Issues**: Report bugs and request features
- **Documentation**: Check the [docs](docs/) directory
- **Discussions**: Use GitHub Discussions for questions

## Version History

| Version | Date | Highlights |
|---------|------|------------|
| v0.21.1 | 2026-02-25 | i18n Enhancement (12 locales) |
| v0.21.0 | 2026-02-25 | Phase 2 Architecture Complete |
| v0.20.0 | 2026-02-24 | Feed-level Sandboxing |
| v0.19.0 | 2026-02-23 | Threat Intelligence Platform |
| v0.18.0 | 2026-02-23 | SSO & Webhooks |

See [CHANGELOG.md](CHANGELOG.md) for complete version history.

## Acknowledgments

- Go community for excellent libraries
- Contributors and maintainers
- Security research community

---

*AegisGate - Enterprise Security Platform*
