package container

// AppConfig holds configuration for application containers
type AppConfig struct {
	ImageName        string            `json:"imageName"`
	HealthEndpoint   string            `json:"healthEndpoint"`
	MetricsEndpoint  string            `json:"metricsEndpoint"`
	ShutdownEndpoint string            `json:"shutdownEndpoint"`
	Environment      map[string]string `json:"environment"`
	Ports            []PortMapping     `json:"ports"`
	Volumes          []VolumeMapping   `json:"volumes"`
}

// PortMapping defines how container ports are mapped
type PortMapping struct {
	Internal int `json:"internal"`
	External int `json:"external"` // 0 means random port
}

// VolumeMapping defines how container volumes are mapped
type VolumeMapping struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// SkeletonConfig holds skeleton-specific configuration
type SkeletonConfig struct {
	ServiceID string                 `json:"serviceId"`
	Plugins   []SkeletonPluginConfig `json:"plugins"`
	Storage   SkeletonStorageConfig  `json:"storage"`
}

// SkeletonPluginConfig defines a skeleton plugin configuration
type SkeletonPluginConfig struct {
	Name    string                 `json:"name"`
	Version string                 `json:"version"`
	Config  map[string]interface{} `json:"config"`
}

// SkeletonStorageConfig defines storage configuration for skeleton applications
type SkeletonStorageConfig struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}
