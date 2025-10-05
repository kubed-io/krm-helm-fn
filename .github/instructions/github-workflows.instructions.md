---
description: 'Python coding conventions and guidelines'
applyTo: '.github/workflows/*.yml, .github/actions/*/*.yml'
---
# GitHub Actions Workflow Instructions

These instructions apply to all files matching the following patterns:
- `.github/workflows/*.yml`
- `.github/actions/*/*.yml`

## YAML Formatting

For all YAML files, please adhere to the guidelines specified in the [yaml.instructions.md](./yaml.instructions.md) file. This ensures consistency and readability across the project.

## Job Naming

All jobs within a workflow must have a `name` key. This name should be a concise and descriptive summary of what the job does.

## Step Naming and Identification

Every step in a job must include both an `id` and a `name` key.

- `name`: The `name` should clearly describe the action the step is performing.
- `id`: The `id` should be a short, lowercase, and hyphenated version of the name, suitable for use in expressions.

Example:
```yaml
- name: Check out repository code
  id: checkout
  uses: actions/checkout@v3
```
