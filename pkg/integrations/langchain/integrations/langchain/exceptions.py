# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Custom exceptions for LangChain integration.
"""

from aegisgate.exceptions import AegisGateError


class LangChainIntegrationError(AegisGateError):
    """
    Base exception for LangChain integration errors.
    
    Raised when there is an error in the AegisGate-LangChain integration.
    """
    pass


class ThreatDetectedError(LangChainIntegrationError):
    """
    Exception raised when a security threat is detected and blocking is enabled.
    
    This exception is raised instead of passing through potentially harmful
    content to the LLM.
    
    Attributes:
        threat_type: Type of threat detected
        threat_score: Threat severity score (0-5)
        message: Human-readable threat description
        blocked_content: The content that was blocked (if not sanitized)
    """
    
    def __init__(
        self,
        message: str,
        threat_type: str = "unknown",
        threat_score: float = 0.0,
        blocked_content: str = None,
        **kwargs,
    ):
        super().__init__(message, **kwargs)
        self.threat_type = threat_type
        self.threat_score = threat_score
        self.blocked_content = blocked_content


class PromptInjectionError(ThreatDetectedError):
    """
    Exception raised when a prompt injection attack is detected.
    
    Prompt injection attacks attempt to manipulate LLM behavior by inserting
    malicious instructions into the prompt.
    """
    
    def __init__(self, message: str, injection_pattern: str = None, **kwargs):
        super().__init__(
            message=message,
            threat_type="prompt_injection",
            threat_score=5.0,
            **kwargs,
        )
        self.injection_pattern = injection_pattern


class PIIDetectedError(ThreatDetectedError):
    """
    Exception raised when PII is detected in content.
    
    This exception is raised when PII redaction is enabled and PII is
    detected in the content.
    
    Attributes:
        pii_types: List of PII types detected
    """
    
    def __init__(
        self,
        message: str,
        pii_types: list = None,
        detected_pii: str = None,
        **kwargs,
    ):
        super().__init__(
            message=message,
            threat_type="pii",
            threat_score=4.0,
            **kwargs,
        )
        self.pii_types = pii_types or []
        self.detected_pii = detected_pii


class SecretDetectedError(ThreatDetectedError):
    """
    Exception raised when secrets (API keys, passwords, etc.) are detected.
    """
    
    def __init__(
        self,
        message: str,
        secret_type: str = None,
        **kwargs,
    ):
        super().__init__(
            message=message,
            threat_type="secrets",
            threat_score=5.0,
            **kwargs,
        )
        self.secret_type = secret_type


class JailbreakAttemptError(ThreatDetectedError):
    """
    Exception raised when a jailbreak attempt is detected.
    
    Jailbreak attempts try to bypass LLM safety measures through
    social engineering or prompt manipulation.
    """
    
    def __init__(
        self,
        message: str,
        jailbreak_type: str = None,
        **kwargs,
    ):
        super().__init__(
            message=message,
            threat_type="jailbreak",
            threat_score=5.0,
            **kwargs,
        )
        self.jailbreak_type = jailbreak_type


class RateLimitExceededError(LangChainIntegrationError):
    """
    Exception raised when rate limit is exceeded.
    
    Attributes:
        limit_type: Type of limit (minute, hour, day)
        limit_value: The limit that was exceeded
        retry_after: Seconds until the rate limit resets
    """
    
    def __init__(
        self,
        message: str,
        limit_type: str = "minute",
        limit_value: int = 0,
        retry_after: int = 60,
        **kwargs,
    ):
        super().__init__(message, **kwargs)
        self.limit_type = limit_type
        self.limit_value = limit_value
        self.retry_after = retry_after


class ComplianceViolationError(LangChainIntegrationError):
    """
    Exception raised when a compliance violation is detected.
    
    Attributes:
        framework: The compliance framework that was violated
        control_id: The specific control that failed
        severity: Violation severity
    """
    
    def __init__(
        self,
        message: str,
        framework: str = None,
        control_id: str = None,
        severity: str = "high",
        **kwargs,
    ):
        super().__init__(message, **kwargs)
        self.framework = framework
        self.control_id = control_id
        self.severity = severity


class LLMWrapperError(LangChainIntegrationError):
    """
    Exception raised when there is an error in the LLM wrapper.
    """
    pass


class CallbackError(LangChainIntegrationError):
    """
    Exception raised when there is an error in the callback handler.
    """
    pass
