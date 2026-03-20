"""
Example 4: Working with AegisGate Services

This example demonstrates usage of various AegisGate services:
- Auth: User management
- Proxy: Request inspection
- Compliance: Framework management
- SIEM: Security event integration
- Webhook: Event notifications
"""

import os
import asyncio
from datetime import datetime

from aegisgate import Client, AsyncClient
from aegisgate.models import (
    SIEMProvider,
    LicenseType,
    ViolationType,
)


API_KEY = os.environ.get("AEGISGATE_API_KEY", "your-api-key-here")


# ============================================================
# Auth Service Examples
# ============================================================

def example_auth_service():
    """Working with authentication and user management."""
    
    client = Client(api_key=API_KEY)
    
    # Login
    print("=== Auth Service: Login ===")
    user = client.auth.login(username="admin", password="secret")
    print(f"Logged in as: {user.username}")
    print(f"Email: {user.email}")
    print(f"Token: {user.token[:20]}..." if hasattr(user, 'token') else "")
    print()
    
    # List users (requires admin privileges)
    print("=== Auth Service: List Users ===")
    users = client.auth.list_users()
    for u in users[:5]:  # Show first 5 users
        print(f"  - {u.username} ({u.email})")
    print()
    
    # Create user
    print("=== Auth Service: Create User ===")
    new_user = client.auth.create_user(
        username="newuser",
        email="newuser@example.com",
        password="secure_password_123",
        role="viewer"
    )
    print(f"Created user: {new_user.username}")
    print()
    
    # Update user
    print("=== Auth Service: Update User ===")
    updated_user = client.auth.update_user(
        user_id=new_user.id,
        email="newemail@example.com"
    )
    print(f"Updated email: {updated_user.email}")
    print()
    
    # Delete user
    print("=== Auth Service: Delete User ===")
    client.auth.delete_user(user_id=new_user.id)
    print("User deleted")
    print()
    
    # Logout
    client.auth.logout()
    print("Logged out")


# ============================================================
# Proxy Service Examples
# ============================================================

def example_proxy_service():
    """Working with proxy and request inspection."""
    
    client = Client(api_key=API_KEY)
    
    # Get proxy statistics
    print("=== Proxy Service: Statistics ===")
    stats = client.proxy.get_stats()
    print(f"Total requests: {stats.total_requests}")
    print(f"Blocked: {stats.blocked_requests}")
    print(f"Allowed: {stats.successful_requests}")
    print(f"Average latency: {stats.avg_latency_ms}ms")
    print()
    
    # Inspect request
    print("=== Proxy Service: Inspect Request ===")
    result = client.proxy.inspect_request(
        content="SELECT * FROM users WHERE id = 1",
        content_type="sql"
    )
    print(f"Has violations: {result.has_violations}")
    if result.has_violations:
        for v in result.violations:
            print(f"  - {v.type}: {v.message}")
    print()
    
    # Detect anomalies
    print("=== Proxy Service: Anomaly Detection ===")
    anomalies = client.proxy.detect_anomalies(
        timeframe="24h",
        sensitivity="high"
    )
    print(f"Anomalies found: {len(anomalies)}")
    for anomaly in anomalies[:3]:
        print(f"  - {anomaly.type}: {anomaly.description}")
    print()
    
    # List violations
    print("=== Proxy Service: List Violations ===")
    violations = client.proxy.list_violations(
        violation_types=[ViolationType.PROMPT_INJECTION, ViolationType.TOXIC_CONTENT],
        limit=10
    )
    for v in violations:
        print(f"  [{v.severity}] {v.type}: {v.message}")
    print()
    
    # Configure content filter
    print("=== Proxy Service: Configure Content Filter ===")
    config = client.proxy.configure_content_filter(
        enabled=True,
        rules=[
            {"type": "prompt_injection", "action": "block"},
            {"type": "pii_exposure", "action": "redact"},
            {"type": "toxic_content", "action": "block", "threshold": 0.8},
        ]
    )
    print(f"Content filter configured: {config.enabled}")


# ============================================================
# Compliance Service Examples
# ============================================================

