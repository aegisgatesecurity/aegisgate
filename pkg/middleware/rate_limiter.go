package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/aegisgatesecurity/aegisgate/pkg/core"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	// RequestsPerMinute - max requests allowed per minute
	RequestsPerMinute int
	// Burst - max burst requests allowed (above steady rate)
	Burst int
	// BlockDuration - how long to block after limit exceeded
	BlockDuration time.Duration
	// KeyFunc - function to generate unique key for rate limiting
	KeyFunc func(r *http.Request) string
}

// DefaultRateLimitConfig returns default configuration based on tier
func DefaultRateLimitConfig(tier core.Tier) RateLimitConfig {
	limits := tier.GetTierLimits()

	// Handle unlimited
	requestsPerMinute := limits.MaxRequestsPerMinute
	if requestsPerMinute == -1 {
		requestsPerMinute = 1000000 // Effectively unlimited
	}

	burst := limits.MaxBurstRequests
	if burst == -1 {
		burst = requestsPerMinute / 10
		if burst < 100 {
			burst = 100
		}
	}

	return RateLimitConfig{
		RequestsPerMinute: requestsPerMinute,
		Burst:            burst,
		BlockDuration:    time.Minute * 5,
		KeyFunc:          DefaultKeyFunc,
	}
}

// DefaultKeyFunc generates a rate limit key from the request
func DefaultKeyFunc(r *http.Request) string {
	// Try to use API key first
	if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
		return "api:" + apiKey
	}

	// Fall back to IP address
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return "ip:" + forwarded
	}

	// Use remote address
	return "ip:" + r.RemoteAddr
}

// RateLimiter implements per-client rate limiting using token bucket
type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*clientLimiter
	config   RateLimitConfig
	cleanup  chan struct{}
	interval time.Duration
}

// clientLimiter holds rate limit state for a single client
type clientLimiter struct {
	tokens   float64
	lastFill time.Time
	blocked  time.Time
	burst    int
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	if config.RequestsPerMinute <= 0 {
		config.RequestsPerMinute = 60
	}
	if config.Burst <= 0 {
		config.Burst = config.RequestsPerMinute / 10
		if config.Burst < 1 {
			config.Burst = 1
		}
	}
	if config.BlockDuration <= 0 {
		config.BlockDuration = time.Minute * 5
	}
	if config.KeyFunc == nil {
		config.KeyFunc = DefaultKeyFunc
	}

	rl := &RateLimiter{
		clients:  make(map[string]*clientLimiter),
		config:   config,
		cleanup:  make(chan struct{}),
		interval: time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanupLoop()

	return rl
}

// cleanupLoop periodically removes stale client entries
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.interval)
	defer ticker.Stop()

	for {
		select {
		case <-rl.cleanup:
			return
		case <-ticker.C:
			rl.cleanupStale()
		}
	}
}

// cleanupStale removes rate limit entries that haven't been accessed recently
func (rl *RateLimiter) cleanupStale() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-rl.interval * 2)
	for key, client := range rl.clients {
		if client.lastFill.Before(cutoff) {
			delete(rl.clients, key)
		}
	}
}

// Allow checks if a request is allowed under rate limits
// Returns true if allowed, false if rate limited
func (rl *RateLimiter) Allow(r *http.Request) bool {
	key := rl.config.KeyFunc(r)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	client, exists := rl.clients[key]
	now := time.Now()

	if !exists {
		// New client - create limiter with full burst
		rl.clients[key] = &clientLimiter{
			tokens:  float64(rl.config.Burst),
			lastFill: now,
			burst:   rl.config.Burst,
		}
		return true
	}

	// Check if client is currently blocked
	if !client.blocked.IsZero() && now.Before(client.blocked.Add(rl.config.BlockDuration)) {
		return false
	}

	// Clear block status if block duration expired
	if !client.blocked.IsZero() {
		client.blocked = time.Time{}
		client.tokens = float64(client.burst)
	}

	// Calculate token refill
	elapsed := now.Sub(client.lastFill).Seconds()
	refillRate := float64(rl.config.RequestsPerMinute) / 60.0
	client.tokens += elapsed * refillRate
	if client.tokens > float64(client.burst) {
		client.tokens = float64(client.burst)
	}
	client.lastFill = now

	// Check if we have tokens available
	if client.tokens >= 1 {
		client.tokens--
		return true
	}

	// Rate limited - block the client
	client.blocked = now
	return false
}

