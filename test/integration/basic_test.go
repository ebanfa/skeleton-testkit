// Package integration provides integration tests for the skeleton-testkit framework.
// These tests verify that the testkit can properly manage containerized skeleton applications
// and their dependencies in realistic testing scenarios.
//
//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/fintechain/skeleton-testkit/pkg/health"
	"github.com/fintechain/skeleton-testkit/pkg/testkit"
	"github.com/fintechain/skeleton-testkit/pkg/verification"
	"github.com/fintechain/skeleton-testkit/test/fixtures"
	"github.com/stretchr/testify/require"
)

// TestBasicSkeletonApp tests the basic functionality of creating and starting
// a skeleton application container. This verifies the core container management
// capabilities of the testkit.
//
// Test Coverage:
// - Container creation with skeleton application image
// - Container startup and lifecycle management
// - Basic error handling and cleanup
// - Container state verification
// - System verification using the verification package
func TestBasicSkeletonApp(t *testing.T) {
	// Create container for skeleton-based application using realistic test image
	app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
	require.NotNil(t, app, "NewSkeletonApp should return a valid container")

	// Verify initial state
	require.False(t, app.IsRunning(), "Container should not be running initially")
	require.Equal(t, fixtures.GetDefaultTestImage(), app.Image(), "Container should have correct image")

	// Start the container
	ctx := context.Background()
	err := app.Start(ctx)
	require.NoError(t, err, "Container should start successfully")

	// Verify running state
	require.True(t, app.IsRunning(), "Container should be running after start")
	require.NotEmpty(t, app.ID(), "Container should have a valid ID")
	require.NotEmpty(t, app.Host(), "Container should have a valid host")

	// Use SystemVerifier to verify skeleton startup
	verifier := verification.NewSystemVerifier(app)
	err = verifier.VerifySkeletonStartup(ctx)
	require.NoError(t, err, "Skeleton system should start up correctly")

	// Verify skeleton health
	err = verifier.VerifySkeletonHealth(ctx)
	require.NoError(t, err, "Skeleton system should be healthy")

	// Ensure cleanup happens
	defer func() {
		// Verify graceful shutdown
		shutdownErr := verifier.VerifySkeletonShutdown(ctx)
		require.NoError(t, shutdownErr, "Skeleton should shutdown gracefully")

		stopErr := app.Stop(ctx)
		require.NoError(t, stopErr, "Container should stop cleanly")
		require.False(t, app.IsRunning(), "Container should not be running after stop")
	}()
}

// TestSkeletonAppLifecycle tests the complete lifecycle of a skeleton application
// container including startup, health verification, and graceful shutdown.
//
// Test Coverage:
// - Full container lifecycle management
// - Health endpoint verification using health monitoring
// - Graceful shutdown behavior
// - Resource cleanup verification
func TestSkeletonAppLifecycle(t *testing.T) {
	// Create skeleton application with health endpoint
	app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage()).
		WithHealthEndpoint("/health")

	ctx := context.Background()

	// Test startup
	startTime := time.Now()
	err := app.Start(ctx)
	require.NoError(t, err, "Container should start successfully")
	startupDuration := time.Since(startTime)

	// Verify startup performance (should be under 30 seconds as per success metrics)
	require.Less(t, startupDuration, 30*time.Second, "Container startup should be under 30 seconds")

	// Verify container is running
	require.True(t, app.IsRunning(), "Container should be running")
	require.Equal(t, "/health", app.HealthEndpoint(), "Health endpoint should be configured")

	// Set up health monitoring
	monitor := health.NewHealthMonitor(app)
	monitor.AddCheck(health.NewHTTPHealthCheck("basic-health", ""))
	monitor.AddCheck(health.NewSkeletonSystemHealthCheck())

	// Start health monitoring
	err = monitor.Start(ctx)
	require.NoError(t, err, "Health monitoring should start successfully")

	// Wait for healthy status
	err = monitor.WaitForHealthy(ctx, 30*time.Second)
	require.NoError(t, err, "Application should become healthy within 30 seconds")

	// Verify health status
	status := monitor.Status()
	require.Equal(t, health.StatusHealthy, status.Overall, "Overall health should be healthy")
	require.NotEmpty(t, status.Checks, "Health checks should be present")

	// Stop health monitoring
	err = monitor.Stop()
	require.NoError(t, err, "Health monitoring should stop cleanly")

	// Test graceful shutdown
	shutdownTime := time.Now()
	err = app.Stop(ctx)
	require.NoError(t, err, "Container should stop gracefully")
	shutdownDuration := time.Since(shutdownTime)

	// Verify shutdown performance
	require.Less(t, shutdownDuration, 10*time.Second, "Container shutdown should be under 10 seconds")
	require.False(t, app.IsRunning(), "Container should not be running after stop")
}

