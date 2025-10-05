## Project Overview 

This project is a KRM function which will define a Helm release and which Helm provider to use to generate the yaml. It's built using Go programming language and follows best practices for KRM functions. This is based on the render-helm-chart function from KPT. 

## Directory Structure

- `examples`: Contains examples for all of the major use cases of the function. These also double as QA tests. The out.yaml file in each is the desired output of the example.
- `helmfn`: Main source directory with the function processor and related code.
- `cmd`: Contains the main entry point for the KRM function which allows running as a CLI tool.
- `providers`: Contains provider implementations (e.g., ArgoCD, Inflate).
- `.github`: Contains GitHub specific files like workflows and issue templates.
- `.vscode`: Contains VSCode specific settings for development.

## Libraries and Frameworks 

- This is a kubernetes Resource Model (KRM) function.
- uses helm cli for the inflate provider
- written in Go programming language
- uses kpt-functions-sdk for KRM function framework

## Important Files

- `README.md`: The main readme file with project overview and instructions. This includes all of the information you need to know about the project including requirements. 
- `CONTRIBUTING.md`: Guidelines for contributing to the project. Includes information on working with the project. Includes coding standards to abide by. 

## Using Web References

The README.md has a `## References` section with links to important related documentation and resources.

- use these as your source of truth to research
- do not make up information, always ground your answers in the references
- use discretion when choosing which references to use
- use the chat history to know what has been used already

## Examples and Testing

- the README.md requirements drive the examples
- the examples drive the tests
- the tests drive the code
- examples define the expected output and behavior

## Helm Providers 

- Inflate with CLI
- ArgoCD
- FluxCD 
- Crossplane Helm Provider
- Rancher K3s Helm Controller

Here is a list of instruction files that contain rules for modifying or creating new code.
These files are important for ensuring that the code is modified or created correctly.
Please make sure to follow the rules specified in these files when working with the codebase.
If the file is not already available as attachment, use the `read_file` tool to acquire it.
Make sure to acquire the instructions before making any changes to the code.
| File | Applies To | Description |
| ------- | --------- | ----------- |
| '/workspaces/krm-helm-fn/.github/instructions/go.instructions.md' | **/*.go,**/go.mod,**/go.sum | Instructions for writing Go code following idiomatic Go practices and community standards |
| '/workspaces/krm-helm-fn/.github/instructions/yaml.instructions.md' | **/*.yaml, **/*.yml | Python coding conventions and guidelines |
| '/workspaces/krm-helm-fn/.github/instructions/github-workflows.instructions.md' | .github/workflows/*.yml,.github/actions/*/*.yml | Instructions for creating and modifying GitHub Actions workflows |
