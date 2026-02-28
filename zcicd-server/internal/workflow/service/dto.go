package service

// Workflow DTOs
type CreateWorkflowRequest struct {
	ProjectID     string                 `json:"project_id" binding:"required,uuid"`
	Name          string                 `json:"name" binding:"required,min=2,max=128"`
	Description   string                 `json:"description"`
	TriggerType   string                 `json:"trigger_type" binding:"omitempty,oneof=manual webhook cron api"`
	TriggerConfig map[string]interface{} `json:"trigger_config"`
	Enabled       *bool                  `json:"enabled"`
	Stages        []CreateStageRequest   `json:"stages"`
}

type CreateStageRequest struct {
	Name      string                 `json:"name" binding:"required,min=2,max=128"`
	StageType string                 `json:"stage_type" binding:"required,oneof=build test deploy custom approval"`
	SortOrder int                    `json:"sort_order"`
	Config    map[string]interface{} `json:"config"`
	Timeout   int                    `json:"timeout"`
	Enabled   *bool                  `json:"enabled"`
	Jobs      []CreateJobRequest     `json:"jobs"`
}

type CreateJobRequest struct {
	Name       string                 `json:"name" binding:"required,min=2,max=128"`
	JobType    string                 `json:"job_type" binding:"required,oneof=build test deploy custom approval"`
	SortOrder  int                    `json:"sort_order"`
	Config     map[string]interface{} `json:"config"`
	TimeoutSec int                    `json:"timeout_sec"`
	Enabled    *bool                  `json:"enabled"`
}

type UpdateWorkflowRequest struct {
	Name          string                 `json:"name" binding:"omitempty,min=2,max=128"`
	Description   string                 `json:"description"`
	TriggerType   string                 `json:"trigger_type" binding:"omitempty,oneof=manual webhook cron api"`
	TriggerConfig map[string]interface{} `json:"trigger_config"`
	Enabled       *bool                  `json:"enabled"`
	Stages        []CreateStageRequest   `json:"stages"`
}

type TriggerWorkflowRequest struct {
	InputParams map[string]string `json:"input_params"`
}

type WorkflowResponse struct {
	ID            string      `json:"id"`
	ProjectID     string      `json:"project_id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	TriggerType   string      `json:"trigger_type"`
	TriggerConfig interface{} `json:"trigger_config"`
	Enabled       bool        `json:"enabled"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
	Stages        []StageResponse `json:"stages,omitempty"`
}

type StageResponse struct {
	ID        string      `json:"id"`
	WorkflowID string     `json:"workflow_id"`
	Name      string      `json:"name"`
	StageType string      `json:"stage_type"`
	SortOrder int         `json:"sort_order"`
	Config    interface{} `json:"config"`
	Timeout   int         `json:"timeout"`
	Enabled   bool        `json:"enabled"`
	Jobs      []JobResponse `json:"jobs,omitempty"`
}

type JobResponse struct {
	ID         string      `json:"id"`
	StageID    string      `json:"stage_id"`
	Name       string      `json:"name"`
	JobType    string      `json:"job_type"`
	SortOrder  int         `json:"sort_order"`
	Config     interface{} `json:"config"`
	TimeoutSec int         `json:"timeout_sec"`
	Enabled    bool        `json:"enabled"`
}

type WorkflowRunResponse struct {
	ID           string      `json:"id"`
	WorkflowID   string      `json:"workflow_id"`
	RunNumber    int         `json:"run_number"`
	Status       string      `json:"status"`
	TriggerType  string      `json:"trigger_type"`
	TriggeredBy  *string     `json:"triggered_by"`
	InputParams  interface{} `json:"input_params"`
	StagesStatus interface{} `json:"stages_status"`
	TektonRefs   interface{} `json:"tekton_refs"`
	StartedAt    *string     `json:"started_at"`
	FinishedAt   *string     `json:"finished_at"`
	DurationSec  *int        `json:"duration_sec"`
	ErrorMessage string      `json:"error_message"`
	CreatedAt    string      `json:"created_at"`
}

