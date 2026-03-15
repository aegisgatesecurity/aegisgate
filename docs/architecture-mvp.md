# AegisGate MVP Architecture

## Overview
AegisGate MVP is a hardened reverse proxy with TLS termination designed for AI chatbot security gateway functionality.

## Components

### Core Packages (3)

1. **pkg/config** - Environment-based configuration
   - Load from environment variables
   - Simplified config structure
   - No file-based config for MVP

2. **pkg/proxy** - Hardened reverse proxy
   - HTTP/HTTPS reverse proxy
   - Request/response inspection hooks
   - Security hardening (rate limits, timeouts)
   - Structured logging

3. **pkg/tls** - TLS termination + self-signed CA
   - Self-signed certificate generation
   - External certificate support
   - HTTPS termination
   - Certificate management interface

### UI (Web-based)

- Single-page web interface
- 3 screens: Dashboard, Certificates, Settings
- Pure HTML/CSS/JS (no frameworks)
- Embedded static files

## Data Flow

Client → HTTPS → AegisGate Proxy → HTTP → Upstream AI Service

## Security Features (MVP)

- Request/response size limits
- Connection timeouts
- Rate limiting
- Graceful shutdown
- Structured logging
- TLS 1.2+ only
- Security headers

## Out of Scope (Phase 2)

- MITRE ATLAS mapping
- Compliance frameworks
- Firecracker orchestration
- Metrics collection
- OPSEC framework
- Policy management
- Alerting system
- Multi-tenant support

## Architecture Diagram

`
┌─────────────┐     HTTPS      ┌─────────────────┐     HTTP      ┌─────────────────┐
│   Client    │ ─────────────→ │  AegisGate Proxy  │ ─────────────→ │  Upstream AI    │
│             │                │                 │                │  Service        │
└─────────────┘                │  • TLS Term     │                └─────────────────┘
                             │  • Inspection   │
                             │  • Rate Limit   │
                             └─────────────────┘
                                      │
                                      ↓
                              ┌─────────────────┐
                              │  Web Admin UI   │
                              │  (3 screens)    │
                              └─────────────────┘
`

