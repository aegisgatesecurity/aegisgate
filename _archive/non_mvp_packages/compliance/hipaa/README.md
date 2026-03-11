# HIPAA Compliance Framework

This directory contains the HIPAA (Health Insurance Portability and Accountability Act) compliance framework implementation for AegisGate.

## Overview

The HIPAA framework helps ensure that AI systems handling healthcare data maintain compliance with HIPAA regulations. It monitors both incoming requests and outgoing responses for potential PHI (Protected Health Information) indicators and violations.

## Features

- **PHI Detection**: Identifies potential PHI indicators in requests
- **Leakage Prevention**: Detects potential PHI disclosure in responses
- **Audit Trail**: Maintains timestamped violation records
- **Configurable Rules**: Enable/disable specific compliance rules

## Compliance Rules

| Rule ID | Name | Severity | Description |
|---------|------|----------|-------------|
| HIPAA_ACCESS_CONTROL | Access Control | CRITICAL | Implement access control measures to protect PHI |
| HIPAA_AUDIT_CONTROL | Audit Control | HIGH | Implement audit controls to record and examine activity |
| HIPAA_AUTHENTICATION | Authentication | CRITICAL | Implement appropriate authentication measures |
| HIPAA_INTEGRITY | Data Integrity | HIGH | Ensure PHI is not altered or destroyed improperly |
| HIPAA_TRANSPORT | Transport Security | CRITICAL | Implement security measures for data transport |
| HIPAA_ENCRYPTION | Encryption | CRITICAL | Implement encryption for PHI at rest and in transit |

## Usage

```go
import (
    "github.com/aegisgatesecurity/aegisgate/pkg/compliance/hipaa"
)

// Create HIPAA framework instance
hipaaFramework := &hipaa.HIPAA{}

// Load rules
err := hipaaFramework.LoadRules()
if err != nil {
    log.Fatal(err)
}

// Check request for HIPAA violations
violations, err := hipaaFramework.CheckRequest(request)
if err != nil {
    log.Fatal(err)
}

// Check response for HIPAA violations
violations, err = hipaaFramework.CheckResponse(response)
if err != nil {
    log.Fatal(err)
}
```

## PHI Indicators

The framework monitors for the following PHI indicators in requests:
- Medical records references
- Patient ID numbers
- Social Security numbers
- Dates of birth
- Diagnoses and treatments
- Prescription information
- Lab results and imaging data

## Integration

This framework integrates with AegisGate's compliance mapper to provide comprehensive healthcare data protection as part of your AI security gateway.
