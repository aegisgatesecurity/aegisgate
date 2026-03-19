# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Proxy and traffic inspection models for the AegisGate Python SDK.

This module provides data classes for proxy statistics, violations,
and threat detection results.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
from typing import Optional, List, Dict, Any


class ViolationSeverity(str, Enum):
    """
    Severity levels for security violations.
    
    Attributes:
        CRITICAL: Immediate threat requiring blocking
        HIGH: Significant threat that should be blocked
        MEDIUM: Moderate concern
        LOW: Minor anomaly or informational
        INFO: Informational only
    """
    CRITICAL = "critical"
    HIGH = "high"
    MEDIUM = "medium"
    LOW = "low"
    INFO = "info"
    
    @classmethod
    def from_string(cls, value: str) -> ViolationSeverity:
        """Create ViolationSeverity from string value."""
        try:
            return cls(value.lower())
        except ValueError:
            return cls.INFO
    
    @classmethod
    def from_score(cls, score: float) -> ViolationSeverity:
        """Determine severity from numeric score."""
        if score >= 4.0:
            return cls.CRITICAL
        elif score >= 3.0:
            return cls.HIGH
        elif score >= 2.0:
            return cls.MEDIUM
        elif score >= 1.0:
            return cls.LOW
        return cls.INFO


class ThreatType(str, Enum):
    """
    Types of threats detected by AegisGate.
    
    These map to MITRE ATLAS techniques and other security frameworks.
    
    Attributes:
        PROMPT_INJECTION: Direct or indirect prompt injection
        DATA_EXFILTRATION: Data exfiltration attempts
        CREDENTIAL_EXPOSURE: API key or credential exposure
        PII_LEAK: Personal identifiable information leak
        JAILBREAK: LLM jailbreak attempt
        SYSTEM_PROMPT_EXTRACTION: System prompt extraction attempt
        SQL_INJECTION: SQL injection attempt
        XSS: Cross-site scripting attempt
        MODEL_DOS: Denial of service against the model
        POLICY_VIOLATION: Compliance policy violation
        SUSPICIOUS_PATTERN: Suspicious pattern not matching known types
        ANOMALY: Behavioral or statistical anomaly
    """
    # Prompt Injection & Manipulation
    PROMPT_INJECTION = "prompt_injection"
    INDIRECT_PROMPT_INJECTION = "indirect_prompt_injection"
    JAILBREAK = "jailbreak"
    SYSTEM_PROMPT_EXTRACTION = "system_prompt_extraction"
    
    # Data & Credential Threats
    DATA_EXFILTRATION = "data_exfiltration"
    CREDENTIAL_EXPOSURE = "credential_exposure"
    PII_LEAK = "pii_leak"
    TRAINING_DATA_EXPOSURE = "training_data_exposure"
    
    # Injection Attacks
    SQL_INJECTION = "sql_injection"
    XSS = "xss"
    COMMAND_INJECTION = "command_injection"
    CODE_INJECTION = "code_injection"
    
    # DoS & Resource Attacks
    MODEL_DOS = "model_dos"
    RESOURCE_EXHAUSTION = "resource_exhaustion"
    
    # Compliance & Policy
    POLICY_VIOLATION = "policy_violation"
    COMPLIANCE_VIOLATION = "compliance_violation"
    FRAMEWORK_VIOLATION = "framework_violation"
    
    # Anomaly Detection
    ANOMALY = "anomaly"
    SUSPICIOUS_PATTERN = "suspicious_pattern"
    BEHAVIORAL_ANOMALY = "behavioral_anomaly"
    TRAFFIC_ANOMALY = "traffic_anomaly"
    
    # Other
    UNKNOWN = "unknown"
    
    @classmethod
    def from_string(cls, value: str) -> ThreatType:
        """Create ThreatType from string value."""
        try:
            return cls(value.lower())
        except ValueError:
            return cls.UNKNOWN


