# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Custom exceptions for the AegisGate Python SDK.

This module defines a hierarchical exception system that allows for granular
error handling when working with the AegisGate API.

Example:
    >>> try:
    ...     client = AegisGateClient(api_key="invalid")
    ...     client.auth.login("user", "pass")
    ... except AuthenticationError as e:
    ...     print(f"Authentication failed: {e}")
    ... except RateLimitError as e:
    ...     print(f"Rate limited: {e.retry_after}s until retry")
"""

from typing import Optional, Any, Dict


class AegisGateError(Exception):
    """
    Base exception for all AegisGate SDK errors.
    
    All other exceptions inherit from this class, allowing you to catch
    all AegisGate-related errors with a single except clause.
    
    Attributes:
        message: Human-readable error message
        code: Optional error code for programmatic handling
        details: Additional error details from the API
        response: Raw API response if available
    """
    
    def __init__(
        self,
        message: str,
        code: Optional[str] = None,
        details: Optional[Dict[str, Any]] = None,
        response: Optional[Dict[str, Any]] = None,
    ):
        super().__init__(message)
        self.message = message
        self.code = code
        self.details = details or {}
        self.response = response
    
    def __repr__(self) -> str:
        parts = [f"AegisGateError({self.message!r}"]
        if self.code:
            parts.append(f", code={self.code!r}")
        parts.append(")")
        return "".join(parts)
    
    def __str__(self) -> str:
        if self.code:
            return f"[{self.code}] {self.message}"
        return self.message


class ConfigurationError(AegisGateError):
    """
    Raised when there is a configuration error.
    
    This exception is raised when:
    - Invalid configuration values are provided
    - Required configuration is missing
    - Configuration values are out of valid range
    
    Example:
        >>> config = ClientConfig(base_url="")  # Empty URL
        Traceback (most recent call last):
            ...
        ConfigurationError: base_url cannot be empty
    """
    pass


class AuthenticationError(AegisGateError):
    """
    Raised when authentication fails.
    
    This exception is raised when:
    - Invalid API key is provided
    - Username/password is incorrect
    - OAuth token is expired or invalid
    - Session has expired
    
    Attributes:
        retry_after: Optional seconds to wait before retry (for token refresh)
    """
    
    def __init__(
        self,
        message: str,
        code: Optional[str] = None,
        retry_after: Optional[int] = None,
        **kwargs,
    ):
        super().__init__(message, code, **kwargs)
        self.retry_after = retry_after


class AuthorizationError(AegisGateError):
    """
    Raised when the authenticated user lacks required permissions.
    
    This exception is raised when:
    - User role doesn't permit the operation
    - License tier doesn't support the feature
    - Resource ownership doesn't allow access
    
    Example:
        >>> try:
        ...     client.admin.delete_user("user123")
        ... except AuthorizationError as e:
        ...     print(f"Insufficient permissions: {e}")
    """
    pass


class ValidationError(AegisGateError):
    """
    Raised when request validation fails.
    
    This exception is raised when:
    - Request body fails validation
    - Required fields are missing
    - Field values are invalid type or range
    
    Attributes:
        field_errors: Dictionary of field names to error messages
    """
    
    def __init__(
        self,
        message: str,
        field_errors: Optional[Dict[str, str]] = None,
        **kwargs,
    ):
        super().__init__(message, **kwargs)
        self.field_errors = field_errors or {}


class ResourceNotFoundError(AegisGateError):
    """
    Raised when a requested resource doesn't exist.
    
    This exception is raised when:
    - The specified resource ID doesn't exist
    - A required configuration is missing
    - The endpoint path is incorrect
    
    Example:
        >>> try:
        ...     client.users.get("nonexistent-id")
        ... except ResourceNotFoundError as e:
        ...     print(f"User not found: {e.resource_id}")
    """
    
    def __init__(
        self,
        message: str,
        resource_type: Optional[str] = None,
        resource_id: Optional[str] = None,
        **kwargs,
    ):
        super().__init__(message, **kwargs)
        self.resource_type = resource_type
        self.resource_id = resource_id


class RateLimitError(AegisGateError):
    """
    Raised when rate limit is exceeded.
    
    This exception is raised when:
    - Too many requests in a time window
    - Burst limit exceeded
    - License tier request limit reached
    
    Attributes:
        retry_after: Seconds to wait before retry
        limit: The rate limit that was exceeded
        remaining: Requests remaining after this error
    """
    
    def __init__(
        self,
        message: str,
        retry_after: Optional[int] = None,
        limit: Optional[int] = None,
        remaining: Optional[int] = None,
        **kwargs,
    ):
        super().__init__(message, **kwargs)
        self.retry_after = retry_after
        self.limit = limit
        self.remaining = remaining


class NetworkError(AegisGateError):
    """
    Raised when a network error occurs.
    
    This exception is raised when:
    - Connection to server fails
    - DNS resolution fails
    - Connection timeout
    - Connection reset by peer
    
    Attributes:
        host: The host that failed to connect
        port: The port that was used
    """
    
    def __init__(
        self,
        message: str,
        host: Optional[str] = None,
        port: Optional[int] = None,
        **kwargs,
    ):
        super().__init__(message, **kwargs)
        self.host = host
        self.port = port


class TimeoutError(AegisGateError):
    """
    Raised when a request times out.
    
    This exception is raised when:
    - Request takes longer than configured timeout
    - Server doesn't respond within time limit
    - gRPC deadline exceeded
    
    Attributes:
        timeout: The timeout value that was exceeded (in seconds)
    """
    
    def __init__(
        self,
        message: str,
        timeout: Optional[float] = None,
        **kwargs,
    ):
        super().__init__(message, **kwargs)
        self.timeout = timeout


class ConnectionError(NetworkError):
    """
    Raised when connection to AegisGate server fails.
    
    This is a more specific variant of NetworkError for connection issues.
    
    Example:
        >>> try:
        ...     client = AegisGateClient(base_url="https://unreachable.example.com")
        ...     client.core.health()
        ... except ConnectionError as e:
        ...     print(f"Cannot connect to AegisGate: {e}")
    """
    pass


class ServerError(AegisGateError):
    """
    Raised when the AegisGate server returns an error.
    
    This exception is raised for:
    - 5xx HTTP status codes
    - Internal server errors
    - Service unavailable
    - Gateway timeout
    
    Attributes:
        status_code: HTTP status code
        request_id: Server-side request ID for debugging
    """
    
    def __init__(
        self,
        message: str,
        status_code: Optional[int] = None,
        request_id: Optional[str] = None,
        **kwargs,
    ):
        super().__init__(message, **kwargs)
        self.status_code = status_code
        self.request_id = request_id


class APIError(AegisGateError):
    """
    Raised when the API returns an error response.
    
    This is the base class for API-level errors that include
    structured error responses from the server.
    
    Attributes:
        error_code: Machine-readable error code from API
        error_message: Human-readable error from API
        request_id: Server-side request ID for debugging
    """
    
    def __init__(
        self,
        message: str,
        error_code: Optional[str] = None,
        error_message: Optional[str] = None,
        request_id: Optional[str] = None,
        **kwargs,
    ):
        super().__init__(message, error_code, **kwargs)
        self.error_code = error_code or self.code
        self.error_message = error_message or message
        self.request_id = request_id


class SSLError(NetworkError):
    """
    Raised when SSL/TLS certificate verification fails.
    
    This exception is raised when:
    - Certificate is expired
    - Certificate hostname mismatch
    - Self-signed certificate not trusted
    - Certificate chain is broken
    
    Example:
        >>> # To handle SSL errors:
        >>> from aegisgate import ClientConfig, TLSConfig
        >>> config = ClientConfig(
        ...     tls=TLSConfig(verify=True, ca_file="/path/to/ca.pem")
        ... )
    """
    pass
