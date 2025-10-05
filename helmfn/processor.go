package helmfn

import (
	"fmt"

	"github.com/kptdev/krm-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/helmfn/types"
	"github.com/kubed-io/krm-helm-fn/providers/argocd"
	"github.com/kubed-io/krm-helm-fn/providers/crossplane"
	"github.com/kubed-io/krm-helm-fn/providers/fluxcd"
	"sigs.k8s.io/yaml"
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

	DebugLog("Processing HelmRelease %s/%s with provider: %s", helmRelease.ObjectMeta.Namespace, helmRelease.ObjectMeta.Name, helmRelease.Spec.Provider)

	// Process based on provider
	switch helmRelease.Spec.Provider {
	case "argocd":
		DebugLog("Processing ArgoCD provider")
		if err := processArgoCDProvider(rl, helmRelease); err != nil {
			return false, fmt.Errorf("failed to process ArgoCD provider: %w", err)
		}
	case "crossplane":
		DebugLog("Processing Crossplane provider")
		if err := processCrossplaneProvider(rl, helmRelease); err != nil {
			return false, fmt.Errorf("failed to process Crossplane provider: %w", err)
		}
	case "fluxcd":
		DebugLog("Processing FluxCD provider")
		if err := processFluxCDProvider(rl, helmRelease); err != nil {
			return false, fmt.Errorf("failed to process FluxCD provider: %w", err)
		}
	default:
		return false, fmt.Errorf("unsupported provider: %s", helmRelease.Spec.Provider)
	}

	// Return true to indicate the function made changes to the resource list
	// (added provider-specific resources like ArgoCD Application)
	return true, nil
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

// processFluxCDProvider handles FluxCD provider processing
func processFluxCDProvider(rl *fn.ResourceList, helmRelease *types.HelmRelease) error {
	// Create FluxCD provider
	provider := fluxcd.NewFluxCDProvider()

	// Generate FluxCD HelmRelease
	fluxHelmRelease, err := provider.GenerateHelmRelease(helmRelease)
	if err != nil {
		return fmt.Errorf("failed to generate FluxCD HelmRelease: %w", err)
	}

	// Convert FluxCD HelmRelease to KubeObject
	helmReleaseBytes, err := yaml.Marshal(fluxHelmRelease)
	if err != nil {
		return fmt.Errorf("failed to marshal FluxCD HelmRelease: %w", err)
	}

	helmReleaseObj, err := fn.ParseKubeObject(helmReleaseBytes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal FluxCD HelmRelease to KubeObject: %w", err)
	}

	// Generate FluxCD HelmRepository
	helmRepo, err := provider.GenerateHelmRepository(helmRelease)
	if err != nil {
		return fmt.Errorf("failed to generate FluxCD HelmRepository: %w", err)
	}

	// Convert FluxCD HelmRepository to KubeObject
	helmRepoBytes, err := yaml.Marshal(helmRepo)
	if err != nil {
		return fmt.Errorf("failed to marshal FluxCD HelmRepository: %w", err)
	}

	helmRepoObj, err := fn.ParseKubeObject(helmRepoBytes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal FluxCD HelmRepository to KubeObject: %w", err)
	}

	// Add both FluxCD resources to the output items
	rl.Items = append(rl.Items, helmReleaseObj, helmRepoObj)

	DebugLog("Added FluxCD HelmRelease and HelmRepository to output")

	return nil
}

// processCrossplaneProvider handles Crossplane provider processing
func processCrossplaneProvider(rl *fn.ResourceList, helmRelease *types.HelmRelease) error {
	// Create Crossplane provider
	provider := crossplane.NewCrossplaneProvider()

	// Generate Crossplane Release
	release, err := provider.GenerateRelease(helmRelease)
	if err != nil {
		return fmt.Errorf("failed to generate Crossplane release: %w", err)
	}

	// Convert to KubeObject
	releaseBytes, err := yaml.Marshal(release)
	if err != nil {
		return fmt.Errorf("failed to marshal Crossplane release: %w", err)
	}

	releaseObj, err := fn.ParseKubeObject(releaseBytes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal to KubeObject: %w", err)
	}

	// Add the Crossplane Release to the output items
	rl.Items = append(rl.Items, releaseObj)

	DebugLog("Added Crossplane Release to output")

	return nil
}
