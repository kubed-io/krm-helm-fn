package argocd

import (
	"fmt"
	"path/filepath"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
	
	"github.com/kptdev/krm-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/helmfn/types"
	"github.com/kubed-io/krm-helm-fn/testutil"
)

func TestArgoCDProvider_GenerateApplication(t *testing.T) {
	provider := NewArgoCDProvider()

	// Create a sample HelmRelease
	helmRelease := &types.HelmRelease{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "krm.kubed.io",
			Kind:       "HelmRelease",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-app",
			Namespace: "my-system",
		},
		Spec: types.HelmReleaseSpec{
			Provider: "argocd",
			Chart: types.ChartSpec{
				Name:    "hello-world",
				Version: "0.1.0",
				Repo:    "https://helm.github.io/examples",
			},
		},
	}

	// Generate ArgoCD Application
	app, err := provider.GenerateApplication(helmRelease)
	if err != nil {
		t.Fatalf("GenerateApplication failed: %v", err)
	}

	// Verify basic structure
	if app.APIVersion != "argoproj.io/v1alpha1" {
		t.Errorf("Expected APIVersion 'argoproj.io/v1alpha1', got '%s'", app.APIVersion)
	}

	if app.Kind != "Application" {
		t.Errorf("Expected Kind 'Application', got '%s'", app.Kind)
	}

	if app.ObjectMeta.Name != "my-app" {
		t.Errorf("Expected Name 'my-app', got '%s'", app.ObjectMeta.Name)
	}

	if app.ObjectMeta.Namespace != "argocd" {
		t.Errorf("Expected Namespace 'argocd', got '%s'", app.ObjectMeta.Namespace)
	}
}

func TestNewArgoCDProvider(t *testing.T) {
	provider := NewArgoCDProvider()
	if provider == nil {
		t.Error("NewArgoCDProvider returned nil")
	}
}

// TestArgoCDProvider_GenerateApplicationFromExample tests the ArgoCD provider using example files
func TestArgoCDProvider_GenerateApplicationFromExample(t *testing.T) {
	// Load the argocd example files
	exampleDir := filepath.Join("..", "..", "examples", "argocd")
	example, err := testutil.LoadExampleFiles(exampleDir)
	if err != nil {
		t.Fatalf("Failed to load example files: %v", err)
	}
	
	// Parse the HelmRelease from the example
	helmRelease, err := parseHelmReleaseFromKubeObject(example.Release)
	if err != nil {
		t.Fatalf("Failed to parse HelmRelease: %v", err)
	}
	
	// Create ArgoCD provider and generate application
	provider := NewArgoCDProvider()
	app, err := provider.GenerateApplication(helmRelease)
	if err != nil {
		t.Fatalf("GenerateApplication failed: %v", err)
	}
	
	// Verify the generated application has correct apiVersion and kind
	if app.APIVersion != "argoproj.io/v1alpha1" {
		t.Errorf("Expected APIVersion 'argoproj.io/v1alpha1', got '%s'", app.APIVersion)
	}
	
	if app.Kind != "Application" {
		t.Errorf("Expected Kind 'Application', got '%s'", app.Kind)
	}
	
	// Verify metadata is set correctly
	if app.ObjectMeta.Name != "my-app" {
		t.Errorf("Expected Name 'my-app', got '%s'", app.ObjectMeta.Name)
	}
	
	if app.ObjectMeta.Namespace != "argocd" {
		t.Errorf("Expected Namespace 'argocd', got '%s'", app.ObjectMeta.Namespace)
	}
	
	// For argocd provider in hello world state, we just check basic structure
	// Future tests will validate the full spec content
}

// parseHelmReleaseFromKubeObject converts a KubeObject to HelmRelease struct for testing
func parseHelmReleaseFromKubeObject(obj *fn.KubeObject) (*types.HelmRelease, error) {
	yamlBytes := obj.String()
	
	var helmRelease types.HelmRelease
	if err := yaml.Unmarshal([]byte(yamlBytes), &helmRelease); err != nil {
		return nil, fmt.Errorf("failed to unmarshal HelmRelease: %w", err)
	}
	
	return &helmRelease, nil
}