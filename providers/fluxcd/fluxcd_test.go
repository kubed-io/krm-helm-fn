package fluxcd

import (
	"path/filepath"
	"testing"

	"github.com/kubed-io/krm-helm-fn/testutil"
)

func TestNewFluxCDProvider(t *testing.T) {
	provider := NewFluxCDProvider()
	if provider == nil {
		t.Error("NewFluxCDProvider returned nil")
	}
}

// TestFluxCDProvider_GenerateHelmReleaseFromExample tests the FluxCD provider using example files
func TestFluxCDProvider_GenerateHelmReleaseFromExample(t *testing.T) {
	// Load the fluxcd example files
	exampleDir := filepath.Join("..", "..", "examples", "fluxcd")
	example, err := testutil.LoadExampleFiles(exampleDir)
	if err != nil {
		t.Fatalf("Failed to load example files: %v", err)
	}

	// Parse the HelmRelease from the example
	helmRelease, err := example.ParseHelmRelease()
	if err != nil {
		t.Fatalf("Failed to parse HelmRelease: %v", err)
	}

	// Create FluxCD provider and generate helm release
	provider := NewFluxCDProvider()
	fluxHelmRelease, err := provider.GenerateHelmRelease(helmRelease)
	if err != nil {
		t.Fatalf("GenerateHelmRelease failed: %v", err)
	}

	// Verify the generated helm release has correct apiVersion and kind
	if fluxHelmRelease.APIVersion != "helm.toolkit.fluxcd.io/v2beta1" {
		t.Errorf("Expected APIVersion 'helm.toolkit.fluxcd.io/v2beta1', got '%s'", fluxHelmRelease.APIVersion)
	}

	if fluxHelmRelease.Kind != "HelmRelease" {
		t.Errorf("Expected Kind 'HelmRelease', got '%s'", fluxHelmRelease.Kind)
	}

	// Verify metadata is set correctly
	if fluxHelmRelease.ObjectMeta.Name != "my-app" {
		t.Errorf("Expected Name 'my-app', got '%s'", fluxHelmRelease.ObjectMeta.Name)
	}

	if fluxHelmRelease.ObjectMeta.Namespace != "my-system" {
		t.Errorf("Expected Namespace 'my-system', got '%s'", fluxHelmRelease.ObjectMeta.Namespace)
	}

	// For fluxcd provider in hello world state, we just check basic structure
	// Future tests will validate the full spec content
}

// TestFluxCDProvider_GenerateHelmRepositoryFromExample tests the FluxCD provider helm repository generation
func TestFluxCDProvider_GenerateHelmRepositoryFromExample(t *testing.T) {
	// Load the fluxcd example files
	exampleDir := filepath.Join("..", "..", "examples", "fluxcd")
	example, err := testutil.LoadExampleFiles(exampleDir)
	if err != nil {
		t.Fatalf("Failed to load example files: %v", err)
	}

	// Parse the HelmRelease from the example
	helmRelease, err := example.ParseHelmRelease()
	if err != nil {
		t.Fatalf("Failed to parse HelmRelease: %v", err)
	}

	// Create FluxCD provider and generate helm repository
	provider := NewFluxCDProvider()
	helmRepo, err := provider.GenerateHelmRepository(helmRelease)
	if err != nil {
		t.Fatalf("GenerateHelmRepository failed: %v", err)
	}

	// Verify the generated helm repository has correct apiVersion and kind
	if helmRepo.APIVersion != "source.toolkit.fluxcd.io/v1beta1" {
		t.Errorf("Expected APIVersion 'source.toolkit.fluxcd.io/v1beta1', got '%s'", helmRepo.APIVersion)
	}

	if helmRepo.Kind != "HelmRepository" {
		t.Errorf("Expected Kind 'HelmRepository', got '%s'", helmRepo.Kind)
	}

	// Verify metadata is set correctly
	if helmRepo.ObjectMeta.Name != "my-app" {
		t.Errorf("Expected Name 'my-app', got '%s'", helmRepo.ObjectMeta.Name)
	}

	if helmRepo.ObjectMeta.Namespace != "my-system" {
		t.Errorf("Expected Namespace 'my-system', got '%s'", helmRepo.ObjectMeta.Namespace)
	}

	// For fluxcd provider in hello world state, we just check basic structure
	// Future tests will validate the full spec content
}
