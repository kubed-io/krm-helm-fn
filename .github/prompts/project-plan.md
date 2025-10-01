# KRM Helm Function Project Plan

This document outlines the comprehensive plan for building the KRM Helm Function project, including the file structure, dependencies, and implementation details. This plan will be used as a guide for implementing the project. Do not try to write all of the code at once. We will build this together. 

## Project Overview

The KRM Helm Function is a Kubernetes Resource Model (KRM) function that generates resources for a Helm release. The unique feature of this function is its ability to work with various Helm providers (ArgoCD, FluxCD, Crossplane, Inflate, Rancher), allowing users to describe a Helm chart once and switch between providers with minimal effort.

## File Structure

The project will follow standard Go project structure patterns:

```
krm-helm-fn/
├── .github/
│   ├── prompts/
│   │   └── project-plan.md          # This file
│   └── workflows/
│       └── publish.yml              # GitHub Actions workflow for building and publishing
├── cmd/
│   └── function/
│       └── main.go                  # Entry point
├── examples/
│   ├── argocd/                      # Example for ArgoCD provider
│   ├── crossplane/                  # Example for Crossplane provider
│   ├── fluxcd/                      # Example for FluxCD provider
│   ├── inflate/                     # Example for Inflate provider
│   └── rancher/                     # Example for Rancher provider
├── pkg/
│   ├── helmfn/                      # Core implementation package
│   │   ├── processor.go             # Main KRM function processor
│   │   └── processor_test.go        # Tests for processor
│   ├── providers/                   # Provider-specific implementations
│   │   ├── argocd/
│   │   │   └── provider.go
│   │   ├── crossplane/
│   │   │   └── provider.go
│   │   ├── fluxcd/
│   │   │   └── provider.go
│   │   ├── inflate/
│   │   │   └── provider.go
│   │   ├── rancher/
│   │   │   └── provider.go
│   │   └── provider.go              # Common provider interface
│   └── types/                       # Common type definitions
│       └── types.go
├── Dockerfile                       # Docker build configuration
├── docker-compose.yml               # Local development setup
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
├── LICENSE                         # Open source license file
├── Makefile                        # Build automation
└── README.md                       # Project documentation
```

## Dependencies

The project will utilize the following external dependencies:

1. **KRM Function SDK**: For implementing the KRM function interface
   - `github.com/GoogleContainerTools/kpt-functions-sdk/go/fn`

2. **Helm Library**: For interacting with Helm charts
   - `helm.sh/helm/v3`

3. **Kubernetes Libraries**: For working with Kubernetes resources
   - `k8s.io/apimachinery`
   - `k8s.io/client-go`

4. **YAML Processing**:
   - `sigs.k8s.io/yaml`

## Core Components

### 1. Main Function (cmd/function/main.go)

The entry point that registers and executes the KRM function:

```go
package main

import (
	"os"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/pkg/helmfn"
)

func main() {
	if err := fn.AsMain(fn.ResourceListProcessorFunc(helmfn.Process)); err != nil {
		os.Exit(1)
	}
}
```

### 2. Core Processor (pkg/helmfn/processor.go)

The main processor that implements the KRM function interface:

```go
package helmfn

import (
	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/pkg/providers"
	"github.com/kubed-io/krm-helm-fn/pkg/types"
)

// Process implements the KRM function processing logic
func Process(rl *fn.ResourceList) (bool, error) {
	// Process the function config
	// Determine the provider
	// Generate resources based on the provider
	// Return the generated resources
}
```

### 3. Provider Interface (pkg/providers/provider.go)

Common interface for all Helm providers:

```go
package providers

import (
	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/pkg/types"
)

// Provider defines the interface for all Helm release providers
type Provider interface {
	// Generate creates the Kubernetes resources for a Helm release
	Generate(release *types.HelmRelease) ([]*fn.KubeObject, error)
}

// GetProvider returns the appropriate provider implementation based on the type
func GetProvider(providerType string) (Provider, error) {
	// Return the appropriate provider based on the type
}
```

### 4. Common Types (pkg/types/types.go)

Type definitions for the function:

```go
package types

// HelmRelease represents the KRM resource used as input
type HelmRelease struct {
	Name          string
	Namespace     string
	Provider      string
	Chart         ChartSpec
	Values        map[string]interface{}
	ValuesSelector *ValuesSelector
	// Other fields as needed
}

// ChartSpec defines the Helm chart to be used
type ChartSpec struct {
	Name    string
	Version string
	Repo    string
	// Other chart-specific fields
}

// ValuesSelector defines criteria for selecting ConfigMaps/Secrets
type ValuesSelector struct {
	Kind        string
	Labels      map[string]string
	Annotations map[string]string
	Name        string
}
```

