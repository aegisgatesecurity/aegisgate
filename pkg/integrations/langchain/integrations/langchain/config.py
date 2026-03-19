# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Configuration for LangChain integration.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Optional, List, Dict, Any, Set
from enum import Enum


class ThreatAction(str, Enum):
    """
    Action to take when a threat is detected.
    
    Attributes:
        BLOCK: Block the request and raise an exception
        ALERT: Allow the request but log an alert
        REDACT: Redact sensitive content and continue
        SANITIZE: Sanitize the input and continue
    """
    BLOCK = "block"
    ALERT = "alert"
    REDACT = "redact"
    SANITIZE = "sanitize"


class ThreatType(str, Enum):
    """
    Types of threats to detect.
    
    Attributes:
        PROMPT_INJECTION: Direct or indirect prompt injection
        PII: Personal Identifiable Information
        SECRETS: API keys, passwords, etc.
        JAILBREAK: LLM jailbreak attempts
        TOXICITY: Toxic or harmful content
        SENSITIVE_DATA: Sensitive business data
    """
    PROMPT_INJECTION = "prompt_injection"
    PII = "pii"
    SECRETS = "secrets"
    JAILBREAK = "jailbreak"
    TOXICITY = "toxicity"
    SENSITIVE_DATA = "sensitive_data"


@dataclass
class ThreatDetectionConfig:
    """
    Configuration for threat detection.
    
    Attributes:
        enabled: Whether threat detection is enabled
        block_on_critical: Block requests with critical threats
        block_on_high: Block requests with high-severity threats
        min_score_to_block: Minimum threat score to trigger blocking (0-5)
        enabled_threats: List of threat types to detect
        custom_patterns: Custom regex patterns to detect
    """
    enabled: bool = True
    block_on_critical: bool = True
    block_on_high: bool = False
    min_score_to_block: float = 3.0
    enabled_threats: Set[ThreatType] = field(
        default_factory=lambda: {
            ThreatType.PROMPT_INJECTION,
            ThreatType.PII,
            ThreatType.SECRETS,
            ThreatType.JAILBREAK,
        }
    )
    custom_patterns: Dict[str, str] = field(default_factory=dict)
    
    def is_enabled_for(self, threat_type: ThreatType) -> bool:
        """Check if a threat type is enabled for detection."""
        return threat_type in self.enabled_threats


@dataclass
class ComplianceConfig:
    """
    Configuration for compliance monitoring.
    
    Attributes:
        enabled: Whether compliance monitoring is enabled
        frameworks: List of compliance frameworks to use
        report_violations: Whether to generate compliance reports
        audit_trail: Whether to maintain full audit trail
    """
    enabled: bool = True
    frameworks: List[str] = field(
        default_factory=lambda: ["mitre_atlas", "nist_ai_rmf"]
    )
    report_violations: bool = True
    audit_trail: bool = True


@dataclass
class RateLimitConfig:
    """
    Configuration for rate limiting.
    
    Attributes:
        enabled: Whether rate limiting is enabled
        requests_per_minute: Maximum requests per minute
        requests_per_hour: Maximum requests per hour
        requests_per_day: Maximum requests per day
        burst_size: Maximum burst size
        block_on_exceed: Whether to block on exceeding limits
    """
    enabled: bool = True
    requests_per_minute: int = 60
    requests_per_hour: int = 1000
    requests_per_day: int = 10000
    burst_size: int = 10
    block_on_exceed: bool = False


@dataclass
class PIIRedactionConfig:
    """
    Configuration for PII redaction.
    
    Attributes:
        enabled: Whether PII redaction is enabled
        redact_ssn: Redact Social Security Numbers
        redact_credit_card: Redact Credit Card Numbers
        redact_email: Redact Email Addresses
        redact_phone: Redact Phone Numbers
        redact_ip: Redact IP Addresses
        redact_date_of_birth: Redact Dates of Birth
        redaction_token: Token to use for redaction
    """
    enabled: bool = True
    redact_ssn: bool = True
    redact_credit_card: bool = True
    redact_email: bool = False
    redact_phone: bool = False
    redact_ip: bool = True
    redact_date_of_birth: bool = False
    redaction_token: str = "[REDACTED]"


@dataclass
class LoggingConfig:
    """
    Configuration for logging.
    
    Attributes:
        enabled: Whether logging is enabled
        log_inputs: Log LLM inputs
        log_outputs: Log LLM outputs
        log_violations: Log security violations
        log_compliance: Log compliance events
        verbose: Enable verbose logging
    """
    enabled: bool = True
    log_inputs: bool = True
    log_outputs: bool = True
    log_violations: bool = True
    log_compliance: bool = True
    verbose: bool = False


