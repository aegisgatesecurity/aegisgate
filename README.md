<div align="center">

# 🛡️ AegisGate™ — Transparent AI API Security Gateway

[![Version](https://img.shields.io/badge/version-v1.1.0-green?logo=semver)](https://github.com/aegisgatesecurity/aegisgate/releases)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25.8+-00ADD8?logo=go)](https://golang.org/)
[![Security](https://img.shields.io/badge/Security-0_CVEs-brightgreen?logo=shield)](SECURITY.md)
[![Test Coverage](https://img.shields.io/badge/Coverage-70%25-brightgreen?logo=codecov)](https://github.com/aegisgatesecurity/aegisgate)

[![Docker](https://img.shields.io/badge/Docker-27MB-2496ED?logo=docker)](https://hub.docker.com/r/aegisgatesecurity/aegisgate)
[![Kubernetes](https://img.shields.io/badge/K8s-Ready-326CE5?logo=kubernetes)](https://kubernetes.io/)
[![CI Status](https://img.shields.io/badge/CI-Passing-brightgreen?logo=github-actions)](https://github.com/aegisgatesecurity/aegisgate/actions)

[📚 Docs](https://aegisgatesecurity.io/docs) &nbsp;•&nbsp; [✨ Features](#-features) &nbsp;•&nbsp; [🚀 Quick Start](#-quick-start) &nbsp;•&nbsp; [🏗️ Architecture](#-architecture) &nbsp;•&nbsp; [⚡ Performance](#-performance-benchmarks) &nbsp;•&nbsp; [🔒 Security](#-security)

</div>

> **30-Second Pitch**: Your AI applications need a security guard — one that speaks HTTP, understands LLM threats, and adds less than 5ms latency. AegisGate™ is that guard. Deploy in 60 seconds. Sleep better tonight.

---

## ⚡ TL;DR

**AegisGate™** is a high-performance, transparent proxy that secures AI API traffic between your applications and LLM providers (OpenAI, Anthropic, Azure, AWS Bedrock, Cohere). Deploy in minutes for:

| 🛡️ **Security** | 📋 **Compliance** | 🚀 **Performance** |
|-----------------|------------------|-------------------|
| Real-time threat blocking | SOC2, HIPAA, PCI-DSS | **<5ms latency** |
| Prompt injection prevention | GDPR, ISO 27001, ISO 42001 | **50,000 req/s** |
| Data leakage protection | NIST AI RMF, MITRE ATLAS | **27MB Docker image** |
| Adversarial attack defense | OWASP LLM Top 10 | **4MB memory footprint** |

**No code changes required.** Point your AI traffic through AegisGate™ — done.

> 🎯 **Reality Check**: A single prompt injection attack can cost millions in data breaches, regulatory fines, and reputation damage. AegisGate™ protects against OWASP LLM Top 10 threats for **free** — the cost of doing nothing is far higher.

---

## 📦 License

### Apache License 2.0

AegisGate™ is released under the Apache License 2.0, a permissive open-source license that allows you to:

- ✅ Use the software for any purpose
- ✅ Modify and distribute the software
- ✅ Use in proprietary software
- ✅ Distribute copies to others

**Full Features Included:**
- Full security scanning features
- Prompt injection prevention (OWASP LLM Top 10)
- Data leakage protection
- MITRE ATLAS threat coverage
- All compliance frameworks
- Self-hosted deployment (Docker, Kubernetes, binary)
- All core proxy and compliance features

See [LICENSE](LICENSE) for full license text.

---

## 🛡️ Features

### Security Protection

| Feature | Description |
|---------|-------------|
| **Prompt Injection Prevention** | Blocks OWASP LLM Top 10 attacks including LLM01-LLM10 |
| **Data Leakage Protection** | Prevents sensitive data (PII, credentials, keys) exfiltration |
| **Adversarial Attack Defense** | Detects jailbreaks, DoS, and model manipulation attempts |
| **Behavioral Analysis** | ML-powered anomaly detection for suspicious patterns |
| **Real-time Threat Blocking** | Sub-millisecond detection and blocking |

### Compliance Frameworks

| Framework | Coverage |
|-----------|----------|
| **OWASP LLM Top 10** | LLM01-LLM10 complete coverage |
| **MITRE ATLAS** | All AI-specific attack patterns |
| **SOC 2** | Security and availability controls |
| **HIPAA** | Healthcare data protection |
| **PCI-DSS** | Payment card security |
| **GDPR** | EU data protection |
| **ISO 27001** | Information security |
| **ISO 42001** | AI management systems |
| **NIST AI RMF** | AI risk management |

### AI Provider Support

| Provider | Status |
|----------|--------|
| OpenAI (GPT-4, GPT-3.5) | ✅ Supported |
| Anthropic (Claude) | ✅ Supported |
| Azure OpenAI | ✅ Supported |
| AWS Bedrock | ✅ Supported |
| Google AI (Gemini) | ✅ Supported |
| Cohere | ✅ Supported |
| Custom Endpoints | ✅ Supported |

---

## 🚀 Quick Start

### Docker (Recommended)

```bash
# Run with default configuration
docker run -d \
  --name aegisgate \
  -p 8080:8080 \
  -p 8443:8443 \
  -e OPENAI_API_KEY=your-api-key \
  ghcr.io/aegisgatesecurity/aegisgate:latest

# Verify it's running
curl http://localhost:8080/health
```

### Docker Compose

```yaml
version: '3.8'
services:
  aegisgate:
    image: ghcr.io/aegisgatesecurity/aegisgate:latest
    ports:
      - "8080:8080"
      - "8443:8443"
    environment:
      - OPENAI_API_KEY=your-api-key
    volumes:
      - ./config:/app/config
```

### Kubernetes (Helm)

```bash
# Add Helm repository
helm repo add aegisgate https://aegisgatesecurity.github.io/helm-charts
helm repo update

# Deploy
helm install aegisgate aegisgate/aegisgate \
  --namespace aegisgate \
  --create-namespace
```

### Build from Source

```bash
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
go build -o aegisgate ./cmd/aegisgate
./aegisgate serve --config ./config/aegisgate.env
```

---

## 🏗️ Architecture

```mermaid
flowchart TB
    subgraph Client["Your Application"]
        A[AI API Request]
    end

    subgraph AegisGate["AegisGate™ Security Gateway"]
        direction LR
        P[Proxy Layer] --> S[Security Engine]
        S --> C[Compliance Engine]
        C --> M[Metrics & Logging]
    end

    subgraph Providers["AI Providers"]
        O[OpenAI]
        AC[Anthropic Claude]
        AZ[Azure OpenAI]
        AWS[AWS Bedrock]
        G[Google Gemini]
    end

    A --> P
    M --> O
    M --> AC
    M --> AZ
    M --> AWS
    M --> G

    style AegisGate fill:#e1f5ff,stroke:#0066cc,stroke-width:3px
    style S fill:#ffeb3b,stroke:#f57f17
    style C fill:#c8e6c9,stroke:#2e7d32
    style P fill:#fff3e0,stroke:#ef6c00
    style M fill:#f3e5f5,stroke:#7b1fa2
```

### Key Components

| Component | Function |
|-----------|----------|
| **Proxy Layer** | Transparent HTTP/HTTPS proxy with TLS termination |
| **Security Engine** | OWASP LLM Top 10, MITRE ATLAS threat detection, prompt injection prevention |
| **Compliance Engine** | Multi-framework policy enforcement (SOC2, HIPAA, GDPR, ISO) |
| **Metrics & Logging** | Prometheus metrics, OpenTelemetry tracing, audit logs |

---

## 📊 Performance Benchmarks

| Metric | Value |
|--------|-------|
| **Latency Overhead** | < 5ms p99 |
| **Throughput** | 50,000 req/s |
| **Docker Image Size** | 27 MB |
| **Memory Footprint** | ~4 MB idle |
| **CPU Usage** | < 2% at 1K req/s |
| **Startup Time** | < 2 seconds |

### Resource Requirements

| Resource | Minimum | Recommended |
|----------|---------|-------------|
| RAM | 256 MB | 512 MB |
| CPU | 1 core | 2+ cores |
| Disk | 100 MB | 1 GB |

---

## 🔒 Security

### Certified Zero Vulnerabilities

| Scan Type | Result | Details |
|-----------|:------:|---------|
| **CVE Scanner** | ✅ 0 CVEs | govulncheck clean |
| **Go 1.25.8** | ✅ Current | All Go stdlib CVEs resolved |
| **Dependency Audit** | ✅ Clean | No vulnerable imports called |
| **Static Analysis** | ✅ Passing | gosec, semgrep clean |

> 🔒 **Zero CVEs, Zero Compromise**: AegisGate™ runs on Go 1.25.8 with zero known vulnerabilities. No patches needed, no security debt to manage — deploy with confidence from day one.

### Defense in Depth

| Layer | Technologies |
|-------|---------------|
| **Transport** | TLS 1.3, mTLS, HTTP/2, HTTP/3 (QUIC) |
| **Authentication** | OAuth 2.0, OIDC, SAML 2.0, API Keys, JWT |
| **Authorization** | RBAC, ABAC, Zero Trust |
| **Data Protection** | AES-256, TLS-in-transit, Key Vault |
| **Runtime** | Seccomp, AppArmor, rootless containers |

### 🐛 Vulnerability Disclosure

Found a security issue? **DO NOT open a public issue.**

- 📧 **Email:** security@aegisgatesecurity.io
- ⏱️ **Response:** Within 48 hours
- 🔧 **Remediation:** 90-day timeline

---

## 🤝 Related Projects

### AegisGuard™ — AI Agent Security Platform

**AegisGuard™** is the companion project for **AI agent security**. While AegisGate secures **API traffic** to LLM providers, AegisGuard secures **agent interactions** through the Model Context Protocol (MCP).

| Aspect | AegisGate™ | AegisGuard™ |
|--------|-------------|-------------|
| **Focus** | API-side security & traffic protection | Agent-side security & compliance |
| **Protocol** | HTTP/HTTPS, HTTP/3 (QUIC) | Model Context Protocol (MCP) |
| **Key Features** | Signature verification, prompt injection protection | RBAC, audit logging, session isolation |
| **Compliance** | SOC 2, HIPAA, PCI-DSS, GDPR, ISO 42001 | GDPR, HIPAA, ISO 27001, NIST AI RMF |

**Together, AegisGate and AegisGuard provide end-to-end AI security — from agent to API.**

👉 **[Learn more about AegisGuard](https://github.com/aegisguardsecurity/aegisguard)**

---

## 💼 Commercial Offerings

AegisGate is available under the Apache 2.0 license with full security and compliance features. Organizations requiring additional capabilities can explore:

| Offering | Highlights |
|----------|------------|
| **Developer** | Multiple AI providers, enhanced threat detection |
| **Professional** | SIEM integration, advanced analytics, priority support |
| **Enterprise** | Custom policies, 99.99% SLA, dedicated support, SSO/SAML |

For details, visit **[aegisgatesecurity.io](https://aegisgatesecurity.io)** or contact sales@aegisgatesecurity.io.

---

## 🤝 Contributing

We welcome contributions! All contributions require signing our CLA.

| Document | Purpose |
|----------|---------|
| [CONTRIBUTING.md](CONTRIBUTING.md) | Development guidelines |
| [CLA.md](CLA.md) | Contributor License Agreement (REQUIRED) |
| [DCO.md](DCO.md) | Developer Certificate of Origin |

### Quick Steps

1. Fork the repository
2. Create a feature branch
3. Sign off all commits: `git commit -s`
4. Submit a Pull Request

📧 Legal questions: **support@aegisgatesecurity.io**

---

## 📝 License

**Apache License 2.0** — Copyright 2025-2026 AegisGate Security. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.

See [LICENSE](LICENSE) for full license text.

> ⚖️ **Legal Note**: All source code is protected by international copyright law, patent rights, and trademark registration.

---

<div align="center">

### 🖤 Love AegisGate™?

**[⭐ Star us on GitHub](https://github.com/aegisgatesecurity/aegisgate/stargazers)**
| **[📢 Share with your team](https://github.com/aegisgatesecurity/aegisgate/discussions)**
| **[❤️ Become a sponsor](https://github.com/sponsors/aegisgatesecurity)**

---

**Built with 🔐 by the [AegisGate™ Security Team](https://github.com/aegisgatesecurity)**

*Enterprise AI Protection — Simplified.*

> 🚀 **What Are You Waiting For?** Start protecting your AI traffic in under 60 seconds. [Deploy Now →](#-quick-start)

</div>
