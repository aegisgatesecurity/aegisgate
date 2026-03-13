# AegisGate - Enterprise AI API Security Platform

<div align="center">

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://docker.com/)

**Enterprise-grade security platform for AI API gateways**

[AegisGate](#features) - [Quick Start](#quick-start) - [Tiers and Pricing](#tiers--pricing) - [Contributing](#contributing)

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

## What is AegisGate?

AegisGate is an enterprise-grade AI API security platform that provides comprehensive protection for organizations deploying AI services at scale. Acting as a secure proxy gateway, AegisGate monitors, filters, and secures all traffic between your applications and AI providers like OpenAI, Anthropic, Azure OpenAI, AWS Bedrock, and Cohere.

### Key Capabilities

- **AI API Proxy** - Transparent proxy with full request/response inspection
- **Security Scanning** - Prompt injection, PII detection, malicious payload blocking
- **Observability** - Prometheus metrics, structured logging, dashboard
- **Compliance** - SOC2, GDPR, HIPAA, PCI-DSS, OWASP, NIST, ISO 27001 ready
- **ML Anomaly Detection** - Traffic pattern analysis, cost anomaly detection
- **Authentication** - JWT validation, API key management, RBAC

---

## Features

### Core Security

| Feature | Description |
|---------|-------------|
| AI Provider Proxy | Unified proxy for OpenAI, Anthropic, Azure, Bedrock, Cohere |
| Prompt Injection Detection | Block malicious prompt injection attacks |
| PII Redaction | Automatic detection and redaction of sensitive data |
| Content Filtering | Scan requests/responses for policy violations |
| Rate Limiting | Token bucket algorithm with tier-based limits |

### Enterprise Features

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| AI Providers | 3 | 5 | Unlimited | Unlimited |
| Rate Limit (req/min) | 100 | 1000 | 10000 | Unlimited |
| Max Users | 5 | 50 | 500 | Unlimited |
| ML Detection | No | Yes | Yes | Yes |
| Compliance Reports | No | No | Yes | Yes |
| SSO/SAML | No | Yes | Yes | Yes |
| Multi-Tenancy | No | No | Yes | Yes |
| Behavioral Analysis | No | No | No | Yes |

---

## Quick Start

### Prerequisites

| Requirement | Minimum Version |
|-------------|-----------------|
| Go | 1.21+ |
| Docker | Latest |

### Docker

docker-compose up -d

Access the dashboard at: http://localhost:8080

### From Source

git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
make build
./bin/aegisgate -tier community

### Configuration

export AEGISGATE_BIND=0.0.0.0:8080
export AEGISGATE_TARGET=https://api.openai.com
export AEGISGATE_TIER=community

---

## Architecture

AegisGate acts as a secure proxy between your applications and AI providers:

- **Proxy Layer** - HTTP/HTTPS proxy with request/response inspection
- **Security Scanner** - Multi-layer threat detection (prompt injection, PII, malware)
- **ML Engine** - Anomaly detection and behavioral analysis
- **Compliance Engine** - Policy enforcement for SOC2, GDPR, HIPAA, etc.
- **Metrics** - Prometheus-compatible metrics export
- **Dashboard** - Web-based administration interface

---

## Kubernetes Deployment

AegisGate v1.0.3+ includes production-ready Kubernetes manifests and Helm charts for enterprise deployments.

### Included Kubernetes Resources

| Resource | Description |
|----------|-------------|
| **Deployment** | Production-ready deployment with resource limits |
| **Service** | ClusterIP service for internal access |
| **ServiceAccount** | Kubernetes service account for pod |
| **ConfigMap** | Configuration management |
| **Secret** | Secure credential storage (mount secrets as volumes) |
| **HorizontalPodAutoscaler (HPA)** | Auto-scaling based on CPU/memory metrics |
| **PodDisruptionBudget (PDB)** | Ensure minimum availability during disruptions |
| **NetworkPolicy** | Pod-level network security controls |

### Deploy to Kubernetes

# Apply Kubernetes manifests
kubectl apply -f deploy/k8s/

# Or use Helm (recommended)
helm install aegisgate ./deploy/helm/aegisgate

### Helm Charts

AegisGate includes two production-ready Helm charts:

#### Basic Chart (aegisgate)

For standard deployments with core security features:

helm install aegisgate ./deploy/helm/aegisgate

#### ML-Enabled Chart (aegisgate-ml)

For deployments requiring ML anomaly detection:

helm install aegisgate-ml ./deploy/helm/aegisgate-ml

| Chart | Purpose | Resources |
|-------|---------|-----------|
| aegisgate | Core proxy + security | Deployment, Service, ConfigMap, ServiceAccount |
| aegisgate-ml | + ML anomaly detection | All basic + HPA, Ingress, PVC, NetworkPolicy |

### Helm Chart Validation Results (v1.0.3)

| Test | Status |
|------|--------|
| helm lint | PASS (0 failures) |
| helm template | Renders successfully |
| helm install --dry-run | Passes API validation |
| Live Kubernetes deployment | Running |
| Version alignment | 1.0.3 across all |

---

## Release Notes

### v1.0.4 (Current)

**Release Date:** March 12, 2026

#### New Features
- **Production-Ready Helm Charts**: Full Kubernetes deployment support with:
  - Basic chart (aegisgate) for core security
  - ML-enabled chart (aegisgate-ml) for anomaly detection
  - Both charts validated and tested in live Kubernetes

#### Improvements
- **Helm Chart**: Fixed template naming (padlock.* to aegisgate.*), array syntax, missing templates
- **Helm Chart**: Fixed version alignment (0.1.0/0.2.0 to 1.0.3)
- **Docker**: Production image built with Go 1.24.0, distroless base, 45.6MB

#### Bug Fixes
- Release workflow now uses anchore/sbom-action/download-syft@v0 with dynamic path resolution
- Added missing serviceaccount.yaml to basic chart
- Added missing service.yaml to basic chart
- Added missing configmap.yaml to basic chart
- Fixed array syntax in ML chart ingress template

---

### v1.0.3

**Release Date:** March 12, 2026

#### New Features
- **Kubernetes Production Manifests**: Added full K8s deployment with:
  - HorizontalPodAutoscaler (HPA) for auto-scaling
  - RBAC (ServiceAccount, ClusterRole, ClusterRoleBinding)
  - PodDisruptionBudget for high availability
  - NetworkPolicy for pod-level security
  - ConfigMap for configuration management

---

### v1.0.2

- Initial stable release
- Core proxy functionality
- Security scanning
- Tier-based access control

---

## Tiers and Pricing

| Tier | Price | Best For |
|------|-------|----------|
| Community | Free | Individuals, learning |
| Developer | $29/mo | Startups |
| Professional | $99/mo | Teams, businesses |
| Enterprise | Custom | Large organizations |

For pricing details, please email sales@aegisgatesecurity.io

---

## Security

Report vulnerabilities to security@aegisgatesecurity.io

See SECURITY.md for full disclosure guidelines.

---

## Documentation

| Document | Description |
|----------|-------------|
| CHANGELOG.md | Release history and changes |
| CONTRIBUTING.md | Contribution guidelines |
| CODE_OF_CONDUCT.md | Community code of conduct |
| SECURITY.md | Security policy and supported versions |
| TODO.md | Future roadmap and feature planning |

### API Reference

# Health check
curl http://localhost:8080/health

# Metrics
curl http://localhost:9090/metrics

# Version info
curl http://localhost:8080/version

---

## Contributing

We welcome contributions! Please read CONTRIBUTING.md before submitting PRs.

git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
git checkout -b feature/your-feature
make test
git commit -m "Add your feature"
git push origin feature/your-feature

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
