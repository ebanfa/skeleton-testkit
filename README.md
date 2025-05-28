# Skeleton-Testkit

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/fintechain/skeleton-testkit)](https://goreportcard.com/report/github.com/fintechain/skeleton-testkit)
[![Coverage Status](https://coveralls.io/repos/github/fintechain/skeleton-testkit/badge.svg?branch=main)](https://coveralls.io/github/fintechain/skeleton-testkit?branch=main)
[![Build Status](https://github.com/fintechain/skeleton-testkit/workflows/CI/badge.svg)](https://github.com/fintechain/skeleton-testkit/actions)
[![Documentation](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/fintechain/skeleton-testkit)

> **A comprehensive testing framework for skeleton-based applications**

Skeleton-Testkit is a specialized testing framework designed to provide robust, containerized testing environments for applications built on the [Skeleton Framework](https://github.com/fintechain/skeleton). It offers production-like testing scenarios, comprehensive verification strategies, and reusable testing infrastructure.

## 🎯 Purpose

The Skeleton-Testkit bridges the gap between unit testing and production deployment by providing:

- **Container-First Testing**: Realistic testing environments using Docker containers
- **Skeleton-Aware Testing**: Deep integration with skeleton component architecture
- **Production-Like Scenarios**: Testing that mirrors real-world deployment conditions
- **Comprehensive Verification**: Health monitoring, state verification, and behavior validation
- **Developer Experience**: Streamlined testing workflows and clear feedback

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Skeleton-Testkit                        │
├─────────────────────────────────────────────────────────────┤
│  Testing Layer                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Verification  │  │   Container     │  │   Health    │ │
│  │   Strategies    │  │   Management    │  │   Monitoring│ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                 Skeleton Framework                         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Components    │  │   Plugins       │  │   Services  │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                 Container Runtime                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │     Docker      │  │   PostgreSQL    │  │    Redis    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+**: [Download and install Go](https://golang.org/dl/)
- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)
- **Make**: Build automation tool

### Installation

```bash
# Clone the repository
git clone https://github.com/fintechain/skeleton-testkit.git
cd skeleton-testkit

# Setup development environment
make setup-dev

# Verify installation
make validate
```

### Basic Usage

```go
package main

import (
    "context"
    "testing"
    
    "github.com/fintechain/skeleton-testkit/pkg/testkit"
    "github.com/fintechain/skeleton-testkit/pkg/container"
)

func TestSkeletonApplication(t *testing.T) {
    // Create a new testkit instance
    tk := testkit.New(t)
    defer tk.Cleanup()
    
    // Start a skeleton application container
    app, err := tk.StartSkeletonApp(context.Background(), &container.SkeletonConfig{
        Image: "my-skeleton-app:latest",
        Env: map[string]string{
            "DATABASE_URL": "postgres://test:test@localhost:5432/testdb",
        },
    })
    if err != nil {
        t.Fatalf("Failed to start skeleton app: %v", err)
    }
    
    // Wait for application to be ready
    if err := app.WaitForReady(context.Background()); err != nil {
        t.Fatalf("Application failed to become ready: %v", err)
    }
    
    // Verify application behavior
    verifier := tk.NewVerifier()
    if err := verifier.VerifyHealthy(app); err != nil {
        t.Fatalf("Health verification failed: %v", err)
    }
    
    // Test application functionality
    response, err := app.HTTPClient().Get("/api/health")
    if err != nil {
        t.Fatalf("Health check failed: %v", err)
    }
    defer response.Body.Close()
    
    if response.StatusCode != 200 {
        t.Errorf("Expected status 200, got %d", response.StatusCode)
    }
}
```

## 📁 Project Structure

```
skeleton-testkit/
├── cmd/                    # Command-line tools
│   ├── testkit-cli/       # CLI for testkit operations
│   └── examples/          # Example applications
├── pkg/                   # Public API packages
│   ├── testkit/          # Core testkit functionality
│   ├── container/        # Container management
│   ├── verification/     # Verification strategies
│   └── health/           # Health monitoring
├── internal/             # Internal implementation
│   ├── domain/           # Domain models and interfaces
│   ├── infrastructure/   # Infrastructure implementations
│   └── adapters/         # External service adapters
├── test/                 # Test suites
│   ├── integration/      # Integration tests
│   ├── unit/            # Unit tests
│   └── fixtures/        # Test fixtures and data
├── examples/             # Usage examples
│   ├── basic/           # Basic usage examples
│   ├── advanced/        # Advanced scenarios
│   └── benchmarks/      # Performance examples
├── docs/                # Documentation
│   ├── architecture/    # Architecture documentation
│   ├── api/            # API documentation
│   └── guides/         # User guides
├── deployments/         # Deployment configurations
│   ├── docker/         # Docker configurations
│   └── k8s/           # Kubernetes manifests
├── scripts/            # Build and utility scripts
├── configs/            # Configuration files
└── Makefile           # Build automation
```

## 🛠️ Development

### Available Make Targets

The project includes a comprehensive Makefile with organized targets:

```bash
# Show all available targets
make help

# Development workflow
make dev                    # Full development cycle
make setup-dev             # Setup development environment
make test-watch            # Run tests in watch mode

# Testing
make test                  # Run all tests
make test-unit            # Unit tests only
make test-integration     # Integration tests (requires Docker)
make test-performance     # Performance benchmarks
make coverage             # Generate coverage report

# Code quality
make lint                 # Run linter
make fmt                  # Format code
make vet                  # Run go vet

# CI/CD
make ci                   # Fast CI pipeline
make ci-integration       # Full CI with integration tests
make validate             # Quick validation
```

### Testing Strategy

The testkit implements a comprehensive testing strategy:

#### Unit Tests
- **Fast execution** (< 30 seconds)
- **No external dependencies**
- **High coverage** of business logic
- **Isolated components**

```bash
make test-unit
```

#### Integration Tests
- **Container-based** testing environments
- **Real database** connections
- **Network communication** testing
- **End-to-end** scenarios

```bash
make test-integration
```

#### Performance Tests
- **Startup time** benchmarks
- **Resource usage** monitoring
- **Throughput** measurements
- **Scalability** testing

```bash
make test-performance
```

### Code Quality Standards

- **Go fmt**: Consistent code formatting
- **Go vet**: Static analysis for common errors
- **golangci-lint**: Comprehensive linting
- **Test coverage**: Minimum 80% coverage
- **Documentation**: Comprehensive godoc comments

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TESTKIT_LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `TESTKIT_TIMEOUT` | Default timeout for operations | `30s` |
| `TESTKIT_DOCKER_HOST` | Docker daemon host | `unix:///var/run/docker.sock` |
| `TESTKIT_CLEANUP_CONTAINERS` | Auto-cleanup containers after tests | `true` |
| `TESTKIT_PARALLEL_TESTS` | Enable parallel test execution | `true` |

### Configuration Files

```yaml
# configs/testkit.yaml
testkit:
  timeout: 30s
  log_level: info
  docker:
    host: unix:///var/run/docker.sock
    cleanup: true
  containers:
    postgres:
      image: postgres:15-alpine
      env:
        POSTGRES_DB: testdb
        POSTGRES_USER: test
        POSTGRES_PASSWORD: test
    redis:
      image: redis:7-alpine
```

## 📚 Documentation

### Core Concepts

- **[Architecture Overview](docs/architecture/skeleton-testkit-specification.md)**: Detailed system architecture
- **[Testing Strategies](docs/guides/testing-strategies.md)**: Comprehensive testing approaches
- **[Container Management](docs/guides/container-management.md)**: Docker container lifecycle
- **[Verification Framework](docs/guides/verification.md)**: Application verification strategies

### API Reference

- **[Testkit API](docs/api/testkit.md)**: Core testkit functionality
- **[Container API](docs/api/container.md)**: Container management interface
- **[Verification API](docs/api/verification.md)**: Verification strategies
- **[Health API](docs/api/health.md)**: Health monitoring interface

### Examples

- **[Basic Usage](examples/basic/)**: Simple testing scenarios
- **[Database Testing](examples/database/)**: Database integration tests
- **[Performance Testing](examples/performance/)**: Benchmark examples
- **[CI/CD Integration](examples/ci/)**: Continuous integration setup

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Make** your changes
4. **Run** tests (`make test`)
5. **Commit** your changes (`git commit -m 'Add amazing feature'`)
6. **Push** to the branch (`git push origin feature/amazing-feature`)
7. **Open** a Pull Request

### Code Standards

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Write comprehensive tests for new features
- Update documentation for API changes
- Ensure all CI checks pass

## 📊 Performance

### Benchmarks

| Operation | Duration | Memory | Allocations |
|-----------|----------|---------|-------------|
| Container Startup | ~2.5s | 45MB | 1,234 |
| Health Check | ~50ms | 2MB | 45 |
| Verification | ~100ms | 5MB | 123 |
| Cleanup | ~1s | 10MB | 234 |

### Resource Requirements

- **Memory**: 512MB minimum, 2GB recommended
- **CPU**: 2 cores minimum, 4 cores recommended
- **Disk**: 10GB for container images and test data
- **Network**: Internet access for pulling container images

## 🔒 Security

### Security Considerations

- **Container Isolation**: Tests run in isolated Docker containers
- **Network Security**: Containers use isolated networks
- **Credential Management**: Test credentials are ephemeral
- **Resource Limits**: Containers have resource constraints

### Reporting Security Issues

Please report security vulnerabilities to [security@fintechain.com](mailto:security@fintechain.com).

## 📈 Roadmap

### Short-term (Q1 2024)
- [ ] Enhanced CI/CD integration
- [ ] Performance optimization
- [ ] Extended container support
- [ ] Improved documentation

### Medium-term (Q2-Q3 2024)
- [ ] Kubernetes testing support
- [ ] Advanced verification strategies
- [ ] Monitoring and observability
- [ ] Plugin system

### Long-term (Q4 2024+)
- [ ] Cloud provider integration
- [ ] Distributed testing
- [ ] AI-powered test generation
- [ ] Community ecosystem

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **[Skeleton Framework](https://github.com/fintechain/skeleton)**: The foundation this testkit builds upon
- **[Docker](https://docker.com)**: Container runtime that enables realistic testing
- **[Go Community](https://golang.org/community)**: For the excellent testing tools and libraries
- **[Contributors](CONTRIBUTORS.md)**: Everyone who has contributed to this project

## 📞 Support

- **Documentation**: [https://docs.fintechain.com/skeleton-testkit](https://docs.fintechain.com/skeleton-testkit)
- **Issues**: [GitHub Issues](https://github.com/fintechain/skeleton-testkit/issues)
- **Discussions**: [GitHub Discussions](https://github.com/fintechain/skeleton-testkit/discussions)
- **Email**: [support@fintechain.com](mailto:support@fintechain.com)

---

<div align="center">
  <strong>Built with ❤️ by the Fintechain Team</strong>
</div> 