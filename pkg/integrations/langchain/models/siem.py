# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
SIEM integration models for the AegisGate Python SDK.

This module provides data classes for SIEM integration, event logging,
and security information management.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
from typing import Optional, List, Dict, Any


class SIEMProvider(str, Enum):
    """
    Supported SIEM providers.
    
    Attributes:
        SPLUNK: Splunk Enterprise Security
        ELASTIC: Elasticsearch/Elastic Security
        SUMOLOGIC: Sumo Logic
        DATADOG: Datadog Security
        PAGERDUTY: PagerDuty alerting
        MICROSOFT_SENTINEL: Microsoft Sentinel
        CHRONICLE: Google Chronicle
        CUSTOM: Custom webhook endpoint
    """
    SPLUNK = "splunk"
    ELASTIC = "elastic"
    SUMOLOGIC = "sumologic"
    DATADOG = "datadog"
    PAGERDUTY = "pagerduty"
    MICROSOFT_SENTINEL = "microsoft_sentinel"
    CHRONICLE = "chronicle"
    CUSTOM = "custom"
    
    @classmethod
    def from_string(cls, value: str) -> SIEMProvider:
        """Create SIEMProvider from string value."""
        try:
            return cls(value.lower())
        except ValueError:
            return cls.CUSTOM


class EventSeverity(str, Enum):
    """
    Event severity levels for SIEM integration.
    
    Attributes:
        EMERGENCY: System is unusable
        ALERT: Immediate action required
        CRITICAL: Critical conditions
        ERROR: Error conditions
        WARNING: Warning conditions
        NOTICE: Normal but significant
        INFO: Informational
        DEBUG: Debug-level messages
    """
    EMERGENCY = "emergency"
    ALERT = "alert"
    CRITICAL = "critical"
    ERROR = "error"
    WARNING = "warning"
    NOTICE = "notice"
    INFO = "info"
    DEBUG = "debug"
    
    @classmethod
    def from_string(cls, value: str) -> EventSeverity:
        """Create EventSeverity from string value."""
        try:
            return cls(value.lower())
        except ValueError:
            return cls.INFO
    
    def to_syslog(self) -> int:
        """Convert to Syslog priority value."""
        mapping = {
            self.EMERGENCY: 0,
            self.ALERT: 1,
            self.CRITICAL: 2,
            self.ERROR: 3,
            self.WARNING: 4,
            self.NOTICE: 5,
            self.INFO: 6,
            self.DEBUG: 7,
        }
        return mapping.get(self, 6)


class EventCategory(str, Enum):
    """
    Event categories for SIEM classification.
    
    Attributes:
        AUTHENTICATION: Authentication events
        AUTHORIZATION: Authorization/access control events
        CONFIGURATION: Configuration changes
        THREAT: Threat detection events
        VIOLATION: Policy violation events
        NETWORK: Network-related events
        SYSTEM: System events
        AUDIT: Audit events
        COMPLIANCE: Compliance-related events
    """
    AUTHENTICATION = "authentication"
    AUTHORIZATION = "authorization"
    CONFIGURATION = "configuration"
    THREAT = "threat"
    VIOLATION = "violation"
    NETWORK = "network"
    SYSTEM = "system"
    AUDIT = "audit"
    COMPLIANCE = "compliance"
    
    @classmethod
    def from_string(cls, value: str) -> EventCategory:
        """Create EventCategory from string value."""
        try:
            return cls(value.lower())
        except ValueError:
            return cls.SYSTEM


