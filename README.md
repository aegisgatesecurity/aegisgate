<!--
AegisGate - Enterprise AI API Security Platform
Comprehensive README Documentation
Version: 1.0.0
Last Updated: March 2026
-->

<div align="center">

# 🔐 AegisGate

### Enterprise AI API Security Platform

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-brightgreen)](https://golang.org)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aegisgatesecurity/aegisgate/actions)
[![Docker](https://img.shields.io/badge/Docker-ready-brightblue)](https://hub.docker.com/r/aegisgate/aegisgate)
[![Enterprise](https://img.shields.io/badge/Edition-Open%20Core-orange)](https://aegisgatesecurity.io/pricing)

*A comprehensive security platform for protecting AI API infrastructure*

[Features](#features) • [Architecture](#architecture) • [Quick Start](#quick-start) • [Documentation](#documentation) • [Community](#community) • [Enterprise](#enterprise)

</div>

---

## Table of Contents

1. [What is AegisGate?](#what-is-aegisgate)
2. [Why AegisGate?](#why-aegisgate)
3. [Features](#features)
4. [Architecture](#architecture)
5. [Quick Start](#quick-start)
6. [Deployment](#deployment)
7. [Configuration](#configuration)
8. [Tier System](#tier-system)
9. [API Reference](#api-reference)
10. [Security](#security)
11. [Contributing](#contributing)
12. [Community](#community)
13. [Enterprise](#enterprise)
14. [Roadmap](#roadmap)
15. [License](#license)

---

## What is AegisGate?

**AegisGate** is an enterprise-grade, open-core AI API security platform designed to protect your AI infrastructure from emerging threats. It serves as a reverse proxy that sits between your applications and AI providers (OpenAI, Anthropic, Azure OpenAI, AWS Bedrock, Google Vertex, and more), providing:

- 🔒 **Threat Detection** - Real-time detection of prompt injections, jailbreaks, and malicious payloads
- 📊 **Compliance** - Built-in frameworks for SOC 2, HIPAA, PCI-DSS, GDPR, NIST, and more
- 🛡️ **Defense in Depth** - Multiple layers of security including mTLS, device attestation, and behavioral analysis
- 📈 **Observability** - Comprehensive metrics, logging, and SIEM integrations
- ⚡ **Performance** - HTTP/2 and HTTP/3 support with intelligent request optimization

### Open Core Model

AegisGate follows an **Open Core** business model:
- **Community Edition** - Free, open-source version with core security features
- **Paid Tiers** - Developer, Professional, and Enterprise tiers with advanced capabilities

---

## Why AegisGate?

### The AI Security Challenge

As AI becomes central to business operations, security teams face unique challenges:

| Challenge | AegisGate Solution |
|-----------|-----------------|
| **Prompt Injection** | Real-time prompt analysis with ML-based detection |
| **Data Exfiltration** | Request/response filtering and DLP capabilities |
| **Compliance** | Pre-built compliance frameworks (SOC 2, HIPAA, PCI-DSS) |
| **API Abuse** | Rate limiting, quota management, and usage analytics |
| **Zero-Day Threats** | Behavioral analysis and anomaly detection |
| **Multi-Provider** | Single integration for OpenAI, Anthropic, Azure, AWS, GCP |

### Key Benefits

```
┌─────────────────────────────────────────────────────────────────┐
│                     AEGISGATE BENEFITS                             │
├─────────────────────────────────────────────────────────────────┤
│  ✓ 5-Minute Integration     - Drop-in reverse proxy           │
│  ✓ Zero Code Changes        - Protect existing apps           │
│  ✓ Enterprise-Grade        - mTLS, device attestation         │
│  ✓ Multi-Cloud Ready       - AWS, Azure, GCP, on-prem         │
│  ✓ Open Source             - Community edition available      │
│  ✓ SOC 2 Compliant         - Security audit ready             │
└─────────────────────────────────────────────────────────────────┘
```

---

## Features

### Core Security Features

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **TLS Termination** | ✅ | ✅ | ✅ | ✅ |
| **Rate Limiting** | 200/min | 1,000/min | 5,000/min | Unlimited |
| **Request Logging** | ✅ | ✅ | ✅ | ✅ |
| **Error Tracking** | ✅ | ✅ | ✅ | ✅ |
| **Metrics (Prometheus)** | ✅ | ✅ | ✅ | ✅ |
| **mTLS Authentication** | - | ✅ | ✅ | ✅ |
| **Device Attestation** | - | - | ✅ | ✅ |
| **OAuth/SSO** | - | ✅ | ✅ | ✅ |
| **SAML/OIDC** | - | - | ✅ | ✅ |
| **Hardware Token (HSM)** | - | - | - | ✅ |

### AI Security Features

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **ML Anomaly Detection** | Basic | Advanced | Advanced | Advanced |
| **Prompt Injection Detection** | - | ✅ | ✅ | ✅ |
| **Content Analysis** | - | - | ✅ | ✅ |
| **Behavioral Analysis** | - | - | - | ✅ |
| **Custom ML Models** | - | - | - | ✅ |
| **Real-time Threat Response** | - | - | - | ✅ |
| **Zero-Day Protection** | - | - | - | ✅ |

### Compliance Frameworks

| Framework | Community | Developer | Professional | Enterprise |
|-----------|:---------:|:---------:|:------------:|:----------:|
| **OWASP** | View | View | Full | Full |
| **SOC 2** | View | View | Full | Full |
| **GDPR** | View | View | Full | Full |
| **NIST** | - | View | Full | Full |
| **HIPAA** | - | - | Full | Full |
| **PCI-DSS** | - | - | Full | Full |
| **ISO 27001** | - | - | Full | Full |
| **FedRAMP** | - | - | - | Full |
| **NIST AI RMF** | - | - | - | Full |
| **HITRUST** | - | - | - | Full |

### Observability & Integrations

| Integration | Community | Developer | Professional | Enterprise |
|-------------|:---------:|:---------:|:------------:|:----------:|
| **Prometheus** | ✅ | ✅ | ✅ | ✅ |
| **Grafana** | - | ✅ | ✅ | ✅ |
| **Datadog** | - | - | ✅ | ✅ |
| **New Relic** | - | - | ✅ | ✅ |
| **Splunk** | - | - | ✅ | ✅ |
| **Elastic** | - | - | ✅ | ✅ |
| **QRadar** | - | - | - | ✅ |
| **Azure Sentinel** | - | - | - | ✅ |
| **AWS CloudWatch** | - | ✅ | ✅ | ✅ |

### Deployment Options

| Option | Community | Developer | Professional | Enterprise |
|--------|:---------:|:---------:|:------------:|:----------:|
| **Docker** | ✅ | ✅ | ✅ | ✅ |
| **Docker Compose** | ✅ | ✅ | ✅ | ✅ |
| **Kubernetes** | - | - | ✅ | ✅ |
| **Helm** | - | - | ✅ | ✅ |
| **Terraform** | - | ✅ | ✅ | ✅ |
| **Service Mesh** | - | - | - | ✅ |
| **Multi-Region** | - | - | - | ✅ |
| **On-Premise** | - | - | - | ✅ |
| **Air-Gapped** | - | - | - | ✅ |

---

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              AEGISGATE ARCHITECTURE                            │
└─────────────────────────────────────────────────────────────────────────────┘

    ┌─────────────┐          ┌─────────────┐          ┌─────────────┐
    │   Client    │          │   Client    │          │   Client    │
    │  Application│          │  Application│          │  Application│
    └──────┬──────┘          └──────┬──────┘          └──────┬──────┘
           │                        │                        │
           │ HTTPS                  │ HTTPS                  │ HTTPS
           │                        │                        │
    ┌──────┴────────────────────────┴────────────────────────┴──────┐
    │                                                                 │
    │                   AEGISGATE PROXY CLUSTER                        │
    │  ┌─────────────────────────────────────────────────────────┐  │
    │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐    │  │
    │  │  │  TLS    │  │   ML    │  │  Rate   │  │Compliance│    │  │
    │  │  │Termination│  │Detection│  │ Limiter │  │  Engine  │    │  │
    │  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘    │  │
    │  │       │            │            │            │          │  │
    │  │       └────────────┼────────────┼────────────┘          │  │
    │  │                    │            │                        │  │
    │  │              ┌─────┴────────────┴─────┐                 │  │
    │  │              │     Request Pipeline    │                 │  │
    │  │              │  (Validation → Transform)│                │  │
    │  │              └───────────┬─────────────┘                 │  │
    │  └──────────────────────────┬┴──────────────────────────────┘  │
    └─────────────────────────────┼──────────────────────────────────┘
                                  │
              ┌───────────────────┼───────────────────┐
              │                   │                   │
              ▼                   ▼                   ▼
    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
    │   OpenAI    │    │  Anthropic  │    │   Azure     │
    │     API     │    │     API     │    │   OpenAI    │
    └─────────────┘    └─────────────┘    └─────────────┘
```

### Component Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      AEGISGATE COMPONENTS                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐       │
│  │    Proxy     │   │   Metrics    │   │  Compliance  │       │
│  │   Server     │   │   Manager    │   │   Manager    │       │
│  │              │   │              │   │              │       │
│  │ • HTTP/1.1   │   │ • Prometheus │   │ • SOC 2      │       │
│  │ • HTTP/2     │   │ • Grafana    │   │ • HIPAA      │       │
│  │ • HTTP/3     │   │ • Datadog    │   │ • PCI-DSS    │       │
│  │ • WebSocket  │   │ • Custom     │   │ • GDPR       │       │
│  └──────────────┘   └──────────────┘   │ • NIST       │       │
│                                        └──────────────┘       │
│  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐       │
│  │  ML Engine   │   │  Auth/Session │   │   Tenant     │       │
│  │              │   │   Manager    │   │   Manager    │       │
│  │ • Anomaly    │   │              │   │              │       │
│  │ • Prompt Inj │   │ • OAuth      │   │ • Isolation  │       │
│  │ • Behavioral │   │ • SAML/OIDC  │   │ • Quotas     │       │
│  │ • Custom     │   │ • Session    │   │ • RBAC       │       │
│  └──────────────┘   └──────────────┘   └──────────────┘       │
│                                                                  │
│  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐       │
│  │   SIEM       │   │   Webhook    │   │   Secrets   │       │
│  │  Integration │   │   Manager    │   │   Manager   │       │
│  │              │   │              │   │              │       │
│  │ • Splunk     │   │ • Alerts     │   │ • Vault     │       │
│  │ • Elastic    │   │ • Events     │   │ • AWS SM    │       │
│  │ • QRadar     │   │ • Filtering  │   │ • Azure Key │       │
│  │ • Syslog     │   │ • Retry      │   │ • Env Vars  │       │
│  └──────────────┘   └──────────────┘   └──────────────┘       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Request Flow

```
CLIENT REQUEST → TLS TERMINATION → RATE LIMIT CHECK → ML DETECTION
       ↓                                                            ↓
  PARSE REQUEST → VALIDATE AUTH → CHECK COMPLIANCE → FORWARD TO AI
       ↑                                                            ↓
RESPONSE LOGGED ← FILTER RESPONSE ← ANALYZE CONTENT ← PROCESS RESPONSE
```

---

## Quick Start

### Prerequisites

| Requirement | Minimum | Recommended |
|-------------|---------|-------------|
| **Go** | 1.21+ | 1.23+ |
| **Docker** | 20.10+ | Latest |
| **Kubernetes** | 1.24+ | 1.28+ |
| **Memory** | 512 MB | 2 GB |
| **CPU** | 1 core | 2+ cores |

### Option 1: Docker (Recommended)

```bash
# Pull the official image
docker pull aegisgate/aegisgate:latest

# Run with basic configuration
docker run -d \
  --name aegisgate \
  -p 8080:8080 \
  -e TARGET_URL=https://api.openai.com \
  -e TIER=community \
  aegisgate/aegisgate:latest

# Test the proxy
curl http://localhost:8080/health
```

### Option 2: Docker Compose

```bash
# Clone the repository
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# Start with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f aegisgate
```

### Option 3: Binary

```bash
# Download the latest release
curl -L -o aegisgate https://github.com/aegisgatesecurity/aegisgate/releases/latest/aegisgate

# Make executable
chmod +x aegisgate

# Run the proxy
./aegisgate --target https://api.openai.com --tier community
```

### Option 4: From Source

```bash
# Clone the repository
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# Build the binary
make build

# Run the proxy
./bin/aegisgate --target https://api.openai.com --tier community
```

### Verify Installation

```bash
# Health check
curl http://localhost:8080/health

# Version info
curl http://localhost:8080/version

# Statistics
curl http://localhost:8080/stats
```

Expected responses:

```json
// Health
{"status":"healthy","tier":"Community","proxy_enabled":true}

// Version
{"version":"1.0.0","commit":"abc123","date":"2026-03-10T12:00:00Z"}

// Stats
{"request_count":0,"enabled":true}
```

---

## Deployment

### Docker Deployment

```bash
# Basic deployment
docker run -d \
  --name aegisgate \
  -p 8080:8080 \
  -p 8443:8443 \
  -e TARGET_URL=https://api.openai.com \
  -e TIER=professional \
  -e API_KEY=$OPENAI_API_KEY \
  -v aegisgate-data:/data \
  aegisgate/aegisgate:latest
```

### Kubernetes Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: aegisgate
spec:
  replicas: 3
  selector:
    matchLabels:
      app: aegisgate
  template:
    metadata:
      labels:
        app: aegisgate
    spec:
      containers:
      - name: aegisgate
        image: aegisgate/aegisgate:latest
        ports:
        - containerPort: 8080
        env:
        - name: TARGET_URL
          value: "https://api.openai.com"
        - name: TIER
          value: "enterprise"
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "2000m"
```

```bash
# Deploy to Kubernetes
kubectl apply -f deploy/k8s/
```

### Helm Chart

```bash
# Add the Helm repository
helm repo add aegisgate https://aegisgatesecurity.github.io/aegisgate

# Install the chart
helm install my-aegisgate aegisgate/aegisgate \
  --set targetUrl=https://api.openai.com \
  --set tier=enterprise
```

---

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `TARGET_URL` | AI provider URL | - | Yes |
| `TIER` | License tier | `community` | No |
| `PORT` | HTTP listen port | `8080` | No |
| `BIND_ADDRESS` | Bind address | `0.0.0.0:8080` | No |
| `RATE_LIMIT` | Requests per minute | Tier-based | No |
| `LOG_LEVEL` | Log level | `info` | No |
| `DATABASE_URL` | Database connection | - | No |

### Configuration File

Create a `config.yaml`:

```yaml
server:
  host: 0.0.0.0
  port: 8080
  tls:
    enabled: true
    cert: /path/to/cert.pem
    key: /path/to/key.pem

upstream:
  url: https://api.openai.com
  timeout: 30s
  retry: 3

security:
  tier: professional
  rate_limit: 5000
  ml_detection: true
  prompt_injection: true

compliance:
  enabled: true
  frameworks:
    - soc2
    - hipaa
    - pci

observability:
  metrics:
    enabled: true
    port: 9090
  logging:
    level: info
    format: json
```

### Tier-Specific Configuration

#### Community Tier

```bash
# .env-community
TIER=community
RATE_LIMIT=200
MAX_USERS=3
MAX_API_KEYS=3
LOG_RETENTION_DAYS=1
ML_DETECTION=false
```

#### Developer Tier

```bash
# .env-developer
TIER=developer
RATE_LIMIT=1000
MAX_USERS=10
MAX_API_KEYS=10
LOG_RETENTION_DAYS=7
ML_DETECTION=true
PROMPT_INJECTION=true
SUPPORT=email
```

#### Professional Tier

```bash
# .env-professional
TIER=professional
RATE_LIMIT=5000
MAX_USERS=25
MAX_API_KEYS=50
LOG_RETENTION_DAYS=30
ML_DETECTION=true
PROMPT_INJECTION=true
CONTENT_ANALYSIS=true
SUPPORT=priority
COMPLIANCE_FRAMEWORKS=soc2,hipaa,pci,nist
```

#### Enterprise Tier

```bash
# .env-enterprise
TIER=enterprise
RATE_LIMIT=0  # Unlimited
MAX_USERS=unlimited
MAX_API_KEYS=unlimited
LOG_RETENTION_DAYS=unlimited
ML_DETECTION=true
PROMPT_INJECTION=true
CONTENT_ANALYSIS=true
BEHAVIORAL_ANALYSIS=true
CUSTOM_ML_MODELS=true
SUPPORT=24x7
COMPLIANCE_FRAMEWORKS=all
```

---

## Tier System

AegisGate uses a tiered licensing model to support both open-source users and enterprise customers.

### Tier Comparison

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **Price** | Free | $29/mo | $99/mo | Custom |
| **Support** | Community | Email | Priority | 24/7 Dedicated |
| **Rate Limit** | 200/min | 1,000/min | 5,000/min | Unlimited |
| **Users** | 3 | 10 | 25 | Unlimited |
| **API Keys** | 3 | 10 | 50 | Unlimited |
| **Log Retention** | 1 day | 7 days | 30 days | Unlimited |
| **ML Features** | Basic | Advanced | Advanced | Advanced |
| **Compliance** | View Only | View Only | Full | Full |
| **Deployment** | Docker | Docker | K8s | Multi-Region |
| **SLA** | - | - | 99.9% | 99.99% |

### Upgrading Tiers

```bash
# Upgrade from Community to Developer
./aegisgate --tier developer --license-key YOUR_LICENSE_KEY

# Upgrade from Developer to Professional  
./aegisgate --tier professional --license-key YOUR_LICENSE_KEY

# Enterprise (contact sales)
# https://aegisgatesecurity.io/contact
```

### Feature Access

Access to features is controlled at the API level:

```go
// Check if feature is available for tier
import "github.com/aegisgatesecurity/aegisgate/pkg/core"

tier := core.GetTierByName("professional")
requiredTier := core.GetRequiredTier("ml_behavioral")

if tier.CanAccess(requiredTier) {
    // Feature is available
}
```

---

## API Reference

### Management Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/version` | GET | Version information |
| `/stats` | GET | Request statistics |
| `/metrics` | GET | Prometheus metrics |
| `/config` | GET | Current configuration |

### Proxy Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/*` | * | Forward to AI provider |
| `/*` | * | Catch-all proxy |

### Example: Proxying OpenAI

```bash
# Instead of calling OpenAI directly:
curl https://api.openai.com/v1/chat/completions \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{"model":"gpt-4","messages":[{"role":"user","content":"Hello"}]}'

# Call through AegisGate:
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "X-AegisGate-Key: $AEGISGATE_API_KEY" \
  -d '{"model":"gpt-4","messages":[{"role":"user","content":"Hello"}]}'
```

### Response Headers

All proxy responses include:

```
X-AegisGate-Tier: Community
X-AegisGate-Request-ID: req_abc123
X-AegisGate-Version: 1.0.0
X-RateLimit-Remaining: 150
```

---

## Security

### Security Features

- ✅ **TLS 1.3** - Modern transport security
- ✅ **mTLS** - Mutual TLS authentication
- ✅ **Device Attestation** - Hardware-backed identity
- ✅ **OAuth 2.0 / OIDC** - Enterprise SSO
- ✅ **SAML** - Enterprise identity providers
- ✅ **HSM Integration** - Hardware security modules
- ✅ **Secret Rotation** - Automated credential management
- ✅ **Audit Encryption** - Encrypted audit logs
- ✅ **FIPS 140-2** - Federal security standards
- ✅ **Runtime Hardening** - System-level protections

### Vulnerability Reporting

If you discover a security vulnerability, please report it responsibly:

```bash
# Do NOT report vulnerabilities in public issues
# Email: security@aegisgatesecurity.ioaegisgatesecurity.io
# PGP: https://aegisgatesecurity.io/security/pgp-key
```

See [SECURITY.md](SECURITY.md) for full security documentation.

---

## Contributing

We welcome contributions! Please see our contributing guidelines:

### Development Setup

```bash
# Clone the repository
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# Install dependencies
make deps

# Run tests
make test

# Build
make build

# Run development server
make run
```

### Code Style

- Follow Go conventions
- Run `golangci-lint` before committing
- Write tests for new features
- Update documentation

### Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Update documentation
6. Submit a pull request

See [CONTRIBUTING.md](CONTRIBUTING.md) for full contribution guidelines.

---

## Community

### Join the Community

| Channel | Link |
|---------|------|
| **GitHub Discussions** | [github.com/aegisgatesecurity/aegisgate/discussions](https://github.com/aegisgatesecurity/aegisgate/discussions) |
| **Discord** | [discord.gg/aegisgate](https://discord.gg/aegisgate) |
| **Twitter** | [@aegisgate_io](https://twitter.com/aegisgate_io) |
| **Blog** | [blog.aegisgatesecurity.io](https://aegisgatesecurity.io) |

### Resources

- 📚 [Documentation](https://aegisgatesecurity.io)
- 💬 [Community Forum](https://aegisgatesecurity.io)
- 🎫 [Issue Tracker](https://github.com/aegisgatesecurity/aegisgate/issues)
- 📊 [Feature Requests](https://github.com/aegisgatesecurity/aegisgate/issues/new?template=feature_request.md)

---

## Enterprise

### AegisGate Enterprise

For enterprise customers, we offer:

- **Custom Deployments** - On-premise, air-gapped, multi-region
- **Dedicated Support** - 24/7 response with SLA
- **Custom Integrations** - Proprietary AI providers and systems
- **Advanced Compliance** - Custom framework implementations
- **Training** - Security and operations training
- **Audit Support** - Third-party audit assistance

### Contact Sales

```
Email: sales@aegisgatesecurity.ioaegisgatesecurity.io
Web: https://aegisgatesecurity.io/contact
Phone: +1 (555) 123-4567
```

### Enterprise Features

| Feature | Description |
|---------|-------------|
| **Custom SLAs** | 99.99% uptime guarantees |
| **Dedicated Infrastructure** | Isolated clusters |
| **Custom Compliance** | Industry-specific frameworks |
| **Security Reviews** | Annual penetration testing |
| **Migration Assistance** | Zero-downtime transitions |
| **Training Programs** | Certified administrator courses |

---

## Roadmap

### Upcoming Features

| Feature | Tier | Status |
|---------|------|--------|
| GraphQL Support | Enterprise | Planned |
| Browser Extension | Enterprise | Planned |
| Serverless Deployment | Enterprise | Planned |
| Custom Domain/Whitelabel | Enterprise | Planned |
| Cross-Tenant Analytics | Enterprise | Planned |
| Usage-Based Billing | Enterprise | Planned |
| Multi-Organization Support | Enterprise | Planned |

### Version History

See [CHANGELOG.md](CHANGELOG.md) for detailed version history.

---

## License

AegisGate is licensed under the **Apache License 2.0**.

### Community Edition

The Community Edition is open source and free to use under the Apache 2.0 license. See [LICENSE](LICENSE) for details.

### Commercial Tiers

Commercial use of Developer, Professional, and Enterprise tiers requires a license agreement. Contact [sales@aegisgatesecurity.ioaegisgatesecurity.io](mailto:sales@aegisgatesecurity.ioaegisgatesecurity.io) for details.

---

## Quick Reference

```bash
# Docker Quick Start
docker run -d -p 8080:8080 -e TARGET_URL=https://api.openai.com aegisgate/aegisgate:latest

# Build from Source
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
make build
./bin/aegisgate --target https://api.openai.com

# View Help
./aegisgate --help

# Check Version
./aegisgate --version
```

---

<div align="center">

### ⭐ Star us on GitHub!

If AegisGate helps you secure your AI infrastructure, please star us!

[![GitHub stars](https://img.shields.io/github/stars/aegisgatesecurity/aegisgate?style=social)](https://github.com/aegisgatesecurity/aegisgate)

---

**Made with 🔐 by the AegisGate Team**

*© 2024-2026 AegisGate Security, Inc. All rights reserved.*

</div>
