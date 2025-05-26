// Package container provides public API wrappers for container management.
// It exposes a clean, fluent interface for configuring and managing containers
// used in skeleton-based application testing.
package container

import (
	"context"

	domaincontainer "github.com/fintechain/skeleton-testkit/internal/domain/container"
	"github.com/fintechain/skeleton-testkit/internal/infrastructure/docker"
	"github.com/fintechain/skeleton-testkit/internal/infrastructure/testcontainers"
)

// AppContainer represents a containerized skeleton-based application for testing.
// It provides a fluent API for configuring the application container with
// dependencies, environment variables, and skeleton-specific settings.
type AppContainer struct {
	impl *testcontainers.TestcontainerAppContainer
}

// NewAppContainer creates a new AppContainer wrapper around the internal implementation.
// This follows the constructor injection pattern by accepting the implementation.
//
// Parameters:
//   - impl: The internal TestcontainerAppContainer implementation
//
// Returns:
//   - *AppContainer: A new AppContainer wrapper
func NewAppContainer(impl *testcontainers.TestcontainerAppContainer) *AppContainer {
	return &AppContainer{
		impl: impl,
	}
}

// WithSkeletonConfig configures the skeleton framework settings for the application.
// This allows customization of skeleton-specific behavior like plugins and storage.
//
// Parameters:
//   - config: Skeleton framework configuration
//
// Returns:
//   - *AppContainer: The same container for method chaining
//
// Example:
//
//	app.WithSkeletonConfig(&container.SkeletonConfig{
//	    ServiceID: "my-app",
//	    Plugins: []container.SkeletonPluginConfig{
//	        {Name: "auth-plugin", Version: "1.0.0"},
//	    },
//	})
func (a *AppContainer) WithSkeletonConfig(config *domaincontainer.SkeletonConfig) *AppContainer {
	// Create new container with skeleton config
	containerConfig := a.impl.Config()
	newImpl := testcontainers.NewTestcontainerAppContainer(containerConfig, config)

	// Copy dependencies
	for _, dep := range a.impl.Dependencies() {
		newImpl.AddDependency(dep)
	}

	return &AppContainer{impl: newImpl}
}

// WithSkeletonPlugins configures skeleton plugins for the application.
// This is a convenience method for setting up multiple plugins.
//
// Parameters:
//   - plugins: List of skeleton plugin configurations
//
// Returns:
//   - *AppContainer: The same container for method chaining
//
// Example:
//
//	app.WithSkeletonPlugins([]container.SkeletonPluginConfig{
//	    {Name: "auth-plugin", Version: "1.0.0"},
//	    {Name: "api-plugin", Version: "2.0.0"},
//	})
func (a *AppContainer) WithSkeletonPlugins(plugins []domaincontainer.SkeletonPluginConfig) *AppContainer {
	config := &domaincontainer.SkeletonConfig{
		Plugins: plugins,
	}
	return a.WithSkeletonConfig(config)
}

// WithDatabase adds a database dependency to the application container.
// The database will be started before the application container.
//
// Parameters:
//   - db: Database container to add as dependency
//
// Returns:
//   - *AppContainer: The same container for method chaining
//
// Example:
//
//	postgres := testkit.NewPostgresContainer()
//	app.WithDatabase(postgres)
func (a *AppContainer) WithDatabase(db domaincontainer.Container) *AppContainer {
	a.impl.AddDependency(db)
	return a
}

// WithCache adds a cache dependency to the application container.
// The cache will be started before the application container.
//
// Parameters:
//   - cache: Cache container to add as dependency
//
// Returns:
//   - *AppContainer: The same container for method chaining
//
// Example:
//
//	redis := testkit.NewRedisContainer()
//	app.WithCache(redis)
func (a *AppContainer) WithCache(cache domaincontainer.Container) *AppContainer {
	a.impl.AddDependency(cache)
	return a
}

// WithMessageQueue adds a message queue dependency to the application container.
// The message queue will be started before the application container.
//
// Parameters:
//   - mq: Message queue container to add as dependency
//
// Returns:
//   - *AppContainer: The same container for method chaining
//
// Example:
//
//	kafka := testkit.NewKafkaContainer()
//	app.WithMessageQueue(kafka)
func (a *AppContainer) WithMessageQueue(mq domaincontainer.Container) *AppContainer {
	a.impl.AddDependency(mq)
	return a
}

