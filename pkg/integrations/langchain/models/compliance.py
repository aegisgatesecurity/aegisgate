# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Compliance framework models for the AegisGate Python SDK.

This module provides data classes for compliance reporting, framework
mappings, and control assessments.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
from typing import Optional, List, Dict, Any


class ComplianceStatus(str, Enum):
    """
    Compliance check status.
    
    Attributes:
        PASS: Control requirement is met
        FAIL: Control requirement is not met
        WARNING: Control has warnings but passes
        NOT_APPLICABLE: Control does not apply
        NOT_ASSESSED: Control has not been assessed
    """
    PASS = "pass"
    FAIL = "fail"
    WARNING = "warning"
    NOT_APPLICABLE = "not_applicable"
    NOT_ASSESSED = "not_assessed"
    
    @classmethod
    def from_string(cls, value: str) -> ComplianceStatus:
        """Create ComplianceStatus from string value."""
        try:
            return cls(value.lower())
        except ValueError:
            return cls.NOT_ASSESSED


class FrameworkType(str, Enum):
    """
    Compliance framework types supported by AegisGate.
    
    Attributes:
        MITRE_ATLAS: MITRE ATLAS (Adversarial Threat Landscape for AI Systems)
        NIST_AI_RMF: NIST AI Risk Management Framework
        OWASP_LLM: OWASP LLM Top 10
        SOC2: SOC 2 Trust Services Criteria
        HIPAA: Health Insurance Portability and Accountability Act
        PCI_DSS: Payment Card Industry Data Security Standard
        GDPR: General Data Protection Regulation
        ISO_27001: ISO/IEC 27001 Information Security
        ISO_42001: ISO/IEC 42001 AI Management System
        CIS: Center for Internet Security Controls
    """
    MITRE_ATLAS = "mitre_atlas"
    NIST_AI_RMF = "nist_ai_rmf"
    OWASP_LLM = "owasp_llm"
    SOC2 = "soc2"
    HIPAA = "hipaa"
    PCI_DSS = "pci_dss"
    GDPR = "gdpr"
    ISO_27001 = "iso_27001"
    ISO_42001 = "iso_42001"
    CIS = "cis"
    
    @classmethod
    def from_string(cls, value: str) -> FrameworkType:
        """Create FrameworkType from string value."""
        try:
            return cls(value.lower().replace("-", "_"))
        except ValueError:
            return cls.MITRE_ATLAS


@dataclass
class Framework:
    """
    Represents a compliance framework.
    
    Attributes:
        id: Framework identifier
        name: Framework name
        version: Framework version
        framework_type: Type of framework
        description: Framework description
        controls_count: Number of controls in framework
        categories: Control categories
        supported: Whether framework is supported
        last_updated: Last update timestamp
    """
    id: str
    name: str
    framework_type: FrameworkType
    version: str = "1.0"
    description: str = ""
    controls_count: int = 0
    categories: List[str] = field(default_factory=list)
    supported: bool = True
    last_updated: Optional[datetime] = None
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> Framework:
        """Create Framework from dictionary."""
        fw_type = FrameworkType.from_string(data.get("framework_type", "mitre_atlas"))
        
        last_updated = None
        if data.get("last_updated"):
            try:
                last_updated = datetime.fromisoformat(
                    data["last_updated"].replace("Z", "+00:00")
                )
            except (ValueError, AttributeError):
                pass
        
        return cls(
            id=data.get("id", ""),
            name=data.get("name", ""),
            framework_type=fw_type,
            version=data.get("version", "1.0"),
            description=data.get("description", ""),
            controls_count=data.get("controls_count", 0),
            categories=data.get("categories", []),
            supported=data.get("supported", True),
            last_updated=last_updated,
        )


@dataclass
class Control:
    """
    Represents a compliance control.
    
    Attributes:
        id: Control identifier (e.g., "GV1", "CC6.2")
        name: Control name
        description: Control description
        category: Control category
        framework_id: Parent framework ID
        status: Current compliance status
        implementation: Implementation description
        evidence: Evidence of implementation
        last_assessed: Last assessment timestamp
        next_assessment: Next scheduled assessment
        findings: List of findings
    """
    id: str
    name: str
    category: str
    framework_id: str
    status: ComplianceStatus = ComplianceStatus.NOT_ASSESSED
    description: str = ""
    implementation: str = ""
    evidence: List[str] = field(default_factory=list)
    last_assessed: Optional[datetime] = None
    next_assessment: Optional[datetime] = None
    findings: List["Finding"] = field(default_factory=list)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> Control:
        """Create Control from dictionary."""
        def parse_datetime(value: Any) -> Optional[datetime]:
            if value is None:
                return None
            if isinstance(value, datetime):
                return value
            try:
                return datetime.fromisoformat(value.replace("Z", "+00:00"))
            except (ValueError, AttributeError):
                return None
        
        findings = []
        if "findings" in data:
            findings = [Finding.from_dict(f) for f in data["findings"]]
        
        return cls(
            id=data.get("id", ""),
            name=data.get("name", ""),
            category=data.get("category", ""),
            framework_id=data.get("framework_id", ""),
            status=ComplianceStatus.from_string(data.get("status", "not_assessed")),
            description=data.get("description", ""),
            implementation=data.get("implementation", ""),
            evidence=data.get("evidence", []),
            last_assessed=parse_datetime(data.get("last_assessed")),
            next_assessment=parse_datetime(data.get("next_assessment")),
            findings=findings,
        )


