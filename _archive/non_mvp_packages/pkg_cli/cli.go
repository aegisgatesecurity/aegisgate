package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aegisgatesecurity/aegisgate/pkg/config"
)

// CLI manages command-line interface for AegisGate
type CLI struct {
	Config    *config.Config
	Version   string
	Debug     bool
	ConfigFile string
}

// New creates a new CLI instance
func New() *CLI {
	return &CLI{
		Version: "0.2.0",
	}
}

// ParseArgs parses command-line arguments
func (c *CLI) ParseArgs(args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-v", "--version":
			fmt.Printf("AegisGate v%s\n", c.Version)
			os.Exit(0)
		case "-d", "--debug":
			c.Debug = true
		case "-c", "--config":
			if i+1 < len(args) {
				c.ConfigFile = args[i+1]
				i++
			}
		}
	}
	return nil
}

// LoadConfig loads configuration
func (c *CLI) LoadConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	c.Config = cfg
	return nil
}

// Run starts the application
func (c *CLI) Run() error {
	fmt.Println("AegisGate Chatbot Security Gateway")
	fmt.Println("=================================")
	
	if c.Debug {
		fmt.Printf("Debug mode: enabled\n")
	}
	
	if c.ConfigFile != "" {
		fmt.Printf("Config file: %s\n", c.ConfigFile)
	}
	
	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Start main application logic
	go c.startApplication()
	
	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\nShutting down...")
	return c.Shutdown()
}

func (c *CLI) startApplication() {
	fmt.Println("Starting AegisGate...")
	fmt.Printf("Listening on: %s\n", c.Config.BindAddress)
}

// Shutdown gracefully shuts down the application
func (c *CLI) Shutdown() error {
	fmt.Println("AegisGate shutdown complete")
	return nil
}

// PrintVersion prints the version information
func (c *CLI) PrintVersion() {
	fmt.Printf("AegisGate v%s\n", c.Version)
	fmt.Printf("Go version: %s\n", "Go 1.21+")
	fmt.Printf("Build time: %s\n", "2026-02-12")
}
