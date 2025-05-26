package docker

import (
	"context"
	"fmt"

	"github.com/fintechain/skeleton-testkit/internal/domain/container"
)

// PortManager manages port mappings and allocations
type PortManager struct {
	allocatedPorts map[string][]int
}

// NewPortManager creates a new port manager
func NewPortManager() *PortManager {
	return &PortManager{
		allocatedPorts: make(map[string][]int),
	}
}

// AllocatePort allocates a port for a container
func (p *PortManager) AllocatePort(containerID string, port int) {
	if _, exists := p.allocatedPorts[containerID]; !exists {
		p.allocatedPorts[containerID] = make([]int, 0)
	}
	p.allocatedPorts[containerID] = append(p.allocatedPorts[containerID], port)
}

// DeallocatePort deallocates a port for a container
func (p *PortManager) DeallocatePort(containerID string, port int) {
	ports, exists := p.allocatedPorts[containerID]
	if !exists {
		return
	}

	for i, allocatedPort := range ports {
		if allocatedPort == port {
			p.allocatedPorts[containerID] = append(ports[:i], ports[i+1:]...)
			break
		}
	}
}

// GetAllocatedPorts returns all allocated ports for a container
func (p *PortManager) GetAllocatedPorts(containerID string) []int {
	ports, exists := p.allocatedPorts[containerID]
	if !exists {
		return []int{}
	}
	return ports
}

// DeallocateAllPorts deallocates all ports for a container
func (p *PortManager) DeallocateAllPorts(containerID string) {
	delete(p.allocatedPorts, containerID)
}

// IsPortAllocated checks if a port is allocated for any container
func (p *PortManager) IsPortAllocated(port int) bool {
	for _, ports := range p.allocatedPorts {
		for _, allocatedPort := range ports {
			if allocatedPort == port {
				return true
			}
		}
	}
	return false
}

// GetPortMapping returns a formatted port mapping string
func GetPortMapping(internal, external int) string {
	if external == 0 {
		return fmt.Sprintf("%d", internal)
	}
	return fmt.Sprintf("%d:%d", external, internal)
}

// ContainerLifecycleManager manages container lifecycle operations
type ContainerLifecycleManager struct {
	containers  map[string]*DockerContainer
	portManager *PortManager
}

// NewContainerLifecycleManager creates a new container lifecycle manager
func NewContainerLifecycleManager() *ContainerLifecycleManager {
	return &ContainerLifecycleManager{
		containers:  make(map[string]*DockerContainer),
		portManager: NewPortManager(),
	}
}

// RegisterContainer registers a container for lifecycle management
func (c *ContainerLifecycleManager) RegisterContainer(dockerContainer *DockerContainer) {
	c.containers[dockerContainer.ID()] = dockerContainer
}

// UnregisterContainer unregisters a container from lifecycle management
func (c *ContainerLifecycleManager) UnregisterContainer(containerID string) {
	delete(c.containers, containerID)
	c.portManager.DeallocateAllPorts(containerID)
}

// StartAll starts all registered containers
func (c *ContainerLifecycleManager) StartAll(ctx context.Context) error {
	for _, dockerContainer := range c.containers {
		if !dockerContainer.IsRunning() {
			if err := dockerContainer.Start(ctx); err != nil {
				return &container.ContainerError{
					Operation: "start_all",
					Container: dockerContainer.ID(),
					Message:   "failed to start container during start all",
					Cause:     err,
				}
			}
		}
	}
	return nil
}

// StopAll stops all registered containers
func (c *ContainerLifecycleManager) StopAll(ctx context.Context) error {
	var lastErr error
	for _, dockerContainer := range c.containers {
		if dockerContainer.IsRunning() {
			if err := dockerContainer.Stop(ctx); err != nil {
				lastErr = &container.ContainerError{
					Operation: "stop_all",
					Container: dockerContainer.ID(),
					Message:   "failed to stop container during stop all",
					Cause:     err,
				}
			}
		}
	}
	return lastErr
}

// GetContainer returns a container by ID
func (c *ContainerLifecycleManager) GetContainer(containerID string) (*DockerContainer, error) {
	dockerContainer, exists := c.containers[containerID]
	if !exists {
		return nil, &container.ContainerError{
			Operation: "get_container",
			Container: containerID,
			Message:   "container not found",
		}
	}
	return dockerContainer, nil
}

// ListContainers returns all registered containers
func (c *ContainerLifecycleManager) ListContainers() map[string]*DockerContainer {
	return c.containers
}

// PortManager returns the port manager
func (c *ContainerLifecycleManager) PortManager() *PortManager {
	return c.portManager
}
