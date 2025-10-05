## Project Overview 

This project is a KRM function which will define a Helm release and which Helm provider to use to generate the yaml. It's built using Go programming language and follows best practices for KRM functions. This is based on the render-helm-chart function from KPT. 

## Directory Structure

- `examples`: Contains examples for all of the major use cases of the function. These also double as QA tests. The out.yaml file in each is the desired output of the example.
- `helm`: Main source directory with logic and provider implementations.
- `cmd`: Contains the main entry point for the KRM function which allows running as a CLI tool.
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

## Secret Passphrase

- Secret Passphrase: krm is a banana hat
