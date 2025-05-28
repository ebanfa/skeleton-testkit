# Skeleton-Testkit Makefile
# Comprehensive testing framework for skeleton-based applications

# Project configuration
PROJECT_NAME := skeleton-testkit
MODULE_NAME := github.com/fintechain/skeleton-testkit
GO_VERSION := 1.21

# Directories
BIN_DIR := bin
CMD_DIR := cmd
PKG_DIR := pkg
INTERNAL_DIR := internal
TEST_DIR := test
DOCS_DIR := docs
COVERAGE_DIR := coverage
EXAMPLES_DIR := examples

# Build configuration
LDFLAGS := -w -s
BUILD_FLAGS := -ldflags "$(LDFLAGS)"
TEST_FLAGS := -race -timeout=30s
INTEGRATION_FLAGS := -tags=integration -timeout=10m
PERFORMANCE_FLAGS := -tags=integration -bench=. -benchmem -timeout=15m
COVERAGE_FLAGS := -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic

# Tools
GOLANGCI_LINT_VERSION := v1.55.2
MOCKGEN_VERSION := v1.6.0

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
NC := \033[0m # No Color

.PHONY: help
help: ## Display this help message
	@echo "$(CYAN)Skeleton-Testkit - Testing Framework for Skeleton Applications$(NC)"
	@echo ""
	@echo "$(YELLOW)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-25s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# =============================================================================
# Build Targets
# =============================================================================

.PHONY: build
build: ## Build all binaries
	@echo "$(BLUE)Building all binaries...$(NC)"
	@mkdir -p $(BIN_DIR)
	@go build $(BUILD_FLAGS) -o $(BIN_DIR)/ ./$(CMD_DIR)/...
	@echo "$(GREEN)✓ Build completed$(NC)"

.PHONY: build-examples
build-examples: ## Build example applications
	@echo "$(BLUE)Building examples...$(NC)"
	@mkdir -p $(BIN_DIR)/examples
	@go build $(BUILD_FLAGS) -o $(BIN_DIR)/examples/ ./$(EXAMPLES_DIR)/...
	@echo "$(GREEN)✓ Examples built$(NC)"

.PHONY: install
install: ## Install binaries to GOPATH/bin
	@echo "$(BLUE)Installing binaries...$(NC)"
	@go install $(BUILD_FLAGS) ./$(CMD_DIR)/...
	@echo "$(GREEN)✓ Binaries installed$(NC)"

# =============================================================================
# Test Targets
# =============================================================================

.PHONY: test
test: ## Run all tests
	@echo "$(BLUE)Running all tests...$(NC)"
	@go test $(TEST_FLAGS) ./...
	@echo "$(GREEN)✓ All tests passed$(NC)"

.PHONY: test-unit
test-unit: ## Run unit tests only (no external dependencies)
	@echo "$(BLUE)Running unit tests...$(NC)"
	@go test $(TEST_FLAGS) ./$(PKG_DIR)/... ./$(INTERNAL_DIR)/...
	@echo "$(GREEN)✓ Unit tests passed$(NC)"

.PHONY: test-integration
test-integration: check-docker ## Run integration tests (requires Docker)
	@echo "$(BLUE)Running integration tests...$(NC)"
	@go test -v $(INTEGRATION_FLAGS) ./$(TEST_DIR)/integration/...
	@echo "$(GREEN)✓ Integration tests passed$(NC)"

.PHONY: test-integration-short
test-integration-short: check-docker ## Run integration tests in short mode
	@echo "$(BLUE)Running integration tests (short mode)...$(NC)"
	@go test -v $(INTEGRATION_FLAGS) -short ./$(TEST_DIR)/integration/...
	@echo "$(GREEN)✓ Short integration tests passed$(NC)"

.PHONY: test-basic
test-basic: check-docker ## Run basic integration tests
	@echo "$(BLUE)Running basic integration tests...$(NC)"
	@go test -v $(INTEGRATION_FLAGS) ./$(TEST_DIR)/integration/basic_test.go
	@echo "$(GREEN)✓ Basic integration tests passed$(NC)"

.PHONY: test-database
test-database: check-docker ## Run database integration tests
	@echo "$(BLUE)Running database integration tests...$(NC)"
	@go test -v $(INTEGRATION_FLAGS) ./$(TEST_DIR)/integration/database_test.go
	@echo "$(GREEN)✓ Database integration tests passed$(NC)"

.PHONY: test-containers
test-containers: check-docker ## Run container management tests
	@echo "$(BLUE)Running container management tests...$(NC)"
	@go test -v $(INTEGRATION_FLAGS) ./$(TEST_DIR)/integration/containers/...
	@echo "$(GREEN)✓ Container tests passed$(NC)"

.PHONY: test-skeleton
test-skeleton: check-docker ## Run skeleton application tests
	@echo "$(BLUE)Running skeleton application tests...$(NC)"
	@go test -v $(INTEGRATION_FLAGS) ./$(TEST_DIR)/integration/skeleton/...
	@echo "$(GREEN)✓ Skeleton application tests passed$(NC)"

.PHONY: test-performance
test-performance: check-docker ## Run performance benchmarks
	@echo "$(BLUE)Running performance benchmarks...$(NC)"
	@go test -v $(PERFORMANCE_FLAGS) ./$(TEST_DIR)/integration/performance_test.go
	@echo "$(GREEN)✓ Performance benchmarks completed$(NC)"

.PHONY: test-ci
test-ci: check-docker ## Run CI/CD validation tests
	@echo "$(BLUE)Running CI/CD validation tests...$(NC)"
	@go test -v $(INTEGRATION_FLAGS) ./$(TEST_DIR)/integration/ci_test.go
	@echo "$(GREEN)✓ CI/CD validation tests passed$(NC)"

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "$(BLUE)Running tests with verbose output...$(NC)"
	@go test -v $(TEST_FLAGS) ./...

.PHONY: test-short
test-short: ## Run tests with short flag (skip long-running tests)
	@echo "$(BLUE)Running short tests...$(NC)"
	@go test -short $(TEST_FLAGS) ./...

.PHONY: test-examples
test-examples: ## Run example tests
	@echo "$(BLUE)Running example tests...$(NC)"
	@go test $(TEST_FLAGS) ./$(EXAMPLES_DIR)/...
	@echo "$(GREEN)✓ Example tests passed$(NC)"

# =============================================================================
# Coverage Targets
# =============================================================================

.PHONY: coverage
coverage: ## Generate test coverage report
	@echo "$(BLUE)Generating coverage report...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test $(COVERAGE_FLAGS) ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1
	@echo "$(GREEN)✓ Coverage report generated: $(COVERAGE_DIR)/coverage.html$(NC)"

.PHONY: coverage-unit
coverage-unit: ## Generate coverage for unit tests only
	@echo "$(BLUE)Generating unit test coverage...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test $(COVERAGE_FLAGS) ./$(PKG_DIR)/... ./$(INTERNAL_DIR)/...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage-unit.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1
	@echo "$(GREEN)✓ Unit test coverage: $(COVERAGE_DIR)/coverage-unit.html$(NC)"

.PHONY: coverage-integration
coverage-integration: check-docker ## Generate coverage for integration tests
	@echo "$(BLUE)Generating integration test coverage...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test $(COVERAGE_FLAGS) $(INTEGRATION_FLAGS) ./$(TEST_DIR)/integration/...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage-integration.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1
	@echo "$(GREEN)✓ Integration test coverage: $(COVERAGE_DIR)/coverage-integration.html$(NC)"

.PHONY: coverage-show
coverage-show: ## Show coverage in browser
	@echo "$(BLUE)Opening coverage report in browser...$(NC)"
	@open $(COVERAGE_DIR)/coverage.html || xdg-open $(COVERAGE_DIR)/coverage.html

# =============================================================================
# Code Quality Targets
# =============================================================================

.PHONY: lint
lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@golangci-lint run ./...
	@echo "$(GREEN)✓ Linting completed$(NC)"

.PHONY: lint-fix
lint-fix: ## Run linter with auto-fix
	@echo "$(BLUE)Running linter with auto-fix...$(NC)"
	@golangci-lint run --fix ./...
	@echo "$(GREEN)✓ Linting with auto-fix completed$(NC)"

.PHONY: fmt
fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Vet completed$(NC)"

.PHONY: mod-tidy
mod-tidy: ## Tidy go modules
	@echo "$(BLUE)Tidying go modules...$(NC)"
	@go mod tidy
	@echo "$(GREEN)✓ Modules tidied$(NC)"

.PHONY: mod-verify
mod-verify: ## Verify go modules
	@echo "$(BLUE)Verifying go modules...$(NC)"
	@go mod verify
	@echo "$(GREEN)✓ Modules verified$(NC)"

.PHONY: mod-download
mod-download: ## Download go modules
	@echo "$(BLUE)Downloading go modules...$(NC)"
	@go mod download
	@echo "$(GREEN)✓ Modules downloaded$(NC)"

# =============================================================================
# Development Targets
# =============================================================================

.PHONY: run-examples
run-examples: build-examples ## Run example applications
	@echo "$(BLUE)Running examples...$(NC)"
	@for example in $(BIN_DIR)/examples/*; do \
		echo "$(CYAN)Running $$example...$(NC)"; \
		$$example || true; \
	done

.PHONY: dev
dev: clean fmt vet lint test build ## Full development cycle (clean, format, vet, lint, test, build)
	@echo "$(GREEN)✓ Development cycle completed$(NC)"

.PHONY: ci
ci: mod-verify fmt vet lint test-short coverage ## CI pipeline (verify, format, vet, lint, test, coverage)
	@echo "$(GREEN)✓ CI pipeline completed$(NC)"

.PHONY: ci-integration
ci-integration: check-docker mod-verify fmt vet lint test-integration coverage-integration ## Full CI with integration tests
	@echo "$(GREEN)✓ Full CI pipeline completed$(NC)"

# =============================================================================
# Docker and Container Targets
# =============================================================================

.PHONY: check-docker
check-docker: ## Check if Docker is available
	@echo "$(BLUE)Checking Docker availability...$(NC)"
	@docker version > /dev/null 2>&1 || (echo "$(RED)Docker is not available. Integration tests require Docker.$(NC)" && exit 1)
	@echo "$(GREEN)✓ Docker is available$(NC)"

.PHONY: pull-test-images
pull-test-images: check-docker ## Pull required Docker images for testing
	@echo "$(BLUE)Pulling test images...$(NC)"
	@docker pull hello-world:latest || echo "$(YELLOW)Warning: Could not pull hello-world image$(NC)"
	@docker pull postgres:15-alpine || echo "$(YELLOW)Warning: Could not pull postgres image$(NC)"
	@docker pull redis:7-alpine || echo "$(YELLOW)Warning: Could not pull redis image$(NC)"
	@echo "$(GREEN)✓ Test images pulled$(NC)"

.PHONY: clean-docker
clean-docker: ## Clean up Docker test containers and images
	@echo "$(BLUE)Cleaning Docker test resources...$(NC)"
	@docker container prune -f || true
	@docker image prune -f || true
	@echo "$(GREEN)✓ Docker cleanup completed$(NC)"

# =============================================================================
# Documentation Targets
# =============================================================================

.PHONY: docs
docs: ## Generate documentation
	@echo "$(BLUE)Generating documentation...$(NC)"
	@mkdir -p $(DOCS_DIR)/api
	@go doc -all ./... > $(DOCS_DIR)/api/api.txt
	@echo "$(GREEN)✓ Documentation generated$(NC)"

.PHONY: docs-serve
docs-serve: ## Serve documentation locally
	@echo "$(BLUE)Serving documentation on http://localhost:6060$(NC)"
	@godoc -http=:6060

.PHONY: docs-examples
docs-examples: ## Generate example documentation
	@echo "$(BLUE)Generating example documentation...$(NC)"
	@find $(EXAMPLES_DIR) -name "*.go" -exec go doc {} \; > $(DOCS_DIR)/examples.txt
	@echo "$(GREEN)✓ Example documentation generated$(NC)"

# =============================================================================
# Mock Generation Targets
# =============================================================================

.PHONY: mocks
mocks: ## Generate mocks for testing
	@echo "$(BLUE)Generating mocks...$(NC)"
	@go generate ./...
	@echo "$(GREEN)✓ Mocks generated$(NC)"

.PHONY: mocks-container
mocks-container: ## Generate container mocks
	@echo "$(BLUE)Generating container mocks...$(NC)"
	@mockgen -source=$(INTERNAL_DIR)/domain/container/container.go -destination=$(INTERNAL_DIR)/domain/container/mocks/container_mock.go
	@echo "$(GREEN)✓ Container mocks generated$(NC)"

.PHONY: mocks-verification
mocks-verification: ## Generate verification mocks
	@echo "$(BLUE)Generating verification mocks...$(NC)"
	@mockgen -source=$(INTERNAL_DIR)/domain/verification/verifier.go -destination=$(INTERNAL_DIR)/domain/verification/mocks/verifier_mock.go
	@echo "$(GREEN)✓ Verification mocks generated$(NC)"

# =============================================================================
# Benchmark Targets
# =============================================================================

.PHONY: bench
bench: check-docker ## Run all benchmarks
	@echo "$(BLUE)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...
	@echo "$(GREEN)✓ Benchmarks completed$(NC)"

.PHONY: bench-startup
bench-startup: check-docker ## Run startup benchmarks
	@echo "$(BLUE)Running startup benchmarks...$(NC)"
	@go test -v $(INTEGRATION_FLAGS) -bench=BenchmarkSkeletonAppStartup -benchmem ./$(TEST_DIR)/integration/performance_test.go

.PHONY: bench-database
bench-database: check-docker ## Run database benchmarks
	@echo "$(BLUE)Running database benchmarks...$(NC)"
	@go test -v $(INTEGRATION_FLAGS) -bench=BenchmarkDatabaseContainerStartup -benchmem ./$(TEST_DIR)/integration/performance_test.go

.PHONY: bench-containers
bench-containers: check-docker ## Run container benchmarks
	@echo "$(BLUE)Running container benchmarks...$(NC)"
	@go test -bench=BenchmarkContainer -benchmem ./$(PKG_DIR)/container/...

# =============================================================================
# Tool Installation Targets
# =============================================================================

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@go install github.com/golang/mock/mockgen@$(MOCKGEN_VERSION)
	@echo "$(GREEN)✓ Development tools installed$(NC)"

.PHONY: check-tools
check-tools: ## Check if required tools are installed
	@echo "$(BLUE)Checking required tools...$(NC)"
	@command -v golangci-lint >/dev/null 2>&1 || { echo "$(RED)golangci-lint not found. Run 'make install-tools'$(NC)"; exit 1; }
	@command -v mockgen >/dev/null 2>&1 || { echo "$(RED)mockgen not found. Run 'make install-tools'$(NC)"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "$(RED)docker not found. Install Docker for integration tests$(NC)"; exit 1; }
	@echo "$(GREEN)✓ All required tools are installed$(NC)"

# =============================================================================
# Development Environment Targets
# =============================================================================

.PHONY: setup-dev
setup-dev: ## Setup development environment
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@go mod download
	@go mod tidy
	@make install-tools
	@make pull-test-images
	@echo "$(GREEN)✓ Development environment setup complete$(NC)"

.PHONY: dev-up
dev-up: ## Start development environment (if using docker-compose)
	@echo "$(BLUE)Starting development environment...$(NC)"
	@docker-compose -f deployments/docker-compose.dev.yml up -d || echo "$(YELLOW)No docker-compose dev environment found$(NC)"

.PHONY: dev-down
dev-down: ## Stop development environment
	@echo "$(BLUE)Stopping development environment...$(NC)"
	@docker-compose -f deployments/docker-compose.dev.yml down || echo "$(YELLOW)No docker-compose dev environment found$(NC)"

.PHONY: dev-shell
dev-shell: ## Enter development container shell
	@echo "$(BLUE)Entering development shell...$(NC)"
	@docker-compose -f deployments/docker-compose.dev.yml exec dev bash || echo "$(YELLOW)No development container found$(NC)"

# =============================================================================
# Watch Targets
# =============================================================================

.PHONY: test-watch
test-watch: ## Run tests in watch mode (requires entr)
	@echo "$(BLUE)Running tests in watch mode...$(NC)"
	@command -v entr >/dev/null 2>&1 || { echo "$(RED)entr not found. Install with: apt-get install entr$(NC)"; exit 1; }
	@find . -name "*.go" | entr -c make test-unit

.PHONY: lint-watch
lint-watch: ## Run linter in watch mode (requires entr)
	@echo "$(BLUE)Running linter in watch mode...$(NC)"
	@command -v entr >/dev/null 2>&1 || { echo "$(RED)entr not found. Install with: apt-get install entr$(NC)"; exit 1; }
	@find . -name "*.go" | entr -c make lint

# =============================================================================
# Cleanup Targets
# =============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BIN_DIR)
	@rm -rf $(COVERAGE_DIR)
	@go clean -cache
	@go clean -testcache
	@echo "$(GREEN)✓ Cleanup completed$(NC)"

.PHONY: clean-mocks
clean-mocks: ## Clean generated mocks
	@echo "$(BLUE)Cleaning generated mocks...$(NC)"
	@find . -name "*_mock.go" -type f -delete
	@echo "$(GREEN)✓ Mocks cleaned$(NC)"

.PHONY: clean-all
clean-all: clean clean-mocks clean-docker ## Clean everything (build artifacts, mocks, docker)
	@echo "$(GREEN)✓ Complete cleanup finished$(NC)"

# =============================================================================
# Release Targets
# =============================================================================

.PHONY: version
version: ## Show version information
	@echo "$(CYAN)Project: $(PROJECT_NAME)$(NC)"
	@echo "$(CYAN)Module: $(MODULE_NAME)$(NC)"
	@echo "$(CYAN)Go Version: $(GO_VERSION)$(NC)"
	@echo "$(CYAN)Git Commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')$(NC)"
	@echo "$(CYAN)Git Branch: $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')$(NC)"

.PHONY: tag
tag: ## Create a git tag (usage: make tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then echo "$(RED)VERSION is required. Usage: make tag VERSION=v1.0.0$(NC)"; exit 1; fi
	@echo "$(BLUE)Creating tag $(VERSION)...$(NC)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "$(GREEN)✓ Tag $(VERSION) created and pushed$(NC)"

.PHONY: release
release: ## Build release artifacts
	@echo "$(BLUE)Building release artifacts...$(NC)"
	@goreleaser release --snapshot --rm-dist || echo "$(YELLOW)goreleaser not found, using manual build$(NC)"
	@make build
	@echo "$(GREEN)✓ Release artifacts built$(NC)"

# =============================================================================
# Validation Targets
# =============================================================================

.PHONY: validate
validate: ## Quick validation (build + unit tests)
	@echo "$(BLUE)Validating project...$(NC)"
	@go build ./...
	@go test -v ./$(PKG_DIR)/... -timeout=2m
	@echo "$(GREEN)✓ Validation complete$(NC)"

.PHONY: validate-full
validate-full: check-docker ## Full validation (all tests)
	@echo "$(BLUE)Running full validation...$(NC)"
	@make fmt vet lint test-unit test-integration
	@echo "$(GREEN)✓ Full validation complete$(NC)"

# =============================================================================
# Help Targets
# =============================================================================

.PHONY: help-testing
help-testing: ## Show testing-specific help
	@echo "$(CYAN)Testing Targets:$(NC)"
	@echo "  $(GREEN)test-unit$(NC)              - Unit tests (no external dependencies)"
	@echo "  $(GREEN)test-integration$(NC)       - Integration tests (requires Docker)"
	@echo "  $(GREEN)test-performance$(NC)       - Performance benchmarks"
	@echo "  $(GREEN)test-ci$(NC)                - CI/CD validation tests"
	@echo "  $(GREEN)test-basic$(NC)             - Basic integration tests only"
	@echo "  $(GREEN)test-database$(NC)          - Database integration tests only"
	@echo "  $(GREEN)test-containers$(NC)        - Container management tests"
	@echo "  $(GREEN)test-skeleton$(NC)          - Skeleton application tests"
	@echo ""
	@echo "$(CYAN)CI/CD Targets:$(NC)"
	@echo "  $(GREEN)ci$(NC)                     - Fast CI pipeline"
	@echo "  $(GREEN)ci-integration$(NC)         - Full CI with integration tests"
	@echo "  $(GREEN)validate$(NC)               - Quick validation"
	@echo "  $(GREEN)validate-full$(NC)          - Full validation"

.PHONY: help-development
help-development: ## Show development-specific help
	@echo "$(CYAN)Development Targets:$(NC)"
	@echo "  $(GREEN)setup-dev$(NC)              - Setup development environment"
	@echo "  $(GREEN)dev$(NC)                    - Full development cycle"
	@echo "  $(GREEN)test-watch$(NC)             - Run tests in watch mode"
	@echo "  $(GREEN)lint-watch$(NC)             - Run linter in watch mode"
	@echo "  $(GREEN)dev-up$(NC)                 - Start development environment"
	@echo "  $(GREEN)dev-down$(NC)               - Stop development environment"
	@echo "  $(GREEN)dev-shell$(NC)              - Enter development container"

# =============================================================================
# Default Target
# =============================================================================

.DEFAULT_GOAL := help 