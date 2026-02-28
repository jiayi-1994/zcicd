package engine

// ArgoApp represents an Argo CD Application.
type ArgoApp struct {
	Name           string
	Namespace      string // argocd namespace
	Project        string // argo project, default "default"
	RepoURL        string
	TargetRevision string
	Path           string // chart path in repo
	DestNamespace  string // target k8s namespace
	DestServer     string // target cluster, default "https://kubernetes.default.svc"
	ValuesOverride map[string]interface{}
	SyncPolicy     string // manual/auto
	AutoSync       bool
	SelfHeal       bool
	Prune          bool
}

// SyncResult from a sync operation.
type SyncResult struct {
	Status     string // Synced/OutOfSync/Unknown
	Health     string // Healthy/Degraded/Progressing/Missing/Suspended/Unknown
	Revision   string
	Message    string
	StartedAt  string
	FinishedAt string
}

// ResourceNode in the resource tree.
type ResourceNode struct {
	Group     string
	Kind      string
	Namespace string
	Name      string
	Status    string
	Health    string
	Message   string
	Children  []ResourceNode
}

// ResourceTree for an application.
type ResourceTree struct {
	Nodes []ResourceNode
}

// AppStatus represents the current status of an Argo CD Application.
type AppStatus struct {
	SyncStatus   string
	HealthStatus string
	Revision     string
	Message      string
	Resources    []ResourceNode
}
