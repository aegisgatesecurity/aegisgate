# RFC 5424 Compliance Implementation

## Overview

This document describes the RFC 5424 compliant syslog implementation in AegisGate.

**Reference**: [RFC 5424 - The Syslog Protocol](https://datatracker.ietf.org/doc/html/rfc5424)

## RFC 5424 Message Format

```
<PRI>VERSION TIMESTAMP HOSTNAME APP-NAME PROCID MSGID SD MSG
```

### Example Output

```
<134>1 2026-03-06T15:07:26.205Z chaos aegisgate 38772 AUTH_SUCCESS [aegisgate@32473 eventId="evt-001" eventType="authentication"] Login successful
```

## Implementation Details

### 1. RFC5424Message (`pkg/siem/rfc5424.go`)

The core message structure includes:

| Field | RFC 5424 Element | Status | Notes |
|-------|-----------------|--------|-------|
| Priority | PRI | ✅ | facility * 8 + severity |
| Version | VERSION | ✅ | Always 1 |
| Timestamp | TIMESTAMP | ✅ | RFC3339 format with microseconds |
| Hostname | HOSTNAME | ✅ | From os.Hostname() or event |
| AppName | APP-NAME | ✅ | "aegisgate" |
| ProcID | PROCID | ✅ | Actual process ID |
| MsgID | MSGID | ✅ | 40+ event types mapped |
| StructuredData | SD | ✅ | Enterprise event data |
| Message | MSG | ✅ | UTF-8 message content |

### 2. MSGID Implementation

Over 40 MSGID values are defined for different event types:

**Authentication:**
- `AUTH_SUCCESS`, `AUTH_FAILURE`, `SESSION_START`, `SESSION_END`
- `TOKEN_REFRESH`, `TOKEN_REVOKE`, `AUTH_LOGOUT`

**Authorization:**
- `AUTHZ_SUCCESS`, `AUTHZ_FAILURE`, `AUTHZ_DENIED`

**Request Handling:**
- `REQUEST_ALLOWED`, `REQUEST_BLOCKED`, `REQUEST_DROPPED`, `REQUEST_THROTTLED`

**Security:**
- `THREAT_DETECTED`, `INTRUSION_ATTEMPT`, `MALWARE_DETECTED`
- `ANOMALY_DETECTED`, `POLICY_VIOLATION`, `RATE_LIMIT_EXCEEDED`

**Proxy:**
- `PROXY_REQUEST`, `PROXY_RESPONSE`, `PROXY_ERROR`
- `MITM_DETECTED`, `TLS_ERROR`

**System:**
- `SYSTEM_START`, `SYSTEM_STOP`, `SYSTEM_ERROR`
- `COMPONENT_FAILURE`, `HEALTH_CHECK`, `METRICS_PUBLISH`

### 3. Structured Data (SD)

Custom SD-ID: `aegisgate@32473`

Parameters include:
- `eventId` - Unique event identifier
- `eventType` - Event type (authentication, threat, etc.)
- `action` - Action taken (block, allow, drop)
- `srcIp` - Source IP address
- `dstIp` - Destination IP address
- `user` - User identifier
- `clientId` - Client/application identifier
- `threatType` - Type of threat detected
- `threatLevel` - Threat severity level
- `pattern` - Matched pattern (if any)
- `framework` - Compliance framework
- `control` - Compliance control ID

### 4. Severity Mapping

| AegisGate Severity | RFC 5424 Severity | Value |
|-----------------|-------------------|-------|
| emergency, emerg, crit, critical, fatal | Critical | 2 |
| alert | Alert | 1 |
| error, err | Error | 3 |
| warning, warn | Warning | 4 |
| notice | Notice | 5 |
| info, informational | Informational | 6 |
| debug, trace, verbose | Debug | 7 |

### 5. Facility Configuration

Default: `local0` (16)

Can be configured via `SyslogOptions`:
```go
type SyslogOptions struct {
    Facility int      // Syslog facility (default: local0)
    AppName  string   // Application name (default: aegisgate)
    Hostname string   // Hostname override
}
```

## Usage

### Basic Usage

```go
formatter := siem.NewSyslogFormatter(siem.PlatformSyslog, siem.SyslogOptions{
    Facility: siem.SyslogFacilityLocal0,
})

event := &siem.Event{
    ID:        "evt-001",
    Type:      "authentication",
    Message:   "User logged in successfully",
    Severity:  siem.SeverityInfo,
    SourceIP:  "192.168.1.100",
    User:      "admin",
}

result, err := formatter.FormatRFC5424(event)
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)
```

### Output

```
<134>1 2026-03-06T15:07:26.205Z hostname aegisgate 12345 AUTH_SUCCESS [aegisgate@32473 eventId="evt-001" eventType="authentication" srcIp="192.168.1.100" user="admin"] User logged in successfully
```

## Benchmarks

```
BenchmarkRFC5424Message_Build-40     	  178677	      6246 ns/op	    1440 B/op	      28 allocs/op
BenchmarkConvertEventToRFC5424-40    	   87471	     15868 ns/op	     824 B/op	      10 allocs/op
```

## Enterprise Integration

### Syslog Servers Tested

- Generic syslog (UDP/TCP)
- RSyslog
- Syslog-ng
- Graylog
- Splunk
- LogRhythm
- AWS CloudWatch Logs

### Transport Options

| Protocol | Port | TLS Support |
|----------|------|-------------|
| UDP | 514 | No |
| TCP | 514 | Optional |
| TLS | 514 | Required |
| RELP | 20514 | Yes |

## Files Added

- `pkg/siem/rfc5424.go` - Core RFC 5424 implementation
- `pkg/siem/rfc5424_test.go` - Comprehensive tests and benchmarks

## Files Modified

- `pkg/siem/types.go` - Added RFC 5424 specific fields:
  - `Action`
  - `SourceIP`
  - `Destination`
  - `User`
  - `ClientID`
  - `ThreatType`
  - `ThreatLevel`
  - `Pattern`
  - `ComplianceFramework`
  - `ComplianceControl`

## Compliance

This implementation meets all RFC 5424 requirements:

- ✅ PRI (Priority)
- ✅ VERSION (always 1)
- ✅ TIMESTAMP (RFC3339)
- ✅ HOSTNAME
- ✅ APP-NAME
- ✅ PROCID
- ✅ MSGID
- ✅ STRUCTURED-DATA
- ✅ MSG (UTF-8)
- ✅ Special character escaping
- ✅ NILVALUE handling