def example_compliance_service():
    """Working with compliance frameworks and checks."""
    
    client = Client(api_key=API_KEY)
    
    # List available frameworks
    print("=== Compliance Service: List Frameworks ===")
    frameworks = client.compliance.list_frameworks()
    for fw in frameworks:
        print(f"  - {fw.name} ({fw.id})")
    print()
    
    # Get specific framework
    print("=== Compliance Service: Get Framework ===")
    framework = client.compliance.get_framework(framework_id="soc2")
    print(f"Framework: {framework.name}")
    print(f"Controls: {len(framework.controls)}")
    print()
    
    # Run compliance check
    print("=== Compliance Service: Run Check ===")
    result = client.compliance.run_check(framework_id="soc2")
    print(f"Status: {result.status}")
    print(f"Passed: {result.passed_controls}")
    print(f"Failed: {result.failed_controls}")
    print()
    
    # Get controls
    print("=== Compliance Service: Get Controls ===")
    controls = client.compliance.get_controls(framework_id="soc2")
    for control in controls[:5]:
        status = "✓" if control.passed else "✗"
        print(f"  {status} {control.id}: {control.name}")
    print()
    
    # Generate report
    print("=== Compliance Service: Generate Report ===")
    report = client.compliance.generate_report(
        framework_id="soc2",
        format="pdf"
    )
    print(f"Report URL: {report.download_url}")
    print()
    
    # Schedule recurring check
    print("=== Compliance Service: Schedule Check ===")
    schedule = client.compliance.schedule_check(
        framework_id="soc2",
        cron_expression="0 0 * * *"  # Daily at midnight
    )
    print(f"Scheduled: {schedule.id} - Next run: {schedule.next_run}")


# ============================================================
# SIEM Service Examples
# ============================================================

def example_siem_service():
    """Working with SIEM integrations."""
    
    client = Client(api_key=API_KEY)
    
    # List integrations
    print("=== SIEM Service: List Integrations ===")
    integrations = client.siem.list_integrations()
    for integ in integrations:
        print(f"  - {integ.name} ({integ.type}): {integ.status}")
    print()
    
    # Create new integration
    print("=== SIEM Service: Create Integration ===")
    integration = client.siem.create_integration(
        name="Splunk Production",
        type=SIEMProvider.SPLUNK,
        config={
            "url": "https://splunk.example.com:8089",
            "token": "splunk-api-token",
            "index": "aegisgate-security",
        }
    )
    print(f"Created integration: {integration.id}")
    print(f"Status: {integration.status}")
    print()
    
    # Test integration
    print("=== SIEM Service: Test Integration ===")
    test_result = client.siem.test_integration(integration_id=integration.id)
    print(f"Test passed: {test_result.success}")
    if not test_result.success:
        print(f"Error: {test_result.error}")
    print()
    
    # Send security event
    print("=== SIEM Service: Send Event ===")
    event = client.siem.send_event(
        event_type="violation_detected",
        severity="high",
        details={
            "violation_type": "prompt_injection",
            "source": "api-gateway",
            "timestamp": datetime.utcnow().isoformat(),
            "user_id": "user-123",
            "request_id": "req-456",
        }
    )
    print(f"Event sent: {event.id}")
    print()
    
    # Send batch events
    print("=== SIEM Service: Batch Events ===")
    events = [
        {
            "event_type": "anomaly_detected",
            "severity": "medium",
            "details": {"score": 0.85}
        },
        {
            "event_type": "pii_exposure",
            "severity": "high",
            "details": {"type": "ssn"}
        },
    ]
    batch_result = client.siem.send_batch_events(events)
    print(f"Events sent: {batch_result.count}")
    print()
    
    # Get events
    print("=== SIEM Service: Get Events ===")
    events = client.siem.get_events(
        start_time="2024-01-01T00:00:00Z",
        end_time="2024-01-31T23:59:59Z",
        limit=10
    )
    for event in events[:3]:
        print(f"  - {event.event_type}: {event.severity}")


# ============================================================
# Webhook Service Examples
# ============================================================

