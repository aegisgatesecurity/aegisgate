# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
#
# This file contains proprietary trade secret information.
# Unauthorized reproduction, distribution, or reverse engineering is prohibited.
# =========================================================================

"""
AegisGate Python SDK

A comprehensive Python SDK for the AegisGate enterprise AI security platform.
Provides both REST and gRPC APIs with automatic retry, connection pooling,
and async support.

Basic Usage:
    from aegisgate import Client

    client = Client(base_url="http://localhost:8080", api_key="your-key")
    
    # Check health
    health = client.core.health()
    print(health)
    
    # Get proxy stats
    stats = client.proxy.get_stats()
    
    # Run compliance check
    result = client.compliance.run_check("hipaa")
    
    client.close()

Async Usage:
    from aegisgate import AsyncClient
    
    async def main():
        async with AsyncClient(base_url="http://localhost:8080") as client:
            health = await client.core.health()
            print(health)
"""

__version__ = "1.0.0"
__author__ = "AegisGate Security"

from aegisgate.client import Client, AsyncClient
from aegisgate.models import (
    Health,
    Version,
    Module,
    License,
    Violation,
    ComplianceResult,
    DetectionResult,
    AnomalyResult,
    SIEMConfig,
    SIEMEvent,
    Webhook,
    User,
)
from aegisgate.services import (
    AuthService,
    ProxyService,
    ComplianceService,
    SIEMService,
    WebhookService,
    CoreService,
)
from aegisgate.langchain import AegisGateCallback, AegisGateFilter

__all__ = [
    # Client
    "Client",
    "AsyncClient",
    # Models
    "Health",
    "Version",
    "Module",
    "License",
    "Violation",
    "ComplianceResult",
    "DetectionResult",
    "AnomalyResult",
    "SIEMConfig",
    "SIEMEvent",
    "Webhook",
    "User",
    # Services
    "AuthService",
    "ProxyService",
    "ComplianceService",
    "SIEMService",
    "WebhookService",
    "CoreService",
    # LangChain Integration
    "AegisGateCallback",
    "AegisGateFilter",
]