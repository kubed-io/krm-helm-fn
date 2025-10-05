# KRM Helm Function  

A KRM Function which generates resources for a helm release. What makes this function unique is it allows you to choose from a number of helm providers which are described below. Each provider has a unique way of specifying a helm chart but in the end it's just a helm chart and values just like using the cli. Now you can describe your helm chart once and then change brands with little effort. Great for trying the different options out there. 

Here is the KRM function minimum manifest to get started. 

```yaml
apiVersion: krm.kubed.io
kind: HelmRelease
metadata:
  name: my-app
  namespace: my-system
  annotations: 
    config.kubernetes.io/function: |
      container: 
        image: kubed/krm-helm-fn:latest
```

## Examples  

Each of the folders in the examples directory contain a kustomization.yaml file which uses the generator to create the resources for the given provider. Each example also includes a file called `output.yaml` which contains the expected output from the generator. All of the examples install minecraft as a simple helm chart to use. 

Simply run with the build command replacing `<name>` with the name of the example directory. 

```bash
kubectl build examples/<name>
```

For example, run the inflator with: 

```bash
kubectl build examples/inflate
```

## Providers Classes

Most of the providers require an operator to be running in the cluster to reconcile the resources. The only exception is the Inflate provider which simply generates all of the resources from the helm chart and values. The provider is decided by setting the `spec.provider` field in the HelmRelease resource to one of; argocd, fluxcd, crossplane, inflate.

### Inflate 

This is a special provider because the provider is this function itself. This is same functionality as the helm inflator for kustomize and kpt where the function itself use the Helm cli to run the helm template command with the given configuration. This way the generator produces all of the resources from the release as resources in cli output which can now be used with kustomize or kpt. 

[Example](./examples/inflate)

### ArgoCD  

This provider generates an ArgoCD Application resource. 

[Example](./examples/argocd)

### FluxCD

This provider generates a FluxCD HelmRelease resource.

[Example](./examples/fluxcd)

### Crossplane

