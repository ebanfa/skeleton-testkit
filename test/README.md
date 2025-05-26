# Skeleton-Testkit Testing

This directory contains the Phase 1 Testing Strategy implementation for the skeleton-testkit framework. The tests are organized to validate the core functionality, performance, and reliability of the testkit in various environments.

## Test Structure

```
test/
├── integration/           # Integration tests
│   ├── basic_test.go     # Basic container lifecycle tests
│   ├── database_test.go  # Database integration tests
│   ├── performance_test.go # Performance benchmarks
│   ├── ci_test.go        # CI/CD validation tests
│   └── integration_test.go # Test setup and configuration
└── README.md             # This file
```

## Test Categories

### 1. Basic Integration Tests (`basic_test.go`)

Tests fundamental container operations and lifecycle management:

- **TestBasicSkeletonApp**: Basic container creation and startup
- **TestSkeletonAppLifecycle**: Complete container lifecycle validation
- **TestSkeletonAppConfiguration**: Configuration options and customization
- **TestSkeletonAppErrorHandling**: Error scenarios and recovery

**Requirements**: Docker daemon accessible

### 2. Database Integration Tests (`database_test.go`)

Tests integration between skeleton applications and database containers:

- **TestSkeletonAppWithDatabase**: Basic app-database integration
- **TestSkeletonAppWithCustomDatabaseConfig**: Custom database configurations
- **TestSkeletonAppMultipleDatabases**: Multiple database dependencies
- **TestDatabaseContainerLifecycle**: Database container lifecycle management

**Requirements**: Docker daemon accessible, sufficient memory for multiple containers

### 3. Performance Tests (`performance_test.go`)

Benchmarks and performance validation:

- **BenchmarkSkeletonAppStartup**: Container startup performance (< 30s requirement)
- **BenchmarkDatabaseContainerStartup**: Database startup performance
- **BenchmarkSkeletonAppWithDatabaseStartup**: Combined startup performance
- **TestContainerStartupPerformance**: Startup time validation
- **TestContainerShutdownPerformance**: Shutdown time validation
- **TestConcurrentContainerOperations**: Concurrent operation performance
- **TestMemoryUsageUnderLoad**: Memory usage patterns

**Requirements**: Docker daemon, sufficient system resources for benchmarking

### 4. CI/CD Validation Tests (`ci_test.go`)

Tests for CI/CD environment compatibility:

- **TestCIEnvironmentCompatibility**: Basic CI environment operation
- **TestCIErrorReporting**: Error message clarity for debugging
- **TestCIParallelExecution**: Parallel test execution support
- **TestCITimeouts**: Timeout handling in CI environments
- **TestCIDockerEnvironment**: Docker environment compatibility

**Requirements**: Docker daemon, CI environment variables (optional)

### 5. Test Configuration (`integration_test.go`)

Common test setup and utilities:

- **TestIntegrationTestSetup**: Test environment validation
- **TestIntegrationTestConfiguration**: Test configuration options
- Helper functions for environment detection and test utilities

## Running Tests

### Prerequisites

1. **Docker**: All integration tests require Docker to be installed and running
2. **Go**: Go 1.19 or later
3. **Dependencies**: Run `go mod download` to install required packages

### Quick Start

```bash
# Validate Docker is available
make check-docker

# Run all tests
make test

# Run only unit tests (fast, no Docker required)
make test-unit

# Run only integration tests
make test-integration
```

### Specific Test Categories

```bash
# Basic integration tests
make test-basic

# Database integration tests
make test-database

# Performance benchmarks
make test-performance

# CI/CD validation tests
make test-ci

# Test environment setup validation
make test-setup
```

### Development and Debugging

```bash
# Run tests with verbose output
make test-integration-verbose

# Run tests in short mode (faster)
make test-integration-short

# Run tests with coverage
make test-coverage

# Run specific benchmarks
make benchmark-startup
make benchmark-database
```

### CI/CD Integration

```bash
# For CI environments (longer timeouts)
make ci-test

# For fast CI validation
make ci-test-fast

# Quick validation (build + unit tests)
make validate
```

## Test Configuration

### Build Tags

Integration tests use the `integration` build tag to separate them from unit tests:

```go
//go:build integration
// +build integration
```

### Environment Variables

Tests can be configured using environment variables:

- `CI`: Set to any value to indicate CI environment
- `GITHUB_ACTIONS`: GitHub Actions CI indicator
- `GITLAB_CI`: GitLab CI indicator
- `SKELETON_TEST_IMAGE`: Override default test image
- Other CI environment variables (see `ci_test.go`)

### Timeouts

Default timeouts are environment-aware:

- **Local development**: 60 seconds for most operations
- **CI environments**: 120 seconds for startup operations
- **Performance tests**: Up to 15 minutes for comprehensive benchmarks

## Success Metrics

The Phase 1 testing strategy validates these success metrics:

### Functional Requirements

- ✅ Container lifecycle management (start, stop, status)
- ✅ Database integration and connectivity
- ✅ Configuration and customization options
- ✅ Error handling and recovery
- ✅ Resource cleanup and isolation

### Performance Requirements

- ✅ Container startup time < 30 seconds
- ✅ Database startup time < 30 seconds
- ✅ Combined startup time < 60 seconds
- ✅ Shutdown time < 10 seconds
- ✅ Memory usage under load monitoring

### Reliability Requirements

- ✅ CI/CD environment compatibility
- ✅ Parallel test execution support
- ✅ Timeout handling and graceful degradation
- ✅ Clear error messages for debugging
- ✅ Resource constraint handling

## Troubleshooting

### Common Issues

1. **Docker not available**
   ```
   Error: Docker is not available. Integration tests require Docker.
   Solution: Install Docker and ensure the daemon is running
   ```

2. **Permission denied accessing Docker**
   ```
   Error: permission denied while trying to connect to Docker daemon
   Solution: Add user to docker group or run with sudo
   ```

3. **Container startup timeout**
   ```
   Error: Container startup timeout
   Solution: Check system resources, increase timeout, or use shorter test mode
   ```

4. **Port conflicts**
   ```
   Error: Port already in use
   Solution: Ensure no other containers are using the same ports
   ```

### Debug Mode

Enable verbose logging for debugging:

```bash
# Run with verbose output
go test -v -tags=integration ./test/integration/...

# Run with race detection
go test -v -race -tags=integration ./test/integration/...

# Run specific test with timeout
go test -v -tags=integration -run TestBasicSkeletonApp ./test/integration/basic_test.go -timeout=5m
```

### Resource Monitoring

Monitor system resources during tests:

```bash
# Monitor Docker containers
docker stats

# Monitor system resources
htop

# Check Docker disk usage
docker system df
```

## Contributing

When adding new tests:

1. Use the `integration` build tag for tests requiring Docker
2. Include appropriate timeouts for CI environments
3. Add cleanup code using `defer` statements
4. Use descriptive test names and error messages
5. Update this README with new test descriptions

### Test Naming Convention

- `Test*`: Standard test functions
- `Benchmark*`: Performance benchmark functions
- `Example*`: Example functions (if any)

### Error Handling

- Use `require` for critical assertions that should stop the test
- Use `assert` for non-critical assertions that allow test continuation
- Provide descriptive error messages for debugging
- Include context about the test environment in error messages

## References

- [Skeleton-Testkit Specification](../docs/architecture/skeleton-testkit-specification.md)
- [Implementation Plan](../docs/implementation-plan.md)
- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify) 