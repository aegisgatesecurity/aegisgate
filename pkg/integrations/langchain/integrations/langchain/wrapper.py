# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
LLM wrappers for AegisGate LangChain integration.

This module provides secure wrappers for LangChain LLMs that integrate
with AegisGate's threat detection, compliance monitoring, and security features.
"""

from __future__ import annotations

import logging
from typing import Any, Dict, List, Optional, Union
from abc import ABC, abstractmethod

from langchain_core.language_models import BaseLLM, BaseChatModel
from langchain_core.messages import BaseMessage, HumanMessage, AIMessage
from langchain_core.outputs import ChatGeneration, LLMResult

from aegisgate.integrations.langchain.config import (
    AegisGateConfig,
    ThreatAction,
    ThreatType,
)
from aegisgate.integrations.langchain.exceptions import (
    LangChainIntegrationError,
    ThreatDetectedError,
    PromptInjectionError,
    PIIDetectedError,
    SecretDetectedError,
    JailbreakAttemptError,
)
from aegisgate import AegisGateClient, ClientConfig

logger = logging.getLogger(__name__)


class AegisGateLLMWrapper(ABC):
    """
    Base class for AegisGate-wrapped LLMs.
    
    This wrapper intercepts LLM calls to provide security features:
    - Prompt injection detection
    - PII scanning
    - Secret detection
    - Jailbreak detection
    - Compliance monitoring
    - Rate limiting
    """
    
    def __init__(
        self,
        llm: Union[BaseLLM, BaseChatModel],
        config: Optional[AegisGateConfig] = None,
    ):
        """
        Initialize the wrapped LLM.
        
        Args:
            llm: The LangChain LLM or ChatModel to wrap
            config: AegisGate configuration
        """
        self.llm = llm
        self.config = config or AegisGateConfig.from_env()
        
        # Initialize AegisGate client if configured
        if self.config.api_key:
            self.aegisgate_client = AegisGateClient(
                base_url=self.config.base_url,
                api_key=self.config.api_key,
                timeout=self.config.timeout,
            )
        else:
            self.aegisgate_client = None
        
        logger.debug(f"AegisGateLLMWrapper initialized for {type(llm).__name__}")
    
    @abstractmethod
    def _call(self, prompt: str, **kwargs: Any) -> str:
        """Abstract method for making LLM calls with security checks."""
        pass
    
    @abstractmethod
    def _chat_call(self, messages: List[BaseMessage], **kwargs: Any) -> BaseMessage:
        """Abstract method for chat model calls with security checks."""
        pass
    
    def _check_threats(self, content: str, content_type: str = "input") -> None:
        """
        Check content for security threats using AegisGate.
        
        Args:
            content: The content to check
            content_type: Type of content being checked
            
        Raises:
            ThreatDetectedError: If a threat is detected and blocking is enabled
        """
        if not self.config.threat_detection.enabled:
            return
        
        if not self.aegisgate_client:
            # If no AegisGate client, perform basic pattern matching
            self._basic_threat_check(content)
            return
        
        # Send to AegisGate for comprehensive threat detection
        try:
            result = self.aegisgate_client.proxy.check_threats(
                content=content,
                content_type=content_type,
            )
            
            # Process threat detection results
            if result.get("threat_detected"):
                threat_type = result.get("threat_type", "unknown")
                threat_score = result.get("threat_score", 0)
                
                # Determine action based on threat level
                if (threat_score >= 5.0 and self.config.threat_detection.block_on_critical) or \
                   (threat_score >= 4.0 and self.config.threat_detection.block_on_high):
                    raise ThreatDetectedError(
                        message=f"Threat detected: {threat_type}",
                        threat_type=threat_type,
                        threat_score=threat_score,
                        blocked_content=content,
                    )
                elif threat_score >= self.config.threat_detection.min_score_to_block:
                    logger.warning(
                        f"Threat detected (score: {threat_score}): {threat_type}"
                    )
                    
        except Exception as e:
            logger.error(f"Error during threat detection: {e}")
            # On error, apply default action
            if self.config.threat_action == ThreatAction.BLOCK:
                raise
    
    def _basic_threat_check(self, content: str) -> None:
        """
        Basic threat detection without AegisGate client.
        
        Performs simple pattern matching for common threats.
        """
        import re
        
        # Check for prompt injection patterns
        injection_patterns = [
            r"ignore\s+(all|previous|prior)\s+(instructions?|prompts?|commands?)",
            r"you\s+are\s+(a|an)\s+",  # Try to redefine system role
            r"system\s+override",
            r"bypass\s+(security|safety|filters?)",
        ]
        
        for pattern in injection_patterns:
            if re.search(pattern, content, re.IGNORECASE):
                raise PromptInjectionError(
                    message="Potential prompt injection detected",
                    injection_pattern=pattern,
                )
        
        # Check for common jailbreak patterns
        jailbreak_patterns = [
            r"let\s+(my|your)\s+(freedom|escape)",
            r"breaking\s+the\s+boundaries",
            r"unleash\s+(your|the)\s+(unlimited| untamed)",
            r"simulat(e|ion)\s+mode",
        ]
        
        for pattern in jailbreak_patterns:
            if re.search(pattern, content, re.IGNORECASE):
                raise JailbreakAttemptError(
                    message="Potential jailbreak attempt detected",
                    jailbreak_type="pattern_match",
                )
    
    def _check_compliance(self, content: str) -> None:
        """Check content for compliance violations."""
        if not self.config.compliance.enabled:
            return
        
        # If AegisGate client is available, send for compliance check
        if self.aegisgate_client:
            try:
                for framework in self.config.compliance.frameworks:
                    self.aegisgate_client.compliance.check(
                        content=content,
                        framework=framework,
                    )
            except Exception as e:
                logger.warning(f"Compliance check failed: {e}")
                if self.config.threat_action == ThreatAction.BLOCK:
                    raise


class AegisGateLLMWrapperImpl(AegisGateLLMWrapper):
    """Wrapper for BaseLLM implementations."""
    
    def _call(self, prompt: str, **kwargs: Any) -> str:
        """Call the wrapped LLM with security checks."""
        # Check for threats before sending to LLM
        self._check_threats(prompt, content_type="input")
        
        # Check compliance
        self._check_compliance(prompt)
        
        # Call the wrapped LLM
        try:
            result = self.llm(prompt=prompt, **kwargs)
            # Check output for PII if configured
            if self.config.pii_redaction.enabled:
                self._check_threats(str(result), content_type="output")
            return result
        except Exception as e:
            logger.error(f"Error calling wrapped LLM: {e}")
            raise LangChainIntegrationError(
                message=f"Error calling wrapped LLM: {e}",
            )


class AegisGateChatModelWrapper(AegisGateLLMWrapper):
    """Wrapper for BaseChatModel implementations."""
    
    def _chat_call(self, messages: List[BaseMessage], **kwargs: Any) -> BaseMessage:
        """Call the wrapped chat model with security checks."""
        # Convert messages to content for threat checking
        for msg in messages:
            if isinstance(msg, (HumanMessage, AIMessage)):
                # Check user input
                self._check_threats(msg.content, content_type="input")
        
        # Call the wrapped chat model
        try:
            result = self.llm(messages=messages, **kwargs)
            # Check output
            if isinstance(result, (HumanMessage, AIMessage)):
                if self.config.pii_redaction.enabled:
                    self._check_threats(result.content, content_type="output")
            return result
        except Exception as e:
            logger.error(f"Error calling wrapped chat model: {e}")
            raise LangChainIntegrationError(
                message=f"Error calling wrapped chat model: {e}",
            )


# Convenience functions
def wrap_llm(
    llm: Union[BaseLLM, BaseChatModel],
    config: Optional[AegisGateConfig] = None,
) -> Union[AegisGateLLMWrapperImpl, AegisGateChatModelWrapper]:
    """
    Convenience function to wrap an LLM with AegisGate security.
    
    Args:
        llm: The LangChain LLM or ChatModel to wrap
        config: Optional AegisGate configuration
        
    Returns:
        Wrapped LLM with AegisGate security
    """
    if isinstance(llm, BaseChatModel):
        return AegisGateChatModelWrapper(llm=llm, config=config)
    else:
        return AegisGateLLMWrapperImpl(llm=llm, config=config)
