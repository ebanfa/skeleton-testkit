// Package integration provides integration tests for the skeleton-testkit framework.
// This file contains database integration tests that verify skeleton applications
// can properly connect to and interact with database containers.
//
//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/fintechain/skeleton-testkit/internal/domain/container"
	"github.com/fintechain/skeleton-testkit/pkg/health"
	"github.com/fintechain/skeleton-testkit/pkg/testkit"
	"github.com/fintechain/skeleton-testkit/pkg/verification"
	"github.com/fintechain/skeleton-testkit/test/fixtures"
	"github.com/stretchr/testify/require"
)

// TestSkeletonAppWithDatabase tests the integration between skeleton applications
// and database containers. This verifies that the testkit can properly manage
// database dependencies and configure skeleton applications to use them.
//
// Test Coverage:
// - Database container creation and startup
// - Skeleton application with database dependency
// - Database connection string configuration
// - Container dependency ordering
// - Cleanup of dependent containers
// - System and component verification with database
func TestSkeletonAppWithDatabase(t *testing.T) {
	// Create database container for skeleton application
	postgres := testkit.NewPostgresContainer()
	require.NotNil(t, postgres, "NewPostgresContainer should return a valid container")

	// Create skeleton application with database dependency
	app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage()).
		WithDatabase(postgres).
		WithSkeletonConfig(&container.SkeletonConfig{
			ServiceID: "test-app-with-db",
			Storage: container.SkeletonStorageConfig{
				Type: "postgres",
				URL:  postgres.ConnectionString(),
			},
		}).
		WithEnvironment(map[string]string{
			"DB_URL": postgres.ConnectionString(),
		})

	require.NotNil(t, app, "Should create skeleton app with database dependency")

	ctx := context.Background()

	// Start the application (should start database first)
	startTime := time.Now()
	err := app.Start(ctx)
	require.NoError(t, err, "Application should start successfully with database")
	startupDuration := time.Since(startTime)

	// Verify startup performance with dependencies
	require.Less(t, startupDuration, 60*time.Second, "Startup with database should be under 60 seconds")

	// Verify both containers are running
	require.True(t, postgres.IsRunning(), "Database container should be running")
	require.True(t, app.IsRunning(), "Application container should be running")

	// Verify database connection details
	require.NotEmpty(t, postgres.ConnectionString(), "Database should have connection string")
	require.NotEmpty(t, postgres.Host(), "Database should have accessible host")

	dbPort, err := postgres.Port(5432)
	require.NoError(t, err, "Should be able to get database port")
	require.Greater(t, dbPort, 0, "Database port should be mapped")

	// Use SystemVerifier to verify skeleton startup with database
	verifier := verification.NewSystemVerifier(app)
	err = verifier.VerifySkeletonStartup(ctx)
	require.NoError(t, err, "Skeleton system should start up correctly with database")

	// Use ComponentVerifier to verify database-related components
	componentVerifier := verification.NewComponentVerifier(app)
	err = componentVerifier.VerifySkeletonComponentRegistered(ctx, "storage")
	require.NoError(t, err, "Storage component should be registered")

	err = componentVerifier.VerifySkeletonComponentInitialized(ctx, "storage")
	require.NoError(t, err, "Storage component should be initialized")

	// Set up health monitoring for both app and database
	monitor := health.NewHealthMonitor(app)
	monitor.AddCheck(health.NewHTTPHealthCheck("app-health", ""))
	monitor.AddCheck(health.NewSkeletonSystemHealthCheck())
	monitor.AddCheck(health.NewSkeletonComponentHealthCheck("storage"))

	err = monitor.Start(ctx)
	require.NoError(t, err, "Health monitoring should start")

	// Wait for healthy status
	err = monitor.WaitForHealthy(ctx, 45*time.Second)
	require.NoError(t, err, "Application with database should become healthy")

	// Verify health status
	status := monitor.Status()
	require.Equal(t, health.StatusHealthy, status.Overall, "Overall health should be healthy")

	// Cleanup - application should stop first, then database
	defer func() {
		monitor.Stop()

		// Verify graceful shutdown
		shutdownErr := verifier.VerifySkeletonShutdown(ctx)
		require.NoError(t, shutdownErr, "Skeleton should shutdown gracefully")

		// Stop application first
		appStopErr := app.Stop(ctx)
		require.NoError(t, appStopErr, "Application should stop cleanly")
		require.False(t, app.IsRunning(), "Application should not be running after stop")

		// Stop database
		dbStopErr := postgres.Stop(ctx)
		require.NoError(t, dbStopErr, "Database should stop cleanly")
		require.False(t, postgres.IsRunning(), "Database should not be running after stop")
	}()
}

