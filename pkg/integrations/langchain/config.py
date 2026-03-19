# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Configuration classes for the AegisGate Python SDK.

This module provides configuration dataclasses that allow flexible setup
of the AegisGate client with various authentication methods, TLS settings,
and connection options.

Example:
    Basic configuration:
    >>> from aegisgate import AegisGateClient, ClientConfig
    >>> config = ClientConfig(
    ...     base_url="http://localhost:8080",
    ...     api_key="AG-your-api-key"
    ... )
    >>> client = AegisGateClient(config=config)

    With TLS:
    >>> from aegisgate import ClientConfig, TLSConfig
    >>> config = ClientConfig(
    ...     base_url="https://aegisgate.example.com",
    ...     tls=TLSConfig(
    ...         verify=True,
    ...         cert_file="/path/to/client.crt",
    ...         key_file="/path/to/client.key",
    ...         ca_file="/path/to/ca.crt"
    ...     )
    ... )

    With gRPC:
    >>> from aegisgate import ClientConfig, GRPCConfig
    >>> config = ClientConfig(
    ...     base_url="http://localhost:8080",
    ...     grpc=GRPCConfig(
    ...         address="localhost:50051",
    ...         use_tls=True
    ...     )
    ... )

    Environment variables:
    >>> # Reads from environment:
    >>> # AEGISGATE_BASE_URL, AEGISGATE_API_KEY, AEGISGATE_GRPC_ADDR, etc.
    >>> client = AegisGateClient.from_env()
"""

from __future__ import annotations

import os
from dataclasses import dataclass, field
from typing import Optional, List, Dict, Any
from urllib.parse import urlparse


@dataclass
class TLSConfig:
    """
    TLS/SSL configuration for secure connections.
    
    Attributes:
        verify: Whether to verify TLS certificates (default: True)
        cert_file: Path to client certificate file (PEM)
        key_file: Path to client private key file (PEM)
        key_password: Password for encrypted private key
        ca_file: Path to CA certificate file for server verification
        server_hostname: Server hostname for SNI (Server Name Indication)
    """
    verify: bool = True
    cert_file: Optional[str] = None
    key_file: Optional[str] = None
    key_password: Optional[str] = None
    ca_file: Optional[str] = None
    server_hostname: Optional[str] = None
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary for HTTP client configuration."""
        result = {"verify": self.verify}
        if self.cert_file:
            result["cert"] = self.cert_file
        if self.key_file:
            if "cert" in result:
                result["cert"] = (self.cert_file, self.key_file, self.key_password)
            else:
                result["cert"] = self.key_file
        if self.ca_file:
            result["ca_cert"] = self.ca_file
        if self.server_hostname:
            result["server_hostname"] = self.server_hostname
        return result


@dataclass
class GRPCConfig:
    """
    gRPC connection configuration.
    
    AegisGate supports gRPC for high-performance streaming and
    low-latency communication.
    
    Attributes:
        address: gRPC server address (host:port)
        use_tls: Whether to use TLS (default: False for localhost)
        timeout: Connection timeout in seconds (default: 10)
        max_receive_message_length: Maximum receive message size in bytes
        max_send_message_length: Maximum send message size in bytes
    """
    address: str = "localhost:50051"
    use_tls: bool = False
    timeout: float = 10.0
    max_receive_message_length: int = 100 * 1024 * 1024  # 100MB
    max_send_message_length: int = 100 * 1024 * 1024  # 100MB
    
    def validate(self) -> List[str]:
        """Validate configuration and return list of errors."""
        errors = []
        if not self.address:
            errors.append("gRPC address cannot be empty")
        if self.timeout <= 0:
            errors.append("gRPC timeout must be positive")
        if ":" not in self.address:
            errors.append("gRPC address must be in host:port format")
        return errors


@dataclass
class AuthConfig:
    """
    Authentication configuration.
    
    Supports multiple authentication methods:
    - API Key (recommended for most use cases)
    - Bearer Token (for session-based auth)
    - Basic Auth (username/password)
    
    Attributes:
        api_key: API key for authentication (recommended)
        token: Bearer token for session authentication
        username: Username for basic authentication
        password: Password for basic authentication
    """
    api_key: Optional[str] = None
    token: Optional[str] = None
    username: Optional[str] = None
    password: Optional[str] = None
    
    def method(self) -> str:
        """Return the authentication method in use."""
        if self.api_key:
            return "api_key"
        elif self.token:
            return "bearer_token"
        elif self.username and self.password:
            return "basic_auth"
        return "none"
    
    def validate(self) -> List[str]:
        """Validate configuration and return list of errors."""
        errors = []
        method = self.method()
        if method == "basic_auth" and not self.password:
            errors.append("Password required for basic authentication")
        return errors


