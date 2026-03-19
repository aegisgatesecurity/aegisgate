# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
SIEM service for the AegisGate Python SDK.

This module provides SIEM integration and event logging functionality.
"""

from __future__ import annotations

from typing import Optional, List, Dict, Any
from datetime import datetime
import logging

from aegisgate.models.siem import (
    SIEMConfig,
    SIEMEvent,
    SIEMProvider,
    EventSeverity,
    EventCategory,
    SIEMStats,
)
from aegisgate.models.proxy import Violation

logger = logging.getLogger(__name__)


class SIEMService:
    """
    Service for SIEM integration management.
    
    This service handles:
    - SIEM configuration
    - Event submission
    - Event querying
    - Stats retrieval
    
    Example:
        >>> client = AegisGateClient(base_url="http://localhost:8080")
        >>> siem = client.siem
        >>> 
        >>> # Get SIEM config
        >>> config = siem.get_config()
        >>> print(f"Provider: {config.provider.value}")
        >>> 
        >>> # Send an event
        >>> event = SIEMEvent(
        ...     timestamp=datetime.now(),
        ...     severity=EventSeverity.CRITICAL,
        ...     category=EventCategory.THREAT,
        ...     title="Threat Detected",
        ...     message="Prompt injection attempt blocked"
        ... )
        >>> siem.send_event(event)
    """
    
    def __init__(self, client: "AegisGateClient"):
        """
        Initialize the SIEM service.
        
        Args:
            client: The AegisGate client instance
        """
        self._client = client
    
    # =====================================================================
    # Configuration
    # =====================================================================
    
    def get_config(self) -> SIEMConfig:
        """
        Get current SIEM configuration.
        
        Returns:
            SIEMConfig object
        
        Example:
            >>> config = siem.get_config()
            >>> print(f"Enabled: {config.enabled}")
        """
        response = self._client._request("GET", "siem/config")
        return SIEMConfig.from_dict(response)
    
    def update_config(self, config: SIEMConfig) -> SIEMConfig:
        """
        Update SIEM configuration.
        
        Args:
            config: New SIEM configuration
        
        Returns:
            Updated SIEMConfig object
        
        Example:
            >>> config = SIEMConfig(
            ...     enabled=True,
            ...     provider=SIEMProvider.SPLUNK,
            ...     endpoint="https://splunk.example.com:8088"
            ... )
            >>> siem.update_config(config)
        """
        data = {
            "enabled": config.enabled,
            "provider": config.provider.value,
            "endpoint": config.endpoint,
            "api_key": config.api_key,
            "batch_size": config.batch_size,
            "flush_interval": config.flush_interval,
            "include_debug": config.include_debug,
        }
        response = self._client._request("PUT", "siem/config", json=data)
        return SIEMConfig.from_dict(response)
    
    def enable(self) -> SIEMConfig:
        """
        Enable SIEM integration.
        
        Returns:
            Updated SIEMConfig object
        
        Example:
            >>> config = siem.enable()
        """
        response = self._client._request("POST", "siem/enable")
        return SIEMConfig.from_dict(response)
    
    def disable(self) -> SIEMConfig:
        """
        Disable SIEM integration.
        
        Returns:
            Updated SIEMConfig object
        
        Example:
            >>> config = siem.disable()
        """
        response = self._client._request("POST", "siem/disable")
        return SIEMConfig.from_dict(response)
    
    def test_connection(self) -> Dict[str, Any]:
        """
        Test SIEM connection.
        
        Returns:
            Test result with success status and details
        
        Example:
            >>> result = siem.test_connection()
            >>> if result["success"]:
            ...     print("Connection successful!")
        """
        response = self._client._request("POST", "siem/test")
        return response
    
    # =====================================================================
    # Events
    # =====================================================================
    
    def send_event(self, event: SIEMEvent) -> str:
        """
        Send a single SIEM event.
        
        Args:
            event: SIEMEvent to send
        
        Returns:
            Event ID
        
        Example:
            >>> event = SIEMEvent(
            ...     timestamp=datetime.now(),
            ...     severity=EventSeverity.WARNING,
            ...     category=EventCategory.AUTHENTICATION,
            ...     title="Failed Login",
            ...     message="User failed to authenticate"
            ... )
            >>> event_id = siem.send_event(event)
        """
        response = self._client._request("POST", "siem/events", json=event.to_dict())
        return response.get("event_id", "")
    
    def send_events(self, events: List[SIEMEvent]) -> Dict[str, Any]:
        """
        Send multiple SIEM events in a batch.
        
        Args:
            events: List of SIEMEvents to send
        
        Returns:
            Batch result with success/failure counts
        
        Example:
            >>> events = [create_event() for _ in range(10)]
            >>> result = siem.send_events(events)
            >>> print(f"Sent: {result['sent']}, Failed: {result['failed']}")
        """
        data = {"events": [e.to_dict() for e in events]}
        response = self._client._request("POST", "siem/events/batch", json=data)
        return response
    
    def send_violation(self, violation: Violation) -> str:
        """
        Send a violation as a SIEM event.
        
        Args:
            violation: Violation to send
        
        Returns:
            Event ID
        
        Example:
            >>> violation = proxy.get_violation("viol-123")
            >>> event_id = siem.send_violation(violation)
        """
        event = SIEMEvent.from_violation(violation)
        return self.send_event(event)
    
    def send_violations(self, violations: List[Violation]) -> Dict[str, Any]:
        """
        Send multiple violations as SIEM events.
        
        Args:
            violations: List of violations to send
        
        Returns:
            Batch result
        
        Example:
            >>> violations = proxy.get_violations(filter=high_severity_filter)
            >>> result = siem.send_violations(violations)
        """
        events = [SIEMEvent.from_violation(v) for v in violations]
        return self.send_events(events)
    
    def query_events(
        self,
        start_time: Optional[datetime] = None,
        end_time: Optional[datetime] = None,
        severity: Optional[EventSeverity] = None,
        category: Optional[EventCategory] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[SIEMEvent]:
        """
        Query SIEM events.
        
        Args:
            start_time: Start of time range
            end_time: End of time range
            severity: Filter by severity
            category: Filter by category
            limit: Maximum results
            offset: Pagination offset
        
        Returns:
            List of SIEMEvent objects
        
        Example:
            >>> events = siem.query_events(
            ...     start_time=datetime.now() - timedelta(days=7),
            ...     severity=EventSeverity.CRITICAL
            ... )
        """
        params = {"limit": limit, "offset": offset}
        if start_time:
            params["start_time"] = start_time.isoformat()
        if end_time:
            params["end_time"] = end_time.isoformat()
        if severity:
            params["severity"] = severity.value
        if category:
            params["category"] = category.value
        
        response = self._client._request("GET", "siem/events", params=params)
        
        events = []
        if "events" in response:
            events = [SIEMEvent.from_dict(e) for e in response["events"]]
        elif "data" in response:
            events = [SIEMEvent.from_dict(e) for e in response["data"]]
        
        return events
    
    def get_event(self, event_id: str) -> SIEMEvent:
        """
        Get a specific SIEM event.
        
        Args:
            event_id: Event identifier
        
        Returns:
            SIEMEvent object
        
        Example:
            >>> event = siem.get_event("event-123")
        """
        response = self._client._request("GET", f"siem/events/{event_id}")
        return SIEMEvent.from_dict(response)
    
    # =====================================================================
    # Statistics
    # =====================================================================
    
    def get_stats(self) -> SIEMStats:
        """
        Get SIEM integration statistics.
        
        Returns:
            SIEMStats object
        
        Example:
            >>> stats = siem.get_stats()
            >>> print(f"Events sent: {stats.events_sent}")
        """
        response = self._client._request("GET", "siem/stats")
        return SIEMStats.from_dict(response)
    
    def get_stats_history(
        self,
        period: str = "24h",
    ) -> Dict[str, Any]:
        """
        Get historical SIEM statistics.
        
        Args:
            period: Time period (e.g., "1h", "24h", "7d")
        
        Returns:
            Historical stats data
        
        Example:
            >>> history = siem.get_stats_history("7d")
        """
        params = {"period": period}
        response = self._client._request("GET", "siem/stats/history", params=params)
        return response
    
    # =====================================================================
    # Providers
    # =====================================================================
    
    def list_providers(self) -> List[Dict[str, Any]]:
        """
        List supported SIEM providers.
        
        Returns:
            List of provider information
        
        Example:
            >>> providers = siem.list_providers()
            >>> for p in providers:
            ...     print(f"{p['name']}: {p['description']}")
        """
        response = self._client._request("GET", "siem/providers")
        return response.get("providers", [])
    
    def configure_splunk(self, endpoint: str, token: str) -> SIEMConfig:
        """
        Configure Splunk SIEM integration.
        
        Args:
            endpoint: Splunk HEC endpoint URL
            token: Splunk HEC token
        
        Returns:
            Updated SIEMConfig object
        
        Example:
            >>> config = siem.configure_splunk(
            ...     "https://splunk.example.com:8088/services/collector",
            ...     "your-hec-token"
            ... )
        """
        config = SIEMConfig(
            enabled=True,
            provider=SIEMProvider.SPLUNK,
            endpoint=endpoint,
            api_key=token,
        )
        return self.update_config(config)
    
    def configure_elastic(
        self,
        endpoint: str,
        api_key: Optional[str] = None,
    ) -> SIEMConfig:
        """
        Configure Elasticsearch SIEM integration.
        
        Args:
            endpoint: Elasticsearch endpoint URL
            api_key: Optional API key
        
        Returns:
            Updated SIEMConfig object
        
        Example:
            >>> config = siem.configure_elastic(
            ...     "https://elastic.example.com:9200",
            ...     "your-api-key"
            ... )
        """
        config = SIEMConfig(
            enabled=True,
            provider=SIEMProvider.ELASTIC,
            endpoint=endpoint,
            api_key=api_key,
        )
        return self.update_config(config)
    
    def configure_custom(self, endpoint: str, headers: Dict[str, str] = None) -> SIEMConfig:
        """
        Configure custom webhook SIEM integration.
        
        Args:
            endpoint: Custom webhook endpoint URL
            headers: Optional custom headers
        
        Returns:
            Updated SIEMConfig object
        
        Example:
            >>> config = siem.configure_custom(
            ...     "https://your-siem.example.com/webhook",
            ...     {"Authorization": "Bearer token"}
            ... )
        """
        config = SIEMConfig(
            enabled=True,
            provider=SIEMProvider.CUSTOM,
            endpoint=endpoint,
        )
        if headers:
            config.filters = {"headers": headers}
        return self.update_config(config)


# Type hint reference - import only for type checking
if False:
    from aegisgate.client import AegisGateClient
