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

// TestProcessFluxCDExample tests the full processor pipeline using the fluxcd example
func TestProcessFluxCDExample(t *testing.T) {
	// Load the fluxcd example files
	exampleDir := filepath.Join("..", "examples", "fluxcd")
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

	// Verify that outputs were generated
	if len(rl.Items) < 3 { // Original ConfigMap + generated HelmRelease + HelmRepository
		t.Errorf("Expected at least 3 items in output (ConfigMap + HelmRelease + HelmRepository), got %d", len(rl.Items))
	}

	// Find the generated FluxCD HelmRelease
	var generatedRelease *fn.KubeObject
	var generatedRepo *fn.KubeObject
	for _, item := range rl.Items {
		if item.GetKind() == "HelmRelease" && item.GetAPIVersion() == "helm.toolkit.fluxcd.io/v2beta1" {
			generatedRelease = item
		}
		if item.GetKind() == "HelmRepository" && item.GetAPIVersion() == "source.toolkit.fluxcd.io/v1beta1" {
			generatedRepo = item
		}
	}

	if generatedRelease == nil {
		t.Fatal("No FluxCD HelmRelease was generated")
	}

	if generatedRepo == nil {
		t.Fatal("No FluxCD HelmRepository was generated")
	}

	// Verify the generated release has correct basic properties
	if generatedRelease.GetName() != "my-app" {
		t.Errorf("Expected release name 'my-app', got '%s'", generatedRelease.GetName())
	}

	if generatedRelease.GetNamespace() != "my-system" {
		t.Errorf("Expected release namespace 'my-system', got '%s'", generatedRelease.GetNamespace())
	}

	// Verify the generated repository has correct basic properties
	if generatedRepo.GetName() != "my-app" {
		t.Errorf("Expected repository name 'my-app', got '%s'", generatedRepo.GetName())
	}

	if generatedRepo.GetNamespace() != "my-system" {
		t.Errorf("Expected repository namespace 'my-system', got '%s'", generatedRepo.GetNamespace())
	}
}

// TestProcessCrossplaneExample tests the full processor pipeline using the crossplane example
func TestProcessCrossplaneExample(t *testing.T) {
	// Load the crossplane example files
	exampleDir := filepath.Join("..", "examples", "crossplane")
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

	// Verify that outputs were generated
	if len(rl.Items) < 2 { // Original ConfigMap + generated Release
		t.Errorf("Expected at least 2 items in output (ConfigMap + Release), got %d", len(rl.Items))
	}

	// Find the generated Crossplane Release
	var generatedRelease *fn.KubeObject
	for _, item := range rl.Items {
		if item.GetKind() == "Release" && item.GetAPIVersion() == "helm.crossplane.io/v1beta1" {
			generatedRelease = item
		}
	}

	if generatedRelease == nil {
		t.Fatal("No Crossplane Release was generated")
	}

	// Verify the generated release has correct basic properties
	if generatedRelease.GetName() != "my-app" {
		t.Errorf("Expected release name 'my-app', got '%s'", generatedRelease.GetName())
	}

	// Crossplane Release is cluster-scoped, so namespace should be empty
	if generatedRelease.GetNamespace() != "" {
		t.Errorf("Expected empty namespace for cluster-scoped resource, got '%s'", generatedRelease.GetNamespace())
	}
}

// TestProcessRancherExample tests the full processor pipeline using the rancher example
func TestProcessRancherExample(t *testing.T) {
	// Load the rancher example files
	exampleDir := filepath.Join("..", "examples", "rancher")
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
	if len(rl.Items) < 2 { // Original ConfigMap + generated Rancher HelmChart
		t.Errorf("Expected at least 2 items in output (ConfigMap + HelmChart), got %d", len(rl.Items))
	}

	// Find the generated Rancher HelmChart
	var generatedChart *fn.KubeObject
	for _, item := range rl.Items {
		if item.GetKind() == "HelmChart" && item.GetAPIVersion() == "helm.cattle.io/v1" {
			generatedChart = item
			break
		}
	}

	if generatedChart == nil {
		t.Fatal("No Rancher HelmChart was generated")
	}

	// Verify the generated chart has correct basic properties
	if generatedChart.GetName() != "my-app" {
		t.Errorf("Expected chart name 'my-app', got '%s'", generatedChart.GetName())
	}

	if generatedChart.GetNamespace() != "my-system" {
		t.Errorf("Expected chart namespace 'my-system', got '%s'", generatedChart.GetNamespace())
	}
}