// Build DTOs
type CreateBuildConfigRequest struct {
	ProjectID      string                 `json:"project_id" binding:"required,uuid"`
	ServiceID      string                 `json:"service_id" binding:"required,uuid"`
	Name           string                 `json:"name" binding:"required,min=2,max=128"`
	TemplateID     *string                `json:"template_id" binding:"omitempty,uuid"`
	RepoURL        string                 `json:"repo_url" binding:"required,url,max=512"`
	Branch         string                 `json:"branch" binding:"omitempty,max=128"`
	BuildEnv       map[string]string      `json:"build_env"`
	BuildScript    string                 `json:"build_script"`
	DockerfilePath string                 `json:"dockerfile_path" binding:"omitempty,max=256"`
	DockerContext  string                 `json:"docker_context" binding:"omitempty,max=256"`
	ImageRepo      string                 `json:"image_repo" binding:"required,max=256"`
	TagStrategy    string                 `json:"tag_strategy" binding:"omitempty,oneof=branch-commit timestamp latest"`
	CacheEnabled   *bool                  `json:"cache_enabled"`
	Variables      map[string]interface{} `json:"variables"`
}

type UpdateBuildConfigRequest struct {
	Name           string                 `json:"name" binding:"omitempty,min=2,max=128"`
	TemplateID     *string                `json:"template_id" binding:"omitempty,uuid"`
	RepoURL        string                 `json:"repo_url" binding:"omitempty,url,max=512"`
	Branch         string                 `json:"branch" binding:"omitempty,max=128"`
	BuildEnv       map[string]string      `json:"build_env"`
	BuildScript    string                 `json:"build_script"`
	DockerfilePath string                 `json:"dockerfile_path" binding:"omitempty,max=256"`
	DockerContext  string                 `json:"docker_context" binding:"omitempty,max=256"`
	ImageRepo      string                 `json:"image_repo" binding:"omitempty,max=256"`
	TagStrategy    string                 `json:"tag_strategy" binding:"omitempty,oneof=branch-commit timestamp latest"`
	CacheEnabled   *bool                  `json:"cache_enabled"`
	Variables      map[string]interface{} `json:"variables"`
}

type TriggerBuildRequest struct {
	Branch        *string           `json:"branch"`
	CommitSHA     *string           `json:"commit_sha"`
	Variables     map[string]string `json:"variables"`
}

type BuildConfigResponse struct {
	ID             string      `json:"id"`
	ProjectID      string      `json:"project_id"`
	ServiceID      string      `json:"service_id"`
	Name           string      `json:"name"`
	TemplateID     *string     `json:"template_id"`
	RepoURL        string      `json:"repo_url"`
	Branch         string      `json:"branch"`
	BuildEnv       interface{} `json:"build_env"`
	BuildScript    string      `json:"build_script"`
	DockerfilePath string      `json:"dockerfile_path"`
	DockerContext  string      `json:"docker_context"`
	ImageRepo      string      `json:"image_repo"`
	TagStrategy    string      `json:"tag_strategy"`
	CacheEnabled   bool        `json:"cache_enabled"`
	Variables      interface{} `json:"variables"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
}

type BuildRunResponse struct {
	ID            string              `json:"id"`
	BuildConfigID string              `json:"build_config_id"`
	RunNumber     int                 `json:"run_number"`
	Status        string              `json:"status"`
	Branch        string              `json:"branch"`
	CommitSHA     string              `json:"commit_sha"`
	CommitMessage string              `json:"commit_message"`
	ImageTag      string              `json:"image_tag"`
	ImageDigest   string              `json:"image_digest"`
	TektonRef     string              `json:"tekton_ref"`
	LogPath       string              `json:"log_path"`
	TriggeredBy   *string             `json:"triggered_by"`
	StartedAt     *string             `json:"started_at"`
	FinishedAt    *string             `json:"finished_at"`
	DurationSec   *int                `json:"duration_sec"`
	CreatedAt     string              `json:"created_at"`
	BuildConfig   *BuildConfigResponse `json:"build_config,omitempty"`
}

// Template DTOs
type CreateTemplateRequest struct {
	Name          string                 `json:"name" binding:"required,min=2,max=128"`
	Language      string                 `json:"language" binding:"required,max=32"`
	Framework     string                 `json:"framework" binding:"omitempty,max=64"`
	Description   string                 `json:"description"`
	BuildEnv      map[string]string      `json:"build_env"`
	BuildScript   string                 `json:"build_script" binding:"required"`
	DockerfileTpl string                 `json:"dockerfile_tpl"`
	TektonTaskTpl string                 `json:"tekton_task_tpl" binding:"required"`
}

type BuildTemplateResponse struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Language      string      `json:"language"`
	Framework     string      `json:"framework"`
	Description   string      `json:"description"`
	BuildEnv      interface{} `json:"build_env"`
	BuildScript   string      `json:"build_script"`
	DockerfileTpl string      `json:"dockerfile_tpl"`
	TektonTaskTpl string      `json:"tekton_task_tpl"`
	IsSystem      bool        `json:"is_system"`
	CreatedBy     *string     `json:"created_by"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
}
