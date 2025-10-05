package crossplane

import (
	"path/filepath"
	"testing"

	"github.com/kubed-io/krm-helm-fn/testutil"
)

func TestNewCrossplaneProvider(t *testing.T) {
	provider := NewCrossplaneProvider()
	if provider == nil {
		t.Error("NewCrossplaneProvider returned nil")
	}
}

// TestCrossplaneProvider_GenerateReleaseFromExample tests the Crossplane provider using example files
func TestCrossplaneProvider_GenerateReleaseFromExample(t *testing.T) {
	// Load the crossplane example files
	exampleDir := filepath.Join("..", "..", "examples", "crossplane")
	example, err := testutil.LoadExampleFiles(exampleDir)
	if err != nil {
		t.Fatalf("Failed to load example files: %v", err)
	}

	// Parse the HelmRelease from the example
	helmRelease, err := example.ParseHelmRelease()
	if err != nil {
		t.Fatalf("Failed to parse HelmRelease: %v", err)
	}

	// Create Crossplane provider and generate release
	provider := NewCrossplaneProvider()
	release, err := provider.GenerateRelease(helmRelease)
	if err != nil {
		t.Fatalf("GenerateRelease failed: %v", err)
	}

	// Verify the generated release has correct apiVersion and kind
	if release.APIVersion != "helm.crossplane.io/v1beta1" {
		t.Errorf("Expected APIVersion 'helm.crossplane.io/v1beta1', got '%s'", release.APIVersion)
	}

	if release.Kind != "Release" {
		t.Errorf("Expected Kind 'Release', got '%s'", release.Kind)
	}

	// Verify metadata is set correctly
	if release.ObjectMeta.Name != "my-app" {
		t.Errorf("Expected Name 'my-app', got '%s'", release.ObjectMeta.Name)
	}

	// Crossplane Release is cluster-scoped, so namespace should be empty
	if release.ObjectMeta.Namespace != "" {
		t.Errorf("Expected empty Namespace for cluster-scoped resource, got '%s'", release.ObjectMeta.Namespace)
	}

	// For crossplane provider in hello world state, we just check basic structure
	// Future tests will validate the full spec content
}