// Middleware returns an HTTP middleware for rate limiting
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow(r) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerMinute))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "rate_limit_exceeded", "message": "Too many requests. Please try again later."}`))
			return
		}

		// Add rate limit headers to response
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerMinute))

		next.ServeHTTP(w, r)
	})
}

// Stop stops the rate limiter cleanup goroutine
func (rl *RateLimiter) Stop() {
	close(rl.cleanup)
}

// GetRateLimitConfig returns the current configuration
func (rl *RateLimiter) GetRateLimitConfig() RateLimitConfig {
	return rl.config
}

// ============================================================================
// Tier-Aware Rate Limiter (uses tier limits automatically)
// ============================================================================

// TierRateLimiter creates a rate limiter based on the tier configuration
type TierRateLimiter struct {
	limiter *RateLimiter
	tier    core.Tier
}

// NewTierRateLimiter creates a new tier-aware rate limiter
func NewTierRateLimiter(tier core.Tier) *TierRateLimiter {
	config := DefaultRateLimitConfig(tier)
	return &TierRateLimiter{
		limiter: NewRateLimiter(config),
		tier:    tier,
	}
}

// Middleware returns the rate limiting middleware
func (trl *TierRateLimiter) Middleware(next http.Handler) http.Handler {
	return trl.limiter.Middleware(next)
}

// UpdateTier updates the rate limiter to use new tier limits
func (trl *TierRateLimiter) UpdateTier(tier core.Tier) {
	trl.tier = tier
	config := DefaultRateLimitConfig(tier)
	trl.limiter = NewRateLimiter(config)
}

// GetTier returns the current tier
func (trl *TierRateLimiter) GetTier() core.Tier {
	return trl.tier
}

// ============================================================================
// Concurrent Connection Limiter
// ============================================================================

// ConnectionLimiter limits concurrent connections
type ConnectionLimiter struct {
	mu           sync.Mutex
	active       map[string]int
	maxPerClient int
	globalMax    int
	globalActive int
}

// NewConnectionLimiter creates a new connection limiter
func NewConnectionLimiter(maxPerClient, globalMax int) *ConnectionLimiter {
	return &ConnectionLimiter{
		active:       make(map[string]int),
		maxPerClient: maxPerClient,
		globalMax:    globalMax,
	}
}

// Acquire attempts to acquire a connection slot
// Returns true if acquired, false if at limit
func (cl *ConnectionLimiter) Acquire(r *http.Request) bool {
	key := DefaultKeyFunc(r)

	cl.mu.Lock()
	defer cl.mu.Unlock()

	// Check global limit
	if cl.globalMax > 0 && cl.globalActive >= cl.globalMax {
		return false
	}

	// Check per-client limit
	if cl.maxPerClient > 0 && cl.active[key] >= cl.maxPerClient {
		return false
	}

	// Increment counters
	cl.active[key]++
	cl.globalActive++

	return true
}

// Release releases a connection slot
func (cl *ConnectionLimiter) Release(r *http.Request) {
	key := DefaultKeyFunc(r)

	cl.mu.Lock()
	defer cl.mu.Unlock()

	if cl.active[key] > 0 {
		cl.active[key]--
	}
	if cl.globalActive > 0 {
		cl.globalActive--
	}
}

// Middleware returns the connection limiting middleware
func (cl *ConnectionLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cl.Acquire(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "too_many_connections", "message": "Too many concurrent connections. Please try again later."}`))
			return
		}

		// Ensure we release on complete
		defer cl.Release(r)

		next.ServeHTTP(w, r)
	})
}

// GetStats returns current connection statistics
func (cl *ConnectionLimiter) GetStats() (perClient map[string]int, global int) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	stats := make(map[string]int)
	for k, v := range cl.active {
		stats[k] = v
	}
	return stats, cl.globalActive
}

// ============================================================================
// Tier-Aware Connection Limiter
// ============================================================================

// TierConnectionLimiter creates a connection limiter based on tier
type TierConnectionLimiter struct {
	limiter *ConnectionLimiter
	tier    core.Tier
}

