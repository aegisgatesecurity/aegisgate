<div align="center">

# 🛡️ AegisGate™ 🔐

Enterprise AI API Security Platform

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/)
[![Version](https://img.shields.io/badge/version-v1.0.13-green?logo=semver)](https://github.com/aegisgatesecurity/aegisgate/releases)
[![Release Date](https://img.shields.io/badge/Release-Date-2026--03--19-brightgreen)](https://github.com/aegisgatesecurity/aegisgate/releases)
[![DCO](https://img.shields.io/badge/DCO-Signed-blue)](CONTRIBUTING.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/aegisgatesecurity/aegisgate)](https://goreportcard.com/report/github.com/aegisgatesecurity/aegisgate)

[![Main CI](https://img.shields.io/github/actions/workflow/status/aegisgatesecurity/aegisgate/main.yml?branch=main&label=Main%20CI)](https://github.com/aegisgatesecurity/aegisgate/actions/workflows/main.yml)
[![Security Tests](https://img.shields.io/github/actions/workflow/status/aegisgatesecurity/aegisgate/security.yml?branch=main&label=Security)](https://github.com/aegisgatesecurity/aegisgate/actions/workflows/security.yml)
[![DCO Check](https://img.shields.io/github/actions/workflow/status/aegisgatesecurity/aegisgate/dco-check.yml?branch=main&label=DCO)](https://github.com/aegisgatesecurity/aegisgate/actions/workflows/dco-check.yml)

</div>

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Licensing](#licensing)
- [Contributing](#contributing)
- [Security](#security)

---

## Overview

AegisGate is an enterprise-grade, Go-based transparent proxy that sits between AI clients and AI service providers. It provides security, governance, and observability for organizational AI API usage.

## Features

🛡️ **Enterprise Security**
- Transparent proxy for all major AI providers
- PII detection and redaction
- Request/response filtering and validation
- Jailbreak and prompt injection detection

🔒 **Data Privacy**
- Local processing - no data leaves your infrastructure
- Configurable data retention policies
- Token-level access control

✅ **Compliance & Governance**
- Request/response logging and auditing
- Usage analytics and reporting
- Rate limiting and quota management
- Cost controls by department, user, or API key

⚡ **Performance**
- Sub-millisecond latency overhead
- Connection pooling
- Caching for repeated queries

🚀 **Easy Integration**
- Drop-in replacement for OpenAI SDKs
- Docker and Kubernetes support
- Prometheus metrics exposed
- OpenTelemetry tracing support

---

## Quick Start

### Prerequisites

- Go 1.24 or later
- Docker (optional)
- Kubernetes (optional)

### Installation

```bash
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
go build -o aegisgate cmd/aegisgate/main.go
./aegisgate
```

### Docker

```bash
docker build -t aegisgate:latest .
docker run -p 8080:8080 aegisgate:latest
```

---

## Architecture

AegisGate operates as a transparent proxy:

```
┌─────────────┐      ┌─────────────┐      ┌─────────────────┐
│   AI Client  │─────▶│  AegisGate   │─────▶│  AI Provider    │
└─────────────┘      └─────────────┘      └─────────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   Storage   │
                    └─────────────┘
```

---

## Configuration

See `config.example.yaml` for all available options.

### Basic Configuration

```yaml
server:
  host: 0.0.0.0
  port: 8080

security:
  pii_detection:
    enabled: true
  prompt_injection:
    enabled: true
```

---

## API Documentation

### Proxy Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /v1/chat/completions | OpenAI-compatible Chat Completions |
| POST | /v1/completions | OpenAI-compatible Completions |
| POST | /v1/embeddings | OpenAI-compatible Embeddings |
| GET | /health | Health check |
| GET | /metrics | Prometheus metrics |

---

## Licensing

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **Price** | Free | $49/mo | $199/mo | Custom |
| **API Calls/mo** | 10,000 | 500,000 | Unlimited | Unlimited |
| **Users** | 1 | 10 | Unlimited | Unlimited |
| **PII Detection** | ❌ | ✅ | ✅ | ✅ |
| **Prompt Injection** | ❌ | ✅ | ✅ | ✅ |
| **SSO/SAML** | ❌ | ❌ | ✅ | ✅ |
| **SLA** | None | None | 99.9% | 99.99% |

Contact [enterprise@aegisgatesecurity.io](mailto:enterprise@aegisgatesecurity.io) for Enterprise licensing.

---

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md).

All contributions must include a DCO sign-off:
```
Signed-off-by: Your Name <your@email.com>
```

---

## Security

For security vulnerabilities, contact security@aegisgatesecurity.io

See [SECURITY.md](SECURITY.md) and [SECURITY_HARDENING.md](SECURITY_HARDENING.md).

---

<div align="center">

**AegisGate** - Enterprise AI API Security Platform

© 2024-2026 AegisGate Security. All rights reserved.

</div>