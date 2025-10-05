package argocd

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/kubed-io/krm-helm-fn/helmfn/types"
)

// ArgoCDApplication represents an ArgoCD Application resource
type ArgoCDApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ArgoCDApplicationSpec `json:"spec,omitempty"`
}

// ArgoCDApplicationSpec defines the desired state of ArgoCD Application
type ArgoCDApplicationSpec struct {
	// For phase 1, we'll keep this minimal
	// Full spec will be added in future phases
}

// ArgoCDProvider handles the transformation of HelmRelease to ArgoCD Application
type ArgoCDProvider struct{}

// NewArgoCDProvider creates a new ArgoCD provider instance
func NewArgoCDProvider() *ArgoCDProvider {
	return &ArgoCDProvider{}
}

// GenerateApplication creates an ArgoCD Application resource from a HelmRelease
func (p *ArgoCDProvider) GenerateApplication(helmRelease *types.HelmRelease) (*ArgoCDApplication, error) {
	// Create the ArgoCD Application with basic structure
	app := &ArgoCDApplication{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      helmRelease.ObjectMeta.Name,
			Namespace: "argocd", // Default ArgoCD namespace
		},
		Spec: ArgoCDApplicationSpec{
			// Minimal spec for phase 1
		},
	}

	return app, nil
}