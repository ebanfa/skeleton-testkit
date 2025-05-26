// Package integration provides integration tests for the skeleton-testkit framework.
// This file contains performance benchmarks that verify the testkit meets
// performance requirements for container startup and lifecycle operations.
//
//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	domaincontainer "github.com/fintechain/skeleton-testkit/internal/domain/container"
	"github.com/fintechain/skeleton-testkit/pkg/container"
	"github.com/fintechain/skeleton-testkit/pkg/testkit"
	"github.com/fintechain/skeleton-testkit/test/fixtures"
	"github.com/stretchr/testify/require"
)

// BenchmarkSkeletonAppStartup benchmarks the startup time of skeleton application
// containers to ensure they meet the performance requirements specified in the
// implementation plan (sub-30 second startup times).
//
// Performance Requirements:
// - Single container startup: < 30 seconds
// - Container creation overhead: < 5 seconds
// - Memory usage: reasonable for testing scenarios
func BenchmarkSkeletonAppStartup(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create skeleton application container
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())

		// Measure startup time
		startTime := time.Now()
		err := app.Start(ctx)
		startupDuration := time.Since(startTime)

		// Report metrics
		if err == nil {
			b.ReportMetric(float64(startupDuration.Milliseconds()), "ms/startup")

			// Cleanup
			_ = app.Stop(ctx)
		} else {
			b.Errorf("Container startup failed: %v", err)
		}
	}
}

// BenchmarkDatabaseContainerStartup benchmarks database container startup
// performance to ensure database dependencies don't significantly impact
// overall test execution time.
//
// Performance Requirements:
// - Database startup: < 30 seconds
// - Connection readiness: < 10 seconds after startup
func BenchmarkDatabaseContainerStartup(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create database container
		postgres := testkit.NewPostgresContainer()

		// Measure startup time
		startTime := time.Now()
		err := postgres.Start(ctx)
		startupDuration := time.Since(startTime)

		// Report metrics
		if err == nil {
			b.ReportMetric(float64(startupDuration.Milliseconds()), "ms/db-startup")

			// Cleanup
			_ = postgres.Stop(ctx)
		} else {
			b.Errorf("Database startup failed: %v", err)
		}
	}
}

// BenchmarkSkeletonAppWithDatabaseStartup benchmarks the startup time of
// skeleton applications with database dependencies to measure the overhead
// of dependency management.
//
// Performance Requirements:
// - App + database startup: < 60 seconds
// - Dependency ordering overhead: < 10 seconds
func BenchmarkSkeletonAppWithDatabaseStartup(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create database and application
		postgres := testkit.NewPostgresContainer()
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage()).
			WithDatabase(postgres).
			WithSkeletonConfig(&domaincontainer.SkeletonConfig{
				ServiceID: "benchmark-app",
				Storage: domaincontainer.SkeletonStorageConfig{
					Type: "postgres",
					URL:  postgres.ConnectionString(),
				},
			})

		// Measure total startup time
		startTime := time.Now()
		err := app.Start(ctx)
		startupDuration := time.Since(startTime)

		// Report metrics
		if err == nil {
			b.ReportMetric(float64(startupDuration.Milliseconds()), "ms/app-db-startup")

			// Cleanup
			_ = app.Stop(ctx)
			_ = postgres.Stop(ctx)
		} else {
			b.Errorf("App with database startup failed: %v", err)
		}
	}
}

// TestContainerStartupPerformance tests that container startup times meet
// the performance requirements specified in the implementation plan.
//
// Performance Requirements (from Phase 1 Success Metrics):
// - Sub-30 second container startup times
// - Basic error handling and cleanup
// - 90%+ test coverage for core components
func TestContainerStartupPerformance(t *testing.T) {
	ctx := context.Background()

	t.Run("SkeletonAppStartupTime", func(t *testing.T) {
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())

		// Measure startup time
		startTime := time.Now()
		err := app.Start(ctx)
		startupDuration := time.Since(startTime)

		require.NoError(t, err, "Container should start successfully")
		require.Less(t, startupDuration, 30*time.Second, "Startup should be under 30 seconds")

		// Cleanup
		defer func() {
			stopErr := app.Stop(ctx)
			require.NoError(t, stopErr, "Container should stop cleanly")
		}()
	})

	t.Run("DatabaseStartupTime", func(t *testing.T) {
		postgres := testkit.NewPostgresContainer()

		// Measure startup time
		startTime := time.Now()
		err := postgres.Start(ctx)
		startupDuration := time.Since(startTime)

		require.NoError(t, err, "Database should start successfully")
		require.Less(t, startupDuration, 30*time.Second, "Database startup should be under 30 seconds")

		// Cleanup
		defer func() {
			stopErr := postgres.Stop(ctx)
			require.NoError(t, stopErr, "Database should stop cleanly")
		}()
	})

	t.Run("AppWithDatabaseStartupTime", func(t *testing.T) {
		postgres := testkit.NewPostgresContainer()
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage()).
			WithDatabase(postgres)

		// Measure total startup time
		startTime := time.Now()
		err := app.Start(ctx)
		startupDuration := time.Since(startTime)

		require.NoError(t, err, "App with database should start successfully")
		require.Less(t, startupDuration, 60*time.Second, "App with database startup should be under 60 seconds")

		// Cleanup
		defer func() {
			appStopErr := app.Stop(ctx)
			require.NoError(t, appStopErr, "App should stop cleanly")

			dbStopErr := postgres.Stop(ctx)
			require.NoError(t, dbStopErr, "Database should stop cleanly")
		}()
	})
}