@dataclass
class SIEMEvent:
    """
    Represents a SIEM event to be sent to external SIEM systems.
    
    Attributes:
        timestamp: Event timestamp
        severity: Event severity level
        category: Event category
        title: Event title/summary
        message: Detailed event message
        source_ip: Source IP address
        destination_ip: Destination IP address
        user_id: User identifier
        session_id: Session identifier
        threat_type: Type of threat (if applicable)
        threat_score: Threat score (0-100)
        atlas_techniques: MITRE ATLAS technique IDs
        metadata: Additional event metadata
    """
    timestamp: datetime
    severity: EventSeverity
    category: EventCategory
    title: str
    message: str
    source_ip: Optional[str] = None
    destination_ip: Optional[str] = None
    user_id: Optional[str] = None
    session_id: Optional[str] = None
    threat_type: Optional[str] = None
    threat_score: Optional[float] = None
    atlas_techniques: List[str] = field(default_factory=list)
    metadata: Dict[str, Any] = field(default_factory=dict)
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary for SIEM submission."""
        result = {
            "timestamp": self.timestamp.isoformat(),
            "severity": self.severity.value,
            "category": self.category.value,
            "title": self.title,
            "message": self.message,
        }
        if self.source_ip:
            result["source_ip"] = self.source_ip
        if self.destination_ip:
            result["destination_ip"] = self.destination_ip
        if self.user_id:
            result["user_id"] = self.user_id
        if self.session_id:
            result["session_id"] = self.session_id
        if self.threat_type:
            result["threat_type"] = self.threat_type
        if self.threat_score is not None:
            result["threat_score"] = self.threat_score
        if self.atlas_techniques:
            result["atlas_techniques"] = self.atlas_techniques
        if self.metadata:
            result["metadata"] = self.metadata
        return result
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> SIEMEvent:
        """Create SIEMEvent from dictionary."""
        timestamp = data.get("timestamp")
        if isinstance(timestamp, str):
            try:
                timestamp = datetime.fromisoformat(timestamp.replace("Z", "+00:00"))
            except ValueError:
                timestamp = datetime.now()
        elif timestamp is None:
            timestamp = datetime.now()
        
        return cls(
            timestamp=timestamp,
            severity=EventSeverity.from_string(data.get("severity", "info")),
            category=EventCategory.from_string(data.get("category", "system")),
            title=data.get("title", ""),
            message=data.get("message", ""),
            source_ip=data.get("source_ip"),
            destination_ip=data.get("destination_ip"),
            user_id=data.get("user_id"),
            session_id=data.get("session_id"),
            threat_type=data.get("threat_type"),
            threat_score=data.get("threat_score"),
            atlas_techniques=data.get("atlas_techniques", []),
            metadata=data.get("metadata", {}),
        )
    
    @classmethod
    def from_violation(
        cls,
        violation: "Violation",
        category: EventCategory = EventCategory.THREAT,
    ) -> SIEMEvent:
        """
        Create a SIEMEvent from a Violation.
        
        Args:
            violation: The violation to convert
            category: Event category (default: threat)
        
        Returns:
            SIEMEvent instance
        """
        return cls(
            timestamp=violation.timestamp,
            severity=EventSeverity.from_string(violation.severity.value),
            category=category,
            title=f"Threat Detected: {violation.threat_type.value}",
            message=violation.description,
            source_ip=violation.client_ip,
            user_id=violation.client_id,
            threat_type=violation.threat_type.value,
            threat_score=violation.score * 20,  # Convert 0-5 to 0-100
            atlas_techniques=violation.atlas_techniques,
            metadata={
                "violation_id": violation.id,
                "blocked": violation.blocked,
                "patterns": violation.patterns,
            },
        )


@dataclass
class SIEMConfig:
    """
    SIEM integration configuration.
    
    Attributes:
        enabled: Whether SIEM integration is enabled
        provider: SIEM provider type
        endpoint: SIEM endpoint URL
        api_key: API key for SIEM (if required)
        batch_size: Number of events to batch before sending
        flush_interval: Seconds between batch flushes
        include_debug: Include debug-level events
        filters: Event filters
    """
    enabled: bool = False
    provider: SIEMProvider = SIEMProvider.CUSTOM
    endpoint: str = ""
    api_key: Optional[str] = None
    batch_size: int = 100
    flush_interval: int = 60
    include_debug: bool = False
    filters: Dict[str, Any] = field(default_factory=dict)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> SIEMConfig:
        """Create SIEMConfig from dictionary."""
        return cls(
            enabled=data.get("enabled", False),
            provider=SIEMProvider.from_string(data.get("provider", "custom")),
            endpoint=data.get("endpoint", ""),
            api_key=data.get("api_key"),
            batch_size=data.get("batch_size", 100),
            flush_interval=data.get("flush_interval", 60),
            include_debug=data.get("include_debug", False),
            filters=data.get("filters", {}),
        )
    
    def validate(self) -> List[str]:
        """Validate configuration and return list of errors."""
        errors = []
        if self.enabled and not self.endpoint:
            errors.append("SIEM endpoint is required when SIEM is enabled")
        if self.batch_size <= 0:
            errors.append("batch_size must be positive")
        if self.flush_interval <= 0:
            errors.append("flush_interval must be positive")
        return errors


@dataclass
class SIEMStats:
    """
    SIEM integration statistics.
    
    Attributes:
        events_sent: Total events sent
        events_failed: Events that failed to send
        bytes_sent: Total bytes sent
        last_sent: Timestamp of last successful send
        last_failed: Timestamp of last failed send
        queue_size: Current event queue size
        providers: Stats per provider
    """
    events_sent: int = 0
    events_failed: int = 0
    bytes_sent: int = 0
    last_sent: Optional[datetime] = None
    last_failed: Optional[datetime] = None
    queue_size: int = 0
    providers: Dict[str, Dict[str, Any]] = field(default_factory=dict)
    
    @property
    def success_rate(self) -> float:
        """Calculate success rate percentage."""
        total = self.events_sent + self.events_failed
        if total == 0:
            return 100.0
        return (self.events_sent / total) * 100
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> SIEMStats:
        """Create SIEMStats from dictionary."""
        def parse_datetime(value: Any) -> Optional[datetime]:
            if value is None:
                return None
            if isinstance(value, datetime):
                return value
            try:
                return datetime.fromisoformat(value.replace("Z", "+00:00"))
            except (ValueError, AttributeError):
                return None
        
        return cls(
            events_sent=data.get("events_sent", 0),
            events_failed=data.get("events_failed", 0),
            bytes_sent=data.get("bytes_sent", 0),
            last_sent=parse_datetime(data.get("last_sent")),
            last_failed=parse_datetime(data.get("last_failed")),
            queue_size=data.get("queue_size", 0),
            providers=data.get("providers", {}),
        )
