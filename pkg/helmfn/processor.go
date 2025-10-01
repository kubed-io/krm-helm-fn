package helmfn

import (
	"fmt"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/pkg/providers"
	"github.com/kubed-io/krm-helm-fn/pkg/types"
)

const (
	apiVersionKey = "apiVersion"
	kindKey       = "kind"
	metadataKey   = "metadata"
	specKey       = "spec"
	
	defaultAPIVersion = "krm.kubed.io"
	defaultKind       = "HelmRelease"

	// spec fields
	providerKey      = "provider"
	chartKey         = "chart"
	valuesKey        = "values"
	valuesSelectorKey = "valuesSelector"
	includecrdsKey   = "includeCRDs"
	apiVersionsKey   = "apiVersions"
	skipTestsKey     = "skipTests"
	releaseNameKey   = "releaseName"
)

// Process implements the KRM function processing logic
func Process(rl *fn.ResourceList) (bool, error) {
	// Check if function config exists
	if rl.FunctionConfig == nil {
		return false, fmt.Errorf("functionConfig is required")
	}

	// Validate function config
	if err := validateFunctionConfig(rl.FunctionConfig); err != nil {
		return false, err
	}

	// Parse HelmRelease from function config
	helmRelease, err := parseHelmRelease(rl.FunctionConfig)
	if err != nil {
		return false, err
	}

	// Determine the provider
	provider, err := providers.GetProvider(helmRelease.Provider)
	if err != nil {
		return false, err
	}

	// Check for ConfigMaps and Secrets matching the ValuesSelector
	if helmRelease.ValuesSelector != nil {
		err = processValuesSelector(helmRelease, rl.Items)
		if err != nil {
			return false, err
		}
	}

	// Generate resources based on the provider
	generatedResources, err := provider.Generate(helmRelease)
	if err != nil {
		return false, err
	}

	// Append generated resources to the ResourceList
	rl.Items = append(rl.Items, generatedResources...)

	return true, nil
}

func validateFunctionConfig(fc *fn.KubeObject) error {
	// Check API version
	apiVersion, found, err := fc.NestedString(apiVersionKey)
	if !found || err != nil {
		return fmt.Errorf("apiVersion is required in functionConfig")
	}
	if apiVersion != defaultAPIVersion {
		return fmt.Errorf("expected apiVersion: %s, got: %s", defaultAPIVersion, apiVersion)
	}

	// Check Kind
	kind, found, err := fc.NestedString(kindKey)
	if !found || err != nil {
		return fmt.Errorf("kind is required in functionConfig")
	}
	if kind != defaultKind {
		return fmt.Errorf("expected kind: %s, got: %s", defaultKind, kind)
	}

	// Check spec
	if !fc.IsNestedMap(specKey) {
		return fmt.Errorf("spec is required in functionConfig")
	}

	// Check provider
	provider, found, err := fc.NestedString(specKey, providerKey)
	if !found || err != nil {
		return fmt.Errorf("spec.provider is required in functionConfig")
	}
	
	// Check chart
	if !fc.IsNestedMap(specKey, chartKey) {
		return fmt.Errorf("spec.chart is required in functionConfig")
	}

	return nil
}