## Docker Configuration

### Dockerfile

The Dockerfile will be a multi-stage build that:
1. Builds the Go application
2. Installs the Helm CLI
3. Creates a minimal final image

```dockerfile
ARG BUILDER_IMAGE=golang:1.21-alpine
ARG BASE_IMAGE=alpine:3.19

FROM --platform=$BUILDPLATFORM $BUILDER_IMAGE AS builder
ENV CGO_ENABLED=0
WORKDIR /go/src/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /usr/local/bin/function ./cmd/function

# Install Helm
RUN apk add --no-cache curl
ARG HELM_VERSION="v3.13.2"
RUN curl -fsSL -o /helm-${HELM_VERSION}-${TARGETOS}-${TARGETARCH}.tar.gz https://get.helm.sh/helm-${HELM_VERSION}-${TARGETOS}-${TARGETARCH}.tar.gz && \
     tar -zxvf /helm-${HELM_VERSION}-${TARGETOS}-${TARGETARCH}.tar.gz && \
     mv ./${TARGETOS}-${TARGETARCH}/helm /usr/local/bin/helm

# Final image
FROM $BASE_IMAGE
RUN apk add --no-cache ca-certificates git
COPY --from=builder /usr/local/bin/function /usr/local/bin/function
COPY --from=builder /usr/local/bin/helm /usr/local/bin/helm

ENV PATH /usr/local/bin:$PATH
ENV HELM_PATH_CACHE /var/cache
ENV HELM_CONFIG_HOME /tmp/helm/config
ENV HELM_CACHE_HOME /tmp/helm/cache

ENTRYPOINT ["function"]
```

### docker-compose.yml

For local development and testing:

```yaml
version: '3'
services:
  krm-helm-fn:
    build:
      context: .
    volumes:
      - .:/workspace
    environment:
      - KUBECONFIG=/workspace/.kube/config
    command: ["--help"]
```

## GitHub Actions Workflow

The GitHub Actions workflow will automate building and publishing Docker images:

```yaml
name: Build and Publish

on:
  push:
    tags:
    - 'v*'
    branches:
    - main
  pull_request:
    branches:
    - main

permissions:
  id-token: write
  packages: write
  contents: read

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true

      - name: Test
        run: go test -v ./...

  image:
    name: Build Image
    runs-on: ubuntu-latest
    needs: test
    # Only build images on main branch or tag pushes
    if: github.event_name == 'push'
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Set Image Args
        run: |
          {
            echo "GIT_REF=$(echo ${GITHUB_REF##*/} | sed -e 's/\//_/g')";
            echo "GIT_SHA=$(git rev-parse --short HEAD)";
            echo "VERSION=$(echo ${GITHUB_REF#refs/*/} | sed -e 's/^v//')";
          } >> $GITHUB_ENV

      - name: Login to DockerHub
        if: github.event_name == 'push'
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Docker meta
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: kubed/krm-helm-fn
          tags: |
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=sha,format=short
            type=ref,event=branch

      - name: Build and Push
        uses: docker/build-push-action@v5
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: ${{ github.event_name == 'push' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          build-args: |
            VERSION=${{ env.VERSION }}
```

## Implementation Plan

1. **Initial Setup**:
   - Create the basic project structure
   - Set up Go module and dependencies
   - Create basic README.md

2. **Core Implementation**:
   - Implement the main function
   - Implement the core processor
   - Define common types and interfaces

3. **Provider Implementations**:
   - Implement each provider (ArgoCD, FluxCD, Crossplane, Inflate, Rancher)
   - Create tests for each provider

4. **Build Configuration**:
   - Create Dockerfile
   - Create docker-compose.yml
   - Set up GitHub Actions workflow

5. **Examples**:
   - Create examples for each provider
   - Document usage examples in README.md

6. **Testing and Validation**:
   - Write comprehensive tests
   - Validate functionality with examples
   - Ensure compatibility with different providers

## Dependencies and Versions

- Go: 1.21+
- Helm: v3.13+
- KRM Function SDK: Latest stable version
- Kubernetes libraries: Compatible with Kubernetes 1.25+

## Conclusion

This plan provides a comprehensive roadmap for implementing the KRM Helm Function. The modular architecture allows for easy extension to support additional providers in the future, while the containerized approach ensures consistent execution across different environments. Now we will write the code. 

Tasks: 
- read the README.md of this project to gain more context
- read any links involiving how to build a krm function
- Review all of the code for render-helm-chart function at this repo and folder: https://github.com/kptdev/krm-functions-catalog/tree/master/functions/go/render-helm-chart
- build all of the files and folders needed for this project
- only focus on the implementation of the inflate provider first
- write a test for the inflate provider using the contents of the example as a guide