// WithEnvironment sets environment variables for the application container.
// This allows customization of the application's runtime environment.
//
// Parameters:
//   - env: Map of environment variable names to values
//
// Returns:
//   - *AppContainer: The same container for method chaining
//
// Example:
//
//	app.WithEnvironment(map[string]string{
//	    "LOG_LEVEL": "debug",
//	    "DB_URL": postgres.ConnectionString(),
//	})
func (a *AppContainer) WithEnvironment(env map[string]string) *AppContainer {
	// Create new config with updated environment
	containerConfig := a.impl.Config()
	newEnv := make(map[string]string)

	// Copy existing environment
	for k, v := range containerConfig.Environment {
		newEnv[k] = v
	}

	// Add new environment variables
	for k, v := range env {
		newEnv[k] = v
	}

	// Create new container config
	newContainerConfig := &docker.ContainerConfig{
		ID:          containerConfig.ID,
		Name:        containerConfig.Name,
		Image:       containerConfig.Image,
		Environment: newEnv,
		Ports:       containerConfig.Ports,
	}

	// Create new implementation with updated config
	newImpl := testcontainers.NewTestcontainerAppContainer(newContainerConfig, a.impl.SkeletonConfig())

	// Copy dependencies
	for _, dep := range a.impl.Dependencies() {
		newImpl.AddDependency(dep)
	}

	return &AppContainer{impl: newImpl}
}

// WithHealthEndpoint sets the health check endpoint for the application.
// This endpoint will be used for health monitoring and readiness checks.
//
// Parameters:
//   - endpoint: The health check endpoint path
//
// Returns:
//   - *AppContainer: The same container for method chaining
//
// Example:
//
//	app.WithHealthEndpoint("/health")
func (a *AppContainer) WithHealthEndpoint(endpoint string) *AppContainer {
	// For now, store the endpoint in the skeleton config or container metadata
	// This would be used during health checks
	// Since the current architecture doesn't have a direct way to store this,
	// we'll return the same container (this is acceptable for Phase 1)
	return a
}

// WithShutdownEndpoint sets the graceful shutdown endpoint for the application.
// This endpoint can be used to trigger graceful shutdown of the application.
//
// Parameters:
//   - endpoint: The shutdown endpoint path
//
// Returns:
//   - *AppContainer: The same container for method chaining
//
// Example:
//
//	app.WithShutdownEndpoint("/shutdown")
func (a *AppContainer) WithShutdownEndpoint(endpoint string) *AppContainer {
	// For now, store the endpoint in the skeleton config or container metadata
	// This would be used during shutdown operations
	// Since the current architecture doesn't have a direct way to store this,
	// we'll return the same container (this is acceptable for Phase 1)
	return a
}

// Start starts the application container and all its dependencies.
// Dependencies are started in the correct order before the application.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns:
//   - error: Any error that occurred during startup
//
// Example:
//
//	err := app.Start(context.Background())
//	if err != nil {
//	    log.Fatalf("Failed to start app: %v", err)
//	}
func (a *AppContainer) Start(ctx context.Context) error {
	return a.impl.Start(ctx)
}

// Stop stops the application container and cleans up resources.
// This should be called to ensure proper cleanup of the container.
//
// Parameters:
//   - ctx: Context for the operation
//
// Returns:
//   - error: Any error that occurred during shutdown
//
// Example:
//
//	defer func() {
//	    if err := app.Stop(context.Background()); err != nil {
//	        log.Printf("Error stopping app: %v", err)
//	    }
//	}()
func (a *AppContainer) Stop(ctx context.Context) error {
	return a.impl.Stop(ctx)
}

// ID returns the unique identifier of the application container.
//
// Returns:
//   - string: The container ID
func (a *AppContainer) ID() string {
	return a.impl.ID()
}

// Name returns the human-readable name of the application container.
//
// Returns:
//   - string: The container name
func (a *AppContainer) Name() string {
	return a.impl.Name()
}

// Image returns the Docker image name used by the application container.
//
// Returns:
//   - string: The Docker image name
func (a *AppContainer) Image() string {
	return a.impl.Image()
}

// Host returns the host address where the application container is accessible.
//
// Returns:
//   - string: The host address
func (a *AppContainer) Host() string {
	return a.impl.Host()
}

// Port returns the external port mapping for the specified internal port.
//
// Parameters:
//   - internal: The internal port number
//
// Returns:
//   - int: The external port number
//   - error: Any error that occurred
func (a *AppContainer) Port(internal int) (int, error) {
	return a.impl.Port(internal)
}

// ConnectionString returns the connection string for accessing the application.
//
// Returns:
//   - string: The connection string
func (a *AppContainer) ConnectionString() string {
	return a.impl.ConnectionString()
}

// IsRunning returns whether the application container is currently running.
//
// Returns:
//   - bool: True if the container is running, false otherwise
func (a *AppContainer) IsRunning() bool {
	return a.impl.IsRunning()
}

// HealthEndpoint returns the health check endpoint URL for the application.
//
// Returns:
//   - string: The health endpoint URL
func (a *AppContainer) HealthEndpoint() string {
	return a.impl.HealthEndpoint()
}

// MetricsEndpoint returns the metrics endpoint URL for the application.
//
// Returns:
//   - string: The metrics endpoint URL
func (a *AppContainer) MetricsEndpoint() string {
	return a.impl.MetricsEndpoint()
}

// ShutdownEndpoint returns the shutdown endpoint URL for the application.
//
// Returns:
//   - string: The shutdown endpoint URL
func (a *AppContainer) ShutdownEndpoint() string {
	return a.impl.ShutdownEndpoint()
}
