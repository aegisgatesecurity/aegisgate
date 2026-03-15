# AegisGate Compliance Framework

Tiered compliance checking for Chatbot Security Gateway.

## Structure

```
pkg/compliance/
├── common/          # Shared interfaces (all tiers)
├── community/       # Free tier frameworks
│   ├── atlas/       # MITRE ATLAS
│   ├── owasp/       # OWASP AI Top 10
│   └── gdpr/        # GDPR compliance
├── enterprise/      # Enterprise tier ($10-15K/mo)
│   ├── nist/        # NIST AI RMF + SP 1500
│   └── iso42001/    # ISO/IEC 42001
└── premium/         # Premium tier ($15-25K/mo)
    ├── soc2/        # SOC 2 Type II
    ├── hipaa/       # HIPAA compliance
    └── pci/         # PCI DSS
```

## Pricing

| Tier      | Monthly Cost | Frameworks                          |
|-----------|--------------|-------------------------------------|
| Community | Free         | ATLAS, OWASP AI, GDPR               |
| Enterprise| $10-15K      | + NIST RMF, ISO 42001               |
| Premium   | $15-25K      | + SOC 2, HIPAA, PCI DSS             |

## Usage

```go
import (
    "github.com/aegisgate/compliance"
    _ "github.com/aegisgate/compliance/community/atlas"
    _ "github.com/aegisgate/compliance/enterprise/nist"
)

func main() {
    registry := compliance.GetRegistry()
    findings, _ := registry.CheckAll(ctx, target)
}
```
