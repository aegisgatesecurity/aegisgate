# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Webhook models for the AegisGate Python SDK.

This module provides data classes for webhook configuration and delivery.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
from typing import Optional, List, Dict, Any, Callable


class WebhookEventType(str, Enum):
    """
    Webhook event types supported by AegisGate.
    
    Attributes:
        VIOLATION_DETECTED: A security violation was detected
        VIOLATION_BLOCKED: A request was blocked
        COMPLIANCE_CHECK: A compliance check was run
        COMPLIANCE_FAIL: A compliance check failed
        USER_LOGIN: A user logged in
        USER_LOGIN_FAILED: A user login failed
        USER_CREATED: A new user was created
        USER_DELETED: A user was deleted
        CONFIG_CHANGED: Configuration was changed
        SYSTEM_ALERT: A system alert was triggered
        THREAT_DETECTED: A threat was detected
        RATE_LIMIT_EXCEEDED: Rate limit was exceeded
        LICENSE_EXPIRING: License is expiring soon
    """
    VIOLATION_DETECTED = "violation.detected"
    VIOLATION_BLOCKED = "violation.blocked"
    COMPLIANCE_CHECK = "compliance.check"
    COMPLIANCE_FAIL = "compliance.fail"
    USER_LOGIN = "user.login"
    USER_LOGIN_FAILED = "user.login_failed"
    USER_CREATED = "user.created"
    USER_DELETED = "user.deleted"
    CONFIG_CHANGED = "config.changed"
    SYSTEM_ALERT = "system.alert"
    THREAT_DETECTED = "threat.detected"
    RATE_LIMIT_EXCEEDED = "rate_limit.exceeded"
    LICENSE_EXPIRING = "license.expiring"


@dataclass
class WebhookEvent:
    """
    Represents a webhook event payload.
    
    Attributes:
        id: Unique event identifier
        type: Event type
        timestamp: Event timestamp
        data: Event-specific data
        webhook_id: Source webhook ID
    """
    id: str
    type: WebhookEventType
    timestamp: datetime
    webhook_id: str
    data: Dict[str, Any] = field(default_factory=dict)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> WebhookEvent:
        """Create WebhookEvent from dictionary."""
        timestamp = data.get("timestamp")
        if isinstance(timestamp, str):
            try:
                timestamp = datetime.fromisoformat(timestamp.replace("Z", "+00:00"))
            except ValueError:
                timestamp = datetime.now()
        elif timestamp is None:
            timestamp = datetime.now()
        
        event_type = WebhookEventType.INFO
        if "type" in data:
            try:
                event_type = WebhookEventType(data["type"])
            except ValueError:
                event_type = WebhookEventType.INFO
        
        return cls(
            id=data.get("id", ""),
            type=event_type,
            timestamp=timestamp,
            webhook_id=data.get("webhook_id", ""),
            data=data.get("data", {}),
        )


@dataclass
class WebhookDelivery:
    """
    Represents a webhook delivery attempt.
    
    Attributes:
        id: Delivery identifier
        webhook_id: Associated webhook ID
        event_id: Associated event ID
        timestamp: Delivery timestamp
        status: Delivery status (success/failed/pending)
        response_code: HTTP response code
        response_body: Response body (truncated)
        error: Error message if failed
        attempts: Number of delivery attempts
    """
    id: str
    webhook_id: str
    event_id: str
    timestamp: datetime
    status: str
    response_code: Optional[int] = None
    response_body: Optional[str] = None
    error: Optional[str] = None
    attempts: int = 1
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> WebhookDelivery:
        """Create WebhookDelivery from dictionary."""
        timestamp = data.get("timestamp")
        if isinstance(timestamp, str):
            try:
                timestamp = datetime.fromisoformat(timestamp.replace("Z", "+00:00"))
            except ValueError:
                timestamp = datetime.now()
        
        return cls(
            id=data.get("id", ""),
            webhook_id=data.get("webhook_id", ""),
            event_id=data.get("event_id", ""),
            timestamp=timestamp,
            status=data.get("status", "pending"),
            response_code=data.get("response_code"),
            response_body=data.get("response_body"),
            error=data.get("error"),
            attempts=data.get("attempts", 1),
        )


