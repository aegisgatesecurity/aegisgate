# SIEM Integration Guide

This document provides comprehensive guidance for integrating AegisGate AI Security Gateway with Security Information and Event Management (SIEM) platforms.

## Table of Contents

1. [Overview](#overview)
2. [Supported Platforms](#supported-platforms)
3. [Configuration](#configuration)
4. [Log Formats](#log-formats)
5. [Event Schema](#event-schema)
6. [Platform-Specific Guides](#platform-specific-guides)
7. [Best Practices](#best-practices)
8. [Troubleshooting](#troubleshooting)

---

## Overview

The AegisGate SIEM integration module enables real-time forwarding of security events to external SIEM platforms. This allows organizations to:

- Centralize security event collection
- Correlate AegisGate events with other security data sources
- Enable advanced threat detection and response
- Meet compliance and audit requirements
- Implement automated incident response workflows

### Architecture

```
┌─────────────────────┐      ┌─────────────────────┐
│   AegisGate Gateway    │      │   SIEM Platform     │
│                     │      │                     │
│  ┌───────────────┐  │      │  ┌───────────────┐  │
│  │ Security      │  │      │  │ Event         │  │
│  │ Events        │──┼──────┼─▶│ Ingestion     │  │
│  └───────────────┘  │      │  └───────────────┘  │
│         │          │      │         │          │
│         ▼          │      │         ▼          │
│  ┌───────────────┐  │      │  ┌───────────────┐  │
│  │ SIEM Manager  │  │      │  │ Dashboards &  │  │
│  │ (Formatter)   │  │      │  │ Alerting      │  │
│  └───────────────┘  │      │  └───────────────┘  │
│         │          │      │                     │
│         ▼          │      │                     │
│  ┌───────────────┐  │      │                     │
│  │ Event Buffer  │  │      │                     │
│  │ (Retry/Queue) │  │      │                     │
│  └───────────────┘  │      │                     │
└─────────────────────┘      └─────────────────────┘
```

---

## Supported Platforms

| Platform | Format | Authentication | Status |
|----------|--------|----------------|--------|
| Splunk | JSON (HEC) | HEC Token | Production Ready |
| Elasticsearch | JSON | Basic, API Key | Production Ready |
| IBM QRadar | LEEF, JSON | API Token | Production Ready |
| Microsoft Sentinel | JSON | Shared Key | Production Ready |
| Sumo Logic | JSON | HTTP Source | Production Ready |
| LogRhythm | Syslog, JSON | API Key | Production Ready |
| ArcSight | CEF | Basic Auth | Production Ready |
| AWS CloudWatch | JSON | IAM | Production Ready |
| AWS Security Hub | JSON | IAM | Production Ready |
| Generic Syslog | RFC 5424 | - | Production Ready |

---

## Configuration

### Quick Start

1. Copy the example configuration:
   ```bash
   cp configs/siem.yaml.example configs/siem.yaml
   ```

2. Edit the configuration file and enable your SIEM platform:
   ```yaml
   platforms:
     - platform: splunk
       enabled: true
       endpoint: "https://your-splunk:8088"
       auth:
         type: "api_key"
         api_key: "your-hec-token"
   ```

3. Restart AegisGate to apply changes:
   ```bash
   systemctl restart aegisgate
   ```

### Configuration Reference

#### Global Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| app_name | string | "aegisgate" | Application identifier |
| environment | string | "production" | Environment tag |
| default_severity | string | "info" | Default event severity |
| include_raw | bool | true | Include raw event data |
| add_hostname | bool | true | Add hostname to events |

#### Filter Settings

| Setting | Type | Description |
|---------|------|-------------|
| min_severity | string | Minimum severity to forward |
| include_categories | []string | Categories to include |
| exclude_categories | []string | Categories to exclude |
| include_types | []string | Event types to include |
| exclude_types | []string | Event types to exclude |

#### Buffer Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| enabled | bool | true | Enable event buffering |
| max_size | int | 10000 | Maximum buffer size |
| flush_interval | duration | 5s | Buffer flush interval |
| persist | bool | false | Persist to disk |

#### Retry Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| enabled | bool | true | Enable retry logic |
| max_attempts | int | 3 | Maximum retry attempts |
| initial_backoff | duration | 1s | Initial backoff duration |
| max_backoff | duration | 30s | Maximum backoff duration |
| backoff_multiplier | float | 2.0 | Backoff multiplier |

---

## Log Formats

### JSON Format

Default format for Splunk, Elasticsearch, Sentinel, Sumo Logic:

```json
{
  "id": "01HXYZ123456",
  "timestamp": "2024-01-15T10:30:00Z",
  "source": "aegisgate",
  "category": "threat",
  "type": "blocked_request",
  "severity": "high",
  "message": "SQL injection attempt blocked",
  "attributes": {
    "source_ip": "192.168.1.100",
    "request_path": "/api/users?id=1' OR '1'='1",
    "http_method": "GET",
    "user_agent": "Mozilla/5.0..."
  },
  "entities": [
    {
      "type": "ip",
      "id": "ip-001",
      "name": "Source IP",
      "value": "192.168.1.100"
    }
  ],
  "mitre": {
    "tactic": "Initial Access",
    "technique": "T1190"
  },
  "compliance": [
    {
      "framework": "SOC2",
      "control": "CC6.1"
    }
  ]
}
```

### CEF Format (ArcSight)

Common Event Format for ArcSight:

```
CEF:0|Block|AegisGate|1.0|sql_injection|SQL injection attempt blocked|8|rt=1705315800000 deviceVendor=Block deviceProduct=AegisGate category=threat eventId=01HXYZ123456 src=192.168.1.100 suser=admin cs1=Initial Attack cs1Label=MitreTactic cs2=T1190 cs2Label=MitreTechnique
```

### LEEF Format (QRadar)

Log Event Extended Format for QRadar:

```
LEEF:2.0|Block|AegisGate|1.0|sql_injection|devTime=2024-01-15T10:30:00Z sev=high cat=threat eventName=SQL injection attempt blocked src=192.168.1.100 srcHost=web-server mitretactic=Initial Access mitreTechnique=T1190
```

### Syslog Format (RFC 5424)

Standard syslog format for LogRhythm and other platforms:

```
<14>1 2024-01-15T10:30:00Z hostname aegisgate - - [event@8732 id="01HXYZ123456" type="blocked_request" category="threat"][mitre@8732 tactic="Initial Access" technique="T1190"][entities@8732 ip="192.168.1.100"] SQL injection attempt blocked
```

---

## Event Schema

### Event Structure

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | string | Yes | Unique event identifier |
| timestamp | time | Yes | Event timestamp (ISO 8601) |
| source | string | No | Event source identifier |
| category | string | Yes | Event category |
| type | string | Yes | Event type |
| severity | string | Yes | Severity level |
| message | string | No | Human-readable description |
| attributes | map | No | Additional key-value pairs |
| entities | array | No | Related entities |
| mitre | object | No | MITRE ATT&CK mapping |
| compliance | array | No | Compliance mapping |

### Event Categories

| Category | Description |
|----------|-------------|
| authentication | Authentication events |
| authorization | Authorization events |
| access | Access control events |
| threat | Security threats detected |
| vulnerability | Vulnerability findings |
| compliance | Compliance violations |
| audit | Audit and compliance events |
| network | Network security events |
| application | Application security events |
| data_loss | Data loss prevention events |
| malware | Malware detection events |
| policy | Policy violation events |

### Severity Levels

| Level | Value | Description |
|-------|-------|-------------|
| critical | 5 | Critical security incident |
| high | 4 | High-priority security event |
| medium | 3 | Medium-priority event |
| low | 2 | Low-priority event |
| info | 1 | Informational event |

---

## Platform-Specific Guides

### Splunk Configuration

1. **Configure HTTP Event Collector (HEC):**
   - Navigate to Settings > Data Inputs > HTTP Event Collector
   - Create a new token with appropriate index access
   - Note the token value for configuration

2. **AegisGate Configuration:**
   ```yaml
   platforms:
     - platform: splunk
       enabled: true
       endpoint: "https://splunk.example.com:8088"
       auth:
         type: "api_key"
         api_key: "YOUR_HEC_TOKEN"
       settings:
         index: "security_events"
         source_type: "aegisgate:security"
         source: "aegisgate-gateway-01"
   ```

3. **Create Splunk Dashboard:**
   ```
   index=security_events source=aegisgate:security
   | stats count by severity, category
   | sort -count
   ```

### Elasticsearch Configuration

1. **Create Index Template:**
   ```json
   PUT _index_template/aegisgate-security
   {
     "index_patterns": ["aegisgate-security-*"],
     "template": {
       "mappings": {
         "properties": {
           "timestamp": { "type": "date" },
           "category": { "type": "keyword" },
           "type": { "type": "keyword" },
           "severity": { "type": "keyword" },
           "source_ip": { "type": "ip" },
           "mitre.tactic": { "type": "keyword" },
           "mitre.technique": { "type": "keyword" }
         }
       }
     }
   }
   ```

2. **AegisGate Configuration:**
   ```yaml
   platforms:
     - platform: elasticsearch
       enabled: true
       endpoint: "https://elasticsearch.example.com:9200"
       auth:
         type: "api_key"
         api_key: "YOUR_BASE64_API_KEY"
       settings:
         index: "aegisgate-security-{yyyy.MM.dd}"
   ```

### Microsoft Sentinel Configuration

1. **Create Log Analytics Workspace:**
   - Note the Workspace ID and Shared Key

2. **AegisGate Configuration:**
   ```yaml
   platforms:
     - platform: sentinel
       enabled: true
       settings:
         workspace_id: "YOUR_WORKSPACE_ID"
         shared_key: "YOUR_SHARED_KEY"
         log_type: "AegisGateSecurity"
   ```

3. **Create Sentinel Analytics Rule:**
   ```kusto
   AegisGateSecurity_CL
   | where severity_s == "critical" or severity_s == "high"
   | summarize count() by type_s, bin(TimeGenerated, 5m)
   | where count_ > 10
   ```

### QRadar Configuration

1. **Create Log Source:**
   - Log Source Type: Universal LEEF
   - Protocol: HTTPS
   - Log Source Identifier: aegisgate-gateway

2. **AegisGate Configuration:**
   ```yaml
   platforms:
     - platform: qradar
       enabled: true
       endpoint: "https://qradar.example.com"
       format: leef
       auth:
         type: "api_key"
         api_key: "YOUR_AUTH_TOKEN"
       settings:
         log_source_id: "LOG_SOURCE_ID"
         use_leef: true
         leef_version: "2.0"
   ```

---

## Best Practices

### Performance

1. **Enable Event Buffering:**
   ```yaml
   buffer:
     enabled: true
     max_size: 10000
     flush_interval: 5s
   ```

2. **Use Batch Processing:**
   ```yaml
   batch:
     enabled: true
     max_size: 100
     max_wait: 5s
   ```

3. **Configure Appropriate Retries:**
   ```yaml
   retry:
     enabled: true
     max_attempts: 3
     initial_backoff: 1s
     max_backoff: 30s
   ```

### Security

1. **Use TLS Always:**
   ```yaml
   tls:
     enabled: true
     insecure_skip_verify: false
   ```

2. **Store Credentials Securely:**
   - Use environment variables for secrets
   - Consider using a secrets manager
   - Rotate API keys regularly

3. **Filter Sensitive Events:**
   ```yaml
   filter:
     exclude_types:
       - "health_check"
       - "debug_event"
   ```

### High Availability

1. **Multiple SIEM Endpoints:**
   Configure multiple platforms for redundancy.

2. **Buffer Persistence:**
   ```yaml
   buffer:
     persist: true
     persist_dir: "/var/lib/aegisgate/siem-buffer"
   ```

3. **Monitor SIEM Health:**
   Enable health checks to detect connectivity issues.

---

## Troubleshooting

### Common Issues

#### Connection Refused

```
Error: siem [splunk] send: connection refused
```

**Solutions:**
- Verify endpoint URL is correct
- Check firewall rules
- Verify SIEM service is running
- Check TLS configuration

#### Authentication Failed

```
Error: siem [elasticsearch] send: status code 401
```

**Solutions:**
- Verify API key/token is valid
- Check credentials haven't expired
- Verify correct authentication type

#### Buffer Overflow

```
Error: siem [custom] buffer: buffer full
```

**Solutions:**
- Increase buffer size
- Check SIEM connectivity
- Review event volume

### Debugging

Enable debug logging:

```yaml
logging:
  level: debug
  siem: trace
```

Check SIEM manager statistics:

```go
stats := siem.GlobalManager().Stats()
log.Printf("Events received: %d", stats.EventsReceived)
log.Printf("Events sent: %d", stats.EventsSent)
log.Printf("Events dropped: %d", stats.EventsDropped)
log.Printf("Errors: %d", stats.Errors)
```

### Health Check Endpoint

```bash
curl http://localhost:8080/api/v1/siem/health
```

Response:
```json
{
  "status": "healthy",
  "platforms": {
    "splunk": {
      "status": "connected",
      "last_send": "2024-01-15T10:30:00Z",
      "events_sent": 15234,
      "errors": 0
    }
  }
}
```

---

## Compliance Mappings

Event types can be automatically mapped to compliance frameworks:

```yaml
compliance:
  enabled: true
  frameworks:
    - name: "SOC2"
      mappings:
        blocked_request: "CC6.1,CC6.6"
        auth_failure: "CC6.1,CC6.7"
```

Supported frameworks:
- SOC 2 Type II
- PCI-DSS
- HIPAA
- NIST 800-53
- ISO 27001
- GDPR

---

## Support

For additional support:
- Documentation: https://docs.aegisgate.security
- Issues: https://github.com/block/aegisgate/issues
- Community: https://community.aegisgate.security
