# SPDX-License-Identifier: MIT
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.

"""
LangChain callback handler for AegisGate security monitoring.

This module provides callback handlers that integrate AegisGate's
AI security capabilities with LangChain's callback system.
"""

import asyncio
from datetime import datetime
from typing import Any, Dict, List, Optional, Union
from uuid import UUID

from aegisgate import Client, AsyncClient
from aegisgate.models import Violation, ViolationSeverity, ViolationType


class AegisGateCallback:
    """
    LangChain callback handler for AegisGate security monitoring.

    This handler inspects LLM inputs and outputs for security violations,
    logs security events, and can block malicious content.

    Synchronous version - use AsyncAegisGateCallback for async LLMs.

    Usage:
        from langchain.llms import OpenAI
        from aegisgate.langchain import AegisGateCallback

        llm = OpenAI(temperature=0)
        callback = AegisGateCallback(base_url="http://localhost:8080")

        result = llm.predict("Hello!", callbacks=[callback])

        if callback.has_violations():
            print(f"Blocked: {callback.get_violations()}")
    """

    def __init__(
        self,
        base_url: Optional[str] = None,
        api_key: Optional[str] = None,
        timeout: float = 30.0,
        block_on_violation: bool = True,
        log_violations: bool = True,
        min_severity: str = "medium",
        violation_types: Optional[List[str]] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ):
        """
        Initialize the AegisGate callback handler.

        Args:
            base_url: AegisGate server URL (defaults to env AEGISGATE_BASE_URL)
            api_key: API key for authentication (defaults to env AEGISGATE_API_KEY)
            timeout: Request timeout in seconds
            block_on_violation: Whether to raise exception on violation
            log_violations: Whether to log violations to AegisGate
            min_severity: Minimum severity level to block ('low', 'medium', 'high', 'critical')
            violation_types: Specific violation types to check (None = all)
            metadata: Additional metadata to include in reports
        """
        self._client = Client(
            base_url=base_url,
            api_key=api_key,
            timeout=timeout,
        )
        self._block_on_violation = block_on_violation
        self._log_violations = log_violations
        self._min_severity = ViolationSeverity(min_severity)
        self._violation_types = (
            [ViolationType(vt) for vt in violation_types] if violation_types else None
        )
        self._metadata = metadata or {}
        self._violations: List[Violation] = []
        self._current_run_id: Optional[str] = None
        self._input_content: Optional[str] = None
        self._output_content: Optional[str] = None

    def __enter__(self):
        """Context manager entry."""
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        """Context manager exit - close client."""
        self._client.close()

    def has_violations(self) -> bool:
        """Check if any violations were detected."""
        return len(self._violations) > 0

    def get_violations(self) -> List[Violation]:
        """Get all detected violations."""
        return self._violations.copy()

    def clear_violations(self) -> None:
        """Clear all recorded violations."""
        self._violations.clear()

    def get_last_input(self) -> Optional[str]:
        """Get the last input content that was checked."""
        return self._input_content

    def get_last_output(self) -> Optional[str]:
        """Get the last output content that was checked."""
        return self._output_content

    def _should_block(self, violation: Violation) -> bool:
        """Determine if a violation should be blocked based on config."""
        # Check severity threshold
        severity_order = ["low", "medium", "high", "critical"]
        if severity_order.index(violation.severity.value) < severity_order.index(
            self._min_severity.value
        ):
            return False

        # Check violation type filter
        if self._violation_types and violation.type not in self._violation_types:
            return False

        return True

    def _check_content(self, content: str, content_type: str) -> List[Violation]:
        """Check content for security violations."""
        if not content:
            return []

        try:
            result = self._client.proxy.inspect_request(
                content=content,
                metadata={
                    "content_type": content_type,
                    "run_id": str(self._current_run_id),
                    **self._metadata,
                },
            )

            violations = []
            for violation in result.violations:
                if self._should_block(violation):
                    violations.append(violation)

            self._violations.extend(violations)
            return violations

        except Exception:
            # On error, don't block unless configured otherwise
            return []

    def on_llm_start(
        self,
        serialized: Dict[str, Any],
        prompts: List[str],
        **kwargs: Any,
    ) -> None:
        """Called when LLM starts running."""
        self._current_run_id = str(kwargs.get("run_id", ""))
        self._violations.clear()

        # Check each prompt for violations
        for prompt in prompts:
            self._input_content = prompt
            violations = self._check_content(prompt, "llm_input")

            if violations and self._block_on_violation:
                violation = violations[0]
                raise SecurityViolationError(
                    f"Input blocked due to security violation: {violation.type.value} - {violation.message}",
                    violations=violations,
                )

    def on_llm_end(
        self,
        response: Dict[str, Any],
        **kwargs: Any,
    ) -> None:
        """Called when LLM finishes running."""
        # Extract output from response
        generations = response.get("generations", [])
        for generation_list in generations:
            for generation in generation_list:
                text = generation.get("text", "")
                if text:
                    self._output_content = text
                    violations = self._check_content(text, "llm_output")

                    if violations and self._block_on_violation:
                        # Log but don't raise - output violations typically logged only
                        if self._log_violations:
                            for v in violations:
                                self._log_violation_event(v, "output")

    def on_llm_error(
        self,
        error: Exception,
        **kwargs: Any,
    ) -> None:
        """Called when LLM encounters an error."""
        # Log the error for security analysis
        if self._log_violations:
            pass  # Could log error details to AegisGate

    def on_chain_start(
        self,
        serialized: Dict[str, Any],
        inputs: Dict[str, Any],
        **kwargs: Any,
    ) -> None:
        """Called when chain starts running."""
        self._current_run_id = str(kwargs.get("run_id", ""))

    def on_tool_start(
        self,
        serialized: Dict[str, Any],
        input_str: str,
        **kwargs: Any,
    ) -> None:
        """Called when tool starts running."""
        violations = self._check_content(input_str, "tool_input")

        if violations and self._block_on_violation:
            violation = violations[0]
            raise SecurityViolationError(
                f"Tool input blocked: {violation.type.value} - {violation.message}",
                violations=violations,
            )

    def on_tool_end(
        self,
        output: str,
        **kwargs: Any,
    ) -> None:
        """Called when tool finishes running."""
        self._check_content(output, "tool_output")

    def on_text(
        self,
        text: str,
        **kwargs: Any,
    ) -> None:
        """Called when agent text is processed."""
        self._check_content(text, "agent_text")

    def _log_violation_event(self, violation: Violation, source: str) -> None:
        """Log a violation event to AegisGate."""
        try:
            self._client.siem.send_event(
                integration_id="langchain-callback",  # Would need proper integration setup
                event_type="security_violation",
                severity=violation.severity.value,
                source=f"langchain:{source}",
                message=violation.message,
                details={
                    "violation_type": violation.type.value,
                    "violation_id": violation.id,
                    "run_id": self._current_run_id,
                    **self._metadata,
                },
            )
        except Exception:
            pass  # Don't fail on logging errors


