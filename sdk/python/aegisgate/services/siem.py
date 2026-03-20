# SPDX-License-Identifier: MIT
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.

"""
SIEM integration service for the AegisGate Python SDK.

Provides SIEM integration and event management for security monitoring.
"""

from typing import Any, Dict, List, Optional
from datetime import datetime

from aegisgate.connection import SyncConnection, AsyncConnection
from aegisgate.models import SIEMConfig, SIEMEvent, SIEMProvider


class SIEMService:
    """Synchronous SIEM integration service."""

    def __init__(self, connection: SyncConnection):
        self._conn = connection

    def list_integrations(self) -> List[Dict[str, Any]]:
        """
        List all SIEM integrations.

        Returns:
            List of SIEM integration configurations
        """
        response = self._conn.get("/api/v1/siem/integrations")
        return response.get("integrations", [])

    def get_integration(self, integration_id: str) -> SIEMConfig:
        """
        Get a specific SIEM integration.

        Args:
            integration_id: Integration identifier

        Returns:
            SIEMConfig object
        """
        response = self._conn.get(f"/api/v1/siem/integrations/{integration_id}")
        return SIEMConfig.from_dict(response)

    def create_integration(
        self,
        provider: SIEMProvider,
        endpoint: str,
        api_key: Optional[str] = None,
        batch_size: int = 100,
        flush_interval_seconds: int = 30,
        custom_headers: Optional[Dict[str, str]] = None,
    ) -> SIEMConfig:
        """
        Create a new SIEM integration.

        Args:
            provider: SIEM provider type
            endpoint: SIEM endpoint URL
            api_key: API key for authentication (optional)
            batch_size: Number of events to batch
            flush_interval_seconds: Seconds between flushes
            custom_headers: Additional HTTP headers

        Returns:
            Created SIEMConfig object
        """
        payload = {
            "provider": provider.value,
            "endpoint": endpoint,
            "batch_size": batch_size,
            "flush_interval_seconds": flush_interval_seconds,
        }
        if api_key:
            payload["api_key"] = api_key
        if custom_headers:
            payload["custom_headers"] = custom_headers

        response = self._conn.post("/api/v1/siem/integrations", json_data=payload)
        return SIEMConfig.from_dict(response)

    def update_integration(
        self,
        integration_id: str,
        **kwargs,
    ) -> SIEMConfig:
        """
        Update a SIEM integration.

        Args:
            integration_id: Integration identifier
            **kwargs: Fields to update (endpoint, api_key, batch_size, etc.)

        Returns:
            Updated SIEMConfig object
        """
        if "provider" in kwargs and isinstance(kwargs["provider"], SIEMProvider):
            kwargs["provider"] = kwargs["provider"].value

        response = self._conn.patch(
            f"/api/v1/siem/integrations/{integration_id}",
            json_data=kwargs,
        )
        return SIEMConfig.from_dict(response)

    def delete_integration(self, integration_id: str) -> None:
        """
        Delete a SIEM integration.

        Args:
            integration_id: Integration identifier
        """
        self._conn.delete(f"/api/v1/siem/integrations/{integration_id}")

    def test_integration(self, integration_id: str) -> Dict[str, Any]:
        """
        Test a SIEM integration connection.

        Args:
            integration_id: Integration identifier

        Returns:
            Test result with connection status
        """
        response = self._conn.post(f"/api/v1/siem/integrations/{integration_id}/test")
        return response

    def send_event(
        self,
        integration_id: str,
        event_type: str,
        severity: str,
        source: str,
        message: str,
        details: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Send an event to a SIEM integration.

        Args:
            integration_id: Integration identifier
            event_type: Type of event
            severity: Event severity ('low', 'medium', 'high', 'critical')
            source: Event source identifier
            message: Event message
            details: Additional event details

        Returns:
            Event submission result
        """
        payload = {
            "event_type": event_type,
            "severity": severity,
            "source": source,
            "message": message,
            "timestamp": datetime.utcnow().isoformat(),
        }
        if details:
            payload["details"] = details

        response = self._conn.post(
            f"/api/v1/siem/integrations/{integration_id}/events",
            json_data=payload,
        )
        return response

    def send_batch_events(
        self,
        integration_id: str,
        events: List[SIEMEvent],
    ) -> Dict[str, Any]:
        """
        Send multiple events to a SIEM integration.

        Args:
            integration_id: Integration identifier
            events: List of SIEMEvent objects

        Returns:
            Batch submission result
        """
        payload = {
            "events": [e.to_dict() for e in events],
        }
        response = self._conn.post(
            f"/api/v1/siem/integrations/{integration_id}/events/batch",
            json_data=payload,
        )
        return response

    def get_events(
        self,
        integration_id: Optional[str] = None,
        event_type: Optional[str] = None,
        severity: Optional[str] = None,
        start_time: Optional[datetime] = None,
        end_time: Optional[datetime] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[SIEMEvent]:
        """
        Query SIEM events with filters.

        Args:
            integration_id: Filter by integration (optional)
            event_type: Filter by event type (optional)
            severity: Filter by severity (optional)
            start_time: Filter by start time (optional)
            end_time: Filter by end time (optional)
            limit: Maximum number of events
            offset: Pagination offset

        Returns:
            List of SIEMEvent objects
        """
        params = {"limit": limit, "offset": offset}
        if integration_id:
            params["integration_id"] = integration_id
        if event_type:
            params["event_type"] = event_type
        if severity:
            params["severity"] = severity
        if start_time:
            params["start_time"] = start_time.isoformat()
        if end_time:
            params["end_time"] = end_time.isoformat()

        response = self._conn.get("/api/v1/siem/events", params=params)
        events = response.get("events", [])
        return [SIEMEvent.from_dict(e) for e in events]

    def flush_events(self, integration_id: str) -> Dict[str, Any]:
        """
        Flush queued events to a SIEM integration.

        Args:
            integration_id: Integration identifier

        Returns:
            Flush result with event count
        """
        response = self._conn.post(
            f"/api/v1/siem/integrations/{integration_id}/flush"
        )
        return response


class AsyncSIEMService:
    """Asynchronous SIEM integration service."""

    def __init__(self, connection: AsyncConnection):
        self._conn = connection

    async def list_integrations(self) -> List[Dict[str, Any]]:
        """
        List all SIEM integrations.

        Returns:
            List of SIEM integration configurations
        """
        response = await self._conn.get("/api/v1/siem/integrations")
        return response.get("integrations", [])

    async def get_integration(self, integration_id: str) -> SIEMConfig:
        """
        Get a specific SIEM integration.

        Args:
            integration_id: Integration identifier

        Returns:
            SIEMConfig object
        """
        response = await self._conn.get(f"/api/v1/siem/integrations/{integration_id}")
        return SIEMConfig.from_dict(response)

    async def create_integration(
        self,
        provider: SIEMProvider,
        endpoint: str,
        api_key: Optional[str] = None,
        batch_size: int = 100,
        flush_interval_seconds: int = 30,
        custom_headers: Optional[Dict[str, str]] = None,
    ) -> SIEMConfig:
        """
        Create a new SIEM integration.

        Args:
            provider: SIEM provider type
            endpoint: SIEM endpoint URL
            api_key: API key for authentication (optional)
            batch_size: Number of events to batch
            flush_interval_seconds: Seconds between flushes
            custom_headers: Additional HTTP headers

        Returns:
            Created SIEMConfig object
        """
        payload = {
            "provider": provider.value,
            "endpoint": endpoint,
            "batch_size": batch_size,
            "flush_interval_seconds": flush_interval_seconds,
        }
        if api_key:
            payload["api_key"] = api_key
        if custom_headers:
            payload["custom_headers"] = custom_headers

        response = await self._conn.post("/api/v1/siem/integrations", json_data=payload)
        return SIEMConfig.from_dict(response)

    async def update_integration(
        self,
        integration_id: str,
        **kwargs,
    ) -> SIEMConfig:
        """
        Update a SIEM integration.

        Args:
            integration_id: Integration identifier
            **kwargs: Fields to update

        Returns:
            Updated SIEMConfig object
        """
        if "provider" in kwargs and isinstance(kwargs["provider"], SIEMProvider):
            kwargs["provider"] = kwargs["provider"].value

        response = await self._conn.patch(
            f"/api/v1/siem/integrations/{integration_id}",
            json_data=kwargs,
        )
        return SIEMConfig.from_dict(response)

    async def delete_integration(self, integration_id: str) -> None:
        """
        Delete a SIEM integration.

        Args:
            integration_id: Integration identifier
        """
        await self._conn.delete(f"/api/v1/siem/integrations/{integration_id}")

    async def test_integration(self, integration_id: str) -> Dict[str, Any]:
        """
        Test a SIEM integration connection.

        Args:
            integration_id: Integration identifier

        Returns:
            Test result with connection status
        """
        response = await self._conn.post(f"/api/v1/siem/integrations/{integration_id}/test")
        return response

    async def send_event(
        self,
        integration_id: str,
        event_type: str,
        severity: str,
        source: str,
        message: str,
        details: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Send an event to a SIEM integration.

        Args:
            integration_id: Integration identifier
            event_type: Type of event
            severity: Event severity
            source: Event source identifier
            message: Event message
            details: Additional event details

        Returns:
            Event submission result
        """
        payload = {
            "event_type": event_type,
            "severity": severity,
            "source": source,
            "message": message,
            "timestamp": datetime.utcnow().isoformat(),
        }
        if details:
            payload["details"] = details

        response = await self._conn.post(
            f"/api/v1/siem/integrations/{integration_id}/events",
            json_data=payload,
        )
        return response

    async def send_batch_events(
        self,
        integration_id: str,
        events: List[SIEMEvent],
    ) -> Dict[str, Any]:
        """
        Send multiple events to a SIEM integration.

        Args:
            integration_id: Integration identifier
            events: List of SIEMEvent objects

        Returns:
            Batch submission result
        """
        payload = {
            "events": [e.to_dict() for e in events],
        }
        response = await self._conn.post(
            f"/api/v1/siem/integrations/{integration_id}/events/batch",
            json_data=payload,
        )
        return response

    async def get_events(
        self,
        integration_id: Optional[str] = None,
        event_type: Optional[str] = None,
        severity: Optional[str] = None,
        start_time: Optional[datetime] = None,
        end_time: Optional[datetime] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[SIEMEvent]:
        """
        Query SIEM events with filters.

        Args:
            integration_id: Filter by integration (optional)
            event_type: Filter by event type (optional)
            severity: Filter by severity (optional)
            start_time: Filter by start time (optional)
            end_time: Filter by end time (optional)
            limit: Maximum number of events
            offset: Pagination offset

        Returns:
            List of SIEMEvent objects
        """
        params = {"limit": limit, "offset": offset}
        if integration_id:
            params["integration_id"] = integration_id
        if event_type:
            params["event_type"] = event_type
        if severity:
            params["severity"] = severity
        if start_time:
            params["start_time"] = start_time.isoformat()
        if end_time:
            params["end_time"] = end_time.isoformat()

        response = await self._conn.get("/api/v1/siem/events", params=params)
        events = response.get("events", [])
        return [SIEMEvent.from_dict(e) for e in events]

    async def flush_events(self, integration_id: str) -> Dict[str, Any]:
        """
        Flush queued events to a SIEM integration.

        Args:
            integration_id: Integration identifier

        Returns:
            Flush result with event count
        """
        response = await self._conn.post(
            f"/api/v1/siem/integrations/{integration_id}/flush"
        )
        return response