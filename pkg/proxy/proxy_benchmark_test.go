// Package proxy provides benchmark tests for the security proxy.
// These tests measure the performance of request forwarding, parsing,
// TLS handshakes, HTTP/2, and MITM proxy operations.
//
// Run benchmarks with: go test -bench=. -benchmem ./pkg/proxy/...
//
//go:build !integration
// +build !integration

package proxy_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aegisgatesecurity/aegisgate/pkg/proxy"
	"github.com/aegisgatesecurity/aegisgate/pkg/resilience"
)

// ============================================================================
// Test Server Helpers
// ============================================================================

// simpleBackendHandler is a minimal HTTP handler that returns OK
func simpleBackendHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

// echoBackendHandler echoes back request body and headers
func echoBackendHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		if r.Body != nil {
			io.Copy(w, r.Body)
		}
	})
}

// createTestProxy creates a proxy instance for benchmarking
func createTestProxy(upstreamURL string, opts *proxy.Options) *proxy.Proxy {
	if opts == nil {
		opts = &proxy.Options{
			BindAddress: "127.0.0.1:0",
			Upstream:    upstreamURL,
			MaxBodySize: 10 * 1024 * 1024,
			Timeout:     30 * time.Second,
			RateLimit:   10000, // High limit for benchmarks
		}
	}
	return proxy.New(opts)
}

// generatePayload creates a payload of specified size
func generatePayload(size int) []byte {
	return make([]byte, size)
}

// ============================================================================
// Request Forwarding Benchmarks
// ============================================================================

// BenchmarkRequestForwarding benchmarks full request lifecycle through proxy
func BenchmarkRequestForwarding(b *testing.B) {
	// Create test backend
	backend := httptest.NewServer(simpleBackendHandler())
	defer backend.Close()

	// Create proxy
	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			b.Fatalf("unexpected status code: %d", rec.Code)
		}
	}
}

// BenchmarkRequestForwarding_PostWithBody benchmarks POST requests with body
func BenchmarkRequestForwarding_PostWithBody(b *testing.B) {
	backend := httptest.NewServer(echoBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	body := []byte("test request body")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		p.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			b.Fatalf("unexpected status code: %d", rec.Code)
		}
	}
}

// BenchmarkRequestForwarding_WithRateLimiting benchmarks with rate limiter active
func BenchmarkRequestForwarding_WithRateLimiting(b *testing.B) {
	backend := httptest.NewServer(simpleBackendHandler())
	defer backend.Close()

	opts := &proxy.Options{
		BindAddress: "127.0.0.1:0",
		Upstream:    backend.URL,
		MaxBodySize: 10 * 1024 * 1024,
		Timeout:     30 * time.Second,
		RateLimit:   100000, // Very high for benchmarks
	}
	p := proxy.New(opts)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			b.Fatalf("unexpected status code: %d", rec.Code)
		}
	}
}

// BenchmarkRequestForwarding_LargeBody benchmarks with large request body
func BenchmarkRequestForwarding_LargeBody(b *testing.B) {
	backend := httptest.NewServer(echoBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	body := generatePayload(100 * 1024) // 100KB

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/octet-stream")
		p.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			b.Fatalf("unexpected status code: %d", rec.Code)
		}
	}
}

// ============================================================================
// Request Parsing Benchmarks
// ============================================================================

// BenchmarkRequestParsing_GET benchmarks GET request parsing
func BenchmarkRequestParsing_GET(b *testing.B) {
	backend := httptest.NewServer(simpleBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/resource?id=123&name=test", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "BenchmarkClient/1.0")
	req.Header.Set("Authorization", "Bearer test-token-12345")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
	}
}

// BenchmarkRequestParsing_POST_JSON benchmarks POST with JSON body
func BenchmarkRequestParsing_POST_JSON(b *testing.B) {
	backend := httptest.NewServer(echoBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	jsonBody := []byte(`{"name":"test","value":123,"data":"example"}`)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/resource", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		p.ServeHTTP(rec, req)
	}
}

// BenchmarkRequestParsing_WithHeaders benchmarks request with many headers
func BenchmarkRequestParsing_WithHeaders(b *testing.B) {
	backend := httptest.NewServer(simpleBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer token")
		req.Header.Set("User-Agent", "TestClient/1.0")
		req.Header.Set("X-Request-ID", "12345")
		req.Header.Set("X-Forwarded-For", "192.168.1.1")
		p.ServeHTTP(rec, req)
	}
}

// ============================================================================
// Response Handling Benchmarks
// ============================================================================

