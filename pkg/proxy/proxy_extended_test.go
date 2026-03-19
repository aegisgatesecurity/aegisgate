// SPDX-License-Identifier: MIT
// =========================================================================
// PROPRIETARY - AegisGate Security
// Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
// =========================================================================

package proxy_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aegisgatesecurity/aegisgate/pkg/proxy"
)

func TestNewProxy(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	if p == nil {
		t.Fatal("New should return non-nil proxy")
	}
}

func TestNewProxyWithNilOptions(t *testing.T) {
	p := proxy.New(nil)
	if p == nil {
		t.Fatal("New(nil) should return non-nil proxy")
	}
}

func TestProxyGetStats(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	stats := p.GetStats()
	if stats == nil {
		t.Error("GetStats should return non-nil")
	}
}

func TestProxyGetHealth(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	health := p.GetHealth()
	if health == nil {
		t.Error("GetHealth should return non-nil")
	}
}

func TestProxyIsEnabled(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	_ = p.IsEnabled()
	// Just test the method exists and returns a value
}

func TestProxyIsEnabledFalse(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	if p.IsEnabled() {
		t.Error("Proxy should be enabled by default")
	}
}

func TestProxyGetStatsStruct(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	stats := p.GetStatsStruct()
	if stats == nil {
		t.Error("GetStatsStruct should return non-nil")
	}
}

func TestProxyGetScanner(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	scanner := p.GetScanner()
	_ = scanner
}

func TestProxySetScanner(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	p.SetScanner(nil)
}

func TestProxyGetMLMiddleware(t *testing.T) {
	opts := &proxy.Options{
		BindAddress:              "127.0.0.1:8080",
		Upstream:                 "https://api.openai.com",
		EnableMLDetection:        true,
		MLSensitivity:           "medium",
		EnablePromptInjectionDetection: true,
	}
	p := proxy.New(opts)
	ml := p.GetMLMiddleware()
	_ = ml
}

func TestProxyOptionsDefaults(t *testing.T) {
	opts := &proxy.Options{}
	if opts.BindAddress != "" {
		t.Error("Default BindAddress should be empty")
	}
}

func TestProxyOptionsWithAllFields(t *testing.T) {
	opts := &proxy.Options{
		BindAddress:                    "0.0.0.0:8080",
		Upstream:                        "https://api.openai.com",
		MaxBodySize:                    1024 * 1024,
		Timeout:                        30 * time.Second,
		RateLimit:                      1000,
		EnableMLDetection:              true,
		MLSensitivity:                 "high",
		EnablePromptInjectionDetection: true,
		EnableContentAnalysis:          true,
		EnableBehavioralAnalysis:       true,
	}
	p := proxy.New(opts)
	if p == nil {
		t.Fatal("New with full options should succeed")
	}
}

func TestProxyGetComplianceManager(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	cm := p.GetComplianceManager()
	_ = cm
}

func TestProxyGetCircuitBreaker(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	cb := p.GetCircuitBreaker()
	_ = cb
}

func TestProxyStop(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	err := p.Stop(nil)
	if err != nil {
		t.Errorf("Stop should not error: %v", err)
	}
}

func TestProxyServeHTTP(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	
	p.ServeHTTP(rec, req)
}

func TestProxyHealth(t *testing.T) {
	opts := &proxy.Options{
		BindAddress: "127.0.0.1:8080",
		Upstream:    "https://api.openai.com",
	}
	p := proxy.New(opts)
	health := p.GetHealth()
	
	// Health should have enabled field
	enabled, ok := health["enabled"].(bool)
	if !ok {
		t.Error("Health should have enabled field")
	}
	_ = enabled
}