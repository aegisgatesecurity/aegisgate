"""
Tests for LangChain content filter integration.
"""

import pytest
from unittest.mock import MagicMock, AsyncMock, patch
from datetime import datetime

from aegisgate.langchain.filter import (
    AegisGateFilter,
    AsyncAegisGateFilter,
    SecurityViolationError,
)
from aegisgate.models import Violation, ViolationSeverity, ViolationType, DetectionResult


class TestAegisGateFilter:
    """Tests for synchronous content filter."""

    @pytest.fixture
    def mock_client(self):
        """Create a mock synchronous client."""
        client = MagicMock()
        return client

    @pytest.fixture
    def mock_detection_result_no_violations(self):
        """Create a mock detection result with no violations."""
        result = MagicMock(spec=DetectionResult)
        result.detected = False
        result.violations = []
        result.confidence = 0.0
        return result

    @pytest.fixture
    def mock_detection_result_with_violations(self):
        """Create a mock detection result with violations."""
        violation = MagicMock(spec=Violation)
        violation.id = "v-001"
        violation.type = ViolationType.PROMPT_INJECTION
        violation.severity = ViolationSeverity.HIGH
        violation.message = "Potential prompt injection"
        violation.confidence = 0.95
        
        result = MagicMock(spec=DetectionResult)
        result.detected = True
        result.violations = [violation]
        result.confidence = 0.95
        return result

    def test_filter_initialization(self, mock_client):
        """Test filter initialization."""
        filter = AegisGateFilter(mock_client)
        
        assert filter.block_on_violation is True
        assert filter.min_severity == ViolationSeverity.MEDIUM
        assert filter.violation_types is None

    def test_filter_configuration(self, mock_client):
        """Test filter configuration options."""
        filter = AegisGateFilter(
            mock_client,
            block_on_violation=False,
            min_severity=ViolationSeverity.HIGH,
            violation_types=[ViolationType.PROMPT_INJECTION],
            metadata={"app": "test"}
        )
        
        assert filter.block_on_violation is False
        assert filter.min_severity == ViolationSeverity.HIGH
        assert filter.violation_types == [ViolationType.PROMPT_INJECTION]
        assert filter.metadata["app"] == "test"

    def test_filter_input_no_violations(self, mock_client, mock_detection_result_no_violations):
        """Test filter_input with no violations."""
        mock_client.proxy.inspect_request.return_value = mock_detection_result_no_violations
        
        filter = AegisGateFilter(mock_client, block_on_violation=True)
        result = filter.filter_input("What is machine learning?")
        
        assert result.detected is False
        mock_client.proxy.inspect_request.assert_called_once()

    def test_filter_input_with_violations_blocking(self, mock_client, mock_detection_result_with_violations):
        """Test filter_input with violations when blocking is enabled."""
        mock_client.proxy.inspect_request.return_value = mock_detection_result_with_violations
        
        filter = AegisGateFilter(mock_client, block_on_violation=True)
        
        with pytest.raises(SecurityViolationError) as exc_info:
            filter.filter_input("malicious prompt")
        
        assert len(exc_info.value.violations) > 0

    def test_filter_input_with_violations_non_blocking(self, mock_client, mock_detection_result_with_violations):
        """Test filter_input with violations when blocking is disabled."""
        mock_detection_result_with_violations.detected = True
        mock_detection_result_with_violations.violations[0].severity = ViolationSeverity.CRITICAL
        mock_client.proxy.inspect_request.return_value = mock_detection_result_with_violations
        
        filter = AegisGateFilter(mock_client, block_on_violation=False)
        result = filter.filter_input("malicious prompt")
        
        # Should return result without raising
        assert result.detected is True
        assert len(result.violations) > 0

    def test_filter_output_no_violations(self, mock_client, mock_detection_result_no_violations):
        """Test filter_output with no violations."""
        mock_client.proxy.inspect_request.return_value = mock_detection_result_no_violations
        
        filter = AegisGateFilter(mock_client)
        result = filter.filter_output("Machine learning is a subset of AI.")
        
        assert result.detected is False
        mock_client.proxy.inspect_request.assert_called_once()

    def test_filter_output_with_violations_blocking(self, mock_client, mock_detection_result_with_violations):
        """Test filter_output with violations when blocking is enabled."""
        mock_detection_result_with_violations.violations[0].severity = ViolationSeverity.CRITICAL
        mock_client.proxy.inspect_request.return_value = mock_detection_result_with_violations
        
        filter = AegisGateFilter(mock_client, block_on_violation=True)
        
        with pytest.raises(SecurityViolationError) as exc_info:
            filter.filter_output("malicious output")
        
        assert len(exc_info.value.violations) > 0

    def test_severity_threshold_filtering(self, mock_client):
        """Test that low severity violations are filtered based on threshold."""
        # Create violation with LOW severity
        low_violation = MagicMock(spec=Violation)
        low_violation.id = "v-low"
        low_violation.type = ViolationType.PROMPT_INJECTION
        low_violation.severity = ViolationSeverity.LOW
        low_violation.message = "Low severity issue"
        
        result = MagicMock(spec=DetectionResult)
        result.detected = True
        result.violations = [low_violation]
        result.confidence = 0.5
        
        mock_client.proxy.inspect_request.return_value = result
        
        # Filter with HIGH threshold should not block LOW violations
        filter = AegisGateFilter(
            mock_client,
            block_on_violation=True,
            min_severity=ViolationSeverity.HIGH
        )
        
        # Should not raise because severity is below threshold
        check_result = filter.filter_input("test prompt")
        assert check_result is not None

    def test_violation_type_filtering(self, mock_client):
        """Test that only specified violation types are checked."""
        # Create PII violation
        pii_violation = MagicMock(spec=Violation)
        pii_violation.id = "v-pii"
        pii_violation.type = ViolationType.PII_EXPOSURE
        pii_violation.severity = ViolationSeverity.HIGH
        pii_violation.message = "PII detected"
        
        result = MagicMock(spec=DetectionResult)
        result.detected = True
        result.violations = [pii_violation]
        result.confidence = 0.9
        
        mock_client.proxy.inspect_request.return_value = result
        
        # Filter for only PROMPT_INJECTION should ignore PII when blocking
        filter = AegisGateFilter(
            mock_client,
            block_on_violation=True,
            violation_types=[ViolationType.PROMPT_INJECTION]
        )
        
        # Should not raise because violation type doesn't match filter
        check_result = filter.filter_input("test prompt")
        assert check_result is not None

    def test_check_violations_error_handling(self, mock_client):
        """Test error handling in check_violations."""
        mock_client.proxy.inspect_request.side_effect = Exception("API Error")
        
        filter = AegisGateFilter(mock_client)
        result = filter.check_violations("test content")
        
        # Should return empty result on error
        assert result.detected is False
        assert len(result.violations) == 0

    def test_validate_llm_input(self, mock_client, mock_detection_result_no_violations):
        """Test validate_llm_input method."""
        mock_client.proxy.inspect_request.return_value = mock_detection_result_no_violations
        
        filter = AegisGateFilter(mock_client)
        result = filter.validate_llm_input(serialized={}, inputs={"prompt": "test"})
        
        # Should return None for valid input
        assert result is None

    def test_validate_llm_output(self, mock_client, mock_detection_result_no_violations):
        """Test validate_llm_output method."""
        mock_client.proxy.inspect_request.return_value = mock_detection_result_no_violations
        
        filter = AegisGateFilter(mock_client)
        result = filter.validate_llm_output("test output")
        
        # Should return None for valid output
        assert result is None


