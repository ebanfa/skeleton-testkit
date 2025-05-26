package errors

import "fmt"

// StartError represents an error that occurred during container start
type StartError struct {
	Container string
	Cause     error
}

// NewStartError creates a new start error
func NewStartError(container string, cause error) *StartError {
	return &StartError{
		Container: container,
		Cause:     cause,
	}
}

// Error implements the error interface
func (e *StartError) Error() string {
	return fmt.Sprintf("failed to start container %s: %v", e.Container, e.Cause)
}

// Unwrap returns the underlying error
func (e *StartError) Unwrap() error {
	return e.Cause
}

// StopError represents an error that occurred during container stop
type StopError struct {
	Container string
	Cause     error
}

// NewStopError creates a new stop error
func NewStopError(container string, cause error) *StopError {
	return &StopError{
		Container: container,
		Cause:     cause,
	}
}

// Error implements the error interface
func (e *StopError) Error() string {
	return fmt.Sprintf("failed to stop container %s: %v", e.Container, e.Cause)
}

// Unwrap returns the underlying error
func (e *StopError) Unwrap() error {
	return e.Cause
}

// WaitError represents an error that occurred while waiting for container readiness
type WaitError struct {
	Container string
	Cause     error
}

// NewWaitError creates a new wait error
func NewWaitError(container string, cause error) *WaitError {
	return &WaitError{
		Container: container,
		Cause:     cause,
	}
}

// Error implements the error interface
func (e *WaitError) Error() string {
	return fmt.Sprintf("failed to wait for container %s to be ready: %v", e.Container, e.Cause)
}

// Unwrap returns the underlying error
func (e *WaitError) Unwrap() error {
	return e.Cause
}
