package health

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// HTTPHealthCheck performs HTTP-based health checks
type HTTPHealthCheck struct {
	name     string
	endpoint string
	interval time.Duration
	timeout  time.Duration
	client   *http.Client
}

// NewHTTPHealthCheck creates a new HTTP health check
func NewHTTPHealthCheck(name, endpoint string) *HTTPHealthCheck {
	return &HTTPHealthCheck{
		name:     name,
		endpoint: endpoint,
		interval: 30 * time.Second,
		timeout:  10 * time.Second,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Name returns the name of the health check
func (h *HTTPHealthCheck) Name() string {
	return h.name
}

// Check performs the health check
func (h *HTTPHealthCheck) Check(ctx context.Context, target HealthTarget) error {
	url := target.HealthEndpoint()
	if h.endpoint != "" {
		url = h.endpoint
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// Interval returns the check interval
func (h *HTTPHealthCheck) Interval() time.Duration {
	return h.interval
}

// Timeout returns the check timeout
func (h *HTTPHealthCheck) Timeout() time.Duration {
	return h.timeout
}

// WithInterval sets the check interval
func (h *HTTPHealthCheck) WithInterval(interval time.Duration) *HTTPHealthCheck {
	h.interval = interval
	return h
}

// WithTimeout sets the check timeout
func (h *HTTPHealthCheck) WithTimeout(timeout time.Duration) *HTTPHealthCheck {
	h.timeout = timeout
	h.client.Timeout = timeout
	return h
}

// SkeletonSystemHealthCheck checks the skeleton system service
type SkeletonSystemHealthCheck struct {
	name     string
	interval time.Duration
	timeout  time.Duration
	client   *http.Client
}

// NewSkeletonSystemHealthCheck creates a new skeleton system health check
func NewSkeletonSystemHealthCheck() *SkeletonSystemHealthCheck {
	return &SkeletonSystemHealthCheck{
		name:     "skeleton-system",
		interval: 30 * time.Second,
		timeout:  10 * time.Second,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Name returns the name of the health check
func (s *SkeletonSystemHealthCheck) Name() string {
	return s.name
}

// Check performs the skeleton system health check
func (s *SkeletonSystemHealthCheck) Check(ctx context.Context, target HealthTarget) error {
	// Check the skeleton system service endpoint
	url := fmt.Sprintf("%s/skeleton/system", target.HealthEndpoint())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create skeleton system request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("skeleton system check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("skeleton system check failed with status %d", resp.StatusCode)
	}

	return nil
}

// Interval returns the check interval
func (s *SkeletonSystemHealthCheck) Interval() time.Duration {
	return s.interval
}

// Timeout returns the check timeout
func (s *SkeletonSystemHealthCheck) Timeout() time.Duration {
	return s.timeout
}

// SkeletonComponentHealthCheck checks skeleton component health
type SkeletonComponentHealthCheck struct {
	name        string
	componentID string
	interval    time.Duration
	timeout     time.Duration
	client      *http.Client
}

// NewSkeletonComponentHealthCheck creates a new skeleton component health check
func NewSkeletonComponentHealthCheck(componentID string) *SkeletonComponentHealthCheck {
	return &SkeletonComponentHealthCheck{
		name:        fmt.Sprintf("skeleton-component-%s", componentID),
		componentID: componentID,
		interval:    30 * time.Second,
		timeout:     10 * time.Second,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Name returns the name of the health check
func (s *SkeletonComponentHealthCheck) Name() string {
	return s.name
}

// Check performs the skeleton component health check
func (s *SkeletonComponentHealthCheck) Check(ctx context.Context, target HealthTarget) error {
	// Check the skeleton component status endpoint
	url := fmt.Sprintf("%s/skeleton/components/%s/status", target.HealthEndpoint(), s.componentID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create component status request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("component status check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("component status check failed with status %d", resp.StatusCode)
	}

	return nil
}

// Interval returns the check interval
func (s *SkeletonComponentHealthCheck) Interval() time.Duration {
	return s.interval
}

// Timeout returns the check timeout
func (s *SkeletonComponentHealthCheck) Timeout() time.Duration {
	return s.timeout
}
