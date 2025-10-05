package helmfn

import (
	"path/filepath"
	"testing"

	"github.com/kptdev/krm-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/testutil"
)

// TestProcessArgoCDExample tests the full processor pipeline using the argocd example
func TestProcessArgoCDExample(t *testing.T) {
	// Load the argocd example files
	exampleDir := filepath.Join("..", "examples", "argocd")
	example, err := testutil.LoadExampleFiles(exampleDir)
	if err != nil {
		t.Fatalf("Failed to load example files: %v", err)
	}

	// Create ResourceList from example
	rl := example.CreateResourceList()

	// Process the ResourceList
	_, err = Process(rl)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Verify that an output was generated
	if len(rl.Items) < 2 { // Original ConfigMap + generated ArgoCD Application
		t.Errorf("Expected at least 2 items in output (ConfigMap + Application), got %d", len(rl.Items))
	}

	// Find the generated ArgoCD Application
	var generatedApp *fn.KubeObject
	for _, item := range rl.Items {
		if item.GetKind() == "Application" && item.GetAPIVersion() == "argoproj.io/v1alpha1" {
			generatedApp = item
			break
		}
	}

	if generatedApp == nil {
		t.Fatal("No ArgoCD Application was generated")
	}

	// Verify the generated application has correct basic properties
	if generatedApp.GetName() != "my-app" {
		t.Errorf("Expected application name 'my-app', got '%s'", generatedApp.GetName())
	}

	if generatedApp.GetNamespace() != "argocd" {
		t.Errorf("Expected application namespace 'argocd', got '%s'", generatedApp.GetNamespace())
	}
}
