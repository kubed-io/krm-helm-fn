---
description: Add a new provider to the KRM Helm function following the established patterns
mode: krm
model: Claude Sonnet 4 (copilot)
tools: ['runCommands', 'runTasks', 'edit', 'search', 'new', 'todos', 'usages', 'vscodeAPI', 'problems', 'changes', 'testFailure', 'openSimpleBrowser', 'fetch', 'githubRepo']
---

# Add New ${input:providerName:Provider name (e.g., fluxcd, crossplane)} Provider

You are a Go developer working on a KRM Helm function. Add a new provider called **${input:providerName}** following the established patterns in this codebase.

## Required Implementation Steps

### 1. Create Provider Implementation
- **File**: `providers/${input:providerName}/${input:providerName}.go`
- **Pattern**: Follow the ArgoCD provider pattern in `providers/argocd/argocd.go`
- **Interface**: Must implement the `Provider` interface from `helmfn/types/types.go`
- **Methods**: 
  - `ProcessResourceList(*fn.ResourceList) error` - Main processing logic
  - Any helper methods for generating resources

### 2. Add Provider Tests
- **File**: `providers/${input:providerName}/${input:providerName}_test.go`
- **Pattern**: Follow the ArgoCD provider tests in `providers/argocd/argocd_test.go`
- **Required Tests**:
  - `TestNew${input:providerName}Provider` - Constructor test
  - `Test${input:providerName}Provider_ProcessResourceList` - Main functionality test
  - Integration test using example files

### 3. Update Main Processor
- **File**: `helmfn/processor.go`
- **Action**: Add new case to the provider switch statement in the `Process` function
- **Pattern**: Follow existing cases for "argocd", "fluxcd", etc.
- **Import**: Add import for the new provider package

### 4. Add Processor Integration Test
- **File**: `helmfn/processor_test.go`
- **Function**: `TestProcess${input:providerName}Example`
- **Pattern**: Follow `TestProcessFluxCDExample` pattern
- **Verification**: Ensure correct resources are generated with proper API versions and metadata

### 5. Create Example Files
- **Directory**: `examples/${input:providerName}/`
- **Required Files**:
  - `kustomization.yaml` - Kustomize configuration
  - `release.yaml` - HelmRelease resource with provider specification
  - `values.yaml` - Helm values file
  - `out.yaml` - Expected output (can start with empty specs like other providers)

### 6. Build and Test
- **Commands to run**:
  ```bash
  # Run tests
  go test -v ./...
  
  # Build clean
  go clean && go build -o bin/function ./cmd/function
  
  # Build Docker image
  docker compose build krm-helm-fn
  
  # Test with kustomize
  cd examples/${input:providerName}
  kustomize build --enable-alpha-plugins --enable-exec --network --enable-helm .
  ```

## Key Implementation Guidelines

### Provider Interface
Reference `helmfn/types/types.go` for the required interface:
- Implement `ProcessResourceList(*fn.ResourceList) error`
- Follow established error handling patterns
- Use debug logging via `helmfn.DebugLog()` for troubleshooting

### Resource Generation
- Generate appropriate Kubernetes resources for the provider
- Use correct API versions and kinds for the target platform
- Start with empty specs (like ArgoCD/FluxCD) - full implementation comes later
- Ensure proper metadata (name, namespace, labels, annotations)

### Testing Approach
- Unit tests for provider functionality
- Integration tests through processor
- Example-based testing using the `testutil` package
- Verify resource types, API versions, and basic metadata

### Debug Logging
- Use the centralized `helmfn.DebugLog()` function
- Enable with `LOG_LEVEL=debug` environment variable
- Log important steps in resource processing

## Success Criteria

✅ Provider implements the required interface  
✅ All tests pass (`go test -v ./...`)  
✅ Clean build succeeds  
✅ Docker image builds successfully  
✅ Kustomize build works without "unsupported provider" errors  
✅ Generated resources have correct structure (even with empty specs)  
✅ Integration test passes in `processor_test.go`  

## Reference Files
- [ArgoCD Provider](../../providers/argocd/argocd.go) - Main implementation pattern
- [FluxCD Provider](../../providers/fluxcd/fluxcd.go) - Recent implementation example  
- [Provider Interface](../../helmfn/types/types.go) - Required interface definition
- [Processor](../../helmfn/processor.go) - Switch statement location
- [Example Structure](../../examples/fluxcd/) - Directory structure and file patterns

## Commands Reference
```bash
# Development workflow
go test -v ./...                                    # Run all tests
go build -o bin/function ./cmd/function            # Build binary
docker compose build krm-helm-fn                   # Build container
kustomize build --enable-alpha-plugins --enable-exec --network --enable-helm .  # Test integration

# Debug mode
LOG_LEVEL=debug go test -v ./helmfn -run TestProcess${input:providerName}Example
LOG_LEVEL=debug kustomize build --enable-alpha-plugins --enable-exec --network --enable-helm .
```

Focus on following the established patterns exactly - this ensures consistency and maintainability across all providers in the codebase. Mainly note how fluxcd and argocd are only partially implemented with empty specs, and follow that same approach for now. Full implementation of provider-specific logic can come later.