# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Data models for the AegisGate Python SDK.

This package contains all data models used by the SDK, organized by service:
- auth: Authentication and user management models
- proxy: Proxy statistics and violation models
- compliance: Compliance framework and check models
- siem: SIEM integration models
- webhook: Webhook configuration models
"""

from aegisgate.models.auth import (
    User,
    UserCreate,
    UserUpdate,
    LoginResult,
    Session,
    SessionInfo,
    Role,
    Provider,
)
from aegisgate.models.proxy import (
    ProxyStats,
    ProxyConfig,
    Violation,
    ViolationFilter,
    ViolationSeverity,
    ThreatType,
    ScanResult,
)
from aegisgate.models.compliance import (
    Framework,
    ComplianceCheck,
    ComplianceReport,
    ComplianceStatus,
    Control,
    Finding,
)
from aegisgate.models.siem import (
    SIEMConfig,
    SIEMEvent,
    SIEMProvider,
    EventSeverity,
    EventCategory,
)
from aegisgate.models.webhook import (
    Webhook,
    WebhookCreate,
    WebhookUpdate,
    WebhookEvent,
    WebhookEventType,
    WebhookDelivery,
)

__all__ = [
    # Auth models
    "User",
    "UserCreate",
    "UserUpdate",
    "LoginResult",
    "Session",
    "SessionInfo",
    "Role",
    "Provider",
    # Proxy models
    "ProxyStats",
    "ProxyConfig",
    "Violation",
    "ViolationFilter",
    "ViolationSeverity",
    "ThreatType",
    "ScanResult",
    # Compliance models
    "Framework",
    "ComplianceCheck",
    "ComplianceReport",
    "ComplianceStatus",
    "Control",
    "Finding",
    # SIEM models
    "SIEMConfig",
    "SIEMEvent",
    "SIEMProvider",
    "EventSeverity",
    "EventCategory",
    # Webhook models
    "Webhook",
    "WebhookCreate",
    "WebhookUpdate",
    "WebhookEvent",
    "WebhookEventType",
    "WebhookDelivery",
]
