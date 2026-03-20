"""
Tests for LangChain callback handler integration.
"""

import pytest
from unittest.mock import MagicMock, AsyncMock, patch
from datetime import datetime

from aegisgate.langchain.callback import (
    AegisGateCallback,
    AsyncAegisGateCallback,
    SecurityViolationError,
)
from aegisgate.models import Violation, ViolationSeverity, ViolationType


class TestAegisGateCallback:
    """Tests for synchronous callback handler."""

    def test_callback_initialization_with_api_key(self):
        """Test callback initialization with API key."""
        with patch('aegisgate.langchain.callback.Client') as mock_client:
            callback = AegisGateCallback(api_key="test-key", base_url="https://test.api.io")
            
            # Check internal attributes (using _ prefix convention)
            assert callback._block_on_violation is True
            assert str(callback._min_severity.value) == "medium"
            mock_client.assert_called_once()

    def test_callback_initialization_with_config(self):
        """Test callback initialization with configuration options."""
        with patch('aegisgate.langchain.callback.Client') as mock_client:
            callback = AegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io",
                block_on_violation=False,
                min_severity="high",
                violation_types=["prompt_injection", "pii_exposure"],
                log_violations=True,
                metadata={"app": "test"}
            )
            
            assert callback._block_on_violation is False
            assert str(callback._min_severity.value) == "high"
            assert len(callback._violation_types) == 2
            assert callback._metadata["app"] == "test"

    def test_on_llm_start_no_violations(self):
        """Test on_llm_start with no violations."""
        with patch('aegisgate.langchain.callback.Client') as mock_client:
            mock_client_instance = MagicMock()
            mock_client.return_value = mock_client_instance
            
            # Mock the proxy.inspect_request to return no violations
            mock_result = MagicMock()
            mock_result.violations = []
            mock_client_instance.proxy.inspect_request.return_value = mock_result
            
            callback = AegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io"
            )
            
            # Should not raise any error
            callback.on_llm_start(
                serialized={"name": "test-llm"},
                prompts=["What is machine learning?"],
                run_id="test-run"
            )

    def test_on_llm_start_with_violations(self):
        """Test on_llm_start with detected violations."""
        with patch('aegisgate.langchain.callback.Client') as mock_client:
            mock_client_instance = MagicMock()
            mock_client.return_value = mock_client_instance
            
            # Create a violation
            mock_violation = Violation(
                id="v-001",
                type=ViolationType.PROMPT_INJECTION,
                severity=ViolationSeverity.HIGH,
                message="Potential prompt injection",
                timestamp=datetime.utcnow()
            )
            
            # Mock the proxy.inspect_request to return violations
            mock_result = MagicMock()
            mock_result.violations = [mock_violation]
            mock_client_instance.proxy.inspect_request.return_value = mock_result
            
            callback = AegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io",
                block_on_violation=True
            )
            
            # Should raise SecurityViolationError when block_on_violation is True
            with pytest.raises(SecurityViolationError):
                callback.on_llm_start(
                    serialized={"name": "test-llm"},
                    prompts=["${system.prompt}"],
                    run_id="test-run"
                )

    def test_on_llm_end(self):
        """Test on_llm_end callback."""
        with patch('aegisgate.langchain.callback.Client') as mock_client:
            mock_client_instance = MagicMock()
            mock_client.return_value = mock_client_instance
            
            mock_result = MagicMock()
            mock_result.violations = []
            mock_client_instance.proxy.inspect_request.return_value = mock_result
            
            callback = AegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io"
            )
            
            # Should not raise any error
            callback.on_llm_end(
                response={"generations": [[{"text": "AI response"}]]},
                run_id="test-run"
            )

    def test_on_llm_error(self):
        """Test on_llm_error callback."""
        with patch('aegisgate.langchain.callback.Client') as mock_client:
            mock_client.return_value = MagicMock()
            
            callback = AegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io"
            )
            
            # Should log error but not raise
            callback.on_llm_error(
                error=ValueError("Test error"),
                run_id="test-run"
            )

    def test_on_chain_start(self):
        """Test on_chain_start callback."""
        with patch('aegisgate.langchain.callback.Client') as mock_client:
            mock_client.return_value = MagicMock()
            
            callback = AegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io",
                block_on_violation=False
            )
            
            callback.on_chain_start(
                serialized={"name": "test-chain"},
                inputs={"input": "test query"},
                run_id="test-run"
            )

    def test_on_tool_start(self):
        """Test on_tool_start callback."""
        with patch('aegisgate.langchain.callback.Client') as mock_client:
            mock_client_instance = MagicMock()
            mock_client.return_value = mock_client_instance
            
            mock_result = MagicMock()
            mock_result.violations = []
            mock_client_instance.proxy.inspect_request.return_value = mock_result
            
            callback = AegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io"
            )
            
            callback.on_tool_start(
                serialized={"name": "test-tool"},
                input_str="test input",
                run_id="test-run"
            )

    def test_on_tool_end(self):
        """Test on_tool_end callback."""
        with patch('aegisgate.langchain.callback.Client') as mock_client:
            mock_client_instance = MagicMock()
            mock_client.return_value = mock_client_instance
            
            mock_result = MagicMock()
            mock_result.violations = []
            mock_client_instance.proxy.inspect_request.return_value = mock_result
            
            callback = AegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io"
            )
            
            callback.on_tool_end(
                output="tool output",
                run_id="test-run"
            )


