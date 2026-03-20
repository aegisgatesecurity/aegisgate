"""
Tests for synchronous and asynchronous clients.
"""

import os
import pytest
from unittest.mock import MagicMock, patch, AsyncMock

from aegisgate import Client, AsyncClient
from aegisgate.connection import ConnectionConfig, SyncConnection, AsyncConnection, APIError, ConnectionError


class TestClient:
    """Tests for the synchronous Client."""

    def test_client_initialization_with_api_key(self):
        """Test client initialization with explicit API key."""
        client = Client(base_url="https://test.api.io", api_key="test-key")
        assert client._config.api_key == "test-key"
        client.close()

    def test_client_initialization_with_env_var(self):
        """Test client initialization using environment variable."""
        os.environ["AEGISGATE_API_KEY"] = "env-test-key"
        os.environ["AEGISGATE_BASE_URL"] = "https://env.api.io"
        client = Client()
        assert client._config.api_key == "env-test-key"
        assert client._config.base_url == "https://env.api.io"
        client.close()
        # Clean up
        del os.environ["AEGISGATE_API_KEY"]
        del os.environ["AEGISGATE_BASE_URL"]

    def test_client_initialization_with_base_url(self):
        """Test client initialization with custom base URL."""
        client = Client(
            api_key="test-key",
            base_url="https://custom.api.io"
        )
        assert client._config.base_url == "https://custom.api.io"
        client.close()

    def test_client_initialization_with_timeout(self):
        """Test client initialization with custom timeout."""
        # Clear environment variables to avoid override
        old_key = os.environ.pop("AEGISGATE_API_KEY", None)
        old_url = os.environ.pop("AEGISGATE_BASE_URL", None)
        try:
            client = Client(
                base_url="https://config.api.io",
                api_key="config-key",
                timeout=60.0,
                max_retries=5
            )
            assert client._config.api_key == "config-key"
            assert client._config.timeout == 60.0
            assert client._config.max_retries == 5
            client.close()
        finally:
            if old_key:
                os.environ["AEGISGATE_API_KEY"] = old_key
            if old_url:
                os.environ["AEGISGATE_BASE_URL"] = old_url

    def test_client_context_manager(self):
        """Test client context manager pattern."""
        with Client(base_url="https://test.api.io", api_key="test-key") as client:
            assert client is not None
            assert isinstance(client, Client)

    def test_client_services_property_access(self):
        """Test accessing all service properties."""
        client = Client(base_url="https://test.api.io", api_key="test-key")
        
        # Test all service properties exist
        assert client.auth is not None
        assert client.proxy is not None
        assert client.compliance is not None
        assert client.siem is not None
        assert client.webhook is not None
        assert client.core is not None
        client.close()


class TestAsyncClient:
    """Tests for the asynchronous AsyncClient."""

    def test_async_client_initialization_with_api_key(self):
        """Test async client initialization with explicit API key."""
        client = AsyncClient(base_url="https://test.api.io", api_key="test-key")
        assert client._config.api_key == "test-key"

    def test_async_client_initialization_with_env_var(self):
        """Test async client initialization using environment variable."""
        os.environ["AEGISGATE_API_KEY"] = "env-async-key"
        os.environ["AEGISGATE_BASE_URL"] = "https://env.api.io"
        client = AsyncClient()
        assert client._config.api_key == "env-async-key"
        # Clean up
        del os.environ["AEGISGATE_API_KEY"]
        del os.environ["AEGISGATE_BASE_URL"]

    def test_async_client_initialization_with_base_url(self):
        """Test async client initialization with custom base URL."""
        client = AsyncClient(
            api_key="test-key",
            base_url="https://custom.api.io"
        )
        assert client._config.base_url == "https://custom.api.io"

    @pytest.mark.asyncio
    async def test_async_client_context_manager(self):
        """Test async client context manager pattern."""
        async with AsyncClient(base_url="https://test.api.io", api_key="test-key") as client:
            assert client is not None
            assert isinstance(client, AsyncClient)

    def test_async_client_services_property_access(self):
        """Test accessing all service properties on async client."""
        client = AsyncClient(base_url="https://test.api.io", api_key="test-key")
        
        # Test all service properties exist
        assert client.auth is not None
        assert client.proxy is not None
        assert client.compliance is not None
        assert client.siem is not None
        assert client.webhook is not None
        assert client.core is not None


class TestConnectionConfig:
    """Tests for ConnectionConfig."""

    def test_default_config(self):
        """Test default connection configuration."""
        config = ConnectionConfig(
            base_url="https://api.aegisgate.io",
            api_key="test-key"
        )
        
        assert config.api_key == "test-key"
        assert config.base_url == "https://api.aegisgate.io"
        assert config.timeout == 30.0
        assert config.max_retries == 3
        assert config.verify_ssl is True

    def test_custom_config(self):
        """Test custom connection configuration."""
        config = ConnectionConfig(
            api_key="test-key",
            base_url="https://custom.api.io",
            timeout=60.0,
            max_retries=5,
            verify_ssl=False,
            proxy="http://proxy.io:8080",
            custom_headers={"X-Custom": "header"}
        )
        
        assert config.base_url == "https://custom.api.io"
        assert config.timeout == 60.0
        assert config.max_retries == 5
        assert config.verify_ssl is False
        assert config.proxy == "http://proxy.io:8080"
        assert config.custom_headers == {"X-Custom": "header"}

    def test_config_strips_trailing_slash(self):
        """Test that base_url trailing slash is stripped."""
        config = ConnectionConfig(
            base_url="https://api.aegisgate.io/",
            api_key="test-key"
        )
        assert config.base_url == "https://api.aegisgate.io"


class TestAPIError:
    """Tests for APIError exception."""

    def test_api_error_creation(self):
        """Test API error creation."""
        error = APIError(
            message="Not Found",
            status_code=404
        )
        
        assert error.status_code == 404
        assert str(error) == "Not Found"
        assert error.details is None

    def test_api_error_with_details(self):
        """Test API error with additional details."""
        error = APIError(
            message="Bad Request",
            status_code=400,
            details={"field": "email", "error": "invalid"}
        )
        
        assert error.status_code == 400
        assert error.details["field"] == "email"
        assert error.details["error"] == "invalid"


class TestConnectionError:
    """Tests for ConnectionError exception."""

    def test_connection_error_creation(self):
        """Test connection error creation."""
        error = ConnectionError("Failed to connect")
        
        assert str(error) == "Failed to connect"