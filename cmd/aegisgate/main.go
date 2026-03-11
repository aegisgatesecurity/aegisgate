// AegisGate - Enterprise AI API Security Platform
// Main entry point for the AegisGate service
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aegisgatesecurity/aegisgate/pkg/core"
	"github.com/aegisgatesecurity/aegisgate/pkg/metrics"
	"github.com/aegisgatesecurity/aegisgate/pkg/middleware"
	"github.com/aegisgatesecurity/aegisgate/pkg/proxy"
)

// Build info - set during build
const version = "1.0.0"
const commit  = "dev"

// date is set at runtime
var date = time.Now().Format(time.RFC3339)

func main() {
	// Parse command line flags
	showVersion := flag.Bool("version", false, "Show version information")
	tierName := flag.String("tier", "community", "License tier (community, developer, professional, enterprise)")
	targetURL := flag.String("target", "https://api.openai.com", "Target AI provider URL")
	bindAddress := flag.String("bind", "0.0.0.0:8080", "Bind address")
	flag.Parse()

	// Show version and exit
	if *showVersion {
		fmt.Printf("AegisGate %s (commit: %s, date: %s)\n", version, commit, date)
		os.Exit(0)
	}

	// Parse tier
	tier := core.GetTierByName(*tierName)
	log.Printf("Starting AegisGate v%s with tier: %s", version, tier.String())

	// Get limits for the tier
	limits := tier.GetTierLimits()
	log.Printf("Rate limit: %s requests/min", limits.FormatLimit("MaxRequestsPerMinute"))
	log.Printf("Max users: %s", limits.FormatLimit("MaxUsers"))

	// Create metrics manager
	metricsMgr := metrics.NewManager(nil)

	// Create proxy server with tier-based configuration
	proxyOpts := &proxy.Options{
		BindAddress:         *bindAddress,
		Upstream:            *targetURL,
		MaxBodySize:         10 * 1024 * 1024, // 10MB
		Timeout:             30 * time.Second,
		RateLimit:           limits.MaxRequestsPerMinute,
		EnableMLDetection:   tier != core.TierCommunity, // ML starts at Developer
		MLSensitivity:       "medium",
		EnablePromptInjectionDetection: tier != core.TierCommunity,
		EnableContentAnalysis:         tier == core.TierProfessional || tier == core.TierEnterprise,
		EnableBehavioralAnalysis:     tier == core.TierEnterprise,
	}

	proxyServer := proxy.New(proxyOpts)

	// Create HTTP server
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		health := proxyServer.GetHealth()
		w.Header().Set("Content-Type", "application/json")
		
		// Check if enabled
		enabled, ok := health["enabled"].(bool)
		if !ok || !enabled {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, `{"status":"unhealthy","tier":"%s","proxy_enabled":false}`, tier.String())
			return
		}
		
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","tier":"%s","proxy_enabled":true}`, tier.String())
	})

	// Version endpoint
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"version":"%s","commit":"%s","date":"%s"}`, version, commit, date)
	})

	// Stats endpoint
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := proxyServer.GetStats()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		reqCount, _ := stats["request_count"].(int64)
		fmt.Fprintf(w, `{"request_count":%d,"enabled":%v}`, reqCount, stats["enabled"])
	})

	// Proxy handler - forward to AI provider
	mux.HandleFunc("/v1/", func(w http.ResponseWriter, r *http.Request) {
		proxyServer.ServeHTTP(w, r)
	})

	// Also handle root path for proxy
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Skip if it's one of our management endpoints
		if r.URL.Path != "/health" && r.URL.Path != "/version" && r.URL.Path != "/stats" {
			proxyServer.ServeHTTP(w, r)
		}
	})

	// Apply middleware
	handler := Chain(mux, tier, metricsMgr)

	httpServer := &http.Server{
		Addr:         *bindAddress,
		Handler:      handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start proxy
	if err := proxyServer.Start(); err != nil {
		log.Fatalf("Failed to start proxy: %v", err)
	}
	log.Printf("Proxy server started, forwarding to: %s", *targetURL)

	// Start HTTP server in goroutine
	go func() {
		log.Printf("HTTP server listening on %s", *bindAddress)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	proxyServer.Stop(ctx)
	
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("AegisGate stopped")
}

// Chain applies middleware in order
func Chain(h http.Handler, tier core.Tier, metricsMgr *metrics.Manager) http.Handler {
	// Feature gating - adds tier headers to responses
	h = middleware.TierBasedResponse(h)

	return h
}
