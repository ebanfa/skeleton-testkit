package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"

	"github.com/fintechain/skeleton-testkit/internal/domain/container"
)

// DockerContainer wraps testcontainers.Container with additional configuration
type DockerContainer struct {
	container testcontainers.Container
	config    *ContainerConfig
}

// ContainerConfig holds basic container configuration
type ContainerConfig struct {
	ID          string
	Name        string
	Image       string
	Environment map[string]string
	Ports       []container.PortMapping
}

// NewDockerContainer creates a new DockerContainer with the given configuration
func NewDockerContainer(config *ContainerConfig) *DockerContainer {
	return &DockerContainer{
		config: config,
	}
}

// SetContainer sets the underlying testcontainers.Container
func (d *DockerContainer) SetContainer(c testcontainers.Container) {
	d.container = c
}

// ID returns the unique identifier of the container
func (d *DockerContainer) ID() string {
	return d.config.ID
}

// Name returns the name of the container
func (d *DockerContainer) Name() string {
	return d.config.Name
}

// Image returns the image of the container
func (d *DockerContainer) Image() string {
	return d.config.Image
}

// Start starts the container
func (d *DockerContainer) Start(ctx context.Context) error {
	if d.container == nil {
		return &container.ContainerError{
			Operation: "start",
			Container: d.ID(),
			Message:   "container not initialized",
		}
	}

	err := d.container.Start(ctx)
	if err != nil {
		return &container.ContainerError{
			Operation: "start",
			Container: d.ID(),
			Message:   "failed to start container",
			Cause:     err,
		}
	}

	return nil
}

// Stop stops the container
func (d *DockerContainer) Stop(ctx context.Context) error {
	if d.container == nil {
		return &container.ContainerError{
			Operation: "stop",
			Container: d.ID(),
			Message:   "container not initialized",
		}
	}

	err := d.container.Stop(ctx, nil)
	if err != nil {
		return &container.ContainerError{
			Operation: "stop",
			Container: d.ID(),
			Message:   "failed to stop container",
			Cause:     err,
		}
	}

	return nil
}

// IsRunning returns true if the container is currently running
func (d *DockerContainer) IsRunning() bool {
	if d.container == nil {
		return false
	}

	ctx := context.Background()
	state, err := d.container.State(ctx)
	if err != nil {
		return false
	}

	return state.Running
}

// Host returns the host address where the container is accessible
func (d *DockerContainer) Host() string {
	if d.container == nil {
		return ""
	}

	ctx := context.Background()
	host, err := d.container.Host(ctx)
	if err != nil {
		return "localhost"
	}

	return host
}

// Port returns the mapped external port for the given internal port
func (d *DockerContainer) Port(internal int) (int, error) {
	if d.container == nil {
		return 0, &container.ContainerError{
			Operation: "port",
			Container: d.config.ID,
			Message:   "container not started",
		}
	}

	ctx := context.Background()
	port, err := d.container.MappedPort(ctx, nat.Port(fmt.Sprintf("%d/tcp", internal)))
	if err != nil {
		return 0, &container.ContainerError{
			Operation: "port",
			Container: d.config.ID,
			Message:   "failed to get mapped port",
			Cause:     err,
		}
	}

	return port.Int(), nil
}

// ConnectionString returns the connection string for the container
func (d *DockerContainer) ConnectionString() string {
	host := d.Host()
	if host == "" {
		return ""
	}

	// For basic containers, return host:port format
	// Specialized containers can override this method
	if len(d.config.Ports) > 0 {
		port, err := d.Port(d.config.Ports[0].Internal)
		if err != nil {
			return ""
		}
		return fmt.Sprintf("%s:%d", host, port)
	}

	return host
}

// WaitForReady waits for the container to be ready with a timeout
func (d *DockerContainer) WaitForReady(ctx context.Context, timeout time.Duration) error {
	if d.container == nil {
		return &container.ContainerError{
			Operation: "wait_for_ready",
			Container: d.config.ID,
			Message:   "container not started",
		}
	}

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Wait for container to be ready
	// This is a basic implementation - specialized containers can override
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutCtx.Done():
			return &container.ContainerError{
				Operation: "wait_for_ready",
				Container: d.config.ID,
				Message:   "timeout waiting for container to be ready",
				Cause:     timeoutCtx.Err(),
			}
		case <-ticker.C:
			if d.IsRunning() {
				return nil
			}
		}
	}
}

// HealthCheck performs a health check on the container
func (d *DockerContainer) HealthCheck(ctx context.Context) error {
	if !d.IsRunning() {
		return &container.ContainerError{
			Operation: "health_check",
			Container: d.config.ID,
			Message:   "container is not running",
		}
	}

	// Basic health check - just verify container is running
	// Specialized containers can override with more sophisticated checks
	return nil
}

// Logs returns the container logs
func (d *DockerContainer) Logs(ctx context.Context) (io.Reader, error) {
	if d.container == nil {
		return nil, &container.ContainerError{
			Operation: "logs",
			Container: d.ID(),
			Message:   "container not initialized",
		}
	}

	logs, err := d.container.Logs(ctx)
	if err != nil {
		return nil, &container.ContainerError{
			Operation: "logs",
			Container: d.ID(),
			Message:   "failed to get container logs",
			Cause:     err,
		}
	}

	return logs, nil
}

// Exec executes a command in the container
func (d *DockerContainer) Exec(ctx context.Context, cmd []string) error {
	if d.container == nil {
		return &container.ContainerError{
			Operation: "exec",
			Container: d.ID(),
			Message:   "container not initialized",
		}
	}

	_, _, err := d.container.Exec(ctx, cmd)
	if err != nil {
		return &container.ContainerError{
			Operation: "exec",
			Container: d.ID(),
			Message:   fmt.Sprintf("failed to execute command: %v", cmd),
			Cause:     err,
		}
	}

	return nil
}

// Config returns the container configuration
func (d *DockerContainer) Config() *ContainerConfig {
	return d.config
}

// Container returns the underlying testcontainers.Container
func (d *DockerContainer) Container() testcontainers.Container {
	return d.container
}
