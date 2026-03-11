// SPDX-License-Identifier: MIT
// AegisGate - Chatbot Security Gateway
// Copyright (c) 2026 John Colvin <john.colvin@securityfirm.com>
// See LICENSE file for details.

package cli

import (
	"testing"

	"github.com/aegisgatesecurity/aegisgate/pkg/config"
)

func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("New() returned nil")
	}
	if c.Version != "0.2.0" {
		t.Errorf("Expected version 0.2.0, got %s", c.Version)
	}
}

func TestParseArgs(t *testing.T) {
	c := New()
	
	tests := []struct {
		name  string
		args  []string
		debug bool
	}{
		{
			name:  "no args",
			args:  []string{},
			debug: false,
		},
		{
			name:  "with debug flag",
			args:  []string{"-d"},
			debug: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.ParseArgs(tt.args)
			if err != nil {
				t.Errorf("ParseArgs() error = %v", err)
			}
			if c.Debug != tt.debug {
				t.Errorf("Debug = %v, want %v", c.Debug, tt.debug)
			}
		})
	}
}

func TestCLIWithConfig(t *testing.T) {
	c := New()
	
	// Set up a test config
	cfg := &config.Config{
		BindAddress: ":8080",
		Proxy: config.ProxyConfig{
			Upstream: "http://localhost:3000",
			Timeout:  30,
		},
	}
	c.Config = cfg

	if c.Config.BindAddress != ":8080" {
		t.Errorf("Config not set correctly")
	}
}

func TestPrintVersion(t *testing.T) {
	c := New()
	
	// Test that PrintVersion doesn't panic
	c.PrintVersion()
}

func TestShutdown(t *testing.T) {
	c := New()
	
	err := c.Shutdown()
	if err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}
