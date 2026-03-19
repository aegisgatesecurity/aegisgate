# SPDX-License-Identifier: MIT
# =========================================================================
# PROPRIETARY - AegisGate Security
# Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
# =========================================================================
"""
Utility functions for the AegisGate Python SDK.

This module provides helper functions for common operations.
"""

from __future__ import annotations

import hashlib
import hmac
import json
import time
from datetime import datetime, timezone
from typing import Any, Dict, Optional


def format_timestamp(dt: Optional[datetime] = None) -> str:
    """
    Format a datetime as ISO 8601 string.
    
    Args:
        dt: Datetime to format (default: now in UTC)
    
    Returns:
        ISO 8601 formatted string
    
    Example:
        >>> format_timestamp()
        '2025-01-15T10:30:00.000000+00:00'
    """
    if dt is None:
        dt = datetime.now(timezone.utc)
    elif dt.tzinfo is None:
        dt = dt.replace(tzinfo=timezone.utc)
    return dt.isoformat()


def parse_timestamp(ts: str) -> datetime:
    """
    Parse an ISO 8601 timestamp string.
    
    Args:
        ts: ISO 8601 formatted string
    
    Returns:
        Datetime object
    
    Example:
        >>> parse_timestamp('2025-01-15T10:30:00Z')
        datetime.datetime(2025, 1, 15, 10, 30, tzinfo=...)
    """
    ts = ts.replace("Z", "+00:00")
    return datetime.fromisoformat(ts)


def generate_webhook_signature(
    payload: str,
    secret: str,
    timestamp: Optional[int] = None,
) -> str:
    """
    Generate a signature for webhook payload verification.
    
    Args:
        payload: The webhook payload (JSON string)
        secret: The webhook secret
        timestamp: Unix timestamp (default: current time)
    
    Returns:
        HMAC-SHA256 signature
    
    Example:
        >>> secret = "your-webhook-secret"
        >>> payload = json.dumps({"event": "test"})
        >>> signature = generate_webhook_signature(payload, secret)
    """
    if timestamp is None:
        timestamp = int(time.time())
    
    signed_payload = f"{timestamp}.{payload}"
    signature = hmac.new(
        secret.encode("utf-8"),
        signed_payload.encode("utf-8"),
        hashlib.sha256,
    ).hexdigest()
    
    return f"t={timestamp},v1={signature}"


def verify_webhook_signature(
    payload: str,
    signature: str,
    secret: str,
    tolerance: int = 300,
) -> bool:
    """
    Verify a webhook signature.
    
    Args:
        payload: The webhook payload (JSON string)
        signature: The signature from the webhook header
        secret: The webhook secret
        tolerance: Time tolerance in seconds (default: 300)
    
    Returns:
        True if signature is valid
    
    Example:
        >>> is_valid = verify_webhook_signature(
        ...     payload, signature, secret
        ... )
    """
    try:
        parts = dict(p.split("=", 1) for p in signature.split(","))
        timestamp = int(parts.get("t", 0))
        received_sig = parts.get("v1", "")
        
        # Check timestamp is within tolerance
        current_time = int(time.time())
        if abs(current_time - timestamp) > tolerance:
            return False
        
        # Generate expected signature
        expected_sig = generate_webhook_signature(payload, secret, timestamp)
        expected_parts = dict(p.split("=", 1) for p in expected_sig.split(","))
        
        return hmac.compare_digest(received_sig, expected_parts.get("v1", ""))
    except (ValueError, KeyError):
        return False


def mask_sensitive_data(data: Dict[str, Any], fields: Optional[list] = None) -> Dict[str, Any]:
    """
    Mask sensitive fields in a dictionary.
    
    Args:
        data: Dictionary to mask
        fields: List of field names to mask (default: common sensitive fields)
    
    Returns:
        Dictionary with masked values
    
    Example:
        >>> data = {"password": "secret123", "email": "test@example.com"}
        >>> mask_sensitive_data(data)
        {'password': '***', 'email': 'te***@example.com'}
    """
    if fields is None:
        fields = [
            "password", "secret", "token", "api_key", "apikey",
            "private_key", "privatekey", "credit_card", "ssn",
        ]
    
    result = {}
    for key, value in data.items():
        if any(f in key.lower() for f in fields):
            result[key] = "***"
        elif isinstance(value, dict):
            result[key] = mask_sensitive_data(value, fields)
        else:
            result[key] = value
    
    return result