// TestSkeletonAppWithCustomDatabaseConfig tests skeleton applications with
// custom database configurations including credentials and database names.
//
// Test Coverage:
// - Custom database configuration
// - Database credentials handling
// - Custom database name configuration
// - Environment variable propagation
// - Configuration validation
// - Component verification with custom config
func TestSkeletonAppWithCustomDatabaseConfig(t *testing.T) {
	// Create database with custom configuration
	dbConfig := &testkit.PostgresConfig{
		Database: "skeleton_test_db",
		Username: "skeleton_user",
		Password: "skeleton_pass",
	}
	postgres := testkit.NewPostgresContainerWithConfig(dbConfig)
	require.NotNil(t, postgres, "Should create database with custom config")

	// Create skeleton application with custom database configuration
	skeletonConfig := &container.SkeletonConfig{
		ServiceID: "custom-db-app",
		Storage: container.SkeletonStorageConfig{
			Type: "postgres",
			URL:  postgres.ConnectionString(),
		},
		Plugins: []container.SkeletonPluginConfig{
			{
				Name:    "database-plugin",
				Version: "1.0.0",
				Config: map[string]interface{}{
					"database": "skeleton_test_db",
					"username": "skeleton_user",
				},
			},
		},
	}

	app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage()).
		WithDatabase(postgres).
		WithSkeletonConfig(skeletonConfig).
		WithEnvironment(map[string]string{
			"DB_NAME":     "skeleton_test_db",
			"DB_USER":     "skeleton_user",
			"DB_PASSWORD": "skeleton_pass",
			"DB_URL":      postgres.ConnectionString(),
		})

	ctx := context.Background()

	// Start and verify configuration
	err := app.Start(ctx)
	require.NoError(t, err, "Application should start with custom database config")

	// Verify containers are running
	require.True(t, postgres.IsRunning(), "Database should be running")
	require.True(t, app.IsRunning(), "Application should be running")

	// Verify connection string contains custom configuration
	connectionString := postgres.ConnectionString()
	require.Contains(t, connectionString, "skeleton_test_db", "Connection string should contain custom database name")
	require.Contains(t, connectionString, "skeleton_user", "Connection string should contain custom username")

	// Use ComponentVerifier to verify plugin components
	componentVerifier := verification.NewComponentVerifier(app)
	err = componentVerifier.VerifySkeletonComponentRegistered(ctx, "database-plugin")
	require.NoError(t, err, "Database plugin component should be registered")

	// Verify component metadata matches configuration
	err = componentVerifier.VerifySkeletonComponentMetadata(ctx, "database-plugin", map[string]interface{}{
		"name":    "database-plugin",
		"version": "1.0.0",
	})
	require.NoError(t, err, "Database plugin metadata should match configuration")

	// Set up health monitoring with custom checks
	monitor := health.NewHealthMonitor(app)
	monitor.AddCheck(health.NewHTTPHealthCheck("custom-health", ""))
	monitor.AddCheck(health.NewSkeletonComponentHealthCheck("database-plugin"))

	err = monitor.Start(ctx)
	require.NoError(t, err, "Health monitoring should start")

	err = monitor.WaitForHealthy(ctx, 30*time.Second)
	require.NoError(t, err, "Custom configured app should become healthy")

	// Cleanup
	defer func() {
		monitor.Stop()
		appStopErr := app.Stop(ctx)
		require.NoError(t, appStopErr, "Application should stop cleanly")

		dbStopErr := postgres.Stop(ctx)
		require.NoError(t, dbStopErr, "Database should stop cleanly")
	}()
}

