# Makefile for skeleton-testkit
# Supports Phase 1 Testing Strategy implementation

.PHONY: help test test-unit test-integration test-performance test-ci clean build lint

# Default target
help:
	@echo "Available targets:"
	@echo "  test              - Run all tests"
	@echo "  test-unit         - Run unit tests only"
	@echo "  test-integration  - Run integration tests"
	@echo "  test-performance  - Run performance benchmarks"
	@echo "  test-ci           - Run CI/CD validation tests"
	@echo "  build             - Build the project"
	@echo "  lint              - Run linters"
	@echo "  clean             - Clean build artifacts"

# Build the project
build:
	@echo "Building skeleton-testkit..."
	go build ./...

# Run all tests
test: test-unit test-integration

# Run unit tests (fast tests without external dependencies)
test-unit:
	@echo "Running unit tests..."
	go test -v ./pkg/...

# Run integration tests (requires Docker)
test-integration:
	@echo "Running integration tests..."
	@echo "Note: This requires Docker to be available"
	go test -v -tags=integration ./test/integration/... -timeout=10m

# Run performance benchmarks
test-performance:
	@echo "Running performance benchmarks..."
	go test -v -tags=integration -bench=. -benchmem ./test/integration/performance_test.go -timeout=15m

# Run CI/CD validation tests
test-ci:
	@echo "Running CI/CD validation tests..."
	go test -v -tags=integration ./test/integration/ci_test.go -timeout=10m

# Run integration tests with short flag (for quick validation)
test-integration-short:
	@echo "Running integration tests (short mode)..."
	go test -v -tags=integration -short ./test/integration/... -timeout=5m

# Run specific integration test files
test-basic:
	@echo "Running basic integration tests..."
	go test -v -tags=integration ./test/integration/basic_test.go -timeout=5m

test-database:
	@echo "Running database integration tests..."
	go test -v -tags=integration ./test/integration/database_test.go -timeout=10m

test-setup:
	@echo "Running integration test setup validation..."
	go test -v -tags=integration ./test/integration/integration_test.go -timeout=5m

# Lint the code
lint:
	@echo "Running linters..."
	go vet ./...
	go fmt ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	go clean ./...
	go clean -testcache

# Development targets

# Run tests in watch mode (requires entr or similar tool)
test-watch:
	@echo "Running tests in watch mode..."
	find . -name "*.go" | entr -c make test-unit

# Run integration tests with verbose output and race detection
test-integration-verbose:
	@echo "Running integration tests with verbose output..."
	go test -v -race -tags=integration ./test/integration/... -timeout=15m

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./pkg/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# CI/CD specific targets

# Target for CI environments - runs all tests with appropriate timeouts
ci-test:
	@echo "Running all tests for CI..."
	go test -v ./pkg/... -timeout=5m
	go test -v -tags=integration ./test/integration/... -timeout=20m

# Target for CI environments - runs only fast tests
ci-test-fast:
	@echo "Running fast tests for CI..."
	go test -v ./pkg/... -timeout=3m
	go test -v -tags=integration -short ./test/integration/... -timeout=10m

# Validate that the project builds and basic tests pass
validate:
	@echo "Validating project..."
	go build ./...
	go test -v ./pkg/... -timeout=2m
	@echo "Validation complete"

# Docker-related targets

# Check if Docker is available
check-docker:
	@echo "Checking Docker availability..."
	@docker version > /dev/null 2>&1 || (echo "Docker is not available. Integration tests require Docker." && exit 1)
	@echo "Docker is available"

# Pull required Docker images for testing
pull-test-images: check-docker
	@echo "Pulling test images..."
	docker pull hello-world:latest || echo "Warning: Could not pull hello-world image"
	@echo "Test images pulled"

# Run integration tests with Docker check
test-integration-docker: check-docker test-integration

# Development setup

# Setup development environment
setup-dev:
	@echo "Setting up development environment..."
	go mod download
	go mod tidy
	@echo "Development environment setup complete"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development tools installed"

# Run comprehensive linting
lint-comprehensive:
	@echo "Running comprehensive linting..."
	golangci-lint run ./...

# Benchmark specific tests
benchmark-startup:
	@echo "Running startup benchmarks..."
	go test -v -tags=integration -bench=BenchmarkSkeletonAppStartup -benchmem ./test/integration/performance_test.go

benchmark-database:
	@echo "Running database benchmarks..."
	go test -v -tags=integration -bench=BenchmarkDatabaseContainerStartup -benchmem ./test/integration/performance_test.go

# Help for specific test categories
help-testing:
	@echo "Testing targets:"
	@echo "  test-unit              - Unit tests (no external dependencies)"
	@echo "  test-integration       - Integration tests (requires Docker)"
	@echo "  test-performance       - Performance benchmarks"
	@echo "  test-ci                - CI/CD validation tests"
	@echo "  test-integration-short - Quick integration tests"
	@echo "  test-basic             - Basic integration tests only"
	@echo "  test-database          - Database integration tests only"
	@echo "  test-setup             - Test environment setup validation"
	@echo ""
	@echo "CI/CD targets:"
	@echo "  ci-test                - All tests for CI (with longer timeouts)"
	@echo "  ci-test-fast           - Fast tests for CI"
	@echo "  validate               - Quick validation (build + unit tests)"
	@echo ""
	@echo "Development targets:"
	@echo "  test-watch             - Run tests in watch mode"
	@echo "  test-coverage          - Run tests with coverage report"
	@echo "  setup-dev              - Setup development environment" 