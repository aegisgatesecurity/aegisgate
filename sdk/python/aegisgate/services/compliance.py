# SPDX-License-Identifier: MIT
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.

"""
Compliance service for the AegisGate Python SDK.

Provides compliance framework management including HIPAA, SOC2, GDPR, and more.
"""

from typing import Any, Dict, List, Optional

from aegisgate.connection import SyncConnection, AsyncConnection
from aegisgate.models import ComplianceResult, ComplianceControl


class ComplianceService:
    """Synchronous compliance service for framework management."""

    def __init__(self, connection: SyncConnection):
        self._conn = connection

    def list_frameworks(self) -> List[Dict[str, Any]]:
        """
        List available compliance frameworks.

        Returns:
            List of framework information dictionaries
        """
        response = self._conn.get("/api/v1/compliance/frameworks")
        return response.get("frameworks", [])

    def get_framework(self, framework: str) -> Dict[str, Any]:
        """
        Get details for a specific compliance framework.

        Args:
            framework: Framework identifier (e.g., 'hipaa', 'soc2', 'gdpr')

        Returns:
            Framework details
        """
        response = self._conn.get(f"/api/v1/compliance/frameworks/{framework}")
        return response

    def run_check(
        self,
        framework: str,
        controls: Optional[List[str]] = None,
    ) -> ComplianceResult:
        """
        Run a compliance check for a framework.

        Args:
            framework: Framework identifier (e.g., 'hipaa', 'soc2', 'gdpr')
            controls: Optional list of specific controls to check

        Returns:
            ComplianceResult with check results
        """
        payload = {"framework": framework}
        if controls:
            payload["controls"] = controls

        response = self._conn.post("/api/v1/compliance/check", json_data=payload)
        return ComplianceResult.from_dict(response)

    def get_controls(self, framework: str) -> List[ComplianceControl]:
        """
        Get compliance controls for a framework.

        Args:
            framework: Framework identifier

        Returns:
            List of ComplianceControl objects
        """
        response = self._conn.get(f"/api/v1/compliance/frameworks/{framework}/controls")
        controls = response.get("controls", [])
        return [ComplianceControl.from_dict(c) for c in controls]

    def get_control(self, framework: str, control_id: str) -> ComplianceControl:
        """
        Get a specific compliance control.

        Args:
            framework: Framework identifier
            control_id: Control identifier

        Returns:
            ComplianceControl object
        """
        response = self._conn.get(
            f"/api/v1/compliance/frameworks/{framework}/controls/{control_id}"
        )
        return ComplianceControl.from_dict(response)

    def update_control_status(
        self,
        framework: str,
        control_id: str,
        status: str,
        evidence: Optional[List[Dict[str, Any]]] = None,
        notes: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Update the status of a compliance control.

        Args:
            framework: Framework identifier
            control_id: Control identifier
            status: New status ('pass', 'fail', 'in_progress', 'not_applicable')
            evidence: Evidence for the status change
            notes: Additional notes

        Returns:
            Updated control status
        """
        payload = {
            "status": status,
        }
        if evidence:
            payload["evidence"] = evidence
        if notes:
            payload["notes"] = notes

        response = self._conn.patch(
            f"/api/v1/compliance/frameworks/{framework}/controls/{control_id}",
            json_data=payload,
        )
        return response

    def get_results(
        self,
        framework: Optional[str] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[ComplianceResult]:
        """
        Get historical compliance check results.

        Args:
            framework: Filter by framework (optional)
            limit: Maximum number of results
            offset: Pagination offset

        Returns:
            List of ComplianceResult objects
        """
        params = {"limit": limit, "offset": offset}
        if framework:
            params["framework"] = framework

        response = self._conn.get("/api/v1/compliance/results", params=params)
        results = response.get("results", [])
        return [ComplianceResult.from_dict(r) for r in results]

    def generate_report(
        self,
        framework: str,
        format: str = "json",
        include_evidence: bool = True,
    ) -> Dict[str, Any]:
        """
        Generate a compliance report.

        Args:
            framework: Framework identifier
            format: Report format ('json', 'pdf', 'html')
            include_evidence: Whether to include evidence in the report

        Returns:
            Report data or download URL
        """
        payload = {
            "framework": framework,
            "format": format,
            "include_evidence": include_evidence,
        }
        response = self._conn.post("/api/v1/compliance/report", json_data=payload)
        return response

    def schedule_check(
        self,
        framework: str,
        schedule: str,
        controls: Optional[List[str]] = None,
    ) -> Dict[str, Any]:
        """
        Schedule a recurring compliance check.

        Args:
            framework: Framework identifier
            schedule: Cron schedule expression
            controls: Optional list of specific controls to check

        Returns:
            Scheduled check configuration
        """
        payload = {
            "framework": framework,
            "schedule": schedule,
        }
        if controls:
            payload["controls"] = controls

        response = self._conn.post("/api/v1/compliance/schedule", json_data=payload)
        return response

    def list_schedules(self) -> List[Dict[str, Any]]:
        """
        List all scheduled compliance checks.

        Returns:
            List of scheduled checks
        """
        response = self._conn.get("/api/v1/compliance/schedule")
        return response.get("schedules", [])

    def delete_schedule(self, schedule_id: str) -> None:
        """
        Delete a scheduled compliance check.

        Args:
            schedule_id: Schedule identifier
        """
        self._conn.delete(f"/api/v1/compliance/schedule/{schedule_id}")


class AsyncComplianceService:
    """Asynchronous compliance service for framework management."""

    def __init__(self, connection: AsyncConnection):
        self._conn = connection

    async def list_frameworks(self) -> List[Dict[str, Any]]:
        """
        List available compliance frameworks.

        Returns:
            List of framework information dictionaries
        """
        response = await self._conn.get("/api/v1/compliance/frameworks")
        return response.get("frameworks", [])

    async def get_framework(self, framework: str) -> Dict[str, Any]:
        """
        Get details for a specific compliance framework.

        Args:
            framework: Framework identifier (e.g., 'hipaa', 'soc2', 'gdpr')

        Returns:
            Framework details
        """
        response = await self._conn.get(f"/api/v1/compliance/frameworks/{framework}")
        return response

    async def run_check(
        self,
        framework: str,
        controls: Optional[List[str]] = None,
    ) -> ComplianceResult:
        """
        Run a compliance check for a framework.

        Args:
            framework: Framework identifier (e.g., 'hipaa', 'soc2', 'gdpr')
            controls: Optional list of specific controls to check

        Returns:
            ComplianceResult with check results
        """
        payload = {"framework": framework}
        if controls:
            payload["controls"] = controls

        response = await self._conn.post("/api/v1/compliance/check", json_data=payload)
        return ComplianceResult.from_dict(response)

    async def get_controls(self, framework: str) -> List[ComplianceControl]:
        """
        Get compliance controls for a framework.

        Args:
            framework: Framework identifier

        Returns:
            List of ComplianceControl objects
        """
        response = await self._conn.get(f"/api/v1/compliance/frameworks/{framework}/controls")
        controls = response.get("controls", [])
        return [ComplianceControl.from_dict(c) for c in controls]

    async def get_control(self, framework: str, control_id: str) -> ComplianceControl:
        """
        Get a specific compliance control.

        Args:
            framework: Framework identifier
            control_id: Control identifier

        Returns:
            ComplianceControl object
        """
        response = await self._conn.get(
            f"/api/v1/compliance/frameworks/{framework}/controls/{control_id}"
        )
        return ComplianceControl.from_dict(response)

    async def update_control_status(
        self,
        framework: str,
        control_id: str,
        status: str,
        evidence: Optional[List[Dict[str, Any]]] = None,
        notes: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Update the status of a compliance control.

        Args:
            framework: Framework identifier
            control_id: Control identifier
            status: New status ('pass', 'fail', 'in_progress', 'not_applicable')
            evidence: Evidence for the status change
            notes: Additional notes

        Returns:
            Updated control status
        """
        payload = {
            "status": status,
        }
        if evidence:
            payload["evidence"] = evidence
        if notes:
            payload["notes"] = notes

        response = await self._conn.patch(
            f"/api/v1/compliance/frameworks/{framework}/controls/{control_id}",
            json_data=payload,
        )
        return response

    async def get_results(
        self,
        framework: Optional[str] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[ComplianceResult]:
        """
        Get historical compliance check results.

        Args:
            framework: Filter by framework (optional)
            limit: Maximum number of results
            offset: Pagination offset

        Returns:
            List of ComplianceResult objects
        """
        params = {"limit": limit, "offset": offset}
        if framework:
            params["framework"] = framework

        response = await self._conn.get("/api/v1/compliance/results", params=params)
        results = response.get("results", [])
        return [ComplianceResult.from_dict(r) for r in results]

    async def generate_report(
        self,
        framework: str,
        format: str = "json",
        include_evidence: bool = True,
    ) -> Dict[str, Any]:
        """
        Generate a compliance report.

        Args:
            framework: Framework identifier
            format: Report format ('json', 'pdf', 'html')
            include_evidence: Whether to include evidence in the report

        Returns:
            Report data or download URL
        """
        payload = {
            "framework": framework,
            "format": format,
            "include_evidence": include_evidence,
        }
        response = await self._conn.post("/api/v1/compliance/report", json_data=payload)
        return response

    async def schedule_check(
        self,
        framework: str,
        schedule: str,
        controls: Optional[List[str]] = None,
    ) -> Dict[str, Any]:
        """
        Schedule a recurring compliance check.

        Args:
            framework: Framework identifier
            schedule: Cron schedule expression
            controls: Optional list of specific controls to check

        Returns:
            Scheduled check configuration
        """
        payload = {
            "framework": framework,
            "schedule": schedule,
        }
        if controls:
            payload["controls"] = controls

        response = await self._conn.post("/api/v1/compliance/schedule", json_data=payload)
        return response

    async def list_schedules(self) -> List[Dict[str, Any]]:
        """
        List all scheduled compliance checks.

        Returns:
            List of scheduled checks
        """
        response = await self._conn.get("/api/v1/compliance/schedule")
        return response.get("schedules", [])

    async def delete_schedule(self, schedule_id: str) -> None:
        """
        Delete a scheduled compliance check.

        Args:
            schedule_id: Schedule identifier
        """
        await self._conn.delete(f"/api/v1/compliance/schedule/{schedule_id}")