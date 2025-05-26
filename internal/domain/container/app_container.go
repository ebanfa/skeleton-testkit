package container

// AppContainer represents a containerized skeleton-based application
type AppContainer struct {
	config         *AppConfig
	dependencies   []Container
	skeletonConfig *SkeletonConfig
}

// NewAppContainer creates a new application container with the given configuration
func NewAppContainer(config *AppConfig) *AppContainer {
	return &AppContainer{
		config:       config,
		dependencies: make([]Container, 0),
	}
}

// Config returns the application configuration
func (a *AppContainer) Config() *AppConfig {
	return a.config
}

// Dependencies returns the list of container dependencies
func (a *AppContainer) Dependencies() []Container {
	return a.dependencies
}

// SkeletonConfig returns the skeleton-specific configuration
func (a *AppContainer) SkeletonConfig() *SkeletonConfig {
	return a.skeletonConfig
}

// WithDependency adds a container dependency
func (a *AppContainer) WithDependency(container Container) *AppContainer {
	a.dependencies = append(a.dependencies, container)
	return a
}

// WithSkeletonConfig sets the skeleton-specific configuration
func (a *AppContainer) WithSkeletonConfig(config *SkeletonConfig) *AppContainer {
	a.skeletonConfig = config
	return a
}
