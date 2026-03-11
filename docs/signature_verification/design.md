Digital Signature Verification - Design Document

Version Information
- Version: 1.0
- Date: 2026-02-25
- Author: AegisGate Team
- Status: Draft for Phase 2 Implementation

1. Executive Summary

This document outlines the architecture and design principles for digital signature verification in AegisGate v0.19.0.
The signature verification system provides cryptographic validation of threat intelligence feed signatures.

1.1 Objectives
- Cryptographic signature validation for feed integrity
- Key management integration
- Support for multiple signature algorithms
- Fast signature verification with caching
- Audit logging for compliance

1.2 Scope
This design covers:
- Signature verification service
- Key management integration
- Multiple algorithm support (RSA, ECDSA, Ed25519)
- Signature caching and performance optimization
- Integration with threat intel and PKI attestation

2. Architecture Overview

2.1 High-Level Architecture
The architecture consists of:
- Verification Service: Core signature verification logic
- Key Manager: Key storage and management
- Algorithm Selector: Algorithm-specific verification logic
- Cache Manager: Signature verification caching

2.2 Key Components

2.2.1 Verification Service
- Interface: VerificationService
- Implementation: signatureVerificationService
- Features:
  - RSA signature verification
  - ECDSA signature verification
  - Ed25519 signature verification
  - Verification result caching

2.2.2 Key Manager
- Interface: KeyManager
- Implementation: keyManagerImpl
- Features:
  - Key storage and retrieval
  - Key rotation support
  - Key validation
  - Multiple key formats (PEM, JWK)

2.2.3 Algorithm Selector
- Interface: AlgorithmSelector
- Implementation: algorithmSelectorImpl
- Features:
  - Algorithm detection from signature
  - Algorithm-specific verification routing
  - Supported algorithms registry

2.2.4 Cache Manager
- Interface: CacheManager
- Implementation: cacheManagerImpl
- Features:
  - Verification result caching
  - Cache expiration policies
  - Cache statistics

3. Signature Verification Definitions

3.1 Signature Verification Result

type VerificationResult struct {
    Valid: bool
    Algorithm: string
    KeyID: string
    Timestamp: time.Time
    Error: error
}

3.2 Verification Service Interface

type VerificationService interface {
    VerifySignature(algorithm string, data []byte, signature []byte, publicKey []byte) (*VerificationResult, error)
    VerifySignatureWithKeyID(data []byte, signature []byte, keyID string) (*VerificationResult, error)
    CacheResult(key string, result *VerificationResult) error
    GetCachedResult(key string) (*VerificationResult, bool)
}

3.3 Key Manager Interface

type KeyManager interface {
    LoadPublicKey(keyID string) ([]byte, error)
    LoadPrivateKey(keyID string) ([]byte, error)
    StorePublicKey(keyID string, publicKey []byte) error
    StorePrivateKey(keyID string, privateKey []byte) error
    RotateKeys(keyID string) error
    ValidateKey(keyID string) error
}

4. Supported Algorithms

4.1 RSA Algorithms
- RSA with SHA-256 (RSASSA-PSS with SHA-256)
- RSA with SHA-384 (RSASSA-PSS with SHA-384)
- RSA with SHA-512 (RSASSA-PSS with SHA-512)

4.2 ECDSA Algorithms
- ECDSA with P-256 and SHA-256
- ECDSA with P-384 and SHA-384
- ECDSA with P-521 and SHA-512

4.3 Ed25519
- Ed25519 signature algorithm

5. Integration Points

5.1 With Threat Intel Package
- Feed signature verification
- Integration with FeedProcessor
- Error handling and logging

5.2 With PKI Attestation Package
- Certificate-based key retrieval
- Chain of trust verification
- Combined validation

5.3 With Trust Domain Package
- Feed-specific verification policies
- Shared audit logging
- Policy enforcement

6. Configuration Schema

6.1 Signature Verification Configuration

YAML Configuration:
signature_verification:
  enabled: true
  algorithms:
    - RSA-SHA256
    - ECDSA-P256-SHA256
    - Ed25519
  key_manager:
    type: filesystem
    path: configs/keys
  cache:
    enabled: true
    ttl: 5m
    max_size: 10000

7. Audit Logging

7.1 Audit Log Entry Structure

type AuditLogEntry struct {
    Timestamp: time.Time
    Event: string
    FeedID: string
    KeyID: string
    Result: VerificationResult
    Details: map[string]interface{}
    Severity: AuditSeverity
}

type AuditSeverity string
const ( AuditSeverityInfo = "info" AuditSeverityWarning = "warning" AuditSeverityCritical = "critical" )

8. Implementation Timeline

8.1 Week 1: Core Implementation
- Implement verification service
- Implement algorithm selector
- Create test infrastructure

8.2 Week 2: Key Manager & Integration
- Implement key manager
- Integrate with trust domain package
- Create configuration loader

8.3 Week 3: Testing & Documentation
- Write comprehensive tests
- Integration testing
- Documentation

9. Security Considerations

9.1 Key Storage
- Keys must be stored securely
- Private keys encrypted at rest
- Access controls for key retrieval

9.2 Algorithm Security
- Only approved algorithms supported
- Algorithm downgrade prevention
- Strong key requirements

9.3 Timing Attacks
- Constant-time operations where possible
- Prevent timing-based side-channel attacks

10. Testing Strategy

10.1 Unit Tests
- Verification with known test vectors
- Key management operations
- Algorithm-specific tests

10.2 Integration Tests
- End-to-end signature verification
- Integration with threat feeds
- Performance and load testing

11. References

- AegisGate Phase 2 Implementation Plan
- Trust Domain Package Design
- PKI Attestation Package Design
- Security Audit Guidelines

Document Status: Draft
Next Review: 2026-03-03
Version Control: git commit history

---

## Next Steps

1. Implement signature verification service in pkg/signature/
2. Create key management integration
3. Integrate with pkg/threatintel/
4. Write unit and integration tests
5. Document key management procedures
