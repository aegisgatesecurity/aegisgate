# API Wrapper Compatibility Layer

## Overview

This document describes the API wrapper compatibility layer created to support integration tests in the AegisGate WAF project. These wrappers provide a unified API surface that matches what the integration tests expect, bridging the gap between existing package implementations and test requirements.

## Background

The integration tests in `tests/integration/integration_test.go` expected certain types and functions that didn't exist in the individual packages. Rather than modifying the tests directly, we created these API wrappers to maintain backward compatibility and provide a consistent interface.

## Design Principles

### Type Aliases
We use type aliases (`type X = Y`) rather than new types to avoid conversion overhead and allow seamless interoperability between the wrapper types and underlying implementations.

### Wrapper Functions
When the test API requires different function signatures or additional functionality, we provide wrapper functions that translate between the test expectations and the actual package implementations.

### Non-Breaking Changes
All wrappers are designed to be non-breaking - existing code continues to work without modification while tests can use the new API signatures.

## Package-by-Package Reference

### pkg/dashboard/api.go

**Purpose:** Provides `Server` type alias and unified constructor.

| Wrapper Type | Underlying Type | Description |
|--------------|-----------------|-------------|
| `Server` | `Dashboard` | Type alias for the dashboard server |

**Functions:**
```go
func NewServer(conf *Config, metricsCollector *metrics.Collector, eventBus *metrics.EventBus) *Server
```
Creates a new dashboard server with the unified API expected by integration tests.

### pkg/metrics/api.go

**Purpose:** Provides Collector type alias and EventBus implementation.

| Wrapper Type | Underlying Type | Description |
|--------------|-----------------|-------------|
| `Collector` | `MetricsCollector` | Type alias for the metrics collector |
| `EventBus` | New implementation | Event subscription/publishing system |

**Structures:**
```go
type EventBus struct {
    subscribers map[string][]chan MetricEvent
    mu          sync.RWMutex
    closed      bool
}
```

**Functions:**
```go
func NewEventBus() *EventBus
func (eb *EventBus) Subscribe(eventType string) chan MetricEvent
func (eb *EventBus) Unsubscribe(eventType string, ch chan MetricEvent)
func (eb *EventBus) Publish(event MetricEvent)
func (eb *EventBus) Close()

func NewCollectorConfig(config *Config) *Collector
func (c *Collector) Start()
func (c *Collector) RecordEvent(event *Event)
```

### pkg/ml/api.go

**Purpose:** Provides AnomalyDetector wrapper with feature vector support.

| Wrapper Type | Underlying Type | Description |
|--------------|-----------------|-------------|
| `AnomalyDetector` | `Detector` | Wrapper for anomaly detection |
| `FeatureVector` | `[]float64` | Type alias for ML features |

**Structures:**
```go
type DetectorConfig struct {
    LearningRate    float64
    Threshold       float64
    WindowSize      int
    FeatureCount    int
    MinSamples      int
    MaxFeatures     int
    // additional config fields...
}
```

**Functions:**
```go
func NewAnomalyDetector(config *DetectorConfig) (*AnomalyDetector, error)
func (ad *AnomalyDetector) UpdateBaseline(vector FeatureVector)
func (ad *AnomalyDetector) IsAnomaly(vector FeatureVector) (bool, float64)
func (ad *AnomalyDetector) GetConfig() DetectorConfig
func (ad *AnomalyDetector) CalculateAnomalyScore(vector FeatureVector) float64
```

### pkg/compliance/api.go

**Purpose:** Provides Framework type aliases and unified constructor.

| Wrapper Type | Underlying Type | Description |
|--------------|-----------------|-------------|
| `FrameworkType` | `Framework` | Type alias for compliance frameworks |
| `Framework` | `PCI_DSS` | Standard compliance constant |

**Structures:**
```go
type APIConfig struct {
    EnabledFrameworks []FrameworkType
    AutoDetect        bool
    GenerateReports   bool
    ActiveFrameworks  []string
}
```

**Functions:**
```go
func NewFramework(config *APIConfig) (*ComplianceManager, error)
```

### pkg/scanner/api.go

**Purpose:** Provides Severity type aliases for test expectations.

| Wrapper Type | Underlying Type | Description |
|--------------|-----------------|-------------|
| `SeverityLevel` | `Severity` | Type alias for severity levels |
| `All` | `Severity(0)` | Match all severities |
| `None` | `Severity(1)` | Match no severities |

**Note:** The functions `AllPatterns()` and `ShouldBlock()` already exist in other scanner files (`patterns_extended.go` and `patterns.go`), so no additional wrappers were needed for these.

## Usage in Tests

The integration tests can now import and use these types:

```go
import (
    "github.com/aegisgatesecurity/aegisgate/pkg/dashboard"
    "github.com/aegisgatesecurity/aegisgate/pkg/metrics"
    "github.com/aegisgatesecurity/aegisgate/pkg/ml"
    "github.com/aegisgatesecurity/aegisgate/pkg/compliance"
)

// Example usage:
server := dashboard.NewServer(config, metricsCollector, eventBus)
anomalyDetector := ml.NewAnomalyDetector(mlConfig)
framework, err := compliance.NewFramework(apiConfig)
```

## Maintenance Notes

### When to Update

1. **New Integration Tests:** When adding new integration tests that require additional API capabilities
2. **Package Refactoring:** When refactoring packages changes function signatures or type names
3. **Test Expectation Changes:** When integration test requirements evolve

### Version Management

- Keep wrappers synchronized with the main package APIs
- Document any API version requirements
- Consider adding version tags for major API changes

### Testing

Run the following to verify wrappers compile correctly:
```bash
go build ./...
go test ./tests/integration/...
```

## Known Issues and Limitations

1. **EventBus buffering:** The current EventBus implementation uses buffered channels. Adjust `BufferSize` in `Config` if you encounter dropped events in high-throughput scenarios.

2. **AnomalyDetector feature vectors:** The `FeatureVector` type is a simple `[]float64` alias. More sophisticated ML feature types may require additional wrapping.

3. **Severity constants:** `All` and `None` are special sentinel values. Ensure they don't conflict with any future severity levels added to the scanner package.

## Future Enhancements

- Add comprehensive unit tests for wrapper functions
- Consider auto-generation tools for maintaining synchronization with package APIs
- Add configuration validation to wrapper constructors
- Document performance characteristics of wrapped functions

## Contributing

When modifying wrappers:
1. Ensure all integration tests continue to pass
2. Update this documentation with new types/functions
3. Consider backward compatibility for existing test code
4. Follow the established naming conventions (e.g., `NewX()` for constructors)

---

**Generated:** February 2026  
**Last Commit:** b732a1c - fix: add API compatibility wrappers for integration tests  
**Maintainer:** Project Contributors
