Feed-specific Trust Domains - Architecture Design Document

Version Information
- Version: 1.0
- Date: 2026-02-25
- Author: AegisGate Team
- Status: Draft for Phase 2 Implementation

1. Executive Summary

This document outlines the architecture and design principles for feed-specific trust domains in AegisGate v0.19.0.

1.1 Objectives
- Feed-level isolation to prevent cascade failures
- Feed-specific trust policies and validation rules
- Centralized trust domain lifecycle management
- Audit logging and compliance tracking
- Integration with existing threat intel and PKI attestation systems

1.2 Scope
This design covers:
- Trust domain definition and boundaries
- Policy engine implementation
- Validation service architecture
- Integration patterns with existing packages
- Configuration management

2. Architecture Overview

2.1 High-Level Architecture
The architecture consists of:
- Trust Domain Manager: Central component managing multiple trust domains
- Policy Engine: Feed-specific policy enforcement
- Validation Engine: Core validation logic
- Audit Service: Compliance logging

2.2 Key Components

2.2.1 Trust Domain Core
- Interface: TrustDomain - Defines essential operations
- Implementation: basicTrustDomain - Thread-safe implementation
- Features:
  - Lifecycle management (enable/disable/destroy)
  - Validation execution
  - Statistics tracking
  - Audit logging

2.2.2 Trust Domain Manager
- Interface: Manager
- Implementation: TrustDomainManager
- Features:
  - Domain registration and deregistration
  - Lifecycle orchestration
  - Audit log distribution

2.2.3 Policy Engine
- Interface: PolicyEngine
- Implementation: defaultPolicyEngine
- Features:
  - Domain assignment per feed
  - Policy configuration
  - Feed validation based on policies

2.2.4 Validation Engine
- Interface: ValidationEngine
- Implementation: basicValidationEngine
- Features:
  - Certificate validation
  - Signature validation
  - Hash chain validation
  - Policy-based validation

3. Trust Domain Definitions

3.1 Trust Domain ID
type TrustDomainID string
const ( TrustDomainIDPrefix = "td_" )
func GenerateTrustDomainID(feedID string) TrustDomainID

3.2 Trust Domain Structure
type TrustDomain struct {
    ID: TrustDomainID
    FeedID: string
    CreatedAt: time.Time
    Enabled: bool
    Statistics: *DomainStatistics
    Policies: []*FeedTrustPolicy
    AuditLog: chan AuditLogEntry
}

3.3 Trust Domain Interface
type TrustDomain interface {
    Enable() error
    Disable() error
    Destroy() error
    IsEnabled() bool
    ValidateCertificate(cert *x509.Certificate) (*ValidationResult, error)
    ValidateSignature(data, signature []byte) (*ValidationResult, error)
    ValidateHashChain(hash, previousHash string) (*ValidationResult, error)
    SetPolicies(policies []*FeedTrustPolicy) error
    AddPolicy(policy *FeedTrustPolicy) error
    RemovePolicy(policyID string) error
    GetStatistics() *DomainStatistics
    ResetStatistics() error
}

4. Policy Engine Architecture

4.1 Policy Structure
type FeedTrustPolicy struct {
    ID: string
    FeedID: string
    TrustDomainID: TrustDomainID
    ValidationMode: ValidationMode
    Parameters: map[string]interface{}
    CreatedAt: time.Time
    UpdatedAt: time.Time
}
type ValidationMode string
const ( ValidationSoft = "soft" ValidationStrict = "strict" ValidationAudit = "audit" )

4.2 Policy Engine Interface
type PolicyEngine interface {
    SetDomain(feedID string, domain TrustDomain) error
    GetDomain(feedID string) (TrustDomain, error)
    RemoveDomain(feedID string) error
    SetPolicy(feedID string, policy *FeedTrustPolicy) error
    GetPolicy(feedID string) (*FeedTrustPolicy, error)
    RemovePolicy(feedID string) error
    ValidateFeed(feedID string, data []byte) (bool, error)
}

5. Validation Engine Details

5.1 Validation Result Structure
type ValidationResult struct {
    Success: bool
    Timestamp: time.Time
    Error: error
    Details: map[string]interface{}
    PolicyID: string
    RuleID: string
}

6. Integration Points

6.1 With Threat Intel Package
- Digital signatures
- Certificate chains
- Hash chain integrity

6.2 With PKI Attestation Package
- Certificate validation
- Identity verification
- Chain of trust verification

7. Configuration Schema

7.1 Trust Domain Configuration
trust_domains:
  feed1:
    domain_id: td_feed1
    enabled: true
    validation_mode: strict
    policies:
      - id: policy1
        parameters:
          timeout: 30s
          max_depth: 5

8. Audit Logging

8.1 Audit Log Entry Structure
type AuditLogEntry struct {
    Timestamp: time.Time
    Event: string
    FeedID: string
    DomainID: string
    Details: map[string]interface{}
    Severity: AuditSeverity
}

9. Implementation Timeline

9.1 Week 1: Core Implementation
- Refactor existing trust domain interface
- Implement validation engine
- Create test infrastructure

9.2 Week 2: Manager & Policy Engine
- Implement trust domain manager
- Implement policy engine
- Create configuration loader

9.3 Week 3: Integration & Testing
- Integrate with threat intel package
- Integrate with PKI attestation package
- Write comprehensive tests

10. Security Considerations

10.1 Isolation Requirements
- Trust domains must be completely isolated
- No shared state between domains
- Secure error handling

10.2 Access Control
- Domain access restricted
- Audit all domain operations
- Secure delegation

11. Testing Strategy

11.1 Unit Tests
- Trust domain lifecycle operations
- Validation algorithms
- Policy evaluation logic

11.2 Integration Tests
- Cross-package integration
- Feed-specific trust domains

12. References

- AegisGate Phase 2 Implementation Plan
- Threat Intel Package Design
- PKI Attestation Package
- Security Audit Guidelines

Document Status: Draft
Next Review: 2026-03-03
Version Control: git commit history
