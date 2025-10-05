package argocd

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/kubed-io/krm-helm-fn/helmfn/types"
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