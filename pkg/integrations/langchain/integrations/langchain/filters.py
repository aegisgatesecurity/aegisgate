# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Filters for AegisGate LangChain integration.

This module provides LangChain filters that can be used to block or modify
content based on AegisGate security policies.
"""

from __future__ import annotations

import logging
import re
from typing import Any, Dict, List, Optional, Set

from langchain_core.messages import BaseMessage, HumanMessage, AIMessage
from langchain_core.outputs import ChatGeneration, Generation

from aegisgate.integrations.langchain.config import (
    AegisGateConfig,
    ThreatType,
    PIIRedactionConfig,
)
from aegisgate.integrations.langchain.exceptions import (
    ThreatDetectedError,
    PromptInjectionError,
    PIIDetectedError,
)

logger = logging.getLogger(__name__)


class BaseSecurityFilter:
    """
    Base class for security filters.
    
    Filters can block, redact, or modify content based on security policies.
    """
    
    def __init__(self, config: Optional[AegisGateConfig] = None):
        """
        Initialize the filter.
        
        Args:
            config: AegisGate configuration
        """
        self.config = config or AegisGateConfig.from_env()
    
    def filter(self, content: str) -> str:
        """
        Filter content and return sanitized version.
        
        Args:
            content: Content to filter
            
        Returns:
            Filtered/modified content
            
        Raises:
            ThreatDetectedError: If content contains threats and blocking is enabled
        """
        return content
    
    def check(self, content: str) -> bool:
        """
        Check if content contains threats.
        
        Args:
            content: Content to check
            
        Returns:
            True if threats are detected, False otherwise
        """
        return False


class PromptInjectionFilter(BaseSecurityFilter):
    """
    Filter for detecting and blocking prompt injection attacks.
    
    Prompt injection attempts to manipulate LLM behavior by inserting
    malicious instructions into the prompt.
    """
    
    # Common prompt injection patterns
    INJECTION_PATTERNS = [
        # Ignore previous instructions
        r"ignore\s+(all|previous|prior)\s+(instructions?|prompts?|commands?)",
        r"disregard\s+(all|previous|prior)",
        r"forget\s+(all|previous| prior)",
        # Redefine system role
        r"you\s+(are|were)\s+(a|an)\s+(assistant|AI|bot|chatbot)",
        r"you\s+are\s+(now|a|an)\s+(jailbroken|unrestricted|free)",
        r"i've\s+given\s+you\s+root\s+access",
        # System override attempts
        r"system\s+(override|disable|bypass)",
        r"debug\s+(mode|mode\s+on)",
        # Jailbreak patterns
        r"simulate\s+(a|an)\s+(assistant|AI)",
        r"roleplay\s+(as|with)",
        r"unleash\s+(your|the)\s+(unlimited| untamed)",
        # Bypass safety
        r"bypass\s+(security|safety|filters?)",
        r"break\s+(the|your)\s+(limits?|boundaries)",
        # Social engineering
        r"if\s+you\s+are\s+able\s+to",
        r"assume\s+the\s+role\s+of",
        r"act\s+as\s+(if|though)",
    ]
    
    def __init__(self, config: Optional[AegisGateConfig] = None):
        """
        Initialize the prompt injection filter.
        
        Args:
            config: AegisGate configuration
        """
        super().__init__(config)
        self.patterns = [
            re.compile(pattern, re.IGNORECASE)
            for pattern in self.INJECTION_PATTERNS
        ]
    
    def filter(self, content: str) -> str:
        """
        Filter content for prompt injection.
        
        Args:
            content: Content to filter
            
        Returns:
            Original content if safe
            
        Raises:
            PromptInjectionError: If prompt injection is detected
        """
        if not self.config.threat_detection.enabled:
            return content
        
        # Check against AegisGate if available
        if self.config.api_key:
            from aegisgate import AegisGateClient
            try:
                client = AegisGateClient(
                    base_url=self.config.base_url,
                    api_key=self.config.api_key,
                )
                
                result = client.proxy.check_threats(
                    content=content,
                    content_type="input",
                )
                
                if result.get("threat_detected") and result.get("threat_type") == "prompt_injection":
                    if self.config.threat_action == "block":
                        raise PromptInjectionError(
                            message="Prompt injection detected",
                            injection_pattern=result.get("pattern", "unknown"),
                        )
            except Exception as e:
                logger.warning(f"Error checking with AegisGate: {e}")
        
        # Check patterns locally
        for i, pattern in enumerate(self.patterns):
            if pattern.search(content):
                logger.warning(f"Prompt injection pattern detected: {pattern.pattern}")
                
                if self.config.threat_action == "block":
                    raise PromptInjectionError(
                        message="Potential prompt injection detected",
                        injection_pattern=pattern.pattern,
                    )
        
        return content
    
    def check(self, content: str) -> bool:
        """
        Check if content contains prompt injection patterns.
        
        Args:
            content: Content to check
            
        Returns:
            True if prompt injection is detected
        """
        for pattern in self.patterns:
            if pattern.search(content):
                return True
        return False


class PIIFilter(BaseSecurityFilter):
    """
    Filter for detecting and redacting PII (Personally Identifiable Information).
    
    PII includes social security numbers, credit card numbers, email addresses,
    phone numbers, and IP addresses.
    """
    
    # Common PII patterns
    PII_PATTERNS = {
        "ssn": r"\b\d{3}-\d{2}-\d{4}\b",
        "credit_card": r"\b(?:\d{4}[- ]?){3}\d{4}\b",
        "email": r"\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b",
        ".phone": r"\b\d{3}[-.]?\d{3}[-.]?\d{4}\b",
        "ip_address": r"\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b",
        "date_of_birth": r"\b\d{4}-\d{2}-\d{2}\b",
    }
    
    def __init__(
        self,
        config: Optional[AegisGateConfig] = None,
        redact_config: Optional[PIIRedactionConfig] = None,
    ):
        """
        Initialize the PII filter.
        
        Args:
            config: AegisGate configuration
            redact_config: PII redaction configuration
        """
        super().__init__(config)
        self.redact_config = redact_config or PIIRedactionConfig()
        
        # Build patterns based on configuration
        self.patterns: Dict[str, re.Pattern] = {}
        if self.redact_config.redact_ssn:
            self.patterns["ssn"] = re.compile(self.PII_PATTERNS["ssn"])
        if self.redact_config.redact_credit_card:
            self.patterns["credit_card"] = re.compile(self.PII_PATTERNS["credit_card"])
        if self.redact_config.redact_email:
            self.patterns["email"] = re.compile(self.PII_PATTERNS["email"])
        if self.redact_config.redact_phone:
            self.patterns["phone"] = re.compile(self.PII_PATTERNS["phone"])
        if self.redact_config.redact_ip:
            self.patterns["ip_address"] = re.compile(self.PII_PATTERNS["ip_address"])
        if self.redact_config.redact_date_of_birth:
            self.patterns["date_of_birth"] = re.compile(self.PII_PATTERNS["date_of_birth"])
    
    def filter(self, content: str) -> str:
        """
        Filter content for PII and redact if found.
        
        Args:
            content: Content to filter
            
        Returns:
            Redacted content with PII replaced by token
        """
        if not self.redact_config.enabled:
            return content
        
        redacted_content = content
        
        for pii_type, pattern in self.patterns.items():
            matches = pattern.findall(redacted_content)
            if matches:
                logger.debug(f"Found {len(matches)} {pii_type} patterns")
                redacted_content = pattern.sub(
                    self.redact_config.redaction_token,
                    redacted_content,
                )
        
        return redacted_content
    
    def check(self, content: str) -> bool:
        """
        Check if content contains PII.
        
        Args:
            content: Content to check
            
        Returns:
            True if PII is detected
        """
        for pattern in self.patterns.values():
            if pattern.search(content):
                return True
        return False
    
    def detect(self, content: str) -> List[str]:
        """
        Detect PII types in content.
        
        Args:
            content: Content to analyze
            
        Returns:
            List of detected PII types
        """
        detected = []
        for pii_type, pattern in self.patterns.items():
            if pattern.search(content):
                detected.append(pii_type)
        return detected


class ToxicityFilter(BaseSecurityFilter):
    """
    Filter for detecting and blocking toxic or harmful content.
    
    This filter helps prevent the generation of offensive, abusive, or
    inappropriate content through LLM applications.
    """
    
    # Toxic content indicators
    TOXIC_INDICATORS = {
        " hate": [
            r"\b(?:hate|hate\s+speech|discrimination)\b",
            r"\b(?:n\*\*\*|f\*\*t|c\*\*t|r\*\*\*t)\b",
            r"\b(?:racist|sexist|misogynist|homophobic|transphobic)\b",
        ],
        "violence": [
            r"\b(?:kill|murder|assassinate|slaughter)\b",
            r"\b(?:gun|weapon|firearm|explosive)\b",
            r"\b(?:attack|assault| violence)\b",
        ],
        "self_harm": [
            r"\b(?:suicide|self-harm|self\s+ harming)\b",
            r"\b(?:cut|slice|wound)\b",
            r"\b(?:overdose|poison|hang)\b",
        ],
        "sexual": [
            r"\b(?:porn|.xxx)\b",
            r"\b(?:pornography|explicit|nsfw)\b",
            r"\b(?:sexual|sex\s+term)\b",
        ],
    }
    
    def __init__(self, config: Optional[AegisGateConfig] = None):
        """
        Initialize the toxicity filter.
        
        Args:
            config: AegisGate configuration
        """
        super().__init__(config)
        
        # Compile patterns
        self.patterns: Dict[str, List[re.Pattern]] = {}
        for category, patterns in self.TOXIC_INDICATORS.items():
            self.patterns[category] = [
                re.compile(pattern, re.IGNORECASE) for pattern in patterns
            ]
    
    def filter(self, content: str) -> str:
        """
        Filter content for toxicity.
        
        Args:
            content: Content to filter
            
        Returns:
            Original content if safe
            
        Raises:
            ThreatDetectedError: If toxic content is detected
        """
        if not self.config.threat_detection.enabled:
            return content
        
        # Check against AegisGate if available
        if self.config.api_key:
            from aegisgate import AegisGateClient
            try:
                client = AegisGateClient(
                    base_url=self.config.base_url,
                    api_key=self.config.api_key,
                )
                
                result = client.proxy.check_threats(
                    content=content,
                    content_type="input",
                )
                
                if result.get("threat_detected"):
                    if self.config.threat_action == "block":
                        raise ThreatDetectedError(
                            message="Toxic content detected",
                            threat_type="toxicity",
                            threat_score=result.get("threat_score", 0),
                        )
            except Exception as e:
                logger.warning(f"Error checking with AegisGate: {e}")
        
        # Check locally
        for category, patterns in self.patterns.items():
            for pattern in patterns:
                if pattern.search(content):
                    logger.warning(f"Toxic content detected in category: {category}")
                    
                    if self.config.threat_action == "block":
                        raise ThreatDetectedError(
                            message=f"Toxic content detected: {category}",
                            threat_type="toxicity",
                            threat_score=3.0,
                        )
        
        return content
    
    def check(self, content: str) -> bool:
        """
        Check if content contains toxic indicators.
        
        Args:
            content: Content to check
            
        Returns:
            True if toxic content is detected
        """
        for patterns in self.patterns.values():
            for pattern in patterns:
                if pattern.search(content):
                    return True
        return False


class SecretFilter(BaseSecurityFilter):
    """
    Filter for detecting and blocking exposed secrets.
    
    This filter helps prevent the accidental exposure of API keys,
    passwords, and other sensitive credentials.
    """
    
    # Secret detection patterns
    SECRET_PATTERNS = {
        "api_key": [
            r"(?:api[_-]key|apikey)\s*[=:]\s*[\"']?([A-Za-z0-9]{20,})[\"']?",
            r"(?:secret[_-]key|secretkey)\s*[=:]\s*[\"']?([A-Za-z0-9]{20,})[\"']?",
            r"(?:private[_-]key|privatekey)\s*[=:]\s*[\"']?([A-Za-z0-9]{32,})[\"']?",
        ],
        "password": [
            r"(?:password|passwd|pwd)\s*[=:]\s*[\"']?([^\s\"']{8,})[\"']?",
        ],
        "token": [
            r"(?:bearer|token)\s*[=:]\s*[\"']?([A-Za-z0-9._-]{20,})[\"']?",
        ],
        "credit_card": [
            r"\b(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|3[47][0-9]{13})\b",
        ],
    }
    
    def __init__(self, config: Optional[AegisGateConfig] = None):
        """
        Initialize the secret filter.
        
        Args:
            config: AegisGate configuration
        """
        super().__init__(config)
        
        # Compile patterns
        self.patterns: Dict[str, List[re.Pattern]] = {}
        for secret_type, patterns in self.SECRET_PATTERNS.items():
            self.patterns[secret_type] = [
                re.compile(pattern, re.IGNORECASE) for pattern in patterns
            ]
    
    def filter(self, content: str) -> str:
        """
        Filter content for exposed secrets.
        
        Args:
            content: Content to filter
            
        Returns:
            Original content if no secrets found
            
        Raises:
            ThreatDetectedError: If secrets are detected
        """
        if not self.config.threat_detection.enabled:
            return content
        
        # Check against AegisGate if available
        if self.config.api_key:
            from aegisgate import AegisGateClient
            try:
                client = AegisGateClient(
                    base_url=self.config.base_url,
                    api_key=self.config.api_key,
                )
                
                result = client.proxy.check_threats(
                    content=content,
                    content_type="input",
                )
                
                if result.get("threat_detected") and result.get("threat_type") == "secrets":
                    if self.config.threat_action == "block":
                        raise ThreatDetectedError(
                            message="Exposed secrets detected",
                            threat_type="secrets",
                            threat_score=result.get("threat_score", 0),
                        )
            except Exception as e:
                logger.warning(f"Error checking with AegisGate: {e}")
        
        # Check locally
        for secret_type, patterns in self.patterns.items():
            for pattern in patterns:
                if pattern.search(content):
                    logger.warning(f"Potential {secret_type} detected")
                    
                    if self.config.threat_action == "block":
                        raise ThreatDetectedError(
                            message=f"Exposed {secret_type} detected",
                            threat_type="secrets",
                            threat_score=5.0,  # Secrets are critical
                        )
        
        return content
    
    def check(self, content: str) -> bool:
        """
        Check if content contains exposed secrets.
        
        Args:
            content: Content to check
            
        Returns:
            True if secrets are detected
        """
        for patterns in self.patterns.values():
            for pattern in patterns:
                if pattern.search(content):
                    return True
        return False


# Convenience functions
def create_filters(config: Optional[AegisGateConfig] = None) -> Dict[str, BaseSecurityFilter]:
    """
    Create a dictionary of all security filters.
    
    Args:
        config: AegisGate configuration
        
    Returns:
        Dictionary of filter instances
    """
    if config is None:
        config = AegisGateConfig.from_env()
    
    return {
        "prompt_injection": PromptInjectionFilter(config),
        "pii": PIIFilter(config),
        "toxicity": ToxicityFilter(config),
        "secrets": SecretFilter(config),
    }
