package service

// DeployConfig DTOs

type CreateDeployConfigReq struct {
	ServiceID      string                 `json:"service_id" binding:"required"`
	EnvironmentID  string                 `json:"environment_id" binding:"required"`
	Name           string                 `json:"name" binding:"required"`
	DeployType     string                 `json:"deploy_type" binding:"required,oneof=helm kustomize yaml"`
	RepoURL        string                 `json:"repo_url" binding:"required"`
	TargetRevision string                 `json:"target_revision"`
	ChartPath      string                 `json:"chart_path"`
	ValuesOverride map[string]interface{} `json:"values_override"`
	SyncPolicy     string                 `json:"sync_policy" binding:"omitempty,oneof=manual auto"`
	AutoSync       bool                   `json:"auto_sync"`
	SelfHeal       bool                   `json:"self_heal"`
	Prune          bool                   `json:"prune"`
	Namespace      string                 `json:"namespace"`
}

type UpdateDeployConfigReq struct {
	Name           string                 `json:"name"`
	DeployType     string                 `json:"deploy_type" binding:"omitempty,oneof=helm kustomize yaml"`
	RepoURL        string                 `json:"repo_url"`
	TargetRevision string                 `json:"target_revision"`
	ChartPath      string                 `json:"chart_path"`
	ValuesOverride map[string]interface{} `json:"values_override"`
	SyncPolicy     string                 `json:"sync_policy" binding:"omitempty,oneof=manual auto"`
	AutoSync       *bool                  `json:"auto_sync"`
	SelfHeal       *bool                  `json:"self_heal"`
	Prune          *bool                  `json:"prune"`
	Namespace      string                 `json:"namespace"`
}

// Deploy Sync/Rollback DTOs

type TriggerSyncReq struct {
	Revision string `json:"revision"`
}

type RollbackReq struct {
	HistoryID string `json:"history_id" binding:"required"`
}

// Approval DTOs

type ApproveReq struct {
	Comment string `json:"comment"`
}

type RejectReq struct {
	Comment string `json:"comment" binding:"required"`
}

// Environment Variable DTOs

type EnvVariableReq struct {
	VarKey      string `json:"var_key" binding:"required"`
	VarValue    string `json:"var_value"`
	IsSecret    bool   `json:"is_secret"`
	Description string `json:"description"`
}

type BatchEnvVariablesReq struct {
	Variables []EnvVariableReq `json:"variables" binding:"required,dive"`
}

// Environment Resource Quota DTOs

type EnvResourceQuotaReq struct {
	CPURequest    string `json:"cpu_request"`
	CPULimit      string `json:"cpu_limit"`
	MemoryRequest string `json:"memory_request"`
	MemoryLimit   string `json:"memory_limit"`
	PodLimit      int    `json:"pod_limit"`
	StorageLimit  string `json:"storage_limit"`
}
