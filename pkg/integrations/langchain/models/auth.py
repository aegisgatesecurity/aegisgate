# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Authentication and user management models for the AegisGate Python SDK.

This module provides data classes for working with AegisGate authentication,
user management, and session handling.
"""

from __future__ import annotations

import re
from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
from typing import Optional, List, Dict, Any


class Role(str, Enum):
    """
    User roles in AegisGate.
    
    Attributes:
        ADMIN: Full administrative access
        USER: Standard user access
        VIEWER: Read-only access
    """
    ADMIN = "admin"
    USER = "user"
    VIEWER = "viewer"
    
    @classmethod
    def from_string(cls, value: str) -> Role:
        """Create Role from string value."""
        try:
            return cls(value.lower())
        except ValueError:
            return cls.USER


class Provider(str, Enum):
    """
    Authentication providers supported by AegisGate.
    
    Attributes:
        LOCAL: Local username/password authentication
        GOOGLE: Google OAuth
        MICROSOFT: Microsoft/Azure AD OAuth
        GITHUB: GitHub OAuth
        OKTA: Okta SSO
        AUTH0: Auth0
        SAML: SAML 2.0
    """
    LOCAL = "local"
    GOOGLE = "google"
    MICROSOFT = "microsoft"
    GITHUB = "github"
    OKTA = "okta"
    AUTH0 = "auth0"
    SAML = "saml"
    
    @classmethod
    def from_string(cls, value: str) -> Provider:
        """Create Provider from string value."""
        try:
            return cls(value.lower())
        except ValueError:
            return cls.LOCAL


@dataclass
class User:
    """
    Represents a user in the AegisGate system.
    
    Attributes:
        id: Unique user identifier
        email: User's email address
        name: User's display name
        role: User's role (admin/user/viewer)
        provider: Authentication provider
        status: Account status
        created_at: Account creation timestamp
        updated_at: Last update timestamp
        last_login: Last successful login timestamp
        locked: Whether account is locked
        mfa_enabled: Whether MFA is enabled
    """
    id: str
    email: str
    name: str
    role: Role = Role.USER
    provider: Provider = Provider.LOCAL
    status: str = "active"
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None
    last_login: Optional[datetime] = None
    locked: bool = False
    mfa_enabled: bool = False
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> User:
        """Create User from dictionary."""
        role = Role.from_string(data.get("role", "user"))
        provider = Provider.from_string(data.get("provider", "local"))
        
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
            email=data.get("email", ""),
            name=data.get("name", ""),
            role=role,
            provider=provider,
            status=data.get("status", "active"),
            created_at=parse_datetime(data.get("created_at")),
            updated_at=parse_datetime(data.get("updated_at")),
            last_login=parse_datetime(data.get("last_login")),
            locked=data.get("locked", False),
            mfa_enabled=data.get("mfa_enabled", False),
        )
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert User to dictionary."""
        result = {
            "id": self.id,
            "email": self.email,
            "name": self.name,
            "role": self.role.value,
            "provider": self.provider.value,
            "status": self.status,
            "locked": self.locked,
            "mfa_enabled": self.mfa_enabled,
        }
        if self.created_at:
            result["created_at"] = self.created_at.isoformat()
        if self.updated_at:
            result["updated_at"] = self.updated_at.isoformat()
        if self.last_login:
            result["last_login"] = self.last_login.isoformat()
        return result


@dataclass
class UserCreate:
    """
    Data for creating a new user.
    
    Attributes:
        email: User's email address (required)
        name: User's display name (required)
        role: User's role (default: user)
        provider: Authentication provider (default: local)
        password: Password for local auth (required if provider is local)
        send_welcome_email: Whether to send welcome email
    """
    email: str
    name: str
    role: Role = Role.USER
    provider: Provider = Provider.LOCAL
    password: Optional[str] = None
    send_welcome_email: bool = False
    
    def __post_init__(self):
        """Validate the user creation data."""
        if not self.email or not self._is_valid_email(self.email):
            raise ValueError(f"Invalid email address: {self.email}")
        
        if self.provider == Provider.LOCAL and not self.password:
            raise ValueError("Password required for local authentication")
    
    @staticmethod
    def _is_valid_email(email: str) -> bool:
        """Validate email format."""
        pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        return bool(re.match(pattern, email))
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary."""
        result = {
            "email": self.email,
            "name": self.name,
            "role": self.role.value,
            "provider": self.provider.value,
            "send_welcome_email": self.send_welcome_email,
        }
        if self.password:
            result["password"] = self.password
        return result


@dataclass
class UserUpdate:
    """
    Data for updating an existing user.
    
    All fields are optional - only provided fields will be updated.
    
    Attributes:
        name: New display name
        role: New role
        status: New account status
        password: New password
        locked: Lock/unlock account
        mfa_enabled: Enable/disable MFA
    """
    name: Optional[str] = None
    role: Optional[Role] = None
    status: Optional[str] = None
    password: Optional[str] = None
    locked: Optional[bool] = None
    mfa_enabled: Optional[bool] = None
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary."""
        result = {}
        if self.name is not None:
            result["name"] = self.name
        if self.role is not None:
            result["role"] = self.role.value
        if self.status is not None:
            result["status"] = self.status
        if self.password is not None:
            result["password"] = self.password
        if self.locked is not None:
            result["locked"] = self.locked
        if self.mfa_enabled is not None:
            result["mfa_enabled"] = self.mfa_enabled
        return result


