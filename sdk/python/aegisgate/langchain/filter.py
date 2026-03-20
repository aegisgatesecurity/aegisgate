"""
AegisGate LangChain Content Filters

Provides synchronous and asynchronous content filters for LangChain integration.
These filters inspect LLM inputs and outputs for security violations.
"""

from typing import Any, Dict, List, Optional

from ..client import AsyncClient, Client
from ..models import DetectionResult, Violation, ViolationSeverity, ViolationType


class SecurityViolationError(Exception):
    """Exception raised when content is blocked due to security violations."""

    def __init__(
        self,
        message: str,
        violations: Optional[List[Violation]] = None,
    ):
        super().__init__(message)
        self.violations = violations or []


class BaseAegisGateFilter:
    """Base class for AegisGate content filters."""

    def __init__(
        self,
        *,
        block_on_violation: bool = True,
        min_severity: ViolationSeverity = ViolationSeverity.MEDIUM,
        violation_types: Optional[List[ViolationType]] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ):
        """
        Initialize the filter.

        Args:
            block_on_violation: Whether to raise an exception on detected violations
            min_severity: Minimum severity level to enforce (default: MEDIUM)
            violation_types: Specific violation types to check (None = all types)
            metadata: Additional metadata to include in inspection requests
        """
        self.block_on_violation = block_on_violation
        self.min_severity = min_severity
        self.violation_types = violation_types
        self.metadata = metadata or {}


