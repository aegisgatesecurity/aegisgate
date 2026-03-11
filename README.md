# AegisGate - Enterprise AI API Security Platform

<div align="center">

# 🔐 AegisGate

### Enterprise AI API Security Platform

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-brightgreen)](https://golang.org)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/aegisgatesecurity/aegisgate/actions)
[![Docker](https://img.shields.io/badge/Docker-ready-brightblue)](https://hub.docker.com/r/aegisgate/aegisgate)
[![Enterprise](https://img.shields.io/badge/Edition-Open%20Core-orange)](https://aegisgatesecurity.io/pricing)

*A comprehensive security platform for protecting AI API infrastructure*

[Features](#features) • [Architecture](#architecture) • [Quick Start](#quick-start) • [Documentation](#documentation) • [Enterprise](#enterprise) • [License](#license)

</div>

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
- **Paid Tiers** - Professional and Enterprise tiers with advanced capabilities

---

## Features

### Core Security Features

| Feature | Community | Professional | Enterprise |
|---------|:---------:|:------------:|:----------:|
| **TLS Termination** | ✅ | ✅ | ✅ |
| **Rate Limiting** | 200/min | 5,000/min | Unlimited |
| **Request Logging** | ✅ | ✅ | ✅ |
| **Error Tracking** | ✅ | ✅ | ✅ |
| **Metrics (Prometheus)** | ✅ | ✅ | ✅ |
| **mTLS Authentication** | - | ✅ | ✅ |
| **Device Attestation** | - | ✅ | ✅ |
| **OAuth/SSO** | - | ✅ | ✅ |
| **SAML/OIDC** | - | ✅ | ✅ |
| **Hardware Token (HSM)** | - | - | ✅ |

### AI Security Features

| Feature | Community | Professional | Enterprise |
|---------|:---------:|:------------:|:----------:|
| **ML Anomaly Detection** | Basic | Advanced | Advanced |
| **Prompt Injection Detection** | - | ✅ | ✅ |
| **Content Analysis** | - | ✅ | ✅ |
| **Behavioral Analysis** | - | - | ✅ |
| **Custom ML Models** | - | - | ✅ |
| **Real-time Threat Response** | - | - | ✅ |
| **Zero-Day Protection** | - | - | ✅ |

### Compliance Frameworks

| Framework | Community | Professional | Enterprise |
|-----------|:---------:|:------------:|:----------:|
| **OWASP** | View | Full | Full |
| **SOC 2** | View | Full | Full |
| **GDPR** | View | Full | Full |
| **NIST** | - | Full | Full |
| **HIPAA** | - | Full | Full |
| **PCI-DSS** | - | Full | Full |
| **ISO 27001** | - | Full | Full |

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
```

### Option 3: From Source

```bash
# Clone the repository
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# Build the binary
make build

# Run the proxy
./bin/aegisgate --target https://api.openai.com --tier community
```

---

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `TARGET_URL` | AI provider URL | - | Yes |
| `TIER` | License tier | `community` | No |
| `PORT` | HTTP listen port | `8080` | No |

### Configuration File

Create a `config.yaml`:

```yaml
server:
  host: 0.0.0.0
  port: 8080

upstream:
  url: https://api.openai.com
  timeout: 30s

security:
  tier: professional
  rate_limit: 5000
  ml_detection: true
```

---

## Enterprise

### Contact Sales

For enterprise licensing and custom deployments:

- **Email**: sales@aegisgatesecurity.io
- **Web**: https://aegisgatesecurity.io/contact

### Enterprise Features

| Feature | Description |
|---------|-------------|
| **Custom SLAs** | 99.99% uptime guarantees |
| **Dedicated Infrastructure** | Isolated clusters |
| **Custom Compliance** | Industry-specific frameworks |
| **24/7 Support** | Dedicated support team |
| **Migration Assistance** | Zero-downtime transitions |

---

## Security

### Vulnerability Reporting

If you discover a security vulnerability, please report it responsibly:

**DO NOT** create a public GitHub issue for security vulnerabilities.

- **Email**: security@aegisgatesecurity.io
- **PGP**: https://aegisgatesecurity.io/security/pgp-key

See [SECURITY.md](SECURITY.md) for full security documentation.

---

## License

AegisGate is licensed under the **Apache License 2.0**.

### Community Edition

The Community Edition is open source and free to use under the Apache 2.0 license. See [LICENSE](LICENSE) for details.

### Commercial Tiers

Commercial use of Professional and Enterprise tiers requires a license agreement. Contact sales@aegisgatesecurity.io for details.

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
```

---

<div align="center">

### ⭐ Star us on GitHub!

If AegisGate helps you secure your AI infrastructure, please star us!

[![GitHub stars](https://img.shields.io/github/stars/aegisgatesecurity/aegisgate?style=social)](https://github.com/aegisgatesecurity/aegisgate)

---

*© 2024-2026 AegisGate Security, Inc. All rights reserved.*

</div>
