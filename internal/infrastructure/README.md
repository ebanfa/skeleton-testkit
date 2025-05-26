# Infrastructure Layer Implementation

This directory contains the infrastructure layer implementation for the skeleton-testkit as defined in section 3.2.3 of the implementation plan and the skeleton-testkit specification.

## Overview

The infrastructure layer provides the technical implementation details for container management, networking, and lifecycle operations. It serves as the foundation for the skeleton-testkit's container-based testing capabilities.

## Specification Compliance

This implementation **strictly adheres** to:
- **Section 3.2.3** of the implementation plan
- **Container Interface** as defined in the skeleton-testkit specification (lines 135-160)
- **Configuration structures** as specified in the skeleton-testkit specification
- **Design principles** including Constructor Injection and Interface Substitution

## Components

### Docker Package (`docker/`)

#### `container.go`
- **DockerContainer**: Wraps `testcontainers.Container` with additional configuration and error handling
- **ContainerConfig**: Holds basic container configuration (ID, name, image, environment, ports)
- **Fully implements** the `container.Container` interface from the domain layer with all required methods:
  - Identity: `ID()`, `Name()`, `Image()`
  - Lifecycle: `Start()`, `Stop()`, `IsRunning()`
  - Network: `Host()`, `Port()`, `ConnectionString()`
  - Health: `WaitForReady()`, `HealthCheck()`
  - Debugging: `Logs()`, `Exec()`
- Provides lifecycle management with proper error handling
- Handles port mapping and host access
- Includes logging and command execution capabilities

#### `network.go`
- **PortManager**: Manages port allocations and mappings for containers
- **ContainerLifecycleManager**: Orchestrates multiple container lifecycles
- **GetPortMapping**: Utility function for port mapping strings
- Provides centralized container and port management

### Testcontainers Package (`testcontainers/`)

#### `app_container.go`
- **TestcontainerAppContainer**: Extends DockerContainer with skeleton-specific functionality
- Manages container dependencies (database, cache, etc.)
- Handles skeleton-specific environment variables and configuration
- Provides connection strings and endpoint URLs for health checks, metrics, and shutdown
- Implements dependency-aware startup and shutdown sequences
- **Specification-compliant** configuration using `SkeletonConfig` and `SkeletonPluginConfig`

#### `postgres_container.go`
- **PostgresContainer**: Specialized container for PostgreSQL databases
- **PostgresConfig**: Configuration for PostgreSQL instances
- Provides **proper PostgreSQL connection string** generation (`postgres://user:pass@host:port/db`)
- Includes proper wait strategies for database readiness

#### `redis_container.go`
- **RedisContainer**: Specialized container for Redis cache instances
- **RedisConfig**: Configuration for Redis instances
- Supports password-protected Redis instances
- Provides **proper Redis connection string** generation (`redis://[password@]host:port`)

## Key Features

### 1. Docker Container Wrapper Implementation ✅
- Wraps testcontainers-go with domain-specific interfaces
- Provides consistent error handling with `ContainerError`
- **Fully implements** the `container.Container` interface as specified
- Supports configuration injection via constructor pattern

### 2. Testcontainers-go Integration ✅
- Full integration with testcontainers-go library
- Specialized containers for common infrastructure (Postgres, Redis)
- Skeleton-specific application container with dependency management
- Proper wait strategies for container readiness

### 3. Basic Networking and Port Management ✅
- Port allocation and deallocation tracking
- Port mapping utilities
- Container lifecycle management
- Support for random port assignment

### 4. Container Lifecycle Management ✅
- Dependency-aware startup sequences
- Graceful shutdown with proper cleanup
- Container registration and discovery
- Bulk operations (start all, stop all)

## Specification Alignment

### Container Interface Compliance
The implementation provides **100% compliance** with the Container interface specification:

```go
// All methods implemented as per specification
type Container interface {
    // Identity
    ID() string                    ✅ Implemented
    Name() string                  ✅ Implemented  
    Image() string                 ✅ Implemented

    // Lifecycle
    Start(ctx context.Context) error    ✅ Implemented
    Stop(ctx context.Context) error     ✅ Implemented
    IsRunning() bool                    ✅ Implemented

    // Network
    Host() string                       ✅ Implemented
    Port(internal int) (int, error)     ✅ Implemented
    ConnectionString() string           ✅ Implemented

    // Health
    WaitForReady(ctx context.Context, timeout time.Duration) error  ✅ Implemented
    HealthCheck(ctx context.Context) error                          ✅ Implemented

    // Logs and Debugging
    Logs(ctx context.Context) (io.Reader, error)  ✅ Implemented
    Exec(ctx context.Context, cmd []string) error ✅ Implemented
}
```

