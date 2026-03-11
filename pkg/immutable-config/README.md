# Immutable Configuration System

A comprehensive, production-ready immutable configuration management system for the AegisGate Chatbot Security Gateway.

## Overview

This package provides a complete immutable, read-only filesystem implementation with the following key features:

- **Versioned Configuration Storage**: All configurations are versioned and stored with full audit trails
- **Integrity Verification**: SHA-256 hash verification for all configuration data
- **Write-Ahead Logging (WAL)**: Atomic operations with crash recovery
- **Point-in-Time Snapshots**: Full state capture with integrity checksums
- **File Change Detection**: Real-time monitoring for unauthorized modifications
- **Seal/Unseal Mechanism**: Lock the filesystem against all modifications
- **Audit Logging**: Complete operation trail with timestamps and signatures

## Architecture

```
pkg/immutable-config/
├── config.go              # ConfigData and ConfigVersion types
├── manager.go             # ConfigManager implementation
├── api.go                 # Public API
├── options.go             # Configuration options
├── filesystem/
│   ├── filesystem.go      # ImmutableFilesystem integration layer
│   ├── provider.go        # FilesystemProvider for disk persistence
│   └── filesystem_test.go # Comprehensive tests
├── readonly/
│   └── readonly.go        # ReadOnlyProvider with seal/unseal
├── snapshot/
│   └── snapshot.go        # SnapshotManager for point-in-time capture
├── wal/
│   └── wal.go             # Write-ahead logging for atomic operations
├── watcher/
│   └── watcher.go         # File change detection
├── integrity/
│   └── integrity.go       # SHA-256 integrity verification
├── logging/
│   └── audit.go           # Audit logging
└── rollback/
    └── manager.go         # Rollback management
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/aegisgatesecurity/aegisgate/pkg/immutable-config/filesystem"
)

func main() {
    // Create immutable filesystem
    fs, err := filesystem.NewImmutableFilesystem(&filesystem.FilesystemConfig{
        BasePath:        "/var/lib/aegisgate/config",
        MaxVersions:     100,
        MaxAuditEntries: 10000,
        EnableWatch:     true,
    })
    if err != nil {
        panic(err)
    }
    defer fs.Close()

    // Initialize
    if err := fs.Initialize(); err != nil {
        panic(err)
    }

    // Save configuration
    config := immutableconfig.NewConfigData("v1.0.0", map[string]interface{}{
        "setting1": "value1",
        "setting2": 42,
    }, map[string]string{
        "author": "admin",
    })

    version, err := fs.Save(config)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Saved version: %s\n", version.Version)

    // Load configuration
    loaded, err := fs.Load("v1.0.0")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Loaded: %v\n", loaded.Data)
}
```

### Seal/Unseal Operations

```go
// Seal the filesystem - no modifications allowed
if err := fs.Seal(); err != nil {
    panic(err)
}

fmt.Printf("Filesystem sealed at: %s\n", fs.SealedAt())

// Any save attempt will fail
_, err = fs.Save(newConfig)
// Error: filesystem is sealed: modifications are not allowed

// Unseal (admin operation)
if err := fs.Unseal(); err != nil {
    panic(err)
}
```

### Snapshots

```go
// Create snapshot
snapshot, err := fs.CreateSnapshot("backup-2026-02-21", "End of month backup")
if err != nil {
    panic(err)
}
fmt.Printf("Snapshot ID: %s\n", snapshot.ID)

// List snapshots
snapshots, err := fs.ListSnapshots()
for _, s := range snapshots {
    fmt.Printf("  %s: %s (%s)\n", s.ID, s.Name, s.Created)
}

// Restore from snapshot
if err := fs.RestoreSnapshot(snapshot.ID); err != nil {
    panic(err)
}

// Verify snapshot integrity
verified, err := fs.snapshotMgr.Verify(snapshot.ID)
if verified {
    fmt.Println("Snapshot integrity verified")
}
```

### Integrity Verification

```go
// Verify all configurations
results, err := fs.VerifyIntegrity()
for version, valid := range results {
    fmt.Printf("Version %s: %v\n", version, valid)
}

// Get audit log
entries := fs.GetAuditLog()
for _, entry := range entries {
    fmt.Printf("[%s] %s: %s\n", entry.Timestamp, entry.EventType, entry.Description)
}
```

### Crash Recovery

```go
// Recover from WAL after crash
if err := fs.Recover(); err != nil {
    panic(err)
}

// Compact WAL to remove old entries
if err := fs.CompactWAL(); err != nil {
    panic(err)
}
```

## Components

### FilesystemProvider

Disk-based storage provider with atomic write operations:

- Atomic writes using temp file + rename
- Automatic loading of existing configurations
- Version management with SHA-256 hashes
- Read-only file permissions after save

### ReadOnlyProvider

Enforcement layer that wraps any Provider:

- Seal/Unseal mechanism
- Modification attempt tracking
- Blocked operation logging

### SnapshotManager

Point-in-time state capture:

- Full configuration snapshots
- Integrity checksums for verification
- Restore capability
- Snapshot listing and deletion

### WAL (Write-Ahead Log)

Atomic operations with recovery:

- Append-only log
- Commit/Rollback semantics
- Crash recovery replay
- Automatic compaction

### Watcher

Real-time file monitoring:

- SHA-256 checksum tracking
- Event-based notifications
- Unauthorized change detection
- Configurable scan intervals

## Configuration Options

```go
type FilesystemConfig struct {
    BasePath        string        // Directory for config storage
    MaxVersions     int           // Maximum versions to retain
    MaxAuditEntries int           // Maximum audit log entries
    WatchInterval   time.Duration // File scan interval
    EnableWatch     bool          // Enable file watching
    AutoSeal        bool          // Auto-seal after writes
}
```

## Security Features

1. **Immutability**: Files set to read-only (0444) after write
2. **Integrity**: SHA-256 hash verification on all reads
3. **Audit Trail**: Complete logging of all operations
4. **Seal Mechanism**: Hard lock against modifications
5. **Change Detection**: Real-time monitoring with Watcher
6. **Crash Recovery**: WAL replay for data consistency

## Performance

- **Read Latency**: ~1-5ms (memory cached)
- **Write Latency**: ~10-50ms (with WAL + fsync)
- **Snapshot Time**: ~100-500ms (depends on config size)
- **Verification Time**: ~1-2ms per config

## Compliance

This implementation supports the following compliance requirements:

- **NIST AI RMF**: Configuration integrity controls
- **ISO/IEC 42001**: Change management requirements
- **SOC 2**: Audit trail requirements
- **PCI-DSS**: File integrity monitoring

## Testing

Run the test suite:

```bash
go test ./pkg/immutable-config/... -v
```

## License

Part of the AegisGate Chatbot Security Gateway project.
