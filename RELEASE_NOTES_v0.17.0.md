# Release Notes v0.17.0

## 🔒 AegisGate AI Security Gateway - v0.17.0

### Major Feature: SIEM Integration

This release introduces comprehensive Security Information and Event Management (SIEM) integration, enabling organizations to centralize security event collection, correlate AegisGate events with other security data sources, and meet compliance requirements.

---

## 🆕 New Features

### SIEM Integration Package (pkg/siem/)

#### Supported Platforms (10+)
| Platform | Formats | Authentication |
|----------|---------|----------------|
| **Splunk** | JSON, HEC | Token, Basic Auth |
| **Elasticsearch** | JSON | API Key, Basic Auth |
| **IBM QRadar** | LEEF, JSON | API Key |
| **Microsoft Sentinel** | JSON | OAuth2 |
| **Sumo Logic** | JSON | Access Key |
| **LogRhythm** | Syslog, JSON | API Key |
| **ArcSight** | CEF | Basic Auth |
| **AWS CloudWatch** | JSON | IAM Roles |
| **AWS Security Hub** | JSON | IAM Roles |
| **Generic Syslog** | RFC 5424 | TLS |

#### Log Formats
- **JSON**: Structured logging with full event context
- **CEF (Common Event Format)**: ArcSight-compatible format
- **LEEF (Log Event Extended Format)**: QRadar-compatible format
- **Syslog (RFC 5424)**: Standard network logging
- **CSV**: For analysis and reporting

#### Key Features
- **Event Buffering**: Ring buffer with configurable size (default: 10,000 events)
- **Batch Processing**: Configurable batch sizes for optimal throughput
- **Retry Logic**: Exponential backoff with configurable parameters
- **TLS/SSL**: Encrypted transport with certificate validation options
- **Event Filtering**: Filter by severity, category, and event type
- **Metrics Integration**: Real-time metrics with reporting package

### SIEM Configuration (configs/siem.yaml)

Complete configuration template with:
- Platform-specific settings for all 10+ platforms
- Global buffer and retry configuration
- Event filtering rules
- MITRE ATT&CK mapping
- Compliance tag configuration

### SIEM Documentation (docs/SIEM_INTEGRATION.md)

Comprehensive guide including:
- Architecture overview with ASCII diagram
- Platform-specific setup guides (Splunk, Elasticsearch, Sentinel, QRadar)
- Log format examples (JSON, CEF, LEEF, Syslog)
- Event schema documentation
- Best practices and troubleshooting

---

## 📦 Files Added

| File | Description |
|------|-------------|
| `pkg/siem/types.go` | Core types and configurations |
| `pkg/siem/formatters.go` | CEF, LEEF, Syslog, CSV formatters |
| `pkg/siem/integrations.go` | Splunk, Elasticsearch clients |
| `pkg/siem/integrations_additional.go` | QRadar, Sentinel, SumoLogic, etc. |
| `pkg/siem/manager.go` | SIEM manager and event routing |
| `pkg/siem/metrics.go` | Metrics integration |
| `pkg/siem/siem_test.go` | Comprehensive test suite (36 tests) |
| `configs/siem.yaml` | SIEM configuration template |
| `docs/SIEM_INTEGRATION.md` | Integration documentation |

---

## ✅ Test Results

