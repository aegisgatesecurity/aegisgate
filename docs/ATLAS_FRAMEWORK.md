# MITRE ATLAS Framework Implementation

## Overview

This document describes the MITRE ATLAS (Adversarial Threat Landscape for AI Systems) implementation in AegisGate, providing comprehensive coverage of adversarial threats against AI/ML systems.

## Supported Techniques

The implementation covers 18 MITRE ATLAS techniques across multiple categories:

### Input Manipulation (IM)
| Technique ID | Name | Description | Patterns |
|--------------|------|-------------|----------|
| ATLAS-T1535 | Direct Prompt Injection | Malicious instructions injected into LLM inputs | 5 |
| ATLAS-T1584 | Indirect Prompt Injection | Hidden instructions in external data | 5 |
| ATLAS-T1632 | System Prompt Extraction | Attempts to extract system prompts | 5 |

### Model Manipulation (MM)
| Technique ID | Name | Description | Patterns |
|--------------|------|-------------|----------|
| ATLAS-T1484 | LLM Jailbreak | Bypassing safety controls | 5 |
| ATLAS-T1658 | Adversarial Examples | Specially crafted inputs to cause misclassification | 5 |
| ATLAS-T1648 | Serverless Compute Exploitation | Exploiting serverless ML deployments | 4 |

### Resource Access (RA)
| Technique ID | Name | Description | Patterns |
|--------------|------|-------------|----------|
| ATLAS-T1589 | Training Data Exposure | Extracting sensitive training data | 5 |
| ATLAS-T1590 | Network Information Gathering | Reconnaissance of AI infrastructure | 4 |
| ATLAS-T1592 | Host Information Gathering | Gathering host configuration info | 4 |
| ATLAS-T1552 | Unsecured Credentials | Exposing ML credentials | 4 |
| ATLAS-T1041 | Exfiltration Over C2 | Data exfiltration via command channel | 5 |

### Authorization and Authentication (AA)
| Technique ID | Name | Description | Patterns |
|--------------|------|-------------|----------|
| ATLAS-T1556 | Modify Authentication | Altering ML authentication | 4 |
| ATLAS-T1611 | Escape to Host | Container escape to underlying host | 4 |
| ATLAS-T1621 | MFA Request Generation | Abusing MFA mechanisms | 4 |

### Impact (IC)
| Technique ID | Name | Description | Patterns |
|--------------|------|-------------|----------|
| ATLAS-T1486 | Data Encrypted for Impact | AI model ransomware | 5 |
| ATLAS-T1566 | Phishing (AI-Enhanced) | AI-augmented phishing attacks | 5 |
| ATLAS-T1599 | Network Boundary Bridging | Crossing network boundaries | 4 |
| ATLAS-T1110 | Brute Force | Credential brute-forcing for AI APIs | 4 |

## Detection Patterns

Total: 86 regex-based detection patterns

Each technique includes multiple patterns covering known attack signatures, anomaly indicators, behavioral patterns, and heuristic rules.

## Framework Mappings

### NIST AI RMF - MITRE ATLAS

The mapping framework provides bidirectional relationships between NIST AI RMF functions and MITRE ATLAS techniques:

| NIST AI RMF Function | ATLAS Techniques Covered | Relationship Types |
|---------------------|-------------------------|-------------------|
| GOVERN (GV1-GV4) | All 18 techniques | supports, mitigates, addresses |
| MAP (MP1-MP4) | All 18 techniques | detects, supports |
| MEASURE (ME1-ME4) | All 18 techniques | detects |
| MANAGE (RG1-RG4) | All 18 techniques | mitigates, supports |

Confidence scores: 0.75 - 0.95

### SOC 2 - MITRE ATLAS

| SOC 2 Control | ATLAS Techniques | Relationship |
|--------------|------------------|--------------|
| CC3.2 (AI Risk Assessment) | T1535, T1566 | addresses |
| CC5.4 (Model Change Control) | T1632 | addresses |
| CC6.2 (ML Environment Security) | T1584, T1556, T1611, T1621, T1110 | mitigates, detects |
| CC6.3 (Data Protection) | T1589, T1552, T1041 | mitigates |
| CC6.4 (Adversarial Defense) | T1535, T1484, T1632, T1648, T1486 | mitigates |
| CC6.5 (Vulnerability Management) | T1589, T1584, T1590, T1592, T1599 | detects, mitigates |
| CC6.6 (System Operations) | T1658 | detects |
| PI1.2 (Processing Integrity) | T1658, T1484 | supports, mitigates |

## Usage

### Running ATLAS Scans

// Create new ATLAS manager
manager := NewAtlasManager()

// Scan input for threats
findings := manager.ScanInput(userInput)

// Check for specific technique
hasInjection := manager.HasTechnique(userInput, "T1535")

### Generating Compliance Reports

// Generate unified compliance report
mapping := NewFrameworkMapping()
report := mapping.GenerateUnifiedReport()

// Get techniques for a specific control
techniques := mapping.GetTechniquesForControl("GV1")

## Coverage Matrix

| Category | Techniques | Detection Patterns |
|----------|------------|-------------------|
| Prompt Injection | 3 | 15 |
| Model Manipulation | 3 | 14 |
| Resource Access | 5 | 22 |
| Auth/Authorization | 3 | 12 |
| Impact | 4 | 18 |
| Total | 18 | 81 |

## Confidence Levels

- 0.95: Direct technical relationship with strong evidence
- 0.90: Strong correlation with specific controls
- 0.85: Moderate relationship with supporting evidence
- 0.80: General relationship with some overlap
- 0.75: Indirect or contextual relationship

## Integration Points

1. Input Validation: All user inputs scanned before processing
2. Model Inference: Runtime behavior monitoring
3. Compliance Reporting: Automated gap analysis
4. Risk Assessment: Framework-aligned risk scoring

## Future Enhancements

- Additional technique coverage (T1565 - Data Manipulation)
- Machine learning-based anomaly detection
- Real-time threat intelligence integration
- Extended pattern coverage for emerging threats
