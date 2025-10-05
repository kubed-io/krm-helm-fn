package testutil

import (
	"path/filepath"
	"testing"
)

func TestLoadExampleFiles(t *testing.T) {
	// Test loading the argocd example
	exampleDir := filepath.Join("..", "examples", "argocd")
	example, err := LoadExampleFiles(exampleDir)
	if err != nil {
		t.Fatalf("Failed to load example files: %v", err)
	}
	
	// Verify release was loaded
	if example.Release == nil {
		t.Error("Release was not loaded")
	}
	
	if example.Release.GetAPIVersion() != "krm.kubed.io" {
		t.Errorf("Expected release APIVersion 'krm.kubed.io', got '%s'", example.Release.GetAPIVersion())
	}
	
	if example.Release.GetKind() != "HelmRelease" {
		t.Errorf("Expected release Kind 'HelmRelease', got '%s'", example.Release.GetKind())
	}
	
	// Verify values ConfigMap was created
	if example.Values == nil {
		t.Error("Values ConfigMap was not created")
	}
	
	if example.Values.GetAPIVersion() != "v1" {
		t.Errorf("Expected values APIVersion 'v1', got '%s'", example.Values.GetAPIVersion())
	}
	
	if example.Values.GetKind() != "ConfigMap" {
		t.Errorf("Expected values Kind 'ConfigMap', got '%s'", example.Values.GetKind())
	}
	
	// Verify expected output was loaded
	if len(example.Expected) == 0 {
		t.Error("Expected output was not loaded")
	}
	
	// For argocd example, we expect one Application resource
	if len(example.Expected) != 1 {
		t.Errorf("Expected 1 output resource, got %d", len(example.Expected))
	}
	
	expectedApp := example.Expected[0]
	if expectedApp.GetAPIVersion() != "argoproj.io/v1alpha1" {
		t.Errorf("Expected output APIVersion 'argoproj.io/v1alpha1', got '%s'", expectedApp.GetAPIVersion())
	}
	
	if expectedApp.GetKind() != "Application" {
		t.Errorf("Expected output Kind 'Application', got '%s'", expectedApp.GetKind())
	}
}

func TestCreateResourceList(t *testing.T) {
	// Load example and create ResourceList
	exampleDir := filepath.Join("..", "examples", "argocd")
	example, err := LoadExampleFiles(exampleDir)
	if err != nil {
		t.Fatalf("Failed to load example files: %v", err)
	}
	
	rl := example.CreateResourceList()
	
	// Verify ResourceList structure
	if rl.FunctionConfig == nil {
		t.Error("FunctionConfig is nil")
	}
	
	if len(rl.Items) != 1 {
		t.Errorf("Expected 1 item in ResourceList, got %d", len(rl.Items))
	}
	
	// Verify the item is the values ConfigMap
	configMap := rl.Items[0]
	if configMap.GetKind() != "ConfigMap" {
		t.Errorf("Expected item to be ConfigMap, got %s", configMap.GetKind())
	}
}