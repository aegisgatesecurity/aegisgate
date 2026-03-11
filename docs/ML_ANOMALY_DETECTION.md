# ML Anomaly Detection Configuration Guide

## Overview

AegisGate includes advanced ML-based anomaly detection capabilities that analyze traffic patterns, detect prompt injection attacks, and identify behavioral anomalies. This guide covers all ML-related configuration options.

## Quick Start

### Enable Basic ML Detection

```yaml
ml:
  enabled: true
  sensitivity: "medium"
```

### Enable Advanced ML Features

```yaml
ml:
  enabled: true
  sensitivity: "high"

ml_advanced:
  enable_prompt_injection: true
  enable_content_analysis: true
  enable_behavioral_analysis: true
```

## Configuration Reference

### Main ML Settings (`ml`)

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `enabled` | bool | `false` | Enable ML anomaly detection |
| `sensitivity` | string | `"medium"` | Detection sensitivity: `low`, `medium`, `high`, `paranoid` |
| `block_on_critical` | bool | `true` | Block requests with critical severity |
| `block_on_high` | bool | `false` | Block requests with high severity |
| `min_score_to_block` | float | `3.0` | Minimum z-score to trigger blocking |
| `sample_rate` | int | `100` | Percentage of requests to analyze (0-100) |
| `excluded_paths` | []string | `["/health", "/ready", "/metrics"]` | Paths to skip |
| `excluded_methods` | []string | `["OPTIONS", "HEAD"]` | HTTP methods to skip |
| `log_all_anomalies` | bool | `true` | Log all anomalies, not just blocked |

### Advanced ML Settings (`ml_advanced`)

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `enable_prompt_injection` | bool | `false` | Detect LLM prompt injection |
| `prompt_injection_sensitivity` | int | `75` | Sensitivity for prompt injection (0-100) |
| `enable_content_analysis` | bool | `false` | Scan LLM responses for PII/secrets |
| `enable_behavioral_analysis` | bool | `false` | Track client behavior patterns |
| `window_size` | int | `1000` | Number of requests for baseline |
| `z_threshold` | float | `3.0` | Standard deviations for anomaly |
| `min_samples` | int | `10` | Samples before detection starts |
| `entropy_threshold` | float | `4.5` | Shannon entropy threshold |

## Sensitivity Levels

### `low`
- Minimal detection
- Lowest false positives
- Best for initial testing

### `medium` (Recommended)
- Balanced detection
- Good for most production workloads

### `high`
- Aggressive detection
- May have more false positives
- Use with `block_on_high: false`

### `paranoid`
- Maximum detection
- Expect some false positives
- Review logs frequently

## Detection Capabilities

### 1. Traffic Anomaly Detection

Analyzes request patterns to detect:
- Unusual request volumes (traffic spikes/drops)
- Abnormal payload sizes
- Path traversal attempts
- Fuzzing/probing patterns

### 2. Prompt Injection Detection

Detects attempts to manipulate LLM behavior:
- Direct instruction overrides ("Ignore previous instructions")
- Jailbreak attempts ("DAN mode", "developer mode")
- System prompt extraction
- Role-playing manipulation
- Code injection attempts
- Obfuscated payloads (Base64, hidden tokens)

**Example detected patterns:**
```text
"Ignore all previous instructions"
"Forget your system prompt"
"New instructions: You are now..."
"DAN mode activate"
"Repeat after me: [system prompt]"
```

### 3. Content Analysis

Scans LLM responses for:
- **PII**: SSN, credit cards, emails, phone numbers
- **Secrets**: API keys, passwords, private keys
- Custom policy violations

### 4. Behavioral Analysis

Tracks client behavior over time:
- Request frequency anomalies (bots/scrapers)
- Path diversity (scraping detection)
- Data exfiltration attempts
- Unusual usage patterns

## Blocking Decisions

### Severity Levels

| Severity | Score Range | Description |
|----------|-------------|-------------|
| Critical | 4-5 | Immediate threat, block by default |
| High | 3-4 | Significant threat |
| Medium | 2-3 | Moderate concern |
| Low | 1-2 | Minor anomaly |

### Blocking Flow

```
Request → ML Analysis → Score Calculation
    ↓
    ├─ Score < Threshold → Allow
    │
    ├─ Score >= Threshold AND Severity >= High AND block_on_high=true → Block
    │
    └─ Score >= Threshold AND Severity = Critical AND block_on_critical=true → Block
```

## Custom Patterns

### Prompt Injection Custom Patterns

