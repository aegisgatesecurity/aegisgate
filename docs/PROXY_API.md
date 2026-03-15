# AegisGate Proxy API Documentation

## Overview

The AegisGate proxy package provides a hardened reverse proxy with built-in security scanning, rate limiting, and compliance checking for AI chatbot APIs.

## Quick Start

package main

import (
    "context"
    "log"
    "github.com/aegisgatesecurity/aegisgate/pkg/proxy"
)

func main() {
    p := proxy.New(&proxy.Options{
        BindAddress: ":8080",
        Upstream:    "http://localhost:3000",
        RateLimit:   100,
    })
    
    go func() {
        if err := p.Start(); err != nil {
            log.Fatal(err)
        }
    }()
    
    p.Stop(context.Background())
}

---

## Configuration

### Options

type Options struct {
    BindAddress string        // Listen address (e.g., ":8080")
    Upstream    string        // Upstream URL
    TLS         *TLSConfig   // TLS configuration (optional)
    MaxBodySize int64        // Max request body size (default: 10MB)
    Timeout     time.Duration // Request timeout (default: 30s)
    RateLimit   int          // Requests per minute (default: 100)
}

---

## API Reference

### proxy.New(opts *Options) *Proxy

Creates a new hardened proxy instance.

### proxy.Start() error

Starts the proxy server.

### proxy.Stop(ctx context.Context) error

Gracefully stops the proxy server.

### proxy.ServeHTTP(w http.ResponseWriter, req *http.Request)

Implements http.Handler interface.

### proxy.GetHealth() map[string]interface{}

Returns health status.

### proxy.GetStats() map[string]interface{}

Returns proxy statistics.

### proxy.GetScanner() *scanner.Scanner

Returns the content scanner instance.

### proxy.SetScanner(s *scanner.Scanner)

Sets a custom content scanner.

### proxy.GetComplianceManager() *compliance.ComplianceManager

Returns the compliance manager.

### proxy.IsEnabled() bool

Checks if proxy server is running.

---

## Security Features

The proxy automatically scans for:

- PII (SSN, Credit Cards, Emails)
- Credentials (API Keys, Passwords)
- Malicious patterns (SQL Injection, XSS)
- MITRE ATLAS compliance

---

## Examples

### Basic Proxy

p := proxy.New(&proxy.Options{
    BindAddress: ":8080",
    Upstream:    "http://localhost:3000",
})

### With TLS

p := proxy.New(&proxy.Options{
    BindAddress: ":8443",
    Upstream:    "http://localhost:3000",
    TLS: &proxy.TLSConfig{
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
})
