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

// PostgresContainer wraps a PostgreSQL container for testing
type PostgresContainer struct {
	*docker.DockerContainer
	database string
	username string
	password string
}

// PostgresConfig holds PostgreSQL container configuration
type PostgresConfig struct {
	Database string
	Username string
	Password string
	Image    string
}

// NewPostgresContainer creates a new PostgreSQL container with default configuration
func NewPostgresContainer() *PostgresContainer {
	return NewPostgresContainerWithConfig(&PostgresConfig{
		Database: "testdb",
		Username: "testuser",
		Password: "testpass",
		Image:    "postgres:15",
	})
}

// NewPostgresContainerWithConfig creates a new PostgreSQL container with custom configuration
func NewPostgresContainerWithConfig(config *PostgresConfig) *PostgresContainer {
	containerConfig := &docker.ContainerConfig{
		ID:    fmt.Sprintf("postgres-%d", time.Now().UnixNano()),
		Name:  "postgres-test",
		Image: config.Image,
		Environment: map[string]string{
			"POSTGRES_DB":       config.Database,
			"POSTGRES_USER":     config.Username,
			"POSTGRES_PASSWORD": config.Password,
		},
		Ports: []container.PortMapping{
			{Internal: 5432, External: 0}, // Random external port
		},
	}

	return &PostgresContainer{
		DockerContainer: docker.NewDockerContainer(containerConfig),
		database:        config.Database,
		username:        config.Username,
		password:        config.Password,
	}
}

// Start starts the PostgreSQL container
func (p *PostgresContainer) Start(ctx context.Context) error {
	if err := p.createContainer(ctx); err != nil {
		return err
	}
	return p.DockerContainer.Start(ctx)
}

// createContainer creates the underlying testcontainer
func (p *PostgresContainer) createContainer(ctx context.Context) error {
	config := p.Config()

	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		Name:         config.Name,
		Env:          config.Environment,
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
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
			Container: p.ID(),
			Message:   "failed to create postgres container",
			Cause:     err,
		}
	}

	p.SetContainer(c)
	return nil
}

// ConnectionString returns the PostgreSQL connection string
func (p *PostgresContainer) ConnectionString() string {
	host := p.Host()
	if host == "" {
		return ""
	}

	port, err := p.Port(5432)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.username, p.password, host, port, p.database)
}

// Database returns the database name
func (p *PostgresContainer) Database() string {
	return p.database
}

// Username returns the username
func (p *PostgresContainer) Username() string {
	return p.username
}

// Password returns the password
func (p *PostgresContainer) Password() string {
	return p.password
}
