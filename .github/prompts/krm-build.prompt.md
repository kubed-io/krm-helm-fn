---
mode: 'agent'
model: GPT-4o
tools: ['githubRepo', 'codebase']
description: 'Generate a new React form component'
---

# KRM Build and Test Prompt

Build the KRM Helm function application and report on any errors or issues found.

## Primary Tasks

1. **Execute Build Task**
   - Run the "go: build" VS Code task
   - Check if the build completes successfully
   - Report the build output and any compilation errors

2. **Execute Test Task**
   - Run the "go: test" VS Code task  
   - Report test results and any test failures
   - Identify any missing test coverage

3. **Error Analysis**
   - Check for any Go compilation errors using get_errors
   - Analyze error messages and provide clear explanations
   - Suggest specific fixes for any issues found

4. **Build Health Report**
   - Confirm the binary was created in bin/function
   - Check file size and permissions
   - Verify the build artifacts are complete

## Expected Workflow

```
1. run_task: "go: build" 
2. get_task_output: analyze build results
3. get_errors: check for compilation issues
4. run_task: "go: test"
5. get_task_output: analyze test results  
6. Report comprehensive build status
```

## Success Criteria

- Build completes without errors
- All tests pass
- Binary is created and executable
- No lint or compilation warnings
- Dependencies are properly resolved

## Failure Handling

If errors are found:
- Clearly explain what went wrong
- Provide specific line numbers and files involved
- Suggest actionable steps to resolve issues
- Offer to run additional diagnostic tasks if needed

Execute this prompt to get a complete build health check of the KRM Helm function project.