This provider generates a Crossplane HelmRelease resource. Requires the [Crossplane Contrib Helm Provider](https://github.com/crossplane-contrib/provider-helm).

[Example](./examples/crossplane)

### Rancher Helm

This provider generates a Rancher HelmChart resource.

[Example](./examples/rancher)

## Spec 

This section describes all that you can do with the HelmRelease spec. Meaning the key `spec` in the KRM resource. Long story short, this seeks to be a one size fits all helm release spec for whatever provider you choose. This covers all the features of Helm that all of the providers support. Any functionality that is unique to a single provider can be further described in the annotations field. 

### Values

This function supports two ways to provide values to the Helm chart: inline values using `spec.values` and values from `ConfigMap`s or `Secret`s.

#### Values from ConfigMaps and Secrets

You can use Kustomize's `configMapGenerator` or `secretGenerator` to create a `ConfigMap` or `Secret` from a values file. This is useful if you want to keep your values in separate files and manage them with Kustomize. The user can decide whether the values should be treated as a secret or not.

The function uses a flexible `valuesSelector` field in the `HelmRelease` spec to identify which `ConfigMap`s and `Secret`s to use for values. This selector gives you multiple ways to match resources.

Here's the schema for the `valuesSelector`:

- `kind`: Specify either "ConfigMap" or "Secret" to restrict the type of resource to select. If omitted, both will be searched.
- `labels`: A map of labels to match on a resource.
- `annotations`: A map of annotations to match on a resource.
- `name`: A specific name to match. Can use regex patterns and wildcards.

Here is an example of how to configure this in your `HelmRelease` resource:

```yaml
spec:
  valuesSelector:
    kind: ConfigMap
    name: "my-app-*-values"
    labels:
      app: my-app
      environment: production
    annotations:
      krm.kubed.io/helm-values: "my-app"
```

And here's how to configure your Kustomize generators to create the matching resources:

```yaml
# ConfigMap example
configMapGenerator:
- name: my-app-prod-values
  files:
  - values.yaml
  options:
    labels:
      app: my-app
      environment: production
    annotations:
      krm.kubed.io/helm-values: "my-app"

# Secret example
secretGenerator:
- name: my-app-secret-values
  files:
  - secret-values.yaml
  options:
    labels:
      app: my-app
      environment: production
    annotations:
      krm.kubed.io/helm-values: "my-app"
```

##### Provider-Specific Notes

- **ArgoCD**: The ArgoCD `Application` CRD does not support referencing `ConfigMap`s for values directly. Therefore, the function will merge the values from any matched `ConfigMap` or `Secret` into the `spec.source.helm.valuesObject` field of the generated `Application` resource.

- **Rancher**: The Rancher `HelmChart` CRD can only reference `Secret`s for values, not `ConfigMap`s. To accommodate this, you should use Kustomize's `secretGenerator` and ensure your valuesSelector is configured accordingly. The function will configure the `HelmChart` to reference matching `Secret`s using the `valuesSecrets` field.

  ```yaml
  # In your HelmRelease spec
  valuesSelector:
    kind: Secret
    annotations:
      krm.kubed.io/helm-values: "my-app"
      
  # In your kustomization.yaml
  secretGenerator:
  - name: my-app-values
    files:
    - values.yaml
    options:
      annotations:
        krm.kubed.io/helm-values: "my-app"
  ```

#### Inline Values

You can also provide values directly in the `HelmRelease` resource using the `spec.values` field. This is ideal if your goal is to reduce the number of files you need to manage.

```yaml
spec:
  values:
    replicaCount: 2
```

#### Provider-Specific Behavior

The way values are embedded in the final resource depends on the provider.

*   **Inflate**: The values are passed directly to the `helm template` command.
*   **ArgoCD**: The values are embedded in the `spec.source.helm.valuesObject` field of the `Application` resource.
    > **Note:** Since the ArgoCD `Application` CRD does not natively support referencing `ConfigMap`s for Helm values, this function provides a workaround. It reads the data from any `ConfigMap` or `Secret` matched by the `valuesSelector` during the `kustomize build` process and merges it into the `spec.source.helm.valuesObject` field of the generated `Application` resource.
*   **FluxCD**: The values are embedded in the `spec.values` field of the `HelmRelease` resource.
*   **Crossplane**: The values are embedded in the `spec.forProvider.values` field of the `Release` resource.
*   **Rancher**: The values are embedded in the `spec.valuesContent` field of the `HelmChart` resource.

Each provider has its own way of handling Helm values, but this function provides a consistent way to specify them across all providers.

## References  

- KPT/KRM Functions
  - [KPT Catalog for render-helm-chart](https://catalog.kpt.dev/render-helm-chart/v0.2/)
  - [krm-function-catalog for render-helm-chart](https://github.com/kptdev/krm-functions-catalog/tree/master/functions/go/render-helm-chart) - This is the inspirational start for this project. The example function which happens to do the inflator for helm charts.
  - [KPT Dev - Developing in Go](https://kpt.dev/book/05-developing-functions/#developing-in-go)
  - [KRM Functions Registry](https://github.com/kubernetes-sigs/krm-functions-registry)
  - [KRM Function Spec](https://github.com/kubernetes-sigs/kustomize/blob/master/cmd/config/docs/api-conventions/functions-spec.md)
- Kustomize
  - [Kustomize Helm Inflator](https://kubectl.docs.kubernetes.io/references/kustomize/builtins/#_helmchartinflationgenerator_)
  - [Kustomize configMapGenerator](https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/configmapgenerator/)
- Helm 
  - [Helm Example Repo](https://github.com/helm/examples) - The chart used in the examples
    - [Example Repo Raw Values](https://raw.githubusercontent.com/helm/examples/refs/heads/main/charts/hello-world/values.yaml) - shows what values can be used in the examples
- ArgoCD
  - [Helm Guide](https://argo-cd.readthedocs.io/en/latest/user-guide/helm/)
  - [CRD Resource Schema Reference](https://raw.githubusercontent.com/argoproj/argo-cd/refs/heads/master/manifests/crds/application-crd.yaml)
- FluxCD
  - [FluxCD Helm Release Documentation](https://fluxcd.io/flux/components/helm/helmreleases/)
  - [Helm Release Example Github Repo for FluxCD](https://github.com/fluxcd/flux2-kustomize-helm-example)
- Crossplane
  - [Crossplane Contrib Helm Provider](https://github.com/crossplane-contrib/provider-helm)
  - [Crossplane Helm Marketplace](https://marketplace.upbound.io/providers/crossplane-contrib/provider-helm/v1.0.2)
- Rancher
  - [helm-controller repo](https://github.com/k3s-io/helm-controller)