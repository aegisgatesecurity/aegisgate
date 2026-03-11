# AegisGate Authentication Package

Enterprise-grade authentication system for AegisGate.

## Features
- OAuth 2.0 (Google, Microsoft, GitHub, Okta, Auth0)
- SAML 2.0 support (Azure, Okta)
- Local username/password authentication
- Role-Based Access Control (RBAC)
- Session management with secure cookies
- Zero external dependencies

## Files
- auth.go - Core types and configuration
- utils.go - Helper functions
- session.go - Session management
- local.go - Local authentication
- oauth.go - OAuth 2.0 implementation
- middleware.go - HTTP middleware
- handlers.go - HTTP handlers
- auth_test.go - Unit tests

## Status
Build: SUCCESS
Vet: PASSED
Dependencies: Zero external (Go standard library only)
