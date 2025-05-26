// Package fixtures provides test fixtures and constants for skeleton-testkit integration tests.
package fixtures

const (
	// TestSkeletonAppImage is a minimal skeleton application image for testing
	// This image implements the required skeleton endpoints:
	// - /health - Basic health check
	// - /api/system/health - Skeleton system service health
	// - /api/components - Skeleton components listing
	// - /metrics - Application metrics
	// - /shutdown - Graceful shutdown
	TestSkeletonAppImage = "fintechain/test-skeleton-app:latest"

	// AlternativeSkeletonAppImage provides an alternative test image
	// for testing different scenarios or versions
	AlternativeSkeletonAppImage = "fintechain/test-skeleton-app:v1.0.0"

	// LegacySkeletonAppImage simulates an older skeleton application
	// for backward compatibility testing
	LegacySkeletonAppImage = "fintechain/legacy-skeleton-app:v0.9.0"
)

// GetDefaultTestImage returns the default test image for skeleton applications
func GetDefaultTestImage() string {
	return TestSkeletonAppImage
}

// GetTestImageWithTag returns a test image with a specific tag
func GetTestImageWithTag(tag string) string {
	return "fintechain/test-skeleton-app:" + tag
}