// TestSkeletonAppMultipleDatabases tests skeleton applications with multiple
// database dependencies to verify complex dependency management.
//
// Test Coverage:
// - Multiple database containers
// - Complex dependency ordering
// - Multiple connection string management
// - Resource cleanup with multiple dependencies
// - Performance with multiple containers
func TestSkeletonAppMultipleDatabases(t *testing.T) {
	// Create primary database
	primaryDB := testkit.NewPostgresContainerWithConfig(&testkit.PostgresConfig{
		Database: "primary_db",
		Username: "primary_user",
		Password: "primary_pass",
	})

	// Create secondary database (could be different type in real scenario)
	secondaryDB := testkit.NewPostgresContainerWithConfig(&testkit.PostgresConfig{
		Database: "secondary_db",
		Username: "secondary_user",
		Password: "secondary_pass",
	})

	// Create skeleton application with multiple database dependencies
	app := testkit.NewSkeletonApp(fixtures.GetDefaultTestImage()).
		WithDatabase(primaryDB).
		WithDatabase(secondaryDB).
		WithSkeletonConfig(&container.SkeletonConfig{
			ServiceID: "multi-db-app",
			Storage: container.SkeletonStorageConfig{
				Type: "postgres",
				URL:  primaryDB.ConnectionString(),
			},
		}).
		WithEnvironment(map[string]string{
			"PRIMARY_DB_URL":   primaryDB.ConnectionString(),
			"SECONDARY_DB_URL": secondaryDB.ConnectionString(),
		})

	ctx := context.Background()

	// Start application with multiple dependencies
	startTime := time.Now()
	err := app.Start(ctx)
	require.NoError(t, err, "Application should start with multiple databases")
	startupDuration := time.Since(startTime)

	// Verify startup performance with multiple dependencies
	require.Less(t, startupDuration, 90*time.Second, "Startup with multiple databases should be under 90 seconds")

	// Verify all containers are running
	require.True(t, primaryDB.IsRunning(), "Primary database should be running")
	require.True(t, secondaryDB.IsRunning(), "Secondary database should be running")
	require.True(t, app.IsRunning(), "Application should be running")

	// Verify both databases are accessible
	require.NotEmpty(t, primaryDB.ConnectionString(), "Primary database should have connection string")
	require.NotEmpty(t, secondaryDB.ConnectionString(), "Secondary database should have connection string")

	// Verify different connection strings
	require.NotEqual(t, primaryDB.ConnectionString(), secondaryDB.ConnectionString(),
		"Databases should have different connection strings")

	// Cleanup in reverse order
	defer func() {
		// Stop application first
		appStopErr := app.Stop(ctx)
		require.NoError(t, appStopErr, "Application should stop cleanly")

		// Stop databases
		primaryStopErr := primaryDB.Stop(ctx)
		require.NoError(t, primaryStopErr, "Primary database should stop cleanly")

		secondaryStopErr := secondaryDB.Stop(ctx)
		require.NoError(t, secondaryStopErr, "Secondary database should stop cleanly")
	}()
}

// TestDatabaseContainerLifecycle tests the lifecycle management of database
// containers independently of skeleton applications.
//
// Test Coverage:
// - Database container creation
// - Database startup and readiness
// - Database connection verification
// - Database shutdown and cleanup
// - Database container configuration
func TestDatabaseContainerLifecycle(t *testing.T) {
	// Test default PostgreSQL container
	postgres := testkit.NewPostgresContainer()
	require.NotNil(t, postgres, "Should create default PostgreSQL container")

	// Verify initial state
	require.False(t, postgres.IsRunning(), "Database should not be running initially")
	require.Equal(t, "postgres:15", postgres.Image(), "Should use default PostgreSQL image")

	ctx := context.Background()

	// Start database
	startTime := time.Now()
	err := postgres.Start(ctx)
	require.NoError(t, err, "Database should start successfully")
	startupDuration := time.Since(startTime)

	// Verify startup performance
	require.Less(t, startupDuration, 30*time.Second, "Database startup should be under 30 seconds")

	// Verify running state
	require.True(t, postgres.IsRunning(), "Database should be running after start")
	require.NotEmpty(t, postgres.ID(), "Database should have valid ID")
	require.NotEmpty(t, postgres.Host(), "Database should have accessible host")

	// Verify connection details
	connectionString := postgres.ConnectionString()
	require.NotEmpty(t, connectionString, "Database should have connection string")
	require.Contains(t, connectionString, "postgres://", "Connection string should be PostgreSQL format")

	// Verify port mapping
	port, err := postgres.Port(5432)
	require.NoError(t, err, "Should be able to get database port")
	require.Greater(t, port, 0, "Database port should be mapped to valid external port")

	// Test graceful shutdown
	shutdownTime := time.Now()
	err = postgres.Stop(ctx)
	require.NoError(t, err, "Database should stop gracefully")
	shutdownDuration := time.Since(shutdownTime)

	// Verify shutdown performance
	require.Less(t, shutdownDuration, 10*time.Second, "Database shutdown should be under 10 seconds")
	require.False(t, postgres.IsRunning(), "Database should not be running after stop")
}
