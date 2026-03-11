# Padlock Makefile

.PHONY: all build test clean dev docker-build docker-run docker-stop help

# Default target
all: help

# Build the binary
build:
	@echo "Building Padlock..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/padlock ./cmd/padlock
	@echo "Build complete: bin/padlock"

# Build for multiple platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/padlock-linux-amd64 ./cmd/padlock
	GOOS=linux GOARCH=arm64 go build -o bin/padlock-linux-arm64 ./cmd/padlock
	GOOS=darwin GOARCH=amd64 go build -o bin/padlock-darwin-amd64 ./cmd/padlock
	GOOS=darwin GOARCH=arm64 go build -o bin/padlock-darwin-arm64 ./cmd/padlock
	GOOS=windows GOARCH=amd64 go build -o bin/padlock-windows-amd64.exe ./cmd/padlock
	@echo "All builds complete in bin/"

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Tests complete"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f *.test
	@echo "Clean complete"

# Run linter
lint:
	@echo "Running linter..."
	go install golang.org/x/lint/golint@latest
	golint -set_exit_status ./...
	@echo "Lint complete"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .
	@echo "Format complete"

# Development mode (run with live reload)
dev:
	@echo "Starting development mode..."
	@echo "Note: Implement live reload as needed"
	go run ./cmd/padlock serve --config ./config/community.env

# Run Docker container (local build)
docker-build:
	@echo "Building Docker image..."
	docker build -t padlock/padlock:latest -f deploy/docker/Dockerfile .

# Run Docker Compose
docker-run:
	@echo "Starting Padlock with Docker Compose..."
	cd deploy/docker && docker-compose up -d

# Stop Docker Compose
docker-stop:
	@echo "Stopping Padlock..."
	cd deploy/docker && docker-compose down

# Docker development
docker-dev:
	docker-compose -f deploy/docker/docker-compose.yml up -d --build

# Run security tests
security:
	@echo "Running security tests..."
	go test -v ./tests/security/...

# Run integration tests
integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Verify dependencies
deps-verify:
	@echo "Verifying dependencies..."
	go mod verify
	go mod why

# Generate code (protobuf, etc.)
generate:
	@echo "Generating code..."
	@# Add code generation commands here as needed
	@echo "Code generation complete"

# Show help
help:
	@echo "Padlock Build System"
	@echo "===================="
	@echo ""
	@echo "Available targets:"
	@echo "  make build          - Build the Padlock binary"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make lint          - Run linter"
	@echo "  make fmt           - Format code"
	@echo "  make dev           - Run in development mode"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Start with Docker Compose"
	@echo "  make docker-stop   - Stop Docker Compose"
	@echo "  make security      - Run security tests"
	@echo "  make integration   - Run integration tests"
	@echo "  make deps          - Install dependencies"
	@echo "  make generate      - Generate code"
	@echo ""
	@echo "Example usage:"
	@echo "  make build         # Build the binary"
	@echo "  make docker-run    # Start with Docker"
	@echo "  make test          # Run all tests"
