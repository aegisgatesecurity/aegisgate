# SPDX-License-Identifier: MIT
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.

"""
Authentication service for the AegisGate Python SDK.
"""

from typing import Any, Dict, List, Optional

from aegisgate.connection import SyncConnection, AsyncConnection
from aegisgate.models import User


class AuthService:
    """Synchronous authentication service."""
    
    def __init__(self, connection: SyncConnection):
        self._conn = connection
    
    def login(self, username: str, password: str) -> Dict[str, Any]:
        """
        Perform user login.
        
        Args:
            username: User's username
            password: User's password
            
        Returns:
            Login result with token
        """
        response = self._conn.post(
            "/api/v1/auth/login",
            json_data={"username": username, "password": password}
        )
        return response
    
    def logout(self) -> None:
        """Perform user logout."""
        self._conn.post("/api/v1/auth/logout")
    
    def list_users(self) -> List[User]:
        """
        List all users.
        
        Returns:
            List of User objects
        """
        response = self._conn.get("/api/v1/users")
        users = response.get("users", [])
        return [User.from_dict(u) for u in users]
    
    def get_user(self, user_id: str) -> User:
        """
        Get a specific user.
        
        Args:
            user_id: User ID
            
        Returns:
            User object
        """
        response = self._conn.get(f"/api/v1/users/{user_id}")
        return User.from_dict(response)
    
    def create_user(self, email: str, name: str, roles: List[str]) -> User:
        """
        Create a new user.
        
        Args:
            email: User's email
            name: User's name
            roles: User's roles
            
        Returns:
            Created User object
        """
        response = self._conn.post(
            "/api/v1/users",
            json_data={"email": email, "name": name, "roles": roles}
        )
        return User.from_dict(response)
    
    def update_user(self, user_id: str, **kwargs) -> User:
        """
        Update a user.
        
        Args:
            user_id: User ID
            **kwargs: Fields to update
            
        Returns:
            Updated User object
        """
        response = self._conn.patch(f"/api/v1/users/{user_id}", json_data=kwargs)
        return User.from_dict(response)
    
    def delete_user(self, user_id: str) -> None:
        """
        Delete a user.
        
        Args:
            user_id: User ID
        """
        self._conn.delete(f"/api/v1/users/{user_id}")


class AsyncAuthService:
    """Asynchronous authentication service."""
    
    def __init__(self, connection: AsyncConnection):
        self._conn = connection
    
    async def login(self, username: str, password: str) -> Dict[str, Any]:
        """
        Perform user login.
        
        Args:
            username: User's username
            password: User's password
            
        Returns:
            Login result with token
        """
        response = await self._conn.post(
            "/api/v1/auth/login",
            json_data={"username": username, "password": password}
        )
        return response
    
    async def logout(self) -> None:
        """Perform user logout."""
        await self._conn.post("/api/v1/auth/logout")
    
    async def list_users(self) -> List[User]:
        """
        List all users.
        
        Returns:
            List of User objects
        """
        response = await self._conn.get("/api/v1/users")
        users = response.get("users", [])
        return [User.from_dict(u) for u in users]
    
    async def get_user(self, user_id: str) -> User:
        """
        Get a specific user.
        
        Args:
            user_id: User ID
            
        Returns:
            User object
        """
        response = await self._conn.get(f"/api/v1/users/{user_id}")
        return User.from_dict(response)
    
    async def create_user(self, email: str, name: str, roles: List[str]) -> User:
        """
        Create a new user.
        
        Args:
            email: User's email
            name: User's name
            roles: User's roles
            
        Returns:
            Created User object
        """
        response = await self._conn.post(
            "/api/v1/users",
            json_data={"email": email, "name": name, "roles": roles}
        )
        return User.from_dict(response)
    
    async def update_user(self, user_id: str, **kwargs) -> User:
        """
        Update a user.
        
        Args:
            user_id: User ID
            **kwargs: Fields to update
            
        Returns:
            Updated User object
        """
        response = await self._conn.patch(f"/api/v1/users/{user_id}", json_data=kwargs)
        return User.from_dict(response)
    
    async def delete_user(self, user_id: str) -> None:
        """
        Delete a user.
        
        Args:
            user_id: User ID
        """
        await self._conn.delete(f"/api/v1/users/{user_id}")
