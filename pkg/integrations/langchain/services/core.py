# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Core service for the AegisGate Python SDK.

This module provides core system operations like health checks, 
version information, and module management.
"""

from __future__ import annotations

from typing import Optional, List, Dict, Any
from dataclasses import dataclass
from datetime import datetime
import logging

from aegisgate.exceptions import ResourceNotFoundError

logger = logging.getLogger(__name__)


@dataclass
class Health:
    """
    System health status.
    
    Attributes:
        status: Overall health status (healthy/degraded/unhealthy)
        checks: Individual health check results
        timestamp: Health check timestamp
        version: Server version
        uptime_seconds: Server uptime in seconds
    """
    status: str
    checks: Dict[str, Any]
    timestamp: datetime
    version: Optional[str] = None
    uptime_seconds: Optional[int] = None
    
    @property
    def is_healthy(self) -> bool:
        """Check if system is healthy."""
        return self.status == "healthy"
    
    @property
    def is_degraded(self) -> bool:
        """Check if system is degraded."""
        return self.status == "degraded"
    
    @property
    def is_unhealthy(self) -> bool:
        """Check if system is unhealthy."""
        return self.status == "unhealthy"
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> Health:
        """Create Health from dictionary."""
        timestamp = data.get("timestamp")
        if isinstance(timestamp, str):
            try:
                timestamp = datetime.fromisoformat(timestamp.replace("Z", "+00:00"))
            except ValueError:
                timestamp = datetime.now()
        elif timestamp is None:
            timestamp = datetime.now()
        
        return cls(
            status=data.get("status", "unknown"),
            checks=data.get("checks", {}),
            timestamp=timestamp,
            version=data.get("version"),
            uptime_seconds=data.get("uptime_seconds"),
        )


@dataclass
class Version:
    """
    Version information.
    
    Attributes:
        version: Semantic version string
        major: Major version number
        minor: Minor version number
        patch: Patch version number
        build_time: Build timestamp
        git_commit: Git commit hash
        go_version: Go version used for build
        platform: Build platform
    """
    version: str
    major: int
    minor: int
    patch: int
    build_time: Optional[str] = None
    git_commit: Optional[str] = None
    go_version: Optional[str] = None
    platform: Optional[str] = None
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> Version:
        """Create Version from dictionary."""
        version_str = data.get("version", "0.0.0")
        parts = version_str.split(".")
        major = int(parts[0]) if len(parts) > 0 else 0
        minor = int(parts[1]) if len(parts) > 1 else 0
        patch = int(parts[2]) if len(parts) > 2 else 0
        
        return cls(
            version=version_str,
            major=major,
            minor=minor,
            patch=patch,
            build_time=data.get("build_time"),
            git_commit=data.get("git_commit"),
            go_version=data.get("go_version"),
            platform=data.get("platform"),
        )


@dataclass
class Module:
    """
    System module information.
    
    Attributes:
        id: Module identifier
        name: Module name
        version: Module version
        description: Module description
        category: Module category
        status: Module status (enabled/disabled)
        enabled: Whether module is enabled
    """
    id: str
    name: str
    version: str
    description: str
    category: str
    status: str
    enabled: bool = True
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> Module:
        """Create Module from dictionary."""
        return cls(
            id=data.get("id", ""),
            name=data.get("name", ""),
            version=data.get("version", ""),
            description=data.get("description", ""),
            category=data.get("category", ""),
            status=data.get("status", "unknown"),
            enabled=data.get("enabled", True),
        )


class CoreService:
    """
    Service for core system operations.
    
    This service handles:
    - Health checks
    - Version information
    - Module management
    - System configuration
    
    Example:
        >>> client = AegisGateClient(base_url="http://localhost:8080")
        >>> core = client.core
        >>> 
        >>> # Health check
        >>> health = core.health()
        >>> print(f"Status: {health.status}")
        >>> 
        >>> # Version info
        >>> version = core.version()
        >>> print(f"Version: {version.version}")
    """
    
    def __init__(self, client: "AegisGateClient"):
        """
        Initialize the core service.
        
        Args:
            client: The AegisGate client instance
        """
        self._client = client
    
    # =====================================================================
    # Health
    # =====================================================================
    
    def health(self) -> Health:
        """
        Get system health status.
        
        Returns:
            Health object
        
        Example:
            >>> health = core.health()
            >>> if health.is_healthy:
            ...     print("System is healthy!")
        """
        response = self._client._request("GET", "health")
        return Health.from_dict(response)
    
    def ready(self) -> bool:
        """
        Check if the system is ready to accept requests.
        
        Returns:
            True if system is ready
        
        Example:
            >>> if core.ready():
            ...     print("Ready to process requests")
        """
        response = self._client._request("GET", "ready")
        return response.get("ready", False)
    
    def live(self) -> bool:
        """
        Check if the system is alive.
        
        This is a lighter check than health() - it only verifies
        that the server is responding.
        
        Returns:
            True if system is alive
        
        Example:
            >>> if core.live():
            ...     print("Server is responding")
        """
        response = self._client._request("GET", "live")
        return response.get("alive", True)
    
    # =====================================================================
    # Version
    # =====================================================================
    
    def version(self) -> Version:
        """
        Get server version information.
        
        Returns:
            Version object
        
        Example:
            >>> version = core.version()
            >>> print(f"AegisGate {version.version}")
        """
        response = self._client._request("GET", "version")
        return Version.from_dict(response)
    
    def check_version(self) -> Dict[str, Any]:
        """
        Check if client and server versions are compatible.
        
        Returns:
            Dictionary with compatibility info
        
        Example:
            >>> info = core.check_version()
            >>> if not info["compatible"]:
            ...     print(f"Warning: {info['message']}")
        """
        from aegisgate._version import __version__
        
        response = self._client._request(
            "GET", "version/check",
            params={"client_version": __version__}
        )
        return response
    
    # =====================================================================
    # Modules
    # =====================================================================
    
    def list_modules(self) -> List[Module]:
        """
        List all available modules.
        
        Returns:
            List of Module objects
        
        Example:
            >>> modules = core.list_modules()
            >>> for mod in modules:
            ...     print(f"{mod.name}: {mod.status}")
        """
        response = self._client._request("GET", "modules")
        
        modules = []
        if "modules" in response:
            modules = [Module.from_dict(m) for m in response["modules"]]
        elif "data" in response:
            modules = [Module.from_dict(m) for m in response["data"]]
        
        return modules
    
    def get_module(self, module_id: str) -> Module:
        """
        Get a specific module.
        
        Args:
            module_id: Module identifier
        
        Returns:
            Module object
        
        Raises:
            ResourceNotFoundError: If module doesn't exist
        
        Example:
            >>> module = core.get_module("atlas")
            >>> print(f"Module: {module.name}")
        """
        response = self._client._request("GET", f"modules/{module_id}")
        return Module.from_dict(response)
    
    def enable_module(self, module_id: str) -> Module:
        """
        Enable a module.
        
        Args:
            module_id: Module identifier
        
        Returns:
            Updated Module object
        
        Example:
            >>> module = core.enable_module("atlas")
        """
        response = self._client._request("POST", f"modules/{module_id}/enable")
        return Module.from_dict(response)
    
    def disable_module(self, module_id: str) -> Module:
        """
        Disable a module.
        
        Args:
            module_id: Module identifier
        
        Returns:
            Updated Module object
        
        Example:
            >>> module = core.disable_module("atlas")
        """
        response = self._client._request("POST", f"modules/{module_id}/disable")
        return Module.from_dict(response)
    
    # =====================================================================
    # Configuration
    # =====================================================================
    
    def get_config(self) -> Dict[str, Any]:
        """
        Get current server configuration.
        
        Returns:
            Configuration dictionary
        
        Example:
            >>> config = core.get_config()
            >>> print(f"Log level: {config['logging']['level']}")
        """
        return self._client._request("GET", "config")
    
    def update_config(self, updates: Dict[str, Any]) -> Dict[str, Any]:
        """
        Update server configuration.
        
        Args:
            updates: Configuration updates
        
        Returns:
            Updated configuration
        
        Example:
            >>> core.update_config({"logging": {"level": "debug"}})
        """
        return self._client._request("PUT", "config", json=updates)
    
    def reload_config(self) -> None:
        """
        Reload configuration from file.
        
        Example:
            >>> core.reload_config()
        """
        self._client._request("POST", "config/reload")
    
    # =====================================================================
    # System Info
    # =====================================================================
    
    def info(self) -> Dict[str, Any]:
        """
        Get system information.
        
        Returns:
            System information dictionary
        
        Example:
            >>> info = core.info()
            >>> print(f"Hostname: {info['hostname']}")
        """
        return self._client._request("GET", "info")
    
    def stats(self) -> Dict[str, Any]:
        """
        Get system statistics.
        
        Returns:
            System statistics dictionary
        
        Example:
            >>> stats = core.stats()
            >>> print(f"Memory: {stats['memory']['used']}")
        """
        return self._client._request("GET", "stats")
    
    def metrics(self) -> str:
        """
        Get Prometheus-format metrics.
        
        Returns:
            Metrics in Prometheus text format
        
        Example:
            >>> metrics = core.metrics()
            >>> with open("metrics.txt", "w") as f:
            ...     f.write(metrics)
        """
        response = self._client._request("GET", "metrics", raw=True)
        if isinstance(response, bytes):
            return response.decode("utf-8")
        return str(response)


# Type hint reference - import only for type checking
if False:
    from aegisgate.client import AegisGateClient