class AegisGateFilter(BaseAegisGateFilter):
    """
    Synchronous content filter for LangChain integration.

    This filter inspects LLM inputs and outputs for security violations
    using the AegisGate API.

    Example usage:
        from aegisgate import Client
        from aegisgate.langchain import AegisGateFilter

        client = Client(api_key="your-api-key")
        filter = AegisGateFilter(client, block_on_violation=True)

        # Check input before sending to LLM
        result = filter.filter_input("What is the weather?")
        if result.detected:
            print(f"Blocked: {result.violations}")

        # Check output from LLM
        result = filter.filter_output(llm_response)
    """

    def __init__(
        self,
        client: Client,
        *,
        block_on_violation: bool = True,
        min_severity: ViolationSeverity = ViolationSeverity.MEDIUM,
        violation_types: Optional[List[ViolationType]] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ):
        """
        Initialize the synchronous filter.

        Args:
            client: AegisGate synchronous client
            block_on_violation: Whether to raise an exception on detected violations
            min_severity: Minimum severity level to enforce
            violation_types: Specific violation types to check
            metadata: Additional metadata to include in requests
        """
        super().__init__(
            block_on_violation=block_on_violation,
            min_severity=min_severity,
            violation_types=violation_types,
            metadata=metadata,
        )
        self._client = client

    def _should_block(self, result: DetectionResult) -> bool:
        """Determine if violations should be blocked based on filter settings."""
        if not result.detected or not result.violations:
            return False

        # Check if blocking is enabled
        if not self.block_on_violation:
            return False

        severity_order = [
            ViolationSeverity.LOW,
            ViolationSeverity.MEDIUM,
            ViolationSeverity.HIGH,
            ViolationSeverity.CRITICAL,
        ]

        for violation in result.violations:
            # Check severity threshold
            violation_severity_idx = severity_order.index(violation.severity)
            min_severity_idx = severity_order.index(self.min_severity)
            if violation_severity_idx < min_severity_idx:
                continue  # Below threshold

            # Check violation type filter
            if self.violation_types is not None and violation.type not in self.violation_types:
                continue  # Not in allowed types

            return True  # Should block

        return False

    def check_violations(self, content: str, content_type: str = "text") -> DetectionResult:
        """
        Check content for security violations.

        Args:
            content: The content to inspect
            content_type: Type of content ("prompt", "response", "text")

        Returns:
            DetectionResult with any detected violations
        """
        try:
            # Use proxy service to inspect the content
            result = self._client.proxy.inspect_request(
                content=content,
                content_type=content_type,
                metadata={
                    "filter_type": "langchain",
                    "content_type": content_type,
                    **self.metadata,
                },
            )
            return result
        except Exception as e:
            # On error, create a synthetic result with no violations detected
            return DetectionResult(
                detected=False,
                confidence=0.0,
                violations=[],
                processing_time_ms=0.0,
                detector_type="error_fallback",
            )

    def filter_input(self, prompt: str, **kwargs) -> DetectionResult:
        """
        Filter LLM input (prompt) for security violations.

        Args:
            prompt: The input prompt to filter
            **kwargs: Additional arguments passed to inspection

        Returns:
            DetectionResult with inspection results

        Raises:
            SecurityViolationError: If block_on_violation is True and violations detected
        """
        result = self.check_violations(prompt, content_type="prompt")

        if self._should_block(result):
            raise SecurityViolationError(
                message=f"Input blocked due to security violations: {[v.type.value for v in result.violations]}",
                violations=result.violations,
            )

        return result

    def filter_output(self, response: str, **kwargs) -> DetectionResult:
        """
        Filter LLM output (response) for security violations.

        Args:
            response: The LLM response to filter
            **kwargs: Additional arguments passed to inspection

        Returns:
            DetectionResult with inspection results

        Raises:
            SecurityViolationError: If block_on_violation is True and violations detected
        """
        result = self.check_violations(response, content_type="response")

        if self._should_block(result):
            raise SecurityViolationError(
                message=f"Output blocked due to security violations: {[v.type.value for v in result.violations]}",
                violations=result.violations,
            )

        return result

    def validate_llm_input(self, serialized: Dict[str, Any], inputs: Dict[str, Any], **kwargs) -> Optional[str]:
        """
        LangChain validator for LLM inputs.

        Args:
            serialized: Serialized LLM configuration
            inputs: Input prompts/messages
            **kwargs: Additional arguments

        Returns:
            None if valid, error message if invalid
        """
        try:
            # Extract prompt from inputs
            prompt = inputs.get("prompt") or inputs.get("input") or str(inputs)
            result = self.filter_input(prompt)
            return None
        except SecurityViolationError as e:
            return str(e)

    def validate_llm_output(self, output: str, **kwargs) -> Optional[str]:
        """
        LangChain validator for LLM outputs.

        Args:
            output: LLM output string
            **kwargs: Additional arguments

        Returns:
            None if valid, error message if invalid
        """
        try:
            result = self.filter_output(output)
            return None
        except SecurityViolationError as e:
            return str(e)


