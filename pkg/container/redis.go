package container

import (
	"context"
	"io"
	"time"

	domaincontainer "github.com/fintechain/skeleton-testkit/internal/domain/container"
	"github.com/fintechain/skeleton-testkit/internal/infrastructure/testcontainers"
)

// RedisContainer represents a Redis cache container for testing.
// It provides a clean interface for managing Redis containers used
// as dependencies in skeleton-based application testing.
type RedisContainer struct {
	impl *testcontainers.RedisContainer
}

// NewRedisContainer creates a new RedisContainer wrapper around the internal implementation.
// This follows the constructor injection pattern by accepting the implementation.
//
// Parameters:
//   - impl: The internal RedisContainer implementation
//
// Returns:
//   - *RedisContainer: A new RedisContainer wrapper
func NewRedisContainer(impl *testcontainers.RedisContainer) *RedisContainer {
	return &RedisContainer{
		impl: impl,
	}
}

// Start starts the Redis container.
// This will pull the Redis image if needed and start the container.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns:
//   - error: Any error that occurred during startup
//
// Example:
//
//	err := redis.Start(context.Background())
//	if err != nil {
//	    log.Fatalf("Failed to start Redis: %v", err)
//	}
func (r *RedisContainer) Start(ctx context.Context) error {
	return r.impl.Start(ctx)
}

// Stop stops the Redis container and cleans up resources.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns:
//   - error: Any error that occurred during shutdown
func (r *RedisContainer) Stop(ctx context.Context) error {
	return r.impl.Stop(ctx)
}

// IsRunning returns whether the Redis container is currently running.
//
// Returns:
//   - bool: True if the container is running, false otherwise
func (r *RedisContainer) IsRunning() bool {
	return r.impl.IsRunning()
}

// ID returns the unique identifier of the Redis container.
//
// Returns:
//   - string: The container ID
func (r *RedisContainer) ID() string {
	return r.impl.ID()
}

// Name returns the human-readable name of the Redis container.
//
// Returns:
//   - string: The container name
func (r *RedisContainer) Name() string {
	return r.impl.Name()
}

// Image returns the Docker image name used by the Redis container.
//
// Returns:
//   - string: The Docker image name
func (r *RedisContainer) Image() string {
	return r.impl.Image()
}

// Host returns the host address where the Redis container is accessible.
//
// Returns:
//   - string: The host address
func (r *RedisContainer) Host() string {
	return r.impl.Host()
}

// Port returns the external port mapping for the specified internal port.
//
// Parameters:
//   - internal: The internal port number
//
// Returns:
//   - int: The external port number
//   - error: Any error that occurred
func (r *RedisContainer) Port(internal int) (int, error) {
	return r.impl.Port(internal)
}

// ConnectionString returns the Redis connection string.
// This can be used to connect to the Redis cache from the application.
//
// Returns:
//   - string: The Redis connection string
//
// Example:
//
//	connStr := redis.ConnectionString()
//	// connStr = "redis://localhost:6379" or "redis://:password@localhost:6379"
func (r *RedisContainer) ConnectionString() string {
	return r.impl.ConnectionString()
}

// WaitForReady waits for the Redis container to be ready to accept connections.
//
// Parameters:
//   - ctx: Context for the operation
//   - timeout: Maximum time to wait for readiness
//
// Returns:
//   - error: Any error that occurred while waiting
func (r *RedisContainer) WaitForReady(ctx context.Context, timeout time.Duration) error {
	return r.impl.WaitForReady(ctx, timeout)
}

// HealthCheck performs a health check on the Redis container.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns:
//   - error: Any error that occurred during the health check
func (r *RedisContainer) HealthCheck(ctx context.Context) error {
	return r.impl.HealthCheck(ctx)
}

// Password returns the Redis password if one is configured.
//
// Returns:
//   - string: The password (empty string if no password)
func (r *RedisContainer) Password() string {
	return r.impl.Password()
}

// Logs returns the container logs for debugging purposes.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns:
//   - io.Reader: Reader for the container logs
//   - error: Any error that occurred while retrieving logs
func (r *RedisContainer) Logs(ctx context.Context) (io.Reader, error) {
	return r.impl.Logs(ctx)
}

// Exec executes a command inside the Redis container.
//
// Parameters:
//   - ctx: Context for the operation
//   - cmd: Command and arguments to execute
//
// Returns:
//   - error: Any error that occurred during command execution
func (r *RedisContainer) Exec(ctx context.Context, cmd []string) error {
	return r.impl.Exec(ctx, cmd)
}

// Ensure RedisContainer implements the Container interface
var _ domaincontainer.Container = (*RedisContainer)(nil)
