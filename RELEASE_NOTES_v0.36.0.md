# AegisGate v0.36.0 Release Notes

<div align="center">

**Release Date**: March 6, 2026  
**Version**: 0.36.0  
**Go Version**: 1.24.0+

</div>

---

## Overview

This release introduces the **Plugin Architecture** - a comprehensive extensibility system for AegisGate that allows developers to extend proxy functionality through a well-defined plugin interface. This marks a significant milestone in AegisGate's evolution as an enterprise-grade security gateway.

---

## Table of Contents

1. [New Features](#new-features)
2. [Breaking Changes](#breaking-changes)
3. [Improvements](#improvements)
4. [Bug Fixes](#bug-fixes)
5. [Documentation](#documentation)
6. [Migration Guide](#migration-guide)
7. [Known Issues](#known-issues)
8. [Contributors](#contributors)
9. [Acknowledgments](#acknowledgments)
10. [Links](#links)

---

## New Features

### 🔌 Plugin System (Major Feature)

The plugin architecture provides a flexible way to extend AegisGate's functionality without modifying the core codebase.

#### Plugin Interface

```go
type Plugin interface {
    Metadata() PluginMetadata
    Init(ctx context.Context, config PluginConfig) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Hooks() []HookType
}
```

#### Hook Points (8 Total)

| Hook | Description | Use Case |
|------|-------------|----------|
| `request_received` | When request is first received | Logging, initial validation |
| `before_forward` | Before forwarding to upstream | Request modification |
| `after_response` | After response from upstream | Response capture |
| `response_sent` | Before response sent to client | Response modification |
| `connection_open` | When connection opens | Connection tracking |
| `connection_close` | When connection closes | Cleanup, stats |
| `error` | When error occurs | Error handling |
| `periodic` | Periodic task execution | Maintenance, cleanup |

#### Plugin Types

| Type | Description |
|------|-------------|
| **Filter** | Request/response processing and modification |
| **Auth** | Custom authentication providers |
| **Analytics** | Metrics collection and reporting |
| **Processor** | Data pipeline processing |
| **Exporter** | Export to external systems (SIEM, databases) |
| **Validator** | Policy and compliance validation |

#### Plugin Manager Features

- **Registration**: Register plugins at runtime
- **Dependency Management**: Declare and resolve plugin dependencies
- **Priority Execution**: Execute plugins in priority order (lower = earlier)
- **Runtime Config Updates**: Update plugin configuration without restart
- **Capability Discovery**: Query plugins by capability
- **Periodic Tasks**: Support for background periodic tasks

#### Configuration

```yaml
plugins:
  enabled: true
  directories:
    - "./plugins"
    - "/etc/aegisgate/plugins"
  timeout: 30s
  enable_periodic: true
  plugin_settings:
    example-plugin:
      enabled: true
      priority: 10
```

### Example Plugins

Three reference implementations are included:

#### 1. ExampleFilterPlugin
- Demonstrates request/response filtering
- Adds custom headers to requests/responses
- Logs request/response details
- Tracks statistics (requests processed, errors)

#### 2. ExampleAnalyticsPlugin
- Collects request/response metrics
- Tracks top paths, status codes
- Monitors bandwidth (bytes in/out)
- Calculates average latency

#### 3. ExamplePeriodicPlugin
- Demonstrates periodic task execution
- Runs maintenance tasks at intervals
- Logs execution timestamps
- Tracks task execution count

---

## Breaking Changes

**None** - This release is fully backward compatible.

---

## Improvements

### Core Improvements

1. **Plugin Configuration Integration**
   - Added `PluginConfig` to main `Config` struct
   - Plugin settings now loadable from `aegisgate.yml`
   - Environment variable support for plugin settings

2. **Code Quality**
   - Fixed deadlock in plugin manager startup
   - Removed unused functions (golangci-lint compliance)
   - Improved test coverage

### Configuration Examples

- Updated `aegisgate.yml.example` with plugin configuration section
- Added example plugin configuration documentation

---

## Bug Fixes

### Plugin System

| Issue | Fix |
|-------|-----|
| Deadlock in Start() | Refactored to collect periodic plugins before releasing lock |
| Unused function warnings | Removed unused `startPeriodicTasks()` function |
| Race conditions | Added proper mutex handling for concurrent access |

### General

- Various minor fixes from previous releases (see git history)

---

## Documentation

### New Documentation

- **README.md** - Comprehensive project documentation
  - Architecture diagrams
  - Getting started guide
  - Configuration reference
  - Plugin system documentation
  - API reference
  - Development guide

- **Plugin System Documentation**
  - Plugin interface specification
  - Hook point descriptions
  - Example plugin implementations
  - Configuration guide

---

## Migration Guide

### Upgrading from v0.35.x

1. **No Breaking Changes**: The upgrade is seamless
2. **Optional Plugin Configuration**: Add plugin settings to `aegisgate.yml` if desired:

```yaml
plugins:
  enabled: true
  directories:
    - "./plugins"
```

### New Plugin Development

To create a custom plugin:

1. **Implement the Plugin interface**:
```go
type MyPlugin struct{}

func (p *MyPlugin) Metadata() plugin.PluginMetadata {
    return plugin.PluginMetadata{
        ID:   "my-custom-plugin",
        Name: "My Custom Plugin",
        Version: "1.0.0",
    }
}

func (p *MyPlugin) Init(ctx context.Context, config plugin.PluginConfig) error {
    return nil
}

func (p *MyPlugin) Start(ctx context.Context) error {
    return nil
}

func (p *MyPlugin) Stop(ctx context.Context) error {
    return nil
}

func (p *MyPlugin) Hooks() []plugin.HookType {
    return []plugin.HookType{plugin.HookRequestReceived}
}
```

2. **Implement optional interfaces**:
   - `RequestProcessor` for request hooks
   - `ResponseProcessor` for response hooks
   - `PeriodicTask` for periodic tasks

3. **Register the plugin**:
```go
mgr := plugin.NewManager(config)
mgr.Register(&MyPlugin{})
```

---

## Known Issues

| Issue | Description | Workaround |
|-------|-------------|------------|
| None known | - | - |

---

## Contributors

This release was developed by the AegisGate team.

---

## Acknowledgments

- Go community for excellent tooling
- Contributors and users of AegisGate

---

## Links

| Resource | URL |
|----------|-----|
| Repository | https://github.com/aegisgatesecurity/aegisgate |
| Issues | https://github.com/aegisgatesecurity/aegisgate/issues |
| Releases | https://github.com/aegisgatesecurity/aegisgate/releases |
| Discussions | https://github.com/aegisgatesecurity/aegisgate/discussions |

---

## Installation

### Binary

Download from [Releases](https://github.com/aegisgatesecurity/aegisgate/releases)

### From Source

```bash
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
go build -o aegisgate ./cmd/aegisgate
```

### Docker

```bash
docker pull aegisgatesecurity/aegisgate:latest
```

---

<div align="center">

**Next Release**: v0.37.0 (Planned: GraphQL API, Penetration Testing)

</div>
