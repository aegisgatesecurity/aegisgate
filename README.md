# AegisGate 🔐

<div align="center">

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/aegisgatesecurity/aegisgate)](go.mod)
[![Release](https://img.shields.io/github/v/release/aegisgatesecurity/aegisgate)](https://github.com/aegisgatesecurity/aegisgate/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/aegisgate/aegisgate.svg)](https://hub.docker.com/r/aegisgate/aegisgate)
[![Tests](https://img.shields.io/github/actions/workflow/status/aegisgatesecurity/aegisgate/test.yml?label=tests)](https://github.com/aegisgatesecurity/aegisgate/actions)
[![Security](https://img.shields.io/github/actions/workflow/status/aegisgatesecurity/aegisgate/security.yml?label=security)](https://github.com/aegisgatesecurity/aegisgate/actions)
[![Discord](https://img.shields.io/discord/123456789?label=Discord)](https://discord.gg/aegisgate)
[![Twitter](https://img.shields.io/twitter/follow/aegisgatesec?style=social)](https://twitter.com/aegisgatesec)

</div>

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Configuration](#configuration)
- [Deployment](#deployment)
- [Security](#security)
- [Compliance](#compliance)
- [API Reference](#api-reference)
- [SDK & Libraries](#sdk--libraries)
- [Contributing](#contributing)
- [Roadmap](#roadmap)
- [Support](#support)
- [License](#license)

---

## Overview

**AegisGate** is an enterprise-grade security platform providing comprehensive API gateway capabilities with built-in AI/ML-powered threat detection, compliance automation, and zero-trust architecture.

Built in Go with a modular architecture, AegisGate delivers high-performance proxy services with support for HTTP/1.1, HTTP/2, HTTP/3, gRPC, and WebSocket traffic. It integrates seamlessly with Kubernetes, Docker, and traditional infrastructure.

### Key Capabilities

| Capability | Description |
|------------|-------------|
| **AI-Powered Security** | ML-based anomaly detection and threat intelligence |
| **Zero-Trust Architecture** | mTLS, PKI attestation, and fine-grained RBAC |
| **Compliance Automation** | Built-in support for SOC 2, HIPAA, PCI-DSS, GDPR, OWASP |
| **Multi-Protocol Support** | HTTP/1.1, HTTP/2, HTTP/3, gRPC, WebSocket |
| **Enterprise SSO** | SAML 2.0, OpenID Connect, OAuth 2.0 |
| **Multi-Tenancy** | Isolated tenant environments with resource quotas |

---

## Features

### Core Features

#### 🔄 **API Gateway**
- High-performance reverse proxy with connection pooling
- Protocol translation (HTTP ↔ gRPC, HTTP ↔ WebSocket)
- Request/response transformation and routing
- Rate limiting with token bucket algorithm
- Circuit breaker and retry policies

#### 🛡️ **Security**
- TLS 1.3 termination with perfect forward secrecy
- Mutual TLS (mTLS) for service-to-service authentication
- Hardware-backed PKI attestation
- Automatic secret rotation
- CSRF/XSS protection headers
- Rate limiting and DDoS mitigation

#### 🤖 **AI/ML Integration**
- Anomaly detection for traffic patterns
- Threat intelligence feed integration (STIX/TAXII)
- Behavioral analysis for suspicious activity
- Predictive capacity planning

#### 📋 **Compliance**
- Automated compliance framework mapping
- Audit logging with immutable records
- Evidence collection for certifications
- Policy enforcement engine

#### 👥 **Identity & Access**
- Role-based access control (RBAC)
- Enterprise SSO (SAML, OIDC)
- API key management with rotation
- Session management and JWT validation

#### 📊 **Observability**
- Prometheus metrics with Grafana dashboards
- Structured logging (JSON format)
- Distributed tracing support
- Real-time health monitoring

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                            AegisGate                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │   Gateway    │  │   Proxy      │  │   Admin      │              │
│  │   Handler    │  │   Engine     │  │   API        │              │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘              │
│         │                 │                 │                      │
│  ┌──────┴─────────────────┴─────────────────┴───────┐              │
│  │              Security Layer                        │              │
│  │  ┌─────────┐  ┌──────────┐  ┌────────────────┐  │              │
│  │  │ TLS/MTLS│  │ Auth/N   │  │ Threat         │  │              │
│  │  │         │  │ RBAC     │  │ Detection      │  │              │
│  │  └─────────┘  └──────────┘  └────────────────┘  │              │
│  └──────────────────────────────────────────────────┘              │
│         │                 │                 │                      │
│  ┌──────┴─────────────────┴─────────────────┴───────┐              │
│  │              Core Modules                        │              │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌───────┐│              │
│  │  │ ML/D    │ │Complian-│ │Webhook  │ │SIEM   ││              │
│  │  │         │ │ce       │ │Manager  │ │Client ││              │
│  │  └─────────┘ └─────────┘ └─────────┘ └───────┘│              │
│  └──────────────────────────────────────────────────┘              │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### Component Overview

| Component | Description | Location |
|-----------|-------------|----------|
| **cmd/aegisgate** | Main entry point | `cmd/aegisgate/main.go` |
| **pkg/proxy** | Core proxy engine | `pkg/proxy/` |
| **pkg/auth** | Authentication & sessions | `pkg/auth/` |
| **pkg/ml** | ML anomaly detection | `pkg/ml/` |
| **pkg/compliance** | Compliance frameworks | `pkg/compliance/` |
| **pkg/threatintel** | Threat intelligence | `pkg/threatintel/` |
| **pkg/webhook** | Webhook management | `pkg/webhook/` |
| **pkg/siem** | SIEM integrations | `pkg/siem/` |

---

## Quick Start

### Using Docker

```bash
# Pull and run the latest release
docker pull aegisgate/aegisgate:latest
docker run -d \
  --name aegisgate \
  -p 8443:8443 \
  -p 8080:8080 \
  -v $(pwd)/config:/config \
  -e AEGISGATE_CONFIG=/config/aegisgate.yml \
  aegisgate/aegisgate:latest
```

### Using Docker Compose

```bash
# Clone the repository
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# Start the development environment
docker-compose -f deploy/docker/docker-compose.yml up -d
```

### Using Homebrew

```bash
# Install on macOS/Linux
brew tap aegisgatesecurity/aegisgate
brew install aegisgate

# Start the server
aegisgate serve
```

### From Source

```bash
# Clone the repository
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# Build the binary
make build

# Run the server
./bin/aegisgate serve
```

---

## Installation

### Prerequisites

| Requirement | Version | Notes |
|-------------|---------|-------|
| Go | 1.21+ | For building from source |
| Docker | Latest | For containerized deployment |
| Kubernetes | 1.24+ | For K8s deployment (optional) |
| PostgreSQL | 14+ | For persistent storage (optional) |
| Redis | 7+ | For caching/sessions (optional) |

### Binary Installation

#### Linux/macOS

```bash
# Download the latest release
curl -LO https://github.com/aegisgatesecurity/aegisgate/releases/latest/download/aegisgate-linux-amd64.tar.gz
tar xzf aegisgate-linux-amd64.tar.gz
sudo mv aegisgate /usr/local/bin/
```

#### Windows

```powershell
# Using PowerShell
Invoke-WebRequest -Uri "https://github.com/aegisgatesecurity/aegisgate/releases/latest/download/aegisgate-windows-amd64.zip" -OutFile aegisgate.zip
Expand-Archive aegisgate.zip -DestinationPath C:\Program Files\AegisGate
$env:PATH += ";C:\Program Files\AegisGate"
```

---

## Configuration

### Configuration File

AegisGate uses YAML for configuration. Copy the example config:

```bash
cp config/aegisgate.yml.example config/aegisgate.yml
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `AEGISGATE_CONFIG` | Path to config file | `./config/aegisgate.yml` |
| `AEGISGATE_LOG_LEVEL` | Log level (debug, info, warn, error) | `info` |
| `AEGISGATE_HOST` | Bind address | `0.0.0.0` |
| `AEGISGATE_PORT` | HTTP port | `8080` |
| `AEGISGATE_TLS_PORT` | HTTPS port | `8443` |
| `AEGISGATE_TLS_ENABLED` | Enable TLS | `false` |
| `AEGISGATE_DB_HOST` | PostgreSQL host | `localhost` |
| `AEGISGATE_DB_PORT` | PostgreSQL port | `5432` |
| `AEGISGATE_REDIS_HOST` | Redis host | `localhost` |
| `AEGISGATE_REDIS_PORT` | Redis port | `6379` |

### Tier-Specific Configuration

AegisGate offers tier-specific configuration presets:

| Tier | Config File | Features |
|------|-------------|----------|
| Community | `config/community.env` | Free, up to 100 req/min, 3 AI providers |
| Developer | `config/developer.env` | $29/mo, 1K req/min, 8 AI providers |
| Professional | `config/professional.env` | $99/mo, unlimited, SSO, multi-tenancy |
| Enterprise | `config/enterprise.env` | Custom, dedicated support, full compliance |

---

## Deployment

### Docker Deployment

```bash
# Single container
docker run -d \
  --name aegisgate \
  -p 8443:8443 \
  -v ./config:/config \
  aegisgate/aegisgate:latest

# With docker-compose
docker-compose -f deploy/docker/docker-compose.yml up -d
```

### Kubernetes Deployment

```bash
# Deploy using Helm
helm repo add aegisgate https://aegisgatesecurity.github.io/helm-charts
helm install aegisgate aegisgate/aegisgate

# Or deploy manifests directly
kubectl apply -f deploy/k8s/
```

### Production Considerations

1. **TLS Termination**: Enable TLS in production
2. **Database**: Use PostgreSQL for persistent data
3. **Caching**: Use Redis for sessions and caching
4. **Monitoring**: Configure Prometheus and Grafana
5. **Secrets**: Use Kubernetes secrets or Vault
6. **Health Checks**: Configure liveness/readiness probes
7. **Scaling**: Use HPA for horizontal pod autoscaling

---

## Security

### Security Features

| Feature | Description |
|---------|-------------|
| **TLS 1.3** | Latest TLS with perfect forward secrecy |
| **mTLS** | Mutual TLS for service authentication |
| **PKI Attestation** | Hardware-backed identity verification |
| **Secret Rotation** | Automated credential management |
| **Audit Logging** | Immutable security event log |
| **RBAC** | Fine-grained role-based access |
| **Rate Limiting** | Configurable per-client limits |

### Security Best Practices

```bash
# Enable TLS in production
AEGISGATE_TLS_ENABLED=true
AEGISGATE_TLS_PORT=8443

# Enable audit logging
AEGISGATE_AUDIT_LOG=true

# Restrict admin access
AEGISGATE_ADMIN_IPS=10.0.0.0/8

# Enable rate limiting
AEGISGATE_RATE_LIMIT=true
```

### Reporting Security Issues

See [SECURITY.md](SECURITY.md) for our security vulnerability reporting process.

---

## Compliance

AegisGate helps meet compliance requirements for:

| Framework | Features |
|-----------|----------|
| **SOC 2** | Audit logging, access control, encryption |
| **HIPAA** | Data encryption, audit trails, access control |
| **PCI-DSS** | TLS, secure defaults, logging |
| **GDPR** | Data encryption, access control, logging |
| **OWASP** | Built-in OWASP Top 10 protection |

### Compliance Automation

```bash
# Enable compliance framework
AEGISGATE_COMPLIANCE_FRAMEWORK=soc2

# Generate compliance report
aegisgate compliance report --framework=soc2
```

---

## API Reference

### REST API

AegisGate provides a comprehensive REST API:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/health` | GET | Health check |
| `/api/v1/config` | GET/PUT | Configuration management |
| `/api/v1/auth/login` | POST | Authenticate |
| `/api/v1/proxy/*` | * | Proxy requests |
| `/api/v1/metrics` | GET | Prometheus metrics |
| `/api/v1/compliance/reports` | GET | Compliance reports |

### GraphQL API

AegisGate also supports GraphQL:

```graphql
type Query {
  health: HealthStatus!
  metrics(since: Time): Metrics!
  config: Config!
}

type Mutation {
  updateConfig(input: ConfigInput!): Config!
}
```

### gRPC API

Service definitions are in `pkg/grpc/`:

```proto
service AegisGateService {
  rpc HealthCheck(HealthRequest) returns (HealthResponse);
  rpc Config(ConfigRequest) returns (ConfigResponse);
  rpc Proxy(stream ProxyRequest) returns (stream ProxyResponse);
}
```

---

## SDK & Libraries

### Go SDK

```go
import "github.com/aegisgatesecurity/aegisgate/sdk/go"

client := aegisgate.NewClient("https://aegisgate.example.com", "api-key")
result, err := client.Health()
```

### API Wrappers

See [docs/API_WRAPPERS.md](docs/API_WRAPPERS.md) for community-maintained wrappers in Python, JavaScript, Ruby, and other languages.

---

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md).

### Development Setup

```bash
# Clone and setup
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
make build
```

### Commit Message Format

```
<type>: <short description>

<long description if needed>

<footer>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

---

## Roadmap

See our [public roadmap](https://github.com/aegisgatesecurity/aegisgate/projects) for upcoming features.

### Upcoming Features

| Feature | Target Version | Status |
|---------|---------------|--------|
| Terraform Provider | v1.1.0 | Planned |
| Additional SIEM Integrations | v1.1.0 | Planned |
| GraphQL Subscriptions | v1.1.0 | In Progress |
| ISO 42001 Compliance | v1.2.0 | Planned |
| FedRAMP Support | v1.2.0 | Planned |

---

## Support

### Resources

| Resource | Link |
|----------|------|
| Documentation | https://docs.aegisgate.com |
| Discord | https://discord.gg/aegisgate |
| GitHub Issues | https://github.com/aegisgatesecurity/aegisgate/issues |
| Twitter | @aegisgatesec |

### Tier Support

| Tier | Support Level | Response Time |
|------|--------------|---------------|
| Community | GitHub Issues | Best effort |
| Developer | Email | 24 hours |
| Professional | Priority | 4 hours |
| Enterprise | Dedicated | 1 hour |

---

## License

AegisGate is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.

---

## Related Projects

- [aegisgate-helm-charts](https://github.com/aegisgatesecurity/helm-charts) - Kubernetes Helm charts
- [aegisgate-sdk-go](https://github.com/aegisgatesecurity/sdk-go) - Go SDK
- [aegisgate-examples](https://github.com/aegisgatesecurity/examples) - Example implementations

---

<div align="center">

**Built with ❤️ by [AegisGate Security](https://aegisgatesecurity.io)**

[![GitHub stars](https://img.shields.io/github/stars/aegisgatesecurity/aegisgate?style=social)](https://github.com/aegisgatesecurity/aegisgate/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/aegisgatesecurity/aegisgate?style=social)](https://github.com/aegisgatesecurity/aegisgate/network)

</div>