// BenchmarkResponseHandling_Small_1KB benchmarks 1KB response
func BenchmarkResponseHandling_Small_1KB(b *testing.B) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(generatePayload(1024))
	}))
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			b.Fatalf("unexpected status code: %d", rec.Code)
		}
	}
}

// BenchmarkResponseHandling_Medium_100KB benchmarks 100KB response
func BenchmarkResponseHandling_Medium_100KB(b *testing.B) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(generatePayload(100 * 1024))
	}))
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			b.Fatalf("unexpected status code: %d", rec.Code)
		}
	}
}

// BenchmarkResponseHandling_Large_1MB benchmarks 1MB response
func BenchmarkResponseHandling_Large_1MB(b *testing.B) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(generatePayload(1024 * 1024))
	}))
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			b.Fatalf("unexpected status code: %d", rec.Code)
		}
	}
}

// BenchmarkResponseHandling_JSON benchmarks JSON response handling
func BenchmarkResponseHandling_JSON(b *testing.B) {
	jsonResponse := []byte(`{"status":"success","data":{"id":123,"name":"test","items":[{"a":1},{"b":2}]}}`)
	
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}))
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	req.Header.Set("Accept", "application/json")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
	}
}

// ============================================================================
// Parallel Request Benchmarks
// ============================================================================

// BenchmarkParallel_RequestForwarding benchmarks concurrent request handling
func BenchmarkParallel_RequestForwarding(b *testing.B) {
	backend := httptest.NewServer(simpleBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			p.ServeHTTP(rec, req)
		}
	})
}

// BenchmarkParallel_RequestForwarding_POST benchmarks concurrent POST handling
func BenchmarkParallel_RequestForwarding_POST(b *testing.B) {
	backend := httptest.NewServer(echoBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	body := []byte("concurrent test data")

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
			req.Header.Set("Content-Type", "text/plain")
			p.ServeHTTP(rec, req)
		}
	})
}

// BenchmarkParallel_ResponseHandling benchmarks concurrent response processing
func BenchmarkParallel_ResponseHandling(b *testing.B) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(generatePayload(10 * 1024)) // 10KB response
	}))
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			p.ServeHTTP(rec, req)
		}
	})
}

// ============================================================================
// Circuit Breaker Benchmarks
// ============================================================================

// BenchmarkCircuitBreaker_Enabled benchmarks proxy with circuit breaker
func BenchmarkCircuitBreaker_Enabled(b *testing.B) {
	backend := httptest.NewServer(simpleBackendHandler())
	defer backend.Close()

	cbConfig := &resilience.CircuitBreakerConfig{
		FailureThreshold: 5,
		Timeout:          30 * time.Second,
	}
	
	opts := &proxy.Options{
		BindAddress:    "127.0.0.1:0",
		Upstream:       backend.URL,
		MaxBodySize:    10 * 1024 * 1024,
		Timeout:        30 * time.Second,
		RateLimit:      10000,
		CircuitBreaker: cbConfig,
	}
	p := proxy.New(opts)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
	}
}

// ============================================================================
// Scanner and Compliance Benchmarks
// ============================================================================

// BenchmarkScanner_RequestBody benchmarks request body scanning
func BenchmarkScanner_RequestBody(b *testing.B) {
	backend := httptest.NewServer(echoBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	// Clean payload without sensitive data
	cleanPayload := generatePayload(1024)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(cleanPayload))
		req.Header.Set("Content-Type", "application/octet-stream")
		p.ServeHTTP(rec, req)
	}
}

// BenchmarkScanner_ResponseBody benchmarks response body scanning
func BenchmarkScanner_ResponseBody(b *testing.B) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(generatePayload(1024))
	}))
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
	}
}

// BenchmarkCompliance_ATLAS benchmarks MITRE ATLAS compliance checking
func BenchmarkCompliance_ATLAS(b *testing.B) {
	backend := httptest.NewServer(echoBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	// Normal request without threats
	payload := []byte("normal request data")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "text/plain")
		p.ServeHTTP(rec, req)
	}
}

// ============================================================================
// Latency and Throughput Benchmarks
// ============================================================================

// BenchmarkLatency_EndToEnd measures end-to-end latency
func BenchmarkLatency_EndToEnd(b *testing.B) {
	backend := httptest.NewServer(simpleBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	
	b.ResetTimer()
	
	var totalLatency time.Duration
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		start := time.Now()
		p.ServeHTTP(rec, req)
		totalLatency += time.Since(start)
	}
	
	b.ReportMetric(float64(totalLatency.Nanoseconds())/float64(b.N), "ns/op")
	b.ReportMetric(float64(b.N)/totalLatency.Seconds(), "req/s")
}

