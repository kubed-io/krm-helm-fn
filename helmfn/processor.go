package helmfn

import (
	"fmt"

	"github.com/kptdev/krm-functions-sdk/go/fn"
	"sigs.k8s.io/yaml"
	"github.com/kubed-io/krm-helm-fn/providers/argocd"
	"github.com/kubed-io/krm-helm-fn/helmfn/types"
)

// Process is the main entry point for the KRM function
func Process(rl *fn.ResourceList) (bool, error) {
	// Check if we have a functionConfig (HelmRelease)
	if rl.FunctionConfig == nil {
		return false, fmt.Errorf("no functionConfig provided")
	}
	
	// Check if this is a HelmRelease resource
	if rl.FunctionConfig.GetAPIVersion() != "krm.kubed.io" || rl.FunctionConfig.GetKind() != "HelmRelease" {
		return false, fmt.Errorf("functionConfig must be a HelmRelease, got %s/%s", 
			rl.FunctionConfig.GetAPIVersion(), rl.FunctionConfig.GetKind())
	}
	
	// Parse the HelmRelease from functionConfig
	helmRelease, err := parseHelmRelease(rl.FunctionConfig)
	if err != nil {
		return false, fmt.Errorf("failed to parse HelmRelease from functionConfig: %w", err)
	}
	
	// Process based on provider
	switch helmRelease.Spec.Provider {
	case "argocd":
		if err := processArgoCDProvider(rl, helmRelease); err != nil {
			return false, fmt.Errorf("failed to process ArgoCD provider: %w", err)
		}
	default:
		return false, fmt.Errorf("unsupported provider: %s", helmRelease.Spec.Provider)
	}
	
	return false, nil
}

// parseHelmRelease converts a KRM object to HelmRelease struct
func parseHelmRelease(obj *fn.KubeObject) (*types.HelmRelease, error) {
	yamlBytes := obj.String()
	
	var helmRelease types.HelmRelease
	if err := yaml.Unmarshal([]byte(yamlBytes), &helmRelease); err != nil {
		return nil, fmt.Errorf("failed to unmarshal HelmRelease: %w", err)
	}
	
	return &helmRelease, nil
}

// processArgoCDProvider handles ArgoCD provider processing
func processArgoCDProvider(rl *fn.ResourceList, helmRelease *types.HelmRelease) error {
	// Create ArgoCD provider
	provider := argocd.NewArgoCDProvider()
	
	// Generate ArgoCD Application
	app, err := provider.GenerateApplication(helmRelease)
	if err != nil {
		return fmt.Errorf("failed to generate ArgoCD application: %w", err)
	}
	
	// Convert to KubeObject
	appBytes, err := yaml.Marshal(app)
	if err != nil {
		return fmt.Errorf("failed to marshal ArgoCD application: %w", err)
	}
	
	appObj, err := fn.ParseKubeObject(appBytes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal to KubeObject: %w", err)
	}
	
	// Add the ArgoCD Application to the output items
	rl.Items = append(rl.Items, appObj)
	
	return nil
}