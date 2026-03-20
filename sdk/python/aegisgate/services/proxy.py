# SPDX-License-Identifier: MIT
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.

"""
Proxy service for the AegisGate Python SDK.

Provides AI proxy operations including request inspection, violation detection,
and traffic analysis.
"""

from typing import Any, Dict, List, Optional

from aegisgate.connection import SyncConnection, AsyncConnection
from aegisgate.models import DetectionResult, AnomalyResult, Violation, ProxyStats


class ProxyService:
    """Synchronous proxy service for AI gateway operations."""

    def __init__(self, connection: SyncConnection):
        self._conn = connection

    def get_stats(self) -> ProxyStats:
        """
        Get proxy statistics.

        Returns:
            ProxyStats object with current statistics
        """
        response = self._conn.get("/api/v1/proxy/stats")
        return ProxyStats.from_dict(response)

    def inspect_request(
        self,
        content: str,
        model: Optional[str] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> DetectionResult:
        """
        Inspect a request for security violations.

        Args:
            content: The content to inspect
            model: Target model identifier (optional)
            metadata: Additional metadata for context

        Returns:
            DetectionResult with any detected violations
        """
        payload = {
            "content": content,
        }
        if model:
            payload["model"] = model
        if metadata:
            payload["metadata"] = metadata

        response = self._conn.post("/api/v1/proxy/inspect", json_data=payload)
        return DetectionResult.from_dict(response)

    def inspect_response(
        self,
        content: str,
        model: Optional[str] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> DetectionResult:
        """
        Inspect a response for security violations.

        Args:
            content: The response content to inspect
            model: Source model identifier (optional)
            metadata: Additional metadata for context

        Returns:
            DetectionResult with any detected violations
        """
        payload = {
            "content": content,
        }
        if model:
            payload["model"] = model
        if metadata:
            payload["metadata"] = metadata

        response = self._conn.post("/api/v1/proxy/inspect-response", json_data=payload)
        return DetectionResult.from_dict(response)

    def detect_anomalies(
        self,
        features: Dict[str, float],
        context: Optional[Dict[str, Any]] = None,
    ) -> AnomalyResult:
        """
        Detect anomalies using ML models.

        Args:
            features: Feature vector for anomaly detection
            context: Additional context for the detection

        Returns:
            AnomalyResult with anomaly detection results
        """
        payload = {
            "features": features,
        }
        if context:
            payload["context"] = context

        response = self._conn.post("/api/v1/proxy/anomaly", json_data=payload)
        return AnomalyResult.from_dict(response)

    def list_violations(
        self,
        limit: int = 100,
        offset: int = 0,
        severity: Optional[str] = None,
        violation_type: Optional[str] = None,
    ) -> List[Violation]:
        """
        List detected violations.

        Args:
            limit: Maximum number of violations to return
            offset: Offset for pagination
            severity: Filter by severity level
            violation_type: Filter by violation type

        Returns:
            List of Violation objects
        """
        params = {"limit": limit, "offset": offset}
        if severity:
            params["severity"] = severity
        if violation_type:
            params["type"] = violation_type

        response = self._conn.get("/api/v1/proxy/violations", params=params)
        violations = response.get("violations", [])
        return [Violation.from_dict(v) for v in violations]

    def get_violation(self, violation_id: str) -> Violation:
        """
        Get a specific violation by ID.

        Args:
            violation_id: The violation ID

        Returns:
            Violation object
        """
        response = self._conn.get(f"/api/v1/proxy/violations/{violation_id}")
        return Violation.from_dict(response)

    def block_request(
        self,
        request_id: str,
        reason: str,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Block a request.

        Args:
            request_id: The request ID to block
            reason: Reason for blocking
            metadata: Additional metadata

        Returns:
            Blocking result
        """
        payload = {
            "request_id": request_id,
            "reason": reason,
        }
        if metadata:
            payload["metadata"] = metadata

        response = self._conn.post("/api/v1/proxy/block", json_data=payload)
        return response

    def allow_request(
        self,
        request_id: str,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Allow a previously blocked request.

        Args:
            request_id: The request ID to allow
            metadata: Additional metadata

        Returns:
            Result of the operation
        """
        payload = {"request_id": request_id}
        if metadata:
            payload["metadata"] = metadata

        response = self._conn.post("/api/v1/proxy/allow", json_data=payload)
        return response

    def configure_content_filter(
        self,
        filters: List[Dict[str, Any]],
        enabled: bool = True,
    ) -> Dict[str, Any]:
        """
        Configure content filtering rules.

        Args:
            filters: List of filter configurations
            enabled: Whether to enable the filters

        Returns:
            Configuration result
        """
        payload = {
            "filters": filters,
            "enabled": enabled,
        }
        response = self._conn.put("/api/v1/proxy/content-filter", json_data=payload)
        return response


class AsyncProxyService:
    """Asynchronous proxy service for AI gateway operations."""

    def __init__(self, connection: AsyncConnection):
        self._conn = connection

    async def get_stats(self) -> ProxyStats:
        """
        Get proxy statistics.

        Returns:
            ProxyStats object with current statistics
        """
        response = await self._conn.get("/api/v1/proxy/stats")
        return ProxyStats.from_dict(response)

    async def inspect_request(
        self,
        content: str,
        model: Optional[str] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> DetectionResult:
        """
        Inspect a request for security violations.

        Args:
            content: The content to inspect
            model: Target model identifier (optional)
            metadata: Additional metadata for context

        Returns:
            DetectionResult with any detected violations
        """
        payload = {
            "content": content,
        }
        if model:
            payload["model"] = model
        if metadata:
            payload["metadata"] = metadata

        response = await self._conn.post("/api/v1/proxy/inspect", json_data=payload)
        return DetectionResult.from_dict(response)

    async def inspect_response(
        self,
        content: str,
        model: Optional[str] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> DetectionResult:
        """
        Inspect a response for security violations.

        Args:
            content: The response content to inspect
            model: Source model identifier (optional)
            metadata: Additional metadata for context

        Returns:
            DetectionResult with any detected violations
        """
        payload = {
            "content": content,
        }
        if model:
            payload["model"] = model
        if metadata:
            payload["metadata"] = metadata

        response = await self._conn.post("/api/v1/proxy/inspect-response", json_data=payload)
        return DetectionResult.from_dict(response)

    async def detect_anomalies(
        self,
        features: Dict[str, float],
        context: Optional[Dict[str, Any]] = None,
    ) -> AnomalyResult:
        """
        Detect anomalies using ML models.

        Args:
            features: Feature vector for anomaly detection
            context: Additional context for the detection

        Returns:
            AnomalyResult with anomaly detection results
        """
        payload = {
            "features": features,
        }
        if context:
            payload["context"] = context

        response = await self._conn.post("/api/v1/proxy/anomaly", json_data=payload)
        return AnomalyResult.from_dict(response)

    async def list_violations(
        self,
        limit: int = 100,
        offset: int = 0,
        severity: Optional[str] = None,
        violation_type: Optional[str] = None,
    ) -> List[Violation]:
        """
        List detected violations.

        Args:
            limit: Maximum number of violations to return
            offset: Offset for pagination
            severity: Filter by severity level
            violation_type: Filter by violation type

        Returns:
            List of Violation objects
        """
        params = {"limit": limit, "offset": offset}
        if severity:
            params["severity"] = severity
        if violation_type:
            params["type"] = violation_type

        response = await self._conn.get("/api/v1/proxy/violations", params=params)
        violations = response.get("violations", [])
        return [Violation.from_dict(v) for v in violations]

    async def get_violation(self, violation_id: str) -> Violation:
        """
        Get a specific violation by ID.

        Args:
            violation_id: The violation ID

        Returns:
            Violation object
        """
        response = await self._conn.get(f"/api/v1/proxy/violations/{violation_id}")
        return Violation.from_dict(response)

    async def block_request(
        self,
        request_id: str,
        reason: str,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Block a request.

        Args:
            request_id: The request ID to block
            reason: Reason for blocking
            metadata: Additional metadata

        Returns:
            Blocking result
        """
        payload = {
            "request_id": request_id,
            "reason": reason,
        }
        if metadata:
            payload["metadata"] = metadata

        response = await self._conn.post("/api/v1/proxy/block", json_data=payload)
        return response

    async def allow_request(
        self,
        request_id: str,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Allow a previously blocked request.

        Args:
            request_id: The request ID to allow
            metadata: Additional metadata

        Returns:
            Result of the operation
        """
        payload = {"request_id": request_id}
        if metadata:
            payload["metadata"] = metadata

        response = await self._conn.post("/api/v1/proxy/allow", json_data=payload)
        return response

    async def configure_content_filter(
        self,
        filters: List[Dict[str, Any]],
        enabled: bool = True,
    ) -> Dict[str, Any]:
        """
        Configure content filtering rules.

        Args:
            filters: List of filter configurations
            enabled: Whether to enable the filters

        Returns:
            Configuration result
        """
        payload = {
            "filters": filters,
            "enabled": enabled,
        }
        response = await self._conn.put("/api/v1/proxy/content-filter", json_data=payload)
        return response