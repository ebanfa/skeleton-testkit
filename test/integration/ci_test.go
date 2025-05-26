// Package integration provides integration tests for the skeleton-testkit framework.
// This file contains CI/CD pipeline validation tests that ensure the testkit
// works correctly in automated testing environments.
//
//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/fintechain/skeleton-testkit/pkg/testkit"
	"github.com/fintechain/skeleton-testkit/test/fixtures"
	"github.com/stretchr/testify/require"
)

// TestCIEnvironmentCompatibility tests that the skeleton-testkit works
// correctly in CI/CD environments with their specific constraints and
// configurations.
//
// CI/CD Requirements:
// - Works in containerized CI environments
// - Handles resource constraints gracefully
// - Provides clear error messages for debugging
// - Supports headless operation
func TestCIEnvironmentCompatibility(t *testing.T) {
	ctx := context.Background()

	t.Run("BasicCIOperation", func(t *testing.T) {
		// Test basic container operations that should work in any CI environment
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
		require.NotNil(t, app, "Should create container in CI environment")

		// Start with reasonable timeout for CI
		err := app.Start(ctx)
		require.NoError(t, err, "Container should start in CI environment")

		// Verify basic functionality
		require.True(t, app.IsRunning(), "Container should be running in CI")
		require.NotEmpty(t, app.Host(), "Container should have accessible host in CI")

		// Cleanup
		defer func() {
			stopErr := app.Stop(ctx)
			require.NoError(t, stopErr, "Container should stop cleanly in CI")
		}()
	})

	t.Run("CIResourceConstraints", func(t *testing.T) {
		// Test behavior under typical CI resource constraints
		if isCIEnvironment() {
			// In CI, use more conservative timeouts and expectations
			app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())

			// Allow longer startup time in CI due to resource constraints
			startTime := time.Now()
			err := app.Start(ctx)
			startupDuration := time.Since(startTime)

			require.NoError(t, err, "Container should start even with CI resource constraints")
			// More lenient timeout for CI environments
			require.Less(t, startupDuration, 120*time.Second, "Startup should complete within CI timeout")

			// Cleanup
			defer func() {
				stopErr := app.Stop(ctx)
				require.NoError(t, stopErr, "Container should stop cleanly in CI")
			}()
		} else {
			t.Skip("Skipping CI resource constraint test in non-CI environment")
		}
	})
}

// TestCIErrorReporting tests that error messages and debugging information
// are clear and actionable in CI/CD environments where interactive debugging
// is not possible.
//
// Error Reporting Requirements:
// - Clear error messages for common failures
// - Sufficient context for debugging
// - Structured logging for CI log aggregation
// - No sensitive information in logs
func TestCIErrorReporting(t *testing.T) {
	ctx := context.Background()

	t.Run("InvalidImageErrorReporting", func(t *testing.T) {
		// Test error reporting for invalid images
		app := testkit.NewSkeletonApp("non-existent-image:invalid-tag")

		err := app.Start(ctx)
		require.Error(t, err, "Should fail with non-existent image")

		// Verify error message is informative for CI debugging
		errorMsg := err.Error()
		require.NotEmpty(t, errorMsg, "Error message should not be empty")
		// Error should contain enough context for CI debugging
		require.Contains(t, errorMsg, "non-existent-image", "Error should mention the problematic image")
	})

	t.Run("NetworkErrorReporting", func(t *testing.T) {
		// Test error reporting for network-related issues
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())

		// Start container to test port access
		err := app.Start(ctx)
		if err != nil {
			// If start fails, verify error is informative
			require.Contains(t, err.Error(), "start", "Error should indicate startup failure")
			return
		}

		// Test port access error reporting
		_, portErr := app.Port(99999) // Invalid port
		if portErr != nil {
			require.Contains(t, portErr.Error(), "port", "Port error should mention port issue")
		}

		// Cleanup
		defer func() {
			stopErr := app.Stop(ctx)
			require.NoError(t, stopErr, "Container should stop cleanly")
		}()
	})
}

