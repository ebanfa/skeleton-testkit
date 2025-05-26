package container

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Container represents a containerized service or application for testing
type Container interface {
	// Identity
	ID() string
	Name() string
	Image() string

	// Lifecycle
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	IsRunning() bool

	// Network
	Host() string
	Port(internal int) (int, error)
	ConnectionString() string

	// Health
	WaitForReady(ctx context.Context, timeout time.Duration) error
	HealthCheck(ctx context.Context) error

	// Logs and Debugging
	Logs(ctx context.Context) (io.Reader, error)
	Exec(ctx context.Context, cmd []string) error
}

// ContainerError represents a container-related error
type ContainerError struct {
	Operation string
	Container string
	Message   string
	Cause     error
}

// Error implements the error interface
func (e *ContainerError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("container %s %s failed: %s: %v", e.Container, e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("container %s %s failed: %s", e.Container, e.Operation, e.Message)
}

// Unwrap returns the underlying error
func (e *ContainerError) Unwrap() error {
	return e.Cause
}
