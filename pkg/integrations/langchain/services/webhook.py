# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Webhook service for the AegisGate Python SDK.

This module provides webhook management functionality.
"""

from __future__ import annotations

from typing import Optional, List, Dict, Any
import logging

from aegisgate.models.webhook import (
    Webhook,
    WebhookCreate,
    WebhookUpdate,
    WebhookEvent,
    WebhookEventType,
    WebhookDelivery,
)
from aegisgate.exceptions import ResourceNotFoundError

logger = logging.getLogger(__name__)


class WebhookService:
    """
    Service for webhook management.
    
    This service handles:
    - Webhook CRUD operations
    - Webhook testing
    - Delivery history
    - Event subscriptions
    
    Example:
        >>> client = AegisGateClient(base_url="http://localhost:8080")
        >>> webhooks = client.webhooks
        >>> 
        >>> # Create a webhook
        >>> webhook = webhooks.create(WebhookCreate(
        ...     name="Security Alerts",
        ...     url="https://my-app.example.com/webhook",
        ...     events=[WebhookEventType.VIOLATION_BLOCKED]
        ... ))
        >>> 
        >>> # List webhooks
        >>> all_webhooks = webhooks.list()
        >>> for wh in all_webhooks:
        ...     print(f"{wh.name}: {wh.url}")
    """
    
    def __init__(self, client: "AegisGateClient"):
        """
        Initialize the webhook service.
        
        Args:
            client: The AegisGate client instance
        """
        self._client = client
    
    # =====================================================================
    # CRUD Operations
    # =====================================================================
    
    def list(self) -> List[Webhook]:
        """
        List all configured webhooks.
        
        Returns:
            List of Webhook objects
        
        Example:
            >>> webhooks = webhooks.list()
        """
        response = self._client._request("GET", "webhooks")
        
        result = []
        if "webhooks" in response:
            result = [Webhook.from_dict(w) for w in response["webhooks"]]
        elif "data" in response:
            result = [Webhook.from_dict(w) for w in response["data"]]
        
        return result
    
    def get(self, webhook_id: str) -> Webhook:
        """
        Get a webhook by ID.
        
        Args:
            webhook_id: Webhook identifier
        
        Returns:
            Webhook object
        
        Raises:
            ResourceNotFoundError: If webhook doesn't exist
        
        Example:
            >>> webhook = webhooks.get("webhook-123")
        """
        response = self._client._request("GET", f"webhooks/{webhook_id}")
        return Webhook.from_dict(response)
    
    def create(self, webhook: WebhookCreate) -> Webhook:
        """
        Create a new webhook.
        
        Args:
            webhook: Webhook creation data
        
        Returns:
            Created Webhook object
        
        Example:
            >>> webhook = webhooks.create(WebhookCreate(
            ...     name="My Webhook",
            ...     url="https://example.com/webhook",
            ...     events=[WebhookEventType.VIOLATION_BLOCKED]
            ... ))
        """
        response = self._client._request(
            "POST", "webhooks", json=webhook.to_dict()
        )
        return Webhook.from_dict(response)
    
    def update(self, webhook_id: str, updates: WebhookUpdate) -> Webhook:
        """
        Update an existing webhook.
        
        Args:
            webhook_id: Webhook identifier
            updates: Fields to update
        
        Returns:
            Updated Webhook object
        
        Raises:
            ResourceNotFoundError: If webhook doesn't exist
        
        Example:
            >>> webhook = webhooks.update(
            ...     "webhook-123",
            ...     WebhookUpdate(enabled=False)
            ... )
        """
        response = self._client._request(
            "PUT", f"webhooks/{webhook_id}", json=updates.to_dict()
        )
        return Webhook.from_dict(response)
    
    def delete(self, webhook_id: str) -> None:
        """
        Delete a webhook.
        
        Args:
            webhook_id: Webhook identifier
        
        Raises:
            ResourceNotFoundError: If webhook doesn't exist
        
        Example:
            >>> webhooks.delete("webhook-123")
        """
        self._client._request("DELETE", f"webhooks/{webhook_id}")
    
    # =====================================================================
    # Webhook Control
    # =====================================================================
    
    def enable(self, webhook_id: str) -> Webhook:
        """
        Enable a webhook.
        
        Args:
            webhook_id: Webhook identifier
        
        Returns:
            Updated Webhook object
        
        Example:
            >>> webhook = webhooks.enable("webhook-123")
        """
        response = self._client._request("POST", f"webhooks/{webhook_id}/enable")
        return Webhook.from_dict(response)
    
    def disable(self, webhook_id: str) -> Webhook:
        """
        Disable a webhook.
        
        Args:
            webhook_id: Webhook identifier
        
        Returns:
            Updated Webhook object
        
        Example:
            >>> webhook = webhooks.disable("webhook-123")
        """
        response = self._client._request("POST", f"webhooks/{webhook_id}/disable")
        return Webhook.from_dict(response)
    
    # =====================================================================
    # Testing
    # =====================================================================
    
    def test(self, webhook_id: str) -> Dict[str, Any]:
        """
        Send a test event to a webhook.
        
        Args:
            webhook_id: Webhook identifier
        
        Returns:
            Test result with delivery status
        
        Raises:
            ResourceNotFoundError: If webhook doesn't exist
        
        Example:
            >>> result = webhooks.test("webhook-123")
            >>> print(f"Status: {result['status']}")
            >>> print(f"Response: {result['response_code']}")
        """
        response = self._client._request("POST", f"webhooks/{webhook_id}/test")
        return response
    
    def test_with_payload(
        self,
        webhook_id: str,
        event_type: WebhookEventType,
        payload: Dict[str, Any],
    ) -> Dict[str, Any]:
        """
        Send a test event with custom payload to a webhook.
        
        Args:
            webhook_id: Webhook identifier
            event_type: Type of event to simulate
            payload: Custom payload data
        
        Returns:
            Test result with delivery status
        
        Example:
            >>> result = webhooks.test_with_payload(
            ...     "webhook-123",
            ...     WebhookEventType.VIOLATION_BLOCKED,
            ...     {"threat_type": "prompt_injection", "blocked": True}
            ... )
        """
        data = {
            "type": event_type.value,
            "payload": payload,
        }
        response = self._client._request(
            "POST", f"webhooks/{webhook_id}/test", json=data
        )
        return response
    
    # =====================================================================
    # Delivery History
    # =====================================================================
    
    def list_deliveries(
        self,
        webhook_id: Optional[str] = None,
        status: Optional[str] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[WebhookDelivery]:
        """
        List webhook delivery attempts.
        
        Args:
            webhook_id: Optional webhook filter
            status: Optional status filter (success, failed, pending)
            limit: Maximum results
            offset: Pagination offset
        
        Returns:
            List of WebhookDelivery objects
        
        Example:
            >>> deliveries = webhooks.list_deliveries(status="failed")
            >>> for d in deliveries:
            ...     print(f"{d.id}: {d.error}")
        """
        params = {"limit": limit, "offset": offset}
        if webhook_id:
            params["webhook_id"] = webhook_id
        if status:
            params["status"] = status
        
        response = self._client._request("GET", "webhooks/deliveries", params=params)
        
        deliveries = []
        if "deliveries" in response:
            deliveries = [WebhookDelivery.from_dict(d) for d in response["deliveries"]]
        elif "data" in response:
            deliveries = [WebhookDelivery.from_dict(d) for d in response["data"]]
        
        return deliveries
    
    def get_delivery(self, delivery_id: str) -> WebhookDelivery:
        """
        Get a specific delivery attempt.
        
        Args:
            delivery_id: Delivery identifier
        
        Returns:
            WebhookDelivery object
        
        Example:
            >>> delivery = webhooks.get_delivery("delivery-123")
        """
        response = self._client._request("GET", f"webhooks/deliveries/{delivery_id}")
        return WebhookDelivery.from_dict(response)
    
    def retry_delivery(self, delivery_id: str) -> Dict[str, Any]:
        """
        Retry a failed delivery.
        
        Args:
            delivery_id: Delivery identifier
        
        Returns:
            Retry result
        
        Example:
            >>> result = webhooks.retry_delivery("delivery-123")
        """
        response = self._client._request(
            "POST", f"webhooks/deliveries/{delivery_id}/retry"
        )
        return response
    
    # =====================================================================
    # Events
    # =====================================================================
    
    def list_event_types(self) -> List[Dict[str, Any]]:
        """
        List all available webhook event types.
        
        Returns:
            List of event type definitions
        
        Example:
            >>> events = webhooks.list_event_types()
            >>> for e in events:
            ...     print(f"{e['type']}: {e['description']}")
        """
        response = self._client._request("GET", "webhooks/events")
        return response.get("event_types", [])
    
    def replay_events(
        self,
        webhook_id: str,
        start_time: str,
        end_time: str,
    ) -> Dict[str, Any]:
        """
        Replay webhook events for a time range.
        
        Args:
            webhook_id: Webhook identifier
            start_time: Start time (ISO format)
            end_time: End time (ISO format)
        
        Returns:
            Replay result with count
        
        Example:
            >>> result = webhooks.replay_events(
            ...     "webhook-123",
            ...     "2025-01-01T00:00:00Z",
            ...     "2025-01-02T00:00:00Z"
            ... )
        """
        data = {
            "start_time": start_time,
            "end_time": end_time,
        }
        response = self._client._request(
            "POST", f"webhooks/{webhook_id}/replay", json=data
        )
        return response


# Type hint reference - import only for type checking
if False:
    from aegisgate.client import AegisGateClient
