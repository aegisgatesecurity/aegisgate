# Trust Domain Package

This package provides feed-specific trust domain management for the AegisGate AI Security Gateway.

## Overview

The trust domain package implements isolation, policy management, and validation services specifically designed for threat intelligence feeds. It prevents cascade failures by maintaining separate trust boundaries for each feed.

## Features

- Feed-Specific Trust Domains: Each feed gets its own dedicated trust domain
- Isolation: Complete isolation between feeds prevents cross-feed contamination
- Policy Engine: Flexible policy management for different feed types
- Validation Services: Certificate, signature, and hash chain validation
- Lifecycle Management: Complete management of trust domain lifecycle
- Audit Logging: Comprehensive audit logging for all operations

## Architecture

### Core Components

1. TrustDomain Interface: Defines the contract for all trust domains
2. ValidationEngine: Provides validation services for trust domains
3. FeedTrustPolicyEngine: Manages feed-specific trust policies
4. TrustDomainManager: Handles lifecycle management of trust domains
5. Integration: Provides integration methods with external systems

### Trust Domain Lifecycle

1. Creation: When a feed is registered, a dedicated trust domain is created
2. Configuration: Feed-specific trust anchors and policies are applied
3. Operation: All validations use the feed's dedicated trust domain
4. Rotation: Trust anchors can be rotated per feed without affecting others
5. Decommission: When feed is removed, trust domain is securely destroyed

## Usage

### Basic Usage

Create trust domain manager:

manager := trustdomain.NewTrustDomainManager(nil)

Create a new trust domain for a feed:

domain, err := manager.CreateDomain("feed1_domain", "feed1")
if err != nil {
    log.Fatal(err)
}

Add trust anchors:

anchor := &trustdomain.TrustAnchor{
    // Trust anchor configuration
}
if err := domain.AddTrustAnchor(anchor); err != nil {
    log.Fatal(err)
}

Validate certificates:

result, err := domain.ValidateCertificate(cert)
if err != nil {
    log.Printf("Validation failed: %v", err)
} else {
    log.Printf("Validation result: %v", result)
}

### Policy-Based Usage

Create policy engine:

policyEngine := trustdomain.DefaultPolicyEngine()

Set domain for a feed:

policyEngine.SetDomain("feed1", domain)

Create and set policy:

policy := &trustdomain.FeedTrustPolicy{
    FeedID:         "feed1",
    TrustDomainID:  "feed1_domain",
    ValidationMode: trustdomain.ValidationStrict,
    Parameters: map[string]interface{}{
        "timeout": 30 * time.Second,
    },
}

policyEngine.SetPolicy("feed1", policy)

Validate request based on policies:

valid, err := policyEngine.ValidateFeed("feed1", request)

## Integration

### With Threat Intel Package

integration := trustdomain.NewIntegration(nil)

// Setup for each feed
for _, feed := range feeds {
    domainID := trustdomain.TrustDomainID(feed.ID + "_domain")
    if err := integration.SetupForFeed(feed.ID, domainID); err != nil {
        log.Printf("Failed to setup trust domain for feed %s: %v", feed.ID, err)
        continue
    }
}

// Validate feed requests
for _, request := range requests {
    valid, err := integration.ValidateFeedRequest(request.FeedID, request)
    if !valid {
        log.Printf("Request validation failed: %v", err)
        continue
    }
    // Process valid request
}

### With PKI Attestation Package

pkiService := pkiattest.NewAttestationService(config)

// Create domain with PKI attestation
validationEngine := trustdomain.NewValidationEngine(domain, config, pkiService, hashStore)

// Validate certificate with PKI attestation
result, err := validationEngine.ValidateCertificate(cert)

## Security Considerations

1. Isolation: Trust domains must be completely isolated
2. Revocation: Each domain handles its own revocation
3. Audit: All trust domain operations should be logged
4. Rotation: Secure key rotation without disruption

## Configuration

See the architecture documentation in docs/trustdomain/architecture.md for detailed configuration options.

## Testing

Unit tests and integration tests are available in the test files:
- validation_test.go
- manager_test.go
- policy_engine_test.go
- integration_test.go

## License

Apache 2.0 - See LICENSE in the project root for details.

---

Version: v0.19.0  
Last Updated: 2026-02-24  
Author: AegisGate Team
