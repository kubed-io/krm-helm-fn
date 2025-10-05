package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kptdev/krm-functions-sdk/go/fn"
)

// ExampleFiles represents the files in an example directory
type ExampleFiles struct {
	Release  *fn.KubeObject
	Values   *fn.KubeObject
	Expected []*fn.KubeObject
}

// LoadExampleFiles loads the example files from a given example directory
func LoadExampleFiles(exampleDir string) (*ExampleFiles, error) {
	example := &ExampleFiles{}

	// Load release.yaml
	releasePath := filepath.Join(exampleDir, "release.yaml")
	releaseBytes, err := os.ReadFile(releasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read release.yaml: %w", err)
	}

	release, err := fn.ParseKubeObject(releaseBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse release.yaml: %w", err)
	}
	example.Release = release

	// Load values.yaml and create a ConfigMap
	valuesPath := filepath.Join(exampleDir, "values.yaml")
	valuesBytes, err := os.ReadFile(valuesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read values.yaml: %w", err)
	}

	// Create ConfigMap with values
	valuesConfigMap, err := createValuesConfigMap(string(valuesBytes), release.GetName())
	if err != nil {
		return nil, fmt.Errorf("failed to create values ConfigMap: %w", err)
	}
	example.Values = valuesConfigMap

	// Load expected output from out.yaml
	outPath := filepath.Join(exampleDir, "out.yaml")
	outBytes, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read out.yaml: %w", err)
	}

	expected, err := parseExpectedOutput(outBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse out.yaml: %w", err)
	}
	example.Expected = expected

	return example, nil
}

// CreateResourceList creates a ResourceList with the release as functionConfig and values as input
func (ef *ExampleFiles) CreateResourceList() *fn.ResourceList {
	rl := &fn.ResourceList{
		FunctionConfig: ef.Release,
		Items:          []*fn.KubeObject{ef.Values},
	}
	return rl
}

// createValuesConfigMap creates a ConfigMap from values.yaml content
func createValuesConfigMap(valuesContent, releaseName string) (*fn.KubeObject, error) {
	configMapYAML := fmt.Sprintf(`apiVersion: v1
kind: ConfigMap
metadata:
  name: %s-values
  annotations:
    krm.kubed.io/helm-values: "%s"
data:
  values.yaml: |
%s`, releaseName, releaseName, indentLines(valuesContent, "    "))

	return fn.ParseKubeObject([]byte(configMapYAML))
}

// parseExpectedOutput parses the expected output YAML which may contain multiple documents
func parseExpectedOutput(yamlBytes []byte) ([]*fn.KubeObject, error) {
	var objects []*fn.KubeObject

	// Split YAML documents
	docs := splitYAMLDocuments(string(yamlBytes))

	for _, doc := range docs {
		if doc == "" {
			continue
		}

		obj, err := fn.ParseKubeObject([]byte(doc))
		if err != nil {
			return nil, fmt.Errorf("failed to parse YAML document: %w", err)
		}
		objects = append(objects, obj)
	}

	return objects, nil
}

// splitYAMLDocuments splits a multi-document YAML string into individual documents
func splitYAMLDocuments(yamlContent string) []string {
	var result []string

	// Split on document separators
	docs := strings.Split(yamlContent, "\n---\n")
	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc != "" && doc != "---" {
			result = append(result, doc)
		}
	}

	return result
}

// indentLines adds the given indentation to each line
func indentLines(content, indent string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}
