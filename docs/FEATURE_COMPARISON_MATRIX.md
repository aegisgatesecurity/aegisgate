# AegisGate Feature Comparison Matrix

## Tiers Overview

| Tier | Target | Price Point | Best For |
|------|--------|-------------|----------|
| **Community** | Individual developers, hobbyists | Free | Learning, personal projects |
| **Developer** | Individual devs, small projects | $29/month | Freelancers, startups building MVPs |
| **Professional** | Small to medium businesses | $99/month | Production apps, teams |
| **Enterprise** | Large organizations | Custom | Global enterprises, regulated industries |

---

## Adoption Funnel Strategy

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                    ADOPTION FUNNEL                                      │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                          │
│    ┌──────────────┐     ┌──────────────┐     ┌────────────────┐     ┌────────────────┐ │
│    │  COMMUNITY   │────▶│  DEVELOPER   │────▶│  PROFESSIONAL  │────▶│   ENTERPRISE    │ │
│    │    (Free)    │     │   ($29/mo)   │     │   ($99/mo)     │     │    (Custom)     │ │
│    └──────────────┘     └──────────────┘     └────────────────┘     └────────────────┘ │
│         │                     │                     │                     │            │
│         ▼                     ▼                     ▼                     ▼            │
│    • Try AegisGate         • Build MVPs         • Go to Production      • Scale &        │
│    • Learn the ropes    • Test Integrations  • Team Collaboration    • Compliance     │
│    • Spread the Word    • Prototype          • Advanced Features     • Enterprise     │
│                                                                  Support               │
│    Conversion: Free    Conversion: $29     Conversion: $99        Support               │
│    → Developer        → Pro              → Enterprise                                 │
│                                                                                          │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 1. Rate & Resource Limits

| Limit Type | Community | Developer | Professional | Enterprise |
|------------|:---------:|:---------:|:------------:|:----------:|
| **Requests/min** | 200 | 1,000 | 5,000 | Unlimited |
| **Concurrent Connections** | 5 | 25 | 100 | Unlimited |
| **Burst Requests** | 50 | 200 | 500 | Unlimited |
| **Max Users** | 3 | 10 | 25 | Unlimited |
| **Max API Keys** | 3 | 10 | 50 | Unlimited |
| **Max Tenants** | 1 | 1 | 25 | Unlimited |
| **Log Retention** | 1 day | 7 days | 30 days | Unlimited |
| **Data Retention** | 7 days | 30 days | 90 days | Unlimited |
| **Log Size Limit** | 100 MB | 1 GB | 10 GB | Unlimited |

---

## 2. AI Proxy & Connectivity

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **Basic AI Proxy** | ✅ | ✅ | ✅ | ✅ |
| OpenAI Integration | ✅ | ✅ | ✅ | ✅ |
| Anthropic Integration | ✅ | ✅ | ✅ | ✅ |
| Cohere Integration | ❌ | ✅ | ✅ | ✅ |
| Azure OpenAI | ❌ | ✅ | ✅ | ✅ |
| AWS Bedrock | ❌ | ❌ | ✅ | ✅ |
| Google Vertex AI | ❌ | ❌ | ❌ | ✅ |
| **Custom Provider Adapter** | ❌ | ❌ | ✅ | ✅ |
| **Internal AI Tool Support** (Cline, etc.) | ❌ | ❌ | ❌ | ✅ |
| Request Caching | ❌ | ✅ | ✅ | ✅ |
| Request Deduplication | ❌ | ✅ | ✅ | ✅ |
| Batch Processing | ❌ | ❌ | ✅ | ✅ |
| Streaming Support | ✅ | ✅ | ✅ | ✅ |
| Connection Pooling | ❌ | ❌ | ✅ | ✅ |

---

## 3. Security Features

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **TLS/SSL Termination** | ✅ | ✅ | ✅ | ✅ |
| Runtime Hardening | ❌ | ✅ | ✅ | ✅ |
| mTLS Support | ❌ | ✅ | ✅ | ✅ |
| PKI Attestation | ❌ | ❌ | ✅ | ✅ |
| Secret Rotation | ❌ | ❌ | ✅ | ✅ |
| API Key Management | ✅ | ✅ | ✅ | ✅ |
| JWT Validation | ✅ | ✅ | ✅ | ✅ |
| OAuth 2.0 / SSO | ❌ | ✅ | ✅ | ✅ |
| SAML Support | ❌ | ❌ | ✅ | ✅ |
| OIDC Support | ❌ | ✅ | ✅ | ✅ |
| LDAP/AD Integration | ❌ | ❌ | ❌ | ✅ |
| Hardware Token Support (YubiKey) | ❌ | ❌ | ❌ | ✅ |
| HSM Integration | ❌ | ❌ | ❌ | ✅ |
| Audit Log Encryption | ❌ | ❌ | ✅ | ✅ |
| FIPS 140-2 Compliance | ❌ | ❌ | ❌ | ✅ |