@dataclass
class RESTConfig:
    """
    REST API configuration.
    
    Attributes:
        base_url: Base URL for the REST API (e.g., "http://localhost:8080")
        api_version: API version prefix (default: "v1")
        user_agent: User-Agent string for requests
        headers: Additional HTTP headers to send with all requests
    """
    base_url: str = "http://localhost:8080"
    api_version: str = "v1"
    user_agent: str = "aegisgate-python-sdk/1.0.0"
    headers: Dict[str, str] = field(default_factory=dict)
    
    def validate(self) -> List[str]:
        """Validate configuration and return list of errors."""
        errors = []
        if not self.base_url:
            errors.append("base_url cannot be empty")
        else:
            try:
                urlparse(self.base_url)
            except Exception:
                errors.append(f"Invalid base_url: {self.base_url}")
        if not self.api_version:
            errors.append("api_version cannot be empty")
        return errors
    
    def api_url(self) -> str:
        """Return the full API URL path."""
        return f"{self.base_url}/api/{self.api_version}"


@dataclass
class ProxyConfig:
    """
    HTTP proxy configuration for outgoing connections.
    
    Use this if AegisGate must be accessed through a proxy.
    
    Attributes:
        http_url: HTTP proxy URL (e.g., "http://proxy.example.com:8080")
        https_url: HTTPS proxy URL (if different from HTTP)
        no_proxy: Comma-separated list of hosts to skip proxying
        auth_username: Proxy authentication username
        auth_password: Proxy authentication password
    """
    http_url: Optional[str] = None
    https_url: Optional[str] = None
    no_proxy: Optional[str] = None
    auth_username: Optional[str] = None
    auth_password: Optional[str] = None
    
    def is_configured(self) -> bool:
        """Return True if proxy is configured."""
        return bool(self.http_url or self.https_url)


@dataclass
class RetryConfig:
    """
    Retry configuration for failed requests.
    
    Attributes:
        max_attempts: Maximum number of retry attempts (default: 3)
        base_delay: Initial delay between retries in seconds (default: 1.0)
        max_delay: Maximum delay between retries in seconds (default: 60.0)
        backoff_factor: Exponential backoff multiplier (default: 2.0)
        retry_on_status: HTTP status codes that trigger retry
        retry_on_timeout: Whether to retry on timeout errors
    """
    max_attempts: int = 3
    base_delay: float = 1.0
    max_delay: float = 60.0
    backoff_factor: float = 2.0
    retry_on_status: List[int] = field(
        default_factory=lambda: [429, 500, 502, 503, 504]
    )
    retry_on_timeout: bool = True
    
    def delay(self, attempt: int) -> float:
        """Calculate delay for given attempt number."""
        import random
        delay = min(self.base_delay * (self.backoff_factor ** attempt), self.max_delay)
        # Add jitter (±10%)
        jitter = delay * 0.1 * (2 * random.random() - 1)
        return delay + jitter


@dataclass
class RateLimitConfig:
    """
    Rate limiting configuration.
    
    These settings control how the SDK manages request rates
    to avoid exceeding server-side limits.
    
    Attributes:
        requests_per_second: Maximum requests per second
        burst_size: Maximum burst size for token bucket
        per_host: Apply rate limiting per-host (default: True)
    """
    requests_per_second: float = 10.0
    burst_size: int = 20
    per_host: bool = True


