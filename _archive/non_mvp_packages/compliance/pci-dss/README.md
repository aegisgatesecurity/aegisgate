# PCI-DSS Compliance Framework

This directory contains the PCI-DSS (Payment Card Industry Data Security Standard) compliance framework implementation for AegisGate.

## Overview

The PCI-DSS framework helps ensure that AI systems handling payment card data maintain compliance with PCI-DSS regulations. It monitors both incoming requests and outgoing responses for potential cardholder data indicators and violations.

## Features

- **Cardholder Data Detection**: Identifies potential cardholder data in requests
- **Leakage Prevention**: Detects potential cardholder data disclosure in responses
- **Audit Trail**: Maintains timestamped violation records
- **Configurable Rules**: Enable/disable specific compliance rules

## Compliance Rules

| Rule ID | Name | Severity | Description |
|---------|------|----------|-------------|
| PCI_BUILD_SECURE | Build Secure Systems | CRITICAL | Build and maintain secure systems and software |
| PCI_PROTECT_CARDHOLDER | Protect Cardholder Data | CRITICAL | Protect stored cardholder data |
| PCI_MAINTAIN_POLICY | Maintain Policy | HIGH | Maintain an information security policy |
| PCI_IDENTITY_MANAGEMENT | Identity Management | HIGH | Implement strong identity management |
| PCI_ACCES_CONTROL | Access Control | CRITICAL | Restrict access to cardholder data |
| PCI_MONITORING | Monitoring & Testing | HIGH | Monitor and test networks |

## Usage

```go
import (
    "github.com/aegisgatesecurity/aegisgate/pkg/compliance/pcidss"
)

// Create PCI-DSS framework instance
pciFramework := &pcidss.PCI_DSS{}

// Load rules
err := pciFramework.LoadRules()
if err != nil {
    log.Fatal(err)
}

// Check request for PCI-DSS violations
violations, err := pciFramework.CheckRequest(request)
if err != nil {
    log.Fatal(err)
}

// Check response for PCI-DSS violations
violations, err = pciFramework.CheckResponse(response)
if err != nil {
    log.Fatal(err)
}
```

## Cardholder Data Indicators

The framework monitors for the following indicators in requests:
- Credit card number references
- Cardholder name references
- CVV or security code references
- Expiration date references
- Billing address references

## Integration

This framework integrates with AegisGate's compliance mapper to provide comprehensive payment card data protection as part of your AI security gateway.
