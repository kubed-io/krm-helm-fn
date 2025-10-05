package fluxcd

import (
	"github.com/kubed-io/krm-helm-fn/helmfn/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FluxCDHelmRelease represents a FluxCD HelmRelease resource
type FluxCDHelmRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FluxCDHelmReleaseSpec `json:"spec,omitempty"`
}

// FluxCDHelmReleaseSpec defines the desired state of FluxCD HelmRelease
type FluxCDHelmReleaseSpec struct {
	// For phase 1, we'll keep this minimal
	// Full spec will be added in future phases
}

// FluxCDHelmRepository represents a FluxCD HelmRepository resource
type FluxCDHelmRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FluxCDHelmRepositorySpec `json:"spec,omitempty"`
}

// FluxCDHelmRepositorySpec defines the desired state of FluxCD HelmRepository
type FluxCDHelmRepositorySpec struct {
	// For phase 1, we'll keep this minimal
	// Full spec will be added in future phases
}

// FluxCDProvider handles the transformation of HelmRelease to FluxCD resources
type FluxCDProvider struct{}

// NewFluxCDProvider creates a new FluxCD provider instance
func NewFluxCDProvider() *FluxCDProvider {
	return &FluxCDProvider{}
}

// GenerateHelmRelease creates a FluxCD HelmRelease resource from a HelmRelease
func (p *FluxCDProvider) GenerateHelmRelease(helmRelease *types.HelmRelease) (*FluxCDHelmRelease, error) {
	// Create the FluxCD HelmRelease with basic structure
	fluxHelmRelease := &FluxCDHelmRelease{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "helm.toolkit.fluxcd.io/v2beta1",
			Kind:       "HelmRelease",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      helmRelease.ObjectMeta.Name,
			Namespace: helmRelease.ObjectMeta.Namespace,
		},
		Spec: FluxCDHelmReleaseSpec{
			// Minimal spec for phase 1
		},
	}

	return fluxHelmRelease, nil
}

// GenerateHelmRepository creates a FluxCD HelmRepository resource from a HelmRelease
func (p *FluxCDProvider) GenerateHelmRepository(helmRelease *types.HelmRelease) (*FluxCDHelmRepository, error) {
	// Create the FluxCD HelmRepository with basic structure
	helmRepo := &FluxCDHelmRepository{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "source.toolkit.fluxcd.io/v1beta1",
			Kind:       "HelmRepository",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      helmRelease.ObjectMeta.Name,
			Namespace: helmRelease.ObjectMeta.Namespace,
		},
		Spec: FluxCDHelmRepositorySpec{
			// Minimal spec for phase 1
		},
	}

	return helmRepo, nil
}
