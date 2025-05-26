// Package health provides health monitoring capabilities for skeleton applications.
// It implements the health monitoring framework specified in the skeleton-testkit specification.
package health

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// HealthTarget represents what is being health checked
type HealthTarget interface {
	HealthEndpoint() string
	ConnectionString() string
}

// HealthCheck represents a health verification
type HealthCheck interface {
	Name() string
	Check(ctx context.Context, target HealthTarget) error
	Interval() time.Duration
	Timeout() time.Duration
}

// HealthStatus represents current health state
type HealthStatus struct {
	Overall   Status                 `json:"overall"`
	Checks    map[string]CheckResult `json:"checks"`
	Timestamp time.Time              `json:"timestamp"`
}

// Status represents health status
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusUnknown   Status = "unknown"
)

// CheckResult represents the result of a health check
type CheckResult struct {
	Status    Status        `json:"status"`
	Duration  time.Duration `json:"duration"`
	Error     string        `json:"error,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// HealthMonitor provides health monitoring capabilities
type HealthMonitor struct {
	target   HealthTarget
	checks   []HealthCheck
	interval time.Duration
	status   HealthStatus
	mutex    sync.RWMutex
	stopCh   chan struct{}
	running  bool
}

// NewHealthMonitor creates a new HealthMonitor for the given target
func NewHealthMonitor(target HealthTarget) *HealthMonitor {
	return &HealthMonitor{
		target:   target,
		checks:   make([]HealthCheck, 0),
		interval: 30 * time.Second,
		status: HealthStatus{
			Overall:   StatusUnknown,
			Checks:    make(map[string]CheckResult),
			Timestamp: time.Now(),
		},
		stopCh: make(chan struct{}),
	}
}

// AddCheck adds a health check to the monitor
func (h *HealthMonitor) AddCheck(check HealthCheck) *HealthMonitor {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.checks = append(h.checks, check)
	return h
}

// Start starts the health monitoring
func (h *HealthMonitor) Start(ctx context.Context) error {
	h.mutex.Lock()
	if h.running {
		h.mutex.Unlock()
		return fmt.Errorf("health monitor is already running")
	}
	h.running = true
	h.mutex.Unlock()

	// Run initial health check
	h.runHealthChecks(ctx)

	// Start monitoring loop
	go h.monitoringLoop(ctx)

	return nil
}

// Stop stops the health monitoring
func (h *HealthMonitor) Stop() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if !h.running {
		return nil
	}

	close(h.stopCh)
	h.running = false
	return nil
}

// Status returns the current health status
func (h *HealthMonitor) Status() HealthStatus {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.status
}

// WaitForHealthy waits for the target to become healthy within the timeout
func (h *HealthMonitor) WaitForHealthy(ctx context.Context, timeout time.Duration) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout waiting for healthy status")
		case <-ticker.C:
			h.runHealthChecks(ctx)
			status := h.Status()
			if status.Overall == StatusHealthy {
				return nil
			}
		}
	}
}

// monitoringLoop runs the health monitoring loop
func (h *HealthMonitor) monitoringLoop(ctx context.Context) {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-h.stopCh:
			return
		case <-ticker.C:
			h.runHealthChecks(ctx)
		}
	}
}

// runHealthChecks executes all health checks and updates status
func (h *HealthMonitor) runHealthChecks(ctx context.Context) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	results := make(map[string]CheckResult)
	overallHealthy := true

	for _, check := range h.checks {
		result := h.executeCheck(ctx, check)
		results[check.Name()] = result

		if result.Status != StatusHealthy {
			overallHealthy = false
		}
	}

	overall := StatusHealthy
	if !overallHealthy {
		overall = StatusUnhealthy
	}

	h.status = HealthStatus{
		Overall:   overall,
		Checks:    results,
		Timestamp: time.Now(),
	}
}

// executeCheck executes a single health check
func (h *HealthMonitor) executeCheck(ctx context.Context, check HealthCheck) CheckResult {
	start := time.Now()

	checkCtx, cancel := context.WithTimeout(ctx, check.Timeout())
	defer cancel()

	err := check.Check(checkCtx, h.target)
	duration := time.Since(start)

	result := CheckResult{
		Duration:  duration,
		Timestamp: time.Now(),
	}

	if err != nil {
		result.Status = StatusUnhealthy
		result.Error = err.Error()
	} else {
		result.Status = StatusHealthy
	}

	return result
}
