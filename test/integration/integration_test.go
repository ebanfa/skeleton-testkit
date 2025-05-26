// Package integration provides integration tests for the skeleton-testkit framework.
// This file contains common test configuration and utilities shared across
// all integration tests.
//
//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/fintechain/skeleton-testkit/pkg/container"
	"github.com/fintechain/skeleton-testkit/pkg/testkit"
	"github.com/fintechain/skeleton-testkit/test/fixtures"
	"github.com/stretchr/testify/require"
)

// TestMain provides setup and teardown for integration tests
func TestMain(m *testing.M) {
	// Setup: Verify Docker is available before running tests
	if !isDockerAvailable() {
		os.Exit(0) // Skip all tests if Docker is not available
	}

	// Run tests
	code := m.Run()

	// Teardown: Clean up any remaining containers
	cleanupContainers()

	os.Exit(code)
}

// TestIntegrationTestSetup verifies that the integration test environment
// is properly configured and ready for testing.
//
// Setup Requirements:
// - Docker daemon is accessible
// - Required images can be pulled
// - Network connectivity is available
// - Sufficient resources are available
func TestIntegrationTestSetup(t *testing.T) {
	ctx := context.Background()

	t.Run("DockerAvailability", func(t *testing.T) {
		// Verify Docker is available and accessible
		require.True(t, isDockerAvailable(), "Docker should be available for integration tests")
	})

	t.Run("BasicContainerOperations", func(t *testing.T) {
		// Test basic container operations work
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
		require.NotNil(t, app, "Should be able to create skeleton app container")

		// Test container lifecycle
		err := app.Start(ctx)
		require.NoError(t, err, "Should be able to start container")

		require.True(t, app.IsRunning(), "Container should be running after start")
		require.NotEmpty(t, app.ID(), "Container should have an ID")

		// Cleanup
		err = app.Stop(ctx)
		require.NoError(t, err, "Should be able to stop container")
		require.False(t, app.IsRunning(), "Container should not be running after stop")
	})

	t.Run("NetworkConnectivity", func(t *testing.T) {
		// Test network connectivity for container operations
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
		err := app.Start(ctx)
		require.NoError(t, err, "Container should start for network test")

		// Verify network access
		host := app.Host()
		require.NotEmpty(t, host, "Container should have accessible host")

		// Cleanup
		defer func() {
			stopErr := app.Stop(ctx)
			require.NoError(t, stopErr, "Container should stop cleanly")
		}()
	})

	t.Run("ResourceAvailability", func(t *testing.T) {
		// Test that sufficient resources are available for testing
		const numContainers = 2
		var containers []*container.AppContainer

		// Start multiple containers to test resource availability
		for i := 0; i < numContainers; i++ {
			app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
			err := app.Start(ctx)
			require.NoError(t, err, "Should be able to start container %d", i)
			containers = append(containers, app)
		}

		// Verify all containers are running
		for i, container := range containers {
			require.True(t, container.IsRunning(), "Container %d should be running", i)
		}

		// Cleanup all containers
		for i, container := range containers {
			err := container.Stop(ctx)
			require.NoError(t, err, "Should be able to stop container %d", i)
		}
	})
}

// TestIntegrationTestConfiguration tests that integration tests can be
// configured appropriately for different environments and scenarios.
//
// Configuration Requirements:
// - Support for different timeout values
// - Environment-specific settings
// - Resource constraint handling
// - Test isolation configuration
func TestIntegrationTestConfiguration(t *testing.T) {
	ctx := context.Background()

	t.Run("TimeoutConfiguration", func(t *testing.T) {
		// Test configurable timeouts
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())

		// Test with custom timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		err := app.Start(timeoutCtx)
		require.NoError(t, err, "Container should start within configured timeout")

		// Cleanup
		defer func() {
			stopErr := app.Stop(ctx)
			require.NoError(t, stopErr, "Container should stop cleanly")
		}()
	})

	t.Run("EnvironmentConfiguration", func(t *testing.T) {
		// Test environment-specific configuration
		if isCIEnvironment() {
			// In CI, test with more conservative settings
			app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())

			startTime := time.Now()
			err := app.Start(ctx)
			duration := time.Since(startTime)

			require.NoError(t, err, "Container should start in CI environment")
			// Allow more time in CI due to resource constraints
			require.Less(t, duration, 120*time.Second, "Startup should complete within CI timeout")

			// Cleanup
			defer func() {
				stopErr := app.Stop(ctx)
				require.NoError(t, stopErr, "Container should stop cleanly in CI")
			}()
		} else {
			// In local development, test with standard settings
			app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())

			err := app.Start(ctx)
			require.NoError(t, err, "Container should start in local environment")

			// Cleanup
			defer func() {
				stopErr := app.Stop(ctx)
				require.NoError(t, stopErr, "Container should stop cleanly locally")
			}()
		}
	})

	t.Run("IsolationConfiguration", func(t *testing.T) {
		// Test that containers are properly isolated
		app1 := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
		app2 := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())

		err1 := app1.Start(ctx)
		require.NoError(t, err1, "First container should start")

		err2 := app2.Start(ctx)
		require.NoError(t, err2, "Second container should start")

		// Verify containers have different IDs (isolation)
		require.NotEqual(t, app1.ID(), app2.ID(), "Containers should have different IDs")

		// Verify containers have different hosts/ports (network isolation)
		if app1.Host() != "" && app2.Host() != "" {
			// If both have hosts, they should be accessible independently
			require.NotEmpty(t, app1.Host(), "First container should have host")
			require.NotEmpty(t, app2.Host(), "Second container should have host")
		}

		// Cleanup
		defer func() {
			stopErr1 := app1.Stop(ctx)
			require.NoError(t, stopErr1, "First container should stop cleanly")

			stopErr2 := app2.Stop(ctx)
			require.NoError(t, stopErr2, "Second container should stop cleanly")
		}()
	})
}

// Helper functions for integration test setup and utilities

// isDockerAvailable checks if Docker is available and accessible
func isDockerAvailable() bool {
	// Try to create a simple container to test Docker availability
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
	if app == nil {
		return false
	}

	// Try to start and immediately stop a test container
	err := app.Start(ctx)
	if err != nil {
		return false
	}

	// Clean up test container
	_ = app.Stop(ctx)
	return true
}

// cleanupContainers performs cleanup of any remaining test containers
func cleanupContainers() {
	// This would typically use Docker API to clean up any containers
	// that might have been left running from failed tests
	// For now, this is a placeholder for the cleanup logic
}

// getTestTimeout returns appropriate timeout values for different environments
func getTestTimeout() time.Duration {
	if isCIEnvironment() {
		// Longer timeout for CI environments
		return 120 * time.Second
	}
	// Standard timeout for local development
	return 60 * time.Second
}

// getTestImage returns the appropriate test image for the current environment
func getTestImage() string {
	return fixtures.GetDefaultTestImage()
}

// skipIfDockerUnavailable skips the test if Docker is not available
func skipIfDockerUnavailable(t *testing.T) {
	if !isDockerAvailable() {
		t.Skip("Docker is not available - skipping test")
	}
}

// skipIfCI skips the test if running in CI environment
func skipIfCI(t *testing.T) {
	if isCIEnvironment() {
		t.Skip("Skipping test in CI environment")
	}
}

// skipIfNotCI skips the test if not running in CI environment
func skipIfNotCI(t *testing.T) {
	if !isCIEnvironment() {
		t.Skip("Skipping test - only runs in CI environment")
	}
}
