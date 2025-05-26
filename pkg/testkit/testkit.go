// Package testkit provides the main entry point for the skeleton-testkit framework.
// It offers a clean, fluent API for creating and managing containerized testing environments
// for skeleton-based applications.
package testkit

import (
	"math/rand"
	"time"

	domaincontainer "github.com/fintechain/skeleton-testkit/internal/domain/container"
	"github.com/fintechain/skeleton-testkit/internal/infrastructure/docker"
	"github.com/fintechain/skeleton-testkit/internal/infrastructure/testcontainers"
	"github.com/fintechain/skeleton-testkit/pkg/container"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewSkeletonApp creates a new container for testing a skeleton-based application
func NewSkeletonApp(imageName string) *container.AppContainer {
	config := &domaincontainer.AppConfig{
		ImageName: imageName,
	}
	return NewSkeletonAppWithConfig(config)
}

// NewSkeletonAppWithConfig creates an app container with custom configuration
func NewSkeletonAppWithConfig(config *domaincontainer.AppConfig) *container.AppContainer {
	// Convert domain config to infrastructure config
	containerConfig := &docker.ContainerConfig{
		ID:          generateContainerID(),
		Name:        "skeleton-app",
		Image:       config.ImageName,
		Environment: config.Environment,
		Ports:       config.Ports,
	}

	// Create testcontainer implementation
	impl := testcontainers.NewTestcontainerAppContainer(containerConfig, nil)

	// Return public API wrapper
	return container.NewAppContainer(impl)
}

// NewPostgresContainer creates a new PostgreSQL container for testing
func NewPostgresContainer() *container.PostgresContainer {
	impl := testcontainers.NewPostgresContainer()
	return container.NewPostgresContainer(impl)
}

// NewPostgresContainerWithConfig creates a PostgreSQL container with custom configuration
func NewPostgresContainerWithConfig(config *PostgresConfig) *container.PostgresContainer {
	postgresConfig := &testcontainers.PostgresConfig{
		Database: config.Database,
		Username: config.Username,
		Password: config.Password,
	}
	impl := testcontainers.NewPostgresContainerWithConfig(postgresConfig)
	return container.NewPostgresContainer(impl)
}

// NewRedisContainer creates a new Redis container for testing
func NewRedisContainer() *container.RedisContainer {
	impl := testcontainers.NewRedisContainer()
	return container.NewRedisContainer(impl)
}

// NewRedisContainerWithConfig creates a Redis container with custom configuration
func NewRedisContainerWithConfig(config *RedisConfig) *container.RedisContainer {
	redisConfig := &testcontainers.RedisConfig{
		Password: config.Password,
	}
	impl := testcontainers.NewRedisContainerWithConfig(redisConfig)
	return container.NewRedisContainer(impl)
}

// generateContainerID generates a unique container ID
func generateContainerID() string {
	// Simple implementation - in practice this would be more sophisticated
	return "container-" + randomString(8)
}

// randomString generates a random string of the given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// PostgresConfig holds configuration for PostgreSQL containers
type PostgresConfig struct {
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Image    string `json:"image"`
}

// RedisConfig holds configuration for Redis containers
type RedisConfig struct {
	Password string `json:"password"`
	Image    string `json:"image"`
}
