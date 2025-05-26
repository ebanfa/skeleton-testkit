// Package verification provides verification strategies for testing skeleton applications.
// It implements the verification framework specified in the skeleton-testkit specification.
package verification

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fintechain/skeleton-testkit/pkg/container"
)

// SystemVerifier verifies skeleton application system-level behavior
type SystemVerifier struct {
	app *container.AppContainer
}

// NewSystemVerifier creates a new SystemVerifier for the given application container
func NewSystemVerifier(app *container.AppContainer) *SystemVerifier {
	return &SystemVerifier{
		app: app,
	}
}

// VerifySkeletonStartup verifies that the skeleton application starts successfully
// and all skeleton components are properly initialized
func (s *SystemVerifier) VerifySkeletonStartup(ctx context.Context) error {
	if !s.app.IsRunning() {
		return fmt.Errorf("skeleton application is not running")
	}

	// Verify skeleton system service endpoint
	if err := s.verifySkeletonSystemService(ctx); err != nil {
		return fmt.Errorf("skeleton system service verification failed: %w", err)
	}

	// Verify basic health endpoint
	if err := s.verifyHealthEndpoint(ctx); err != nil {
		return fmt.Errorf("health endpoint verification failed: %w", err)
	}

	return nil
}

// VerifySkeletonShutdown verifies that the skeleton application shuts down gracefully
func (s *SystemVerifier) VerifySkeletonShutdown(ctx context.Context) error {
	if !s.app.IsRunning() {
		return nil // Already stopped
	}

	// Stop the application
	if err := s.app.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop skeleton application: %w", err)
	}

	// Verify it's no longer running
	if s.app.IsRunning() {
		return fmt.Errorf("skeleton application is still running after stop")
	}

	return nil
}

// VerifySkeletonHealth verifies the health status of the skeleton application
func (s *SystemVerifier) VerifySkeletonHealth(ctx context.Context) error {
	if !s.app.IsRunning() {
		return fmt.Errorf("skeleton application is not running")
	}

	// Check health endpoint
	if err := s.verifyHealthEndpoint(ctx); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	return nil
}

// VerifySkeletonSystemService verifies the skeleton system service is operational
func (s *SystemVerifier) VerifySkeletonSystemService(ctx context.Context) error {
	if !s.app.IsRunning() {
		return fmt.Errorf("skeleton application is not running")
	}

	return s.verifySkeletonSystemService(ctx)
}

// verifySkeletonSystemService checks the skeleton system service endpoint
func (s *SystemVerifier) verifySkeletonSystemService(ctx context.Context) error {
	baseURL := s.app.ConnectionString()
	if baseURL == "" {
		return fmt.Errorf("unable to get application connection string")
	}

	// Check skeleton system service endpoint
	systemURL := baseURL + "/api/system/health"

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", systemURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to reach skeleton system service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("skeleton system service returned status %d", resp.StatusCode)
	}

	return nil
}

// verifyHealthEndpoint checks the basic health endpoint
func (s *SystemVerifier) verifyHealthEndpoint(ctx context.Context) error {
	baseURL := s.app.ConnectionString()
	if baseURL == "" {
		return fmt.Errorf("unable to get application connection string")
	}

	healthURL := baseURL + s.app.HealthEndpoint()

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to reach health endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health endpoint returned status %d", resp.StatusCode)
	}

	return nil
}
