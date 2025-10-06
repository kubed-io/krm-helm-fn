package rancher

import (
	"github.com/kubed-io/krm-helm-fn/helmfn/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RancherHelmChart represents a Rancher K3s HelmChart resource
type RancherHelmChart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RancherHelmChartSpec `json:"spec,omitempty"`
}

// RancherHelmChartSpec defines the desired state of Rancher K3s HelmChart
type RancherHelmChartSpec struct {
	// For phase 1, we'll keep this minimal
	// Full spec will be added in future phases
}

// RancherProvider handles the transformation of HelmRelease to Rancher K3s HelmChart
type RancherProvider struct{}

// NewRancherProvider creates a new Rancher provider instance
func NewRancherProvider() *RancherProvider {
	return &RancherProvider{}
}

// GenerateHelmChart creates a Rancher K3s HelmChart resource from a HelmRelease
func (p *RancherProvider) GenerateHelmChart(helmRelease *types.HelmRelease) (*RancherHelmChart, error) {
	// Create the Rancher K3s HelmChart with basic structure
	helmChart := &RancherHelmChart{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "helm.cattle.io/v1",
			Kind:       "HelmChart",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      helmRelease.ObjectMeta.Name,
			Namespace: helmRelease.ObjectMeta.Namespace,
		},
		Spec: RancherHelmChartSpec{
			// Minimal spec for phase 1
		},
	}

	return helmChart, nil
}