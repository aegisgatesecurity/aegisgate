"""
Example 2: Async AegisGate Client Usage

This example demonstrates asynchronous usage of the AegisGate client
for high-performance concurrent operations.
"""

import asyncio
import os
from aegisgate import AsyncClient


# Set up credentials
API_KEY = os.environ.get("AEGISGATE_API_KEY", "your-api-key-here")
BASE_URL = os.environ.get("AEGISGATE_BASE_URL", "https://api.aegisgate.io")


async def main():
    """Demonstrate async client usage."""
    
    # Initialize async client
    async with AsyncClient(
        base_url=BASE_URL,
        api_key=API_KEY
    ) as client:
        # 1. Async Health Check
        print("=== Async Health Check ===")
        health = await client.core.health()
        print(f"Status: {health.status}")
        print()
        
        # 2. Concurrent Operations
        print("=== Concurrent Operations ===")
        
        # Run multiple operations concurrently
        results = await asyncio.gather(
            client.core.health(),
            client.proxy.get_stats(),
            client.core.version(),
            return_exceptions=True
        )
        
        print(f"Health: {results[0].status if hasattr(results[0], 'status') else results[0]}")
        print(f"Stats: {results[1].total_requests if hasattr(results[1], 'total_requests') else results[1]}")
        print(f"Version: {results[2].version if hasattr(results[2], 'version') else results[2]}")
        print()


async def example_batch_request_inspection():
    """Example of batch request inspection."""
    
    prompts = [
        "What is Python?",
        "Explain machine learning",
        "How does blockchain work?",
        "What is cloud computing?",
        "Explain DevOps"
    ]
    
    async with AsyncClient(api_key=API_KEY) as client:
        # Inspect all prompts concurrently
        inspection_tasks = [
            client.proxy.inspect_request(content=prompt, content_type="prompt")
            for prompt in prompts
        ]
        
        results = await asyncio.gather(*inspection_tasks, return_exceptions=True)
        
        for prompt, result in zip(prompts, results):
            if isinstance(result, Exception):
                print(f"Error inspecting '{prompt[:30]}...': {result}")
            else:
                has_violations = getattr(result, 'has_violations', False)
                print(f"Prompt: '{prompt[:30]}...' - Violations: {has_violations}")


async def example_concurrent_health_checks():
    """Example of concurrent health monitoring."""
    
    services = [
        ("https://api1.aegisgate.io", "Service 1"),
        ("https://api2.aegisgate.io", "Service 2"),
        ("https://api3.aegisgate.io", "Service 3"),
    ]
    
    async def check_service_health(base_url: str, name: str):
        try:
            async with AsyncClient(base_url=base_url, api_key=API_KEY) as client:
                health = await client.core.health()
                return (name, health.status, None)
        except Exception as e:
            return (name, "error", str(e))
    
    # Check all services concurrently
    tasks = [
        check_service_health(url, name)
        for url, name in services
    ]
    
    results = await asyncio.gather(*tasks)
    
    for name, status, error in results:
        if error:
            print(f"{name}: Error - {error}")
        else:
            print(f"{name}: {status}")


async def example_stream_processing():
    """Example of handling streamed responses."""
    
    async with AsyncClient(api_key=API_KEY) as client:
        # Simulate processing a stream of requests
        async def process_request(request_id: int):
            result = await client.proxy.inspect_request(
                content=f"Request {request_id}",
                content_type="prompt"
            )
            return request_id, result.has_violations
        
        # Process stream of requests
        request_ids = range(1, 11)  # 10 requests
        tasks = [process_request(rid) for rid in request_ids]
        
        async for coro in asyncio.as_completed(tasks):
            request_id, has_violations = await coro
            print(f"Request {request_id}: violations={has_violations}")


async def example_async_context_managers():
    """Example of different context manager patterns."""
    
    # Pattern 1: Single client for multiple operations
    async with AsyncClient(api_key=API_KEY) as client:
        health = await client.core.health()
        stats = await client.proxy.get_stats()
        # Connection stays open for both operations
    
    # Pattern 2: Client per operation (not recommended for production)
    async def single_operation():
        async with AsyncClient(api_key=API_KEY) as client:
            return await client.core.health()
    
    result = await single_operation()
    print(f"Health: {result.status}")


if __name__ == "__main__":
    # Run main async function
    asyncio.run(main())
    
    # Run other examples
    print("\n=== Batch Inspection ===")
    asyncio.run(example_batch_request_inspection())
    
    print("\n=== Concurrent Health Checks ===")
    asyncio.run(example_concurrent_health_checks())