// TestSkeletonAppConfiguration tests various configuration options for skeleton
// applications including environment variables and endpoints.
//
// Test Coverage:
// - Environment variable configuration
// - Health and metrics endpoint configuration
// - Shutdown endpoint configuration
// - Configuration persistence across container operations
// - Component verification using the verification package
func TestSkeletonAppConfiguration(t *testing.T) {
	// Create skeleton application with comprehensive configuration
	app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage()).
		WithEnvironment(map[string]string{
			"LOG_LEVEL":    "debug",
			"SERVICE_NAME": "test-skeleton-app",
			"PORT":         "8080",
		}).
		WithHealthEndpoint("/api/health").
		WithShutdownEndpoint("/api/shutdown")

	// Verify configuration is applied
	require.Equal(t, fixtures.GetDefaultTestImage(), app.Image(), "Image should be configured correctly")
	require.Equal(t, "/api/health", app.HealthEndpoint(), "Health endpoint should be configured")
	require.Equal(t, "/api/shutdown", app.ShutdownEndpoint(), "Shutdown endpoint should be configured")

	ctx := context.Background()

	// Start container and verify it respects configuration
	err := app.Start(ctx)
	require.NoError(t, err, "Container should start with custom configuration")

	// Verify container is accessible
	require.True(t, app.IsRunning(), "Container should be running")
	require.NotEmpty(t, app.Host(), "Container should have accessible host")

	// Get port mapping for verification
	port, err := app.Port(8080)
	require.NoError(t, err, "Should be able to get port mapping")
	require.Greater(t, port, 0, "Port should be mapped to a valid external port")

	// Use ComponentVerifier to verify skeleton components
	componentVerifier := verification.NewComponentVerifier(app)

	// Verify that core components are registered and initialized
	err = componentVerifier.VerifySkeletonComponentRegistered(ctx, "system")
	require.NoError(t, err, "System component should be registered")

	err = componentVerifier.VerifySkeletonComponentInitialized(ctx, "system")
	require.NoError(t, err, "System component should be initialized")

	// Set up comprehensive health monitoring
	monitor := health.NewHealthMonitor(app)
	monitor.AddCheck(health.NewHTTPHealthCheck("api-health", "/api/health"))
	monitor.AddCheck(health.NewSkeletonSystemHealthCheck())
	monitor.AddCheck(health.NewSkeletonComponentHealthCheck("system"))

	err = monitor.Start(ctx)
	require.NoError(t, err, "Health monitoring should start")

	// Wait for all checks to pass
	err = monitor.WaitForHealthy(ctx, 30*time.Second)
	require.NoError(t, err, "All health checks should pass")

	// Cleanup
	defer func() {
		monitor.Stop()
		stopErr := app.Stop(ctx)
		require.NoError(t, stopErr, "Container should stop cleanly")
	}()
}

// TestSkeletonAppErrorHandling tests error scenarios and recovery behavior
// for skeleton application containers.
//
// Test Coverage:
// - Invalid image handling
// - Startup failure scenarios
// - Stop operation on non-running containers
// - Resource cleanup on errors
func TestSkeletonAppErrorHandling(t *testing.T) {
	t.Run("InvalidImage", func(t *testing.T) {
		// Test with non-existent image
		app := testkit.NewSkeletonApp("non-existent-image:latest")
		require.NotNil(t, app, "Should create container even with invalid image")

		ctx := context.Background()
		err := app.Start(ctx)
		require.Error(t, err, "Should fail to start with non-existent image")
		require.False(t, app.IsRunning(), "Container should not be running after failed start")
	})

	t.Run("StopNonRunningContainer", func(t *testing.T) {
		// Test stopping a container that was never started
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
		require.False(t, app.IsRunning(), "Container should not be running initially")

		ctx := context.Background()
		err := app.Stop(ctx)
		// Should handle gracefully - either no error or specific "not running" error
		if err != nil {
			require.Contains(t, err.Error(), "not running", "Error should indicate container is not running")
		}
	})

	t.Run("DoubleStart", func(t *testing.T) {
		// Test starting a container twice
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
		ctx := context.Background()

		// First start should succeed
		err := app.Start(ctx)
		require.NoError(t, err, "First start should succeed")
		require.True(t, app.IsRunning(), "Container should be running")

		// Second start should handle gracefully
		err = app.Start(ctx)
		// Should either be no-op or return appropriate error
		if err != nil {
			require.Contains(t, err.Error(), "already", "Error should indicate container is already running")
		}

		// Cleanup
		defer func() {
			stopErr := app.Stop(ctx)
			require.NoError(t, stopErr, "Container should stop cleanly")
		}()
	})
}
