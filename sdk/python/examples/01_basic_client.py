"""
Example 1: Basic AegisGate Client Usage

This example demonstrates basic usage of the synchronous AegisGate client
for authentication, health checks, and proxy statistics.
"""

import os
from aegisgate import Client

# Set up credentials (or use environment variables)
API_KEY = os.environ.get("AEGISGATE_API_KEY", "your-api-key-here")
BASE_URL = os.environ.get("AEGISGATE_BASE_URL", "https://api.aegisgate.io")


def main():
    """Demonstrate basic client usage."""
    
    # Initialize client with API key
    # You can also set AEGISGATE_API_KEY environment variable
    # and use: client = Client()
    client = Client(
        base_url=BASE_URL,
        api_key=API_KEY
    )
    
    # Using context manager automatically handles connection cleanup
    with client:
        # 1. Health Check
        print("=== Health Check ===")
        health = client.core.health()
        print(f"Status: {health.status}")
        print(f"Version: {health.version}")
        print()
        
        # 2. Detailed Health Check
        print("=== Detailed Health ===")
        detailed_health = client.core.health_detailed()
        print(f"Status: {detailed_health.status}")
        if hasattr(detailed_health, 'modules'):
            for module in detailed_health.modules:
                print(f"  {module.name}: {module.status}")
        print()
        
        # 3. Get Proxy Statistics
        print("=== Proxy Statistics ===")
        stats = client.proxy.get_stats()
        print(f"Total Requests: {stats.total_requests}")
        print(f"Blocked Requests: {stats.blocked_requests}")
        print(f"Successful Requests: {stats.successful_requests}")
        print()
        
        # 4. Inspect a Request
        print("=== Request Inspection ===")
        test_prompt = "What is machine learning?"
        result = client.proxy.inspect_request(
            content=test_prompt,
            content_type="prompt"
        )
        print(f"Prompt: {test_prompt}")
        print(f"Has Violations: {result.has_violations}")
        if result.has_violations:
            for violation in result.violations:
                print(f"  - {violation.type}: {violation.message}")
        print()
        
        # 5. List Recent Violations
        print("=== Recent Violations ===")
        violations = client.proxy.list_violations(limit=5)
        for v in violations:
            print(f"  [{v.severity}] {v.type}: {v.message}")
        print()


def example_with_context_manager():
    """Example showing context manager pattern."""
    
    # Client automatically manages connections
    with Client(api_key=API_KEY) as client:
        health = client.core.health()
        print(f"Status: {health.status}")


def example_without_context_manager():
    """Example showing manual cleanup pattern."""
    
    client = Client(api_key=API_KEY)
    try:
        health = client.core.health()
        print(f"Status: {health.status}")
    finally:
        # Manually close connection
        client.close()


if __name__ == "__main__":
    main()