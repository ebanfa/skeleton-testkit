# Skeleton-Testkit

> A comprehensive testing framework for applications built on the skeleton component system

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/fintechain/skeleton-testkit/actions)

## 🎯 Vision

The **skeleton-testkit** is a testing framework designed specifically for applications built on the [skeleton component system](https://github.com/ebanfa/skeleton). It provides reusable testing infrastructure, container management, and verification patterns that enable developers to build robust, production-ready applications with confidence.

### What Makes This Different

- **Skeleton-Native**: Built specifically for skeleton-based applications
- **Container-First**: Uses Docker containers for realistic testing environments
- **Verification-Driven**: Comprehensive verification of skeleton components, plugins, and services
- **Developer-Friendly**: Simple, intuitive API that reduces testing complexity
- **Production-Ready**: Robust error handling, monitoring, and CI/CD integration

## 🚀 Quick Start

```go
func TestMySkeletonApp(t *testing.T) {
    // Create infrastructure containers
    postgres := testkit.NewPostgresContainer()
    redis := testkit.NewRedisContainer()
    
    // Create skeleton application with dependencies
    app := testkit.NewSkeletonApp("my-app:latest").
        WithDatabase(postgres).
        WithCache(redis).
        WithSkeletonConfig(&testkit.SkeletonConfig{
            ServiceID: "my-app",
            Plugins: []testkit.SkeletonPluginConfig{
                {Name: "auth-plugin", Version: "1.0.0"},
            },
        })
    
    // Verify skeleton application behavior
    verifier := testkit.NewSystemVerifier(app)
    err := verifier.VerifySkeletonStartup(context.Background())
    require.NoError(t, err)
    
    // Test skeleton components and plugins
    componentVerifier := testkit.NewComponentVerifier(app)
    err = componentVerifier.VerifySkeletonComponentRegistered(ctx, "auth-service")
    require.NoError(t, err)
}
```

## 📋 Features

### 🐳 Container Management
- **Application Containers**: Skeleton-based application containers with fluent configuration
- **Infrastructure Containers**: Postgres, Redis, Kafka, and other dependencies
- **Multi-Container Orchestration**: Complex scenarios with dependency management
- **Network Management**: Container networking and service discovery

### ✅ Verification Framework
- **System Verification**: Startup, shutdown, and health verification for skeleton applications
- **Component Verification**: Verify skeleton component registration and lifecycle
- **Plugin Verification**: Test skeleton plugin loading, unloading, and functionality
- **Service Verification**: Verify skeleton service management and health
- **Operation Verification**: Test skeleton operation execution and results

### 🏥 Health Monitoring
- **Health Checks**: Built-in health checks for skeleton applications
- **Monitoring**: Continuous health monitoring during tests
- **Status Reporting**: Comprehensive health status aggregation
- **Custom Checks**: Support for custom health verification strategies

### ⚙️ Configuration & Integration
- **YAML Configuration**: Declarative configuration for complex test scenarios
- **CI/CD Integration**: Ready-to-use templates for GitHub Actions, GitLab CI
- **Performance Optimization**: Container pooling, parallel execution
- **Monitoring**: Prometheus metrics, Grafana dashboards

## 📖 Documentation

### Getting Started
- **[Quick Start Guide](docs/architecture/QUICK_START_GUIDE.md)** - Immediate actionable steps to begin implementation
- **[Implementation Plan](docs/architecture/IMPLEMENTATION_PLAN.md)** - Comprehensive 12-week development roadmap
- **[Architecture Specification](docs/architecture/skeleton-testkit-specification.md)** - Complete API specification and design

### Examples
- **[Basic Usage](examples/basic/)** - Simple skeleton application testing
- **[Database Integration](examples/database/)** - Testing with database dependencies
- **[Plugin Development](examples/plugins/)** - Plugin lifecycle testing
- **[Advanced Scenarios](examples/advanced/)** - Multi-container integration testing

### API Reference
- **[Container API](docs/api/container.md)** - Container management and orchestration
- **[Verification API](docs/api/verification.md)** - Verification strategies and patterns
- **[Health API](docs/api/health.md)** - Health monitoring and checks

## 🏗️ Implementation Roadmap

### Phase 1: Foundation (Weeks 1-3)
**Goal**: Basic container management for skeleton applications

- [x] Project structure and build system
- [x] Core container interfaces and implementations
- [x] Basic skeleton application containers
- [x] Infrastructure containers (Postgres, Redis)
- [x] Integration testing framework

### Phase 2: Verification Framework (Weeks 4-6)
**Goal**: Comprehensive verification of skeleton applications

- [ ] System-level verification strategies
- [ ] Health monitoring and status reporting
- [ ] Skeleton-specific verification patterns
- [ ] Error handling and debugging support

### Phase 3: Advanced Features (Weeks 7-9)
**Goal**: Plugin verification, orchestration, examples

- [ ] Plugin and service verification
- [ ] Multi-container orchestration
- [ ] Configuration management system
- [ ] Comprehensive examples and documentation

### Phase 4: Production Ready (Weeks 10-12)
**Goal**: Performance, CI/CD, monitoring

- [ ] Performance optimization and container pooling
- [ ] CI/CD integration templates
- [ ] Monitoring and observability
- [ ] CLI tools and production readiness

## 🛠️ Development

### Prerequisites
- Go 1.21+
- Docker and Docker Compose
- Make (optional, for build automation)

### Setup
```bash
# Clone the repository
git clone https://github.com/fintechain/skeleton-testkit.git
cd skeleton-testkit

# Install dependencies
go mod download

# Run tests
make test

# Run integration tests (requires Docker)
make test-integration
```

### Project Structure
```
skeleton-testkit/
├── pkg/                    # Public API
│   ├── testkit/           # Main entry point
│   ├── container/         # Container management
│   ├── verification/      # Verification framework
│   └── health/           # Health monitoring
├── internal/              # Private implementation
│   ├── domain/           # Core concepts
│   ├── application/      # Use cases
│   ├── infrastructure/   # Technical implementations
│   └── interfaces/       # External interfaces
├── examples/             # Usage examples
├── test/                 # Testing infrastructure
└── docs/                # Documentation
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow
1. **Design**: Create interface and API design
2. **Test**: Write integration tests first
3. **Implement**: Build minimal viable implementation
4. **Example**: Create working example
5. **Document**: Update API documentation
6. **Review**: Code review and testing

### Code Quality Standards
- Comprehensive error handling with context
- Structured logging throughout
- 90%+ test coverage
- Clean, testable interfaces
- Extensive documentation and examples

## 📊 Success Metrics

### Developer Experience
- **50% reduction** in test setup time compared to manual container management
- **Sub-30 second** container startup times
- **99%+ test success rate** in CI/CD environments

### Ecosystem Adoption
- **80% adoption rate** among skeleton-based projects
- **Active community** contributions and feedback
- **Production-ready** CI/CD integration

## 🔗 Related Projects

- **[Skeleton Framework](https://github.com/ebanfa/skeleton)** - The component system that this testkit supports
- **[Testcontainers-Go](https://github.com/testcontainers/testcontainers-go)** - Underlying container management library
- **[Component-Testkit Concept](docs/architecture/component-testkit.md)** - Broader vision for component-based testing

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Skeleton Framework Team** - For creating the component system this testkit supports
- **Testcontainers Community** - For the excellent container testing foundation
- **Go Community** - For the robust ecosystem and testing tools

---

**Ready to get started?** Check out the [Quick Start Guide](docs/architecture/QUICK_START_GUIDE.md) and begin building robust tests for your skeleton-based applications today!

## 📞 Support

- **GitHub Issues**: [Report bugs and request features](https://github.com/fintechain/skeleton-testkit/issues)
- **Documentation**: [Complete guides and API reference](docs/)
- **Community**: [Join the discussion](https://community.fintechain.com/skeleton-testkit)
- **Examples**: [Real-world usage patterns](examples/)

**The skeleton-testkit is a testing utility that enhances skeleton-based application development by providing containerized testing environments, comprehensive verification strategies, and reusable testing patterns. It works WITH the skeleton framework, not as a replacement for it.** 