---

## 4. Compliance Frameworks

### Community Tier (Free)
| Framework | Status | Notes |
|-----------|--------|-------|
| **OWASP Top 10** | ✅ | Basic checks included |
| **SOC 2** | ⚠️ | View-only (no enforcement) |
| **GDPR** | ⚠️ | View-only (no enforcement) |

### Developer Tier ($29/month)
| Framework | Status | Notes |
|-----------|--------|-------|
| **All Community Frameworks** | ✅ | Included |
| **Basic Security Compliance** | ✅ | Enforced |
| **NIST (View-only)** | ✅ | Basic coverage |
| **Custom Compliance Rules** | ⚠️ | Limited (5 rules) |

### Professional Tier ($99/month)
| Framework | Status | Notes |
|-----------|--------|-------|
| **HIPAA** | ✅ | Full compliance checking |
| **PCI-DSS** | ✅ | Full compliance checking |
| **SOC 2** | ✅ | Full compliance checking |
| **GDPR** | ✅ | Full compliance checking |
| **OWASP** | ✅ | Full compliance checking |
| **NIST 800-53** | ✅ | Basic coverage |
| **ISO 27001** | ✅ | Basic coverage |
| Custom Framework Builder | ✅ | Create custom rules |

### Enterprise Tier (Custom Pricing)
| Framework | Status | Notes |
|-----------|--------|-------|
| **All Professional Frameworks** | ✅ | Included |
| **ISO 42001** | ✅ | AI-specific standard |
| **NIST AI RMF** | ✅ | AI risk management |
| **HITRUST** | ✅ | Healthcare framework |
| **FedRAMP** | ✅ | US Government |
| **SOC 2 Type II** | ✅ | Advanced reporting |
| **Custom Framework API** | ✅ | Programmatic creation |
| **Atlas Framework** | ✅ | MITRE ATT&CK for AI |
| **COBIT** | ✅ | IT governance |
| **NIST CSF** | ✅ | Cybersecurity framework |
| **GLBA** | ✅ | Financial services |
| **SOX** | ✅ | Financial reporting |

---

## 5. Machine Learning & Anomaly Detection

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **Basic Anomaly Detection** | ✅ | ✅ | ✅ | ✅ |
| **Traffic Pattern Analysis** | ✅ | ✅ | ✅ | ✅ |
| **Cost Anomaly Alerts** | ❌ | ✅ | ✅ | ✅ |
| **Usage Anomaly Alerts** | ❌ | ✅ | ✅ | ✅ |
| **Behavioral Analysis** | ❌ | ❌ | ✅ | ✅ |
| **Predictive Cost Modeling** | ❌ | ❌ | ✅ | ✅ |
| **Threat Detection ML** | ❌ | ❌ | ✅ | ✅ |
| **Custom ML Models** | ❌ | ❌ | ❌ | ✅ |
| **Real-time Threat Response** | ❌ | ❌ | ❌ | ✅ |
| **Zero-day Attack Detection** | ❌ | ❌ | ❌ | ✅ |

---

## 6. Multi-Tenancy & Access Control

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **Single Tenant** | ✅ | ✅ | ✅ | ✅ |
| **Multi-Tenancy** | ❌ | ❌ | ✅ | ✅ |
| Role-Based Access Control (RBAC) | ✅ | ✅ | ✅ | ✅ |
| Granular Permissions | ❌ | ✅ | ✅ | ✅ |
| Department/Team Separation | ❌ | ❌ | ✅ | ✅ |
| Custom Roles | ❌ | ✅ | ✅ | ✅ |
| Policy Engine | ❌ | ❌ | ✅ | ✅ |
| Cross-Tenant Analytics | ❌ | ❌ | ❌ | ✅ |
| White-Labeling | ❌ | ❌ | ❌ | ✅ |
| Custom Domains | ❌ | ❌ | ❌ | ✅ |

---

## 7. Observability & Monitoring

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **Basic Metrics** | ✅ | ✅ | ✅ | ✅ |
| **Request Logging** | ✅ | ✅ | ✅ | ✅ |
| **Error Tracking** | ✅ | ✅ | ✅ | ✅ |
| **Dashboard** | Basic | Advanced | Advanced | Custom |
| Grafana Integration | ❌ | ✅ | ✅ | ✅ |
| Datadog Integration | ❌ | ❌ | ✅ | ✅ |
| New Relic Integration | ❌ | ❌ | ✅ | ✅ |
| **SIEM Integration** | ❌ | ❌ | ✅ | ✅ |
| Splunk / Elastic | ❌ | ❌ | ✅ | ✅ |
| AWS CloudWatch | ❌ | ✅ | ✅ | ✅ |
| **Log Retention** | 1 day | 7 days | 30 days | Unlimited |
| **Alert Webhooks** | 3 | 10 | 25 | Unlimited |