class TestAsyncAegisGateFilter:
    """Tests for asynchronous content filter."""

    @pytest.fixture
    def mock_async_client(self):
        """Create a mock asynchronous client."""
        client = AsyncMock()
        return client

    @pytest.fixture
    def mock_detection_result_no_violations(self):
        """Create a mock detection result with no violations."""
        result = MagicMock(spec=DetectionResult)
        result.detected = False
        result.violations = []
        result.confidence = 0.0
        return result

    @pytest.fixture
    def mock_detection_result_with_violations(self):
        """Create a mock detection result with violations."""
        violation = MagicMock(spec=Violation)
        violation.id = "v-001"
        violation.type = ViolationType.PROMPT_INJECTION
        violation.severity = ViolationSeverity.HIGH
        violation.message = "Potential prompt injection"
        
        result = MagicMock(spec=DetectionResult)
        result.detected = True
        result.violations = [violation]
        result.confidence = 0.95
        return result

    def test_async_filter_initialization(self, mock_async_client):
        """Test async filter initialization."""
        filter = AsyncAegisGateFilter(mock_async_client)
        
        assert filter.block_on_violation is True
        assert filter.min_severity == ViolationSeverity.MEDIUM
        assert filter.violation_types is None

    @pytest.mark.asyncio
    async def test_async_filter_input_no_violations(self, mock_async_client, mock_detection_result_no_violations):
        """Test async filter_input with no violations."""
        mock_async_client.proxy.inspect_request.return_value = mock_detection_result_no_violations
        
        filter = AsyncAegisGateFilter(mock_async_client)
        result = await filter.filter_input("What is machine learning?")
        
        assert result.detected is False
        mock_async_client.proxy.inspect_request.assert_called_once()

    @pytest.mark.asyncio
    async def test_async_filter_input_with_violations(self, mock_async_client, mock_detection_result_with_violations):
        """Test async filter_input with violations."""
        mock_detection_result_with_violations.violations[0].severity = ViolationSeverity.CRITICAL
        mock_async_client.proxy.inspect_request.return_value = mock_detection_result_with_violations
        
        filter = AsyncAegisGateFilter(mock_async_client, block_on_violation=True)
        
        with pytest.raises(SecurityViolationError):
            await filter.filter_input("malicious prompt")

    @pytest.mark.asyncio
    async def test_async_filter_output_no_violations(self, mock_async_client, mock_detection_result_no_violations):
        """Test async filter_output with no violations."""
        mock_async_client.proxy.inspect_request.return_value = mock_detection_result_no_violations
        
        filter = AsyncAegisGateFilter(mock_async_client)
        result = await filter.filter_output("AI response")
        
        assert result.detected is False

    @pytest.mark.asyncio
    async def test_async_filter_output_with_violations(self, mock_async_client, mock_detection_result_with_violations):
        """Test async filter_output with violations."""
        mock_detection_result_with_violations.violations[0].severity = ViolationSeverity.CRITICAL
        mock_async_client.proxy.inspect_request.return_value = mock_detection_result_with_violations
        
        filter = AsyncAegisGateFilter(mock_async_client, block_on_violation=True)
        
        with pytest.raises(SecurityViolationError):
            await filter.filter_output("malicious output")

    @pytest.mark.asyncio
    async def test_async_check_violations(self, mock_async_client, mock_detection_result_no_violations):
        """Test async check_violations method."""
        mock_async_client.proxy.inspect_request.return_value = mock_detection_result_no_violations
        
        filter = AsyncAegisGateFilter(mock_async_client)
        result = await filter.check_violations("test content", content_type="prompt")
        
        assert result.detected is False
        mock_async_client.proxy.inspect_request.assert_called_once()

    @pytest.mark.asyncio
    async def test_async_check_violations_error_handling(self, mock_async_client):
        """Test async error handling."""
        mock_async_client.proxy.inspect_request.side_effect = Exception("API Error")
        
        filter = AsyncAegisGateFilter(mock_async_client)
        result = await filter.check_violations("test content")
        
        assert result.detected is False

    @pytest.mark.asyncio
    async def test_async_validate_llm_input(self, mock_async_client, mock_detection_result_no_violations):
        """Test async validate_llm_input method."""
        mock_async_client.proxy.inspect_request.return_value = mock_detection_result_no_violations
        
        filter = AsyncAegisGateFilter(mock_async_client)
        result = await filter.validate_llm_input(serialized={}, inputs={"prompt": "test"})
        
        assert result is None

    @pytest.mark.asyncio
    async def test_async_validate_llm_output(self, mock_async_client, mock_detection_result_no_violations):
        """Test async validate_llm_output method."""
        mock_async_client.proxy.inspect_request.return_value = mock_detection_result_no_violations
        
        filter = AsyncAegisGateFilter(mock_async_client)
        result = await filter.validate_llm_output("test output")
        
        assert result is None