func parseHelmRelease(fc *fn.KubeObject) (*types.HelmRelease, error) {
	// Get metadata.name
	name, found, err := fc.NestedString(metadataKey, "name")
	if !found || err != nil {
		return nil, fmt.Errorf("metadata.name is required in functionConfig")
	}
	
	// Get metadata.namespace
	namespace, found, err := fc.NestedString(metadataKey, "namespace")
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata.namespace: %v", err)
	}
	
	// Get provider
	provider, found, err := fc.NestedString(specKey, providerKey)
	if !found || err != nil {
		return nil, fmt.Errorf("spec.provider is required in functionConfig")
	}

	// Get chart
	chart := types.ChartSpec{}
	chartName, found, err := fc.NestedString(specKey, chartKey, "name")
	if !found || err != nil {
		return nil, fmt.Errorf("spec.chart.name is required in functionConfig")
	}
	chart.Name = chartName

	chartVersion, found, err := fc.NestedString(specKey, chartKey, "version")
	if err != nil {
		return nil, fmt.Errorf("failed to get spec.chart.version: %v", err)
	}
	chart.Version = chartVersion

	chartRepo, found, err := fc.NestedString(specKey, chartKey, "repo")
	if err != nil {
		return nil, fmt.Errorf("failed to get spec.chart.repo: %v", err)
	}
	chart.Repo = chartRepo

	// Create HelmRelease
	helmRelease := &types.HelmRelease{
		Name:      name,
		Namespace: namespace,
		Provider:  provider,
		Chart:     chart,
	}
	
	// Get values
	if fc.IsNestedMap(specKey, valuesKey) {
		values, err := fc.NestedMap(specKey, valuesKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get spec.values: %v", err)
		}
		helmRelease.Values = values
	}
	
	// Get valuesSelector
	if fc.IsNestedMap(specKey, valuesSelectorKey) {
		valuesSelector := &types.ValuesSelector{}
		
		// Kind
		if kind, found, err := fc.NestedString(specKey, valuesSelectorKey, "kind"); found && err == nil {
			valuesSelector.Kind = kind
		}
		
		// Labels
		if fc.IsNestedMap(specKey, valuesSelectorKey, "labels") {
			if labels, err := fc.NestedMap(specKey, valuesSelectorKey, "labels"); err == nil {
				valuesSelector.Labels = labels
			}
		}
		
		// Annotations
		if fc.IsNestedMap(specKey, valuesSelectorKey, "annotations") {
			if annotations, err := fc.NestedMap(specKey, valuesSelectorKey, "annotations"); err == nil {
				valuesSelector.Annotations = annotations
			}
		}
		
		// Name
		if name, found, err := fc.NestedString(specKey, valuesSelectorKey, "name"); found && err == nil {
			valuesSelector.Name = name
		}
		
		helmRelease.ValuesSelector = valuesSelector
	}

	// Get includeCRDs
	if includeCRDs, found, err := fc.NestedBool(specKey, includecrdsKey); found && err == nil {
		helmRelease.IncludeCRDs = includeCRDs
	}

	// Get apiVersions
	if fc.IsNestedSlice(specKey, apiVersionsKey) {
		apiVersionsSlice, err := fc.NestedStringSlice(specKey, apiVersionsKey)
		if err == nil {
			helmRelease.ApiVersions = apiVersionsSlice
		}
	}

	// Get skipTests
	if skipTests, found, err := fc.NestedBool(specKey, skipTestsKey); found && err == nil {
		helmRelease.SkipTests = skipTests
	}

	// Get releaseName (optional)
	if releaseName, found, err := fc.NestedString(specKey, releaseNameKey); found && err == nil {
		helmRelease.ReleaseName = releaseName
	} else {
		// Default to metadata.name if not specified
		helmRelease.ReleaseName = name
	}

	return helmRelease, nil
}

func processValuesSelector(release *types.HelmRelease, items []*fn.KubeObject) error {
	// If no valuesSelector, nothing to do
	if release.ValuesSelector == nil {
		return nil
	}
	
	// Initialize values map if nil
	if release.Values == nil {
		release.Values = make(map[string]interface{})
	}
	
	// Find matching ConfigMaps/Secrets and merge their values
	for _, obj := range items {
		kind := obj.GetKind()
		
		// Filter by kind if specified
		if release.ValuesSelector.Kind != "" && release.ValuesSelector.Kind != kind {
			continue
		}
		
		// Only process ConfigMap or Secret
		if kind != "ConfigMap" && kind != "Secret" {
			continue
		}
		
		// Match by name if specified
		if release.ValuesSelector.Name != "" {
			if release.ValuesSelector.Name != obj.GetName() {
				continue
			}
		}
		
		// Match by labels if specified
		if len(release.ValuesSelector.Labels) > 0 {
			match := true
			for k, v := range release.ValuesSelector.Labels {
				if labelVal, found := obj.GetLabel(k); !found || labelVal != v {
					match = false
					break
				}
			}
			if !match {
				continue
			}
		}
		
		// Match by annotations if specified
		if len(release.ValuesSelector.Annotations) > 0 {
			match := true
			for k, v := range release.ValuesSelector.Annotations {
				if annoVal, found := obj.GetAnnotation(k); !found || annoVal != v {
					match = false
					break
				}
			}
			if !match {
				continue
			}
		}
		
		// Found a match, extract values from data
		if obj.IsNestedMap("data") {
			data, err := obj.NestedMap("data")
			if err != nil {
				return fmt.Errorf("failed to get data from %s/%s: %v", kind, obj.GetName(), err)
			}
			
			// Merge data into values
			for k, v := range data {
				// In case of ConfigMap or Secret, values are always strings
				// We need to parse them if they're intended to be used as Helm values
				// For simplicity, we're just using the string values directly here
				release.Values[k] = v
			}
		}
	}
	
	return nil
}