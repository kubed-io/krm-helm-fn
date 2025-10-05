package crossplane

import (
	"github.com/kubed-io/krm-helm-fn/helmfn/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CrossplaneRelease represents a Crossplane Helm Release resource
type CrossplaneRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CrossplaneReleaseSpec `json:"spec,omitempty"`
}

// CrossplaneReleaseSpec defines the desired state of Crossplane Release
type CrossplaneReleaseSpec struct {
	// For phase 1, we'll keep this minimal
	// Full spec will be added in future phases
}

// CrossplaneProvider handles the transformation of HelmRelease to Crossplane Release
type CrossplaneProvider struct{}

// NewCrossplaneProvider creates a new Crossplane provider instance
func NewCrossplaneProvider() *CrossplaneProvider {
	return &CrossplaneProvider{}
}

// GenerateRelease creates a Crossplane Release resource from a HelmRelease
func (p *CrossplaneProvider) GenerateRelease(helmRelease *types.HelmRelease) (*CrossplaneRelease, error) {
	// Create the Crossplane Release with basic structure
	release := &CrossplaneRelease{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "helm.crossplane.io/v1beta1",
			Kind:       "Release",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: helmRelease.ObjectMeta.Name,
			// Crossplane Helm Release is cluster-scoped, so no namespace
		},
		Spec: CrossplaneReleaseSpec{
			// Minimal spec for phase 1
		},
	}

	return release, nil
}