@dataclass
class LoginResult:
    """
    Result of a successful login attempt.
    
    Attributes:
        success: Whether login was successful
        token: Authentication token (for session auth)
        expires_at: Token expiration timestamp
        user: User object (if available)
        error: Error message (if login failed)
    """
    success: bool
    token: Optional[str] = None
    expires_at: Optional[datetime] = None
    user: Optional[User] = None
    error: Optional[str] = None
    
    @property
    def is_expired(self) -> bool:
        """Check if the token is expired."""
        from datetime import timezone
        if self.expires_at is None:
            return False
        # Ensure we're comparing timezone-aware datetimes
        now = datetime.now(timezone.utc)
        expires = self.expires_at
        if expires.tzinfo is None:
            expires = expires.replace(tzinfo=timezone.utc)
        return now >= expires
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> LoginResult:
        """Create LoginResult from dictionary."""
        user = None
        if "user" in data and data["user"]:
            user = User.from_dict(data["user"])
        
        expires_at = None
        if data.get("expires_at"):
            try:
                expires_at = datetime.fromisoformat(
                    data["expires_at"].replace("Z", "+00:00")
                )
            except (ValueError, AttributeError):
                pass
        
        return cls(
            success=data.get("success", False),
            token=data.get("token"),
            expires_at=expires_at,
            user=user,
            error=data.get("error"),
        )


@dataclass
class Session:
    """
    Represents an active user session.
    
    Attributes:
        id: Session identifier
        user_id: Associated user ID
        user_email: User's email
        ip_address: Client IP address
        user_agent: Client user agent
        created_at: Session creation time
        expires_at: Session expiration time
        last_activity: Last activity timestamp
        provider: Authentication provider used
    """
    id: str
    user_id: str
    user_email: str
    ip_address: Optional[str] = None
    user_agent: Optional[str] = None
    created_at: Optional[datetime] = None
    expires_at: Optional[datetime] = None
    last_activity: Optional[datetime] = None
    provider: Provider = Provider.LOCAL
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> Session:
        """Create Session from dictionary."""
        provider = Provider.from_string(data.get("provider", "local"))
        
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
            user_id=data.get("user_id", ""),
            user_email=data.get("user_email", ""),
            ip_address=data.get("ip_address"),
            user_agent=data.get("user_agent"),
            created_at=parse_datetime(data.get("created_at")),
            expires_at=parse_datetime(data.get("expires_at")),
            last_activity=parse_datetime(data.get("last_activity")),
            provider=provider,
        )
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert Session to dictionary."""
        result = {
            "id": self.id,
            "user_id": self.user_id,
            "user_email": self.user_email,
            "provider": self.provider.value,
        }
        if self.ip_address:
            result["ip_address"] = self.ip_address
        if self.user_agent:
            result["user_agent"] = self.user_agent
        if self.created_at:
            result["created_at"] = self.created_at.isoformat()
        if self.expires_at:
            result["expires_at"] = self.expires_at.isoformat()
        if self.last_activity:
            result["last_activity"] = self.last_activity.isoformat()
        return result
    
    @property
    def is_expired(self) -> bool:
        """Check if session is expired."""
        from datetime import timezone
        if self.expires_at is None:
            return False
        # Ensure we're comparing timezone-aware datetimes
        now = datetime.now(timezone.utc)
        expires = self.expires_at
        if expires.tzinfo is None:
            expires = expires.replace(tzinfo=timezone.utc)
        return now >= expires


@dataclass
class SessionInfo:
    """
    Information about a session for display purposes.
    
    Attributes:
        total_active: Total number of active sessions
        by_user: Sessions grouped by user
        by_provider: Sessions grouped by provider
        oldest_session: Timestamp of oldest active session
        average_duration: Average session duration in seconds
    """
    total_active: int
    by_user: Dict[str, int] = field(default_factory=dict)
    by_provider: Dict[str, int] = field(default_factory=dict)
    oldest_session: Optional[datetime] = None
    average_duration: Optional[float] = None
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> SessionInfo:
        """Create SessionInfo from dictionary."""
        oldest = None
        if data.get("oldest_session"):
            try:
                oldest = datetime.fromisoformat(
                    data["oldest_session"].replace("Z", "+00:00")
                )
            except (ValueError, AttributeError):
                pass
        
        return cls(
            total_active=data.get("total_active", 0),
            by_user=data.get("by_user", {}),
            by_provider=data.get("by_provider", {}),
            oldest_session=oldest,
            average_duration=data.get("average_duration"),
        )
