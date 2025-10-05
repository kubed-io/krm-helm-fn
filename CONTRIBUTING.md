# Contributing

Follow these steps to be proper. There is a lot of very specific steps and automations, please read this entire document before starting. Many questions will be answered by simply getting everything setup exactly the same way in the instructions.


## Build Mod File

Regen the mod with:
```sh
go mod init github.com/kubed-io/krm-helm-fn
```

Gets the SDK
```sh
go get github.com/kptdev/krm-functions-sdk/go/fn
```

## Examples

The `/examples` directory contains a set of use cases for the KRM Helm function. Each subdirectory represents a specific provider or scenario (e.g., `argocd`, `fluxcd`). These examples serve a dual purpose:
1.  **Documentation**: They provide clear, working examples of how to use the function.
2.  **Testing**: They are the foundation of our automated testing infrastructure.

Each example directory typically contains:
- `release.yaml`: The `HelmRelease` custom resource that acts as the input to the function.
- `values.yaml`: Helm values that are provided to the chart.
- `out.yaml`: The expected YAML output after the function runs. This is used as the success criteria in our tests.
- `kustomization.yaml`: A Kustomize configuration to allow for easy manual testing of the example.

When adding a new feature or fixing a bug, you should add or update an example to cover the use case.

## Testing

This project follows the principle that "tests are driven by examples and code is driven by tests." The testing infrastructure is designed to use the resources in the `/examples` directory to validate the function's behavior.

A dedicated utility package, `testutil`, loads the example files to create test data. This ensures that our tests are always running against the same resources that we provide as documentation. The `testutil` package provides helpers to load `release.yaml`, `values.yaml`, and the expected `out.yaml` for use in tests.

This approach provides several benefits:
- **Consistency**: Examples and tests always stay in sync.
- **Real Data**: Tests use actual example data instead of hardcoded mocks.
- **Maintainability**: Changes to examples automatically update test expectations.

### Running Tests

To run all the tests in the project, use the following Go command:
```sh
go test -v ./...
```