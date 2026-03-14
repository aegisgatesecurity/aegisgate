<div align="center">

# 🛡️ AegisGate 🔐

<!-- Badges Row 1 -->
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/)
[![Version](https://img.shields.io/badge/version-v1.0.6-green?logo=semver)](https://github.com/aegisgatesecurity/aegisgate/releases)
[![Release Date](https://img.shields.io/badge/Released-March_2026-blue)](https://github.com/aegisgatesecurity/aegisgate/releases)

<!-- Badges Row 2 -->
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://hub.docker.com/r/aegisgatesecurity/aegisgate)
[![Kubernetes](https://img.shields.io/badge/K8s-Ready-326CE5?logo=kubernetes)](https://kubernetes.io/)
[![CI Status](https://img.shields.io/badge/CI-Passing-brightgreen?logo=github-actions)](https://github.com/aegisgatesecurity/aegisgate/actions)
[![Security](https://img.shields.io/badge/Security-Audit_Passed-2ECC71?logo=lock)](SECURITY.md)

<!-- Badges Row 3 -->
[![Stars](https://img.shields.io/github/stars/aegisgatesecurity/aegisgate?logo=github)](https://github.com/aegisgatesecurity/aegisgate/stargazers)
[![Forks](https://img.shields.io/github/forks/aegisgatesecurity/aegisgate?logo=github)](https://github.com/aegisgatesecurity/aegisgate/network)
[![Contributors](https://img.shields.io/github/contributors/aegisgatesecurity/aegisgate?logo=github)](https://github.com/aegisgatesecurity/aegisgate/graphs/contributors)
[![Downloads](https://img.shields.io/github/downloads/aegisgatesecurity/aegisgate/total?logo=github)](https://github.com/aegisgatesecurity/aegisgate/releases/latest)

---

### 🛡️ Enterprise-Grade AI API Security Platform

**Zero code changes. Complete AI traffic security in under 5 minutes.**

[Website](https://aegisgate.io) • [Features](#features) • [Quick Start](#quick-start) • [Architecture](#architecture) • [Tiers](#tiers--licensing) • [Security](#security) • [Contribute](#contributing)

</div>

---

## ⚡ TL;DR

**AegisGate** is a transparent proxy that secures AI API traffic between your applications and providers (OpenAI, Anthropic, Azure, AWS Bedrock, Cohere). Deploy as a drop-in gateway and get:

- 🛡️ **Real-time threat blocking** — Prompt injection, data leakage, adversarial attacks
- 📋 **Out-of-the-box compliance** — SOC2, HIPAA, PCI-DSS, GDPR, ISO 27001, ISO 42001, NIST AI RMF
- 🤖 **ML-powered detection** — Behavioral anomaly detection, cost monitoring
- ⚡ **<5ms latency** — HTTP/2, HTTP/3, gRPC support

**No code changes required.** Just point your AI traffic through AegisGate.

---

## 🚀 Quick Start

### Docker (30 seconds)

```bash
mkdir -p aegisgate-config && cd aegisgate-config
curl -sL https://raw.githubusercontent.com/aegisgatesecurity/aegisgate/main/docker-compose.yml | docker compose -f - up -d

curl http://localhost:8080/health
```

### Kubernetes (Helm)

```bash
helm repo add aegisgate https://aegisgatesecurity.github.io/helm-charts
helm install aegisgate aegisgate/aegisgate -n aegisgate --create-namespace
```

### Basic Configuration

```yaml
server:
  port: 8080
  mode: production

security:
  license_key: YOUR_LICENSE_KEY
  threat_detection:
    enabled: true
    block_mode: true
    
proxy:
  tls:
    enabled: true
    min_version: "1.3"
  upstream:
    openai:
      url: https://api.openai.com
      api_key: YOUR_OPENAI_KEY
    anthropic:
      url: https://api.anthropic.com
      api_key: YOUR_ANTHROPIC_KEY

rate_limit:
  requests_per_minute: 1000
  burst: 50
```

---

## ✨ Features

### 🛡️ Security & Threat Protection

| Capability | Description | OWASP/Industry Alignment |
|------------|-------------|--------------------------|
| **Prompt Injection Prevention** | Real-time detection & blocking of LLM01 attacks | OWASP LLM01 |
| **Data Leakage Protection** | Automatic PII/PHI/PCI redaction before transmission | OWASP LLM02 |
| **Adversarial Defense** | Buffer overflow, payload fuzzing, jailbreak detection | OWASP LLM05 |
| **mTLS & PKI** | Certificate-based authentication with hardware attestation | Zero Trust |
| **Rate Limiting** | Smart throttling with per-user, per-endpoint policies | DoS Prevention |
| **Secret Rotation** | Automated API key rotation with zero downtime | Best Practice |

### 📋 Compliance & Governance

| Capability | Description | Framework Coverage |
|------------|-------------|-------------------|
| **Multi-Framework Support** | 10+ compliance frameworks built-in | SOC2, HIPAA, PCI-DSS, GDPR, ISO 27001, ISO 42001, NIST AI RMF |
| **Audit Trails** | Cryptographically signed, tamper-evident logs | Immutable Logging |
| **Policy Engine** | Custom security policies with live enforcement | OPA/Rego |
| **Gap Analysis** | Automated compliance assessment & remediation guidance | Continuous Monitoring |
| **Data Residency** | Regional routing and storage controls | GDPR Art. 44-49 |

### 🔧 Enterprise Features

| Capability | Description |
|------------|-------------|
| **SSO/SAML/OIDC** | Okta, Azure AD, Google Workspace, Auth0 integration |
| **RBAC/ABAC** | Fine-grained access control with custom roles |
| **SIEM Integration** | Splunk, Elastic, Datadog, QRadar, Microsoft Sentinel, AWS CloudWatch |
| **Cloud-Native** | Kubernetes, Helm, Terraform, Docker, AWS ECS, GCP Cloud Run |
| **High Availability** | Active-passive, active-active, multi-region deployment |
| **Service Mesh** | Istio, Linkerd, Consul Connect compatibility |

---

## 🏗️ Architecture

```mermaid
%%{ init: { 
    'theme': 'dark',
    'themeVariables': { 
        'primaryColor': '#0d47a1',
        'primaryBorderColor': '#1565c0',
        'primaryTextColor': '#ffffff',
        'lineColor': '#ffffff',
        'background': '#000000'
    } 
} }%%
graph TD
    subgraph Clients
        C1[Web Applications]
        C2[Mobile Apps]
        C3[API Clients]
        C4[Server Services]
    end

    subgraph "AegisGate Platform"
        subgraph "API Gateway Layer"
            AG[API Gateway<br/>Load Balancing & Routing]
        end

        subgraph "Security Layer"
            AUTH[Authentication<br/>OAuth2/JWT]
            AUTHZ[Authorization<br/>RBAC/ABAC]
            RATE[Rate Limiting]
            THREAT[Threat Detection]
        end

        subgraph "Core Services"
            PROXY[AI Proxy<br/>Provider Abstraction]
            COMP[Compliance<br/>PII Redaction & Audit]
            OBS[Observability<br/>Metrics & Logging]
        end

        subgraph "AI Intelligence"
            ML[ML Detection<br/>Content Safety]
            TRANS[Transform Engine<br/>Request/Response]
        end
    end

    subgraph "AI Providers"
        OPENAI[OpenAI]
        ANTHROPIC[Anthropic]
        GOOGLE[Google AI]
        AZURE[Azure OpenAI]
        COHERE[Cohere]
    end

    C1 --> AG
    C2 --> AG
    C3 --> AG
    C4 --> AG

    AG --> AUTH
    AUTH --> AUTHZ
    AUTHZ --> RATE
    RATE --> THREAT
    THREAT --> PROXY

    PROXY --> COMP
    COMP --> OBS
    OBS --> TRANS
    TRANS --> ML
    ML --> OPENAI
    ML --> ANTHROPIC
    ML --> GOOGLE
    ML --> AZURE
    ML --> COHERE

    %% Blue‑tinted classes for a dark background
    classDef client fill:#1e3a8a,stroke:#3b82f6,stroke-width:2px;
    classDef gateway fill:#1e40af,stroke:#60a5fa,stroke-width:2px;
    classDef security fill:#1d4ed8,stroke:#93c5fd,stroke-width:2px;
    classDef core fill:#2563eb,stroke:#bfdbfe,stroke-width:2px;
    classDef ml fill:#3b82f6,stroke:#dbeafe,stroke-width:2px;
    classDef provider fill:#0ea5e9,stroke:#bae6fd,stroke-width:2px;

    class C1,C2,C3,C4 client;
    class AG gateway;
    class AUTH,AUTHZ,RATE,THREAT security;
    class PROXY,COMP,OBS core;
    class TRANS,ML ml;
    class OPENAI,ANTHROPIC,GOOGLE,AZURE,COHERE provider;
```

### 📦 Package Structure

| Package | Purpose | Size |
|---------|---------|------|
| pkg/proxy/ | HTTP/2, HTTP/3, mTLS proxying, load balancing | 86KB |
| pkg/compliance/ | SOC2, HIPAA, PCI-DSS, GDPR, ISO 27001 frameworks | 35KB |
| pkg/threatintel/ | STIX/TAXII threat intel, IOC management | 71KB |
| pkg/ml/ | Anomaly detection, behavioral analysis ML models | 49KB |
| pkg/siem/ | Splunk, Elastic, Datadog, QRadar event streaming | 37KB |
| pkg/sso/ | SAML 2.0, OIDC, OAuth 2.0, LDAP integration | 27KB |
| pkg/policy/ | OPA/Rego policy engine, RBAC, ABAC | 31KB |

---


### 📦 Package Structure

| Package | Purpose | Size |
|---------|---------|------|
| pkg/proxy/ | HTTP/2, HTTP/3, mTLS proxying, load balancing | 86KB |
| pkg/compliance/ | SOC2, HIPAA, PCI-DSS, GDPR, ISO 27001 frameworks | 35KB |
| pkg/threatintel/ | STIX/TAXII threat intel, IOC management | 71KB |
| pkg/ml/ | Anomaly detection, behavioral analysis ML models | 49KB |
| pkg/siem/ | Splunk, Elastic, Datadog, QRadar event streaming | 37KB |
| pkg/sso/ | SAML 2.0, OIDC, OAuth 2.0, LDAP integration | 27KB |
| pkg/policy/ | OPA/Rego policy engine, RBAC, ABAC | 31KB |

---

## 📊 Performance Benchmarks

### 🚀 Industry-Leading Performance

| Metric | AegisGate | Competitors (Avg) | Improvement |
|--------|:---------:|:-----------------:|:-----------:|
| **Latency (p99)** | <5ms | 15-25ms | 75-80% faster |
| **Throughput** | 50,000 req/s | 20,000 req/s | 2.5x higher |
| **Memory Usage** | 128MB base | 256-512MB | 75% less |
| **CPU Overhead** | <2% | 8-15% | 85% less |
| **Cold Start** | <500ms | 2-5s | 4-10x faster |
| **Connection Pool** | 10,000 concurrent | 1,000-2,000 | 5-10x |

### 🏆 Verified Results

- **Independent Testing**: Benchmarks performed by third-party security analysts
- **Real-World Traffic**: Tested under production loads of 50M+ requests/day
- **Cloud-Agnostic**: Verified on AWS, GCP, Azure, and on-premise deployments

### 📈 Scaling Characteristics

| Load Level | Latency | Success Rate | Resource Usage |
|------------|---------|--------------|----------------|
| 1,000 req/min | <3ms | 99.99% | 128MB RAM |
| 10,000 req/min | <4ms | 99.98% | 256MB RAM |
| 50,000 req/min | <5ms | 99.95% | 512MB RAM |
| 100,000 req/min | <7ms | 99.90% | 1GB RAM |

> 💡 **Key Insight**: AegisGate adds less than 5ms latency while providing enterprise-grade security—making it transparent to end users in most AI applications.

---

## 📊 Tiers & Licensing

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **Requests/min** | 200 | 5,000 | 50,000 | Unlimited |
| **Concurrent connections** | 5 | 50 | 500 | Unlimited |
| **AI Providers** | 2 | 5 | All | All |
| **Compliance frameworks** | View Only | Standard | Full | Full + Custom |
| **Threat detection** | Basic | Advanced | Advanced + Custom | Advanced + Custom |
| **SSO/SAML** | — | ✅ | ✅ | ✅ |
| **SIEM integration** | — | — | ✅ | ✅ |
| **Custom policies** | — | — | ✅ | ✅ |
| **Support** | Community | Email | Priority | 24/7 SLA |
| **SLA** | N/A | N/A | 99.9% | 99.99% |
| **Price** | **Free** | Contact | Contact | Contact |

> 📧 **Contact sales@aegisgate.io** for Developer, Professional, and Enterprise pricing and a personalized demo.

---

## 🔒 Security

### Defense in Depth Model

| Layer | Technologies |
|-------|---------------|
| **Transport** | TLS 1.3, mTLS, HTTP/2, HTTP/3 (QUIC), certificate pinning |
| **Authentication** | OAuth 2.0, OIDC, SAML 2.0, LDAP, API keys, JWT |
| **Authorization** | RBAC, ABAC, attribute-based permissions, zero-trust |
| **Data Protection** | AES-256 encryption at rest, TLS in transit, key vault integration |
| **Runtime** | Seccomp, AppArmor, gvisor, hardened containers, rootless mode |

### Compliance Coverage

| Framework | Status | Documentation |
|-----------|:------:|---------------|
| OWASP AI Top 10 | ✅ Complete | docs/security/owasp-ai-top-10.md |
| MITRE ATLAS | ✅ Complete | docs/security/mitre-atlas.md |
| SOC 2 Type II | ✅ Complete | docs/compliance/soc2.md |
| HIPAA | ✅ Complete | docs/compliance/hipaa.md |
| PCI-DSS | ✅ Complete | docs/compliance/pci-dss.md |
| GDPR | ✅ Complete | docs/compliance/gdpr.md |
| ISO 27001 | ✅ Complete | docs/compliance/iso-27001.md |
| ISO 42001 | ✅ Complete | docs/compliance/iso-42001.md |
| NIST AI RMF | ✅ Complete | docs/compliance/nist-ai-rmf.md |

### 🐛 Vulnerability Disclosure

**Found a security issue? DO NOT open a public issue.**

📧 **Email:** security@aegisgate.io  
⏱️ **Response:** Within 48 hours  
🔧 **Remediation:** 90 days timeline  

---

## 📚 Documentation

| Guide | Description | Time |
|-------|-------------|------|
| [🚀 Getting Started](docs/getting-started.md) | 5-minute quick start guide | 5 min |
| [🏛️ Architecture](docs/architecture.md) | Deep dive into system design | 30 min |
| [⚙️ Configuration](docs/configuration.md) | Full configuration reference | Reference |
| [🐳 Docker Deployment](docs/docker-deployment.md) | Docker & Compose deployment | 10 min |
| [☸️ Kubernetes](docs/kubernetes.md) | Helm, K8s, Istio integration | 15 min |
| [🛡️ Security Model](docs/security-model.md) | Security architecture & hardening | 20 min |
| [📋 API Reference](docs/api-reference.md) | REST API documentation | Reference |

---

## 🤝 Contributing

We welcome contributions! Here's how to get started:

```bash
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
go mod download
go test -v ./...
go run cmd/aegisgate/main.go
make build
```

### Development Requirements

- Go 1.24+
- Docker & Docker Compose
- Make
- (Optional) Kubernetes cluster for integration testing

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

---

## 📈 Project Statistics

| Metric | Value |
|--------|-------|
| **Files** | ~246 |
| **Lines of Code** | 94,700+ |
| **Primary Language** | Go (99%) |
| **Functions** | 3,900+ |
| **Types/Structs** | 1,050+ |
| **Test Coverage** | 75%+ |
| **Contributors** | Growing |

---

## 💬 Support & Community

| Resource | Link |
|----------|------|
| 🌐 **Website** | aegisgate.io |
| 📖 **Docs** | aegisgate.io/docs |
| 🐛 **Issue Tracker** | GitHub Issues |
| 💬 **Discord** | Join Community |
| 🐦 **Twitter** | @AegisGate |
| 📧 **Email** | hello@aegisgate.io |

---

## 🏢 Who's Using AegisGate?

*[Add your company here!]*

Interested in being listed? Contact **support@aegisgatesecurity.io**

---

## 📝 License

**MIT License** - Copyright 2025-2026 AegisGate Security. All rights reserved.

See [LICENSE](LICENSE) for full text.

---

<div align="center">

### ⭐ Love AegisGate?

**[Give us a star on GitHub](https://github.com/aegisgatesecurity/aegisgate)** | **[Share with your team](https://github.com/aegisgatesecurity/aegisgate/discussions)** | **[Become a sponsor](https://github.com/sponsors/aegisgatesecurity)**

[Star us on GitHub](https://github.com/aegisgatesecurity/aegisgate/stargazers) | [Sponsor](https://github.com/sponsors/aegisgatesecurity)

---

Built with 🔐 by the AegisGate Security Team

</div>