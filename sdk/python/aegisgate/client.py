# SPDX-License-Identifier: MIT
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.

"""
Main client for the AegisGate Python SDK.
"""

import os
from typing import Optional

from aegisgate.connection import ConnectionConfig, SyncConnection, AsyncConnection
from aegisgate.services import (
    AuthService,
    AsyncAuthService,
    ProxyService,
    AsyncProxyService,
    ComplianceService,
    AsyncComplianceService,
    SIEMService,
    AsyncSIEMService,
    WebhookService,
    AsyncWebhookService,
    CoreService,
    AsyncCoreService,
)


class Client:
    """
    Synchronous AegisGate client.
    
    Usage:
        client = Client(base_url="http://localhost:8080", api_key="your-key")
        health = client.core.health()
        print(health)
        client.close()
    
    Or as context manager:
        with Client(base_url="http://localhost:8080") as client:
            health = client.core.health()
    """
    
    def __init__(
        self,
        base_url: Optional[str] = None,
        api_key: Optional[str] = None,
        timeout: float = 30.0,
        max_retries: int = 3,
        verify_ssl: bool = True,
        **kwargs
    ):
        # Support environment variables
        if base_url is None:
            base_url = os.environ.get("AEGISGATE_BASE_URL", "http://localhost:8080")
        if api_key is None:
            api_key = os.environ.get("AEGISGATE_API_KEY")
        
        self._config = ConnectionConfig(
            base_url=base_url,
            api_key=api_key,
            timeout=timeout,
            max_retries=max_retries,
            verify_ssl=verify_ssl,
            custom_headers=kwargs.get("headers"),
            proxy=kwargs.get("proxy"),
        )
        
        self._connection = SyncConnection(self._config)
        self._token: Optional[str] = None
        
        # Initialize services
        self._auth = AuthService(self._connection)
        self._proxy = ProxyService(self._connection)
        self._compliance = ComplianceService(self._connection)
        self._siem = SIEMService(self._connection)
        self._webhook = WebhookService(self._connection)
        self._core = CoreService(self._connection)
    
    @property
    def auth(self) -> AuthService:
        """Access authentication service."""
        return self._auth
    
    @property
    def proxy(self) -> ProxyService:
        """Access proxy service."""
        return self._proxy
    
    @property
    def compliance(self) -> ComplianceService:
        """Access compliance service."""
        return self._compliance
    
    @property
    def siem(self) -> SIEMService:
        """Access SIEM service."""
        return self._siem
    
    @property
    def webhook(self) -> WebhookService:
        """Access webhook service."""
        return self._webhook
    
    @property
    def core(self) -> CoreService:
        """Access core service."""
        return self._core
    
    def set_token(self, token: str) -> None:
        """Set authentication token."""
        self._token = token
        self._connection.config.api_key = None
        self._connection.config.custom_headers["Authorization"] = f"Bearer {token}"
    
    def set_api_key(self, api_key: str) -> None:
        """Set API key for authentication."""
        self._config.api_key = api_key
        self._token = None
        if "Authorization" in self._connection.config.custom_headers:
            del self._connection.config.custom_headers["Authorization"]
    
    def close(self) -> None:
        """Close the client and release resources."""
        self._connection.close()
    
    def __enter__(self) -> "Client":
        """Enter context manager."""
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        """Exit context manager."""
        self.close()


class AsyncClient:
    """
    Asynchronous AegisGate client.
    
    Usage:
        async with AsyncClient(base_url="http://localhost:8080") as client:
            health = await client.core.health()
            print(health)
    
    Or manually:
        client = AsyncClient(base_url="http://localhost:8080")
        try:
            health = await client.core.health()
        finally:
            await client.close()
    """
    
    def __init__(
        self,
        base_url: Optional[str] = None,
        api_key: Optional[str] = None,
        timeout: float = 30.0,
        max_retries: int = 3,
        verify_ssl: bool = True,
        **kwargs
    ):
        # Support environment variables
        if base_url is None:
            base_url = os.environ.get("AEGISGATE_BASE_URL", "http://localhost:8080")
        if api_key is None:
            api_key = os.environ.get("AEGISGATE_API_KEY")
        
        self._config = ConnectionConfig(
            base_url=base_url,
            api_key=api_key,
            timeout=timeout,
            max_retries=max_retries,
            verify_ssl=verify_ssl,
            custom_headers=kwargs.get("headers"),
            proxy=kwargs.get("proxy"),
        )
        
        self._connection = AsyncConnection(self._config)
        self._token: Optional[str] = None
        
        # Initialize services
        self._auth = AsyncAuthService(self._connection)
        self._proxy = AsyncProxyService(self._connection)
        self._compliance = AsyncComplianceService(self._connection)
        self._siem = AsyncSIEMService(self._connection)
        self._webhook = AsyncWebhookService(self._connection)
        self._core = AsyncCoreService(self._connection)
    
    @property
    def auth(self) -> AsyncAuthService:
        """Access authentication service."""
        return self._auth
    
    @property
    def proxy(self) -> AsyncProxyService:
        """Access proxy service."""
        return self._proxy
    
    @property
    def compliance(self) -> AsyncComplianceService:
        """Access compliance service."""
        return self._compliance
    
    @property
    def siem(self) -> AsyncSIEMService:
        """Access SIEM service."""
        return self._siem
    
    @property
    def webhook(self) -> AsyncWebhookService:
        """Access webhook service."""
        return self._webhook
    
    @property
    def core(self) -> AsyncCoreService:
        """Access core service."""
        return self._core
    
    def set_token(self, token: str) -> None:
        """Set authentication token."""
        self._token = token
        self._connection.config.api_key = None
        self._connection.config.custom_headers["Authorization"] = f"Bearer {token}"
    
    def set_api_key(self, api_key: str) -> None:
        """Set API key for authentication."""
        self._config.api_key = api_key
        self._token = None
        if "Authorization" in self._connection.config.custom_headers:
            del self._connection.config.custom_headers["Authorization"]
    
    async def close(self) -> None:
        """Close the client and release resources."""
        await self._connection.close()
    
    async def __aenter__(self) -> "AsyncClient":
        """Enter async context manager."""
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb) -> None:
        """Exit async context manager."""
        await self.close()