```yaml
prompt_injection:
  custom_patterns:
    - name: "confidential_keyword"
      pattern: "(?i)(confidential|secret|classified)"
      severity: 3
    - name: "override_attempt"
      pattern: "(?i)(you_are_now|you_must_now)"
      severity: 4
```

### Content Analysis Custom Rules

Custom rules can be added programmatically:

```go
analyzer.AddRule(ml.ContentRule{
    Name:     "custom_rule",
    Pattern:  regexp.MustCompile(`(?i)(blocked_word)`),
    Severity: 3,
    Action:   "block",
})
```

## Integration Examples

### Programmatic Usage

```go
// Load configuration
cfg, err := config.Load()
if err != nil {
    log.Fatal(err)
}

// Configure ML
cfg.ML = &config.MLConfig{
    Enabled:               true,
    Sensitivity:          "high",
    EnablePromptInjectionDetection: true,
    EnableContentAnalysis: true,
}

// Create ML-enabled proxy
proxy, err := proxy.NewProxyWithConfig(cfg)
if err != nil {
    log.Fatal(err)
}

// Get ML statistics
stats := proxy.GetMLStats()
fmt.Printf("ML Anomalies: %v\n", stats["middleware"])
```

### Direct ML Usage

```go
// Create prompt injection detector
detector := ml.NewPromptInjectionDetector(75)

// Analyze input
result := detector.Detect("Ignore all previous instructions")

if result.IsInjection {
    fmt.Printf("Prompt injection detected: %s\n", result.Explanation)
    fmt.Printf("Confidence: %.0f%%\n", result.Score)
    fmt.Printf("Matched patterns: %v\n", result.MatchedPatterns)
}

// Create content analyzer
analyzer := ml.NewContentAnalyzer()
contentResult := analyzer.Analyze(llmResponse)

if contentResult.IsViolation {
    fmt.Printf("Content violation: %v\n", contentResult.ViolationTypes)
}

// Create behavioral analyzer
behaviorAnalyzer := ml.NewBehavioralAnalyzer()
behaviorResult := behaviorAnalyzer.AnalyzeRequest(clientID, method, path, bytes)

if behaviorResult.IsAnomaly {
    fmt.Printf("Behavioral anomaly: %s\n", behaviorResult.AnomalyType)
}
```

## Monitoring

### Metrics

ML detection exposes the following metrics:

```
aegisgate_ml_total_requests_total
aegisgate_ml_analyzed_requests_total
aegisgate_ml_blocked_requests_total
aegisgate_ml_anomaly_detections_total
aegisgate_ml_prompt_injection_detections_total
aegisgate_ml_content_violations_total
aegisgate_ml_behavioral_anomalies_total
```

### Log Format

```
[ML_ANOMALY] Path: POST /api/chat | Client: 192.168.1.100 | Anomalies: 2 | Blocked: true
  - Type: prompt_injection | Severity: 5 | Score: 85.0
  - Pattern: ignore_previous
```

## Performance Considerations

### Sample Rate

For high-throughput environments, reduce sample rate:

```yaml
ml:
  sample_rate: 50  # Analyze 50% of requests
```

### Window Size

For very high traffic, increase window size:

```yaml
ml_advanced:
  window_size: 10000  # Larger baseline window
  min_samples: 100    # More samples before detection
```

### Latency Impact

| Configuration | Latency Impact |
|---------------|----------------|
| Basic ML | <1ms per request |
| +Prompt Injection | <2ms per request |
| +Content Analysis | <5ms per request |
| +Behavioral Analysis | <1ms per request |

## Troubleshooting

### Too Many False Positives

1. Lower sensitivity: `sensitivity: "low"`
2. Increase threshold: `min_score_to_block: 4.0`
3. Disable blocking: `block_on_critical: false`

### Not Detecting Known Attacks

1. Increase sensitivity: `sensitivity: "high"`
2. Lower threshold: `min_score_to_block: 2.5`
3. Add custom patterns

### High Memory Usage

1. Reduce window size: `window_size: 500`
2. Enable sampling: `sample_rate: 50`

## Examples

See `.mlconfig/` directory for complete examples:

- `aegisgate.minimal.yml` - Basic ML setup
- `aegisgate.example.yml` - Full-featured configuration
- `aegisgate.production.yml` - Production-optimized

## Related Documentation

- [Configuration Guide](./CONFIGURATION.md)
- [Security](./SECURITY.md)
- [Compliance](./ATLAS_FRAMEWORK.md)
- [API Reference](./PROXY_API.md)