class TestBaseAegisGateFilter:
    """Tests for base filter functionality."""

    def test_should_block_logic(self):
        """Test _should_block internal method logic."""
        mock_client = MagicMock()
        
        # Create filter with specific configuration
        filter = AegisGateFilter(
            mock_client,
            block_on_violation=True,
            min_severity=ViolationSeverity.MEDIUM
        )
        
        # Test with violations below threshold
        low_violation = MagicMock(spec=Violation)
        low_violation.severity = ViolationSeverity.LOW
        low_violation.type = ViolationType.PROMPT_INJECTION
        
        result_low = MagicMock(spec=DetectionResult)
        result_low.detected = True
        result_low.violations = [low_violation]
        
        # Should not block LOW severity when min is MEDIUM
        assert filter._should_block(result_low) is False
        
        # Test with violations at threshold
        high_violation = MagicMock(spec=Violation)
        high_violation.severity = ViolationSeverity.HIGH
        high_violation.type = ViolationType.PROMPT_INJECTION
        
        result_high = MagicMock(spec=DetectionResult)
        result_high.detected = True
        result_high.violations = [high_violation]
        
        # Should block HIGH severity when min is MEDIUM
        assert filter._should_block(result_high) is True

    def test_violation_type_configuration(self):
        """Test violation type filtering configuration."""
        mock_client = MagicMock()
        
        # Create filter with specific violation types
        filter = AegisGateFilter(
            mock_client,
            violation_types=[
                ViolationType.PROMPT_INJECTION,
                ViolationType.TOXIC_CONTENT
            ]
        )
        
        assert ViolationType.PROMPT_INJECTION in filter.violation_types
        assert ViolationType.TOXIC_CONTENT in filter.violation_types
        assert ViolationType.PII_EXPOSURE not in filter.violation_types


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

    def test_error_empty_violations(self):
        """Test error with no violations."""
        error = SecurityViolationError(message="No violations")
        
        assert len(error.violations) == 0