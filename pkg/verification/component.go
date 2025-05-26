package verification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fintechain/skeleton-testkit/pkg/container"
)

// ComponentVerifier verifies skeleton component behavior
type ComponentVerifier struct {
	app *container.AppContainer
}

// NewComponentVerifier creates a new ComponentVerifier for the given application container
func NewComponentVerifier(app *container.AppContainer) *ComponentVerifier {
	return &ComponentVerifier{
		app: app,
	}
}

// VerifySkeletonComponentRegistered verifies that a skeleton component is registered
func (c *ComponentVerifier) VerifySkeletonComponentRegistered(ctx context.Context, componentID string) error {
	if !c.app.IsRunning() {
		return fmt.Errorf("skeleton application is not running")
	}

	components, err := c.getRegisteredComponents(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registered components: %w", err)
	}

	for _, comp := range components {
		if comp == componentID {
			return nil
		}
	}

	return fmt.Errorf("component %s is not registered", componentID)
}

// VerifySkeletonComponentInitialized verifies that a skeleton component is initialized
func (c *ComponentVerifier) VerifySkeletonComponentInitialized(ctx context.Context, componentID string) error {
	if !c.app.IsRunning() {
		return fmt.Errorf("skeleton application is not running")
	}

	// First verify it's registered
	if err := c.VerifySkeletonComponentRegistered(ctx, componentID); err != nil {
		return err
	}

	// Check component status endpoint
	baseURL := c.app.ConnectionString()
	if baseURL == "" {
		return fmt.Errorf("unable to get application connection string")
	}

	statusURL := fmt.Sprintf("%s/api/components/%s/status", baseURL, componentID)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", statusURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to reach component status endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("component %s status endpoint returned status %d", componentID, resp.StatusCode)
	}

	return nil
}

// VerifySkeletonComponentDisposed verifies that a skeleton component is properly disposed
func (c *ComponentVerifier) VerifySkeletonComponentDisposed(ctx context.Context, componentID string) error {
	if !c.app.IsRunning() {
		return fmt.Errorf("skeleton application is not running")
	}

	components, err := c.getRegisteredComponents(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registered components: %w", err)
	}

	for _, comp := range components {
		if comp == componentID {
			return fmt.Errorf("component %s is still registered (not disposed)", componentID)
		}
	}

	return nil
}

// VerifySkeletonComponentMetadata verifies component metadata matches expected values
func (c *ComponentVerifier) VerifySkeletonComponentMetadata(ctx context.Context, componentID string, expected map[string]interface{}) error {
	if !c.app.IsRunning() {
		return fmt.Errorf("skeleton application is not running")
	}

	// Get component metadata
	baseURL := c.app.ConnectionString()
	if baseURL == "" {
		return fmt.Errorf("unable to get application connection string")
	}

	metadataURL := fmt.Sprintf("%s/api/components/%s/metadata", baseURL, componentID)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", metadataURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to reach component metadata endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("component %s metadata endpoint returned status %d", componentID, resp.StatusCode)
	}

	var metadata map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return fmt.Errorf("failed to decode component metadata: %w", err)
	}

	// Verify expected metadata
	for key, expectedValue := range expected {
		actualValue, exists := metadata[key]
		if !exists {
			return fmt.Errorf("component %s missing metadata key %s", componentID, key)
		}
		if actualValue != expectedValue {
			return fmt.Errorf("component %s metadata key %s: expected %v, got %v", componentID, key, expectedValue, actualValue)
		}
	}

	return nil
}

// getRegisteredComponents retrieves the list of registered components
func (c *ComponentVerifier) getRegisteredComponents(ctx context.Context) ([]string, error) {
	baseURL := c.app.ConnectionString()
	if baseURL == "" {
		return nil, fmt.Errorf("unable to get application connection string")
	}

	componentsURL := baseURL + "/api/components"

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", componentsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach components endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("components endpoint returned status %d", resp.StatusCode)
	}

	var components []string
	if err := json.NewDecoder(resp.Body).Decode(&components); err != nil {
		return nil, fmt.Errorf("failed to decode components list: %w", err)
	}

	return components, nil
}
