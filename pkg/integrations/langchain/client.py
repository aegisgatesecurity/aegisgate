# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
AegisGate Python SDK - Main Client

This module provides the main AegisGateClient class for synchronous
communication with the AegisGate API.
"""

from __future__ import annotations

import json
import logging
import time
from typing import Optional, Dict, Any, Union
from urllib.parse import urlencode

import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

from aegisgate._version import __version__
from aegisgate.config import ClientConfig, TLSConfig
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


class AegisGateClient:
    """
    Main client for interacting with the AegisGate API.
    
    This client provides access to all AegisGate services including:
    - Authentication and user management (auth)
    - Proxy statistics and violations (proxy)
    - Compliance framework management (compliance)
    - SIEM integration (siem)
    - Webhook management (webhooks)
    - Core system operations (core)
    
    Example:
        Basic usage:
        >>> from aegisgate import AegisGateClient
        >>> client = AegisGateClient(
        ...     base_url="http://localhost:8080",
        ...     api_key="AG-your-api-key"
        ... )
        >>> 
        >>> # Check health
        >>> health = client.core.health()
        >>> print(f"Status: {health.status}")
        >>> 
        >>> # Get proxy stats
        >>> stats = client.proxy.get_stats()
        >>> print(f"Requests: {stats.requests_total}")
        >>> 
        >>> # Close client
        >>> client.close()
        
        With environment variables:
        >>> # Set AEGISGATE_BASE_URL and AEGISGATE_API_KEY
        >>> client = AegisGateClient.from_env()
        
        With custom configuration:
        >>> from aegisgate import ClientConfig, TLSConfig
        >>> config = ClientConfig(
        ...     base_url="https://aegisgate.example.com",
        ...     api_key="AG-your-api-key",
        ...     timeout=60.0,
        ...     tls=TLSConfig(verify=True, ca_file="/path/to/ca.pem"),
        ...     retry=RetryConfig(max_attempts=5)
        ... )
        >>> client = AegisGateClient(config=config)
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
        Initialize the AegisGate client.
        
        Args:
            base_url: Base URL for the API (default: http://localhost:8080)
            api_key: API key for authentication
            token: Bearer token for session authentication
            timeout: Request timeout in seconds (default: 30)
            config: ClientConfig object for full configuration
            **kwargs: Additional configuration options
            
        Raises:
            ConfigurationError: If configuration is invalid
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
        
        # Initialize HTTP session
        self._session = self._create_session()
        
        # Token storage
        self._token: Optional[str] = None
        
        # Initialize services
        self._init_services()
        
        logger.debug(f"AegisGate client initialized: {self._config.base_url}")
    
    def _create_session(self) -> requests.Session:
        """
        Create configured HTTP session.
        
        Returns:
            Configured requests.Session
        """
        session = requests.Session()
        
        # Configure retry strategy
        retry_config = self._config.retry
        retry_strategy = Retry(
            total=retry_config.max_attempts,
            backoff_factor=retry_config.backoff_factor,
            status_forcelist=retry_config.retry_on_status,
            allowed_methods=["GET", "POST", "PUT", "DELETE", "PATCH"],
            raise_on_status=False,
        )
        
        # Configure adapter
        adapter = HTTPAdapter(
            max_retries=retry_strategy,
            pool_connections=self._config.max_keepalive_connections,
            pool_maxsize=self._config.max_connections,
        )
        
        session.mount("http://", adapter)
        session.mount("https://", adapter)
        
        # Set default headers
        session.headers.update({
            "User-Agent": f"aegisgate-python-sdk/{__version__}",
            "Accept": "application/json",
        })
        
        # Configure TLS if specified
        tls_config = self._config.tls
        if tls_config and not tls_config.verify:
            import urllib3
            urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
        
        return session
    
    def _init_services(self) -> None:
        """Initialize all service clients."""
        self.auth = AuthService(self)
        self.proxy = ProxyService(self)
        self.compliance = ComplianceService(self)
        self.siem = SIEMService(self)
        self.webhooks = WebhookService(self)
        self.core = CoreService(self)
    
    def _get_auth_headers(self) -> Dict[str, str]:
        """Get authentication headers."""
        headers = {}
        
        # API key takes precedence
        if self._config.auth and self._config.auth.api_key:
            headers["X-API-Key"] = self._config.auth.api_key
        # Then bearer token
        elif self._token:
            headers["Authorization"] = f"Bearer {self._token}"
        elif self._config.auth and self._config.auth.token:
            headers["Authorization"] = f"Bearer {self._config.auth.token}"
        
        return headers
    
    def _request(
        self,
        method: str,
        endpoint: str,
        params: Optional[Dict[str, Any]] = None,
        json: Optional[Dict[str, Any]] = None,
        data: Optional[Any] = None,
        raw: bool = False,
    ) -> Union[Dict[str, Any], bytes]:
        """
        Make an HTTP request to the AegisGate API.
        
        Args:
            method: HTTP method (GET, POST, PUT, DELETE, etc.)
            endpoint: API endpoint path
            params: URL query parameters
            json: JSON request body
            data: Raw request data
            raw: If True, return raw response bytes
        
        Returns:
            Response data (dict or bytes if raw=True)
        
        Raises:
            AuthenticationError: If authentication fails
            RateLimitError: If rate limit is exceeded
            ServerError: If server returns 5xx error
            APIError: If API returns an error response
        """
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
            response = self._session.request(
                method=method.upper(),
                url=url,
                headers=headers,
                json=json,
                data=data,
                timeout=self._config.timeout,
                verify=self._config.tls.verify if self._config.tls else True,
                cert=self._config.tls.cert_file if self._config.tls else None,
            )
        except requests.exceptions.SSLError as e:
            raise SSLError(
                message=f"SSL certificate verification failed: {e}",
                code="SSL_ERROR",
            )
        except requests.exceptions.ConnectionError as e:
            raise ConnectionError(
                message=f"Failed to connect to {self._config.base_url}: {e}",
                code="CONNECTION_ERROR",
                host=self._config.base_url,
            )
        except requests.exceptions.Timeout as e:
            raise TimeoutError(
                message=f"Request timed out after {self._config.timeout}s: {e}",
                timeout=self._config.timeout,
            )
        
        # Handle response
        return self._handle_response(response, raw)
    
    def _handle_response(
        self,
        response: requests.Response,
        raw: bool = False,
    ) -> Union[Dict[str, Any], bytes]:
        """
        Handle HTTP response and raise appropriate exceptions.
        
        Args:
            response: requests.Response object
            raw: If True, return raw response bytes
        
        Returns:
            Response data
        
        Raises:
            AuthenticationError: For 401 responses
            RateLimitError: For 429 responses
            ResourceNotFoundError: For 404 responses
            ValidationError: For 400 responses
            ServerError: For 5xx responses
            APIError: For other error responses
        """
        # Return raw bytes if requested
        if raw:
            return response.content
        
        # Parse JSON response
        try:
            data = response.json()
        except json.JSONDecodeError:
            if response.status_code == 200:
                return {}
            raise ServerError(
                message=f"Invalid JSON response: {response.text[:200]}",
                status_code=response.status_code,
            )
        
        # Check for API errors in response body
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
        
        # Handle HTTP status codes
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
    
    # =====================================================================
    # Token Management
    # =====================================================================
    
    def _set_token(self, token: str) -> None:
        """Set the authentication token."""
        self._token = token
        logger.debug("Authentication token set")
    
    def _clear_token(self) -> None:
        """Clear the authentication token."""
        self._token = None
        logger.debug("Authentication token cleared")
    
    @property
    def token(self) -> Optional[str]:
        """Get the current authentication token."""
        return self._token
    
    # =====================================================================
    # Context Manager
    # =====================================================================
    
    def __enter__(self) -> "AegisGateClient":
        """Context manager entry."""
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        """Context manager exit."""
        self.close()
    
    def close(self) -> None:
        """Close the client and release resources."""
        if self._session:
            self._session.close()
        logger.debug("AegisGate client closed")
    
    # =====================================================================
    # Class Methods
    # =====================================================================
    
    @classmethod
    def from_env(cls) -> "AegisGateClient":
        """
        Create a client from environment variables.
        
        Environment variables:
            AEGISGATE_BASE_URL: Base URL for the API
            AEGISGATE_API_KEY: API key for authentication
            AEGISGATE_TOKEN: Bearer token for authentication
            AEGISGATE_TIMEOUT: Request timeout in seconds
            AEGISGATE_TLS_VERIFY: Whether to verify TLS (true/false)
            AEGISGATE_DEBUG: Enable debug mode (true/false)
        
        Returns:
            AegisGateClient instance
        
        Example:
            >>> import os
            >>> os.environ["AEGISGATE_BASE_URL"] = "http://localhost:8080"
            >>> os.environ["AEGISGATE_API_KEY"] = "AG-your-key"
            >>> client = AegisGateClient.from_env()
        """
        config = ClientConfig.from_env()
        return cls(config=config)
    
    @classmethod
    def from_url(cls, url: str, api_key: Optional[str] = None) -> "AegisGateClient":
        """
        Create a client from a single URL with optional API key.
        
        Args:
            url: Full URL (will extract base URL and path)
            api_key: Optional API key
        
        Returns:
            AegisGateClient instance
        
        Example:
            >>> client = AegisGateClient.from_url(
            ...     "https://aegisgate.example.com/api/v1",
            ...     "AG-your-key"
            ... )
        """
        from urllib.parse import urlparse
        
        parsed = urlparse(url)
        base_url = f"{parsed.scheme}://{parsed.netloc}"
        
        return cls(base_url=base_url, api_key=api_key)


# Expose client at module level
AegisGate = AegisGateClient