// BenchmarkThroughput_Concurrent measures throughput under concurrent load
func BenchmarkThroughput_Concurrent(b *testing.B) {
	backend := httptest.NewServer(simpleBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			p.ServeHTTP(rec, req)
		}
	})
	
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "req/s")
}

// ============================================================================
// Memory Allocation Benchmarks
// ============================================================================

// BenchmarkMemoryAllocation_RequestProcessing measures memory per request
func BenchmarkMemoryAllocation_RequestProcessing(b *testing.B) {
	backend := httptest.NewServer(echoBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	body := []byte("test data for memory measurement")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "text/plain")
		p.ServeHTTP(rec, req)
	}
}

// BenchmarkMemoryAllocation_LargePayload measures memory with large payloads
func BenchmarkMemoryAllocation_LargePayload(b *testing.B) {
	backend := httptest.NewServer(echoBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	body := generatePayload(100 * 1024) // 100KB

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/octet-stream")
		p.ServeHTTP(rec, req)
	}
}

// ============================================================================
// Error Handling Benchmarks
// ============================================================================

// BenchmarkErrorHandling_BodyTooLarge benchmarks rejection of oversized bodies
func BenchmarkErrorHandling_BodyTooLarge(b *testing.B) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:0",
		Upstream:    "http://127.0.0.1:3000",
		MaxBodySize: 1024, // 1KB limit
		Timeout:     30 * time.Second,
		RateLimit:   10000,
	}
	
	p := proxy.New(opts)
	defer p.Stop(context.Background())

	body := generatePayload(10 * 1024) // 10KB - exceeds limit

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/octet-stream")
		p.ServeHTTP(rec, req)
		if rec.Code != http.StatusRequestEntityTooLarge {
			b.Fatalf("expected 413, got: %d", rec.Code)
		}
	}
}

// BenchmarkErrorHandling_RateLimitExceeded benchmarks rate limiting
func BenchmarkErrorHandling_RateLimitExceeded(b *testing.B) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:0",
		Upstream:    "http://127.0.0.1:3000",
		MaxBodySize: 10 * 1024 * 1024,
		Timeout:     30 * time.Second,
		RateLimit:   1, // Very low limit
	}
	
	p := proxy.New(opts)
	defer p.Stop(context.Background())

	// Drain the rate limiter first
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	p.ServeHTTP(rec, req)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		p.ServeHTTP(rec, req)
		if rec.Code != http.StatusTooManyRequests {
			b.Fatalf("expected 429, got: %d", rec.Code)
		}
	}
}

// ============================================================================
// HTTP Method Variation Benchmarks
// ============================================================================

// BenchmarkMethod_GET benchmarks GET requests
func BenchmarkMethod_GET(b *testing.B) {
	backend := httptest.NewServer(simpleBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
	}
}

// BenchmarkMethod_POST benchmarks POST requests
func BenchmarkMethod_POST(b *testing.B) {
	backend := httptest.NewServer(echoBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	body := []byte("post data")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "text/plain")
		p.ServeHTTP(rec, req)
	}
}

// BenchmarkMethod_PUT benchmarks PUT requests
func BenchmarkMethod_PUT(b *testing.B) {
	backend := httptest.NewServer(echoBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	body := []byte("put data")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "text/plain")
		p.ServeHTTP(rec, req)
	}
}

// BenchmarkMethod_DELETE benchmarks DELETE requests
func BenchmarkMethod_DELETE(b *testing.B) {
	backend := httptest.NewServer(simpleBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodDelete, "/test", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
	}
}

// ============================================================================
// Security Header Benchmarks
// ============================================================================

// BenchmarkSecurityHeaders_Outbound measures security headers on outbound requests
func BenchmarkSecurityHeaders_Outbound(b *testing.B) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify security headers are present
		if r.Header.Get("X-Forwarded-Proto") != "https" {
			b.Error("X-Forwarded-Proto not set")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
	}
}

// BenchmarkSecurityHeaders_Inbound measures security headers on responses
func BenchmarkSecurityHeaders_Inbound(b *testing.B) {
	backend := httptest.NewServer(simpleBackendHandler())
	defer backend.Close()

	p := createTestProxy(backend.URL, nil)
	defer p.Stop(context.Background())

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		p.ServeHTTP(rec, req)
		
		// Verify security headers
		headers := rec.Result().Header
		if headers.Get("X-Frame-Options") != "DENY" {
			b.Error("X-Frame-Options not set")
		}
		if headers.Get("X-Content-Type-Options") != "nosniff" {
			b.Error("X-Content-Type-Options not set")
		}
	}
}
