package types

// HelmRelease represents the KRM resource used as input
type HelmRelease struct {
	Name          string
	Namespace     string
	Provider      string
	Chart         ChartSpec
	Values        map[string]interface{}
	ValuesSelector *ValuesSelector
	// Additional fields
	IncludeCRDs   bool
	ApiVersions   []string
	SkipTests     bool
	ReleaseName   string
}

// ChartSpec defines the Helm chart to be used
type ChartSpec struct {
	Name    string
	Version string
	Repo    string
	// Additional fields
	Auth   *ChartAuth
}

// ChartAuth defines authentication details for private Helm repos
type ChartAuth struct {
	Username string
	Password string
}

// ValuesSelector defines criteria for selecting ConfigMaps/Secrets
type ValuesSelector struct {
	Kind        string
	Labels      map[string]string
	Annotations map[string]string
	Name        string
}