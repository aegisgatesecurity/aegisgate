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
├── enterprise/      # Enterprise tier
│   ├── nist/        # NIST AI RMF + SP 1500
│   └── iso42001/    # ISO/IEC 42001
└── premium/         # Premium tier
    ├── soc2/        # SOC 2 Type II
    ├── hipaa/       # HIPAA compliance
    └── pci/         # PCI DSS
```

## Framework Support

| Tier      | Frameworks                           |
|-----------|--------------------------------------|
| Community | ATLAS, OWASP AI, GDPR                |
| Enterprise| + NIST RMF, ISO 42001                |
| Premium   | + SOC 2, HIPAA, PCI DSS              |

*Enterprise and Premium tier frameworks require appropriate licensing.*

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
