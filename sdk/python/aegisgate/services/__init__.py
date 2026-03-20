# SPDX-License-Identifier: MIT
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.

"""
Service modules for the AegisGate Python SDK.
"""

from aegisgate.services.auth import AuthService, AsyncAuthService
from aegisgate.services.proxy import ProxyService, AsyncProxyService
from aegisgate.services.compliance import ComplianceService, AsyncComplianceService
from aegisgate.services.siem import SIEMService, AsyncSIEMService
from aegisgate.services.webhook import WebhookService, AsyncWebhookService
from aegisgate.services.core import CoreService, AsyncCoreService

__all__ = [
    # Sync services
    "AuthService",
    "ProxyService",
    "ComplianceService",
    "SIEMService",
    "WebhookService",
    "CoreService",
    # Async services
    "AsyncAuthService",
    "AsyncProxyService",
    "AsyncComplianceService",
    "AsyncSIEMService",
    "AsyncWebhookService",
    "AsyncCoreService",
]
