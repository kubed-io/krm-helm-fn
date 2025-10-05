The goal is to use the resources from the examples to run tests. This way our tests can be driven by the examples and the code will be driven by the tests. 

Requirements: 
- loads the out.yaml from the examples to use as a success criteria for only the release objects
- uses the release.yaml from the examples to create the release objects
- manually add the confimap objects with values into the resource list versus trying to run the configmap generator

Tasks: 
- research examples structure and files
- research the code to understand the ins and outs of the release object creation
- make a common resource for building mock ResourceList objects using the example files
- implement this for the argocd provider
