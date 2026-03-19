# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Compliance service for the AegisGate Python SDK.

This module provides compliance framework management and reporting functionality.
"""

from __future__ import annotations

from typing import Optional, List, Dict, Any
from datetime import datetime
import logging

from aegisgate.models.compliance import (
    Framework,
    ComplianceCheck,
    ComplianceReport,
    ComplianceStatus,
    Control,
    Finding,
    FrameworkType,
)
from aegisgate.exceptions import ResourceNotFoundError

logger = logging.getLogger(__name__)


class ComplianceService:
    """
    Service for compliance framework management.
    
    This service handles:
    - Framework enumeration and details
    - Compliance checking and reporting
    - Control assessment
    - Finding management
    
    Example:
        >>> client = AegisGateClient(base_url="http://localhost:8080")
        >>> compliance = client.compliance
        >>> 
        >>> # List available frameworks
        >>> frameworks = compliance.list_frameworks()
        >>> for fw in frameworks:
        ...     print(f"{fw.name} ({fw.version})")
        >>> 
        >>> # Run a compliance check
        >>> check = compliance.run_check(FrameworkType.MITRE_ATLAS)
        >>> print(f"Compliance score: {check.compliance_score}%")
    """
    
    def __init__(self, client: "AegisGateClient"):
        """
        Initialize the compliance service.
        
        Args:
            client: The AegisGate client instance
        """
        self._client = client
    
    # =====================================================================
    # Frameworks
    # =====================================================================
    
    def list_frameworks(self) -> List[Framework]:
        """
        List all available compliance frameworks.
        
        Returns:
            List of Framework objects
        
        Example:
            >>> frameworks = compliance.list_frameworks()
            >>> for fw in frameworks:
            ...     print(f"{fw.name}: {fw.controls_count} controls")
        """
        response = self._client._request("GET", "compliance/frameworks")
        
        frameworks = []
        if "frameworks" in response:
            frameworks = [Framework.from_dict(f) for f in response["frameworks"]]
        elif "data" in response:
            frameworks = [Framework.from_dict(f) for f in response["data"]]
        
        return frameworks
    
    def get_framework(self, framework_id: str) -> Framework:
        """
        Get a specific framework by ID.
        
        Args:
            framework_id: Framework identifier
        
        Returns:
            Framework object
        
        Raises:
            ResourceNotFoundError: If framework doesn't exist
        
        Example:
            >>> fw = compliance.get_framework("mitre_atlas")
            >>> print(f"Framework: {fw.name}")
        """
        response = self._client._request("GET", f"compliance/frameworks/{framework_id}")
        return Framework.from_dict(response)
    
    def get_framework_controls(
        self,
        framework_id: str,
        category: Optional[str] = None,
    ) -> List[Control]:
        """
        Get all controls for a framework.
        
        Args:
            framework_id: Framework identifier
            category: Optional filter by category
        
        Returns:
            List of Control objects
        
        Example:
            >>> controls = compliance.get_framework_controls("nist_ai_rmf")
            >>> for ctrl in controls:
            ...     print(f"{ctrl.id}: {ctrl.name}")
        """
        params = {"category": category} if category else None
        response = self._client._request(
            "GET", f"compliance/frameworks/{framework_id}/controls", params=params
        )
        
        controls = []
        if "controls" in response:
            controls = [Control.from_dict(c) for c in response["controls"]]
        
        return controls
    
    # =====================================================================
    # Compliance Checks
    # =====================================================================
    
    def run_check(
        self,
        framework: FrameworkType,
        period: Optional[str] = None,
    ) -> ComplianceCheck:
        """
        Run a compliance check for a framework.
        
        Args:
            framework: Framework to check
            period: Optional time period for the check
        
        Returns:
            ComplianceCheck result
        
        Example:
            >>> check = compliance.run_check(FrameworkType.MITRE_ATLAS)
            >>> print(f"Status: {check.status.value}")
            >>> print(f"Score: {check.compliance_score}%")
        """
        data = {"framework": framework.value}
        if period:
            data["period"] = period
        
        response = self._client._request(
            "POST", "compliance/check", json=data
        )
        return ComplianceCheck.from_dict(response)
    
    def run_all_checks(self) -> List[ComplianceCheck]:
        """
        Run compliance checks for all supported frameworks.
        
        Returns:
            List of ComplianceCheck results
        
        Example:
            >>> checks = compliance.run_all_checks()
            >>> for check in checks:
            ...     print(f"{check.framework_name}: {check.compliance_score}%")
        """
        response = self._client._request("POST", "compliance/check/all")
        
        checks = []
        if "checks" in response:
            checks = [ComplianceCheck.from_dict(c) for c in response["checks"]]
        
        return checks
    
    def get_check(self, check_id: str) -> ComplianceCheck:
        """
        Get a specific compliance check result.
        
        Args:
            check_id: Check identifier
        
        Returns:
            ComplianceCheck object
        
        Raises:
            ResourceNotFoundError: If check doesn't exist
        
        Example:
            >>> check = compliance.get_check("check-123")
        """
        response = self._client._request("GET", f"compliance/checks/{check_id}")
        return ComplianceCheck.from_dict(response)
    
    def list_checks(
        self,
        framework: Optional[FrameworkType] = None,
        status: Optional[ComplianceStatus] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[ComplianceCheck]:
        """
        List compliance check history.
        
        Args:
            framework: Optional framework filter
            status: Optional status filter
            limit: Maximum number of results
            offset: Pagination offset
        
        Returns:
            List of ComplianceCheck objects
        
        Example:
            >>> checks = compliance.list_checks(status=ComplianceStatus.FAIL)
        """
        params = {"limit": limit, "offset": offset}
        if framework:
            params["framework"] = framework.value
        if status:
            params["status"] = status.value
        
        response = self._client._request("GET", "compliance/checks", params=params)
        
        checks = []
        if "checks" in response:
            checks = [ComplianceCheck.from_dict(c) for c in response["checks"]]
        elif "data" in response:
            checks = [ComplianceCheck.from_dict(c) for c in response["data"]]
        
        return checks
    
    # =====================================================================
    # Reports
    # =====================================================================
    
    def generate_report(
        self,
        framework_ids: Optional[List[str]] = None,
        period_start: Optional[datetime] = None,
        period_end: Optional[datetime] = None,
    ) -> ComplianceReport:
        """
        Generate a compliance report.
        
        Args:
            framework_ids: List of frameworks to include (all if None)
            period_start: Report period start
            period_end: Report period end
        
        Returns:
            ComplianceReport object
        
        Example:
            >>> report = compliance.generate_report(
            ...     framework_ids=["mitre_atlas", "nist_ai_rmf"],
            ...     period_end=datetime.now()
            ... )
            >>> print(f"Overall score: {report.overall_score}%")
        """
        data = {}
        if framework_ids:
            data["frameworks"] = framework_ids
        if period_start:
            data["period_start"] = period_start.isoformat()
        if period_end:
            data["period_end"] = period_end.isoformat()
        
        response = self._client._request(
            "POST", "compliance/reports", json=data
        )
        return ComplianceReport.from_dict(response)
    
    def get_report(self, report_id: str) -> ComplianceReport:
        """
        Get a specific compliance report.
        
        Args:
            report_id: Report identifier
        
        Returns:
            ComplianceReport object
        
        Raises:
            ResourceNotFoundError: If report doesn't exist
        
        Example:
            >>> report = compliance.get_report("report-123")
        """
        response = self._client._request("GET", f"compliance/reports/{report_id}")
        return ComplianceReport.from_dict(response)
    
    def list_reports(
        self,
        framework: Optional[FrameworkType] = None,
        limit: int = 50,
        offset: int = 0,
    ) -> List[ComplianceReport]:
        """
        List compliance reports.
        
        Args:
            framework: Optional framework filter
            limit: Maximum number of results
            offset: Pagination offset
        
        Returns:
            List of ComplianceReport objects
        
        Example:
            >>> reports = compliance.list_reports()
            >>> for report in reports:
            ...     print(f"{report.id}: {report.generated_at}")
        """
        params = {"limit": limit, "offset": offset}
        if framework:
            params["framework"] = framework.value
        
        response = self._client._request("GET", "compliance/reports", params=params)
        
        reports = []
        if "reports" in response:
            reports = [ComplianceReport.from_dict(r) for r in response["reports"]]
        elif "data" in response:
            reports = [ComplianceReport.from_dict(r) for r in response["data"]]
        
        return reports
    
    # =====================================================================
    # Controls
    # =====================================================================
    
    def get_control(self, control_id: str) -> Control:
        """
        Get a specific control.
        
        Args:
            control_id: Control identifier
        
        Returns:
            Control object
        
        Example:
            >>> control = compliance.get_control("GV1")
        """
        response = self._client._request("GET", f"compliance/controls/{control_id}")
        return Control.from_dict(response)
    
    def update_control_status(
        self,
        control_id: str,
        status: ComplianceStatus,
        notes: Optional[str] = None,
    ) -> Control:
        """
        Update a control's compliance status.
        
        Args:
            control_id: Control identifier
            status: New compliance status
            notes: Optional notes
        
        Returns:
            Updated Control object
        
        Example:
            >>> control = compliance.update_control_status(
            ...     "GV1",
            ...     ComplianceStatus.PASS,
            ...     notes="Implemented and tested"
            ... )
        """
        data = {"status": status.value}
        if notes:
            data["notes"] = notes
        
        response = self._client._request(
            "PUT", f"compliance/controls/{control_id}", json=data
        )
        return Control.from_dict(response)
    
    # =====================================================================
    # Findings
    # =====================================================================
    
    def list_findings(
        self,
        control_id: Optional[str] = None,
        severity: Optional[str] = None,
        status: Optional[str] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[Finding]:
        """
        List compliance findings.
        
        Args:
            control_id: Optional control filter
            severity: Optional severity filter
            status: Optional status filter (open/resolved)
            limit: Maximum number of results
            offset: Pagination offset
        
        Returns:
            List of Finding objects
        
        Example:
            >>> findings = compliance.list_findings(status="open")
        """
        params = {"limit": limit, "offset": offset}
        if control_id:
            params["control_id"] = control_id
        if severity:
            params["severity"] = severity
        if status:
            params["status"] = status
        
        response = self._client._request("GET", "compliance/findings", params=params)
        
        findings = []
        if "findings" in response:
            findings = [Finding.from_dict(f) for f in response["findings"]]
        elif "data" in response:
            findings = [Finding.from_dict(f) for f in response["data"]]
        
        return findings
    
    def resolve_finding(
        self,
        finding_id: str,
        resolution: str,
    ) -> Finding:
        """
        Mark a finding as resolved.
        
        Args:
            finding_id: Finding identifier
            resolution: Resolution notes
        
        Returns:
            Updated Finding object
        
        Example:
            >>> finding = compliance.resolve_finding(
            ...     "finding-123",
            ...     "Patched vulnerability"
            ... )
        """
        data = {"resolution": resolution}
        response = self._client._request(
            "POST", f"compliance/findings/{finding_id}/resolve", json=data
        )
        return Finding.from_dict(response)


# Type hint reference - import only for type checking
if False:
    from aegisgate.client import AegisGateClient