def truncate_string(s: str, max_length: int = 100, suffix: str = "...") -> str:
    """
    Truncate a string to a maximum length.
    
    Args:
        s: String to truncate
        max_length: Maximum length
        suffix: Suffix to append if truncated
    
    Returns:
        Truncated string
    
    Example:
        >>> truncate_string("Hello World", 8)
        'Hello...'
    """
    if len(s) <= max_length:
        return s
    # Reserve space for suffix
    suffix_len = len(suffix)
    if suffix_len >= max_length:
        return suffix[:max_length]
    # Truncate and add suffix
    return s[:max_length - suffix_len] + suffix


def sanitize_filename(filename: str) -> str:
    """
    Sanitize a filename by removing unsafe characters.
    
    Args:
        filename: Filename to sanitize
    
    Returns:
        Sanitized filename
    
    Example:
        >>> sanitize_filename("test<script>.pdf")
        'testscript.pdf'
    """
    import re
    # Remove unsafe characters
    sanitized = re.sub(r'[<>:"/\\|?*\x00-\x1f]', "_", filename)
    # Remove leading/trailing dots and spaces
    sanitized = sanitized.strip(". ")
    # Limit length
    if len(sanitized) > 255:
        name, ext = sanitized.rsplit(".", 1) if "." in sanitized else (sanitized, "")
        max_name = 255 - len(ext) - 1 if ext else 255
        sanitized = name[:max_name] + ("." + ext if ext else "")
    return sanitized


def format_bytes(num_bytes: int) -> str:
    """
    Format bytes as human-readable string.
    
    Args:
        num_bytes: Number of bytes
    
    Returns:
        Human-readable string
    
    Example:
        >>> format_bytes(1024 * 1024 * 100)
        '100.0 MB'
    """
    for unit in ["B", "KB", "MB", "GB", "TB"]:
        if abs(num_bytes) < 1024.0:
            return f"{num_bytes:.1f} {unit}"
        num_bytes /= 1024.0
    return f"{num_bytes:.1f} PB"


def parse_duration(duration_str: str) -> int:
    """
    Parse duration string to seconds.
    
    Args:
        duration_str: Duration string (e.g., "1h30m", "30s", "1d")
    
    Returns:
        Duration in seconds
    
    Example:
        >>> parse_duration("1h30m")
        5400
    """
    units = {
        "s": 1,
        "m": 60,
        "h": 3600,
        "d": 86400,
        "w": 604800,
    }
    
    total = 0
    current = ""
    
    for char in duration_str:
        if char.isdigit():
            current += char
        elif char.lower() in units:
            if current:
                total += int(current) * units[char.lower()]
            current = ""
    
    # Handle remaining digits (without unit)
    if current:
        total += int(current)
    
    return total


def deep_merge(dict1: Dict[str, Any], dict2: Dict[str, Any]) -> Dict[str, Any]:
    """
    Deep merge two dictionaries.
    
    Args:
        dict1: First dictionary
        dict2: Second dictionary (takes precedence)
    
    Returns:
        Merged dictionary
    
    Example:
        >>> deep_merge({"a": 1}, {"b": 2})
        {'a': 1, 'b': 2}
    """
    result = dict1.copy()
    
    for key, value in dict2.items():
        if key in result and isinstance(result[key], dict) and isinstance(value, dict):
            result[key] = deep_merge(result[key], value)
        else:
            result[key] = value
    
    return result


def flatten_dict(d: Dict[str, Any], parent_key: str = "", sep: str = ".") -> Dict[str, Any]:
    """
    Flatten a nested dictionary.
    
    Args:
        d: Dictionary to flatten
        parent_key: Parent key prefix
        sep: Separator for nested keys
    
    Returns:
        Flattened dictionary
    
    Example:
        >>> flatten_dict({"a": {"b": {"c": 1}}})
        {'a.b.c': 1}
    """
    items = []
    for k, v in d.items():
        new_key = f"{parent_key}{sep}{k}" if parent_key else k
        if isinstance(v, dict):
            items.extend(flatten_dict(v, new_key, sep).items())
        else:
            items.append((new_key, v))
    return dict(items)
