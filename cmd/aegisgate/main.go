package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/AegisGate/AegisGate/pkg/adapters/primary/rpc"
	"github.com/AegisGate/AegisGate/pkg/adapters/primary/ui"
	"github.com/AegisGate/AegisGate/pkg/adapters/storage"
	"github.com/AegisGate/AegisGate/pkg/core/domain"
	"github.com/AegisGate/AegisGate/pkg/core/services"
	"github.com/AegisGate/AegisGate/pkg/infrastructure/logging"
	"github.com/AegisGate/AegisGate/pkg/plugins"
	"github.com/AegisGate/AegisGate/pkg/server"
)

const version = "1.0.5"

var (
	showVersion = flag.Bool("version", false, "show version")
	configPath  = flag.String("config", "", "path to configuration file")
	logLevel    = flag.String("log-level", "info", "set log level (debug, info, warn, error)")
	listenAddr  = flag.String("listen", ":8443", "listen address")
	dataDir     = flag.String("data", "./data", "data directory path")
	pluginsDir  = flag.String("plugins", "./plugins", "plugins directory path")
	dbType      = flag.String("db", "sqlite3", "database type (sqlite3, postgres, mysql)")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("AegisGate version %s
", version)
		os.Exit(0)
	}

	logger := logging.NewLogger(*logLevel)

	// Validate required directories
	if err := validatePaths(*dataDir, *pluginsDir); err != nil {
		logger.Error("Path validation failed", "error", err)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := loadConfig(*configPath)
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	cfg.ListenAddr = *listenAddr

	// Initialize storage
	storageAdapter, err := storage.NewAdapter(*dbType, *dataDir)
	if err != nil {
		logger.Error("Failed to initialize storage", "error", err)
		os.Exit(1)
	}

	// Initialize services
	authService := services.NewAuthService(storageAdapter)
	siemService := services.NewSIEMService(storageAdapter)
	proxyService := services.NewProxyService(storageAdapter)
	adminService := services.NewAdminService(storageAdapter, version)

	// Load plugins
	pluginManager := plugins.NewManager(*pluginsDir, logger)
	if err := pluginManager.Load(); err != nil {
		logger.Warn("Plugin loading failed", "error", err)
	}

	// Create server
	srv := server.NewServer(cfg, logger, storageAdapter, authService, siemService, proxyService, adminService, pluginManager)

	// Start RPC server in background
	go func() {
		if err := rpc.StartRPCServer(authService, siemService, proxyService, adminService); err != nil {
			logger.Error("RPC server failed", "error", err)
		}
	}()

	logger.Info("Starting AegisGate", "version", version, "listen", *listenAddr)
	if err := srv.Start(); err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func validatePaths(dataDir, pluginsDir string) error {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	return nil
}

func loadConfig(configPath string) (*domain.Config, error) {
	cfg := domain.DefaultConfig()

	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// Simple TOML parsing would go here
		if strings.Contains(string(data), "[server]") {
			// Parse TOML config
		}
	}

	return cfg, nil
}
