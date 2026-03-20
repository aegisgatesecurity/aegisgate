# SPDX-License-Identifier: MIT
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.

"""
Data models for the AegisGate Python SDK.
"""

from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
from typing import Any, Dict, List, Optional


class ViolationSeverity(str, Enum):
    """Security violation severity levels."""
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"
    CRITICAL = "critical"


class ViolationType(str, Enum):
    """Types of security violations."""
    PROMPT_INJECTION = "prompt_injection"
    DATA_EXFILTRATION = "data_exfiltration"
    MODEL_MANIPULATION = "model_manipulation"
    PII_EXPOSURE = "pii_exposure"
    TOXIC_CONTENT = "toxic_content"
    ADVERSARIAL_ATTACK = "adversarial_attack"
    COMPLIANCE_VIOLATION = "compliance_violation"


class LicenseType(str, Enum):
    """License tier types."""
    COMMUNITY = "community"
    DEVELOPER = "developer"
    PROFESSIONAL = "professional"
    ENTERPRISE = "enterprise"
    CUSTOM = "custom"


class SIEMProvider(str, Enum):
    """Supported SIEM providers."""
    SPLUNK = "splunk"
    ELASTICSEARCH = "elasticsearch"
    SENTINEL = "azure_sentinel"
    QRADAR = "qradar"
    CHRONICLE = "chronicle"
    CUSTOM = "custom"


@dataclass
class Health:
    """Health check response."""
    status: str
    version: str
    uptime_seconds: int
    components: Dict[str, str] = field(default_factory=dict)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Health":
        return cls(
            status=data.get("status", "unknown"),
            version=data.get("version", ""),
            uptime_seconds=data.get("uptime_seconds", 0),
            components=data.get("components", {}),
        )


@dataclass
class Version:
    """Version information."""
    version: str
    go_version: str
    platform: str
    build_date: Optional[datetime] = None
    git_commit: Optional[str] = None
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Version":
        build_date = None
        if data.get("build_date"):
            build_date = datetime.fromisoformat(data["build_date"])
        return cls(
            version=data.get("version", ""),
            go_version=data.get("go_version", ""),
            platform=data.get("platform", ""),
            build_date=build_date,
            git_commit=data.get("git_commit"),
        )


@dataclass
class Module:
    """Module information."""
    name: str
    version: str
    enabled: bool
    tier: int
    description: Optional[str] = None
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Module":
        return cls(
            name=data.get("name", ""),
            version=data.get("version", ""),
            enabled=data.get("enabled", False),
            tier=data.get("tier", 0),
            description=data.get("description"),
        )


@dataclass
class License:
    """License information."""
    type: LicenseType
    customer_id: str
    expires_at: Optional[datetime] = None
    features: List[str] = field(default_factory=list)
    tier: str = "community"
    valid: bool = False
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "License":
        expires_at = None
        if data.get("expires_at"):
            expires_at = datetime.fromisoformat(data["expires_at"])
        return cls(
            type=LicenseType(data.get("type", "community")),
            customer_id=data.get("customer_id", ""),
            expires_at=expires_at,
            features=data.get("features", []),
            tier=data.get("tier", "community"),
            valid=data.get("valid", False),
        )


@dataclass
class Violation:
    """Security violation details."""
    id: str
    type: ViolationType
    severity: ViolationSeverity
    message: str
    timestamp: datetime
    source_ip: Optional[str] = None
    user_id: Optional[str] = None
    request_id: Optional[str] = None
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Violation":
        timestamp = datetime.fromisoformat(data.get("timestamp", datetime.utcnow().isoformat()))
        return cls(
            id=data.get("id", ""),
            type=ViolationType(data.get("type", "prompt_injection")),
            severity=ViolationSeverity(data.get("severity", "medium")),
            message=data.get("message", ""),
            timestamp=timestamp,
            source_ip=data.get("source_ip"),
            user_id=data.get("user_id"),
            request_id=data.get("request_id"),
            metadata=data.get("metadata", {}),
        )


@dataclass
class DetectionResult:
    """Security detection result."""
    detected: bool
    confidence: float
    violations: List[Violation] = field(default_factory=list)
    processing_time_ms: float = 0.0
    detector_type: str = "unknown"
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "DetectionResult":
        violations = [
            Violation.from_dict(v) for v in data.get("violations", [])
        ]
        return cls(
            detected=data.get("detected", False),
            confidence=data.get("confidence", 0.0),
            violations=violations,
            processing_time_ms=data.get("processing_time_ms", 0.0),
            detector_type=data.get("detector_type", "unknown"),
        )


@dataclass
class AnomalyResult:
    """ML anomaly detection result."""
    is_anomaly: bool
    anomaly_score: float
    baseline_score: float
    deviation: float
    features: Dict[str, float] = field(default_factory=dict)
    recommendation: Optional[str] = None
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "AnomalyResult":
        return cls(
            is_anomaly=data.get("is_anomaly", False),
            anomaly_score=data.get("anomaly_score", 0.0),
            baseline_score=data.get("baseline_score", 0.0),
            deviation=data.get("deviation", 0.0),
            features=data.get("features", {}),
            recommendation=data.get("recommendation"),
        )


@dataclass
class ComplianceControl:
    """Compliance control definition."""
    id: str
    name: str
    description: str
    category: str
    severity: str
    automated: bool
    framework: str
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "ComplianceControl":
        return cls(
            id=data.get("id", ""),
            name=data.get("name", ""),
            description=data.get("description", ""),
            category=data.get("category", ""),
            severity=data.get("severity", "medium"),
            automated=data.get("automated", False),
            framework=data.get("framework", ""),
        )


