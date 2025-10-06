package rancher

import (
	"path/filepath"
	"testing"

	"github.com/kubed-io/krm-helm-fn/testutil"
)

func TestNewRancherProvider(t *testing.T) {
	provider := NewRancherProvider()
	if provider == nil {
		t.Error("NewRancherProvider returned nil")
	}
}

// TestRancherProvider_GenerateHelmChartFromExample tests the Rancher provider using example files
func TestRancherProvider_GenerateHelmChartFromExample(t *testing.T) {
	// Load the rancher example files
	exampleDir := filepath.Join("..", "..", "examples", "rancher")
	example, err := testutil.LoadExampleFiles(exampleDir)
	if err != nil {
		t.Fatalf("Failed to load example files: %v", err)
	}

	// Parse the HelmRelease from the example
	helmRelease, err := example.ParseHelmRelease()
	if err != nil {
		t.Fatalf("Failed to parse HelmRelease: %v", err)
	}

	// Create Rancher provider and generate helm chart
	provider := NewRancherProvider()
	helmChart, err := provider.GenerateHelmChart(helmRelease)
	if err != nil {
		t.Fatalf("GenerateHelmChart failed: %v", err)
	}

	// Verify the generated helm chart has correct apiVersion and kind
	if helmChart.APIVersion != "helm.cattle.io/v1" {
		t.Errorf("Expected APIVersion 'helm.cattle.io/v1', got '%s'", helmChart.APIVersion)
	}

	if helmChart.Kind != "HelmChart" {
		t.Errorf("Expected Kind 'HelmChart', got '%s'", helmChart.Kind)
	}

	// Verify metadata is set correctly
	if helmChart.ObjectMeta.Name != "my-app" {
		t.Errorf("Expected Name 'my-app', got '%s'", helmChart.ObjectMeta.Name)
	}

	if helmChart.ObjectMeta.Namespace != "my-system" {
		t.Errorf("Expected Namespace 'my-system', got '%s'", helmChart.ObjectMeta.Namespace)
	}

	// For rancher provider in hello world state, we just check basic structure
	// Future tests will validate the full spec content
}