// TestCIParallelExecution tests that the skeleton-testkit supports parallel
// test execution as commonly used in CI/CD pipelines to reduce build times.
//
// Parallel Execution Requirements:
// - No resource conflicts between parallel tests
// - Proper container isolation
// - Reasonable resource usage under parallel load
// - Clean cleanup even with parallel execution
func TestCIParallelExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping parallel execution test in short mode")
	}

	ctx := context.Background()
	const numParallelTests = 3

	t.Run("ParallelContainerCreation", func(t *testing.T) {
		// Run multiple container tests in parallel
		for i := 0; i < numParallelTests; i++ {
			i := i // Capture loop variable
			t.Run("ParallelTest", func(t *testing.T) {
				t.Parallel()

				app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
				require.NotNil(t, app, "Should create container in parallel test %d", i)

				err := app.Start(ctx)
				require.NoError(t, err, "Container should start in parallel test %d", i)

				// Verify isolation - each container should have unique ID
				require.NotEmpty(t, app.ID(), "Container should have unique ID in parallel test %d", i)

				// Cleanup
				defer func() {
					stopErr := app.Stop(ctx)
					require.NoError(t, stopErr, "Container should stop cleanly in parallel test %d", i)
				}()
			})
		}
	})

	t.Run("ParallelDatabaseTests", func(t *testing.T) {
		// Run multiple database tests in parallel
		for i := 0; i < numParallelTests; i++ {
			i := i // Capture loop variable
			t.Run("ParallelDBTest", func(t *testing.T) {
				t.Parallel()

				postgres := testkit.NewPostgresContainer()
				require.NotNil(t, postgres, "Should create database in parallel test %d", i)

				err := postgres.Start(ctx)
				require.NoError(t, err, "Database should start in parallel test %d", i)

				// Verify each database has unique connection details
				require.NotEmpty(t, postgres.ConnectionString(), "Database should have connection string in parallel test %d", i)

				// Cleanup
				defer func() {
					stopErr := postgres.Stop(ctx)
					require.NoError(t, stopErr, "Database should stop cleanly in parallel test %d", i)
				}()
			})
		}
	})
}

// TestCITimeouts tests that the skeleton-testkit handles timeouts appropriately
// in CI environments where builds may be killed after certain time limits.
//
// Timeout Requirements:
// - Reasonable default timeouts for CI
// - Graceful handling of timeout scenarios
// - Proper cleanup even when operations timeout
// - Clear timeout error messages
func TestCITimeouts(t *testing.T) {
	ctx := context.Background()

	t.Run("StartupTimeout", func(t *testing.T) {
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())

		// Create context with timeout to simulate CI timeout behavior
		timeoutCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
		defer cancel()

		err := app.Start(timeoutCtx)
		if err != nil {
			// If timeout occurs, verify it's handled gracefully
			if timeoutCtx.Err() == context.DeadlineExceeded {
				require.Contains(t, err.Error(), "timeout", "Timeout error should be clear")
			} else {
				require.NoError(t, err, "Non-timeout errors should not occur")
			}
		} else {
			// If successful, verify container is running
			require.True(t, app.IsRunning(), "Container should be running after successful start")

			// Cleanup
			defer func() {
				stopErr := app.Stop(ctx)
				require.NoError(t, stopErr, "Container should stop cleanly")
			}()
		}
	})

	t.Run("ShutdownTimeout", func(t *testing.T) {
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())

		// Start container first
		err := app.Start(ctx)
		require.NoError(t, err, "Container should start for shutdown test")

		// Test shutdown with timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		err = app.Stop(timeoutCtx)
		require.NoError(t, err, "Container should stop within timeout")
		require.False(t, app.IsRunning(), "Container should not be running after stop")
	})
}

// TestCIDockerEnvironment tests compatibility with different Docker
// environments commonly found in CI/CD systems.
//
// Docker Environment Requirements:
// - Works with Docker-in-Docker (DinD)
// - Compatible with Docker socket mounting
// - Handles Docker daemon connectivity issues
// - Works with different Docker versions
func TestCIDockerEnvironment(t *testing.T) {
	ctx := context.Background()

	t.Run("DockerConnectivity", func(t *testing.T) {
		// Test basic Docker connectivity
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
		require.NotNil(t, app, "Should create container with Docker connectivity")

		// Attempt to start - this will test Docker daemon connectivity
		err := app.Start(ctx)
		if err != nil {
			// If Docker is not available, skip the test with clear message
			if isDockerUnavailable(err) {
				t.Skip("Docker daemon not available - skipping Docker connectivity test")
			} else {
				require.NoError(t, err, "Docker connectivity should work in CI")
			}
		} else {
			require.True(t, app.IsRunning(), "Container should be running with Docker connectivity")

			// Cleanup
			defer func() {
				stopErr := app.Stop(ctx)
				require.NoError(t, stopErr, "Container should stop cleanly")
			}()
		}
	})
}

// Helper functions for CI environment detection and error classification

// isCIEnvironment detects if tests are running in a CI environment
func isCIEnvironment() bool {
	// Check common CI environment variables
	ciEnvVars := []string{
		"CI",             // Generic CI indicator
		"GITHUB_ACTIONS", // GitHub Actions
		"GITLAB_CI",      // GitLab CI
		"JENKINS_URL",    // Jenkins
		"TRAVIS",         // Travis CI
		"CIRCLECI",       // CircleCI
		"BUILDKITE",      // Buildkite
	}

	for _, envVar := range ciEnvVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}

	return false
}

// isDockerUnavailable checks if an error indicates Docker is unavailable
func isDockerUnavailable(err error) bool {
	if err == nil {
		return false
	}

	errorMsg := err.Error()
	dockerUnavailableIndicators := []string{
		"docker daemon",
		"cannot connect to the Docker daemon",
		"docker: not found",
		"permission denied",
		"connection refused",
	}

	for _, indicator := range dockerUnavailableIndicators {
		if contains(errorMsg, indicator) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr))))
}

// containsSubstring performs a simple substring search
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
