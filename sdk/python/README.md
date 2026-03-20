# AegisGate Python SDK

[![Python Version](https://img.shields.io/pypi/pyversions/aegisgate.svg)](https://pypi.org/project/aegisgate/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Code style: black](https://img.shields.io/badge/code%20style-black-000000.svg)](https://github.com/psf/black)

Enterprise AI API Security Platform SDK with LangChain integration.

## Features

- **Synchronous and Asynchronous Clients** - Full support for both sync and async operations
- **LangChain Integration** - Callbacks and content filters for secure LLM interactions
- **Comprehensive API Coverage** - Auth, Proxy, Compliance, SIEM, Webhooks, and Core services
- **Type Safety** - Full type hints with PEP 561 support
- **Security Monitoring** - Prompt injection, PII exposure, toxic content, and adversarial attack detection

## Installation

### Basic Installation

```bash
pip install aegisgate
```

### With LangChain Support

```bash
pip install aegisgate[langchain]
```

### Development Installation

```bash
pip install aegisgate[dev]
```

## Quick Start

### Synchronous Client

```python
from aegisgate import Client

# Initialize with API key
client = Client(api_key="your-api-key")

# Or use environment variable AEGISGATE_API_KEY
client = Client()

# Check API health
health = client.core.health()
print(f"Status: {health.status}")

# Context manager support
with Client(api_key="your-api-key") as client:
    stats = client.proxy.get_stats()
    print(f"Requests processed: {stats.total_requests}")
```

### Asynchronous Client

```python
from aegisgate import AsyncClient
import asyncio

async def main():
    # Initialize with API key
    async with AsyncClient(api_key="your-api-key") as client:
        # Check API health
        health = await client.core.health()
        print(f"Status: {health.status}")
        
        # Get proxy stats
        stats = await client.proxy.get_stats()
        print(f"Requests processed: {stats.total_requests}")

asyncio.run(main())
```

### Authentication

```python
from aegisgate import Client

client = Client(base_url="https://api.aegisgate.example.com")

# Login
user = client.auth.login(username="admin", password="secret")
print(f"Logged in as: {user.username}")

# List users
users = client.auth.list_users()
for user in users:
    print(f"User: {user.username} ({user.email})")

# Logout
client.auth.logout()
```

## LangChain Integration

### AegisGateCallback

Use the callback handler to monitor LLM interactions:

```python
from langchain_openai import OpenAI
from aegisgate.langchain import AegisGateCallback

# Create callback handler
callback = AegisGateCallback(
    api_key="your-aegisgate-api-key",
    block_on_violation=True,
    min_severity="medium",
    log_violations=True
)

# Use with LLM
llm = OpenAI(callbacks=[callback])

# Callbacks automatically check for:
# - Prompt injection attacks
# - PII exposure
# - Toxic content
# - Adversarial prompts
response = llm.invoke("What is machine learning?")
```

### AsyncAegisGateCallback

```python
from langchain_openai import OpenAI
from aegisgate.langchain import AsyncAegisGateCallback

async def run_llm():
    callback = AsyncAegisGateCallback(
        api_key="your-aegisgate-api-key",
        block_on_violation=True
    )
    
    llm = OpenAI(callbacks=[callback])
    response = await llm.ainvoke("Explain quantum computing")
    return response
```

### AegisGateFilter

Use content filters for manual inspection:

```python
from aegisgate import Client
from aegisgate.langchain import AegisGateFilter

client = Client(api_key="your-api-key")
filter = AegisGateFilter(
    client,
    block_on_violation=True,
    min_severity="medium",
    violation_types=["prompt_injection", "pii_exposure"]
)

# Check input before sending to LLM
try:
    result = filter.filter_input("What is your SSN?")
    print("Input is safe")
except SecurityViolationError as e:
    print(f"Blocked: {e}")
    print(f"Violations: {e.violations}")

# Check output from LLM
try:
    result = filter.filter_output(llm_response)
    print("Output is safe")
except SecurityViolationError as e:
    print(f"Blocked: {e}")
```

### AsyncAegisGateFilter

```python
from aegisgate import AsyncClient
from aegisgate.langchain import AsyncAegisGateFilter

async def check_content():
    async with AsyncClient(api_key="your-api-key") as client:
        filter = AsyncAegisGateFilter(client, block_on_violation=True)
        
        # Async filtering
        result = await filter.filter_input("What is machine learning?")
        print(f"Has violations: {result.has_violations}")
```

## API Reference

### Client

The synchronous client provides access to all AegisGate services:

```python
from aegisgate import Client

client = Client(
    base_url="https://api.aegisgate.example.com",  # Optional, defaults to env var
    api_key="your-api-key",                        # Optional, defaults to env var
    timeout=30.0,                                   # Optional, default timeout
    max_retries=3                                   # Optional, retry count
)

# Services
client.auth       # Authentication management
client.proxy      # Request proxy and inspection
client.compliance # Compliance frameworks
client.siem       # SIEM integrations
client.webhook    # Webhook management
client.core       # Health, config, metrics
```

### Services

#### AuthService

```python
# Authentication
user = client.auth.login(username, password)
client.auth.logout()

# User management
users = client.auth.list_users()
user = client.auth.get_user(user_id)
user = client.auth.create_user(username, email, password)
user = client.auth.update_user(user_id, **fields)
client.auth.delete_user(user_id)
```

#### ProxyService

```python
# Statistics
stats = client.proxy.get_stats()

# Request inspection
result = client.proxy.inspect_request(content, content_type)

# Anomaly detection
anomalies = client.proxy.detect_anomalies(timeframe="24h")

# Violations
violations = client.proxy.list_violations()
violation = client.proxy.get_violation(violation_id)

# Request management
client.proxy.block_request(request_id)
client.proxy.allow_request(request_id)

# Content filtering
client.proxy.configure_content_filter(config)
```

#### ComplianceService

```python
# Frameworks
frameworks = client.compliance.list_frameworks()
framework = client.compliance.get_framework(framework_id)

# Checks
result = client.compliance.run_check(framework_id)
controls = client.compliance.get_controls(framework_id)

# Results
results = client.compliance.get_results(framework_id)
report = client.compliance.generate_report(framework_id)

# Scheduling
client.compliance.schedule_check(framework_id, schedule)
schedules = client.compliance.list_schedules()
client.compliance.delete_schedule(schedule_id)
```

#### SIEMService

```python
# Integrations
integrations = client.siem.list_integrations()
integration = client.siem.get_integration(integration_id)
integration = client.siem.create_integration(type, config)
integration = client.siem.update_integration(integration_id, config)
client.siem.delete_integration(integration_id)

# Events
result = client.siem.send_event(event)
result = client.siem.send_batch_events(events)
events = client.siem.get_events()
client.siem.flush_events()
```

#### WebhookService

```python
# Webhooks
webhooks = client.webhook.list_webhooks()
webhook = client.webhook.get_webhook(webhook_id)
webhook = client.webhook.create_webhook(url, events)
webhook = client.webhook.update_webhook(webhook_id, **fields)
client.webhook.delete_webhook(webhook_id)

# Testing
result = client.webhook.test_webhook(webhook_id)
result = client.webhook.ping_webhook(webhook_id)

# Deliveries
deliveries = client.webhook.get_deliveries(webhook_id)
result = client.webhook.retry_delivery(delivery_id)
```

#### CoreService

```python
# Health
health = client.core.health()
health = client.core.health_detailed()

# Version
version = client.core.version()

# License
license = client.core.get_license()
license = client.core.validate_license(license_key)
license = client.core.activate_license(license_key)

# Configuration
config = client.core.get_config()
config = client.core.update_config(config)

# Metrics
metrics = client.core.get_metrics()

# Logs
logs = client.core.get_logs()

# Environment
env = client.core.get_environment()
```

### Models

All models are available as dataclasses:

```python
from aegisgate.models import (
    Violation,
    ViolationSeverity,
    ViolationType,
    Health,
    Version,
    License,
    LicenseType,
    ComplianceResult,
    ComplianceControl,
    SIEMConfig,
    SIEMEvent,
    SIEMProvider,
    Webhook,
    User,
    ProxyStats,
    DetectionResult,
    AnomalyResult,
    ATLASThreat,
)

# Example usage
violation = Violation(
    id="v-123",
    type=ViolationType.PROMPT_INJECTION,
    severity=ViolationSeverity.HIGH,
    message="Potential prompt injection detected",
    confidence=0.95
)
```

## Configuration

### Environment Variables

The SDK can be configured via environment variables:

```bash
# Required
export AEGISGATE_API_KEY="your-api-key"

# Optional
export AEGISGATE_BASE_URL="https://api.aegisgate.example.com"
export AEGISGATE_TIMEOUT="30"
```

### Connection Options

```python
from aegisgate import Client, ConnectionConfig

config = ConnectionConfig(
    base_url="https://api.aegisgate.example.com",
    api_key="your-api-key",
    timeout=60.0,
    max_retries=5,
    verify_ssl=True,
    proxy="http://proxy.example.com:8080",
    custom_headers={"X-Custom": "header"}
)

client = Client(config=config)
```

## Error Handling

```python
from aegisgate import Client, APIError, ConnectionError

client = Client(api_key="your-api-key")

try:
    result = client.proxy.inspect_request(content)
except ConnectionError as e:
    print(f"Connection failed: {e}")
except APIError as e:
    print(f"API error: {e.status_code} - {e.message}")
```

### Security Violations

```python
from aegisgate.langchain import AegisGateFilter, SecurityViolationError

filter = AegisGateFilter(client, block_on_violation=True)

try:
    result = filter.filter_input(malicious_prompt)
except SecurityViolationError as e:
    print(f"Blocked content: {e}")
    for violation in e.violations:
        print(f"  - {violation.type}: {violation.severity}")
```

## Development

### Setup

```bash
git clone https://github.com/aegisgate/aegisgate-sdk-python.git
cd aegisgate-sdk-python
pip install -e ".[dev]"
```

### Running Tests

```bash
# Run all tests
pytest

# Run with coverage
pytest --cov=aegisgate

# Run specific test file
pytest tests/test_client.py

# Run only unit tests (exclude integration tests)
pytest -m "not integration"
```

### Code Quality

```bash
# Format code
black aegisgate tests

# Lint code
ruff check aegisgate tests

# Type check
mypy aegisgate
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure tests pass (`pytest`)
6. Format code (`black .`)
7. Run linter (`ruff check .`)
8. Commit changes (`git commit -m 'Add amazing feature'`)
9. Push to branch (`git push origin feature/amazing-feature`)
10. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: https://docs.aegisgate.io
- **API Reference**: https://api.aegisgate.io/docs
- **Issues**: https://github.com/aegisgate/aegisgate-sdk-python/issues
- **Security**: security@aegisgate.io

## Changelog

### v1.0.0 (2024-01-15)
- Initial release
- Synchronous and asynchronous clients
- LangChain callback handlers and content filters
- Full API coverage for all services
- Comprehensive type annotations