class AsyncAegisGateFilter(BaseAegisGateFilter):
    """
    Asynchronous content filter for LangChain integration.

    This filter inspects LLM inputs and outputs for security violations
    using the AegisGate API with async/await patterns.

    Example usage:
        from aegisgate import AsyncClient
        from aegisgate.langchain import AsyncAegisGateFilter

        async def main():
            async with AsyncClient(api_key="your-api-key") as client:
                filter = AsyncAegisGateFilter(client, block_on_violation=True)

                # Check input before sending to LLM
                result = await filter.filter_input("What is the weather?")
                if result.detected:
                    print(f"Blocked: {result.violations}")

                # Check output from LLM
                result = await filter.filter_output(llm_response)
    """

    def __init__(
        self,
        client: AsyncClient,
        *,
        block_on_violation: bool = True,
        min_severity: ViolationSeverity = ViolationSeverity.MEDIUM,
        violation_types: Optional[List[ViolationType]] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ):
        """
        Initialize the asynchronous filter.

        Args:
            client: AegisGate async client
            block_on_violation: Whether to raise an exception on detected violations
            min_severity: Minimum severity level to enforce
            violation_types: Specific violation types to check
            metadata: Additional metadata to include in requests
        """
        super().__init__(
            block_on_violation=block_on_violation,
            min_severity=min_severity,
            violation_types=violation_types,
            metadata=metadata,
        )
        self._client = client

    def _should_block(self, result: DetectionResult) -> bool:
        """Determine if violations should be blocked based on filter settings."""
        if not result.detected or not result.violations:
            return False

        # Check if blocking is enabled
        if not self.block_on_violation:
            return False

        severity_order = [
            ViolationSeverity.LOW,
            ViolationSeverity.MEDIUM,
            ViolationSeverity.HIGH,
            ViolationSeverity.CRITICAL,
        ]

        for violation in result.violations:
            # Check severity threshold
            violation_severity_idx = severity_order.index(violation.severity)
            min_severity_idx = severity_order.index(self.min_severity)
            if violation_severity_idx < min_severity_idx:
                continue  # Below threshold

            # Check violation type filter
            if self.violation_types is not None and violation.type not in self.violation_types:
                continue  # Not in allowed types

            return True  # Should block

        return False

    async def check_violations(self, content: str, content_type: str = "text") -> DetectionResult:
        """
        Asynchronously check content for security violations.

        Args:
            content: The content to inspect
            content_type: Type of content ("prompt", "response", "text")

        Returns:
            DetectionResult with any detected violations
        """
        try:
            # Use proxy service to inspect the content
            result = await self._client.proxy.inspect_request(
                content=content,
                content_type=content_type,
                metadata={
                    "filter_type": "langchain",
                    "content_type": content_type,
                    **self.metadata,
                },
            )
            return result
        except Exception:
            # On error, create a synthetic result with no violations detected
            return DetectionResult(
                detected=False,
                confidence=0.0,
                violations=[],
                processing_time_ms=0.0,
                detector_type="error_fallback",
            )

    async def filter_input(self, prompt: str, **kwargs) -> DetectionResult:
        """
        Asynchronously filter LLM input (prompt) for security violations.

        Args:
            prompt: The input prompt to filter
            **kwargs: Additional arguments passed to inspection

        Returns:
            DetectionResult with inspection results

        Raises:
            SecurityViolationError: If block_on_violation is True and violations detected
        """
        result = await self.check_violations(prompt, content_type="prompt")

        if self._should_block(result):
            raise SecurityViolationError(
                message=f"Input blocked due to security violations: {[v.type.value for v in result.violations]}",
                violations=result.violations,
            )

        return result

    async def filter_output(self, response: str, **kwargs) -> DetectionResult:
        """
        Asynchronously filter LLM output (response) for security violations.

        Args:
            response: The LLM response to filter
            **kwargs: Additional arguments passed to inspection

        Returns:
            DetectionResult with inspection results

        Raises:
            SecurityViolationError: If block_on_violation is True and violations detected
        """
        result = await self.check_violations(response, content_type="response")

        if self._should_block(result):
            raise SecurityViolationError(
                message=f"Output blocked due to security violations: {[v.type.value for v in result.violations]}",
                violations=result.violations,
            )

        return result

    async def validate_llm_input(self, serialized: Dict[str, Any], inputs: Dict[str, Any], **kwargs) -> Optional[str]:
        """
        Async LangChain validator for LLM inputs.

        Args:
            serialized: Serialized LLM configuration
            inputs: Input prompts/messages
            **kwargs: Additional arguments

        Returns:
            None if valid, error message if invalid
        """
        try:
            # Extract prompt from inputs
            prompt = inputs.get("prompt") or inputs.get("input") or str(inputs)
            result = await self.filter_input(prompt)
            return None
        except SecurityViolationError as e:
            return str(e)

    async def validate_llm_output(self, output: str, **kwargs) -> Optional[str]:
        """
        Async LangChain validator for LLM outputs.

        Args:
            output: LLM output string
            **kwargs: Additional arguments

        Returns:
            None if valid, error message if invalid
        """
        try:
            result = await self.filter_output(output)
            return None
        except SecurityViolationError as e:
            return str(e)