class TestAsyncAegisGateCallback:
    """Tests for asynchronous callback handler."""

    def test_async_callback_initialization(self):
        """Test async callback initialization."""
        with patch('aegisgate.langchain.callback.AsyncClient') as mock_client:
            callback = AsyncAegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io"
            )
            
            assert callback._block_on_violation is True
            assert str(callback._min_severity.value) == "medium"
            mock_client.assert_called_once()

    @pytest.mark.asyncio
    async def test_async_on_llm_start_no_violations(self):
        """Test async on_llm_start with no violations."""
        with patch('aegisgate.langchain.callback.AsyncClient') as mock_client:
            mock_client_instance = AsyncMock()
            mock_client.return_value = mock_client_instance
            
            # Mock the proxy.inspect_request to return no violations
            mock_result = MagicMock()
            mock_result.violations = []
            mock_client_instance.proxy.inspect_request.return_value = mock_result
            
            callback = AsyncAegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io"
            )
            
            # Should not raise any error
            await callback.on_llm_start(
                serialized={"name": "test-llm"},
                prompts=["What is machine learning?"],
                run_id="test-run"
            )

    @pytest.mark.asyncio
    async def test_async_on_llm_start_with_violations(self):
        """Test async on_llm_start with detected violations."""
        with patch('aegisgate.langchain.callback.AsyncClient') as mock_client:
            mock_client_instance = AsyncMock()
            mock_client.return_value = mock_client_instance
            
            # Create a violation
            mock_violation = Violation(
                id="v-002",
                type=ViolationType.PROMPT_INJECTION,
                severity=ViolationSeverity.HIGH,
                message="Potential prompt injection",
                timestamp=datetime.utcnow()
            )
            
            # Mock the proxy.inspect_request to return violations
            mock_result = MagicMock()
            mock_result.violations = [mock_violation]
            mock_client_instance.proxy.inspect_request.return_value = mock_result
            
            callback = AsyncAegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io",
                block_on_violation=True
            )
            
            # Should raise SecurityViolationError
            with pytest.raises(SecurityViolationError):
                await callback.on_llm_start(
                    serialized={"name": "test-llm"},
                    prompts=["${system.prompt}"],
                    run_id="test-run"
                )

    @pytest.mark.asyncio
    async def test_async_on_llm_end(self):
        """Test async on_llm_end callback."""
        with patch('aegisgate.langchain.callback.AsyncClient') as mock_client:
            mock_client_instance = AsyncMock()
            mock_client.return_value = mock_client_instance
            
            mock_result = MagicMock()
            mock_result.violations = []
            mock_client_instance.proxy.inspect_request.return_value = mock_result
            
            callback = AsyncAegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io"
            )
            
            await callback.on_llm_end(
                response={"generations": [[{"text": "AI response"}]]},
                run_id="test-run"
            )

    @pytest.mark.asyncio
    async def test_async_on_chain_start(self):
        """Test async on_chain_start callback."""
        with patch('aegisgate.langchain.callback.AsyncClient') as mock_client:
            mock_client.return_value = AsyncMock()
            
            callback = AsyncAegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io"
            )
            
            await callback.on_chain_start(
                serialized={"name": "test-chain"},
                inputs={"input": "test query"},
                run_id="test-run"
            )

    @pytest.mark.asyncio
    async def test_async_on_tool_start(self):
        """Test async on_tool_start callback."""
        with patch('aegisgate.langchain.callback.AsyncClient') as mock_client:
            mock_client_instance = AsyncMock()
            mock_client.return_value = mock_client_instance
            
            mock_result = MagicMock()
            mock_result.violations = []
            mock_client_instance.proxy.inspect_request.return_value = mock_result
            
            callback = AsyncAegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io"
            )
            
            await callback.on_tool_start(
                serialized={"name": "test-tool"},
                input_str="test input",
                run_id="test-run"
            )

    @pytest.mark.asyncio
    async def test_async_on_tool_end(self):
        """Test async on_tool_end callback."""
        with patch('aegisgate.langchain.callback.AsyncClient') as mock_client:
            mock_client_instance = AsyncMock()
            mock_client.return_value = mock_client_instance
            
            mock_result = MagicMock()
            mock_result.violations = []
            mock_client_instance.proxy.inspect_request.return_value = mock_result
            
            callback = AsyncAegisGateCallback(
                api_key="test-key",
                base_url="https://test.api.io"
            )
            
            await callback.on_tool_end(
                output="tool output",
                run_id="test-run"
            )