---

## 8. API & Integrations

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **REST API** | ✅ | ✅ | ✅ | ✅ |
| **gRPC API** | ❌ | ✅ | ✅ | ✅ |
| GraphQL | ❌ | ❌ | ❌ | ✅ |
| Webhook Triggers | ❌ | ✅ | ✅ | ✅ |
| Terraform Provider | ❌ | ✅ | ✅ | ✅ |
| Kubernetes Operator | ❌ | ❌ | ✅ | ✅ |
| Helm Charts | ❌ | ❌ | ✅ | ✅ |
| SDK (Go, Python, JS) | ❌ | ❌ | ✅ | ✅ |

---

## 9. Storage & Data

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **In-Memory Storage** | ✅ | ✅ | ✅ | ✅ |
| **File-based Storage** | ✅ | ✅ | ✅ | ✅ |
| PostgreSQL | ❌ | ✅ | ✅ | ✅ |
| MySQL/MariaDB | ❌ | ✅ | ✅ | ✅ |
| Redis | ❌ | ✅ | ✅ | ✅ |
| S3-compatible Storage | ❌ | ❌ | ✅ | ✅ |
| Data Encryption at Rest | ❌ | ✅ | ✅ | ✅ |
| Data Retention Policies | ❌ | ❌ | ✅ | ✅ |

---

## 10. Deployment & Infrastructure

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **Docker** | ✅ | ✅ | ✅ | ✅ |
| **Docker Compose** | ✅ | ✅ | ✅ | ✅ |
| Kubernetes | ❌ | ❌ | ✅ | ✅ |
| Terraform | ❌ | ✅ | ✅ | ✅ |
| Helm Charts | ❌ | ❌ | ✅ | ✅ |
| Service Mesh | ❌ | ❌ | ❌ | ✅ |
| Auto-scaling | ❌ | ❌ | ❌ | ✅ |
| On-Premise | ❌ | ❌ | ❌ | ✅ |
| Air-gapped Installation | ❌ | ❌ | ❌ | ✅ |

---

## 11. Support & SLA

| Feature | Community | Developer | Professional | Enterprise |
|---------|:---------:|:---------:|:------------:|:----------:|
| **Documentation** | ✅ | ✅ | ✅ | ✅ |
| **Community Forum** | ✅ | ✅ | ✅ | ✅ |
| **Knowledge Base** | ❌ | ✅ | ✅ | ✅ |
| **Email Support** | ❌ | ✅ | ✅ | ✅ |
| Priority Support | ❌ | ❌ | ✅ | ✅ |
| Dedicated Support | ❌ | ❌ | ❌ | ✅ |
| 24/7 Support | ❌ | ❌ | ❌ | ✅ |
| **SLA** | None | 99.5% | 99.5% | 99.99% |
| Response Time | N/A | 48 hours | 24 hours | 1 hour |

---

## Quick Comparison Summary

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                   TIER AT A GLANCE                                          │
├──────────────────────┬──────────────┬──────────────┬─────────────────┬────────────────────┤
│                      │  COMMUNITY   │  DEVELOPER   │  PROFESSIONAL  │     ENTERPRISE     │
├──────────────────────┼──────────────┼──────────────┼─────────────────┼────────────────────┤
│ Price                │    FREE      │   $29/mo     │    $99/month    │      Custom        │
│ Requests/min         │     200      │    1,000     │     5,000       │      Unlimited     │
│ Concurrent           │       5      │      25      │       100       │      Unlimited     │
│ AI Providers         │      2       │      4       │        6        │        All         │
│ Compliance           │   OWASP      │  OWASP+NIST  │     10+         │       All + ISO     │
│ Database             │   File       │   Postgres+  │    Postgres+    │      Custom         │
│ Tenancy              │   Single     │    Single    │    Multi        │      Multi + Custom│
│ Support              │   Forum      │    Email     │    Priority     │    24/7 Dedicated   │
│ SLA                  │     --       │    99.5%     │    99.5%        │       99.99%        │
│ Max Users            │      3       │      10      │       25        │      Unlimited      │
│ Log Retention        │   1 day      │   7 days     │    30 days      │      Unlimited      │
└──────────────────────┴──────────────┴──────────────┴─────────────────┴────────────────────┘
```

---

*Document Version: 2.1*  
*Last Updated: 2026-03-10*  
*For the latest feature list, visit: https://aegisgate.example.com/features*
