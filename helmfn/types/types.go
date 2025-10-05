package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HelmRelease represents the KRM HelmRelease resource
type HelmRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              HelmReleaseSpec `json:"spec,omitempty"`
}

// HelmReleaseSpec defines the desired state of HelmRelease
type HelmReleaseSpec struct {
	Provider string    `json:"provider,omitempty"`
	Chart    ChartSpec `json:"chart,omitempty"`
	// Values will be added in a future phase
	// ValuesSelector will be added in a future phase
}

// ChartSpec defines the Helm chart details
type ChartSpec struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Repo    string `json:"repo,omitempty"`
}
