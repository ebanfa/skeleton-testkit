package container

import (
	"context"
	"io"
	"time"

	domaincontainer "github.com/fintechain/skeleton-testkit/internal/domain/container"
	"github.com/fintechain/skeleton-testkit/internal/infrastructure/testcontainers"
)

// PostgresContainer represents a PostgreSQL database container for testing.
// It provides a clean interface for managing PostgreSQL containers used
// as dependencies in skeleton-based application testing.
type PostgresContainer struct {
	impl *testcontainers.PostgresContainer
}

// NewPostgresContainer creates a new PostgresContainer wrapper around the internal implementation.
// This follows the constructor injection pattern by accepting the implementation.
//
// Parameters:
//   - impl: The internal PostgresContainer implementation
//
// Returns:
//   - *PostgresContainer: A new PostgresContainer wrapper
func NewPostgresContainer(impl *testcontainers.PostgresContainer) *PostgresContainer {
	return &PostgresContainer{
		impl: impl,
	}
}

// Start starts the PostgreSQL container.
// This will pull the PostgreSQL image if needed and start the container.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns:
//   - error: Any error that occurred during startup
//
// Example:
//
//	err := postgres.Start(context.Background())
//	if err != nil {
//	    log.Fatalf("Failed to start PostgreSQL: %v", err)
//	}
func (p *PostgresContainer) Start(ctx context.Context) error {
	return p.impl.Start(ctx)
}

// Stop stops the PostgreSQL container and cleans up resources.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns:
//   - error: Any error that occurred during shutdown
func (p *PostgresContainer) Stop(ctx context.Context) error {
	return p.impl.Stop(ctx)
}

// IsRunning returns whether the PostgreSQL container is currently running.
//
// Returns:
//   - bool: True if the container is running, false otherwise
func (p *PostgresContainer) IsRunning() bool {
	return p.impl.IsRunning()
}

// ID returns the unique identifier of the PostgreSQL container.
//
// Returns:
//   - string: The container ID
func (p *PostgresContainer) ID() string {
	return p.impl.ID()
}

// Name returns the human-readable name of the PostgreSQL container.
//
// Returns:
//   - string: The container name
func (p *PostgresContainer) Name() string {
	return p.impl.Name()
}

// Image returns the Docker image name used by the PostgreSQL container.
//
// Returns:
//   - string: The Docker image name
func (p *PostgresContainer) Image() string {
	return p.impl.Image()
}

// Host returns the host address where the PostgreSQL container is accessible.
//
// Returns:
//   - string: The host address
func (p *PostgresContainer) Host() string {
	return p.impl.Host()
}

// Port returns the external port mapping for the specified internal port.
//
// Parameters:
//   - internal: The internal port number
//
// Returns:
//   - int: The external port number
//   - error: Any error that occurred
func (p *PostgresContainer) Port(internal int) (int, error) {
	return p.impl.Port(internal)
}

// ConnectionString returns the PostgreSQL connection string.
// This can be used to connect to the PostgreSQL database from the application.
//
// Returns:
//   - string: The PostgreSQL connection string
//
// Example:
//
//	connStr := postgres.ConnectionString()
//	// connStr = "postgres://testuser:testpass@localhost:5432/testdb"
func (p *PostgresContainer) ConnectionString() string {
	return p.impl.ConnectionString()
}

// WaitForReady waits for the PostgreSQL container to be ready to accept connections.
//
// Parameters:
//   - ctx: Context for the operation
//   - timeout: Maximum time to wait for readiness
//
// Returns:
//   - error: Any error that occurred while waiting
func (p *PostgresContainer) WaitForReady(ctx context.Context, timeout time.Duration) error {
	return p.impl.WaitForReady(ctx, timeout)
}

// HealthCheck performs a health check on the PostgreSQL container.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns:
//   - error: Any error that occurred during the health check
func (p *PostgresContainer) HealthCheck(ctx context.Context) error {
	return p.impl.HealthCheck(ctx)
}

// Database returns the name of the PostgreSQL database.
//
// Returns:
//   - string: The database name
func (p *PostgresContainer) Database() string {
	return p.impl.Database()
}

// Username returns the PostgreSQL username.
//
// Returns:
//   - string: The username
func (p *PostgresContainer) Username() string {
	return p.impl.Username()
}

// Password returns the PostgreSQL password.
//
// Returns:
//   - string: The password
func (p *PostgresContainer) Password() string {
	return p.impl.Password()
}

// Logs returns the container logs for debugging purposes.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns:
//   - io.Reader: Reader for the container logs
//   - error: Any error that occurred while retrieving logs
func (p *PostgresContainer) Logs(ctx context.Context) (io.Reader, error) {
	return p.impl.Logs(ctx)
}

// Exec executes a command inside the PostgreSQL container.
//
// Parameters:
//   - ctx: Context for the operation
//   - cmd: Command and arguments to execute
//
// Returns:
//   - error: Any error that occurred during command execution
func (p *PostgresContainer) Exec(ctx context.Context, cmd []string) error {
	return p.impl.Exec(ctx, cmd)
}

// Ensure PostgresContainer implements the Container interface
var _ domaincontainer.Container = (*PostgresContainer)(nil)