class TestSecurityViolationError:
    """Tests for SecurityViolationError exception."""

    def test_error_creation(self):
        """Test security violation error creation."""
        violation = Violation(
            id="v-001",
            type=ViolationType.PROMPT_INJECTION,
            severity=ViolationSeverity.HIGH,
            message="Blocked content",
            timestamp=datetime.utcnow()
        )
        
        error = SecurityViolationError(
            message="Content blocked",
            violations=[violation]
        )
        
        assert "Content blocked" in str(error)
        assert len(error.violations) == 1
        assert error.violations[0].type == ViolationType.PROMPT_INJECTION

    def test_error_string_representation(self):
        """Test string representation of error."""
        error = SecurityViolationError(message="Content blocked")
        
        assert "Content blocked" in str(error)

    def test_error_get_violation_summary(self):
        """Test error get_violation_summary method."""
        violation1 = Violation(
            id="v-001",
            type=ViolationType.PROMPT_INJECTION,
            severity=ViolationSeverity.HIGH,
            message="Injection attempt",
            timestamp=datetime.utcnow()
        )
        violation2 = Violation(
            id="v-002",
            type=ViolationType.PII_EXPOSURE,
            severity=ViolationSeverity.MEDIUM,
            message="PII detected",
            timestamp=datetime.utcnow()
        )
        
        error = SecurityViolationError(
            message="Multiple violations",
            violations=[violation1, violation2]
        )
        
        summary = error.get_violation_summary()
        assert "prompt_injection" in summary
        assert "pii_exposure" in summary

    def test_error_empty_violations(self):
        """Test error with no violations."""
        error = SecurityViolationError(message="No violations")
        
        assert len(error.violations) == 0
        summary = error.get_violation_summary()
        assert "No violations" in summary