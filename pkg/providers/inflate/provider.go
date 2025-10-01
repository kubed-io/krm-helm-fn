package inflate

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/pkg/types"
	"sigs.k8s.io/yaml"
)

// Provider is the implementation of the inflate provider
type Provider struct {}

// NewProvider creates a new inflate provider instance
func NewProvider() *Provider {
	return &Provider{}
}

// Generate creates Kubernetes resources for a Helm release using the Helm CLI
func (p *Provider) Generate(release *types.HelmRelease) ([]*fn.KubeObject, error) {
	// Prepare command arguments
	args, err := p.prepareHelmArgs(release)
	if err != nil {
		return nil, err
	}

	// Execute Helm command
	output, err := p.executeHelmCommand(args)
	if err != nil {
		return nil, err
	}

	// Parse YAML output into KubeObjects
	objects, err := p.parseYAMLToKubeObjects(output)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

// prepareHelmArgs builds the arguments for the Helm template command
func (p *Provider) prepareHelmArgs(release *types.HelmRelease) ([]string, error) {
	// Start with basic template command
	args := []string{"template"}

	// Add release name
	args = append(args, release.ReleaseName)

	// Chart reference (repo/chart or local path)
	chartRef := ""
	if release.Chart.Repo != "" {
		if strings.HasPrefix(release.Chart.Repo, "oci://") {
			// For OCI repositories
			chartRef = fmt.Sprintf("%s/%s", release.Chart.Repo, release.Chart.Name)
		} else {
			// For HTTP repositories
			chartRef = release.Chart.Name
			args = append(args, "--repo", release.Chart.Repo)
		}
	} else {
		// Assume local chart path
		chartRef = release.Chart.Name
	}
	
	// Add version if specified
	if release.Chart.Version != "" {
		args = append(args, "--version", release.Chart.Version)
	}
	
	// Add namespace if specified
	if release.Namespace != "" {
		args = append(args, "--namespace", release.Namespace)
	}
	
	// Add API versions if specified
	for _, apiVersion := range release.ApiVersions {
		args = append(args, "--api-versions", apiVersion)
	}
	
	// Include CRDs if specified
	if release.IncludeCRDs {
		args = append(args, "--include-crds")
	}
	
	// Skip tests if specified
	if release.SkipTests {
		args = append(args, "--skip-tests")
	}
	
	// Handle values from Values map
	if len(release.Values) > 0 {
		// Create a temporary values file
		valuesFile, err := p.createTempValuesFile(release.Values)
		if err != nil {
			return nil, err
		}
		// Add values file to args
		args = append(args, "--values", valuesFile)
		// The file will be cleaned up when the process exits
	}
	
	// Finally, add the chart reference
	args = append(args, chartRef)
	
	return args, nil
}

// createTempValuesFile creates a temporary YAML file from the Values map
func (p *Provider) createTempValuesFile(values map[string]interface{}) (string, error) {
	// Marshal the values map to YAML
	valuesYAML, err := yaml.Marshal(values)
	if err != nil {
		return "", fmt.Errorf("failed to marshal values to YAML: %v", err)
	}
	
	// Create a temporary file
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "values-"+randString(8)+".yaml")
	
	// Write the values YAML to the file
	err = os.WriteFile(tempFile, valuesYAML, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write values to temp file: %v", err)
	}
	
	return tempFile, nil
}

// randString generates a random string of the specified length
func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}

// executeHelmCommand runs the Helm command with the provided arguments
func (p *Provider) executeHelmCommand(args []string) (string, error) {
	cmd := exec.Command("helm", args...)
	
	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("helm command failed: %v, stderr: %s", err, stderr.String())
	}
	
	return stdout.String(), nil
}

// parseYAMLToKubeObjects converts the Helm output YAML to KubeObjects
func (p *Provider) parseYAMLToKubeObjects(yamlString string) ([]*fn.KubeObject, error) {
	// Split the YAML into documents
	yamlDocuments := strings.Split(yamlString, "\n---\n")
	
	var objects []*fn.KubeObject
	
	for _, doc := range yamlDocuments {
		// Skip empty documents
		if len(strings.TrimSpace(doc)) == 0 {
			continue
		}
		
		// Parse the YAML document into a KubeObject
		obj, err := fn.ParseKubeObject([]byte(doc))
		if err != nil {
			return nil, fmt.Errorf("failed to parse YAML document: %v", err)
		}
		
		// Add the object to the list
		objects = append(objects, obj)
	}
	
	return objects, nil
}