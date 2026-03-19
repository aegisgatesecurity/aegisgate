# Frequently Asked Questions

---

## General Questions

### What is AegisGate?

AegisGate is an enterprise-grade security platform that provides comprehensive protection for AI API infrastructure. It offers API security, compliance automation, observability, and ML-powered threat detection.

### Why do I need AegisGate?

If you're using AI APIs (OpenAI, Anthropic, etc.) in production, you need:
- **Security** - Protect against attacks and unauthorized access
- **Compliance** - Meet regulatory requirements (GDPR, SOC 2, OWASP)
- **Observability** - Understand your AI usage patterns
- **Cost Control** - Monitor and control AI spending

### Is AegisGate a SaaS product?

No. AegisGate is **self-hosted** software. You run it on your own infrastructure (Docker, Kubernetes, bare metal). Your data never leaves your infrastructure.

---

## Technical Questions

### What languages/frameworks does AegisGate support?

AegisGate is language-agnostic because it works as a proxy. Any application that can make HTTP requests can use AegisGate.

| Integration | Support |
|-------------|---------|
| **REST API** | ✅ Full |
| **gRPC** | ✅ Available |
| **Webhooks** | ✅ Available |
| **SDK (Go, Python, JS)** | ✅ Available |

### Which AI providers does AegisGate support?

| Provider | Support |
|----------|---------|
| OpenAI | ✅ |
| Anthropic | ✅ |
| Cohere | ✅ |
| Azure OpenAI | ✅ |
| AWS Bedrock | ✅ |
| Google Vertex | ✅ |
| Custom/Internal | ✅ |

### What compliance frameworks does AegisGate support?

| Framework | Description | Availability |
|-----------|-------------|--------------|
| OWASP Top 10 | Security vulnerabilities | Community |
| SOC 2 | Service organization control | Community |
| GDPR | Data protection | Community |
| ISO 27001 | Information security | Community |
| NIST AI RMF | AI risk management | Enterprise |
| HIPAA | Healthcare compliance | Enterprise |
| PCI-DSS | Payment card security | Enterprise |
| ISO 42001 | AI-specific standard | Enterprise |
| FedRAMP | US government security | Enterprise |

> 💡 **Enterprise Compliance**: Additional compliance frameworks (HIPAA, PCI-DSS, ISO 42001, NIST AI RMF) available with Enterprise license. Contact sales@aegisgatesecurity.io

### Can I run AegisGate on my own infrastructure?

Yes! AegisGate is designed for self-hosting.

| Deployment | Support |
|------------|---------|
| Docker | ✅ |
| Kubernetes | ✅ |
| Bare Metal | ✅ |
| Air-gapped | ✅ |

### What are the system requirements?

| Requirement | Minimum | Recommended |
|-------------|---------|-------------|
| **CPU** | 1 core | 2+ cores |
| **RAM** | 512 MB | 2 GB |
| **Disk** | 5 GB | 20 GB |
| **OS** | Linux/Windows/Mac | Linux |

For production with PostgreSQL:
| Requirement | Minimum | Recommended |
|-------------|---------|-------------|
| **CPU** | 2 cores | 4+ cores |
| **RAM** | 2 GB | 4+ GB |
| **PostgreSQL** | 14+ | 14+ |

---

## Licensing

### How does licensing work?

AegisGate uses license keys to activate different feature tiers:

- **Community**: Base features, no license required
- **Developer**: Additional features with self-hosted license
- **Professional**: Enhanced features with self-hosted license
- **Enterprise**: Full features with custom license and support

### Can I use AegisGate commercially?

Yes! The Apache 2.0 license allows commercial use for the Community edition. Paid tiers include commercial licensing with support.

### What happens if my license expires?

Features may revert to Community tier during any transition period.

---

## Security

### Is my data secure?

Yes! AegisGate is designed with security in mind:

- **Self-hosted**: Data never leaves your infrastructure
- **Encryption**: TLS 1.3, mTLS support
- **Audit Logging**: Comprehensive audit trails
- **No telemetry**: We don't collect your usage data

### Can AegisGate decrypt my traffic?

AegisGate can decrypt traffic for security scanning, but:
- **HTTPS/TLS termination** is optional and configurable
- **mTLS** provides mutual authentication
- You control what's inspected

### Does AegisGate store my API keys?

No. AegisGate proxies requests but doesn't store your AI provider API keys. They're passed through to the AI provider.

---

## Troubleshooting

### AegisGate won't start

```bash
# Check logs
docker logs aegisgate

# Common issues:
# - Port already in use: Check other processes on 8080
# - Permission denied: Check file permissions
# - Missing config: Ensure config file exists
```

### Can't connect to AI provider

```bash
# Test connectivity
docker exec aegisgate curl https://api.openai.com

# Check API key
echo $OPENAI_API_KEY

# Check network
docker network ls
```

### Performance issues

```bash
# Check resources
docker stats aegisgate

# Increase limits
# Edit docker-compose.yml
```

### Rate limiting triggered

You're hitting your tier's rate limits. Options:
1. Wait and retry
2. Upgrade your tier
3. Contact support for adjustments

---

## Getting Help

### How do I get support?

| Tier | Support Channel |
|------|-----------------|
| Community | [Discord](https://discord.gg/aegisgate), [Forum](https://community.aegisgate.example.com) |
| Developer | Email support |
| Professional | Priority email |
| Enterprise | 24/7 dedicated support |

### Where can I find documentation?

- [Getting Started](getting-started.md)
- [Configuration](CONFIGURATION.md)
- [API Reference](docs/API.md)
- [GitHub Issues](https://github.com/aegisgatesecurity/aegisgate/issues)

### How do I report bugs?

See our [Issue Template](.github/ISSUE_TEMPLATE/ISSUE_TEMPLATE.md) or email security@aegisgatesecurity.example.com for security issues.

---

## Contributing

### How can I contribute?

1. **Star the repo** ⭐
2. **Report bugs** 🐛
3. **Request features** 💡
4. **Submit PRs** 📝
5. **Write docs** 📚

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Is there a bug bounty program?

Yes! See [bugbounty.aegisgate.example.com](https://bugbounty.aegisgate.example.com) for details.

---

*Can't find your answer? Contact us at support@aegisgatesecurity.example.com*
