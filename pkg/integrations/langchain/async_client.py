# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
AegisGate Python SDK - Async Client

This module provides the AsyncAegisGateClient class for asynchronous
communication with the AegisGate API.
"""

from __future__ import annotations

import json
import logging
from typing import Optional, Dict, Any, Union, AsyncIterator
from urllib.parse import urlencode

import httpx

from aegisgate._version import __version__
from aegisgate.config import ClientConfig
from aegisgate.exceptions import (
    AegisGateError,
    APIError,
    AuthenticationError,
    ConfigurationError,
    ConnectionError,
    RateLimitError,
    ResourceNotFoundError,
    ServerError,
    TimeoutError,
    ValidationError,
    SSLError,
)
from aegisgate.services.auth import AuthService
from aegisgate.services.proxy import ProxyService
from aegisgate.services.compliance import ComplianceService
from aegisgate.services.siem import SIEMService
from aegisgate.services.webhook import WebhookService
from aegisgate.services.core import CoreService

logger = logging.getLogger(__name__)


class AsyncAegisGateClient:
    """
    Async client for interacting with the AegisGate API.
    
    This client provides the same functionality as AegisGateClient but
    with full async/await support for non-blocking operations.
    
    Example:
        >>> import asyncio
        >>> from aegisgate import AsyncAegisGateClient
        >>> 
        >>> async def main():
        ...     async with AsyncAegisGateClient(
        ...         base_url="http://localhost:8080",
        ...         api_key="AG-your-api-key"
        ...     ) as client:
        ...         # Check health
        ...         health = await client.core.health()
        ...         print(f"Status: {health.status}")
        ...         
        ...         # Get proxy stats
        ...         stats = await client.proxy.get_stats()
        ...         print(f"Requests: {stats.requests_total}")
        >>> 
        >>> asyncio.run(main())
    """
    
    def __init__(
        self,
        base_url: Optional[str] = None,
        api_key: Optional[str] = None,
        token: Optional[str] = None,
        timeout: float = 30.0,
        config: Optional[ClientConfig] = None,
        **kwargs,
    ):
        """
        Initialize the async AegisGate client.
        
        Args:
            base_url: Base URL for the API (default: http://localhost:8080)
            api_key: API key for authentication
            token: Bearer token for session authentication
            timeout: Request timeout in seconds (default: 30)
            config: ClientConfig object for full configuration
            **kwargs: Additional configuration options
        """
        # Initialize configuration
        if config is not None:
            self._config = config
        else:
            self._config = ClientConfig(
                base_url=base_url or "http://localhost:8080",
                api_key=api_key,
                token=token,
                timeout=timeout,
                **kwargs,
            )
        
        # Validate configuration
        errors = self._config.validate()
        if errors:
            raise ConfigurationError(
                message=f"Configuration errors: {', '.join(errors)}",
                code="CONFIG_ERROR",
            )
        
        # Initialize HTTP client
        self._client: Optional[httpx.AsyncClient] = None
        
        # Token storage
        self._token: Optional[str] = None
        
        # Initialize services
        self._init_services()
        
        logger.debug(f"AsyncAegisGate client initialized: {self._config.base_url}")
    
    def _init_services(self) -> None:
        """Initialize all service clients."""
        self.auth = AuthService(self)
        self.proxy = ProxyService(self)
        self.compliance = ComplianceService(self)
        self.siem = SIEMService(self)
        self.webhooks = WebhookService(self)
        self.core = CoreService(self)
    
    async def __aenter__(self) -> "AsyncAegisGateClient":
        """Context manager entry."""
        await self.open()
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb) -> None:
        """Context manager exit."""
        await self.close()
    
    async def open(self) -> None:
        """Open the HTTP client connection."""
        if self._client is None:
            # Configure TLS
            verify = True
            cert = None
            if self._config.tls:
                verify = self._config.tls.verify
                if self._config.tls.cert_file:
                    cert = self._config.tls.cert_file
            
            # Create async client with retry configuration
            retry_config = self._config.retry
            self._client = httpx.AsyncClient(
                timeout=httpx.Timeout(self._config.timeout),
                limits=httpx.Limits(
                    max_connections=self._config.max_connections,
                    max_keepalive_connections=self._config.max_keepalive_connections,
                    keepalive_expiry=self._config.keepalive_expiry,
                ),
                verify=verify,
                cert=cert,
            )
            
            # Set default headers
            self._client.headers["User-Agent"] = f"aegisgate-python-sdk/{__version__}"
            self._client.headers["Accept"] = "application/json"
    
    async def close(self) -> None:
        """Close the HTTP client connection."""
        if self._client:
            await self._client.aclose()
            self._client = None
        logger.debug("AsyncAegisGate client closed")
    
    def _get_auth_headers(self) -> Dict[str, str]:
        """Get authentication headers."""
        headers = {}
        
        if self._config.auth and self._config.auth.api_key:
            headers["X-API-Key"] = self._config.auth.api_key
        elif self._token:
            headers["Authorization"] = f"Bearer {self._token}"
        elif self._config.auth and self._config.auth.token:
            headers["Authorization"] = f"Bearer {self._config.auth.token}"
        
        return headers
    
    async def _request(
        self,
        method: str,
        endpoint: str,
        params: Optional[Dict[str, Any]] = None,
        json: Optional[Dict[str, Any]] = None,
        data: Optional[Any] = None,
        raw: bool = False,
    ) -> Union[Dict[str, Any], bytes]:
        """
        Make an async HTTP request to the AegisGate API.
        
        Args:
            method: HTTP method (GET, POST, PUT, DELETE, etc.)
            endpoint: API endpoint path
            params: URL query parameters
            json: JSON request body
            data: Raw request data
            raw: If True, return raw response bytes
        
        Returns:
            Response data (dict or bytes if raw=True)
        """
        if self._client is None:
            await self.open()
        
        # Build URL
        base_url = self._config.base_url.rstrip("/")
        url = f"{base_url}/api/v1/{endpoint.lstrip('/')}"
        
        # Add query parameters
        if params:
            query_string = urlencode(params, doseq=True)
            url = f"{url}?{query_string}"
        
        # Get headers
        headers = self._get_auth_headers()
        if json is not None:
            headers["Content-Type"] = "application/json"
        
        try:
            response = await self._client.request(
                method=method.upper(),
                url=url,
                headers=headers,
                json=json,
                content=data,
            )
        except httpx.SSLError as e:
            raise SSLError(
                message=f"SSL certificate verification failed: {e}",
                code="SSL_ERROR",
            )
        except httpx.ConnectError as e:
            raise ConnectionError(
                message=f"Failed to connect to {self._config.base_url}: {e}",
                code="CONNECTION_ERROR",
                host=self._config.base_url,
            )
        except httpx.TimeoutException as e:
            raise TimeoutError(
                message=f"Request timed out after {self._config.timeout}s: {e}",
                timeout=self._config.timeout,
            )
        
        return await self._handle_response(response, raw)
    
    async def _handle_response(
        self,
        response: httpx.Response,
        raw: bool = False,
    ) -> Union[Dict[str, Any], bytes]:
        """
        Handle HTTP response and raise appropriate exceptions.
        
        Args:
            response: httpx.Response object
            raw: If True, return raw response bytes
        
        Returns:
            Response data
        """
        if raw:
            return response.content
        
        try:
            data = response.json()
        except json.JSONDecodeError:
            if response.status_code == 200:
                return {}
            raise ServerError(
                message=f"Invalid JSON response: {response.text[:200]}",
                status_code=response.status_code,
            )
        
        if isinstance(data, dict):
            if data.get("error"):
                error_data = data["error"]
                error_code = error_data.get("code", "UNKNOWN")
                error_message = error_data.get("message", "Unknown error")
                
                if response.status_code == 401:
                    raise AuthenticationError(
                        message=error_message,
                        code=error_code,
                    )
                elif response.status_code == 429:
                    raise RateLimitError(
                        message=error_message,
                        code=error_code,
                        retry_after=error_data.get("retry_after"),
                    )
                elif response.status_code == 400:
                    raise ValidationError(
                        message=error_message,
                        code=error_code,
                        field_errors=error_data.get("field_errors", {}),
                    )
                else:
                    raise APIError(
                        message=error_message,
                        error_code=error_code,
                        error_message=error_message,
                        request_id=response.headers.get("X-Request-ID"),
                    )
        
        if response.status_code == 401:
            raise AuthenticationError(
                message="Authentication failed",
                code="UNAUTHORIZED",
            )
        elif response.status_code == 403:
            raise AuthenticationError(
                message="Access forbidden",
                code="FORBIDDEN",
            )
        elif response.status_code == 404:
            raise ResourceNotFoundError(
                message="Resource not found",
                code="NOT_FOUND",
            )
        elif response.status_code == 429:
            retry_after = int(response.headers.get("Retry-After", 60))
            raise RateLimitError(
                message="Rate limit exceeded",
                code="RATE_LIMITED",
                retry_after=retry_after,
            )
        elif 400 <= response.status_code < 500:
            raise APIError(
                message=f"Client error: {response.status_code}",
                status_code=response.status_code,
            )
        elif response.status_code >= 500:
            raise ServerError(
                message=f"Server error: {response.status_code}",
                status_code=response.status_code,
                request_id=response.headers.get("X-Request-ID"),
            )
        
        return data
    
    # Token management
    async def _set_token(self, token: str) -> None:
        """Set the authentication token."""
        self._token = token
    
    async def _clear_token(self) -> None:
        """Clear the authentication token."""
        self._token = None
    
    @property
    def token(self) -> Optional[str]:
        """Get the current authentication token."""
        return self._token
    
    @classmethod
    async def from_env(cls) -> "AsyncAegisGateClient":
        """
        Create a client from environment variables.
        
        Returns:
            AsyncAegisGateClient instance
        """
        config = ClientConfig.from_env()
        return cls(config=config)


# Alias
AegisGateAsync = AsyncAegisGateClient
