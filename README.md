# AegisGate - Enterprise AI API Security Platform

**Enterprise-grade security for AI API gateways**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://golang.org/)
[![Version](https://img.shields.io/badge/version-v1.0.3-green.svg)](https://github.com/aegisgatesecurity/aegisgate/releases)

## Overview

AegisGate is a comprehensive security platform for AI API gateways, providing real-time threat detection, compliance monitoring, and ML-powered behavioral analytics.

### Key Capabilities

- AI-Native Security: OpenAI, Anthropic, Azure OpenAI, AWS Bedrock support
- Compliance: MITRE ATLAS, OWASP AI Top 10, HIPAA, PCI-DSS, SOC 2, GDPR, ISO 27001
- ML-Powered Detection: Real-time anomaly detection, prompt injection prevention
- Enterprise SSO: SAML, OIDC, OAuth 2.0, LDAP
- Observability: Splunk, Elastic, Datadog, QRadar integration

## Quick Start

### Docker Compose (Recommended)

```bash
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
docker-compose -f deploy/docker/docker-compose.yml up -d
```

### Kubernetes (Helm)

```bash
helm install aegisgate ./deploy/helm/aegisgate --namespace aegisgate --create-namespace
```

### Binary Installation

```bash
curl -L https://github.com/aegisgatesecurity/aegisgate/releases/download/v1.0.3/aegisgate_v1.0.3.tar.gz | tar xz
./aegisgate --version
```

## Tiers and Licensing

### Free: Community
- 200 requests/min, 3 users, basic features
- No license required

### Paid Tiers
| Tier | Price | Features |
|------|-------|----------|
| Developer | $29/mo | 1K req/min, 10 users, 4 AI providers |
| Professional | $99/mo | 5K req/min, 25 users, all frameworks |
| Enterprise | Custom | Unlimited, air-gapped, HSM support |

### License Format
Cryptographically signed licenses:
- Developer/Professional: base64(JSON).base64(HMAC-SHA256)
- Enterprise: base64(JSON).base64(RSA-SIGNATURE) + hardware binding

## Security

### Security-First Design
- TLS 1.3, mTLS, HTTP/2, HTTP/3
- OAuth 2.0, OIDC, SAML 2.0, LDAP
- RBAC, ABAC with fine-grained permissions
- Encrypted data at rest and in transit

### Compliance Frameworks
Complete support for: OWASP AI Top 10, MITRE ATLAS, SOC 2, HIPAA, PCI-DSS, GDPR, ISO 27001, ISO 42001, NIST AI RMF

### Reporting Vulnerabilities
Email security@aegisgate.io - DO NOT open public issues.

## Documentation

- [Quick Start](docs/quickstart.md)
- [Architecture](docs/architecture.md)
- [Configuration](docs/configuration.md)
- [Security Guide](SECURITY.md)
- [Migration Guide](MIGRATION.md)

## Contributing

See [Contributing Guide](CONTRIBUTING.md) for details.

```bash
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
go mod download
go test ./...
```

## Project Statistics

- **246 Files**, **94K+ Lines**, **99% Go**
- **3,900+ Functions**, **1,050+ Types**

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.

Enterprise features require commercial license. Contact sales@aegisgate.io.

## Support

- Docs: https://docs.aegisgate.io
- Discord: https://discord.gg/aegisgate
- Email: support@aegisgate.io
- Twitter: [@AegisGateIO](https://twitter.com/AegisGateIO)

---

**Made by the AegisGate Team**
