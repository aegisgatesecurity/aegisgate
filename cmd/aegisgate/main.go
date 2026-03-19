// SPDX-License-Identifier: MIT
// =========================================================================
// PROPRIETARY - AegisGate Security
// Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
// =========================================================================
//
// This file contains proprietary trade secret information.
// Unauthorized reproduction, distribution, or reverse engineering is prohibited.
// =========================================================================

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
const version = "v1.0.14"
const commit = "dev"

var date = time.Now().Format(time.RFC3339)

// Management endpoints that should NOT be proxied
var managementEndpoints = map[string]bool{
	"/health":  true,
	"/version": true,
	"/stats":   true,
}

func main() {
	showVersion := flag.Bool("version", false, "Show version information")
	tierName := flag.String("tier", "community", "License tier (requires LICENSE_KEY for Developer/Professional/Enterprise)")
	licenseKey := flag.String("license-key", os.Getenv("LICENSE_KEY"), "License key (required for non-community tiers)")
	targetURL := flag.String("target", "https://api.openai.com", "Target AI provider URL")
	bindAddress := flag.String("bind", "0.0.0.0:8080", "Bind address")
	flag.Parse()

	if *showVersion {
		fmt.Printf("AegisGate %s (commit: %s, date: %s)\n", version, commit, date)
		os.Exit(0)
	}

	tier := core.GetTierByName(*tierName)
	
	// SECURITY: Require license key for non-community tiers
	if *tierName != "community" {
		if *licenseKey == "" {
			log.Fatal("ERROR: License key is required for non-community tiers. Set LICENSE_KEY environment variable or use -license-key flag.")
		}
		// Validate the license with the admin panel
		validator := middleware.GetGlobalValidator()
		ctx := context.Background()
		result, err := validator.Validate(ctx)
		if err != nil || !result.Valid {
			log.Fatalf("ERROR: Invalid license key: %v", err)
		}
		// Use tier from license, not CLI flag (prevents bypass)
		tier = result.Tier
		log.Printf("License validated: tier=%s, expires=%s", tier.String(), result.ExpiresAt.Format("2006-01-02"))
	}
	
	log.Printf("Starting AegisGate v%s with tier: %s", version, tier.String())

	limits := tier.GetTierLimits()
	log.Printf("Rate limit: %s requests/min", limits.FormatLimit("MaxRequestsPerMinute"))
	log.Printf("Max users: %s", limits.FormatLimit("MaxUsers"))

	metricsMgr := metrics.NewManager(nil)

	proxyOpts := &proxy.Options{
		BindAddress:                    *bindAddress,
		Upstream:                       *targetURL,
		MaxBodySize:                    10 * 1024 * 1024,
		Timeout:                        30 * time.Second,
		RateLimit:                      limits.MaxRequestsPerMinute,
		EnableMLDetection:              tier != core.TierCommunity,
		MLSensitivity:                  "medium",
		EnablePromptInjectionDetection: tier != core.TierCommunity,
		EnableContentAnalysis:          tier == core.TierProfessional || tier == core.TierEnterprise,
		EnableBehavioralAnalysis:       tier == core.TierEnterprise,
	}

	proxyServer := proxy.New(proxyOpts)

	// Create single mux that handles both management and proxy routes
	mux := http.NewServeMux()

	// Management endpoints - handle locally
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		health := proxyServer.GetHealth()
		w.Header().Set("Content-Type", "application/json")

		enabled, ok := health["enabled"].(bool)
		if !ok || !enabled {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "{\"status\":\"unhealthy\"}")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"status\":\"healthy\",\"tier\":\"%s\"}", tier.String())
	})

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-AegisGate-Version", version)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"version\":\"%s\",\"commit\":\"%s\"}", version, commit)
	})

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := proxyServer.GetStats()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		reqCount, _ := stats["request_count"].(int64)
		fmt.Fprintf(w, "{\"request_count\":%d}", reqCount)
	})

	// Forward everything else to proxy handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Don't proxy management endpoints
		if managementEndpoints[r.URL.Path] {
			// Already handled above, but safeguard
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Use proxy as handler
		proxyServer.ServeHTTP(w, r)
	})

	handler := Chain(mux, tier, metricsMgr)

	httpServer := &http.Server{
		Addr:         *bindAddress,
		Handler:      handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Proxy server configured, forwarding to: %s", *targetURL)
	log.Printf("HTTP server listening on %s", *bindAddress)

	// Start single HTTP server (NOT proxyServer.Start() which creates a separate server!)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop proxy (cleanup rate limiters, etc)
	proxyServer.Stop(ctx)

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("AegisGate stopped")
}

// Chain chains together the HTTP handler with tier-based middleware and metrics.
func Chain(h http.Handler, tier core.Tier, metricsMgr *metrics.Manager) http.Handler {
	h = middleware.TierBasedResponse(h)
	return h
}
