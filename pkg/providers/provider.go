package providers

import (
	"fmt"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/pkg/providers/inflate"
	"github.com/kubed-io/krm-helm-fn/pkg/types"
)

// Provider defines the interface for all Helm release providers
type Provider interface {
	// Generate creates the Kubernetes resources for a Helm release
	Generate(release *types.HelmRelease) ([]*fn.KubeObject, error)
}

// GetProvider returns the appropriate provider implementation based on the type
func GetProvider(providerType string) (Provider, error) {
	switch providerType {
	case "inflate":
		return inflate.NewProvider(), nil
	case "argocd":
		// TODO: Implement ArgoCD provider
		return nil, fmt.Errorf("argocd provider not yet implemented")
	case "fluxcd":
		// TODO: Implement FluxCD provider
		return nil, fmt.Errorf("fluxcd provider not yet implemented")
	case "crossplane":
		// TODO: Implement Crossplane provider
		return nil, fmt.Errorf("crossplane provider not yet implemented")
	case "rancher":
		// TODO: Implement Rancher provider
		return nil, fmt.Errorf("rancher provider not yet implemented")
	default:
		return nil, fmt.Errorf("unknown provider type: %s", providerType)
	}
}