// NewTierConnectionLimiter creates a new tier-aware connection limiter
func NewTierConnectionLimiter(tier core.Tier) *TierConnectionLimiter {
	limits := tier.GetTierLimits()

	maxConcurrent := limits.MaxConcurrentConnections
	if maxConcurrent == -1 {
		maxConcurrent = 10000 // Effectively unlimited
	}

	return &TierConnectionLimiter{
		limiter: NewConnectionLimiter(maxConcurrent/10, maxConcurrent),
		tier:    tier,
	}
}

// Middleware returns the connection limiting middleware
func (tcl *TierConnectionLimiter) Middleware(next http.Handler) http.Handler {
	return tcl.limiter.Middleware(next)
}

// UpdateTier updates the connection limiter to use new tier limits
func (tcl *TierConnectionLimiter) UpdateTier(tier core.Tier) {
	tcl.tier = tier
	limits := tier.GetTierLimits()

	maxConcurrent := limits.MaxConcurrentConnections
	if maxConcurrent == -1 {
		maxConcurrent = 10000
	}

	tcl.limiter = NewConnectionLimiter(maxConcurrent/10, maxConcurrent)
}

// ============================================================================
// Combined Rate and Connection Limiter
// ============================================================================

// CombinedLimiter combines rate limiting and connection limiting
type CombinedLimiter struct {
	rateLimiter      *TierRateLimiter
	connectionLimiter *TierConnectionLimiter
}

// NewCombinedLimiter creates a new combined limiter
func NewCombinedLimiter(tier core.Tier) *CombinedLimiter {
	return &CombinedLimiter{
		rateLimiter:      NewTierRateLimiter(tier),
		connectionLimiter: NewTierConnectionLimiter(tier),
	}
}

// Middleware returns the combined limiting middleware
// Apply connection limit first, then rate limit
func (cl *CombinedLimiter) Middleware(next http.Handler) http.Handler {
	return cl.connectionLimiter.Middleware(cl.rateLimiter.Middleware(next))
}

// UpdateTier updates both limiters to use new tier limits
func (cl *CombinedLimiter) UpdateTier(tier core.Tier) {
	cl.rateLimiter.UpdateTier(tier)
	cl.connectionLimiter.UpdateTier(tier)
}

// GetTier returns the current tier
func (cl *CombinedLimiter) GetTier() core.Tier {
	return cl.rateLimiter.GetTier()
}

// ============================================================================
// HTTP Handlers for Rate Limit Status
// ============================================================================

// RateLimitStatusHandler returns JSON status about rate limits
type RateLimitStatusHandler struct {
	limiter *CombinedLimiter
}

// NewRateLimitStatusHandler creates a new status handler
func NewRateLimitStatusHandler(limiter *CombinedLimiter) *RateLimitStatusHandler {
	return &RateLimitStatusHandler{limiter: limiter}
}

// ServeHTTP returns rate limit status as JSON
func (h *RateLimitStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stats, globalConns := h.limiter.connectionLimiter.limiter.GetStats()
	config := h.limiter.rateLimiter.limiter.GetRateLimitConfig()

	response := fmt.Sprintf(`{
	"rate_limit": {
		"requests_per_minute": %d,
		"burst": %d
	},
	"connections": {
		"active": %d,
		"by_client": %s
	}
}`, config.RequestsPerMinute, config.Burst, globalConns, formatClientStats(stats))

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}

func formatClientStats(stats map[string]int) string {
	result := "{"
	first := true
	for k, v := range stats {
		if !first {
			result += ","
		}
		first = false
		result += fmt.Sprintf(`"%s":%d`, k, v)
	}
	result += "}"
	return result
}

// ============================================================================
// Usage Examples
// ============================================================================

/*
// Example: Basic usage with tier-based limits
func Example() {
	// Create a rate limiter for Community tier (200/min, 5 concurrent)
	limiter := NewCombinedLimiter(core.TierCommunity)

	// Wrap your HTTP handler
	handler := limiter.Middleware(yourHandler)

	// Or use with tier updates (e.g., when user upgrades)
	limiter.UpdateTier(core.TierDeveloper) // Now 1000/min, 25 concurrent
	limiter.UpdateTier(core.TierProfessional) // Now 5000/min, 100 concurrent
}
*/
