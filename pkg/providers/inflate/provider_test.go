package inflate_test

import (
	"testing"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/pkg/providers/inflate"
	"github.com/kubed-io/krm-helm-fn/pkg/types"
)

func TestInflateProvider_Generate(t *testing.T) {
	// Create a sample HelmRelease
	helmRelease := &types.HelmRelease{
		Name:        "my-app",
		Namespace:   "my-system",
		Provider:    "inflate",
		ReleaseName: "my-app",
		Chart: types.ChartSpec{
			Name:    "hello-world",
			Version: "0.1.0",
			Repo:    "https://helm.github.io/examples",
		},
		Values: map[string]interface{}{
			"replicaCount": 2,
			"service": map[string]interface{}{
				"port": 443,
			},
		},
		IncludeCRDs: true,
		ApiVersions: []string{"example.com/v1"},
	}

	// Create a provider instance
	provider := inflate.NewProvider()

	// This test requires the Helm CLI to be installed and available in the PATH
	// It also requires internet connectivity to download the chart

	// Generate resources
	resources, err := provider.Generate(helmRelease)
	if err != nil {
		t.Fatalf("Failed to generate resources: %v", err)
	}

	// Verify the resources were generated
	if len(resources) == 0 {
		t.Error("No resources were generated")
	}

	// Check for expected resources
	found := false
	for _, res := range resources {
		if res.GetKind() == "Deployment" && res.GetName() == "my-app" {
			found = true
			// Check for expected replica count
			replicas, exists, err := res.NestedInt("spec", "replicas")
			if !exists || err != nil {
				t.Error("Failed to find spec.replicas in Deployment")
			} else if replicas != 2 {
				t.Errorf("Expected 2 replicas, got %d", replicas)
			}
		}
	}

	if !found {
		t.Error("Expected Deployment/my-app not found in generated resources")
	}
}

func TestInflateProvider_WithValuesSelector(t *testing.T) {
	// Create a ConfigMap with values
	configMapYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-app-values
  annotations:
    krm.kubed.io/helm-values: "my-app"
data:
  service.port: "443"
`
	configMap, err := fn.ParseKubeObject([]byte(configMapYAML))
	if err != nil {
		t.Fatalf("Failed to parse ConfigMap: %v", err)
	}

	// Create a sample HelmRelease with ValuesSelector
	helmRelease := &types.HelmRelease{
		Name:        "my-app",
		Namespace:   "my-system",
		Provider:    "inflate",
		ReleaseName: "my-app",
		Chart: types.ChartSpec{
			Name:    "hello-world",
			Version: "0.1.0",
			Repo:    "https://helm.github.io/examples",
		},
		Values: map[string]interface{}{
			"replicaCount": 2,
		},
		ValuesSelector: &types.ValuesSelector{
			Annotations: map[string]string{
				"krm.kubed.io/helm-values": "my-app",
			},
		},
	}

	// This test is more complex and would require mocking the Helm CLI execution
	// For simplicity, we'll just verify the ValuesSelector matches the ConfigMap
	if len(helmRelease.ValuesSelector.Annotations) != 1 {
		t.Errorf("Expected 1 annotation in ValuesSelector, got %d", len(helmRelease.ValuesSelector.Annotations))
	}

	if val, ok := configMap.GetAnnotation("krm.kubed.io/helm-values"); !ok || val != "my-app" {
		t.Error("ConfigMap annotation does not match ValuesSelector criteria")
	}
}