@dataclass
class CacheConfig:
    """
    Configuration for caching.
    
    Attributes:
        enabled: Whether caching is enabled
        ttl_seconds: Cache TTL in seconds
        max_size: Maximum cache size
        cache_responses: Cache LLM responses
        cache_violations: Cache violation results
    """
    enabled: bool = False
    ttl_seconds: int = 300
    max_size: int = 1000
    cache_responses: bool = True
    cache_violations: bool = True


@dataclass
class AegisGateConfig:
    """
    Main configuration for AegisGate LangChain integration.
    
    This class provides comprehensive configuration for securing LangChain
    LLM applications with AegisGate.
    
    Attributes:
        base_url: AegisGate server URL
        api_key: AegisGate API key
        timeout: Request timeout in seconds
        
        # Threat Detection
        threat_detection: Threat detection configuration
        threat_action: Default action for detected threats
        
        # Compliance
        compliance: Compliance monitoring configuration
        
        # Rate Limiting
        rate_limit: Rate limiting configuration
        
        # PII Handling
        pii_redaction: PII redaction configuration
        
        # Logging
        logging: Logging configuration
        
        # Caching
        cache: Cache configuration
        
    Example:
        Basic configuration:
        >>> config = AegisGateConfig(
        ...     base_url="http://localhost:8080",
        ...     api_key="AG-your-key"
        ... )
        
        Strict security:
        >>> config = AegisGateConfig(
        ...     base_url="https://aegisgate.example.com",
        ...     api_key="AG-your-key",
        ...     threat_detection=ThreatDetectionConfig(
        ...         block_on_critical=True,
        ...         block_on_high=True,
        ...     ),
        ...     rate_limit=RateLimitConfig(
        ...         requests_per_minute=10,
        ...         block_on_exceed=True,
        ...     )
        ... )
        
        Compliance-focused:
        >>> config = AegisGateConfig(
        ...     base_url="http://localhost:8080",
        ...     compliance=ComplianceConfig(
        ...         frameworks=["hipaa", "pci_dss", "gdpr"],
        ...         audit_trail=True,
        ...     )
        ... )
    """
    # Connection settings
    base_url: str = "http://localhost:8080"
    api_key: Optional[str] = None
    timeout: float = 30.0
    verify_ssl: bool = True
    proxy_url: Optional[str] = None
    
    # Sub-configurations
    threat_detection: ThreatDetectionConfig = field(
        default_factory=ThreatDetectionConfig
    )
    compliance: ComplianceConfig = field(
        default_factory=ComplianceConfig
    )
    rate_limit: RateLimitConfig = field(
        default_factory=RateLimitConfig
    )
    pii_redaction: PIIRedactionConfig = field(
        default_factory=PIIRedactionConfig
    )
    logging: LoggingConfig = field(
        default_factory=LoggingConfig
    )
    cache: CacheConfig = field(
        default_factory=CacheConfig
    )
    
    # Default action for threats
    threat_action: ThreatAction = ThreatAction.BLOCK
    
    # Legacy support
    block_on_critical: bool = True
    block_on_high: bool = False
    
    def __post_init__(self):
        """Apply legacy settings if provided."""
        # Apply legacy settings to threat detection config
        if not self.threat_detection:
            self.threat_detection = ThreatDetectionConfig()
        
        self.threat_detection.block_on_critical = self.block_on_critical
        self.threat_detection.block_on_high = self.block_on_high
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert configuration to dictionary."""
        return {
            "base_url": self.base_url,
            "api_key": self.api_key,
            "timeout": self.timeout,
            "verify_ssl": self.verify_ssl,
            "threat_detection": {
                "enabled": self.threat_detection.enabled,
                "block_on_critical": self.threat_detection.block_on_critical,
                "block_on_high": self.threat_detection.block_on_high,
                "min_score_to_block": self.threat_detection.min_score_to_block,
                "enabled_threats": [t.value for t in self.threat_detection.enabled_threats],
            },
            "compliance": {
                "enabled": self.compliance.enabled,
                "frameworks": self.compliance.frameworks,
            },
            "rate_limit": {
                "enabled": self.rate_limit.enabled,
                "requests_per_minute": self.rate_limit.requests_per_minute,
            },
            "threat_action": self.threat_action.value,
        }
    
    @classmethod
    def from_env(cls) -> "AegisGateConfig":
        """Create configuration from environment variables."""
        import os
        
        return cls(
            base_url=os.environ.get("AEGISGATE_BASE_URL", "http://localhost:8080"),
            api_key=os.environ.get("AEGISGATE_API_KEY"),
            timeout=float(os.environ.get("AEGISGATE_TIMEOUT", "30")),
        )
