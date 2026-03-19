# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Authentication service for the AegisGate Python SDK.

This module provides authentication and user management functionality.
"""

from __future__ import annotations

from typing import Optional, List, Dict, Any
import logging

from aegisgate.models.auth import (
    User,
    UserCreate,
    UserUpdate,
    LoginResult,
    Session,
    SessionInfo,
    Role,
    Provider,
)
from aegisgate.exceptions import (
    AuthenticationError,
    AuthorizationError,
    ValidationError,
    ResourceNotFoundError,
)

logger = logging.getLogger(__name__)


class AuthService:
    """
    Service for authentication and user management.
    
    This service handles:
    - User login/logout
    - Session management
    - User CRUD operations
    - Role and permission management
    
    Example:
        >>> client = AegisGateClient(base_url="http://localhost:8080")
        >>> auth = client.auth
        >>> 
        >>> # Login
        >>> result = auth.login("admin@example.com", "password123")
        >>> print(f"Logged in: {result.user.name}")
        >>> 
        >>> # List users
        >>> users = auth.list_users()
        >>> for user in users:
        ...     print(f"User: {user.email} ({user.role.value})")
    """
    
    def __init__(self, client: "AegisGateClient"):
        """
        Initialize the auth service.
        
        Args:
            client: The AegisGate client instance
        """
        self._client = client
    
    def login(self, username: str, password: str) -> LoginResult:
        """
        Authenticate a user with username and password.
        
        Args:
            username: User's email or username
            password: User's password
        
        Returns:
            LoginResult containing the authentication token
        
        Raises:
            AuthenticationError: If credentials are invalid
        
        Example:
            >>> result = auth.login("user@example.com", "password123")
            >>> print(f"Token: {result.token}")
        """
        data = {
            "username": username,
            "password": password,
        }
        
        response = self._client._request("POST", "auth/login", json=data)
        
        result = LoginResult.from_dict(response)
        
        if not result.success:
            raise AuthenticationError(
                message=result.error or "Authentication failed",
                code="AUTH_FAILED",
            )
        
        # Store token in client
        if result.token:
            self._client._set_token(result.token)
        
        return result
    
    def logout(self) -> None:
        """
        Log out the current user and invalidate the session.
        
        Raises:
            AuthenticationError: If logout fails
        
        Example:
            >>> auth.logout()
        """
        try:
            self._client._request("POST", "auth/logout")
        except Exception as e:
            logger.warning(f"Logout request failed: {e}")
        finally:
            self._client._clear_token()
    
    def refresh_token(self) -> LoginResult:
        """
        Refresh the current authentication token.
        
        Returns:
            LoginResult with the new token
        
        Raises:
            AuthenticationError: If refresh fails
        """
        response = self._client._request("POST", "auth/refresh")
        result = LoginResult.from_dict(response)
        
        if not result.success:
            raise AuthenticationError(
                message=result.error or "Token refresh failed",
                code="REFRESH_FAILED",
            )
        
        if result.token:
            self._client._set_token(result.token)
        
        return result
    
    # =====================================================================
    # User Management
    # =====================================================================
    
    def list_users(
        self,
        role: Optional[Role] = None,
        provider: Optional[Provider] = None,
        status: Optional[str] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[User]:
        """
        List all users with optional filtering.
        
        Args:
            role: Filter by user role
            provider: Filter by authentication provider
            status: Filter by account status
            limit: Maximum number of results
            offset: Pagination offset
        
        Returns:
            List of User objects
        
        Raises:
            AuthorizationError: If user lacks admin privileges
        
        Example:
            >>> users = auth.list_users(role=Role.ADMIN)
            >>> for user in users:
            ...     print(f"Admin: {user.email}")
        """
        params = {
            "limit": limit,
            "offset": offset,
        }
        if role:
            params["role"] = role.value
        if provider:
            params["provider"] = provider.value
        if status:
            params["status"] = status
        
        response = self._client._request("GET", "users", params=params)
        
        users = []
        if "users" in response:
            users = [User.from_dict(u) for u in response["users"]]
        elif "data" in response:
            users = [User.from_dict(u) for u in response["data"]]
        
        return users
    
    def get_user(self, user_id: str) -> User:
        """
        Get a user by ID.
        
        Args:
            user_id: User identifier
        
        Returns:
            User object
        
        Raises:
            ResourceNotFoundError: If user doesn't exist
        
        Example:
            >>> user = auth.get_user("user-123")
            >>> print(f"User: {user.name}")
        """
        response = self._client._request("GET", f"users/{user_id}")
        return User.from_dict(response)
    
    def create_user(self, user_data: UserCreate) -> User:
        """
        Create a new user.
        
        Args:
            user_data: User creation data
        
        Returns:
            Created User object
        
        Raises:
            AuthorizationError: If user lacks admin privileges
            ValidationError: If user data is invalid
        
        Example:
            >>> new_user = auth.create_user(UserCreate(
            ...     email="newuser@example.com",
            ...     name="New User",
            ...     password="SecurePass123!"
            ... ))
            >>> print(f"Created: {new_user.id}")
        """
        response = self._client._request(
            "POST", "users", json=user_data.to_dict()
        )
        return User.from_dict(response)
    
    def update_user(self, user_id: str, updates: UserUpdate) -> User:
        """
        Update an existing user.
        
        Args:
            user_id: User identifier
            updates: Fields to update
        
        Returns:
            Updated User object
        
        Raises:
            ResourceNotFoundError: If user doesn't exist
            ValidationError: If update data is invalid
        
        Example:
            >>> updated = auth.update_user(
            ...     "user-123",
            ...     UserUpdate(name="New Name", role=Role.ADMIN)
            ... )
        """
        response = self._client._request(
            "PUT", f"users/{user_id}", json=updates.to_dict()
        )
        return User.from_dict(response)
    
    def delete_user(self, user_id: str, soft: bool = True) -> None:
        """
        Delete a user.
        
        Args:
            user_id: User identifier
            soft: If True, perform soft delete (default: True)
        
        Raises:
            ResourceNotFoundError: If user doesn't exist
            AuthorizationError: If deleting the last admin
        
        Example:
            >>> auth.delete_user("user-123", soft=True)  # Soft delete
            >>> auth.delete_user("user-456", soft=False)  # Hard delete
        """
        params = {"soft": soft} if soft else None
        self._client._request("DELETE", f"users/{user_id}", params=params)
    
    def reset_password(self, user_id: str) -> str:
        """
        Reset a user's password.
        
        Args:
            user_id: User identifier
        
        Returns:
            Temporary password
        
        Raises:
            ResourceNotFoundError: If user doesn't exist
        
        Example:
            >>> temp_password = auth.reset_password("user-123")
            >>> print(f"Temporary password: {temp_password}")
        """
        response = self._client._request(
            "POST", f"users/{user_id}/reset-password"
        )
        return response.get("password", "")
    
    def change_password(self, old_password: str, new_password: str) -> None:
        """
        Change the current user's password.
        
        Args:
            old_password: Current password
            new_password: New password
        
        Raises:
            AuthenticationError: If old password is incorrect
            ValidationError: If new password doesn't meet requirements
        
        Example:
            >>> auth.change_password("oldpass", "newpass123!")
        """
        data = {
            "old_password": old_password,
            "new_password": new_password,
        }
        self._client._request("POST", "auth/change-password", json=data)
    
    def lock_user(self, user_id: str) -> User:
        """
        Lock a user account.
        
        Args:
            user_id: User identifier
        
        Returns:
            Updated User object
        
        Example:
            >>> auth.lock_user("user-123")
        """
        response = self._client._request(
            "POST", f"users/{user_id}/lock"
        )
        return User.from_dict(response)
    
    def unlock_user(self, user_id: str) -> User:
        """
        Unlock a user account.
        
        Args:
            user_id: User identifier
        
        Returns:
            Updated User object
        
        Example:
            >>> auth.unlock_user("user-123")
        """
        response = self._client._request(
            "POST", f"users/{user_id}/unlock"
        )
        return User.from_dict(response)
    
    # =====================================================================
    # Session Management
    # =====================================================================
    
    def list_sessions(
        self,
        user_id: Optional[str] = None,
        limit: int = 100,
        offset: int = 0,
    ) -> List[Session]:
        """
        List active sessions.
        
        Args:
            user_id: Filter by user ID
            limit: Maximum number of results
            offset: Pagination offset
        
        Returns:
            List of Session objects
        
        Example:
            >>> sessions = auth.list_sessions()
            >>> for session in sessions:
            ...     print(f"Session: {session.id} - {session.ip_address}")
        """
        params = {"limit": limit, "offset": offset}
        if user_id:
            params["user_id"] = user_id
        
        response = self._client._request("GET", "sessions", params=params)
        
        sessions = []
        if "sessions" in response:
            sessions = [Session.from_dict(s) for s in response["sessions"]]
        elif "data" in response:
            sessions = [Session.from_dict(s) for s in response["data"]]
        
        return sessions
    
    def get_session(self, session_id: str) -> Session:
        """
        Get a session by ID.
        
        Args:
            session_id: Session identifier
        
        Returns:
            Session object
        
        Example:
            >>> session = auth.get_session("sess-123")
            >>> print(f"User: {session.user_email}")
        """
        response = self._client._request("GET", f"sessions/{session_id}")
        return Session.from_dict(response)
    
    def revoke_session(self, session_id: str) -> None:
        """
        Revoke a specific session.
        
        Args:
            session_id: Session identifier
        
        Example:
            >>> auth.revoke_session("sess-123")
        """
        self._client._request("DELETE", f"sessions/{session_id}")
    
    def revoke_all_sessions(self, user_id: Optional[str] = None) -> int:
        """
        Revoke all sessions, optionally for a specific user.
        
        Args:
            user_id: If provided, only revoke sessions for this user
        
        Returns:
            Number of sessions revoked
        
        Example:
            >>> count = auth.revoke_all_sessions()
            >>> print(f"Revoked {count} sessions")
            >>> auth.revoke_all_sessions("user-123")  # Only for specific user
        """
        params = {"user_id": user_id} if user_id else None
        response = self._client._request("DELETE", "sessions", params=params)
        return response.get("revoked_count", 0)
    
    def get_session_stats(self) -> SessionInfo:
        """
        Get session statistics.
        
        Returns:
            SessionInfo object
        
        Example:
            >>> stats = auth.get_session_stats()
            >>> print(f"Active sessions: {stats.total_active}")
        """
        response = self._client._request("GET", "sessions/stats")
        return SessionInfo.from_dict(response)
    
    # =====================================================================
    # Role & Permission Management
    # =====================================================================
    
    def change_role(self, user_id: str, role: Role) -> User:
        """
        Change a user's role.
        
        Args:
            user_id: User identifier
            role: New role
        
        Returns:
            Updated User object
        
        Example:
            >>> auth.change_role("user-123", Role.ADMIN)
        """
        response = self._client._request(
            "PATCH", f"users/{user_id}/role", json={"role": role.value}
        )
        return User.from_dict(response)
    
    def list_roles(self) -> List[Dict[str, Any]]:
        """
        List available roles and their permissions.
        
        Returns:
            List of role definitions
        
        Example:
            >>> roles = auth.list_roles()
            >>> for role in roles:
            ...     print(f"Role: {role['name']}")
        """
        response = self._client._request("GET", "roles")
        return response.get("roles", [])
    
    def get_current_user(self) -> User:
        """
        Get the currently authenticated user.
        
        Returns:
            User object for the current user
        
        Example:
            >>> me = auth.get_current_user()
            >>> print(f"Logged in as: {me.email}")
        """
        response = self._client._request("GET", "auth/me")
        if "user" in response:
            return User.from_dict(response["user"])
        return User.from_dict(response)


# Type hint reference - import only for type checking
if False:
    from aegisgate.client import AegisGateClient