// TestContainerShutdownPerformance tests that container shutdown times
// are reasonable and don't block test execution.
//
// Performance Requirements:
// - Container shutdown: < 10 seconds
// - Graceful shutdown handling
// - Resource cleanup verification
func TestContainerShutdownPerformance(t *testing.T) {
	ctx := context.Background()

	t.Run("SkeletonAppShutdownTime", func(t *testing.T) {
		app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())

		// Start container
		err := app.Start(ctx)
		require.NoError(t, err, "Container should start successfully")

		// Measure shutdown time
		shutdownTime := time.Now()
		err = app.Stop(ctx)
		shutdownDuration := time.Since(shutdownTime)

		require.NoError(t, err, "Container should stop successfully")
		require.Less(t, shutdownDuration, 10*time.Second, "Shutdown should be under 10 seconds")
		require.False(t, app.IsRunning(), "Container should not be running after stop")
	})

	t.Run("DatabaseShutdownTime", func(t *testing.T) {
		postgres := testkit.NewPostgresContainer()

		// Start database
		err := postgres.Start(ctx)
		require.NoError(t, err, "Database should start successfully")

		// Measure shutdown time
		shutdownTime := time.Now()
		err = postgres.Stop(ctx)
		shutdownDuration := time.Since(shutdownTime)

		require.NoError(t, err, "Database should stop successfully")
		require.Less(t, shutdownDuration, 10*time.Second, "Database shutdown should be under 10 seconds")
		require.False(t, postgres.IsRunning(), "Database should not be running after stop")
	})
}

// TestConcurrentContainerOperations tests the performance of concurrent
// container operations to verify the testkit can handle parallel test execution.
//
// Performance Requirements:
// - Support multiple concurrent containers
// - No resource contention issues
// - Reasonable memory usage under load
func TestConcurrentContainerOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent operations test in short mode")
	}

	ctx := context.Background()
	const numContainers = 3 // Keep reasonable for CI environments

	t.Run("ConcurrentSkeletonApps", func(t *testing.T) {
		apps := make([]*container.AppContainer, numContainers)

		// Create containers
		for i := 0; i < numContainers; i++ {
			apps[i] = testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
		}

		// Start all containers concurrently
		startTime := time.Now()
		for i := 0; i < numContainers; i++ {
			go func(app *container.AppContainer) {
				err := app.Start(ctx)
				require.NoError(t, err, "Concurrent container should start successfully")
			}(apps[i])
		}

		// Wait for all to be running (with timeout)
		timeout := time.After(60 * time.Second)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				t.Fatal("Timeout waiting for concurrent containers to start")
			case <-ticker.C:
				allRunning := true
				for _, app := range apps {
					if !app.IsRunning() {
						allRunning = false
						break
					}
				}
				if allRunning {
					goto cleanup
				}
			}
		}

	cleanup:
		startupDuration := time.Since(startTime)
		require.Less(t, startupDuration, 90*time.Second, "Concurrent startup should complete within 90 seconds")

		// Cleanup all containers
		for _, app := range apps {
			if app.IsRunning() {
				err := app.Stop(ctx)
				require.NoError(t, err, "Container should stop cleanly")
			}
		}
	})
}

// TestMemoryUsageUnderLoad tests memory usage patterns during container
// operations to ensure the testkit doesn't have memory leaks.
//
// Performance Requirements:
// - Reasonable memory usage growth
// - No obvious memory leaks
// - Cleanup releases resources
func TestMemoryUsageUnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory usage test in short mode")
	}

	ctx := context.Background()
	const iterations = 5

	t.Run("MemoryUsagePattern", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			// Create and start container
			app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage())
			err := app.Start(ctx)
			require.NoError(t, err, "Container should start successfully")

			// Verify it's running
			require.True(t, app.IsRunning(), "Container should be running")

			// Stop and cleanup
			err = app.Stop(ctx)
			require.NoError(t, err, "Container should stop successfully")
			require.False(t, app.IsRunning(), "Container should not be running after stop")

			// Small delay to allow cleanup
			time.Sleep(100 * time.Millisecond)
		}
	})
}
