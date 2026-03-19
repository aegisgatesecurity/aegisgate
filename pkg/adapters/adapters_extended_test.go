// SPDX-License-Identifier: MIT
// =========================================================================
// PROPRIETARY - AegisGate Security
// Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
// =========================================================================

package adapters_test

import (
	"testing"
	"time"

	"github.com/aegisgatesecurity/aegisgate/pkg/adapters"
)

func TestNewAuthModule(t *testing.T) {
	module := adapters.NewAuthModule()
	if module == nil {
		t.Fatal("NewAuthModule should return non-nil")
	}
	
	provides := module.Provides()
	if len(provides) == 0 {
		t.Error("Provides should return non-empty list")
	}
}

func TestAuthModuleGetManager(t *testing.T) {
	module := adapters.NewAuthModule()
	manager := module.GetManager()
	// Manager may be nil if not initialized - just verify method exists
	_ = manager
}

func TestNewDashboardModule(t *testing.T) {
	module := adapters.NewDashboardModule()
	if module == nil {
		t.Fatal("NewDashboardModule should return non-nil")
	}
	
	provides := module.Provides()
	if len(provides) == 0 {
		t.Error("Provides should return non-empty list")
	}
	
	deps := module.Dependencies()
	_ = deps
	
	optDeps := module.OptionalDependencies()
	_ = optDeps
}

func TestDashboardModuleSetI18nManager(t *testing.T) {
	module := adapters.NewDashboardModule()
	module.SetI18nManager(nil)
}

func TestDashboardModuleGetDashboard(t *testing.T) {
	module := adapters.NewDashboardModule()
	dashboard := module.GetDashboard()
	_ = dashboard
}

func TestNewI18nModule(t *testing.T) {
	module := adapters.NewI18nModule()
	if module == nil {
		t.Fatal("NewI18nModule should return non-nil")
	}
	
	provides := module.Provides()
	if len(provides) == 0 {
		t.Error("Provides should return non-empty list")
	}
}

func TestI18nModuleGetManager(t *testing.T) {
	module := adapters.NewI18nModule()
	manager := module.GetManager()
	_ = manager
}

func TestNewMetricsModule(t *testing.T) {
	module := adapters.NewMetricsModule()
	if module == nil {
		t.Fatal("NewMetricsModule should return non-nil")
	}
	
	provides := module.Provides()
	if len(provides) == 0 {
		t.Error("Provides should return non-empty list")
	}
}

func TestMetricsModuleGetCollector(t *testing.T) {
	module := adapters.NewMetricsModule()
	collector := module.GetCollector()
	_ = collector
}

func TestNewProxyModule(t *testing.T) {
	module := adapters.NewProxyModule()
	if module == nil {
		t.Fatal("NewProxyModule should return non-nil")
	}
	
	provides := module.Provides()
	if len(provides) == 0 {
		t.Error("Provides should return non-empty list")
	}
	
	deps := module.Dependencies()
	_ = deps
	
	optDeps := module.OptionalDependencies()
	_ = optDeps
}

func TestProxyModuleGetHandler(t *testing.T) {
	module := adapters.NewProxyModule()
	handler := module.GetHandler()
	_ = handler
}

func TestNewScannerModule(t *testing.T) {
	module := adapters.NewScannerModule()
	if module == nil {
		t.Fatal("NewScannerModule should return non-nil")
	}
	
	provides := module.Provides()
	if len(provides) == 0 {
		t.Error("Provides should return non-empty list")
	}
}

func TestScannerModuleGetScanner(t *testing.T) {
	module := adapters.NewScannerModule()
	scanner := module.GetScanner()
	_ = scanner
}

func TestNewTLSModule(t *testing.T) {
	module := adapters.NewTLSModule()
	if module == nil {
		t.Fatal("NewTLSModule should return non-nil")
	}
	
	provides := module.Provides()
	if len(provides) == 0 {
		t.Error("Provides should return non-empty list")
	}
}

func TestTLSModuleGetManager(t *testing.T) {
	module := adapters.NewTLSModule()
	manager := module.GetManager()
	_ = manager
}

func TestModuleNowVariable(t *testing.T) {
	// Test the Now variable for time dependency injection
	originalNow := adapters.Now
	defer func() { adapters.Now = originalNow }()
	
	// Set a fixed time
	fixedTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	adapters.Now = func() time.Time { return fixedTime }
	
	// Now() should return our fixed time
	if !adapters.Now().Equal(fixedTime) {
		t.Error("Now variable should return fixed time")
	}
}

func TestAllModuleTypes(t *testing.T) {
	// Verify all module types can be created
	modules := []struct {
		name   string
		create func() interface{}
	}{
		{"AuthModule", func() interface{} { return adapters.NewAuthModule() }},
		{"DashboardModule", func() interface{} { return adapters.NewDashboardModule() }},
		{"I18nModule", func() interface{} { return adapters.NewI18nModule() }},
		{"MetricsModule", func() interface{} { return adapters.NewMetricsModule() }},
		{"ProxyModule", func() interface{} { return adapters.NewProxyModule() }},
		{"ScannerModule", func() interface{} { return adapters.NewScannerModule() }},
		{"TLSModule", func() interface{} { return adapters.NewTLSModule() }},
	}
	
	for _, m := range modules {
		t.Run(m.name, func(t *testing.T) {
			module := m.create()
			if module == nil {
				t.Errorf("%s should not be nil", m.name)
			}
		})
	}
}

func TestAuthModuleProvidesContent(t *testing.T) {
	module := adapters.NewAuthModule()
	provides := module.Provides()
	
	// Check that at least one capability is returned
	found := false
	for _, p := range provides {
		if len(p) > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Provides should return at least one non-empty capability")
	}
}

func TestDashboardModuleProvidesContent(t *testing.T) {
	module := adapters.NewDashboardModule()
	provides := module.Provides()
	
	for _, p := range provides {
		if len(p) == 0 {
			t.Error("Each provided capability should have a name")
		}
	}
}

func TestProxyModuleDependenciesContent(t *testing.T) {
	module := adapters.NewProxyModule()
	deps := module.Dependencies()
	optDeps := module.OptionalDependencies()
	
	// Dependencies may be empty or have values
	_ = deps
	_ = optDeps
}