def example_webhook_service():
    """Working with webhooks for event notifications."""
    
    client = Client(api_key=API_KEY)
    
    # List webhooks
    print("=== Webhook Service: List Webhooks ===")
    webhooks = client.webhook.list_webhooks()
    for wh in webhooks:
        print(f"  - {wh.name}: {wh.url} ({wh.status})")
    print()
    
    # Create webhook
    print("=== Webhook Service: Create Webhook ===")
    webhook = client.webhook.create_webhook(
        name="Security Alerts",
        url="https://hooks.example.com/aegisgate",
        events=["violation_detected", "anomaly_detected", "compliance_failed"],
        secret="webhook-secret-key"
    )
    print(f"Created webhook: {webhook.id}")
    print()
    
    # Test webhook
    print("=== Webhook Service: Test Webhook ===")
    test_result = client.webhook.test_webhook(webhook_id=webhook.id)
    print(f"Test status: {test_result.status}")
    print(f"Response time: {test_result.response_time_ms}ms")
    print()
    
    # Get deliveries
    print("=== Webhook Service: Get Deliveries ===")
    deliveries = client.webhook.get_deliveries(webhook_id=webhook.id)
    for delivery in deliveries[:5]:
        status = "✓" if delivery.success else "✗"
        print(f"  {status} {delivery.event_type}: {delivery.status}")
    print()
    
    # Retry failed delivery
    print("=== Webhook Service: Retry Delivery ===")
    if deliveries and not deliveries[0].success:
        result = client.webhook.retry_delivery(delivery_id=deliveries[0].id)
        print(f"Retry status: {result.status}")
    print()
    
    # Ping webhook
    print("=== Webhook Service: Ping Webhook ===")
    ping_result = client.webhook.ping_webhook(webhook_id=webhook.id)
    print(f"Ping response: {ping_result.status}")
    print()
    
    # Update webhook
    print("=== Webhook Service: Update Webhook ===")
    updated = client.webhook.update_webhook(
        webhook_id=webhook.id,
        events=["violation_detected", "compliance_failed", "pii_exposure"]
    )
    print(f"Updated events: {updated.events}")
    print()
    
    # Delete webhook
    print("=== Webhook Service: Delete Webhook ===")
    client.webhook.delete_webhook(webhook_id=webhook.id)
    print("Webhook deleted")


# ============================================================
# Core Service Examples
# ============================================================

def example_core_service():
    """Working with core system functions."""
    
    client = Client(api_key=API_KEY)
    
    # Health check
    print("=== Core Service: Health ===")
    health = client.core.health()
    print(f"Status: {health.status}")
    print()
    
    # Detailed health
    print("=== Core Service: Health Detailed ===")
    detailed = client.core.health_detailed()
    print(f"Status: {detailed.status}")
    for module in detailed.modules:
        status_icon = "✓" if module.status == "healthy" else "✗"
        print(f"  {status_icon} {module.name}: {module.status}")
    print()
    
    # Version
    print("=== Core Service: Version ===")
    version = client.core.version()
    print(f"Version: {version.version}")
    print(f"Build: {version.build}")
    print(f"Commit: {version.commit}")
    print()
    
    # Get license
    print("=== Core Service: License ===")
    license = client.core.get_license()
    print(f"Status: {license.status}")
    print(f"Type: {license.type}")
    print(f"Expires: {license.expires_at}")
    print()
    
    # Validate license
    print("=== Core Service: Validate License ===")
    validation = client.core.validate_license(license_key="YOUR-LICENSE-KEY")
    print(f"Valid: {validation.valid}")
    print()
    
    # Get metrics
    print("=== Core Service: Metrics ===")
    metrics = client.core.get_metrics()
    print(f"Requests: {metrics.requests_total}")
    print(f"Errors: {metrics.errors_total}")
    print(f"Avg latency: {metrics.latency_avg_ms}ms")
    print()
    
    # Get logs
    print("=== Core Service: Logs ===")
    logs = client.core.get_logs(level="error", limit=5)
    for log in logs:
        print(f"  [{log.timestamp}] {log.level}: {log.message}")
    print()
    
    # Get environment info
    print("=== Core Service: Environment ===")
    env = client.core.get_environment()
    print(f"Environment: {env.environment}")
    print(f"Region: {env.region}")
    print(f"Data center: {env.datacenter}")


# ============================================================
# Run Examples
# ============================================================

if __name__ == "__main__":
    print("=" * 60)
    print("AegisGate Service Examples")
    print("=" * 60)
    
    # Uncomment to run examples:
    
    # example_auth_service()
    # example_proxy_service()
    # example_compliance_service()
    # example_siem_service()
    # example_webhook_service()
    # example_core_service()
    
    print("\nNote: Uncomment examples to run with actual API credentials")