@dataclass
class ComplianceResult:
    """Compliance check result."""
    framework: str
    compliant: bool
    score: float
    controls_passed: int
    controls_failed: int
    controls_total: int
    failures: List[Dict[str, Any]] = field(default_factory=list)
    checked_at: datetime = field(default_factory=datetime.utcnow)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "ComplianceResult":
        checked_at = datetime.utcnow()
        if data.get("checked_at"):
            checked_at = datetime.fromisoformat(data["checked_at"])
        return cls(
            framework=data.get("framework", ""),
            compliant=data.get("compliant", False),
            score=data.get("score", 0.0),
            controls_passed=data.get("controls_passed", 0),
            controls_failed=data.get("controls_failed", 0),
            controls_total=data.get("controls_total", 0),
            failures=data.get("failures", []),
            checked_at=checked_at,
        )


@dataclass
class SIEMConfig:
    """SIEM integration configuration."""
    provider: SIEMProvider
    endpoint: str
    api_key: Optional[str] = None
    enabled: bool = True
    batch_size: int = 100
    flush_interval_seconds: int = 30
    custom_headers: Dict[str, str] = field(default_factory=dict)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "provider": self.provider.value,
            "endpoint": self.endpoint,
            "api_key": self.api_key,
            "enabled": self.enabled,
            "batch_size": self.batch_size,
            "flush_interval_seconds": self.flush_interval_seconds,
            "custom_headers": self.custom_headers,
        }
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "SIEMConfig":
        return cls(
            provider=SIEMProvider(data.get("provider", "custom")),
            endpoint=data.get("endpoint", ""),
            api_key=data.get("api_key"),
            enabled=data.get("enabled", True),
            batch_size=data.get("batch_size", 100),
            flush_interval_seconds=data.get("flush_interval_seconds", 30),
            custom_headers=data.get("custom_headers", {}),
        )


@dataclass
class SIEMEvent:
    """SIEM event for logging."""
    event_type: str
    timestamp: datetime
    severity: str
    source: str
    message: str
    details: Dict[str, Any] = field(default_factory=dict)
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "event_type": self.event_type,
            "timestamp": self.timestamp.isoformat(),
            "severity": self.severity,
            "source": self.source,
            "message": self.message,
            "details": self.details,
        }


@dataclass
class Webhook:
    """Webhook configuration."""
    id: str
    url: str
    events: List[str] = field(default_factory=list)
    secret: Optional[str] = None
    enabled: bool = True
    created_at: Optional[datetime] = None
    last_triggered: Optional[datetime] = None
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Webhook":
        created_at = None
        if data.get("created_at"):
            created_at = datetime.fromisoformat(data["created_at"])
        last_triggered = None
        if data.get("last_triggered"):
            last_triggered = datetime.fromisoformat(data["last_triggered"])
        return cls(
            id=data.get("id", ""),
            url=data.get("url", ""),
            events=data.get("events", []),
            secret=data.get("secret"),
            enabled=data.get("enabled", True),
            created_at=created_at,
            last_triggered=last_triggered,
        )
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "id": self.id,
            "url": self.url,
            "events": self.events,
            "secret": self.secret,
            "enabled": self.enabled,
        }


@dataclass
class User:
    """User information."""
    id: str
    email: str
    name: Optional[str] = None
    roles: List[str] = field(default_factory=list)
    created_at: Optional[datetime] = None
    last_login: Optional[datetime] = None
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "User":
        created_at = None
        if data.get("created_at"):
            created_at = datetime.fromisoformat(data["created_at"])
        last_login = None
        if data.get("last_login"):
            last_login = datetime.fromisoformat(data["last_login"])
        return cls(
            id=data.get("id", ""),
            email=data.get("email", ""),
            name=data.get("name"),
            roles=data.get("roles", []),
            created_at=created_at,
            last_login=last_login,
        )


@dataclass
class ProxyStats:
    """Proxy statistics."""
    requests_total: int
    requests_blocked: int
    requests_allowed: int
    avg_latency_ms: float
    violations_detected: int
    uptime_seconds: int
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "ProxyStats":
        return cls(
            requests_total=data.get("requests_total", 0),
            requests_blocked=data.get("requests_blocked", 0),
            requests_allowed=data.get("requests_allowed", 0),
            avg_latency_ms=data.get("avg_latency_ms", 0.0),
            violations_detected=data.get("violations_detected", 0),
            uptime_seconds=data.get("uptime_seconds", 0),
        )


@dataclass
class ATLASThreat:
    """MITRE ATLAS threat definition."""
    id: str
    name: str
    description: str
    tactics: List[str] = field(default_factory=list)
    platforms: List[str] = field(default_factory=list)
    detection_patterns: List[str] = field(default_factory=list)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "ATLASThreat":
        return cls(
            id=data.get("id", ""),
            name=data.get("name", ""),
            description=data.get("description", ""),
            tactics=data.get("tactics", []),
            platforms=data.get("platforms", []),
            detection_patterns=data.get("detection_patterns", []),
        )


__all__ = [
    "ViolationSeverity",
    "ViolationType",
    "LicenseType",
    "SIEMProvider",
    "Health",
    "Version",
    "Module",
    "License",
    "Violation",
    "DetectionResult",
    "AnomalyResult",
    "ComplianceControl",
    "ComplianceResult",
    "SIEMConfig",
    "SIEMEvent",
    "Webhook",
    "User",
    "ProxyStats",
    "ATLASThreat",
]