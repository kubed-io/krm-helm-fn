# Getting Project Started

Make a plan for building this project. The goal is to make a report in the directory .github/prompts/project-plan.md. This report should be clear and concise on what all the files that are needed for a go project, the structure of the project, and any other relevant information. This will be used by another chat to actually build the project. In the meant time the goal is to chat with me to come up with a good plan.

Research the following first:
- read the README.md of this project to gain more context
  - read any links involiving how to build a krm function
- Review all of the code for render-helm-chart function at this repo and folder: https://github.com/kptdev/krm-functions-catalog/tree/master/functions/go/render-helm-chart
- make sure we have all the necessary files for a go project 
- up to date versions for any dependencies
- Dockerfile and docker-compose.yml
- a github action to build the image
  - example of the gha workflow ot use: https://raw.githubusercontent.com/kubed-io/krm-py/refs/heads/main/.github/workflows/publish.yml
- use standard go patterns for the project structure

