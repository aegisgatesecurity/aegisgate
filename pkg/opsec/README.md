# OPSEC Module

Operational Security Module for AegisGate LLM Gateway

## Overview

The OPSEC (Operational Security) module provides comprehensive operational security features for the AegisGate security gateway, implementing defense in depth strategies for handling sensitive data and maintaining security postures.

## Features

### 🔒 Audit Logging
- Thread-safe, tamper-evident audit trail with SHA-256 hash chains
- Configurable severity levels (debug, info, warning, error, fatal, critical)
- JSON export/import for compliance reporting
- Entry pruning to prevent memory exhaustion
- Chain integrity verification

### 🔑 Secret Rotation
- Auto/manual secret rotation with configurable periods
- Cryptographic randomness via `crypto/rand`
- Memory-safe secret storage with secure wiping
- Rotation analytics and metrics

### 🧹 Memory Scrubbing
- Secure memory wiping using `crypto/subtle`
- Multi-pass secure deletion (DoD 5220.22-M inspired)
- Thread-safe operations
- Prevents compiler optimization from scrubbing

### 🛡️ Threat Modeling
- 7 pre-defined LLM/AI threat vectors
- OWASP AI Top 10 mapping
- Real-time pattern analysis
- Custom threat registration

### 🔧 Runtime Hardening
- ASLR detection
- Linux capability management (stubs)
- Resource limit configuration
- Seccomp profile support (stubs)

## Quick Start

```go
import "github.com/aegisgatesecurity/aegisgate/pkg/opsec"

// Create with default configuration
manager := opsec.New()

// Initialize all components
err := manager.Initialize()
if err != nil {
    log.Fatal(err)
}

// Start background processes
defer manager.Start()
defer manager.Stop()

// Log audit events
manager.LogAudit("request_processed", map[string]string{
    "user_id": "12345",
    "action":  "prompt_validation",
})

// Analyze for threats
threats := manager.AnalyzeThreats(input, output)
if len(threats) > 0 {
    // Handle threats
}
```

## Configuration

```go
config := opsec.OPSECConfig{
    AuditEnabled:     true,
    AuditMaxEntries:  10000,
    RotationEnabled:  true,
    RotationPeriod:   24 * time.Hour,
    SecretLength:     32,
    MemoryScrubbing:  true,
    ThreatModeling:   true,
    RuntimeHardening: true,
}

manager := opsec.NewWithConfig(config)
```

## Security Considerations

- **Thread Safety**: All operations are thread-safe using `sync.RWMutex`
- **Memory Safety**: Secrets stored as `[]byte` for scrubbing capability
- **Zero Dependencies**: No external dependencies beyond stdlib
- **Audit Integrity**: SHA-256 hash chains for tamper detection

## Compliance

- NIST SP 800-53 Rev. 5 (Audit Controls)
- OWASP ASVS V7 (Logging and Monitoring)
- MITRE ATLAS (AI Threat Modeling)

## Testing

```bash
cd pkg/opsec
go test -v
go test -bench=.
```

## License

MIT License - See LICENSE file in repository root