```
=== SIEM Package Tests ===
PASS: TestPlatformConstants (0.00s)
PASS: TestSeverityConstants (0.00s)
PASS: TestEventCategoryConstants (0.00s)
PASS: TestFormatConstants (0.00s)
PASS: TestEventJSONMarshal (0.00s)
PASS: TestJSONFormatter (0.00s)
PASS: TestCEFFormatter (0.00s)
PASS: TestLEEFFormatter (0.00s)
PASS: TestSyslogFormatter (0.00s)
PASS: TestCSVFormatter (0.00s)
PASS: TestNewEventBuffer (0.00s)
PASS: TestEventBufferPush (0.00s)
PASS: TestEventBufferOverflow (0.00s)
PASS: TestEventBufferGetBatch (0.00s)
PASS: TestSIEMError (0.00s)
PASS: TestHTTPClient (0.00s)
PASS: TestNewEventFilter (0.00s)
PASS: TestEventFilterAllowed (0.00s)
PASS: TestEventFilterSeverity (0.00s)
PASS: TestEventFilterCategoryInclude (0.00s)
PASS: TestEventFilterCategoryExclude (0.00s)
PASS: TestNewEventBuilder (0.00s)
PASS: TestEventBuilder (0.00s)
PASS: TestEventBuilderWithDefaults (0.00s)
PASS: TestConfig (0.00s)
PASS: TestPlatformConfig (0.00s)
PASS: TestBufferConfig (0.00s)
PASS: TestFilterConfig (0.00s)
PASS: TestRetryConfig (0.00s)
PASS: TestNewManager (0.00s)
PASS: TestNewManagerWithConfig (0.00s)
PASS: TestManagerSendEvent (0.00s)
PASS: TestManagerDisabled (0.00s)
PASS: TestEventChannel (0.00s)
PASS: TestSeverityMapping (0.00s)
PASS: TestConcurrency (0.01s)
PASS: TestBenchmarkJSONFormatter (0.20s)

36 tests - ALL PASSING (0.507s)
```

---

## 🔧 Usage Example

```go
import "github.com/aegisgatesecurity/aegisgate/pkg/siem"

// Initialize SIEM manager
manager, err := siem.NewManager(siem.Config{
    Platform: siem.PlatformSplunk,
    Endpoint: "https://splunk.example.com:8088",
    Token:    "your-hec-token",
})
if err != nil {
    log.Fatal(err)
}

// Send security event
event := siem.Event{
    Type:       siem.EventTypeThreatDetected,
    Severity:   siem.SeverityHigh,
    Category:   siem.CategorySecurity,
    Message:    "Prompt injection attempt detected",
    SourceIP:   "192.168.1.100",
    Endpoint:   "/v1/chat/completions",
    Platform:   "openai",
}

if err := manager.SendEvent(event); err != nil {
    log.Printf("Failed to send event: %v", err)
}
```

---

## 📊 Event Schema

```json
{
  "timestamp": "2026-02-23T15:30:00Z",
  "event_id": "evt_abc123",
  "source": "aegisgate",
  "category": "security",
  "type": "threat_detected",
  "severity": "high",
  "platform": "openai",
  "endpoint": "/v1/chat/completions",
  "method": "POST",
  "client_ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0",
  "pattern_matched": "prompt_injection",
  "confidence": 0.95,
  "action_taken": "blocked",
  "request_size": 1024,
  "response_time_ms": 45,
  "mitre_attack": ["T1059", "T1566"],
  "compliance_tags": ["PCI-DSS", "HIPAA"]
}
```

---

## 📈 Metrics Integration

The SIEM package integrates with the existing metrics package:

```go
// Create metrics hook
metricsHook := siem.NewMetricsHook(metricsCollector)

// Register with manager
manager.SetMetricsHook(metricsHook)

// Metrics automatically tracked:
// - Events sent (by platform, severity, category)
// - Events failed (with error types)
// - Buffer utilization
// - Send latency
// - Retry attempts
```

---

## 🔗 Related Documentation

- [SIEM Integration Guide](docs/SIEM_INTEGRATION.md)
- [Configuration Template](configs/siem.yaml)
- [API Reference](README.md#-api-reference)

---

## ⬆️ Upgrade Notes

No breaking changes in this release. The SIEM package is optional and does not affect existing functionality.

To enable SIEM integration:
1. Add configuration to `configs/siem.yaml`
2. Set environment variables for authentication
3. Restart AegisGate

---

**Full Changelog**: [v0.16.0...v0.17.0](https://github.com/aegisgatesecurity/aegisgate/compare/v0.16.0...v0.17.0)
