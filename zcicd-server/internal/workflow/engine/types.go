package engine

import "time"

// WorkflowModel represents the platform's workflow model for Tekton rendering.
type WorkflowModel struct {
	WorkflowID   string
	WorkflowName string
	RunID        string
	RunNumber    int
	ProjectID    string
	Namespace    string
	Stages       []StageModel
	Params       map[string]string
}

// StageModel represents a single stage within a workflow.
type StageModel struct {
	ID        string
	Name      string
	StageType string // build, test, deploy, custom
	SortOrder int
	Jobs      []JobModel
	Config    map[string]interface{}
}

// JobModel represents a single job within a stage.
type JobModel struct {
	ID      string
	Name    string
	JobType string // build, test, deploy, custom, approval
	Config  map[string]interface{}
	Timeout int
}

// BuildModel represents a build configuration for Tekton Task rendering.
type BuildModel struct {
	BuildConfigID  string
	RunID          string
	RunNumber      int
	ProjectID      string
	ServiceName    string
	Namespace      string
	RepoURL        string
	Branch         string
	CommitSHA      string
	BuildScript    string
	DockerfilePath string
	DockerContext   string
	ImageRepo      string
	ImageTag       string
	BuildEnv       map[string]string
	Variables      map[string]string
	CacheEnabled   bool
}

// RunStatus represents the status of a Tekton run.
type RunStatus struct {
	Name       string
	Status     string // pending, running, succeeded, failed, cancelled
	StartedAt  *time.Time
	FinishedAt *time.Time
	Message    string
	Steps      []StepStatus
}

// StepStatus represents the status of a single step within a run.
type StepStatus struct {
	Name      string
	Status    string
	Container string
	LogURL    string
}
