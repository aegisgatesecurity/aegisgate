# SPDX-License-Identifier: MIT
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.

"""
Webhook management service for the AegisGate Python SDK.

Provides webhook configuration and management for event notifications.
"""

from typing import Any, Dict, List, Optional
from datetime import datetime

from aegisgate.connection import SyncConnection, AsyncConnection
from aegisgate.models import Webhook


class WebhookService:
    """Synchronous webhook management service."""

    def __init__(self, connection: SyncConnection):
        self._conn = connection

    def list_webhooks(self) -> List[Webhook]:
        """
        List all webhooks.

        Returns:
            List of Webhook objects
        """
        response = self._conn.get("/api/v1/webhooks")
        webhooks = response.get("webhooks", [])
        return [Webhook.from_dict(w) for w in webhooks]

    def get_webhook(self, webhook_id: str) -> Webhook:
        """
        Get a specific webhook.

        Args:
            webhook_id: Webhook identifier

        Returns:
            Webhook object
        """
        response = self._conn.get(f"/api/v1/webhooks/{webhook_id}")
        return Webhook.from_dict(response)

    def create_webhook(
        self,
        url: str,
        events: List[str],
        secret: Optional[str] = None,
        enabled: bool = True,
    ) -> Webhook:
        """
        Create a new webhook.

        Args:
            url: Webhook URL endpoint
            events: List of event types to trigger the webhook
            secret: Secret for webhook signature verification
            enabled: Whether the webhook is enabled

        Returns:
            Created Webhook object
        """
        payload = {
            "url": url,
            "events": events,
            "enabled": enabled,
        }
        if secret:
            payload["secret"] = secret

        response = self._conn.post("/api/v1/webhooks", json_data=payload)
        return Webhook.from_dict(response)

    def update_webhook(
        self,
        webhook_id: str,
        **kwargs,
    ) -> Webhook:
        """
        Update a webhook.

        Args:
            webhook_id: Webhook identifier
            **kwargs: Fields to update (url, events, secret, enabled)

        Returns:
            Updated Webhook object
        """
        response = self._conn.patch(
            f"/api/v1/webhooks/{webhook_id}",
            json_data=kwargs,
        )
        return Webhook.from_dict(response)

    def delete_webhook(self, webhook_id: str) -> None:
        """
        Delete a webhook.

        Args:
            webhook_id: Webhook identifier
        """
        self._conn.delete(f"/api/v1/webhooks/{webhook_id}")

    def test_webhook(self, webhook_id: str) -> Dict[str, Any]:
        """
        Test a webhook by sending a test event.

        Args:
            webhook_id: Webhook identifier

        Returns:
            Test result with delivery status
        """
        response = self._conn.post(f"/api/v1/webhooks/{webhook_id}/test")
        return response

    def get_deliveries(
        self,
        webhook_id: Optional[str] = None,
        status: Optional[str] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[Dict[str, Any]]:
        """
        Get webhook delivery history.

        Args:
            webhook_id: Filter by webhook (optional)
            status: Filter by delivery status ('success', 'failed', 'pending')
            limit: Maximum number of deliveries
            offset: Pagination offset

        Returns:
            List of delivery records
        """
        params = {"limit": limit, "offset": offset}
        if webhook_id:
            params["webhook_id"] = webhook_id
        if status:
            params["status"] = status

        response = self._conn.get("/api/v1/webhooks/deliveries", params=params)
        return response.get("deliveries", [])

    def retry_delivery(self, delivery_id: str) -> Dict[str, Any]:
        """
        Retry a failed webhook delivery.

        Args:
            delivery_id: Delivery identifier

        Returns:
            Retry result
        """
        response = self._conn.post(f"/api/v1/webhooks/deliveries/{delivery_id}/retry")
        return response

    def get_event_types(self) -> List[Dict[str, Any]]:
        """
        Get available webhook event types.

        Returns:
            List of event type definitions
        """
        response = self._conn.get("/api/v1/webhooks/event-types")
        return response.get("event_types", [])

    def ping_webhook(self, webhook_id: str) -> Dict[str, Any]:
        """
        Ping a webhook to check connectivity.

        Args:
            webhook_id: Webhook identifier

        Returns:
            Ping result with latency
        """
        response = self._conn.post(f"/api/v1/webhooks/{webhook_id}/ping")
        return response


class AsyncWebhookService:
    """Asynchronous webhook management service."""

    def __init__(self, connection: AsyncConnection):
        self._conn = connection

    async def list_webhooks(self) -> List[Webhook]:
        """
        List all webhooks.

        Returns:
            List of Webhook objects
        """
        response = await self._conn.get("/api/v1/webhooks")
        webhooks = response.get("webhooks", [])
        return [Webhook.from_dict(w) for w in webhooks]

    async def get_webhook(self, webhook_id: str) -> Webhook:
        """
        Get a specific webhook.

        Args:
            webhook_id: Webhook identifier

        Returns:
            Webhook object
        """
        response = await self._conn.get(f"/api/v1/webhooks/{webhook_id}")
        return Webhook.from_dict(response)

    async def create_webhook(
        self,
        url: str,
        events: List[str],
        secret: Optional[str] = None,
        enabled: bool = True,
    ) -> Webhook:
        """
        Create a new webhook.

        Args:
            url: Webhook URL endpoint
            events: List of event types to trigger the webhook
            secret: Secret for webhook signature verification
            enabled: Whether the webhook is enabled

        Returns:
            Created Webhook object
        """
        payload = {
            "url": url,
            "events": events,
            "enabled": enabled,
        }
        if secret:
            payload["secret"] = secret

        response = await self._conn.post("/api/v1/webhooks", json_data=payload)
        return Webhook.from_dict(response)

    async def update_webhook(
        self,
        webhook_id: str,
        **kwargs,
    ) -> Webhook:
        """
        Update a webhook.

        Args:
            webhook_id: Webhook identifier
            **kwargs: Fields to update

        Returns:
            Updated Webhook object
        """
        response = await self._conn.patch(
            f"/api/v1/webhooks/{webhook_id}",
            json_data=kwargs,
        )
        return Webhook.from_dict(response)

    async def delete_webhook(self, webhook_id: str) -> None:
        """
        Delete a webhook.

        Args:
            webhook_id: Webhook identifier
        """
        await self._conn.delete(f"/api/v1/webhooks/{webhook_id}")

    async def test_webhook(self, webhook_id: str) -> Dict[str, Any]:
        """
        Test a webhook by sending a test event.

        Args:
            webhook_id: Webhook identifier

        Returns:
            Test result with delivery status
        """
        response = await self._conn.post(f"/api/v1/webhooks/{webhook_id}/test")
        return response

    async def get_deliveries(
        self,
        webhook_id: Optional[str] = None,
        status: Optional[str] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[Dict[str, Any]]:
        """
        Get webhook delivery history.

        Args:
            webhook_id: Filter by webhook (optional)
            status: Filter by delivery status
            limit: Maximum number of deliveries
            offset: Pagination offset

        Returns:
            List of delivery records
        """
        params = {"limit": limit, "offset": offset}
        if webhook_id:
            params["webhook_id"] = webhook_id
        if status:
            params["status"] = status

        response = await self._conn.get("/api/v1/webhooks/deliveries", params=params)
        return response.get("deliveries", [])

    async def retry_delivery(self, delivery_id: str) -> Dict[str, Any]:
        """
        Retry a failed webhook delivery.

        Args:
            delivery_id: Delivery identifier

        Returns:
            Retry result
        """
        response = await self._conn.post(f"/api/v1/webhooks/deliveries/{delivery_id}/retry")
        return response

    async def get_event_types(self) -> List[Dict[str, Any]]:
        """
        Get available webhook event types.

        Returns:
            List of event type definitions
        """
        response = await self._conn.get("/api/v1/webhooks/event-types")
        return response.get("event_types", [])

    async def ping_webhook(self, webhook_id: str) -> Dict[str, Any]:
        """
        Ping a webhook to check connectivity.

        Args:
            webhook_id: Webhook identifier

        Returns:
            Ping result with latency
        """
        response = await self._conn.post(f"/api/v1/webhooks/{webhook_id}/ping")
        return response