@dataclass
class Finding:
    """
    Represents a compliance finding.
    
    Attributes:
        id: Finding identifier
        control_id: Associated control ID
        severity: Finding severity
        title: Finding title
        description: Finding description
        recommendation: Remediation recommendation
        affected_resources: List of affected resources
        detected_at: When the finding was detected
        resolved_at: When the finding was resolved
        status: Finding status (open/resolved)
    """
    id: str
    control_id: str
    severity: str
    title: str
    description: str = ""
    recommendation: str = ""
    affected_resources: List[str] = field(default_factory=list)
    detected_at: Optional[datetime] = None
    resolved_at: Optional[datetime] = None
    status: str = "open"
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> Finding:
        """Create Finding from dictionary."""
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
            id=data.get("id", ""),
            control_id=data.get("control_id", ""),
            severity=data.get("severity", "medium"),
            title=data.get("title", ""),
            description=data.get("description", ""),
            recommendation=data.get("recommendation", ""),
            affected_resources=data.get("affected_resources", []),
            detected_at=parse_datetime(data.get("detected_at")),
            resolved_at=parse_datetime(data.get("resolved_at")),
            status=data.get("status", "open"),
        )


@dataclass
class ComplianceCheck:
    """
    Result of a compliance check.
    
    Attributes:
        framework_id: Framework being checked
        framework_name: Name of the framework
        timestamp: Check timestamp
        status: Overall compliance status
        controls_passed: Number of controls that passed
        controls_failed: Number of controls that failed
        controls_warning: Number of controls with warnings
        controls_not_applicable: Number of not applicable controls
        controls_total: Total number of controls assessed
        compliance_score: Overall compliance score (0-100)
        controls: List of control results
        findings: All findings from the check
        duration_ms: Check duration in milliseconds
    """
    framework_id: str
    framework_name: str
    timestamp: datetime
    status: ComplianceStatus
    controls_passed: int = 0
    controls_failed: int = 0
    controls_warning: int = 0
    controls_not_applicable: int = 0
    controls_total: int = 0
    compliance_score: float = 0.0
    controls: List[Control] = field(default_factory=list)
    findings: List[Finding] = field(default_factory=list)
    duration_ms: float = 0.0
    
    @property
    def is_compliant(self) -> bool:
        """Check if the compliance check passed overall."""
        return self.status == ComplianceStatus.PASS
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> ComplianceCheck:
        """Create ComplianceCheck from dictionary."""
        timestamp = data.get("timestamp")
        if isinstance(timestamp, str):
            try:
                timestamp = datetime.fromisoformat(timestamp.replace("Z", "+00:00"))
            except ValueError:
                timestamp = datetime.now()
        elif timestamp is None:
            timestamp = datetime.now()
        
        controls = []
        if "controls" in data:
            controls = [Control.from_dict(c) for c in data["controls"]]
        
        findings = []
        if "findings" in data:
            findings = [Finding.from_dict(f) for f in data["findings"]]
        
        return cls(
            framework_id=data.get("framework_id", ""),
            framework_name=data.get("framework_name", ""),
            timestamp=timestamp,
            status=ComplianceStatus.from_string(data.get("status", "not_assessed")),
            controls_passed=data.get("controls_passed", 0),
            controls_failed=data.get("controls_failed", 0),
            controls_warning=data.get("controls_warning", 0),
            controls_not_applicable=data.get("controls_not_applicable", 0),
            controls_total=data.get("controls_total", 0),
            compliance_score=data.get("compliance_score", 0.0),
            controls=controls,
            findings=findings,
            duration_ms=data.get("duration_ms", 0.0),
        )


@dataclass
class ComplianceReport:
    """
    Comprehensive compliance report for multiple frameworks.
    
    Attributes:
        id: Report identifier
        generated_at: Report generation timestamp
        period_start: Report period start
        period_end: Report period end
        frameworks: List of framework check results
        summary: Executive summary
        recommendations: List of recommendations
        next_scheduled: Next scheduled report
    """
    id: str
    generated_at: datetime
    period_start: datetime
    period_end: datetime
    frameworks: List[ComplianceCheck] = field(default_factory=list)
    summary: str = ""
    recommendations: List[str] = field(default_factory=list)
    next_scheduled: Optional[datetime] = None
    
    @property
    def overall_score(self) -> float:
        """Calculate overall compliance score across frameworks."""
        if not self.frameworks:
            return 0.0
        return sum(f.compliance_score for f in self.frameworks) / len(self.frameworks)
    
    @property
    def total_findings(self) -> int:
        """Count total findings across all frameworks."""
        return sum(len(f.findings) for f in self.frameworks)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> ComplianceReport:
        """Create ComplianceReport from dictionary."""
        def parse_datetime(value: Any) -> Optional[datetime]:
            if value is None:
                return None
            if isinstance(value, datetime):
                return value
            try:
                return datetime.fromisoformat(value.replace("Z", "+00:00"))
            except (ValueError, AttributeError):
                return None
        
        generated_at = parse_datetime(data.get("generated_at")) or datetime.now()
        period_start = parse_datetime(data.get("period_start")) or datetime.now()
        period_end = parse_datetime(data.get("period_end")) or datetime.now()
        
        frameworks = []
        if "frameworks" in data:
            frameworks = [ComplianceCheck.from_dict(f) for f in data["frameworks"]]
        elif "checks" in data:
            frameworks = [ComplianceCheck.from_dict(f) for f in data["checks"]]
        
        return cls(
            id=data.get("id", ""),
            generated_at=generated_at,
            period_start=period_start,
            period_end=period_end,
            frameworks=frameworks,
            summary=data.get("summary", ""),
            recommendations=data.get("recommendations", []),
            next_scheduled=parse_datetime(data.get("next_scheduled")),
        )
