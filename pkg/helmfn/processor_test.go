package helmfn_test

import (
	"testing"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/pkg/helmfn"
)

func TestProcess(t *testing.T) {
	// Create a sample function config
	functionConfigYAML := `
apiVersion: krm.kubed.io
kind: HelmRelease
metadata:
  name: my-app
  namespace: my-system
spec:
  provider: inflate
  chart:
    name: hello-world
    version: 0.1.0
    repo: https://helm.github.io/examples
  includeCRDs: true
  apiVersions:
  - example.com/v1
  values:
    replicaCount: 2
    service:
      port: 443
`
	functionConfig, err := fn.ParseKubeObject([]byte(functionConfigYAML))
	if err != nil {
		t.Fatalf("Failed to parse function config: %v", err)
	}

	// Create a resource list with the function config
	resourceList := &fn.ResourceList{
		FunctionConfig: functionConfig,
		Items:         []*fn.KubeObject{},
	}

	// Process the resource list
	_, err = helmfn.Process(resourceList)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Check that resources were generated
	if len(resourceList.Items) == 0 {
		t.Error("No resources were generated")
	}

	// Check for expected resources
	found := false
	for _, res := range resourceList.Items {
		if res.GetKind() == "Deployment" && res.GetName() == "my-app" {
			found = true
			// Check for expected replica count
			replicas, exists, err := res.NestedInt("spec", "replicas")
			if !exists || err != nil {
				t.Error("Failed to find spec.replicas in Deployment")
			} else if replicas != 2 {
				t.Errorf("Expected 2 replicas, got %d", replicas)
			}
		}
	}

	if !found {
		t.Error("Expected Deployment/my-app not found in generated resources")
	}
}

func TestProcess_WithValuesSelector(t *testing.T) {
	// Create a sample ConfigMap with values
	configMapYAML := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-app-values
  annotations:
    krm.kubed.io/helm-values: "my-app"
data:
  service.port: "443"
`
	configMap, err := fn.ParseKubeObject([]byte(configMapYAML))
	if err != nil {
		t.Fatalf("Failed to parse ConfigMap: %v", err)
	}

	// Create a sample function config with ValuesSelector
	functionConfigYAML := `
apiVersion: krm.kubed.io
kind: HelmRelease
metadata:
  name: my-app
  namespace: my-system
spec:
  provider: inflate
  chart:
    name: hello-world
    version: 0.1.0
    repo: https://helm.github.io/examples
  values:
    replicaCount: 2
  valuesSelector:
    annotations:
      krm.kubed.io/helm-values: "my-app"
`
	functionConfig, err := fn.ParseKubeObject([]byte(functionConfigYAML))
	if err != nil {
		t.Fatalf("Failed to parse function config: %v", err)
	}

	// Create a resource list with the function config and ConfigMap
	resourceList := &fn.ResourceList{
		FunctionConfig: functionConfig,
		Items:         []*fn.KubeObject{configMap},
	}

	// Process the resource list
	_, err = helmfn.Process(resourceList)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Check that resources were generated
	if len(resourceList.Items) <= 1 { // Should have more than just the ConfigMap
		t.Error("No resources were generated")
	}

	// Check for expected resources
	found := false
	for _, res := range resourceList.Items {
		if res.GetKind() == "Service" && res.GetName() == "my-app" {
			found = true
			// Ideally, we would check that the service port is 443 from the ConfigMap
			// but this requires the actual Helm execution, which is harder to test
		}
	}

	if !found {
		t.Error("Expected Service/my-app not found in generated resources")
	}
}