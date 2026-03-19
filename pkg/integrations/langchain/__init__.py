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
AegisGate Python SDK - Enterprise AI API Security Platform

This SDK provides a Python interface to the AegisGate enterprise security platform
for securing AI chatbot/agent environments. It supports:

- MITRE ATLAS threat detection and compliance
- NIST AI RMF framework alignment
- OWASP LLM Top 10 security controls
- Real-time traffic inspection and policy enforcement
- Multi-framework compliance reporting (SOC2, HIPAA, PCI-DSS, GDPR, ISO 27001/42001)

Basic Usage:
    >>> from aegisgate import AegisGateClient
    >>> client = AegisGateClient(base_url="http://localhost:8080", api_key="your-api-key")
    >>> health = client.core.health()
    >>> print(f"Status: {health.status}")

Async Usage:
    >>> import asyncio
    >>> from aegisgate import AsyncAegisGateClient
    >>> async def main():
    ...     async with AsyncAegisGateClient(base_url="http://localhost:8080") as client:
    ...         health = await client.core.health()
    ...         print(f"Status: {health.status}")
    >>> asyncio.run(main())

Configuration:
    >>> from aegisgate import ClientConfig, TLSConfig
    >>> config = ClientConfig(
    ...     base_url="https://aegisgate.example.com",
    ...     api_key="your-api-key",
    ...     timeout=30.0,
    ...     tls=TLSConfig(verify=True, cert_file="/path/to/cert.pem")
    ... )
    >>> client = AegisGateClient(config=config)
"""

from aegisgate._version import (
    __version__,
    __version_info__,
    __author__,
    __license__,
    __copyright__,
)
from aegisgate.client import AegisGateClient
from aegisgate.async_client import AsyncAegisGateClient
from aegisgate.exceptions import (
    AegisGateError,
    AuthenticationError,
    APIError,
    ConfigurationError,
    ConnectionError,
    NetworkError,
    RateLimitError,
    ResourceNotFoundError,
    ServerError,
    TimeoutError,
    ValidationError,
)
from aegisgate.config import (
    ClientConfig,
    RESTConfig,
    GRPCConfig,
    AuthConfig,
    TLSConfig,
)
from aegisgate.models.auth import User, LoginResult, Session
from aegisgate.models.proxy import ProxyStats, Violation, ViolationFilter
from aegisgate.models.compliance import Framework, ComplianceCheck, ComplianceReport
from aegisgate.models.siem import SIEMConfig, SIEMEvent
from aegisgate.models.webhook import Webhook, WebhookEvent

__all__ = [
    # Version info
    "__version__",
    "__version_info__",
    "__author__",
    "__license__",
    "__copyright__",
    # Main clients
    "AegisGateClient",
    "AsyncAegisGateClient",
    # Configuration
    "ClientConfig",
    "RESTConfig",
    "GRPCConfig",
    "AuthConfig",
    "TLSConfig",
    # Exceptions
    "AegisGateError",
    "AuthenticationError",
    "APIError",
    "ConfigurationError",
    "ConnectionError",
    "NetworkError",
    "RateLimitError",
    "ResourceNotFoundError",
    "ServerError",
    "TimeoutError",
    "ValidationError",
    # Models - Auth
    "User",
    "LoginResult",
    "Session",
    # Models - Proxy
    "ProxyStats",
    "Violation",
    "ViolationFilter",
    # Models - Compliance
    "Framework",
    "ComplianceCheck",
    "ComplianceReport",
    # Models - SIEM
    "SIEMConfig",
    "SIEMEvent",
    # Models - Webhook
    "Webhook",
    "WebhookEvent",
]
