# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Service modules for the AegisGate Python SDK.

This package contains service classes that provide the main API functionality:
- auth: Authentication and user management
- proxy: Proxy statistics and violations
- compliance: Compliance framework and checks
- siem: SIEM integration
- webhook: Webhook management
- core: Core system operations
"""

from aegisgate.services.auth import AuthService
from aegisgate.services.proxy import ProxyService
from aegisgate.services.compliance import ComplianceService
from aegisgate.services.siem import SIEMService
from aegisgate.services.webhook import WebhookService
from aegisgate.services.core import CoreService

__all__ = [
    "AuthService",
    "ProxyService",
    "ComplianceService",
    "SIEMService",
    "WebhookService",
    "CoreService",
]