class AsyncAegisGateCallback:
    """
    Async LangChain callback handler for AegisGate security monitoring.

    This is the async version of AegisGateCallback for use with async LLMs.

    Usage:
        from langchain.llms import OpenAI
        from aegisgate.langchain import AsyncAegisGateCallback

        callback = AsyncAegisGateCallback(base_url="http://localhost:8080")

        async def run():
            async with callback:
                result = await llm.apredict("Hello!", callbacks=[callback])

                if callback.has_violations():
                    print(f"Violations: {callback.get_violations()}")
    """

    def __init__(
        self,
        base_url: Optional[str] = None,
        api_key: Optional[str] = None,
        timeout: float = 30.0,
        block_on_violation: bool = True,
        log_violations: bool = True,
        min_severity: str = "medium",
        violation_types: Optional[List[str]] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ):
        """
        Initialize the async AegisGate callback handler.

        Args: Same as AegisGateCallback
        """
        self._client = AsyncClient(
            base_url=base_url,
            api_key=api_key,
            timeout=timeout,
        )
        self._block_on_violation = block_on_violation
        self._log_violations = log_violations
        self._min_severity = ViolationSeverity(min_severity)
        self._violation_types = (
            [ViolationType(vt) for vt in violation_types] if violation_types else None
        )
        self._metadata = metadata or {}
        self._violations: List[Violation] = []
        self._current_run_id: Optional[str] = None
        self._input_content: Optional[str] = None
        self._output_content: Optional[str] = None

    async def __aenter__(self):
        """Async context manager entry."""
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Async context manager exit - close client."""
        await self._client.close()

    def has_violations(self) -> bool:
        """Check if any violations were detected."""
        return len(self._violations) > 0

    def get_violations(self) -> List[Violation]:
        """Get all detected violations."""
        return self._violations.copy()

    def clear_violations(self) -> None:
        """Clear all recorded violations."""
        self._violations.clear()

    def get_last_input(self) -> Optional[str]:
        """Get the last input content that was checked."""
        return self._input_content

    def get_last_output(self) -> Optional[str]:
        """Get the last output content that was checked."""
        return self._output_content

    def _should_block(self, violation: Violation) -> bool:
        """Determine if a violation should be blocked based on config."""
        severity_order = ["low", "medium", "high", "critical"]
        if severity_order.index(violation.severity.value) < severity_order.index(
            self._min_severity.value
        ):
            return False

        if self._violation_types and violation.type not in self._violation_types:
            return False

        return True

    async def _check_content(self, content: str, content_type: str) -> List[Violation]:
        """Check content for security violations asynchronously."""
        if not content:
            return []

        try:
            result = await self._client.proxy.inspect_request(
                content=content,
                metadata={
                    "content_type": content_type,
                    "run_id": str(self._current_run_id),
                    **self._metadata,
                },
            )

            violations = []
            for violation in result.violations:
                if self._should_block(violation):
                    violations.append(violation)

            self._violations.extend(violations)
            return violations

        except Exception:
            return []

    async def on_llm_start(
        self,
        serialized: Dict[str, Any],
        prompts: List[str],
        **kwargs: Any,
    ) -> None:
        """Called when LLM starts running."""
        self._current_run_id = str(kwargs.get("run_id", ""))
        self._violations.clear()

        for prompt in prompts:
            self._input_content = prompt
            violations = await self._check_content(prompt, "llm_input")

            if violations and self._block_on_violation:
                violation = violations[0]
                raise SecurityViolationError(
                    f"Input blocked due to security violation: {violation.type.value} - {violation.message}",
                    violations=violations,
                )

    async def on_llm_end(
        self,
        response: Dict[str, Any],
        **kwargs: Any,
    ) -> None:
        """Called when LLM finishes running."""
        generations = response.get("generations", [])
        for generation_list in generations:
            for generation in generation_list:
                text = generation.get("text", "")
                if text:
                    self._output_content = text
                    violations = await self._check_content(text, "llm_output")

                    if violations and self._log_violations:
                        for v in violations:
                            await self._log_violation_event(v, "output")

    async def on_llm_error(
        self,
        error: Exception,
        **kwargs: Any,
    ) -> None:
        """Called when LLM encounters an error."""
        pass

    async def on_chain_start(
        self,
        serialized: Dict[str, Any],
        inputs: Dict[str, Any],
        **kwargs: Any,
    ) -> None:
        """Called when chain starts running."""
        self._current_run_id = str(kwargs.get("run_id", ""))

    async def on_tool_start(
        self,
        serialized: Dict[str, Any],
        input_str: str,
        **kwargs: Any,
    ) -> None:
        """Called when tool starts running."""
        violations = await self._check_content(input_str, "tool_input")

        if violations and self._block_on_violation:
            violation = violations[0]
            raise SecurityViolationError(
                f"Tool input blocked: {violation.type.value} - {violation.message}",
                violations=violations,
            )

    async def on_tool_end(
        self,
        output: str,
        **kwargs: Any,
    ) -> None:
        """Called when tool finishes running."""
        await self._check_content(output, "tool_output")

    async def _log_violation_event(self, violation: Violation, source: str) -> None:
        """Log a violation event to AegisGate asynchronously."""
        try:
            await self._client.siem.send_event(
                integration_id="langchain-callback",
                event_type="security_violation",
                severity=violation.severity.value,
                source=f"langchain:{source}",
                message=violation.message,
                details={
                    "violation_type": violation.type.value,
                    "violation_id": violation.id,
                    "run_id": self._current_run_id,
                    **self._metadata,
                },
            )
        except Exception:
            pass


class SecurityViolationError(Exception):
    """Exception raised when a security violation is detected."""

    def __init__(
        self,
        message: str,
        violations: Optional[List[Violation]] = None,
    ):
        super().__init__(message)
        self.violations = violations or []

    def get_violation_summary(self) -> str:
        """Get a summary of all violations."""
        if not self.violations:
            return "No violations recorded"
        return "; ".join(
            f"{v.type.value}({v.severity.value}): {v.message}"
            for v in self.violations
        )