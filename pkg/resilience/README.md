# Resilience Package

This package provides resilience patterns for the AegisGate proxy, including circuit breakers, retry logic, and timeouts.

## Circuit Breaker

The circuit breaker pattern prevents cascading failures by wrapping operations and monitoring their success/failure rates.

### Configuration

```go
config := resilience.CircuitBreakerConfig{
    FailureThreshold: 5,      // Number of failures before opening circuit
    SuccessThreshold: 3,      // Successes needed to close circuit from half-open
    Timeout: 30 * time.Second, // Time circuit stays open before half-open
    RequestTimeout: 10 * time.Second, // Timeout for each request
    MaxRequests: 3,         // Max requests in half-open state
}
```

### Usage

```go
cb := resilience.NewCircuitBreaker(config)

err := cb.Execute(ctx, func(ctx context.Context) error {
    // Your operation here
    return doSomething()
})
if err != nil {
    // Handle circuit breaker error
}
```

### States

- **Closed**: Normal operation, requests pass through
- **Open**: Failures exceeded threshold, requests fail fast
- **Half-Open**: Testing if service recovered

### Metrics

The circuit breaker exposes metrics via atomic counters:
- `TotalRequests()`: Total requests attempted
- `FailedRequests()`: Failed requests
- `RejectedRequests()`: Rejected requests (circuit open)
- `GetState()`: Current state as string

## Retry with Exponential Backoff

```go
config := resilience.RetryConfig{
    MaxAttempts: 3,
    InitialDelay: 100 * time.Millisecond,
    MaxDelay: 5 * time.Second,
    Multiplier: 2.0,
    Jitter: true,
}

executor := resilience.NewRetryExecutor(config)
err := executor.Execute(ctx, operation)
```

## Timeout Executor

```go
executor := resilience.NewTimeoutExecutor(10 * time.Second)
err := executor.Execute(ctx, operation)
```

## ResilientClient

Combine all patterns:

```go
client := resilience.NewResilientClient(resilience.ResilientClientConfig{
    CircuitBreaker: cbConfig,
    Retry: retryConfig,
    Timeout: 10 * time.Second,
})
```