### Configuration Compliance
All configuration structures match the specification exactly:

- **AppConfig**: Includes all specified fields (ImageName, HealthEndpoint, MetricsEndpoint, ShutdownEndpoint, Environment, Ports, Volumes)
- **SkeletonConfig**: Uses correct naming (`SkeletonPluginConfig`, `SkeletonStorageConfig`)
- **PortMapping**: Supports internal/external port mapping with random port assignment
- **VolumeMapping**: Supports source/target volume mapping

## Usage Examples

### Basic Container Usage
```go
config := &docker.ContainerConfig{
    ID:    "test-container",
    Name:  "test",
    Image: "nginx:latest",
    Environment: map[string]string{
        "ENV": "test",
    },
    Ports: []container.PortMapping{
        {Internal: 80, External: 0},
    },
}

dockerContainer := docker.NewDockerContainer(config)
err := dockerContainer.Start(ctx)

// Access all interface methods
fmt.Printf("Container ID: %s\n", dockerContainer.ID())
fmt.Printf("Container Name: %s\n", dockerContainer.Name())
fmt.Printf("Container Image: %s\n", dockerContainer.Image())
fmt.Printf("Connection String: %s\n", dockerContainer.ConnectionString())
```

### Skeleton Application Container
```go
skeletonConfig := &container.SkeletonConfig{
    ServiceID: "test-app",
    Storage: container.SkeletonStorageConfig{
        Type: "postgres",
        URL:  "postgres://...",
    },
    Plugins: []container.SkeletonPluginConfig{
        {Name: "auth-plugin", Version: "1.0.0"},
    },
}

appContainer := testcontainers.NewTestcontainerAppContainer(config, skeletonConfig)
err := appContainer.Start(ctx)

// Access skeleton-specific endpoints
healthURL := appContainer.HealthEndpoint()      // "/health"
metricsURL := appContainer.MetricsEndpoint()    // "/metrics"
shutdownURL := appContainer.ShutdownEndpoint()  // "/shutdown"
```

### Infrastructure Containers
```go
// PostgreSQL with proper connection string
postgres := testcontainers.NewPostgresContainer()
err := postgres.Start(ctx)
connectionString := postgres.ConnectionString()  // "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"

// Redis with proper connection string
redis := testcontainers.NewRedisContainer()
err := redis.Start(ctx)
connectionString := redis.ConnectionString()     // "redis://localhost:6379" or "redis://:password@localhost:6379"
```

## Design Principles

1. **Constructor Injection**: All external dependencies are injected via constructors
2. **Interface Compliance**: Implements domain interfaces consistently and completely
3. **Error Handling**: Uses domain-specific error types with context
4. **Configuration-Driven**: Supports flexible configuration via structs
5. **Dependency Management**: Handles container dependencies automatically
6. **Specification Adherence**: Strictly follows the skeleton-testkit specification

## Dependencies

- `github.com/testcontainers/testcontainers-go`: Core testcontainers functionality
- `github.com/docker/go-connections/nat`: Docker port handling
- Internal domain packages for interfaces and error types

## Testing

The infrastructure layer is designed to be testable with:
- Mock implementations of external dependencies
- Constructor injection for all external systems
- Clear separation of concerns between configuration and implementation

## Specification Compliance Verification

✅ **Section 3.2.3 Deliverables**:
- [x] Docker container wrapper implementation
- [x] Testcontainers-go integration  
- [x] Basic networking and port management
- [x] Container lifecycle management

✅ **Container Interface Specification**:
- [x] All 12 interface methods implemented
- [x] Proper error handling with ContainerError
- [x] Specification-compliant method signatures
- [x] Correct return types and behaviors

✅ **Configuration Specification**:
- [x] AppConfig with all required fields
- [x] SkeletonConfig with correct naming
- [x] PortMapping and VolumeMapping support
- [x] Skeleton-specific plugin and storage configuration

This implementation provides a solid, specification-compliant foundation for the skeleton-testkit's container-based testing capabilities. 