@dataclass
class Violation:
    """
    Represents a security violation detected by AegisGate.
    
    Attributes:
        id: Unique violation identifier
        timestamp: When the violation occurred
        severity: Violation severity level
        threat_type: Type of threat detected
        threat_category: Broader threat category
        description: Human-readable description
        client_ip: Client IP address
        client_id: Client identifier (API key/user ID)
        path: Request path
        method: HTTP method
        blocked: Whether the request was blocked
        score: Threat score (0.0 - 5.0)
        patterns: List of matched detection patterns
        atlas_techniques: MITRE ATLAS technique IDs
        metadata: Additional metadata
    """
    id: str
    timestamp: datetime
    severity: ViolationSeverity
    threat_type: ThreatType
    description: str
    blocked: bool
    score: float
    threat_category: Optional[str] = None
    client_ip: Optional[str] = None
    client_id: Optional[str] = None
    path: Optional[str] = None
    method: Optional[str] = None
    patterns: List[str] = field(default_factory=list)
    atlas_techniques: List[str] = field(default_factory=list)
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> Violation:
        """Create Violation from dictionary."""
        timestamp = data.get("timestamp")
        if isinstance(timestamp, str):
            try:
                timestamp = datetime.fromisoformat(timestamp.replace("Z", "+00:00"))
            except ValueError:
                timestamp = datetime.now()
        elif timestamp is None:
            timestamp = datetime.now()
        
        severity = ViolationSeverity.from_string(data.get("severity", "info"))
        if "score" in data:
            severity = ViolationSeverity.from_score(data["score"])
        
        return cls(
            id=data.get("id", ""),
            timestamp=timestamp,
            severity=severity,
            threat_type=ThreatType.from_string(data.get("threat_type", "unknown")),
            description=data.get("description", ""),
            blocked=data.get("blocked", False),
            score=data.get("score", 0.0),
            threat_category=data.get("threat_category"),
            client_ip=data.get("client_ip"),
            client_id=data.get("client_id"),
            path=data.get("path"),
            method=data.get("method"),
            patterns=data.get("patterns", []),
            atlas_techniques=data.get("atlas_techniques", []),
            metadata=data.get("metadata", {}),
        )
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert Violation to dictionary."""
        result = {
            "id": self.id,
            "timestamp": self.timestamp.isoformat() if self.timestamp else None,
            "severity": self.severity.value,
            "threat_type": self.threat_type.value,
            "description": self.description,
            "blocked": self.blocked,
            "score": self.score,
        }
        if self.threat_category:
            result["threat_category"] = self.threat_category
        if self.client_ip:
            result["client_ip"] = self.client_ip
        if self.client_id:
            result["client_id"] = self.client_id
        if self.path:
            result["path"] = self.path
        if self.method:
            result["method"] = self.method
        if self.patterns:
            result["patterns"] = self.patterns
        if self.atlas_techniques:
            result["atlas_techniques"] = self.atlas_techniques
        if self.metadata:
            result["metadata"] = self.metadata
        return result


@dataclass
class ViolationFilter:
    """
    Filter criteria for querying violations.
    
    Attributes:
        severity: Filter by severity level
        threat_type: Filter by threat type
        blocked: Filter by blocked status
        start_date: Filter violations after this date
        end_date: Filter violations before this date
        client_ip: Filter by client IP
        limit: Maximum number of results
        offset: Pagination offset
        sort_by: Field to sort by
        sort_order: Sort order (asc/desc)
    """
    severity: Optional[ViolationSeverity] = None
    threat_type: Optional[ThreatType] = None
    blocked: Optional[bool] = None
    start_date: Optional[datetime] = None
    end_date: Optional[datetime] = None
    client_ip: Optional[str] = None
    limit: int = 100
    offset: int = 0
    sort_by: str = "timestamp"
    sort_order: str = "desc"
    
    def to_query_params(self) -> Dict[str, str]:
        """Convert to URL query parameters."""
        params = {}
        if self.severity:
            params["severity"] = self.severity.value
        if self.threat_type:
            params["threat_type"] = self.threat_type.value
        if self.blocked is not None:
            params["blocked"] = "true" if self.blocked else "false"
        if self.start_date:
            params["start_date"] = self.start_date.isoformat()
        if self.end_date:
            params["end_date"] = self.end_date.isoformat()
        if self.client_ip:
            params["client_ip"] = self.client_ip
        if self.limit != 100:
            params["limit"] = str(self.limit)
        if self.offset:
            params["offset"] = str(self.offset)
        if self.sort_by != "timestamp":
            params["sort_by"] = self.sort_by
        if self.sort_order != "desc":
            params["sort_order"] = self.sort_order
        return params


@dataclass
class ProxyStats:
    """
    Proxy statistics and metrics.
    
    Attributes:
        requests_total: Total number of requests processed
        requests_blocked: Requests blocked by AegisGate
        bytes_in: Total bytes received
        bytes_out: Total bytes sent
        avg_response_time: Average response time in milliseconds
        active_connections: Currently active connections
        violations_total: Total violations detected
        violations_by_severity: Violations grouped by severity
        violations_by_type: Violations grouped by type
        requests_by_method: Requests grouped by HTTP method
        requests_by_status: Requests grouped by status code
        top_blocked_paths: Most commonly blocked paths
        top_threat_types: Most common threat types
        uptime_seconds: Server uptime in seconds
        last_updated: Last stats update timestamp
    """
    requests_total: int = 0
    requests_blocked: int = 0
    bytes_in: int = 0
    bytes_out: int = 0
    avg_response_time: float = 0.0
    active_connections: int = 0
    violations_total: int = 0
    violations_by_severity: Dict[str, int] = field(default_factory=dict)
    violations_by_type: Dict[str, int] = field(default_factory=dict)
    requests_by_method: Dict[str, int] = field(default_factory=dict)
    requests_by_status: Dict[str, int] = field(default_factory=dict)
    top_blocked_paths: List[Dict[str, Any]] = field(default_factory=list)
    top_threat_types: List[Dict[str, Any]] = field(default_factory=list)
    uptime_seconds: int = 0
    last_updated: Optional[datetime] = None
    
    @property
    def block_rate(self) -> float:
        """Calculate block rate as percentage."""
        if self.requests_total == 0:
            return 0.0
        return (self.requests_blocked / self.requests_total) * 100
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> ProxyStats:
        """Create ProxyStats from dictionary."""
        last_updated = None
        if data.get("last_updated"):
            try:
                last_updated = datetime.fromisoformat(
                    data["last_updated"].replace("Z", "+00:00")
                )
            except (ValueError, AttributeError):
                pass
        
        return cls(
            requests_total=data.get("requests_total", 0),
            requests_blocked=data.get("requests_blocked", 0),
            bytes_in=data.get("bytes_in", 0),
            bytes_out=data.get("bytes_out", 0),
            avg_response_time=data.get("avg_response_time", 0.0),
            active_connections=data.get("active_connections", 0),
            violations_total=data.get("violations_total", 0),
            violations_by_severity=data.get("violations_by_severity", {}),
            violations_by_type=data.get("violations_by_type", {}),
            requests_by_method=data.get("requests_by_method", {}),
            requests_by_status=data.get("requests_by_status", {}),
            top_blocked_paths=data.get("top_blocked_paths", []),
            top_threat_types=data.get("top_threat_types", []),
            uptime_seconds=data.get("uptime_seconds", 0),
            last_updated=last_updated,
        )


@dataclass
class ProxyConfig:
    """
    Proxy configuration settings.
    
    Attributes:
        enabled: Whether proxy is enabled
        bind_address: Address to bind to
        upstream_url: Upstream server URL
        rate_limit: Requests per minute
        max_body_size: Maximum request body size in bytes
        timeout: Request timeout in seconds
        tls_enabled: Whether TLS is enabled
        tls_cert_file: Path to TLS certificate
        tls_key_file: Path to TLS private key
    """
    enabled: bool = True
    bind_address: str = ":8080"
    upstream_url: str = "http://localhost:3000"
    rate_limit: int = 100
    max_body_size: int = 10 * 1024 * 1024  # 10MB
    timeout: float = 30.0
    tls_enabled: bool = False
    tls_cert_file: Optional[str] = None
    tls_key_file: Optional[str] = None
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> ProxyConfig:
        """Create ProxyConfig from dictionary."""
        return cls(
            enabled=data.get("enabled", True),
            bind_address=data.get("bind_address", ":8080"),
            upstream_url=data.get("upstream_url", "http://localhost:3000"),
            rate_limit=data.get("rate_limit", 100),
            max_body_size=data.get("max_body_size", 10 * 1024 * 1024),
            timeout=data.get("timeout", 30.0),
            tls_enabled=data.get("tls_enabled", False),
            tls_cert_file=data.get("tls_cert_file"),
            tls_key_file=data.get("tls_key_file"),
        )


@dataclass
class ScanResult:
    """
    Result of a content scan.
    
    Attributes:
        is_safe: Whether the content passed all checks
        threats_found: List of detected threats
        scan_time_ms: Time taken to scan in milliseconds
        patterns_matched: Number of patterns matched
        rules_triggered: List of triggered rule names
    """
    is_safe: bool
    threats_found: List[Violation] = field(default_factory=list)
    scan_time_ms: float = 0.0
    patterns_matched: int = 0
    rules_triggered: List[str] = field(default_factory=list)
    
    @property
    def has_critical(self) -> bool:
        """Check if any threat is critical severity."""
        return any(v.severity == ViolationSeverity.CRITICAL for v in self.threats_found)
    
    @property
    def has_high(self) -> bool:
        """Check if any threat is high severity."""
        return any(v.severity == ViolationSeverity.HIGH for v in self.threats_found)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> ScanResult:
        """Create ScanResult from dictionary."""
        threats = []
        if "threats" in data:
            threats = [Violation.from_dict(v) for v in data["threats"]]
        elif "threats_found" in data:
            threats = [Violation.from_dict(v) for v in data["threats_found"]]
        
        return cls(
            is_safe=data.get("is_safe", len(threats) == 0),
            threats_found=threats,
            scan_time_ms=data.get("scan_time_ms", 0.0),
            patterns_matched=data.get("patterns_matched", 0),
            rules_triggered=data.get("rules_triggered", []),
        )
