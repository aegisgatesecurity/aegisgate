# ATLAS Compliance Test Makefile

.PHONY: help test test-verbose test-coverage test-race bench clean install

# Default target
help:
	@echo "ATLAS Compliance Test Suite"
	@echo ""
	@echo "Available targets:"
	@echo "  make test          - Run ATLAS compliance tests"
	@echo "  make test-verbose  - Run tests with verbose output"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make test-race     - Run tests with race detector"
	@echo "  make bench         - Run benchmarks"
	@echo "  make clean         - Clean test artifacts"
	@echo "  make install       - Install test dependencies"

# Run ATLAS compliance tests
test:
	@echo "Running ATLAS compliance tests..."
	cd $(shell pwd) && go test -v -run Atlas ./tests/integration/...

# Run tests with verbose output
test-verbose:
	@echo "Running ATLAS compliance tests (verbose)..."
	cd $(shell pwd) && go test -v -count=1 -run Atlas ./tests/integration/...

# Run tests with coverage
test-coverage:
	@echo "Running ATLAS compliance tests with coverage..."
	cd $(shell pwd) && go test -coverprofile=coverage.out -covermode=atomic -run Atlas ./tests/integration/...
	@echo ""
	@echo "Coverage Report:"
	@go tool cover -func=coverage.out

# Run tests with race detector
test-race:
	@echo "Running ATLAS compliance tests with race detector..."
	cd $(shell pwd) && go test -race -run Atlas ./tests/integration/...

# Run benchmarks
bench:
	@echo "Running ATLAS compliance benchmarks..."
	cd $(shell pwd) && go test -bench=. -benchmem ./tests/integration/...

# Clean test artifacts
clean:
	@echo "Cleaning test artifacts..."
	rm -f coverage.out
	rm -f test_*.log
	rm -f test_*.tmp

# Install test dependencies
install:
	@echo "Installing test dependencies..."
	go get -t ./tests/integration/...
	go mod tidy