@dataclass
class Webhook:
    """
    Represents a webhook configuration.
    
    Attributes:
        id: Unique webhook identifier
        name: Webhook name
        url: Webhook endpoint URL
        secret: Webhook signing secret
        enabled: Whether webhook is enabled
        events: List of event types to send
        headers: Custom headers to include
        timeout: Request timeout in seconds
        retry_policy: Retry policy configuration
        created_at: Creation timestamp
        updated_at: Last update timestamp
        last_triggered: Last trigger timestamp
        delivery_stats: Delivery statistics
    """
    id: str
    name: str
    url: str
    enabled: bool = True
    events: List[WebhookEventType] = field(default_factory=list)
    secret: Optional[str] = None
    headers: Dict[str, str] = field(default_factory=dict)
    timeout: float = 10.0
    retry_policy: Dict[str, Any] = field(default_factory=dict)
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None
    last_triggered: Optional[datetime] = None
    delivery_stats: Dict[str, int] = field(default_factory=dict)
    
    def __post_init__(self):
        """Convert string events to enum if necessary."""
        if self.events:
            converted = []
            for event in self.events:
                if isinstance(event, str):
                    try:
                        converted.append(WebhookEventType(event))
                    except ValueError:
                        pass
                else:
                    converted.append(event)
            self.events = converted
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> Webhook:
        """Create Webhook from dictionary."""
        def parse_datetime(value: Any) -> Optional[datetime]:
            if value is None:
                return None
            if isinstance(value, datetime):
                return value
            try:
                return datetime.fromisoformat(value.replace("Z", "+00:00"))
            except (ValueError, AttributeError):
                return None
        
        events = []
        if "events" in data:
            for event_str in data["events"]:
                try:
                    events.append(WebhookEventType(event_str))
                except ValueError:
                    pass
        
        return cls(
            id=data.get("id", ""),
            name=data.get("name", ""),
            url=data.get("url", ""),
            secret=data.get("secret"),
            enabled=data.get("enabled", True),
            events=events,
            headers=data.get("headers", {}),
            timeout=data.get("timeout", 10.0),
            retry_policy=data.get("retry_policy", {}),
            created_at=parse_datetime(data.get("created_at")),
            updated_at=parse_datetime(data.get("updated_at")),
            last_triggered=parse_datetime(data.get("last_triggered")),
            delivery_stats=data.get("delivery_stats", {}),
        )
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert Webhook to dictionary."""
        result = {
            "name": self.name,
            "url": self.url,
            "enabled": self.enabled,
            "events": [e.value for e in self.events],
            "headers": self.headers,
            "timeout": self.timeout,
            "retry_policy": self.retry_policy,
        }
        if self.secret:
            result["secret"] = self.secret
        return result
    
    @property
    def event_types(self) -> List[str]:
        """Return event types as strings."""
        return [e.value for e in self.events]


@dataclass
class WebhookCreate:
    """
    Data for creating a new webhook.
    
    Attributes:
        name: Webhook name (required)
        url: Webhook endpoint URL (required)
        events: List of event types to subscribe to (required)
        secret: Optional signing secret
        headers: Custom headers
        timeout: Request timeout
        retry_policy: Retry configuration
    """
    name: str
    url: str
    events: List[WebhookEventType]
    secret: Optional[str] = None
    headers: Dict[str, str] = field(default_factory=dict)
    timeout: float = 10.0
    retry_policy: Optional[Dict[str, Any]] = None
    
    def __post_init__(self):
        """Validate the webhook creation data."""
        if not self.name:
            raise ValueError("Webhook name is required")
        if not self.url:
            raise ValueError("Webhook URL is required")
        if not self.events:
            raise ValueError("At least one event type is required")
        
        # Validate URL format
        if not self.url.startswith(("http://", "https://")):
            raise ValueError("Webhook URL must start with http:// or https://")
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary."""
        return {
            "name": self.name,
            "url": self.url,
            "events": [e.value for e in self.events],
            "headers": self.headers,
            "timeout": self.timeout,
            "retry_policy": self.retry_policy or {
                "max_attempts": 3,
                "initial_delay": 1.0,
                "max_delay": 60.0,
            },
        }


@dataclass
class WebhookUpdate:
    """
    Data for updating an existing webhook.
    
    All fields are optional - only provided fields will be updated.
    """
    name: Optional[str] = None
    url: Optional[str] = None
    enabled: Optional[bool] = None
    events: Optional[List[WebhookEventType]] = None
    secret: Optional[str] = None
    headers: Optional[Dict[str, str]] = None
    timeout: Optional[float] = None
    retry_policy: Optional[Dict[str, Any]] = None
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary."""
        result = {}
        if self.name is not None:
            result["name"] = self.name
        if self.url is not None:
            result["url"] = self.url
        if self.enabled is not None:
            result["enabled"] = self.enabled
        if self.events is not None:
            result["events"] = [e.value for e in self.events]
        if self.secret is not None:
            result["secret"] = self.secret
        if self.headers is not None:
            result["headers"] = self.headers
        if self.timeout is not None:
            result["timeout"] = self.timeout
        if self.retry_policy is not None:
            result["retry_policy"] = self.retry_policy
        return result
