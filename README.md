# AegisGate - Enterprise AI API Security Platform

<div align="center">

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://golang.org/)
[![Version](https://img.shields.io/badge/version-v1.0.4-green.svg)](https://github.com/aegisgatesecurity/aegisgate/releases)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://docker.com/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-326CE5.svg)](https://kubernetes.io/)

**Enterprise-grade security platform for AI API gateways**

[Features](#features) • [Quick Start](#quick-start) • [Tiers](#tiers-and-licensing) • [Security](#security) • [Documentation](#documentation) • [Contributing](#contributing)

</div>

---

## Why AegisGate?

In an era where AI powers critical business operations, security isn't optional—it's foundational. AegisGate provides enterprise-grade protection for your AI infrastructure without compromising performance or usability.

### The Business Case

| Challenge | AegisGate Solution |
|-----------|-------------------|
| **Prompt Injection Attacks** | Multi-layer detection blocks malicious prompts before they reach your AI |
| **Data Leakage / PII Exposure** | Automatic PII redaction keeps sensitive data out of AI logs |
| **Shadow AI Usage** | Centralized proxy gives you visibility into all AI traffic |
| **Compliance Overhead** | Out-of-the-box support for SOC2, HIPAA, GDPR, PCI-DSS, ISO 27001 |
| **Cost Anomalies** | ML-powered detection identifies unusual spending patterns |
| **Vendor Lock-in** | Unified API supports OpenAI, Anthropic, Azure, AWS, Cohere—switch providers without code changes |

### Key Differentiators

- **Zero-Latency Security**: Inline processing adds <5ms to request latency
- **Transparent Deployment**: No code changes required—deploy as a drop-in proxy
- **Enterprise-Ready**: SSO/SAML, RBAC, audit logging, mTLS support
- **Cost-Effective**: Up to 70% cheaper than competitors like Palo Alto AI Security or Microsoft Copilot Security

### Who Uses AegisGate?

- **Financial Services**: Protect LLM-powered trading algorithms and customer service bots
- **Healthcare**: Ensure HIPAA compliance for AI-assisted diagnosis tools
- **Enterprise**: Centralize AI governance across hundreds of internal applications
- **SaaS Providers**: Add security layer to AI-powered products without development overhead

---

## Overview

AegisGate is a comprehensive, enterprise-grade security platform designed specifically for AI API gateways. It provides real-time threat detection, compliance monitoring, secure proxying, and advanced ML-based behavioral analytics to protect your AI infrastructure.

### What's New in v1.0.4 - Security Hardening Release

This release includes critical security improvements:

- **✅ DEV_MODE Bypass Removed** - Eliminated environment variable bypass that allowed unlimited access
- **🔐 Cryptographic License Validation** - HMAC-SHA256 (Developer/Professional) and RSA-PKCS1v15 (Enterprise) signatures
- **🛡️ Hardware ID Binding** - Enterprise licenses bound to specific machines via hardware fingerprinting
- **📋 Unified 4-Tier Model** - Consistent Community→Developer→Professional→Enterprise licensing
- **🧹 Sensitive Data Cleanup** - Removed internal documentation from public repository

**Existing installations**: See [MIGRATION.md](MIGRATION.md) for license upgrade instructions.

### Key Capabilities

- AI-Native Security: Purpose-built for OpenAI, Anthropic, Azure OpenAI, AWS Bedrock, and custom AI providers
- Compliance Frameworks: MITRE ATLAS, OWASP AI Top 10, HIPAA, PCI-DSS, SOC 2, GDPR, ISO 27001, ISO 42001, NIST AI RMF
- ML-Powered Detection: Real-time anomaly detection, prompt injection prevention, behavioral analytics
- Enterprise SSO: SAML, OIDC, OAuth 2.0, LDAP integration
- Observability: Full SIEM integration (Splunk, Elastic, Datadog, QRadar)
- Cloud-Native: Kubernetes, Helm, Terraform, Docker deployment

---

## Features

### Security & Threat Detection
- **Prompt Injection Prevention**: Real-time detection of malicious prompts (OWASP LLM01)
- **Data Leakage Protection**: PII/PHI/PCI scanning before and after AI interactions
- **Rate Limiting**: Intelligent request throttling with burst protection
- **mTLS & PKI**: Certificate-based authentication and attestation
- **Secret Rotation**: Automated key rotation with zero downtime

### Compliance & Governance
- **Framework Coverage**: 10+ compliance frameworks with automated evidence generation
- **Audit Trails**: Cryptographically signed, tamper-evident logging
- **Policy Engine**: Custom security policies with real-time enforcement
- **Gap Analysis**: Automated compliance assessment and remediation guidance

### AI-Specific Protections
- **Request/Response Scanning**: Deep content inspection for AI interactions
- **Cost Anomaly Detection**: ML-powered spending pattern analysis
- **Multi-Tenant Isolation**: Department-level separation and access control
- **Token Usage Monitoring**: Real-time quota enforcement and alerting

### Deployment Options
- **Self-Hosted**: Docker, Kubernetes, Helm charts
- **Hybrid**: On-premise with cloud analytics
- **Air-Gapped**: Fully offline operation with local license validation

---

## Architecture

AegisGate follows a modular, plugin-based architecture with 246 files, 94,701+ lines of code, 99% Go.

### Package Structure

| Package | Description | Size |
|---------|-------------|------|
| `pkg/proxy/` | HTTP/2, HTTP/3, mTLS proxying | 86KB |
| `pkg/compliance/` | Compliance frameworks (9+ standards) | 35KB |
| `pkg/threatintel/` | STIX/TAXII, IOC correlation | 71KB |
| `pkg/ml/` | Anomaly detection, ML models | 49KB |
| `pkg/siem/` | Security event integration | 37KB |
| `pkg/sso/` | SAML, OIDC, OAuth | 27KB |
| `pkg/webhook/` | Event system | 33KB |

---

## Quick Start

### Prerequisites
- Go 1.24+ (for building from source)
- Docker 20.10+ (for containerized deployment)
- Kubernetes 1.25+ (optional, for K8s deployment)

### Installation

#### Option 1: Docker Compose (Recommended for evaluation)

```bash
# Clone repository
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# Start with Docker Compose (includes PostgreSQL, Redis, Grafana)
docker-compose -f deploy/docker/docker-compose.yml up -d

# Verify installation
curl http://localhost:8080/health
```

#### Option 2: Kubernetes (Helm)

```bash
# Add Helm repository
helm repo add aegisgate https://charts.aegisgate.io
helm repo update

# Install with default values
helm install aegisgate aegisgate/aegisgate \
  --namespace aegisgate \
  --create-namespace
```

#### Option 3: Binary Installation

```bash
# Download latest release
curl -L https://github.com/aegisgatesecurity/aegisgate/releases/download/v1.0.4/aegisgate-linux-amd64.tar.gz | tar xz

# Run with default config
./aegisgate --config config/aegisgate.yml.example
```

### Configuration

Create your configuration file:

```yaml
# config/aegisgate.yml
server:
  port: 8080
  tls:
    enabled: true
    cert: /path/to/cert.pem
    key: /path/to/key.pem

security:
  # License key (required for paid features)
  license_key: ${AEGISGATE_LICENSE_KEY}
  
  # Threat detection
  threat_detection:
    enabled: true
    ml_models: ["prompt_injection", "data_leakage", "cost_anomaly"]
    
  # Compliance
  compliance:
    enabled: true
    frameworks: ["owasp", "atlas", "gdpr"]

proxy:
  upstream:
    openai:
      url: https://api.openai.com
      api_key: ${OPENAI_API_KEY}
    anthropic:
      url: https://api.anthropic.com
      api_key: ${ANTHROPIC_API_KEY}
```

### Environment Variables

```bash
# Required
export AEGISGATE_LICENSE_KEY="your-license-key"

# Optional
export OPENAI_API_KEY="your-openai-key"
export ANTHROPIC_API_KEY="your-anthropic-key"
export DATABASE_URL="postgresql://user:pass@localhost/aegisgate"
export REDIS_URL="redis://localhost:6379"
```

---

## Tiers and Licensing

AegisGate uses a unified 4-tier licensing model with cryptographic validation.

### Free Tier: Community
**Perfect for evaluation and personal projects**

- 200 requests/minute
- 5 concurrent connections  
- 3 users
- OpenAI + Anthropic support
- OWASP, GDPR (view-only)
- Basic threat detection

**License**: Not required (defaults to Community)

### Paid Tiers

| Feature | Developer ($29/mo) | Professional ($99/mo) | Enterprise (Custom) |
|---------|-------------------|----------------------|---------------------|
| Requests/min | 1,000 | 5,000 | Unlimited |
| Users | 10 | 25 | Unlimited |
| AI Providers | 4 | 6 | All |
| Compliance | OWASP + NIST | All frameworks | All + ISO 42001 |
| SSO | OAuth | SAML/OIDC | Full LDAP |
| ML Features | Basic | Advanced | Custom models |
| SIEM | Dashboard only | Full integration | All platforms |
| Support | Email | Priority | 24/7 Dedicated |
| Deployment | Docker/K8s | K8s/Helm | Air-gapped/HSM |

### License Format

Licenses use cryptographic signing:

- Developer/Professional: `base64(JSON).base64(HMAC-SHA256)`
- Enterprise: `base64(JSON).base64(RSA-SIGNATURE)` + hardware binding

### Generating a License

```bash
# Generate Developer/Professional license
go run cmd/licensegen/main.go \
  -tier=professional \
  -email=user@company.com \
  -days=365 \
  -secret=/path/to/.license-secret

# Generate Enterprise license with hardware binding
go run cmd/licensegen/main.go \
  -tier=enterprise \
  -email=admin@enterprise.com \
  -days=365 \
  -hardware-id="$(cat /sys/class/net/eth0/address | tr -d ':')$(cat /proc/cpuinfo | grep Serial | cut -d' ' -f2)" \
  -secret=/path/to/.license-secret
```

---

## Security

### Security-First Design

AegisGate implements defense in depth:

1. **Transport Security**: TLS 1.3, mTLS, HTTP/2, HTTP/3
2. **Authentication**: OAuth 2.0, OIDC, SAML 2.0, LDAP
3. **Authorization**: RBAC, ABAC, fine-grained permissions
4. **Data Protection**: Encryption at rest and in transit
5. **Runtime Security**: Seccomp, AppArmor, hardened containers

### Compliance Frameworks

| Framework | Status |
|-----------|--------|
| OWASP AI Top 10 | Complete |
| MITRE ATLAS | Complete |
| SOC 2 Type II | Complete |
| HIPAA | Complete |
| PCI-DSS | Complete |
| GDPR | Complete |
| ISO 27001 | Complete |
| ISO 42001 | Complete |
| NIST AI RMF | Complete |

### Vulnerability Disclosure

If you discover a vulnerability:

1. DO NOT open a public issue
2. Email security@aegisgate.io with details
3. Include steps to reproduce
4. Allow 90 days for remediation before public disclosure

---

## Documentation

### Getting Started
- [Quick Start Guide](docs/quickstart.md)
- [Architecture Overview](docs/architecture.md)
- [Configuration Reference](docs/configuration.md)

### Deployment
- [Docker Deployment](docs/deploy/docker.md)
- [Kubernetes/Helm](docs/deploy/kubernetes.md)
- [Terraform](docs/deploy/terraform.md)
- [Air-Gapped](docs/deploy/airgapped.md)

### Security
- [Threat Model](docs/security/threat-model.md)
- [Compliance Guide](docs/compliance/README.md)
- [License Security](docs/security/license-security.md)

---

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone and build
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# Install dependencies
go mod download

# Run tests
go test ./...

# Run locally
go run cmd/aegisgate/main.go
```

---

## Project Statistics

- **246 Files**
- **94,701+ Lines of Code**
- **99% Go**
- **3,900+ Functions**
- **1,050+ Types/Structs**

---

## Support

| Resource | Link |
|----------|------|
| Documentation | https://github.com/aegisgatesecurity/aegisgate/tree/main/docs |
| GitHub Issues | https://github.com/aegisgatesecurity/aegisgate/issues |
| Website | Coming Soon |

---

## License

MIT License - Copyright 2025-2026 AegisGate Security. All rights reserved.