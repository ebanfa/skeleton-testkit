# Skeleton-Testkit

A simple testing framework for applications built on the skeleton component system.

## Quick Start

```go
// Create a skeleton application container
app := testkit.NewSkeletonApp("my-app:latest")

// Start the container
err := app.Start(context.Background())
if err != nil {
    log.Fatal(err)
}
defer app.Stop(context.Background())

// Your tests here...
```

## Installation

```bash
go get github.com/fintechain/skeleton-testkit
```

## Basic Usage

### Testing with Database

```go
func TestWithDatabase(t *testing.T) {
    // testcontainers-go creates and manages containers
    postgres := testkit.NewPostgresContainer()
    app := testkit.NewSkeletonApp("my-app:latest").
        WithDatabase(postgres)
    
    err := app.Start(context.Background())
    require.NoError(t, err)
    defer app.Stop(context.Background())
    
    // Your database tests here...
    // Containers are automatically cleaned up
}
```

## Development

### Prerequisites

- Docker
- Docker Compose
- Make

### Development Environment

The project includes a containerized development environment that provides Go 1.21+ and Docker access for testcontainers:

```bash
# Start development environment
make dev-up

# Enter development container
make dev-shell

# Inside the container, you can:
go test ./...
go build ./...
make lint
```

### Local Development (Alternative)

If you prefer local development:

- Go 1.21+
- Docker (for testcontainers)

```bash
make build
make test
```

### Testing

```bash
# Run tests locally
make test

# Run tests in development container
make dev-test

# Run integration tests (uses testcontainers)
make test-integration
```

### Key Points

- **External dependencies** (PostgreSQL, Redis, etc.) are created and managed by testcontainers-go during tests
- **No pre-defined services** - each test gets isolated, clean containers
- **Development container** provides consistent Go environment with Docker access

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make test` and `make lint`
6. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details. 