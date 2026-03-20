"""
Pytest configuration and fixtures for AegisGate SDK tests.
"""

import os
import pytest
from unittest.mock import MagicMock, AsyncMock
from typing import Generator

# Set environment variables for testing
os.environ.setdefault("AEGISGATE_API_KEY", "test-api-key")
os.environ.setdefault("AEGISGATE_BASE_URL", "https://test.aegisgate.io")


@pytest.fixture
def mock_api_key() -> str:
    """Return a test API key."""
    return "test-api-key-12345"


@pytest.fixture
def mock_base_url() -> str:
    """Return a test base URL."""
    return "https://test.aegisgate.io"


@pytest.fixture
def mock_response() -> dict:
    """Return a mock API response."""
    return {
        "id": "test-id",
        "status": "success",
        "data": {"key": "value"}
    }


@pytest.fixture
def mock_health_response() -> dict:
    """Return a mock health check response."""
    return {
        "status": "healthy",
        "version": "1.0.0",
        "uptime": 86400
    }


@pytest.fixture
def mock_violation_response() -> dict:
    """Return a mock violation detection response."""
    return {
        "has_violations": True,
        "violations": [
            {
                "id": "v-001",
                "type": "prompt_injection",
                "severity": "high",
                "message": "Potential prompt injection detected",
                "confidence": 0.95
            }
        ],
        "confidence": 0.95,
        "metadata": {}
    }


@pytest.fixture
def sync_connection_mock() -> MagicMock:
    """Return a mock synchronous connection."""
    mock = MagicMock()
    mock.request.return_value = {"status": "success"}
    mock.get.return_value = {"status": "success"}
    mock.post.return_value = {"status": "success"}
    mock.put.return_value = {"status": "success"}
    mock.delete.return_value = {"status": "success"}
    return mock


@pytest.fixture
def async_connection_mock() -> AsyncMock:
    """Return a mock asynchronous connection."""
    mock = AsyncMock()
    mock.request.return_value = {"status": "success"}
    mock.get.return_value = {"status": "success"}
    mock.post.return_value = {"status": "success"}
    mock.put.return_value = {"status": "success"}
    mock.delete.return_value = {"status": "success"}
    return mock


@pytest.fixture
def mock_client(sync_connection_mock: MagicMock) -> "Client":
    """Return a mock synchronous client."""
    from aegisgate import Client
    from aegisgate.connection import ConnectionConfig
    
    client = Client.__new__(Client)
    client._config = ConnectionConfig(api_key="test-key")
    client._connection = sync_connection_mock
    return client


@pytest.fixture
def mock_async_client(async_connection_mock: AsyncMock) -> "AsyncClient":
    """Return a mock asynchronous client."""
    from aegisgate import AsyncClient
    from aegisgate.connection import ConnectionConfig
    
    client = AsyncClient.__new__(AsyncClient)
    client._config = ConnectionConfig(api_key="test-key")
    client._connection = async_connection_mock
    return client


# Async test markers
@pytest.fixture
def anyio_backend():
    """Set the async backend for anyio tests."""
    return 'asyncio'