package testcontainers

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/fintechain/skeleton-testkit/internal/domain/container"
	"github.com/fintechain/skeleton-testkit/internal/infrastructure/docker"
)

// RedisContainer wraps a Redis container for testing
type RedisContainer struct {
	*docker.DockerContainer
	password string
}

// RedisConfig holds Redis container configuration
type RedisConfig struct {
	Password string
	Image    string
}

// NewRedisContainer creates a new Redis container with default configuration
func NewRedisContainer() *RedisContainer {
	return NewRedisContainerWithConfig(&RedisConfig{
		Password: "",
		Image:    "redis:7",
	})
}

// NewRedisContainerWithConfig creates a new Redis container with custom configuration
func NewRedisContainerWithConfig(config *RedisConfig) *RedisContainer {
	env := make(map[string]string)
	if config.Password != "" {
		env["REDIS_PASSWORD"] = config.Password
	}

	containerConfig := &docker.ContainerConfig{
		ID:          fmt.Sprintf("redis-%d", time.Now().UnixNano()),
		Name:        "redis-test",
		Image:       config.Image,
		Environment: env,
		Ports: []container.PortMapping{
			{Internal: 6379, External: 0}, // Random external port
		},
	}

	return &RedisContainer{
		DockerContainer: docker.NewDockerContainer(containerConfig),
		password:        config.Password,
	}
}

// Start starts the Redis container
func (r *RedisContainer) Start(ctx context.Context) error {
	if err := r.createContainer(ctx); err != nil {
		return err
	}
	return r.DockerContainer.Start(ctx)
}

// createContainer creates the underlying testcontainer
func (r *RedisContainer) createContainer(ctx context.Context) error {
	config := r.Config()

	cmd := []string{"redis-server"}
	if r.password != "" {
		cmd = append(cmd, "--requirepass", r.password)
	}

	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		Name:         config.Name,
		Env:          config.Environment,
		ExposedPorts: []string{"6379/tcp"},
		Cmd:          cmd,
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("6379/tcp"),
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(30*time.Second),
		),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})
	if err != nil {
		return &container.ContainerError{
			Operation: "create",
			Container: r.ID(),
			Message:   "failed to create redis container",
			Cause:     err,
		}
	}

	r.SetContainer(c)
	return nil
}

// ConnectionString returns the Redis connection string
func (r *RedisContainer) ConnectionString() string {
	host := r.Host()
	if host == "" {
		return ""
	}

	port, err := r.Port(6379)
	if err != nil {
		return ""
	}

	if r.password != "" {
		return fmt.Sprintf("redis://:%s@%s:%d", r.password, host, port)
	}
	return fmt.Sprintf("redis://%s:%d", host, port)
}

// Password returns the password
func (r *RedisContainer) Password() string {
	return r.password
}
