const fs = require(
fs);
const readmeContent = `# AegisGate

[![CI Status](https://github.com/aegisgatesecurity/aegisgate/workflows/CI/badge.svg)](https://github.com/aegisgatesecurity/aegisgate/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/aegisgatesecurity/aegisgate)](https://goreportcard.com/report/github.com/aegisgatesecurity/aegisgate)
[![Go Version](https://img.shields.io/github/go-mod/go-version/aegisgatesecurity/aegisgate)](https://go.dev)
[![License](https://img.shields.io/github/license/aegisgatesecurity/aegisgate)](LICENSE)
[![Version](https://img.shields.io/github/v/release/aegisgatesecurity/aegisgate?include_prereleases)](https://github.com/aegisgatesecurity/aegisgate/releases)

**Enterprise AI Security Gateway - A comprehensive security framework for AI/LLM API protection**

---

## Overview

AegisGate is a high-performance, production-ready security gateway written in **Go 1.24** that provides comprehensive protection for AI/LLM APIs. It acts as a secure proxy between client applications and upstream AI services (OpenAI, Anthropic, Cohere, Azure OpenAI, etc.), providing security, compliance, and observability at the edge.

### Why AegisGate?

- **Zero-Trust Architecture**: Every request is validated, scanned, and logged
- **AI-First Design**: Purpose-built for LLM APIs with prompt injection detection, PII filtering, and content classification
- **Enterprise Ready**: mTLS, OAuth2/OIDC/SAML SSO, SOC 2, HIPAA, PCI-DSS compliance frameworks
- **Observable**: Prometheus metrics, OpenTelemetry tracing, Kubernetes-ready health probes
- **Production Hardened**: 99.9% uptime design with graceful degradation

---

## Key Features

### Security

| Feature | Description |
|---------|-------------|
| **MITM Proxy** | Full HTTP/HTTPS interception with certificate handling |
| **mTLS** | Mutual TLS authentication between clients and upstream |
| **Content Scanning** | Real-time threat detection (SQL injection, XSS, prompt injection) |
| **PII Detection** | Automatic detection and filtering of sensitive data |
| **Secrets Management** | Secure storage and retrieval of credentials |
| **PKI Attestation** | Certificate chain validation and backdoor prevention |
| **Security Middleware** | CSRF protection, XSS prevention, panic recovery, audit logging |

### Authentication and Access Control

| Feature | Description |
|---------|-------------|
| **OAuth2/OIDC** | Industry-standard authorization framework |
| **SAML 2.0** | Enterprise SSO integration |
| **JWT Validation** | Token verification and claims inspection |
| **Trust Domains** | Multi-tenant isolation policies |
| **Rate Limiting** | Per-client, per-endpoint rate limits |

### Observability

| Feature | Description |
|---------|-------------|
| **Prometheus Metrics** | /metrics endpoint for Prometheus scraping |
| **Health Endpoints** | /health/live, /health/ready, /health, /health/components |
| **Real-time Dashboard** | Web UI for monitoring and configuration |
| **SIEM Integration** | Splunk, ELK Stack, Azure Sentinel connectors |
| **Webhooks** | Real-time event notifications |
| **WebSocket SSE** | Streaming event updates |

### Compliance

| Tier | Frameworks |
|------|------------|
| **Community (Free)** | MITRE ATLAS, OWASP AI Top 10, GDPR |
| **Enterprise** | + NIST AI RMF, ISO/IEC 42001 |
| **Premium** | + SOC 2, HIPAA, PCI DSS |

---

## Quick Start

### Installation

\`\`\`bash
# Install via Go
go install github.com/aegisgatesecurity/aegisgate/cmd/aegisgate@latest

# Or clone and build
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
go build -o aegisgate.exe ./cmd/aegisgate
\`\`\`

### Docker

\`\`\`bash
# Build Docker image
docker build -t aegisgate:latest .

# Run with default settings
docker run -p 8080:8080 aegisgate:latest
\`\`\`

### Run in 60 Seconds

\`\`\`bash
# Start the security gateway (default: localhost:8080)
./aegisgate.exe

# Test the health endpoints
curl http://localhost:8080/health
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready

# Get Prometheus metrics
curl http://localhost:8080/metrics

# Enable MITM mode
export AEGISGATE_MITM_ENABLED=true
./aegisgate.exe
\`\`\`

---

## Architecture

AegisGate consists of **30+ modular packages** organized into cohesive layers:

\`\`\`
Client Applications
        |
   AegisGate Gateway
 +---------------+ +---------------+ +---------------+
 |    Proxy      | |     TLS/      | |      SSO      |
 |    Layer      | |     mTLS      | |     Auth      |
 +---------------+ +---------------+ +---------------+
 +---------------+ +---------------+ +---------------+
 |   Scanner     | |  Compliance   | |   Dashboard   |
 |  (Threats)    | |  (14 Modules) | |     (UI)      |
 +---------------+ +---------------+ +---------------+
 +---------------+ +---------------+ +---------------+
 |     SIEM      | |   Threat      | |  Signature    |
 | Integration   | |    Intel      | |   Verify      |
 +---------------+ +---------------+ +---------------+
        |
   Upstream Services
 (OpenAI, Anthropic, Cohere, Custom LLMs)
\`\`\`

---

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| AEGISGATE_LISTEN_ADDR | :8080 | Server listen address |
| AEGISGATE_UPSTREAM_URL | - | Upstream AI API URL |
| AEGISGATE_API_KEY | - | API key for upstream |
| AEGISGATE_MITM_ENABLED | false | Enable MITM proxy |
| AEGISGATE_MTLS_MODE | disabled | mTLS mode (disabled/optional/required) |
| AEGISGATE_LOG_LEVEL | info | Log level (debug/info/warn/error) |
| AEGISGATE_METRICS_ENABLED | true | Enable Prometheus metrics |
| AEGISGATE_CSRF_ENABLED | true | Enable CSRF protection |
| AEGISGATE_RATE_LIMIT | 100 | Requests per minute |

### Health Endpoints

| Endpoint | Purpose | Use Case |
|----------|---------|----------|
| /health/live | Liveness probe | Kubernetes liveness |
| /health/ready | Readiness probe | Kubernetes readiness |
| /health | Enhanced health | Full system status |
| /health/components | Component status | Detailed diagnostics |
| /metrics | Prometheus metrics | Monitoring scraping |

---

## Package Structure

\`\`\`
aegisgate/
├── cmd/
│   └── aegisgate/              # Main entry point
├── pkg/
│   ├── auth/                 # Authentication handlers
│   ├── certificate/          # Certificate management
│   ├── compliance/           # Compliance frameworks (14 modules)
│   ├── config/               # Configuration management
│   ├── dashboard/            # Web UI and API server
│   │   └── observability.go # Prometheus and health endpoints
│   ├── immutable-config/    # Immutable configuration
│   ├── metrics/              # Internal metrics collector
│   ├── ml/                   # Machine learning (threat detection)
│   ├── opsec/                # Operational security
│   ├── proxy/                # HTTP/HTTPS proxy
│   │   └── mitm.go          # MITM proxy implementation
│   ├── reporting/           # Report generation
│   ├── scanner/             # Content/Threat scanning
│   ├── secrets/              # Secrets management
│   ├── security/             # Security middleware
│   ├── siem/                 # SIEM integrations
│   ├── sso/                  # SSO (OIDC, SAML)
│   ├── threatintel/          # Threat intelligence (STIX/TAXII)
│   ├── tls/                  # TLS/mTLS handling
│   ├── trustdomain/          # Trust domain policies webhook/              #
│   ├── Webhook notifications
│   └── websocket/            # WebSocket support
├── tests/
│   └── integration/          # Integration tests
└── docs/                     # Documentation
\`\`\`

---

## Testing

\`\`\`bash
# Run all tests
go test ./...

# Run security tests
go test ./pkg/security/... -v

# Run integration tests
go test ./tests/integration/... -v

# Run benchmarks
go test -bench=. -benchmem ./pkg/security/...

# Run with coverage
go test -cover ./...
\`\`\`

---

## Performance

### Benchmark Results

| Operation | Latency | Throughput |
|-----------|---------|------------|
| Baseline (no middleware) | 460 ns/op | 2.7M ops/sec |
| Security Headers | 6,481 ns/op | 154K ops/sec |
| CSRF Protection | 24,430 ns/op | 41K ops/sec |
| Full Security Suite | ~28,456 ns/op | 35K ops/sec |

### System Metrics

- **Memory**: ~50MB baseline
- **Goroutines**: ~20-50 under normal load
- **Latency overhead**: less than 5ms for security scanning

---

## Roadmap

### v0.29.0 (Next)
- Advanced ML-based threat detection
- Enhanced SIEM integrations
- Kubernetes operator
- Helm chart repository

### v1.0.0 (GA)
- Production stability guarantees
- Enterprise SLA support
- Multi-region deployment guides
- Advanced rate limiting

---

## Security

See SECURITY.md for vulnerability disclosure and security features.

### Security Features

- **Panic Recovery**: Graceful error handling without crashes
- **CSRF Protection**: Token-based request forgery prevention
- **XSS Prevention**: Content sanitization and CSP
- **Audit Logging**: Comprehensive request/response logging
- **mTLS**: Mutual TLS authentication
- **PKI Attestation**: Certificate chain validation

---

## License

Released under the **MIT License** - see LICENSE

---

## Contributing

Contributions are welcome. Please read our contributing guidelines before submitting PRs.

### Development Setup

\`\`\`bash
# Clone the repository
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o aegisgate.exe ./cmd/aegisgate
\`\`\`

---

## Support

- **Issues**: https://github.com/aegisgatesecurity/aegisgate/issues
- **Discussions**: https://github.com/aegisgatesecurity/aegisgate/discussions
- **Documentation**: https://pkg.go.dev/github.com/aegisgatesecurity/aegisgate

---

Built with care by the AegisGate Team
`;

fs.writeFileSync(aegisgate/README.md, readmeContent);
console.log(Wrote
README.md
-
 + readmeContent.length + 
characters);