@dataclass
class ClientConfig:
    """
    Main configuration class for the AegisGate client.
    
    This is the primary configuration object that combines all
    configuration options. It can be created directly or loaded
    from environment variables.
    
    Attributes:
        base_url: Base URL for the AegisGate REST API
        timeout: Request timeout in seconds (default: 30)
        max_connections: Maximum number of concurrent connections
        max_keepalive_connections: Maximum keepalive connections
        keepalive_expiry: Keepalive connection expiry in seconds
        tls: TLS configuration
        grpc: gRPC configuration
        auth: Authentication configuration
        proxy: HTTP proxy configuration
        retry: Retry configuration
        rate_limit: Rate limiting configuration
        debug: Enable debug mode (verbose logging)
        
    Example:
        >>> config = ClientConfig(
        ...     base_url="https://aegisgate.example.com",
        ...     api_key="AG-your-api-key",
        ...     timeout=60.0,
        ...     tls=TLSConfig(verify=True),
        ...     retry=RetryConfig(max_attempts=5)
        ... )
        >>> client = AegisGateClient(config=config)
    """
    # Connection settings
    base_url: str = "http://localhost:8080"
    timeout: float = 30.0
    max_connections: int = 100
    max_keepalive_connections: int = 20
    keepalive_expiry: float = 30.0
    
    # Sub-configurations
    tls: Optional[TLSConfig] = None
    grpc: Optional[GRPCConfig] = None
    auth: Optional[AuthConfig] = None
    proxy: Optional[ProxyConfig] = None
    retry: Optional[RetryConfig] = None
    rate_limit: Optional[RateLimitConfig] = None
    
    # Feature flags
    debug: bool = False
    
    # For convenience - these are extracted to auth
    api_key: Optional[str] = field(default=None, repr=False)
    token: Optional[str] = field(default=None, repr=False)
    
    def __post_init__(self):
        """Process configuration after initialization."""
        # Handle convenience auth parameters
        if self.api_key or self.token:
            if self.auth is None:
                self.auth = AuthConfig()
            if self.api_key:
                self.auth.api_key = self.api_key
            if self.token:
                self.auth.token = self.token
        
        # Set defaults for sub-configurations
        if self.tls is None:
            self.tls = TLSConfig()
        if self.grpc is None:
            self.grpc = GRPCConfig()
        if self.retry is None:
            self.retry = RetryConfig()
        if self.rate_limit is None:
            self.rate_limit = RateLimitConfig()
    
    def validate(self) -> List[str]:
        """Validate all configuration and return list of errors."""
        errors = []
        
        # Validate REST config
        rest_config = RESTConfig(base_url=self.base_url)
        errors.extend(rest_config.validate())
        
        # Validate sub-configurations
        if self.grpc:
            errors.extend(self.grpc.validate())
        if self.auth:
            errors.extend(self.auth.validate())
        
        # Validate timeout
        if self.timeout <= 0:
            errors.append("timeout must be positive")
        
        return errors
    
    def is_valid(self) -> bool:
        """Return True if configuration is valid."""
        return len(self.validate()) == 0
    
    @classmethod
    def from_env(cls) -> ClientConfig:
        """
        Create configuration from environment variables.
        
        Environment variables:
            AEGISGATE_BASE_URL: Base URL for the API
            AEGISGATE_API_KEY: API key for authentication
            AEGISGATE_TOKEN: Bearer token for authentication
            AEGISGATE_GRPC_ADDR: gRPC server address
            AEGISGATE_TIMEOUT: Request timeout in seconds
            AEGISGATE_TLS_VERIFY: Whether to verify TLS (true/false)
            AEGISGATE_DEBUG: Enable debug mode (true/false)
        
        Returns:
            ClientConfig with values from environment
        
        Example:
            >>> import os
            >>> os.environ["AEGISGATE_BASE_URL"] = "https://aegisgate.example.com"
            >>> os.environ["AEGISGATE_API_KEY"] = "AG-your-api-key"
            >>> config = ClientConfig.from_env()
        """
        config = cls()
        
        if base_url := os.environ.get("AEGISGATE_BASE_URL"):
            config.base_url = base_url
        
        if api_key := os.environ.get("AEGISGATE_API_KEY"):
            config.api_key = api_key
        
        if token := os.environ.get("AEGISGATE_TOKEN"):
            config.token = token
        
        if grpc_addr := os.environ.get("AEGISGATE_GRPC_ADDR"):
            if config.grpc is None:
                config.grpc = GRPCConfig()
            config.grpc.address = grpc_addr
        
        if timeout := os.environ.get("AEGISGATE_TIMEOUT"):
            try:
                config.timeout = float(timeout)
            except ValueError:
                pass
        
        if tls_verify := os.environ.get("AEGISGATE_TLS_VERIFY"):
            if config.tls is None:
                config.tls = TLSConfig()
            config.tls.verify = tls_verify.lower() in ("true", "1", "yes")
        
        if debug := os.environ.get("AEGISGATE_DEBUG"):
            config.debug = debug.lower() in ("true", "1", "yes")
        
        return config
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert configuration to dictionary."""
        result = {
            "base_url": self.base_url,
            "timeout": self.timeout,
            "max_connections": self.max_connections,
            "max_keepalive_connections": self.max_keepalive_connections,
            "keepalive_expiry": self.keepalive_expiry,
            "debug": self.debug,
        }
        if self.tls:
            result["tls"] = self.tls.to_dict()
        if self.grpc:
            result["grpc"] = {
                "address": self.grpc.address,
                "use_tls": self.grpc.use_tls,
                "timeout": self.grpc.timeout,
            }
        if self.auth:
            result["auth"] = {
                "method": self.auth.method(),
            }
        if self.proxy and self.proxy.is_configured():
            result["proxy"] = {
                "http_url": self.proxy.http_url,
                "https_url": self.proxy.https_url,
            }
        return result
