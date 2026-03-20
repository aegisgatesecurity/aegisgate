# SPDX-License-Identifier: MIT
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.

"""
Core service for the AegisGate Python SDK.

Provides system health, version, licensing, and module management operations.
"""

from typing import Any, Dict, List, Optional

from aegisgate.connection import SyncConnection, AsyncConnection
from aegisgate.models import Health, Version, Module, License, LicenseType


class CoreService:
    """Synchronous core system service."""

    def __init__(self, connection: SyncConnection):
        self._conn = connection

    def health(self) -> Health:
        """
        Get system health status.

        Returns:
            Health object with system status
        """
        response = self._conn.get("/api/v1/health")
        return Health.from_dict(response)

    def health_detailed(self) -> Dict[str, Any]:
        """
        Get detailed system health information.

        Returns:
            Detailed health status including all components
        """
        response = self._conn.get("/api/v1/health/detailed")
        return response

    def version(self) -> Version:
        """
        Get system version information.

        Returns:
            Version object with build details
        """
        response = self._conn.get("/api/v1/version")
        return Version.from_dict(response)

    def get_license(self) -> License:
        """
        Get current license information.

        Returns:
            License object with license details
        """
        response = self._conn.get("/api/v1/license")
        return License.from_dict(response)

    def validate_license(self, license_key: str) -> Dict[str, Any]:
        """
        Validate a license key.

        Args:
            license_key: License key to validate

        Returns:
            Validation result
        """
        response = self._conn.post(
            "/api/v1/license/validate",
            json_data={"license_key": license_key},
        )
        return response

    def activate_license(self, license_key: str) -> License:
        """
        Activate a license key.

        Args:
            license_key: License key to activate

        Returns:
            Activated License object
        """
        response = self._conn.post(
            "/api/v1/license/activate",
            json_data={"license_key": license_key},
        )
        return License.from_dict(response)

    def deactivate_license(self) -> None:
        """
        Deactivate the current license.
        """
        self._conn.post("/api/v1/license/deactivate")

    def list_modules(self) -> List[Module]:
        """
        List all available modules.

        Returns:
            List of Module objects
        """
        response = self._conn.get("/api/v1/modules")
        modules = response.get("modules", [])
        return [Module.from_dict(m) for m in modules]

    def get_module(self, module_name: str) -> Module:
        """
        Get a specific module's information.

        Args:
            module_name: Module name

        Returns:
            Module object
        """
        response = self._conn.get(f"/api/v1/modules/{module_name}")
        return Module.from_dict(response)

    def enable_module(self, module_name: str) -> Module:
        """
        Enable a module.

        Args:
            module_name: Module name to enable

        Returns:
            Updated Module object
        """
        response = self._conn.post(f"/api/v1/modules/{module_name}/enable")
        return Module.from_dict(response)

    def disable_module(self, module_name: str) -> Module:
        """
        Disable a module.

        Args:
            module_name: Module name to disable

        Returns:
            Updated Module object
        """
        response = self._conn.post(f"/api/v1/modules/{module_name}/disable")
        return Module.from_dict(response)

    def get_config(self) -> Dict[str, Any]:
        """
        Get current system configuration.

        Returns:
            System configuration dictionary
        """
        response = self._conn.get("/api/v1/config")
        return response

    def update_config(self, config: Dict[str, Any]) -> Dict[str, Any]:
        """
        Update system configuration.

        Args:
            config: Configuration dictionary to update

        Returns:
            Updated configuration
        """
        response = self._conn.patch("/api/v1/config", json_data=config)
        return response

    def get_metrics(self) -> Dict[str, Any]:
        """
        Get system metrics.

        Returns:
            Metrics dictionary
        """
        response = self._conn.get("/api/v1/metrics")
        return response

    def get_logs(
        self,
        level: Optional[str] = None,
        module: Optional[str] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[Dict[str, Any]]:
        """
        Get system logs.

        Args:
            level: Filter by log level ('debug', 'info', 'warn', 'error')
            module: Filter by module name
            limit: Maximum number of logs
            offset: Pagination offset

        Returns:
            List of log entries
        """
        params = {"limit": limit, "offset": offset}
        if level:
            params["level"] = level
        if module:
            params["module"] = module

        response = self._conn.get("/api/v1/logs", params=params)
        return response.get("logs", [])

    def get_status(self) -> Dict[str, Any]:
        """
        Get overall system status summary.

        Returns:
            Status summary including health, metrics, and alerts
        """
        response = self._conn.get("/api/v1/status")
        return response

    def restart_service(self, service_name: str) -> Dict[str, Any]:
        """
        Restart a specific service.

        Args:
            service_name: Name of the service to restart

        Returns:
            Restart result
        """
        response = self._conn.post(f"/api/v1/services/{service_name}/restart")
        return response

    def get_environment(self) -> Dict[str, Any]:
        """
        Get environment information.

        Returns:
            Environment details (platform, go version, build info)
        """
        response = self._conn.get("/api/v1/environment")
        return response


class AsyncCoreService:
    """Asynchronous core system service."""

    def __init__(self, connection: AsyncConnection):
        self._conn = connection

    async def health(self) -> Health:
        """
        Get system health status.

        Returns:
            Health object with system status
        """
        response = await self._conn.get("/api/v1/health")
        return Health.from_dict(response)

    async def health_detailed(self) -> Dict[str, Any]:
        """
        Get detailed system health information.

        Returns:
            Detailed health status including all components
        """
        response = await self._conn.get("/api/v1/health/detailed")
        return response

    async def version(self) -> Version:
        """
        Get system version information.

        Returns:
            Version object with build details
        """
        response = await self._conn.get("/api/v1/version")
        return Version.from_dict(response)

    async def get_license(self) -> License:
        """
        Get current license information.

        Returns:
            License object with license details
        """
        response = await self._conn.get("/api/v1/license")
        return License.from_dict(response)

    async def validate_license(self, license_key: str) -> Dict[str, Any]:
        """
        Validate a license key.

        Args:
            license_key: License key to validate

        Returns:
            Validation result
        """
        response = await self._conn.post(
            "/api/v1/license/validate",
            json_data={"license_key": license_key},
        )
        return response

    async def activate_license(self, license_key: str) -> License:
        """
        Activate a license key.

        Args:
            license_key: License key to activate

        Returns:
            Activated License object
        """
        response = await self._conn.post(
            "/api/v1/license/activate",
            json_data={"license_key": license_key},
        )
        return License.from_dict(response)

    async def deactivate_license(self) -> None:
        """
        Deactivate the current license.
        """
        await self._conn.post("/api/v1/license/deactivate")

    async def list_modules(self) -> List[Module]:
        """
        List all available modules.

        Returns:
            List of Module objects
        """
        response = await self._conn.get("/api/v1/modules")
        modules = response.get("modules", [])
        return [Module.from_dict(m) for m in modules]

    async def get_module(self, module_name: str) -> Module:
        """
        Get a specific module's information.

        Args:
            module_name: Module name

        Returns:
            Module object
        """
        response = await self._conn.get(f"/api/v1/modules/{module_name}")
        return Module.from_dict(response)

    async def enable_module(self, module_name: str) -> Module:
        """
        Enable a module.

        Args:
            module_name: Module name to enable

        Returns:
            Updated Module object
        """
        response = await self._conn.post(f"/api/v1/modules/{module_name}/enable")
        return Module.from_dict(response)

    async def disable_module(self, module_name: str) -> Module:
        """
        Disable a module.

        Args:
            module_name: Module name to disable

        Returns:
            Updated Module object
        """
        response = await self._conn.post(f"/api/v1/modules/{module_name}/disable")
        return Module.from_dict(response)

    async def get_config(self) -> Dict[str, Any]:
        """
        Get current system configuration.

        Returns:
            System configuration dictionary
        """
        response = await self._conn.get("/api/v1/config")
        return response

    async def update_config(self, config: Dict[str, Any]) -> Dict[str, Any]:
        """
        Update system configuration.

        Args:
            config: Configuration dictionary to update

        Returns:
            Updated configuration
        """
        response = await self._conn.patch("/api/v1/config", json_data=config)
        return response

    async def get_metrics(self) -> Dict[str, Any]:
        """
        Get system metrics.

        Returns:
            Metrics dictionary
        """
        response = await self._conn.get("/api/v1/metrics")
        return response

    async def get_logs(
        self,
        level: Optional[str] = None,
        module: Optional[str] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[Dict[str, Any]]:
        """
        Get system logs.

        Args:
            level: Filter by log level
            module: Filter by module name
            limit: Maximum number of logs
            offset: Pagination offset

        Returns:
            List of log entries
        """
        params = {"limit": limit, "offset": offset}
        if level:
            params["level"] = level
        if module:
            params["module"] = module

        response = await self._conn.get("/api/v1/logs", params=params)
        return response.get("logs", [])

    async def get_status(self) -> Dict[str, Any]:
        """
        Get overall system status summary.

        Returns:
            Status summary including health, metrics, and alerts
        """
        response = await self._conn.get("/api/v1/status")
        return response

    async def restart_service(self, service_name: str) -> Dict[str, Any]:
        """
        Restart a specific service.

        Args:
            service_name: Name of the service to restart

        Returns:
            Restart result
        """
        response = await self._conn.post(f"/api/v1/services/{service_name}/restart")
        return response

    async def get_environment(self) -> Dict[str, Any]:
        """
        Get environment information.

        Returns:
            Environment details (platform, go version, build info)
        """
        response = await self._conn.get("/api/v1/environment")
        return response