package testcontainers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/fintechain/skeleton-testkit/internal/domain/container"
	"github.com/fintechain/skeleton-testkit/internal/infrastructure/docker"
)

// TestcontainerAppContainer extends DockerContainer with skeleton-specific functionality
type TestcontainerAppContainer struct {
	*docker.DockerContainer
	skeletonConfig *container.SkeletonConfig
	dependencies   []container.Container
}

// NewTestcontainerAppContainer creates a new TestcontainerAppContainer
func NewTestcontainerAppContainer(config *docker.ContainerConfig, skeletonConfig *container.SkeletonConfig) *TestcontainerAppContainer {
	return &TestcontainerAppContainer{
		DockerContainer: docker.NewDockerContainer(config),
		skeletonConfig:  skeletonConfig,
		dependencies:    make([]container.Container, 0),
	}
}

// AddDependency adds a container dependency
func (t *TestcontainerAppContainer) AddDependency(dep container.Container) {
	t.dependencies = append(t.dependencies, dep)
}

// Start starts the container and its dependencies
func (t *TestcontainerAppContainer) Start(ctx context.Context) error {
	// Start dependencies first
	for _, dep := range t.dependencies {
		if !dep.IsRunning() {
			if err := dep.Start(ctx); err != nil {
				return &container.ContainerError{
					Operation: "start_dependency",
					Container: t.ID(),
					Message:   fmt.Sprintf("failed to start dependency %s", dep.ID()),
					Cause:     err,
				}
			}
		}
	}

	// Create the testcontainer
	if err := t.createContainer(ctx); err != nil {
		return err
	}

	// Start the container
	return t.DockerContainer.Start(ctx)
}

// createContainer creates the underlying testcontainer
func (t *TestcontainerAppContainer) createContainer(ctx context.Context) error {
	config := t.Config()

	// Build environment variables
	env := make(map[string]string)
	for k, v := range config.Environment {
		env[k] = v
	}

	// Add skeleton-specific environment variables
	if t.skeletonConfig != nil {
		// Serialize complete skeleton config as JSON
		skeletonConfigJSON, err := json.Marshal(t.skeletonConfig)
		if err != nil {
			return &container.ContainerError{
				Operation: "serialize_skeleton_config",
				Container: t.ID(),
				Message:   "failed to serialize skeleton configuration",
				Cause:     err,
			}
		}
		env["SKELETON_CONFIG"] = string(skeletonConfigJSON)

		// Keep backward compatibility with individual fields
		if t.skeletonConfig.ServiceID != "" {
			env["SKELETON_SERVICE_ID"] = t.skeletonConfig.ServiceID
		}
		if t.skeletonConfig.Storage.Type != "" {
			env["SKELETON_STORAGE_TYPE"] = t.skeletonConfig.Storage.Type
		}
		if t.skeletonConfig.Storage.URL != "" {
			env["SKELETON_STORAGE_URL"] = t.skeletonConfig.Storage.URL
		}
	}

	// Build exposed ports
	exposedPorts := make([]string, 0)
	for _, port := range config.Ports {
		exposedPorts = append(exposedPorts, fmt.Sprintf("%d/tcp", port.Internal))
	}

	// Create container request
	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		Name:         config.Name,
		Env:          env,
		ExposedPorts: exposedPorts,
		WaitingFor:   wait.ForListeningPort("8080/tcp").WithStartupTimeout(30 * time.Second),
	}

	// Create the container
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false, // We'll start it manually
	})
	if err != nil {
		return &container.ContainerError{
			Operation: "create",
			Container: t.ID(),
			Message:   "failed to create testcontainer",
			Cause:     err,
		}
	}

	t.SetContainer(c)
	return nil
}

// Stop stops the container
func (t *TestcontainerAppContainer) Stop(ctx context.Context) error {
	// Stop the main container first
	if err := t.DockerContainer.Stop(ctx); err != nil {
		return err
	}

	// Stop dependencies in reverse order
	for i := len(t.dependencies) - 1; i >= 0; i-- {
		dep := t.dependencies[i]
		if dep.IsRunning() {
			if err := dep.Stop(ctx); err != nil {
				// Log error but continue stopping other dependencies
				continue
			}
		}
	}

	return nil
}

// WaitForReady waits for the container and its dependencies to be ready
func (t *TestcontainerAppContainer) WaitForReady(ctx context.Context, timeout time.Duration) error {
	// Wait for dependencies first
	for _, dep := range t.dependencies {
		if err := dep.WaitForReady(ctx, timeout); err != nil {
			return &container.ContainerError{
				Operation: "wait_dependency",
				Container: t.ID(),
				Message:   fmt.Sprintf("dependency %s not ready", dep.ID()),
				Cause:     err,
			}
		}
	}

	// Wait for main container
	if err := t.DockerContainer.WaitForReady(ctx, timeout); err != nil {
		return err
	}

	// Validate skeleton endpoints if this is a skeleton application
	if t.skeletonConfig != nil {
		if err := t.validateSkeletonEndpoints(ctx); err != nil {
			return &container.ContainerError{
				Operation: "validate_skeleton_endpoints",
				Container: t.ID(),
				Message:   "skeleton endpoints validation failed",
				Cause:     err,
			}
		}
	}

	return nil
}

// validateSkeletonEndpoints validates that skeleton-specific endpoints are accessible
func (t *TestcontainerAppContainer) validateSkeletonEndpoints(ctx context.Context) error {
	baseURL := t.ConnectionString()
	if baseURL == "" {
		return fmt.Errorf("unable to get connection string for skeleton endpoint validation")
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Validate skeleton system service endpoint
	systemURL := baseURL + "/api/system/health"
	if err := t.validateEndpoint(ctx, client, systemURL, "skeleton system service"); err != nil {
		return err
	}

	// Validate skeleton components endpoint
	componentsURL := baseURL + "/api/components"
	if err := t.validateEndpoint(ctx, client, componentsURL, "skeleton components"); err != nil {
		return err
	}

	return nil
}

// validateEndpoint validates that a specific endpoint is accessible
func (t *TestcontainerAppContainer) validateEndpoint(ctx context.Context, client *http.Client, url, name string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for %s endpoint: %w", name, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to reach %s endpoint at %s: %w", name, url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s endpoint returned status %d", name, resp.StatusCode)
	}

	return nil
}

// SkeletonConfig returns the skeleton configuration
func (t *TestcontainerAppContainer) SkeletonConfig() *container.SkeletonConfig {
	return t.skeletonConfig
}

// Dependencies returns the container dependencies
func (t *TestcontainerAppContainer) Dependencies() []container.Container {
	return t.dependencies
}

// ConnectionString returns a connection string for the container
func (t *TestcontainerAppContainer) ConnectionString() string {
	host := t.Host()
	port, err := t.Port(8080) // Default application port
	if err != nil {
		return ""
	}
	return fmt.Sprintf("http://%s:%d", host, port)
}

// HealthEndpoint returns the health check endpoint URL
func (t *TestcontainerAppContainer) HealthEndpoint() string {
	if t.skeletonConfig != nil {
		// Default health endpoint for skeleton applications
		return "/health"
	}
	return "/health"
}

// MetricsEndpoint returns the metrics endpoint URL
func (t *TestcontainerAppContainer) MetricsEndpoint() string {
	if t.skeletonConfig != nil {
		// Default metrics endpoint for skeleton applications
		return "/metrics"
	}
	return "/metrics"
}

func (t *TestcontainerAppContainer) ShutdownEndpoint() string {
	if t.skeletonConfig != nil {
		// Default shutdown endpoint for skeleton applications
		return "/shutdown"
	}
	return "/shutdown"
}
